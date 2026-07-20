package usecase

import (
	"context"
	"errors"
	"fmt"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"sort"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TimeOffRequestUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	TimeOffRequestRepo   *repository.TimeOffRequestRepository
	TimeOffTypeRepo      *repository.TimeOffTypeRepository
	TimeOffBalanceRepo   *repository.TimeOffBalanceRepository
	TimeOffApprovalRepo  *repository.TimeOffApprovalRepository
	EmployeeContractRepo *repository.EmployeeContractRepository
}

func NewTimeOffRequestUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	timeOffRequestRepo *repository.TimeOffRequestRepository,
	timeOffTypeRepo *repository.TimeOffTypeRepository,
	timeOffBalanceRepo *repository.TimeOffBalanceRepository,
	timeOffApprovalRepo *repository.TimeOffApprovalRepository,
	employeeContractRepo *repository.EmployeeContractRepository,
) *TimeOffRequestUseCase {
	return &TimeOffRequestUseCase{
		DB:                   db,
		Log:                  log,
		Validate:             validate,
		TimeOffRequestRepo:   timeOffRequestRepo,
		TimeOffTypeRepo:      timeOffTypeRepo,
		TimeOffBalanceRepo:   timeOffBalanceRepo,
		TimeOffApprovalRepo:  timeOffApprovalRepo,
		EmployeeContractRepo: employeeContractRepo,
	}
}

// TODO: Validate business rules (quota, overlapping dates) before insert.
func (c *TimeOffRequestUseCase) CreateRequest(ctx context.Context, employeeID string, request *model.CreateTimeOffRequest) (*model.TimeOffRequestResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	status := strings.ToUpper(strings.TrimSpace(request.RequestStatus))
	if status == "" {
		status = "PENDING"
	}
	request.RequestStatus = status

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	if status != "PENDING" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Pengajuan cuti harus berstatus PENDING")
	}

	startDate, err := lib.ParseDateToUnixMilli(request.StartDate)
	if err != nil || startDate == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid start_date")
	}
	endDate, err := lib.ParseDateToUnixMilli(request.EndDate)
	if err != nil || endDate == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid end_date")
	}
	if startDate > endDate {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid start_date or end_date")
	}
	// Derive requested days from the date range (inclusive).
	const dayMillis = 24 * 60 * 60 * 1000
	request.RequestedDays = int((endDate-startDate)/dayMillis) + 1
	if request.RequestedDays <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid date range")
	}

	// TODO: Consider excluding rejected/cancelled requests from overlap check.
	var overlapCount int64
	if err := tx.Table("time_off_requests").
		Where("employee_id = ?", employeeID).
		Where("request_status IN ?", []string{"PENDING", "APPROVED"}).
		Where("NOT (end_date < ? OR start_date > ?)", startDate, endDate).
		Where("request_status NOT IN ?", []string{"REJECTED", "CANCELLED"}).
		Count(&overlapCount).Error; err != nil {
		c.Log.WithError(err).Error("Failed to check overlap dates")
		return nil, fiber.ErrInternalServerError
	}
	if overlapCount > 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tanggal pengajuan cuti bertabrakan dengan pengajuan lain yang sudah ada")
	}
	// Limit to 5 concurrent leave requests per day across all employees for any day in the requested range.
	// The limit is bases on time off type, so different types of leave (e.g. annual vs sick) can have separate limits.
	var concurrentCount int64
	if err := tx.Table("time_off_requests").
		Where("employee_id != ?", employeeID).
		Where("time_off_type_id = ?", request.TimeOffTypeID).
		Where("request_status IN ?", []string{"PENDING", "APPROVED"}).
		Where("NOT (end_date < ? OR start_date > ?)", startDate, endDate).
		Count(&concurrentCount).Error; err != nil {
		c.Log.WithError(err).Error("Failed to check concurrent leave count")
		return nil, fiber.ErrInternalServerError
	}
	if concurrentCount >= 5 {
		return nil, fiber.NewError(fiber.StatusTooManyRequests, "Kuota Pengajuan Cuti Harian Sudah Penuh")
	}

	timeOffType, err := c.TimeOffTypeRepo.FindByID(tx, request.TimeOffTypeID)
	if err != nil {
		c.Log.WithError(err).Error("Tipe cuti tidak ditemukan")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tipe cuti tidak ditemukan")
	}
	if timeOffType.IsQuotaBased {
		// TODO: Use company timezone when deriving period year.
		periodYear := time.UnixMilli(startDate).UTC().Year()
		balance, err := c.TimeOffBalanceRepo.FindByEmployeeTypeYear(tx, employeeID, request.TimeOffTypeID, periodYear)
		if err != nil {
			c.Log.WithError(err).Error("Kuota cuti tidak ditemukan")
			return nil, fiber.NewError(fiber.StatusBadRequest, "Kuota cuti tidak ditemukan")
		}
		if request.RequestedDays > balance.RemainingDays {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Jumlah hari yang diminta melebihi kuota yang tersedia")
		}
	}

	item := &entity.TimeOffRequest{
		EmployeeID:    employeeID,
		TimeOffTypeId: request.TimeOffTypeID,
		RequestedDays: request.RequestedDays,
		StartDate:     startDate,
		EndDate:       &endDate,
		RequestReason: &request.RequestReason,
		RequestStatus: &status,
		FileUrl:       request.FileUrl,
	}

	if err := c.TimeOffRequestRepo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Gagal membuat pengajuan cuti")
		return nil, fiber.ErrInternalServerError
	}
	// Ensure requested dates are persisted (entity hook sets start_date to now).
	if err := tx.Table("time_off_requests").
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"start_date": startDate,
			"end_date":   endDate,
		}).Error; err != nil {
		c.Log.WithError(err).Error("Gagal memperbarui tanggal pengajuan cuti")
		return nil, fiber.ErrInternalServerError
	}

	// Build approval chain from position hierarchy.
	// TODO: Consider caching org structure to reduce DB roundtrips.
	approvals, err := c.buildApprovalsFromPositionChain(tx, employeeID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to build approval chain")
		return nil, err
	}
	if len(approvals) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tidak ada approver yang tersedia untuk pengajuan cuti ini")
	}
	// Bind approvals to the newly created request; approval order is assigned
	// here (and only here) so it is always sequential 1..N without gaps.
	for i := range approvals {
		approvals[i].TimeOffRequestId = item.ID
		approvals[i].Status = "PENDING"
		approvals[i].ApprovalOrder = i + 1
	}

	if err := c.TimeOffApprovalRepo.CreateMany(tx, approvals); err != nil {
		c.Log.WithError(err).Error("Gagal membuat data approval untuk pengajuan cuti")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.TimeOffRequestToResponse(item), nil
}

// TODO: Add authorization scoping for admin vs current-user list.
func (c *TimeOffRequestUseCase) ListRequests(ctx context.Context, request *model.SearchTimeOffRequest) ([]model.TimeOffRequestResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.TimeOffRequestRepo.List(tx, request, true)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off requests")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffRequestResponse, len(items))
	for i, item := range items {
		responses[i] = *model.TimeOffRequestToResponse(&item)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}

// TODO: Add ownership checks for non-admin users.
func (c *TimeOffRequestUseCase) GetRequestByID(ctx context.Context, id string) (*model.TimeOffRequestResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item, err := c.TimeOffRequestRepo.FindByID(tx, id, true)
	if err != nil {
		c.Log.WithError(err).Error("Cuti tidak ditemukan")
		return nil, fiber.ErrNotFound
	}

	response := model.TimeOffRequestToResponse(item)

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return response, nil
}

func (c *TimeOffRequestUseCase) UpdateRequest(ctx context.Context, id string, request *model.UpdateTimeOffRequest) (*model.TimeOffRequestResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	item, err := c.TimeOffRequestRepo.FindByID(tx, id, true)
	if err != nil {
		c.Log.WithError(err).Error("Cuti tidak ditemukan")
		return nil, fiber.ErrNotFound
	}

	if item.RequestStatus == nil || strings.ToUpper(strings.TrimSpace(*item.RequestStatus)) != "PENDING" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Hanya pengajuan yang berstatus PENDING yang dapat diperbarui")
	}
	if request.RequestStatus != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "request_status cannot be updated from this endpoint")
	}

	startDate := item.StartDate
	endDate := int64(0)
	if item.EndDate != nil {
		endDate = *item.EndDate
	}
	timeOffTypeID := item.TimeOffTypeId

	if request.StartDate != nil {
		startDate, err = lib.ParseDateToUnixMilli(strings.TrimSpace(*request.StartDate))
		if err != nil || startDate == 0 {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid start_date")
		}
	}
	if request.EndDate != nil {
		endDate, err = lib.ParseDateToUnixMilli(strings.TrimSpace(*request.EndDate))
		if err != nil || endDate == 0 {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid end_date")
		}
	}
	if request.TimeOffTypeID != nil {
		timeOffTypeID = strings.TrimSpace(*request.TimeOffTypeID)
		if timeOffTypeID == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "time_off_type_id cannot be empty")
		}
	}

	if startDate > endDate {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tanggal mulai tidak boleh setelah tanggal selesai")
	}

	const dayMillis = 24 * 60 * 60 * 1000
	requestedDays := int((endDate-startDate)/dayMillis) + 1
	if request.RequestedDays != nil && *request.RequestedDays > 0 {
		requestedDays = *request.RequestedDays
	}
	if requestedDays <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tanggal tidak valid atau jumlah hari yang diminta tidak valid")
	}

	var overlapCount int64
	if err := tx.Table("time_off_requests").
		Where("employee_id = ?", item.EmployeeID).
		Where("id <> ?", item.ID).
		Where("request_status IN ?", []string{"PENDING", "APPROVED"}).
		Where("NOT (end_date < ? OR start_date > ?)", startDate, endDate).
		Count(&overlapCount).Error; err != nil {
		c.Log.WithError(err).Error("Failed to check overlap dates")
		return nil, fiber.ErrInternalServerError
	}
	if overlapCount > 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Pengajuan cuti tumpang tindih dengan pengajuan yang sudah ada")
	}

	timeOffType, err := c.TimeOffTypeRepo.FindByID(tx, timeOffTypeID)
	if err != nil {
		c.Log.WithError(err).Error("Jenis cuti tidak ditemukan")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Jenis cuti tidak ditemukan")
	}
	if timeOffType.IsQuotaBased {
		periodYear := time.UnixMilli(startDate).UTC().Year()
		balance, err := c.TimeOffBalanceRepo.FindByEmployeeTypeYear(tx, item.EmployeeID, timeOffTypeID, periodYear)
		if err != nil {
			c.Log.WithError(err).Error("Saldo cuti tidak ditemukan")
			return nil, fiber.NewError(fiber.StatusBadRequest, "Saldo cuti tidak ditemukan")
		}
		if requestedDays > balance.RemainingDays {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Jumlah hari yang diminta melebihi saldo yang tersedia")
		}
	}

	item.TimeOffTypeId = timeOffTypeID
	item.StartDate = startDate
	item.EndDate = &endDate
	item.RequestedDays = requestedDays
	if request.RequestReason != nil {
		reason := strings.TrimSpace(*request.RequestReason)
		item.RequestReason = &reason
	}

	if err := c.TimeOffRequestRepo.Update(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to update time off request")
		return nil, fiber.ErrInternalServerError
	}

	result, err := c.TimeOffRequestRepo.FindByID(tx, id, true)
	if err != nil {
		c.Log.WithError(err).Error("Failed to reload time off request")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.TimeOffRequestToResponse(result), nil
}

func (c *TimeOffRequestUseCase) DeleteRequest(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item, err := c.TimeOffRequestRepo.FindByID(tx, id, false)
	if err != nil {
		c.Log.WithError(err).Error("Time off request not found")
		return fiber.ErrNotFound
	}

	if item.RequestStatus == nil || strings.ToUpper(strings.TrimSpace(*item.RequestStatus)) != "PENDING" {
		return fiber.NewError(fiber.StatusBadRequest, "Hanya pengajuan yang berstatus PENDING yang dapat dihapus")
	}

	if err := c.TimeOffRequestRepo.Delete(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to delete time off request")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

// TODO: Add company scoping if needed.
func (c *TimeOffRequestUseCase) GetRequestOwner(ctx context.Context, id string) (string, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var row struct {
		EmployeeID string `gorm:"column:employee_id"`
	}
	if err := tx.Table("time_off_requests").
		Select("employee_id").
		Where("id = ?", id).
		Take(&row).Error; err != nil {
		c.Log.WithError(err).Error("Time off request not found")
		return "", fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return "", fiber.ErrInternalServerError
	}

	return row.EmployeeID, nil
}

// TODO: Add cycle detection limit configuration.
func (c *TimeOffRequestUseCase) buildApprovalsFromPositionChain(tx *gorm.DB, employeeID string) ([]entity.TimeOffApproval, error) {
	// Fetch active contract to get current position + division.
	contract, err := c.EmployeeContractRepo.FindLatestActiveByEmployee(tx, employeeID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Active contract not found")
	}

	// Load current position and start from its parent.
	type positionRow struct {
		ID         string  `gorm:"column:id"`
		ParentID   *string `gorm:"column:parent_id"`
		Name       string  `gorm:"column:name"`
		IsApprover bool    `gorm:"column:is_approver"`
		CompanyID  string  `gorm:"column:company_id"`
	}

	var current positionRow
	if err := tx.Table("positions").
		Select("id, parent_id, name, is_approver, company_id").
		Where("id = ?", contract.PositionID).
		Take(&current).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Position not found")
	}

	// Intermediate struct to collect approvers from both sources before sorting.
	type approverCandidate struct {
		EmployeeID string
		IsRequired bool
		Depth      int // depth from root: 0 = root, higher = closer to employee
	}

	candidates := make([]approverCandidate, 0, 8)
	// Cache depth-from-root per positionID to avoid repeated walks.
	depthCache := make(map[string]int)

	parentID := current.ParentID

	const maxDepth = 20
	// Walk up the position hierarchy until root or depth limit.
	for depth := 0; depth < maxDepth && parentID != nil; depth++ {
		var parent positionRow
		// Resolve parent position.
		if err := tx.Table("positions").
			Select("id, parent_id, name, is_approver").
			Where("id = ?", *parentID).
			Take(&parent).Error; err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Parent position not found")
		}

		var approver struct {
			EmployeeID string `gorm:"column:employee_id"`
		}
		// Find the approver (employee) who holds the parent position in same division.
		err := tx.Table("employee_contracts").
			Select("employee_id").
			Where("is_active = ?", true).
			Where("position_id = ?", parent.ID).
			// Where("division_id = ?", contract.DivisionID).
			// Where("end_date IS NULL OR end_date >= ?", nowEpoch()).
			Order("start_date DESC").
			Limit(1).
			Take(&approver).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Vacant position: skip it and keep walking up the hierarchy.
				parentID = parent.ParentID
				continue
			}
			c.Log.WithError(err).Error("Failed to resolve approver for position " + parent.Name)
			return nil, fiber.ErrInternalServerError
		}

		// Resolve depth-from-root for this parent position.
		posDepth, err := c.resolveDepthFromRoot(tx, parent.ID, depthCache)
		if err != nil {
			c.Log.WithError(err).Error("Failed to resolve depth for position " + parent.Name)
			return nil, fiber.ErrInternalServerError
		}

		candidates = append(candidates, approverCandidate{
			EmployeeID: approver.EmployeeID,
			IsRequired: parent.IsApprover,
			Depth:      posDepth,
		})
		parentID = parent.ParentID
	}

	// Always include all company-wide required approvers (is_approver = true),
	// ordered deterministically so the approval order is stable across requests.
	var globalApprovers []struct {
		EmployeeID string `gorm:"column:employee_id"`
		PositionID string `gorm:"column:position_id"`
	}
	if err := tx.Table("employee_contracts ec").
		Select("ec.employee_id, ec.position_id").
		Joins("JOIN positions p ON p.id = ec.position_id").
		Where("ec.is_active = ?", true).
		Where("p.is_approver = ?", true).
		Where("p.company_id = ?", current.CompanyID).
		Where("ec.employee_id != ?", employeeID).
		Order("ec.employee_id ASC").
		Find(&globalApprovers).Error; err != nil {
		c.Log.WithError(err).Error("Failed to fetch global required approvers")
		return nil, fiber.ErrInternalServerError
	}

	for _, ga := range globalApprovers {
		posDepth, err := c.resolveDepthFromRoot(tx, ga.PositionID, depthCache)
		if err != nil {
			c.Log.WithError(err).Error("Failed to resolve depth for global approver")
			return nil, fiber.ErrInternalServerError
		}
		candidates = append(candidates, approverCandidate{
			EmployeeID: ga.EmployeeID,
			IsRequired: true,
			Depth:      posDepth,
		})
	}

	// Sort candidates: identify the root position first (Depth == 0),
	// then split the rest into "required" (is_approver = true) and
	// "non-required" groups. Non-required approvers keep depth-based order
	// (closest to employee first). Required approvers are always grouped
	// immediately before root, regardless of their actual depth. Root — if
	// reached — always comes last.
	var rootCandidate *approverCandidate
	nonRequired := make([]approverCandidate, 0, len(candidates))
	required := make([]approverCandidate, 0, len(candidates))

	for i := range candidates {
		cand := candidates[i]
		if cand.Depth == 0 {
			// First candidate found at depth 0 is treated as root.
			if rootCandidate == nil {
				rootCandidate = &candidates[i]
			}
			continue
		}
		if cand.IsRequired {
			required = append(required, cand)
		} else {
			nonRequired = append(nonRequired, cand)
		}
	}

	// Dedup within required group by EmployeeID before sorting, since the
	// same position (and its approver) may be reached both via chain walk
	// and the global query.
	requiredSeen := make(map[string]bool)
	dedupedRequired := make([]approverCandidate, 0, len(required))
	for _, cand := range required {
		if requiredSeen[cand.EmployeeID] {
			continue
		}
		requiredSeen[cand.EmployeeID] = true
		dedupedRequired = append(dedupedRequired, cand)
	}
	required = dedupedRequired

	sort.Slice(required, func(i, j int) bool {
		return required[i].EmployeeID < required[j].EmployeeID
	})

	sort.Slice(nonRequired, func(i, j int) bool {
		return nonRequired[i].Depth > nonRequired[j].Depth
	})
	// Deterministic order within the required group.
	sort.Slice(required, func(i, j int) bool {
		return required[i].EmployeeID < required[j].EmployeeID
	})

	ordered := make([]approverCandidate, 0, len(candidates))
	ordered = append(ordered, nonRequired...)
	ordered = append(ordered, required...)
	if rootCandidate != nil {
		ordered = append(ordered, *rootCandidate)
	}

	// ApprovalOrder is intentionally not set here; CreateRequest assigns it
	// sequentially after the full chain is built.
	approvals := make([]entity.TimeOffApproval, 0, len(ordered))
	seen := make(map[string]bool)
	// Dedup + exclude self AFTER sorting so final order reflects hierarchy depth.
	for _, cand := range ordered {
		if cand.EmployeeID == employeeID || seen[cand.EmployeeID] {
			continue
		}
		seen[cand.EmployeeID] = true
		approvals = append(approvals, entity.TimeOffApproval{
			ApproverId: cand.EmployeeID,
			IsRequired: cand.IsRequired,
			Status:     "PENDING",
		})
	}
	return approvals, nil
}

// resolveDepthFromRoot walks parent_id from positionID up to the root (parent_id IS NULL)
// and returns the depth (root = 0). Results are cached in depthCache to avoid
// redundant queries when the same position appears in both chain and global sources.
func (c *TimeOffRequestUseCase) resolveDepthFromRoot(tx *gorm.DB, positionID string, depthCache map[string]int) (int, error) {
	if d, ok := depthCache[positionID]; ok {
		return d, nil
	}

	// Collect the path from positionID to root so we can cache every node along the way.
	path := []string{positionID}
	currentID := positionID

	const maxDepth = 20
	for i := 0; i < maxDepth; i++ {
		var row struct {
			ID       string  `gorm:"column:id"`
			ParentID *string `gorm:"column:parent_id"`
		}
		if err := tx.Table("positions").
			Select("id, parent_id").
			Where("id = ?", currentID).
			Take(&row).Error; err != nil {
			return 0, fmt.Errorf("position %s not found: %w", currentID, err)
		}

		if row.ParentID == nil {
			// Reached root. path[0] = original positionID, path[last] = root.
			// Root depth = 0, original position depth = len(path)-1.
			for idx, pid := range path {
				depthCache[pid] = len(path) - 1 - idx
			}
			return depthCache[positionID], nil
		}

		// Check if parent is already cached — shortcut the walk.
		if cachedDepth, ok := depthCache[*row.ParentID]; ok {
			// Parent's depth is known; assign depths to the rest of the path.
			for idx, pid := range path {
				depthCache[pid] = cachedDepth + len(path) - idx
			}
			return depthCache[positionID], nil
		}

		path = append(path, *row.ParentID)
		currentID = *row.ParentID
	}

	return 0, fmt.Errorf("position hierarchy exceeds max depth (%d) — possible cycle", maxDepth)
}

package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
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
		return nil, fiber.NewError(fiber.StatusBadRequest, "New request status must be PENDING")
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
		Count(&overlapCount).Error; err != nil {
		c.Log.WithError(err).Error("Failed to check overlap dates")
		return nil, fiber.ErrInternalServerError
	}
	if overlapCount > 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Time off request overlaps with existing request")
	}

	timeOffType, err := c.TimeOffTypeRepo.FindByID(tx, request.TimeOffTypeID)
	if err != nil {
		c.Log.WithError(err).Error("Time off type not found")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Time off type not found")
	}
	if timeOffType.IsQuotaBased {
		// TODO: Use company timezone when deriving period year.
		periodYear := time.UnixMilli(startDate).UTC().Year()
		balance, err := c.TimeOffBalanceRepo.FindByEmployeeTypeYear(tx, employeeID, request.TimeOffTypeID, periodYear)
		if err != nil {
			c.Log.WithError(err).Error("Time off balance not found")
			return nil, fiber.NewError(fiber.StatusBadRequest, "Time off balance not found")
		}
		if request.RequestedDays > balance.RemainingDays {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Requested days exceed remaining balance")
		}
	}

	item := &entity.TimeOffRequest{
		EmployeeId:    employeeID,
		TimeOffTypeId: request.TimeOffTypeID,
		RequestedDays: request.RequestedDays,
		StartDate:     startDate,
		EndDate:       &endDate,
		RequestReason: &request.RequestReason,
		RequestStatus: &status,
	}

	if err := c.TimeOffRequestRepo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create time off request")
		return nil, fiber.ErrInternalServerError
	}
	// Ensure requested dates are persisted (entity hook sets start_date to now).
	if err := tx.Table("time_off_requests").
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"start_date": startDate,
			"end_date":   endDate,
		}).Error; err != nil {
		c.Log.WithError(err).Error("Failed to update time off request dates")
		return nil, fiber.ErrInternalServerError
	}

	// Build approval chain from position hierarchy.
	// TODO: Consider caching org structure to reduce DB roundtrips.
	approvals, err := c.buildApprovalsFromPositionChain(tx, employeeID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to build approval chain")
		return nil, err
	}
	// Bind approvals to the newly created request and set initial status.
	for i := range approvals {
		approvals[i].TimeOffRequestId = item.ID
		approvals[i].Status = "PENDING"
	}
	// Persist approval records.
	if err := c.TimeOffApprovalRepo.CreateMany(tx, approvals); err != nil {
		c.Log.WithError(err).Error("Failed to create time off approvals")
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
		c.Log.WithError(err).Error("Time off request not found")
		return nil, fiber.ErrNotFound
	}

	response := model.TimeOffRequestToResponse(item)

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return response, nil
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
	}

	approvals := make([]entity.TimeOffApproval, 0, 4)
	var current positionRow
	if err := tx.Table("positions").
		Select("id, parent_id, name, is_approver").
		Where("id = ?", contract.PositionID).
		Take(&current).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Position not found")
	}

	currentName := strings.TrimSpace(current.Name)
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

		parentName := strings.TrimSpace(parent.Name)

		var approver struct {
			EmployeeID string `gorm:"column:employee_id"`
		}
		// Find the approver (employee) who holds the parent position in same division.
		if err := tx.Table("employee_contracts").
			Select("employee_id").
			Where("position_id = ?", parent.ID).
			Where("division_id = ?", contract.DivisionID).
			Where("end_date IS NULL OR end_date >= ?", nowEpoch()).
			Order("start_date DESC").
			Limit(1).
			Take(&approver).Error; err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Approver not found for position")
		}

		// Avoid self-approval.
		if approver.EmployeeID != employeeID {
			approvals = append(approvals, entity.TimeOffApproval{
				ApproverId: approver.EmployeeID,
				IsRequired: parent.IsApprover,
				Status:     "PENDING",
			})
		}

		// Stop rules:
		// 1) Below Direktur Operasional -> stop at Direktur Operasional.
		// 2) Direktur Operasional -> only need Direktur Utama.
		// 3) Direktur Utama -> only need Komisaris Utama.
		if strings.EqualFold(parentName, "Direktur Operasional") && !strings.EqualFold(currentName, "Direktur Operasional") {
			break
		}
		if strings.EqualFold(currentName, "Direktur Operasional") && strings.EqualFold(parentName, "Direktur Utama") {
			break
		}
		if strings.EqualFold(currentName, "Direktur Utama") && strings.EqualFold(parentName, "Komisaris Utama") {
			break
		}
		// Fallback stop: if we hit a required approver, stop here unless current is Direktur Operasional.
		if parent.IsApprover && !strings.EqualFold(currentName, "Direktur Operasional") {
			break
		}

		// Move to the next parent.
		currentName = parentName
		parentID = parent.ParentID
	}

	// Require at least one approver.
	if len(approvals) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Approval chain is empty")
	}

	return approvals, nil
}

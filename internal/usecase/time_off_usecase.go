package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TimeOffUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	Validate              *validator.Validate
	TimeOffRequestRepo    *repository.TimeOffRequestRepository
	TimeOffTypeRepo       *repository.TimeOffTypeRepository
	TimeOffBalanceRepo    *repository.TimeOffBalanceRepository
	TimeOffApprovalRepo   *repository.TimeOffApprovalRepository
	TimeOffAttachmentRepo *repository.TimeOffAttachmentRepository
}

func NewTimeOffUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	timeOffRequestRepo *repository.TimeOffRequestRepository,
	timeOffTypeRepo *repository.TimeOffTypeRepository,
	timeOffBalanceRepo *repository.TimeOffBalanceRepository,
	timeOffApprovalRepo *repository.TimeOffApprovalRepository,
	timeOffAttachmentRepo *repository.TimeOffAttachmentRepository,
) *TimeOffUseCase {
	return &TimeOffUseCase{
		DB:                    db,
		Log:                   log,
		Validate:              validate,
		TimeOffRequestRepo:    timeOffRequestRepo,
		TimeOffTypeRepo:       timeOffTypeRepo,
		TimeOffBalanceRepo:    timeOffBalanceRepo,
		TimeOffApprovalRepo:   timeOffApprovalRepo,
		TimeOffAttachmentRepo: timeOffAttachmentRepo,
	}
}

// TODO: Validate business rules (quota, overlapping dates) before insert.
func (c *TimeOffUseCase) CreateRequest(ctx context.Context, employeeID string, request *model.CreateTimeOffRequest) (*model.TimeOffRequestResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	startDate := mustParseEpoch(request.StartDate)
	endDate := mustParseEpoch(request.EndDate)
	if startDate == 0 || endDate == 0 || startDate > endDate {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid start_date or end_date")
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
		if int(request.RequestedDays) > balance.RemainingDays {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Requested days exceed remaining balance")
		}
	}

	item := &entity.Time_Off_Requests{
		EmployeeId:    employeeID,
		TimeOffTypeId: request.TimeOffTypeID,
		RequestedDays: int(request.RequestedDays),
		StartDate:     startDate,
		EndDate:       endDate,
		RequestReason: &request.RequestReason,
		RequestStatus: "PENDING",
		CreatedAt:     nowEpoch(),
	}

	if err := c.TimeOffRequestRepo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create time off request")
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

	return &model.TimeOffRequestResponse{
		ID:            item.ID,
		EmployeeID:    item.EmployeeId,
		TimeOffTypeID: item.TimeOffTypeId,
		StartDate:     request.StartDate,
		EndDate:       request.EndDate,
		RequestedDays: request.RequestedDays,
		RequestReason: request.RequestReason,
		RequestStatus: item.RequestStatus,
		CreatedAt:     time.UnixMilli(item.CreatedAt),
	}, nil
}

// TODO: Add authorization scoping for admin vs current-user list.
func (c *TimeOffUseCase) ListRequests(ctx context.Context, request *model.SearchTimeOffRequest) ([]model.TimeOffRequestResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.TimeOffRequestRepo.List(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off requests")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffRequestResponse, len(items))
	for i, item := range items {
		requestReason := ""
		if item.RequestReason != nil {
			requestReason = *item.RequestReason
		}

		responses[i] = model.TimeOffRequestResponse{
			ID:            item.ID,
			EmployeeID:    item.EmployeeId,
			TimeOffTypeID: item.TimeOffTypeId,
			StartDate:     epochToDateString(item.StartDate),
			EndDate:       epochToDateString(item.EndDate),
			RequestedDays: float64(item.RequestedDays),
			RequestReason: requestReason,
			RequestStatus: item.RequestStatus,
			CreatedAt:     time.UnixMilli(item.CreatedAt),
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}

// TODO: Add ownership checks for non-admin users.
func (c *TimeOffUseCase) GetRequestByID(ctx context.Context, id string) (*model.TimeOffRequestResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item, err := c.TimeOffRequestRepo.FindByID(tx, id)
	if err != nil {
		c.Log.WithError(err).Error("Time off request not found")
		return nil, fiber.ErrNotFound
	}

	requestReason := ""
	if item.RequestReason != nil {
		requestReason = *item.RequestReason
	}

	response := &model.TimeOffRequestResponse{
		ID:            item.ID,
		EmployeeID:    item.EmployeeId,
		TimeOffTypeID: item.TimeOffTypeId,
		StartDate:     epochToDateString(item.StartDate),
		EndDate:       epochToDateString(item.EndDate),
		RequestedDays: float64(item.RequestedDays),
		RequestReason: requestReason,
		RequestStatus: item.RequestStatus,
		CreatedAt:     time.UnixMilli(item.CreatedAt),
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return response, nil
}

// TODO: Add company scoping if needed.
func (c *TimeOffUseCase) GetRequestOwner(ctx context.Context, id string) (string, error) {
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

// TODO: Add caching if types rarely change.
func (c *TimeOffUseCase) ListTypes(ctx context.Context) ([]model.TimeOffTypeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	items, err := c.TimeOffTypeRepo.List(tx)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off types")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffTypeResponse, len(items))
	for i, item := range items {
		responses[i] = model.TimeOffTypeResponse{
			ID:               item.ID,
			Name:             item.Name,
			Category:         item.Category,
			IsQuotaBased:     item.IsQuotaBased,
			DefaultQuotaDays: float64(item.DefaultQuotaDays),
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return responses, nil
}

// TODO: Add Create Type use case and restrict to admin users.
func (c *TimeOffUseCase) CreateType(ctx context.Context, request *model.CreateTimeOffTypeRequest) (*model.TimeOffTypeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	item := &entity.Time_Off_Type{
		Name:             request.Name,
		Category:         request.Category,
		IsQuotaBased:     request.IsQuotaBased,
		DefaultQuotaDays: int(request.DefaultQuotaDays),
	}

	if err := c.TimeOffTypeRepo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create time off type")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return &model.TimeOffTypeResponse{
		ID:               item.ID,
		Name:             item.Name,
		Category:         item.Category,
		IsQuotaBased:     item.IsQuotaBased,
		DefaultQuotaDays: float64(item.DefaultQuotaDays),
	}, nil
}

// TODO: Enforce company scoping if balances are shared across tenants.
func (c *TimeOffUseCase) ListBalances(ctx context.Context, employeeID string, request *model.SearchTimeOffBalanceRequest) ([]model.TimeOffBalanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, fiber.ErrBadRequest
	}

	items, err := c.TimeOffBalanceRepo.ListByEmployee(tx, employeeID, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off balances")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffBalanceResponse, len(items))
	for i, item := range items {
		responses[i] = model.TimeOffBalanceResponse{
			ID:            item.ID,
			EmployeeID:    item.EmployeeId,
			TimeOffTypeID: item.TimeOffTypeId,
			PeriodYear:    item.PeriodYear,
			EntitledDays:  float64(item.EntitledDays),
			UsedDays:      float64(item.UsedDays),
			RemainingDays: float64(item.RemainingDays),
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return responses, nil
}

// TODO: Add audit log for manual balance overrides.
func (c *TimeOffUseCase) SetBalance(ctx context.Context, request *model.SetTimeOffBalanceRequest) (*model.TimeOffBalanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	if _, err := c.TimeOffTypeRepo.FindByID(tx, request.TimeOffTypeID); err != nil {
		c.Log.WithError(err).Error("Time off type not found")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Time off type not found")
	}

	remainingDays := request.RemainingDays
	if remainingDays == nil {
		calculated := request.EntitledDays - request.UsedDays
		if calculated < 0 {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Remaining days cannot be negative")
		}
		remainingDays = &calculated
	}

	item, err := c.TimeOffBalanceRepo.FindByEmployeeTypeYear(tx, request.EmployeeID, request.TimeOffTypeID, request.PeriodYear)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.WithError(err).Error("Failed to check existing balance")
		return nil, fiber.ErrInternalServerError
	}

	if item == nil || err == gorm.ErrRecordNotFound {
		item = &entity.Time_Off_Balance{
			EmployeeId:    request.EmployeeID,
			TimeOffTypeId: request.TimeOffTypeID,
			PeriodYear:    request.PeriodYear,
			EntitledDays:  request.EntitledDays,
			UsedDays:      request.UsedDays,
			RemainingDays: *remainingDays,
		}
		if err := c.TimeOffBalanceRepo.Create(tx, item); err != nil {
			c.Log.WithError(err).Error("Failed to create time off balance")
			return nil, fiber.ErrInternalServerError
		}
	} else {
		item.EntitledDays = request.EntitledDays
		item.UsedDays = request.UsedDays
		item.RemainingDays = *remainingDays
		if err := c.TimeOffBalanceRepo.Update(tx, item); err != nil {
			c.Log.WithError(err).Error("Failed to update time off balance")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return &model.TimeOffBalanceResponse{
		ID:            item.ID,
		EmployeeID:    item.EmployeeId,
		TimeOffTypeID: item.TimeOffTypeId,
		PeriodYear:    item.PeriodYear,
		EntitledDays:  float64(item.EntitledDays),
		UsedDays:      float64(item.UsedDays),
		RemainingDays: float64(item.RemainingDays),
	}, nil
}

// TODO: Include approver position and division when those joins are available.
func (c *TimeOffUseCase) ListApprovals(ctx context.Context, requestID string) ([]model.TimeOffApprovalResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	type approvalRow struct {
		ID               string  `gorm:"column:id"`
		TimeOffRequestID string  `gorm:"column:time_off_request_id"`
		ApproverID       string  `gorm:"column:approver_id"`
		ApproverName     *string `gorm:"column:approver_name"`
		Status           string  `gorm:"column:approval_status"`
		ActionAt         *int64  `gorm:"column:action_at"`
		ActionReason     *string `gorm:"column:action_reason"`
	}

	var rows []approvalRow
	if err := tx.
		Table("time_off_approvals AS a").
		Select(`
			a.id,
			a.time_off_request_id,
			a.approver_id,
			a.approval_status,
			a.action_at,
			a.action_reason,
			e.fullname AS approver_name
		`).
		Joins("LEFT JOIN employees e ON e.id = a.approver_id").
		Where("a.time_off_request_id = ?", requestID).
		Order("a.id ASC").
		Find(&rows).Error; err != nil {
		c.Log.WithError(err).Error("Failed to list time off approvals")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffApprovalResponse, len(rows))
	for i, row := range rows {
		var actionAt *time.Time
		if row.ActionAt != nil && *row.ActionAt > 0 {
			t := time.UnixMilli(*row.ActionAt)
			actionAt = &t
		}

		approverName := ""
		if row.ApproverName != nil {
			approverName = *row.ApproverName
		}

		responses[i] = model.TimeOffApprovalResponse{
			ID:                 row.ID,
			TimeOffRequestID:   row.TimeOffRequestID,
			ApproverEmployeeID: row.ApproverID,
			ApproverName:       approverName,
			ApproverPosition:   "",
			ApproverDivision:   "",
			Status:             row.Status,
			ActionAt:           actionAt,
			ActionReason:       row.ActionReason,
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return responses, nil
}

// TODO: Add permission checks for approver role if needed.
func (c *TimeOffUseCase) Approve(ctx context.Context, requestID, approvalID, approverID string, request *model.ApproveTimeOffRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return fiber.ErrBadRequest
	}

	var approval entity.Time_Off_Approval
	if err := tx.
		Where("id = ? AND time_off_request_id = ? AND approver_id = ?", approvalID, requestID, approverID).
		Take(&approval).Error; err != nil {
		c.Log.WithError(err).Error("Approval not found")
		return fiber.ErrNotFound
	}

	if approval.Status == "APPROVED" || approval.Status == "REJECTED" {
		return fiber.NewError(fiber.StatusBadRequest, "Approval already processed")
	}

	updates := map[string]any{
		"approval_status": "APPROVED",
		"action_reason":   request.ActionReason,
		"action_at":       nowEpoch(),
	}
	if err := tx.Table("time_off_approvals").Where("id = ?", approval.ID).Updates(updates).Error; err != nil {
		c.Log.WithError(err).Error("Failed to approve time off request")
		return fiber.ErrInternalServerError
	}

	var pendingCount int64
	if err := tx.
		Table("time_off_approvals").
		Where("time_off_request_id = ? AND approval_status = ?", requestID, "PENDING").
		Count(&pendingCount).Error; err != nil {
		c.Log.WithError(err).Error("Failed to count pending approvals")
		return fiber.ErrInternalServerError
	}

	if pendingCount == 0 {
		if err := tx.Table("time_off_requests").
			Where("id = ?", requestID).
			Update("request_status", "APPROVED").Error; err != nil {
			c.Log.WithError(err).Error("Failed to update time off request status")
			return fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

// TODO: Consider audit logging for short approve path.
func (c *TimeOffUseCase) ApproveByApprovalID(ctx context.Context, approvalID, approverID string, request *model.ApproveTimeOffRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return fiber.ErrBadRequest
	}

	var approval entity.Time_Off_Approval
	if err := tx.
		Where("id = ? AND approver_id = ?", approvalID, approverID).
		Take(&approval).Error; err != nil {
		c.Log.WithError(err).Error("Approval not found")
		return fiber.ErrNotFound
	}

	if approval.Status == "APPROVED" || approval.Status == "REJECTED" {
		return fiber.NewError(fiber.StatusBadRequest, "Approval already processed")
	}

	updates := map[string]any{
		"approval_status": "APPROVED",
		"action_reason":   request.ActionReason,
		"action_at":       nowEpoch(),
	}
	if err := tx.Table("time_off_approvals").Where("id = ?", approval.ID).Updates(updates).Error; err != nil {
		c.Log.WithError(err).Error("Failed to approve time off request")
		return fiber.ErrInternalServerError
	}

	var pendingCount int64
	if err := tx.
		Table("time_off_approvals").
		Where("time_off_request_id = ? AND approval_status = ?", approval.TimeOffRequestId, "PENDING").
		Count(&pendingCount).Error; err != nil {
		c.Log.WithError(err).Error("Failed to count pending approvals")
		return fiber.ErrInternalServerError
	}

	if pendingCount == 0 {
		if err := tx.Table("time_off_requests").
			Where("id = ?", approval.TimeOffRequestId).
			Update("request_status", "APPROVED").Error; err != nil {
			c.Log.WithError(err).Error("Failed to update time off request status")
			return fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

// TODO: Add permission checks for approver role if needed.
func (c *TimeOffUseCase) Reject(ctx context.Context, requestID, approvalID, approverID string, request *model.RejectTimeOffRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return fiber.ErrBadRequest
	}

	var approval entity.Time_Off_Approval
	if err := tx.
		Where("id = ? AND time_off_request_id = ? AND approver_id = ?", approvalID, requestID, approverID).
		Take(&approval).Error; err != nil {
		c.Log.WithError(err).Error("Approval not found")
		return fiber.ErrNotFound
	}

	if approval.Status == "APPROVED" || approval.Status == "REJECTED" {
		return fiber.NewError(fiber.StatusBadRequest, "Approval already processed")
	}

	updates := map[string]any{
		"approval_status": "REJECTED",
		"action_reason":   request.ActionReason,
		"action_at":       nowEpoch(),
	}
	if err := tx.Table("time_off_approvals").Where("id = ?", approval.ID).Updates(updates).Error; err != nil {
		c.Log.WithError(err).Error("Failed to reject time off request")
		return fiber.ErrInternalServerError
	}

	if err := tx.Table("time_off_requests").
		Where("id = ?", requestID).
		Update("request_status", "REJECTED").Error; err != nil {
		c.Log.WithError(err).Error("Failed to update time off request status")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

// TODO: Consider audit logging for short reject path.
func (c *TimeOffUseCase) RejectByApprovalID(ctx context.Context, approvalID, approverID string, request *model.RejectTimeOffRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return fiber.ErrBadRequest
	}

	var approval entity.Time_Off_Approval
	if err := tx.
		Where("id = ? AND approver_id = ?", approvalID, approverID).
		Take(&approval).Error; err != nil {
		c.Log.WithError(err).Error("Approval not found")
		return fiber.ErrNotFound
	}

	if approval.Status == "APPROVED" || approval.Status == "REJECTED" {
		return fiber.NewError(fiber.StatusBadRequest, "Approval already processed")
	}

	updates := map[string]any{
		"approval_status": "REJECTED",
		"action_reason":   request.ActionReason,
		"action_at":       nowEpoch(),
	}
	if err := tx.Table("time_off_approvals").Where("id = ?", approval.ID).Updates(updates).Error; err != nil {
		c.Log.WithError(err).Error("Failed to reject time off request")
		return fiber.ErrInternalServerError
	}

	if err := tx.Table("time_off_requests").
		Where("id = ?", approval.TimeOffRequestId).
		Update("request_status", "REJECTED").Error; err != nil {
		c.Log.WithError(err).Error("Failed to update time off request status")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

// TODO: Add authorization checks when admin views all approvals.
func (c *TimeOffUseCase) ListApprovalsByApprover(ctx context.Context, approverID string, request *model.SearchTimeOffApprovalRequest) ([]model.TimeOffApprovalResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, 0, fiber.ErrBadRequest
	}

	type approvalRow struct {
		ID               string  `gorm:"column:id"`
		TimeOffRequestID string  `gorm:"column:time_off_request_id"`
		ApproverID       string  `gorm:"column:approver_id"`
		ApproverName     *string `gorm:"column:approver_name"`
		Status           string  `gorm:"column:approval_status"`
		ActionAt         *int64  `gorm:"column:action_at"`
		ActionReason     *string `gorm:"column:action_reason"`
	}

	query := tx.Table("time_off_approvals AS a").
		Select(`
			a.id,
			a.time_off_request_id,
			a.approver_id,
			a.approval_status,
			a.action_at,
			a.action_reason,
			e.fullname AS approver_name
		`).
		Joins("LEFT JOIN employees e ON e.id = a.approver_id").
		Where("a.approver_id = ?", approverID)

	if request.Status != "" {
		query = query.Where("a.approval_status = ?", request.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.Log.WithError(err).Error("Failed to count approvals")
		return nil, 0, fiber.ErrInternalServerError
	}

	var rows []approvalRow
	if err := query.
		Order("a.action_at IS NULL DESC, a.action_at DESC, a.id DESC").
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&rows).Error; err != nil {
		c.Log.WithError(err).Error("Failed to list approvals")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffApprovalResponse, len(rows))
	for i, row := range rows {
		var actionAt *time.Time
		if row.ActionAt != nil && *row.ActionAt > 0 {
			t := time.UnixMilli(*row.ActionAt)
			actionAt = &t
		}

		approverName := ""
		if row.ApproverName != nil {
			approverName = *row.ApproverName
		}

		responses[i] = model.TimeOffApprovalResponse{
			ID:                 row.ID,
			TimeOffRequestID:   row.TimeOffRequestID,
			ApproverEmployeeID: row.ApproverID,
			ApproverName:       approverName,
			ApproverPosition:   "",
			ApproverDivision:   "",
			Status:             row.Status,
			ActionAt:           actionAt,
			ActionReason:       row.ActionReason,
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}

// TODO: Consider file validation (mime/size) and ownership checks.
func (c *TimeOffUseCase) ListAttachments(ctx context.Context, requestID string) ([]model.TimeOffAttachmentResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	items, err := c.TimeOffAttachmentRepo.ListByRequestID(tx, requestID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off attachments")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffAttachmentResponse, len(items))
	for i, item := range items {
		responses[i] = model.TimeOffAttachmentResponse{
			ID:               item.ID,
			TimeOffRequestID: item.TimeOffRequestId,
			FileName:         item.FileName,
			MimeType:         item.MimeType,
			FileSize:         item.FileSize,
			FileUrl:          item.FileUrl,
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return responses, nil
}

// TODO: Enforce request ownership and approval status before allowing attachments.
func (c *TimeOffUseCase) CreateAttachment(ctx context.Context, requestID string, request *model.CreateTimeOffAttachmentRequest) (*model.TimeOffAttachmentResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	if _, err := c.TimeOffRequestRepo.FindByID(tx, requestID); err != nil {
		c.Log.WithError(err).Error("Time off request not found")
		return nil, fiber.ErrNotFound
	}

	item := &entity.Time_Off_Attachment{
		TimeOffRequestId: requestID,
		FileName:         request.FileName,
		MimeType:         request.MimeType,
		FileSize:         request.FileSize,
		FileUrl:          request.FileUrl,
	}

	if err := c.TimeOffAttachmentRepo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create time off attachment")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return &model.TimeOffAttachmentResponse{
		ID:               item.ID,
		TimeOffRequestID: item.TimeOffRequestId,
		FileName:         item.FileName,
		MimeType:         item.MimeType,
		FileSize:         item.FileSize,
		FileUrl:          item.FileUrl,
	}, nil
}

// TODO: Replace with a proper date parsing strategy shared across the codebase.
// TODO: Handle parse error explicitly instead of returning 0.
func mustParseEpoch(date string) int64 {
	parsed, _ := lib.ParseDateToUnixMilli(date)
	return parsed
}

// TODO: Centralize time source for deterministic tests.
func nowEpoch() int64 {
	return time.Now().UnixMilli()
}

// TODO: Consider timezone handling if date is company-local.
func epochToDateString(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.UnixMilli(ts).Format("2006-01-02")
}

// TODO: Add cycle detection limit configuration.
func (c *TimeOffUseCase) buildApprovalsFromPositionChain(tx *gorm.DB, employeeID string) ([]entity.Time_Off_Approval, error) {
	// Fetch active contract to get current position + division.
	type contractRow struct {
		EmployeeID string `gorm:"column:employee_id"`
		DivisionID string `gorm:"column:division_id"`
		PositionID string `gorm:"column:position_id"`
		StartDate  int64  `gorm:"column:start_date"`
		EndDate    *int64 `gorm:"column:end_date"`
	}

	var contract contractRow
	now := nowEpoch()
	if err := tx.Table("employee_contracts").
		Select("employee_id, division_id, position_id, start_date, end_date").
		Where("employee_id = ?", employeeID).
		Where("end_date IS NULL OR end_date >= ?", now).
		Order("start_date DESC").
		Limit(1).
		Take(&contract).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Active contract not found")
	}

	// Load current position and start from its parent.
	type positionRow struct {
		ID       string  `gorm:"column:id"`
		ParentID *string `gorm:"column:parent_id"`
	}

	approvals := make([]entity.Time_Off_Approval, 0, 4)
	var current positionRow
	if err := tx.Table("positions").
		Select("id, parent_id").
		Where("id = ?", contract.PositionID).
		Take(&current).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Position not found")
	}

	parentID := current.ParentID

	const maxDepth = 20
	// Walk up the position hierarchy until root or depth limit.
	for depth := 0; depth < maxDepth && parentID != nil; depth++ {
		var parent positionRow
		// Resolve parent position.
		if err := tx.Table("positions").
			Select("id, parent_id").
			Where("id = ?", *parentID).
			Take(&parent).Error; err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Parent position not found")
		}

		var approver struct {
			EmployeeID string `gorm:"column:employee_id"`
		}
		// Find the approver (employee) who holds the parent position in same division.
		if err := tx.Table("employee_contracts").
			Select("employee_id").
			Where("position_id = ?", parent.ID).
			Where("division_id = ?", contract.DivisionID).
			Where("end_date IS NULL OR end_date >= ?", now).
			Order("start_date DESC").
			Limit(1).
			Take(&approver).Error; err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Approver not found for position")
		}

		// Avoid self-approval.
		if approver.EmployeeID != employeeID {
			approvals = append(approvals, entity.Time_Off_Approval{
				ApproverId: approver.EmployeeID,
			})
		}

		// Move to the next parent.
		parentID = parent.ParentID
	}

	// Require at least one approver.
	if len(approvals) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Approval chain is empty")
	}

	return approvals, nil
}

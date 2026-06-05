package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TimeOffApprovalUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	TimeOffRequestRepo *repository.TimeOffRequestRepository
	TimeOffTypeRepo    *repository.TimeOffTypeRepository
	TimeOffBalanceRepo *repository.TimeOffBalanceRepository
}

func NewTimeOffApprovalUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	timeOffRequestRepo *repository.TimeOffRequestRepository,
	timeOffTypeRepo *repository.TimeOffTypeRepository,
	timeOffBalanceRepo *repository.TimeOffBalanceRepository,
) *TimeOffApprovalUseCase {
	return &TimeOffApprovalUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		TimeOffRequestRepo: timeOffRequestRepo,
		TimeOffTypeRepo:    timeOffTypeRepo,
		TimeOffBalanceRepo: timeOffBalanceRepo,
	}
}

// TODO: Include approver position and division when those joins are available.
func (c *TimeOffApprovalUseCase) ListApprovals(ctx context.Context, requestID string) ([]model.TimeOffApprovalResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var items []entity.TimeOffApproval
	if err := tx.
		Preload("TimeOffRequest").
		Preload("TimeOffRequest.Employee").
		Preload("TimeOffRequest.TimeOffType").
		Where("time_off_request_id = ?", requestID).
		Order("id ASC").
		Find(&items).Error; err != nil {
		c.Log.WithError(err).Error("Failed to list time off approvals")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffApprovalResponse, len(items))
	for i := range items {
		responses[i] = *model.TimeOffApprovalToResponse(&items[i])
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return responses, nil
}

// TODO: Add permission checks for approver role if needed.
func (c *TimeOffApprovalUseCase) Decide(ctx context.Context, requestID, approvalID, approverID string, request *model.DecideTimeOffApprovalRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.ActionAt == 0 {
		request.ActionAt = nowEpoch()
	}
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return fiber.ErrBadRequest
	}

	action := strings.ToUpper(strings.TrimSpace(request.Action))
	if action != "APPROVE" && action != "REJECT" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid action")
	}
	if action == "REJECT" && strings.TrimSpace(request.ActionReason) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Action reason is required")
	}

	var approval entity.TimeOffApproval
	if err := tx.
		Where("id = ? AND time_off_request_id = ? AND approver_id = ?", approvalID, requestID, approverID).
		Take(&approval).Error; err != nil {
		c.Log.WithError(err).Error("Approval not found")
		return fiber.ErrNotFound
	}

	if approval.Status == "APPROVED" || approval.Status == "REJECTED" {
		return fiber.NewError(fiber.StatusBadRequest, "Approval already processed")
	}

	actionStatus := "APPROVED"
	if action == "REJECT" {
		actionStatus = "REJECTED"
	}
	updates := map[string]any{
		"approval_status": actionStatus,
		"action_reason":   request.ActionReason,
		"action_at":       request.ActionAt,
	}
	if err := tx.Table("time_off_approvals").Where("id = ?", approval.ID).Updates(updates).Error; err != nil {
		c.Log.WithError(err).Error("Failed to update time off approval")
		return fiber.ErrInternalServerError
	}

	if action == "APPROVE" {
		var isApprover bool
		if err := tx.Raw(`
			SELECT p.is_approver
			FROM employee_contracts ec
			JOIN positions p ON p.id = ec.position_id
			WHERE ec.employee_id = ? AND ec.is_active = true
			LIMIT 1
		`, approverID).Scan(&isApprover).Error; err != nil {
			c.Log.WithError(err).Error("Failed to check approver position")
			return fiber.ErrInternalServerError
		}

		if isApprover {
			var pendingRequired int64
			if err := tx.Table("time_off_approvals").
				Where("time_off_request_id = ?", requestID).
				Where("is_required = ?", true).
				Where("approval_status != ?", "APPROVED").
				Count(&pendingRequired).Error; err != nil {
				c.Log.WithError(err).Error("Failed to count pending approvals")
				return fiber.ErrInternalServerError
			}
			if pendingRequired == 0 {
				if err := tx.Table("time_off_requests").
					Where("id = ?", requestID).
					Update("request_status", "APPROVED").Error; err != nil {
					c.Log.WithError(err).Error("Failed to update time off request status")
					return fiber.ErrInternalServerError
				}
				if err := c.applyBalanceOnApproval(tx, requestID); err != nil {
					c.Log.WithError(err).Error("Failed to update time off balance")
					return err
				}
			}
		}
	} else {
		if err := tx.Table("time_off_requests").
			Where("id = ?", requestID).
			Update("request_status", "REJECTED").Error; err != nil {
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

// TODO: Consider audit logging for short decide path.
func (c *TimeOffApprovalUseCase) DecideByApprovalID(ctx context.Context, approvalID, approverID string, request *model.DecideTimeOffApprovalRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.ActionAt == 0 {
		request.ActionAt = nowEpoch()
	}
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return fiber.ErrBadRequest
	}

	action := strings.ToUpper(strings.TrimSpace(request.Action))
	if action != "APPROVE" && action != "REJECT" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid action")
	}
	if action == "REJECT" && strings.TrimSpace(request.ActionReason) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Action reason is required")
	}

	var approval entity.TimeOffApproval
	if err := tx.
		Where("id = ? AND approver_id = ?", approvalID, approverID).
		Take(&approval).Error; err != nil {
		c.Log.WithError(err).Error("Approval not found")
		return fiber.ErrNotFound
	}

	if approval.Status == "APPROVED" || approval.Status == "REJECTED" {
		return fiber.NewError(fiber.StatusBadRequest, "Approval already processed")
	}

	actionStatus := "APPROVED"
	if action == "REJECT" {
		actionStatus = "REJECTED"
	}
	updates := map[string]any{
		"approval_status": actionStatus,
		"action_reason":   request.ActionReason,
		"action_at":       request.ActionAt,
	}
	if err := tx.Table("time_off_approvals").Where("id = ?", approval.ID).Updates(updates).Error; err != nil {
		c.Log.WithError(err).Error("Failed to update time off approval")
		return fiber.ErrInternalServerError
	}

	if action == "APPROVE" {
		var isApprover bool
		if err := tx.Raw(`
			SELECT p.is_approver
			FROM employee_contracts ec
			JOIN positions p ON p.id = ec.position_id
			WHERE ec.employee_id = ? AND ec.is_active = true
			LIMIT 1
		`, approverID).Scan(&isApprover).Error; err != nil {
			c.Log.WithError(err).Error("Failed to check approver position")
			return fiber.ErrInternalServerError
		}

		if isApprover {
			var pendingRequired int64
			if err := tx.Table("time_off_approvals").
				Where("time_off_request_id = ?", approval.TimeOffRequestId).
				Where("is_required = ?", true).
				Where("approval_status != ?", "APPROVED").
				Count(&pendingRequired).Error; err != nil {
				c.Log.WithError(err).Error("Failed to count pending approvals")
				return fiber.ErrInternalServerError
			}
			if pendingRequired == 0 {
				if err := tx.Table("time_off_requests").
					Where("id = ?", approval.TimeOffRequestId).
					Update("request_status", "APPROVED").Error; err != nil {
					c.Log.WithError(err).Error("Failed to update time off request status")
					return fiber.ErrInternalServerError
				}
				if err := c.applyBalanceOnApproval(tx, approval.TimeOffRequestId); err != nil {
					c.Log.WithError(err).Error("Failed to update time off balance")
					return err
				}
			}
		}
	} else {
		if err := tx.Table("time_off_requests").
			Where("id = ?", approval.TimeOffRequestId).
			Update("request_status", "REJECTED").Error; err != nil {
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

// TODO: Add authorization checks when admin views all approvals.
func (c *TimeOffApprovalUseCase) ListApprovalsByApprover(ctx context.Context, approverID string, request *model.SearchTimeOffApprovalRequest) ([]model.TimeOffApprovalResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, 0, fiber.ErrBadRequest
	}

	query := tx.Model(&entity.TimeOffApproval{}).
		Where("approver_id = ?", approverID).
		Where(`NOT EXISTS (
			SELECT 1 FROM time_off_approvals prev
			WHERE prev.time_off_request_id = time_off_approvals.time_off_request_id
			AND prev.approval_order < time_off_approvals.approval_order
			AND prev.approval_status != 'APPROVED')`).
		Preload("Employee").
		Preload("TimeOffRequest").
		Preload("TimeOffRequest.Employee").
		Preload("TimeOffRequest.Employee.EmployeeContract").
		Preload("TimeOffRequest.Employee.EmployeeContract.Position").
		Preload("TimeOffRequest.Employee.EmployeeContract.Division").
		Preload("TimeOffRequest.TimeOffType").
		Preload("TimeOffRequest.Approvals").
		Preload("TimeOffRequest.Approvals.Employee")

	if request.Status != "" {
		query = query.Where("approval_status = ?", request.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.Log.WithError(err).Error("Failed to count approvals")
		return nil, 0, fiber.ErrInternalServerError
	}

	var rows []entity.TimeOffApproval
	if err := query.
		Order("action_at IS NULL DESC, action_at DESC, id DESC").
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&rows).Error; err != nil {
		c.Log.WithError(err).Error("Failed to list approvals")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffApprovalResponse, len(rows))
	for i := range rows {
		responses[i] = *model.TimeOffApprovalToResponse(&rows[i])
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}

func (c *TimeOffApprovalUseCase) applyBalanceOnApproval(tx *gorm.DB, requestID string) error {
	request, err := c.TimeOffRequestRepo.FindByID(tx, requestID, false)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Time off request not found")
	}

	timeOffType, err := c.TimeOffTypeRepo.FindByID(tx, request.TimeOffTypeId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Time off type not found")
	}
	if !timeOffType.IsQuotaBased {
		return nil
	}

	periodYear := time.UnixMilli(request.StartDate).UTC().Year()
	balance, err := c.TimeOffBalanceRepo.FindByEmployeeTypeYear(tx, request.EmployeeID, request.TimeOffTypeId, periodYear)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Time off balance not found")
	}

	newUsed := balance.UsedDays + request.RequestedDays
	newRemaining := balance.EntitledDays - newUsed
	if newRemaining < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Requested days exceed remaining balance")
	}

	balance.UsedDays = newUsed
	balance.RemainingDays = newRemaining
	if err := c.TimeOffBalanceRepo.Update(tx, balance); err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

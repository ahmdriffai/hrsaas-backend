package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TimeOffApprovalController struct {
	ApprovalUseCase *usecase.TimeOffApprovalUseCase
	RequestUseCase  *usecase.TimeOffRequestUseCase
	Log             *logrus.Logger
}

func NewTimeOffApprovalController(
	approvalUseCase *usecase.TimeOffApprovalUseCase,
	requestUseCase *usecase.TimeOffRequestUseCase,
	log *logrus.Logger,
) *TimeOffApprovalController {
	return &TimeOffApprovalController{
		ApprovalUseCase: approvalUseCase,
		RequestUseCase:  requestUseCase,
		Log:             log,
	}
}

// TODO: Restrict access to request owner and approvers.
func (c *TimeOffApprovalController) ListApprovals(ctx *fiber.Ctx) error {
	requestID := ctx.Params("id")
	if requestID == "" {
		return fiber.ErrBadRequest
	}

	if err := ensureOwnerOrAdmin(ctx, c.RequestUseCase, requestID); err != nil {
		return err
	}

	responses, err := c.ApprovalUseCase.ListApprovals(ctx.UserContext(), requestID)
	if err != nil {
		c.Log.WithError(err).Error("failed to list time off approvals")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffApprovalResponse]{
		Data: responses,
	})
}

// TODO: Add audit logging for approval actions.
func (c *TimeOffApprovalController) Decide(ctx *fiber.Ctx) error {
	requestID := ctx.Params("id")
	approvalID := ctx.Params("approval_id")
	if requestID == "" || approvalID == "" {
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request := new(model.DecideTimeOffApprovalRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	if err := c.ApprovalUseCase.Decide(ctx.UserContext(), requestID, approvalID, user.Employee.ID, request); err != nil {
		c.Log.WithError(err).Error("failed to update time off approval")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

// TODO: Add filters for request status and date range if needed.
func (c *TimeOffApprovalController) ListMyApprovals(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request := new(model.SearchTimeOffApprovalRequest)
	request.Status = ctx.Query("status", "")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.ApprovalUseCase.ListApprovalsByApprover(ctx.UserContext(), user.Employee.ID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list approvals for current user")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffApprovalResponse]{
		Data:   responses,
		Paging: paging,
	})
}

// TODO: Add audit logging for approval actions.
func (c *TimeOffApprovalController) DecideShort(ctx *fiber.Ctx) error {
	approvalID := ctx.Params("approval_id")
	if approvalID == "" {
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request := new(model.DecideTimeOffApprovalRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	if err := c.ApprovalUseCase.DecideByApprovalID(ctx.UserContext(), approvalID, user.Employee.ID, request); err != nil {
		c.Log.WithError(err).Error("failed to update time off approval")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

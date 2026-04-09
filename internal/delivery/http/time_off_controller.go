package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TimeOffController struct {
	UseCase *usecase.TimeOffUseCase
	Log     *logrus.Logger
}

func NewTimeOffController(useCase *usecase.TimeOffUseCase, log *logrus.Logger) *TimeOffController {
	return &TimeOffController{
		UseCase: useCase,
		Log:     log,
	}
}

// TODO: Enforce role-based access (employee only) at middleware or here.
func (c *TimeOffController) CreateRequest(ctx *fiber.Ctx) error {
	request := new(model.CreateTimeOffRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	response, err := c.UseCase.CreateRequest(ctx.UserContext(), user.Employee.ID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create time off request")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TimeOffRequestResponse]{
		Data: response,
	})
}

func (c *TimeOffController) GetRequestByID(ctx *fiber.Ctx) error {
	requestID := ctx.Params("id")
	if requestID == "" {
		return fiber.ErrBadRequest
	}

	if err := c.ensureOwnerOrAdmin(ctx, requestID); err != nil {
		return err
	}

	response, err := c.UseCase.GetRequestByID(ctx.UserContext(), requestID)
	if err != nil {
		c.Log.WithError(err).Error("failed to get time off request detail")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TimeOffRequestResponse]{
		Data: response,
	})
}

// TODO: Add admin-only filters and company scoping.
func (c *TimeOffController) ListRequests(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	if !strings.EqualFold(user.Role, "ADMIN") {
		return fiber.NewError(fiber.StatusForbidden, "Forbidden")
	}

	request := new(model.SearchTimeOffRequest)
	request.EmployeeID = ctx.Query("employee_id", "")
	request.TimeOffTypeID = ctx.Query("time_off_type_id", "")
	request.RequestStatus = ctx.Query("request_status", "")
	request.StartDate = ctx.Query("start_date", "")
	request.EndDate = ctx.Query("end_date", "")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.ListRequests(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list time off requests")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffRequestResponse]{
		Data:   responses,
		Paging: paging,
	})
}

// TODO: Return 403 if user has no employee profile.
func (c *TimeOffController) ListCurrentRequests(ctx *fiber.Ctx) error {
	request := new(model.SearchTimeOffRequest)
	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request.EmployeeID = user.Employee.ID
	request.TimeOffTypeID = ctx.Query("time_off_type_id", "")
	request.RequestStatus = ctx.Query("request_status", "")
	request.StartDate = ctx.Query("start_date", "")
	request.EndDate = ctx.Query("end_date", "")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.ListRequests(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list current user time off requests")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffRequestResponse]{
		Data:   responses,
		Paging: paging,
	})
}

// TODO: Support pagination if types grow large.
func (c *TimeOffController) ListTypes(ctx *fiber.Ctx) error {
	responses, err := c.UseCase.ListTypes(ctx.UserContext())
	if err != nil {
		c.Log.WithError(err).Error("failed to list time off types")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffTypeResponse]{
		Data: responses,
	})
}

// TODO: Enforce admin-only access with middleware at router.
func (c *TimeOffController) CreateType(ctx *fiber.Ctx) error {
	request := new(model.CreateTimeOffTypeRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.CreateType(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create time off type")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TimeOffTypeResponse]{
		Data: response,
	})
}

// TODO: Support filtering by period range if needed.
func (c *TimeOffController) ListCurrentBalances(ctx *fiber.Ctx) error {
	request := new(model.SearchTimeOffBalanceRequest)
	request.TimeOffTypeID = ctx.Query("time_off_type_id", "")
	request.PeriodYear = ctx.QueryInt("period_year", 0)

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	responses, err := c.UseCase.ListBalances(ctx.UserContext(), user.Employee.ID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list time off balances")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffBalanceResponse]{
		Data: responses,
	})
}

// TODO: Restrict access to request owner and approvers.
func (c *TimeOffController) ListApprovals(ctx *fiber.Ctx) error {
	requestID := ctx.Params("id")
	if requestID == "" {
		return fiber.ErrBadRequest
	}

	if err := c.ensureOwnerOrAdmin(ctx, requestID); err != nil {
		return err
	}

	responses, err := c.UseCase.ListApprovals(ctx.UserContext(), requestID)
	if err != nil {
		c.Log.WithError(err).Error("failed to list time off approvals")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffApprovalResponse]{
		Data: responses,
	})
}

// TODO: Add audit logging for approval actions.
func (c *TimeOffController) Approve(ctx *fiber.Ctx) error {
	requestID := ctx.Params("id")
	approvalID := ctx.Params("approval_id")
	if requestID == "" || approvalID == "" {
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request := new(model.ApproveTimeOffRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	if err := c.UseCase.Approve(ctx.UserContext(), requestID, approvalID, user.Employee.ID, request); err != nil {
		c.Log.WithError(err).Error("failed to approve time off request")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

// TODO: Add filters for request status and date range if needed.
func (c *TimeOffController) ListMyApprovals(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request := new(model.SearchTimeOffApprovalRequest)
	request.Status = ctx.Query("status", "")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.ListApprovalsByApprover(ctx.UserContext(), user.Employee.ID, request)
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
func (c *TimeOffController) ApproveShort(ctx *fiber.Ctx) error {
	approvalID := ctx.Params("approval_id")
	if approvalID == "" {
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request := new(model.ApproveTimeOffRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	if err := c.UseCase.ApproveByApprovalID(ctx.UserContext(), approvalID, user.Employee.ID, request); err != nil {
		c.Log.WithError(err).Error("failed to approve time off request")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

// TODO: Add audit logging for reject actions.
func (c *TimeOffController) Reject(ctx *fiber.Ctx) error {
	requestID := ctx.Params("id")
	approvalID := ctx.Params("approval_id")
	if requestID == "" || approvalID == "" {
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request := new(model.RejectTimeOffRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	if err := c.UseCase.Reject(ctx.UserContext(), requestID, approvalID, user.Employee.ID, request); err != nil {
		c.Log.WithError(err).Error("failed to reject time off request")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

// TODO: Add audit logging for reject actions.
func (c *TimeOffController) RejectShort(ctx *fiber.Ctx) error {
	approvalID := ctx.Params("approval_id")
	if approvalID == "" {
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request := new(model.RejectTimeOffRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	if err := c.UseCase.RejectByApprovalID(ctx.UserContext(), approvalID, user.Employee.ID, request); err != nil {
		c.Log.WithError(err).Error("failed to reject time off request")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

// TODO: Restrict to request owner and approvers.
func (c *TimeOffController) ListAttachments(ctx *fiber.Ctx) error {
	requestID := ctx.Params("id")
	if requestID == "" {
		return fiber.ErrBadRequest
	}

	if err := c.ensureOwnerOrAdmin(ctx, requestID); err != nil {
		return err
	}

	responses, err := c.UseCase.ListAttachments(ctx.UserContext(), requestID)
	if err != nil {
		c.Log.WithError(err).Error("failed to list time off attachments")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffAttachmentResponse]{
		Data: responses,
	})
}

// TODO: Validate request status before allowing upload.
func (c *TimeOffController) CreateAttachment(ctx *fiber.Ctx) error {
	requestID := ctx.Params("id")
	if requestID == "" {
		return fiber.ErrBadRequest
	}

	if err := c.ensureOwnerOrAdmin(ctx, requestID); err != nil {
		return err
	}

	request := new(model.CreateTimeOffAttachmentRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.CreateAttachment(ctx.UserContext(), requestID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create time off attachment")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TimeOffAttachmentResponse]{
		Data: response,
	})
}

func (c *TimeOffController) ensureOwnerOrAdmin(ctx *fiber.Ctx, requestID string) error {
	user := middleware.GetUser(ctx)
	if strings.EqualFold(user.Role, "ADMIN") {
		return nil
	}
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	ownerID, err := c.UseCase.GetRequestOwner(ctx.UserContext(), requestID)
	if err != nil {
		return err
	}
	if ownerID != user.Employee.ID {
		return fiber.NewError(fiber.StatusForbidden, "Forbidden")
	}
	return nil
}

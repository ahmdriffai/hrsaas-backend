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

type TimeOffRequestController struct {
	RequestUseCase *usecase.TimeOffRequestUseCase
	Log            *logrus.Logger
}

func NewTimeOffRequestController(requestUseCase *usecase.TimeOffRequestUseCase, log *logrus.Logger) *TimeOffRequestController {
	return &TimeOffRequestController{
		RequestUseCase: requestUseCase,
		Log:            log,
	}
}

// TODO: Enforce role-based access (employee only) at middleware or here.
func (c *TimeOffRequestController) CreateRequest(ctx *fiber.Ctx) error {
	request := new(model.CreateTimeOffRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	response, err := c.RequestUseCase.CreateRequest(ctx.UserContext(), user.Employee.ID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create time off request")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TimeOffRequestResponse]{
		Data: response,
	})
}

func (c *TimeOffRequestController) GetRequestByID(ctx *fiber.Ctx) error {
	requestID := ctx.Params("id")
	if requestID == "" {
		return fiber.ErrBadRequest
	}

	if err := ensureOwnerOrAdmin(ctx, c.RequestUseCase, requestID); err != nil {
		return err
	}

	response, err := c.RequestUseCase.GetRequestByID(ctx.UserContext(), requestID)
	if err != nil {
		c.Log.WithError(err).Error("failed to get time off request detail")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TimeOffRequestResponse]{
		Data: response,
	})
}

// TODO: Add admin-only filters and company scoping.
func (c *TimeOffRequestController) ListRequests(ctx *fiber.Ctx) error {
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

	responses, total, err := c.RequestUseCase.ListRequests(ctx.UserContext(), request)
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
func (c *TimeOffRequestController) ListCurrentRequests(ctx *fiber.Ctx) error {
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

	responses, total, err := c.RequestUseCase.ListRequests(ctx.UserContext(), request)
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

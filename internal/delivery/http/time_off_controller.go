package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

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

func (c *TimeOffController) ListRequests(ctx *fiber.Ctx) error {
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

package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TimeOffBalanceController struct {
	BalanceUseCase *usecase.TimeOffBalanceUseCase
	Log            *logrus.Logger
}

func NewTimeOffBalanceController(balanceUseCase *usecase.TimeOffBalanceUseCase, log *logrus.Logger) *TimeOffBalanceController {
	return &TimeOffBalanceController{
		BalanceUseCase: balanceUseCase,
		Log:            log,
	}
}

// TODO: Support filtering by period range if needed.
func (c *TimeOffBalanceController) ListCurrentBalances(ctx *fiber.Ctx) error {
	request := new(model.SearchTimeOffBalanceRequest)
	request.TimeOffTypeID = ctx.Query("time_off_type_id", "")
	request.PeriodYear = ctx.QueryInt("period_year", 0)

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	responses, err := c.BalanceUseCase.ListBalances(ctx.UserContext(), user.Employee.ID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list time off balances")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffBalanceResponse]{
		Data: responses,
	})
}

// TODO: Enforce admin-only access with middleware at router.
func (c *TimeOffBalanceController) ListBalancesByEmployee(ctx *fiber.Ctx) error {
	request := new(model.SearchTimeOffBalanceRequest)
	request.TimeOffTypeID = ctx.Query("time_off_type_id", "")
	request.PeriodYear = ctx.QueryInt("period_year", 0)

	employeeID := ctx.Query("employee_id", "")
	if employeeID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "employee_id is required")
	}

	responses, err := c.BalanceUseCase.ListBalances(ctx.UserContext(), employeeID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list time off balances")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffBalanceResponse]{
		Data: responses,
	})
}

// TODO: Enforce admin-only access with middleware at router.
func (c *TimeOffBalanceController) SetBalance(ctx *fiber.Ctx) error {
	request := new(model.SetTimeOffBalanceRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.BalanceUseCase.SetBalance(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to set time off balance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TimeOffBalanceResponse]{
		Data: response,
	})
}

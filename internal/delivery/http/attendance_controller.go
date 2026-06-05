package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AttendanceController struct {
	UseCase *usecase.AttendanceUseCase
	Log     *logrus.Logger
}

func NewAttendanceController(useCase *usecase.AttendanceUseCase, log *logrus.Logger) *AttendanceController {
	return &AttendanceController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *AttendanceController) parseRequest(ctx *fiber.Ctx) (*model.CheckInAttendanceRequest, error) {
	user := middleware.GetUser(ctx)
	request := new(model.CheckInAttendanceRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return nil, fiber.ErrBadRequest
	}
	request.CompanyID = user.CompanyID
	request.EmployeeID = user.Employee.ID
	return request, nil
}

func (c *AttendanceController) CheckIn(ctx *fiber.Ctx) error {
	request, err := c.parseRequest(ctx)
	if err != nil {
		return err
	}

	response, err := c.UseCase.CheckIn(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to check in attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

func (c *AttendanceController) CheckOut(ctx *fiber.Ctx) error {
	request, err := c.parseRequest(ctx)
	if err != nil {
		return err
	}

	response, err := c.UseCase.CheckOut(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to check out attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

func (c *AttendanceController) BreakIn(ctx *fiber.Ctx) error {
	request, err := c.parseRequest(ctx)
	if err != nil {
		return err
	}

	response, err := c.UseCase.BreakIn(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to break in attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

func (c *AttendanceController) BreakOut(ctx *fiber.Ctx) error {
	request, err := c.parseRequest(ctx)
	if err != nil {
		return err
	}

	response, err := c.UseCase.BreakOut(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to break out attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

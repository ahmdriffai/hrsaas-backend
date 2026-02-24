package http

import (
	"fmt"
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ShiftController struct {
	UseCase *usecase.ShiftUseCase
	Log     *logrus.Logger
}

func NewShifController(usecase *usecase.ShiftUseCase, log *logrus.Logger) *ShiftController {
	return &ShiftController{
		UseCase: usecase,
		Log:     log,
	}
}

func (c *ShiftController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateShiftRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.CompanyID = middleware.GetCompanyId(ctx)

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create shift")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ShiftResponse]{
		Data: response,
	})
}

func (c *ShiftController) AssignEmployee(ctx *fiber.Ctx) error {
	request := new(model.AssignEmployeeToShiftRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		fmt.Println("error y ages")
		return fiber.ErrBadRequest
	}

	request.CompanyID = middleware.GetCompanyId(ctx)

	if err := c.UseCase.AssignEmployee(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("failed to assign employee to shift")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

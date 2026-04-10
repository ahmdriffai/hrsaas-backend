package http

import (
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TimeOffTypeController struct {
	TypeUseCase *usecase.TimeOffTypeUseCase
	Log         *logrus.Logger
}

func NewTimeOffTypeController(typeUseCase *usecase.TimeOffTypeUseCase, log *logrus.Logger) *TimeOffTypeController {
	return &TimeOffTypeController{
		TypeUseCase: typeUseCase,
		Log:         log,
	}
}

// TODO: Support pagination if types grow large.
func (c *TimeOffTypeController) ListTypes(ctx *fiber.Ctx) error {
	responses, err := c.TypeUseCase.ListTypes(ctx.UserContext())
	if err != nil {
		c.Log.WithError(err).Error("failed to list time off types")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffTypeResponse]{
		Data: responses,
	})
}

// TODO: Enforce admin-only access with middleware at router.
func (c *TimeOffTypeController) CreateType(ctx *fiber.Ctx) error {
	request := new(model.CreateTimeOffTypeRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.TypeUseCase.CreateType(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create time off type")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TimeOffTypeResponse]{
		Data: response,
	})
}

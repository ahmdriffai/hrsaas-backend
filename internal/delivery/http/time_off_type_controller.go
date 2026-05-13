package http

import (
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

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
	request := &model.SearchTimeOffTypeRequest{
		Name: ctx.Query("name", ""),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}

	responses, total, err := c.TypeUseCase.ListTypes(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list time off types")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.TimeOffTypeResponse]{
		Data: responses,
		Paging: &model.PageMetadata{
			Page:      request.Page,
			Size:      request.Size,
			TotalItem: total,
			TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
		},
	})
}

// TODO: Enforce admin-only access with middleware at router.
func (c *TimeOffTypeController) CreateType(ctx *fiber.Ctx) error {
	request := new(model.CreateTimeOffTypeRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.TypeUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create time off type")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TimeOffTypeResponse]{
		Data: response,
	})
}

func (c *TimeOffTypeController) Detail(ctx *fiber.Ctx) error {
	response, err := c.TypeUseCase.Detail(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		c.Log.WithError(err).Error("failed to get time off type detail")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TimeOffTypeResponse]{
		Data: response,
	})
}

func (c *TimeOffTypeController) UpdateType(ctx *fiber.Ctx) error {
	request := new(model.UpdateTimeOffTypeRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.TypeUseCase.UpdateType(ctx.UserContext(), ctx.Params("id"), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to update time off type")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TimeOffTypeResponse]{
		Data: response,
	})
}

func (c *TimeOffTypeController) DeleteType(ctx *fiber.Ctx) error {
	if err := c.TypeUseCase.Delete(ctx.UserContext(), ctx.Params("id")); err != nil {
		c.Log.WithError(err).Error("failed to delete time off type")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

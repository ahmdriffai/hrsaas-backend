package http

import (
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type PermissionController struct {
	UseCase *usecase.PermissionUseCase
	Log     *logrus.Logger
}

func NewPermissionController(useCase *usecase.PermissionUseCase, log *logrus.Logger) *PermissionController {
	return &PermissionController{UseCase: useCase, Log: log}
}

func (c *PermissionController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreatePermissionRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PermissionResponse]{Data: response})
}

func (c *PermissionController) List(ctx *fiber.Ctx) error {
	request := &model.SearchPermissionRequest{
		Name: ctx.Query("name", ""),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.PermissionResponse]{
		Data: responses,
		Paging: &model.PageMetadata{
			Page:      request.Page,
			Size:      request.Size,
			TotalItem: total,
			TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
		},
	})
}

func (c *PermissionController) Detail(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	response, err := c.UseCase.Detail(ctx.UserContext(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PermissionResponse]{Data: response})
}

func (c *PermissionController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdatePermissionRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.ID = ctx.Params("id")

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PermissionResponse]{Data: response})
}

func (c *PermissionController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.UseCase.Delete(ctx.UserContext(), id); err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[any]{Data: nil})
}

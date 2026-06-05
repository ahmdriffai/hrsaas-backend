package http

import (
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EmployeeIdentityController struct {
	UseCase *usecase.EmployeeIdentityUseCase
	Log     *logrus.Logger
}

func NewEmployeeIdentityController(useCase *usecase.EmployeeIdentityUseCase, log *logrus.Logger) *EmployeeIdentityController {
	return &EmployeeIdentityController{UseCase: useCase, Log: log}
}

func (c *EmployeeIdentityController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateEmployeeIdentityRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeIdentityResponse]{Data: response})
}

func (c *EmployeeIdentityController) List(ctx *fiber.Ctx) error {
	request := &model.SearchEmployeeIdentityRequest{
		EmployeeID: ctx.Query("employee_id", ""),
		Page:       ctx.QueryInt("page", 1),
		Size:       ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.EmployeeIdentityResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *EmployeeIdentityController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateEmployeeIdentityRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Update(ctx.UserContext(), ctx.Params("identity_id"), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeIdentityResponse]{Data: response})
}

func (c *EmployeeIdentityController) Delete(ctx *fiber.Ctx) error {
	if err := c.UseCase.Delete(ctx.UserContext(), ctx.Params("identity_id")); err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

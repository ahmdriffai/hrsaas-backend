package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EmployeeTrainingController struct {
	UseCase *usecase.EmployeeTrainingUseCase
	Log     *logrus.Logger
}

func NewEmployeeTrainingController(useCase *usecase.EmployeeTrainingUseCase, log *logrus.Logger) *EmployeeTrainingController {
	return &EmployeeTrainingController{UseCase: useCase, Log: log}
}

func (c *EmployeeTrainingController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateEmployeeTrainingRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}
	request.CompanyID = middleware.GetCompanyId(ctx)

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeTrainingResponse]{Data: response})
}

func (c *EmployeeTrainingController) List(ctx *fiber.Ctx) error {
	request := &model.SearchEmployeeTrainingRequest{
		CompanyID:  middleware.GetCompanyId(ctx),
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

	return ctx.JSON(model.WebResponse[[]model.EmployeeTrainingResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *EmployeeTrainingController) Detail(ctx *fiber.Ctx) error {
	response, err := c.UseCase.Detail(ctx.UserContext(), ctx.Params("training_id"))
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeTrainingResponse]{Data: response})
}

func (c *EmployeeTrainingController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateEmployeeTrainingRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}
	request.ID = ctx.Params("training_id")

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeTrainingResponse]{Data: response})
}

func (c *EmployeeTrainingController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("training_id")
	if err := c.UseCase.Delete(ctx.UserContext(), id); err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

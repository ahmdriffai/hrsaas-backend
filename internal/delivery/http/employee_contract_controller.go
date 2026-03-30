package http

import (
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EmployeeContractController struct {
	UseCase *usecase.EmployeeContractUseCase
	Log     *logrus.Logger
}

func NewEmployeeContractController(useCase *usecase.EmployeeContractUseCase, log *logrus.Logger) *EmployeeContractController {
	return &EmployeeContractController{
		UseCase: useCase,
		Log:     log,
	}
}

// TODO: Enforce admin-only access.
func (c *EmployeeContractController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateEmployeeContractRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create employee contract")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeContractResponse]{
		Data: response,
	})
}

// TODO: List employee contract should be accessible by employee and admin.
func (c *EmployeeContractController) List(ctx *fiber.Ctx) error {
	request := new(model.SearchEmployeeContractRequest)
	request.EmployeeID = ctx.Query("employee_id", "")
	request.DivisionID = ctx.Query("division_id", "")
	request.PositionID = ctx.Query("position_id", "")
	request.ActiveOnly = ctx.QueryBool("active_only", false)
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list employee contracts")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.EmployeeContractResponse]{
		Data:   responses,
		Paging: paging,
	})
}

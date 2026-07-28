package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

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

func (c *ShiftController) List(ctx *fiber.Ctx) error {
	request := new(model.SearchShiftRequest)
	request.CompanyID = middleware.GetCompanyId(ctx)
	request.Key = ctx.Query("key", "")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list shifts")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.ShiftResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *ShiftController) ListCurrent(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	company := middleware.GetCompanyId(ctx)
	request := &model.SearchShiftRequest{
		CompanyID:  company,
		EmployeeID: user.Employee.ID,
		Page:       ctx.QueryInt("page", 1),
		Size:       ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.ShiftResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *ShiftController) Detail(ctx *fiber.Ctx) error {
	request := new(model.DetailShifRequest)
	request.CompanyID = middleware.GetCompanyId(ctx)
	request.ShiftID = ctx.Params("shiftID")

	response, err := c.UseCase.DetailShift(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list shifts")
		return err
	}

	return ctx.JSON(model.WebResponse[model.ShiftResponse]{
		Data: *response,
	})
}

func (c *ShiftController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateShiftRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	companyID := middleware.GetCompanyId(ctx)

	response, err := c.UseCase.Update(ctx.UserContext(), ctx.Params("shiftID"), companyID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to update shift")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ShiftResponse]{
		Data: response,
	})
}

func (c *ShiftController) BulkAssignEmployees(ctx *fiber.Ctx) error {
	request := new(model.BulkAssignEmployeesToShiftRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.CompanyID = middleware.GetCompanyId(ctx)
	request.ShiftID = ctx.Params("shiftID")

	response, err := c.UseCase.BulkAssignEmployees(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to bulk assign employees to shift")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ShiftResponse]{Data: response})
}

func (c *ShiftController) Delete(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)

	if err := c.UseCase.Delete(ctx.UserContext(), ctx.Params("shiftID"), companyID); err != nil {
		c.Log.WithError(err).Error("failed to delete shift")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

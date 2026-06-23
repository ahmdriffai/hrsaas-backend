package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OfficeLocationController struct {
	UseCase *usecase.OfficeLocationUseCase
	Log     *logrus.Logger
}

func NewOfficeLocationController(useCase *usecase.OfficeLocationUseCase, log *logrus.Logger) *OfficeLocationController {
	return &OfficeLocationController{
		UseCase: useCase,
		Log:     log,
	}
}

// Create Office Location Controller
func (c *OfficeLocationController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateOfficeLocationRequest)
	companyID := middleware.GetCompanyId(ctx)
	request.CompanyID = companyID

	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create office location")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.OfficeLocationResponse]{
		Data: response,
	})
}

/* Search Office Location Controller
 */
func (c *OfficeLocationController) List(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)
	request := new(model.SearchOfficeLocationRequest)
	request.Key = ctx.Query("key", "")
	request.CompanyID = companyID
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching office locations")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.OfficeLocationResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *OfficeLocationController) AssignEmployee(ctx *fiber.Ctx) error {
	request := new(model.AssignEmployeeToOfficeLocationRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.CompanyID = middleware.GetCompanyId(ctx)

	if err := c.UseCase.AssignEmployee(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("failed to assign employee to office location")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

func (c *OfficeLocationController) Detail(ctx *fiber.Ctx) error {
	request := new(model.DetailOfficeLocationRequest)
	request.CompanyID = middleware.GetCompanyId(ctx)
	request.OfficeLocationID = ctx.Params("officeLocationID")

	response, err := c.UseCase.DetailOfficeLocation(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list shifts")
		return err
	}

	return ctx.JSON(model.WebResponse[model.OfficeLocationResponse]{
		Data: *response,
	})
}

func (c *OfficeLocationController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateOfficeLocationRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	companyID := middleware.GetCompanyId(ctx)
	response, err := c.UseCase.Update(ctx.UserContext(), ctx.Params("officeLocationID"), companyID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to update office location")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.OfficeLocationResponse]{
		Data: response,
	})
}

func (c *OfficeLocationController) BulkAssignEmployees(ctx *fiber.Ctx) error {
	request := new(model.BulkAssignEmployeesToOfficeLocationRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.CompanyID = middleware.GetCompanyId(ctx)
	request.OfficeLocationID = ctx.Params("officeLocationID")

	response, err := c.UseCase.BulkAssignEmployees(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to bulk assign employees to office location")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.OfficeLocationResponse]{Data: response})
}

func (c *OfficeLocationController) Delete(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)

	if err := c.UseCase.Delete(ctx.UserContext(), ctx.Params("officeLocationID"), companyID); err != nil {
		c.Log.WithError(err).Error("failed to delete office location")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

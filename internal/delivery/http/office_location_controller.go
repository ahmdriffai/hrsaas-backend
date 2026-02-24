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

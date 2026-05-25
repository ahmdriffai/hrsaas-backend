package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type HolidayController struct {
	UseCase *usecase.HolidayUseCase
	Log     *logrus.Logger
}

func NewHolidayController(useCase *usecase.HolidayUseCase, log *logrus.Logger) *HolidayController {
	return &HolidayController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *HolidayController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateHolidayRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("Failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	request.CompanyID = middleware.GetCompanyId(ctx)

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.HolidayResponse]{
		Data: response,
	})
}

func (c *HolidayController) List(ctx *fiber.Ctx) error {
	request := new(model.SearchHolidayRequest)
	request.CompanyID = middleware.GetCompanyId(ctx)
	request.Key = ctx.Query("key", "")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

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
	return ctx.JSON(model.WebResponse[[]model.HolidayResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *HolidayController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateHolidayRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("Failed to parse body request")
		return fiber.ErrBadRequest
	}

	request.ID = ctx.Params("id")
	request.CompanyID = middleware.GetCompanyId(ctx)

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.HolidayResponse]{
		Data: response,
	})
}

func (c *HolidayController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	companyID := middleware.GetCompanyId(ctx)

	if err := c.UseCase.Delete(ctx.UserContext(), id, companyID); err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

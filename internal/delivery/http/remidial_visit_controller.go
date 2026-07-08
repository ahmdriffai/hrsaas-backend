package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RemidialVisitController struct {
	UseCase *usecase.RemidialVisitUseCase
	Log     *logrus.Logger
}

func NewRemidialVisitController(useCase *usecase.RemidialVisitUseCase, log *logrus.Logger) *RemidialVisitController {
	return &RemidialVisitController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *RemidialVisitController) SearchNasabah(ctx *fiber.Ctx) error {
	request := new(model.SearchNasabahRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	responses, err := c.UseCase.SearchNasabah(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to search nasabah")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.NasabahData]{
		Data: responses,
	})
}

func (c *RemidialVisitController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateRemidialVisitRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request.CompanyID = user.CompanyID
	request.EmployeeID = user.Employee.ID

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create remidial visit")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RemidialVisitResponse]{
		Data: response,
	})
}

func (c *RemidialVisitController) ListCurrent(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request := new(model.SearchRemidialVisitRequest)
	request.EmployeeID = user.Employee.ID
	request.NasabahName = ctx.Query("nama", "")
	request.StartDate = ctx.Query("start_date", "")
	request.EndDate = ctx.Query("end_date", "")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list current remidial visits")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.RemidialVisitResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *RemidialVisitController) ListAdmin(ctx *fiber.Ctx) error {
	request := new(model.SearchRemidialVisitRequest)
	request.EmployeeID = ctx.Query("employee_id", "")
	request.EmployeeName = ctx.Query("employee_name", "")
	request.NasabahName = ctx.Query("nama", "")
	request.StartDate = ctx.Query("start_date", "")
	request.EndDate = ctx.Query("end_date", "")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list admin remidial visits")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.RemidialVisitResponse]{
		Data:   responses,
		Paging: paging,
	})
}

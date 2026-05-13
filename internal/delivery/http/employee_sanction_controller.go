package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EmSancController struct {
	EmSancUseCase *usecase.EmSancUseCase
	Log           *logrus.Logger
}

func NewEmSancController(emSancUseCase *usecase.EmSancUseCase, log *logrus.Logger) *EmSancController {
	return &EmSancController{
		EmSancUseCase: emSancUseCase,
		Log:           log,
	}
}

func (c *EmSancController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateEmSancRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("Failed to parse request body")
		return fiber.ErrBadRequest
	}

	companyID := middleware.GetCompanyId(ctx)
	request.CompanyID = companyID
	response, err := c.EmSancUseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to create employee sanction")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmSancResponse]{
		Data: response,
	})
}

func (c *EmSancController) CurrentSearch(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)
	user := middleware.GetUser(ctx)
	sanctionID := ctx.Query("sanction_id")

	req := new(model.SearchEmSancRequest)

	req.EmployeeID = user.Employee.ID
	req.SanctionID = sanctionID
	req.CompanyID = companyID

	// start_date
	req.StartDate = ctx.Query("start_date", "")
	req.EndDate = ctx.Query("end_date", "")

	req.Reason = ctx.Query("reason", "")
	req.Status = ctx.Query("status", "")

	// pagination default
	req.Page = ctx.QueryInt("page", 1)
	req.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.EmSancUseCase.Search(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("error searching employee sanction")
		return err
	}

	paging := &model.PageMetadata{
		Page:      req.Page,
		Size:      req.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(req.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.EmSancResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *EmSancController) Search(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)
	req := new(model.SearchEmSancRequest)
	req.CompanyID = companyID

	// start_date
	req.StartDate = ctx.Query("start_date", "")
	req.EndDate = ctx.Query("end_date", "")

	// sanction id
	req.SanctionID = ctx.Query("sanction_id", "")
	req.Reason = ctx.Query("reason", "")
	req.Status = ctx.Query("status", "")

	// pagination default
	req.Page = ctx.QueryInt("page", 1)
	req.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.EmSancUseCase.Search(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("error searching contact")
		return err
	}

	paging := &model.PageMetadata{
		Page:      req.Page,
		Size:      req.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(req.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.EmSancResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *EmSancController) Detail(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)

	response, err := c.EmSancUseCase.Detail(ctx.UserContext(), ctx.Params("id"), companyID)
	if err != nil {
		c.Log.WithError(err).Error("failed to get employee sanction detail")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmSancResponse]{
		Data: response,
	})
}

func (c *EmSancController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateEmSancRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("Failed to parse request body")
		return fiber.ErrBadRequest
	}

	companyID := middleware.GetCompanyId(ctx)

	response, err := c.EmSancUseCase.Update(ctx.UserContext(), ctx.Params("id"), companyID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to update employee sanction")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmSancResponse]{
		Data: response,
	})
}

func (c *EmSancController) Delete(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)

	if err := c.EmSancUseCase.Delete(ctx.UserContext(), ctx.Params("id"), companyID); err != nil {
		c.Log.WithError(err).Error("failed to delete employee sanction")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/lib"
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
	req.UserID = user.ID
	req.SanctionID = sanctionID
	req.CompanyID = companyID

	// start_date
	if v := ctx.Query("start_date"); v != "" {
		d := &lib.DateOnly{}
		if err := d.UnmarshalJSON([]byte(v)); err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"error": "invalid end_date format (YYYY-MM-DD)",
			})
		}
		req.StartDate = d
	}

	// end_date
	if v := ctx.Query("end_date"); v != "" {
		d := &lib.DateOnly{}
		if err := d.UnmarshalJSON([]byte(v)); err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"error": "invalid end_date format (YYYY-MM-DD)",
			})
		}
		req.EndDate = d
	}
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

func (c *EmSancController) Search(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)
	req := new(model.SearchEmSancRequest)
	req.CompanyID = companyID

	// start_date
	if v := ctx.Query("start_date"); v != "" {
		d := &lib.DateOnly{}
		if err := d.UnmarshalJSON([]byte(v)); err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"error": "invalid end_date format (YYYY-MM-DD)",
			})
		}
		req.StartDate = d
	}

	// end_date
	if v := ctx.Query("end_date"); v != "" {
		d := &lib.DateOnly{}
		if err := d.UnmarshalJSON([]byte(v)); err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"error": "invalid end_date format (YYYY-MM-DD)",
			})
		}
		req.EndDate = d
	}
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

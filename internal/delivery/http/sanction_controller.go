package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type SanctionController struct {
	SanctionUseCase *usecase.SanctionUseCase
	Log             *logrus.Logger
}

func NewSanctionController(sanctionUseCase *usecase.SanctionUseCase, log *logrus.Logger) *SanctionController {
	return &SanctionController{
		SanctionUseCase: sanctionUseCase,
		Log:             log,
	}
}

/* Create Sanction Controller
 */
func (c *SanctionController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateSanctionRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	companyID := middleware.GetCompanyId(ctx)
	request.CompanyID = companyID

	response, err := c.SanctionUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create sanction")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SanctionResponse]{
		Data: response,
	})
}

/* Search Sanction
 */
func (c *SanctionController) ListSanction(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)
	request := new(model.SearchSanctionRequest)
	request.Key = ctx.Query("key", "")
	request.CompanyID = companyID
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.SanctionUseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching sanction")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.SanctionResponse]{
		Data:   responses,
		Paging: paging,
	})
}

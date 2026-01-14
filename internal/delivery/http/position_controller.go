package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type PositionController struct {
	PositionUseCase *usecase.PositionUseCase
	Log             *logrus.Logger
}

func NewPositionController(positionUseCase *usecase.PositionUseCase, log *logrus.Logger) *PositionController {
	return &PositionController{
		Log:             log,
		PositionUseCase: positionUseCase,
	}
}

/* List Position Controller
 */
func (c *PositionController) ListPosition(ctx *fiber.Ctx) error {
	request := new(model.SeachPositionRequest)
	companyID := middleware.GetCompanyId(ctx)
	request.CompanyID = companyID
	request.Name = ctx.Query("name", "")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.PositionUseCase.Search(ctx.UserContext(), request)
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

	return ctx.JSON(model.WebResponse[[]model.PositionResponse]{
		Data:   responses,
		Paging: paging,
	})
}

/* Create Position
 */
func (c *PositionController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreatePositionRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	companyID := middleware.GetCompanyId(ctx)
	request.CompanyID = companyID

	response, err := c.PositionUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create sanction")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PositionResponse]{
		Data: response,
	})
}

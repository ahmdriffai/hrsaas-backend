package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CompanyController struct {
	UseCase *usecase.CompanyUseCase
	Log     *logrus.Logger
}

func NewCompanyController(useCase *usecase.CompanyUseCase, log *logrus.Logger) *CompanyController {
	return &CompanyController{
		UseCase: useCase,
		Log:     log,
	}
}

/*
Create Controller
*/
func (c *CompanyController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateCompanyRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create address")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CompanyResponse]{
		Data: response,
	})
}

/*
Register Company Controller
*/
func (c *CompanyController) Register(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	request := new(model.RegisterCompanyRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.UserID = user.ID

	response, err := c.UseCase.Register(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to Register company")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CompanyResponse]{
		Data: response,
	})
}

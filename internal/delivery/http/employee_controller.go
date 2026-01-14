package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EmployeeController struct {
	EmployeeUseCase *usecase.EmployeeUseCase
	Log             *logrus.Logger
}

func NewEmployeeController(employeeUseCase *usecase.EmployeeUseCase, log *logrus.Logger) *EmployeeController {
	return &EmployeeController{
		EmployeeUseCase: employeeUseCase,
		Log:             log,
	}
}

/* Create Employee Controller
 */
func (c *EmployeeController) CreateEmployee(ctx *fiber.Ctx) error {

	request := new(model.CreateEmployeeRequest)
	if err := ctx.BodyParser(&request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	request.CompanyID = user.CompanyID

	response, err := c.EmployeeUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to create employee")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeResponse]{
		Data: response,
	})
}

/* Search Employee Controller
 */
func (c *EmployeeController) ListEmployee(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)
	request := new(model.SearchEmployeeRequest)
	request.Key = ctx.Query("key", "")
	request.CompanyID = companyID
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.EmployeeUseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching contact")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.EmployeeResponse]{
		Data:   responses,
		Paging: paging,
	})
}

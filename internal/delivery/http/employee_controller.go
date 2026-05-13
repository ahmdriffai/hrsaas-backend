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

// Create Employee Controller
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

// Search Employee Controller
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

// Import Excel Employee Controller
func (c *EmployeeController) ImportExcel(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)
	file, err := ctx.FormFile("file")
	if err != nil {
		c.Log.WithError(err).Error("Failed to get file from form data")
		return fiber.ErrBadRequest
	}

	request := &model.ImportExcelEmployeeRequest{
		CompanyID: companyID,
		File:      file,
	}

	totalData, err := c.EmployeeUseCase.ImportExcel(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to import excel employee")
		return err
	}

	return ctx.JSON(model.WebResponse[int64]{
		Data: totalData,
	})
}

// Detail Employee Controller
func (c *EmployeeController) DetailEmployee(ctx *fiber.Ctx) error {
	employeeID := ctx.Params("id")
	companyID := middleware.GetCompanyId(ctx)

	request := &model.DetailEmployeeRequest{
		ID:        employeeID,
		CompanyID: companyID,
	}

	response, err := c.EmployeeUseCase.Detail(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to get employee detail")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeResponse]{
		Data: response,
	})
}

func (c *EmployeeController) UpdateEmployee(ctx *fiber.Ctx) error {
	request := new(model.UpdateEmployeeRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.ID = ctx.Params("id")
	companyID := middleware.GetCompanyId(ctx)

	response, err := c.EmployeeUseCase.Update(ctx.UserContext(), companyID, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to update employee")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeResponse]{
		Data: response,
	})
}

func (c *EmployeeController) DeleteEmployee(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)

	if err := c.EmployeeUseCase.Delete(ctx.UserContext(), ctx.Params("id"), companyID); err != nil {
		c.Log.WithError(err).Error("Failed to delete employee")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

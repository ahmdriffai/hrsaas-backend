package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AttendanceController struct {
	UseCase *usecase.AttendanceUseCase
	Log     *logrus.Logger
}

func NewAttendanceController(useCase *usecase.AttendanceUseCase, log *logrus.Logger) *AttendanceController {
	return &AttendanceController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *AttendanceController) parseRequest(ctx *fiber.Ctx) (*model.CheckInAttendanceRequest, error) {
	user := middleware.GetUser(ctx)
	request := new(model.CheckInAttendanceRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return nil, fiber.ErrBadRequest
	}
	request.CompanyID = user.CompanyID
	request.EmployeeID = user.Employee.ID
	return request, nil
}

func (c *AttendanceController) List(ctx *fiber.Ctx) error {
	request := new(model.SearchAttendanceRequest)
	request.CompanyID = middleware.GetCompanyId(ctx)
	request.EmployeeID = ctx.Query("employee_id", "")
	request.Date = ctx.Query("date", "")
	request.Status = ctx.Query("status", "")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list attendances")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.AttendanceResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *AttendanceController) ListCurrent(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	request := &model.SearchAttendanceRequest{
		CompanyID:  user.CompanyID,
		EmployeeID: user.Employee.ID,
		Page:       ctx.QueryInt("page", 1),
		Size:       ctx.QueryInt("size", 10),
	}

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

	return ctx.JSON(model.WebResponse[[]model.AttendanceResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *AttendanceController) ListLog(ctx *fiber.Ctx) error {
	request := new(model.SearchAttendanceLogRequest)
	request.CompanyID = middleware.GetCompanyId(ctx)
	request.AttendanceID = ctx.Query("attendance_id", "")
	request.EmployeeID = ctx.Query("employee_id", "")
	request.Type = ctx.Query("type", "")
	request.Date = ctx.Query("date", "")
	if isApproved := ctx.Query("is_approved", ""); isApproved != "" {
		val := isApproved == "true"
		request.IsApproved = &val
	}
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.SearchLog(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list attendance logs")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.AttendanceLogResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *AttendanceController) Detail(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)

	response, err := c.UseCase.Detail(ctx.UserContext(), ctx.Params("attendanceID"), companyID)
	if err != nil {
		c.Log.WithError(err).Error("failed to get attendance detail")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

func (c *AttendanceController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateAttendanceRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	companyID := middleware.GetCompanyId(ctx)
	response, err := c.UseCase.Update(ctx.UserContext(), ctx.Params("attendanceID"), companyID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to update attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

func (c *AttendanceController) Delete(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)

	if err := c.UseCase.Delete(ctx.UserContext(), ctx.Params("attendanceID"), companyID); err != nil {
		c.Log.WithError(err).Error("failed to delete attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{Data: nil})
}

func (c *AttendanceController) CheckIn(ctx *fiber.Ctx) error {
	request, err := c.parseRequest(ctx)
	if err != nil {
		return err
	}

	response, err := c.UseCase.CheckIn(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to check in attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

func (c *AttendanceController) CheckOut(ctx *fiber.Ctx) error {
	request, err := c.parseRequest(ctx)
	if err != nil {
		return err
	}

	response, err := c.UseCase.CheckOut(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to check out attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

func (c *AttendanceController) BreakIn(ctx *fiber.Ctx) error {
	request, err := c.parseRequest(ctx)
	if err != nil {
		return err
	}

	response, err := c.UseCase.BreakIn(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to break in attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

func (c *AttendanceController) BreakOut(ctx *fiber.Ctx) error {
	request, err := c.parseRequest(ctx)
	if err != nil {
		return err
	}

	response, err := c.UseCase.BreakOut(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to break out attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

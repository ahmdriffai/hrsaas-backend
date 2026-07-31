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

	file, err := ctx.FormFile("file")
	if err != nil {
		c.Log.WithError(err).Error("failed to read face image file")
		return nil, fiber.NewError(fiber.StatusBadRequest, "File gambar wajah wajib diisi")
	}
	request.File = file

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

func (c *AttendanceController) DetailToday(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	response, err := c.UseCase.DetailToday(ctx.UserContext(), user.Employee.ID, user.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("failed to get today's attendace detail")
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

func (c *AttendanceController) RegisterFaceCurrent(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	request := new(model.RegisterFaceRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	request.EmployeeID = user.Employee.ID
	request.CompanyID = user.CompanyID

	response, err := c.UseCase.RegisterFace(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to register face")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RegisterFaceResponse]{Data: response})
}

func (c *AttendanceController) FaceStatusCurrent(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	response, err := c.UseCase.FaceStatus(ctx.UserContext(), user.Employee.ID, user.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to get face status")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.FaceStatusResponse]{Data: response})
}

func (c *AttendanceController) RegisterFace(ctx *fiber.Ctx) error {
	request := &model.RegisterFaceRequest{
		EmployeeID: ctx.Params("employee_id"),
		CompanyID:  middleware.GetCompanyId(ctx),
		ObjectKey:  ctx.Params("object_key"),
	}

	response, err := c.UseCase.RegisterFace(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to register face")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RegisterFaceResponse]{Data: response})
}

func (c *AttendanceController) FaceStatus(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)

	response, err := c.UseCase.FaceStatus(ctx.UserContext(), ctx.Params("employee_id"), companyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to get face status")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.FaceStatusResponse]{Data: response})
}

func (c *AttendanceController) DeleteFace(ctx *fiber.Ctx) error {
	companyID := middleware.GetCompanyId(ctx)

	if err := c.UseCase.DeleteFace(ctx.UserContext(), ctx.Params("employee_id"), companyID); err != nil {
		c.Log.WithError(err).Error("Failed to delete face")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{Data: nil})
}

func (c *AttendanceController) parseLendRequest(ctx *fiber.Ctx) (*model.CheckInAttendanceRequest, error) {
	request := new(model.CheckInAttendanceRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return nil, fiber.ErrBadRequest
	}

	if request.EmployeeID == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "employee_id wajib diisi")
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		c.Log.WithError(err).Error("failed to read face image file")
		return nil, fiber.NewError(fiber.StatusBadRequest, "File gambar wajah wajib diisi")
	}
	request.File = file

	request.CompanyID = middleware.GetCompanyId(ctx)
	return request, nil
}

func (c *AttendanceController) LendCheckIn(ctx *fiber.Ctx) error {
	request, err := c.parseLendRequest(ctx)
	if err != nil {
		return err
	}

	response, err := c.UseCase.CheckIn(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to lend check in attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

func (c *AttendanceController) LendCheckOut(ctx *fiber.Ctx) error {
	request, err := c.parseLendRequest(ctx)
	if err != nil {
		return err
	}

	response, err := c.UseCase.CheckOut(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to lend check out attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{Data: response})
}

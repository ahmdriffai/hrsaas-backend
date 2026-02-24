package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"

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

/*
Check In Controller
*/
func (c *AttendanceController) CheckIn(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	request := new(model.CheckInAttendanceRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}
	request.CompanyID = user.CompanyID
	request.EmployeeID = user.Employee.ID

	response, err := c.UseCase.CheckIn(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to check in attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{
		Data: response,
	})
}

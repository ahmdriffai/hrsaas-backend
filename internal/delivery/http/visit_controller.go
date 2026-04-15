package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type VisitController struct {
	UseCase *usecase.VisitUseCase
	Log     *logrus.Logger
}

func NewVisitController(useCase *usecase.VisitUseCase, log *logrus.Logger) *VisitController {
	return &VisitController{
		UseCase: useCase,
		Log:     log,
	}
}

// TODO: Enforce employee-only access (current user).
func (c *VisitController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateVisitRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	response, err := c.UseCase.Create(ctx.UserContext(), user.Employee.ID, user.CompanyID, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create visit")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.VisitResponse]{
		Data: response,
	})
}

// TODO: Enforce admin-only access.
func (c *VisitController) List(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	if !strings.EqualFold(user.Role, "ADMIN") {
		return fiber.NewError(fiber.StatusForbidden, "Forbidden")
	}

	request := new(model.SearchVisitRequest)
	request.EmployeeID = ctx.Query("employee_id", "")
	request.VisitType = ctx.Query("visit_type", "")
	request.StartDate = ctx.Query("start_date", "")
	request.EndDate = ctx.Query("end_date", "")
	request.SortBy = ctx.Query("sort_by", "newest")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list visits")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.VisitResponse]{
		Data:   responses,
		Paging: paging,
	})
}

// TODO: Return 403 if user has no employee profile.
func (c *VisitController) ListCurrent(ctx *fiber.Ctx) error {
	request := new(model.SearchVisitRequest)
	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	request.EmployeeID = user.Employee.ID
	request.VisitType = ctx.Query("visit_type", "")
	request.StartDate = ctx.Query("start_date", "")
	request.EndDate = ctx.Query("end_date", "")
	request.SortBy = ctx.Query("sort_by", "newest")
	request.Page = ctx.QueryInt("page", 1)
	request.Size = ctx.QueryInt("size", 10)

	responses, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list visits")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.VisitResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *VisitController) GetByID(ctx *fiber.Ctx) error {
	requestID := ctx.Params("id")
	if requestID == "" {
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if !strings.EqualFold(user.Role, "ADMIN") {
		if user.Employee == nil || user.Employee.ID == "" {
			return fiber.NewError(fiber.StatusForbidden, "Forbidden")
		}

		ownerID, err := c.UseCase.GetVisitOwner(ctx.UserContext(), requestID)
		if err != nil {
			c.Log.WithError(err).Error("failed to get visit owner")
			return err
		}
		if ownerID != user.Employee.ID {
			return fiber.NewError(fiber.StatusForbidden, "Forbidden")
		}
	}

	response, err := c.UseCase.GetByID(ctx.UserContext(), requestID)
	if err != nil {
		c.Log.WithError(err).Error("failed to get visit detail")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.VisitResponse]{
		Data: response,
	})
}

func (c *VisitController) CanDoVisit(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	visitType := ctx.Query("visit_type", "")
	if visitType == "" {
		return fiber.NewError(fiber.StatusBadRequest, "visit_type is required")
	}

	response, err := c.UseCase.CanDoVisit(ctx.UserContext(), user.Employee.ID, visitType)
	if err != nil {
		c.Log.WithError(err).Error("failed to check visit permission")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CanDoVisitResponse]{
		Data: response,
	})
}

func (c *VisitController) GetUnclosedVisit(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	response, err := c.UseCase.GetUnclosedVisit(ctx.UserContext(), user.Employee.ID)
	if err != nil {
		c.Log.WithError(err).Error("failed to get unclosed visit")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.VisitResponse]{
		Data: response,
	})
}

// TODO: Decide whether delete should be admin-only or owner-only.
func (c *VisitController) Delete(ctx *fiber.Ctx) error {
	requestID := ctx.Params("id")
	if requestID == "" {
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(ctx)
	if !strings.EqualFold(user.Role, "ADMIN") {
		if user.Employee == nil || user.Employee.ID == "" {
			return fiber.NewError(fiber.StatusForbidden, "Forbidden")
		}

		ownerID, err := c.UseCase.GetVisitOwner(ctx.UserContext(), requestID)
		if err != nil {
			c.Log.WithError(err).Error("failed to get visit owner")
			return err
		}
		if ownerID != user.Employee.ID {
			return fiber.NewError(fiber.StatusForbidden, "Forbidden")
		}
	}

	if err := c.UseCase.Delete(ctx.UserContext(), requestID); err != nil {
		c.Log.WithError(err).Error("failed to delete visit")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{Data: nil})
}

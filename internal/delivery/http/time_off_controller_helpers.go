package http

import (
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ensureOwnerOrAdmin(ctx *fiber.Ctx, requestUseCase *usecase.TimeOffRequestUseCase, requestID string) error {
	user := middleware.GetUser(ctx)
	if strings.EqualFold(user.Role, "ADMIN") {
		return nil
	}
	if user.Employee == nil || user.Employee.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee not found")
	}

	ownerID, err := requestUseCase.GetRequestOwner(ctx.UserContext(), requestID)
	if err != nil {
		return err
	}
	if ownerID != user.Employee.ID {
		return fiber.NewError(fiber.StatusForbidden, "Forbidden")
	}
	return nil
}

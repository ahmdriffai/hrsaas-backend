package middleware

import (
	"hr-sas/internal/model"

	"github.com/gofiber/fiber/v2"
)

func NewEmployee() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*model.UserResponse)
		if !ok || user == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		if user.Employee == nil {
			return fiber.NewError(fiber.StatusForbidden, "Forbidden: Employee access required")
		}

		return ctx.Next()
	}
}

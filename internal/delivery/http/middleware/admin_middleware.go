package middleware

import (
	"strings"

	"hr-sas/internal/model"

	"github.com/gofiber/fiber/v2"
)

func NewAdmin() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*model.UserResponse)
		if !ok || user == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		if !strings.EqualFold(user.Role, "ADMIN") {
			return fiber.NewError(fiber.StatusForbidden, "Forbidden: Admin access required")
		}

		return ctx.Next()
	}
}

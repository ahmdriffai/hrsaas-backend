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

		for _, role := range user.Roles {
			if strings.EqualFold(role.Name, "ADMIN") {
				return ctx.Next()
			}
		}

		return fiber.NewError(fiber.StatusForbidden, "Forbidden: Admin access required")
	}
}

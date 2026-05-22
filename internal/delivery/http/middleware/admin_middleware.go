package middleware

import (
	"strings"

	"hr-sas/internal/model"

	"github.com/gofiber/fiber/v2"
)

func NewAdmin(permission string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*model.UserResponse)
		if !ok || user == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		for _, p := range user.Permissions {
			if strings.EqualFold(p.Name, permission) {
				return ctx.Next()
			}
		}

		return fiber.NewError(fiber.StatusForbidden, "Forbidden: insufficient permissions")
	}
}

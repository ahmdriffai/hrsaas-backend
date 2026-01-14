package middleware

import (
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUseCase *usecase.UserUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Ambil header Authorization
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Missing Authorization header")
		}

		// Harus format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token format")
		}

		token := parts[1]
		request := &model.VerifyUserRequest{Token: token}
		userUseCase.Log.Debugf("Authorization : %s", request.Token)

		user, err := userUseCase.Verify(ctx.UserContext(), request)
		if err != nil {
			userUseCase.Log.Warnf("Failed find user by token : %+v", err)
			return err
		}

		// Simpan user ke context agar bisa diakses di handler
		ctx.Locals("user", user)

		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.UserResponse {
	return ctx.Locals("user").(*model.UserResponse)
}

func GetCompanyId(ctx *fiber.Ctx) string {
	user := GetUser(ctx)
	return user.CompanyID
}

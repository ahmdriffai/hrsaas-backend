package http

import (
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UserController struct {
	UserUseCase *usecase.UserUseCase
	Log         *logrus.Logger
	Viper       *viper.Viper
}

func NewUserController(userUseCase *usecase.UserUseCase, log *logrus.Logger, viper *viper.Viper) *UserController {
	return &UserController{
		UserUseCase: userUseCase,
		Log:         log,
		Viper:       viper,
	}
}

/*
Cookies
*/
func (c *UserController) setTokenCookie(ctx *fiber.Ctx, token string) {
	secure := c.Viper.GetBool("app.cookie_secure")
	domain := c.Viper.GetString("app.cookie_domain")
	sameSite := "Lax"
	if secure {
		sameSite = "None"
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Domain:   domain,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
	})
}

func (c *UserController) clearTokenCookie(ctx *fiber.Ctx) {
	secure := c.Viper.GetBool("app.cookie_secure")
	domain := c.Viper.GetString("app.cookie_domain")
	sameSite := "Lax"
	if secure {
		sameSite = "None"
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Domain:   domain,
		Path:     "/",
		MaxAge:   -1,
	})
}

/*
Register User Controller
*/
func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	response, err := c.UserUseCase.Register(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to register user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{
		Data: response,
	})
}

/*
Login User Controller
*/
func (c *UserController) Login(ctx *fiber.Ctx) error {
	userAgent := ctx.Get(fiber.HeaderUserAgent)
	ip := ctx.IP()

	c.Log.Infof("Login attempt from IP: %s, User-Agent: %s", ip, userAgent)

	request := new(model.LoginUserRequest)
	request.UserAgent = userAgent
	request.Ip = ip

	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	response, err := c.UserUseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to login user")
		return err
	}

	c.setTokenCookie(ctx, response.Token)

	return ctx.JSON(model.WebResponse[*model.LoginUserResponse]{
		Data: response,
	})
}

/* Get Current User Controller
 */
func (c *UserController) GetCurrentUser(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*model.UserResponse)
	return ctx.JSON(model.WebResponse[*model.UserResponse]{
		Data: user,
	})
}

func (c *UserController) ResetPassword(ctx *fiber.Ctx) error {
	request := new(model.ResetPasswordRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := c.UserUseCase.ResetPassword(ctx.UserContext(), ctx.Params("id"), request); err != nil {
		c.Log.WithError(err).Error("failed to reset password")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

func (c *UserController) ChangePassword(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*model.UserResponse)

	request := new(model.ChangePasswordRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := c.UserUseCase.ChangePassword(ctx.UserContext(), user.ID, request); err != nil {
		c.Log.WithError(err).Error("failed to change password")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

func (c *UserController) List(ctx *fiber.Ctx) error {
	request := &model.SearchUserRequest{
		Key:  ctx.Query("key", ""),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UserUseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list users")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.UserResponse]{
		Data: responses,
		Paging: &model.PageMetadata{
			Page:      request.Page,
			Size:      request.Size,
			TotalItem: total,
			TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
		},
	})
}

func (c *UserController) Detail(ctx *fiber.Ctx) error {
	response, err := c.UserUseCase.Detail(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		c.Log.WithError(err).Error("failed to get user detail")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{
		Data: response,
	})
}

func (c *UserController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.ID = ctx.Params("id")

	response, err := c.UserUseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to update user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{
		Data: response,
	})
}

func (c *UserController) Delete(ctx *fiber.Ctx) error {
	if err := c.UserUseCase.Delete(ctx.UserContext(), ctx.Params("id")); err != nil {
		c.Log.WithError(err).Error("failed to delete user")
		return err
	}

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

/*
Logout User Controller
*/
func (c *UserController) Logout(ctx *fiber.Ctx) error {
	token := ctx.Cookies("token")
	if token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "No token found in cookies")
	}

	if err := c.UserUseCase.Logout(ctx.UserContext(), token); err != nil {
		c.Log.WithError(err).Error("Gagal Logout")
		return err
	}

	c.clearTokenCookie(ctx)
	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

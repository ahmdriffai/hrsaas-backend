package http

import (
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"time"

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

	isProduction := viper.GetString("app.env") == "production"

	sameSite := "Lax"
	if isProduction {
		sameSite = "None"
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    response.Token,
		HTTPOnly: true,
		Secure:   c.Viper.GetBool("app.cookie_secure"), // true kalau production (HTTPS)
		SameSite: sameSite,
		Domain:   c.Viper.GetString("app.cookie_domain"), // ".bankwonosobo.co.id"
		Path:     "/",
	})

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

/*
Logout User Controller
*/
func (c *UserController) Logout(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*model.UserResponse)

	err := c.UserUseCase.Logout(ctx.UserContext(), user.ID)
	if err != nil {
		c.Log.WithError(err).Error("failed to logout user")
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})

	return ctx.JSON(model.WebResponse[any]{
		Data: nil,
	})
}

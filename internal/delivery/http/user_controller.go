package http

import (
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UserUseCase *usecase.UserUseCase
	Log         *logrus.Logger
}

func NewUserController(userUseCase *usecase.UserUseCase, log *logrus.Logger) *UserController {
	return &UserController{
		UserUseCase: userUseCase,
		Log:         log,
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

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    response.Token,
		HTTPOnly: true,
		Secure:   false, // true kalau production (HTTPS)
		SameSite: "None",
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
Assign Roles to User Controller
*/
func (c *UserController) AssignRoles(ctx *fiber.Ctx) error {
	request := new(model.AssignRoleRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.UserID = ctx.Params("id")

	response, err := c.UserUseCase.AssignRoles(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

/*
Remove Roles from User Controller
*/
func (c *UserController) RemoveRoles(ctx *fiber.Ctx) error {
	request := new(model.RemoveRoleRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.UserID = ctx.Params("id")

	response, err := c.UserUseCase.RemoveRoles(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
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

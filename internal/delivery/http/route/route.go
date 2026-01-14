package route

import (
	"hr-sas/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	AuthMiddleware     fiber.Handler
	CompanyController  *http.CompanyController
	UserController     *http.UserController
	EmployeeController *http.EmployeeController
	SanctionController *http.SanctionController
	EmSancController   *http.EmSancController
	PositionController *http.PositionController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRouter()
	c.SetupCompanyRouter()
	c.SetupUserRouter()
	c.SetupEmployeeRouter()
	c.SetupSanctionRouter()
	c.SetupEmployeeSanctionRouter()
	c.SetupPositionRouter()
}

/*
Guest Router
*/
func (c *RouteConfig) SetupGuestRouter() {
	c.App.Post("/api/_login", c.UserController.Login)
	c.App.Post("/api/_register", c.UserController.Register)
}

func (c *RouteConfig) SetupCompanyRouter() {
	route := c.App.Group("/api/companies", c.AuthMiddleware)
	route.Post("/", c.CompanyController.Create)
	route.Post("/_register", c.CompanyController.Register)
}

func (c *RouteConfig) SetupUserRouter() {
	route := c.App.Group("/api/users", c.AuthMiddleware)
	route.Get("/_current", c.UserController.GetCurrentUser)
	route.Delete("/_logout", c.UserController.Logout)
}

func (c *RouteConfig) SetupEmployeeRouter() {
	route := c.App.Group("/api/employees", c.AuthMiddleware)
	route.Post("/", c.EmployeeController.CreateEmployee)
	route.Get("/", c.EmployeeController.ListEmployee)
}

func (c *RouteConfig) SetupSanctionRouter() {
	route := c.App.Group("/api/sanctions", c.AuthMiddleware)
	route.Post("/", c.SanctionController.Create)
	route.Get("/", c.SanctionController.ListSanction)
}

func (c *RouteConfig) SetupEmployeeSanctionRouter() {
	route := c.App.Group("/api/employee-sanctions", c.AuthMiddleware)
	route.Post("/", c.EmSancController.Create)
	route.Get("/", c.EmSancController.Search)
	route.Get("/_current", c.EmSancController.CurrentSearch)
}

func (c *RouteConfig) SetupPositionRouter() {
	route := c.App.Group("/api/positions", c.AuthMiddleware)
	route.Get("/", c.PositionController.ListPosition)
	route.Post("/", c.PositionController.Create)
}

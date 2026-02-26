package route

import (
	"hr-sas/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                      *fiber.App
	AuthMiddleware           fiber.Handler
	CompanyController        *http.CompanyController
	UserController           *http.UserController
	EmployeeController       *http.EmployeeController
	OfficeLocationController *http.OfficeLocationController
	SanctionController       *http.SanctionController
	EmSancController         *http.EmSancController
	PositionController       *http.PositionController
	AttendanceController     *http.AttendanceController
	ShiftController          *http.ShiftController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRouter()
	c.SetupCompanyRouter()
	c.SetupUserRouter()
	c.SetupEmployeeRouter()
	c.SetupSanctionRouter()
	c.SetupEmployeeSanctionRouter()
	c.SetupPositionRouter()
	c.SetupOfficeLocationRouter()
	c.SetupAttendanceRouter()
	c.SetupShiftRouter()
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
	c.App.Delete("/api/users/_logout", c.UserController.Logout)
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

func (c *RouteConfig) SetupOfficeLocationRouter() {
	route := c.App.Group("/api/office-locations", c.AuthMiddleware)
	route.Get("/", c.OfficeLocationController.List)
	route.Post("/", c.OfficeLocationController.Create)
	route.Post("/assign-employee", c.OfficeLocationController.AssignEmployee)
}

func (c *RouteConfig) SetupAttendanceRouter() {

	route := c.App.Group("/api/attendances", c.AuthMiddleware)
	route.Post("/check-in", c.AttendanceController.CheckIn)
}

func (c *RouteConfig) SetupShiftRouter() {
	route := c.App.Group("/api/shifts", c.AuthMiddleware)
	route.Get("/", c.ShiftController.List)
	route.Post("/", c.ShiftController.Create)
	route.Post("/assign-employee", c.ShiftController.AssignEmployee)
}

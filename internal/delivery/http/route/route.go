package route

import (
	"hr-sas/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App *fiber.App

	AuthMiddleware     fiber.Handler
	AdminMiddleware    fiber.Handler
	EmployeeMiddleware fiber.Handler

	CompanyController        *http.CompanyController
	UserController           *http.UserController
	EmployeeController       *http.EmployeeController
	OfficeLocationController *http.OfficeLocationController
	SanctionController       *http.SanctionController
	EmSancController         *http.EmSancController
	PositionController       *http.PositionController
	AttendanceController     *http.AttendanceController
	ShiftController          *http.ShiftController
	TimeOffController        *http.TimeOffController
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
	c.SetupTimeOffRouter()
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
	route.Get("/", c.EmployeeController.ListEmployee)

	adminRoute := route.Group("/", c.AdminMiddleware)
	adminRoute.Post("/", c.EmployeeController.CreateEmployee)
}

func (c *RouteConfig) SetupSanctionRouter() {
	route := c.App.Group("/api/sanctions", c.AuthMiddleware, c.AdminMiddleware)
	route.Post("/", c.SanctionController.Create)
	route.Get("/", c.SanctionController.ListSanction)
}

func (c *RouteConfig) SetupEmployeeSanctionRouter() {
	route := c.App.Group("/api/employee-sanctions", c.AuthMiddleware)
	route.Get("/_current", c.EmployeeMiddleware, c.EmSancController.CurrentSearch)

	adminRouter := route.Group("/", c.AuthMiddleware)
	adminRouter.Post("/", c.EmSancController.Create)
	adminRouter.Get("/", c.EmSancController.Search)
}

func (c *RouteConfig) SetupPositionRouter() {
	route := c.App.Group("/api/positions", c.AuthMiddleware, c.AdminMiddleware)
	route.Get("/", c.PositionController.ListPosition)
	route.Post("/", c.PositionController.Create)
}

func (c *RouteConfig) SetupOfficeLocationRouter() {
	route := c.App.Group("/api/office-locations", c.AuthMiddleware, c.AdminMiddleware)
	route.Get("/", c.OfficeLocationController.List)
	route.Post("/", c.OfficeLocationController.Create)
	route.Post("/assign-employee", c.OfficeLocationController.AssignEmployee)
}

func (c *RouteConfig) SetupAttendanceRouter() {
	route := c.App.Group("/api/attendances", c.AuthMiddleware)
	route.Post("/check-in", c.EmployeeMiddleware, c.AttendanceController.CheckIn)
}

func (c *RouteConfig) SetupShiftRouter() {
	route := c.App.Group("/api/shifts", c.AuthMiddleware, c.AdminMiddleware)
	route.Get("/", c.ShiftController.List)
	route.Post("/", c.ShiftController.Create)
	route.Post("/assign-employee", c.ShiftController.AssignEmployee)
}

func (c *RouteConfig) SetupTimeOffRouter() {
	route := c.App.Group("/api/time-off-requests", c.AuthMiddleware)
	route.Get("/", c.TimeOffController.ListRequests)
	route.Post("/", c.EmployeeMiddleware, c.TimeOffController.CreateRequest)
	route.Get("/_current", c.EmployeeMiddleware, c.TimeOffController.ListCurrentRequests)
	route.Get("/:id/approvals", c.TimeOffController.ListApprovals)
	route.Patch("/:id/approvals/:approval_id/approve", c.EmployeeMiddleware, c.TimeOffController.Approve)
	route.Patch("/:id/approvals/:approval_id/reject", c.EmployeeMiddleware, c.TimeOffController.Reject)

	typeRoute := c.App.Group("/api/time-off-types", c.AuthMiddleware)
	typeRoute.Get("/", c.TimeOffController.ListTypes)

	balanceRoute := c.App.Group("/api/time-off-balances", c.AuthMiddleware)
	balanceRoute.Get("/_current", c.EmployeeMiddleware, c.TimeOffController.ListCurrentBalances)
}

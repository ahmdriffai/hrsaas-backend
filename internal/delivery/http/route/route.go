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

	CompanyController          *http.CompanyController
	UserController             *http.UserController
	EmployeeController         *http.EmployeeController
	OfficeLocationController   *http.OfficeLocationController
	SanctionController         *http.SanctionController
	EmSancController           *http.EmSancController
	PositionController         *http.PositionController
	AttendanceController       *http.AttendanceController
	ShiftController            *http.ShiftController
	TimeOffRequestController   *http.TimeOffRequestController
	TimeOffTypeController      *http.TimeOffTypeController
	TimeOffBalanceController   *http.TimeOffBalanceController
	TimeOffApprovalController  *http.TimeOffApprovalController
	UploadController           *http.UploadController
	EmployeeContractController *http.EmployeeContractController
	DivisionController         *http.DivisionController
	VisitController            *http.VisitController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRouter()
	c.SetupCompanyRouter()
	c.SetupUserRouter()
	c.SetupEmployeeRouter()
	c.SetupSanctionRouter()
	c.SetupEmployeeSanctionRouter()
	c.SetupPositionRouter()
	c.SetupDivisionRouter()
	c.SetupOfficeLocationRouter()
	c.SetupAttendanceRouter()
	c.SetupShiftRouter()
	c.SetupTimeOffRouter()
	c.SetupCommonRouter()
	c.SetupEmployeeContractRouter()
	c.SetupTimeOffApprovalRouter()
	c.SetupVisitRouter()
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

func (c *RouteConfig) SetupEmployeeContractRouter() {
	route := c.App.Group("/api/employee-contracts", c.AuthMiddleware)
	route.Get("/", c.EmployeeContractController.List)
	adminRoute := route.Group("/", c.AdminMiddleware)
	adminRoute.Post("/", c.EmployeeContractController.Create)
}

func (c *RouteConfig) SetupDivisionRouter() {
	route := c.App.Group("/api/divisions", c.AuthMiddleware)
	route.Get("/", c.DivisionController.List)

	adminRoute := route.Group("/", c.AdminMiddleware)
	adminRoute.Post("/", c.DivisionController.Create)
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
	route.Get("/:officeLocationID", c.OfficeLocationController.Detail)
	route.Post("/assign-employee", c.OfficeLocationController.AssignEmployee)
}

func (c *RouteConfig) SetupAttendanceRouter() {
	route := c.App.Group("/api/attendances", c.AuthMiddleware)
	employeeRoute := route.Group("/", c.EmployeeMiddleware)
	employeeRoute.Post("/check-in", c.AttendanceController.CheckIn)
}

func (c *RouteConfig) SetupShiftRouter() {
	route := c.App.Group("/api/shifts", c.AuthMiddleware, c.AdminMiddleware)
	route.Get("/", c.ShiftController.List)
	route.Post("/", c.ShiftController.Create)
	route.Get("/:shiftID", c.ShiftController.Detail)
	route.Post("/assign-employee", c.ShiftController.AssignEmployee)
}

func (c *RouteConfig) SetupTimeOffRouter() {
	route := c.App.Group("/api/time-off-requests", c.AuthMiddleware)
	route.Get("/", c.TimeOffRequestController.ListRequests)

	employeeRoute := route.Group("/", c.EmployeeMiddleware)
	employeeRoute.Post("/", c.TimeOffRequestController.CreateRequest)
	employeeRoute.Get("/_current", c.TimeOffRequestController.ListCurrentRequests)

	route.Get("/:id", c.TimeOffRequestController.GetRequestByID)
	route.Get("/:id/approvals", c.TimeOffApprovalController.ListApprovals)

	employeeRoute.Patch("/:id/approvals/:approval_id", c.TimeOffApprovalController.Decide)

	typeRoute := c.App.Group("/api/time-off-types", c.AuthMiddleware)
	typeRoute.Get("/", c.TimeOffTypeController.ListTypes)
	typeAdminRoute := typeRoute.Group("/", c.AdminMiddleware)
	typeAdminRoute.Post("/", c.TimeOffTypeController.CreateType)

	balanceRoute := c.App.Group("/api/time-off-balances", c.AuthMiddleware)
	balanceRoute.Post("/_set", c.AdminMiddleware, c.TimeOffBalanceController.SetBalance)
	balanceRoute.Get("/", c.AdminMiddleware, c.TimeOffBalanceController.ListBalancesByEmployee)
	balanceRoute.Get("/_current", c.EmployeeMiddleware, c.TimeOffBalanceController.ListCurrentBalances)
}

func (c *RouteConfig) SetupCommonRouter() {
	route := c.App.Group("/api", c.AuthMiddleware)
	route.Post("/upload", c.UploadController.Upload)
}

func (c *RouteConfig) SetupTimeOffApprovalRouter() {
	route := c.App.Group("/api/time-off-approvals", c.AuthMiddleware)
	employeeRoute := route.Group("/", c.EmployeeMiddleware)
	employeeRoute.Get("/_current", c.TimeOffApprovalController.ListMyApprovals)
	employeeRoute.Patch("/:approval_id", c.TimeOffApprovalController.DecideShort)
}

func (c *RouteConfig) SetupVisitRouter() {
	route := c.App.Group("/api/visits", c.AuthMiddleware)
	route.Get("/", c.VisitController.List)

	employeeRoute := route.Group("/", c.EmployeeMiddleware)
	employeeRoute.Post("/", c.VisitController.Create)
	employeeRoute.Get("/_current", c.VisitController.ListCurrent)
	employeeRoute.Get("/_current/can-do", c.VisitController.CanDoVisit)
	employeeRoute.Get("/_current/unclosed", c.VisitController.GetUnclosedVisit)

	route.Get("/:id", c.VisitController.GetByID)
	route.Delete("/:id", c.VisitController.Delete)
}

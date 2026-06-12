package route

import (
	"hr-sas/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App *fiber.App

	AuthMiddleware     fiber.Handler
	AdminMiddleware    func(permission string) fiber.Handler
	EmployeeMiddleware fiber.Handler

	CompanyController           *http.CompanyController
	UserController              *http.UserController
	EmployeeController          *http.EmployeeController
	OfficeLocationController    *http.OfficeLocationController
	SanctionController          *http.SanctionController
	EmSancController            *http.EmSancController
	HolidayController           *http.HolidayController
	PositionController          *http.PositionController
	AttendanceController        *http.AttendanceController
	ShiftController             *http.ShiftController
	TimeOffRequestController    *http.TimeOffRequestController
	TimeOffTypeController       *http.TimeOffTypeController
	TimeOffBalanceController    *http.TimeOffBalanceController
	TimeOffApprovalController   *http.TimeOffApprovalController
	UploadController            *http.UploadController
	EmployeeContractController  *http.EmployeeContractController
	DivisionController          *http.DivisionController
	VisitController             *http.VisitController
	PermissionController        *http.PermissionController
	RoleController              *http.RoleController
	EmployeeDocumentController  *http.EmployeeDocumentController
	EmployeeEducationController *http.EmployeeEducationController
	EmployeeTrainingController  *http.EmployeeTrainingController
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
	c.SetupPermissionRouter()
	c.SetupRoleRouter()
	c.SetupEmployeeDocumentRouter()
	c.SetupHolidayRouter()
	c.SetupEmployeeEducationRouter()
	c.SetupEmployeeTrainingRouter()
}

func (c *RouteConfig) SetupGuestRouter() {
	c.App.Post("/api/_login", c.UserController.Login)
	c.App.Post("/api/_register", c.UserController.Register)
	c.App.Delete("/api/users/_logout", c.UserController.Logout)
}

func (c *RouteConfig) SetupCompanyRouter() {
	route := c.App.Group("/api/companies", c.AuthMiddleware)
	route.Post("/", c.CompanyController.Create)
	route.Post("/_register", c.CompanyController.Register)
}

func (c *RouteConfig) SetupUserRouter() {
	route := c.App.Group("/api/users", c.AuthMiddleware)
	route.Get("/_current", c.UserController.GetCurrentUser)
	route.Patch("/_change-password", c.UserController.ChangePassword)

	adminMW := c.AdminMiddleware("USERS")
	route.Get("/", adminMW, c.UserController.List)
	route.Get("/:id", adminMW, c.UserController.Detail)
	route.Put("/:id", adminMW, c.UserController.Update)
	route.Delete("/:id", adminMW, c.UserController.Delete)
	route.Patch("/:id/_reset-password", adminMW, c.UserController.ResetPassword)
}

func (c *RouteConfig) SetupEmployeeRouter() {
	route := c.App.Group("/api/employees", c.AuthMiddleware)
	route.Get("/", c.EmployeeController.ListEmployee)

	adminMW := c.AdminMiddleware("EMPLOYEES")
	route.Post("/", adminMW, c.EmployeeController.CreateEmployee)
	route.Post("/import-excel", adminMW, c.EmployeeController.ImportExcel)
	route.Get("/:id", adminMW, c.EmployeeController.DetailEmployee)
	route.Put("/:id", adminMW, c.EmployeeController.UpdateEmployee)
	route.Delete("/:id", adminMW, c.EmployeeController.DeleteEmployee)
}

func (c *RouteConfig) SetupEmployeeContractRouter() {
	route := c.App.Group("/api/employee-contracts", c.AuthMiddleware)
	route.Get("/", c.EmployeeContractController.List)

	adminMW := c.AdminMiddleware("EMPLOYEE_CONTRACTS")
	route.Post("/", adminMW, c.EmployeeContractController.Create)
	route.Get("/:id", adminMW, c.EmployeeContractController.Detail)
	route.Put("/:id", adminMW, c.EmployeeContractController.Update)
	route.Delete("/:id", adminMW, c.EmployeeContractController.Delete)
}

func (c *RouteConfig) SetupDivisionRouter() {
	route := c.App.Group("/api/divisions", c.AuthMiddleware)
	route.Get("/", c.DivisionController.List)
	route.Get("/:id", c.DivisionController.Detail)

	adminMW := c.AdminMiddleware("DIVISIONS")
	route.Post("/", adminMW, c.DivisionController.Create)
	route.Put("/:id", adminMW, c.DivisionController.Update)
	route.Delete("/:id", adminMW, c.DivisionController.Delete)
}

func (c *RouteConfig) SetupSanctionRouter() {
	route := c.App.Group("/api/sanctions", c.AuthMiddleware)
	adminMW := c.AdminMiddleware("SANCTIONS")
	route.Post("/", adminMW, c.SanctionController.Create)
	route.Get("/", adminMW, c.SanctionController.ListSanction)
	route.Get("/:id", adminMW, c.SanctionController.Detail)
	route.Put("/:id", adminMW, c.SanctionController.Update)
	route.Delete("/:id", adminMW, c.SanctionController.Delete)
}

func (c *RouteConfig) SetupEmployeeSanctionRouter() {
	route := c.App.Group("/api/employee-sanctions", c.AuthMiddleware)
	route.Get("/_current", c.EmployeeMiddleware, c.EmSancController.CurrentSearch)

	adminMW := c.AdminMiddleware("EMPLOYEE_SANCTIONS")
	route.Post("/", adminMW, c.EmSancController.Create)
	route.Get("/", adminMW, c.EmSancController.Search)
	route.Get("/:id", adminMW, c.EmSancController.Detail)
	route.Put("/:id", adminMW, c.EmSancController.Update)
	route.Delete("/:id", adminMW, c.EmSancController.Delete)
}

func (c *RouteConfig) SetupPositionRouter() {
	route := c.App.Group("/api/positions", c.AuthMiddleware)
	adminMW := c.AdminMiddleware("POSITIONS")
	route.Get("/", adminMW, c.PositionController.ListPosition)
	route.Post("/", adminMW, c.PositionController.Create)
	route.Get("/:id", adminMW, c.PositionController.Detail)
	route.Put("/:id", adminMW, c.PositionController.Update)
	route.Delete("/:id", adminMW, c.PositionController.Delete)
}

func (c *RouteConfig) SetupOfficeLocationRouter() {
	route := c.App.Group("/api/office-locations", c.AuthMiddleware)
	adminMW := c.AdminMiddleware("OFFICE_LOCATIONS")
	route.Get("/", adminMW, c.OfficeLocationController.List)
	route.Post("/", adminMW, c.OfficeLocationController.Create)
	route.Get("/:officeLocationID", adminMW, c.OfficeLocationController.Detail)
	route.Put("/:officeLocationID", adminMW, c.OfficeLocationController.Update)
	route.Delete("/:officeLocationID", adminMW, c.OfficeLocationController.Delete)
	route.Post("/assign-employee", adminMW, c.OfficeLocationController.AssignEmployee)
}

func (c *RouteConfig) SetupAttendanceRouter() {
	route := c.App.Group("/api/attendances", c.AuthMiddleware)
	route.Post("/check-in", c.EmployeeMiddleware, c.AttendanceController.CheckIn)
	route.Post("/check-out", c.EmployeeMiddleware, c.AttendanceController.CheckOut)
	route.Post("/break-in", c.EmployeeMiddleware, c.AttendanceController.BreakIn)
	route.Post("/break-out", c.EmployeeMiddleware, c.AttendanceController.BreakOut)
}

func (c *RouteConfig) SetupShiftRouter() {
	route := c.App.Group("/api/shifts", c.AuthMiddleware)
	adminMW := c.AdminMiddleware("SHIFTS")
	route.Get("/", adminMW, c.ShiftController.List)
	route.Post("/", adminMW, c.ShiftController.Create)
	route.Get("/:shiftID", adminMW, c.ShiftController.Detail)
	route.Put("/:shiftID", adminMW, c.ShiftController.Update)
	route.Delete("/:shiftID", adminMW, c.ShiftController.Delete)
	route.Post("/assign-employee", adminMW, c.ShiftController.AssignEmployee)
}

func (c *RouteConfig) SetupTimeOffRouter() {
	// TIME OFF REQUEST ROUTES
	route := c.App.Group("/api/time-off-requests", c.AuthMiddleware)
	timeOffAdminMW := c.AdminMiddleware("TIME_OFF_REQUESTS")

	// employee only
	route.Post("/", c.EmployeeMiddleware, c.TimeOffRequestController.CreateRequest)
	route.Get("/_current", c.EmployeeMiddleware, c.TimeOffRequestController.ListCurrentRequests)

	// admin only
	route.Get("/", timeOffAdminMW, c.TimeOffRequestController.ListRequests)
	route.Delete("/:id", timeOffAdminMW, c.TimeOffRequestController.DeleteRequest)
	route.Post("/:employee_id", timeOffAdminMW, c.TimeOffRequestController.AdminCreateRequest)

	// shared
	route.Get("/:id/detail", c.TimeOffRequestController.GetRequestByID)
	route.Put("/:id", c.TimeOffRequestController.UpdateRequest)
	route.Get("/:id/approvals", c.TimeOffApprovalController.ListApprovals)

	// TIME OFF TYPE ROUTES
	typeRoute := c.App.Group("/api/time-off-types", c.AuthMiddleware)
	typeRoute.Get("/", c.TimeOffTypeController.ListTypes)
	typeRoute.Get("/:id", c.TimeOffTypeController.Detail)

	typeAdminMW := c.AdminMiddleware("TIME_OFF_TYPES")
	typeRoute.Post("/", typeAdminMW, c.TimeOffTypeController.CreateType)
	typeRoute.Put("/:id", typeAdminMW, c.TimeOffTypeController.UpdateType)
	typeRoute.Delete("/:id", typeAdminMW, c.TimeOffTypeController.DeleteType)

	// TIME OFF BALANCE ROUTES
	balanceRoute := c.App.Group("/api/time-off-balances", c.AuthMiddleware)
	balanceRoute.Get("/_current", c.EmployeeMiddleware, c.TimeOffBalanceController.ListCurrentBalances)
	balanceRoute.Get("/:id", c.TimeOffBalanceController.Detail)

	balanceAdminMW := c.AdminMiddleware("TIME_OFF_BALANCES")
	balanceRoute.Post("/_set", balanceAdminMW, c.TimeOffBalanceController.SetBalance)
	balanceRoute.Get("/", balanceAdminMW, c.TimeOffBalanceController.ListBalancesByEmployee)
	balanceRoute.Put("/:id", balanceAdminMW, c.TimeOffBalanceController.Update)
	balanceRoute.Delete("/:id", balanceAdminMW, c.TimeOffBalanceController.Delete)
}

func (c *RouteConfig) SetupCommonRouter() {
	route := c.App.Group("/api", c.AuthMiddleware)
	route.Post("/upload", c.UploadController.Upload)
	route.Post("/uploads", c.UploadController.Uploads)
}

func (c *RouteConfig) SetupTimeOffApprovalRouter() {
	route := c.App.Group("/api/time-off-approvals", c.AuthMiddleware)
	route.Get("/_current", c.EmployeeMiddleware, c.TimeOffApprovalController.ListMyApprovals)
	route.Patch("/:approval_id", c.EmployeeMiddleware, c.TimeOffApprovalController.DecideShort)
}

func (c *RouteConfig) SetupPermissionRouter() {
	route := c.App.Group("/api/permissions", c.AuthMiddleware)
	adminMW := c.AdminMiddleware("PERMISSIONS")
	route.Get("/", adminMW, c.PermissionController.List)
	route.Post("/", adminMW, c.PermissionController.Create)
	route.Get("/:id", adminMW, c.PermissionController.Detail)
	route.Put("/:id", adminMW, c.PermissionController.Update)
	route.Delete("/:id", adminMW, c.PermissionController.Delete)
}

func (c *RouteConfig) SetupRoleRouter() {
	route := c.App.Group("/api/roles", c.AuthMiddleware)
	adminMW := c.AdminMiddleware("ROLES")
	route.Get("/", adminMW, c.RoleController.List)
	route.Post("/", adminMW, c.RoleController.Create)
	route.Get("/:id", adminMW, c.RoleController.Detail)
	route.Put("/:id", adminMW, c.RoleController.Update)
	route.Delete("/:id", adminMW, c.RoleController.Delete)
	route.Post("/:id/permissions", adminMW, c.RoleController.AssignPermissions)
}

func (c *RouteConfig) SetupEmployeeDocumentRouter() {
	route := c.App.Group("/api/employee-docs", c.AuthMiddleware)
	adminMW := c.AdminMiddleware("EMPLOYEE_DOCUMENTS")
	route.Post("/", adminMW, c.EmployeeDocumentController.Create)
	route.Get("/", adminMW, c.EmployeeDocumentController.List)
	route.Put("/:id", adminMW, c.EmployeeDocumentController.Update)
	route.Delete("/:id", adminMW, c.EmployeeDocumentController.Delete)
}

func (c *RouteConfig) SetupVisitRouter() {
	route := c.App.Group("/api/visits", c.AuthMiddleware)
	visitAdminMW := c.AdminMiddleware("VISITS")

	// shared
	route.Get("/:id/detail", c.VisitController.GetByID)

	// employee only
	route.Post("/", c.EmployeeMiddleware, c.VisitController.Create)
	route.Get("/_current", c.EmployeeMiddleware, c.VisitController.ListCurrent)
	route.Get("/_current/can-do", c.EmployeeMiddleware, c.VisitController.CanDoVisit)
	route.Get("/_current/unclosed", c.EmployeeMiddleware, c.VisitController.GetUnclosedVisit)

	// admin only
	route.Get("/", visitAdminMW, c.VisitController.List)
	route.Put("/:id", visitAdminMW, c.VisitController.Update)
	route.Delete("/:id", visitAdminMW, c.VisitController.Delete)
}

func (c *RouteConfig) SetupEmployeeEducationRouter() {
	route := c.App.Group("/api/employee-educations", c.AuthMiddleware)
	adminMW := c.AdminMiddleware("EMPLOYEE_EDUCATIONS")
	route.Post("/", adminMW, c.EmployeeEducationController.Create)
	route.Get("/", adminMW, c.EmployeeEducationController.List)
	route.Get("/:education_id", adminMW, c.EmployeeEducationController.Detail)
	route.Put("/:education_id", adminMW, c.EmployeeEducationController.Update)
	route.Delete("/:education_id", adminMW, c.EmployeeEducationController.Delete)
}

func (c *RouteConfig) SetupEmployeeTrainingRouter() {
	route := c.App.Group("/api/employee-trainings", c.AuthMiddleware)
	adminMW := c.AdminMiddleware("EMPLOYEE_TRAININGS")
	route.Post("/", adminMW, c.EmployeeTrainingController.Create)
	route.Get("/", adminMW, c.EmployeeTrainingController.List)
	route.Get("/:training_id", adminMW, c.EmployeeTrainingController.Detail)
	route.Put("/:training_id", adminMW, c.EmployeeTrainingController.Update)
	route.Delete("/:training_id", adminMW, c.EmployeeTrainingController.Delete)
}

func (c *RouteConfig) SetupHolidayRouter() {
	route := c.App.Group("/api/holidays", c.AuthMiddleware)
	adminMW := c.AdminMiddleware("HOLIDAYS")
	route.Post("/", adminMW, c.HolidayController.Create)
	route.Get("/", adminMW, c.HolidayController.List)
	route.Put("/:id", adminMW, c.HolidayController.Update)
	route.Delete("/:id", adminMW, c.HolidayController.Delete)
}

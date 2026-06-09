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

/*
Guest Router
*/
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

	adminRoute := route.Group("/", c.AdminMiddleware("USERS"))
	adminRoute.Get("/", c.UserController.List)
	adminRoute.Get("/:id", c.UserController.Detail)
	adminRoute.Put("/:id", c.UserController.Update)
	adminRoute.Delete("/:id", c.UserController.Delete)
	adminRoute.Patch("/:id/_reset-password", c.UserController.ResetPassword)
}

func (c *RouteConfig) SetupEmployeeRouter() {
	route := c.App.Group("/api/employees", c.AuthMiddleware)
	route.Get("/", c.EmployeeController.ListEmployee)

	adminRoute := route.Group("/", c.AdminMiddleware("EMPLOYEES"))
	adminRoute.Post("/", c.EmployeeController.CreateEmployee)
	adminRoute.Post("/import-excel", c.EmployeeController.ImportExcel)
	adminRoute.Get("/:id", c.EmployeeController.DetailEmployee)
	adminRoute.Put("/:id", c.EmployeeController.UpdateEmployee)
	adminRoute.Delete("/:id", c.EmployeeController.DeleteEmployee)
}

func (c *RouteConfig) SetupEmployeeContractRouter() {
	route := c.App.Group("/api/employee-contracts", c.AuthMiddleware)
	route.Get("/", c.EmployeeContractController.List)
	adminRoute := route.Group("/", c.AdminMiddleware("EMPLOYEE_CONTRACTS"))
	adminRoute.Post("/", c.EmployeeContractController.Create)
	adminRoute.Get("/:id", c.EmployeeContractController.Detail)
	adminRoute.Put("/:id", c.EmployeeContractController.Update)
	adminRoute.Delete("/:id", c.EmployeeContractController.Delete)
}

func (c *RouteConfig) SetupDivisionRouter() {
	route := c.App.Group("/api/divisions", c.AuthMiddleware)
	route.Get("/", c.DivisionController.List)
	route.Get("/:id", c.DivisionController.Detail)

	adminRoute := route.Group("/", c.AdminMiddleware("DIVISIONS"))
	adminRoute.Post("/", c.DivisionController.Create)
	adminRoute.Put("/:id", c.DivisionController.Update)
	adminRoute.Delete("/:id", c.DivisionController.Delete)
}

func (c *RouteConfig) SetupSanctionRouter() {
	route := c.App.Group("/api/sanctions", c.AuthMiddleware, c.AdminMiddleware("SANCTIONS"))
	route.Post("/", c.SanctionController.Create)
	route.Get("/", c.SanctionController.ListSanction)
	route.Get("/:id", c.SanctionController.Detail)
	route.Put("/:id", c.SanctionController.Update)
	route.Delete("/:id", c.SanctionController.Delete)
}

func (c *RouteConfig) SetupEmployeeSanctionRouter() {
	route := c.App.Group("/api/employee-sanctions", c.AuthMiddleware)
	route.Get("/_current", c.EmployeeMiddleware, c.EmSancController.CurrentSearch)

	adminRouter := route.Group("/", c.AdminMiddleware("EMPLOYEE_SANCTIONS"))
	adminRouter.Post("/", c.EmSancController.Create)
	adminRouter.Get("/", c.EmSancController.Search)
	adminRouter.Get("/:id", c.EmSancController.Detail)
	adminRouter.Put("/:id", c.EmSancController.Update)
	adminRouter.Delete("/:id", c.EmSancController.Delete)
}

func (c *RouteConfig) SetupPositionRouter() {
	route := c.App.Group("/api/positions", c.AuthMiddleware, c.AdminMiddleware("POSITIONS"))
	route.Get("/", c.PositionController.ListPosition)
	route.Post("/", c.PositionController.Create)
	route.Get("/:id", c.PositionController.Detail)
	route.Put("/:id", c.PositionController.Update)
	route.Delete("/:id", c.PositionController.Delete)
}

func (c *RouteConfig) SetupOfficeLocationRouter() {
	route := c.App.Group("/api/office-locations", c.AuthMiddleware, c.AdminMiddleware("OFFICE_LOCATIONS"))
	route.Get("/", c.OfficeLocationController.List)
	route.Post("/", c.OfficeLocationController.Create)
	route.Get("/:officeLocationID", c.OfficeLocationController.Detail)
	route.Put("/:officeLocationID", c.OfficeLocationController.Update)
	route.Delete("/:officeLocationID", c.OfficeLocationController.Delete)
	route.Post("/assign-employee", c.OfficeLocationController.AssignEmployee)
}

func (c *RouteConfig) SetupAttendanceRouter() {
	route := c.App.Group("/api/attendances", c.AuthMiddleware)
	employeeRoute := route.Group("/", c.EmployeeMiddleware)
	employeeRoute.Post("/check-in", c.AttendanceController.CheckIn)
	employeeRoute.Post("/check-out", c.AttendanceController.CheckOut)
	employeeRoute.Post("/break-in", c.AttendanceController.BreakIn)
	employeeRoute.Post("/break-out", c.AttendanceController.BreakOut)
}

func (c *RouteConfig) SetupShiftRouter() {
	route := c.App.Group("/api/shifts", c.AuthMiddleware, c.AdminMiddleware("SHIFTS"))
	route.Get("/", c.ShiftController.List)
	route.Post("/", c.ShiftController.Create)
	route.Get("/:shiftID", c.ShiftController.Detail)
	route.Put("/:shiftID", c.ShiftController.Update)
	route.Delete("/:shiftID", c.ShiftController.Delete)
	route.Post("/assign-employee", c.ShiftController.AssignEmployee)
}

func (c *RouteConfig) SetupTimeOffRouter() {
	// TIME OFF REQUEST ROUTES
	// admin and employee
	timeOffRequest := c.App.Group("/api/time-off-requests", c.AuthMiddleware)
	timeOffRequest.Get("/:id/detail", c.TimeOffRequestController.GetRequestByID)
	timeOffRequest.Put("/:id", c.TimeOffRequestController.UpdateRequest)
	timeOffRequest.Get("/:id/approvals", c.TimeOffApprovalController.ListApprovals)

	// admin only
	adminRoute := timeOffRequest.Group("/", c.AdminMiddleware("TIME_OFF_REQUESTS"))
	adminRoute.Get("/", c.TimeOffRequestController.ListRequests)
	adminRoute.Delete("/:id", c.TimeOffRequestController.DeleteRequest)
	adminRoute.Post("/:employee_id", c.TimeOffRequestController.AdminCreateRequest)

	// employee only
	employeeRoute := timeOffRequest.Group("/", c.EmployeeMiddleware)
	employeeRoute.Post("/", c.TimeOffRequestController.CreateRequest)
	employeeRoute.Get("/_current", c.TimeOffRequestController.ListCurrentRequests)

	// TIME OFF REQUEST ROUTES
	// admin and employee
	typeRoute := c.App.Group("/api/time-off-types", c.AuthMiddleware)
	typeRoute.Get("/", c.TimeOffTypeController.ListTypes)
	typeRoute.Get("/:id", c.TimeOffTypeController.Detail)

	// admin only
	typeAdminRoute := typeRoute.Group("/", c.AdminMiddleware("TIME_OFF_TYPES"))
	typeAdminRoute.Post("/", c.TimeOffTypeController.CreateType)
	typeAdminRoute.Put("/:id", c.TimeOffTypeController.UpdateType)
	typeAdminRoute.Delete("/:id", c.TimeOffTypeController.DeleteType)

	// TIME OFF BALANCE ROUTES
	balanceRoute := c.App.Group("/api/time-off-balances", c.AuthMiddleware)
	balanceRoute.Get("/_current", c.EmployeeMiddleware, c.TimeOffBalanceController.ListCurrentBalances)
	balanceRoute.Get("/:id", c.TimeOffBalanceController.Detail)

	balanceAdminRoute := balanceRoute.Group("/", c.AdminMiddleware("TIME_OFF_BALANCES"))
	balanceAdminRoute.Post("/_set", c.TimeOffBalanceController.SetBalance)
	balanceAdminRoute.Get("/", c.TimeOffBalanceController.ListBalancesByEmployee)
	balanceAdminRoute.Put("/:id", c.TimeOffBalanceController.Update)
	balanceAdminRoute.Delete("/:id", c.TimeOffBalanceController.Delete)
}

func (c *RouteConfig) SetupCommonRouter() {
	route := c.App.Group("/api", c.AuthMiddleware)
	route.Post("/upload", c.UploadController.Upload)
	route.Post("/uploads", c.UploadController.Uploads)
}

func (c *RouteConfig) SetupTimeOffApprovalRouter() {
	route := c.App.Group("/api/time-off-approvals", c.AuthMiddleware)
	employeeRoute := route.Group("/", c.EmployeeMiddleware)
	employeeRoute.Get("/_current", c.TimeOffApprovalController.ListMyApprovals)
	employeeRoute.Patch("/:approval_id", c.TimeOffApprovalController.DecideShort)
}

func (c *RouteConfig) SetupPermissionRouter() {
	route := c.App.Group("/api/permissions", c.AuthMiddleware, c.AdminMiddleware("PERMISSIONS"))
	route.Get("/", c.PermissionController.List)
	route.Post("/", c.PermissionController.Create)
	route.Get("/:id", c.PermissionController.Detail)
	route.Put("/:id", c.PermissionController.Update)
	route.Delete("/:id", c.PermissionController.Delete)
}

func (c *RouteConfig) SetupRoleRouter() {
	route := c.App.Group("/api/roles", c.AuthMiddleware, c.AdminMiddleware("ROLES"))
	route.Get("/", c.RoleController.List)
	route.Post("/", c.RoleController.Create)
	route.Get("/:id", c.RoleController.Detail)
	route.Put("/:id", c.RoleController.Update)
	route.Delete("/:id", c.RoleController.Delete)
	route.Post("/:id/permissions", c.RoleController.AssignPermissions)
}

func (c *RouteConfig) SetupEmployeeDocumentRouter() {
	route := c.App.Group("/api/employee-docs", c.AuthMiddleware, c.AdminMiddleware("EMPLOYEE_DOCUMENTS"))
	route.Post("/", c.EmployeeDocumentController.Create)
	route.Get("/", c.EmployeeDocumentController.List)
	route.Put("/:id", c.EmployeeDocumentController.Update)
	route.Delete("/:id", c.EmployeeDocumentController.Delete)
}

func (c *RouteConfig) SetupVisitRouter() {
	route := c.App.Group("/api/visits", c.AuthMiddleware)
	route.Get("/:id/detail", c.VisitController.GetByID)

	employeeRoute := route.Group("/", c.EmployeeMiddleware)
	employeeRoute.Post("/", c.VisitController.Create)
	employeeRoute.Get("/_current", c.VisitController.ListCurrent)
	employeeRoute.Get("/_current/can-do", c.VisitController.CanDoVisit)
	employeeRoute.Get("/_current/unclosed", c.VisitController.GetUnclosedVisit)

	adminRoute := route.Group("/", c.AdminMiddleware("VISITS"))
	adminRoute.Get("/", c.VisitController.List)
	adminRoute.Put("/:id", c.VisitController.Update)
	adminRoute.Delete("/:id", c.VisitController.Delete)
}

func (c *RouteConfig) SetupEmployeeEducationRouter() {
	route := c.App.Group("/api/employee-educations", c.AuthMiddleware, c.AdminMiddleware("EMPLOYEE_EDUCATIONS"))
	route.Post("/", c.EmployeeEducationController.Create)
	route.Get("/", c.EmployeeEducationController.List)
	route.Get("/:education_id", c.EmployeeEducationController.Detail)
	route.Put("/:education_id", c.EmployeeEducationController.Update)
	route.Delete("/:education_id", c.EmployeeEducationController.Delete)
}

func (c *RouteConfig) SetupEmployeeTrainingRouter() {
	route := c.App.Group("/api/employee-trainings", c.AuthMiddleware, c.AdminMiddleware("EMPLOYEE_TRAININGS"))
	route.Post("/", c.EmployeeTrainingController.Create)
	route.Get("/", c.EmployeeTrainingController.List)
	route.Get("/:training_id", c.EmployeeTrainingController.Detail)
	route.Put("/:training_id", c.EmployeeTrainingController.Update)
	route.Delete("/:training_id", c.EmployeeTrainingController.Delete)
}

func (c *RouteConfig) SetupHolidayRouter() {
	route := c.App.Group("/api/holidays", c.AuthMiddleware, c.AdminMiddleware("HOLIDAYS"))
	route.Post("/", c.HolidayController.Create)
	route.Get("/", c.HolidayController.List)
	route.Put("/:id", c.HolidayController.Update)
	route.Delete("/:id", c.HolidayController.Delete)
}

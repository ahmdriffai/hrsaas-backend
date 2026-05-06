package config

import (
	"hr-sas/internal/delivery/http"
	"hr-sas/internal/delivery/http/middleware"
	"hr-sas/internal/delivery/http/route"
	"hr-sas/internal/repository"
	"hr-sas/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup repository
	companyRepository := repository.NewCompanyRepository(config.Log)
	userRepository := repository.NewUserRepository(config.Log)
	sessionRepository := repository.NewSessionRepository(config.Log)
	employeeRepository := repository.NewEmployeeRepository(config.Log)
	sanctionRepository := repository.NewSanctionRepository(config.Log)
	emSancRepository := repository.NewEmSancRepository(config.Log)
	positionRepository := repository.NewPositionRepository(config.Log)
	officeLocationRepositoruy := repository.NewOfficeLocationRepository(config.Log)
	attendaceRepositpry := repository.NewAttendanceRepository(config.Log)
	attendanceLogRepository := repository.NewAttendanceLogRepository(config.Log)
	shifRepository := repository.NewShiftRepository(config.Log)
	shiftDayRepository := repository.NewShiftDayRepository(config.Log)
	timeOffTypeRepository := repository.NewTimeOffTypeRepository(config.Log)
	timeOffRequestRepository := repository.NewTimeOffRequestRepository(config.Log)
	timeOffBalanceRepository := repository.NewTimeOffBalanceRepository(config.Log)
	timeOffApprovalRepository := repository.NewTimeOffApprovalRepository(config.Log)
	employeeContractRepository := repository.NewEmployeeContractRepository(config.Log)
	divisionRepository := repository.NewDivisionRepository(config.Log)
	visitRepository := repository.NewVisitRepository(config.Log)
	permissionRepository := repository.NewPermissionRepository(config.Log)
	roleRepository := repository.NewRoleRepository(config.Log)

	// setup producer

	// setup usecase
	companyUsecase := usecase.NewCompanyUseCase(config.DB, config.Log, config.Validate, companyRepository, userRepository)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, sessionRepository, companyRepository, roleRepository)
	employeeUseCase := usecase.NewEmployeeUseCase(
		config.DB,
		config.Log,
		config.Validate,
		employeeRepository,
		userRepository,
		employeeContractRepository,
		positionRepository,
		divisionRepository,
	)
	sanctionUseCase := usecase.NewSantionUseCase(config.DB, config.Log, config.Validate, sanctionRepository)
	emSancUseCase := usecase.NewEmSancUseCase(config.DB, config.Log, config.Validate, emSancRepository, sanctionRepository, employeeRepository)
	positionUseCase := usecase.NewPositionUseCase(config.DB, config.Log, config.Validate, positionRepository)
	officeLocationUseCase := usecase.NewOfficeLocationUseCase(config.DB, config.Log, config.Validate, officeLocationRepositoruy)
	attendanceUseCase := usecase.NewAttendanceUseCase(config.DB, config.Log, config.Validate, attendaceRepositpry, officeLocationRepositoruy, shifRepository, shiftDayRepository, attendanceLogRepository)
	shiftUseCase := usecase.NewShiftUseCase(config.DB, config.Log, config.Validate, shifRepository, shiftDayRepository)
	timeOffRequestUseCase := usecase.NewTimeOffRequestUseCase(
		config.DB,
		config.Log,
		config.Validate,
		timeOffRequestRepository,
		timeOffTypeRepository,
		timeOffBalanceRepository,
		timeOffApprovalRepository,
		employeeContractRepository,
	)
	timeOffTypeUseCase := usecase.NewTimeOffTypeUseCase(
		config.DB,
		config.Log,
		config.Validate,
		timeOffTypeRepository,
	)
	timeOffBalanceUseCase := usecase.NewTimeOffBalanceUseCase(
		config.DB,
		config.Log,
		config.Validate,
		timeOffBalanceRepository,
		timeOffTypeRepository,
	)
	timeOffApprovalUseCase := usecase.NewTimeOffApprovalUseCase(
		config.DB,
		config.Log,
		config.Validate,
		timeOffRequestRepository,
		timeOffTypeRepository,
		timeOffBalanceRepository,
	)
	employeeContractUseCase := usecase.NewEmployeeContractUseCase(
		config.DB,
		config.Log,
		config.Validate,
		employeeContractRepository,
		timeOffTypeRepository,
		timeOffBalanceRepository,
	)
	divisionUseCase := usecase.NewDivisionUseCase(config.DB, config.Log, config.Validate, divisionRepository)
	visitUseCase := usecase.NewVisitUseCase(config.DB, config.Log, config.Validate, visitRepository)
	permissionUseCase := usecase.NewPermissionUseCase(config.DB, config.Log, config.Validate, permissionRepository)
	roleUseCase := usecase.NewRoleUseCase(config.DB, config.Log, config.Validate, roleRepository, permissionRepository)

	// setup controller
	companyController := http.NewCompanyController(companyUsecase, config.Log)
	userController := http.NewUserController(userUseCase, config.Log)
	employeeController := http.NewEmployeeController(employeeUseCase, config.Log)
	santionController := http.NewSanctionController(sanctionUseCase, config.Log)
	emSangController := http.NewEmSancController(emSancUseCase, config.Log)
	positionController := http.NewPositionController(positionUseCase, config.Log)
	officeLocationController := http.NewOfficeLocationController(officeLocationUseCase, config.Log)
	attendaceController := http.NewAttendanceController(attendanceUseCase, config.Log)
	shiftController := http.NewShifController(shiftUseCase, config.Log)
	timeOffRequestController := http.NewTimeOffRequestController(timeOffRequestUseCase, config.Log)
	timeOffTypeController := http.NewTimeOffTypeController(timeOffTypeUseCase, config.Log)
	timeOffBalanceController := http.NewTimeOffBalanceController(timeOffBalanceUseCase, config.Log)
	timeOffApprovalController := http.NewTimeOffApprovalController(timeOffApprovalUseCase, timeOffRequestUseCase, config.Log)
	uploadController := http.NewUploadController(config.Log)
	employeeContractController := http.NewEmployeeContractController(employeeContractUseCase, config.Log)
	divisionController := http.NewDivisionController(divisionUseCase, config.Log)
	visitController := http.NewVisitController(visitUseCase, config.Log)
	permissionController := http.NewPermissionController(permissionUseCase, config.Log)
	roleController := http.NewRoleController(roleUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)
	adminMiddleware := middleware.NewAdmin()
	employeeMiddleware := middleware.NewEmployee()

	// route config
	routeConfig := route.RouteConfig{
		App: config.App,

		AuthMiddleware:     authMiddleware,
		AdminMiddleware:    adminMiddleware,
		EmployeeMiddleware: employeeMiddleware,

		CompanyController:          companyController,
		UserController:             userController,
		EmployeeController:         employeeController,
		SanctionController:         santionController,
		EmSancController:           emSangController,
		PositionController:         positionController,
		OfficeLocationController:   officeLocationController,
		AttendanceController:       attendaceController,
		ShiftController:            shiftController,
		TimeOffRequestController:   timeOffRequestController,
		TimeOffTypeController:      timeOffTypeController,
		TimeOffBalanceController:   timeOffBalanceController,
		TimeOffApprovalController:  timeOffApprovalController,
		UploadController:           uploadController,
		EmployeeContractController: employeeContractController,
		DivisionController:         divisionController,
		VisitController:            visitController,
		PermissionController:       permissionController,
		RoleController:             roleController,
	}

	routeConfig.Setup()
}

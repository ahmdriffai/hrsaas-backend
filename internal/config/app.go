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
	timeOffAttachmentRepository := repository.NewTimeOffAttachmentRepository(config.Log)

	// setup producer

	// setup usecase
	companyUsecase := usecase.NewCompanyUseCase(config.DB, config.Log, config.Validate, companyRepository, userRepository)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, sessionRepository, companyRepository)
	employeeUseCase := usecase.NewEmployeeUseCase(config.DB, config.Log, config.Validate, employeeRepository, userRepository)
	sanctionUseCase := usecase.NewSantionUseCase(config.DB, config.Log, config.Validate, sanctionRepository)
	emSancUseCase := usecase.NewEmSancUseCase(config.DB, config.Log, config.Validate, emSancRepository, sanctionRepository, employeeRepository)
	positionUseCase := usecase.NewPositionUseCase(config.DB, config.Log, config.Validate, positionRepository)
	officeLocationUseCase := usecase.NewOfficeLocationUseCase(config.DB, config.Log, config.Validate, officeLocationRepositoruy)
	attendanceUseCase := usecase.NewAttendanceUseCase(config.DB, config.Log, config.Validate, attendaceRepositpry, officeLocationRepositoruy, shifRepository, shiftDayRepository, attendanceLogRepository)
	shiftUseCase := usecase.NewShiftUseCase(config.DB, config.Log, config.Validate, shifRepository, shiftDayRepository)
	timeOffUseCase := usecase.NewTimeOffUseCase(
		config.DB,
		config.Log,
		config.Validate,
		timeOffRequestRepository,
		timeOffTypeRepository,
		timeOffBalanceRepository,
		timeOffApprovalRepository,
		timeOffAttachmentRepository,
	)

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
	timeOffController := http.NewTimeOffController(timeOffUseCase, config.Log)

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

		CompanyController:        companyController,
		UserController:           userController,
		EmployeeController:       employeeController,
		SanctionController:       santionController,
		EmSancController:         emSangController,
		PositionController:       positionController,
		OfficeLocationController: officeLocationController,
		AttendanceController:     attendaceController,
		ShiftController:          shiftController,
		TimeOffController:        timeOffController,
	}

	routeConfig.Setup()
}

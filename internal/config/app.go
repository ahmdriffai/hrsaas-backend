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

	// setup producer

	// setup usecase
	companyUsecase := usecase.NewCompanyUseCase(config.DB, config.Log, config.Validate, companyRepository, userRepository)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, sessionRepository, companyRepository)
	employeeUseCase := usecase.NewEmployeeUseCase(config.DB, config.Log, config.Validate, employeeRepository)
	sanctionUseCase := usecase.NewSantionUseCase(config.DB, config.Log, config.Validate, sanctionRepository)
	emSancUseCase := usecase.NewEmSancUseCase(config.DB, config.Log, config.Validate, emSancRepository, sanctionRepository, employeeRepository)
	positionUseCase := usecase.NewPositionUseCase(config.DB, config.Log, config.Validate, positionRepository)

	// setup controller
	companyController := http.NewCompanyController(companyUsecase, config.Log)
	userController := http.NewUserController(userUseCase, config.Log)
	employeeController := http.NewEmployeeController(employeeUseCase, config.Log)
	santionController := http.NewSanctionController(sanctionUseCase, config.Log)
	emSangController := http.NewEmSancController(emSancUseCase, config.Log)
	positionController := http.NewPositionController(positionUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)

	// route config
	routeConfig := route.RouteConfig{
		App:                config.App,
		CompanyController:  companyController,
		AuthMiddleware:     authMiddleware,
		UserController:     userController,
		EmployeeController: employeeController,
		SanctionController: santionController,
		EmSancController:   emSangController,
		PositionController: positionController,
	}

	routeConfig.Setup()
}

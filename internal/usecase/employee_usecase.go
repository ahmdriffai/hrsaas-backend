package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/model/converter"
	"hr-sas/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type EmployeeUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	EmployeeRepository *repository.EmployeeRepository
	UserRepository     *repository.UserRepository
}

func NewEmployeeUseCase(db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	employeeRepository *repository.EmployeeRepository,
	userRepository *repository.UserRepository) *EmployeeUseCase {
	return &EmployeeUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		EmployeeRepository: employeeRepository,
		UserRepository:     userRepository,
	}
}

/* Create Employee Usecase
 */
func (c *EmployeeUseCase) Create(ctx context.Context, request *model.CreateEmployeeRequest) (*model.EmployeeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	// validate employee number
	total, err := c.EmployeeRepository.CountByEmployeeNumberAndCompanyID(tx, request.EmployeeNumber, request.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to count employee by number and company")
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.WithError(err).Error("Failed to count employee by number and company")
		return nil, fiber.NewError(fiber.StatusConflict, "Employee number already usage.")
	}

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.WithError(err).Error("Failed to hash password")
		return nil, fiber.ErrInternalServerError
	}

	// create user entity
	user := &entity.User{
		Name:      request.Fullname,
		Email:     request.Email,
		Password:  string(passwordHash),
		Role:      "USER",
		CompanyID: request.CompanyID,
	}

	// create user in database
	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.WithError(err).Error("Failed to create user")
		return nil, fiber.ErrInternalServerError
	}

	// create employee entity
	employee := &entity.Employee{
		Fullname:       request.Fullname,
		BirthPlace:     request.BirthPlace,
		BirthDate:      request.BirthDate,
		BlodType:       request.BlodType,
		MaritalStatus:  request.MaritalStatus,
		Religion:       request.Religion,
		Phone:          request.Phone,
		Timezone:       request.Timezone,
		CompanyID:      request.CompanyID,
		EmployeeNumber: request.EmployeeNumber,
		UserID:         user.ID,
	}

	if err := c.EmployeeRepository.Create(tx, employee); err != nil {
		c.Log.WithError(err).Error("Failed to create employee")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.EmployeeToResponse(employee), nil
}

/* Search Employee
 */

func (c *EmployeeUseCase) Search(ctx context.Context, request *model.SearchEmployeeRequest) ([]model.EmployeeResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	employees, total, err := c.EmployeeRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting employee")
		return nil, 0, fiber.ErrInternalServerError
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.EmployeeResponse, len(employees))
	for i, employee := range employees {
		responses[i] = *converter.EmployeeToResponse(&employee)
	}

	return responses, total, nil
}

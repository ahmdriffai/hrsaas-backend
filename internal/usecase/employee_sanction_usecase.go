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
	"gorm.io/gorm"
)

type EmSancUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	EmSancRepository   *repository.EmSancRepository
	SanctionRepository *repository.SanctionRepository
	EmployeeRepository *repository.EmployeeRepository
}

func NewEmSancUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, emSancRepository *repository.EmSancRepository, sanctionRepository *repository.SanctionRepository, employeeRepository *repository.EmployeeRepository) *EmSancUseCase {
	return &EmSancUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		EmSancRepository:   emSancRepository,
		SanctionRepository: sanctionRepository,
		EmployeeRepository: employeeRepository,
	}
}

/* Create Employee Sanction Usecase
 */
func (c *EmSancUseCase) Create(ctx context.Context, request *model.CreateEmSancRequest) (*model.EmSancResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	// check if sanction exists
	total, err := c.SanctionRepository.CountById(tx, request.SanctionID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to check sanction existence")
		return nil, fiber.ErrBadRequest
	}
	if total == 0 {
		c.Log.Error("Sanction not found")
		return nil, fiber.ErrBadRequest
	}

	// check if employee exists
	total, err = c.EmployeeRepository.CountById(tx, request.EmployeeID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to check employee existence")
		return nil, fiber.ErrBadRequest
	}
	if total == 0 {
		c.Log.Error("Employee not found")
		return nil, fiber.ErrBadRequest
	}

	status := "active"
	employeeSanction := &entity.EmployeeSanction{
		EmployeeID:  request.EmployeeID,
		SanctionID:  request.SanctionID,
		CompanyID:   request.CompanyID,
		Reason:      &request.Reason,
		StartDate:   request.StartDate.Time,
		EndDate:     &request.EndDate.Time,
		Status:      &status,
		DocumentUrl: request.DocumentUrl,
	}

	err = c.EmSancRepository.Create(tx, employeeSanction)
	if err != nil {
		c.Log.Error("Error creating employee sanction:", err)
		return nil, err
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.EmSancToResponse(employeeSanction), nil
}

/* Search Employee Sanction Usecase
 */
func (c *EmSancUseCase) Search(ctx context.Context, request *model.SearchEmSancRequest) ([]model.EmSancResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	emSancs, total, err := c.EmSancRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting sanction")
		return nil, 0, fiber.ErrInternalServerError
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.EmSancResponse, len(emSancs))
	for i, contact := range emSancs {
		responses[i] = *converter.EmSancToResponse(&contact)
	}

	return responses, total, nil
}

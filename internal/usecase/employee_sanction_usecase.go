package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strings"
	"time"

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

	startDate, _ := lib.ParseDateToUnixMilli(request.StartDate)
	endDate, _ := lib.ParseDateToUnixMilli(request.EndDate)

	status := "active"
	employeeSanction := &entity.EmployeeSanction{
		EmployeeID:  request.EmployeeID,
		SanctionID:  request.SanctionID,
		CompanyID:   request.CompanyID,
		Reason:      &request.Reason,
		StartDate:   startDate,
		EndDate:     &endDate,
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

	return model.EmSancToResponse(employeeSanction), nil
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
		responses[i] = *model.EmSancToResponse(&contact)
	}

	return responses, total, nil
}

func (c *EmSancUseCase) Detail(ctx context.Context, id string, companyID string) (*model.EmSancResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item := new(entity.EmployeeSanction)
	if err := c.EmSancRepository.FindByIdAndCompany(tx, item, id, companyID, "Employee", "Sanction"); err != nil {
		c.Log.WithError(err).Error("Employee sanction not found")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmSancToResponse(item), nil
}

func (c *EmSancUseCase) Update(ctx context.Context, id string, companyID string, request *model.UpdateEmSancRequest) (*model.EmSancResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	item := new(entity.EmployeeSanction)
	if err := c.EmSancRepository.FindByIdAndCompany(tx, item, id, companyID, "Employee", "Sanction"); err != nil {
		c.Log.WithError(err).Error("Employee sanction not found")
		return nil, fiber.ErrNotFound
	}

	var nextStartDate int64 = item.StartDate
	var nextEndDate *int64 = item.EndDate

	if request.Reason != nil {
		reason := strings.TrimSpace(*request.Reason)
		if reason == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "reason cannot be empty")
		}
		item.Reason = &reason
	}

	if request.StartDate != nil {
		startDate := strings.TrimSpace(*request.StartDate)
		if startDate == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "start_date cannot be empty")
		}

		parsedStartDate, err := lib.ParseDateToUnixMilli(startDate)
		if err != nil {
			c.Log.WithError(err).Error("Failed to parse sanction start date")
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid start_date format")
		}

		nextStartDate = parsedStartDate
		item.StartDate = parsedStartDate
	}

	if request.EndDate != nil {
		endDate := strings.TrimSpace(*request.EndDate)
		if endDate == "" {
			nextEndDate = nil
			item.EndDate = nil
		} else {
			parsedEndDate, err := lib.ParseDateToUnixMilli(endDate)
			if err != nil {
				c.Log.WithError(err).Error("Failed to parse sanction end date")
				return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid end_date format")
			}
			nextEndDate = &parsedEndDate
			item.EndDate = &parsedEndDate
		}
	}

	if nextEndDate != nil && *nextEndDate < nextStartDate {
		return nil, fiber.NewError(fiber.StatusBadRequest, "end_date cannot be earlier than start_date")
	}

	if request.Status != nil {
		status := strings.ToLower(strings.TrimSpace(*request.Status))
		if status == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "status cannot be empty")
		}
		item.Status = &status
	}

	if request.DocumentUrl != nil {
		item.DocumentUrl = strings.TrimSpace(*request.DocumentUrl)
	}

	item.UpdatedAt = time.Now().UnixMilli()

	if err := c.EmSancRepository.Update(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to update employee sanction")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmSancToResponse(item), nil
}

func (c *EmSancUseCase) Delete(ctx context.Context, id string, companyID string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item := new(entity.EmployeeSanction)
	if err := c.EmSancRepository.FindByIdAndCompany(tx, item, id, companyID); err != nil {
		c.Log.WithError(err).Error("Employee sanction not found")
		return fiber.ErrNotFound
	}

	if err := c.EmSancRepository.Delete(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to delete employee sanction")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

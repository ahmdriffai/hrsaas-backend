package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ShiftUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	ShiftRepository *repository.ShiftRepository
}

func NewShiftUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	shiftRepository *repository.ShiftRepository,
) *ShiftUseCase {
	return &ShiftUseCase{
		DB:              db,
		Log:             log,
		Validate:        validate,
		ShiftRepository: shiftRepository,
	}
}

// Create Shift
func (c *ShiftUseCase) Create(ctx context.Context, request *model.CreateShiftRequest) (*model.ShiftResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		startTime, err = time.Parse("15:04:05", request.StartTime)
		if err != nil {
			c.Log.WithError(err).Error("Failed to parse start_time")
			return nil, fiber.ErrBadRequest
		}
	}

	endTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		endTime, err = time.Parse("15:04:05", request.EndTime)
		if err != nil {
			c.Log.WithError(err).Error("Failed to parse end_time")
			return nil, fiber.ErrBadRequest
		}
	}

	shift := &entity.Shift{
		ID:            uuid.NewString(),
		CompanyID:     request.CompanyID,
		Name:          request.Name,
		StartTime:     startTime,
		EndTime:       endTime,
		LateTolerance: request.LateTolerance,
	}

	if err := c.ShiftRepository.Create(tx, shift); err != nil {
		c.Log.WithError(err).Error("Failed to create shift")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return &model.ShiftResponse{
		ID:            shift.ID,
		CompanyID:     shift.CompanyID,
		Name:          shift.Name,
		StartTime:     shift.StartTime,
		EndTime:       shift.EndTime,
		LateTolerance: shift.LateTolerance,
		CreatedAt:     shift.CreatedAt.String(),
		UpdatedAt:     shift.UpdatedAt.String(),
	}, nil
}

func (c *ShiftUseCase) AssignEmployee(ctx context.Context, request *model.AssignEmployeeToShiftRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return fiber.ErrBadRequest
	}

	var shiftTotal int64
	shiftTotal, err := c.ShiftRepository.CountByIDAndCompanyID(tx, request.ShiftID, request.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to check shift existence")
		return fiber.ErrInternalServerError
	}
	if shiftTotal == 0 {
		c.Log.Error("Shift not found")
		return fiber.ErrBadRequest
	}

	employeeTotal, err := c.ShiftRepository.CountEmployeeByIDAndCompanyID(tx, request.EmployeeID, request.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to check employee existence")
		return fiber.ErrInternalServerError
	}
	if employeeTotal == 0 {
		c.Log.Error("Employee not found")
		return fiber.ErrBadRequest
	}

	if err := c.ShiftRepository.DeleteEmployeeShiftsByEmployeeID(tx, request.EmployeeID); err != nil {
		c.Log.WithError(err).Error("Failed to remove previous employee shift")
		return fiber.ErrInternalServerError
	}

	if err := c.ShiftRepository.AssignEmployeeToShift(tx, request.EmployeeID, request.ShiftID); err != nil {
		c.Log.WithError(err).Error("Failed to assign employee to shift")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

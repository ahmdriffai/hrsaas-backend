package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TimeOffBalanceUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	TimeOffBalanceRepo *repository.TimeOffBalanceRepository
	TimeOffTypeRepo    *repository.TimeOffTypeRepository
}

func NewTimeOffBalanceUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	timeOffBalanceRepo *repository.TimeOffBalanceRepository,
	timeOffTypeRepo *repository.TimeOffTypeRepository,
) *TimeOffBalanceUseCase {
	return &TimeOffBalanceUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		TimeOffBalanceRepo: timeOffBalanceRepo,
		TimeOffTypeRepo:    timeOffTypeRepo,
	}
}

// TODO: Enforce company scoping if balances are shared across tenants.
func (c *TimeOffBalanceUseCase) ListBalances(ctx context.Context, employeeID string, request *model.SearchTimeOffBalanceRequest) ([]model.TimeOffBalanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, fiber.ErrBadRequest
	}

	items, err := c.TimeOffBalanceRepo.ListByEmployee(tx, employeeID, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off balances")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffBalanceResponse, len(items))
	for i, item := range items {
		responses[i] = *model.TimeOffBalanceToResponse(&item)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return responses, nil
}

// TODO: Add audit log for manual balance overrides.
func (c *TimeOffBalanceUseCase) SetBalance(ctx context.Context, request *model.SetTimeOffBalanceRequest) (*model.TimeOffBalanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	if _, err := c.TimeOffTypeRepo.FindByID(tx, request.TimeOffTypeID); err != nil {
		c.Log.WithError(err).Error("Time off type not found")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Time off type not found")
	}

	remainingDays := request.RemainingDays
	if remainingDays == nil {
		calculated := request.EntitledDays - request.UsedDays
		if calculated < 0 {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Remaining days cannot be negative")
		}
		remainingDays = &calculated
	}

	item, err := c.TimeOffBalanceRepo.FindByEmployeeTypeYear(tx, request.EmployeeID, request.TimeOffTypeID, request.PeriodYear)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.WithError(err).Error("Failed to check existing balance")
		return nil, fiber.ErrInternalServerError
	}

	if item == nil || err == gorm.ErrRecordNotFound {
		item = &entity.TimeOffBalance{
			EmployeeId:    request.EmployeeID,
			TimeOffTypeId: request.TimeOffTypeID,
			PeriodYear:    request.PeriodYear,
			EntitledDays:  request.EntitledDays,
			UsedDays:      request.UsedDays,
			RemainingDays: *remainingDays,
		}
		if err := c.TimeOffBalanceRepo.Create(tx, item); err != nil {
			c.Log.WithError(err).Error("Failed to create time off balance")
			return nil, fiber.ErrInternalServerError
		}
	} else {
		item.EntitledDays = request.EntitledDays
		item.UsedDays = request.UsedDays
		item.RemainingDays = *remainingDays
		if err := c.TimeOffBalanceRepo.Update(tx, item); err != nil {
			c.Log.WithError(err).Error("Failed to update time off balance")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.TimeOffBalanceToResponse(item), nil
}

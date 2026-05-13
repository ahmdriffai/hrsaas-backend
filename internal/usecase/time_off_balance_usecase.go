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
			EmployeeID:    request.EmployeeID,
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

func (c *TimeOffBalanceUseCase) Detail(ctx context.Context, id string) (*model.TimeOffBalanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item, err := c.TimeOffBalanceRepo.FindByID(tx, id)
	if err != nil {
		c.Log.WithError(err).Error("Time off balance not found")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.TimeOffBalanceToResponse(item), nil
}

func (c *TimeOffBalanceUseCase) Update(ctx context.Context, id string, request *model.UpdateTimeOffBalanceRequest) (*model.TimeOffBalanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	item, err := c.TimeOffBalanceRepo.FindByID(tx, id)
	if err != nil {
		c.Log.WithError(err).Error("Time off balance not found")
		return nil, fiber.ErrNotFound
	}

	if request.PeriodYear != nil {
		item.PeriodYear = *request.PeriodYear
	}
	if request.EntitledDays != nil {
		item.EntitledDays = *request.EntitledDays
	}
	if request.UsedDays != nil {
		item.UsedDays = *request.UsedDays
	}

	if request.RemainingDays != nil {
		item.RemainingDays = *request.RemainingDays
	} else {
		item.RemainingDays = item.EntitledDays - item.UsedDays
	}

	if item.RemainingDays < 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Remaining days cannot be negative")
	}

	if err := c.TimeOffBalanceRepo.Update(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to update time off balance")
		return nil, fiber.ErrInternalServerError
	}

	reloaded, err := c.TimeOffBalanceRepo.FindByID(tx, id)
	if err != nil {
		c.Log.WithError(err).Error("Failed to reload time off balance")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.TimeOffBalanceToResponse(reloaded), nil
}

func (c *TimeOffBalanceUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item, err := c.TimeOffBalanceRepo.FindByID(tx, id)
	if err != nil {
		c.Log.WithError(err).Error("Time off balance not found")
		return fiber.ErrNotFound
	}

	if err := c.TimeOffBalanceRepo.Delete(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to delete time off balance")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

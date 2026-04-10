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

type TimeOffTypeUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	TimeOffTypeRepo *repository.TimeOffTypeRepository
}

func NewTimeOffTypeUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	timeOffTypeRepo *repository.TimeOffTypeRepository,
) *TimeOffTypeUseCase {
	return &TimeOffTypeUseCase{
		DB:              db,
		Log:             log,
		Validate:        validate,
		TimeOffTypeRepo: timeOffTypeRepo,
	}
}

// TODO: Add caching if types rarely change.
func (c *TimeOffTypeUseCase) ListTypes(ctx context.Context) ([]model.TimeOffTypeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	items, err := c.TimeOffTypeRepo.List(tx)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off types")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffTypeResponse, len(items))
	for i, item := range items {
		responses[i] = *model.TimeOffTypeToResponse(&item)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return responses, nil
}

// TODO: Add Create Type use case and restrict to admin users.
func (c *TimeOffTypeUseCase) CreateType(ctx context.Context, request *model.CreateTimeOffTypeRequest) (*model.TimeOffTypeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	item := &entity.TimeOffType{
		Name:             request.Name,
		Category:         request.Category,
		IsQuotaBased:     request.IsQuotaBased,
		DefaultQuotaDays: request.DefaultQuotaDays,
	}

	if err := c.TimeOffTypeRepo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create time off type")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.TimeOffTypeToResponse(item), nil
}

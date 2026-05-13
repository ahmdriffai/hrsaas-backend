package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strings"

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
func (c *TimeOffTypeUseCase) ListTypes(ctx context.Context, request *model.SearchTimeOffTypeRequest) ([]model.TimeOffTypeResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.TimeOffTypeRepo.List(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off types")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffTypeResponse, len(items))
	for i, item := range items {
		responses[i] = *model.TimeOffTypeToResponse(&item)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}

// TODO: Add Create Type use case and restrict to admin users.
func (c *TimeOffTypeUseCase) Create(ctx context.Context, request *model.CreateTimeOffTypeRequest) (*model.TimeOffTypeResponse, error) {
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

func (c *TimeOffTypeUseCase) Detail(ctx context.Context, id string) (*model.TimeOffTypeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item, err := c.TimeOffTypeRepo.FindByID(tx, id)
	if err != nil {
		c.Log.WithError(err).Error("Failed to find time off type by ID")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.TimeOffTypeToResponse(item), nil
}

func (c *TimeOffTypeUseCase) UpdateType(ctx context.Context, id string, request *model.UpdateTimeOffTypeRequest) (*model.TimeOffTypeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	item, err := c.TimeOffTypeRepo.FindByID(tx, id)
	if err != nil {
		c.Log.WithError(err).Error("Failed to find time off type by ID")
		return nil, fiber.ErrNotFound
	}

	if request.Name != nil {
		name := strings.TrimSpace(*request.Name)
		if name == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		item.Name = name
	}
	if request.Category != nil {
		item.Category = *request.Category
	}
	if request.IsQuotaBased != nil {
		item.IsQuotaBased = *request.IsQuotaBased
	}
	if request.DefaultQuotaDays != nil {
		item.DefaultQuotaDays = *request.DefaultQuotaDays
	}

	if err := c.TimeOffTypeRepo.Update(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to update time off type")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.TimeOffTypeToResponse(item), nil
}

func (c *TimeOffTypeUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if _, err := c.TimeOffTypeRepo.FindByID(tx, id); err != nil {
		c.Log.WithError(err).Error("Failed to find time off type by ID")
		return fiber.ErrNotFound
	}

	if err := c.TimeOffTypeRepo.DeleteByID(tx, id); err != nil {
		c.Log.WithError(err).Error("Failed to delete time off type")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

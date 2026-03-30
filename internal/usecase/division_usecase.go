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

type DivisionUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	DivisionRepository *repository.DivisionRepository
}

func NewDivisionUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, divisionRepository *repository.DivisionRepository) *DivisionUseCase {
	return &DivisionUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		DivisionRepository: divisionRepository,
	}
}

// TODO: Add uniqueness check for division name per company.
func (c *DivisionUseCase) Create(ctx context.Context, request *model.CreateDivisionRequest) (*model.DivisionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	division := &entity.Division{
		CompanyID:   request.CompanyID,
		Name:        request.Name,
		Description: request.Description,
	}

	if err := c.DivisionRepository.Create(tx, division); err != nil {
		c.Log.WithError(err).Error("Failed to create division")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return &model.DivisionResponse{
		ID:          division.ID,
		CompanyID:   division.CompanyID,
		Name:        division.Name,
		Description: division.Description,
	}, nil
}

func (c *DivisionUseCase) Search(ctx context.Context, request *model.SearchDivisionRequest) ([]model.DivisionResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.DivisionRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error searching divisions")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.DivisionResponse, len(items))
	for i, item := range items {
		responses[i] = model.DivisionResponse{
			ID:          item.ID,
			CompanyID:   item.CompanyID,
			Name:        item.Name,
			Description: item.Description,
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}

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

type SanctionUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	SanctionRepository *repository.SanctionRepository
}

func NewSantionUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, sanctionRepository *repository.SanctionRepository) *SanctionUseCase {
	return &SanctionUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		SanctionRepository: sanctionRepository,
	}
}

/* Create Sanction Usecase
 */

func (c *SanctionUseCase) Create(ctx context.Context, request *model.CreateSanctionRequest) (*model.SanctionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	sanction := &entity.Sanction{
		CompanyID:   request.CompanyID,
		Name:        request.Name,
		Description: request.Description,
		Level:       request.Level,
		Note:        request.Note,
	}

	err := c.SanctionRepository.Create(tx, sanction)
	if err != nil {
		c.Log.Error("Error creating sanction:", err)
		return nil, err
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	response := &model.SanctionResponse{
		ID:          sanction.ID,
		CompanyID:   sanction.CompanyID,
		Name:        sanction.Name,
		Description: sanction.Description,
		Level:       sanction.Level,
		Note:        sanction.Note,
		CreatedAt:   sanction.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   sanction.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

/*
Get All Sanction
*/

func (c *SanctionUseCase) Search(ctx context.Context, request *model.SearchSanctionRequest) ([]model.SanctionResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	sanctions, total, err := c.SanctionRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting sanctions")
		return nil, 0, fiber.ErrInternalServerError
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.SanctionResponse, len(sanctions))
	for i, sanction := range sanctions {
		responses[i] = *converter.SanctionToResponse(&sanction)
	}

	return responses, total, nil
}

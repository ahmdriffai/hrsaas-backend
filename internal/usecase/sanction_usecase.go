package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strings"
	"time"

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

	return model.SanctionToResponse(sanction), nil
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
		responses[i] = *model.SanctionToResponse(&sanction)
	}

	return responses, total, nil
}

func (c *SanctionUseCase) Detail(ctx context.Context, id string, companyID string) (*model.SanctionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	sanction := new(entity.Sanction)
	if err := c.SanctionRepository.FindByIdAndCompany(tx, sanction, id, companyID); err != nil {
		c.Log.WithError(err).Error("Sanction not found")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.SanctionToResponse(sanction), nil
}

func (c *SanctionUseCase) Update(ctx context.Context, id string, companyID string, request *model.UpdateSanctionRequest) (*model.SanctionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	sanction := new(entity.Sanction)
	if err := c.SanctionRepository.FindByIdAndCompany(tx, sanction, id, companyID); err != nil {
		c.Log.WithError(err).Error("Sanction not found")
		return nil, fiber.ErrNotFound
	}

	if request.Name != nil {
		name := strings.TrimSpace(*request.Name)
		if name == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		sanction.Name = name
	}

	if request.Description != nil {
		sanction.Description = request.Description
	}

	if request.Level != nil {
		sanction.Level = *request.Level
	}

	if request.Note != nil {
		sanction.Note = request.Note
	}

	sanction.UpdatedAt = time.Now().UnixMilli()

	if err := c.SanctionRepository.Update(tx, sanction); err != nil {
		c.Log.WithError(err).Error("Failed to update sanction")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.SanctionToResponse(sanction), nil
}

func (c *SanctionUseCase) Delete(ctx context.Context, id string, companyID string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	sanction := new(entity.Sanction)
	if err := c.SanctionRepository.FindByIdAndCompany(tx, sanction, id, companyID); err != nil {
		c.Log.WithError(err).Error("Sanction not found")
		return fiber.ErrNotFound
	}

	if err := c.SanctionRepository.Delete(tx, sanction); err != nil {
		c.Log.WithError(err).Error("Failed to delete sanction")
		if strings.Contains(strings.ToLower(err.Error()), "foreign key") || strings.Contains(err.Error(), "SQLSTATE 23503") {
			return fiber.NewError(fiber.StatusConflict, "Sanction cannot be deleted because it is still used by employee sanctions")
		}
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

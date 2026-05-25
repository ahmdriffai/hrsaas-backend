package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HolidayUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	HolidayRepository *repository.HolidayRepository
}

func NewHolidayUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, holidayRepository *repository.HolidayRepository) *HolidayUseCase {
	return &HolidayUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		HolidayRepository: holidayRepository,
	}
}

func (c *HolidayUseCase) Create(ctx context.Context, request *model.CreateHolidayRequest) (*model.HolidayResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	if request.EndDate == 0 {
		request.EndDate = request.Date
	}

	holiday := &entity.Holiday{
		CompanyID:  request.CompanyID,
		Name:       request.Name,
		IsNational: request.IsNational,
		Date:       request.Date,
		EndDate:    request.EndDate,
	}

	if err := c.HolidayRepository.Create(tx, holiday); err != nil {
		c.Log.WithError(err).Error("Failed to create holiday")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.HolidayToResponse(holiday), nil
}

func (c *HolidayUseCase) Search(ctx context.Context, request *model.SearchHolidayRequest) ([]model.HolidayResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, 0, fiber.ErrBadRequest
	}
	holiday, total, err := c.HolidayRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to search holiday")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.HolidayResponse, len(holiday))
	for _, h := range holiday {
		responses = append(responses, *model.HolidayToResponse(&h))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}

func (c *HolidayUseCase) Update(ctx context.Context, request *model.UpdateHolidayRequest) (*model.HolidayResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	count, err := c.HolidayRepository.CountByIDAndCompanyID(tx, request.ID, request.CompanyID)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if count == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Holiday not found")
	}

	holiday := new(entity.Holiday)
	if err := c.HolidayRepository.FindById(tx, holiday, request.ID); err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if request.Name != nil {
		holiday.Name = *request.Name
	}

	if request.IsNational != nil {
		holiday.IsNational = *request.IsNational
	}

	if request.Date != nil {
		holiday.Date = *request.Date
	}

	if request.EndDate != nil {
		holiday.EndDate = *request.EndDate
	}

	holiday.UpdatedAt = time.Now().UnixMilli()

	if err := c.HolidayRepository.Update(tx, holiday); err != nil {
		c.Log.WithError(err).Error("Failed to update holiday")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.HolidayToResponse(holiday), nil
}

func (c *HolidayUseCase) Delete(ctx context.Context, id, companyID string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	holiday := new(entity.Holiday)
	if err := c.HolidayRepository.FindById(tx, holiday, id); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Holiday not found")
	}

	if err := c.HolidayRepository.Delete(tx, holiday); err != nil {
		c.Log.WithError(err).Error("Failed to delete holiday")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

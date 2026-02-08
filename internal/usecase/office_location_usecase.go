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

type OfficeLocationUseCase struct {
	DB                       *gorm.DB
	Log                      *logrus.Logger
	Validate                 *validator.Validate
	OfficeLocationRepository *repository.OfficeLocationRepository
}

func NewOfficeLocationUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	officeLocationRepository *repository.OfficeLocationRepository,
) *OfficeLocationUseCase {
	return &OfficeLocationUseCase{
		DB:                       db,
		Log:                      log,
		Validate:                 validate,
		OfficeLocationRepository: officeLocationRepository,
	}
}

/*Create Office Location Usecase */
func (c *OfficeLocationUseCase) Create(ctx context.Context, request *model.CreateOfficeLocationRequest) (*model.OfficeLocationResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	officeLocation := &entity.OfficeLocation{
		Name:      request.Name,
		Address:   request.Address,
		Lat:       request.Lat,
		Lng:       request.Lng,
		Radius:    request.Radius,
		IsActive:  true,
		CompanyID: request.CompanyID,
	}

	if err := c.OfficeLocationRepository.Create(tx, officeLocation); err != nil {
		c.Log.WithError(err).Error("Failed to create office location")
		return nil, err
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.OfficeLocationToResponse(officeLocation), nil
}


/* Search Office Location
 */

func (c *OfficeLocationUseCase) Search(ctx context.Context, request *model.SearchOfficeLocationRequest) ([]model.OfficeLocationResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	officeLocations, total, err := c.OfficeLocationRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting office locations")
		return nil, 0, fiber.ErrInternalServerError
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.OfficeLocationResponse, len(officeLocations))
	for i, officeLocation := range officeLocations {
		responses[i] = *model.OfficeLocationToResponse(&officeLocation)
	}

	for i, officeLocation := range officeLocations {
		responses[i] = *model.OfficeLocationToResponse(&officeLocation)
	}

	return responses, total, nil
}

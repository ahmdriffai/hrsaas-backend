package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strconv"
	"strings"
	"time"

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
		Lat:       strconv.FormatFloat(request.Lat, 'f', -1, 64),
		Lng:       strconv.FormatFloat(request.Lng, 'f', -1, 64),
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

func (c *OfficeLocationUseCase) AssignEmployee(ctx context.Context, request *model.AssignEmployeeToOfficeLocationRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return fiber.ErrBadRequest
	}

	officeLocationTotal, err := c.OfficeLocationRepository.CountByIDAndCompanyID(tx, request.OfficeLocationID, request.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to check office location existence")
		return fiber.ErrInternalServerError
	}
	if officeLocationTotal == 0 {
		c.Log.Error("Office location not found")
		return fiber.ErrBadRequest
	}

	employeeTotal, err := c.OfficeLocationRepository.CountEmployeeByIDAndCompanyID(tx, request.EmployeeID, request.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to check employee existence")
		return fiber.ErrInternalServerError
	}
	if employeeTotal == 0 {
		c.Log.Error("Employee not found")
		return fiber.ErrBadRequest
	}

	if err := c.OfficeLocationRepository.DeleteEmployeeOfficeLocationsByEmployeeID(tx, request.EmployeeID); err != nil {
		c.Log.WithError(err).Error("Failed to remove previous employee office location")
		return fiber.ErrInternalServerError
	}

	if err := c.OfficeLocationRepository.AssignEmployeeToOfficeLocation(tx, request.EmployeeID, request.OfficeLocationID); err != nil {
		c.Log.WithError(err).Error("Failed to assign employee to office location")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *OfficeLocationUseCase) DetailOfficeLocation(ctx context.Context, request *model.DetailOfficeLocationRequest) (*model.OfficeLocationResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// cek user exist
	count, err := c.OfficeLocationRepository.CountByIDAndCompanyID(tx, request.OfficeLocationID, request.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to count user by email")
		return nil, fiber.ErrInternalServerError
	}

	if count == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Ofi not found.")
	}

	var officeLocation = new(entity.OfficeLocation)

	err = c.OfficeLocationRepository.FindByIdAndCompany(tx, officeLocation, request.OfficeLocationID, request.CompanyID, "Employees")
	if err != nil {
		c.Log.WithError(err).Error("Failed to find ")
		return nil, fiber.ErrInternalServerError
	}

	return model.OfficeLocationToResponse(officeLocation), nil
}

func (c *OfficeLocationUseCase) Update(ctx context.Context, requestID string, companyID string, request *model.UpdateOfficeLocationRequest) (*model.OfficeLocationResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	officeLocation := new(entity.OfficeLocation)
	if err := c.OfficeLocationRepository.FindByIdAndCompany(tx, officeLocation, requestID, companyID); err != nil {
		c.Log.WithError(err).Error("Office location not found")
		return nil, fiber.ErrNotFound
	}

	if request.Name != nil {
		name := strings.TrimSpace(*request.Name)
		if name == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		officeLocation.Name = name
	}
	if request.Address != nil {
		officeLocation.Address = strings.TrimSpace(*request.Address)
	}
	if request.Lat != nil {
		officeLocation.Lat = strconv.FormatFloat(*request.Lat, 'f', -1, 64)
	}
	if request.Lng != nil {
		officeLocation.Lng = strconv.FormatFloat(*request.Lng, 'f', -1, 64)
	}
	if request.Radius != nil {
		officeLocation.Radius = *request.Radius
	}
	if request.IsActive != nil {
		officeLocation.IsActive = *request.IsActive
	}
	officeLocation.UpdatedAt = time.Now().UnixMilli()

	if err := c.OfficeLocationRepository.Update(tx, officeLocation); err != nil {
		c.Log.WithError(err).Error("Failed to update office location")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.OfficeLocationToResponse(officeLocation), nil
}

func (c *OfficeLocationUseCase) Delete(ctx context.Context, requestID string, companyID string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	officeLocation := new(entity.OfficeLocation)
	if err := c.OfficeLocationRepository.FindByIdAndCompany(tx, officeLocation, requestID, companyID); err != nil {
		c.Log.WithError(err).Error("Office location not found")
		return fiber.ErrNotFound
	}

	if err := c.OfficeLocationRepository.Delete(tx, officeLocation); err != nil {
		c.Log.WithError(err).Error("Failed to delete office location")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

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

type EmployeeIdentityUseCase struct {
	DB                         *gorm.DB
	Log                        *logrus.Logger
	Validate                   *validator.Validate
	EmployeeIdentityRepository *repository.EmployeeIdentityRepository
	EmployeeRepository         *repository.EmployeeRepository
}

func NewEmployeeIdentityUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, employeeIdentityRepository *repository.EmployeeIdentityRepository, employeeRepository *repository.EmployeeRepository) *EmployeeIdentityUseCase {
	return &EmployeeIdentityUseCase{
		DB:                         db,
		Log:                        log,
		Validate:                   validate,
		EmployeeIdentityRepository: employeeIdentityRepository,
		EmployeeRepository:         employeeRepository,
	}
}

func (c *EmployeeIdentityUseCase) Create(ctx context.Context, request *model.CreateEmployeeIdentityRequest) (*model.EmployeeIdentityResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	item := &entity.EmployeeIdentity{
		EmployeeID:                 request.EmployeeID,
		IdentityType:               request.IdentityType,
		IdentityNumber:             request.IdentityNumber,
		Address:                    request.Address,
		City:                       request.City,
		PostalCode:                 request.PostalCode,
		DomicililyAddress:          request.DomicililyAddress,
		IsDomicililySameAsIdentity: request.IsDomicililySameAsIdentity,
		IsDefault:                  request.IsDefault,
	}

	if err := c.EmployeeIdentityRepository.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create employee identity")
		return nil, fiber.ErrInternalServerError
	}

	if err := c.EmployeeIdentityRepository.FindById(tx, item, item.ID, "Employee"); err != nil {
		c.Log.WithError(err).Error("Failed to load employee identity")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeIdentityToResponse(item), nil
}

func (c *EmployeeIdentityUseCase) Update(ctx context.Context, id string, request *model.UpdateEmployeeIdentityRequest) (*model.EmployeeIdentityResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item := new(entity.EmployeeIdentity)
	if err := c.EmployeeIdentityRepository.FindById(tx, item, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.ErrNotFound
		}
		c.Log.WithError(err).Error("Failed to find employee identity")
		return nil, fiber.ErrInternalServerError
	}

	if request.IdentityType != nil {
		item.IdentityType = *request.IdentityType
	}
	if request.IdentityNumber != nil {
		item.IdentityNumber = *request.IdentityNumber
	}
	if request.Address != nil {
		item.Address = *request.Address
	}
	if request.City != nil {
		item.City = *request.City
	}
	if request.PostalCode != nil {
		item.PostalCode = *request.PostalCode
	}
	if request.DomicililyAddress != nil {
		item.DomicililyAddress = *request.DomicililyAddress
	}
	if request.IsDomicililySameAsIdentity != nil {
		item.IsDomicililySameAsIdentity = *request.IsDomicililySameAsIdentity
	}
	if request.IsDefault != nil {
		item.IsDefault = *request.IsDefault
	}

	if err := c.EmployeeIdentityRepository.Update(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to update employee identity")
		return nil, fiber.ErrInternalServerError
	}

	if err := c.EmployeeIdentityRepository.FindById(tx, item, item.ID, "Employee"); err != nil {
		c.Log.WithError(err).Error("Failed to load employee identity")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeIdentityToResponse(item), nil
}

func (c *EmployeeIdentityUseCase) List(ctx context.Context, request *model.SearchEmployeeIdentityRequest) ([]model.EmployeeIdentityResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.EmployeeIdentityRepository.List(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list employee identities")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return model.EmployeeIdentitiesToResponse(items), total, nil
}

func (c *EmployeeIdentityUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var item entity.EmployeeIdentity
	if err := c.EmployeeIdentityRepository.FindById(tx, &item, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.ErrNotFound
		}
		c.Log.WithError(err).Error("Failed to find employee identity")
		return fiber.ErrInternalServerError
	}

	if err := c.EmployeeIdentityRepository.Delete(tx, &item); err != nil {
		c.Log.WithError(err).Error("Failed to delete employee identity")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

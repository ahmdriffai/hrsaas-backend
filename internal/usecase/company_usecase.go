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

type CompanyUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	CompanyRepository *repository.CompanyRepository
	UserRepository    *repository.UserRepository
}

func NewCompanyUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	companyRepository *repository.CompanyRepository,
	userRepository *repository.UserRepository,
) *CompanyUseCase {
	return &CompanyUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		CompanyRepository: companyRepository,
		UserRepository:    userRepository,
	}
}

/*
Create Company
*/
func (c *CompanyUseCase) Create(ctx context.Context, request *model.CreateCompanyRequest) (*model.CompanyResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
 
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	company := &entity.Company{
		Name:           request.Name,
		LogoUrl:        &request.LogoUrl,
		BussinessField: &request.BussinessField,
		Address:        &request.Address,
		Province:       &request.Province,
		City:           &request.City,
		District:       &request.District,
		Village:        &request.Village,
		ZipCode:        &request.ZipCode,
		PhoneNumber:    &request.PhoneNumber,
		FaxNumber:      &request.FaxNumber,
		Email:          &request.Email,
		Website:        &request.Website,
	}
	if err := c.CompanyRepository.Create(tx, company); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CompanyToResponse(company), nil
}

func (c *CompanyUseCase) Register(ctx context.Context, request *model.RegisterCompanyRequest) (*model.CompanyResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	err := c.UserRepository.FindById(tx, user, request.UserID)
	if err != nil {
		c.Log.WithError(err).Error("User not found")
		return nil, fiber.ErrNotFound
	}

	// create company
	company := &entity.Company{
		Name:           request.Name,
		LogoUrl:        request.LogoUrl,
		BussinessField: request.BussinessField,
		Address:        request.Address,
		Province:       request.Province,
		City:           request.City,
		District:       request.District,
		Village:        request.Village,
		ZipCode:        request.ZipCode,
		PhoneNumber:    request.PhoneNumber,
		FaxNumber:      request.FaxNumber,
		Email:          request.Email,
		Website:        request.Website,
	}
	if err := c.CompanyRepository.Create(tx, company); err != nil {
		c.Log.WithError(err).Error("Failed to create company")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CompanyToResponse(company), nil
}

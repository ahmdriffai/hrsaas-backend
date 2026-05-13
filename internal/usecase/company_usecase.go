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

func (c *CompanyUseCase) List(ctx context.Context) ([]model.CompanyResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var companies []entity.Company
	if err := tx.Find(&companies).Error; err != nil {
		c.Log.WithError(err).Error("Failed to list companies")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.CompanyResponse, len(companies))
	for i, company := range companies {
		responses[i] = *converter.CompanyToResponse(&company)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return responses, nil
}

func (c *CompanyUseCase) Update(ctx context.Context, companyID string, request *model.UpdateCompanyRequest) (*model.CompanyResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	var company entity.Company
	if err := c.CompanyRepository.FindById(tx, &company, companyID); err != nil {
		c.Log.WithError(err).Error("Company not found")
		return nil, fiber.ErrNotFound
	}

	if request.Name != nil {
		company.Name = *request.Name
	}
	if request.LogoUrl != nil {
		company.LogoUrl = request.LogoUrl
	}
	if request.BussinessField != nil {
		company.BussinessField = request.BussinessField
	}
	if request.Address != nil {
		company.Address = request.Address
	}
	if request.Province != nil {
		company.Province = request.Province
	}
	if request.City != nil {
		company.City = request.City
	}
	if request.District != nil {
		company.District = request.District
	}
	if request.Village != nil {
		company.Village = request.Village
	}
	if request.ZipCode != nil {
		company.ZipCode = request.ZipCode
	}
	if request.PhoneNumber != nil {
		company.PhoneNumber = request.PhoneNumber
	}
	if request.FaxNumber != nil {
		company.FaxNumber = request.FaxNumber
	}
	if request.Email != nil {
		company.Email = request.Email
	}
	if request.Website != nil {
		company.Website = request.Website
	}

	if err := c.CompanyRepository.Update(tx, &company); err != nil {
		c.Log.WithError(err).Error("Failed to update company")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CompanyToResponse(&company), nil
}

func (c *CompanyUseCase) Detail(ctx context.Context, companyID string) (*model.CompanyResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var company entity.Company
	if err := c.CompanyRepository.FindById(tx, &company, companyID); err != nil {
		c.Log.WithError(err).Error("Company not found")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CompanyToResponse(&company), nil
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

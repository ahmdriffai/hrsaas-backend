package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeEducationUseCase struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     *repository.EmployeeEducationRepository
}

func NewEmployeeEducationUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	repo *repository.EmployeeEducationRepository,
) *EmployeeEducationUseCase {
	return &EmployeeEducationUseCase{DB: db, Log: log, Validate: validate, Repo: repo}
}

func (c *EmployeeEducationUseCase) Create(ctx context.Context, request *model.CreateEmployeeEducationRequest) (*model.EmployeeEducationResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	graduationYear, err := lib.ParseDateToUnixMilli(request.GraduationYear)
	if err != nil {
		c.Log.WithError(err).Error("Failed to parse graduation year")
		return nil, fiber.ErrBadRequest
	}

	item := &entity.EmployeeEducation{
		CompanyID:       request.CompanyID,
		EmployeeID:      request.EmployeeID,
		EducationLevel:  request.EducationLevel,
		InstitutionName: request.InstitutionName,
		Major:           request.Major,
		GraduationYear:  graduationYear,
		GPA:             request.GPA,
		StartYear:       request.StartYear,
		EndYear:         request.EndYear,
	}

	if err := c.Repo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create employee education")
		return nil, fiber.ErrInternalServerError
	}

	if err := c.Repo.FindById(tx, item, item.ID, "Employee"); err != nil {
		c.Log.WithError(err).Error("Failed to load employee education")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeEducationToResponse(item), nil
}

func (c *EmployeeEducationUseCase) List(ctx context.Context, request *model.SearchEmployeeEducationRequest) ([]model.EmployeeEducationResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.Repo.List(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list employee educations")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return model.EmployeeEducationsToResponse(items), total, nil
}

func (c *EmployeeEducationUseCase) Detail(ctx context.Context, id string) (*model.EmployeeEducationResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var item entity.EmployeeEducation
	if err := c.Repo.FindById(tx, &item, id, "Employee"); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.ErrNotFound
		}
		c.Log.WithError(err).Error("Failed to find employee education")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeEducationToResponse(&item), nil
}

func (c *EmployeeEducationUseCase) Update(ctx context.Context, request *model.UpdateEmployeeEducationRequest) (*model.EmployeeEducationResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	var item entity.EmployeeEducation
	if err := c.Repo.FindById(tx, &item, request.ID, "Employee"); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.ErrNotFound
		}
		c.Log.WithError(err).Error("Failed to find employee education")
		return nil, fiber.ErrInternalServerError
	}

	if request.EducationLevel != nil {
		item.EducationLevel = *request.EducationLevel
	}
	if request.InstitutionName != nil {
		item.InstitutionName = *request.InstitutionName
	}
	if request.Major != nil {
		item.Major = *request.Major
	}
	if request.GraduationYear != nil {
		t, err := lib.ParseDateToUnixMilli(*request.GraduationYear)
		if err != nil {
			c.Log.WithError(err).Error("Failed to parse graduation year")
			return nil, fiber.ErrBadRequest
		}
		item.GraduationYear = t
	}
	if request.GPA != nil {
		item.GPA = request.GPA
	}
	if request.StartYear != nil {
		item.StartYear = request.StartYear
	}
	if request.EndYear != nil {
		item.EndYear = request.EndYear
	}

	if err := c.Repo.Update(tx, &item); err != nil {
		c.Log.WithError(err).Error("Failed to update employee education")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeEducationToResponse(&item), nil
}

func (c *EmployeeEducationUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var item entity.EmployeeEducation
	if err := c.Repo.FindById(tx, &item, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.ErrNotFound
		}
		c.Log.WithError(err).Error("Failed to find employee education")
		return fiber.ErrInternalServerError
	}

	if err := c.Repo.Delete(tx, &item); err != nil {
		c.Log.WithError(err).Error("Failed to delete employee education")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

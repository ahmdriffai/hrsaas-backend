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

type EmployeeTrainingUseCase struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     *repository.EmployeeTrainingRepository
}

func NewEmployeeTrainingUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	repo *repository.EmployeeTrainingRepository,
) *EmployeeTrainingUseCase {
	return &EmployeeTrainingUseCase{DB: db, Log: log, Validate: validate, Repo: repo}
}

func (c *EmployeeTrainingUseCase) Create(ctx context.Context, request *model.CreateEmployeeTrainingRequest) (*model.EmployeeTrainingResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	startDate, err := lib.ParseDateToUnixMilli(request.StartDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Format start_date tidak valid, gunakan YYYY-MM-DD")
	}

	item := &entity.EmployeeTraining{
		CompanyID:      request.CompanyID,
		EmployeeID:     request.EmployeeID,
		TrainingName:   request.TrainingName,
		Organizer:      request.Organizer,
		StartDate:      startDate,
		CertificateURL: request.CertificateURL,
	}

	if request.EndDate != nil {
		endDate, err := lib.ParseDateToUnixMilli(*request.EndDate)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Format end_date tidak valid, gunakan YYYY-MM-DD")
		}
		item.EndDate = &endDate
	}

	if err := c.Repo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create employee training")
		return nil, fiber.ErrInternalServerError
	}

	if err := c.Repo.FindById(tx, item, item.ID, "Employee"); err != nil {
		c.Log.WithError(err).Error("Failed to load employee training")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeTrainingToResponse(item), nil
}

func (c *EmployeeTrainingUseCase) List(ctx context.Context, request *model.SearchEmployeeTrainingRequest) ([]model.EmployeeTrainingResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.Repo.List(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list employee trainings")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return model.EmployeeTrainingsToResponse(items), total, nil
}

func (c *EmployeeTrainingUseCase) Detail(ctx context.Context, id string) (*model.EmployeeTrainingResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var item entity.EmployeeTraining
	if err := c.Repo.FindById(tx, &item, id, "Employee"); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.ErrNotFound
		}
		c.Log.WithError(err).Error("Failed to find employee training")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeTrainingToResponse(&item), nil
}

func (c *EmployeeTrainingUseCase) Update(ctx context.Context, request *model.UpdateEmployeeTrainingRequest) (*model.EmployeeTrainingResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	var item entity.EmployeeTraining
	if err := c.Repo.FindById(tx, &item, request.ID, "Employee"); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.ErrNotFound
		}
		c.Log.WithError(err).Error("Failed to find employee training")
		return nil, fiber.ErrInternalServerError
	}

	if request.TrainingName != nil {
		item.TrainingName = *request.TrainingName
	}
	if request.Organizer != nil {
		item.Organizer = *request.Organizer
	}
	if request.StartDate != nil {
		startDate, err := lib.ParseDateToUnixMilli(*request.StartDate)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Format start_date tidak valid, gunakan YYYY-MM-DD")
		}
		item.StartDate = startDate
	}
	if request.EndDate != nil {
		endDate, err := lib.ParseDateToUnixMilli(*request.EndDate)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Format end_date tidak valid, gunakan YYYY-MM-DD")
		}
		item.EndDate = &endDate
	}
	if request.CertificateURL != nil {
		item.CertificateURL = request.CertificateURL
	}

	if err := c.Repo.Update(tx, &item); err != nil {
		c.Log.WithError(err).Error("Failed to update employee training")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeTrainingToResponse(&item), nil
}

func (c *EmployeeTrainingUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var item entity.EmployeeTraining
	if err := c.Repo.FindById(tx, &item, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.ErrNotFound
		}
		c.Log.WithError(err).Error("Failed to find employee training")
		return fiber.ErrInternalServerError
	}

	if err := c.Repo.Delete(tx, &item); err != nil {
		c.Log.WithError(err).Error("Failed to delete employee training")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

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

type EmployeeDocumentUseCase struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     *repository.EmployeeDocumentRepository
}

func NewEmployeeDocumentUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	repo *repository.EmployeeDocumentRepository,
) *EmployeeDocumentUseCase {
	return &EmployeeDocumentUseCase{DB: db, Log: log, Validate: validate, Repo: repo}
}

func (c *EmployeeDocumentUseCase) Create(ctx context.Context, request *model.CreateEmployeeDocumentRequest) (*model.EmployeeDocumentResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	item := &entity.EmployeeDocument{
		EmployeeID: request.EmployeeID,
		DocType:    request.DocType,
		DocNumber:  request.DocNumber,
		FileURL:    request.FileURL,
	}

	if err := c.Repo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create employee document")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeDocumentToResponse(item), nil
}

func (c *EmployeeDocumentUseCase) List(ctx context.Context, request *model.SearchEmployeeDocumentRequest) ([]model.EmployeeDocumentResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.Repo.List(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list employee documents")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return model.EmployeeDocumentsToResponse(items), total, nil
}

func (c *EmployeeDocumentUseCase) Update(ctx context.Context, request *model.UpdateEmployeeDocumentRequest) (*model.EmployeeDocumentResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	var item entity.EmployeeDocument
	if err := c.Repo.FindById(tx, &item, request.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.ErrNotFound
		}
		c.Log.WithError(err).Error("Failed to find employee document")
		return nil, fiber.ErrInternalServerError
	}

	if request.DocType != nil {
		item.DocType = *request.DocType
	}
	if request.DocNumber != nil {
		item.DocNumber = *request.DocNumber
	}
	if request.FileURL != nil {
		item.FileURL = *request.FileURL
	}

	if err := c.Repo.Update(tx, &item); err != nil {
		c.Log.WithError(err).Error("Failed to update employee document")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeDocumentToResponse(&item), nil
}

func (c *EmployeeDocumentUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var item entity.EmployeeDocument
	if err := c.Repo.FindById(tx, &item, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.ErrNotFound
		}
		c.Log.WithError(err).Error("Failed to find employee document")
		return fiber.ErrInternalServerError
	}

	if err := c.Repo.Delete(tx, &item); err != nil {
		c.Log.WithError(err).Error("Failed to delete employee document")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

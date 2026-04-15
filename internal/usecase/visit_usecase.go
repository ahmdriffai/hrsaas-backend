package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VisitUseCase struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     *repository.VisitRepository
}

func NewVisitUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, repo *repository.VisitRepository) *VisitUseCase {
	return &VisitUseCase{
		DB:       db,
		Log:      log,
		Validate: validate,
		Repo:     repo,
	}
}

// TODO: Validate visit_type, open visit rules, and required fields per business rule.
func (c *VisitUseCase) Create(ctx context.Context, employeeID, companyID string, request *model.CreateVisitRequest) (*model.VisitResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	visitType := strings.ToUpper(strings.TrimSpace(request.VisitType))
	if visitType != "IN" && visitType != "OUT" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "visit_type must be IN or OUT")
	}

	lastVisit, err := c.Repo.FindLastByEmployee(tx, employeeID, false)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.WithError(err).Error("Failed to check last visit")
		return nil, fiber.ErrInternalServerError
	}

	if visitType == "IN" && lastVisit != nil && lastVisit.VisitType == "IN" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tidak dapat memulai kunjungan. Lakukan selesai kunjungan terlebih dahulu")
	}
	if visitType == "OUT" {
		if lastVisit == nil || lastVisit.VisitType != "IN" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Tidak dapat menyelesaikan kunjungan. Anda harus memulai kunjungan terlebih dahulu.")
		}
	}

	item := &entity.Visit{
		EmployeeID: employeeID,
		CompanyID:  companyID,
		VisitType:  visitType,
		Note:       request.Note,
		Latitude:   request.Latitude,
		Longitude:  request.Longitude,
		Address:    request.Address,
		FileURL:    request.FileUrl,
		CreatedAt:  nowEpoch(),
	}

	if err := c.Repo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create visit")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Preload("Employee").Preload("Company").First(item, "id = ?", item.ID).Error; err != nil {
		c.Log.WithError(err).Error("Failed to load visit relations")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.VisitToResponse(item), nil
}

func (c *VisitUseCase) List(ctx context.Context, request *model.SearchVisitRequest) ([]model.VisitResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}
	if request.SortBy != "" && request.SortBy != "newest" && request.SortBy != "oldest" {
		return nil, 0, fiber.NewError(fiber.StatusBadRequest, "sort_by must be newest or oldest")
	}
	if request.StartDate != "" {
		if _, err := lib.ParseDateToUnixMilli(request.StartDate); err != nil {
			return nil, 0, fiber.NewError(fiber.StatusBadRequest, "Invalid start_date")
		}
	}
	if request.EndDate != "" {
		if _, err := lib.ParseDateToUnixMilli(request.EndDate); err != nil {
			return nil, 0, fiber.NewError(fiber.StatusBadRequest, "Invalid end_date")
		}
	}

	items, total, err := c.Repo.List(tx, request, true)
	if err != nil {
		c.Log.WithError(err).Error("failed to list visits")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.VisitResponse, len(items))
	for i := range items {
		responses[i] = *model.VisitToResponse(&items[i])
	}

	return responses, total, nil
}

func (c *VisitUseCase) GetByID(ctx context.Context, id string) (*model.VisitResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item, err := c.Repo.FindByID(tx, id, true)
	if err != nil {
		c.Log.WithError(err).Error("visit not found")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.VisitToResponse(item), nil
}

func (c *VisitUseCase) GetVisitOwner(ctx context.Context, id string) (string, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item, err := c.Repo.FindByID(tx, id, false)
	if err != nil {
		c.Log.WithError(err).Error("visit not found")
		return "", fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return "", fiber.ErrInternalServerError
	}

	return item.EmployeeID, nil
}

func (c *VisitUseCase) CanDoVisit(ctx context.Context, employeeID, visitType string) (*model.CanDoVisitResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	visitType = strings.ToUpper(strings.TrimSpace(visitType))
	if visitType != "IN" && visitType != "OUT" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "visit_type must be IN or OUT")
	}

	lastVisit, err := c.Repo.FindLastByEmployee(tx, employeeID, false)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.WithError(err).Error("Failed to check last visit")
		return nil, fiber.ErrInternalServerError
	}

	response := &model.CanDoVisitResponse{
		CanDoVisit: true,
		Message:    "You can do visit " + visitType,
	}

	if visitType == "IN" {
		if lastVisit != nil && lastVisit.VisitType == "IN" {
			response.CanDoVisit = false
			response.Message = "Anda tidak dapat melakukan memulai kunjungan sebelum menyelesaikan kunjungan sebelumnya"
		}
	} else {
		if lastVisit == nil || lastVisit.VisitType != "IN" {
			response.CanDoVisit = false
			response.Message = "Anda tidak dapat menyelesaikan kunjungan sebelum memulai kunjungan baru"
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return response, nil
}

func (c *VisitUseCase) GetUnclosedVisit(ctx context.Context, employeeID string) (*model.VisitResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	lastVisit, err := c.Repo.FindLastByEmployee(tx, employeeID, true)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.WithError(err).Error("Failed to check last visit")
		return nil, fiber.ErrInternalServerError
	}

	if lastVisit == nil || lastVisit.VisitType != "IN" {
		if err := tx.Commit().Error; err != nil {
			c.Log.WithError(err).Error("Failed to commit transaction")
			return nil, fiber.ErrInternalServerError
		}
		return nil, nil
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.VisitToResponse(lastVisit), nil
}

// TODO: Consider soft delete if audits are required.
func (c *VisitUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item, err := c.Repo.FindByID(tx, id, false)
	if err != nil {
		c.Log.WithError(err).Error("visit not found")
		return fiber.ErrNotFound
	}

	if err := c.Repo.Delete(tx, item); err != nil {
		c.Log.WithError(err).Error("failed to delete visit")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

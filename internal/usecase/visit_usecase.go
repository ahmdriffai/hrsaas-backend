package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strings"
	"time"

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

func (c *VisitUseCase) Create(ctx context.Context, employeeID, companyID string, request *model.CreateVisitRequest) (*model.VisitResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	visitType := strings.ToUpper(strings.TrimSpace(request.VisitType))

	now := time.Now()
	nowMilli := now.UnixMilli()
	todayStr := now.Format("2006-01-02")
	timeStr := now.Format("15:04:05")

	var location *string
	if request.Latitude != nil && request.Longitude != nil {
		loc := *request.Latitude + ", " + *request.Longitude
		location = &loc
	}

	var visitID string

	if visitType == "IN" {
		if request.ClientName == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "client_name is required for IN visit")
		}

		_, err := c.Repo.FindLastOpenByEmployee(tx, employeeID)
		if err == nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Tidak dapat memulai kunjungan. Lakukan selesai kunjungan terlebih dahulu")
		}
		if err != gorm.ErrRecordNotFound {
			c.Log.WithError(err).Error("Failed to check open visit")
			return nil, fiber.ErrInternalServerError
		}

		visit := &entity.Visit{
			EmployeeID: employeeID,
			CompanyID:  companyID,
			Date:       todayStr,
			ClientName: request.ClientName,
		}
		if err := c.Repo.Create(tx, visit); err != nil {
			c.Log.WithError(err).Error("Failed to create visit")
			return nil, fiber.ErrInternalServerError
		}
		visitID = visit.ID

	} else {
		openVisit, err := c.Repo.FindLastOpenByEmployee(tx, employeeID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fiber.NewError(fiber.StatusBadRequest, "Tidak dapat menyelesaikan kunjungan. Anda harus memulai kunjungan terlebih dahulu.")
			}
			c.Log.WithError(err).Error("Failed to find open visit")
			return nil, fiber.ErrInternalServerError
		}
		visitID = openVisit.ID

		if err := tx.Model(&entity.Visit{}).Where("id = ?", visitID).Update("updated_at", nowMilli).Error; err != nil {
			c.Log.WithError(err).Error("Failed to update visit updated_at")
			return nil, fiber.ErrInternalServerError
		}
	}

	detail := &entity.VisitDetail{
		VisitID:   visitID,
		VisitType: visitType,
		VisitAt:   timeStr,
		DateVisit: todayStr,
		FileUrl:   request.FileUrl,
		Location:  location,
		Address:   request.Address,
		Note:      request.Note,
	}
	if err := tx.Create(detail).Error; err != nil {
		c.Log.WithError(err).Error("Failed to create visit detail")
		return nil, fiber.ErrInternalServerError
	}

	var result entity.Visit
	if err := tx.Preload("Employee").Preload("Details").Where("id = ?", visitID).Take(&result).Error; err != nil {
		c.Log.WithError(err).Error("Failed to load visit with relations")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.VisitToResponse(&result), nil
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
		if _, err := time.Parse("2006-01-02", request.StartDate); err != nil {
			return nil, 0, fiber.NewError(fiber.StatusBadRequest, "Invalid start_date")
		}
	}
	if request.EndDate != "" {
		if _, err := time.Parse("2006-01-02", request.EndDate); err != nil {
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

func (c *VisitUseCase) Update(ctx context.Context, id string, request *model.UpdateVisitRequest) (*model.VisitResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	if _, err := c.Repo.FindByID(tx, id, true); err != nil {
		c.Log.WithError(err).Error("visit not found")
		return nil, fiber.ErrNotFound
	}

	if request.Latitude != nil || request.Longitude != nil {
		if request.Latitude == nil || request.Longitude == nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "latitude and longitude must be provided together")
		}
	}

	nowMilli := time.Now().UnixMilli()

	visitUpdates := map[string]any{
		"updated_at": nowMilli,
	}
	if request.ClientName != nil {
		clientName := strings.TrimSpace(*request.ClientName)
		if clientName == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "client_name cannot be empty")
		}
		visitUpdates["client_name"] = clientName
	}

	if len(visitUpdates) > 1 {
		if err := tx.Model(&entity.Visit{}).Where("id = ?", id).Updates(visitUpdates).Error; err != nil {
			c.Log.WithError(err).Error("Failed to update visit")
			return nil, fiber.ErrInternalServerError
		}
	}

	detail, err := c.Repo.FindLatestDetailByVisitID(tx, id)
	if err != nil {
		c.Log.WithError(err).Error("visit detail not found")
		return nil, fiber.ErrInternalServerError
	}

	detailUpdates := map[string]any{
		"updated_at": nowMilli,
	}
	if request.FileUrl != nil {
		detailUpdates["file_url"] = request.FileUrl
	}
	if request.Address != nil {
		detailUpdates["address"] = request.Address
	}
	if request.Note != nil {
		detailUpdates["note"] = request.Note
	}
	if request.Latitude != nil && request.Longitude != nil {
		location := strings.TrimSpace(*request.Latitude) + ", " + strings.TrimSpace(*request.Longitude)
		detailUpdates["location"] = location
	}

	if len(detailUpdates) > 1 {
		if err := tx.Model(&entity.VisitDetail{}).Where("id = ?", detail.ID).Updates(detailUpdates).Error; err != nil {
			c.Log.WithError(err).Error("Failed to update visit detail")
			return nil, fiber.ErrInternalServerError
		}
	}

	result, err := c.Repo.FindByID(tx, id, true)
	if err != nil {
		c.Log.WithError(err).Error("Failed to load updated visit")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.VisitToResponse(result), nil
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

	_, err := c.Repo.FindLastOpenByEmployee(tx, employeeID)
	hasOpenVisit := err == nil

	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.WithError(err).Error("Failed to check open visit")
		return nil, fiber.ErrInternalServerError
	}

	response := &model.CanDoVisitResponse{CanDoVisit: true, Message: "You can do visit " + visitType}

	if visitType == "IN" && hasOpenVisit {
		response.CanDoVisit = false
		response.Message = "Anda tidak dapat melakukan memulai kunjungan sebelum menyelesaikan kunjungan sebelumnya"
	} else if visitType == "OUT" && !hasOpenVisit {
		response.CanDoVisit = false
		response.Message = "Anda tidak dapat menyelesaikan kunjungan sebelum memulai kunjungan baru"
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

	openVisit, err := c.Repo.FindLastOpenByEmployee(tx, employeeID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.WithError(err).Error("Failed to check open visit")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return model.VisitToResponse(openVisit), nil
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

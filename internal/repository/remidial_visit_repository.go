package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RemidialVisitRepository struct {
	Repository[entity.RemidialVisit]
	Log        *logrus.Logger
	httpClient *http.Client
	baseURL    string
}

func NewRemidialVisitRepository(log *logrus.Logger, baseURL string) *RemidialVisitRepository {
	return &RemidialVisitRepository{
		Log:        log,
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

func (r *RemidialVisitRepository) SearchNasabah(ctx context.Context, req *model.SearchNasabahRequest) ([]model.NasabahData, error) {
	data, err := json.Marshal(req)
	if err != nil {
		r.Log.WithError(err).Error("Failed to marshal request")
		return nil, fiber.ErrInternalServerError
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/bwakses/apirem", bytes.NewBuffer(data))
	if err != nil {
		r.Log.WithError(err).Error("Failed to create request")
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		r.Log.WithError(err).Error("Failed to call external API")
		return nil, fiber.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r.Log.Errorf("External API returned status: %d", resp.StatusCode)
		return nil, fiber.ErrInternalServerError
	}

	var result model.ExternalNasabahResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		r.Log.WithError(err).Error("Failed to decode external API response")
		return nil, fiber.ErrInternalServerError
	}

	if result.RC != "00" {
		r.Log.Errorf("External API returned RC: %s", result.RC)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal mendapatkan data dari server core")
	}

	return result.Data, nil
}

func (r *RemidialVisitRepository) List(db *gorm.DB, request *model.SearchRemidialVisitRequest) ([]entity.RemidialVisit, int64, error) {
	var items []entity.RemidialVisit

	query, err := r.FilterSearch(db.Model(&entity.RemidialVisit{}), request)
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("created_at DESC")
	if request.Page > 0 && request.Size > 0 {
		query = query.Offset((request.Page - 1) * request.Size).Limit(request.Size)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *RemidialVisitRepository) FilterSearch(db *gorm.DB, request *model.SearchRemidialVisitRequest) (*gorm.DB, error) {
	query := db.Joins("LEFT JOIN employees e ON e.id = remidial_visit.employee_id")
	if request.EmployeeID != "" {
		query = query.Where("remidial_visit.employee_id = ?", request.EmployeeID)
	}
	if request.EmployeeName != "" {
		query = query.Where("e.fullname ILIKE ?", "%"+request.EmployeeName+"%")
	}
	if request.NasabahName != "" {
		query = query.Where("remidial_visit.nasabah_name ILIKE ?", "%"+request.NasabahName+"%")
	}
	if request.StartDate != "" {
		query = query.Where("remidial_visit.created_at >= ?", request.StartDate)
	}
	if request.EndDate != "" {
		query = query.Where("remidial_visit.created_at <= ?", request.EndDate)
	}
	if request.NoPjm != "" {
		query = query.Where("remidial_visit.no_pjm ILIKE ?", "%"+request.NoPjm+"%")
	}
	return query, nil
}

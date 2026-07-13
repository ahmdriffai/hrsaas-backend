package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AnnouncementRepository struct {
	Repository[entity.Announcement]
	Log *logrus.Logger
}

func NewAnnouncementRepository(log *logrus.Logger) *AnnouncementRepository {
	return &AnnouncementRepository{Log: log}
}

func (r *AnnouncementRepository) List(db *gorm.DB, request *model.SearchAnnouncementRequest) ([]entity.Announcement, int64, error) {
	var items []entity.Announcement

	query := db.Model(&entity.Announcement{}).Preload("Employee")
	if request.Title != "" {
		query = query.Where("title LIKE ?", "%"+request.Title+"%")
	}
	if request.Category != "" {
		query = query.Where("category = ?", request.Category)
	}
	if request.CompanyID != "" {
		query = query.Where("company_id = ?", request.CompanyID)
	}
	if request.EmployeeID != "" {
		query = query.Where("employee_id = ?", request.EmployeeID)
	}
	if request.CreatedAt != 0 {
		query = query.Where("created_at >= ?", request.CreatedAt)
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

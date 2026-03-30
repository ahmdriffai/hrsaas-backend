package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VisitRepository struct {
	Repository[entity.Visit]
	Log *logrus.Logger
}

func NewVisitRepository(log *logrus.Logger) *VisitRepository {
	return &VisitRepository{Log: log}
}

func (r *VisitRepository) List(db *gorm.DB, request *model.SearchVisitRequest) ([]entity.Visit, int64, error) {
	var items []entity.Visit

	query := db.Model(&entity.Visit{})
	if request.EmployeeID != "" {
		query = query.Where("employee_id = ?", request.EmployeeID)
	}
	if request.VisitType != "" {
		query = query.Where("visit_type = ?", request.VisitType)
	}
	if request.StartDate != "" {
		startDate, _ := lib.ParseDateToUnixMilli(request.StartDate)
		query = query.Where("created_at >= ?", startDate)
	}
	if request.EndDate != "" {
		endDate, _ := lib.ParseDateToUnixMilli(request.EndDate)
		endDate = endDate + int64(24*time.Hour/time.Millisecond) - 1
		query = query.Where("created_at <= ?", endDate)
	}

	if request.SortBy == "oldest" {
		query = query.Order("created_at ASC")
	} else {
		query = query.Order("created_at DESC")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *VisitRepository) FindByID(db *gorm.DB, id string) (*entity.Visit, error) {
	var item entity.Visit
	if err := db.Where("id = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *VisitRepository) FindLastByEmployee(db *gorm.DB, employeeID string) (*entity.Visit, error) {
	var item entity.Visit
	if err := db.Where("employee_id = ?", employeeID).Order("created_at DESC").Limit(1).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

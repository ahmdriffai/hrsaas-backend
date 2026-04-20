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

func (r *VisitRepository) List(db *gorm.DB, request *model.SearchVisitRequest, withRelations bool) ([]entity.Visit, int64, error) {
	var items []entity.Visit

	countQuery, err := r.applyFilters(db.Model(&entity.Visit{}), request)
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	dataQuery, err := r.applyFilters(db.Model(&entity.Visit{}), request)
	if err != nil {
		return nil, 0, err
	}
	if withRelations {
		dataQuery = dataQuery.Preload("Employee").Preload("Company")
	}
	if request.SortBy == "oldest" {
		dataQuery = dataQuery.Order("created_at ASC")
	} else {
		dataQuery = dataQuery.Order("created_at DESC")
	}

	if err := dataQuery.Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *VisitRepository) FindByID(db *gorm.DB, id string, withRelations bool) (*entity.Visit, error) {
	var item entity.Visit
	query := db.Model(&entity.Visit{})
	if withRelations {
		query = query.Preload("Employee").Preload("Company")
	}
	if err := query.Where("id = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *VisitRepository) FindLastByEmployee(db *gorm.DB, employeeID string, withRelations bool) (*entity.Visit, error) {
	var item entity.Visit
	query := db.Model(&entity.Visit{})
	if withRelations {
		query = query.Preload("Employee").Preload("Company")
	}
	if err := query.Where("employee_id = ?", employeeID).Order("created_at DESC").Limit(1).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *VisitRepository) applyFilters(query *gorm.DB, request *model.SearchVisitRequest) (*gorm.DB, error) {
	if request == nil {
		return query, nil
	}
	if request.EmployeeID != "" {
		query = query.Where("employee_id = ?", request.EmployeeID)
	}
	if request.VisitType != "" {
		query = query.Where("visit_type = ?", request.VisitType)
	}
	if request.StartDate != "" {
		startDate, err := lib.ParseDateToUnixMilli(request.StartDate)
		if err != nil {
			return nil, err
		}
		query = query.Where("created_at >= ?", startDate)
	}
	if request.EndDate != "" {
		endDate, err := lib.ParseDateToUnixMilli(request.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = endDate + int64(24*time.Hour/time.Millisecond) - 1
		query = query.Where("created_at <= ?", endDate)
	}
	return query, nil
}

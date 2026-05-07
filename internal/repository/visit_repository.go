package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

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
		dataQuery = dataQuery.Preload("Employee").Preload("Details")
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
		query = query.Preload("Employee").Preload("Details")
	}
	if err := query.Where("id = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// FindLastOpenByEmployee returns the most recent visit that has an IN detail but no OUT detail.
func (r *VisitRepository) FindLastOpenByEmployee(db *gorm.DB, employeeID string) (*entity.Visit, error) {
	var item entity.Visit
	err := db.Model(&entity.Visit{}).
		Preload("Details").
		Where("employee_id = ?", employeeID).
		Where("EXISTS (SELECT 1 FROM visit_details WHERE visit_details.visit_id = visits.id AND visit_details.visit_type = 'IN')").
		Where("NOT EXISTS (SELECT 1 FROM visit_details WHERE visit_details.visit_id = visits.id AND visit_details.visit_type = 'OUT')").
		Order("created_at DESC").
		Limit(1).
		Take(&item).Error
	if err != nil {
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
	if request.StartDate != "" {
		query = query.Where("date >= ?", request.StartDate)
	}
	if request.EndDate != "" {
		query = query.Where("date <= ?", request.EndDate)
	}
	return query, nil
}

package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TimeOffRequestRepository struct {
	Repository[entity.TimeOffRequest]
	Log *logrus.Logger
}

func NewTimeOffRequestRepository(log *logrus.Logger) *TimeOffRequestRepository {
	return &TimeOffRequestRepository{Log: log}
}

func (r *TimeOffRequestRepository) List(db *gorm.DB, request *model.SearchTimeOffRequest, withRelations bool) ([]entity.TimeOffRequest, int64, error) {
	var items []entity.TimeOffRequest

	query := db.Model(&entity.TimeOffRequest{})
	query = r.applyFilters(query, request)
	if withRelations {
		query = query.
			Preload("Employee").
			Preload("TimeOffType").
			Preload("Approvals").
			Preload("Approvals.Employee").
			Preload("Approvals.TimeOffRequest").
			Preload("Employee.EmployeeContract").
			Preload("Employee.EmployeeContract.Position").
			Preload("Employee.EmployeeContract.Division")
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

func (r *TimeOffRequestRepository) FindByID(db *gorm.DB, id string, withRelations bool) (*entity.TimeOffRequest, error) {
	var item entity.TimeOffRequest
	query := db.Model(&entity.TimeOffRequest{})
	if withRelations {
		query = query.
			Preload("Employee").
			Preload("TimeOffType").
			Preload("Approvals").
			Preload("Approvals.Employee").
			Preload("Approvals.TimeOffRequest")
	}
	if err := query.Where("id = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *TimeOffRequestRepository) UpdateStatus(db *gorm.DB, id string, status string) error {
	return db.Model(&entity.TimeOffRequest{}).
		Where("id = ?", id).
		Update("request_status", status).Error
}

func (r *TimeOffRequestRepository) applyFilters(query *gorm.DB, request *model.SearchTimeOffRequest) *gorm.DB {
	if request == nil {
		return query
	}
	if request.EmployeeID != "" {
		query = query.Where("employee_id = ?", request.EmployeeID)
	}
	if request.TimeOffTypeID != "" {
		query = query.Where("time_off_type_id = ?", request.TimeOffTypeID)
	}
	if request.RequestStatus != "" {
		query = query.Where("request_status = ?", request.RequestStatus)
	}
	if request.StartDate != "" {
		startDate, _ := lib.ParseDateToUnixMilli(request.StartDate)
		query = query.Where("start_date >= ?", startDate)
	}
	if request.EndDate != "" {
		endDate, _ := lib.ParseDateToUnixMilli(request.EndDate)
		query = query.Where("end_date <= ?", endDate)
	}
	if request.TimeOffTypeID != "" {
		query = query.Where("time_off_type_id = ?", request.TimeOffTypeID)
	}

	return query
}

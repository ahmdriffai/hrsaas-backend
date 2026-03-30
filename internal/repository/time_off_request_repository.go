package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TimeOffRequestRepository struct {
	Repository[entity.Time_Off_Requests]
	Log *logrus.Logger
}

func NewTimeOffRequestRepository(log *logrus.Logger) *TimeOffRequestRepository {
	return &TimeOffRequestRepository{Log: log}
}

func (r *TimeOffRequestRepository) List(db *gorm.DB, request *model.SearchTimeOffRequest) ([]entity.Time_Off_Requests, int64, error) {
	var items []entity.Time_Off_Requests

	query := db.Model(&entity.Time_Off_Requests{})
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

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *TimeOffRequestRepository) FindByID(db *gorm.DB, id string) (*entity.Time_Off_Requests, error) {
	var item entity.Time_Off_Requests
	if err := db.Where("id = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

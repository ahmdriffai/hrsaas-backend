package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmSancRepository struct {
	Repository[entity.EmployeeSanction]
	Log *logrus.Logger
}

func NewEmSancRepository(log *logrus.Logger) *EmSancRepository {
	return &EmSancRepository{
		Log: log,
	}
}

func (r *EmSancRepository) Search(db *gorm.DB, request *model.SearchEmSancRequest) ([]entity.EmployeeSanction, int64, error) {
	var employeeSanc []entity.EmployeeSanction
	if err := db.Preload("Employee").Preload("Sanction").Scopes(r.FilterSearch(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&employeeSanc).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.EmployeeSanction{}).Scopes(r.FilterSearch(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return employeeSanc, total, nil
}

func (r *EmSancRepository) FilterSearch(request *model.SearchEmSancRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("company_id = ?", request.CompanyID)

		if request.EmployeeID != "" {
			tx = tx.Where("employee_id = ?", request.EmployeeID)
		}

		if reason := request.Reason; reason != "" {
			reason = "%" + reason + "%"
			tx = tx.Where("reason LIKE ?", reason)
		}

		if sanctionID := request.SanctionID; sanctionID != "" {
			tx = tx.Where("sanction_id = ?", sanctionID)
		}

		if status := request.Status; status != "" {
			tx = tx.Where("status = ?", status)
		}

		startDate, _ := lib.ParseDateToUnixMilli(request.StartDate)
		endDate, _ := lib.ParseDateToUnixMilli(request.EndDate)

		if request.StartDate != "" && request.EndDate != "" {
			tx = tx.Where(
				"start_date >= ? AND end_date <= ?",
				startDate,
				endDate,
			)
		} else if request.StartDate != "" {
			tx = tx.Where(
				"start_date >= ?",
				startDate,
			)
		} else if request.EndDate != "" {
			tx = tx.Where(
				"end_date <= ?",
				endDate,
			)
		}

		return tx
	}
}

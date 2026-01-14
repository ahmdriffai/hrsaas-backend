package repository

import (
	"hr-sas/internal/entity"
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
		tx = tx.Joins("JOIN employees e ON e.id = employee_sanctions.employee_id").
			Joins("JOIN sanctions s ON s.id = employee_sanctions.sanction_id").
			Where("e.company_id = ?", request.CompanyID)

		if request.UserID != "" {
			tx = tx.Where("e.user_id = ?", request.UserID)
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

		// ✅ start_date filter
		if request.StartDate != nil && request.EndDate != nil {
			tx = tx.Where(
				"start_date >= ? AND end_date <= ?",
				request.StartDate.Time,
				request.EndDate.Time,
			)
		} else if request.StartDate != nil {
			tx = tx.Where("start_date >= ?", request.StartDate.Time)
		} else if request.EndDate != nil {
			tx = tx.Where("end_date <= ?", request.EndDate.Time)
		}

		return tx
	}
}

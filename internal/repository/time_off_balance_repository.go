package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TimeOffBalanceRepository struct {
	Repository[entity.Time_Off_Balance]
	Log *logrus.Logger
}

func NewTimeOffBalanceRepository(log *logrus.Logger) *TimeOffBalanceRepository {
	return &TimeOffBalanceRepository{Log: log}
}

func (r *TimeOffBalanceRepository) ListByEmployee(db *gorm.DB, employeeID string, request *model.SearchTimeOffBalanceRequest) ([]entity.Time_Off_Balance, error) {
	var items []entity.Time_Off_Balance

	query := db.Model(&entity.Time_Off_Balance{}).Where("employee_id = ?", employeeID)
	if request.TimeOffTypeID != "" {
		query = query.Where("time_off_type_id = ?", request.TimeOffTypeID)
	}
	if request.PeriodYear > 0 {
		query = query.Where("period_year = ?", request.PeriodYear)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

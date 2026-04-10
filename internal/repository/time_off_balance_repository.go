package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TimeOffBalanceRepository struct {
	Repository[entity.TimeOffBalance]
	Log *logrus.Logger
}

func NewTimeOffBalanceRepository(log *logrus.Logger) *TimeOffBalanceRepository {
	return &TimeOffBalanceRepository{Log: log}
}

func (r *TimeOffBalanceRepository) ListByEmployee(db *gorm.DB, employeeID string, request *model.SearchTimeOffBalanceRequest) ([]entity.TimeOffBalance, error) {
	var items []entity.TimeOffBalance

	query := db.Model(&entity.TimeOffBalance{}).
		Preload("Employee").
		Preload("TimeOffType").
		Where("employee_id = ?", employeeID)
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

func (r *TimeOffBalanceRepository) FindByEmployeeTypeYear(db *gorm.DB, employeeID, timeOffTypeID string, periodYear int) (*entity.TimeOffBalance, error) {
	var item entity.TimeOffBalance
	if err := db.
		Where("employee_id = ?", employeeID).
		Where("time_off_type_id = ?", timeOffTypeID).
		Where("period_year = ?", periodYear).
		Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// TODO : create update balance function, which will be used to update balance after approval or cancellation of time off request. Need to check if balance is sufficient before approval, and update balance after approval or cancellation.
func (r *TimeOffBalanceRepository) UpdateBalance(db *gorm.DB, employeeID, timeOffTypeID string, periodYear int, delta int) error {
	var item entity.TimeOffBalance
	if err := db.
		Where("employee_id = ?", employeeID).
		Where("time_off_type_id = ?", timeOffTypeID).
		Where("period_year = ?", periodYear).
		Take(&item).Error; err != nil {
		return err
	}

	newBalance := item.RemainingDays + delta
	if newBalance < 0 {
		return gorm.ErrInvalidData
	}

	item.RemainingDays = newBalance
	if err := db.Save(&item).Error; err != nil {
		return err
	}

	return nil
}

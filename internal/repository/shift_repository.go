package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ShiftRepository struct {
	Repository[entity.Shift]
	Log *logrus.Logger
}

func NewShiftRepository(log *logrus.Logger) *ShiftRepository {
	return &ShiftRepository{
		Log: log,
	}
}

func (r *ShiftRepository) FindByEmployeeID(db *gorm.DB, employeeId string) ([]entity.Shift, error) {
	var shifts []entity.Shift
	if err := db.
		Model(&entity.Shift{}).
		Select(`
			shifts.id,
			shifts.company_id,
			shifts.name,
			shifts.late_tolerance,
			shifts.created_at,
			shifts.updated_at
		`).
		Joins("JOIN employee_shifts ON employee_shifts.shift_id = shifts.id").
		Where("employee_shifts.employee_id = ?", employeeId).
		Find(&shifts).Error; err != nil {
		return nil, err
	}
	return shifts, nil
}

func (r *ShiftRepository) CountByIDAndCompanyID(db *gorm.DB, shiftID, companyID string) (int64, error) {
	var total int64
	err := db.Model(&entity.Shift{}).
		Where("id = ?", shiftID).
		Where("company_id = ?", companyID).
		Count(&total).Error
	return total, err
}

func (r *ShiftRepository) CountEmployeeByIDAndCompanyID(db *gorm.DB, employeeID, companyID string) (int64, error) {
	var total int64
	err := db.Model(&entity.Employee{}).
		Where("id = ?", employeeID).
		Where("company_id = ?", companyID).
		Count(&total).Error
	return total, err
}

func (r *ShiftRepository) DeleteEmployeeShiftsByEmployeeID(db *gorm.DB, employeeID string) error {
	return db.Exec("DELETE FROM employee_shifts WHERE employee_id = ?", employeeID).Error
}

func (r *ShiftRepository) AssignEmployeeToShift(db *gorm.DB, employeeID, shiftID string) error {
	return db.Exec(
		"INSERT INTO employee_shifts (employee_id, shift_id) VALUES (?, ?)",
		employeeID,
		shiftID,
	).Error
}

func (r *ShiftRepository) Search(db *gorm.DB, request *model.SearchShiftRequest) ([]entity.Shift, int64, error) {
	var shifts []entity.Shift
	if err := db.Scopes(r.FilterSearch(request)).
		Select(`
			shifts.id,
			shifts.company_id,
			shifts.name,
			shifts.late_tolerance,
			shifts.created_at,
			shifts.updated_at
		`).
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&shifts).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&entity.Shift{}).Scopes(r.FilterSearch(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return shifts, total, nil
}

func (r *ShiftRepository) FilterSearch(request *model.SearchShiftRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("company_id = ?", request.CompanyID)
		if request.Key != "" {
			tx = tx.Where("name LIKE ?", "%"+request.Key+"%")
		}
		return tx
	}
}

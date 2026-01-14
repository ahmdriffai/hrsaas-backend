package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeRepository struct {
	Repository[entity.Employee]
	Log *logrus.Logger
}

func NewEmployeeRepository(log *logrus.Logger) *EmployeeRepository {
	return &EmployeeRepository{
		Log: log,
	}
}

func (r *EmployeeRepository) CountByEmployeeNumberAndCompanyID(db *gorm.DB, employeeNumber string, CompanyID string) (int64, error) {
	var total int64
	err := db.Model(new(entity.Employee)).Where("employee_number = ?", employeeNumber).Where("company_id = ?", CompanyID).Count(&total).Error
	return total, err
}

func (r *EmployeeRepository) Search(db *gorm.DB, request *model.SearchEmployeeRequest) ([]entity.Employee, int64, error) {
	var employee []entity.Employee
	if err := db.Scopes(r.FilterSearch(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&employee).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Employee{}).Scopes(r.FilterSearch(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return employee, total, nil
}

func (r *EmployeeRepository) FilterSearch(request *model.SearchEmployeeRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("company_id = ?", request.CompanyID)

		if key := request.Key; key != "" {
			key = "%" + key + "%"
			tx = tx.Where("fullname LIKE ?", key).Or("employee_number LIKE ?", key)
		}

		return tx
	}
}

package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeEducationRepository struct {
	Repository[entity.EmployeeEducation]
	Log *logrus.Logger
}

func NewEmployeeEducationRepository(log *logrus.Logger) *EmployeeEducationRepository {
	return &EmployeeEducationRepository{Log: log}
}

func (r *EmployeeEducationRepository) List(db *gorm.DB, request *model.SearchEmployeeEducationRequest) ([]entity.EmployeeEducation, int64, error) {
	var items []entity.EmployeeEducation

	query := db.Model(&entity.EmployeeEducation{}).Preload("Employee").Where("company_id = ?", request.CompanyID)
	if request.EmployeeID != "" {
		query = query.Where("employee_id = ?", request.EmployeeID)
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

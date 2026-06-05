package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeIdentityRepository struct {
	Repository[entity.EmployeeIdentity]
	Log *logrus.Logger
}

func NewEmployeeIdentityRepository(log *logrus.Logger) *EmployeeIdentityRepository {
	return &EmployeeIdentityRepository{Log: log}
}

func (r *EmployeeIdentityRepository) List(db *gorm.DB, request *model.SearchEmployeeIdentityRequest) ([]entity.EmployeeIdentity, int64, error) {
	var items []entity.EmployeeIdentity

	query := db.Model(&entity.EmployeeIdentity{}).Where("employee_id = ?", request.EmployeeID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

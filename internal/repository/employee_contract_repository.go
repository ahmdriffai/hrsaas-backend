package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeContractRepository struct {
	Repository[entity.EmployeeContract]
	Log *logrus.Logger
}

func NewEmployeeContractRepository(log *logrus.Logger) *EmployeeContractRepository {
	return &EmployeeContractRepository{Log: log}
}

func (r *EmployeeContractRepository) List(db *gorm.DB, request *model.SearchEmployeeContractRequest) ([]entity.EmployeeContract, int64, error) {
	var items []entity.EmployeeContract

	query := db.Model(&entity.EmployeeContract{})
	if request.EmployeeID != "" {
		query = query.Where("employee_id = ?", request.EmployeeID)
	}
	if request.DivisionID != "" {
		query = query.Where("division_id = ?", request.DivisionID)
	}
	if request.PositionID != "" {
		query = query.Where("position_id = ?", request.PositionID)
	}
	if request.ActiveOnly {
		now := time.Now().UnixMilli()
		query = query.Where("end_date IS NULL OR end_date >= ?", now)
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

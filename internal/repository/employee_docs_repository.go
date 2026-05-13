package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeDocumentRepository struct {
	Repository[entity.EmployeeDocument]
	Log *logrus.Logger
}

func NewEmployeeDocumentRepository(log *logrus.Logger) *EmployeeDocumentRepository {
	return &EmployeeDocumentRepository{Log: log}
}

func (r *EmployeeDocumentRepository) List(db *gorm.DB, request *model.SearchEmployeeDocumentRequest) ([]entity.EmployeeDocument, int64, error) {
	var items []entity.EmployeeDocument

	query := db.Model(&entity.EmployeeDocument{}).Preload("Employee")
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

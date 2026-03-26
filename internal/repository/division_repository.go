package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DivisionRepository struct {
	Repository[entity.Division]
	Log *logrus.Logger
}

func NewDivisionRepository(log *logrus.Logger) *DivisionRepository {
	return &DivisionRepository{Log: log}
}

func (r *DivisionRepository) Search(db *gorm.DB, request *model.SearchDivisionRequest) ([]entity.Division, int64, error) {
	var items []entity.Division

	if err := db.Scopes(r.FilterSearch(request)).
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&entity.Division{}).Scopes(r.FilterSearch(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *DivisionRepository) FilterSearch(request *model.SearchDivisionRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("company_id = ?", request.CompanyID)
		if key := request.Name; key != "" {
			key = "%" + key + "%"
			tx = tx.Where("name LIKE ?", key).Or("description LIKE ?", key)
		}
		return tx
	}
}

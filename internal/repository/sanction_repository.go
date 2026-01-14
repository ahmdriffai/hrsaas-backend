package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SanctionRepository struct {
	Repository[entity.Sanction]
	Log *logrus.Logger
}

func NewSanctionRepository(log *logrus.Logger) *SanctionRepository {
	return &SanctionRepository{
		Log: log,
	}
}

func (r *SanctionRepository) Search(db *gorm.DB, request *model.SearchSanctionRequest) ([]entity.Sanction, int64, error) {
	var sanctions []entity.Sanction
	if err := db.Scopes(r.FilterSearch(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&sanctions).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Sanction{}).Scopes(r.FilterSearch(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return sanctions, total, nil
}

func (r *SanctionRepository) FilterSearch(request *model.SearchSanctionRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("company_id = ?", request.CompanyID)

		if key := request.Key; key != "" {
			key = "%" + key + "%"
			tx = tx.Where("name LIKE ?", key).Or("description LIKE ?", key)
		}

		return tx
	}
}

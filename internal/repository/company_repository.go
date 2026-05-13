package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	Repository[entity.Company]
	Log *logrus.Logger
}

func NewCompanyRepository(log *logrus.Logger) *CompanyRepository {
	return &CompanyRepository{
		Log: log,
	}
}

func (r *CompanyRepository) Search(db *gorm.DB, request *model.SearchCompanyRequest) ([]entity.Company, int64, error) {
	var items []entity.Company
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

func (r *CompanyRepository) FilterSearch(request *model.SearchCompanyRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("company_id = ?", request.CompanyId)
		if key := request.Name; key != "" {
			key = "%" + key + "%"
			tx = tx.Where("name LIKE ?", key).Or("description LIKE ?", key)
		}
		return tx
	}
}

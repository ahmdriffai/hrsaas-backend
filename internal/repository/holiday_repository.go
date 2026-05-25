package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HolidayRepository struct {
	Repository[entity.Holiday]
	Log *logrus.Logger
}

func NewHolidayRepository(log *logrus.Logger) *HolidayRepository {
	return &HolidayRepository{
		Log: log,
	}
}

func (r *HolidayRepository) Search(db *gorm.DB, request *model.SearchHolidayRequest) ([]entity.Holiday, int64, error) {
	var holidays []entity.Holiday
	if err := db.Scopes(r.FilterSearch(request)).
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&holidays).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&entity.Holiday{}).Scopes(r.FilterSearch(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return holidays, total, nil
}

func (r *HolidayRepository) FilterSearch(request *model.SearchHolidayRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("company_id = ? OR is_national = true", request.CompanyID)
		if request.Key != "" {
			tx = tx.Where("name LIKE ?", "%"+request.Key+"%")
		}
		return tx
	}
}

func (r *HolidayRepository) CountByIDAndCompanyID(db *gorm.DB, id, companyID string) (int64, error) {
	var total int64
	err := db.Model(&entity.Holiday{}).
		Where("id = ?", id).
		Where("company_id = ? OR is_national = true", companyID).
		Count(&total).Error
	return total, err
}

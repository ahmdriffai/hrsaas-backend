package repository

import (
	"hr-sas/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TimeOffTypeRepository struct {
	Repository[entity.TimeOffType]
	Log *logrus.Logger
}

func NewTimeOffTypeRepository(log *logrus.Logger) *TimeOffTypeRepository {
	return &TimeOffTypeRepository{Log: log}
}

func (r *TimeOffTypeRepository) List(db *gorm.DB) ([]entity.TimeOffType, error) {
	var items []entity.TimeOffType
	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TimeOffTypeRepository) ListQuotaBased(db *gorm.DB) ([]entity.TimeOffType, error) {
	var items []entity.TimeOffType
	if err := db.Where("is_quota_based = ?", true).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TimeOffTypeRepository) FindByID(db *gorm.DB, id string) (*entity.TimeOffType, error) {
	var item entity.TimeOffType
	if err := db.Where("id = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

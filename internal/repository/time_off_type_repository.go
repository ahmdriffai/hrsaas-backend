package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

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

func (r *TimeOffTypeRepository) List(db *gorm.DB, request *model.SearchTimeOffTypeRequest) ([]entity.TimeOffType, int64, error) {
	var items []entity.TimeOffType

	query := db.Model(&entity.TimeOffType{})
	if request != nil && request.Name != "" {
		query = query.Where("name ILIKE ?", "%"+request.Name+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if request != nil {
		query = query.Offset((request.Page - 1) * request.Size).Limit(request.Size)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
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

func (r *TimeOffTypeRepository) DeleteByID(db *gorm.DB, id string) error {
	return db.Where("id = ?", id).Delete(&entity.TimeOffType{}).Error
}

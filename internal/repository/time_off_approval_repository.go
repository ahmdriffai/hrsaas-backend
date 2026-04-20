package repository

import (
	"hr-sas/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TimeOffApprovalRepository struct {
	Repository[entity.TimeOffApproval]
	Log *logrus.Logger
}

func NewTimeOffApprovalRepository(log *logrus.Logger) *TimeOffApprovalRepository {
	return &TimeOffApprovalRepository{Log: log}
}

func (r *TimeOffApprovalRepository) ListByRequestID(db *gorm.DB, requestID string) ([]entity.TimeOffApproval, error) {
	var items []entity.TimeOffApproval
	if err := db.Where("time_off_request_id = ?", requestID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TimeOffApprovalRepository) FindByID(db *gorm.DB, id string) (*entity.TimeOffApproval, error) {
	var item entity.TimeOffApproval
	if err := db.Where("id = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *TimeOffApprovalRepository) CreateMany(db *gorm.DB, items []entity.TimeOffApproval) error {
	if len(items) == 0 {
		return nil
	}
	return db.Create(&items).Error
}

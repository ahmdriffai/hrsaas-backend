package repository

import (
	"hr-sas/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TimeOffAttachmentRepository struct {
	Repository[entity.Time_Off_Attachment]
	Log *logrus.Logger
}

func NewTimeOffAttachmentRepository(log *logrus.Logger) *TimeOffAttachmentRepository {
	return &TimeOffAttachmentRepository{Log: log}
}

func (r *TimeOffAttachmentRepository) ListByRequestID(db *gorm.DB, requestID string) ([]entity.Time_Off_Attachment, error) {
	var items []entity.Time_Off_Attachment
	if err := db.Where("time_off_request_id = ?", requestID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

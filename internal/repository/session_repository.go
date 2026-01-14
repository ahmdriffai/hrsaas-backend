package repository

import (
	"hr-sas/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SessionRepository struct {
	Repository[entity.Session]
	Log *logrus.Logger
}

func NewSessionRepository(log *logrus.Logger) *SessionRepository {
	return &SessionRepository{
		Log: log,
	}
}

func (c *SessionRepository) FindByToken(db *gorm.DB, session *entity.Session, token string) error {
	return db.Where("token = ?", token).First(session).Error
}

func (r *SessionRepository) CountByUserId(db *gorm.DB, userId string) (int64, error) {
	var total int64
	err := db.Model(new(entity.Session)).Where("user_id = ?", userId).Count(&total).Error
	return total, err
}

func (r *SessionRepository) DeleteByUserId(db *gorm.DB, userId string) error {
	return db.Where("user_id = ?", userId).Delete(new(entity.Session)).Error
}

package repository

import (
	"hr-sas/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AttendanceLogRepository struct {
	Repository[entity.AttendanceLog]
	Log *logrus.Logger
}

func NewAttendanceLogRepository(log *logrus.Logger) *AttendanceLogRepository {
	return &AttendanceLogRepository{
		Log: log,
	}
}

func (r *AttendanceLogRepository) CountByAttendanceIDAndType(db *gorm.DB, attendanceID string, logType string) (int64, error) {
	var count int64
	err := db.Model(&entity.AttendanceLog{}).
		Where("attendance_id = ? AND type = ?", attendanceID, logType).
		Count(&count).Error
	return count, err
}

func (r *AttendanceLogRepository) FindLastByAttendanceIDAndType(db *gorm.DB, entity *entity.AttendanceLog, attendanceID string, logType string) error {
	return db.
		Where("attendance_id = ? AND type = ?", attendanceID, logType).
		Order("time DESC").
		Take(entity).Error
}

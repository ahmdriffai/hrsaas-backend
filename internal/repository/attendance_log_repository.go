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

func (r *AttendanceLogRepository) FindByAttendanceID(db *gorm.DB, attendanceID string) ([]entity.AttendanceLog, error) {
	var logs []entity.AttendanceLog
	if err := db.
		Where("attendance_id = ?", attendanceID).
		Order("time ASC").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *AttendanceLogRepository) DeleteByAttendanceID(db *gorm.DB, attendanceID string) error {
	return db.Where("attendance_id = ?", attendanceID).Delete(&entity.AttendanceLog{}).Error
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

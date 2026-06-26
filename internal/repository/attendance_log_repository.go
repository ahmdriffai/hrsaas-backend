package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"time"

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

func (r *AttendanceLogRepository) Search(db *gorm.DB, request *model.SearchAttendanceLogRequest) ([]entity.AttendanceLog, int64, error) {
	var logs []entity.AttendanceLog
	if err := db.Scopes(r.FilterSearch(request)).
		Order("attendance_logs.time DESC").
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.AttendanceLog{}).Scopes(r.FilterSearch(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *AttendanceLogRepository) FilterSearch(request *model.SearchAttendanceLogRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Joins("JOIN attendances ON attendances.id = attendance_logs.attendance_id").
			Where("attendances.company_id = ?", request.CompanyID)

		if request.AttendanceID != "" {
			tx = tx.Where("attendance_logs.attendance_id = ?", request.AttendanceID)
		}

		if request.EmployeeID != "" {
			tx = tx.Where("attendances.employee_id = ?", request.EmployeeID)
		}

		if request.Type != "" {
			tx = tx.Where("attendance_logs.type = ?", request.Type)
		}

		if request.IsApproved != nil {
			tx = tx.Where("attendance_logs.is_approved = ?", *request.IsApproved)
		}

		if request.Date != "" {
			if t, err := time.Parse("2006-01-02", request.Date); err == nil {
				startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).UnixMilli()
				tx = tx.Where("attendances.date = ?", startOfDay)
			}
		}

		return tx
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

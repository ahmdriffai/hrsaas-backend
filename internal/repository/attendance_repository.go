package repository

import (
	"hr-sas/internal/entity"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AttendanceRepository struct {
	Repository[entity.Attendance]
	Log *logrus.Logger
}

func NewAttendanceRepository(log *logrus.Logger) *AttendanceRepository {
	return &AttendanceRepository{
		Log: log,
	}
}
func (r *AttendanceRepository) Update(db *gorm.DB, attendance *entity.Attendance) error {
	attendance.UpdatedAt = time.Now().UnixMilli()

	updates := map[string]any{
		"check_out_time":      attendance.CheckOutTime,
		"total_work_minutes":  attendance.TotalWorkMinutes,
		"total_break_minutes": attendance.TotalBreakMinutes,
		"status":              attendance.Status,
		"updated_at":          attendance.UpdatedAt,
	}

	return db.Table(attendance.TableName()).
		Where("id = ?", attendance.ID).
		Updates(updates).
		Error
}

func (r *AttendanceRepository) FindByEmployeeIDAndDate(db *gorm.DB, entity *entity.Attendance, employeeId string, date int64) error {
	// Pastikan jamnya 00:00:00
	t := time.UnixMilli(date)

	// tanggal jam 00:00:00
	startOfDay := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0,
		t.Location(),
	).UnixMilli()

	return db.
		Where("employee_id = ? AND date = ?", employeeId, startOfDay).
		Take(entity).
		Error
}

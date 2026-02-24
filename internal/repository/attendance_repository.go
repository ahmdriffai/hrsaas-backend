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
	attendance.UpdatedAt = time.Now()

	updates := map[string]any{
		"company_id":          attendance.CompanyID,
		"employee_id":         attendance.EmployeeID,
		"date":                attendance.Date,
		"check_in_time":       attendance.CheckInTime,
		"total_work_minutes":  gorm.Expr("make_interval(mins => ?)", attendance.TotalWorkMinutes),
		"total_break_minutes": gorm.Expr("make_interval(mins => ?)", attendance.TotalBreakMinutes),
		"status":              attendance.Status,
		"updated_at":          attendance.UpdatedAt,
	}

	if attendance.CheckOutTime.IsZero() {
		updates["check_out_time"] = nil
	} else {
		updates["check_out_time"] = attendance.CheckOutTime
	}

	return db.Table(attendance.TableName()).
		Where("id = ?", attendance.ID).
		Updates(updates).
		Error
}

func (r *AttendanceRepository) FindByEmployeeIDAndDate(db *gorm.DB, entity *entity.Attendance, employeeId string, date time.Time) error {
	// Pastikan jamnya 00:00:00
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	return db.
		Where("employee_id = ? AND date = ?", employeeId, dateOnly).
		Take(entity).
		Error
}

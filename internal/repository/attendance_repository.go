package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
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
		"date":                attendance.Date,
		"check_in_time":       attendance.CheckInTime,
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

func (r *AttendanceRepository) Search(db *gorm.DB, request *model.SearchAttendanceRequest) ([]entity.Attendance, int64, error) {
	var attendances []entity.Attendance
	if err := db.Scopes(r.FilterSearch(request)).
		Order("date DESC").
		Preload("Employee").
		Preload("Employee.EmployeeContract").
		Preload("Employee.EmployeeContract.Position").
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&attendances).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Attendance{}).Scopes(r.FilterSearch(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

func (r *AttendanceRepository) FilterSearch(request *model.SearchAttendanceRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("company_id = ?", request.CompanyID)

		if request.EmployeeID != "" {
			tx = tx.Where("employee_id = ?", request.EmployeeID)
		}

		if request.Status != "" {
			tx = tx.Where("status = ?", request.Status)
		}

		if request.IsApproved != nil {
			tx = tx.Where("is_approved = ?", *request.IsApproved)
		}

		if request.Date != "" {
			if t, err := time.Parse("2006-01-02", request.Date); err == nil {
				startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).UnixMilli()
				tx = tx.Where("date = ?", startOfDay)
			}
		}

		return tx
	}
}

func (r *AttendanceRepository) FindByEmployeeIDAndDate(db *gorm.DB, entity *entity.Attendance, employeeId string, date int64) error {
	t := time.UnixMilli(date)

	startOfDay := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0,
		t.Location(),
	).UnixMilli()

	return db.
		Preload("Employee").
		Preload("Employee.EmployeeContract").
		Preload("Employee.EmployeeContract.Position").
		Where("employee_id = ? AND date = ?", employeeId, startOfDay).
		Take(entity).
		Error
}

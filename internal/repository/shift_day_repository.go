package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ShiftDayRepository struct {
	Repository[entity.ShiftDays]
	Log *logrus.Logger
}

func NewShiftDayRepository(log *logrus.Logger) *ShiftDayRepository {
	return &ShiftDayRepository{
		Log: log,
	}
}

func (r *ShiftDayRepository) CreateBatch(db *gorm.DB, shiftDays []entity.ShiftDays) error {
	return db.Create(&shiftDays).Error
}

func (r *ShiftDayRepository) FindByShiftIDAndWeekday(db *gorm.DB, shiftDay *entity.ShiftDays, shiftID string, weekday int) error {
	type shiftDayRow struct {
		ID              int    `gorm:"column:id"`
		ShiftID         string `gorm:"column:shift_id"`
		Weekday         int    `gorm:"column:weekday"`
		DayType         string `gorm:"column:day_type"`
		CheckIn         string `gorm:"column:check_in"`
		CheckOut        string `gorm:"column:check_out"`
		BreakStart      string `gorm:"column:break_start"`
		BreakEnd        string `gorm:"column:break_end"`
		MaxBreakMinutes int    `gorm:"column:max_break_minutes"`
		CreatedAt       int64  `gorm:"column:created_at"`
		UpdatedAt       int64  `gorm:"column:updated_at"`
	}

	var row shiftDayRow
	if err := db.Table("shift_days").
		Select("id, shift_id, weekday, day_type, check_in, check_out, break_start, break_end, max_break_minutes, created_at, updated_at").
		Where("shift_id = ? AND weekday = ?", shiftID, weekday).
		Take(&row).Error; err != nil {
		return err
	}

	checkIn, _ := lib.ParseTimeToUnixMilli(row.CheckIn)

	checkOut, _ := lib.ParseTimeToUnixMilli(row.CheckOut)

	breakStart, _ := lib.ParseTimeToUnixMilli(row.BreakStart)

	breakEnd, _ := lib.ParseTimeToUnixMilli(row.BreakEnd)

	*shiftDay = entity.ShiftDays{
		ID:              row.ID,
		ShiftID:         row.ShiftID,
		Weekday:         row.Weekday,
		DayType:         row.DayType,
		CheckIn:         checkIn,
		CheckOut:        checkOut,
		BreakStart:      breakStart,
		BreakEnd:        breakEnd,
		MaxBreakMinutes: row.MaxBreakMinutes,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}

	return nil
}

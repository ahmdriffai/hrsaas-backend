package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"time"

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
		CreatedAt       time.Time `gorm:"column:created_at"`
		UpdatedAt       time.Time `gorm:"column:updated_at"`
	}

	var row shiftDayRow
	if err := db.Table("shift_days").
		Select("id, shift_id, weekday, day_type, check_in, check_out, break_start, break_end, max_break_minutes, created_at, updated_at").
		Where("shift_id = ? AND weekday = ?", shiftID, weekday).
		Take(&row).Error; err != nil {
		return err
	}

	checkIn, err := lib.ParseTimeHHMMOrHHMMSS(row.CheckIn)
	if err != nil {
		return err
	}
	checkOut, err := lib.ParseTimeHHMMOrHHMMSS(row.CheckOut)
	if err != nil {
		return err
	}
	breakStart, err := lib.ParseTimeHHMMOrHHMMSS(row.BreakStart)
	if err != nil {
		return err
	}
	breakEnd, err := lib.ParseTimeHHMMOrHHMMSS(row.BreakEnd)
	if err != nil {
		return err
	}

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

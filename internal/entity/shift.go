package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shift struct {
	ID            string `gorm:"column:id;primaryKey"`
	CompanyID     string `gorm:"column:company_id"`
	Name          string `gorm:"column:name"`
	LateTolerance int    `gorm:"column:late_tolerance"`
	CreatedAt     int64  `gorm:"column:created_at"`
	UpdatedAt     int64  `gorm:"column:updated_at"`

	Employees []Employee `gorm:"many2many:employee_shifts;joinForeignKey:shift_id;joinReferences:employee_id"`
	ShiftDays []ShiftDay `gorm:"foreignKey:shift_id;references:ID"`
}

func (s *Shift) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.NewString()
	s.CreatedAt = int64(time.Now().UnixMilli())
	s.UpdatedAt = int64(time.Now().UnixMilli())
	return nil
}

func (s *Shift) TableName() string {
	return "shifts"
}

type ShiftDay struct {
	ID              int    `gorm:"column:id;primaryKey"`
	ShiftID         string `gorm:"column:shift_id"`
	Weekday         int    `gorm:"column:weekday"`
	DayType         string `gorm:"column:day_type"`
	CheckIn         int64  `gorm:"column:check_in;type:time"`
	CheckOut        int64  `gorm:"column:check_out;type:time"`
	BreakStart      int64  `gorm:"column:break_start;type:time"`
	BreakEnd        int64  `gorm:"column:break_end;type:time"`
	MaxBreakMinutes int    `gorm:"column:max_break_minutes"`

	CreatedAt int64 `gorm:"column:created_at"`
	UpdatedAt int64 `gorm:"column:updated_at"`

	Shift Shift `gorm:"foreignKey:ShiftID;references:ID;constraint:OnDelete:CASCADE"`
}

func (s *ShiftDay) TableName() string {
	return "shift_days"
}

func (s *ShiftDay) BeforeCreate(tx *gorm.DB) (err error) {
	s.CreatedAt = int64(time.Now().UnixMilli())
	s.UpdatedAt = int64(time.Now().UnixMilli())
	return nil
}

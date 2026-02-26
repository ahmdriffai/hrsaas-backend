package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shift struct {
	ID            string    `gorm:"column:id;primaryKey"`
	CompanyID     string    `gorm:"column:company_id"`
	Name          string    `gorm:"column:name"`
	LateTolerance int       `gorm:"column:late_tolerance"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Employee []Employee `gorm:"many2many:employee_shifts;joinForeignKey:shift_id;joinReferences:employee_id"`
}

func (s *Shift) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.NewString()
	return nil
}

func (s *Shift) TableName() string {
	return "shifts"
}

type ShiftDays struct {
	ID              int       `gorm:"column:id;primaryKey"`
	ShiftID         string    `gorm:"column:shift_id"`
	Weekday         int       `gorm:"column:weekday"`
	DayType         string    `gorm:"column:day_type"`
	CheckIn         time.Time `gorm:"column:check_in;type:time"`
	CheckOut        time.Time `gorm:"column:check_out;type:time"`
	BreakStart      time.Time `gorm:"column:break_start;type:time"`
	BreakEnd        time.Time `gorm:"column:break_end;type:time"`
	MaxBreakMinutes int       `gorm:"column:max_break_minutes"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Shift Shift `gorm:"foreignKey:ShiftID;references:ID;constraint:OnDelete:CASCADE"`
}

func (ShiftDays) TableName() string {
	return "shift_days"
}

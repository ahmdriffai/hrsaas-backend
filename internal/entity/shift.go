package entity

import "time"

type Shift struct {
	ID            string    `json:"id"`
	CompanyID     string    `json:"company_id"`
	Name          string    `json:"name"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	LateTolerance int       `json:"late_tolerance"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Employee []Employee `gorm:"many2many:employee_shifts;joinForeignKey:shift_id;joinReferences:employee_id"`
}

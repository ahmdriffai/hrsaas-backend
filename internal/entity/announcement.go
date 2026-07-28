package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Announcement struct {
	ID         string   `gorm:"column:id;primaryKey"`
	CompanyID  string   `gorm:"column:company_id;not null"`
	EmployeeID string   `gorm:"column:employee_id;not null"`
	Employee   Employee `gorm:"foreignKey:EmployeeID;references:ID"`
	Title      string   `gorm:"column:title;not null"`
	Category   string   `gorm:"column:category;not null"`
	Content    string   `gorm:"column:content;not null"`
	CreatedAt  int64    `gorm:"column:created_at;autoCreateTime"`
}

func (a *Announcement) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	a.CreatedAt = int64(time.Now().UnixMilli())
	return nil
}

func (Announcement) TableName() string {
	return "announcements"
}

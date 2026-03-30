package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Division struct {
	ID          string  `gorm:"column:id;primaryKey"`
	CompanyID   string  `gorm:"column:company_id;not null"`
	Name        string  `gorm:"column:name;not null"`
	Description *string `gorm:"column:description"`
	CreatedAt   int64   `gorm:"column:created_at"`
	UpdatedAt   int64   `gorm:"column:updated_at"`
}

func (d *Division) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.NewString()
	d.CreatedAt = int64(time.Now().UnixMilli())
	d.UpdatedAt = int64(time.Now().UnixMilli())
	return nil
}

func (d *Division) TableName() string {
	return "divisions"
}

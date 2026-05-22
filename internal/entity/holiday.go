package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Holiday struct {
	ID         string `gorm:"column:id;primaryKey"`
	CompanyID  string `gorm:"column:company_id"`
	Name       string `gorm:"column:name"`
	IsNational bool   `gorm:"column:is_national"`
	Date       int64  `gorm:"column:date"`
	EndDate    int64  `gorm:"column:end_date"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdatedAt  int64  `gorm:"column:updated_at"`
}

func (h *Holiday) BeforeCreate(tx *gorm.DB) (err error) {
	h.ID = uuid.NewString()
	h.CreatedAt = int64(time.Now().UnixMilli())
	h.UpdatedAt = int64(time.Now().UnixMilli())
	return nil
}

func (h *Holiday) TableName() string {
	return "holidays"
}

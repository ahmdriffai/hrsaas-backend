package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID        string `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name;not null"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	p.CreatedAt = int64(time.Now().UnixMilli())
	p.UpdatedAt = int64(time.Now().UnixMilli())
	return nil
}

func (p *Permission) TableName() string {
	return "permissions"
}

package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Position struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;not null"`
	CompanyID string    `gorm:"column:company_id"`
	ParentID  *string   `gorm:"column:parent_id"`
	Parent    *Position `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Company   Company   `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:DELETE"`
	CreatedAt int64     `gorm:"column:created_at"`
	UpdatedAt int64     `gorm:"column:updated_at"`
}

// BeforeCreate hook to set UUID
func (u *Position) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	u.CreatedAt = int64(time.Now().UnixMilli())
	u.UpdatedAt = int64(time.Now().UnixMilli())
	return nil
}

func (c *Position) TableName() string {
	return "positions"
}

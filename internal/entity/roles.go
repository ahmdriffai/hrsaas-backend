package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          string       `gorm:"column:id;primaryKey"`
	Name        string       `gorm:"column:name;not null"`
	CompanyID   string       `gorm:"column:company_id;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	CreatedAt   int64        `gorm:"column:created_at"`
	UpdatedAt   int64        `gorm:"column:updated_at"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewString()
	r.CreatedAt = int64(time.Now().UnixMilli())
	r.UpdatedAt = int64(time.Now().UnixMilli())
	return nil
}

func (r *Role) TableName() string {
	return "roles"
}

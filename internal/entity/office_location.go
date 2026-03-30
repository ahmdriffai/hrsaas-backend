package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfficeLocation struct {
	ID        string `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name;not null"`
	Address   string `gorm:"column:address"`
	Lat       string `gorm:"column:lat"`
	Lng       string `gorm:"column:lng"`
	Radius    int    `gorm:"column:radius_meters"`
	IsActive  bool   `gorm:"column:is_active"`
	CompanyID string `gorm:"column:company_id;not null"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`

	Employee []Employee `gorm:"many2many:employee_office_locations;joinForeignKey:office_location_id;joinReferences:employee_id"`
}

// BeforeCreate hook to set UUID
func (u *OfficeLocation) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	u.CreatedAt = int64(time.Now().UnixMilli())
	u.UpdatedAt = int64(time.Now().UnixMilli())
	return nil
}

func (OfficeLocation) TableName() string {
	return "office_locations"
}

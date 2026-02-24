package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfficeLocation struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;not null"`
	Address   string    `gorm:"column:address"`
	Lat       string    `gorm:"column:lat"`
	Lng       string    `gorm:"column:lng"`
	Radius    int       `gorm:"column:radius_meters"`
	IsActive  bool      `gorm:"column:is_active"`
	CompanyID string    `gorm:"column:company_id;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Employee []Employee `gorm:"many2many:employee_office_locations;joinForeignKey:office_location_id;joinReferences:employee_id"`
}

// BeforeCreate hook to set UUID
func (u *OfficeLocation) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (OfficeLocation) TableName() string {
	return "office_locations"
}

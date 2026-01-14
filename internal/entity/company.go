package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	ID             string  `gorm:"column:id;primaryKey"`
	Name           string  `gorm:"column:name;not null"`
	LogoUrl        *string `gorm:"column:logo_url"`
	BussinessField *string `gorm:"column:bussiness_field"`
	Address        *string `gorm:"column:address"`
	Province       *string `gorm:"column:province"`
	City           *string `gorm:"column:city"`
	District       *string `gorm:"column:district"`
	Village        *string `gorm:"column:village"`
	ZipCode        *string `gorm:"column:zip_code"`
	PhoneNumber    *string `gorm:"column:phone_number"`
	FaxNumber      *string `gorm:"column:fax_number"`
	Email          *string `gorm:"column:email"`
	Website        *string `gorm:"column:website"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// BeforeCreate hook to set UUID
func (u *Company) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (c *Company) TableName() string {
	return "companies"
}

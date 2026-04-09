package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Visit struct {
	ID         string  `gorm:"column:id;primaryKey"`
	EmployeeID string  `gorm:"column:employee_id;not null"`
	CompanyID  string  `gorm:"column:company_id;not null"`
	VisitType  string  `gorm:"column:visit_type;not null"`
	Note       *string `gorm:"column:note"`
	Latitude   *string `gorm:"column:latitude"`
	Longitude  *string `gorm:"column:longitude"`
	Address    *string `gorm:"column:address"`
	FileURL    *string `gorm:"column:file_url"`
	CreatedAt  int64   `gorm:"column:created_at;not null"`
}

func (u *Visit) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	if u.CreatedAt == 0 {
		u.CreatedAt = time.Now().UnixMilli()
	}
	return nil
}

func (Visit) TableName() string {
	return "visits"
}

package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Visit struct {
	ID         string        `gorm:"column:id;primaryKey"`
	EmployeeID string        `gorm:"column:employee_id;not null"`
	CompanyID  string        `gorm:"column:company_id;not null"`
	Date       string        `gorm:"column:date;not null"`
	ClientName string        `gorm:"column:client_name;not null"`
	CreatedAt  int64         `gorm:"column:created_at;not null"`
	UpdatedAt  int64         `gorm:"column:updated_at;not null"`
	Employee   Employee      `gorm:"foreignKey:EmployeeID;references:ID"`
	Company    Company       `gorm:"foreignKey:CompanyID;references:ID"`
	Details    []VisitDetail `gorm:"foreignKey:VisitID;references:ID"`
}

func (v *Visit) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.NewString()
	v.CreatedAt = time.Now().UnixMilli()
	v.UpdatedAt = time.Now().UnixMilli()
	return nil
}

func (Visit) TableName() string {
	return "visits"
}

type VisitDetail struct {
	ID        string  `gorm:"column:id;primaryKey"`
	VisitID   string  `gorm:"column:visit_id;not null"`
	VisitType string  `gorm:"column:visit_type;not null"`
	VisitAt   string  `gorm:"column:visit_at;not null"`
	DateVisit string  `gorm:"column:date_visit;not null"`
	FileUrl   *string `gorm:"column:file_url"`
	Location  *string `gorm:"column:location"`
	Address   *string `gorm:"column:address"`
	Note      *string `gorm:"column:note"`
	CreatedAt int64   `gorm:"column:created_at;not null"`
	UpdatedAt int64   `gorm:"column:updated_at;not null"`
}

func (vd *VisitDetail) BeforeCreate(tx *gorm.DB) (err error) {
	vd.ID = uuid.NewString()
	vd.CreatedAt = time.Now().UnixMilli()
	vd.UpdatedAt = time.Now().UnixMilli()
	return nil
}

func (VisitDetail) TableName() string {
	return "visit_details"
}

package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeDocument struct {
	ID         string   `gorm:"column:id;primaryKey"`
	CompanyID  string   `gorm:"column:company_id;not null"`
	EmployeeID string   `gorm:"column:employee_id;not null"`
	DocType    string   `gorm:"column:doc_type;not null"`
	DocNumber  string   `gorm:"column:doc_number;not null"`
	DocName    string   `gorm:"column:doc_name;not null"`
	FileURL    string   `gorm:"column:file_url;not null"`
	IssuedAt   int64    `gorm:"column:issued_at;not null"`
	CreatedAt  int64    `gorm:"column:created_at;not null"`
	UpdatedAt  int64    `gorm:"column:updated_at;not null"`
	Employee   Employee `gorm:"foreignKey:EmployeeID"`
}

func (u *EmployeeDocument) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	u.CreatedAt = int64(time.Now().UnixMilli())
	u.UpdatedAt = int64(time.Now().UnixMilli())
	return nil
}

func (u *EmployeeDocument) TableName() string {
	return "employee_docs"
}

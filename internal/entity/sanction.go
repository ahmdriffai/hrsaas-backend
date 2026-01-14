package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sanction struct {
	ID          string    `gorm:"column:id;primaryKey"`
	CompanyID   string    `gorm:"column:company_id;not null"`
	Name        string    `gorm:"column:name;not null"`
	Description *string   `gorm:"column:description"`
	Level       int       `gorm:"column:level;not null"`
	Note        *string   `gorm:"column:note"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// BeforeCreate hook to set UUID
func (s *Sanction) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.NewString()
	return nil
}

func (s *Sanction) TableName() string {
	return "sanctions"
}

type EmployeeSanction struct {
	ID          string     `gorm:"column:id;primaryKey"`
	EmployeeID  string     `gorm:"column:employee_id;not null"`
	SanctionID  string     `gorm:"column:sanction_id;not null"`
	CompanyID   string     `gorm:"column:company_id;not null"`
	Reason      *string    `gorm:"column:reason"`
	StartDate   time.Time  `gorm:"column:start_date;not null"`
	EndDate     *time.Time `gorm:"column:end_date"`
	Status      *string    `gorm:"column:status"`
	DocumentUrl string     `gorm:"column:document_url"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`

	Employee Employee `gorm:"foreignKey:EmployeeID"`
	Sanction Sanction `gorm:"foreignKey:SanctionID"`
}

// BeforeCreate hook to set UUID
func (s *EmployeeSanction) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.NewString()
	return nil
}

func (s *EmployeeSanction) TableName() string {
	return "employee_sanctions"
}

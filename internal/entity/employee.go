package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	ID             string `gorm:"column:id;primaryKey"`
	CompanyID      string `gorm:"column:company_id;not null"`
	UserID         string `gorm:"column:user_id;not null"`
	EmployeeNumber string `gorm:"column:employee_number;uniqueIndex"`
	Fullname       string `gorm:"column:fullname;not null"`
	Gender         string `gorm:"column:gender;not null"`
	BirthPlace     string `gorm:"column:birth_place;not null"`
	BirthDate      string `gorm:"column:birth_date;not null"`
	BlodType       string `gorm:"column:blood_type;not null"`
	MaritalStatus  string `gorm:"column:marital_status;not null"`
	Religion       string `gorm:"column:religion;not null"`
	Phone          string `gorm:"column:phone;not null"`
	Timezone       string `gorm:"column:timezone;not null"`

	OfficeLocations []OfficeLocation `gorm:"many2many:employee_office_locations"`
}

// BeforeCreate hook to set UUID
func (u *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (c *Employee) TableName() string {
	return "employees"
}

type EmployeeIdentification struct {
	ID                         string `gorm:"column:id;primaryKey"`
	EmployeeID                 string `gorm:"column:employee_id;not null"`
	IdentityType               string `gorm:"column:identity_type;not null"`
	IdentityNumber             string `gorm:"column:identity_number;not null"`
	Address                    string `gorm:"column:address;not null"`
	City                       string `gorm:"column:city;not null"`
	PostalCode                 string `gorm:"column:postal_code;not null"`
	DomicililyAddress          string `gorm:"column:domicility_address;"`
	DomicililyCity             string `gorm:"column:domicility_city;"`
	DomicililyPostalCode       string `gorm:"column:domicility_postal_code"`
	IsDomicililySameAsIdentity bool   `gorm:"column:is_domicility_same_as_identity;not null"`
	IsDefault                  bool   `gorm:"column:is_default;not null"`
}

// BeforeCreate hook to set UUID
func (u *EmployeeIdentification) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (c *EmployeeIdentification) TableName() string {
	return "employee_identifications"
}

type EmployeeContract struct {
	ID           string  `gorm:"column:id;primaryKey"`
	EmployeeID   string  `gorm:"column:employee_id;not null"`
	ContractType string  `gorm:"column:contract_type;not null"`
	StartDate    string  `gorm:"column:start_date;not null"`
	EndDate      string  `gorm:"column:end_date;not null"`
	DivisionID   string  `gorm:"column:division_id;not null"`
	PositionID   string  `gorm:"column:position_id;not null"`
	Salary       float64 `gorm:"column:salary;not null"`
	IsDefault    bool    `gorm:"column:is_default;not null"`
}

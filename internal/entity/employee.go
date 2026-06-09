package entity

import (
	"time"

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
	IdentityNumber string `gorm:"column:identity_number;not null"`
	BirthPlace     string `gorm:"column:birth_place;not null"`
	BirthDate      int64  `gorm:"column:birth_date;not null"`
	BlodType       string `gorm:"column:blood_type;not null"`
	MaritalStatus  string `gorm:"column:marital_status;not null"`
	Religion       string `gorm:"column:religion;not null"`
	Phone          string `gorm:"column:phone;not null"`
	Address        string `gorm:"column:address;not null"`
	City           string `gorm:"column:city;not null"`
	Timezone       string `gorm:"column:timezone;not null"`
	CreatedAt      int64  `gorm:"column:created_at"`
	UpdatedAt      int64  `gorm:"column:updated_at"`

	User               User
	EmployeeContract   []EmployeeContract  `gorm:"foreignKey:EmployeeID;references:ID"`
	OfficeLocations    []OfficeLocation    `gorm:"many2many:employee_office_locations"`
	EmployeeDocuments  []EmployeeDocument  `gorm:"foreignKey:EmployeeID;references:ID"`
	EmployeeEducations []EmployeeEducation `gorm:"foreignKey:EmployeeID;references:ID"`
	EmployeeTrainings  []EmployeeTraining  `gorm:"foreignKey:EmployeeID;references:ID"`
}

// BeforeCreate hook to set UUID
func (u *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	u.CreatedAt = int64(time.Now().UnixMilli())
	u.UpdatedAt = int64(time.Now().UnixMilli())
	return nil
}

func (c *Employee) TableName() string {
	return "employees"
}

type EmployeeContract struct {
	ID           string   `gorm:"column:id;primaryKey"`
	EmployeeID   string   `gorm:"column:employee_id;not null"`
	ContractType string   `gorm:"column:contract_type;not null"`
	StartDate    int64    `gorm:"column:start_date;not null"`
	EndDate      *int64   `gorm:"column:end_date"`
	DivisionID   string   `gorm:"column:division_id;not null"`
	PositionID   string   `gorm:"column:position_id;not null"`
	Salary       float64  `gorm:"column:salary;not null"`
	Employee     Employee `gorm:"foreignKey:EmployeeID;references:ID"`
	Division     Division `gorm:"foreignKey:DivisionID;references:ID"`
	Position     Position `gorm:"foreignKey:PositionID;references:ID"`
	IsActive     bool     `gorm:"column:is_active;not null"`
}

// BeforeCreate hook to set UUID
func (c *EmployeeContract) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()
	return nil
}

type EmployeeEducation struct {
	ID              string   `gorm:"column:id;primaryKey"`
	CompanyID       string   `gorm:"column:company_id;not null"`
	EmployeeID      string   `gorm:"column:employee_id;not null"`
	EducationLevel  string   `gorm:"column:education_level;not null"`
	InstitutionName string   `gorm:"column:institution_name;not null"`
	Major           string   `gorm:"column:major;not null"`
	GraduationYear  int64    `gorm:"column:graduation_year"`
	GPA             *float64 `gorm:"column:gpa"`
	StartYear       *int     `gorm:"column:start_year"`
	EndYear         *int     `gorm:"column:end_year"`
	CreatedAt       int64    `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt       int64    `gorm:"column:updated_at;autoUpdateTime:milli"`
	Employee        Employee `gorm:"foreignKey:EmployeeID;references:ID"`
}

// BeforeCreate hook to set UUID
func (c *EmployeeEducation) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()
	return nil
}

func (c *EmployeeEducation) TableName() string {
	return "employee_educations"
}

type EmployeeTraining struct {
	ID             string   `gorm:"column:id;primaryKey"`
	CompanyID      string   `gorm:"column:company_id;not null"`
	EmployeeID     string   `gorm:"column:employee_id;not null"`
	TrainingName   string   `gorm:"column:training_name;not null"`
	Organizer      string   `gorm:"column:organizer;not null"`
	StartDate      int64    `gorm:"column:start_date;not null"`
	EndDate        *int64   `gorm:"column:end_date"`
	CertificateURL *string  `gorm:"column:certificate_url"`
	CreatedAt      int64    `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt      int64    `gorm:"column:updated_at;autoUpdateTime:milli"`
	Employee       Employee `gorm:"foreignKey:EmployeeID;references:ID"`
}

func (c *EmployeeTraining) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()
	return nil
}

func (c *EmployeeTraining) TableName() string {
	return "employee_trainings"
}

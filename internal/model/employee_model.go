package model

import "hr-sas/internal/entity"

type EmployeeResponse struct {
	ID             string `json:"id,omitempty"`
	CompanyID      string `json:"company_id,omitempty"`
	UserID         string `json:"user_id,omitempty"`
	EmployeeNumber string `json:"employee_number,omitempty"`
	Fullname       string `json:"fullname,omitempty"`
	BirthPlace     string `json:"birth_place,omitempty"`
	BirthDate      int64  `json:"birth_date,omitempty"`
	BlodType       string `json:"blood_type,omitempty"`
	MaritalStatus  string `json:"marital_status,omitempty"`
	Religion       string `json:"religion,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
	CreatedAt      int64  `json:"created_at,omitempty"`
	UpdatedAt      int64  `json:"updated_at,omitempty"`
}

type CreateEmployeeRequest struct {
	CompanyID      string `json:"-" validate:"required"`
	Fullname       string `json:"fullname" validate:"required"`
	EmployeeNumber string `json:"employee_number" validate:"required"`
	BirthPlace     string `json:"birth_place" validate:"required"`
	BirthDate      string `json:"birth_date" validate:"required"`
	BlodType       string `json:"blood_type" validate:"required"`
	MaritalStatus  string `json:"marital_status" validate:"required"`
	Religion       string `json:"religion" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	Timezone       string `json:"timezone" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=3"`
}

type SearchEmployeeRequest struct {
	CompanyID string `json:"company_id" validate:"required"`
	Key       string `json:"key" validate:"max=100"`
	Page      int    `json:"page" validate:"min=1"`
	Size      int    `json:"size" validate:"min=1,max=100"`
}

func EmployeeToResponse(employee *entity.Employee) *EmployeeResponse {
	if employee == nil {
		return nil
	}

	return &EmployeeResponse{
		ID:             employee.ID,
		CompanyID:      employee.CompanyID,
		UserID:         employee.UserID,
		Fullname:       employee.Fullname,
		BirthPlace:     employee.BirthPlace,
		BirthDate:      employee.BirthDate,
		BlodType:       employee.BlodType,
		MaritalStatus:  employee.MaritalStatus,
		Religion:       employee.Religion,
		Phone:          employee.Phone,
		Timezone:       employee.Timezone,
		EmployeeNumber: employee.EmployeeNumber,
	}
}

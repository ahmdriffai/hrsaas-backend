package model

import (
	"hr-sas/internal/entity"
	"mime/multipart"
)

type EmployeeResponse struct {
	ID                 string                      `json:"id,omitempty"`
	CompanyID          string                      `json:"company_id,omitempty"`
	UserID             string                      `json:"user_id,omitempty"`
	EmployeeNumber     string                      `json:"employee_number,omitempty"`
	Fullname           string                      `json:"fullname,omitempty"`
	Gender             string                      `json:"gender,omitempty"`
	BirthPlace         string                      `json:"birth_place,omitempty"`
	BirthDate          int64                       `json:"birth_date,omitempty"`
	IdentityNumber     string                      `json:"identity_number,omitempty"`
	BlodType           string                      `json:"blood_type,omitempty"`
	MaritalStatus      string                      `json:"marital_status,omitempty"`
	Religion           string                      `json:"religion,omitempty"`
	Phone              string                      `json:"phone,omitempty"`
	Timezone           string                      `json:"timezone,omitempty"`
	Contracts          []EmployeeContractResponse  `json:"contracts,omitempty"`
	EmployeeDocs       []EmployeeDocumentResponse  `json:"employee_docs,omitempty"`
	EmployeeEducations []EmployeeEducationResponse `json:"employee_educations,omitempty"`
	EmployeeTrainings  []EmployeeTrainingResponse  `json:"employee_trainings,omitempty"`
	User               *UserResponse               `json:"user,omitempty"`
	CreatedAt          int64                       `json:"created_at,omitempty"`
	UpdatedAt          int64                       `json:"updated_at,omitempty"`
}

type CreateEmployeeRequest struct {
	CompanyID      string `json:"-" validate:"required"`
	Fullname       string `json:"fullname" validate:"required"`
	Gender         string `json:"gender" validate:"required"`
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

type UpdateEmployeeRequest struct {
	ID             string  `json:"-"`
	Fullname       *string `json:"fullname,omitempty"`
	Gender         *string `json:"gender,omitempty"`
	EmployeeNumber *string `json:"employee_number,omitempty"`
	BirthPlace     *string `json:"birth_place,omitempty"`
	BirthDate      *string `json:"birth_date,omitempty"`
	BlodType       *string `json:"blood_type,omitempty"`
	MaritalStatus  *string `json:"marital_status,omitempty"`
	Religion       *string `json:"religion,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Timezone       *string `json:"timezone,omitempty"`
	Email          *string `json:"email,omitempty" validate:"omitempty,email"`
}

type SearchEmployeeRequest struct {
	CompanyID string `json:"company_id" validate:"required"`
	Key       string `json:"key" validate:"max=100"`
	Page      int    `json:"page" validate:"min=1"`
	Size      int    `json:"size" validate:"min=1,max=1000"`
}

type ImportExcelEmployeeRequest struct {
	CompanyID string                `json:"company_id" validate:"required"`
	File      *multipart.FileHeader `form:"file" validate:"required"`
}

type DetailEmployeeRequest struct {
	ID        string `json:"id" validate:"required"`
	CompanyID string `json:"company_id" validate:"required"`
}

func EmployeeToResponse(employee *entity.Employee) *EmployeeResponse {
	if employee == nil {
		return nil
	}

	var identityNumber string
	if employee.EmployeeIdentifications != nil {
		identityNumber = employee.EmployeeIdentifications.IdentityNumber
	}

	contracts := EmployeeContractsToResponse(employee.EmployeeContract)
	employeeDocs := EmployeeDocumentsToResponse(employee.EmployeeDocuments)
	user := UserToResponse(&employee.User)

	response := &EmployeeResponse{
		ID:             employee.ID,
		CompanyID:      employee.CompanyID,
		UserID:         employee.UserID,
		Fullname:       employee.Fullname,
		Gender:         employee.Gender,
		BirthPlace:     employee.BirthPlace,
		BirthDate:      employee.BirthDate,
		IdentityNumber: identityNumber,
		BlodType:       employee.BlodType,
		MaritalStatus:  employee.MaritalStatus,
		Religion:       employee.Religion,
		Phone:          employee.Phone,
		Timezone:       employee.Timezone,
		EmployeeNumber: employee.EmployeeNumber,
		Contracts:      contracts,
		EmployeeDocs:   employeeDocs,
		CreatedAt:      employee.CreatedAt,
		UpdatedAt:      employee.UpdatedAt,
	}

	response.User = user

	return response
}

func EmployeeSummaryToResponse(employee *entity.Employee) *EmployeeResponse {
	if employee == nil {
		return nil
	}

	response := &EmployeeResponse{
		ID:             employee.ID,
		UserID:         employee.UserID,
		EmployeeNumber: employee.EmployeeNumber,
		Fullname:       employee.Fullname,
		BirthPlace:     employee.BirthPlace,
		Timezone:       employee.Timezone,
	}
	return response
}

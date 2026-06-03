package model

import "hr-sas/internal/entity"

type EmployeeEducationResponse struct {
	ID              string   `json:"id"`
	CompanyID       string   `json:"company_id"`
	EmployeeID      string   `json:"employee_id"`
	EmployeeName    string   `json:"employee_name"`
	EducationLevel  string   `json:"education_level"`
	InstitutionName string   `json:"institution_name"`
	Major           string   `json:"major"`
	GraduationYear  int64    `json:"graduation_year"`
	GPA             *float64 `json:"gpa,omitempty"`
	StartYear       *int     `json:"start_year,omitempty"`
	EndYear         *int     `json:"end_year,omitempty"`
}

type CreateEmployeeEducationRequest struct {
	CompanyID       string   `json:"-" validate:"required"`
	EmployeeID      string   `json:"employee_id" validate:"required"`
	EducationLevel  string   `json:"education_level" validate:"required"`
	InstitutionName string   `json:"institution_name" validate:"required"`
	Major           string   `json:"major" validate:"required"`
	GraduationYear  string   `json:"graduation_year" validate:"required"`
	GPA             *float64 `json:"gpa,omitempty" validate:"omitempty,min=0,max=4"`
	StartYear       *int     `json:"start_year,omitempty"`
	EndYear         *int     `json:"end_year,omitempty"`
}

type UpdateEmployeeEducationRequest struct {
	ID              string   `json:"-" validate:"required"`
	EducationLevel  *string  `json:"education_level,omitempty"`
	InstitutionName *string  `json:"institution_name,omitempty"`
	Major           *string  `json:"major,omitempty"`
	GraduationYear  *string  `json:"graduation_year,omitempty"`
	GPA             *float64 `json:"gpa,omitempty" validate:"omitempty,min=0,max=4"`
	StartYear       *int     `json:"start_year,omitempty"`
	EndYear         *int     `json:"end_year,omitempty"`
}

type SearchEmployeeEducationRequest struct {
	CompanyID  string `json:"-" validate:"required"`
	EmployeeID string `json:"employee_id,omitempty"`
	Page       int    `json:"page" validate:"min=1"`
	Size       int    `json:"size" validate:"min=1,max=100"`
}

func EmployeeEducationToResponse(e *entity.EmployeeEducation) *EmployeeEducationResponse {
	if e == nil {
		return nil
	}
	return &EmployeeEducationResponse{
		ID:              e.ID,
		CompanyID:       e.CompanyID,
		EmployeeID:      e.EmployeeID,
		EmployeeName:    e.Employee.Fullname,
		EducationLevel:  e.EducationLevel,
		InstitutionName: e.InstitutionName,
		Major:           e.Major,
		GraduationYear:  e.GraduationYear,
		GPA:             e.GPA,
		StartYear:       e.StartYear,
		EndYear:         e.EndYear,
	}
}

func EmployeeEducationsToResponse(list []entity.EmployeeEducation) []EmployeeEducationResponse {
	if len(list) == 0 {
		return nil
	}
	responses := make([]EmployeeEducationResponse, 0, len(list))
	for i := range list {
		if r := EmployeeEducationToResponse(&list[i]); r != nil {
			responses = append(responses, *r)
		}
	}
	return responses
}

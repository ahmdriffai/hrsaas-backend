package model

import "hr-sas/internal/entity"

type EmployeeTrainingResponse struct {
	ID             string  `json:"id"`
	CompanyID      string  `json:"company_id"`
	EmployeeID     string  `json:"employee_id"`
	EmployeeName   string  `json:"employee_name"`
	TrainingName   string  `json:"training_name"`
	Organizer      string  `json:"organizer"`
	StartDate      int64   `json:"start_date"`
	EndDate        *int64  `json:"end_date,omitempty"`
	CertificateURL *string `json:"certificate_url,omitempty"`
}

type CreateEmployeeTrainingRequest struct {
	CompanyID      string  `json:"-" validate:"required"`
	EmployeeID     string  `json:"employee_id" validate:"required"`
	TrainingName   string  `json:"training_name" validate:"required"`
	Organizer      string  `json:"organizer" validate:"required"`
	StartDate      string  `json:"start_date" validate:"required"`
	EndDate        *string `json:"end_date,omitempty"`
	CertificateURL *string `json:"certificate_url,omitempty" validate:"omitempty,url"`
}

type UpdateEmployeeTrainingRequest struct {
	ID             string  `json:"-" validate:"required"`
	TrainingName   *string `json:"training_name,omitempty"`
	Organizer      *string `json:"organizer,omitempty"`
	StartDate      *string `json:"start_date,omitempty"`
	EndDate        *string `json:"end_date,omitempty"`
	CertificateURL *string `json:"certificate_url,omitempty" validate:"omitempty,url"`
}

type SearchEmployeeTrainingRequest struct {
	CompanyID  string `json:"-" validate:"required"`
	EmployeeID string `json:"employee_id,omitempty"`
	Page       int    `json:"page" validate:"min=1"`
	Size       int    `json:"size" validate:"min=1,max=100"`
}

func EmployeeTrainingToResponse(t *entity.EmployeeTraining) *EmployeeTrainingResponse {
	if t == nil {
		return nil
	}
	return &EmployeeTrainingResponse{
		ID:             t.ID,
		CompanyID:      t.CompanyID,
		EmployeeID:     t.EmployeeID,
		EmployeeName:   t.Employee.Fullname,
		TrainingName:   t.TrainingName,
		Organizer:      t.Organizer,
		StartDate:      t.StartDate,
		EndDate:        t.EndDate,
		CertificateURL: t.CertificateURL,
	}
}

func EmployeeTrainingsToResponse(list []entity.EmployeeTraining) []EmployeeTrainingResponse {
	if len(list) == 0 {
		return nil
	}
	responses := make([]EmployeeTrainingResponse, 0, len(list))
	for i := range list {
		if r := EmployeeTrainingToResponse(&list[i]); r != nil {
			responses = append(responses, *r)
		}
	}
	return responses
}

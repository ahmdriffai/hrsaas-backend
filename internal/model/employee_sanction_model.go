package model

import (
	"hr-sas/internal/entity"
)

type EmSancResponse struct {
	ID          string           `json:"id"`
	EmployeeID  string           `json:"employee_id"`
	SanctionID  string           `json:"sanction_id"`
	CompanyID   string           `json:"company_id"`
	Reason      *string          `json:"reason,omitempty"`
	StartDate   int64            `json:"start_date,omitempty"`
	EndDate     *int64           `json:"end_date,omitempty"`
	Status      string           `json:"status,omitempty"`
	DocumentUrl string           `json:"document_url"`
	Employee    EmployeeResponse `json:"employee"`
	Sanction    SanctionResponse `json:"sanction"`
	CreatedAt   int64            `json:"created_at"`
	UpdatedAt   int64            `json:"updated_at"`
}

type CreateEmSancRequest struct {
	EmployeeID  string `json:"employee_id" validate:"required"`
	SanctionID  string `json:"sanction_id" validate:"required"`
	CompanyID   string `json:"-" validate:"required"`
	Reason      string `json:"reason" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date,omitempty"`
	DocumentUrl string `json:"document_url"`
}

type SearchEmSancRequest struct {
	EmployeeID string `json:"-"`
	SanctionID string `json:"sanction_id" validate:"max=100"`
	CompanyID  string `json:"company_id" validate:"required"`
	Reason     string `json:"reason" validate:"max=100"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Status     string `json:"status" validate:"max=10"`
	Page       int    `json:"page" validate:"min=1"`
	Size       int    `json:"size" validate:"min=1,max=100"`
}

func EmSancToResponse(emSanc *entity.EmployeeSanction) *EmSancResponse {
	return &EmSancResponse{
		ID:          emSanc.ID,
		EmployeeID:  emSanc.EmployeeID,
		SanctionID:  emSanc.SanctionID,
		Reason:      emSanc.Reason,
		StartDate:   emSanc.StartDate,
		EndDate:     emSanc.EndDate,
		CompanyID:   emSanc.CompanyID,
		Employee:    *EmployeeToResponse(&emSanc.Employee),
		DocumentUrl: emSanc.DocumentUrl,
		Sanction:    *SanctionToResponse(&emSanc.Sanction),
		CreatedAt:   emSanc.CreatedAt,
		UpdatedAt:   emSanc.UpdatedAt,
	}
}

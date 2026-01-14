package model

import (
	"hr-sas/internal/lib"
	"time"
)

type EmSancResponse struct {
	ID          string           `json:"id"`
	EmployeeID  string           `json:"employee_id"`
	SanctionID  string           `json:"sanction_id"`
	CompanyID   string           `json:"company_id"`
	Reason      *string          `json:"reason,omitempty"`
	StartDate   time.Time        `json:"start_date,omitempty"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
	Status      string           `json:"status,omitempty"`
	DocumentUrl string           `json:"document_url"`
	Employee    EmployeeResponse `json:"employee"`
	Sanction    SanctionResponse `json:"sanction"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
}

type CreateEmSancRequest struct {
	EmployeeID  string       `json:"employee_id" validate:"required"`
	SanctionID  string       `json:"sanction_id" validate:"required"`
	CompanyID   string       `json:"-" validate:"required"`
	Reason      string       `json:"reason" validate:"required"`
	StartDate   lib.DateOnly `json:"start_date" validate:"required"`
	EndDate     lib.DateOnly `json:"end_date,omitempty"`
	DocumentUrl string       `json:"document_url"`
}

type SearchEmSancRequest struct {
	UserID     string        `json:"user_id" validate:"max=100"`
	SanctionID string        `json:"sanction_id" validate:"max=100"`
	CompanyID  string        `json:"company_id" validate:"required"`
	Reason     string        `json:"reason" validate:"max=100"`
	StartDate  *lib.DateOnly `json:"start_date"`
	EndDate    *lib.DateOnly `json:"end_date"`
	Status     string        `json:"status" validate:"max=10"`
	Page       int           `json:"page" validate:"min=1"`
	Size       int           `json:"size" validate:"min=1,max=100"`
}

package model

import "hr-sas/internal/entity"

type EmployeeContractResponse struct {
	ID           string  `json:"id"`
	EmployeeID   string  `json:"employee_id"`
	ContractType string  `json:"contract_type"`
	StartDate    int64   `json:"start_date"`
	EndDate      *int64  `json:"end_date,omitempty"`
	DivisionID   string  `json:"division_id"`
	PositionID   string  `json:"position_id"`
	Salary       float64 `json:"salary"`
	// Employee     EmployeeResponse `json:"employee"`
	Division DivisionResponse `json:"division"`
	Position PositionResponse `json:"position"`
}

type CreateEmployeeContractRequest struct {
	EmployeeID   string  `json:"employee_id" validate:"required"`
	ContractType string  `json:"contract_type" validate:"required"`
	StartDate    string  `json:"start_date" validate:"required"`
	EndDate      *string `json:"end_date"`
	DivisionID   string  `json:"division_id" validate:"required"`
	PositionID   string  `json:"position_id" validate:"required"`
	Salary       float64 `json:"salary" validate:"required,min=0"`
}

type SearchEmployeeContractRequest struct {
	EmployeeID string `json:"employee_id" validate:"max=100"`
	DivisionID string `json:"division_id" validate:"max=100"`
	PositionID string `json:"position_id" validate:"max=100"`
	ActiveOnly bool   `json:"active_only"`
	Page       int    `json:"page" validate:"min=1"`
	Size       int    `json:"size" validate:"min=1,max=100"`
}

func EmployeeContractToResponse(contract *entity.EmployeeContract) *EmployeeContractResponse {
	return &EmployeeContractResponse{
		ID:           contract.ID,
		EmployeeID:   contract.EmployeeID,
		ContractType: contract.ContractType,
		StartDate:    contract.StartDate,
		EndDate:      contract.EndDate,
		DivisionID:   contract.DivisionID,
		PositionID:   contract.PositionID,
		Salary:       contract.Salary,
		// Employee:     *EmployeeToResponse(&contract.Employee),
		Division: *DivisionToResponse(&contract.Division),
		Position: *PositionToResponse(&contract.Position),
	}
}

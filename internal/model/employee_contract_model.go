package model

import "hr-sas/internal/entity"

type EmployeeContractResponse struct {
	ID           string           `json:"id"`
	EmployeeID   string           `json:"employee_id"`
	ContractType string           `json:"contract_type"`
	StartDate    int64            `json:"start_date"`
	EndDate      *int64           `json:"end_date,omitempty"`
	DivisionID   string           `json:"division_id"`
	PositionID   string           `json:"position_id"`
	Salary       float64          `json:"salary"`
	IsActive     bool             `json:"is_active"`
	Employee     EmployeeResponse `json:"employee,omitempty"`
	User         *UserResponse    `json:"user,omitempty"`
	Division     DivisionResponse `json:"division,omitempty"`
	Position     PositionResponse `json:"position,omitempty"`
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

type UpdateEmployeeContractRequest struct {
	ContractType *string  `json:"contract_type,omitempty"`
	StartDate    *string  `json:"start_date,omitempty"`
	EndDate      *string  `json:"end_date,omitempty"`
	DivisionID   *string  `json:"division_id,omitempty"`
	PositionID   *string  `json:"position_id,omitempty"`
	Salary       *float64 `json:"salary,omitempty" validate:"omitempty,min=0"`
	IsActive     *bool    `json:"is_active,omitempty"`
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
	if contract == nil {
		return nil
	}

	user := UserToResponse(&contract.Employee.User)

	response := &EmployeeContractResponse{
		ID:           contract.ID,
		EmployeeID:   contract.EmployeeID,
		ContractType: contract.ContractType,
		StartDate:    contract.StartDate,
		EndDate:      contract.EndDate,
		DivisionID:   contract.DivisionID,
		PositionID:   contract.PositionID,
		Salary:       contract.Salary,
		IsActive:     contract.IsActive,
	}

	response.User = user

	if employee := EmployeeSummaryToResponse(&contract.Employee); employee != nil {
		response.Employee = *employee
	}
	if division := DivisionToResponse(&contract.Division); division != nil {
		response.Division = *division
	}
	if position := PositionToResponse(&contract.Position); position != nil {
		response.Position = *position
	}

	return response
}

func EmployeeContractsToResponse(contracts []entity.EmployeeContract) []EmployeeContractResponse {
	if len(contracts) == 0 {
		return nil
	}

	responses := make([]EmployeeContractResponse, 0, len(contracts))
	for i := range contracts {
		if response := EmployeeContractToResponse(&contracts[i]); response != nil {
			responses = append(responses, *response)
		}
	}

	return responses
}

package model

type EmployeeContractResponse struct {
	ID           string  `json:"id"`
	EmployeeID   string  `json:"employee_id"`
	ContractType string  `json:"contract_type"`
	StartDate    string  `json:"start_date"`
	EndDate      *string `json:"end_date,omitempty"`
	DivisionID   string  `json:"division_id"`
	PositionID   string  `json:"position_id"`
	Salary       float64 `json:"salary"`
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

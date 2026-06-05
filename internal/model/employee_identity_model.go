package model

import "hr-sas/internal/entity"

type EmployeeIdentityResponse struct {
	ID                         string `json:"id"`
	EmployeeID                 string `json:"employee_id"`
	IdentityType               string `json:"identity_type"`
	IdentityNumber             string `json:"identity_number"`
	Address                    string `json:"address"`
	City                       string `json:"city"`
	PostalCode                 string `json:"postal_code"`
	DomicililyAddress          string `json:"domicile_address,omitempty"`
	IsDomicililySameAsIdentity bool   `json:"domicily_as_ktp"`
	IsDefault                  bool   `json:"is_default"`
}

type CreateEmployeeIdentityRequest struct {
	EmployeeID                 string `json:"employee_id" validate:"required"`
	IdentityType               string `json:"identity_type" validate:"required"`
	IdentityNumber             string `json:"identity_number" validate:"required"`
	Address                    string `json:"address" validate:"required"`
	City                       string `json:"city" validate:"required"`
	PostalCode                 string `json:"postal_code" validate:"required"`
	DomicililyAddress          string `json:"domicile_address,omitempty"`
	IsDomicililySameAsIdentity bool   `json:"domicily_as_ktp"`
	IsDefault                  bool   `json:"is_default"`
}

type UpdateEmployeeIdentityRequest struct {
	IdentityType               *string `json:"identity_type,omitempty"`
	IdentityNumber             *string `json:"identity_number,omitempty"`
	Address                    *string `json:"address,omitempty"`
	City                       *string `json:"city,omitempty"`
	PostalCode                 *string `json:"postal_code,omitempty"`
	DomicililyAddress          *string `json:"domicile_address,omitempty"`
	IsDomicililySameAsIdentity *bool   `json:"domicily_as_ktp,omitempty"`
	IsDefault                  *bool   `json:"is_default,omitempty"`
}

type SearchEmployeeIdentityRequest struct {
	EmployeeID string `json:"employee_id" validate:"required"`
	Page       int    `json:"page" validate:"min=1"`
	Size       int    `json:"size" validate:"min=1,max=100"`
}

func EmployeeIdentityToResponse(e *entity.EmployeeIdentity) *EmployeeIdentityResponse {
	if e == nil {
		return nil
	}
	return &EmployeeIdentityResponse{
		ID:                         e.ID,
		EmployeeID:                 e.EmployeeID,
		IdentityType:               e.IdentityType,
		IdentityNumber:             e.IdentityNumber,
		Address:                    e.Address,
		City:                       e.City,
		PostalCode:                 e.PostalCode,
		DomicililyAddress:          e.DomicililyAddress,
		IsDomicililySameAsIdentity: e.IsDomicililySameAsIdentity,
		IsDefault:                  e.IsDefault,
	}
}

func EmployeeIdentitiesToResponse(items []entity.EmployeeIdentity) []EmployeeIdentityResponse {
	responses := make([]EmployeeIdentityResponse, len(items))
	for i, item := range items {
		responses[i] = *EmployeeIdentityToResponse(&item)
	}
	return responses
}

func EmployeeIdentitiesNumberResponse(e *entity.EmployeeIdentity) string {
	if e == nil {
		return ""
	}
	return e.IdentityNumber
}

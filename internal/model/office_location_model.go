package model

import "hr-sas/internal/entity"

type OfficeLocationResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Lat      string `json:"lat"`
	Lng      string `json:"lng"`
	Radius   int    `json:"radius_meters"`
	IsActive bool   `json:"is_active"`
}

type CreateOfficeLocationRequest struct {
	Name      string `json:"name" validate:"required"`
	Address   string `json:"address" validate:"required"`
	Lat       string `json:"lat" validate:"required"`
	Lng       string `json:"lng" validate:"required"`
	Radius    int    `json:"radius" validate:"required,min=0"`
	CompanyID string `json:"-" validate:"required"`
}

type SearchOfficeLocationRequest struct {
	CompanyID string `json:"-" validate:"required"`
	Key       string `json:"key" validate:"max=100"`
	Page      int    `json:"page" validate:"min=1"`
	Size      int    `json:"size" validate:"min=1,max=100"`
}

// converter
func OfficeLocationToResponse(officeLocation *entity.OfficeLocation) *OfficeLocationResponse {
	return &OfficeLocationResponse{
		ID:       officeLocation.ID,
		Name:     officeLocation.Name,
		Address:  officeLocation.Address,
		Lat:      officeLocation.Lat,
		Lng:      officeLocation.Lng,
		Radius:   officeLocation.Radius,
		IsActive: officeLocation.IsActive,
	}
}

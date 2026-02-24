package model

import (
	"hr-sas/internal/entity"
	"strconv"
)

type OfficeLocationResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Radius   int     `json:"radius_meters"`
	IsActive bool    `json:"is_active"`
}

type CreateOfficeLocationRequest struct {
	Name      string  `json:"name" validate:"required"`
	Address   string  `json:"address" validate:"required"`
	Lat       float64 `json:"lat" validate:"required"`
	Lng       float64 `json:"lng" validate:"required"`
	Radius    int     `json:"radius" validate:"required,min=0"`
	CompanyID string  `json:"-" validate:"required"`
}

type SearchOfficeLocationRequest struct {
	CompanyID string `json:"-" validate:"required"`
	Key       string `json:"key" validate:"max=100"`
	Page      int    `json:"page" validate:"min=1"`
	Size      int    `json:"size" validate:"min=1,max=100"`
}

type AssignEmployeeToOfficeLocationRequest struct {
	CompanyID        string `json:"-" validate:"required"`
	EmployeeID       string `json:"employee_id" validate:"required"`
	OfficeLocationID string `json:"office_location_id" validate:"required"`
}

// converter
func OfficeLocationToResponse(officeLocation *entity.OfficeLocation) *OfficeLocationResponse {
	lat, _ := strconv.ParseFloat(officeLocation.Lat, 64)
	lng, _ := strconv.ParseFloat(officeLocation.Lng, 64)

	return &OfficeLocationResponse{
		ID:       officeLocation.ID,
		Name:     officeLocation.Name,
		Address:  officeLocation.Address,
		Lat:      lat,
		Lng:      lng,
		Radius:   officeLocation.Radius,
		IsActive: officeLocation.IsActive,
	}
}

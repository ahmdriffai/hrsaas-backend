package model

import (
	"hr-sas/internal/entity"
)

type HolidayResponse struct {
	ID         string `json:"id,omitempty"`
	CompanyID  string `json:"company_id,omitempty"`
	Name       string `json:"name,omitempty"`
	IsNational bool   `json:"is_national,omitempty"`
	Date       int64  `json:"date,omitempty"`
	EndDate    int64  `json:"end_date,omitempty"`
	CreatedAt  int64  `json:"created_at,omitempty"`
	UpdatedAt  int64  `json:"updated_at,omitempty"`
}

type CreateHolidayRequest struct {
	CompanyID  string `json:"-" validate:"required"`
	Name       string `json:"name" validate:"required"`
	IsNational bool   `json:"is_national"`
	Date       int64  `json:"date" validate:"required"`
	EndDate    int64  `json:"end_date"`
}

type UpdateHolidayRequest struct {
	ID         string  `json:"-"`
	CompanyID  string  `json:"-" validate:"required"`
	Name       *string `json:"name,omitempty"`
	IsNational *bool   `json:"is_national,omitempty"`
	Date       *int64  `json:"date,omitempty"`
	EndDate    *int64  `json:"end_date,omitempty"`
}

type SearchHolidayRequest struct {
	CompanyID string `json:"-" validate:"required"`
	Key       string `json:"key" validate:"max=100"`
	Page      int    `json:"page" validate:"min=1"`
	Size      int    `json:"size" validate:"min=1,max=100"`
}

type DetailHolidayRequest struct {
	CompanyID string `json:"-" validate:"required"`
	HolidayID string `json:"-" validate:"required"`
}

func HolidayToResponse(holiday *entity.Holiday) *HolidayResponse {
	response := &HolidayResponse{
		ID:         holiday.ID,
		CompanyID:  holiday.CompanyID,
		Name:       holiday.Name,
		IsNational: holiday.IsNational,
		Date:       holiday.Date,
		EndDate:    holiday.EndDate,
		CreatedAt:  holiday.CreatedAt,
		UpdatedAt:  holiday.UpdatedAt,
	}
	return response
}

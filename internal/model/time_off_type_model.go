package model

import "hr-sas/internal/entity"

type TimeOffTypeResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Category         string `json:"category"`
	IsQuotaBased     bool   `json:"is_quota_based"`
	DefaultQuotaDays int    `json:"default_quota_days"`
}

type CreateTimeOffTypeRequest struct {
	Name             string `json:"name" validate:"required,max=100"`
	Category         string `json:"category" validate:"required,oneof=IZIN SAKIT CUTI"`
	IsQuotaBased     bool   `json:"is_quota_based"`
	DefaultQuotaDays int    `json:"default_quota_days" validate:"min=0"`
}

type UpdateTimeOffTypeRequest struct {
	Name             *string `json:"name" validate:"required,max=100,omitempty"`
	Category         *string `json:"category" validate:"required,oneof=IZIN SAKIT CUTI,omitempty"`
	IsQuotaBased     *bool   `json:"is_quota_based,omitempty"`
	DefaultQuotaDays *int    `json:"default_quota_days" validate:"min=0,omitempty"`
}

type SearchTimeOffTypeRequest struct {
	Name string `json:"name,omitempty" validate:"max=100"`
	Page int    `json:"page,omitempty" validate:"min=1"`
	Size int    `json:"size,omitempty" validate:"min=1,max=100"`
}

func TimeOffTypeToResponse(timeOffType *entity.TimeOffType) *TimeOffTypeResponse {
	return &TimeOffTypeResponse{
		ID:               timeOffType.ID,
		Name:             timeOffType.Name,
		Category:         timeOffType.Category,
		IsQuotaBased:     timeOffType.IsQuotaBased,
		DefaultQuotaDays: timeOffType.DefaultQuotaDays,
	}
}

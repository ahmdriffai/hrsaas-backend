package model

import "hr-sas/internal/entity"

type SanctionResponse struct {
	ID          string  `json:"id"`
	CompanyID   string  `json:"company_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Level       int     `json:"level,omitempty"`
	Note        *string `json:"note,omitempty"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

type CreateSanctionRequest struct {
	CompanyID   string  `json:"-" validate:"required"` // Added CompanyID field
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Level       int     `json:"level,omitempty"`
	Note        *string `json:"note,omitempty"`
}

type SearchSanctionRequest struct {
	CompanyID string `json:"-" validate:"required"` // Added CompanyID field
	Key       string `json:"key"`
	Page      int    `json:"page" validate:"min=1"`
	Size      int    `json:"size" validate:"min=1,max=100"`
}

func SanctionToResponse(sanction *entity.Sanction) *SanctionResponse {
	return &SanctionResponse{
		ID:          sanction.ID,
		CompanyID:   sanction.CompanyID,
		Description: sanction.Description,
		Note:        sanction.Note,
		Name:        sanction.Name,
		CreatedAt:   sanction.CreatedAt,
		UpdatedAt:   sanction.UpdatedAt,
	}
}

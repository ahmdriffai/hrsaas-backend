package model

type SanctionResponse struct {
	ID          string  `json:"id"`
	CompanyID   string  `json:"company_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Level       int     `json:"level,omitempty"`
	Note        *string `json:"note,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
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

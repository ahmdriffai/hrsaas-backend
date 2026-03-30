package model

type DivisionResponse struct {
	ID          string  `json:"id"`
	CompanyID   string  `json:"company_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type CreateDivisionRequest struct {
	CompanyID   string  `json:"-"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type SearchDivisionRequest struct {
	CompanyID string `json:"-" validate:"required"`
	Name      string `json:"name"`
	Page      int    `json:"page" validate:"min=1"`
	Size      int    `json:"size" validate:"min=1,max=100"`
}

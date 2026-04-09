package model

type PositionResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	CompanyID string            `json:"company_id"`
	ParentID  *string           `json:"parent_id,omitempty"`
	Parent    *PositionResponse `json:"parent,omitempty"`
	IsApprover bool             `json:"is_approver"`
}

type CreatePositionRequest struct {
	CompanyID string  `json:"-"`
	Name      string  `json:"name" validate:"required"`
	ParentID  *string `json:"parent_id,omitempty"`
	IsApprover bool   `json:"is_approver"`
}

type SeachPositionRequest struct {
	CompanyID string `json:"-" validate:"required"`
	Name      string `json:"name"`
	Page      int    `json:"page" validate:"min=1"`
	Size      int    `json:"size" validate:"min=1,max=100"`
}

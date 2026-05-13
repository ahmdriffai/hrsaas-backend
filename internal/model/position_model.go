package model

import "hr-sas/internal/entity"

type PositionResponse struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	IsApprover bool              `json:"is_approver"`
	CompanyID  string            `json:"company_id"`
	ParentID   *string           `json:"parent_id,omitempty"`
	Parent     *PositionResponse `json:"parent,omitempty"`
}

type CreatePositionRequest struct {
	CompanyID  string  `json:"-"`
	Name       string  `json:"name" validate:"required"`
	ParentID   *string `json:"parent_id,omitempty"`
	IsApprover bool    `json:"is_approver"`
}

type UpdatePositionRequest struct {
	Name       *string `json:"name,omitempty"`
	ParentID   *string `json:"parent_id,omitempty"`
	IsApprover *bool   `json:"is_approver,omitempty"`
}

type SearchPositionRequest struct {
	CompanyID string `json:"-" validate:"required"`
	Name      string `json:"name"`
	Page      int    `json:"page" validate:"min=1"`
	Size      int    `json:"size" validate:"min=1,max=100"`
}

func PositionToResponse(position *entity.Position) *PositionResponse {
	if position == nil {
		return nil
	}
	return &PositionResponse{
		ID:         position.ID,
		Name:       position.Name,
		CompanyID:  position.CompanyID,
		ParentID:   position.ParentID,
		IsApprover: position.IsApprover,
	}
}

package model

import (
	"hr-sas/internal/entity"
)

type AnnouncementResponse struct {
	ID           string `json:"id"`
	CompanyID    string `json:"company_id"`
	EmployeeID   string `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	Title        string `json:"title"`
	Category     string `json:"category"`
	Content      string `json:"content"`
	CreatedAt    int64  `json:"created_at"`
}

type CreateAnnouncementRequest struct {
	CompanyID  string `json:"-" validate:"required,uuid4"`
	EmployeeID string `json:"-" validate:"required,uuid4"`
	Title      string `json:"title" validate:"required"`
	Category   string `json:"category" validate:"required"`
	Content    string `json:"content" validate:"required"`
}

type UpdateAnnouncementRequest struct {
	Title    *string `json:"title,omitempty"`
	Category *string `json:"category,omitempty"`
	Content  *string `json:"content,omitempty"`
}

type SearchAnnouncementRequest struct {
	CompanyID  string `json:"-" validate:"required,uuid4"`
	EmployeeID string `json:"employee_id,omitempty" validate:"omitempty,uuid4"`
	Title      string `json:"title,omitempty"`
	Category   string `json:"category,omitempty"`
	Page       int    `json:"page" validate:"min=1"`
	Size       int    `json:"size" validate:"min=1,max=100"`
	CreatedAt  int64  `json:"created_at,omitempty" validate:"omitempty"`
}

func NewAnnouncementResponse(a *entity.Announcement) *AnnouncementResponse {
	if a == nil {
		return nil
	}
	resp := &AnnouncementResponse{
		ID:           a.ID,
		CompanyID:    a.CompanyID,
		EmployeeID:   a.EmployeeID,
		EmployeeName: a.Employee.Fullname,
		Title:        a.Title,
		Category:     a.Category,
		Content:      a.Content,
		CreatedAt:    a.CreatedAt,
	}
	return resp
}

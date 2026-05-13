package model

import (
	"hr-sas/internal/entity"
	"strings"
)

type VisitResponse struct {
	ID           string                `json:"id"`
	EmployeeID   string                `json:"employee_id"`
	EmployeeName string                `json:"employee_name"`
	Date         string                `json:"date"`
	ClientName   string                `json:"client_name"`
	CreatedAt    int64                 `json:"created_at"`
	Details      []VisitDetailResponse `json:"details"`
}

type VisitDetailResponse struct {
	ID        string  `json:"id"`
	VisitType string  `json:"visit_type"`
	VisitAt   string  `json:"visit_at"`
	DateVisit string  `json:"date_visit"`
	FileUrl   *string `json:"file_url"`
	Latitude  *string `json:"latitude"`
	Longitude *string `json:"longitude"`
	Address   *string `json:"address"`
	Note      *string `json:"note"`
	VisitID   string  `json:"visit_id"`
	CreatedAt int64   `json:"created_at"`
}

type CreateVisitRequest struct {
	VisitType  string  `json:"visit_type" validate:"required,oneof=IN OUT"`
	ClientName string  `json:"client_name"`
	Note       *string `json:"note,omitempty"`
	Latitude   *string `json:"latitude,omitempty"`
	Longitude  *string `json:"longitude,omitempty"`
	Address    *string `json:"address,omitempty"`
	FileUrl    *string `json:"file_url,omitempty"`
}

type UpdateVisitRequest struct {
	ClientName *string `json:"client_name,omitempty"`
	Note       *string `json:"note,omitempty"`
	Latitude   *string `json:"latitude,omitempty"`
	Longitude  *string `json:"longitude,omitempty"`
	Address    *string `json:"address,omitempty"`
	FileUrl    *string `json:"file_url,omitempty"`
}

type SearchVisitRequest struct {
	EmployeeID string `json:"employee_id,omitempty" validate:"max=100"`
	StartDate  string `json:"start_date,omitempty" validate:"max=20"`
	EndDate    string `json:"end_date,omitempty" validate:"max=20"`
	SortBy     string `json:"sort_by,omitempty" validate:"max=10"`
	Page       int    `json:"page,omitempty" validate:"min=1"`
	Size       int    `json:"size,omitempty" validate:"min=1,max=100"`
}

type CanDoVisitResponse struct {
	CanDoVisit bool   `json:"can_do_visit"`
	Message    string `json:"message,omitempty"`
}

func VisitDetailToResponse(detail *entity.VisitDetail) VisitDetailResponse {
	var lat, lon *string
	if detail.Location != nil {
		parts := strings.SplitN(*detail.Location, ", ", 2)
		if len(parts) == 2 {
			lat = &parts[0]
			lon = &parts[1]
		}
	}
	return VisitDetailResponse{
		ID:        detail.ID,
		VisitType: detail.VisitType,
		VisitAt:   detail.VisitAt,
		DateVisit: detail.DateVisit,
		FileUrl:   detail.FileUrl,
		Latitude:  lat,
		Longitude: lon,
		Address:   detail.Address,
		Note:      detail.Note,
		VisitID:   detail.VisitID,
		CreatedAt: detail.CreatedAt,
	}
}

func VisitToResponse(visit *entity.Visit) *VisitResponse {
	if visit == nil {
		return nil
	}
	details := make([]VisitDetailResponse, len(visit.Details))
	for i := range visit.Details {
		details[i] = VisitDetailToResponse(&visit.Details[i])
	}
	return &VisitResponse{
		ID:           visit.ID,
		EmployeeID:   visit.EmployeeID,
		EmployeeName: visit.Employee.Fullname,
		Date:         visit.Date,
		ClientName:   visit.ClientName,
		CreatedAt:    visit.CreatedAt,
		Details:      details,
	}
}

package model

type VisitResponse struct {
	ID         string  `json:"id"`
	EmployeeID string  `json:"employee_id"`
	CompanyID  string  `json:"company_id"`
	VisitType  string  `json:"visit_type"`
	Note       *string `json:"note,omitempty"`
	Latitude   *string `json:"latitude,omitempty"`
	Longitude  *string `json:"longitude,omitempty"`
	Address    *string `json:"address,omitempty"`
	FileName   *string `json:"file_name,omitempty"`
	MimeType   *string `json:"mime_type,omitempty"`
	FileSize   *int    `json:"file_size,omitempty"`
	FileUrl    *string `json:"file_url,omitempty"`
	CreatedAt  int64   `json:"created_at"`
}

type CreateVisitRequest struct {
	VisitType string  `json:"visit_type" validate:"required,oneof=IN OUT"`
	Note      *string `json:"note,omitempty"`
	Latitude  *string `json:"latitude,omitempty"`
	Longitude *string `json:"longitude,omitempty"`
	Address   *string `json:"address,omitempty"`
	FileName  *string `json:"file_name,omitempty"`
	MimeType  *string `json:"mime_type,omitempty"`
	FileSize  *int    `json:"file_size,omitempty"`
	FileUrl   *string `json:"file_url,omitempty"`
}

type SearchVisitRequest struct {
	EmployeeID string `json:"employee_id,omitempty" validate:"max=100"`
	VisitType  string `json:"visit_type,omitempty" validate:"max=10"`
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

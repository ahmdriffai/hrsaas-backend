package model

import "hr-sas/internal/entity"

type EmployeeDocumentResponse struct {
	ID           string `json:"id"`
	EmployeeID   string `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	DocType      string `json:"doc_type"`
	DocNumber    string `json:"doc_number"`
	DocName      string `json:"doc_name"`
	FileURL      string `json:"file_url"`
	IssuedAt     int64  `json:"issued"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

type CreateEmployeeDocumentRequest struct {
	EmployeeID string `json:"employee_id" validate:"required"`
	DocType    string `json:"doc_type" validate:"required"`
	DocNumber  string `json:"doc_number" validate:"required"`
	DocName    string `json:"doc_name" validate:"required"`
	Issued     string `json:"issued" validate:"required"`
	FileURL    string `json:"file_url" validate:"required"`
}

type UpdateEmployeeDocumentRequest struct {
	ID        string  `json:"-" validate:"required"`
	DocType   *string `json:"doc_type,omitempty"`
	DocNumber *string `json:"doc_number,omitempty"`
	DocName   *string `json:"doc_name,omitempty"`
	Issued    *string `json:"issued,omitempty"`
	FileURL   *string `json:"file_url,omitempty"`
}

type SearchEmployeeDocumentRequest struct {
	EmployeeID string `json:"employee_id"`
	Page       int    `json:"page" validate:"min=1"`
	Size       int    `json:"size" validate:"min=1,max=100"`
}

func EmployeeDocumentToResponse(doc *entity.EmployeeDocument) *EmployeeDocumentResponse {
	if doc == nil {
		return nil
	}
	resp := &EmployeeDocumentResponse{
		ID:           doc.ID,
		EmployeeID:   doc.EmployeeID,
		EmployeeName: doc.Employee.Fullname,
		DocType:      doc.DocType,
		DocNumber:    doc.DocNumber,
		DocName:      doc.DocName,
		FileURL:      doc.FileURL,
		IssuedAt:     doc.IssuedAt,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
	return resp
}

func EmployeeDocumentsToResponse(docs []entity.EmployeeDocument) []EmployeeDocumentResponse {
	if len(docs) == 0 {
		return nil
	}
	responses := make([]EmployeeDocumentResponse, 0, len(docs))
	for i := range docs {
		if r := EmployeeDocumentToResponse(&docs[i]); r != nil {
			responses = append(responses, *r)
		}
	}
	return responses
}

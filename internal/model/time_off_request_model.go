package model

import (
	"hr-sas/internal/entity"
)

type ApprovalResponse struct {
	EmployeeName string `json:"employee_name"`
	IsRequired   bool   `json:"is_required"`
	Status       string `json:"status"`
	ActionAt     int64  `json:"action_at"`
	ActionReason string `json:"action_reason"`
}

type TimeOffRequestResponse struct {
	ID            string              `json:"id"`
	CompanyID     string              `json:"company_id"`
	EmployeeID    string              `json:"employee_id"`
	TimeOffTypeID string              `json:"time_off_type_id"`
	StartDate     int64               `json:"start_date"`
	EndDate       *int64              `json:"end_date"`
	RequestedDays int                 `json:"requested_days"`
	RequestReason *string             `json:"request_reason"`
	RequestStatus *string             `json:"request_status"`
	FileUrl       *string             `json:"file_url,omitempty"`
	CreatedAt     int64               `json:"created_at"`
	Employee      EmployeeResponse    `json:"employee"`
	TimeOffType   TimeOffTypeResponse `json:"time_off_type"`
	Approvals     []ApprovalResponse  `json:"approvals"`
}

type CreateTimeOffRequest struct {
	TimeOffTypeID string `json:"time_off_type_id" validate:"required"`
	StartDate     string `json:"start_date" validate:"required"`
	EndDate       string `json:"end_date" validate:"required"`
	RequestedDays int    `json:"requested_days"`
	RequestReason string `json:"request_reason" validate:"required,max=255"`
	RequestStatus string `json:"request_status" validate:"required,oneof=PENDING APPROVED REJECTED"`
	CreatedAt     int64  `json:"created_at"`
}

type SearchTimeOffRequest struct {
	EmployeeID    string `json:"employee_id" validate:"max=100"`
	TimeOffTypeID string `json:"time_off_type_id" validate:"max=100"`
	RequestStatus string `json:"request_status" validate:"max=20"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Page          int    `json:"page" validate:"min=1"`
	Size          int    `json:"size" validate:"min=1,max=100"`
}

type UpdateTimeOffRequest struct {
	TimeOffTypeID *string `json:"time_off_type_id,omitempty"`
	StartDate     *string `json:"start_date,omitempty"`
	EndDate       *string `json:"end_date,omitempty"`
	RequestedDays *int    `json:"requested_days,omitempty"`
	RequestReason *string `json:"request_reason,omitempty" validate:"max=255"`
	RequestStatus *string `json:"request_status,omitempty" validate:"omitempty,oneof=PENDING APPROVED REJECTED"`
}

func TimeOffRequestToResponse(request *entity.TimeOffRequest) *TimeOffRequestResponse {
	approvals := make([]ApprovalResponse, len(request.Approvals))
	for i, v := range request.Approvals {
		approval := ApprovalResponse{
			EmployeeName: v.Employee.Fullname,
			IsRequired:   v.IsRequired,
			Status:       v.Status,
			ActionAt:     v.ActionAt,
			ActionReason: v.ActionReason,
		}
		approvals[i] = approval
	}

	return &TimeOffRequestResponse{
		ID:            request.ID,
		CompanyID:     request.CompanyID,
		EmployeeID:    request.EmployeeID,
		TimeOffTypeID: request.TimeOffTypeId,
		StartDate:     request.StartDate,
		EndDate:       request.EndDate,
		RequestedDays: request.RequestedDays,
		RequestReason: request.RequestReason,
		RequestStatus: request.RequestStatus,
		FileUrl:       request.FileUrl,
		CreatedAt:     request.CreatedAt,
		Employee:      *EmployeeToResponse(&request.Employee),
		TimeOffType:   *TimeOffTypeToResponse(&request.TimeOffType),
		Approvals:     approvals,
	}
}

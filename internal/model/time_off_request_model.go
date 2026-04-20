package model

import "hr-sas/internal/entity"

type TimeOffRequestResponse struct {
	ID            string                    `json:"id"`
	EmployeeID    string                    `json:"employee_id"`
	TimeOffTypeID string                    `json:"time_off_type_id"`
	StartDate     int64                     `json:"start_date"`
	EndDate       *int64                    `json:"end_date"`
	RequestedDays int                       `json:"requested_days"`
	RequestReason *string                   `json:"request_reason"`
	RequestStatus *string                   `json:"request_status"`
	FileUrl       *string                   `json:"file_url,omitempty"`
	CreatedAt     int64                     `json:"created_at"`
	Employee      EmployeeResponse          `json:"employee"`
	TimeOffType   TimeOffTypeResponse       `json:"time_off_type"`
	Approvals     []TimeOffApprovalResponse `json:"approvals"`
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

func TimeOffRequestToResponse(request *entity.TimeOffRequest) *TimeOffRequestResponse {
	approvals := make([]TimeOffApprovalResponse, len(request.Approvals))
	for i := range request.Approvals {
		approvals[i] = *TimeOffApprovalToResponse(&request.Approvals[i])
	}

	return &TimeOffRequestResponse{
		ID:            request.ID,
		EmployeeID:    request.EmployeeId,
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

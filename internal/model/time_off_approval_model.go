package model

import "hr-sas/internal/entity"

type TimeOffApprovalResponse struct {
	ID                 string                  `json:"id"`
	TimeOffRequestID   string                  `json:"time_off_request_id"`
	ApproverEmployeeID string                  `json:"approver_employee_id"`
	IsRequired         bool                    `json:"is_required"`
	Status             string                  `json:"status"`
	ActionAt           int64                   `json:"action_at,omitempty"`
	ActionReason       *string                 `json:"action_reason,omitempty"`
	Employee           EmployeeResponse        `json:"employee"`
	TimeOffRequest     *TimeOffRequestResponse `json:"time_off_request,omitempty"`
}

type DecideTimeOffApprovalRequest struct {
	Action       string `json:"action" validate:"required,oneof=APPROVE REJECT"`
	ActionReason string `json:"action_reason" validate:"max=255"`
	ActionAt     int64  `json:"action_at" validate:"required"`
}

type SearchTimeOffApprovalRequest struct {
	Status string `json:"status" validate:"max=20"`
	Page   int    `json:"page" validate:"min=1"`
	Size   int    `json:"size" validate:"min=1,max=100"`
}

type TimeOffApprovalPolicyResponse struct {
	ID            string              `json:"id"`
	CompanyID     string              `json:"company_id"`
	TimeOffTypeID string              `json:"time_off_type_id"`
	DivisionID    string              `json:"division_id"`
	IsActive      bool                `json:"is_active"`
	CreatedAt     int64               `json:"created_at"`
	TimeOffType   TimeOffTypeResponse `json:"time_off_type"`
	Company       CompanyResponse     `json:"company"`
	Division      DivisionResponse    `json:"division"`
}

type CreateTimeOffApprovalPolicyRequest struct {
	TimeOffTypeID string           `json:"time_off_type_id" validate:"required"`
	DivisionID    string           `json:"division_id"`
	IsActive      bool             `json:"is_active"`
	Division      DivisionResponse `json:"division"`
}

type TimeOffApprovalPolicyStepResponse struct {
	ID                 string            `json:"id"`
	PolicyID           string            `json:"policy_id"`
	StepNo             int               `json:"step_no"`
	ApproverEmployeeID *string           `json:"approver_employee_id,omitempty"`
	ApproverPositionID *string           `json:"approver_position_id,omitempty"`
	Employee           *EmployeeResponse `json:"employee,omitempty"`
	Posisition         *PositionResponse `json:"position,omitempty"`
}

type CreateTimeOffApprovalPolicyStepRequest struct {
	StepNo             int               `json:"step_no" validate:"required,min=1"`
	ApproverEmployeeID *string           `json:"approver_employee_id,omitempty"`
	ApproverPositionID *string           `json:"approver_position_id,omitempty"`
	Employee           *EmployeeResponse `json:"employee,omitempty"`
	Position           *PositionResponse `json:"position,omitempty"`
}

func TimeOffApprovalToResponse(approval *entity.TimeOffApproval) *TimeOffApprovalResponse {
	timeoffRequest := TimeOffRequestToResponse(&approval.TimeOffRequest)
	return &TimeOffApprovalResponse{
		ID:                 approval.ID,
		TimeOffRequestID:   approval.TimeOffRequestId,
		ApproverEmployeeID: approval.ApproverId,
		IsRequired:         approval.IsRequired,
		Status:             approval.Status,
		ActionAt:           approval.ActionAt,
		ActionReason:       &approval.ActionReason,
		Employee:           *EmployeeToResponse(&approval.Employee),
		TimeOffRequest:     timeoffRequest,
	}
}

package model

import "time"

type TimeOffRequestResponse struct {
	ID            string    `json:"id"`
	EmployeeID    string    `json:"employee_id"`
	TimeOffTypeID string    `json:"time_off_type_id"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	RequestedDays float64   `json:"requested_days"`
	RequestReason string    `json:"request_reason"`
	RequestStatus string    `json:"request_status"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateTimeOffRequest struct {
	TimeOffTypeID string  `json:"time_off_type_id" validate:"required"`
	StartDate     string  `json:"start_date" validate:"required"`
	EndDate       string  `json:"end_date" validate:"required"`
	RequestedDays float64 `json:"requested_days" validate:"required,gt=0"`
	RequestReason string  `json:"request_reason" validate:"required,max=255"`
}

type SearchTimeOffRequest struct {
	TimeOffTypeID string `json:"time_off_type_id" validate:"max=100"`
	RequestStatus string `json:"request_status" validate:"max=20"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Page          int    `json:"page" validate:"min=1"`
	Size          int    `json:"size" validate:"min=1,max=100"`
}

type TimeOffTypeResponse struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Category         string  `json:"category"`
	IsQuotaBased     bool    `json:"is_quota_based"`
	DefaultQuotaDays float64 `json:"default_quota_days"`
}

type CreateTimeOffTypeRequest struct {
	Name             string  `json:"name" validate:"required,max=100"`
	Category         string  `json:"category" validate:"required,oneof=IZIN SAKIT CUTI"`
	IsQuotaBased     bool    `json:"is_quota_based"`
	DefaultQuotaDays float64 `json:"default_quota_days" validate:"min=0"`
}

type TimeOffBalanceResponse struct {
	ID            string  `json:"id"`
	EmployeeID    string  `json:"employee_id"`
	TimeOffTypeID string  `json:"time_off_type_id"`
	PeriodYear    int     `json:"period_year"`
	EntitledDays  float64 `json:"entitled_days"`
	UsedDays      float64 `json:"used_days"`
	RemainingDays float64 `json:"remaining_days"`
}

type SearchTimeOffBalanceRequest struct {
	TimeOffTypeID string `json:"time_off_type_id" validate:"max=100"`
	PeriodYear    int    `json:"period_year" validate:"min=2000"`
}

type TimeOffApprovalResponse struct {
	ID                 string     `json:"id"`
	TimeOffRequestID   string     `json:"time_off_request_id"`
	ApproverEmployeeID string     `json:"approver_employee_id"`
	ApproverName       string     `json:"approver_name"`
	ApproverPosition   string     `json:"approver_position"`
	ApproverDivision   string     `json:"approver_division"`
	Status             string     `json:"status"`
	ActionAt           *time.Time `json:"action_at,omitempty"`
	ActionReason       *string    `json:"action_reason,omitempty"`
}

type ApproveTimeOffRequest struct {
	ActionReason string `json:"action_reason" validate:"max=255"`
}

type RejectTimeOffRequest struct {
	ActionReason string `json:"action_reason" validate:"required,max=255"`
}

type TimeOffApprovalPolicyResponse struct {
	ID            string    `json:"id"`
	CompanyID     string    `json:"company_id"`
	TimeOffTypeID string    `json:"time_off_type_id"`
	DivisionID    *string   `json:"division_id,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateTimeOffApprovalPolicyRequest struct {
	TimeOffTypeID string  `json:"time_off_type_id" validate:"required"`
	DivisionID    *string `json:"division_id"`
	IsActive      bool    `json:"is_active"`
}

type TimeOffApprovalPolicyStepResponse struct {
	ID                 string  `json:"id"`
	PolicyID           string  `json:"policy_id"`
	StepNo             int     `json:"step_no"`
	ApproverEmployeeID *string `json:"approver_employee_id,omitempty"`
	ApproverPositionID *string `json:"approver_position_id,omitempty"`
}

type CreateTimeOffApprovalPolicyStepRequest struct {
	StepNo             int     `json:"step_no" validate:"required,min=1"`
	ApproverEmployeeID *string `json:"approver_employee_id,omitempty"`
	ApproverPositionID *string `json:"approver_position_id,omitempty"`
}

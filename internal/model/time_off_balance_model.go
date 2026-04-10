package model

import "hr-sas/internal/entity"

type TimeOffBalanceResponse struct {
	ID            string              `json:"id"`
	EmployeeID    string              `json:"employee_id"`
	TimeOffTypeID string              `json:"time_off_type_id"`
	PeriodYear    int                 `json:"period_year"`
	EntitledDays  int                 `json:"entitled_days"`
	UsedDays      int                 `json:"used_days"`
	RemainingDays int                 `json:"remaining_days"`
	Employee      EmployeeResponse    `json:"employee"`
	TimeOffType   TimeOffTypeResponse `json:"time_off_type"`
}

type SearchTimeOffBalanceRequest struct {
	TimeOffTypeID string `json:"time_off_type_id" validate:"max=100"`
	PeriodYear    int    `json:"period_year" validate:"min=2000"`
}

type SetTimeOffBalanceRequest struct {
	EmployeeID    string `json:"employee_id" validate:"required"`
	TimeOffTypeID string `json:"time_off_type_id" validate:"required"`
	PeriodYear    int    `json:"period_year" validate:"required,min=2000"`
	EntitledDays  int    `json:"entitled_days" validate:"min=0"`
	UsedDays      int    `json:"used_days" validate:"min=0"`
	RemainingDays *int   `json:"remaining_days,omitempty" validate:"omitempty,min=0"`
}

func TimeOffBalanceToResponse(balance *entity.TimeOffBalance) *TimeOffBalanceResponse {
	return &TimeOffBalanceResponse{
		ID:            balance.ID,
		EmployeeID:    balance.EmployeeId,
		TimeOffTypeID: balance.TimeOffTypeId,
		PeriodYear:    balance.PeriodYear,
		EntitledDays:  balance.EntitledDays,
		UsedDays:      balance.UsedDays,
		RemainingDays: balance.RemainingDays,
		Employee:      *EmployeeToResponse(&balance.Employee),
		TimeOffType:   *TimeOffTypeToResponse(&balance.TimeOffType),
	}
}

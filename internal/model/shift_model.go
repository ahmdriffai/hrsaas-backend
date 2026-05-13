package model

import (
	"hr-sas/internal/entity"
)

type ShiftResponse struct {
	ID            string             `json:"id"`
	CompanyID     string             `json:"company_id"`
	Name          string             `json:"name"`
	LateTolerance int                `json:"late_tolerance"`
	ShiftDays     []ShiftDayResponse `json:"shift_days,omitempty"`
	Employees     []EmployeeResponse `json:"employees,omitempty"`
	CreatedAt     int64              `json:"created_at"`
	UpdatedAt     int64              `json:"updated_at"`
}

type ShiftDayResponse struct {
	Weekday         int    `json:"weekday"`
	DayType         string `json:"day_type"`
	CheckIn         int64  `json:"check_in"`
	CheckOut        int64  `json:"check_out"`
	BreakStart      int64  `json:"break_start"`
	BreakEnd        int64  `json:"break_end"`
	MaxBreakMinutes int    `json:"max_break_minutes"`
}

type CreateShiftRequest struct {
	CompanyID        string            `json:"-"`
	Name             string            `json:"name" validate:"required"`
	LateTolerance    int               `json:"late_tolerance"`
	ShiftDayRequests []ShiftDayRequest `json:"shift_days" validate:"dive"`
}

type UpdateShiftRequest struct {
	Name             *string            `json:"name,omitempty"`
	LateTolerance    *int               `json:"late_tolerance,omitempty"`
	ShiftDayRequests *[]ShiftDayRequest `json:"shift_days,omitempty" validate:"omitempty,dive"`
}

type ShiftDayRequest struct {
	Weekday         int    `json:"weekday" validate:"required,min=1,max=7"`
	DayType         string `json:"day_type" validate:"required,oneof=workday offday"`
	CheckIn         string `json:"check_in"`
	CheckOut        string `json:"check_out"`
	BreakStart      string `json:"break_start"`
	BreakEnd        string `json:"break_end"`
	MaxBreakMinutes int    `json:"max_break_minutes" validate:"min=0"`
}

type AssignEmployeeToShiftRequest struct {
	CompanyID  string `json:"-" validate:"required"`
	EmployeeID string `json:"employee_id" validate:"required"`
	ShiftID    string `json:"shift_id" validate:"required"`
}

type SearchShiftRequest struct {
	CompanyID string `json:"-" validate:"required"`
	Key       string `json:"key" validate:"max=100"`
	Page      int    `json:"page" validate:"min=1"`
	Size      int    `json:"size" validate:"min=1,max=100"`
}

type DetailShifRequest struct {
	CompanyID string `json:"-" validate:"required"`
	ShiftID   string `json:"-" validate:"required"`
}

func ShiftDayToResponse(shifDays *entity.ShiftDay) *ShiftDayResponse {
	return &ShiftDayResponse{
		Weekday:         shifDays.Weekday,
		DayType:         shifDays.DayType,
		CheckIn:         shifDays.CheckIn,
		CheckOut:        shifDays.CheckOut,
		BreakStart:      shifDays.BreakStart,
		BreakEnd:        shifDays.BreakEnd,
		MaxBreakMinutes: shifDays.MaxBreakMinutes,
	}
}

// converter
func ShiftToResponse(shift *entity.Shift) *ShiftResponse {
	shifDays := make([]ShiftDayResponse, len(shift.ShiftDays))
	for i, shifDay := range shift.ShiftDays {
		shifDays[i] = *ShiftDayToResponse(&shifDay)
	}

	employees := make([]EmployeeResponse, len(shift.Employees))
	for i, employee := range shift.Employees {
		employees[i] = *EmployeeToResponse(&employee)
	}

	return &ShiftResponse{
		ID:            shift.ID,
		Name:          shift.Name,
		CompanyID:     shift.CompanyID,
		LateTolerance: shift.LateTolerance,
		Employees:     employees,
		ShiftDays:     shifDays,
		CreatedAt:     shift.CreatedAt,
		UpdatedAt:     shift.UpdatedAt,
	}
}

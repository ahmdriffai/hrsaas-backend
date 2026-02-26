package model

type ShiftResponse struct {
	ID            string `json:"id"`
	CompanyID     string `json:"company_id"`
	Name          string `json:"name"`
	LateTolerance int    `json:"late_tolerance"`
	ShiftDays     []ShiftDayResponse `json:"shift_days"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ShiftDayResponse struct {
	Weekday         int    `json:"weekday"`
	DayType         string `json:"day_type"`
	CheckIn         string `json:"check_in"`
	CheckOut        string `json:"check_out"`
	BreakStart      string `json:"break_start"`
	BreakEnd        string `json:"break_end"`
	MaxBreakMinutes int    `json:"max_break_minutes"`
}

type CreateShiftRequest struct {
	CompanyID        string            `json:"-"`
	Name             string            `json:"name" validate:"required"`
	LateTolerance    int               `json:"late_tolerance"`
	ShiftDayRequests []ShiftDayRequest `json:"shift_days" validate:"dive"`
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

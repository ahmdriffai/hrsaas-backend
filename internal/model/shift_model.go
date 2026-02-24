package model

import "time"

type ShiftResponse struct {
	ID            string    `json:"id"`
	CompanyID     string    `json:"company_id"`
	Name          string    `json:"name"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	LateTolerance int       `json:"late_tolerance"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

type CreateShiftRequest struct {
	CompanyID     string `json:"-"`
	Name          string `json:"name" validate:"required"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	LateTolerance int    `json:"late_tolerance"`
}

type AssignEmployeeToShiftRequest struct {
	CompanyID  string `json:"-" validate:"required"`
	EmployeeID string `json:"employee_id" validate:"required"`
	ShiftID    string `json:"shift_id" validate:"required"`
}

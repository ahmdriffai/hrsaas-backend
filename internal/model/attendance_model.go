package model

import "hr-sas/internal/entity"

type AttendanceResponse struct {
	ID                string `json:"id"`
	CompanyID         string `json:"company_id"`
	EmployeeID        string `json:"employee_id"`
	Date              string `json:"date"`
	CheckInTime       string `json:"check_in_time,omitempty"`
	CheckOutTime      string `json:"check_out_time,omitempty"`
	TotalWorkMinutes  int    `json:"total_work_minutes,omitempty"`
	TotalBreakMinutes int    `json:"total_break_minutes,omitempty"`
	Status            string `json:"status"`
}

type CheckInAttendanceRequest struct {
	CompanyID    string  `json:"-" validate:"required,uuid4"`
	EmployeeID   string  `json:"-" validate:"required,uuid4"`
	Lat          float64 `json:"lat" validate:"required"`
	Lng          float64 `json:"lng" validate:"required"`
	FaceImageUrl string  `json:"face_image_url" validate:"required,url"`
	DeviceInfo   string  `json:"device_info" validate:"required"`
}

// converter
func AttendandeToResponse(officeLocation *entity.Attendance) *AttendanceResponse {
	return &AttendanceResponse{
		ID:                officeLocation.ID,
		CompanyID:         officeLocation.CompanyID,
		EmployeeID:        officeLocation.EmployeeID,
		Date:              officeLocation.Date.String(),
		CheckInTime:       officeLocation.CheckInTime.String(),
		CheckOutTime:      officeLocation.CheckOutTime.String(),
		TotalWorkMinutes:  officeLocation.TotalWorkMinutes,
		TotalBreakMinutes: officeLocation.TotalBreakMinutes,
		Status:            officeLocation.Status,
	}
}

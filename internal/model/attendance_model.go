package model

import "hr-sas/internal/entity"

type AttendanceResponse struct {
	ID                string                  `json:"id"`
	CompanyID         string                  `json:"company_id"`
	EmployeeID        string                  `json:"employee_id"`
	EmployeeName      string                  `json:"employee_name"`
	EmployeePosition  string                  `json:"employee_position,omitempty"`
	Date              int64                   `json:"date"`
	CheckInTime       int64                   `json:"check_in_time,omitempty"`
	CheckOutTime      int64                   `json:"check_out_time,omitempty"`
	TotalWorkMinutes  int                     `json:"total_work_minutes,omitempty"`
	TotalBreakMinutes int                     `json:"total_break_minutes,omitempty"`
	Status            string                  `json:"status"`
	IsApproved        bool                    `json:"is_approved,omitempty"`
	Logs              []AttendanceLogResponse `json:"logs,omitempty"`
}

type AttendanceLogResponse struct {
	ID                 string  `json:"id"`
	AttendanceID       string  `json:"attendance_id"`
	Type               string  `json:"type"`
	Time               int64   `json:"time"`
	Lat                float64 `json:"lat"`
	Lng                float64 `json:"lng"`
	LocationDistance   float64 `json:"location_distance"`
	IsLocationVerified bool    `json:"is_location_verified"`
	IsFaceVerified     bool    `json:"is_face_verified"`
	FaceConfidence     float64 `json:"face_confidence"`
	FaceImageURL       string  `json:"face_image_url"`
	DeviceInfo         string  `json:"device_info"`
}

type CheckInAttendanceRequest struct {
	CompanyID    string  `json:"-" validate:"required,uuid4"`
	EmployeeID   string  `json:"-" validate:"required,uuid4"`
	Lat          float64 `json:"lat" validate:"required"`
	Lng          float64 `json:"lng" validate:"required"`
	FaceImageUrl string  `json:"face_image_url" validate:"required,url"`
	DeviceInfo   string  `json:"device_info" validate:"required"`
}

type UpdateAttendanceRequest struct {
	Date              *int64  `json:"date,omitempty"`
	CheckInTime       *int64  `json:"check_in_time,omitempty"`
	CheckOutTime      *int64  `json:"check_out_time,omitempty"`
	TotalWorkMinutes  *int    `json:"total_work_minutes,omitempty"`
	TotalBreakMinutes *int    `json:"total_break_minutes,omitempty"`
	IsApproved        *bool   `json:"is_approved,omitempty"`
	Status            *string `json:"status,omitempty" validate:"omitempty,oneof=HADIR TERLAMBAT ALPHA IZIN SAKIT"`
}

type SearchAttendanceRequest struct {
	CompanyID  string `json:"-" validate:"required,uuid4"`
	EmployeeID string `json:"employee_id,omitempty" validate:"omitempty,uuid4"`
	Date       string `json:"date,omitempty" validate:"omitempty,max=20"`
	Status     string `json:"status,omitempty" validate:"omitempty,oneof=HADIR TERLAMBAT ALPHA IZIN SAKIT"`
	IsApproved *bool  `json:"is_approved,omitempty"`
	Page       int    `json:"page,omitempty" validate:"min=1"`
	Size       int    `json:"size,omitempty" validate:"min=1,max=100"`
}

// converter
func AttendandeToResponse(attendance *entity.Attendance) *AttendanceResponse {
	employeePosition := ""
	for _, contract := range attendance.Employee.EmployeeContract {
		if contract.IsActive {
			employeePosition = contract.Position.Name
			break
		}
	}

	return &AttendanceResponse{
		ID:                attendance.ID,
		CompanyID:         attendance.CompanyID,
		EmployeeID:        attendance.EmployeeID,
		EmployeeName:      attendance.Employee.Fullname,
		EmployeePosition:  employeePosition,
		Date:              attendance.Date,
		CheckInTime:       attendance.CheckInTime,
		CheckOutTime:      attendance.CheckOutTime,
		TotalWorkMinutes:  attendance.TotalWorkMinutes,
		TotalBreakMinutes: attendance.TotalBreakMinutes,
		Status:            attendance.Status,
		IsApproved:        attendance.IsApproved,
	}
}

func AttendanceLogToResponse(log *entity.AttendanceLog) *AttendanceLogResponse {
	return &AttendanceLogResponse{
		ID:                 log.ID,
		AttendanceID:       log.AttendanceID,
		Type:               log.Type,
		Time:               log.Time,
		Lat:                log.Lat,
		Lng:                log.Lng,
		LocationDistance:   log.LocationDistance,
		IsLocationVerified: log.IsLocationVerified,
		IsFaceVerified:     log.IsFaceVerified,
		FaceConfidence:     log.FaceConfidence,
		FaceImageURL:       log.FaceImageURL,
		DeviceInfo:         log.DeviceInfo,
	}
}

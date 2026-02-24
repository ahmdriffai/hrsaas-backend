package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attendance struct {
	ID                string    `gorm:"column:id;primaryKey"`
	CompanyID         string    `gorm:"column:company_id;not null"`
	EmployeeID        string    `gorm:"column:employee_id;not null"`
	Date              time.Time `gorm:"column:date;not null"`
	CheckInTime       time.Time `gorm:"column:check_in_time"`
	CheckOutTime      time.Time `gorm:"column:check_out_time"`
	TotalWorkMinutes  int       `gorm:"column:total_work_minutes"`
	TotalBreakMinutes int       `gorm:"column:total_break_minutes"`
	Status            string    `gorm:"column:status;not null"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// BeforeCreate hook to set UUID
func (u *Attendance) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (Attendance) TableName() string {
	return "attendances"
}

type AttendanceLog struct {
	ID                 string    `gorm:"column:id;primaryKey"`
	AttendanceID       string    `gorm:"column:attendance_id;not null"`
	Type               string    `gorm:"column:type;not null"`
	Time               time.Time `gorm:"column:time;not null"`
	Lat                float64   `gorm:"column:lat"`               // Use float64 for latitude
	Lng                float64   `gorm:"column:lng"`               // Use float64 for longitude
	LocationDistance   float64   `gorm:"column:location_distance"` // Distance in meters
	IsLocationVerified bool      `gorm:"column:is_location_verified"`
	IsFaceVerified     bool      `gorm:"column:is_face_verified"`
	FaceConfidence     float64   `gorm:"column:face_confidence"` // Confidence level of face verification
	FaceImageURL       string    `gorm:"column:face_image_url"`  // URL to the face image
	DeviceInfo         string    `gorm:"column:device_info"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Foreign key constraint to attendances table
	Attendance Attendance `gorm:"-"`
}

// BeforeCreate hook to set UUID
func (u *AttendanceLog) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (AttendanceLog) TableName() string {
	return "attendance_logs"
}

package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Time_Off_Type struct {
	ID               string `gorm:"column:id;primaryKey"`
	Name             string `gorm:"column:name;not null"`
	Category         string `gorm:"column:category;not null"`
	IsQuotaBased     bool   `gorm:"column:is_quota_based;not null"`
	DefaultQuotaDays int    `gorm:"column:default_quota_days;not null"`
}

func (u *Time_Off_Type) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (c *Time_Off_Type) TableName() string {
	return "time_off_types"
}

type Time_Off_Requests struct {
	ID         string `gorm:"column:id;primaryKey"`
	EmployeeId string `gorm:"column:employee_id;not null"`
	// CompanyId     string `gorm:"column:company_id;not null"`
	TimeOffTypeId string  `gorm:"column:time_off_type_id;not null"`
	StartDate     int64   `gorm:"column:start_date;not null"`
	EndDate       int64   `gorm:"column:end_date;not null"`
	RequestReason *string `gorm:"column:request_reason"`
	RequestStatus string  `gorm:"column:request_status;not null"`
	CreatedAt     int64   `gorm:"column:created_at;not null"`
}

func (u *Time_Off_Requests) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (c *Time_Off_Requests) TableName() string {
	return "time_off_requests"
}

type Time_Off_Balance struct {
	ID            string `gorm:"column:id;primaryKey"`
	EmployeeId    string `gorm:"column:employee_id;not null"`
	TimeOffTypeId string `gorm:"column:time_off_type_id;not null"`
	PeriodYear    int    `gorm:"column:period_year;not null"`
	EntitledDays  int    `gorm:"column:entitled_days;not null"`
	UsedDays      int    `gorm:"column:used_days;not null"`
	RemainingDays int    `gorm:"column:remaining_days;not null"`
}

func (u *Time_Off_Balance) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (c *Time_Off_Balance) TableName() string {
	return "time_off_balances"
}

type Time_Off_Attachment struct {
	ID               string `gorm:"column:id;primaryKey"`
	TimeOffRequestId string `gorm:"column:time_off_request_id;not null"`
	FileName         string `gorm:"column:file_name;not null"`
	MimeType         string `gorm:"column:mime_type;not null"`
	FileSize         int    `gorm:"column:file_size;not null"`
	FileUrl          string `gorm:"column:file_url;not null"`
}

func (u *Time_Off_Attachment) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (c *Time_Off_Attachment) TableName() string {
	return "time_off_attachments"
}

type Time_Off_Approval struct {
	ID               string `gorm:"column:id;primaryKey"`
	TimeOffRequestId string `gorm:"column:time_off_request_id;not null"`
	ApproverId       string `gorm:"column:approver_id;not null"`
	Status           string `gorm:"column:approval_status;not null"`
	ActionReason     string `gorm:"column:action_reason"`
	ActionAt         int64  `gorm:"column:action_at"`
	// CreatedAt        int64   `gorm:"column:created_at;not null"`
}

func (u *Time_Off_Approval) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (c *Time_Off_Approval) TableName() string {
	return "time_off_approvals"
}

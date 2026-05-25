package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeOffType struct {
	ID               string `gorm:"column:id;primaryKey"`
	Name             string `gorm:"column:name;not null"`
	Category         string `gorm:"column:category;not null"`
	IsQuotaBased     bool   `gorm:"column:is_quota_based;not null"`
	DefaultQuotaDays int    `gorm:"column:default_quota_days;not null"`
}

func (u *TimeOffType) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (c *TimeOffType) TableName() string {
	return "time_off_types"
}

type TimeOffRequest struct {
	ID            string            `gorm:"column:id;primaryKey"`
	CompanyID     string            `gorm:"column:company_id;not null"`
	EmployeeID    string            `gorm:"column:employee_id;not null"`
	TimeOffTypeId string            `gorm:"column:time_off_type_id;not null"`
	RequestedDays int               `gorm:"column:requested_days;not null"`
	StartDate     int64             `gorm:"column:start_date;not null"`
	EndDate       *int64            `gorm:"column:end_date;not null"`
	RequestReason *string           `gorm:"column:request_reason"`
	RequestStatus *string           `gorm:"column:request_status;not null"`
	FileUrl       *string           `gorm:"column:file_url"`
	CreatedAt     int64             `gorm:"column:created_at;not null"`
	Employee      Employee          `gorm:"foreignKey:EmployeeID"`
	TimeOffType   TimeOffType       `gorm:"foreignKey:TimeOffTypeID"`
	Approvals     []TimeOffApproval `gorm:"foreignKey:TimeOffRequestId"`
}

func (u *TimeOffRequest) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	u.CreatedAt = int64(time.Now().UnixMilli())
	u.StartDate = int64(time.Now().UnixMilli())
	return nil
}

func (c *TimeOffRequest) TableName() string {
	return "time_off_requests"
}

type TimeOffBalance struct {
	ID            string      `gorm:"column:id;primaryKey"`
	CompanyID     string      `gorm:"column:company_id;not null"`
	EmployeeID    string      `gorm:"column:employee_id;not null"`
	TimeOffTypeId string      `gorm:"column:time_off_type_id;not null"`
	PeriodYear    int         `gorm:"column:period_year;not null"`
	EntitledDays  int         `gorm:"column:entitled_days;not null"`
	UsedDays      int         `gorm:"column:used_days;not null"`
	RemainingDays int         `gorm:"column:remaining_days;not null"`
	Employee      Employee    `gorm:"foreignKey:EmployeeId"`
	TimeOffType   TimeOffType `gorm:"foreignKey:TimeOffTypeId"`
}

func (u *TimeOffBalance) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (c *TimeOffBalance) TableName() string {
	return "time_off_balances"
}

type TimeOffAttachment struct {
	ID               string         `gorm:"column:id;primaryKey"`
	TimeOffRequestId string         `gorm:"column:time_off_request_id;not null"`
	FileName         string         `gorm:"column:file_name;not null"`
	MimeType         string         `gorm:"column:mime_type;not null"`
	FileSize         int            `gorm:"column:file_size;not null"`
	FileUrl          string         `gorm:"column:file_url;not null"`
	TimeOffRequest   TimeOffRequest `gorm:"foreignKey:TimeOffRequestId"`
}

func (u *TimeOffAttachment) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

func (c *TimeOffAttachment) TableName() string {
	return "time_off_attachments"
}

type TimeOffApproval struct {
	ID               string         `gorm:"column:id;primaryKey"`
	TimeOffRequestId string         `gorm:"column:time_off_request_id;not null"`
	ApproverId       string         `gorm:"column:approver_id;not null"`
	ApprovalOrder    int            `gorm:"column:approval_order;not null"`
	Status           string         `gorm:"column:approval_status;not null"`
	IsRequired       bool           `gorm:"column:is_required;not null"`
	ActionReason     string         `gorm:"column:action_reason"`
	ActionAt         int64          `gorm:"column:action_at"`
	Employee         Employee       `gorm:"foreignKey:ApproverId"`
	TimeOffRequest   TimeOffRequest `gorm:"foreignKey:TimeOffRequestId"`
}

func (u *TimeOffApproval) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	u.ActionAt = int64(time.Now().UnixMilli())
	return nil
}

func (c *TimeOffApproval) TableName() string {
	return "time_off_approvals"
}

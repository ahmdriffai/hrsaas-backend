package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RemidialVisit struct {
	ID                        string `gorm:"column:id;primaryKey"`
	CompanyID                 string `gorm:"column:company_id;not null"`
	EmployeeID                string `gorm:"column:employee_id;not null"`
	NasabahID                 string `gorm:"column:nasabah_id;not null"`
	NasabahName               string `gorm:"column:nasabah_name;not null"`
	NoPjm                     string `gorm:"column:no_pjm;not null"`
	LoanType                  string `gorm:"column:loan_type;not null"`
	Unit                      string `gorm:"column:unit;not null"`
	Collectibility            string `gorm:"column:collectibility;not null"`
	LoanLimit                 int64  `gorm:"column:loan_limit;not null"`
	OutstandingBalance        int64  `gorm:"column:outstanding_balance;not null"`
	OverduePrincipal          int64  `gorm:"column:overdue_principal;not null"`
	OverdueInterest           int64  `gorm:"column:overdue_interest;not null"`
	OverdueTotal              int64  `gorm:"column:overdue_total;not null"`
	OverduePrincipalFrequency int64  `gorm:"column:overdue_principal_frequency;not null"`
	OverdueInterestFrequency  int64  `gorm:"column:overdue_interest_frequency;not null"`
	OverduePrincipalDays      int64  `gorm:"column:overdue_principal_days;not null"`
	OverdueInterestDays       int64  `gorm:"column:overdue_interest_days;not null"`
	LoanStatus                string `gorm:"column:loan_status;not null"`
	TotalPaid                 int64  `gorm:"column:total_paid;not null"`
	Note                      string `gorm:"column:note"`
	CreatedAt                 int64  `gorm:"column:created_at;not null"`
}

func (r *RemidialVisit) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewString()
	r.CreatedAt = time.Now().UnixMilli()
	return nil
}

func (RemidialVisit) TableName() string {
	return "remidial_visit"
}

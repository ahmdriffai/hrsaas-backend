package model

type SearchNasabahRequest struct {
	BranchCode  string `json:"kodecabang"`
	CIF         string `json:"cif"`
	NamaNasabah string `json:"nama"`
	NIK         string `json:"nik"`
	LoanNumber  string `json:"nopjm"`
}

type ExternalNasabahResponse struct {
	RC   string        `json:"RC"`
	Data []NasabahData `json:"data"`
}

type NasabahData struct {
	KodeCabang         string `json:"KodeCabang"`
	NoPjm              string `json:"NoPjm"`
	NasabahID          string `json:"NasabahID"`
	Nama               string `json:"Nama"`
	NIK                string `json:"NIK"`
	TmpLahir           string `json:"TmpLahir"`
	TglLahir           string `json:"TglLahir"`
	Alamat             string `json:"Alamat"`
	Phone              string `json:"Phone"`
	Email              string `json:"Email"`
	JnsPjm             string `json:"JnsPjm"`
	Unit               string `json:"Unit"`
	Col                string `json:"Col"`
	Plafond            int64  `json:"Plafond"`
	Saldo              int64  `json:"Saldo"`
	TunggakanPokok     int64  `json:"TgkPokok"`
	TunggakanBunga     int64  `json:"TgkBunga"`
	TunggakanTotal     int64  `json:"TgkTotal"`
	TunggakanPokokFrek int    `json:"TgkPokokFrek"`
	TunggakanBungaFrek int    `json:"TgkBungaFrek"`
	TunggakanPokokHari int    `json:"TgkPokokHari"`
	TunggakanBungaHari int    `json:"TgkBungaHAri"`
	StatusPinjaman     string `json:"StsPjm"`
}

type RemidialVisitResponse struct {
	ID                        string `json:"id"`
	CompanyID                 string `json:"company_id"`
	EmployeeID                string `json:"employee_id"`
	EmployeeName              string `json:"employee_name,omitempty"`
	NasabahID                 string `json:"nasabah_id"`
	NasabahName               string `json:"nasabah_name,omitempty"`
	NoPjm                     string `json:"no_pjm"`
	LoanType                  string `json:"loan_type"`
	Unit                      string `json:"unit"`
	Collectibility            string `json:"collectibility"`
	LoanLimit                 int64  `json:"loan_limit"`
	OutstandingBalance        int64  `json:"outstanding_balance"`
	OverduePrincipal          int64  `json:"overdue_principal"`
	OverdueInterest           int64  `json:"overdue_interest"`
	OverdueTotal              int64  `json:"overdue_total"`
	OverduePrincipalFrequency int64  `json:"overdue_principal_frequency"`
	OverdueInterestFrequency  int64  `json:"overdue_interest_frequency"`
	OverduePrincipalDays      int64  `json:"overdue_principal_days"`
	OverdueInterestDays       int64  `json:"overdue_interest_days"`
	LoanStatus                string `json:"loan_status"`
	TotalPaid                 int64  `json:"total_paid"`
	Note                      string `json:"note"`
	CreatedAt                 int64  `json:"created_at"`
}

type CreateRemidialVisitRequest struct {
	CompanyID                 string `json:"-" form:"-" validate:"required,uuid4"`
	EmployeeID                string `json:"-" form:"-" validate:"required,uuid4"`
	NasabahID                 string `json:"nasabah_id" validate:"required"`
	NasabahName               string `json:"nasabah_name"`
	NoPjm                     string `json:"no_pjm"`
	LoanType                  string `json:"loan_type"`
	Unit                      string `json:"unit"`
	Collectibility            string `json:"collectibility"`
	LoanLimit                 int64  `json:"loan_limit"`
	OutstandingBalance        int64  `json:"outstanding_balance"`
	OverduePrincipal          int64  `json:"overdue_principal"`
	OverdueInterest           int64  `json:"overdue_interest"`
	OverdueTotal              int64  `json:"overdue_total"`
	OverduePrincipalFrequency int64  `json:"overdue_principal_frequency"`
	OverdueInterestFrequency  int64  `json:"overdue_interest_frequency"`
	OverduePrincipalDays      int64  `json:"overdue_principal_days"`
	OverdueInterestDays       int64  `json:"overdue_interest_days"`
	LoanStatus                string `json:"loan_status"`
	TotalPaid                 int64  `json:"total_paid"`
	Note                      string `json:"note"`
	CreatedAt                 int64  `json:"created_at"`
}

type SearchRemidialVisitRequest struct {
	EmployeeID   string `json:"employee_id,omitempty" validate:"omitempty,uuid4"`
	EmployeeName string `json:"employee_name,omitempty" validate:"max=100"`
	NasabahName  string `json:"nama,omitempty" validate:"max=100"`
	StartDate    string `json:"start_date,omitempty" validate:"max=20"`
	EndDate      string `json:"end_date,omitempty" validate:"max=20"`
	Page         int    `json:"page,omitempty" validate:"min=1"`
	Size         int    `json:"size,omitempty" validate:"min=1,max=100"`
}

type UpdateRemidialVisitRequest struct {
	TotalPaid int64  `json:"total_paid" validate:"required"`
	Note      string `json:"note,omitempty"`
}

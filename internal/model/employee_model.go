package model

type EmployeeResponse struct {
	ID             string  `json:"id,omitempty"`
	CompanyID      string  `json:"company_id,omitempty"`
	UserID         *string `json:"user_id,omitempty"`
	EmployeeNumber string  `json:"employee_number,omitempty"`
	Fullname       string  `json:"fullname,omitempty"`
	BirthPlace     string  `json:"birth_place,omitempty"`
	BirthDate      string  `json:"birth_date,omitempty"`
	BlodType       string  `json:"blood_type,omitempty"`
	MaritalStatus  string  `json:"marital_status,omitempty"`
	Religion       string  `json:"religion,omitempty"`
	Phone          string  `json:"phone,omitempty"`
	Timezone       string  `json:"timezone,omitempty"`
}

type CreateEmployeeRequest struct {
	CompanyID      string `json:"-" validate:"required"`
	Fullname       string `json:"fullname" validate:"required"`
	EmployeeNumber string `json:"employee_number" validate:"required"`
	BirthPlace     string `json:"birth_place" validate:"required"`
	BirthDate      string `json:"birth_date" validate:"required"`
	BlodType       string `json:"blood_type" validate:"required"`
	MaritalStatus  string `json:"marital_status" validate:"required"`
	Religion       string `json:"religion" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	Timezone       string `json:"timezone" validate:"required"`
}

type SearchEmployeeRequest struct {
	CompanyID string `json:"company_id" validate:"required"`
	Key       string `json:"key" validate:"max=100"`
	Page      int    `json:"page" validate:"min=1"`
	Size      int    `json:"size" validate:"min=1,max=100"`
}

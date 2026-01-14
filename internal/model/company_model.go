package model

import "time"

type CompanyResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	LogoUrl        *string   `json:"logo_url"`
	BussinessField *string   `json:"bussiness_field"`
	Address        *string   `json:"address"`
	Province       *string   `json:"province"`
	City           *string   `json:"city"`
	District       *string   `json:"district"`
	Village        *string   `json:"village"`
	ZipCode        *string   `json:"zip_code"`
	PhoneNumber    *string   `json:"phone_number"`
	FaxNumber      *string   `json:"fax_number"`
	Email          *string   `json:"email"`
	Website        *string   `json:"website"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateCompanyRequest struct {
	Name           string `json:"name" validate:"required"`
	LogoUrl        string `json:"logo_url,omitempty"`
	BussinessField string `json:"bussiness_field,omitempty"`
	Address        string `json:"address,omitempty"`
	Province       string `json:"province,omitempty"`
	City           string `json:"city,omitempty"`
	District       string `json:"district,omitempty"`
	Village        string `json:"village,omitempty"`
	ZipCode        string `json:"zip_code,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	FaxNumber      string `json:"fax_number,omitempty"`
	Email          string `json:"email,omitempty"`
	Website        string `json:"website,omitempty"`
}

type RegisterCompanyRequest struct {
	UserID         string  `json:"-" validate:"required"`
	Name           string  `json:"name" validate:"required"`
	LogoUrl        *string `json:"logo_url,omitempty"`
	BussinessField *string `json:"bussiness_field,omitempty"`
	Address        *string `json:"address,omitempty"`
	Province       *string `json:"province,omitempty"`
	City           *string `json:"city,omitempty"`
	District       *string `json:"district,omitempty"`
	Village        *string `json:"village,omitempty"`
	ZipCode        *string `json:"zip_code,omitempty"`
	PhoneNumber    *string `json:"phone_number,omitempty"`
	FaxNumber      *string `json:"fax_number,omitempty"`
	Email          *string `json:"email,omitempty"`
	Website        *string `json:"website,omitempty"`
}

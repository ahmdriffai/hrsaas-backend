package model

import (
	"time"
)

type UserResponse struct {
	ID            string            `json:"id,omitempty"`
	Name          string            `json:"name,omitempty"`
	Email         string            `json:"email,omitempty"`
	EmailVerified bool              `json:"email_verified,omitempty"`
	Image         *string           `json:"image,omitempty"`
	CompanyID     string            `json:"company_id,omitempty"`
	Role          string            `json:"role,omitempty"`
	Employee      *EmployeeResponse `json:"employee,omitempty"`
	CreatedAt     time.Time         `json:"created_at,omitempty"`
	UpdatedAt     time.Time         `json:"updated_at,omitempty"`
}

type LoginUserResponse struct {
	User  UserResponse `json:"user,omitempty"`
	Token string       `json:"token,omitempty"`
}

type VerifyUserRequest struct {
	Token string `validate:"required"`
}

type RegisterUserRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	CompanyName string `json:"company_name" validate:"required"`
}

type LoginUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Ip        string `json:"-"`
	UserAgent string `json:"-"`
}

package model

import (
	"hr-sas/internal/entity"
)

type UserResponse struct {
	ID            string               `json:"id,omitempty"`
	Name          string               `json:"name,omitempty"`
	Email         string               `json:"email,omitempty"`
	EmailVerified bool                 `json:"email_verified,omitempty"`
	Image         *string              `json:"image,omitempty"`
	CompanyID     string               `json:"company_id,omitempty"`
	Roles         []RoleResponse       `json:"roles,omitempty"`
	Permissions   []PermissionResponse `json:"permissions,omitempty"`
	Employee      *EmployeeResponse    `json:"employee,omitempty"`
	CreatedAt     int64                `json:"created_at,omitempty"`
	UpdatedAt     int64                `json:"updated_at,omitempty"`
}

type LoginUserResponse struct {
	User  UserResponse `json:"user,omitempty"`
	Token string       `json:"token,omitempty"`
}

type SearchUserRequest struct {
	Key  string `json:"key" validate:"max=100"`
	Page int    `json:"page" validate:"min=1"`
	Size int    `json:"size" validate:"min=1,max=100"`
}

type UpdateUserRequest struct {
	ID            string    `json:"-"`
	Name          *string   `json:"name,omitempty"`
	Email         *string   `json:"email,omitempty" validate:"omitempty,email"`
	Image         *string   `json:"image,omitempty"`
	EmailVerified *bool     `json:"email_verified,omitempty"`
	RoleIDs       *[]string `json:"role_ids,omitempty"`
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

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type AssignRoleRequest struct {
	UserID string   `json:"-"`
	Roles  []string `json:"roles" validate:"required"`
}

type RemoveRoleRequest struct {
	UserID string   `json:"-"`
	Roles  []string `json:"roles" validate:"required"`
}

func UserToResponse(user *entity.User) *UserResponse {
	if user == nil || user.ID == "" {
		return nil
	}

	var employeeResponse *EmployeeResponse
	var roles []RoleResponse
	var permissions []PermissionResponse

	if user.Employee != nil {
		employeeResponse = EmployeeToResponse(user.Employee)
	}
	for _, r := range user.Roles {
		roles = append(roles, RoleResponse{
			Name: r.Name,
		})
		for _, p := range r.Permissions {
			permissions = append(permissions, PermissionResponse{
				Name: p.Name,
			})
		}
	}

	return &UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Image:         user.Image,
		Roles:         roles,
		Permissions:   permissions,
		CompanyID:     user.CompanyID,
		EmailVerified: user.EmailVerified,
		Employee:      employeeResponse,
	}
}

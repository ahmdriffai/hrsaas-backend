package model

import "hr-sas/internal/entity"

type RoleResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   int64                `json:"created_at,omitempty"`
	UpdatedAt   int64                `json:"updated_at,omitempty"`
}

type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required"`
	Permissions []string `json:"permissions"`
}

type UpdateRoleRequest struct {
	ID   string `json:"-"`
	Name string `json:"name" validate:"required"`
}

type SearchRoleRequest struct {
	Name string `json:"name"`
	Page int    `json:"page" validate:"min=1"`
	Size int    `json:"size" validate:"min=1,max=100"`
}

type AssignPermissionRequest struct {
	RoleID      string   `json:"-"`
	Permissions []string `json:"permissions" validate:"required"`
}

func RoleToResponse(role *entity.Role) *RoleResponse {
	if role == nil {
		return nil
	}

	permissions := make([]PermissionResponse, len(role.Permissions))
	for i, permission := range role.Permissions {
		permissions[i] = PermissionResponse{
			ID:   permission.ID,
			Name: permission.Name,
		}
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

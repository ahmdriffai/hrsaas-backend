package model

type PermissionResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type CreatePermissionRequest struct {
	Name  string   `json:"name" validate:"required"`
	Roles []string `json:"roles"`
}

type UpdatePermissionRequest struct {
	ID   string `json:"-"`
	Name string `json:"name" validate:"required"`
}

type SearchPermissionRequest struct {
	Name string `json:"name"`
	Page int    `json:"page" validate:"min=1"`
	Size int    `json:"size" validate:"min=1,max=100"`
}

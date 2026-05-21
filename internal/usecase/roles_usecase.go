package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	RoleRepository       *repository.RoleRepository
	PermissionRepository *repository.PermissionRepository
}

func NewRoleUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	roleRepository *repository.RoleRepository,
	permissionRepository *repository.PermissionRepository,
) *RoleUseCase {
	return &RoleUseCase{
		DB:                   db,
		Log:                  log,
		Validate:             validate,
		RoleRepository:       roleRepository,
		PermissionRepository: permissionRepository,
	}
}

func (u *RoleUseCase) Create(ctx context.Context, request *model.CreateRoleRequest) (*model.RoleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	role := &entity.Role{Name: request.Name}

	if err := u.RoleRepository.Create(tx, role); err != nil {
		u.Log.WithError(err).Error("Failed to create role")
		return nil, fiber.ErrInternalServerError
	}

	if len(request.Permissions) > 0 {
		permissions, err := u.PermissionRepository.FindByIDs(tx, request.Permissions)
		if err != nil {
			u.Log.WithError(err).Error("Failed to find permissions")
			return nil, fiber.ErrInternalServerError
		}

		if _, err := u.RoleRepository.AssignPermissions(ctx, tx, role, permissions); err != nil {
			u.Log.WithError(err).Error("Failed to assign permissions to role")
			return nil, fiber.ErrInternalServerError
		}
		role.Permissions = permissions
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.RoleToResponse(role), nil
}

func (u *RoleUseCase) List(ctx context.Context, request *model.SearchRoleRequest) ([]model.RoleResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("Failed to validate request body")
		return nil, 0, fiber.ErrBadRequest
	}

	roles, total, err := u.RoleRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Error("Failed to search roles")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	var responses []model.RoleResponse
	for _, r := range roles {
		responses = append(responses, *model.RoleToResponse(&r))
	}

	return responses, total, nil
}

func (u *RoleUseCase) Detail(ctx context.Context, id string) (*model.RoleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var role entity.Role
	if err := u.RoleRepository.FindById(tx, &role, id, "Permissions"); err != nil {
		u.Log.WithError(err).Error("Role not found")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.RoleToResponse(&role), nil
}

func (u *RoleUseCase) Update(ctx context.Context, request *model.UpdateRoleRequest) (*model.RoleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	var role entity.Role
	if err := u.RoleRepository.FindById(tx, &role, request.ID); err != nil {
		u.Log.WithError(err).Error("Role not found")
		return nil, fiber.ErrNotFound
	}

	role.Name = request.Name
	role.UpdatedAt = time.Now().UnixMilli()

	if err := u.RoleRepository.Update(tx, &role); err != nil {
		u.Log.WithError(err).Error("Failed to update role")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.RoleToResponse(&role), nil
}

func (u *RoleUseCase) Delete(ctx context.Context, id string) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var role entity.Role
	if err := u.RoleRepository.FindById(tx, &role, id); err != nil {
		u.Log.WithError(err).Error("Role not found")
		return fiber.ErrNotFound
	}

	if err := u.RoleRepository.Delete(tx, &role); err != nil {
		u.Log.WithError(err).Error("Failed to delete role")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (u *RoleUseCase) AssignPermissions(ctx context.Context, request *model.AssignPermissionRequest) (*model.RoleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	permissions, err := u.PermissionRepository.FindByIDs(tx, request.Permissions)
	if err != nil {
		u.Log.WithError(err).Error("Failed to fetch permissions")
		return nil, fiber.ErrInternalServerError
	}

	if len(permissions) != len(request.Permissions) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "some permissions not found")
	}

	var role entity.Role
	if err := u.RoleRepository.FindById(tx, &role, request.RoleID, "Permissions"); err != nil {
		u.Log.WithError(err).Error("Role not found")
		return nil, fiber.ErrNotFound
	}

	if _, err := u.RoleRepository.AssignPermissions(ctx, tx, &role, permissions); err != nil {
		u.Log.WithError(err).Error("Failed to assign permissions to role")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.RoleToResponse(&role), nil
}

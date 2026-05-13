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

type PermissionUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	PermissionRepository *repository.PermissionRepository
}

func NewPermissionUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	permissionRepository *repository.PermissionRepository,
) *PermissionUseCase {
	return &PermissionUseCase{
		DB:                   db,
		Log:                  log,
		Validate:             validate,
		PermissionRepository: permissionRepository,
	}
}

func (u *PermissionUseCase) Create(ctx context.Context, request *model.CreatePermissionRequest) (*model.PermissionResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	permission := &entity.Permission{Name: request.Name}

	if err := u.PermissionRepository.Create(tx, permission); err != nil {
		u.Log.WithError(err).Error("Failed to create permission")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return permissionToResponse(permission), nil
}

func (u *PermissionUseCase) List(ctx context.Context, request *model.SearchPermissionRequest) ([]model.PermissionResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("Failed to validate request body")
		return nil, 0, fiber.ErrBadRequest
	}

	permissions, total, err := u.PermissionRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Error("Failed to search permissions")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	var responses []model.PermissionResponse
	for _, p := range permissions {
		responses = append(responses, *permissionToResponse(&p))
	}

	return responses, total, nil
}

func (u *PermissionUseCase) Detail(ctx context.Context, id string) (*model.PermissionResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var permission entity.Permission
	if err := u.PermissionRepository.FindById(tx, &permission, id); err != nil {
		u.Log.WithError(err).Error("Permission not found")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return permissionToResponse(&permission), nil
}

func (u *PermissionUseCase) Update(ctx context.Context, request *model.UpdatePermissionRequest) (*model.PermissionResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	var permission entity.Permission
	if err := u.PermissionRepository.FindById(tx, &permission, request.ID); err != nil {
		u.Log.WithError(err).Error("Permission not found")
		return nil, fiber.ErrNotFound
	}

	permission.Name = request.Name
	permission.UpdatedAt = time.Now().UnixMilli()

	if err := u.PermissionRepository.Update(tx, &permission); err != nil {
		u.Log.WithError(err).Error("Failed to update permission")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return permissionToResponse(&permission), nil
}

func (u *PermissionUseCase) Delete(ctx context.Context, id string) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var permission entity.Permission
	if err := u.PermissionRepository.FindById(tx, &permission, id); err != nil {
		u.Log.WithError(err).Error("Permission not found")
		return fiber.ErrNotFound
	}

	if err := u.PermissionRepository.Delete(tx, &permission); err != nil {
		u.Log.WithError(err).Error("Failed to delete permission")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

func permissionToResponse(p *entity.Permission) *model.PermissionResponse {
	return &model.PermissionResponse{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

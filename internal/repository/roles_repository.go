package repository

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleRepository struct {
	Repository[entity.Role]
	Log *logrus.Logger
}

func NewRoleRepository(log *logrus.Logger) *RoleRepository {
	return &RoleRepository{Log: log}
}

func (r *RoleRepository) FindByName(db *gorm.DB, name string) (*entity.Role, error) {
	var result entity.Role
	if err := db.Where("name = ?", name).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *RoleRepository) FindByNames(db *gorm.DB, names []string) ([]entity.Role, error) {
	var result []entity.Role
	if err := db.Where("name IN ?", names).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RoleRepository) FindByIDs(db *gorm.DB, ids []string) ([]entity.Role, error) {
	var roles []entity.Role
	if err := db.Where("id IN ?", ids).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) AssignPermissions(ctx context.Context, db *gorm.DB, role *entity.Role, permissions []entity.Permission) (*entity.Role, error) {
	err := db.WithContext(ctx).Omit("Permissions.*").Model(role).Association("Permissions").Replace(permissions)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) FindAll(db *gorm.DB) ([]entity.Role, error) {
	var result []entity.Role
	if err := db.Preload("Permissions").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RoleRepository) Search(db *gorm.DB, request *model.SearchRoleRequest) ([]entity.Role, int64, error) {
	var roles []entity.Role

	query := db.Preload("Permissions")
	if key := request.Name; key != "" {
		query = query.Where("name LIKE ?", "%"+key+"%")
	}
	if err := query.Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := db.Model(&entity.Role{})
	if key := request.Name; key != "" {
		countQuery = countQuery.Where("name LIKE ?", "%"+key+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

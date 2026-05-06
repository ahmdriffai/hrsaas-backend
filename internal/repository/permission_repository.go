package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	Repository[entity.Permission]
	Log *logrus.Logger
}

func NewPermissionRepository(log *logrus.Logger) *PermissionRepository {
	return &PermissionRepository{Log: log}
}

func (r *PermissionRepository) FindByName(db *gorm.DB, name string) (*entity.Permission, error) {
	var result entity.Permission
	if err := db.Where("name = ?", name).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *PermissionRepository) FindByNames(db *gorm.DB, names []string) ([]entity.Permission, error) {
	var result []entity.Permission
	if err := db.Where("name IN ?", names).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PermissionRepository) FindByIDs(db *gorm.DB, ids []string) ([]entity.Permission, error) {
	var permissions []entity.Permission
	if err := db.Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *PermissionRepository) FindAll(db *gorm.DB) ([]entity.Permission, error) {
	var result []entity.Permission
	if err := db.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PermissionRepository) Search(db *gorm.DB, request *model.SearchPermissionRequest) ([]entity.Permission, int64, error) {
	var permissions []entity.Permission

	query := db
	if key := request.Name; key != "" {
		query = query.Where("name LIKE ?", "%"+key+"%")
	}
	if err := query.Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := db.Model(&entity.Permission{})
	if key := request.Name; key != "" {
		countQuery = countQuery.Where("name LIKE ?", "%"+key+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

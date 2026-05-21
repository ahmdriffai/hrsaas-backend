package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) CountByEmail(db *gorm.DB, email string) (int64, error) {
	var total int64
	err := db.Model(new(entity.User)).Where("email = ?", email).Count(&total).Error
	return total, err
}

func (r *UserRepository) CountByEmailExcludeID(db *gorm.DB, email string, userID string) (int64, error) {
	var total int64
	err := db.Model(new(entity.User)).
		Where("email = ?", email).
		Where("id <> ?", userID).
		Count(&total).Error
	return total, err
}

func (r *UserRepository) FindByEmail(db *gorm.DB, entity *entity.User, email string, preloads ...string) error {
	query := db

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	return query.Where("email = ?", email).Take(entity).Error
}

func (r *UserRepository) Search(db *gorm.DB, request *model.SearchUserRequest) ([]entity.User, int64, error) {
	var users []entity.User

	query := db.Model(&entity.User{}).Preload("Roles").Preload("Employee")
	if request.Key != "" {
		key := "%" + request.Key + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", key, key)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) AssignRoles(db *gorm.DB, user *entity.User, roles []entity.Role) error {
	return db.Omit("Roles.*").Model(user).Association("Roles").Replace(roles)
}

func (r *UserRepository) RemoveRoles(db *gorm.DB, user *entity.User, roles []entity.Role) error {
	return db.Model(user).Association("Roles").Delete(roles)
}

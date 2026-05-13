package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PositionRepository struct {
	Repository[entity.Position]
	Log *logrus.Logger
}

func NewPositionRepository(log *logrus.Logger) *PositionRepository {
	return &PositionRepository{
		Log: log,
	}
}

/* Find All By Company
 */
func (c *PositionRepository) FindAllByCompany(tx *gorm.DB, companyID string) ([]entity.Position, error) {
	var positions []entity.Position
	if err := tx.Where("company_id = ?", companyID).Find(&positions).Error; err != nil {
		return nil, err
	}
	return positions, nil
}

func (r *PositionRepository) Search(db *gorm.DB, request *model.SearchPositionRequest) ([]entity.Position, int64, error) {
	var positions []entity.Position
	if err := db.Scopes(r.FilterSearch(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&positions).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Position{}).Scopes(r.FilterSearch(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return positions, total, nil
}

func (r *PositionRepository) FilterSearch(request *model.SearchPositionRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("company_id = ?", request.CompanyID)

		if key := request.Name; key != "" {
			key = "%" + key + "%"
			tx = tx.Where("name LIKE ?", key).Or("description LIKE ?", key)
		}

		return tx
	}
}

func (r *PositionRepository) FindByName(db *gorm.DB, entity *entity.Position, name string, preloads ...string) error {
	query := db

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	return query.Where("name = ?", name).Take(entity).Error
}

func (r *PositionRepository) FindByID(db *gorm.DB, id string, preloads ...string) (*entity.Position, error) {
	var item entity.Position
	query := db

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Where("id = ?", id).Take(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

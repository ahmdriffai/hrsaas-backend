package repository

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OfficeLocationRepository struct {
	Repository[entity.OfficeLocation]
	Log *logrus.Logger
}

func NewOfficeLocationRepository(log *logrus.Logger) *OfficeLocationRepository {
	return &OfficeLocationRepository{
		Log: log,
	}
}

func (r *OfficeLocationRepository) Search(db *gorm.DB, request *model.SearchOfficeLocationRequest) ([]entity.OfficeLocation, int64, error) {
	var officeLocations []entity.OfficeLocation
	if err := db.Scopes(r.FilterSearch(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&officeLocations).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.OfficeLocation{}).Scopes(r.FilterSearch(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return officeLocations, total, nil
}

func (r *OfficeLocationRepository) FilterSearch(request *model.SearchOfficeLocationRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("company_id = ?", request.CompanyID)

		if key := request.Key; key != "" {
			key = "%" + key + "%"
			tx = tx.Where("name LIKE ?", key)
		}

		return tx
	}
}

func (r *OfficeLocationRepository) GetByEmployeeID(db *gorm.DB, employeeId string) ([]entity.OfficeLocation, error) {
	var officeLocations []entity.OfficeLocation
	if err := db.Where("id IN (SELECT office_location_id FROM employee_office_locations WHERE employee_id = ?)", employeeId).Find(&officeLocations).Error; err != nil {
		return nil, err
	}
	return officeLocations, nil
}

func (r *OfficeLocationRepository) CountByIDAndCompanyID(db *gorm.DB, officeLocationID, companyID string) (int64, error) {
	var total int64
	err := db.Model(&entity.OfficeLocation{}).
		Where("id = ?", officeLocationID).
		Where("company_id = ?", companyID).
		Count(&total).Error
	return total, err
}

func (r *OfficeLocationRepository) CountEmployeeByIDAndCompanyID(db *gorm.DB, employeeID, companyID string) (int64, error) {
	var total int64
	err := db.Model(&entity.Employee{}).
		Where("id = ?", employeeID).
		Where("company_id = ?", companyID).
		Count(&total).Error
	return total, err
}

func (r *OfficeLocationRepository) DeleteEmployeeOfficeLocationsByEmployeeID(db *gorm.DB, employeeID string) error {
	return db.Exec("DELETE FROM employee_office_locations WHERE employee_id = ?", employeeID).Error
}

func (r *OfficeLocationRepository) AssignEmployeeToOfficeLocation(db *gorm.DB, employeeID, officeLocationID string) error {
	return db.Exec(
		"INSERT INTO employee_office_locations (employee_id, office_location_id) VALUES (?, ?)",
		employeeID,
		officeLocationID,
	).Error
}

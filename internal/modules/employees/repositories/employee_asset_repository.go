package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type EmployeeAssetRepository struct {
	db *gorm.DB
}

func NewEmployeeAssetRepository() *EmployeeAssetRepository {
	return &EmployeeAssetRepository{
		db: database.GetDB(),
	}
}

// Create creates a new asset record
func (r *EmployeeAssetRepository) Create(asset *models.EmployeeAsset) error {
	return r.db.Create(asset).Error
}

// FindByDraftID finds all assets by draft ID
func (r *EmployeeAssetRepository) FindByDraftID(draftID uint) ([]models.EmployeeAsset, error) {
	var assets []models.EmployeeAsset
	err := r.db.Where("draft_id = ?", draftID).Find(&assets).Error
	return assets, err
}

// FindByEmployeeID finds all assets by employee ID
func (r *EmployeeAssetRepository) FindByEmployeeID(employeeID uint) ([]models.EmployeeAsset, error) {
	var assets []models.EmployeeAsset
	err := r.db.Where("employee_id = ?", employeeID).Find(&assets).Error
	return assets, err
}

// Update updates an asset
func (r *EmployeeAssetRepository) Update(asset *models.EmployeeAsset) error {
	return r.db.Save(asset).Error
}

// DeleteByDraftID deletes all assets by draft ID
func (r *EmployeeAssetRepository) DeleteByDraftID(draftID uint) error {
	return r.db.Where("draft_id = ?", draftID).Delete(&models.EmployeeAsset{}).Error
}

func (r *EmployeeAssetRepository) DeleteByEmployeeID(employeeID uint) error {
	return r.db.Where("employee_id = ?", employeeID).Delete(&models.EmployeeAsset{}).Error
}

// MigrateToEmployee migrates assets from draft to employee
func (r *EmployeeAssetRepository) MigrateToEmployee(draftID uint, employeeID uint) error {
	return r.db.Model(&models.EmployeeAsset{}).
		Where("draft_id = ?", draftID).
		Updates(map[string]interface{}{
			"employee_id": employeeID,
			"draft_id":    nil,
		}).Error
}

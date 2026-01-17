package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type EmployeeAddressRepository struct {
	db *gorm.DB
}

func NewEmployeeAddressRepository() *EmployeeAddressRepository {
	return &EmployeeAddressRepository{
		db: database.GetDB(),
	}
}

// Create creates a new employee address
func (r *EmployeeAddressRepository) Create(address *models.EmployeeAddress) error {
	return r.db.Create(address).Error
}

// FindByDraftID finds addresses by draft ID
func (r *EmployeeAddressRepository) FindByDraftID(draftID uint) ([]models.EmployeeAddress, error) {
	var addresses []models.EmployeeAddress
	err := r.db.Where("draft_id = ?", draftID).Find(&addresses).Error
	return addresses, err
}

// FindByEmployeeID finds addresses by employee ID
func (r *EmployeeAddressRepository) FindByEmployeeID(employeeID uint) ([]models.EmployeeAddress, error) {
	var addresses []models.EmployeeAddress
	err := r.db.Where("employee_id = ?", employeeID).Find(&addresses).Error
	return addresses, err
}

// Update updates an address
func (r *EmployeeAddressRepository) Update(address *models.EmployeeAddress) error {
	return r.db.Save(address).Error
}

// DeleteByDraftID deletes all addresses for a draft
func (r *EmployeeAddressRepository) DeleteByDraftID(draftID uint) error {
	return r.db.Where("draft_id = ?", draftID).Delete(&models.EmployeeAddress{}).Error
}

// MigrateToEmployee migrates addresses from draft to employee
func (r *EmployeeAddressRepository) MigrateToEmployee(draftID uint, employeeID uint) error {
	return r.db.Model(&models.EmployeeAddress{}).
		Where("draft_id = ? AND employee_id IS NULL", draftID).
		Updates(map[string]interface{}{
			"employee_id": employeeID,
			"draft_id":    nil,
		}).Error
}

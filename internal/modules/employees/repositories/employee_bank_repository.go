package repositories

import (
	"errors"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type EmployeeBankRepository struct {
	db *gorm.DB
}

func NewEmployeeBankRepository() *EmployeeBankRepository {
	return &EmployeeBankRepository{
		db: database.GetDB(),
	}
}

// Create creates a new bank account record
func (r *EmployeeBankRepository) Create(bank *models.EmployeeBankAccount) error {
	return r.db.Create(bank).Error
}

// FindByDraftID finds bank account by draft ID
func (r *EmployeeBankRepository) FindByDraftID(draftID uint) (*models.EmployeeBankAccount, error) {
	var bank models.EmployeeBankAccount
	err := r.db.Where("draft_id = ?", draftID).First(&bank).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found (expected)
		}
		return nil, err
	}
	return &bank, nil
}

// FindByEmployeeID finds bank accounts by employee ID
func (r *EmployeeBankRepository) FindByEmployeeID(employeeID uint) ([]models.EmployeeBankAccount, error) {
	var banks []models.EmployeeBankAccount
	err := r.db.Where("employee_id = ?", employeeID).Find(&banks).Error
	return banks, err
}

// Update updates bank account
func (r *EmployeeBankRepository) Update(bank *models.EmployeeBankAccount) error {
	return r.db.Save(bank).Error
}

// DeleteByDraftID deletes bank account for a draft
func (r *EmployeeBankRepository) DeleteByDraftID(draftID uint) error {
	return r.db.Where("draft_id = ?", draftID).Delete(&models.EmployeeBankAccount{}).Error
}

// MigrateToEmployee migrates bank account from draft to employee
func (r *EmployeeBankRepository) MigrateToEmployee(draftID uint, employeeID uint) error {
	return r.db.Model(&models.EmployeeBankAccount{}).
		Where("draft_id = ? AND employee_id IS NULL", draftID).
		Updates(map[string]interface{}{
			"employee_id": employeeID,
			"draft_id":    nil,
		}).Error
}

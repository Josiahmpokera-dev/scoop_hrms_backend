package repositories

import (
	"errors"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type EmployeePolicyRepository struct {
	db *gorm.DB
}

func NewEmployeePolicyRepository() *EmployeePolicyRepository {
	return &EmployeePolicyRepository{
		db: database.GetDB(),
	}
}

// Create creates a new policy assignment
func (r *EmployeePolicyRepository) Create(policy *models.EmployeePolicy) error {
	return r.db.Create(policy).Error
}

// FindByDraftID finds policy by draft ID
func (r *EmployeePolicyRepository) FindByDraftID(draftID uint) (*models.EmployeePolicy, error) {
	var policy models.EmployeePolicy
	err := r.db.Where("draft_id = ?", draftID).First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found (expected)
		}
		return nil, err
	}
	return &policy, nil
}

// FindByEmployeeID finds policy by employee ID
func (r *EmployeePolicyRepository) FindByEmployeeID(employeeID uint) (*models.EmployeePolicy, error) {
	var policy models.EmployeePolicy
	err := r.db.Where("employee_id = ?", employeeID).First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// Update updates policy assignment
func (r *EmployeePolicyRepository) Update(policy *models.EmployeePolicy) error {
	return r.db.Save(policy).Error
}

// DeleteByDraftID deletes policy for a draft
func (r *EmployeePolicyRepository) DeleteByDraftID(draftID uint) error {
	return r.db.Where("draft_id = ?", draftID).Delete(&models.EmployeePolicy{}).Error
}

// MigrateToEmployee migrates policy from draft to employee
func (r *EmployeePolicyRepository) MigrateToEmployee(draftID uint, employeeID uint) error {
	return r.db.Model(&models.EmployeePolicy{}).
		Where("draft_id = ? AND employee_id IS NULL", draftID).
		Updates(map[string]interface{}{
			"employee_id": employeeID,
			"draft_id":    nil,
		}).Error
}

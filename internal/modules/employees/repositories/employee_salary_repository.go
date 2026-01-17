package repositories

import (
	"errors"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type EmployeeSalaryRepository struct {
	db *gorm.DB
}

func NewEmployeeSalaryRepository() *EmployeeSalaryRepository {
	return &EmployeeSalaryRepository{
		db: database.GetDB(),
	}
}

// Create creates a new salary component record
func (r *EmployeeSalaryRepository) Create(salary *models.EmployeeSalaryComponent) error {
	return r.db.Create(salary).Error
}

// FindByDraftID finds salary by draft ID
func (r *EmployeeSalaryRepository) FindByDraftID(draftID uint) (*models.EmployeeSalaryComponent, error) {
	var salary models.EmployeeSalaryComponent
	err := r.db.Where("draft_id = ?", draftID).First(&salary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found (expected)
		}
		return nil, err
	}
	return &salary, nil
}

// FindByEmployeeID finds salary by employee ID
func (r *EmployeeSalaryRepository) FindByEmployeeID(employeeID uint) (*models.EmployeeSalaryComponent, error) {
	var salary models.EmployeeSalaryComponent
	err := r.db.Where("employee_id = ?", employeeID).First(&salary).Error
	if err != nil {
		return nil, err
	}
	return &salary, nil
}

// Update updates salary component
func (r *EmployeeSalaryRepository) Update(salary *models.EmployeeSalaryComponent) error {
	return r.db.Save(salary).Error
}

// DeleteByDraftID deletes salary for a draft
func (r *EmployeeSalaryRepository) DeleteByDraftID(draftID uint) error {
	return r.db.Where("draft_id = ?", draftID).Delete(&models.EmployeeSalaryComponent{}).Error
}

// MigrateToEmployee migrates salary from draft to employee
func (r *EmployeeSalaryRepository) MigrateToEmployee(draftID uint, employeeID uint) error {
	return r.db.Model(&models.EmployeeSalaryComponent{}).
		Where("draft_id = ? AND employee_id IS NULL", draftID).
		Updates(map[string]interface{}{
			"employee_id": employeeID,
			"draft_id":    nil,
		}).Error
}

package repositories

import (
	"errors"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type EmployeeEmploymentDetailsRepository struct {
	db *gorm.DB
}

func NewEmployeeEmploymentDetailsRepository() *EmployeeEmploymentDetailsRepository {
	return &EmployeeEmploymentDetailsRepository{
		db: database.GetDB(),
	}
}

// Create creates a new employment details record
func (r *EmployeeEmploymentDetailsRepository) Create(details *models.EmployeeEmploymentDetails) error {
	return r.db.Create(details).Error
}

// FindByDraftID finds employment details by draft ID
func (r *EmployeeEmploymentDetailsRepository) FindByDraftID(draftID uint) (*models.EmployeeEmploymentDetails, error) {
	var details models.EmployeeEmploymentDetails
	err := r.db.Where("draft_id = ?", draftID).First(&details).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found (expected)
		}
		return nil, err
	}
	return &details, nil
}

// FindByEmployeeID finds employment details by employee ID
func (r *EmployeeEmploymentDetailsRepository) FindByEmployeeID(employeeID uint) (*models.EmployeeEmploymentDetails, error) {
	var details models.EmployeeEmploymentDetails
	err := r.db.Where("employee_id = ?", employeeID).First(&details).Error
	if err != nil {
		return nil, err
	}
	return &details, nil
}

// Update updates employment details
func (r *EmployeeEmploymentDetailsRepository) Update(details *models.EmployeeEmploymentDetails) error {
	return r.db.Save(details).Error
}

// DeleteByDraftID deletes employment details by draft ID
func (r *EmployeeEmploymentDetailsRepository) DeleteByDraftID(draftID uint) error {
	return r.db.Where("draft_id = ?", draftID).Delete(&models.EmployeeEmploymentDetails{}).Error
}

// MigrateToEmployee migrates employment details from draft to employee
func (r *EmployeeEmploymentDetailsRepository) MigrateToEmployee(draftID uint, employeeID uint) error {
	return r.db.Model(&models.EmployeeEmploymentDetails{}).
		Where("draft_id = ?", draftID).
		Updates(map[string]interface{}{
			"employee_id": employeeID,
			"draft_id":    nil,
		}).Error
}

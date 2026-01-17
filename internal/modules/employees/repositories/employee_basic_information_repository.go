package repositories

import (
	"errors"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type EmployeeBasicInformationRepository struct {
	db *gorm.DB
}

func NewEmployeeBasicInformationRepository() *EmployeeBasicInformationRepository {
	return &EmployeeBasicInformationRepository{
		db: database.GetDB(),
	}
}

// Create creates a new basic information record
func (r *EmployeeBasicInformationRepository) Create(info *models.EmployeeBasicInformation) error {
	return r.db.Create(info).Error
}

// FindByDraftID finds basic information by draft ID
func (r *EmployeeBasicInformationRepository) FindByDraftID(draftID uint) (*models.EmployeeBasicInformation, error) {
	var info models.EmployeeBasicInformation
	err := r.db.Where("draft_id = ?", draftID).First(&info).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found (expected)
		}
		return nil, err
	}
	return &info, nil
}

// FindByEmployeeID finds basic information by employee ID
func (r *EmployeeBasicInformationRepository) FindByEmployeeID(employeeID uint) (*models.EmployeeBasicInformation, error) {
	var info models.EmployeeBasicInformation
	err := r.db.Where("employee_id = ?", employeeID).First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// Update updates basic information
func (r *EmployeeBasicInformationRepository) Update(info *models.EmployeeBasicInformation) error {
	return r.db.Save(info).Error
}

// DeleteByDraftID deletes basic information by draft ID
func (r *EmployeeBasicInformationRepository) DeleteByDraftID(draftID uint) error {
	return r.db.Where("draft_id = ?", draftID).Delete(&models.EmployeeBasicInformation{}).Error
}

// MigrateToEmployee migrates basic information from draft to employee
func (r *EmployeeBasicInformationRepository) MigrateToEmployee(draftID uint, employeeID uint) error {
	return r.db.Model(&models.EmployeeBasicInformation{}).
		Where("draft_id = ?", draftID).
		Updates(map[string]interface{}{
			"employee_id": employeeID,
			"draft_id":    nil,
		}).Error
}

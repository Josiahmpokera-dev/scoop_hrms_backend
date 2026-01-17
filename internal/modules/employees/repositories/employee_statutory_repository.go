package repositories

import (
	"errors"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type EmployeeStatutoryRepository struct {
	db *gorm.DB
}

func NewEmployeeStatutoryRepository() *EmployeeStatutoryRepository {
	return &EmployeeStatutoryRepository{
		db: database.GetDB(),
	}
}

// Create creates a new statutory info record
func (r *EmployeeStatutoryRepository) Create(statutory *models.EmployeeStatutoryInfo) error {
	return r.db.Create(statutory).Error
}

// FindByDraftID finds statutory info by draft ID
func (r *EmployeeStatutoryRepository) FindByDraftID(draftID uint) (*models.EmployeeStatutoryInfo, error) {
	var statutory models.EmployeeStatutoryInfo
	err := r.db.Where("draft_id = ?", draftID).First(&statutory).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found (expected)
		}
		return nil, err
	}
	return &statutory, nil
}

// FindByEmployeeID finds statutory info by employee ID
func (r *EmployeeStatutoryRepository) FindByEmployeeID(employeeID uint) (*models.EmployeeStatutoryInfo, error) {
	var statutory models.EmployeeStatutoryInfo
	err := r.db.Where("employee_id = ?", employeeID).First(&statutory).Error
	if err != nil {
		return nil, err
	}
	return &statutory, nil
}

// Update updates statutory info
func (r *EmployeeStatutoryRepository) Update(statutory *models.EmployeeStatutoryInfo) error {
	return r.db.Save(statutory).Error
}

// DeleteByDraftID deletes statutory info for a draft
func (r *EmployeeStatutoryRepository) DeleteByDraftID(draftID uint) error {
	return r.db.Where("draft_id = ?", draftID).Delete(&models.EmployeeStatutoryInfo{}).Error
}

// MigrateToEmployee migrates statutory info from draft to employee
func (r *EmployeeStatutoryRepository) MigrateToEmployee(draftID uint, employeeID uint) error {
	return r.db.Model(&models.EmployeeStatutoryInfo{}).
		Where("draft_id = ? AND employee_id IS NULL", draftID).
		Updates(map[string]interface{}{
			"employee_id": employeeID,
			"draft_id":    nil,
		}).Error
}

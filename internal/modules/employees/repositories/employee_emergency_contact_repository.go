package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type EmployeeEmergencyContactRepository struct {
	db *gorm.DB
}

func NewEmployeeEmergencyContactRepository() *EmployeeEmergencyContactRepository {
	return &EmployeeEmergencyContactRepository{
		db: database.GetDB(),
	}
}

// Create creates a new emergency contact
func (r *EmployeeEmergencyContactRepository) Create(contact *models.EmployeeEmergencyContact) error {
	return r.db.Create(contact).Error
}

// FindByDraftID finds emergency contacts by draft ID
func (r *EmployeeEmergencyContactRepository) FindByDraftID(draftID uint) ([]models.EmployeeEmergencyContact, error) {
	var contacts []models.EmployeeEmergencyContact
	err := r.db.Where("draft_id = ?", draftID).Find(&contacts).Error
	return contacts, err
}

// FindByEmployeeID finds emergency contacts by employee ID
func (r *EmployeeEmergencyContactRepository) FindByEmployeeID(employeeID uint) ([]models.EmployeeEmergencyContact, error) {
	var contacts []models.EmployeeEmergencyContact
	err := r.db.Where("employee_id = ?", employeeID).Find(&contacts).Error
	return contacts, err
}

// Update updates an emergency contact
func (r *EmployeeEmergencyContactRepository) Update(contact *models.EmployeeEmergencyContact) error {
	return r.db.Save(contact).Error
}

// Delete deletes an emergency contact
func (r *EmployeeEmergencyContactRepository) Delete(id uint) error {
	return r.db.Delete(&models.EmployeeEmergencyContact{}, id).Error
}

// DeleteByDraftID deletes all emergency contacts for a draft
func (r *EmployeeEmergencyContactRepository) DeleteByDraftID(draftID uint) error {
	return r.db.Where("draft_id = ?", draftID).Delete(&models.EmployeeEmergencyContact{}).Error
}

// MigrateToEmployee migrates emergency contacts from draft to employee
func (r *EmployeeEmergencyContactRepository) MigrateToEmployee(draftID uint, employeeID uint) error {
	return r.db.Model(&models.EmployeeEmergencyContact{}).
		Where("draft_id = ? AND employee_id IS NULL", draftID).
		Updates(map[string]interface{}{
			"employee_id": employeeID,
			"draft_id":    nil,
		}).Error
}

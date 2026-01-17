package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type EmployeeDocumentRepository struct {
	db *gorm.DB
}

func NewEmployeeDocumentRepository() *EmployeeDocumentRepository {
	return &EmployeeDocumentRepository{
		db: database.GetDB(),
	}
}

// Create creates a new document record
func (r *EmployeeDocumentRepository) Create(doc *models.EmployeeDocument) error {
	return r.db.Create(doc).Error
}

// FindByDraftID finds documents by draft ID
func (r *EmployeeDocumentRepository) FindByDraftID(draftID uint) ([]models.EmployeeDocument, error) {
	var docs []models.EmployeeDocument
	err := r.db.Where("draft_id = ?", draftID).Find(&docs).Error
	return docs, err
}

// FindByDraftIDAndType finds a document by draft ID and document type
func (r *EmployeeDocumentRepository) FindByDraftIDAndType(draftID uint, documentType models.DocumentType) (*models.EmployeeDocument, error) {
	var doc models.EmployeeDocument
	err := r.db.Where("draft_id = ? AND document_type = ?", draftID, documentType).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// FindByEmployeeID finds documents by employee ID
func (r *EmployeeDocumentRepository) FindByEmployeeID(employeeID uint) ([]models.EmployeeDocument, error) {
	var docs []models.EmployeeDocument
	err := r.db.Where("employee_id = ?", employeeID).Find(&docs).Error
	return docs, err
}

// Update updates a document
func (r *EmployeeDocumentRepository) Update(doc *models.EmployeeDocument) error {
	return r.db.Save(doc).Error
}

// Delete deletes a document
func (r *EmployeeDocumentRepository) Delete(id uint) error {
	return r.db.Delete(&models.EmployeeDocument{}, id).Error
}

// DeleteByDraftID deletes all documents for a draft
func (r *EmployeeDocumentRepository) DeleteByDraftID(draftID uint) error {
	return r.db.Where("draft_id = ?", draftID).Delete(&models.EmployeeDocument{}).Error
}

// MigrateToEmployee migrates documents from draft to employee
func (r *EmployeeDocumentRepository) MigrateToEmployee(draftID uint, employeeID uint) error {
	return r.db.Model(&models.EmployeeDocument{}).
		Where("draft_id = ? AND employee_id IS NULL", draftID).
		Updates(map[string]interface{}{
			"employee_id": employeeID,
			"draft_id":    nil,
		}).Error
}

// FindByID finds a document by ID
func (r *EmployeeDocumentRepository) FindByID(id uint) (*models.EmployeeDocument, error) {
	var doc models.EmployeeDocument
	err := r.db.First(&doc, id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

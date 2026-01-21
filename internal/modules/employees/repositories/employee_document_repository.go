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

// FindByEmployeeIDWithPagination finds documents by employee ID with pagination
func (r *EmployeeDocumentRepository) FindByEmployeeIDWithPagination(employeeID uint, documentType *string, page, pageSize int) ([]models.EmployeeDocument, int64, error) {
	var docs []models.EmployeeDocument
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.EmployeeDocument{}).Where("employee_id = ?", employeeID)

	if documentType != nil && *documentType != "" {
		query = query.Where("document_type = ?", *documentType)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&docs).Error
	return docs, total, err
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

// FindByEmployeeIDString finds documents by employee ID string (e.g., "EMP001")
func (r *EmployeeDocumentRepository) FindByEmployeeIDString(employeeID string) ([]models.EmployeeDocument, error) {
	var docs []models.EmployeeDocument
	err := r.db.Where("employee_id_string = ? OR employee_id IN (SELECT id FROM employees WHERE employee_id = ?)", employeeID, employeeID).
		Order("created_at DESC").
		Find(&docs).Error
	return docs, err
}

// GetAllDocuments gets all documents (for statistics)
func (r *EmployeeDocumentRepository) GetAllDocuments(tenantID *uint) ([]models.EmployeeDocument, error) {
	var docs []models.EmployeeDocument
	query := r.db.Model(&models.EmployeeDocument{})
	
	// Filter by tenant if provided
	if tenantID != nil {
		// Filter documents where employee belongs to tenant
		query = query.Where("employee_id IN (SELECT id FROM employees WHERE tenant_id = ?) OR employee_id_string IN (SELECT employee_id FROM employees WHERE tenant_id = ?)", *tenantID, *tenantID)
	}
	
	err := query.Find(&docs).Error
	return docs, err
}

// CountByType counts documents grouped by type
func (r *EmployeeDocumentRepository) CountByType(tenantID *uint) (map[string]int64, error) {
	type Result struct {
		DocumentType string
		Count        int64
	}
	
	var results []Result
	query := r.db.Model(&models.EmployeeDocument{}).
		Select("document_type, COUNT(*) as count").
		Group("document_type")
	
	// Filter by tenant if provided
	if tenantID != nil {
		query = query.Where("employee_id IN (SELECT id FROM employees WHERE tenant_id = ?) OR employee_id_string IN (SELECT employee_id FROM employees WHERE tenant_id = ?)", *tenantID, *tenantID)
	}
	
	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}
	
	resultMap := make(map[string]int64)
	for _, r := range results {
		resultMap[r.DocumentType] = r.Count
	}
	
	return resultMap, nil
}

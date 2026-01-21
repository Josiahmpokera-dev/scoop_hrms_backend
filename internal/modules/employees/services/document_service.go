package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
)

// DocumentService handles document-related business logic
type DocumentService struct {
	documentRepo  *employeeRepos.EmployeeDocumentRepository
	employeeRepo  *employeeRepos.EmployeeRepository
	statutoryRepo *employeeRepos.EmployeeStatutoryRepository
}

// NewDocumentService creates a new document service
func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentRepo:  employeeRepos.NewEmployeeDocumentRepository(),
		employeeRepo:  employeeRepos.NewEmployeeRepository(),
		statutoryRepo: employeeRepos.NewEmployeeStatutoryRepository(),
	}
}

// ListEmployeesWithDocuments lists all employees with their uploaded documents
func (s *DocumentService) ListEmployeesWithDocuments(tenantID *uint, page, pageSize int) ([]map[string]interface{}, int64, error) {
	// Get employees with pagination
	offset := (page - 1) * pageSize
	employees, total, err := s.employeeRepo.List(pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get employees: %w", err)
	}

	// Filter by tenant if needed
	filteredEmployees := make([]models.Employee, 0)
	for _, emp := range employees {
		if tenantID == nil || (emp.TenantID != nil && *emp.TenantID == *tenantID) {
			filteredEmployees = append(filteredEmployees, emp)
		}
	}

	// Get documents for all employees
	result := make([]map[string]interface{}, len(filteredEmployees))
	for i, emp := range filteredEmployees {
		employeeData := map[string]interface{}{
			"employee_id":   emp.EmployeeID,
			"first_name":     emp.FirstName,
			"last_name":      emp.LastName,
			"email":          emp.WorkEmail,
			"department_id":  emp.DepartmentID,
			"position_id":    emp.PositionID,
			"status":         emp.Status,
			"hire_date":      emp.HireDate,
		}

		// Get documents for this employee (by employee_id string or employee DB ID)
		docs, err := s.documentRepo.FindByEmployeeIDString(emp.EmployeeID)
		if err == nil && len(docs) > 0 {
			// Create comma-separated list of document types
			documentTypes := make([]string, len(docs))
			for j, doc := range docs {
				documentTypes[j] = string(doc.DocumentType)
			}
			// Join document types with comma and space
			documentsList := ""
			if len(documentTypes) > 0 {
				documentsList = strings.Join(documentTypes, ", ")
			}
			employeeData["documents"] = documentsList
			employeeData["document_count"] = len(docs)
		} else {
			employeeData["documents"] = ""
			employeeData["document_count"] = 0
		}

		result[i] = employeeData
	}

	return result, total, nil
}

// GetEmployeeDocuments gets all documents for a specific employee
func (s *DocumentService) GetEmployeeDocuments(employeeID string, tenantID *uint) ([]models.EmployeeDocument, error) {
	// Verify employee exists and belongs to tenant
	employee, err := s.employeeRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, fmt.Errorf("employee not found")
	}

	// Verify tenant ownership
	if tenantID != nil && employee.TenantID != nil && *employee.TenantID != *tenantID {
		return nil, fmt.Errorf("employee does not belong to your tenant")
	}

	// Get documents
	docs, err := s.documentRepo.FindByEmployeeIDString(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	return docs, nil
}

// GetDocumentByID gets a specific document by ID
func (s *DocumentService) GetDocumentByID(documentID uint, tenantID *uint) (*models.EmployeeDocument, error) {
	doc, err := s.documentRepo.FindByID(documentID)
	if err != nil {
		return nil, fmt.Errorf("document not found")
	}

	// Verify tenant ownership through employee
	if doc.EmployeeID != nil {
		employee, err := s.employeeRepo.FindByID(*doc.EmployeeID)
		if err == nil && employee != nil {
			if tenantID != nil && employee.TenantID != nil && *employee.TenantID != *tenantID {
				return nil, fmt.Errorf("document does not belong to your tenant")
			}
		}
	} else if doc.EmployeeIDString != nil {
		// Check by employee ID string
		employee, err := s.employeeRepo.FindByEmployeeID(*doc.EmployeeIDString)
		if err == nil && employee != nil {
			if tenantID != nil && employee.TenantID != nil && *employee.TenantID != *tenantID {
				return nil, fmt.Errorf("document does not belong to your tenant")
			}
		}
	}

	return doc, nil
}

// GetDocumentStatistics gets statistics about employee documents
func (s *DocumentService) GetDocumentStatistics(tenantID *uint) (map[string]interface{}, error) {
	// Get all documents for the tenant
	allDocs, err := s.documentRepo.GetAllDocuments(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	// Total documents
	totalDocuments := int64(len(allDocs))

	// Count by document type
	documentTypes, err := s.documentRepo.CountByType(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to count by type: %w", err)
	}

	// Convert to array format for response
	documentTypesArray := make([]map[string]interface{}, 0, len(documentTypes))
	for docType, count := range documentTypes {
		documentTypesArray = append(documentTypesArray, map[string]interface{}{
			"type":  docType,
			"count": count,
		})
	}

	// Calculate expired and expiring soon
	now := time.Now()
	expiringSoonThreshold := now.AddDate(0, 0, 30) // 30 days from now
	defaultValidityPeriod := 2 * 365 * 24 * time.Hour // 2 years for documents without explicit expiry

	var expiredCount int64
	var expiringSoonCount int64

	// Track employees we've already checked for statutory info
	employeeStatutoryCache := make(map[uint]*models.EmployeeStatutoryInfo)

	for _, doc := range allDocs {
		var expiryDate *time.Time

		// Determine expiry date based on document type
		if doc.DocumentType == models.DocumentTypeWorkPermit {
			// For work permits, check statutory info
			if doc.EmployeeID != nil {
				statutory, exists := employeeStatutoryCache[*doc.EmployeeID]
				if !exists {
					statutory, _ = s.statutoryRepo.FindByEmployeeID(*doc.EmployeeID)
					if statutory != nil {
						employeeStatutoryCache[*doc.EmployeeID] = statutory
					}
				}
				if statutory != nil && statutory.WorkPermitExpiryDate != nil {
					expiryDate = statutory.WorkPermitExpiryDate
				}
			}
		} else if doc.DocumentType == models.DocumentTypeIdentity {
			// For identity documents (passports), check statutory info
			if doc.EmployeeID != nil {
				statutory, exists := employeeStatutoryCache[*doc.EmployeeID]
				if !exists {
					statutory, _ = s.statutoryRepo.FindByEmployeeID(*doc.EmployeeID)
					if statutory != nil {
						employeeStatutoryCache[*doc.EmployeeID] = statutory
					}
				}
				if statutory != nil && statutory.PassportExpiryDate != nil {
					expiryDate = statutory.PassportExpiryDate
				}
			}
		}

		// If no explicit expiry date, use default validity period from upload date
		if expiryDate == nil {
			defaultExpiry := doc.CreatedAt.Add(defaultValidityPeriod)
			expiryDate = &defaultExpiry
		}

		// Check if expired or expiring soon
		if expiryDate.Before(now) {
			expiredCount++
		} else if expiryDate.Before(expiringSoonThreshold) && expiryDate.After(now) {
			expiringSoonCount++
		}
	}

	return map[string]interface{}{
		"total_documents": totalDocuments,
		"expired":         expiredCount,
		"expiring_soon":   expiringSoonCount,
		"document_types":  documentTypesArray,
	}, nil
}

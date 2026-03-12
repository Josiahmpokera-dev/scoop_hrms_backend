package repositories

import (
	"fmt"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

// ServiceRequestRepository handles service request database operations
type ServiceRequestRepository struct {
	db *gorm.DB
}

// NewServiceRequestRepository creates a new service request repository
func NewServiceRequestRepository() *ServiceRequestRepository {
	return &ServiceRequestRepository{
		db: database.GetDB(),
	}
}

// Create creates a new service request
func (r *ServiceRequestRepository) Create(request *models.ServiceRequest) error {
	return r.db.Create(request).Error
}

// FindByID finds a service request by ID
func (r *ServiceRequestRepository) FindByID(id uint) (*models.ServiceRequest, error) {
	var request models.ServiceRequest
	err := r.db.Where("id = ?", id).First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// FindByEmployeeID finds service requests by employee ID
func (r *ServiceRequestRepository) FindByEmployeeID(employeeID uint, tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.ServiceRequest, int64, error) {
	var requests []models.ServiceRequest
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.ServiceRequest{}).Where("employee_id = ?", employeeID)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if requestType, ok := filters["type"].(string); ok && requestType != "" {
		query = query.Where("type = ?", requestType)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if priority, ok := filters["priority"].(string); ok && priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error
	return requests, total, err
}

// FindAll finds all service requests with pagination and filters (for HR/admin use)
func (r *ServiceRequestRepository) FindAll(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.ServiceRequest, int64, error) {
	var requests []models.ServiceRequest
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.ServiceRequest{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if requestType, ok := filters["type"].(string); ok && requestType != "" {
		query = query.Where("type = ?", requestType)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if priority, ok := filters["priority"].(string); ok && priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if letterType, ok := filters["letter_type"].(string); ok && letterType != "" {
		query = query.Where("letter_type = ?", letterType)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error
	return requests, total, err
}

// Update updates a service request
func (r *ServiceRequestRepository) Update(request *models.ServiceRequest) error {
	return r.db.Save(request).Error
}

// GenerateRequestNumber generates a unique request number
func (r *ServiceRequestRepository) GenerateRequestNumber(tenantID *uint) (string, error) {
	var count int64
	query := r.db.Model(&models.ServiceRequest{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Count(&count).Error; err != nil {
		return "", err
	}
	
	// Format: SR-YYYY-XXX (e.g., SR-2026-001)
	year := "2026" // TODO: Get current year dynamically
	sequence := int(count) + 1
	return "SR-" + year + "-" + formatServiceRequestSequence(sequence), nil
}

// formatServiceRequestSequence formats sequence number with leading zeros
func formatServiceRequestSequence(seq int) string {
	if seq < 10 {
		return "00" + fmt.Sprintf("%d", seq)
	} else if seq < 100 {
		return "0" + fmt.Sprintf("%d", seq)
	}
	return fmt.Sprintf("%d", seq)
}

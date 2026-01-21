package repositories

import (
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

// ProfileUpdateRequestRepository handles profile update request database operations
type ProfileUpdateRequestRepository struct {
	db *gorm.DB
}

// NewProfileUpdateRequestRepository creates a new profile update request repository
func NewProfileUpdateRequestRepository() *ProfileUpdateRequestRepository {
	return &ProfileUpdateRequestRepository{
		db: database.GetDB(),
	}
}

// Create creates a new profile update request
func (r *ProfileUpdateRequestRepository) Create(request *models.ProfileUpdateRequest) error {
	return r.db.Create(request).Error
}

// FindByEmployeeID finds profile update requests by employee ID
func (r *ProfileUpdateRequestRepository) FindByEmployeeID(employeeID uint, tenantID *uint) ([]models.ProfileUpdateRequest, error) {
	var requests []models.ProfileUpdateRequest
	query := r.db.Where("employee_id = ?", employeeID)
	
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	err := query.Order("submitted_at DESC").Find(&requests).Error
	return requests, err
}

// FindPendingByEmployeeID finds pending profile update requests for an employee
func (r *ProfileUpdateRequestRepository) FindPendingByEmployeeID(employeeID uint, tenantID *uint) ([]models.ProfileUpdateRequest, error) {
	var requests []models.ProfileUpdateRequest
	query := r.db.Where("employee_id = ? AND status = ?", employeeID, models.ProfileUpdateStatusPending)
	
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	err := query.Order("submitted_at DESC").Find(&requests).Error
	return requests, err
}

// GenerateUpdateRequestID generates a unique update request ID
func (r *ProfileUpdateRequestRepository) GenerateUpdateRequestID(tenantID *uint) (string, error) {
	var count int64
	query := r.db.Model(&models.ProfileUpdateRequest{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Count(&count).Error; err != nil {
		return "", err
	}
	
	// Format: PRU-YYYY-XXX (e.g., PRU-2026-001)
	year := time.Now().Format("2006")
	sequence := int(count) + 1
	return "PRU-" + year + "-" + formatProfileUpdateSequence(sequence), nil
}

// formatProfileUpdateSequence formats sequence number with leading zeros
func formatProfileUpdateSequence(seq int) string {
	if seq < 10 {
		return "00" + fmt.Sprintf("%d", seq)
	} else if seq < 100 {
		return "0" + fmt.Sprintf("%d", seq)
	}
	return fmt.Sprintf("%d", seq)
}

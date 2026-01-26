package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	"gorm.io/gorm"
)

type RosterChangeRequestRepository struct {
	db *gorm.DB
}

func NewRosterChangeRequestRepository() *RosterChangeRequestRepository {
	return &RosterChangeRequestRepository{
		db: database.GetDB(),
	}
}

// Create creates a new roster change request
func (r *RosterChangeRequestRepository) Create(request *models.RosterChangeRequest) error {
	return r.db.Create(request).Error
}

// FindByID finds a roster change request by ID
func (r *RosterChangeRequestRepository) FindByID(id uint) (*models.RosterChangeRequest, error) {
	var request models.RosterChangeRequest
	err := r.db.Preload("Assignment").Preload("Assignment.Shift").
		Preload("Assignment.Location").
		Preload("RequestedShift").
		First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// Update updates a roster change request
func (r *RosterChangeRequestRepository) Update(request *models.RosterChangeRequest) error {
	return r.db.Save(request).Error
}

// Delete soft deletes a roster change request
func (r *RosterChangeRequestRepository) Delete(id uint) error {
	return r.db.Delete(&models.RosterChangeRequest{}, id).Error
}

// List returns roster change requests with pagination and filters
func (r *RosterChangeRequestRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.RosterChangeRequest, int64, error) {
	var requests []models.RosterChangeRequest
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.RosterChangeRequest{}).
		Preload("Assignment").Preload("Assignment.Shift").
		Preload("RequestedShift")

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if requestedBy, ok := filters["requested_by"].(string); ok && requestedBy != "" {
		query = query.Where("requested_by = ?", requestedBy)
	}

	if assignmentID, ok := filters["assignment_id"].(uint); ok && assignmentID > 0 {
		query = query.Where("assignment_id = ?", assignmentID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("requested_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error
	return requests, total, err
}

// HasPendingRequest checks if there's a pending change request for an assignment
func (r *RosterChangeRequestRepository) HasPendingRequest(assignmentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.RosterChangeRequest{}).
		Where("assignment_id = ? AND status = ?", assignmentID, "pending").
		Count(&count).Error
	return count > 0, err
}

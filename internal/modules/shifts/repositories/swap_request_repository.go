package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	"gorm.io/gorm"
)

type SwapRequestRepository struct {
	db *gorm.DB
}

func NewSwapRequestRepository() *SwapRequestRepository {
	return &SwapRequestRepository{
		db: database.GetDB(),
	}
}

// Create creates a new swap request
func (r *SwapRequestRepository) Create(request *models.SwapRequest) error {
	return r.db.Create(request).Error
}

// FindByID finds a swap request by ID
func (r *SwapRequestRepository) FindByID(id uint) (*models.SwapRequest, error) {
	var request models.SwapRequest
	err := r.db.Preload("Assignment").Preload("Assignment.Shift").
		Preload("Assignment.Location").
		Preload("SwapAssignment").Preload("SwapAssignment.Shift").
		Preload("SwapAssignment.Location").
		First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// Update updates a swap request
func (r *SwapRequestRepository) Update(request *models.SwapRequest) error {
	return r.db.Save(request).Error
}

// Delete soft deletes a swap request
func (r *SwapRequestRepository) Delete(id uint) error {
	return r.db.Delete(&models.SwapRequest{}, id).Error
}

// List returns swap requests with pagination and filters
func (r *SwapRequestRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.SwapRequest, int64, error) {
	var requests []models.SwapRequest
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.SwapRequest{}).
		Preload("Assignment").Preload("Assignment.Shift").
		Preload("SwapAssignment").Preload("SwapAssignment.Shift")

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

	if requestedWith, ok := filters["requested_with"].(string); ok && requestedWith != "" {
		query = query.Where("requested_with = ?", requestedWith)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("requested_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error
	return requests, total, err
}

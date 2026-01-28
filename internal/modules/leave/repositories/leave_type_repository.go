package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"gorm.io/gorm"
)

type LeaveTypeRepository struct {
	db *gorm.DB
}

func NewLeaveTypeRepository() *LeaveTypeRepository {
	return &LeaveTypeRepository{
		db: database.GetDB(),
	}
}

// Create creates a new leave type
func (r *LeaveTypeRepository) Create(leaveType *models.LeaveType) error {
	return r.db.Create(leaveType).Error
}

// FindByID finds a leave type by ID
func (r *LeaveTypeRepository) FindByID(id uint) (*models.LeaveType, error) {
	var leaveType models.LeaveType
	err := r.db.First(&leaveType, id).Error
	if err != nil {
		return nil, err
	}
	return &leaveType, nil
}

// FindByCode finds a leave type by code
func (r *LeaveTypeRepository) FindByCode(code string) (*models.LeaveType, error) {
	var leaveType models.LeaveType
	err := r.db.Where("code = ?", code).First(&leaveType).Error
	if err != nil {
		return nil, err
	}
	return &leaveType, nil
}

// Update updates a leave type
func (r *LeaveTypeRepository) Update(leaveType *models.LeaveType) error {
	return r.db.Save(leaveType).Error
}

// Delete soft deletes a leave type
func (r *LeaveTypeRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeaveType{}, id).Error
}

// List returns leave types with pagination and filters
func (r *LeaveTypeRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.LeaveType, int64, error) {
	var leaveTypes []models.LeaveType
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.LeaveType{})

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("category ASC, name ASC").Offset(offset).Limit(pageSize).Find(&leaveTypes).Error
	return leaveTypes, total, err
}

// ExistsByCode checks if a leave type with the given code exists
func (r *LeaveTypeRepository) ExistsByCode(code string) bool {
	var count int64
	r.db.Model(&models.LeaveType{}).Where("code = ?", code).Count(&count)
	return count > 0
}

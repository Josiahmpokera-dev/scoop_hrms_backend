package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"gorm.io/gorm"
)

type LeavePolicyRepository struct {
	db *gorm.DB
}

func NewLeavePolicyRepository() *LeavePolicyRepository {
	return &LeavePolicyRepository{
		db: database.GetDB(),
	}
}

// Create creates a new leave policy
func (r *LeavePolicyRepository) Create(policy *models.LeavePolicy) error {
	return r.db.Create(policy).Error
}

// FindByID finds a leave policy by ID
func (r *LeavePolicyRepository) FindByID(id uint) (*models.LeavePolicy, error) {
	var policy models.LeavePolicy
	err := r.db.Preload("LeaveType").First(&policy, id).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// FindByLeaveTypeCode finds policies by leave type code
func (r *LeavePolicyRepository) FindByLeaveTypeCode(leaveTypeCode string, tenantID *uint) ([]models.LeavePolicy, error) {
	var policies []models.LeavePolicy
	query := r.db.Preload("LeaveType").Where("leave_type_code = ?", leaveTypeCode)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.Find(&policies).Error
	return policies, err
}

// FindByCountryAndLeaveType finds policy by country and leave type
func (r *LeavePolicyRepository) FindByCountryAndLeaveType(country, leaveTypeCode string, tenantID *uint) (*models.LeavePolicy, error) {
	var policy models.LeavePolicy
	query := r.db.Preload("LeaveType").Where("country = ? AND leave_type_code = ?", country, leaveTypeCode)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// Update updates a leave policy
func (r *LeavePolicyRepository) Update(policy *models.LeavePolicy) error {
	return r.db.Save(policy).Error
}

// Delete soft deletes a leave policy
func (r *LeavePolicyRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeavePolicy{}, id).Error
}

// List returns leave policies with pagination and filters
func (r *LeavePolicyRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.LeavePolicy, int64, error) {
	var policies []models.LeavePolicy
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.LeavePolicy{}).Preload("LeaveType")

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if country, ok := filters["country"].(string); ok && country != "" {
		query = query.Where("country = ?", country)
	}

	if leaveTypeCode, ok := filters["leave_type_code"].(string); ok && leaveTypeCode != "" {
		query = query.Where("leave_type_code = ?", leaveTypeCode)
	}

	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("country ASC, leave_type_code ASC").Offset(offset).Limit(pageSize).Find(&policies).Error
	return policies, total, err
}

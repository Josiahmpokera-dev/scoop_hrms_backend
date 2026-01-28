package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"gorm.io/gorm"
)

type LeaveApprovalRepository struct {
	db *gorm.DB
}

func NewLeaveApprovalRepository() *LeaveApprovalRepository {
	return &LeaveApprovalRepository{
		db: database.GetDB(),
	}
}

// Create creates a new leave approval
func (r *LeaveApprovalRepository) Create(approval *models.LeaveApproval) error {
	return r.db.Create(approval).Error
}

// FindByID finds a leave approval by ID
func (r *LeaveApprovalRepository) FindByID(id uint) (*models.LeaveApproval, error) {
	var approval models.LeaveApproval
	err := r.db.First(&approval, id).Error
	if err != nil {
		return nil, err
	}
	return &approval, nil
}

// FindByLeaveRequestID finds all approvals for a leave request
func (r *LeaveApprovalRepository) FindByLeaveRequestID(requestID uint) ([]models.LeaveApproval, error) {
	var approvals []models.LeaveApproval
	err := r.db.Where("leave_request_id = ?", requestID).
		Order("level ASC").
		Find(&approvals).Error
	return approvals, err
}

// FindPendingByLeaveRequestID finds pending approvals for a leave request
func (r *LeaveApprovalRepository) FindPendingByLeaveRequestID(requestID uint) ([]models.LeaveApproval, error) {
	var approvals []models.LeaveApproval
	err := r.db.Where("leave_request_id = ? AND status = ? AND is_applicable = ?", requestID, "pending", true).
		Order("level ASC").
		Find(&approvals).Error
	return approvals, err
}

// Update updates a leave approval
func (r *LeaveApprovalRepository) Update(approval *models.LeaveApproval) error {
	return r.db.Save(approval).Error
}

// CreateApprovalWorkflow creates the approval workflow for a leave request
func (r *LeaveApprovalRepository) CreateApprovalWorkflow(requestID uint, tenantID *uint) error {
	// Create default approval levels
	approvals := []models.LeaveApproval{
		{
			LeaveRequestID: requestID,
			TenantID:       tenantID,
			Level:          1,
			ApproverType:   "head_of_department",
			Status:         "pending",
			IsApplicable:   true,
		},
		{
			LeaveRequestID: requestID,
			TenantID:       tenantID,
			Level:          2,
			ApproverType:   "hr_department",
			Status:         "pending",
			IsApplicable:   true,
		},
		{
			LeaveRequestID: requestID,
			TenantID:       tenantID,
			Level:          3,
			ApproverType:   "director_ceo",
			Status:         "pending",
			IsApplicable:   false, // May not be applicable for all requests
		},
	}

	for _, approval := range approvals {
		if err := r.db.Create(&approval).Error; err != nil {
			return err
		}
	}

	return nil
}

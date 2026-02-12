package services

import (
	"errors"
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/repositories"
)

type LeaveTypeService struct {
	repo *repositories.LeaveTypeRepository
}

func NewLeaveTypeService() *LeaveTypeService {
	return &LeaveTypeService{
		repo: repositories.NewLeaveTypeRepository(),
	}
}

// CreateLeaveType creates a new leave type
func (s *LeaveTypeService) CreateLeaveType(req *models.CreateLeaveTypeRequest, tenantID *uint, createdBy *uint) (*models.LeaveType, error) {
	_ = tenantID
	// Check if code already exists
	if s.repo.ExistsByCode(req.Code) {
		return nil, errors.New("leave type with this code already exists")
	}

	leaveType := &models.LeaveType{
		Code:                 req.Code,
		Name:                 req.Name,
		Icon:                 req.Icon,
		Category:             req.Category,
		PaidLeave:            true,
		RequiresDocumentation: false,
		IsActive:             true,
		Description:          req.Description,
		CreatedBy:            createdBy,
		UpdatedBy:            createdBy,
	}

	if req.PaidLeave != nil {
		leaveType.PaidLeave = *req.PaidLeave
	}
	if req.RequiresDocumentation != nil {
		leaveType.RequiresDocumentation = *req.RequiresDocumentation
	}
	if req.IsActive != nil {
		leaveType.IsActive = *req.IsActive
	}

	if err := s.repo.Create(leaveType); err != nil {
		return nil, fmt.Errorf("failed to create leave type: %w", err)
	}

	return s.repo.FindByID(leaveType.ID)
}

// GetLeaveTypeByID retrieves a leave type by ID
func (s *LeaveTypeService) GetLeaveTypeByID(id uint) (*models.LeaveType, error) {
	return s.repo.FindByID(id)
}

// GetLeaveTypeByCode retrieves a leave type by code
func (s *LeaveTypeService) GetLeaveTypeByCode(code string) (*models.LeaveType, error) {
	return s.repo.FindByCode(code)
}

// UpdateLeaveType updates a leave type
func (s *LeaveTypeService) UpdateLeaveType(id uint, req *models.UpdateLeaveTypeRequest, updatedBy *uint) (*models.LeaveType, error) {
	leaveType, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("leave type not found")
	}

	if req.Name != nil {
		leaveType.Name = *req.Name
	}
	if req.Icon != nil {
		leaveType.Icon = req.Icon
	}
	if req.Category != nil {
		leaveType.Category = *req.Category
	}
	if req.PaidLeave != nil {
		leaveType.PaidLeave = *req.PaidLeave
	}
	if req.RequiresDocumentation != nil {
		leaveType.RequiresDocumentation = *req.RequiresDocumentation
	}
	if req.IsActive != nil {
		leaveType.IsActive = *req.IsActive
	}
	if req.Description != nil {
		leaveType.Description = req.Description
	}

	leaveType.UpdatedBy = updatedBy

	if err := s.repo.Update(leaveType); err != nil {
		return nil, fmt.Errorf("failed to update leave type: %w", err)
	}

	return s.repo.FindByID(id)
}

// DeleteLeaveType deletes a leave type
func (s *LeaveTypeService) DeleteLeaveType(id uint) error {
	// TODO: Check if leave type is used in any policies or requests
	return s.repo.Delete(id)
}

// ListLeaveTypes lists leave types with pagination and filters
func (s *LeaveTypeService) ListLeaveTypes(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.LeaveType, int64, error) {
	return s.repo.List(tenantID, page, pageSize, filters)
}

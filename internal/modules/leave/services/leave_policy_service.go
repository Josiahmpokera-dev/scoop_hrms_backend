package services

import (
	"errors"
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/repositories"
)

type LeavePolicyService struct {
	repo         *repositories.LeavePolicyRepository
	leaveTypeRepo *repositories.LeaveTypeRepository
}

func NewLeavePolicyService() *LeavePolicyService {
	return &LeavePolicyService{
		repo:         repositories.NewLeavePolicyRepository(),
		leaveTypeRepo: repositories.NewLeaveTypeRepository(),
	}
}

// CreateLeavePolicy creates a new leave policy
func (s *LeavePolicyService) CreateLeavePolicy(req *models.CreateLeavePolicyRequest, tenantID *uint, createdBy *uint) (*models.LeavePolicy, error) {
	_ = tenantID
	// Validate leave type exists
	_, err := s.leaveTypeRepo.FindByCode(req.LeaveTypeCode)
	if err != nil {
		return nil, errors.New("leave type not found")
	}

	policy := &models.LeavePolicy{
		PolicyName:            req.PolicyName,
		Country:               req.Country,
		LeaveTypeCode:         req.LeaveTypeCode,
		Entitlement:           req.Entitlement,
		AccrualFrequency:      req.AccrualFrequency,
		ProrationOnJoin:       true,
		ProrationOnExit:       true,
		CarryForward:          false,
		EncashmentAllowed:     false,
		NegativeBalanceAllowed: false,
		SandwichRules:         false,
		HalfDayAllowed:        false,
		MinimumNoticeDays:     0,
		IsActive:              true,
		CreatedBy:             createdBy,
		UpdatedBy:             createdBy,
	}

	// Apply optional fields
	if req.ProrationOnJoin != nil {
		policy.ProrationOnJoin = *req.ProrationOnJoin
	}
	if req.ProrationOnExit != nil {
		policy.ProrationOnExit = *req.ProrationOnExit
	}
	if req.CarryForward != nil {
		policy.CarryForward = *req.CarryForward
	}
	if req.CarryForwardLimit != nil {
		policy.CarryForwardLimit = req.CarryForwardLimit
	}
	if req.CarryForwardExpiry != nil {
		policy.CarryForwardExpiry = req.CarryForwardExpiry
	}
	if req.EncashmentAllowed != nil {
		policy.EncashmentAllowed = *req.EncashmentAllowed
	}
	if req.EncashmentLimit != nil {
		policy.EncashmentLimit = req.EncashmentLimit
	}
	if req.NegativeBalanceAllowed != nil {
		policy.NegativeBalanceAllowed = *req.NegativeBalanceAllowed
	}
	if req.SandwichRules != nil {
		policy.SandwichRules = *req.SandwichRules
	}
	if req.HalfDayAllowed != nil {
		policy.HalfDayAllowed = *req.HalfDayAllowed
	}
	if req.MinimumNoticeDays != nil {
		policy.MinimumNoticeDays = *req.MinimumNoticeDays
	}
	if req.MaximumDaysPerRequest != nil {
		policy.MaximumDaysPerRequest = req.MaximumDaysPerRequest
	}
	if req.IsActive != nil {
		policy.IsActive = *req.IsActive
	}

	if err := s.repo.Create(policy); err != nil {
		return nil, fmt.Errorf("failed to create leave policy: %w", err)
	}

	return s.repo.FindByID(policy.ID)
}

// GetLeavePolicyByID retrieves a leave policy by ID
func (s *LeavePolicyService) GetLeavePolicyByID(id uint) (*models.LeavePolicy, error) {
	return s.repo.FindByID(id)
}

// GetPolicyGuidelines retrieves policy guidelines for a leave type
func (s *LeavePolicyService) GetPolicyGuidelines(leaveTypeCode, country string, tenantID *uint) (*models.LeavePolicy, error) {
	return s.repo.FindByCountryAndLeaveType(country, leaveTypeCode, tenantID)
}

// UpdateLeavePolicy updates a leave policy
func (s *LeavePolicyService) UpdateLeavePolicy(id uint, req *models.UpdateLeavePolicyRequest, updatedBy *uint) (*models.LeavePolicy, error) {
	policy, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("leave policy not found")
	}

	if req.PolicyName != nil {
		policy.PolicyName = *req.PolicyName
	}
	if req.Entitlement != nil {
		policy.Entitlement = *req.Entitlement
	}
	if req.AccrualFrequency != nil {
		policy.AccrualFrequency = *req.AccrualFrequency
	}
	if req.ProrationOnJoin != nil {
		policy.ProrationOnJoin = *req.ProrationOnJoin
	}
	if req.ProrationOnExit != nil {
		policy.ProrationOnExit = *req.ProrationOnExit
	}
	if req.CarryForward != nil {
		policy.CarryForward = *req.CarryForward
	}
	if req.CarryForwardLimit != nil {
		policy.CarryForwardLimit = req.CarryForwardLimit
	}
	if req.CarryForwardExpiry != nil {
		policy.CarryForwardExpiry = req.CarryForwardExpiry
	}
	if req.EncashmentAllowed != nil {
		policy.EncashmentAllowed = *req.EncashmentAllowed
	}
	if req.EncashmentLimit != nil {
		policy.EncashmentLimit = req.EncashmentLimit
	}
	if req.NegativeBalanceAllowed != nil {
		policy.NegativeBalanceAllowed = *req.NegativeBalanceAllowed
	}
	if req.SandwichRules != nil {
		policy.SandwichRules = *req.SandwichRules
	}
	if req.HalfDayAllowed != nil {
		policy.HalfDayAllowed = *req.HalfDayAllowed
	}
	if req.MinimumNoticeDays != nil {
		policy.MinimumNoticeDays = *req.MinimumNoticeDays
	}
	if req.MaximumDaysPerRequest != nil {
		policy.MaximumDaysPerRequest = req.MaximumDaysPerRequest
	}
	if req.IsActive != nil {
		policy.IsActive = *req.IsActive
	}

	policy.UpdatedBy = updatedBy

	if err := s.repo.Update(policy); err != nil {
		return nil, fmt.Errorf("failed to update leave policy: %w", err)
	}

	return s.repo.FindByID(id)
}

// DeleteLeavePolicy deletes a leave policy
func (s *LeavePolicyService) DeleteLeavePolicy(id uint) error {
	count, err := s.repo.CountEmployeeAssignments(id)
	if err != nil {
		return fmt.Errorf("failed to check leave policy usage: %w", err)
	}
	if count > 0 {
		return errors.New("Cannot delete leave policy in use")
	}
	return s.repo.Delete(id)
}

// ListLeavePolicies lists leave policies with pagination and filters
func (s *LeavePolicyService) ListLeavePolicies(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.LeavePolicy, int64, error) {
	return s.repo.List(tenantID, page, pageSize, filters)
}

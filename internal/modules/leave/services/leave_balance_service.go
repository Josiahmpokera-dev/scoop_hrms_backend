package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
)

type LeaveBalanceService struct {
	repo         *repositories.LeaveBalanceRepository
	policyRepo   *repositories.LeavePolicyRepository
	leaveTypeRepo *repositories.LeaveTypeRepository
	employeeRepo *employeeRepos.EmployeeRepository
}

func NewLeaveBalanceService() *LeaveBalanceService {
	return &LeaveBalanceService{
		repo:         repositories.NewLeaveBalanceRepository(),
		policyRepo:   repositories.NewLeavePolicyRepository(),
		leaveTypeRepo: repositories.NewLeaveTypeRepository(),
		employeeRepo: employeeRepos.NewEmployeeRepository(),
	}
}

// GetEmployeeLeaveBalances gets leave balances for an employee
func (s *LeaveBalanceService) GetEmployeeLeaveBalances(employeeID string, year int, tenantID *uint) ([]models.LeaveBalance, error) {
	if year == 0 {
		year = time.Now().Year()
	}

	balances, err := s.repo.FindByEmployeeAndYear(employeeID, year, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get leave balances: %w", err)
	}

	// If no balances exist, initialize them based on employee's leave policies
	if len(balances) == 0 {
		// Get employee's leave policy assignments
		// This would require integration with employee_policies table
		// For now, return empty array - balances will be created when first leave is requested
		return []models.LeaveBalance{}, nil
	}

	return balances, nil
}

// InitializeBalance initializes leave balance for an employee based on policy
func (s *LeaveBalanceService) InitializeBalance(employeeID, leaveTypeCode string, year int, tenantID *uint) (*models.LeaveBalance, error) {
	// Get policy for this leave type (would need country from employee)
	// For now, use default policy
	policies, err := s.policyRepo.FindByLeaveTypeCode(leaveTypeCode, tenantID)
	if err != nil || len(policies) == 0 {
		return nil, errors.New("no policy found for this leave type")
	}

	policy := policies[0] // Use first policy found

	balance := &models.LeaveBalance{
		EmployeeID:  employeeID,
		LeaveTypeCode: leaveTypeCode,
		Year:        year,
		Entitlement: float64(policy.Entitlement),
		Used:        0,
		Pending:     0,
		Available:  float64(policy.Entitlement),
		CarriedForward: 0,
	}

	// Set expiry if carry forward is enabled
	if policy.CarryForward && policy.CarryForwardExpiry != nil {
		balance.ExpiresOn = policy.CarryForwardExpiry
	}

	if err := s.repo.CreateOrUpdate(balance); err != nil {
		return nil, fmt.Errorf("failed to initialize balance: %w", err)
	}

	return balance, nil
}

// UpdateBalanceAfterApproval updates balance when leave is approved
func (s *LeaveBalanceService) UpdateBalanceAfterApproval(employeeID, leaveTypeCode string, year int, days float64) error {
	balance, err := s.repo.FindByEmployeeAndTypeAndYear(employeeID, leaveTypeCode, year)
	if err != nil {
		// Balance doesn't exist, initialize it
		// This would need tenantID - for now, skip
		return fmt.Errorf("balance not found: %w", err)
	}

	// Decrement pending and increment used
	if err := s.repo.DecrementPending(employeeID, leaveTypeCode, year, days); err != nil {
		return err
	}
	if err := s.repo.IncrementUsed(employeeID, leaveTypeCode, year, days); err != nil {
		return err
	}

	// Recalculate available
	balance.Used += days
	balance.Pending -= days
	balance.Available = balance.Entitlement + balance.CarriedForward - balance.Used - balance.Pending

	return s.repo.Update(balance)
}

// AddPendingBalance adds days to pending when request is submitted
func (s *LeaveBalanceService) AddPendingBalance(employeeID, leaveTypeCode string, year int, days float64) error {
	balance, err := s.repo.FindByEmployeeAndTypeAndYear(employeeID, leaveTypeCode, year)
	if err != nil {
		// Balance doesn't exist - this should be initialized first
		return fmt.Errorf("balance not found: %w", err)
	}

	if err := s.repo.IncrementPending(employeeID, leaveTypeCode, year, days); err != nil {
		return err
	}

	// Recalculate available
	balance.Pending += days
	balance.Available = balance.Entitlement + balance.CarriedForward - balance.Used - balance.Pending

	return s.repo.Update(balance)
}

// RemovePendingBalance removes days from pending when request is cancelled/rejected
func (s *LeaveBalanceService) RemovePendingBalance(employeeID, leaveTypeCode string, year int, days float64) error {
	return s.repo.DecrementPending(employeeID, leaveTypeCode, year, days)
}

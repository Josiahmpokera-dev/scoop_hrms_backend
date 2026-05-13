package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/repositories"
)

// Default annual leave when no policy row exists (post-probation full year).
const (
	defaultAnnualLeaveCode        = "AL"
	defaultAnnualLeaveEntitlement = 28.0
)

type LeaveBalanceService struct {
	repo             *repositories.LeaveBalanceRepository
	policyRepo       *repositories.LeavePolicyRepository
	leaveTypeRepo    *repositories.LeaveTypeRepository
	employeeRepo     *employeeRepos.EmployeeRepository
	leaveRequestRepo *repositories.LeaveRequestRepository
}

func NewLeaveBalanceService() *LeaveBalanceService {
	return &LeaveBalanceService{
		repo:             repositories.NewLeaveBalanceRepository(),
		policyRepo:       repositories.NewLeavePolicyRepository(),
		leaveTypeRepo:    repositories.NewLeaveTypeRepository(),
		employeeRepo:     employeeRepos.NewEmployeeRepository(),
		leaveRequestRepo: repositories.NewLeaveRequestRepository(),
	}
}

// pickPolicy chooses the best matching active policy for an employee (e.g. by nationality / country).
func (s *LeaveBalanceService) pickPolicy(policies []models.LeavePolicy, leaveTypeCode string, employee *employeeModels.Employee) *models.LeavePolicy {
	var candidates []*models.LeavePolicy
	for i := range policies {
		p := &policies[i]
		if p.LeaveTypeCode != leaveTypeCode {
			continue
		}
		candidates = append(candidates, p)
	}
	if len(candidates) == 0 {
		return nil
	}
	if employee != nil && employee.Nationality != nil {
		nat := strings.ToLower(strings.TrimSpace(*employee.Nationality))
		if nat != "" {
			for _, p := range candidates {
				pc := strings.ToLower(strings.TrimSpace(p.Country))
				if pc != "" && (strings.EqualFold(nat, pc) || strings.Contains(nat, pc) || strings.Contains(pc, nat)) {
					return p
				}
			}
		}
	}
	return candidates[0]
}

// GetEmployeeLeaveBalances returns per–leave-type balances for the year: entitlement from policy
// (full annual entitlement, assuming probation is complete), used/pending from leave requests,
// and available = entitlement + carried_forward - used - pending.
func (s *LeaveBalanceService) GetEmployeeLeaveBalances(employeeID string, year int, tenantID *uint) ([]models.LeaveBalance, error) {
	if year == 0 {
		year = time.Now().Year()
	}

	employee, _ := s.employeeRepo.FindByEmployeeID(employeeID)

	policies, err := s.policyRepo.FindActivePolicies(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to load leave policies: %w", err)
	}

	existingBalances, err := s.repo.FindByEmployeeAndYear(employeeID, year, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get leave balances: %w", err)
	}
	existingByCode := make(map[string]models.LeaveBalance, len(existingBalances))
	for i := range existingBalances {
		b := existingBalances[i]
		existingByCode[b.LeaveTypeCode] = b
	}

	seen := make(map[string]struct{})
	var typeCodes []string
	for i := range policies {
		p := &policies[i]
		if _, ok := seen[p.LeaveTypeCode]; ok {
			continue
		}
		seen[p.LeaveTypeCode] = struct{}{}
		typeCodes = append(typeCodes, p.LeaveTypeCode)
	}

	// If there are stored balance rows but no policy records, still expose those types using DB entitlement.
	for _, b := range existingBalances {
		if _, ok := seen[b.LeaveTypeCode]; !ok {
			seen[b.LeaveTypeCode] = struct{}{}
			typeCodes = append(typeCodes, b.LeaveTypeCode)
		}
	}

	// Always include annual leave so balances are never empty when the employee has no usage yet.
	if _, ok := seen[defaultAnnualLeaveCode]; !ok {
		seen[defaultAnnualLeaveCode] = struct{}{}
		typeCodes = append([]string{defaultAnnualLeaveCode}, typeCodes...)
	}

	out := make([]models.LeaveBalance, 0, len(typeCodes))
	for _, code := range typeCodes {
		policy := s.pickPolicy(policies, code, employee)
		eb, hasExisting := existingByCode[code]

		var entitlement float64
		if policy != nil {
			entitlement = float64(policy.Entitlement)
		} else if hasExisting {
			entitlement = eb.Entitlement
		} else if code == defaultAnnualLeaveCode {
			entitlement = defaultAnnualLeaveEntitlement
		} else {
			continue
		}

		used, err := s.leaveRequestRepo.SumConsumedDaysForYear(employeeID, code, year, tenantID)
		if err != nil {
			return nil, fmt.Errorf("failed to sum used leave days: %w", err)
		}
		pending, err := s.leaveRequestRepo.SumPendingDaysForYear(employeeID, code, year, tenantID)
		if err != nil {
			return nil, fmt.Errorf("failed to sum pending leave days: %w", err)
		}

		carried := 0.0
		if hasExisting {
			carried = eb.CarriedForward
		}

		available := entitlement + carried - used - pending
		negAllowed := policy != nil && policy.NegativeBalanceAllowed
		if !negAllowed && available < 0 {
			available = 0
		}

		bal := models.LeaveBalance{
			EmployeeID:                   employeeID,
			LeaveTypeCode:                code,
			Year:                         year,
			Entitlement:                  entitlement,
			Used:                         used,
			Pending:                      pending,
			Available:                    available,
			CarriedForward:               carried,
			LeaveTakenFromPreviousYear:   0,
			LeaveBalanceFromPreviousYear: 0,
			PreviousDaysUsed:             0,
		}
		if hasExisting {
			bal.ID = eb.ID
			bal.ExpiresOn = eb.ExpiresOn
			bal.LeaveTakenFromPreviousYear = eb.LeaveTakenFromPreviousYear
			bal.LeaveBalanceFromPreviousYear = eb.LeaveBalanceFromPreviousYear
			bal.PreviousDaysUsed = eb.PreviousDaysUsed
			bal.LastAccrualDate = eb.LastAccrualDate
			bal.CreatedAt = eb.CreatedAt
			bal.UpdatedAt = eb.UpdatedAt
		}

		if lt, err := s.leaveTypeRepo.FindByCode(code); err == nil && lt != nil {
			bal.LeaveType = lt
		}

		out = append(out, bal)
	}

	return out, nil
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

	// Calculate entitlement based on probation period for annual leave
	entitlement := float64(policy.Entitlement)

	// For annual leave, check if employee is still in probation
	if leaveTypeCode == "AL" {
		employee, err := s.employeeRepo.FindByEmployeeID(employeeID)
		if err == nil && employee != nil {
			// Check if employee is still in probation
			if employee.ProbationPeriodDays != nil && employee.ExpectedConfirmationDate != nil {
				probationEndDate := employee.ExpectedConfirmationDate
				currentTime := time.Now()

				// If still in probation, set entitlement to 0 (annual leave starts after probation)
				if currentTime.Before(*probationEndDate) {
					entitlement = 0
				}
			}
		}
	}

	balance := &models.LeaveBalance{
		EmployeeID:     employeeID,
		LeaveTypeCode:  leaveTypeCode,
		Year:           year,
		Entitlement:    entitlement,
		Used:           0,
		Pending:        0,
		Available:      entitlement,
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

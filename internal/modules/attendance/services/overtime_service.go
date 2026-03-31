package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
)

// OvertimeService handles overtime business logic
type OvertimeService struct {
	overtimeRepo       *repositories.OvertimeRepository
	employeeRepo       *employeeRepos.EmployeeRepository
	employeeSalaryRepo *employeeRepos.EmployeeSalaryRepository
}

// NewOvertimeService creates a new overtime service
func NewOvertimeService() *OvertimeService {
	return &OvertimeService{
		overtimeRepo:       repositories.NewOvertimeRepository(),
		employeeRepo:       employeeRepos.NewEmployeeRepository(),
		employeeSalaryRepo: employeeRepos.NewEmployeeSalaryRepository(),
	}
}

// --- Policy Management ---

// CreatePolicy creates a new overtime policy (Admin/HR)
func (s *OvertimeService) CreatePolicy(user *userModels.User, req models.CreateOvertimePolicyRequest) (*models.OvertimePolicy, error) {
	if user.Role != userModels.RoleAdmin && user.Role != userModels.RoleHR {
		return nil, errors.New("only Admin or HR can create overtime policies")
	}

	policy := &models.OvertimePolicy{
		Name:                 req.Name,
		DepartmentID:         req.DepartmentID,
		IsActive:             true,
		DailyThresholdHours:  req.DailyThresholdHours,
		WeeklyThresholdHours: req.WeeklyThresholdHours,
		WeekdayMultiplier:    req.WeekdayMultiplier,
		WeekendMultiplier:    req.WeekendMultiplier,
		HolidayMultiplier:    req.HolidayMultiplier,
		MaxDailyOTHours:      req.MaxDailyOTHours,
		MaxWeeklyOTHours:     req.MaxWeeklyOTHours,
		MaxMonthlyOTHours:    req.MaxMonthlyOTHours,
		RequiresPreApproval:  req.RequiresPreApproval,
		AllowCompOff:         req.AllowCompOff,
		CompOffRatio:         req.CompOffRatio,
		CompOffExpiryDays:    req.CompOffExpiryDays,
		MonthlyBudgetCap:     req.MonthlyBudgetCap,
		CreatedBy:            &user.ID,
	}

	if err := s.overtimeRepo.CreatePolicy(policy); err != nil {
		return nil, fmt.Errorf("failed to create overtime policy: %w", err)
	}

	return policy, nil
}

// UpdatePolicy updates an existing OT policy (Admin/HR)
func (s *OvertimeService) UpdatePolicy(user *userModels.User, policyID uint, req models.CreateOvertimePolicyRequest) (*models.OvertimePolicy, error) {
	if user.Role != userModels.RoleAdmin && user.Role != userModels.RoleHR {
		return nil, errors.New("only Admin or HR can update overtime policies")
	}

	policy, err := s.overtimeRepo.FindPolicyByID(policyID)
	if err != nil {
		return nil, errors.New("overtime policy not found")
	}

	policy.Name = req.Name
	policy.DepartmentID = req.DepartmentID
	policy.DailyThresholdHours = req.DailyThresholdHours
	policy.WeeklyThresholdHours = req.WeeklyThresholdHours
	policy.WeekdayMultiplier = req.WeekdayMultiplier
	policy.WeekendMultiplier = req.WeekendMultiplier
	policy.HolidayMultiplier = req.HolidayMultiplier
	policy.MaxDailyOTHours = req.MaxDailyOTHours
	policy.MaxWeeklyOTHours = req.MaxWeeklyOTHours
	policy.MaxMonthlyOTHours = req.MaxMonthlyOTHours
	policy.RequiresPreApproval = req.RequiresPreApproval
	policy.AllowCompOff = req.AllowCompOff
	policy.CompOffRatio = req.CompOffRatio
	policy.CompOffExpiryDays = req.CompOffExpiryDays
	policy.MonthlyBudgetCap = req.MonthlyBudgetCap
	policy.UpdatedBy = &user.ID

	if err := s.overtimeRepo.UpdatePolicy(policy); err != nil {
		return nil, fmt.Errorf("failed to update overtime policy: %w", err)
	}

	return policy, nil
}

// DeletePolicy soft deletes a policy
func (s *OvertimeService) DeletePolicy(user *userModels.User, policyID uint) error {
	if user.Role != userModels.RoleAdmin && user.Role != userModels.RoleHR {
		return errors.New("only Admin or HR can delete overtime policies")
	}
	return s.overtimeRepo.DeletePolicy(policyID)
}

// GetPolicy returns a policy by ID
func (s *OvertimeService) GetPolicy(policyID uint) (*models.OvertimePolicy, error) {
	return s.overtimeRepo.FindPolicyByID(policyID)
}

// ListPolicies returns all policies
func (s *OvertimeService) ListPolicies(page, pageSize int) ([]models.OvertimePolicy, int64, error) {
	return s.overtimeRepo.ListPolicies(page, pageSize)
}

// --- OT Request Management ---

// CreateOTRequest creates a new overtime request
func (s *OvertimeService) CreateOTRequest(user *userModels.User, req models.CreateOvertimeRequestDTO) (*models.OvertimeRequest, error) {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, errors.New("employee record not found for this user")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	// Find applicable policy
	policy, err := s.overtimeRepo.FindPolicyByDepartment(employee.DepartmentID)
	if err != nil {
		return nil, errors.New("no overtime policy found for your department; please contact HR")
	}

	// Validate hours against policy limits
	if req.Hours > policy.MaxDailyOTHours {
		return nil, fmt.Errorf("requested hours (%.1f) exceed daily limit (%.1f)", req.Hours, policy.MaxDailyOTHours)
	}

	// Check weekly limit
	weekday := date.Weekday()
	daysToMonday := int(weekday - time.Monday)
	if daysToMonday < 0 {
		daysToMonday += 7
	}
	weekStart := date.AddDate(0, 0, -daysToMonday)
	weekEnd := weekStart.AddDate(0, 0, 6)
	weeklyHours, err := s.overtimeRepo.GetEmployeeOTHoursInWeek(employee.ID, weekStart, weekEnd)
	if err == nil && (weeklyHours+req.Hours) > policy.MaxWeeklyOTHours {
		return nil, fmt.Errorf("total weekly OT hours (%.1f + %.1f = %.1f) would exceed weekly limit (%.1f)",
			weeklyHours, req.Hours, weeklyHours+req.Hours, policy.MaxWeeklyOTHours)
	}

	// Check monthly limit
	monthlyHours, err := s.overtimeRepo.GetEmployeeOTHoursInMonth(employee.ID, date.Year(), date.Month())
	if err == nil && (monthlyHours+req.Hours) > policy.MaxMonthlyOTHours {
		return nil, fmt.Errorf("total monthly OT hours (%.1f + %.1f = %.1f) would exceed monthly limit (%.1f)",
			monthlyHours, req.Hours, monthlyHours+req.Hours, policy.MaxMonthlyOTHours)
	}

	// Validate compensation type
	if req.CompensationType == "comp_off" && !policy.AllowCompOff {
		return nil, errors.New("comp-off is not allowed under the current policy")
	}

	// Determine multiplier
	otType := models.OvertimeType(req.OvertimeType)
	var multiplier float64
	switch otType {
	case models.OTTypeWeekday:
		multiplier = policy.WeekdayMultiplier
	case models.OTTypeWeekend:
		multiplier = policy.WeekendMultiplier
	case models.OTTypeHoliday:
		multiplier = policy.HolidayMultiplier
	default:
		multiplier = 1.0
	}

	// Calculate compensation
	compType := models.CompensationType(req.CompensationType)
	var payoutAmount, compOffHours *float64
	var compOffExpiry *time.Time

	switch compType {
	case models.CompTypePayout:
		// Payout amount will be calculated later when integrated with payroll
		val := req.Hours * multiplier
		payoutAmount = &val
	case models.CompTypeCompOff:
		val := req.Hours * policy.CompOffRatio
		compOffHours = &val
		expiry := date.AddDate(0, 0, policy.CompOffExpiryDays)
		compOffExpiry = &expiry
	case models.CompTypeBoth:
		// 50-50 split
		payVal := (req.Hours / 2) * multiplier
		compVal := (req.Hours / 2) * policy.CompOffRatio
		payoutAmount = &payVal
		compOffHours = &compVal
		expiry := date.AddDate(0, 0, policy.CompOffExpiryDays)
		compOffExpiry = &expiry
	}

	otRequest := &models.OvertimeRequest{
		EmployeeID:       employee.ID,
		PolicyID:         &policy.ID,
		Date:             date,
		OvertimeType:     otType,
		Hours:            req.Hours,
		Reason:           req.Reason,
		CompensationType: compType,
		Status:           models.OTStatusPending,
		Multiplier:       multiplier,
		PayoutAmount:     payoutAmount,
		CompOffHours:     compOffHours,
		CompOffExpiry:    compOffExpiry,
		CreatedBy:        &user.ID,
	}

	if err := s.overtimeRepo.CreateRequest(otRequest); err != nil {
		return nil, fmt.Errorf("failed to create overtime request: %w", err)
	}

	return otRequest, nil
}

// ApproveOTRequest approves or rejects an OT request (Admin/HR/Manager)
func (s *OvertimeService) ApproveOTRequest(user *userModels.User, requestID uint, req models.ApproveOvertimeRequestDTO) (*models.OvertimeRequest, error) {
	if user.Role != userModels.RoleAdmin && user.Role != userModels.RoleHR {
		return nil, errors.New("only Admin or HR can approve overtime requests")
	}

	otRequest, err := s.overtimeRepo.FindRequestByID(requestID)
	if err != nil {
		return nil, errors.New("overtime request not found")
	}

	if otRequest.Status != models.OTStatusPending {
		return nil, fmt.Errorf("cannot process request in '%s' status", otRequest.Status)
	}

	now := time.Now()
	approverEmployee, _ := s.employeeRepo.FindByUserID(user.ID)
	var approverEmployeeID uint
	if approverEmployee != nil {
		approverEmployeeID = approverEmployee.ID
	}

	if req.Status == "approved" {
		otRequest.Status = models.OTStatusApproved
		otRequest.ApprovedBy = &user.ID
		otRequest.ApprovedAt = &now
	} else {
		otRequest.Status = models.OTStatusRejected
		otRequest.RejectionReason = req.Comments
	}

	// Create approval record
	approval := &models.OvertimeApproval{
		OvertimeRequestID: requestID,
		ApproverID:        approverEmployeeID,
		Level:             1,
		Status:            req.Status,
		Comments:          req.Comments,
		ActionAt:          &now,
	}
	_ = s.overtimeRepo.CreateApproval(approval)

	if err := s.overtimeRepo.UpdateRequest(otRequest); err != nil {
		return nil, fmt.Errorf("failed to update overtime request: %w", err)
	}

	return otRequest, nil
}

// CancelOTRequest cancels an OT request (by employee)
func (s *OvertimeService) CancelOTRequest(user *userModels.User, requestID uint) (*models.OvertimeRequest, error) {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, errors.New("employee record not found for this user")
	}

	otRequest, err := s.overtimeRepo.FindRequestByID(requestID)
	if err != nil {
		return nil, errors.New("overtime request not found")
	}

	if otRequest.EmployeeID != employee.ID {
		return nil, errors.New("you can only cancel your own overtime requests")
	}

	if otRequest.Status != models.OTStatusPending {
		return nil, fmt.Errorf("can only cancel pending requests (current: %s)", otRequest.Status)
	}

	otRequest.Status = models.OTStatusCancelled
	if err := s.overtimeRepo.UpdateRequest(otRequest); err != nil {
		return nil, fmt.Errorf("failed to cancel overtime request: %w", err)
	}

	return otRequest, nil
}

// GetMyOTRequests returns an employee's OT requests
func (s *OvertimeService) GetMyOTRequests(user *userModels.User, status string, page, pageSize int) ([]models.OvertimeRequest, int64, error) {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, 0, errors.New("employee record not found for this user")
	}
	return s.overtimeRepo.ListRequestsByEmployee(employee.ID, status, page, pageSize)
}

// GetOTRequestByID returns an OT request by ID
func (s *OvertimeService) GetOTRequestByID(user *userModels.User, requestID uint) (*models.OvertimeRequest, error) {
	otRequest, err := s.overtimeRepo.FindRequestByID(requestID)
	if err != nil {
		return nil, errors.New("overtime request not found")
	}

	// Check permission
	if user.Role != userModels.RoleAdmin && user.Role != userModels.RoleHR {
		employee, err := s.employeeRepo.FindByUserID(user.ID)
		if err != nil || employee == nil || otRequest.EmployeeID != employee.ID {
			return nil, errors.New("access denied")
		}
	}

	return otRequest, nil
}

// GetPendingApprovals returns pending OT requests (Admin/HR)
func (s *OvertimeService) GetPendingApprovals(page, pageSize int) ([]models.OvertimeRequest, int64, error) {
	return s.overtimeRepo.ListPendingRequests(page, pageSize)
}

// ListAllRequests returns all OT requests with filters (Admin/HR)
func (s *OvertimeService) ListAllRequests(status, startDateStr, endDateStr string, page, pageSize int) ([]models.OvertimeRequest, int64, error) {
	var startDate, endDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &t
		}
	}
	return s.overtimeRepo.ListAllRequests(status, startDate, endDate, page, pageSize)
}

// GetOvertimeStats returns overtime statistics for the current period
func (s *OvertimeService) GetOvertimeStats(startDateStr, endDateStr string) (map[string]interface{}, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = t
		}
	}

	return s.overtimeRepo.GetOvertimeStats(startDate, endDate)
}

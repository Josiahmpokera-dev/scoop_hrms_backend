package services

import (
	"errors"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
)

// ManualPunchService handles manual punch business logic
type ManualPunchService struct {
	manualPunchRepo *repositories.ManualPunchRepository
	employeeRepo    *employeeRepos.EmployeeRepository
}

// NewManualPunchService creates a new manual punch service
func NewManualPunchService() *ManualPunchService {
	return &ManualPunchService{
		manualPunchRepo: repositories.NewManualPunchRepository(),
		employeeRepo:    employeeRepos.NewEmployeeRepository(),
	}
}

// CreateManualPunch creates a new manual punch request
func (s *ManualPunchService) CreateManualPunch(employeeID uint, request models.ManualPunchRequest) (*models.ManualPunch, error) {
	// Validate that employee exists
	employee, err := s.employeeRepo.FindByID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Check if employee already has pending requests (limit to prevent spam)
	pendingCount, err := s.manualPunchRepo.CountPendingByEmployee(employeeID)
	if err != nil {
		return nil, err
	}
	
	if pendingCount >= 5 {
		return nil, errors.New("you have too many pending manual punch requests. Please wait for existing requests to be processed")
	}

	// Use current date and time for the punch
	currentTime := time.Now()
	punchDate := currentTime
	punchTime := currentTime

	// Create the manual punch
	manualPunch := &models.ManualPunch{
		EmployeeID:  employeeID,
		PunchDate:   punchDate,
		PunchTime:   punchTime,
		PunchType:   request.PunchType,
		Reason:      request.Reason,
		SupportingDocument: request.SupportingDocument,
		Status:      models.ManualPunchStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Add employee name for response
	manualPunch.EmployeeName = employee.FirstName + " " + employee.LastName

	err = s.manualPunchRepo.Create(manualPunch)
	if err != nil {
		return nil, err
	}

	return manualPunch, nil
}

// GetManualPunchByID gets a manual punch by ID
func (s *ManualPunchService) GetManualPunchByID(id uint) (*models.ManualPunch, error) {
	return s.manualPunchRepo.FindByID(id)
}

// GetEmployeeManualPunches gets manual punches for an employee with optional filters
func (s *ManualPunchService) GetEmployeeManualPunches(employeeID uint, status models.ManualPunchStatus, startDate, endDate *time.Time) ([]models.ManualPunch, error) {
	return s.manualPunchRepo.FindByEmployee(employeeID, status, startDate, endDate)
}

// GetAllManualPunches gets all manual punches with optional filters and pagination
func (s *ManualPunchService) GetAllManualPunches(filter models.ManualPunchFilter) ([]models.ManualPunch, int64, error) {
	return s.manualPunchRepo.FindAll(filter)
}

// GetPendingManualPunches gets all pending manual punch requests for approval
func (s *ManualPunchService) GetPendingManualPunches() ([]models.ManualPunch, error) {
	return s.manualPunchRepo.GetPendingRequests()
}

// UpdateManualPunchStatus updates the status of a manual punch request
func (s *ManualPunchService) UpdateManualPunchStatus(punchID uint, approver *userModels.User, request models.ManualPunchApprovalRequest) (*models.ManualPunch, error) {
	// Get the manual punch
	manualPunch, err := s.manualPunchRepo.FindByID(punchID)
	if err != nil {
		return nil, errors.New("manual punch request not found")
	}

	// Validate that punch is in pending status
	if manualPunch.Status != models.ManualPunchStatusPending {
		return nil, errors.New("manual punch request has already been processed")
	}

	// Validate rejection reason if status is rejected
	if request.Status == models.ManualPunchStatusRejected && (request.RejectionReason == nil || *request.RejectionReason == "") {
		return nil, errors.New("rejection reason is required when rejecting a manual punch request")
	}

	// Update the status
	var approverID *uint
	if approver != nil {
		approverID = &approver.ID
	}

	err = s.manualPunchRepo.UpdateStatus(punchID, request.Status, approverID, request.RejectionReason)
	if err != nil {
		return nil, err
	}

	// Get the updated punch with approver name if available
	updatedPunch, err := s.manualPunchRepo.FindByID(punchID)
	if err != nil {
		return nil, err
	}

	// Add approver name if available
	if approverID != nil {
		approverEmployee, err := s.employeeRepo.FindByUserID(*approverID)
		if err == nil && approverEmployee != nil {
			updatedPunch.ApproverName = approverEmployee.FirstName + " " + approverEmployee.LastName
		}
	}

	return updatedPunch, nil
}

// CancelManualPunch allows an employee to cancel their own pending manual punch request
func (s *ManualPunchService) CancelManualPunch(punchID uint, employeeID uint) error {
	// Get the manual punch
	manualPunch, err := s.manualPunchRepo.FindByID(punchID)
	if err != nil {
		return errors.New("manual punch request not found")
	}

	// Validate ownership
	if manualPunch.EmployeeID != employeeID {
		return errors.New("you can only cancel your own manual punch requests")
	}

	// Validate that punch is in pending status
	if manualPunch.Status != models.ManualPunchStatusPending {
		return errors.New("only pending manual punch requests can be cancelled")
	}

	// Update status to cancelled
	return s.manualPunchRepo.UpdateStatus(punchID, models.ManualPunchStatusCancelled, nil, nil)
}

// DeleteManualPunch deletes a manual punch request (admin/supervisor only)
func (s *ManualPunchService) DeleteManualPunch(punchID uint) error {
	return s.manualPunchRepo.Delete(punchID)
}

// ValidateManualPunchTime validates if a manual punch time is reasonable
func (s *ManualPunchService) ValidateManualPunchTime(punchDateStr string, punchTimeStr string) error {
	// Parse punch date from string (format: "2006-01-02")
	punchDate, err := time.Parse("2006-01-02", punchDateStr)
	if err != nil {
		return errors.New("invalid punch date format. Use YYYY-MM-DD")
	}

	// Parse punch time from string (format: "15:04:05")
	punchTime, err := time.Parse("15:04:05", punchTimeStr)
	if err != nil {
		return errors.New("invalid punch time format. Use HH:MM:SS")
	}

	// Combine date and time
	punchDateTime := time.Date(
		punchDate.Year(),
		punchDate.Month(),
		punchDate.Day(),
		punchTime.Hour(),
		punchTime.Minute(),
		punchTime.Second(),
		0,
		time.UTC,
	)

	// Check if punch is in the future
	if punchDateTime.After(time.Now()) {
		return errors.New("manual punch cannot be in the future")
	}

	// Check if punch is too far in the past (e.g., more than 30 days)
	if punchDateTime.Before(time.Now().AddDate(0, 0, -30)) {
		return errors.New("manual punch cannot be more than 30 days in the past")
	}

	return nil
}

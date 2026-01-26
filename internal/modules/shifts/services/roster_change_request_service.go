package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
)

type RosterChangeRequestService struct {
	repo         *repositories.RosterChangeRequestRepository
	rosterRepo   *repositories.RosterRepository
	shiftRepo    *repositories.ShiftRepository
	employeeRepo *employeeRepos.EmployeeRepository
}

func NewRosterChangeRequestService() *RosterChangeRequestService {
	return &RosterChangeRequestService{
		repo:         repositories.NewRosterChangeRequestRepository(),
		rosterRepo:   repositories.NewRosterRepository(),
		shiftRepo:    repositories.NewShiftRepository(),
		employeeRepo: employeeRepos.NewEmployeeRepository(),
	}
}

// CreateRosterChangeRequest creates a new roster change request
func (s *RosterChangeRequestService) CreateRosterChangeRequest(req *models.CreateRosterChangeRequestRequest, requestedBy string, tenantID *uint) (*models.RosterChangeRequest, error) {
	// Validate assignment exists
	assignment, err := s.rosterRepo.FindByID(req.AssignmentID)
	if err != nil {
		return nil, errors.New("assignment not found")
	}

	// Validate requester owns the assignment
	if assignment.EmployeeID != requestedBy {
		return nil, errors.New("you can only request changes for your own assignments")
	}

	// Validate assignment is not already changed/swapped
	if assignment.Status == "swapped" {
		return nil, errors.New("assignment has been swapped and cannot be changed")
	}

	// Validate requested shift if provided
	if req.RequestedShiftID != nil {
		_, err := s.shiftRepo.FindByID(*req.RequestedShiftID)
		if err != nil {
			return nil, errors.New("requested shift not found")
		}
	}

	// Parse requested date if provided
	var requestedDate *time.Time
	if req.RequestedDate != nil && *req.RequestedDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.RequestedDate)
		if err != nil {
			return nil, errors.New("invalid requested_date format. Use YYYY-MM-DD")
		}
		requestedDate = &parsedDate
	}

	// Check if there's already a pending request for this assignment
	hasPending, err := s.repo.HasPendingRequest(req.AssignmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing requests: %w", err)
	}
	if hasPending {
		return nil, errors.New("you already have a pending change request for this assignment")
	}

	changeRequest := &models.RosterChangeRequest{
		TenantID:          tenantID,
		RequestedBy:       requestedBy,
		AssignmentID:     req.AssignmentID,
		RequestedShiftID: req.RequestedShiftID,
		RequestedDate:    requestedDate,
		RequestedLocationID: req.RequestedLocationID,
		Reason:           req.Reason,
		Status:           "pending",
	}

	if err := s.repo.Create(changeRequest); err != nil {
		return nil, fmt.Errorf("failed to create roster change request: %w", err)
	}

	return s.repo.FindByID(changeRequest.ID)
}

// GetRosterChangeRequestByID retrieves a roster change request by ID
func (s *RosterChangeRequestService) GetRosterChangeRequestByID(id uint) (*models.RosterChangeRequest, error) {
	return s.repo.FindByID(id)
}

// ListRosterChangeRequests lists roster change requests with pagination and filters
func (s *RosterChangeRequestService) ListRosterChangeRequests(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.RosterChangeRequest, int64, error) {
	return s.repo.List(tenantID, page, pageSize, filters)
}

// ApproveRosterChangeRequest approves a roster change request and updates the assignment
func (s *RosterChangeRequestService) ApproveRosterChangeRequest(id uint, req *models.ApproveRosterChangeRequestRequest, reviewedBy *uint) (*models.RosterChangeRequest, error) {
	changeRequest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("roster change request not found")
	}

	if changeRequest.Status != "pending" {
		return nil, errors.New("only pending change requests can be approved")
	}

	// Get assignment
	assignment, err := s.rosterRepo.FindByID(changeRequest.AssignmentID)
	if err != nil {
		return nil, errors.New("assignment not found")
	}

	// Update assignment with requested changes
	updated := false
	if changeRequest.RequestedShiftID != nil {
		// Validate shift exists
		_, err := s.shiftRepo.FindByID(*changeRequest.RequestedShiftID)
		if err != nil {
			return nil, errors.New("requested shift not found")
		}
		assignment.ShiftID = *changeRequest.RequestedShiftID
		updated = true
	}

	if changeRequest.RequestedDate != nil {
		assignment.Date = *changeRequest.RequestedDate
		updated = true
	}

	if changeRequest.RequestedLocationID != nil {
		assignment.LocationID = changeRequest.RequestedLocationID
		updated = true
	}

	if !updated {
		return nil, errors.New("no changes specified in the request")
	}

	// Update assignment
	if err := s.rosterRepo.Update(assignment); err != nil {
		return nil, fmt.Errorf("failed to update assignment: %w", err)
	}

	// Update change request
	now := time.Now()
	changeRequest.Status = "approved"
	changeRequest.ReviewedAt = &now
	changeRequest.ReviewedBy = reviewedBy

	if err := s.repo.Update(changeRequest); err != nil {
		return nil, fmt.Errorf("failed to update change request: %w", err)
	}

	// TODO: Send notifications if requested
	// if req.NotifyEmployee != nil && *req.NotifyEmployee {
	//     // Implement notification logic
	// }

	return s.repo.FindByID(id)
}

// RejectRosterChangeRequest rejects a roster change request
func (s *RosterChangeRequestService) RejectRosterChangeRequest(id uint, req *models.RejectRosterChangeRequestRequest, reviewedBy *uint) (*models.RosterChangeRequest, error) {
	changeRequest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("roster change request not found")
	}

	if changeRequest.Status != "pending" {
		return nil, errors.New("only pending change requests can be rejected")
	}

	now := time.Now()
	changeRequest.Status = "rejected"
	changeRequest.ReviewedAt = &now
	changeRequest.ReviewedBy = reviewedBy

	if req.RejectionReason != nil {
		changeRequest.RejectionReason = req.RejectionReason
	}

	if err := s.repo.Update(changeRequest); err != nil {
		return nil, fmt.Errorf("failed to update change request: %w", err)
	}

	// TODO: Send notifications if requested
	// if req.NotifyEmployee != nil && *req.NotifyEmployee {
	//     // Implement notification logic
	// }

	return s.repo.FindByID(id)
}

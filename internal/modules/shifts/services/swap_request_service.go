package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
)

type SwapRequestService struct {
	repo         *repositories.SwapRequestRepository
	rosterRepo   *repositories.RosterRepository
	employeeRepo *employeeRepos.EmployeeRepository
}

func NewSwapRequestService() *SwapRequestService {
	return &SwapRequestService{
		repo:         repositories.NewSwapRequestRepository(),
		rosterRepo:   repositories.NewRosterRepository(),
		employeeRepo: employeeRepos.NewEmployeeRepository(),
	}
}

// CreateSwapRequest creates a new swap request
func (s *SwapRequestService) CreateSwapRequest(req *models.CreateSwapRequestRequest, requestedBy string, tenantID *uint) (*models.SwapRequest, error) {
	// Validate assignments exist
	assignment, err := s.rosterRepo.FindByID(req.AssignmentID)
	if err != nil {
		return nil, errors.New("assignment not found")
	}

	swapAssignment, err := s.rosterRepo.FindByID(req.SwapAssignmentID)
	if err != nil {
		return nil, errors.New("swap assignment not found")
	}

	// Validate assignments belong to different employees
	if assignment.EmployeeID == swapAssignment.EmployeeID {
		return nil, errors.New("cannot swap with the same employee")
	}

	// Validate requester owns the assignment
	if assignment.EmployeeID != requestedBy {
		return nil, errors.New("you can only request swaps for your own assignments")
	}

	// Validate assignments are not already swapped
	if assignment.Status == "swapped" {
		return nil, errors.New("assignment is already swapped")
	}
	if swapAssignment.Status == "swapped" {
		return nil, errors.New("swap assignment is already swapped")
	}

	swapRequest := &models.SwapRequest{
		RequestedBy:      requestedBy,
		RequestedWith:    swapAssignment.EmployeeID,
		AssignmentID:     req.AssignmentID,
		SwapAssignmentID: req.SwapAssignmentID,
		Reason:           req.Reason,
		Status:           "pending",
	}

	if err := s.repo.Create(swapRequest); err != nil {
		return nil, fmt.Errorf("failed to create swap request: %w", err)
	}

	return s.repo.FindByID(swapRequest.ID)
}

// GetSwapRequestByID retrieves a swap request by ID
func (s *SwapRequestService) GetSwapRequestByID(id uint) (*models.SwapRequest, error) {
	return s.repo.FindByID(id)
}

// ListSwapRequests lists swap requests with pagination and filters
func (s *SwapRequestService) ListSwapRequests(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.SwapRequest, int64, error) {
	return s.repo.List(tenantID, page, pageSize, filters)
}

// ApproveSwapRequest approves a swap request and updates roster assignments
func (s *SwapRequestService) ApproveSwapRequest(id uint, req *models.ApproveSwapRequestRequest, reviewedBy *uint) (*models.SwapRequest, error) {
	swapRequest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("swap request not found")
	}

	if swapRequest.Status != "pending" {
		return nil, errors.New("only pending swap requests can be approved")
	}

	// Get assignments
	assignment, err := s.rosterRepo.FindByID(swapRequest.AssignmentID)
	if err != nil {
		return nil, errors.New("assignment not found")
	}

	swapAssignment, err := s.rosterRepo.FindByID(swapRequest.SwapAssignmentID)
	if err != nil {
		return nil, errors.New("swap assignment not found")
	}

	// Swap the employees
	tempEmployeeID := assignment.EmployeeID
	assignment.EmployeeID = swapAssignment.EmployeeID
	swapAssignment.EmployeeID = tempEmployeeID

	// Update status
	assignment.Status = "swapped"
	swapAssignment.Status = "swapped"

	// Update assignments
	if err := s.rosterRepo.Update(assignment); err != nil {
		return nil, fmt.Errorf("failed to update assignment: %w", err)
	}

	if err := s.rosterRepo.Update(swapAssignment); err != nil {
		return nil, fmt.Errorf("failed to update swap assignment: %w", err)
	}

	// Update swap request
	now := time.Now()
	swapRequest.Status = "approved"
	swapRequest.ReviewedAt = &now
	swapRequest.ReviewedBy = reviewedBy

	if err := s.repo.Update(swapRequest); err != nil {
		return nil, fmt.Errorf("failed to update swap request: %w", err)
	}

	// TODO: Send notifications if requested
	// if req.NotifyEmployees != nil && *req.NotifyEmployees {
	//     // Implement notification logic
	// }

	return s.repo.FindByID(id)
}

// RejectSwapRequest rejects a swap request
func (s *SwapRequestService) RejectSwapRequest(id uint, req *models.RejectSwapRequestRequest, reviewedBy *uint) (*models.SwapRequest, error) {
	swapRequest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("swap request not found")
	}

	if swapRequest.Status != "pending" {
		return nil, errors.New("only pending swap requests can be rejected")
	}

	now := time.Now()
	swapRequest.Status = "rejected"
	swapRequest.ReviewedAt = &now
	swapRequest.ReviewedBy = reviewedBy

	if req.RejectionReason != nil {
		swapRequest.RejectionReason = req.RejectionReason
	}

	if err := s.repo.Update(swapRequest); err != nil {
		return nil, fmt.Errorf("failed to update swap request: %w", err)
	}

	// TODO: Send notifications if requested
	// if req.NotifyEmployees != nil && *req.NotifyEmployees {
	//     // Implement notification logic
	// }

	return s.repo.FindByID(id)
}

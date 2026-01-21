package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	helpdeskRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/repositories"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
)

// TicketAgentService handles agent/admin ticket management
type TicketAgentService struct {
	ticketRepo  *helpdeskRepos.TicketRepository
	commentRepo *helpdeskRepos.CommentRepository
	userRepo    *userRepos.UserRepository
}

// NewTicketAgentService creates a new ticket agent service
func NewTicketAgentService() *TicketAgentService {
	return &TicketAgentService{
		ticketRepo:  helpdeskRepos.NewTicketRepository(),
		commentRepo: helpdeskRepos.NewCommentRepository(),
		userRepo:    userRepos.NewUserRepository(),
	}
}

// ListTickets lists all tickets for agents/admin
func (s *TicketAgentService) ListTickets(tenantID *uint, userID uint, page, pageSize int, filters map[string]interface{}) ([]models.Ticket, int64, error) {
	// Add user_id to filters for "me" filter
	if assignedTo, ok := filters["assigned_to"].(string); ok && assignedTo == "me" {
		filters["user_id"] = userID
	}

	tickets, total, err := s.ticketRepo.ListAll(tenantID, page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tickets: %w", err)
	}

	return tickets, total, err
}

// AssignTicket assigns or reassigns a ticket
func (s *TicketAgentService) AssignTicket(ticketID uint, assigneeUserID uint, tenantID *uint, assignedBy uint, note *string) (*models.Ticket, error) {
	// Get ticket
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	// Verify tenant ownership
	if tenantID != nil && ticket.TenantID != nil && *ticket.TenantID != *tenantID {
		return nil, errors.New("ticket does not belong to your tenant")
	}

	// Get assignee user
	assignee, err := s.userRepo.FindByID(assigneeUserID)
	if err != nil {
		return nil, errors.New("assignee user not found")
	}

	// Assign ticket
	ticket.AssignedToID = &assigneeUserID
	assigneeName := fmt.Sprintf("%s %s", assignee.FirstName, assignee.LastName)
	if assigneeName == " " {
		assigneeName = assignee.Username
	}
	ticket.AssignedToName = &assigneeName
	ticket.AssignedToEmail = &assignee.Email
	ticket.UpdatedBy = &assignedBy

	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, fmt.Errorf("failed to assign ticket: %w", err)
	}

	// Add assignment note if provided
	if note != nil && *note != "" {
		_ = s.commentRepo.Create(&models.Comment{
			TenantID:   tenantID,
			TicketID:   ticketID,
			AuthorID:   assignedBy,
			AuthorName: "System",
			AuthorType: models.CommentAuthorTypeAgent,
			Text:       *note,
			IsPrivate:  true,
			Timestamp:  time.Now(),
		})
	}

	return ticket, nil
}

// UpdateTicketStatus updates ticket status
func (s *TicketAgentService) UpdateTicketStatus(ticketID uint, status string, tenantID *uint, updatedBy uint, note *string) (*models.Ticket, error) {
	// Get ticket
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	// Verify tenant ownership
	if tenantID != nil && ticket.TenantID != nil && *ticket.TenantID != *tenantID {
		return nil, errors.New("ticket does not belong to your tenant")
	}

	// Validate status transition
	validStatus := models.TicketStatus(status)
	if !isValidStatusTransition(ticket.Status, validStatus) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", ticket.Status, validStatus)
	}

	// Update status
	ticket.Status = validStatus
	ticket.UpdatedBy = &updatedBy

	// Set resolved_at if status is Resolved
	if validStatus == models.TicketStatusResolved {
		now := time.Now()
		ticket.ResolvedAt = &now
		ticket.ResolvedByID = &updatedBy
		if ticket.ResolutionTime == nil {
			ticket.ResolutionTime = &now
		}
	}

	// Set closed_at if status is Closed
	if validStatus == models.TicketStatusClosed {
		now := time.Now()
		ticket.ClosedAt = &now
		ticket.ClosedByID = &updatedBy
	}

	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, fmt.Errorf("failed to update ticket status: %w", err)
	}

	// Add status change note if provided
	if note != nil && *note != "" {
		_ = s.commentRepo.Create(&models.Comment{
			TenantID:   tenantID,
			TicketID:   ticketID,
			AuthorID:   updatedBy,
			AuthorName: "System",
			AuthorType: models.CommentAuthorTypeAgent,
			Text:       *note,
			IsPrivate:  false,
			Timestamp:  time.Now(),
		})
	}

	return ticket, nil
}

// ResolveTicket resolves a ticket
func (s *TicketAgentService) ResolveTicket(ticketID uint, tenantID *uint, resolvedBy uint, resolutionSummary *string, internalNote *string) (*models.Ticket, error) {
	// Get ticket
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	// Verify tenant ownership
	if tenantID != nil && ticket.TenantID != nil && *ticket.TenantID != *tenantID {
		return nil, errors.New("ticket does not belong to your tenant")
	}

	// Resolve ticket
	now := time.Now()
	ticket.Status = models.TicketStatusResolved
	ticket.ResolvedAt = &now
	ticket.ResolvedByID = &resolvedBy
	ticket.ResolutionTime = &now
	ticket.ResolutionSummary = resolutionSummary
	ticket.InternalNote = internalNote
	ticket.UpdatedBy = &resolvedBy

	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, fmt.Errorf("failed to resolve ticket: %w", err)
	}

	return ticket, nil
}

// GetStatistics gets helpdesk statistics
func (s *TicketAgentService) GetStatistics(tenantID *uint, fromDate, toDate *time.Time) (map[string]interface{}, error) {
	return s.ticketRepo.GetStatistics(tenantID, fromDate, toDate)
}

// GetRecentTickets gets recent tickets
func (s *TicketAgentService) GetRecentTickets(tenantID *uint, limit int) ([]models.Ticket, error) {
	return s.ticketRepo.GetRecentTickets(tenantID, limit)
}

// isValidStatusTransition validates status transitions
func isValidStatusTransition(currentStatus, newStatus models.TicketStatus) bool {
	validTransitions := map[models.TicketStatus][]models.TicketStatus{
		models.TicketStatusOpen:       {models.TicketStatusInProgress, models.TicketStatusPending, models.TicketStatusCancelled},
		models.TicketStatusPending:    {models.TicketStatusInProgress, models.TicketStatusResolved},
		models.TicketStatusInProgress: {models.TicketStatusPending, models.TicketStatusResolved},
		models.TicketStatusResolved:   {models.TicketStatusClosed, models.TicketStatusOpen}, // Reopen
		models.TicketStatusClosed:     {}, // Cannot transition from closed
		models.TicketStatusCancelled:  {}, // Cannot transition from cancelled
	}

	allowed, ok := validTransitions[currentStatus]
	if !ok {
		return false
	}

	for _, status := range allowed {
		if status == newStatus {
			return true
		}
	}

	return false
}

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	helpdeskRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	"gorm.io/gorm"
)

// TicketService handles ticket business logic
type TicketService struct {
	ticketRepo     *helpdeskRepos.TicketRepository
	commentRepo    *helpdeskRepos.CommentRepository
	attachmentRepo *helpdeskRepos.AttachmentRepository
	employeeRepo   *employeeRepos.EmployeeRepository
	userRepo       *userRepos.UserRepository
}

// NewTicketService creates a new ticket service
func NewTicketService() *TicketService {
	return &TicketService{
		ticketRepo:     helpdeskRepos.NewTicketRepository(),
		commentRepo:    helpdeskRepos.NewCommentRepository(),
		attachmentRepo: helpdeskRepos.NewAttachmentRepository(),
		employeeRepo:   employeeRepos.NewEmployeeRepository(),
		userRepo:       userRepos.NewUserRepository(),
	}
}

// CreateTicket creates a new ticket for an employee
func (s *TicketService) CreateTicket(userID uint, tenantID *uint, title, description, category string, subCategory *string, priority string, channel string, tags []string, updatedBy *uint) (*models.Ticket, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee record not found for this user")
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Validate priority
	validPriority := models.TicketPriority(priority)
	if validPriority != models.TicketPriorityLow && validPriority != models.TicketPriorityMedium && validPriority != models.TicketPriorityHigh && validPriority != models.TicketPriorityCritical {
		validPriority = models.TicketPriorityMedium // Default
	}

	// Validate channel
	validChannel := models.TicketChannel(channel)
	if validChannel != models.TicketChannelPortal && validChannel != models.TicketChannelEmail && validChannel != models.TicketChannelWhatsApp && validChannel != models.TicketChannelPhone && validChannel != models.TicketChannelChat {
		validChannel = models.TicketChannelPortal // Default
	}

	// Generate ticket number
	ticketNumber, err := s.ticketRepo.GenerateTicketNumber(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ticket number: %w", err)
	}

	// Store tags as JSON
	var tagsJSON *string
	if len(tags) > 0 {
		tagsBytes, _ := json.Marshal(tags)
		tagsStr := string(tagsBytes)
		tagsJSON = &tagsStr
	}

	// Calculate due date based on category (default 48 hours for first response)
	dueDate := time.Now().Add(48 * time.Hour)

	// Create ticket
	ticket := &models.Ticket{
		TenantID:          tenantID,
		TicketNumber:      ticketNumber,
		Title:             title,
		Description:       description,
		Category:          category,
		SubCategory:       subCategory,
		Priority:          validPriority,
		Status:            models.TicketStatusOpen,
		Channel:           validChannel,
		RequesterID:       employee.ID,
		RequesterEmployeeID: &employee.EmployeeID,
		DueDate:           &dueDate,
		TagsJSON:          tagsJSON,
		UpdatedBy:         updatedBy,
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticket, nil
}

// ListMyTickets lists tickets for the authenticated employee
func (s *TicketService) ListMyTickets(userID uint, tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Ticket, int64, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, 0, errors.New("employee record not found for this user")
	}

	// Get tickets
	tickets, total, err := s.ticketRepo.FindByRequesterID(employee.ID, tenantID, page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tickets: %w", err)
	}

	return tickets, total, err
}

// GetTicketDetails gets ticket details with comments and attachments
func (s *TicketService) GetTicketDetails(ticketID uint, userID uint, tenantID *uint, isAgentOrAdmin bool) (*models.Ticket, []models.Comment, []models.Attachment, error) {
	// Get ticket
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, nil, nil, errors.New("ticket not found")
	}

	// Verify tenant ownership
	if tenantID != nil && ticket.TenantID != nil && *ticket.TenantID != *tenantID {
		return nil, nil, nil, errors.New("ticket does not belong to your tenant")
	}

	// Verify ownership (if not agent/admin)
	if !isAgentOrAdmin {
		employee, err := s.employeeRepo.FindByUserID(userID)
		if err != nil || employee.ID != ticket.RequesterID {
			return nil, nil, nil, errors.New("unauthorized: this ticket does not belong to you")
		}
	}

	// Get comments (include private if agent/admin)
	comments, err := s.commentRepo.FindByTicketID(ticketID, isAgentOrAdmin, &userID)
	if err != nil {
		comments = []models.Comment{} // Continue even if comments fail
	}

	// Get attachments
	attachments, err := s.attachmentRepo.FindByTicketID(ticketID)
	if err != nil {
		attachments = []models.Attachment{} // Continue even if attachments fail
	}

	return ticket, comments, attachments, nil
}

// AddComment adds a comment to a ticket
func (s *TicketService) AddComment(ticketID uint, userID uint, tenantID *uint, text string, isPrivate bool, isAgentOrAdmin bool) (*models.Comment, error) {
	// Get ticket
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	// Verify tenant ownership
	if tenantID != nil && ticket.TenantID != nil && *ticket.TenantID != *tenantID {
		return nil, errors.New("ticket does not belong to your tenant")
	}

	// Verify ownership (if not agent/admin)
	if !isAgentOrAdmin {
		employee, err := s.employeeRepo.FindByUserID(userID)
		if err != nil || employee.ID != ticket.RequesterID {
			return nil, errors.New("unauthorized: this ticket does not belong to you")
		}
		// Employees cannot create private comments
		isPrivate = false
	}

	// Get user for author name
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Determine author type
	authorType := models.CommentAuthorTypeRequester
	if isAgentOrAdmin {
		if user.Role == "admin" {
			authorType = models.CommentAuthorTypeAdmin
		} else {
			authorType = models.CommentAuthorTypeAgent
		}
	}

	authorName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	if authorName == " " {
		authorName = user.Username
	}

	// Create comment
	comment := &models.Comment{
		TenantID:   tenantID,
		TicketID:   ticketID,
		AuthorID:   userID,
		AuthorName: authorName,
		AuthorType: authorType,
		Text:       text,
		IsPrivate:  isPrivate,
		Timestamp:  time.Now(),
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Update ticket's first response time if this is the first comment by an agent/admin
	if isAgentOrAdmin && ticket.FirstResponseTime == nil {
		now := time.Now()
		ticket.FirstResponseTime = &now
		s.ticketRepo.Update(ticket)
	}

	return comment, nil
}

// CloseTicket closes a ticket
func (s *TicketService) CloseTicket(ticketID uint, userID uint, tenantID *uint, comment string) error {
	// Get ticket
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	// Verify tenant ownership
	if tenantID != nil && ticket.TenantID != nil && *ticket.TenantID != *tenantID {
		return errors.New("ticket does not belong to your tenant")
	}

	// Verify ownership
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil || employee.ID != ticket.RequesterID {
		return errors.New("unauthorized: this ticket does not belong to you")
	}

	// Check if can be closed
	if ticket.Status == models.TicketStatusClosed || ticket.Status == models.TicketStatusCancelled {
		return errors.New("ticket is already closed or cancelled")
	}

	// Close ticket
	now := time.Now()
	ticket.Status = models.TicketStatusClosed
	ticket.ClosedAt = &now
	ticket.ClosedByID = &userID

	if err := s.ticketRepo.Update(ticket); err != nil {
		return fmt.Errorf("failed to close ticket: %w", err)
	}

	// Add closing comment if provided
	if comment != "" {
		_, _ = s.AddComment(ticketID, userID, tenantID, comment, false, false)
	}

	return nil
}

// SubmitCSAT submits CSAT rating for a ticket
func (s *TicketService) SubmitCSAT(ticketID uint, userID uint, tenantID *uint, rating int, comment *string) error {
	// Get ticket
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	// Verify tenant ownership
	if tenantID != nil && ticket.TenantID != nil && *ticket.TenantID != *tenantID {
		return errors.New("ticket does not belong to your tenant")
	}

	// Verify ownership
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil || employee.ID != ticket.RequesterID {
		return errors.New("unauthorized: this ticket does not belong to you")
	}

	// Validate rating
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	// Update CSAT
	ticket.CSATRating = &rating
	ticket.CSATComment = comment

	if err := s.ticketRepo.Update(ticket); err != nil {
		return fmt.Errorf("failed to submit CSAT: %w", err)
	}

	return nil
}

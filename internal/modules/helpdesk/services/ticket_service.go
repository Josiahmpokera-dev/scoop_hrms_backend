package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	helpdeskRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/repositories"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
)

// AttachmentInput is used when adding attachments to a ticket
type AttachmentInput struct {
	Name string
	URL  string
	Size int64
	Type string
}

// TicketService handles ticket business logic
type TicketService struct {
	ticketRepo     *helpdeskRepos.TicketRepository
	commentRepo    *helpdeskRepos.CommentRepository
	attachmentRepo *helpdeskRepos.AttachmentRepository
	categoryRepo   *helpdeskRepos.TicketCategoryRepository
	employeeRepo   *employeeRepos.EmployeeRepository
	userRepo       *userRepos.UserRepository
}

// NewTicketService creates a new ticket service
func NewTicketService() *TicketService {
	return &TicketService{
		ticketRepo:     helpdeskRepos.NewTicketRepository(),
		commentRepo:    helpdeskRepos.NewCommentRepository(),
		attachmentRepo: helpdeskRepos.NewAttachmentRepository(),
		categoryRepo:   helpdeskRepos.NewTicketCategoryRepository(),
		employeeRepo:   employeeRepos.NewEmployeeRepository(),
		userRepo:       userRepos.NewUserRepository(),
	}
}

// CreateTicket creates a new ticket for any authenticated user (employee or user-only e.g. Admin)
func (s *TicketService) CreateTicket(userID uint, tenantID *uint, title, description, category string, subCategory *string, priority string, channel string, tags []string, updatedBy *uint) (*models.Ticket, error) {
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

	ticket := &models.Ticket{
		TicketNumber: ticketNumber,
		Title:        title,
		Description:  description,
		Category:     category,
		SubCategory:  subCategory,
		Priority:     validPriority,
		Status:       models.TicketStatusOpen,
		Channel:      validChannel,
		DueDate:      &dueDate,
		TagsJSON:     tagsJSON,
		UpdatedBy:    updatedBy,
	}

	// Requester: use Employee if user has one, else User-only (e.g. Admin)
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err == nil {
		ticket.RequesterID = employee.ID
		ticket.RequesterEmployeeID = &employee.EmployeeID
	} else {
		user, uErr := s.userRepo.FindByID(userID)
		if uErr != nil {
			return nil, errors.New("user not found")
		}
		ticket.RequesterID = 0
		ticket.RequesterUserID = &userID
		requesterName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		if requesterName == " " {
			requesterName = user.Username
		}
		ticket.RequesterName = &requesterName
		ticket.RequesterEmail = &user.Email
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticket, nil
}

// ListMyTickets lists tickets for the authenticated user (as requester: employee or user-only)
func (s *TicketService) ListMyTickets(userID uint, tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Ticket, int64, error) {
	var employeeID *uint
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err == nil {
		employeeID = &employee.ID
	}
	tickets, total, err := s.ticketRepo.FindByRequesterUserIDOrEmployee(userID, employeeID, tenantID, page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tickets: %w", err)
	}
	return tickets, total, err
}

// isRequester returns true if the user is the ticket requester (as employee or user-only)
func (s *TicketService) isRequester(ticket *models.Ticket, userID uint) bool {
	if ticket.RequesterUserID != nil && *ticket.RequesterUserID == userID {
		return true
	}
	if ticket.RequesterID == 0 {
		return false
	}
	employee, err := s.employeeRepo.FindByUserID(userID)
	return err == nil && employee.ID == ticket.RequesterID
}

// GetTicketDetails gets ticket details with comments and attachments
func (s *TicketService) GetTicketDetails(ticketID uint, userID uint, tenantID *uint, isAgentOrAdmin bool) (*models.Ticket, []models.Comment, []models.Attachment, error) {
	// Get ticket
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, nil, nil, errors.New("ticket not found")
	}


	// Verify ownership (if not agent/admin): requester is this user (employee or user-only)
	if !isAgentOrAdmin && !s.isRequester(ticket, userID) {
		return nil, nil, nil, errors.New("unauthorized: this ticket does not belong to you")
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

	// Verify ownership (if not agent/admin)
	if !isAgentOrAdmin {
		if !s.isRequester(ticket, userID) {
			return nil, errors.New("unauthorized: this ticket does not belong to you")
		}
		// Non-agent cannot create private comments
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

	// Verify ownership (requester only can close their own ticket)
	if !s.isRequester(ticket, userID) {
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

// GetTicketAttachments returns attachments for a ticket
func (s *TicketService) GetTicketAttachments(ticketID uint) ([]models.Attachment, error) {
	return s.attachmentRepo.FindByTicketID(ticketID)
}

// GetTicketCommentsCount returns the count of (non-private) comments for a ticket
func (s *TicketService) GetTicketCommentsCount(ticketID uint) (int, error) {
	count, err := s.commentRepo.CountByTicketID(ticketID)
	return int(count), err
}

// ListTicketCategories returns all ticket categories for the tenant (for Create Ticket form)
func (s *TicketService) ListTicketCategories(tenantID *uint) ([]models.TicketCategory, error) {
	return s.categoryRepo.ListAll(tenantID)
}

// AddAttachments adds one or more attachments to a ticket (after files are uploaded to storage)
func (s *TicketService) AddAttachments(ticketID uint, userID uint, tenantID *uint, inputs []AttachmentInput) ([]models.Attachment, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket not found")
	}
	if !s.isRequester(ticket, userID) {
		return nil, errors.New("unauthorized: this ticket does not belong to you")
	}

	out := make([]models.Attachment, 0, len(inputs))
	now := time.Now()
	for _, in := range inputs {
		att := &models.Attachment{
			TicketID:   ticketID,
			Name:       in.Name,
			URL:        in.URL,
			Size:       in.Size,
			Type:       in.Type,
			UploadedBy: userID,
			UploadedAt: now,
		}
		if err := s.attachmentRepo.Create(att); err != nil {
			return nil, fmt.Errorf("failed to create attachment: %w", err)
		}
		out = append(out, *att)
	}
	return out, nil
}

// GetRequesterInfo returns requester details for a ticket by employee ID (when requester is an employee)
func (s *TicketService) GetRequesterInfo(employeeID uint) (employeeIDStr, name, email, phone, department string, err error) {
	employee, err := s.employeeRepo.FindByID(employeeID)
	if err != nil {
		return "", "", "", "", "", err
	}
	employeeIDStr = employee.EmployeeID
	name = employee.FirstName + " " + employee.LastName
	if name == " " {
		name = employee.EmployeeID
	}
	if employee.WorkEmail != nil && *employee.WorkEmail != "" {
		email = *employee.WorkEmail
	} else if employee.PersonalEmail != nil && *employee.PersonalEmail != "" {
		email = *employee.PersonalEmail
	}
	if employee.PhoneNumber != nil {
		phone = *employee.PhoneNumber
	}
	return employeeIDStr, name, email, phone, department, nil
}

// GetRequesterInfoForTicket returns requester details from the ticket (employee or user-only e.g. Admin)
func (s *TicketService) GetRequesterInfoForTicket(ticket *models.Ticket) (id, name, email, phone, department string, err error) {
	if ticket.RequesterUserID != nil {
		id = fmt.Sprintf("USER-%d", *ticket.RequesterUserID)
		if ticket.RequesterName != nil {
			name = *ticket.RequesterName
		}
		if ticket.RequesterEmail != nil {
			email = *ticket.RequesterEmail
		}
		return id, name, email, "", "", nil
	}
	if ticket.RequesterID != 0 {
		return s.GetRequesterInfo(ticket.RequesterID)
	}
	return "", "", "", "", "", errors.New("no requester on ticket")
}

// SubmitCSAT submits CSAT rating for a ticket
func (s *TicketService) SubmitCSAT(ticketID uint, userID uint, tenantID *uint, rating int, comment *string) error {
	// Get ticket
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	// Verify ownership (requester only)
	if !s.isRequester(ticket, userID) {
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

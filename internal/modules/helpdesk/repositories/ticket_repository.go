package repositories

import (
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	"gorm.io/gorm"
)

// TicketRepository handles ticket database operations
type TicketRepository struct {
	db *gorm.DB
}

// NewTicketRepository creates a new ticket repository
func NewTicketRepository() *TicketRepository {
	return &TicketRepository{
		db: database.GetDB(),
	}
}

// Create creates a new ticket
func (r *TicketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

// FindByID finds a ticket by ID
func (r *TicketRepository) FindByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.Where("id = ?", id).First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// FindByTicketNumber finds a ticket by ticket number
func (r *TicketRepository) FindByTicketNumber(ticketNumber string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.Where("ticket_number = ?", ticketNumber).First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// FindByRequesterID finds tickets by requester (employee) ID
func (r *TicketRepository) FindByRequesterID(requesterID uint, tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Ticket{}).Where("requester_id = ?", requesterID)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR ticket_number ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if priority, ok := filters["priority"].(string); ok && priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "created_at"
	sortDir := "DESC"
	if sb, ok := filters["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	if sd, ok := filters["sort_dir"].(string); ok && sd != "" {
		sortDir = sd
	}
	orderBy := fmt.Sprintf("%s %s", sortBy, sortDir)

	// Get paginated results
	err := query.Order(orderBy).Offset(offset).Limit(pageSize).Find(&tickets).Error
	return tickets, total, err
}

// ListAll lists all tickets (for agents/admin)
func (r *TicketRepository) ListAll(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Ticket{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR ticket_number ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if priority, ok := filters["priority"].(string); ok && priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if assignedTo, ok := filters["assigned_to"].(string); ok && assignedTo != "" {
		if assignedTo == "me" {
			if userID, ok := filters["user_id"].(uint); ok {
				query = query.Where("assigned_to_id = ?", userID)
			}
		} else if assignedTo == "unassigned" {
			query = query.Where("assigned_to_id IS NULL")
		} else {
			// Assume it's a user ID string
			if userID, ok := filters["assigned_to_user_id"].(uint); ok {
				query = query.Where("assigned_to_id = ?", userID)
			}
		}
	}
	if queue, ok := filters["queue"].(string); ok && queue != "" {
		query = query.Where("queue = ?", queue)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "created_at"
	sortDir := "DESC"
	if sb, ok := filters["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	if sd, ok := filters["sort_dir"].(string); ok && sd != "" {
		sortDir = sd
	}
	orderBy := fmt.Sprintf("%s %s", sortBy, sortDir)

	// Get paginated results
	err := query.Order(orderBy).Offset(offset).Limit(pageSize).Find(&tickets).Error
	return tickets, total, err
}

// Update updates a ticket
func (r *TicketRepository) Update(ticket *models.Ticket) error {
	return r.db.Save(ticket).Error
}

// GenerateTicketNumber generates a unique ticket number
func (r *TicketRepository) GenerateTicketNumber(tenantID *uint) (string, error) {
	var count int64
	query := r.db.Model(&models.Ticket{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Count(&count).Error; err != nil {
		return "", err
	}

	// Format: HD-YYYY-XXX (e.g., HD-2026-001)
	year := time.Now().Format("2006")
	sequence := int(count) + 1
	return "HD-" + year + "-" + formatTicketSequence(sequence), nil
}

// formatTicketSequence formats sequence number with leading zeros
func formatTicketSequence(seq int) string {
	if seq < 10 {
		return "00" + fmt.Sprintf("%d", seq)
	} else if seq < 100 {
		return "0" + fmt.Sprintf("%d", seq)
	}
	return fmt.Sprintf("%d", seq)
}

// GetRecentTickets gets recent tickets
func (r *TicketRepository) GetRecentTickets(tenantID *uint, limit int) ([]models.Ticket, error) {
	var tickets []models.Ticket
	query := r.db.Model(&models.Ticket{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.Order("created_at DESC").Limit(limit).Find(&tickets).Error
	return tickets, err
}

// GetStatistics gets ticket statistics
func (r *TicketRepository) GetStatistics(tenantID *uint, fromDate, toDate *time.Time) (map[string]interface{}, error) {
	query := r.db.Model(&models.Ticket{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if fromDate != nil {
		query = query.Where("created_at >= ?", *fromDate)
	}
	if toDate != nil {
		query = query.Where("created_at <= ?", *toDate)
	}

	stats := make(map[string]interface{})

	// Total tickets
	var totalTickets int64
	query.Count(&totalTickets)
	stats["total_tickets"] = totalTickets

	// Tickets by status
	var statusCounts []struct {
		Status string
		Count  int64
	}
	query.Select("status, COUNT(*) as count").Group("status").Scan(&statusCounts)
	statusMap := make(map[string]int64)
	for _, sc := range statusCounts {
		statusMap[sc.Status] = sc.Count
	}
	stats["tickets_by_status"] = statusMap

	// Tickets by category
	var categoryCounts []struct {
		Category string
		Count    int64
	}
	query.Select("category, COUNT(*) as count").Group("category").Scan(&categoryCounts)
	categoryMap := make(map[string]int64)
	for _, cc := range categoryCounts {
		categoryMap[cc.Category] = cc.Count
	}
	stats["tickets_by_category"] = categoryMap

	// Tickets by priority
	var priorityCounts []struct {
		Priority string
		Count    int64
	}
	query.Select("priority, COUNT(*) as count").Group("priority").Scan(&priorityCounts)
	priorityMap := make(map[string]int64)
	for _, pc := range priorityCounts {
		priorityMap[pc.Priority] = pc.Count
	}
	stats["tickets_by_priority"] = priorityMap

	return stats, nil
}

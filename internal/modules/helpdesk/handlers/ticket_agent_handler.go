package handlers

import (
	"encoding/csv"
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// TicketAgentHandler handles agent/admin ticket HTTP requests
type TicketAgentHandler struct {
	service *services.TicketAgentService
}

// NewTicketAgentHandler creates a new ticket agent handler
func NewTicketAgentHandler() *TicketAgentHandler {
	return &TicketAgentHandler{
		service: services.NewTicketAgentService(),
	}
}

// ListTickets handles listing all tickets for agents
// @Summary List all tickets (Agent/Admin)
// @Description Retrieve a paginated list of all tickets for agent/admin review
// @Tags Helpdesk (Agent/Admin)
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Param search query string false "Search term"
// @Param status query string false "Status filter"
// @Param priority query string false "Priority filter"
// @Param category query string false "Category filter"
// @Param assigned_to query string false "Assigned to (me, unassigned, or user_id)"
// @Param queue query string false "Queue filter"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/agent/tickets [get]
func (h *TicketAgentHandler) ListTickets(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userObj, ok := user.(*userModels.User)
	if !ok {
		response.Unauthorized(c, "Invalid user context")
		return
	}

	tenantID := middleware.GetTenantID(c)

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	filters := map[string]interface{}{}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if priority := c.Query("priority"); priority != "" {
		filters["priority"] = priority
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		filters["assigned_to"] = assignedTo
	}
	if queue := c.Query("queue"); queue != "" {
		filters["queue"] = queue
	}

	tickets, total, err := h.service.ListTickets(tenantID, userObj.ID, page, pageSize, filters)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	formattedTickets := make([]map[string]interface{}, 0, len(tickets))
	for i := range tickets {
		ticket := &tickets[i]
		formattedTickets = append(formattedTickets, FormatTicketResponse(ticket, false, nil, BuildAssignedToFromTicket(ticket)))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	payload := map[string]interface{}{
		"data": formattedTickets,
		"meta": map[string]interface{}{
			"page":        page,
			"per_page":    pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	}
	response.Success(c, "Tickets retrieved successfully", payload)
}

// AssignTicket handles assigning a ticket
// @Summary Assign ticket
// @Description Assign or reassign a ticket to an agent
// @Tags Helpdesk (Agent/Admin)
// @Accept json
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Param request body map[string]interface{} true "Assignment"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/agent/tickets/:ticket_id/assign [post]
func (h *TicketAgentHandler) AssignTicket(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userObj, ok := user.(*userModels.User)
	if !ok {
		response.Unauthorized(c, "Invalid user context")
		return
	}

	tenantID := middleware.GetTenantID(c)

	ticketIDStr := c.Param("ticket_id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ticket ID", nil)
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	assigneeUserIDFloat, ok := req["assignee_user_id"].(float64)
	if !ok {
		response.BadRequest(c, "assignee_user_id is required", nil)
		return
	}
	assigneeUserID := uint(assigneeUserIDFloat)

	var note *string
	if n, ok := req["note"].(string); ok && n != "" {
		note = &n
	}

	ticket, err := h.service.AssignTicket(uint(ticketID), assigneeUserID, tenantID, userObj.ID, note)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	at := BuildAssignedToFromTicket(ticket)
	assignedToMap := map[string]interface{}{"name": at.Name, "email": at.Email, "team": at.Team, "id": at.ID}
	responseData := map[string]interface{}{
		"ticket_id":  ticket.ID,
		"assignedTo": assignedToMap,
		"updatedAt":  ticket.UpdatedAt,
	}
	response.Success(c, "Ticket assigned successfully", responseData)
}

// UpdateTicketStatus handles updating ticket status
// @Summary Update ticket status
// @Description Update the status of a ticket
// @Tags Helpdesk (Agent/Admin)
// @Accept json
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Param request body map[string]interface{} true "Status update"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/agent/tickets/:ticket_id/status [patch]
func (h *TicketAgentHandler) UpdateTicketStatus(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userObj, ok := user.(*userModels.User)
	if !ok {
		response.Unauthorized(c, "Invalid user context")
		return
	}

	tenantID := middleware.GetTenantID(c)

	ticketIDStr := c.Param("ticket_id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ticket ID", nil)
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	status, ok := req["status"].(string)
	if !ok || status == "" {
		response.BadRequest(c, "status is required", nil)
		return
	}

	var note *string
	if n, ok := req["note"].(string); ok && n != "" {
		note = &n
	}

	ticket, err := h.service.UpdateTicketStatus(uint(ticketID), status, tenantID, userObj.ID, note)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	responseData := map[string]interface{}{
		"ticket_id": ticket.ID,
		"status":    string(ticket.Status),
		"updatedAt": ticket.UpdatedAt,
	}
	response.Success(c, "Ticket status updated successfully", responseData)
}

// ResolveTicket handles resolving a ticket
// @Summary Resolve ticket
// @Description Resolve a ticket with resolution summary
// @Tags Helpdesk (Agent/Admin)
// @Accept json
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Param request body map[string]interface{} true "Resolution"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/agent/tickets/:ticket_id/resolve [post]
func (h *TicketAgentHandler) ResolveTicket(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userObj, ok := user.(*userModels.User)
	if !ok {
		response.Unauthorized(c, "Invalid user context")
		return
	}

	tenantID := middleware.GetTenantID(c)

	ticketIDStr := c.Param("ticket_id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ticket ID", nil)
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	var resolutionSummary, internalNote *string
	if rs, ok := req["resolution_summary"].(string); ok && rs != "" {
		resolutionSummary = &rs
	}
	if in, ok := req["internal_note"].(string); ok && in != "" {
		internalNote = &in
	}

	ticket, err := h.service.ResolveTicket(uint(ticketID), tenantID, userObj.ID, resolutionSummary, internalNote)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	responseData := map[string]interface{}{
		"ticket_id":  ticket.ID,
		"status":     string(ticket.Status),
		"resolvedAt": ticket.ResolvedAt,
	}
	response.Success(c, "Ticket resolved successfully", responseData)
}

// GetStatistics handles getting helpdesk statistics
// @Summary Get helpdesk statistics
// @Description Retrieve helpdesk statistics and KPIs
// @Tags Helpdesk (Agent/Admin)
// @Produce json
// @Param date_range query string false "Date range (today, thisWeek, thisMonth, custom)"
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/dashboard/statistics [get]
func (h *TicketAgentHandler) GetStatistics(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	dateRange := c.Query("date_range")
	var fromDate, toDate *time.Time

	now := time.Now()
	switch dateRange {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		fromDate = &start
		toDate = &now
	case "thisWeek":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := now.AddDate(0, 0, -weekday+1)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		fromDate = &start
		toDate = &now
	case "thisMonth":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		fromDate = &start
		toDate = &now
	case "custom":
		if fromStr := c.Query("from"); fromStr != "" {
			if parsed, err := time.Parse("2006-01-02", fromStr); err == nil {
				fromDate = &parsed
			}
		}
		if toStr := c.Query("to"); toStr != "" {
			if parsed, err := time.Parse("2006-01-02", toStr); err == nil {
				toDate = &parsed
			}
		}
	}

	stats, err := h.service.GetStatistics(tenantID, fromDate, toDate)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Map to camelCase and doc shape: totalTickets, openTickets, inProgress, resolved, closed, avgFirstResponseTime, avgResolutionTime, slaCompliance, csatScore, ticketsByCategory, ticketsByPriority, topAgents
	byStatus, _ := stats["tickets_by_status"].(map[string]int64)
	getStatus := func(s string) int {
		if byStatus == nil {
			return 0
		}
		c, _ := byStatus[s]
		return int(c)
	}
	payload := map[string]interface{}{
		"totalTickets":         stats["total_tickets"],
		"openTickets":          getStatus("Open"),
		"inProgress":           getStatus("In Progress"),
		"resolved":             getStatus("Resolved"),
		"closed":               getStatus("Closed"),
		"avgFirstResponseTime": "0 minutes",
		"avgResolutionTime":    "0 hours",
		"slaCompliance":        0,
		"csatScore":            0,
		"ticketsByCategory":    stats["tickets_by_category"],
		"ticketsByPriority":    stats["tickets_by_priority"],
		"topAgents":            []interface{}{},
	}
	response.Success(c, "Helpdesk statistics retrieved successfully", payload)
}

// GetRecentTickets handles getting recent tickets
// @Summary Get recent tickets
// @Description Retrieve recent tickets for dashboard
// @Tags Helpdesk (Agent/Admin)
// @Produce json
// @Param limit query int false "Limit (default: 6)"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/tickets/recent [get]
func (h *TicketAgentHandler) GetRecentTickets(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	limit := 6
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	tickets, err := h.service.GetRecentTickets(tenantID, limit)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	formattedTickets := make([]map[string]interface{}, 0, len(tickets))
	for i := range tickets {
		ticket := &tickets[i]
		ticketMap := FormatTicketResponse(ticket, false, nil, BuildAssignedToFromTicket(ticket))
		formattedTickets = append(formattedTickets, ticketMap)
	}
	response.Success(c, "Recent tickets retrieved successfully", formattedTickets)
}

// ExportReport exports helpdesk report as CSV or XLSX
// GET /helpdesk/reports/export?format=csv&date_range=thisMonth&from=...&to=...
func (h *TicketAgentHandler) ExportReport(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	format := strings.ToLower(c.DefaultQuery("format", "csv"))
	if format != "csv" && format != "xlsx" {
		format = "csv"
	}

	dateRange := c.Query("date_range")
	var fromDate, toDate *time.Time
	now := time.Now()
	switch dateRange {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		fromDate = &start
		toDate = &now
	case "thisWeek":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := now.AddDate(0, 0, -weekday+1)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		fromDate = &start
		toDate = &now
	case "thisMonth":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		fromDate = &start
		toDate = &now
	case "custom":
		if fromStr := c.Query("from"); fromStr != "" {
			if parsed, err := time.Parse("2006-01-02", fromStr); err == nil {
				fromDate = &parsed
			}
		}
		if toStr := c.Query("to"); toStr != "" {
			if parsed, err := time.Parse("2006-01-02", toStr); err == nil {
				toDate = &parsed
			}
		}
	}

	tickets, err := h.service.GetTicketsForExport(tenantID, fromDate, toDate)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	if format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=\"helpdesk_report.csv\"")
		w := csv.NewWriter(c.Writer)
		_ = w.Write([]string{"ID", "Ticket Number", "Title", "Category", "Priority", "Status", "Channel", "Created At", "Updated At"})
		for _, t := range tickets {
			_ = w.Write([]string{
				strconv.FormatUint(uint64(t.ID), 10),
				t.TicketNumber,
				escapeCSV(t.Title),
				t.Category,
				string(t.Priority),
				string(t.Status),
				string(t.Channel),
				t.CreatedAt.Format(time.RFC3339),
				t.UpdatedAt.Format(time.RFC3339),
			})
		}
		w.Flush()
		return
	}

	// XLSX: return CSV with .xlsx filename for now (real XLSX would need excelize)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\"helpdesk_report.xlsx\"")
	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{"ID", "Ticket Number", "Title", "Category", "Priority", "Status", "Channel", "Created At", "Updated At"})
	for _, t := range tickets {
		_ = w.Write([]string{
			strconv.FormatUint(uint64(t.ID), 10),
			t.TicketNumber,
			escapeCSV(t.Title),
			t.Category,
			string(t.Priority),
			string(t.Status),
			string(t.Channel),
			t.CreatedAt.Format(time.RFC3339),
			t.UpdatedAt.Format(time.RFC3339),
		})
	}
	w.Flush()
}

func escapeCSV(s string) string {
	if strings.Contains(s, ",") || strings.Contains(s, "\"") || strings.Contains(s, "\n") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}

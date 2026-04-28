package handlers

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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
// @Router/helpdesk/agent/tickets [get]
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
// @Router/helpdesk/agent/tickets/:ticket_id/assign [post]
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
// @Router/helpdesk/agent/tickets/:ticket_id/status [patch]
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
// @Router/helpdesk/agent/tickets/:ticket_id/resolve [post]
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
// @Param date_range query string false "Date range: today, yesterday, thisWeek, lastWeek, thisMonth, lastMonth, thisYear, lastYear"
// @Param from query string false "Custom start date (YYYY-MM-DD)"
// @Param to query string false "Custom end date (YYYY-MM-DD)"
// @Success 200 {object} response.APIResponse
// @Router/helpdesk/dashboard/statistics [get]
func (h *TicketAgentHandler) GetStatistics(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	fromDate, toDate := parseDateRange(c)

	stats, err := h.service.GetStatistics(tenantID, fromDate, toDate)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Map to camelCase response
	byStatus, _ := stats["tickets_by_status"].(map[string]int64)
	getStatus := func(s string) int {
		if byStatus == nil {
			return 0
		}
		cnt, _ := byStatus[s]
		return int(cnt)
	}

	// Format average times
	avgFirstResponseTime := "0m"
	if minutes, ok := stats["avg_first_response_minutes"].(float64); ok && minutes > 0 {
		avgFirstResponseTime = formatMinutesToDuration(minutes)
	}
	avgResolutionTime := "0h"
	if minutes, ok := stats["avg_resolution_minutes"].(float64); ok && minutes > 0 {
		avgResolutionTime = formatMinutesToDuration(minutes)
	}

	slaCompliance := float64(0)
	if v, ok := stats["sla_compliance"].(float64); ok {
		slaCompliance = roundToOneDecimal(v)
	}
	csatScore := float64(0)
	if v, ok := stats["csat_score"].(float64); ok {
		csatScore = roundToOneDecimal(v)
	}

	// Format top agents
	topAgentsRaw, _ := stats["top_agents"].([]repositories.TopAgentStats)
	topAgents := make([]map[string]interface{}, 0, len(topAgentsRaw))
	for _, a := range topAgentsRaw {
		agentMap := map[string]interface{}{
			"id":       a.ID,
			"name":     a.Name,
			"resolved": a.ResolvedCount,
			"avgTime":  formatMinutesToDuration(a.AvgTimeMinutes),
			"csat":     roundToOneDecimal(a.CSATScore),
		}
		topAgents = append(topAgents, agentMap)
	}

	payload := map[string]interface{}{
		"totalTickets":         stats["total_tickets"],
		"openTickets":          getStatus("Open"),
		"inProgress":           getStatus("In Progress"),
		"resolved":             getStatus("Resolved"),
		"closed":               getStatus("Closed"),
		"avgFirstResponseTime": avgFirstResponseTime,
		"avgResolutionTime":    avgResolutionTime,
		"slaCompliance":        slaCompliance,
		"csatScore":            csatScore,
		"ticketsByCategory":    stats["tickets_by_category"],
		"ticketsByPriority":    stats["tickets_by_priority"],
		"topAgents":            topAgents,
	}
	response.Success(c, "Dashboard statistics retrieved successfully", payload)
}

// parseDateRange parses date_range, from, to query params into time pointers
func parseDateRange(c *gin.Context) (*time.Time, *time.Time) {
	dateRange := c.Query("date_range")
	now := time.Now()
	var fromDate, toDate *time.Time

	switch dateRange {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
		fromDate, toDate = &start, &end
	case "yesterday":
		y := now.AddDate(0, 0, -1)
		start := time.Date(y.Year(), y.Month(), y.Day(), 0, 0, 0, 0, y.Location())
		end := time.Date(y.Year(), y.Month(), y.Day(), 23, 59, 59, 999999999, y.Location())
		fromDate, toDate = &start, &end
	case "thisWeek":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := now.AddDate(0, 0, -weekday+1)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		fromDate = &start
		toDate = &now
	case "lastWeek":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		thisWeekStart := now.AddDate(0, 0, -weekday+1)
		start := thisWeekStart.AddDate(0, 0, -7)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		end := thisWeekStart.AddDate(0, 0, -1)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
		fromDate, toDate = &start, &end
	case "thisMonth":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		fromDate = &start
		toDate = &now
	case "lastMonth":
		firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		start := firstOfThisMonth.AddDate(0, -1, 0)
		end := firstOfThisMonth.AddDate(0, 0, -1)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
		fromDate, toDate = &start, &end
	case "thisYear":
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		fromDate = &start
		toDate = &now
	case "lastYear":
		start := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
		end := time.Date(now.Year()-1, 12, 31, 23, 59, 59, 999999999, now.Location())
		fromDate, toDate = &start, &end
	default:
		// Custom or no date_range: use from/to params
		if fromStr := c.Query("from"); fromStr != "" {
			if parsed, err := time.Parse("2006-01-02", fromStr); err == nil {
				fromDate = &parsed
			}
		}
		if toStr := c.Query("to"); toStr != "" {
			if parsed, err := time.Parse("2006-01-02", toStr); err == nil {
				end := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 999999999, parsed.Location())
				toDate = &end
			}
		}
	}
	return fromDate, toDate
}

// formatMinutesToDuration formats minutes to human-readable string: "45m", "2.4h", "1d 3h"
func formatMinutesToDuration(minutes float64) string {
	if minutes <= 0 {
		return "0m"
	}
	if minutes < 60 {
		return strconv.Itoa(int(minutes)) + "m"
	}
	hours := minutes / 60
	if hours < 24 {
		if hours == float64(int(hours)) {
			return strconv.Itoa(int(hours)) + "h"
		}
		return strconv.FormatFloat(hours, 'f', 1, 64) + "h"
	}
	days := int(hours / 24)
	remainingHours := int(hours) % 24
	if remainingHours == 0 {
		return strconv.Itoa(days) + "d"
	}
	return strconv.Itoa(days) + "d " + strconv.Itoa(remainingHours) + "h"
}

// roundToOneDecimal rounds a float to one decimal place
func roundToOneDecimal(v float64) float64 {
	return float64(int(v*10+0.5)) / 10
}

// GetRecentTickets handles getting recent tickets
// @Summary Get recent tickets
// @Description Retrieve recent tickets for dashboard with requester info
// @Tags Helpdesk (Agent/Admin)
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 50)"
// @Success 200 {object} response.APIResponse
// @Router/helpdesk/tickets/recent [get]
func (h *TicketAgentHandler) GetRecentTickets(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	limit := 10
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

	// Format tickets to match dashboard API spec
	formattedTickets := make([]map[string]interface{}, 0, len(tickets))
	for i := range tickets {
		ticket := &tickets[i]
		
		// Build requester info
		requester := map[string]interface{}{
			"name": "",
		}
		if ticket.RequesterName != nil {
			requester["name"] = *ticket.RequesterName
		}
		// Department from assigned team or empty (would require employee lookup for full dept)
		requester["department"] = ""

		// Build assignedTo info
		var assignedTo interface{}
		if ticket.AssignedToName != nil {
			assignedTo = map[string]interface{}{
				"name": *ticket.AssignedToName,
			}
		}

		ticketMap := map[string]interface{}{
			"id":           ticket.ID,
			"ticketNumber": ticket.TicketNumber,
			"title":        ticket.Title,
			"requester":    requester,
			"category":     ticket.Category,
			"priority":     string(ticket.Priority),
			"status":       string(ticket.Status),
			"assignedTo":   assignedTo,
			"createdAt":    ticket.CreatedAt,
		}
		formattedTickets = append(formattedTickets, ticketMap)
	}
	response.Success(c, "Recent tickets retrieved successfully", formattedTickets)
}

// ExportReport exports helpdesk report as CSV or XLSX with multiple sheets
// GET /helpdesk/reports/export?format=xlsx&date_range=thisMonth&from=...&to=...
func (h *TicketAgentHandler) ExportReport(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	format := strings.ToLower(c.DefaultQuery("format", "xlsx"))
	if format != "csv" && format != "xlsx" {
		format = "xlsx"
	}

	fromDate, toDate := parseDateRange(c)

	tickets, err := h.service.GetTicketsForExport(tenantID, fromDate, toDate)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Get statistics for summary sheet
	stats, _ := h.service.GetStatistics(tenantID, fromDate, toDate)

	filename := fmt.Sprintf("helpdesk-report-%s", time.Now().Format("2006-01-02"))

	if format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.csv\"", filename))
		w := csv.NewWriter(c.Writer)
		_ = w.Write([]string{"Ticket #", "Title", "Requester", "Department", "Category", "Priority", "Status", "Assigned To", "Created At", "Resolved At"})
		for _, t := range tickets {
			requester := ""
			if t.RequesterName != nil {
				requester = *t.RequesterName
			}
			assignedTo := ""
			if t.AssignedToName != nil {
				assignedTo = *t.AssignedToName
			}
			resolvedAt := ""
			if t.ResolvedAt != nil {
				resolvedAt = t.ResolvedAt.Format("2006-01-02 15:04:05")
			}
			_ = w.Write([]string{
				t.TicketNumber,
				escapeCSV(t.Title),
				requester,
				"", // Department would require employee lookup
				t.Category,
				string(t.Priority),
				string(t.Status),
				assignedTo,
				t.CreatedAt.Format("2006-01-02 15:04:05"),
				resolvedAt,
			})
		}
		w.Flush()
		return
	}

	// XLSX with multiple sheets
	f := excelize.NewFile()
	defer f.Close()

	// Sheet 1: Summary
	f.SetSheetName("Sheet1", "Summary")
	summaryData := [][]string{
		{"Metric", "Value"},
		{"Total Tickets", fmt.Sprintf("%v", stats["total_tickets"])},
	}
	if byStatus, ok := stats["tickets_by_status"].(map[string]int64); ok {
		summaryData = append(summaryData, []string{"Open Tickets", fmt.Sprintf("%d", byStatus["Open"])})
		summaryData = append(summaryData, []string{"In Progress", fmt.Sprintf("%d", byStatus["In Progress"])})
		summaryData = append(summaryData, []string{"Resolved", fmt.Sprintf("%d", byStatus["Resolved"])})
		summaryData = append(summaryData, []string{"Closed", fmt.Sprintf("%d", byStatus["Closed"])})
	}
	if v, ok := stats["avg_first_response_minutes"].(float64); ok {
		summaryData = append(summaryData, []string{"Avg First Response", formatMinutesToDuration(v)})
	}
	if v, ok := stats["avg_resolution_minutes"].(float64); ok {
		summaryData = append(summaryData, []string{"Avg Resolution Time", formatMinutesToDuration(v)})
	}
	if v, ok := stats["sla_compliance"].(float64); ok {
		summaryData = append(summaryData, []string{"SLA Compliance", fmt.Sprintf("%.1f%%", v)})
	}
	if v, ok := stats["csat_score"].(float64); ok {
		summaryData = append(summaryData, []string{"CSAT Score", fmt.Sprintf("%.1f", v)})
	}
	for i, row := range summaryData {
		for j, cell := range row {
			cellRef, _ := excelize.CoordinatesToCellName(j+1, i+1)
			f.SetCellValue("Summary", cellRef, cell)
		}
	}

	// Sheet 2: Tickets Detail
	f.NewSheet("Tickets Detail")
	ticketHeaders := []string{"Ticket #", "Title", "Requester", "Department", "Category", "Priority", "Status", "Assigned To", "Created At", "Resolved At"}
	for j, h := range ticketHeaders {
		cellRef, _ := excelize.CoordinatesToCellName(j+1, 1)
		f.SetCellValue("Tickets Detail", cellRef, h)
	}
	for i, t := range tickets {
		requester := ""
		if t.RequesterName != nil {
			requester = *t.RequesterName
		}
		assignedTo := ""
		if t.AssignedToName != nil {
			assignedTo = *t.AssignedToName
		}
		resolvedAt := ""
		if t.ResolvedAt != nil {
			resolvedAt = t.ResolvedAt.Format("2006-01-02 15:04:05")
		}
		row := []string{t.TicketNumber, t.Title, requester, "", t.Category, string(t.Priority), string(t.Status), assignedTo, t.CreatedAt.Format("2006-01-02 15:04:05"), resolvedAt}
		for j, cell := range row {
			cellRef, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue("Tickets Detail", cellRef, cell)
		}
	}

	// Sheet 3: By Category
	f.NewSheet("By Category")
	f.SetCellValue("By Category", "A1", "Category")
	f.SetCellValue("By Category", "B1", "Count")
	f.SetCellValue("By Category", "C1", "Percentage")
	if byCat, ok := stats["tickets_by_category"].(map[string]int64); ok {
		total, _ := stats["total_tickets"].(int64)
		row := 2
		for cat, cnt := range byCat {
			pct := float64(0)
			if total > 0 {
				pct = float64(cnt) / float64(total) * 100
			}
			f.SetCellValue("By Category", fmt.Sprintf("A%d", row), cat)
			f.SetCellValue("By Category", fmt.Sprintf("B%d", row), cnt)
			f.SetCellValue("By Category", fmt.Sprintf("C%d", row), fmt.Sprintf("%.1f%%", pct))
			row++
		}
	}

	// Sheet 4: By Priority
	f.NewSheet("By Priority")
	f.SetCellValue("By Priority", "A1", "Priority")
	f.SetCellValue("By Priority", "B1", "Count")
	f.SetCellValue("By Priority", "C1", "Percentage")
	if byPri, ok := stats["tickets_by_priority"].(map[string]int64); ok {
		total, _ := stats["total_tickets"].(int64)
		row := 2
		for pri, cnt := range byPri {
			pct := float64(0)
			if total > 0 {
				pct = float64(cnt) / float64(total) * 100
			}
			f.SetCellValue("By Priority", fmt.Sprintf("A%d", row), pri)
			f.SetCellValue("By Priority", fmt.Sprintf("B%d", row), cnt)
			f.SetCellValue("By Priority", fmt.Sprintf("C%d", row), fmt.Sprintf("%.1f%%", pct))
			row++
		}
	}

	// Sheet 5: Agent Performance
	f.NewSheet("Agent Performance")
	f.SetCellValue("Agent Performance", "A1", "Agent")
	f.SetCellValue("Agent Performance", "B1", "Tickets Resolved")
	f.SetCellValue("Agent Performance", "C1", "Avg Resolution Time")
	f.SetCellValue("Agent Performance", "D1", "CSAT Score")
	if topAgents, ok := stats["top_agents"].([]repositories.TopAgentStats); ok {
		for i, a := range topAgents {
			row := i + 2
			f.SetCellValue("Agent Performance", fmt.Sprintf("A%d", row), a.Name)
			f.SetCellValue("Agent Performance", fmt.Sprintf("B%d", row), a.ResolvedCount)
			f.SetCellValue("Agent Performance", fmt.Sprintf("C%d", row), formatMinutesToDuration(a.AvgTimeMinutes))
			f.SetCellValue("Agent Performance", fmt.Sprintf("D%d", row), fmt.Sprintf("%.1f", a.CSATScore))
		}
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.xlsx\"", filename))
	if err := f.Write(c.Writer); err != nil {
		response.InternalServerError(c, "Failed to generate report", err.Error())
	}
}

func escapeCSV(s string) string {
	if strings.Contains(s, ",") || strings.Contains(s, "\"") || strings.Contains(s, "\n") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}

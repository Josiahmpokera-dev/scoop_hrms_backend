package handlers

import (
	"strconv"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// TicketHandler handles ticket HTTP requests for employees
type TicketHandler struct {
	service *services.TicketService
}

// NewTicketHandler creates a new ticket handler
func NewTicketHandler() *TicketHandler {
	return &TicketHandler{
		service: services.NewTicketService(),
	}
}

// CreateTicket handles creating a new ticket
// @Summary Create ticket
// @Description Create a new helpdesk ticket
// @Tags Helpdesk (Employee)
// @Accept multipart/form-data
// @Produce json
// @Param category formData string true "Category"
// @Param sub_category formData string false "Sub Category"
// @Param priority formData string true "Priority"
// @Param title formData string true "Title"
// @Param description formData string true "Description"
// @Param channel formData string false "Channel"
// @Param tags formData string false "Tags (comma-separated)"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/tickets [post]
func (h *TicketHandler) CreateTicket(c *gin.Context) {
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

	category := c.PostForm("category")
	subCategory := c.PostForm("sub_category")
	priority := c.PostForm("priority")
	title := c.PostForm("title")
	description := c.PostForm("description")
	channel := c.PostForm("channel")
	tagsStr := c.PostForm("tags")

	if category == "" || priority == "" || title == "" || description == "" {
		response.BadRequest(c, "category, priority, title, and description are required", nil)
		return
	}

	var subCategoryPtr *string
	if subCategory != "" {
		subCategoryPtr = &subCategory
	}

	if channel == "" {
		channel = "Portal"
	}

	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	ticket, err := h.service.CreateTicket(userObj.ID, tenantID, title, description, category, subCategoryPtr, priority, channel, tags, &userObj.ID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Format response
	responseData := map[string]interface{}{
		"id":           ticket.ID,
		"ticket_number": ticket.TicketNumber,
		"status":       string(ticket.Status),
		"created_at":   ticket.CreatedAt,
	}

	response.Success(c, "Ticket created successfully", responseData)
}

// ListMyTickets handles listing tickets for the authenticated employee
// @Summary List my tickets
// @Description Retrieve a paginated list of tickets for the authenticated employee
// @Tags Helpdesk (Employee)
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Param search query string false "Search term"
// @Param status query string false "Status filter"
// @Param priority query string false "Priority filter"
// @Param category query string false "Category filter"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/tickets [get]
func (h *TicketHandler) ListMyTickets(c *gin.Context) {
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
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filters["sort_by"] = sortBy
	}
	if sortDir := c.Query("sort_dir"); sortDir != "" {
		filters["sort_dir"] = sortDir
	}

	tickets, total, err := h.service.ListMyTickets(userObj.ID, tenantID, page, pageSize, filters)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Format response
	formattedTickets := []map[string]interface{}{}
	for _, ticket := range tickets {
		ticketMap := formatTicketResponse(&ticket, false)
		formattedTickets = append(formattedTickets, ticketMap)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Tickets retrieved successfully", formattedTickets, meta)
}

// GetTicketDetails handles getting ticket details
// @Summary Get ticket details
// @Description Retrieve detailed information about a specific ticket
// @Tags Helpdesk (Employee)
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/tickets/:ticket_id [get]
func (h *TicketHandler) GetTicketDetails(c *gin.Context) {
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

	ticket, comments, attachments, err := h.service.GetTicketDetails(uint(ticketID), userObj.ID, tenantID, false)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	// Format response
	responseData := formatTicketResponse(ticket, true)
	
	// Add comments
	formattedComments := []map[string]interface{}{}
	for _, comment := range comments {
		formattedComments = append(formattedComments, map[string]interface{}{
			"id":          comment.ID,
			"author":      comment.AuthorName,
			"author_type": string(comment.AuthorType),
			"text":        comment.Text,
			"timestamp":   comment.Timestamp,
			"is_private":  comment.IsPrivate,
		})
	}
	responseData["comments"] = formattedComments
	responseData["comments_count"] = len(formattedComments)

	// Add attachments
	formattedAttachments := []map[string]interface{}{}
	for _, att := range attachments {
		formattedAttachments = append(formattedAttachments, map[string]interface{}{
			"id":          att.ID,
			"name":        att.Name,
			"url":         att.URL,
			"size":        att.Size,
			"type":        att.Type,
			"uploaded_at": att.UploadedAt,
		})
	}
	responseData["attachments"] = formattedAttachments

	response.Success(c, "Ticket details retrieved successfully", responseData)
}

// AddComment handles adding a comment to a ticket
// @Summary Add comment
// @Description Add a comment to a ticket
// @Tags Helpdesk (Employee)
// @Accept json
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Param request body map[string]interface{} true "Comment"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/tickets/:ticket_id/comments [post]
func (h *TicketHandler) AddComment(c *gin.Context) {
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

	text, ok := req["text"].(string)
	if !ok || text == "" {
		response.BadRequest(c, "text is required", nil)
		return
	}

	isPrivate := false
	if ip, ok := req["is_private"].(bool); ok {
		isPrivate = ip
	}

	comment, err := h.service.AddComment(uint(ticketID), userObj.ID, tenantID, text, isPrivate, false)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	responseData := map[string]interface{}{
		"id":          comment.ID,
		"ticket_id":   comment.TicketID,
		"author":      comment.AuthorName,
		"author_type": string(comment.AuthorType),
		"text":        comment.Text,
		"timestamp":   comment.Timestamp,
		"is_private":  comment.IsPrivate,
	}

	response.Success(c, "Comment posted successfully", responseData)
}

// CloseTicket handles closing a ticket
// @Summary Close ticket
// @Description Close a ticket
// @Tags Helpdesk (Employee)
// @Accept json
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Param request body map[string]interface{} true "Close request"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/tickets/:ticket_id/close [post]
func (h *TicketHandler) CloseTicket(c *gin.Context) {
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
	comment := ""
	if err := c.ShouldBindJSON(&req); err == nil {
		if c, ok := req["comment"].(string); ok {
			comment = c
		}
	}

	err = h.service.CloseTicket(uint(ticketID), userObj.ID, tenantID, comment)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Ticket closed successfully", nil)
}

// SubmitCSAT handles submitting CSAT rating
// @Summary Submit CSAT
// @Description Submit customer satisfaction rating for a ticket
// @Tags Helpdesk (Employee)
// @Accept json
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Param request body map[string]interface{} true "CSAT"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/helpdesk/tickets/:ticket_id/csat [post]
func (h *TicketHandler) SubmitCSAT(c *gin.Context) {
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

	ratingFloat, ok := req["rating"].(float64)
	if !ok {
		response.BadRequest(c, "rating is required", nil)
		return
	}
	rating := int(ratingFloat)

	var comment *string
	if c, ok := req["comment"].(string); ok && c != "" {
		comment = &c
	}

	err = h.service.SubmitCSAT(uint(ticketID), userObj.ID, tenantID, rating, comment)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "CSAT submitted successfully", nil)
}

// formatTicketResponse formats a ticket for response (exported for use in other handlers)
func formatTicketResponse(ticket *models.Ticket, includeDetails bool) map[string]interface{} {
	return formatTicketResponseHelper(ticket, includeDetails)
}

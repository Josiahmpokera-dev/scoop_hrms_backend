package handlers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/storage"
	"github.com/gin-gonic/gin"
)

// TicketHandler handles ticket HTTP requests for employees
type TicketHandler struct {
	service        *services.TicketService
	storageService *storage.StorageService
}

// NewTicketHandler creates a new ticket handler
func NewTicketHandler() *TicketHandler {
	return &TicketHandler{
		service:        services.NewTicketService(),
		storageService: storage.NewStorageService(),
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
// @Router/helpdesk/tickets [post]
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

	responseData := map[string]interface{}{
		"id":           ticket.ID,
		"ticketNumber": ticket.TicketNumber,
		"status":       string(ticket.Status),
		"createdAt":    ticket.CreatedAt,
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
// @Router/helpdesk/tickets [get]
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
	// Map camelCase sort_by from doc to DB column names
	if sortBy := c.Query("sort_by"); sortBy != "" {
		switch sortBy {
		case "createdAt":
			filters["sort_by"] = "created_at"
		case "updatedAt":
			filters["sort_by"] = "updated_at"
		case "dueDate":
			filters["sort_by"] = "due_date"
		default:
			filters["sort_by"] = sortBy
		}
	}
	if sortDir := c.Query("sort_dir"); sortDir != "" {
		filters["sort_dir"] = sortDir
	}

	tickets, total, err := h.service.ListMyTickets(userObj.ID, tenantID, page, pageSize, filters)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Format response with requester enrichment and camelCase
	formattedTickets := make([]map[string]interface{}, 0, len(tickets))
	for i := range tickets {
		ticket := &tickets[i]
		var requester *RequesterInfo
		if eid, name, email, phone, dept, err := h.service.GetRequesterInfoForTicket(ticket); err == nil {
			requester = &RequesterInfo{ID: eid, Name: name, Email: email, Phone: phone, Department: dept}
		}
		ticketMap := FormatTicketResponse(ticket, false, requester, BuildAssignedToFromTicket(ticket))
		// Add list-only fields: attachments summary, commentsCount
		attachments, _ := h.service.GetTicketAttachments(ticket.ID)
		commentsCount, _ := h.service.GetTicketCommentsCount(ticket.ID)
		ticketMap["attachments"] = formatAttachmentsForList(attachments)
		ticketMap["commentsCount"] = commentsCount
		formattedTickets = append(formattedTickets, ticketMap)
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

// GetTicketDetails handles getting ticket details
// @Summary Get ticket details
// @Description Retrieve detailed information about a specific ticket
// @Tags Helpdesk (Employee)
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Success 200 {object} response.APIResponse
// @Router/helpdesk/tickets/:ticket_id [get]
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

	var requester *RequesterInfo
	if eid, name, email, phone, dept, err := h.service.GetRequesterInfoForTicket(ticket); err == nil {
		requester = &RequesterInfo{ID: eid, Name: name, Email: email, Phone: phone, Department: dept}
	}
	responseData := FormatTicketResponse(ticket, true, requester, BuildAssignedToFromTicket(ticket))

	formattedComments := make([]map[string]interface{}, 0, len(comments))
	for _, comment := range comments {
		formattedComments = append(formattedComments, map[string]interface{}{
			"id":         comment.ID,
			"author":     comment.AuthorName,
			"authorType": string(comment.AuthorType),
			"text":       comment.Text,
			"timestamp":  comment.Timestamp,
			"isPrivate":  comment.IsPrivate,
		})
	}
	responseData["comments"] = formattedComments
	responseData["attachments"] = formatAttachmentsForList(attachments)

	response.Success(c, "Ticket retrieved successfully", responseData)
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
// @Router/helpdesk/tickets/:ticket_id/comments [post]
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
		"id":         comment.ID,
		"ticket_id":  comment.TicketID,
		"author":     comment.AuthorName,
		"authorType": string(comment.AuthorType),
		"text":       comment.Text,
		"timestamp":  comment.Timestamp,
		"isPrivate":  comment.IsPrivate,
	}
	response.Success(c, "Comment posted successfully", responseData)
}

// UploadAttachments handles uploading attachments to a ticket (multipart/form-data, attachments[] or file)
// @Summary Upload attachments to ticket
// @Tags Helpdesk (Employee)
// @Accept multipart/form-data
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Param attachments[] formData file true "Files"
// @Success 200 {object} response.APIResponse
// @Router/helpdesk/tickets/:ticket_id/attachments [post]
func (h *TicketHandler) UploadAttachments(c *gin.Context) {
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

	form, err := c.MultipartForm()
	if err != nil {
		response.BadRequest(c, "multipart form required", nil)
		return
	}
	files := form.File["attachments[]"]
	if len(files) == 0 {
		files = form.File["attachments"]
	}
	if len(files) == 0 {
		// Single file as "file"
		if f, _ := c.FormFile("file"); f != nil {
			files = []*multipart.FileHeader{f}
		}
	}
	if len(files) == 0 {
		response.BadRequest(c, "at least one file (attachments[] or file) is required", nil)
		return
	}

	var inputs []services.AttachmentInput
	for _, fileHeader := range files {
		url, size, contentType, uploadErr := h.storageService.UploadFile(fileHeader, "helpdesk", fmt.Sprintf("%d", ticketID))
		if uploadErr != nil {
			response.BadRequest(c, uploadErr.Error(), nil)
			return
		}
		inputs = append(inputs, services.AttachmentInput{
			Name: fileHeader.Filename,
			URL:  h.storageService.ResolveURL(c.Request, url),
			Size: size,
			Type: contentType,
		})
	}

	attachments, err := h.service.AddAttachments(uint(ticketID), userObj.ID, tenantID, inputs)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	formatted := formatAttachmentsForList(attachments)
	response.Success(c, "Attachments uploaded successfully", formatted)
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
// @Router/helpdesk/tickets/:ticket_id/close [post]
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
// @Router/helpdesk/tickets/:ticket_id/csat [post]
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

// GetTicketCategories returns ticket categories for Create Ticket form (dropdowns, SLA info)
// @Summary Get ticket categories
// @Tags Helpdesk (Employee)
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router/helpdesk/ticket-categories [get]
func (h *TicketHandler) GetTicketCategories(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	categories, err := h.service.ListTicketCategories(tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	formatted := make([]map[string]interface{}, 0, len(categories))
	for i := range categories {
		cat := &categories[i]
		sla := map[string]interface{}{}
		if cat.SLAFirstResponseMinutes != nil {
			sla["firstResponse"] = strconv.Itoa(*cat.SLAFirstResponseMinutes)
		}
		if cat.SLAResolutionHours != nil {
			sla["resolution"] = strconv.Itoa(*cat.SLAResolutionHours)
		}
		item := map[string]interface{}{
			"id":              cat.ID,
			"name":            cat.Name,
			"description":     nil,
			"sla":             sla,
			"defaultPriority": nil,
			"requiredFields":  []string{},
			"templates":       []interface{}{},
		}
		if cat.Description != nil {
			item["description"] = *cat.Description
		}
		if cat.DefaultPriority != nil {
			item["defaultPriority"] = *cat.DefaultPriority
		}
		if cat.RequiredFieldsJSON != nil && *cat.RequiredFieldsJSON != "" {
			var rf []string
			_ = json.Unmarshal([]byte(*cat.RequiredFieldsJSON), &rf)
			item["requiredFields"] = rf
		}
		if cat.TemplatesJSON != nil && *cat.TemplatesJSON != "" {
			var tpl []map[string]interface{}
			_ = json.Unmarshal([]byte(*cat.TemplatesJSON), &tpl)
			item["templates"] = tpl
		}
		formatted = append(formatted, item)
	}
	response.Success(c, "Ticket categories retrieved successfully", formatted)
}

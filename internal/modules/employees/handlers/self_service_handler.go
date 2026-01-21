package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// SelfServiceHandler handles self-service HTTP requests
type SelfServiceHandler struct {
	service *services.SelfServiceService
}

// NewSelfServiceHandler creates a new self-service handler
func NewSelfServiceHandler() *SelfServiceHandler {
	return &SelfServiceHandler{
		service: services.NewSelfServiceService(),
	}
}

// GetProfile handles getting employee profile
// @Summary Get employee profile
// @Description Retrieve complete profile information for the authenticated employee
// @Tags Self-Service
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/self-service/profile [get]
func (h *SelfServiceHandler) GetProfile(c *gin.Context) {
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

	profile, err := h.service.GetEmployeeProfile(userObj.ID, tenantID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Profile retrieved successfully", profile)
}

// UpdateProfile handles profile update request
// @Summary Update profile information
// @Description Submit a profile update request (requires manager approval)
// @Tags Self-Service
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Update request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/self-service/profile/update [post]
func (h *SelfServiceHandler) UpdateProfile(c *gin.Context) {
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

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	section, ok := req["section"].(string)
	if !ok {
		response.BadRequest(c, "section is required", nil)
		return
	}

	updates, ok := req["updates"].(map[string]interface{})
	if !ok {
		response.BadRequest(c, "updates are required", nil)
		return
	}

	var reason *string
	if r, ok := req["reason"].(string); ok {
		reason = &r
	}

	tenantID := middleware.GetTenantID(c)
	updatedBy := &userObj.ID

	updateRequest, err := h.service.CreateProfileUpdateRequest(userObj.ID, tenantID, section, updates, reason, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Build response
	responseData := map[string]interface{}{
		"update_request_id": updateRequest.UpdateRequestID,
		"section":          string(updateRequest.Section),
		"status":           string(updateRequest.Status),
		"submitted_at":     updateRequest.SubmittedAt,
	}

	// Get approver info
	if updateRequest.ApproverID != nil {
		employeeService := services.NewEmployeeService()
		manager, _ := employeeService.GetEmployeeByID(*updateRequest.ApproverID)
		if manager != nil {
			responseData["approver"] = map[string]interface{}{
				"id":    manager.ID,
				"name":  manager.FullName(),
				"email": manager.WorkEmail,
			}
		}
	}

	responseData["estimated_processing_time"] = "24-48 hours"

	response.Success(c, "Profile update request submitted successfully", responseData)
}

// GetProfileUpdateStatus handles getting profile update status
// @Summary Get profile update status
// @Description Get the status of pending profile update requests
// @Tags Self-Service
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/profile/update-status [get]
func (h *SelfServiceHandler) GetProfileUpdateStatus(c *gin.Context) {
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

	status, err := h.service.GetProfileUpdateStatus(userObj.ID, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Update status retrieved successfully", status)
}

// GetDocuments handles getting employee documents
// @Summary Get employee documents
// @Description Retrieve all documents associated with the employee
// @Tags Self-Service
// @Produce json
// @Param document_type query string false "Document type filter"
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/profile/documents [get]
func (h *SelfServiceHandler) GetDocuments(c *gin.Context) {
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

	// Parse query parameters
	documentType := c.Query("document_type")
	var docType *string
	if documentType != "" {
		docType = &documentType
	}

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

	documents, total, err := h.service.GetEmployeeDocuments(userObj.ID, tenantID, docType, page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Build response with pagination
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Documents retrieved successfully", documents, meta)
}

// ListServiceRequests handles getting service requests
// @Summary Get service requests
// @Description Retrieve list of service requests (HR Letters, IT Requests, Facilities)
// @Tags Self-Service
// @Produce json
// @Param type query string false "Request type"
// @Param status query string false "Request status"
// @Param priority query string false "Request priority"
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/requests [get]
func (h *SelfServiceHandler) ListServiceRequests(c *gin.Context) {
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

	// Parse query parameters
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
	if requestType := c.Query("type"); requestType != "" {
		filters["type"] = requestType
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if priority := c.Query("priority"); priority != "" {
		filters["priority"] = priority
	}

	requests, total, err := h.service.ListServiceRequests(userObj.ID, tenantID, page, pageSize, filters)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Format response
	formattedRequests := []map[string]interface{}{}
	for _, req := range requests {
		reqMap := map[string]interface{}{
			"id":             req.ID,
			"request_number": req.RequestNumber,
			"type":           string(req.Type),
			"category":       req.Category,
			"subject":        req.Subject,
			"status":         string(req.Status),
			"priority":       string(req.Priority),
			"requested_date": req.RequestedDate,
			"completed_date": req.CompletedDate,
			"sla_hours":      req.SLAHours,
			"is_downloadable": req.DocumentURL != nil && *req.DocumentURL != "",
		}

		if req.AssignedTo != nil {
			reqMap["assigned_to"] = map[string]interface{}{
				"name": *req.AssignedTo,
			}
		}

		formattedRequests = append(formattedRequests, reqMap)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Service requests retrieved successfully", formattedRequests, meta)
}

// GetServiceRequestDetails handles getting service request details
// @Summary Get service request details
// @Description Retrieve detailed information about a specific service request
// @Tags Self-Service
// @Produce json
// @Param request_id path int true "Request ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/requests/:request_id [get]
func (h *SelfServiceHandler) GetServiceRequestDetails(c *gin.Context) {
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

	requestIDStr := c.Param("request_id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	request, err := h.service.GetServiceRequestDetails(userObj.ID, uint(requestID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	// Format response
	responseData := map[string]interface{}{
		"id":             request.ID,
		"request_number": request.RequestNumber,
		"type":           string(request.Type),
		"category":       request.Category,
		"subject":        request.Subject,
		"description":    request.Description,
		"status":         string(request.Status),
		"priority":       string(request.Priority),
		"requested_date": request.RequestedDate,
		"completed_date": request.CompletedDate,
		"sla_hours":      request.SLAHours,
		"is_downloadable": request.DocumentURL != nil && *request.DocumentURL != "",
	}

	if request.AssignedTo != nil {
		responseData["assigned_to"] = map[string]interface{}{
			"name": *request.AssignedTo,
		}
	}

	if request.DocumentURL != nil {
		responseData["document"] = map[string]interface{}{
			"file_url": *request.DocumentURL,
			"file_name": request.FileName,
		}
	}

	response.Success(c, "Request details retrieved successfully", responseData)
}

// CreateServiceRequest handles creating a service request
// @Summary Create service request
// @Description Create a new service request (HR Letter, IT Request, or Facilities)
// @Tags Self-Service
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Service request"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/requests [post]
func (h *SelfServiceHandler) CreateServiceRequest(c *gin.Context) {
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

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	requestType, ok := req["type"].(string)
	if !ok {
		response.BadRequest(c, "type is required", nil)
		return
	}

	category, ok := req["category"].(string)
	if !ok {
		response.BadRequest(c, "category is required", nil)
		return
	}

	subject, ok := req["subject"].(string)
	if !ok {
		response.BadRequest(c, "subject is required", nil)
		return
	}

	var description *string
	if d, ok := req["description"].(string); ok {
		description = &d
	}

	priority := "medium"
	if p, ok := req["priority"].(string); ok {
		priority = p
	}

	// Extract additional data based on type
	additionalData := map[string]interface{}{}
	if requestType == "hr_letter" {
		if letterType, ok := req["letter_type"].(string); ok {
			additionalData["letter_type"] = letterType
		}
		if purpose, ok := req["purpose"].(string); ok {
			additionalData["purpose"] = purpose
		}
		if addressedTo, ok := req["addressed_to"].(string); ok {
			additionalData["addressed_to"] = addressedTo
		}
		if notes, ok := req["additional_notes"].(string); ok {
			additionalData["additional_notes"] = notes
		}
	} else {
		if items, ok := req["requested_items"].([]interface{}); ok {
			additionalData["requested_items"] = items
		}
	}

	tenantID := middleware.GetTenantID(c)
	updatedBy := &userObj.ID

	serviceRequest, err := h.service.CreateServiceRequest(userObj.ID, tenantID, requestType, category, subject, description, priority, additionalData, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Format response
	responseData := map[string]interface{}{
		"id":             serviceRequest.ID,
		"request_number": serviceRequest.RequestNumber,
		"type":           string(serviceRequest.Type),
		"category":       serviceRequest.Category,
		"subject":        serviceRequest.Subject,
		"status":         string(serviceRequest.Status),
		"priority":       string(serviceRequest.Priority),
		"requested_date": serviceRequest.RequestedDate,
		"sla_hours":      serviceRequest.SLAHours,
	}

	if serviceRequest.AssignedTo != nil {
		responseData["assigned_to"] = map[string]interface{}{
			"name": *serviceRequest.AssignedTo,
		}
	}

	if serviceRequest.EstimatedCompletionDate != nil {
		responseData["estimated_completion_date"] = *serviceRequest.EstimatedCompletionDate
	}

	response.Success(c, "Service request created successfully", responseData)
}

// CancelServiceRequest handles cancelling a service request
// @Summary Cancel service request
// @Description Cancel a pending service request
// @Tags Self-Service
// @Accept json
// @Produce json
// @Param request_id path int true "Request ID"
// @Param request body map[string]interface{} true "Cancel request"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/requests/:request_id/cancel [post]
func (h *SelfServiceHandler) CancelServiceRequest(c *gin.Context) {
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

	requestIDStr := c.Param("request_id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	reason := "No reason provided"
	if r, ok := req["reason"].(string); ok {
		reason = r
	}

	err = h.service.CancelServiceRequest(userObj.ID, uint(requestID), reason)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Service request cancelled successfully", nil)
}

// SearchDirectory handles searching employees in the directory
// @Summary Search employees directory
// @Description Search and filter employees in the organization directory
// @Tags Self-Service
// @Produce json
// @Param search query string false "Search term"
// @Param department_id query int false "Department ID"
// @Param position_id query int false "Position ID"
// @Param location_id query int false "Location ID"
// @Param status query string false "Status"
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/directory [get]
func (h *SelfServiceHandler) SearchDirectory(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	// Parse query parameters
	search := c.Query("search")
	var searchPtr *string
	if search != "" {
		searchPtr = &search
	}

	var departmentID *uint
	if dID := c.Query("department_id"); dID != "" {
		if parsed, err := strconv.ParseUint(dID, 10, 32); err == nil {
			id := uint(parsed)
			departmentID = &id
		}
	}

	var positionID *uint
	if pID := c.Query("position_id"); pID != "" {
		if parsed, err := strconv.ParseUint(pID, 10, 32); err == nil {
			id := uint(parsed)
			positionID = &id
		}
	}

	var locationID *uint
	if lID := c.Query("location_id"); lID != "" {
		if parsed, err := strconv.ParseUint(lID, 10, 32); err == nil {
			id := uint(parsed)
			locationID = &id
		}
	}

	status := c.Query("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

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

	employees, total, err := h.service.SearchEmployees(tenantID, searchPtr, departmentID, positionID, locationID, statusPtr, page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Format response
	formattedEmployees := []map[string]interface{}{}
	for _, emp := range employees {
		empMap := map[string]interface{}{
			"id":          emp.ID,
			"employee_id": emp.EmployeeID,
			"full_name":   emp.FullName(),
			"first_name":  emp.FirstName,
			"last_name":   emp.LastName,
			"status":      string(emp.Status),
		}

		if emp.PhotoURL != nil {
			empMap["photo"] = *emp.PhotoURL
		}

		if emp.WorkEmail != nil {
			empMap["official_email"] = *emp.WorkEmail
		}

		if emp.WorkPhone != nil {
			empMap["work_phone"] = *emp.WorkPhone
		}

		formattedEmployees = append(formattedEmployees, empMap)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Directory search completed successfully", formattedEmployees, meta)
}

// GetDirectoryEmployeeDetails handles getting employee directory details
// @Summary Get employee directory details
// @Description Get detailed information about a specific employee from the directory
// @Tags Self-Service
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/directory/:employee_id [get]
func (h *SelfServiceHandler) GetDirectoryEmployeeDetails(c *gin.Context) {
	employeeID := c.Param("employee_id")
	tenantID := middleware.GetTenantID(c)

	details, err := h.service.GetEmployeeDirectoryDetails(employeeID, tenantID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Employee directory details retrieved successfully", details)
}

// DownloadIDCard handles downloading employee ID card
// @Summary Download employee ID card
// @Description Generate and download employee ID card (placeholder - not yet implemented)
// @Tags Self-Service
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/profile/id-card [get]
func (h *SelfServiceHandler) DownloadIDCard(c *gin.Context) {
	response.BadRequest(c, "ID card generation is not yet implemented. This feature will be available in a future update.", nil)
}

// ListPayslips handles getting payslips list
// @Summary Get payslips list
// @Description Retrieve list of payslips for the authenticated employee (placeholder - requires payroll module)
// @Tags Self-Service
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/payslips [get]
func (h *SelfServiceHandler) ListPayslips(c *gin.Context) {
	response.BadRequest(c, "Payslip functionality is not yet implemented. This feature requires the payroll module to be integrated.", nil)
}

// GetPayslipDetails handles getting payslip details
// @Summary Get payslip details
// @Description Retrieve detailed information about a specific payslip (placeholder - requires payroll module)
// @Tags Self-Service
// @Produce json
// @Param payslip_id path int true "Payslip ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/payslips/:payslip_id [get]
func (h *SelfServiceHandler) GetPayslipDetails(c *gin.Context) {
	response.BadRequest(c, "Payslip functionality is not yet implemented. This feature requires the payroll module to be integrated.", nil)
}

// DownloadPayslip handles downloading payslip PDF
// @Summary Download payslip PDF
// @Description Download payslip as PDF (placeholder - requires payroll module)
// @Tags Self-Service
// @Produce application/pdf
// @Param payslip_id path int true "Payslip ID"
// @Success 200 {file} application/pdf
// @Router /api/v1/self-service/payslips/:payslip_id/download [get]
func (h *SelfServiceHandler) DownloadPayslip(c *gin.Context) {
	response.BadRequest(c, "Payslip download is not yet implemented. This feature requires the payroll module to be integrated.", nil)
}

// EmailPayslip handles emailing payslip
// @Summary Email payslip
// @Description Send payslip via email (placeholder - requires payroll module)
// @Tags Self-Service
// @Accept json
// @Produce json
// @Param payslip_id path int true "Payslip ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/payslips/:payslip_id/email [post]
func (h *SelfServiceHandler) EmailPayslip(c *gin.Context) {
	response.BadRequest(c, "Payslip email functionality is not yet implemented. This feature requires the payroll module to be integrated.", nil)
}

// GetYTDSummary handles getting year-to-date summary
// @Summary Get YTD summary
// @Description Retrieve year-to-date salary summary (placeholder - requires payroll module)
// @Tags Self-Service
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/payslips/ytd-summary [get]
func (h *SelfServiceHandler) GetYTDSummary(c *gin.Context) {
	response.BadRequest(c, "YTD summary is not yet implemented. This feature requires the payroll module to be integrated.", nil)
}

// RaiseSalaryQuery handles raising salary query
// @Summary Raise salary query
// @Description Submit a query about salary/payslip (placeholder - requires payroll module)
// @Tags Self-Service
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Query request"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/payslips/query [post]
func (h *SelfServiceHandler) RaiseSalaryQuery(c *gin.Context) {
	response.BadRequest(c, "Salary query functionality is not yet implemented. This feature requires the payroll module to be integrated.", nil)
}

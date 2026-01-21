package handlers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// AssetSelfServiceHandler handles self-service asset HTTP requests
type AssetSelfServiceHandler struct {
	service *services.AssetSelfService
}

// NewAssetSelfServiceHandler creates a new asset self-service handler
func NewAssetSelfServiceHandler() *AssetSelfServiceHandler {
	return &AssetSelfServiceHandler{
		service: services.NewAssetSelfService(),
	}
}

// GetMyAssignedAssets handles getting assigned assets
// @Summary Get my assigned assets
// @Description Retrieve all assets assigned to the authenticated employee
// @Tags Self-Service Assets
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/self-service/assets [get]
func (h *AssetSelfServiceHandler) GetMyAssignedAssets(c *gin.Context) {
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

	assets, err := h.service.GetMyAssignedAssets(userObj.ID, tenantID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	// Format response
	formattedAssets := []map[string]interface{}{}
	for _, asset := range assets {
		assetMap := map[string]interface{}{
			"id":             asset.ID,
			"asset_code":    asset.AssetCode,
			"asset_type":    asset.AssetType,
			"brand":         asset.Brand,
			"model":         asset.Model,
			"serial_number": asset.SerialNumber,
			"status":        asset.Status,
			"condition":     asset.Condition,
			"assigned_date": asset.AssignedDate,
		}

		if asset.PurchaseDate != nil {
			assetMap["purchase_date"] = asset.PurchaseDate
		}
		if asset.WarrantyExpiry != nil {
			assetMap["warranty_expiry"] = asset.WarrantyExpiry
		}
		if asset.Value != nil {
			assetMap["value"] = asset.Value
		}
		if asset.Notes != nil {
			assetMap["notes"] = asset.Notes
		}

		formattedAssets = append(formattedAssets, assetMap)
	}

	response.Success(c, "Assigned assets retrieved successfully", formattedAssets)
}

// CreateAssetRequest handles creating an asset request
// @Summary Create asset request
// @Description Request a new asset, replacement, or additional asset from HR
// @Tags Self-Service Assets
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Asset request"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/assets/requests [post]
func (h *AssetSelfServiceHandler) CreateAssetRequest(c *gin.Context) {
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

	requestType, ok := req["request_type"].(string)
	if !ok {
		response.BadRequest(c, "request_type is required", nil)
		return
	}

	assetType, ok := req["asset_type"].(string)
	if !ok {
		response.BadRequest(c, "asset_type is required", nil)
		return
	}

	justification, ok := req["justification"].(string)
	if !ok || justification == "" {
		response.BadRequest(c, "justification is required", nil)
		return
	}

	var brand, model *string
	if b, ok := req["brand"].(string); ok && b != "" {
		brand = &b
	}
	if m, ok := req["model"].(string); ok && m != "" {
		model = &m
	}

	priority := "medium"
	if p, ok := req["priority"].(string); ok {
		priority = p
	}

	var replacingAssetID *uint
	if rID, ok := req["replacing_asset_id"].(float64); ok && rID > 0 {
		id := uint(rID)
		replacingAssetID = &id
	}

	var notes *string
	if n, ok := req["notes"].(string); ok {
		notes = &n
	}

	tenantID := middleware.GetTenantID(c)
	updatedBy := &userObj.ID

	assetRequest, err := h.service.CreateAssetRequest(userObj.ID, tenantID, requestType, assetType, justification, brand, model, priority, replacingAssetID, notes, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Format response
	responseData := map[string]interface{}{
		"id":             assetRequest.ID,
		"request_number": assetRequest.RequestNumber,
		"request_type":   string(assetRequest.RequestType),
		"asset_type":     assetRequest.AssetType,
		"status":         string(assetRequest.Status),
		"priority":       string(assetRequest.Priority),
		"requested_date": assetRequest.RequestedDate,
	}

	if assetRequest.Brand != nil {
		responseData["brand"] = *assetRequest.Brand
	}
	if assetRequest.Model != nil {
		responseData["model"] = *assetRequest.Model
	}
	if assetRequest.ReplacingAssetCode != nil {
		responseData["replacing_asset_code"] = *assetRequest.ReplacingAssetCode
	}

	response.Success(c, "Asset request created successfully", responseData)
}

// ListAssetRequests handles getting asset requests
// @Summary Get asset requests
// @Description Retrieve list of asset requests for the authenticated employee
// @Tags Self-Service Assets
// @Produce json
// @Param request_type query string false "Request type"
// @Param status query string false "Request status"
// @Param asset_type query string false "Asset type"
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/assets/requests [get]
func (h *AssetSelfServiceHandler) ListAssetRequests(c *gin.Context) {
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
	if requestType := c.Query("request_type"); requestType != "" {
		filters["request_type"] = requestType
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if assetType := c.Query("asset_type"); assetType != "" {
		filters["asset_type"] = assetType
	}

	requests, total, err := h.service.ListAssetRequests(userObj.ID, tenantID, page, pageSize, filters)
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
			"request_type":   string(req.RequestType),
			"asset_type":     req.AssetType,
			"status":         string(req.Status),
			"priority":       string(req.Priority),
			"justification": req.Justification,
			"requested_date": req.RequestedDate,
		}

		if req.Brand != nil {
			reqMap["brand"] = *req.Brand
		}
		if req.Model != nil {
			reqMap["model"] = *req.Model
		}
		if req.ReplacingAssetCode != nil {
			reqMap["replacing_asset_code"] = *req.ReplacingAssetCode
		}
		if req.ApprovedAt != nil {
			reqMap["approved_at"] = *req.ApprovedAt
		}
		if req.RejectedAt != nil {
			reqMap["rejected_at"] = *req.RejectedAt
			if req.RejectionReason != nil {
				reqMap["rejection_reason"] = *req.RejectionReason
			}
		}
		if req.FulfilledAt != nil {
			reqMap["fulfilled_at"] = *req.FulfilledAt
			if req.AssignedAssetCode != nil {
				reqMap["assigned_asset_code"] = *req.AssignedAssetCode
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
	response.SuccessWithMeta(c, "Asset requests retrieved successfully", formattedRequests, meta)
}

// GetAssetRequestDetails handles getting asset request details
// @Summary Get asset request details
// @Description Retrieve detailed information about a specific asset request
// @Tags Self-Service Assets
// @Produce json
// @Param request_id path int true "Request ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/assets/requests/:request_id [get]
func (h *AssetSelfServiceHandler) GetAssetRequestDetails(c *gin.Context) {
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

	request, err := h.service.GetAssetRequestDetails(userObj.ID, uint(requestID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	// Format response
	responseData := map[string]interface{}{
		"id":             request.ID,
		"request_number": request.RequestNumber,
		"request_type":   string(request.RequestType),
		"asset_type":     request.AssetType,
		"status":         string(request.Status),
		"priority":       string(request.Priority),
		"justification":  request.Justification,
		"requested_date": request.RequestedDate,
	}

	if request.Brand != nil {
		responseData["brand"] = *request.Brand
	}
	if request.Model != nil {
		responseData["model"] = *request.Model
	}
	if request.Notes != nil {
		responseData["notes"] = *request.Notes
	}
	if request.ReplacingAssetCode != nil {
		responseData["replacing_asset_code"] = *request.ReplacingAssetCode
	}
	if request.ApprovedAt != nil {
		responseData["approved_at"] = *request.ApprovedAt
	}
	if request.RejectedAt != nil {
		responseData["rejected_at"] = *request.RejectedAt
		if request.RejectionReason != nil {
			responseData["rejection_reason"] = *request.RejectionReason
		}
	}
	if request.FulfilledAt != nil {
		responseData["fulfilled_at"] = *request.FulfilledAt
		if request.AssignedAssetCode != nil {
			responseData["assigned_asset_code"] = *request.AssignedAssetCode
		}
	}

	response.Success(c, "Asset request details retrieved successfully", responseData)
}

// CancelAssetRequest handles cancelling an asset request
// @Summary Cancel asset request
// @Description Cancel a pending asset request
// @Tags Self-Service Assets
// @Accept json
// @Produce json
// @Param request_id path int true "Request ID"
// @Param request body map[string]interface{} true "Cancel request"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/assets/requests/:request_id/cancel [post]
func (h *AssetSelfServiceHandler) CancelAssetRequest(c *gin.Context) {
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

	reason := "No reason provided"
	// Try to get reason from request body (optional)
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err == nil {
		if r, ok := req["reason"].(string); ok && r != "" {
			reason = r
		}
	}

	err = h.service.CancelAssetRequest(userObj.ID, uint(requestID), reason)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset request cancelled successfully", nil)
}

// CreateAssetIssue handles creating an asset issue report
// @Summary Report asset issue
// @Description Report an issue with an assigned asset (malfunction, damage, lost, stolen, etc.)
// @Tags Self-Service Assets
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Asset issue"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/assets/issues [post]
func (h *AssetSelfServiceHandler) CreateAssetIssue(c *gin.Context) {
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

	assetIDFloat, ok := req["asset_id"].(float64)
	if !ok {
		response.BadRequest(c, "asset_id is required", nil)
		return
	}
	assetID := uint(assetIDFloat)

	issueType, ok := req["issue_type"].(string)
	if !ok {
		response.BadRequest(c, "issue_type is required", nil)
		return
	}

	title, ok := req["title"].(string)
	if !ok || title == "" {
		response.BadRequest(c, "title is required", nil)
		return
	}

	description, ok := req["description"].(string)
	if !ok || description == "" {
		response.BadRequest(c, "description is required", nil)
		return
	}

	priority := "medium"
	if p, ok := req["priority"].(string); ok {
		priority = p
	}

	var incidentDate *time.Time
	if idStr, ok := req["incident_date"].(string); ok && idStr != "" {
		parsed, err := time.Parse("2006-01-02", idStr)
		if err == nil {
			incidentDate = &parsed
		}
	}

	var incidentLocation, policeReportNumber *string
	if loc, ok := req["incident_location"].(string); ok && loc != "" {
		incidentLocation = &loc
	}
	if prn, ok := req["police_report_number"].(string); ok && prn != "" {
		policeReportNumber = &prn
	}

	var attachmentURLs []string
	if atts, ok := req["attachment_urls"].([]interface{}); ok {
		for _, att := range atts {
			if url, ok := att.(string); ok {
				attachmentURLs = append(attachmentURLs, url)
			}
		}
	}

	tenantID := middleware.GetTenantID(c)
	updatedBy := &userObj.ID

	issue, err := h.service.CreateAssetIssue(userObj.ID, tenantID, assetID, issueType, title, description, priority, incidentDate, incidentLocation, policeReportNumber, attachmentURLs, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Format response
	responseData := map[string]interface{}{
		"id":           issue.ID,
		"issue_number": issue.IssueNumber,
		"asset_id":     issue.AssetID,
		"asset_code":   issue.AssetCode,
		"issue_type":   string(issue.IssueType),
		"status":       string(issue.Status),
		"priority":     string(issue.Priority),
		"title":        issue.Title,
		"description": issue.Description,
		"reported_date": issue.ReportedDate,
	}

	if issue.IncidentDate != nil {
		responseData["incident_date"] = *issue.IncidentDate
	}
	if issue.IncidentLocation != nil {
		responseData["incident_location"] = *issue.IncidentLocation
	}

	response.Success(c, "Asset issue reported successfully", responseData)
}

// ListAssetIssues handles getting asset issues
// @Summary Get asset issues
// @Description Retrieve list of asset issues reported by the authenticated employee
// @Tags Self-Service Assets
// @Produce json
// @Param issue_type query string false "Issue type"
// @Param status query string false "Issue status"
// @Param asset_id query int false "Asset ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/assets/issues [get]
func (h *AssetSelfServiceHandler) ListAssetIssues(c *gin.Context) {
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
	if issueType := c.Query("issue_type"); issueType != "" {
		filters["issue_type"] = issueType
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if assetIDStr := c.Query("asset_id"); assetIDStr != "" {
		if assetID, err := strconv.ParseUint(assetIDStr, 10, 32); err == nil {
			filters["asset_id"] = uint(assetID)
		}
	}

	issues, total, err := h.service.ListAssetIssues(userObj.ID, tenantID, page, pageSize, filters)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Format response
	formattedIssues := []map[string]interface{}{}
	for _, issue := range issues {
		issueMap := map[string]interface{}{
			"id":            issue.ID,
			"issue_number":  issue.IssueNumber,
			"asset_id":      issue.AssetID,
			"asset_code":    issue.AssetCode,
			"issue_type":    string(issue.IssueType),
			"status":        string(issue.Status),
			"priority":      string(issue.Priority),
			"title":         issue.Title,
			"description":   issue.Description,
			"reported_date": issue.ReportedDate,
		}

		if issue.IncidentDate != nil {
			issueMap["incident_date"] = *issue.IncidentDate
		}
		if issue.IncidentLocation != nil {
			issueMap["incident_location"] = *issue.IncidentLocation
		}
		if issue.PoliceReportNumber != nil {
			issueMap["police_report_number"] = *issue.PoliceReportNumber
		}
		if issue.ResolvedAt != nil {
			issueMap["resolved_at"] = *issue.ResolvedAt
			if issue.ResolutionNotes != nil {
				issueMap["resolution_notes"] = *issue.ResolutionNotes
			}
			if issue.ActionTaken != nil {
				issueMap["action_taken"] = *issue.ActionTaken
			}
		}

		formattedIssues = append(formattedIssues, issueMap)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Asset issues retrieved successfully", formattedIssues, meta)
}

// GetAssetIssueDetails handles getting asset issue details
// @Summary Get asset issue details
// @Description Retrieve detailed information about a specific asset issue
// @Tags Self-Service Assets
// @Produce json
// @Param issue_id path int true "Issue ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/assets/issues/:issue_id [get]
func (h *AssetSelfServiceHandler) GetAssetIssueDetails(c *gin.Context) {
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

	issueIDStr := c.Param("issue_id")
	issueID, err := strconv.ParseUint(issueIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid issue ID", nil)
		return
	}

	issue, err := h.service.GetAssetIssueDetails(userObj.ID, uint(issueID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	// Parse attachment URLs from JSON
	var attachmentURLs []string
	if issue.AttachmentURLsJSON != nil && *issue.AttachmentURLsJSON != "" {
		_ = json.Unmarshal([]byte(*issue.AttachmentURLsJSON), &attachmentURLs)
	}

	// Format response
	responseData := map[string]interface{}{
		"id":            issue.ID,
		"issue_number":  issue.IssueNumber,
		"asset_id":      issue.AssetID,
		"asset_code":    issue.AssetCode,
		"issue_type":    string(issue.IssueType),
		"status":        string(issue.Status),
		"priority":      string(issue.Priority),
		"title":         issue.Title,
		"description":   issue.Description,
		"reported_date": issue.ReportedDate,
		"attachment_urls": attachmentURLs,
	}

	if issue.IncidentDate != nil {
		responseData["incident_date"] = *issue.IncidentDate
	}
	if issue.IncidentLocation != nil {
		responseData["incident_location"] = *issue.IncidentLocation
	}
	if issue.PoliceReportNumber != nil {
		responseData["police_report_number"] = *issue.PoliceReportNumber
	}
	if issue.ResolvedAt != nil {
		responseData["resolved_at"] = *issue.ResolvedAt
		if issue.ResolutionNotes != nil {
			responseData["resolution_notes"] = *issue.ResolutionNotes
		}
		if issue.ActionTaken != nil {
			responseData["action_taken"] = *issue.ActionTaken
		}
	}

	response.Success(c, "Asset issue details retrieved successfully", responseData)
}

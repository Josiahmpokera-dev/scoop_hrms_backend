package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	assetRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/repositories"
	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// AssetHRHandler handles HR/Admin asset management HTTP requests
type AssetHRHandler struct {
	service *services.AssetHRService
}

// NewAssetHRHandler creates a new asset HR handler
func NewAssetHRHandler() *AssetHRHandler {
	return &AssetHRHandler{
		service: services.NewAssetHRService(),
	}
}

// ListAssetRequests handles listing all asset requests (for HR/Admin)
// @Summary List all asset requests
// @Description Retrieve a paginated list of all asset requests for HR/Admin review
// @Tags Assets (HR/Admin)
// @Produce json
// @Param employee_id query string false "Filter by employee ID"
// @Param request_type query string false "Filter by request type"
// @Param status query string false "Filter by status"
// @Param asset_type query string false "Filter by asset type"
// @Param priority query string false "Filter by priority"
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Success 200 {object} response.APIResponse
// @Router/assets/requests [get]
func (h *AssetHRHandler) ListAssetRequests(c *gin.Context) {
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
	if employeeID := c.Query("employee_id"); employeeID != "" {
		filters["employee_id"] = employeeID
	}
	if requestType := c.Query("request_type"); requestType != "" {
		filters["request_type"] = requestType
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if assetType := c.Query("asset_type"); assetType != "" {
		filters["asset_type"] = assetType
	}
	if priority := c.Query("priority"); priority != "" {
		filters["priority"] = priority
	}

	requests, total, err := h.service.ListAssetRequests(tenantID, page, pageSize, filters)
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
			"employee_id":    req.EmployeeID,
			"request_type":   string(req.RequestType),
			"asset_type":     req.AssetType,
			"status":         string(req.Status),
			"priority":       string(req.Priority),
			"justification":  req.Justification,
			"requested_date": req.RequestedDate,
		}

		if req.Brand != nil {
			reqMap["brand"] = *req.Brand
		}
		if req.Model != nil {
			reqMap["model"] = *req.Model
		}
		if req.Notes != nil {
			reqMap["notes"] = *req.Notes
		}
		if req.ReplacingAssetCode != nil {
			reqMap["replacing_asset_code"] = *req.ReplacingAssetCode
		}
		if req.ApprovedAt != nil {
			reqMap["approved_at"] = *req.ApprovedAt
			if req.ApprovedBy != nil {
				reqMap["approved_by"] = *req.ApprovedBy
			}
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
// @Tags Assets (HR/Admin)
// @Produce json
// @Param request_id path int true "Request ID"
// @Success 200 {object} response.APIResponse
// @Router/assets/requests/:request_id [get]
func (h *AssetHRHandler) GetAssetRequestDetails(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	requestIDStr := c.Param("request_id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	request, err := h.service.GetAssetRequestDetails(uint(requestID), tenantID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	// Format response
	responseData := map[string]interface{}{
		"id":             request.ID,
		"request_number": request.RequestNumber,
		"employee_id":    request.EmployeeID,
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
		if request.ApprovedBy != nil {
			responseData["approved_by"] = *request.ApprovedBy
		}
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

// ApproveAssetRequest handles approving an asset request
// @Summary Approve asset request
// @Description Approve a pending asset request
// @Tags Assets (HR/Admin)
// @Produce json
// @Param request_id path int true "Request ID"
// @Success 200 {object} response.APIResponse
// @Router/assets/requests/:request_id/approve [post]
func (h *AssetHRHandler) ApproveAssetRequest(c *gin.Context) {
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

	requestIDStr := c.Param("request_id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	request, err := h.service.ApproveAssetRequest(uint(requestID), tenantID, &userObj.ID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	responseData := map[string]interface{}{
		"id":             request.ID,
		"request_number": request.RequestNumber,
		"status":         string(request.Status),
		"approved_at":    request.ApprovedAt,
	}

	response.Success(c, "Asset request approved successfully", responseData)
}

// RejectAssetRequest handles rejecting an asset request
// @Summary Reject asset request
// @Description Reject a pending asset request
// @Tags Assets (HR/Admin)
// @Accept json
// @Produce json
// @Param request_id path int true "Request ID"
// @Param request body map[string]interface{} true "Rejection reason"
// @Success 200 {object} response.APIResponse
// @Router/assets/requests/:request_id/reject [post]
func (h *AssetHRHandler) RejectAssetRequest(c *gin.Context) {
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

	rejectionReason := "No reason provided"
	if r, ok := req["rejection_reason"].(string); ok && r != "" {
		rejectionReason = r
	}

	request, err := h.service.RejectAssetRequest(uint(requestID), tenantID, rejectionReason, &userObj.ID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	responseData := map[string]interface{}{
		"id":              request.ID,
		"request_number":  request.RequestNumber,
		"status":          string(request.Status),
		"rejected_at":     request.RejectedAt,
		"rejection_reason": request.RejectionReason,
	}

	response.Success(c, "Asset request rejected successfully", responseData)
}

// FulfillAssetRequest handles fulfilling an approved asset request
// @Summary Fulfill asset request
// @Description Assign an asset to fulfill an approved asset request
// @Tags Assets (HR/Admin)
// @Accept json
// @Produce json
// @Param request_id path int true "Request ID"
// @Param request body map[string]interface{} true "Asset assignment"
// @Success 200 {object} response.APIResponse
// @Router/assets/requests/:request_id/fulfill [post]
func (h *AssetHRHandler) FulfillAssetRequest(c *gin.Context) {
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

	assetIDFloat, ok := req["asset_id"].(float64)
	if !ok {
		response.BadRequest(c, "asset_id is required", nil)
		return
	}
	assetID := uint(assetIDFloat)

	request, asset, err := h.service.FulfillAssetRequest(uint(requestID), assetID, tenantID, &userObj.ID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	responseData := map[string]interface{}{
		"request": map[string]interface{}{
			"id":             request.ID,
			"request_number": request.RequestNumber,
			"status":         string(request.Status),
			"fulfilled_at":   request.FulfilledAt,
			"assigned_asset_code": request.AssignedAssetCode,
		},
		"asset": map[string]interface{}{
			"id":          asset.ID,
			"asset_code":  asset.AssetCode,
			"asset_type":  asset.AssetType,
			"brand":       asset.Brand,
			"model":       asset.Model,
			"employee_id": asset.EmployeeID,
			"status":      asset.Status,
			"assigned_date": asset.AssignedDate,
		},
	}

	response.Success(c, "Asset request fulfilled successfully", responseData)
}

// ReassignAsset handles reassigning an asset to another employee
// @Summary Reassign asset
// @Description Reassign an asset from one employee to another
// @Tags Assets (HR/Admin)
// @Accept json
// @Produce json
// @Param asset_id path int true "Asset ID"
// @Param request body map[string]interface{} true "Reassignment details"
// @Success 200 {object} response.APIResponse
// @Router/assets/:asset_id/reassign [post]
func (h *AssetHRHandler) ReassignAsset(c *gin.Context) {
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

	assetIDStr := c.Param("asset_id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid asset ID", nil)
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	newEmployeeID, ok := req["employee_id"].(string)
	if !ok || newEmployeeID == "" {
		response.BadRequest(c, "employee_id is required", nil)
		return
	}

	var notes *string
	if n, ok := req["notes"].(string); ok && n != "" {
		notes = &n
	}

	asset, err := h.service.ReassignAsset(uint(assetID), newEmployeeID, tenantID, &userObj.ID, notes)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	responseData := map[string]interface{}{
		"id":          asset.ID,
		"asset_code":  asset.AssetCode,
		"asset_type":  asset.AssetType,
		"brand":       asset.Brand,
		"model":       asset.Model,
		"employee_id": asset.EmployeeID,
		"assigned_to": asset.AssignedTo,
		"status":      asset.Status,
		"assigned_date": asset.AssignedDate,
	}

	response.Success(c, "Asset reassigned successfully", responseData)
}

// GetAssignmentHistory handles listing assignment/reassignment/return history.
// @Summary Get assets assignment history
// @Description Retrieve assignment timeline with filters for asset, employee, action, and date range
// @Tags Assets (HR/Admin)
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 20, max: 100)"
// @Param asset_id query int false "Filter by asset ID"
// @Param employee_id query string false "Filter by employee ID (from/to)"
// @Param action query string false "assign | reassign | return"
// @Param start_date query string false "YYYY-MM-DD"
// @Param end_date query string false "YYYY-MM-DD"
// @Success 200 {object} response.APIResponse
// @Router /assets/assignment-history [get]
func (h *AssetHRHandler) GetAssignmentHistory(c *gin.Context) {
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

	var assetID *uint
	if rawAssetID := strings.TrimSpace(c.Query("asset_id")); rawAssetID != "" {
		parsed, err := strconv.ParseUint(rawAssetID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid asset_id", nil)
			return
		}
		id := uint(parsed)
		assetID = &id
	}

	var employeeID *string
	if rawEmployeeID := strings.TrimSpace(c.Query("employee_id")); rawEmployeeID != "" {
		employeeID = &rawEmployeeID
	}

	var action *string
	if rawAction := strings.TrimSpace(c.Query("action")); rawAction != "" {
		lc := strings.ToLower(rawAction)
		action = &lc
	}

	var startDate *time.Time
	if rawStart := strings.TrimSpace(c.Query("start_date")); rawStart != "" {
		t, err := time.Parse("2006-01-02", rawStart)
		if err != nil {
			response.BadRequest(c, "Invalid start_date format. Use YYYY-MM-DD", nil)
			return
		}
		startDate = &t
	}

	var endDate *time.Time
	if rawEnd := strings.TrimSpace(c.Query("end_date")); rawEnd != "" {
		t, err := time.Parse("2006-01-02", rawEnd)
		if err != nil {
			response.BadRequest(c, "Invalid end_date format. Use YYYY-MM-DD", nil)
			return
		}
		dayEnd := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
		endDate = &dayEnd
	}

	rows, total, err := h.service.ListAssignmentHistory(tenantID, page, pageSize, assetID, employeeID, action, startDate, endDate)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	assetRepo := assetRepos.NewAssetRepository()
	deptRepo := departmentRepos.NewDepartmentRepository()
	modelCache := map[uint]string{}
	deptCache := map[string]string{}

	formatted := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		assetModel := ""
		if cached, ok := modelCache[row.AssetID]; ok {
			assetModel = cached
		} else if asset, assetErr := assetRepo.FindByID(row.AssetID); assetErr == nil && asset != nil {
			assetModel = asset.Model
			modelCache[row.AssetID] = assetModel
		}

		departmentName := ""
		if row.Department != nil && strings.TrimSpace(*row.Department) != "" {
			rawDept := strings.TrimSpace(*row.Department)
			if cached, ok := deptCache[rawDept]; ok {
				departmentName = cached
			} else {
				departmentName = rawDept
				if deptID, parseErr := strconv.ParseUint(rawDept, 10, 32); parseErr == nil {
					if dept, deptErr := deptRepo.FindByID(uint(deptID)); deptErr == nil && dept != nil && strings.TrimSpace(dept.Name) != "" {
						departmentName = dept.Name
					}
				}
				deptCache[rawDept] = departmentName
			}
		}

		formatted = append(formatted, map[string]interface{}{
			"id":               row.ID,
			"asset_code":       row.AssetCode,
			"asset_model":      assetModel,
			"action":           string(row.Action),
			"from_employee_id": row.FromEmployeeID,
			"to_employee_id":   row.ToEmployeeID,
			"from_employee":    row.FromEmployee,
			"to_employee":      row.ToEmployee,
			"department":       departmentName,
			"notes":            row.Notes,
			"occurred_at":      row.OccurredAt,
			"actor_user_id":    row.ActorUserID,
			"created_at":       row.CreatedAt,
		})
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Asset assignment history retrieved successfully", formatted, meta)
}

// GetAssetAssignmentHistory handles listing assignment history for one asset.
// @Summary Get assignment history by asset
// @Description Retrieve assignment timeline for a specific asset id
// @Tags Assets (HR/Admin)
// @Produce json
// @Param asset_id path int true "Asset ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 20, max: 100)"
// @Param employee_id query string false "Filter by employee ID (from/to)"
// @Param action query string false "assign | reassign | return"
// @Param start_date query string false "YYYY-MM-DD"
// @Param end_date query string false "YYYY-MM-DD"
// @Success 200 {object} response.APIResponse
// @Router /assets/{asset_id}/assignment-history [get]
func (h *AssetHRHandler) GetAssetAssignmentHistory(c *gin.Context) {
	assetIDStr := strings.TrimSpace(c.Param("asset_id"))
	if assetIDStr == "" {
		response.BadRequest(c, "asset_id is required", nil)
		return
	}
	if _, err := strconv.ParseUint(assetIDStr, 10, 32); err != nil {
		response.BadRequest(c, "Invalid asset_id", nil)
		return
	}

	// Reuse shared history endpoint behavior by injecting asset_id filter.
	q := c.Request.URL.Query()
	q.Set("asset_id", assetIDStr)
	c.Request.URL.RawQuery = q.Encode()
	h.GetAssignmentHistory(c)
}

package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// OvertimeHandler handles overtime API endpoints
type OvertimeHandler struct {
	overtimeService *services.OvertimeService
}

// NewOvertimeHandler creates a new overtime handler
func NewOvertimeHandler() *OvertimeHandler {
	return &OvertimeHandler{
		overtimeService: services.NewOvertimeService(),
	}
}

// --- Policy Endpoints ---

// CreatePolicy godoc
// @Summary Create an overtime policy (Admin/HR)
func (h *OvertimeHandler) CreatePolicy(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	var req models.CreateOvertimePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	policy, err := h.overtimeService.CreatePolicy(user, req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Overtime policy created successfully", policy)
}

// UpdatePolicy godoc
// @Summary Update an overtime policy (Admin/HR)
func (h *OvertimeHandler) UpdatePolicy(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	policyID, err := strconv.ParseUint(c.Param("policy_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid policy ID", nil)
		return
	}

	var req models.CreateOvertimePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	policy, err := h.overtimeService.UpdatePolicy(user, uint(policyID), req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Overtime policy updated successfully", policy)
}

// DeletePolicy godoc
// @Summary Delete an overtime policy (Admin/HR)
func (h *OvertimeHandler) DeletePolicy(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	policyID, err := strconv.ParseUint(c.Param("policy_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid policy ID", nil)
		return
	}

	if err := h.overtimeService.DeletePolicy(user, uint(policyID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Overtime policy deleted successfully", nil)
}

// GetPolicy godoc
// @Summary Get an overtime policy by ID
func (h *OvertimeHandler) GetPolicy(c *gin.Context) {
	policyID, err := strconv.ParseUint(c.Param("policy_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid policy ID", nil)
		return
	}

	policy, err := h.overtimeService.GetPolicy(uint(policyID))
	if err != nil {
		response.NotFound(c, "Overtime policy not found")
		return
	}

	response.Success(c, "Overtime policy retrieved successfully", policy)
}

// ListPolicies godoc
// @Summary List all overtime policies
func (h *OvertimeHandler) ListPolicies(c *gin.Context) {
	page, pageSize := getPagination(c)

	policies, total, err := h.overtimeService.ListPolicies(page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response.SuccessWithMeta(c, "Overtime policies retrieved successfully", policies, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// --- OT Request Endpoints ---

// CreateOTRequest godoc
// @Summary Submit an overtime request
func (h *OvertimeHandler) CreateOTRequest(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	var req models.CreateOvertimeRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	otRequest, err := h.overtimeService.CreateOTRequest(user, req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Overtime request submitted successfully", otRequest)
}

// GetMyOTRequests godoc
// @Summary Get my overtime requests
func (h *OvertimeHandler) GetMyOTRequests(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	page, pageSize := getPagination(c)
	status := c.Query("status")

	requests, total, err := h.overtimeService.GetMyOTRequests(user, status, page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response.SuccessWithMeta(c, "Overtime requests retrieved successfully", requests, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetOTRequestByID godoc
// @Summary Get an overtime request by ID
func (h *OvertimeHandler) GetOTRequestByID(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	otRequest, err := h.overtimeService.GetOTRequestByID(user, uint(requestID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Overtime request retrieved successfully", otRequest)
}

// ApproveOTRequest godoc
// @Summary Approve or reject an OT request (Admin/HR)
func (h *OvertimeHandler) ApproveOTRequest(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	var req models.ApproveOvertimeRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	otRequest, err := h.overtimeService.ApproveOTRequest(user, uint(requestID), req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Overtime request processed successfully", otRequest)
}

// CancelOTRequest godoc
// @Summary Cancel an overtime request
func (h *OvertimeHandler) CancelOTRequest(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	otRequest, err := h.overtimeService.CancelOTRequest(user, uint(requestID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Overtime request cancelled successfully", otRequest)
}

// GetPendingApprovals godoc
// @Summary Get pending OT approvals (Admin/HR)
func (h *OvertimeHandler) GetPendingApprovals(c *gin.Context) {
	page, pageSize := getPagination(c)

	requests, total, err := h.overtimeService.GetPendingApprovals(page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response.SuccessWithMeta(c, "Pending overtime requests retrieved successfully", requests, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// ListAllRequests godoc
// @Summary List all overtime requests (Admin/HR)
func (h *OvertimeHandler) ListAllRequests(c *gin.Context) {
	page, pageSize := getPagination(c)
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	requests, total, err := h.overtimeService.ListAllRequests(status, startDate, endDate, page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response.SuccessWithMeta(c, "Overtime requests retrieved successfully", requests, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetOvertimeStats godoc
// @Summary Get overtime statistics
func (h *OvertimeHandler) GetOvertimeStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.overtimeService.GetOvertimeStats(startDate, endDate)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Overtime statistics retrieved successfully", stats)
}

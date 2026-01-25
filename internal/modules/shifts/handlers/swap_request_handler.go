package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/services"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type SwapRequestHandler struct {
	service     *services.SwapRequestService
	employeeRepo *employeeRepos.EmployeeRepository
}

func NewSwapRequestHandler() *SwapRequestHandler {
	return &SwapRequestHandler{
		service:     services.NewSwapRequestService(),
		employeeRepo: employeeRepos.NewEmployeeRepository(),
	}
}

// toSwapRequestResponse converts SwapRequest to response format
func (h *SwapRequestHandler) toSwapRequestResponse(swapRequest *models.SwapRequest) map[string]interface{} {
	requestedByEmp, _ := h.employeeRepo.FindByEmployeeID(swapRequest.RequestedBy)
	requestedWithEmp, _ := h.employeeRepo.FindByEmployeeID(swapRequest.RequestedWith)

	response := map[string]interface{}{
		"id":                    swapRequest.ID,
		"requested_by":          swapRequest.RequestedBy,
		"requested_by_name":     getEmployeeFullName(requestedByEmp),
		"requested_by_photo":    getEmployeePhoto(requestedByEmp),
		"requested_by_department": getEmployeeDepartment(requestedByEmp),
		"requested_with":        swapRequest.RequestedWith,
		"requested_with_name":   getEmployeeFullName(requestedWithEmp),
		"requested_with_photo":  getEmployeePhoto(requestedWithEmp),
		"requested_with_department": getEmployeeDepartment(requestedWithEmp),
		"assignment_id":         swapRequest.AssignmentID,
		"swap_assignment_id":    swapRequest.SwapAssignmentID,
		"reason":                swapRequest.Reason,
		"status":                swapRequest.Status,
		"requested_at":          swapRequest.RequestedAt.Format("2006-01-02T15:04:05Z07:00"),
		"reviewed_at":           nil,
		"reviewed_by":           swapRequest.ReviewedBy,
		"reviewer_name":         nil,
	}

	if swapRequest.ReviewedAt != nil {
		response["reviewed_at"] = swapRequest.ReviewedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if swapRequest.RejectionReason != nil {
		response["rejection_reason"] = *swapRequest.RejectionReason
	}

	// Add assignment details if loaded
	if swapRequest.Assignment != nil {
		response["assignment_date"] = swapRequest.Assignment.Date.Format("2006-01-02")
		if swapRequest.Assignment.Shift != nil {
			response["assignment_shift"] = swapRequest.Assignment.Shift.ShiftName
			response["assignment_shift_timing"] = swapRequest.Assignment.Shift.StartTime + " - " + swapRequest.Assignment.Shift.EndTime
		}
	}

	// Add swap assignment details if loaded
	if swapRequest.SwapAssignment != nil {
		response["swap_assignment_date"] = swapRequest.SwapAssignment.Date.Format("2006-01-02")
		if swapRequest.SwapAssignment.Shift != nil {
			response["swap_assignment_shift"] = swapRequest.SwapAssignment.Shift.ShiftName
			response["swap_assignment_shift_timing"] = swapRequest.SwapAssignment.Shift.StartTime + " - " + swapRequest.SwapAssignment.Shift.EndTime
		}
	}

	return response
}

// Helper function to get employee name from ID
func getEmployeeNameFromID(employeeID string, employeeRepo *employeeRepos.EmployeeRepository) string {
	employee, err := employeeRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return ""
	}
	return employee.FullName()
}

// ListSwapRequests handles listing swap requests
func (h *SwapRequestHandler) ListSwapRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)

	// Build filters
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if requestedBy := c.Query("requested_by"); requestedBy != "" {
		filters["requested_by"] = requestedBy
	}
	if requestedWith := c.Query("requested_with"); requestedWith != "" {
		filters["requested_with"] = requestedWith
	}

	swapRequests, total, err := h.service.ListSwapRequests(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve swap requests", err.Error())
		return
	}

	// Convert to response format
	swapRequestResponses := make([]map[string]interface{}, len(swapRequests))
	for i, swapRequest := range swapRequests {
		swapRequestResponses[i] = h.toSwapRequestResponse(&swapRequest)
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}

	response.Success(c, "Swap requests retrieved successfully", map[string]interface{}{
		"data": swapRequestResponses,
		"meta": map[string]interface{}{
			"page":       page,
			"per_page":   pageSize,
			"total":      total,
			"total_pages": totalPages,
		},
	})
}

// GetSwapRequest handles getting a swap request by ID
func (h *SwapRequestHandler) GetSwapRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	swapRequest, err := h.service.GetSwapRequestByID(uint(id))
	if err != nil {
		response.NotFound(c, "Swap request not found")
		return
	}

	response.Success(c, "Swap request retrieved successfully", h.toSwapRequestResponse(swapRequest))
}

// CreateSwapRequest handles creating a swap request
func (h *SwapRequestHandler) CreateSwapRequest(c *gin.Context) {
	var req models.CreateSwapRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	
	// Get employee ID from user - this would need to be implemented based on your user-employee relationship
	var requestedBy string
	if userObj, ok := user.(*userModels.User); ok {
		// Find employee by user ID
		employee, err := h.employeeRepo.FindByUserID(userObj.ID)
		if err == nil && employee != nil {
			requestedBy = employee.EmployeeID
		} else {
			response.BadRequest(c, "Employee not found for user", nil)
			return
		}
	} else {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	swapRequest, err := h.service.CreateSwapRequest(&req, requestedBy, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Swap request created successfully", h.toSwapRequestResponse(swapRequest))
}

// ApproveSwapRequest handles approving a swap request
func (h *SwapRequestHandler) ApproveSwapRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	var req models.ApproveSwapRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// JSON is optional
		_ = err
	}

	user, _ := c.Get("user")
	var reviewedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		reviewedBy = &userObj.ID
	}

	swapRequest, err := h.service.ApproveSwapRequest(uint(id), &req, reviewedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Get updated assignments for response - reload from service
	rosterService := services.NewRosterService()
	assignment, _ := rosterService.GetRosterAssignmentByID(swapRequest.AssignmentID)
	swapAssignment, _ := rosterService.GetRosterAssignmentByID(swapRequest.SwapAssignmentID)

	response.Success(c, "Swap request approved successfully", map[string]interface{}{
		"id":              swapRequest.ID,
		"status":           swapRequest.Status,
		"reviewed_at":      swapRequest.ReviewedAt.Format("2006-01-02T15:04:05Z07:00"),
		"reviewed_by":      swapRequest.ReviewedBy,
		"updated_assignments": []map[string]interface{}{
			{
				"assignment_id":      assignment.ID,
				"new_employee_id":    assignment.EmployeeID,
				"new_employee_name":  getEmployeeNameFromID(assignment.EmployeeID, h.employeeRepo),
			},
			{
				"assignment_id":      swapAssignment.ID,
				"new_employee_id":    swapAssignment.EmployeeID,
				"new_employee_name":  getEmployeeNameFromID(swapAssignment.EmployeeID, h.employeeRepo),
			},
		},
	})
}

// RejectSwapRequest handles rejecting a swap request
func (h *SwapRequestHandler) RejectSwapRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	var req models.RejectSwapRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var reviewedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		reviewedBy = &userObj.ID
	}

	swapRequest, err := h.service.RejectSwapRequest(uint(id), &req, reviewedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Swap request rejected successfully", map[string]interface{}{
		"id":              swapRequest.ID,
		"status":           swapRequest.Status,
		"rejection_reason": swapRequest.RejectionReason,
		"reviewed_at":      swapRequest.ReviewedAt.Format("2006-01-02T15:04:05Z07:00"),
		"reviewed_by":      swapRequest.ReviewedBy,
	})
}

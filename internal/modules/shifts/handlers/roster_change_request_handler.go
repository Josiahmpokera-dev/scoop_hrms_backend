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

type RosterChangeRequestHandler struct {
	service     *services.RosterChangeRequestService
	employeeRepo *employeeRepos.EmployeeRepository
}

func NewRosterChangeRequestHandler() *RosterChangeRequestHandler {
	return &RosterChangeRequestHandler{
		service:     services.NewRosterChangeRequestService(),
		employeeRepo: employeeRepos.NewEmployeeRepository(),
	}
}

// toRosterChangeRequestResponse converts RosterChangeRequest to response format
func (h *RosterChangeRequestHandler) toRosterChangeRequestResponse(changeRequest *models.RosterChangeRequest) map[string]interface{} {
	requestedByEmp, _ := h.employeeRepo.FindByEmployeeID(changeRequest.RequestedBy)

	response := map[string]interface{}{
		"id":                    changeRequest.ID,
		"requested_by":          changeRequest.RequestedBy,
		"requested_by_name":     getEmployeeFullName(requestedByEmp),
		"requested_by_photo":    getEmployeePhoto(requestedByEmp),
		"requested_by_department": getEmployeeDepartment(requestedByEmp),
		"assignment_id":         changeRequest.AssignmentID,
		"requested_shift_id":     changeRequest.RequestedShiftID,
		"requested_date":         nil,
		"requested_location_id":  changeRequest.RequestedLocationID,
		"reason":                changeRequest.Reason,
		"status":                changeRequest.Status,
		"requested_at":          changeRequest.RequestedAt.Format("2006-01-02T15:04:05Z07:00"),
		"reviewed_at":           nil,
		"reviewed_by":           changeRequest.ReviewedBy,
		"reviewer_name":         nil,
	}

	if changeRequest.RequestedDate != nil {
		response["requested_date"] = changeRequest.RequestedDate.Format("2006-01-02")
	}

	if changeRequest.ReviewedAt != nil {
		response["reviewed_at"] = changeRequest.ReviewedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if changeRequest.RejectionReason != nil {
		response["rejection_reason"] = *changeRequest.RejectionReason
	}

	// Add assignment details if loaded
	if changeRequest.Assignment != nil {
		response["assignment_date"] = changeRequest.Assignment.Date.Format("2006-01-02")
		if changeRequest.Assignment.Shift != nil {
			response["current_shift_id"] = changeRequest.Assignment.ShiftID
			response["current_shift_name"] = changeRequest.Assignment.Shift.ShiftName
			response["current_shift_timing"] = changeRequest.Assignment.Shift.StartTime + " - " + changeRequest.Assignment.Shift.EndTime
		}
		if changeRequest.Assignment.Location != nil {
			response["current_location_id"] = changeRequest.Assignment.LocationID
			response["current_location_name"] = changeRequest.Assignment.Location.Name
		}
	}

	// Add requested shift details if loaded
	if changeRequest.RequestedShift != nil {
		response["requested_shift_name"] = changeRequest.RequestedShift.ShiftName
		response["requested_shift_timing"] = changeRequest.RequestedShift.StartTime + " - " + changeRequest.RequestedShift.EndTime
	}

	return response
}

// ListRosterChangeRequests handles listing roster change requests
func (h *RosterChangeRequestHandler) ListRosterChangeRequests(c *gin.Context) {
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
	
	// For employee self-service, automatically filter by their employee ID
	// For HR/Admin, they can see all or filter by requested_by
	user, _ := c.Get("user")
	if userObj, ok := user.(*userModels.User); ok {
		employee, err := h.employeeRepo.FindByUserID(userObj.ID)
		if err == nil && employee != nil {
			// Check if user is HR/Admin - if not, filter by their employee ID
			if userObj.Role != userModels.RoleAdmin && userObj.Role != userModels.RoleHR {
				filters["requested_by"] = employee.EmployeeID
			} else {
				// HR/Admin can filter by requested_by if provided
				if requestedBy := c.Query("requested_by"); requestedBy != "" {
					filters["requested_by"] = requestedBy
				}
			}
		}
	} else {
		// If requested_by is explicitly provided, use it
		if requestedBy := c.Query("requested_by"); requestedBy != "" {
			filters["requested_by"] = requestedBy
		}
	}
	
	if assignmentIDStr := c.Query("assignment_id"); assignmentIDStr != "" {
		if assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32); err == nil {
			filters["assignment_id"] = uint(assignmentID)
		}
	}

	changeRequests, total, err := h.service.ListRosterChangeRequests(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve roster change requests", err.Error())
		return
	}

	// Convert to response format
	changeRequestResponses := make([]map[string]interface{}, len(changeRequests))
	for i, changeRequest := range changeRequests {
		changeRequestResponses[i] = h.toRosterChangeRequestResponse(&changeRequest)
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}

	response.Success(c, "Roster change requests retrieved successfully", map[string]interface{}{
		"data": changeRequestResponses,
		"meta": map[string]interface{}{
			"page":       page,
			"per_page":   pageSize,
			"total":      total,
			"total_pages": totalPages,
		},
	})
}

// GetRosterChangeRequest handles getting a roster change request by ID
func (h *RosterChangeRequestHandler) GetRosterChangeRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	changeRequest, err := h.service.GetRosterChangeRequestByID(uint(id))
	if err != nil {
		response.NotFound(c, "Roster change request not found")
		return
	}

	response.Success(c, "Roster change request retrieved successfully", h.toRosterChangeRequestResponse(changeRequest))
}

// CreateRosterChangeRequest handles creating a roster change request (Employee endpoint)
func (h *RosterChangeRequestHandler) CreateRosterChangeRequest(c *gin.Context) {
	var req models.CreateRosterChangeRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	
	// Get employee ID from user
	var requestedBy string
	if userObj, ok := user.(*userModels.User); ok {
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

	changeRequest, err := h.service.CreateRosterChangeRequest(&req, requestedBy, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Roster change request created successfully", h.toRosterChangeRequestResponse(changeRequest))
}

// ApproveRosterChangeRequest handles approving a roster change request (HR/Admin endpoint)
func (h *RosterChangeRequestHandler) ApproveRosterChangeRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	var req models.ApproveRosterChangeRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// JSON is optional
		_ = err
	}

	user, _ := c.Get("user")
	var reviewedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		reviewedBy = &userObj.ID
	}

	changeRequest, err := h.service.ApproveRosterChangeRequest(uint(id), &req, reviewedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Get updated assignment for response
	rosterService := services.NewRosterService()
	assignment, _ := rosterService.GetRosterAssignmentByID(changeRequest.AssignmentID)

	assignmentData := map[string]interface{}{}
	if assignment != nil {
		assignmentData = map[string]interface{}{
			"assignment_id":      assignment.ID,
			"employee_id":        assignment.EmployeeID,
			"date":               assignment.Date.Format("2006-01-02"),
			"shift_id":           assignment.ShiftID,
			"location_id":        assignment.LocationID,
			"status":             assignment.Status,
		}
		if assignment.Shift != nil {
			assignmentData["shift_name"] = assignment.Shift.ShiftName
		}
		if assignment.Location != nil {
			assignmentData["location_name"] = assignment.Location.Name
		}
	}

	response.Success(c, "Roster change request approved successfully", map[string]interface{}{
		"id":              changeRequest.ID,
		"status":           changeRequest.Status,
		"reviewed_at":      changeRequest.ReviewedAt.Format("2006-01-02T15:04:05Z07:00"),
		"reviewed_by":      changeRequest.ReviewedBy,
		"updated_assignment": assignmentData,
	})
}

// RejectRosterChangeRequest handles rejecting a roster change request (HR/Admin endpoint)
func (h *RosterChangeRequestHandler) RejectRosterChangeRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	var req models.RejectRosterChangeRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var reviewedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		reviewedBy = &userObj.ID
	}

	changeRequest, err := h.service.RejectRosterChangeRequest(uint(id), &req, reviewedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Roster change request rejected successfully", map[string]interface{}{
		"id":              changeRequest.ID,
		"status":           changeRequest.Status,
		"rejection_reason": changeRequest.RejectionReason,
		"reviewed_at":      changeRequest.ReviewedAt.Format("2006-01-02T15:04:05Z07:00"),
		"reviewed_by":      changeRequest.ReviewedBy,
	})
}

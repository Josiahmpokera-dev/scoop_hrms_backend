package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/services"
	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type RosterHandler struct {
	service       *services.RosterService
	employeeRepo  *employeeRepos.EmployeeRepository
	departmentRepo *departmentRepos.DepartmentRepository
}

func NewRosterHandler() *RosterHandler {
	return &RosterHandler{
		service:       services.NewRosterService(),
		employeeRepo:  employeeRepos.NewEmployeeRepository(),
		departmentRepo: departmentRepos.NewDepartmentRepository(),
	}
}

// toRosterResponse converts RosterAssignment to response format
func (h *RosterHandler) toRosterResponse(assignment *models.RosterAssignment) map[string]interface{} {
	employee, _ := h.employeeRepo.FindByEmployeeID(assignment.EmployeeID)

	response := map[string]interface{}{
		"id":            assignment.ID,
		"employee_id":   assignment.EmployeeID,
		"employee_name": getEmployeeFullName(employee),
		"date":          assignment.Date.Format("2006-01-02"),
		"shift_id":      assignment.ShiftID,
		"status":        assignment.Status,
		"created_at":    assignment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updated_at":    assignment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"created_by":    assignment.CreatedBy,
		"updated_by":    assignment.UpdatedBy,
	}

	if assignment.Shift != nil {
		response["shift_name"] = assignment.Shift.ShiftName
		response["shift_code"] = assignment.Shift.ShiftCode
		response["shift_timing"] = assignment.Shift.StartTime + " - " + assignment.Shift.EndTime
	}

	if assignment.Location != nil {
		response["location_id"] = assignment.Location.ID
		response["location_name"] = assignment.Location.Name
	}

	if employee != nil {
		response["employee_photo"] = getEmployeePhoto(employee)
		// Get department name
		if employee.DepartmentID != nil {
			dept, err := h.departmentRepo.FindByID(*employee.DepartmentID)
			if err == nil && dept != nil {
				response["department"] = dept.Name
			}
		}
	}

	response["swap_request"] = nil // Will be populated if there's a swap request

	return response
}

// ListRosterAssignments handles listing roster assignments
func (h *RosterHandler) ListRosterAssignments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)

	// Build filters
	filters := make(map[string]interface{})
	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filters["end_date"] = endDate
	}
	if employeeID := c.Query("employee_id"); employeeID != "" {
		filters["employee_id"] = employeeID
	}
	if shiftIDStr := c.Query("shift_id"); shiftIDStr != "" {
		if shiftID, err := strconv.ParseUint(shiftIDStr, 10, 32); err == nil {
			filters["shift_id"] = uint(shiftID)
		}
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if locationIDStr := c.Query("location_id"); locationIDStr != "" {
		if locationID, err := strconv.ParseUint(locationIDStr, 10, 32); err == nil {
			filters["location_id"] = uint(locationID)
		}
	}

	assignments, total, err := h.service.ListRosterAssignments(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve roster assignments", err.Error())
		return
	}

	// Convert to response format
	assignmentResponses := make([]map[string]interface{}, len(assignments))
	for i, assignment := range assignments {
		assignmentResponses[i] = h.toRosterResponse(&assignment)
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}

	response.Success(c, "Roster assignments retrieved successfully", map[string]interface{}{
		"data": assignmentResponses,
		"meta": map[string]interface{}{
			"page":       page,
			"per_page":   pageSize,
			"total":      total,
			"total_pages": totalPages,
		},
	})
}

// GetRosterAssignment handles getting a roster assignment by ID
func (h *RosterHandler) GetRosterAssignment(c *gin.Context) {
	idStr := c.Param("assignment_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid assignment ID", nil)
		return
	}

	assignment, err := h.service.GetRosterAssignmentByID(uint(id))
	if err != nil {
		response.NotFound(c, "Roster assignment not found")
		return
	}

	response.Success(c, "Roster assignment retrieved successfully", h.toRosterResponse(assignment))
}

// CreateRosterAssignment handles creating a roster assignment
func (h *RosterHandler) CreateRosterAssignment(c *gin.Context) {
	var req models.CreateRosterAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var createdBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		createdBy = &userObj.ID
	}

	assignment, err := h.service.CreateRosterAssignment(&req, tenantID, createdBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Roster assignment created successfully", h.toRosterResponse(assignment))
}

// BulkCreateRosterAssignments handles bulk creating roster assignments
func (h *RosterHandler) BulkCreateRosterAssignments(c *gin.Context) {
	var req models.BulkCreateRosterAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var createdBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		createdBy = &userObj.ID
	}

	assignments, created, failed, err := h.service.BulkCreateRosterAssignments(&req, tenantID, createdBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert to response format
	assignmentResponses := make([]map[string]interface{}, len(assignments))
	for i, assignment := range assignments {
		assignmentResponses[i] = map[string]interface{}{
			"id":          assignment.ID,
			"employee_id": assignment.EmployeeID,
			"date":        assignment.Date.Format("2006-01-02"),
			"shift_id":    assignment.ShiftID,
			"status":      assignment.Status,
		}
	}

	response.Created(c, "Roster assignments created successfully", map[string]interface{}{
		"created":      created,
		"failed":       failed,
		"assignments":  assignmentResponses,
	})
}

// UpdateRosterAssignment handles updating a roster assignment
func (h *RosterHandler) UpdateRosterAssignment(c *gin.Context) {
	idStr := c.Param("assignment_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid assignment ID", nil)
		return
	}

	var req models.UpdateRosterAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	assignment, err := h.service.UpdateRosterAssignment(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Roster assignment updated successfully", h.toRosterResponse(assignment))
}

// DeleteRosterAssignment handles deleting a roster assignment
func (h *RosterHandler) DeleteRosterAssignment(c *gin.Context) {
	idStr := c.Param("assignment_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid assignment ID", nil)
		return
	}

	if err := h.service.DeleteRosterAssignment(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Roster assignment deleted successfully", map[string]interface{}{
		"assignment_id": id,
	})
}

// GetWeeklyRosterView handles getting weekly roster view
func (h *RosterHandler) GetWeeklyRosterView(c *gin.Context) {
	var req models.WeeklyRosterRequest
	req.Week = c.Query("week")
	if req.Week == "" {
		response.BadRequest(c, "week parameter is required (format: YYYY-WW)", nil)
		return
	}

	// Parse optional filters
	if employeeIDsStr := c.Query("employee_ids"); employeeIDsStr != "" {
		// Parse comma-separated employee IDs
		employeeIDs := []string{}
		parts := strings.Split(employeeIDsStr, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				employeeIDs = append(employeeIDs, trimmed)
			}
		}
		req.EmployeeIDs = employeeIDs
	}
	if departmentIDStr := c.Query("department_id"); departmentIDStr != "" {
		if departmentID, err := strconv.ParseUint(departmentIDStr, 10, 32); err == nil {
			req.DepartmentID = new(uint)
			*req.DepartmentID = uint(departmentID)
		}
	}
	if locationIDStr := c.Query("location_id"); locationIDStr != "" {
		if locationID, err := strconv.ParseUint(locationIDStr, 10, 32); err == nil {
			req.LocationID = new(uint)
			*req.LocationID = uint(locationID)
		}
	}

	tenantID := middleware.GetTenantID(c)
	result, err := h.service.GetWeeklyRosterView(&req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Weekly roster retrieved successfully", result)
}

// AutoSchedule handles auto-scheduling rosters
func (h *RosterHandler) AutoSchedule(c *gin.Context) {
	var req models.AutoScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var createdBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		createdBy = &userObj.ID
	}

	assignments, count, err := h.service.AutoSchedule(&req, tenantID, createdBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert to response format
	assignmentResponses := make([]map[string]interface{}, len(assignments))
	for i, assignment := range assignments {
		assignmentResponses[i] = map[string]interface{}{
			"id":          assignment.ID,
			"employee_id": assignment.EmployeeID,
			"date":        assignment.Date.Format("2006-01-02"),
			"shift_id":    assignment.ShiftID,
			"status":      assignment.Status,
		}
	}

	response.Created(c, "Roster auto-scheduled successfully", map[string]interface{}{
		"created_assignments": count,
		"assignments":         assignmentResponses,
	})
}

// PublishRoster handles publishing rosters
func (h *RosterHandler) PublishRoster(c *gin.Context) {
	var req models.PublishRosterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	count, err := h.service.PublishRosters(&req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	user, _ := c.Get("user")
	var publishedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		publishedBy = &userObj.ID
	}

	response.Success(c, "Roster published successfully", map[string]interface{}{
		"published_count": count,
		"published_at":    time.Now().Format("2006-01-02T15:04:05Z07:00"),
		"published_by":    publishedBy,
	})
}

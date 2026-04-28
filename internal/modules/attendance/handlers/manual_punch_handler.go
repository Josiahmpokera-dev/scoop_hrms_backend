package handlers

import (
	"strconv"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// ManualPunchHandler handles manual punch API endpoints
type ManualPunchHandler struct {
	manualPunchService *services.ManualPunchService
}

// NewManualPunchHandler creates a new manual punch handler
func NewManualPunchHandler() *ManualPunchHandler {
	return &ManualPunchHandler{
		manualPunchService: services.NewManualPunchService(),
	}
}

// CreateManualPunch godoc
// @Summary Create a manual punch request
// @Description Create a manual attendance punch request for when biometric/fingerprint fails
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ManualPunchRequest true "Manual punch request details"
// @Success 201 {object} response.APIResponse "Manual punch created successfully"
// @Failure 400 {object} response.APIResponse "Invalid request"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /attendance/manual-punches [post]
func (h *ManualPunchHandler) CreateManualPunch(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	var request models.ManualPunchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Invalid request format", err)
		return
	}



	manualPunch, err := h.manualPunchService.CreateManualPunch(user.ID, request)
	if err != nil {
		response.InternalServerError(c, "Failed to create manual punch", err)
		return
	}

	response.Created(c, "Manual punch request submitted successfully", manualPunch)
}

// GetMyManualPunches godoc
// @Summary Get my manual punch requests
// @Description Get all manual punch requests for the authenticated employee
// @Tags Attendance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param status query string false "Filter by status (pending, approved, rejected, cancelled)"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Success 200 {object} response.APIResponse{data=[]models.ManualPunch} "Manual punches retrieved successfully"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /attendance/manual-punches/my [get]
func (h *ManualPunchHandler) GetMyManualPunches(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	status := models.ManualPunchStatus(c.Query("status"))
	var startDate, endDate *time.Time

	// Parse start date filter
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsedDate
		}
	}

	// Parse end date filter
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsedDate
		}
	}

	punches, err := h.manualPunchService.GetEmployeeManualPunches(user.ID, status, startDate, endDate)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve manual punches", err)
		return
	}

	response.Success(c, "Manual punches retrieved successfully", punches)
}

// GetManualPunchByID godoc
// @Summary Get manual punch by ID
// @Description Get a specific manual punch request by ID
// @Tags Attendance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Manual punch ID"
// @Success 200 {object} response.APIResponse{data=models.ManualPunch} "Manual punch retrieved successfully"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 404 {object} response.APIResponse "Manual punch not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /attendance/manual-punches/{id} [get]
func (h *ManualPunchHandler) GetManualPunchByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid manual punch ID", err)
		return
	}

	manualPunch, err := h.manualPunchService.GetManualPunchByID(uint(id))
	if err != nil {
		response.NotFound(c, "Manual punch not found")
		return
	}

	response.Success(c, "Manual punch retrieved successfully", manualPunch)
}

// GetAllManualPunches godoc
// @Summary Get all manual punch requests (Admin/HR)
// @Description Get all manual punch requests with optional filtering and pagination
// @Tags Attendance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param employee_id query int false "Filter by employee ID"
// @Param status query string false "Filter by status (pending, approved, rejected, cancelled)"
// @Param punch_type query string false "Filter by punch type (in, out)"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} response.APIResponse{data=[]models.ManualPunch} "Manual punches retrieved successfully"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /attendance/manual-punches [get]
func (h *ManualPunchHandler) GetAllManualPunches(c *gin.Context) {
	page, limit := getPagination(c)

	filter := models.ManualPunchFilter{
		Page:  page,
		Limit: limit,
	}

	// Parse employee ID filter
	if employeeIDStr := c.Query("employee_id"); employeeIDStr != "" {
		if employeeID, err := strconv.Atoi(employeeIDStr); err == nil {
			employeeIDUint := uint(employeeID)
			filter.EmployeeID = &employeeIDUint
		}
	}

	// Parse status filter
	if status := c.Query("status"); status != "" {
		filter.Status = models.ManualPunchStatus(status)
	}

	// Parse punch type filter
	if punchType := c.Query("punch_type"); punchType != "" {
		filter.PunchType = models.ManualPunchType(punchType)
	}

	// Parse start date filter
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &parsedDate
		}
	}

	// Parse end date filter
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &parsedDate
		}
	}

	punches, total, err := h.manualPunchService.GetAllManualPunches(filter)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve manual punches", err)
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}
	response.SuccessWithMeta(c, "Manual punches retrieved successfully", punches, meta)
}

// GetPendingManualPunches godoc
// @Summary Get pending manual punch requests (Admin/HR)
// @Description Get all pending manual punch requests for approval
// @Tags Attendance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.APIResponse{data=[]models.ManualPunch} "Pending manual punches retrieved successfully"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /attendance/manual-punches/pending [get]
func (h *ManualPunchHandler) GetPendingManualPunches(c *gin.Context) {
	punches, err := h.manualPunchService.GetPendingManualPunches()
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve pending manual punches", err)
		return
	}

	response.Success(c, "Pending manual punches retrieved successfully", punches)
}

// UpdateManualPunchStatus godoc
// @Summary Approve/Reject manual punch request (Admin/HR)
// @Description Approve or reject a manual punch request
// @Tags Attendance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Manual punch ID"
// @Param request body models.ManualPunchApprovalRequest true "Approval decision"
// @Success 200 {object} response.APIResponse{data=models.ManualPunch} "Manual punch status updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid request"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden"
// @Failure 404 {object} response.APIResponse "Manual punch not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /attendance/manual-punches/{id}/status [patch]
func (h *ManualPunchHandler) UpdateManualPunchStatus(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid manual punch ID", err)
		return
	}

	var request models.ManualPunchApprovalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Invalid request format", err)
		return
	}

	manualPunch, err := h.manualPunchService.UpdateManualPunchStatus(uint(id), user, request)
	if err != nil {
		response.InternalServerError(c, "Failed to update manual punch status", err)
		return
	}

	response.Success(c, "Manual punch status updated successfully", manualPunch)
}

// CancelManualPunch godoc
// @Summary Cancel manual punch request
// @Description Cancel a pending manual punch request
// @Tags Attendance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Manual punch ID"
// @Success 200 {object} response.APIResponse "Manual punch cancelled successfully"
// @Failure 400 {object} response.APIResponse "Invalid request"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden"
// @Failure 404 {object} response.APIResponse "Manual punch not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /attendance/manual-punches/{id}/cancel [post]
func (h *ManualPunchHandler) CancelManualPunch(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid manual punch ID", err)
		return
	}

	err = h.manualPunchService.CancelManualPunch(uint(id), user.ID)
	if err != nil {
		response.InternalServerError(c, "Failed to cancel manual punch", err)
		return
	}

	response.Success(c, "Manual punch cancelled successfully", nil)
}

// DeleteManualPunch godoc
// @Summary Delete manual punch request (Admin/HR)
// @Description Delete a manual punch request (soft delete)
// @Tags Attendance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Manual punch ID"
// @Success 200 {object} response.APIResponse "Manual punch deleted successfully"
// @Failure 400 {object} response.APIResponse "Invalid request"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden"
// @Failure 404 {object} response.APIResponse "Manual punch not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /attendance/manual-punches/{id} [delete]
func (h *ManualPunchHandler) DeleteManualPunch(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid manual punch ID", err)
		return
	}

	err = h.manualPunchService.DeleteManualPunch(uint(id))
	if err != nil {
		response.InternalServerError(c, "Failed to delete manual punch", err)
		return
	}

	response.Success(c, "Manual punch deleted successfully", nil)
}

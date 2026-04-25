package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type LeaveRequestHandler struct {
	service        *services.LeaveRequestService
	balanceService *services.LeaveBalanceService
	employeeRepo   *employeeRepos.EmployeeRepository
}

type leaveRequestListItem struct {
	models.LeaveRequest
	EmployeeFirstName  string  `json:"employee_first_name"`
	EmployeeMiddleName *string `json:"employee_middle_name,omitempty"`
	EmployeeLastName   string  `json:"employee_last_name"`
}

func NewLeaveRequestHandler() *LeaveRequestHandler {
	return &LeaveRequestHandler{
		service:        services.NewLeaveRequestService(),
		balanceService: services.NewLeaveBalanceService(),
		employeeRepo:   employeeRepos.NewEmployeeRepository(),
	}
}

// GetEmployeeInfo handles getting employee information for form
func (h *LeaveRequestHandler) GetEmployeeInfo(c *gin.Context) {
	user, _ := c.Get("user")
	var employeeID string

	if userObj, ok := user.(*userModels.User); ok {
		employee, err := h.employeeRepo.FindByUserID(userObj.ID)
		if err != nil || employee == nil {
			response.BadRequest(c, "Employee not found for user", nil)
			return
		}
		employeeID = employee.EmployeeID
	} else {
		// Allow query parameter for HR/Admin
		employeeID = c.Query("employee_id")
		if employeeID == "" {
			response.Unauthorized(c, "User not authenticated")
			return
		}
	}

	employee, err := h.employeeRepo.FindByEmployeeID(employeeID)
	if err != nil {
		response.NotFound(c, "Employee not found")
		return
	}

	// Get emergency contact
	emergencyContact := ""
	emergencyContactNumber := ""
	// TODO: Get from employee_emergency_contacts table

	response.Success(c, "Employee information retrieved successfully", map[string]interface{}{
		"employee_id":              employee.EmployeeID,
		"first_name":               employee.FirstName,
		"middle_name":              employee.MiddleName,
		"last_name":                employee.LastName,
		"job_position":             "", // TODO: Get from position
		"department":               "", // TODO: Get from department
		"contact_number":           employee.PhoneNumber,
		"emergency_contact_person": emergencyContact,
		"emergency_contact_number": emergencyContactNumber,
	})
}

// GetEmployeeLeaveBalances handles getting employee leave balances
func (h *LeaveRequestHandler) GetEmployeeLeaveBalances(c *gin.Context) {
	user, _ := c.Get("user")
	var employeeID string
	year := 0

	if userObj, ok := user.(*userModels.User); ok {
		employee, err := h.employeeRepo.FindByUserID(userObj.ID)
		if err != nil || employee == nil {
			response.BadRequest(c, "Employee not found for user", nil)
			return
		}
		employeeID = employee.EmployeeID
	} else {
		employeeID = c.Query("employee_id")
		if employeeID == "" {
			response.Unauthorized(c, "User not authenticated")
			return
		}
	}

	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	tenantID := middleware.GetTenantID(c)
	balances, err := h.balanceService.GetEmployeeLeaveBalances(employeeID, year, tenantID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve leave balances", err.Error())
		return
	}

	response.Success(c, "Leave balances retrieved successfully", balances)
}

// CalculateLeaveDays handles calculating leave days
func (h *LeaveRequestHandler) CalculateLeaveDays(c *gin.Context) {
	var req struct {
		FromDate   string `json:"from_date" binding:"required"`
		ToDate     string `json:"to_date" binding:"required"`
		HalfDay    bool   `json:"half_day"`
		EmployeeID string `json:"employee_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	employeeID := req.EmployeeID
	if employeeID == "" {
		if userObj, ok := user.(*userModels.User); ok {
			employee, err := h.employeeRepo.FindByUserID(userObj.ID)
			if err == nil && employee != nil {
				employeeID = employee.EmployeeID
			}
		}
	}

	// Parse dates
	fromDate, err := time.Parse("2006-01-02", req.FromDate)
	if err != nil {
		response.BadRequest(c, "Invalid from_date format. Use YYYY-MM-DD", nil)
		return
	}

	toDate, err := time.Parse("2006-01-02", req.ToDate)
	if err != nil {
		response.BadRequest(c, "Invalid to_date format. Use YYYY-MM-DD", nil)
		return
	}

	tenantID := middleware.GetTenantID(c)
	totalDays, weekends, holidays, err := h.service.CalculateLeaveDays(
		fromDate, toDate, req.HalfDay, employeeID, tenantID,
	)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Days calculated successfully", map[string]interface{}{
		"total_days":   totalDays,
		"working_days": totalDays,
		"weekends":     weekends,
		"holidays":     holidays,
		"half_day":     req.HalfDay,
	})
}

// CreateLeaveRequest handles creating a leave request
func (h *LeaveRequestHandler) CreateLeaveRequest(c *gin.Context) {
	var req models.CreateLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var employeeID string
	var createdBy *uint

	if userObj, ok := user.(*userModels.User); ok {
		createdBy = &userObj.ID
		employee, err := h.employeeRepo.FindByUserID(userObj.ID)
		if err != nil || employee == nil {
			response.BadRequest(c, "Employee not found for user", nil)
			return
		}
		employeeID = employee.EmployeeID
	} else {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	tenantID := middleware.GetTenantID(c)
	leaveRequest, err := h.service.CreateLeaveRequest(&req, employeeID, tenantID, createdBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Leave application submitted successfully", leaveRequest)
}

// ListLeaveRequests handles listing leave requests
func (h *LeaveRequestHandler) ListLeaveRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)

	// Build filters
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if employeeID := c.Query("employee_id"); employeeID != "" {
		filters["employee_id"] = employeeID
	}
	if leaveTypeCode := c.Query("leave_type_code"); leaveTypeCode != "" {
		filters["leave_type_code"] = leaveTypeCode
	}
	if fromDate := c.Query("from_date"); fromDate != "" {
		filters["from_date"] = fromDate
	}
	if toDate := c.Query("to_date"); toDate != "" {
		filters["to_date"] = toDate
	}
	if search := strings.TrimSpace(c.Query("search")); search != "" {
		filters["search"] = search
	}

	// For employees, filter by their own requests
	user, _ := c.Get("user")
	if userObj, ok := user.(*userModels.User); ok {
		if userObj.Role != userModels.RoleAdmin && userObj.Role != userModels.RoleHR {
			employee, err := h.employeeRepo.FindByUserID(userObj.ID)
			if err == nil && employee != nil {
				filters["employee_id"] = employee.EmployeeID
			}
		}
	}

	requests, total, err := h.service.ListLeaveRequests(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve leave requests", err.Error())
		return
	}

	employeeIDs := make([]string, 0)
	seen := make(map[string]struct{})
	for _, r := range requests {
		if r.EmployeeID == "" {
			continue
		}
		if _, ok := seen[r.EmployeeID]; ok {
			continue
		}
		seen[r.EmployeeID] = struct{}{}
		employeeIDs = append(employeeIDs, r.EmployeeID)
	}

	employees, err := h.employeeRepo.FindByEmployeeIDs(employeeIDs)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve employee details", err.Error())
		return
	}
	employeeByID := make(map[string]employeeModels.Employee, len(employees))
	for _, e := range employees {
		employeeByID[e.EmployeeID] = e
	}

	items := make([]leaveRequestListItem, 0, len(requests))
	for i := range requests {
		r := requests[i]
		e, ok := employeeByID[r.EmployeeID]
		if ok {
			items = append(items, leaveRequestListItem{
				LeaveRequest:       r,
				EmployeeFirstName:  e.FirstName,
				EmployeeMiddleName: e.MiddleName,
				EmployeeLastName:   e.LastName,
			})
		} else {
			items = append(items, leaveRequestListItem{
				LeaveRequest:      r,
				EmployeeFirstName: "",
				EmployeeLastName:  "",
			})
		}
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}

	response.SuccessWithMeta(c, "Leave requests retrieved successfully", items, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int(totalPages),
	})
}

// GetLeaveRequest handles getting a leave request by ID
func (h *LeaveRequestHandler) GetLeaveRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	request, err := h.service.GetLeaveRequestByID(uint(id))
	if err != nil {
		response.NotFound(c, "Leave request not found")
		return
	}

	response.Success(c, "Leave request details retrieved successfully", request)
}

// ApproveLeaveRequest handles approving a leave request
func (h *LeaveRequestHandler) ApproveLeaveRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	var req models.ApproveLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var approverID string
	if userObj, ok := user.(*userModels.User); ok {
		employee, err := h.employeeRepo.FindByUserID(userObj.ID)
		if err == nil && employee != nil {
			approverID = employee.EmployeeID
		}
	}

	tenantID := middleware.GetTenantID(c)
	request, err := h.service.ApproveLeaveRequest(uint(id), &req, approverID, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Leave request approved successfully", request)
}

// RejectLeaveRequest handles rejecting a leave request
func (h *LeaveRequestHandler) RejectLeaveRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	var req models.RejectLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var approverID string
	if userObj, ok := user.(*userModels.User); ok {
		employee, err := h.employeeRepo.FindByUserID(userObj.ID)
		if err == nil && employee != nil {
			approverID = employee.EmployeeID
		}
	}

	request, err := h.service.RejectLeaveRequest(uint(id), &req, approverID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Leave request rejected successfully", request)
}

// CancelLeaveRequest handles cancelling a leave request
func (h *LeaveRequestHandler) CancelLeaveRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	var req models.CancelLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	request, err := h.service.CancelLeaveRequest(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Leave request cancelled successfully", request)
}

// UpdateLeaveRequest handles updating a leave request
func (h *LeaveRequestHandler) UpdateLeaveRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	var req models.UpdateLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	request, err := h.service.UpdateLeaveRequest(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Leave request modified successfully", request)
}

// DeleteLeaveRequest handles deleting a draft leave request
func (h *LeaveRequestHandler) DeleteLeaveRequest(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	if err := h.service.DeleteLeaveRequest(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Draft leave request deleted successfully", nil)
}

// ReturnForInfo handles returning a leave request to the employee for more information
// @Summary Return leave request for information
// @Description HR/Admin can return a pending leave request to the employee requesting additional information
// @Tags Leave Management (HR/Admin)
// @Accept json
// @Produce json
// @Param request_id path int true "Leave Request ID"
// @Param request body models.ReturnForInfoRequest true "Information required"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/leave/requests/{request_id}/return-for-info [post]
func (h *LeaveRequestHandler) ReturnForInfo(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid request ID", nil)
		return
	}

	var req models.ReturnForInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	request, err := h.service.ReturnForInfo(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Leave request returned for information successfully", request)
}

// ListAllLeaveRequests handles listing all leave requests for HR/Admin
// @Summary List all leave requests (HR/Admin)
// @Description Retrieve a paginated list of all leave requests across all employees for HR/Admin review
// @Tags Leave Management (HR/Admin)
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 20, max: 100)"
// @Param status query string false "Filter: pending, approved, rejected, cancelled, draft, returned_for_info"
// @Param employee_id query string false "Filter by employee ID"
// @Param leave_type_code query string false "Filter by leave type code"
// @Param from_date query string false "Filter from date (YYYY-MM-DD)"
// @Param to_date query string false "Filter to date (YYYY-MM-DD)"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/leave/admin/requests [get]
func (h *LeaveRequestHandler) ListAllLeaveRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)

	// Build filters — do NOT restrict by employee_id (HR sees all)
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if employeeID := c.Query("employee_id"); employeeID != "" {
		filters["employee_id"] = employeeID
	}
	if leaveTypeCode := c.Query("leave_type_code"); leaveTypeCode != "" {
		filters["leave_type_code"] = leaveTypeCode
	}
	if fromDate := c.Query("from_date"); fromDate != "" {
		filters["from_date"] = fromDate
	}
	if toDate := c.Query("to_date"); toDate != "" {
		filters["to_date"] = toDate
	}
	if search := strings.TrimSpace(c.Query("search")); search != "" {
		filters["search"] = search
	}

	requests, total, err := h.service.ListLeaveRequests(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve leave requests", err.Error())
		return
	}

	employeeIDs := make([]string, 0)
	seen := make(map[string]struct{})
	for _, r := range requests {
		if r.EmployeeID == "" {
			continue
		}
		if _, ok := seen[r.EmployeeID]; ok {
			continue
		}
		seen[r.EmployeeID] = struct{}{}
		employeeIDs = append(employeeIDs, r.EmployeeID)
	}

	employees, err := h.employeeRepo.FindByEmployeeIDs(employeeIDs)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve employee details", err.Error())
		return
	}
	employeeByID := make(map[string]employeeModels.Employee, len(employees))
	for _, e := range employees {
		employeeByID[e.EmployeeID] = e
	}

	items := make([]leaveRequestListItem, 0, len(requests))
	for i := range requests {
		r := requests[i]
		e, ok := employeeByID[r.EmployeeID]
		if ok {
			items = append(items, leaveRequestListItem{
				LeaveRequest:       r,
				EmployeeFirstName:  e.FirstName,
				EmployeeMiddleName: e.MiddleName,
				EmployeeLastName:   e.LastName,
			})
		} else {
			items = append(items, leaveRequestListItem{
				LeaveRequest:      r,
				EmployeeFirstName: "",
				EmployeeLastName:  "",
			})
		}
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}

	response.SuccessWithMeta(c, "Leave requests retrieved successfully", items, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int(totalPages),
	})
}

// GetActiveLeaveTypes returns active leave types for the employee leave form dropdown
// @Summary Get active leave types (Employee)
// @Description Returns active leave types for the leave request form dropdown
// @Tags Leave (Employee)
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 20)"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/leave/types [get]
func (h *LeaveRequestHandler) GetActiveLeaveTypes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	tenantID := middleware.GetTenantID(c)
	filters := map[string]interface{}{"is_active": true}

	leaveTypes, total, err := h.service.ListActiveLeaveTypes(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve leave types", err.Error())
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}

	response.Success(c, "Leave types retrieved successfully", map[string]interface{}{
		"data": leaveTypes,
		"meta": map[string]interface{}{
			"page":        page,
			"per_page":    pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetLeaveHolidays returns holidays for the leave calendar (employee use)
// @Summary Get holidays (Employee)
// @Description Returns holidays for the current or specified year, used for leave day calculation
// @Tags Leave (Employee)
// @Produce json
// @Param year query int false "Year (default: current year)"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/leave/holidays [get]
func (h *LeaveRequestHandler) GetLeaveHolidays(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	year := time.Now().Year()
	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil && y > 2000 && y < 2100 {
			year = y
		}
	}

	holidays, err := h.service.GetHolidaysByYear(year, tenantID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve holidays", err.Error())
		return
	}

	response.Success(c, "Holidays retrieved successfully", holidays)
}

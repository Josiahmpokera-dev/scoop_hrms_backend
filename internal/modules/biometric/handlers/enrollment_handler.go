package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type EnrollmentHandler struct {
	service *services.EnrollmentService
}

func NewEnrollmentHandler() *EnrollmentHandler {
	return &EnrollmentHandler{service: services.NewEnrollmentService()}
}

func callerID(c *gin.Context) uint {
	if v, exists := c.Get("user_id"); exists {
		switch id := v.(type) {
		case uint:
			return id
		case float64:
			return uint(id)
		}
	}
	return 0
}

// POST /biometric/enrollments/link — link one employee to a biometric emp_code
func (h *EnrollmentHandler) Link(c *gin.Context) {
	var req models.LinkEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	enrollment, err := h.service.LinkEmployee(&req, callerID(c))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Employee linked to biometric device successfully", enrollment)
}

// POST /biometric/enrollments/bulk-link — link multiple employees at once
func (h *EnrollmentHandler) BulkLink(c *gin.Context) {
	var req models.BulkLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	linked, errs, _ := h.service.BulkLink(&req, callerID(c))
	response.Success(c, "Bulk link completed", gin.H{
		"linked": linked,
		"errors": errs,
	})
}

// DELETE /biometric/enrollments/:employeeId/unlink — unlink employee
func (h *EnrollmentHandler) Unlink(c *gin.Context) {
	employeeID, err := strconv.ParseUint(c.Param("employeeId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid employee ID", nil)
		return
	}

	if err := h.service.UnlinkEmployee(uint(employeeID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Employee unlinked from biometric device successfully", nil)
}

// POST /biometric/enrollments/auto-link — auto-link by matching EmployeeID = EmpCode
func (h *EnrollmentHandler) AutoLink(c *gin.Context) {
	result, err := h.service.AutoLink(callerID(c))
	if err != nil {
		response.InternalServerError(c, "Auto-link failed", err.Error())
		return
	}
	response.Success(c, "Auto-link completed", result)
}

// GET /biometric/enrollments — list all linked employees
func (h *EnrollmentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	enrollments, total, err := h.service.ListEnrollments(page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list enrollments", err.Error())
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response.SuccessWithMeta(c, "Enrollments retrieved successfully", enrollments, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GET /biometric/enrollments/unlinked — list employees not yet linked
func (h *EnrollmentHandler) GetUnlinked(c *gin.Context) {
	employees, err := h.service.GetUnlinkedEmployees()
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve unlinked employees", err.Error())
		return
	}
	response.Success(c, "Unlinked employees retrieved successfully", employees)
}

// GET /biometric/enrollments/device-users — list biometric device users (with linked status)
func (h *EnrollmentHandler) GetDeviceUsers(c *gin.Context) {
	users, err := h.service.GetDeviceUsers()
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve device users", err.Error())
		return
	}
	response.Success(c, "Device users retrieved successfully", users)
}

// GET /biometric/enrollments/statistics — enrollment coverage stats
func (h *EnrollmentHandler) GetStatistics(c *gin.Context) {
	stats, err := h.service.GetStatistics()
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve statistics", err.Error())
		return
	}
	response.Success(c, "Enrollment statistics retrieved successfully", stats)
}

// GET /biometric/attendance/merged — merged attendance with employee data
func (h *EnrollmentHandler) GetMergedAttendance(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")
	if startDate == "" || endDate == "" {
		response.BadRequest(c, "start_date and end_date query parameters are required (YYYY-MM-DD)", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var employeeID *uint
	if eid := c.Query("employee_id"); eid != "" {
		id, err := strconv.ParseUint(eid, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid employee_id", nil)
			return
		}
		uid := uint(id)
		employeeID = &uid
	}

	records, total, totalPages, err := h.service.GetMergedAttendance(startDate, endDate, employeeID, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve merged attendance", err.Error())
		return
	}

	response.SuccessWithMeta(c, "Merged attendance retrieved successfully", records, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GET /biometric/attendance/merged/:employeeId — merged attendance for a single employee
func (h *EnrollmentHandler) GetEmployeeAttendance(c *gin.Context) {
	eid, err := strconv.ParseUint(c.Param("employeeId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid employee ID", nil)
		return
	}

	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")
	if startDate == "" || endDate == "" {
		response.BadRequest(c, "start_date and end_date query parameters are required (YYYY-MM-DD)", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	uid := uint(eid)
	records, total, totalPages, err := h.service.GetMergedAttendance(startDate, endDate, &uid, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve employee attendance", err.Error())
		return
	}

	response.SuccessWithMeta(c, "Employee attendance retrieved successfully", records, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

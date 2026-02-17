package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type EmployeeTargetHandler struct {
	service *services.EmployeeTargetService
}

func NewEmployeeTargetHandler() *EmployeeTargetHandler {
	return &EmployeeTargetHandler{service: services.NewEmployeeTargetService()}
}

func (h *EmployeeTargetHandler) ListTargets(c *gin.Context) {
	employeeID := c.Query("employee_id")
	department := c.Query("department")
	deptTargetID := c.Query("department_target_id")
	status := c.Query("status")
	period := c.Query("period")
	assignedBy := c.Query("assigned_by")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	targets, total, err := h.service.List(employeeID, department, deptTargetID, status, period, assignedBy, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list employee targets", err.Error())
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{Page: page, PerPage: pageSize, Total: total, TotalPages: totalPages}
	response.SuccessWithMeta(c, "Employee targets retrieved successfully", targets, meta)
}

func (h *EmployeeTargetHandler) GetTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid target ID", nil)
		return
	}
	target, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "Employee target not found")
		return
	}
	response.Success(c, "Employee target retrieved successfully", target)
}

func (h *EmployeeTargetHandler) CreateTarget(c *gin.Context) {
	var req models.CreateEmployeeTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	userName := ""
	target, err := h.service.Create(callerID, userName, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "Employee target created successfully", target)
}

func (h *EmployeeTargetHandler) BulkAssignTargets(c *gin.Context) {
	var req models.BulkAssignTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	userName := ""
	targets, err := h.service.BulkAssign(callerID, userName, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "Targets assigned successfully", targets)
}

func (h *EmployeeTargetHandler) UpdateTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid target ID", nil)
		return
	}
	var req models.UpdateEmployeeTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	target, err := h.service.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Employee target updated successfully", target)
}

func (h *EmployeeTargetHandler) DeleteTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid target ID", nil)
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Employee target deleted successfully", nil)
}

func (h *EmployeeTargetHandler) UpdateProgress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid target ID", nil)
		return
	}
	var req models.UpdateEmployeeTargetProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	target, err := h.service.UpdateProgress(uint(id), req.CurrentValue, "", "")
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Progress updated successfully", target)
}

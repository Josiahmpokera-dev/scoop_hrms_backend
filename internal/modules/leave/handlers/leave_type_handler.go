package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
)

type LeaveTypeHandler struct {
	service *services.LeaveTypeService
}

func NewLeaveTypeHandler() *LeaveTypeHandler {
	return &LeaveTypeHandler{
		service: services.NewLeaveTypeService(),
	}
}

// ListLeaveTypes handles listing leave types
func (h *LeaveTypeHandler) ListLeaveTypes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)

	// Build filters
	filters := make(map[string]interface{})
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filters["is_active"] = isActive
		}
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}

	leaveTypes, total, err := h.service.ListLeaveTypes(tenantID, page, pageSize, filters)
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
			"page":       page,
			"per_page":   pageSize,
			"total":      total,
			"total_pages": totalPages,
		},
	})
}

// GetLeaveType handles getting a leave type by ID
func (h *LeaveTypeHandler) GetLeaveType(c *gin.Context) {
	idStr := c.Param("type_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid leave type ID", nil)
		return
	}

	leaveType, err := h.service.GetLeaveTypeByID(uint(id))
	if err != nil {
		response.NotFound(c, "Leave type not found")
		return
	}

	response.Success(c, "Leave type retrieved successfully", leaveType)
}

// CreateLeaveType handles creating a leave type
func (h *LeaveTypeHandler) CreateLeaveType(c *gin.Context) {
	var req models.CreateLeaveTypeRequest
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

	leaveType, err := h.service.CreateLeaveType(&req, tenantID, createdBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Leave type created successfully", leaveType)
}

// UpdateLeaveType handles updating a leave type
func (h *LeaveTypeHandler) UpdateLeaveType(c *gin.Context) {
	idStr := c.Param("type_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid leave type ID", nil)
		return
	}

	var req models.UpdateLeaveTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	leaveType, err := h.service.UpdateLeaveType(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Leave type updated successfully", leaveType)
}

// DeleteLeaveType handles deleting a leave type
func (h *LeaveTypeHandler) DeleteLeaveType(c *gin.Context) {
	idStr := c.Param("type_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid leave type ID", nil)
		return
	}

	if err := h.service.DeleteLeaveType(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Leave type deleted successfully", nil)
}

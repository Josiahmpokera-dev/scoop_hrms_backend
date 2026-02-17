package handlers

import (
	"fmt"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type DepartmentTargetHandler struct {
	service *services.DepartmentTargetService
}

func NewDepartmentTargetHandler() *DepartmentTargetHandler {
	return &DepartmentTargetHandler{service: services.NewDepartmentTargetService()}
}

func (h *DepartmentTargetHandler) ListTargets(c *gin.Context) {
	department := c.Query("department")
	category := c.Query("category")
	status := c.Query("status")
	period := c.Query("period")
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
	targets, total, err := h.service.List(department, category, status, period, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list department targets", err.Error())
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{Page: page, PerPage: pageSize, Total: total, TotalPages: totalPages}
	response.SuccessWithMeta(c, "Department targets retrieved successfully", targets, meta)
}

func (h *DepartmentTargetHandler) GetTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid target ID", nil)
		return
	}
	target, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "Department target not found")
		return
	}
	response.Success(c, "Department target retrieved successfully", target)
}

func (h *DepartmentTargetHandler) CreateTarget(c *gin.Context) {
	var req models.CreateDepartmentTargetRequest
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
	response.Created(c, "Department target created successfully", target)
}

func (h *DepartmentTargetHandler) UpdateTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid target ID", nil)
		return
	}
	var req models.UpdateDepartmentTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	target, err := h.service.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Department target updated successfully", target)
}

func (h *DepartmentTargetHandler) DeleteTarget(c *gin.Context) {
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
	response.Success(c, "Department target deleted successfully", nil)
}

func (h *DepartmentTargetHandler) UpdateProgress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid target ID", nil)
		return
	}
	var req models.UpdateDepartmentTargetProgressRequest
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

func (h *DepartmentTargetHandler) CompleteMilestone(c *gin.Context) {
	idStr := c.Param("id")
	milestoneIDStr := c.Param("milestoneId")
	id, err1 := strconv.ParseUint(idStr, 10, 32)
	milestoneID, err2 := strconv.ParseUint(milestoneIDStr, 10, 32)
	if err1 != nil || err2 != nil {
		response.BadRequest(c, "Invalid target or milestone ID", nil)
		return
	}
	milestone, err := h.service.CompleteMilestone(uint(id), uint(milestoneID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Milestone completed successfully", milestone)
}

func (h *DepartmentTargetHandler) LinkGoal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid target ID", nil)
		return
	}
	var req models.LinkDepartmentTargetGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	goalIDStr := fmt.Sprintf("%d", req.GoalID)
	target, err := h.service.LinkGoal(uint(id), goalIDStr)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Goal linked successfully", target)
}

func (h *DepartmentTargetHandler) LinkProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid target ID", nil)
		return
	}
	var req models.LinkDepartmentTargetProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	linkReq := &models.LinkProjectRequest{
		ProjectID:   req.ProjectID,
		ProjectCode: req.ProjectCode,
		ProjectName: req.ProjectName,
	}
	target, err := h.service.LinkProject(uint(id), linkReq)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Project linked successfully", target)
}

func (h *DepartmentTargetHandler) UnlinkProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid target ID", nil)
		return
	}
	target, err := h.service.UnlinkProject(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Project unlinked successfully", target)
}

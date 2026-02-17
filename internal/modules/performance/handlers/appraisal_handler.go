package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type AppraisalHandler struct {
	service *services.AppraisalService
}

func NewAppraisalHandler() *AppraisalHandler {
	return &AppraisalHandler{service: services.NewAppraisalService()}
}

func (h *AppraisalHandler) ListCycles(c *gin.Context) {
	status := c.Query("status")
	cycleType := c.Query("type")
	yearStr := c.Query("year")
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
	var year int
	if yearStr != "" {
		year, _ = strconv.Atoi(yearStr)
	}
	cycles, total, err := h.service.ListCycles(status, cycleType, year, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list appraisal cycles", err.Error())
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{Page: page, PerPage: pageSize, Total: total, TotalPages: totalPages}
	response.SuccessWithMeta(c, "Appraisal cycles retrieved successfully", cycles, meta)
}

func (h *AppraisalHandler) GetActiveCycle(c *gin.Context) {
	cycle, err := h.service.GetActiveCycle()
	if err != nil {
		response.NotFound(c, "No active appraisal cycle")
		return
	}
	response.Success(c, "Active cycle retrieved successfully", cycle)
}

func (h *AppraisalHandler) GetCycle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid cycle ID", nil)
		return
	}
	cycle, err := h.service.GetCycleByID(uint(id))
	if err != nil {
		response.NotFound(c, "Appraisal cycle not found")
		return
	}
	response.Success(c, "Appraisal cycle retrieved successfully", cycle)
}

func (h *AppraisalHandler) CreateCycle(c *gin.Context) {
	var req models.CreateAppraisalCycleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	cycle, err := h.service.CreateCycle(callerID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "Appraisal cycle created successfully", cycle)
}

func (h *AppraisalHandler) UpdateCycle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid cycle ID", nil)
		return
	}
	var req models.UpdateAppraisalCycleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	cycle, err := h.service.UpdateCycle(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Appraisal cycle updated successfully", cycle)
}

func (h *AppraisalHandler) ChangeCycleStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid cycle ID", nil)
		return
	}
	var req struct {
		Status string `json:"status" binding:"required,oneof=Draft Active Calibration Completed Closed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	cycle, err := h.service.ChangeCycleStatus(uint(id), req.Status)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Cycle status updated successfully", cycle)
}

func (h *AppraisalHandler) DeleteCycle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid cycle ID", nil)
		return
	}
	if err := h.service.DeleteCycle(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Appraisal cycle deleted successfully", nil)
}

func (h *AppraisalHandler) ListAppraisals(c *gin.Context) {
	cycleIDStr := c.Query("cycle_id")
	employeeIDStr := c.Query("employee_id")
	department := c.Query("department")
	status := c.Query("status")
	mineStr := c.Query("mine")
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

	var cycleIDPtr *uint
	if cycleIDStr != "" {
		if cid, e := strconv.ParseUint(cycleIDStr, 10, 32); e == nil {
			cidVal := uint(cid)
			cycleIDPtr = &cidVal
		}
	}
	var employeeIDPtr *uint
	if employeeIDStr != "" {
		if eid, e := strconv.ParseUint(employeeIDStr, 10, 32); e == nil {
			eidVal := uint(eid)
			employeeIDPtr = &eidVal
		}
	}
	mine := mineStr == "true" || mineStr == "1"

	appraisals, total, err := h.service.ListAppraisals(employeeIDPtr, cycleIDPtr, status, department, mine, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list appraisals", err.Error())
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{Page: page, PerPage: pageSize, Total: total, TotalPages: totalPages}
	response.SuccessWithMeta(c, "Appraisals retrieved successfully", appraisals, meta)
}

func (h *AppraisalHandler) GetAppraisalSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	summary, err := h.service.GetAppraisalSummary(callerID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve appraisal summary", err.Error())
		return
	}
	response.Success(c, "Appraisal summary retrieved successfully", summary)
}

func (h *AppraisalHandler) GetAppraisal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appraisal ID", nil)
		return
	}
	appraisal, err := h.service.GetAppraisalByID(uint(id))
	if err != nil {
		response.NotFound(c, "Appraisal not found")
		return
	}
	response.Success(c, "Appraisal retrieved successfully", appraisal)
}

func (h *AppraisalHandler) SubmitSelfReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appraisal ID", nil)
		return
	}
	var req models.SubmitSelfReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	userName := ""
	appraisal, err := h.service.SubmitSelfReview(uint(id), &req, userName)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Self review submitted successfully", appraisal)
}

func (h *AppraisalHandler) SubmitManagerReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appraisal ID", nil)
		return
	}
	var req models.SubmitManagerReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	userName := ""
	appraisal, err := h.service.SubmitManagerReview(uint(id), &req, userName)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Manager review submitted successfully", appraisal)
}

func (h *AppraisalHandler) SubmitCalibration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appraisal ID", nil)
		return
	}
	var req models.SubmitCalibrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	appraisal, err := h.service.SubmitCalibration(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Calibration submitted successfully", appraisal)
}

func (h *AppraisalHandler) SendBack(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appraisal ID", nil)
		return
	}
	var req models.SendBackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	appraisal, err := h.service.SendBack(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Appraisal sent back successfully", appraisal)
}

func (h *AppraisalHandler) Finalize(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appraisal ID", nil)
		return
	}
	var req models.FinalizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	appraisal, err := h.service.Finalize(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Appraisal finalized successfully", appraisal)
}

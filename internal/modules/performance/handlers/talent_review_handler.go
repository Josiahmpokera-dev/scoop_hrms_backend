package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type TalentReviewHandler struct {
	service *services.TalentReviewService
}

func NewTalentReviewHandler() *TalentReviewHandler {
	return &TalentReviewHandler{service: services.NewTalentReviewService()}
}

func (h *TalentReviewHandler) ListReviews(c *gin.Context) {
	department := c.Query("department")
	riskOfLoss := c.Query("risk_of_loss")
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
	reviews, total, err := h.service.List(department, nil, nil, riskOfLoss, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list talent reviews", err.Error())
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{Page: page, PerPage: pageSize, Total: total, TotalPages: totalPages}
	response.SuccessWithMeta(c, "Talent reviews retrieved successfully", reviews, meta)
}

func (h *TalentReviewHandler) GetReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid employee ID", nil)
		return
	}
	review, err := h.service.GetByEmployeeID(uint(id))
	if err != nil {
		response.NotFound(c, "Talent review not found")
		return
	}
	response.Success(c, "Talent review retrieved successfully", review)
}

func (h *TalentReviewHandler) UpdateReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid employee ID", nil)
		return
	}
	var req models.UpdateTalentReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	review, err := h.service.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Talent review updated successfully", review)
}

func (h *TalentReviewHandler) CreateCalibrationSession(c *gin.Context) {
	var req models.CreateCalibrationSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	session, err := h.service.CreateCalibrationSession(callerID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "Calibration session created successfully", session)
}

func (h *TalentReviewHandler) AddSuccessionPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid employee ID", nil)
		return
	}
	var req models.AddSuccessionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	plan, err := h.service.AddSuccessionPlan(uint(id), callerID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "Succession plan added successfully", plan)
}

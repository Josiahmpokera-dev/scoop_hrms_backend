package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type Feedback360Handler struct {
	service *services.Feedback360Service
}

func NewFeedback360Handler() *Feedback360Handler {
	return &Feedback360Handler{service: services.NewFeedback360Service()}
}

func (h *Feedback360Handler) ListCampaigns(c *gin.Context) {
	status := c.Query("status")
	department := c.Query("department")
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
	campaigns, total, err := h.service.List(status, department, nil, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list campaigns", err.Error())
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{Page: page, PerPage: pageSize, Total: total, TotalPages: totalPages}
	response.SuccessWithMeta(c, "Campaigns retrieved successfully", campaigns, meta)
}

func (h *Feedback360Handler) GetCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid campaign ID", nil)
		return
	}
	campaign, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "Campaign not found")
		return
	}
	response.Success(c, "Campaign retrieved successfully", campaign)
}

func (h *Feedback360Handler) LaunchCampaign(c *gin.Context) {
	var req models.LaunchFeedback360Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	campaign, err := h.service.Launch(callerID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "Campaign launched successfully", campaign)
}

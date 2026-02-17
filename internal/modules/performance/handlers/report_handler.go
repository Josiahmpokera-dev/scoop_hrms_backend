package handlers

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{service: services.NewReportService()}
}

func (h *ReportHandler) GetSummary(c *gin.Context) {
	department := c.Query("department")
	summary, err := h.service.GetSummary(department)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve report summary", err.Error())
		return
	}
	response.Success(c, "Report summary retrieved successfully", summary)
}

func (h *ReportHandler) GetRatingDistribution(c *gin.Context) {
	department := c.Query("department")
	data, err := h.service.GetRatingDistribution(department)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve rating distribution", err.Error())
		return
	}
	response.Success(c, "Rating distribution retrieved successfully", data)
}

func (h *ReportHandler) GetDepartmentSummary(c *gin.Context) {
	data, err := h.service.GetDepartmentSummary()
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve department summary", err.Error())
		return
	}
	response.Success(c, "Department summary retrieved successfully", data)
}

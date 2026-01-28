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

type HolidayHandler struct {
	service *services.HolidayService
}

func NewHolidayHandler() *HolidayHandler {
	return &HolidayHandler{
		service: services.NewHolidayService(),
	}
}

// ListHolidays handles listing holidays
func (h *HolidayHandler) ListHolidays(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)

	// Build filters
	filters := make(map[string]interface{})
	if yearStr := c.Query("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			filters["year"] = year
		}
	}
	if holidayType := c.Query("type"); holidayType != "" {
		filters["type"] = holidayType
	}
	if isFloaterStr := c.Query("is_floater"); isFloaterStr != "" {
		if isFloater, err := strconv.ParseBool(isFloaterStr); err == nil {
			filters["is_floater"] = isFloater
		}
	}

	holidays, total, err := h.service.ListHolidays(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve holidays", err.Error())
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}

	response.Success(c, "Holidays retrieved successfully", map[string]interface{}{
		"data": holidays,
		"meta": map[string]interface{}{
			"page":       page,
			"per_page":   pageSize,
			"total":      total,
			"total_pages": totalPages,
		},
	})
}

// GetHoliday handles getting a holiday by ID
func (h *HolidayHandler) GetHoliday(c *gin.Context) {
	idStr := c.Param("holiday_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid holiday ID", nil)
		return
	}

	holiday, err := h.service.GetHolidayByID(uint(id))
	if err != nil {
		response.NotFound(c, "Holiday not found")
		return
	}

	response.Success(c, "Holiday details retrieved successfully", holiday)
}

// CreateHoliday handles creating a holiday
func (h *HolidayHandler) CreateHoliday(c *gin.Context) {
	var req models.CreateHolidayRequest
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

	holiday, err := h.service.CreateHoliday(&req, tenantID, createdBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Holiday created successfully", holiday)
}

// UpdateHoliday handles updating a holiday
func (h *HolidayHandler) UpdateHoliday(c *gin.Context) {
	idStr := c.Param("holiday_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid holiday ID", nil)
		return
	}

	var req models.UpdateHolidayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	holiday, err := h.service.UpdateHoliday(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Holiday updated successfully", holiday)
}

// DeleteHoliday handles deleting a holiday
func (h *HolidayHandler) DeleteHoliday(c *gin.Context) {
	idStr := c.Param("holiday_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid holiday ID", nil)
		return
	}

	if err := h.service.DeleteHoliday(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Holiday deleted successfully", nil)
}

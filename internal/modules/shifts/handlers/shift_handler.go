package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type ShiftHandler struct {
	service *services.ShiftService
}

func NewShiftHandler() *ShiftHandler {
	return &ShiftHandler{
		service: services.NewShiftService(),
	}
}

// toShiftResponse converts Shift model to response format
func (h *ShiftHandler) toShiftResponse(shift *models.Shift) map[string]interface{} {
	// Parse weekly_off from JSON
	var weeklyOff []string
	if shift.WeeklyOff != "" {
		json.Unmarshal([]byte(shift.WeeklyOff), &weeklyOff)
	}

	// Extract location IDs and names
	locationIDs := make([]uint, 0)
	locationNames := make([]string, 0)
	for _, loc := range shift.Locations {
		locationIDs = append(locationIDs, loc.ID)
		locationNames = append(locationNames, loc.Name)
	}

	// Format locations for detailed view
	locations := make([]map[string]interface{}, 0)
	for _, loc := range shift.Locations {
		locations = append(locations, map[string]interface{}{
			"id":   loc.ID,
			"name": loc.Name,
		})
	}

	return map[string]interface{}{
		"id":                 shift.ID,
		"shift_name":         shift.ShiftName,
		"shift_code":         shift.ShiftCode,
		"shift_type":         shift.ShiftType,
		"start_time":         shift.StartTime,
		"end_time":           shift.EndTime,
		"working_hours":      shift.WorkingHours,
		"break_duration":     shift.BreakDuration,
		"grace_minutes":      shift.GraceMinutes,
		"late_mark_after":    shift.LateMarkAfter,
		"early_going_minutes": shift.EarlyGoingMinutes,
		"half_day_hours":     shift.HalfDayHours,
		"minimum_hours":      shift.MinimumHours,
		"cross_day":          shift.CrossDay,
		"night_shift":        shift.NightShift,
		"weekly_off":         weeklyOff,
		"locations":          locationIDs,
		"location_names":     locationNames,
		"is_active":          shift.IsActive,
		"created_at":         shift.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updated_at":         shift.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"created_by":         shift.CreatedBy,
		"updated_by":         shift.UpdatedBy,
	}
}

// toShiftDetailResponse converts Shift model to detailed response format with full location objects
func (h *ShiftHandler) toShiftDetailResponse(shift *models.Shift) map[string]interface{} {
	// Parse weekly_off from JSON
	var weeklyOff []string
	if shift.WeeklyOff != "" {
		json.Unmarshal([]byte(shift.WeeklyOff), &weeklyOff)
	}

	// Format locations for detailed view
	locations := make([]map[string]interface{}, 0)
	for _, loc := range shift.Locations {
		locations = append(locations, map[string]interface{}{
			"id":   loc.ID,
			"name": loc.Name,
		})
	}

	return map[string]interface{}{
		"id":                 shift.ID,
		"shift_name":         shift.ShiftName,
		"shift_code":         shift.ShiftCode,
		"shift_type":         shift.ShiftType,
		"start_time":         shift.StartTime,
		"end_time":           shift.EndTime,
		"working_hours":      shift.WorkingHours,
		"break_duration":     shift.BreakDuration,
		"grace_minutes":      shift.GraceMinutes,
		"late_mark_after":    shift.LateMarkAfter,
		"early_going_minutes": shift.EarlyGoingMinutes,
		"half_day_hours":     shift.HalfDayHours,
		"minimum_hours":      shift.MinimumHours,
		"cross_day":          shift.CrossDay,
		"night_shift":        shift.NightShift,
		"weekly_off":         weeklyOff,
		"locations":          locations,
		"is_active":          shift.IsActive,
		"created_at":         shift.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updated_at":         shift.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"created_by":         shift.CreatedBy,
		"updated_by":         shift.UpdatedBy,
	}
}

// GetStatistics handles getting shifts and rosters statistics
func (h *ShiftHandler) GetStatistics(c *gin.Context) {
	statsService := services.NewStatisticsService()
	tenantID := middleware.GetTenantID(c)

	stats, err := statsService.GetStatistics(tenantID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve statistics", err.Error())
		return
	}

	response.Success(c, "Statistics retrieved successfully", stats)
}

// ListShifts handles listing shifts
func (h *ShiftHandler) ListShifts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)

	// Build filters
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if shiftType := c.Query("shift_type"); shiftType != "" {
		filters["shift_type"] = shiftType
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	shifts, total, err := h.service.ListShifts(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve shifts", err.Error())
		return
	}

	// Convert to response format
	shiftResponses := make([]map[string]interface{}, len(shifts))
	for i, shift := range shifts {
		shiftResponses[i] = h.toShiftResponse(&shift)
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}

	response.Success(c, "Shifts retrieved successfully", map[string]interface{}{
		"data": shiftResponses,
		"meta": map[string]interface{}{
			"page":       page,
			"per_page":   pageSize,
			"total":      total,
			"total_pages": totalPages,
		},
	})
}

// GetShift handles getting a shift by ID
func (h *ShiftHandler) GetShift(c *gin.Context) {
	idStr := c.Param("shift_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid shift ID", nil)
		return
	}

	shift, err := h.service.GetShiftByID(uint(id))
	if err != nil {
		response.NotFound(c, "Shift not found")
		return
	}

	response.Success(c, "Shift retrieved successfully", h.toShiftDetailResponse(shift))
}

// CreateShift handles creating a shift
func (h *ShiftHandler) CreateShift(c *gin.Context) {
	var req models.CreateShiftRequest
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

	shift, err := h.service.CreateShift(&req, tenantID, createdBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Shift created successfully", h.toShiftDetailResponse(shift))
}

// UpdateShift handles updating a shift
func (h *ShiftHandler) UpdateShift(c *gin.Context) {
	idStr := c.Param("shift_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid shift ID", nil)
		return
	}

	var req models.UpdateShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	shift, err := h.service.UpdateShift(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Shift updated successfully", h.toShiftDetailResponse(shift))
}

// DuplicateShift handles duplicating a shift
func (h *ShiftHandler) DuplicateShift(c *gin.Context) {
	idStr := c.Param("shift_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid shift ID", nil)
		return
	}

	var req models.DuplicateShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// JSON is optional, so we don't fail if it's empty
		_ = err
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var createdBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		createdBy = &userObj.ID
	}

	shift, err := h.service.DuplicateShift(uint(id), &req, tenantID, createdBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Shift duplicated successfully", h.toShiftDetailResponse(shift))
}

// DeleteShift handles deleting a shift
func (h *ShiftHandler) DeleteShift(c *gin.Context) {
	idStr := c.Param("shift_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid shift ID", nil)
		return
	}

	if err := h.service.DeleteShift(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Shift deleted successfully", map[string]interface{}{
		"shift_id": id,
	})
}

package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// TimesheetHandler handles timesheet API endpoints
type TimesheetHandler struct {
	timesheetService *services.TimesheetService
}

// NewTimesheetHandler creates a new timesheet handler
func NewTimesheetHandler() *TimesheetHandler {
	return &TimesheetHandler{
		timesheetService: services.NewTimesheetService(),
	}
}

// helper to extract user from context
func getUserFromContext(c *gin.Context) *userModels.User {
	user, _ := c.Get("user")
	if userObj, ok := user.(*userModels.User); ok {
		return userObj
	}
	return nil
}

// helper to extract pagination params
func getPagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}

// CreateEntry godoc
// @Summary Create a timesheet entry
// @Description Create a single timesheet entry for the authenticated employee
func (h *TimesheetHandler) CreateEntry(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	var req models.CreateTimesheetEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	entry, err := h.timesheetService.CreateEntry(user, req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Timesheet entry created successfully", entry)
}

// BulkCreateEntries godoc
// @Summary Create multiple timesheet entries
// @Description Create multiple entries for a week
func (h *TimesheetHandler) BulkCreateEntries(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	var req models.BulkCreateTimesheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	week, err := h.timesheetService.BulkCreateEntries(user, req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Timesheet entries created successfully", week)
}

// UpdateEntry godoc
// @Summary Update a timesheet entry
func (h *TimesheetHandler) UpdateEntry(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	entryID, err := strconv.ParseUint(c.Param("entry_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid entry ID", nil)
		return
	}

	var req models.CreateTimesheetEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	entry, err := h.timesheetService.UpdateEntry(user, uint(entryID), req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Timesheet entry updated successfully", entry)
}

// DeleteEntry godoc
// @Summary Delete a timesheet entry
func (h *TimesheetHandler) DeleteEntry(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	entryID, err := strconv.ParseUint(c.Param("entry_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid entry ID", nil)
		return
	}

	if err := h.timesheetService.DeleteEntry(user, uint(entryID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Timesheet entry deleted successfully", nil)
}

// GetMyTimesheets godoc
// @Summary Get my timesheets
// @Description Get the authenticated employee's timesheet weeks
func (h *TimesheetHandler) GetMyTimesheets(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	page, pageSize := getPagination(c)
	status := c.Query("status")

	weeks, total, err := h.timesheetService.GetMyTimesheets(user, status, page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response.SuccessWithMeta(c, "Timesheets retrieved successfully", weeks, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetTimesheetByID godoc
// @Summary Get a timesheet by ID
func (h *TimesheetHandler) GetTimesheetByID(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	timesheetID, err := strconv.ParseUint(c.Param("timesheet_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid timesheet ID", nil)
		return
	}

	week, err := h.timesheetService.GetTimesheetByID(user, uint(timesheetID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Timesheet retrieved successfully", week)
}

// SubmitTimesheet godoc
// @Summary Submit a timesheet for approval
func (h *TimesheetHandler) SubmitTimesheet(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	timesheetID, err := strconv.ParseUint(c.Param("timesheet_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid timesheet ID", nil)
		return
	}

	week, err := h.timesheetService.SubmitTimesheet(user, uint(timesheetID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Timesheet submitted for approval", week)
}

// RecallTimesheet godoc
// @Summary Recall a submitted timesheet
func (h *TimesheetHandler) RecallTimesheet(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	timesheetID, err := strconv.ParseUint(c.Param("timesheet_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid timesheet ID", nil)
		return
	}

	week, err := h.timesheetService.RecallTimesheet(user, uint(timesheetID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Timesheet recalled successfully", week)
}

// ApproveTimesheet godoc
// @Summary Approve or reject a timesheet (Admin/HR)
func (h *TimesheetHandler) ApproveTimesheet(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	timesheetID, err := strconv.ParseUint(c.Param("timesheet_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid timesheet ID", nil)
		return
	}

	var req models.ApproveTimesheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	week, err := h.timesheetService.ApproveTimesheet(user, uint(timesheetID), req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Timesheet processed successfully", week)
}

// GetPendingApprovals godoc
// @Summary Get timesheets pending approval (Admin/HR)
func (h *TimesheetHandler) GetPendingApprovals(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	page, pageSize := getPagination(c)

	weeks, total, err := h.timesheetService.GetPendingApprovals(user, page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response.SuccessWithMeta(c, "Pending timesheets retrieved successfully", weeks, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetTeamTimesheets godoc
// @Summary Get team timesheets (Manager view)
func (h *TimesheetHandler) GetTeamTimesheets(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	page, pageSize := getPagination(c)
	weekStart := c.Query("week_start")
	status := c.Query("status")

	weeks, total, err := h.timesheetService.GetTeamTimesheets(user, weekStart, status, page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response.SuccessWithMeta(c, "Team timesheets retrieved successfully", weeks, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// CopyLastWeek godoc
// @Summary Copy entries from last week
func (h *TimesheetHandler) CopyLastWeek(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	var req models.CopyLastWeekRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	week, err := h.timesheetService.CopyLastWeek(user, req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Entries copied from last week", week)
}

// GetEmployeeStats godoc
// @Summary Get timesheet utilization stats for the authenticated employee
func (h *TimesheetHandler) GetEmployeeStats(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		response.BadRequest(c, "start_date and end_date query parameters are required", nil)
		return
	}

	stats, err := h.timesheetService.GetEmployeeStats(user, startDate, endDate)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Timesheet statistics retrieved successfully", stats)
}

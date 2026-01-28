package handlers

import (
	"strconv"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type LeaveCalendarHandler struct {
	service *services.LeaveCalendarService
}

func NewLeaveCalendarHandler() *LeaveCalendarHandler {
	return &LeaveCalendarHandler{
		service: services.NewLeaveCalendarService(),
	}
}

// GetLeaveCalendar handles getting monthly leave calendar
func (h *LeaveCalendarHandler) GetLeaveCalendar(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	
	if yearStr == "" || monthStr == "" {
		now := time.Now()
		yearStr = strconv.Itoa(now.Year())
		monthStr = strconv.Itoa(int(now.Month()))
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.BadRequest(c, "Invalid year", nil)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		response.BadRequest(c, "Invalid month. Must be 1-12", nil)
		return
	}

	tenantID := middleware.GetTenantID(c)

	// Build filters
	filters := make(map[string]interface{})
	if departmentID := c.Query("department_id"); departmentID != "" {
		// TODO: Filter by department
	}
	if employeeID := c.Query("employee_id"); employeeID != "" {
		filters["employee_id"] = employeeID
	}

	calendar, err := h.service.GetMonthlyCalendar(year, month, tenantID, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve calendar", err.Error())
		return
	}

	response.Success(c, "Leave calendar retrieved successfully", calendar)
}

// GetEmployeesOnLeaveToday handles getting employees on leave today
func (h *LeaveCalendarHandler) GetEmployeesOnLeaveToday(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	// Build filters
	filters := make(map[string]interface{})
	if departmentID := c.Query("department_id"); departmentID != "" {
		// TODO: Filter by department
	}
	if locationID := c.Query("location_id"); locationID != "" {
		// TODO: Filter by location
	}

	employees, err := h.service.GetEmployeesOnLeaveToday(tenantID, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve employees on leave", err.Error())
		return
	}

	response.Success(c, "Today's leave list retrieved successfully", map[string]interface{}{
		"date":            time.Now().Format("2006-01-02"),
		"total_on_leave": len(employees),
		"employees":      employees,
	})
}

// GetWeeklyCalendar handles getting weekly leave calendar
func (h *LeaveCalendarHandler) GetWeeklyCalendar(c *gin.Context) {
	yearStr := c.Query("year")
	weekStr := c.Query("week")
	
	if yearStr == "" || weekStr == "" {
		now := time.Now()
		yearStr = strconv.Itoa(now.Year())
		_, week := now.ISOWeek()
		weekStr = strconv.Itoa(week)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.BadRequest(c, "Invalid year", nil)
		return
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 || week > 52 {
		response.BadRequest(c, "Invalid week. Must be 1-52", nil)
		return
	}

	tenantID := middleware.GetTenantID(c)

	// Build filters
	filters := make(map[string]interface{})
	if departmentID := c.Query("department_id"); departmentID != "" {
		// TODO: Filter by department
	}

	calendar, err := h.service.GetWeeklyCalendar(year, week, tenantID, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve weekly calendar", err.Error())
		return
	}

	response.Success(c, "Week calendar retrieved successfully", calendar)
}

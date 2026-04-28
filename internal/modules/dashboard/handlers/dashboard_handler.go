package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/dashboard/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard HTTP requests
type DashboardHandler struct {
	dashboardService *services.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		dashboardService: services.NewDashboardService(),
	}
}

// GetStatistics returns dashboard statistics based on user role
// @Summary Get dashboard statistics
// @Description Returns KPI statistics based on user role. Admins get organization-wide stats, employees get personal stats.
// @Tags Dashboard
// @Produce json
// @Param view query string false "View type: admin or employee (defaults based on user role)"
// @Param date query string false "Date for statistics (default: today, format: YYYY-MM-DD)"
// @Success 200 {object} response.APIResponse
// @Router/dashboard/statistics [get]
func (h *DashboardHandler) GetStatistics(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	uid := userID.(uint)

	// Parse date parameter
	dateStr := c.Query("date")
	date := time.Now()
	if dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = parsed
		}
	}

	// Determine view based on user role or parameter
	view := strings.ToLower(c.Query("view"))
	userRoles, _ := c.Get("user_roles")
	roles, ok := userRoles.([]string)
	isAdmin := false
	if ok {
		for _, role := range roles {
			if role == "admin" || role == "hr" {
				isAdmin = true
				break
			}
		}
	}

	// Default view based on role
	if view == "" {
		if isAdmin {
			view = "admin"
		} else {
			view = "employee"
		}
	}

	// Admin view
	if view == "admin" {
		if !isAdmin {
			response.Forbidden(c, "You do not have permission to view admin statistics")
			return
		}

		stats, err := h.dashboardService.GetAdminStatistics(tenantID, date)
		if err != nil {
			response.InternalServerError(c, "Failed to retrieve statistics", err.Error())
			return
		}
		response.Success(c, "Dashboard statistics retrieved successfully", stats)
		return
	}

	// Employee view
	stats, err := h.dashboardService.GetEmployeeStatistics(uid, date)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve statistics", err.Error())
		return
	}
	response.Success(c, "Dashboard statistics retrieved successfully", stats)
}

// GetAnnouncements returns company announcements
// @Summary Get announcements
// @Description Returns company announcements for the dashboard
// @Tags Dashboard
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 5, max: 20)"
// @Param type query string false "Filter by type: info, warning, success, urgent"
// @Param active_only query bool false "Only show active announcements (default: true)"
// @Success 200 {object} response.APIResponse
// @Router/dashboard/announcements [get]
func (h *DashboardHandler) GetAnnouncements(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	uid := userID.(uint)

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 5
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 20 {
			pageSize = parsed
		}
	}

	announcementType := strings.TrimSpace(c.Query("type"))

	activeOnly := true
	if ao := c.Query("active_only"); ao == "false" || ao == "0" {
		activeOnly = false
	}

	result, err := h.dashboardService.GetAnnouncements(tenantID, uid, page, pageSize, announcementType, activeOnly)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve announcements", err.Error())
		return
	}

	response.Success(c, "Announcements retrieved successfully", result)
}

// MarkAnnouncementAsRead marks an announcement as read
// @Summary Mark announcement as read
// @Description Marks an announcement as read for the current user
// @Tags Dashboard
// @Produce json
// @Param id path int true "Announcement ID"
// @Success 200 {object} response.APIResponse
// @Router/dashboard/announcements/{id}/read [post]
func (h *DashboardHandler) MarkAnnouncementAsRead(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	uid := userID.(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid announcement ID", nil)
		return
	}

	readAt, err := h.dashboardService.MarkAnnouncementAsRead(uint(id), uid, tenantID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Announcement marked as read", gin.H{
		"announcementId": id,
		"readAt":         readAt,
	})
}

// GetEvents returns upcoming events (birthdays, anniversaries)
// @Summary Get upcoming events
// @Description Returns upcoming birthdays, work anniversaries, and other employee events
// @Tags Dashboard
// @Produce json
// @Param days_ahead query int false "Number of days to look ahead (default: 30, max: 90)"
// @Param event_type query string false "Filter: birthday, anniversary, probation_end, all"
// @Param department_id query int false "Filter by department"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10)"
// @Success 200 {object} response.APIResponse
// @Router/dashboard/events [get]
func (h *DashboardHandler) GetEvents(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	daysAhead := 30
	if d := c.Query("days_ahead"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 90 {
			daysAhead = parsed
		}
	}

	eventType := strings.TrimSpace(c.Query("event_type"))

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 50 {
			pageSize = parsed
		}
	}

	result, err := h.dashboardService.GetUpcomingEvents(tenantID, daysAhead, eventType, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve events", err.Error())
		return
	}

	response.Success(c, "Upcoming events retrieved successfully", result)
}

// GetMyActivity returns the logged-in employee's recent activities
// @Summary Get my recent activity
// @Description Returns the logged-in employee's recent activities
// @Tags Dashboard
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10, max: 50)"
// @Param activity_type query string false "Filter: leave, timesheet, payslip, request, all"
// @Param days query int false "Activities from last N days (default: 30)"
// @Success 200 {object} response.APIResponse
// @Router/dashboard/my-activity [get]
func (h *DashboardHandler) GetMyActivity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	uid := userID.(uint)

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 50 {
			pageSize = parsed
		}
	}

	activityType := strings.TrimSpace(c.Query("activity_type"))

	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 365 {
			days = parsed
		}
	}

	result, err := h.dashboardService.GetMyActivity(uid, days, page, pageSize, activityType)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve activities", err.Error())
		return
	}

	response.Success(c, "Recent activities retrieved successfully", result)
}

// GetQuickActions returns personalized quick actions based on user role
// @Summary Get quick actions
// @Description Returns personalized quick action links based on user role and permissions
// @Tags Dashboard
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router/dashboard/quick-actions [get]
func (h *DashboardHandler) GetQuickActions(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	uid := userID.(uint)

	// Check if user is admin
	userRoles, _ := c.Get("user_roles")
	roles, ok := userRoles.([]string)
	isAdmin := false
	if ok {
		for _, role := range roles {
			if role == "admin" || role == "hr" {
				isAdmin = true
				break
			}
		}
	}

	// Get pending approvals count for badge
	pendingApprovals := 0
	if isAdmin {
		approvals, _ := h.dashboardService.GetPendingApprovals(tenantID, "", 1, 1)
		if approvals != nil {
			pendingApprovals = approvals.Summary.Total
		}
	}

	result := h.dashboardService.GetQuickActions(uid, isAdmin, pendingApprovals)
	response.Success(c, "Quick actions retrieved successfully", result)
}

// GetPendingApprovals returns pending approvals for managers/admins
// @Summary Get pending approvals
// @Description Returns detailed breakdown of pending approvals for managers/admins
// @Tags Dashboard
// @Produce json
// @Param type query string false "Filter: leave, timesheet, expense, all"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10)"
// @Success 200 {object} response.APIResponse
// @Router/dashboard/pending-approvals [get]
func (h *DashboardHandler) GetPendingApprovals(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	approvalType := strings.TrimSpace(c.Query("type"))

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 50 {
			pageSize = parsed
		}
	}

	result, err := h.dashboardService.GetPendingApprovals(tenantID, approvalType, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve pending approvals", err.Error())
		return
	}

	response.Success(c, "Pending approvals retrieved successfully", result)
}

// GetApprovalCounts returns counts of pending approvals for dashboard badges
// @Summary Get approval counts
// @Description Returns counts of pending approvals for dashboard badges (simplified version of pending-approvals)
// @Tags Dashboard
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router/dashboard/approval-counts [get]
func (h *DashboardHandler) GetApprovalCounts(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	// Get user role from context
	userRoleInterface, exists := c.Get("user_role")
	if !exists {
		response.Unauthorized(c, "User role not found")
		return
	}
	userRole := userRoleInterface.(string)

	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found")
		return
	}
	userID := userIDInterface.(uint)

	// Get approval counts
	counts, err := h.dashboardService.GetApprovalCounts(tenantID, userRole, userID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve approval counts", err.Error())
		return
	}

	response.Success(c, "Approval counts retrieved successfully", counts)
}

// ────────────────────────── Employee Dashboard ──────────────────────────

// GetEmployeeStatistics returns personal dashboard statistics for the logged-in employee
// @Summary Get employee dashboard statistics
// @Description Returns personal KPI stats: leave balance, hours this week, pending requests, next payday, attendance
// @Tags Employee Dashboard
// @Produce json
// @Param date query string false "Date for statistics (default: today, format: YYYY-MM-DD)"
// @Success 200 {object} response.APIResponse
// @Router/dashboard/employee/statistics [get]
func (h *DashboardHandler) GetEmployeeStatistics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	uid := userID.(uint)

	// Parse date
	date := time.Now()
	if dateStr := c.Query("date"); dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = parsed
		}
	}

	stats, err := h.dashboardService.GetEmployeeStatistics(uid, date)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve employee statistics", err.Error())
		return
	}

	response.Success(c, "Employee dashboard statistics retrieved successfully", stats)
}

// GetEmployeeActivity returns recent activity feed for the logged-in employee
// @Summary Get employee recent activity
// @Description Returns the employee's recent activities: leave requests, timesheets, payslips, service requests
// @Tags Employee Dashboard
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10, max: 50)"
// @Param activity_type query string false "Filter: leave, timesheet, payslip, request, all"
// @Param days query int false "Activities from last N days (default: 30, max: 365)"
// @Success 200 {object} response.APIResponse
// @Router/dashboard/employee/my-activity [get]
func (h *DashboardHandler) GetEmployeeActivity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	uid := userID.(uint)

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 50 {
			pageSize = parsed
		}
	}

	activityType := strings.TrimSpace(c.Query("activity_type"))

	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 365 {
			days = parsed
		}
	}

	result, err := h.dashboardService.GetMyActivity(uid, days, page, pageSize, activityType)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve employee activity", err.Error())
		return
	}

	response.Success(c, "Employee recent activity retrieved successfully", result)
}

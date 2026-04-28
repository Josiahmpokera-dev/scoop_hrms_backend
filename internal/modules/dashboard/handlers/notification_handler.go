package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	attendanceModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	dashboardRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/dashboard/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NotificationHandler handles notification count endpoints
type NotificationHandler struct {
	db               *gorm.DB
	dashboardRepo    *dashboardRepos.DashboardRepository
	announcementRepo *dashboardRepos.AnnouncementRepository
	employeeRepo     *employeeRepos.EmployeeRepository
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		db:               database.GetDB(),
		dashboardRepo:    dashboardRepos.NewDashboardRepository(),
		announcementRepo: dashboardRepos.NewAnnouncementRepository(),
		employeeRepo:     employeeRepos.NewEmployeeRepository(),
	}
}

// GetNotificationCount returns aggregated notification counts for the authenticated user.
// Admin/HR users see pending approvals counts; employees see their own pending items.
// @Summary Get notification counts
// @Description Returns unread/pending notification counts aggregated from all modules
// @Tags Notifications
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router/notifications/count [get]
func (h *NotificationHandler) GetNotificationCount(c *gin.Context) {
	user, _ := c.Get("user")
	userObj, ok := user.(*userModels.User)
	if !ok {
		response.Unauthorized(c, "Authentication required")
		return
	}

	tenantID := middleware.GetTenantID(c)
	isAdmin := userObj.Role == userModels.RoleAdmin || userObj.Role == userModels.RoleHR

	// Counts
	var pendingLeaveApprovals int64
	var pendingTimesheetApprovals int64
	var pendingOvertimeApprovals int64
	var unreadAnnouncements int64
	var myPendingLeave int64
	var myPendingTimesheet int64
	var myPendingOvertime int64

	// Unread announcements (for all users)
	unreadAnnouncements, _ = h.announcementRepo.GetUnreadCount(tenantID, userObj.ID)

	if isAdmin {
		// Admin/HR: count pending approvals across the organization
		pendingLeaveApprovals, _ = h.dashboardRepo.GetPendingLeaveRequests(tenantID)

		// Pending timesheet approvals (submitted status)
		h.db.Model(&attendanceModels.TimesheetWeek{}).
			Where("status = ?", attendanceModels.TimesheetStatusSubmitted).
			Count(&pendingTimesheetApprovals)

		// Pending overtime approvals
		h.db.Model(&attendanceModels.OvertimeRequest{}).
			Where("status = ?", attendanceModels.OTStatusPending).
			Count(&pendingOvertimeApprovals)
	}

	// Employee's own pending items
	employee, err := h.employeeRepo.FindByUserID(userObj.ID)
	if err == nil && employee != nil {
		myPendingLeave, _ = h.dashboardRepo.GetEmployeePendingLeaveRequests(employee.EmployeeID)

		// My pending timesheet submissions (draft timesheets that need to be submitted)
		h.db.Model(&attendanceModels.TimesheetWeek{}).
			Where("employee_id = ? AND status = ?", employee.ID, attendanceModels.TimesheetStatusDraft).
			Count(&myPendingTimesheet)

		// My pending overtime requests
		h.db.Model(&attendanceModels.OvertimeRequest{}).
			Where("employee_id = ? AND status = ?", employee.ID, attendanceModels.OTStatusPending).
			Count(&myPendingOvertime)
	}

	// Calculate totals
	totalApprovals := pendingLeaveApprovals + pendingTimesheetApprovals + pendingOvertimeApprovals
	totalMyPending := myPendingLeave + myPendingTimesheet + myPendingOvertime
	totalCount := unreadAnnouncements + totalMyPending
	if isAdmin {
		totalCount += totalApprovals
	}

	result := gin.H{
		"total": totalCount,
		"unread_announcements": unreadAnnouncements,
		"my_pending": gin.H{
			"total":             totalMyPending,
			"leave_requests":    myPendingLeave,
			"draft_timesheets":  myPendingTimesheet,
			"overtime_requests": myPendingOvertime,
		},
	}

	// Only include approval counts for admin/HR
	if isAdmin {
		result["pending_approvals"] = gin.H{
			"total":              totalApprovals,
			"leave_requests":     pendingLeaveApprovals,
			"timesheet_approvals": pendingTimesheetApprovals,
			"overtime_requests":  pendingOvertimeApprovals,
		}
	}

	response.Success(c, "Notification counts retrieved successfully", result)
}

// GetNotifications returns a list of recent notifications for the user (unread announcements + pending items)
// @Summary Get notifications list
// @Description Returns recent notification items
// @Tags Notifications
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router/notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	user, _ := c.Get("user")
	userObj, ok := user.(*userModels.User)
	if !ok {
		response.Unauthorized(c, "Authentication required")
		return
	}

	tenantID := middleware.GetTenantID(c)
	isAdmin := userObj.Role == userModels.RoleAdmin || userObj.Role == userModels.RoleHR

	var notifications []gin.H

	// Unread announcements
	unreadCount, _ := h.announcementRepo.GetUnreadCount(tenantID, userObj.ID)
	if unreadCount > 0 {
		notifications = append(notifications, gin.H{
			"type":    "announcement",
			"title":   "Unread Announcements",
			"message": formatCount(unreadCount, "unread announcement"),
			"count":   unreadCount,
			"icon":    "megaphone",
			"color":   "blue",
			"href":    "/dashboard",
		})
	}

	// Employee's own pending items
	employee, err := h.employeeRepo.FindByUserID(userObj.ID)
	if err == nil && employee != nil {
		pendingLeave, _ := h.dashboardRepo.GetEmployeePendingLeaveRequests(employee.EmployeeID)
		if pendingLeave > 0 {
			notifications = append(notifications, gin.H{
				"type":    "leave",
				"title":   "Pending Leave Requests",
				"message": formatCount(pendingLeave, "leave request awaiting approval"),
				"count":   pendingLeave,
				"icon":    "calendar",
				"color":   "orange",
				"href":    "/leave/requests",
			})
		}

		var draftTimesheets int64
		h.db.Model(&attendanceModels.TimesheetWeek{}).
			Where("employee_id = ? AND status = ?", employee.ID, attendanceModels.TimesheetStatusDraft).
			Count(&draftTimesheets)
		if draftTimesheets > 0 {
			notifications = append(notifications, gin.H{
				"type":    "timesheet",
				"title":   "Draft Timesheets",
				"message": formatCount(draftTimesheets, "timesheet to submit"),
				"count":   draftTimesheets,
				"icon":    "clock",
				"color":   "yellow",
				"href":    "/attendance/timesheets",
			})
		}

		var pendingOT int64
		h.db.Model(&attendanceModels.OvertimeRequest{}).
			Where("employee_id = ? AND status = ?", employee.ID, attendanceModels.OTStatusPending).
			Count(&pendingOT)
		if pendingOT > 0 {
			notifications = append(notifications, gin.H{
				"type":    "overtime",
				"title":   "Pending Overtime Requests",
				"message": formatCount(pendingOT, "overtime request pending"),
				"count":   pendingOT,
				"icon":    "zap",
				"color":   "orange",
				"href":    "/attendance/overtime",
			})
		}
	}

	// Admin/HR approval notifications
	if isAdmin {
		pendingLeaveApprovals, _ := h.dashboardRepo.GetPendingLeaveRequests(tenantID)
		if pendingLeaveApprovals > 0 {
			notifications = append(notifications, gin.H{
				"type":    "approval_leave",
				"title":   "Leave Approvals Needed",
				"message": formatCount(pendingLeaveApprovals, "leave request needs your approval"),
				"count":   pendingLeaveApprovals,
				"icon":    "check-circle",
				"color":   "red",
				"href":    "/leave/requests",
			})
		}

		var pendingTSApprovals int64
		h.db.Model(&attendanceModels.TimesheetWeek{}).
			Where("status = ?", attendanceModels.TimesheetStatusSubmitted).
			Count(&pendingTSApprovals)
		if pendingTSApprovals > 0 {
			notifications = append(notifications, gin.H{
				"type":    "approval_timesheet",
				"title":   "Timesheet Approvals Needed",
				"message": formatCount(pendingTSApprovals, "timesheet needs your approval"),
				"count":   pendingTSApprovals,
				"icon":    "clipboard-check",
				"color":   "red",
				"href":    "/attendance/timesheets",
			})
		}

		var pendingOTApprovals int64
		h.db.Model(&attendanceModels.OvertimeRequest{}).
			Where("status = ?", attendanceModels.OTStatusPending).
			Count(&pendingOTApprovals)
		if pendingOTApprovals > 0 {
			notifications = append(notifications, gin.H{
				"type":    "approval_overtime",
				"title":   "Overtime Approvals Needed",
				"message": formatCount(pendingOTApprovals, "overtime request needs your approval"),
				"count":   pendingOTApprovals,
				"icon":    "alert-circle",
				"color":   "red",
				"href":    "/attendance/overtime",
			})
		}
	}

	if notifications == nil {
		notifications = []gin.H{}
	}

	response.Success(c, "Notifications retrieved successfully", gin.H{
		"notifications": notifications,
		"total":         len(notifications),
	})
}

// formatCount returns a human-readable count string
func formatCount(count int64, singular string) string {
	if count == 1 {
		return "1 " + singular
	}
	// Simple plural: just add "s"
	return strconv.FormatInt(count, 10) + " " + singular + "s"
}

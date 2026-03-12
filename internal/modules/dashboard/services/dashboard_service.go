package services

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/dashboard/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/dashboard/repositories"
)

// DashboardService handles dashboard business logic
type DashboardService struct {
	dashboardRepo    *repositories.DashboardRepository
	announcementRepo *repositories.AnnouncementRepository
}

// NewDashboardService creates a new dashboard service
func NewDashboardService() *DashboardService {
	return &DashboardService{
		dashboardRepo:    repositories.NewDashboardRepository(),
		announcementRepo: repositories.NewAnnouncementRepository(),
	}
}

// ========================================
// Statistics
// ========================================

// GetAdminStatistics returns statistics for admin dashboard
func (s *DashboardService) GetAdminStatistics(tenantID *uint, date time.Time) (*models.AdminStatisticsResponse, error) {
	// Get headcount
	currentHeadcount, err := s.dashboardRepo.GetTotalHeadcount(tenantID)
	if err != nil {
		return nil, err
	}

	lastMonthHeadcount, _ := s.dashboardRepo.GetHeadcountLastMonth(tenantID)
	change := int(currentHeadcount - lastMonthHeadcount)
	var changePercent float64
	var trend models.TrendDirection
	if lastMonthHeadcount > 0 {
		changePercent = float64(change) / float64(lastMonthHeadcount) * 100
	}
	if change > 0 {
		trend = models.TrendUp
	} else if change < 0 {
		trend = models.TrendDown
	} else {
		trend = models.TrendStable
	}

	// Get employees on leave today
	onLeave, _ := s.dashboardRepo.GetEmployeesOnLeaveToday(tenantID, date)
	// Present = Total - OnLeave (simplified; in real app, calculate from attendance)
	present := currentHeadcount - onLeave
	absent := int64(0) // Would come from attendance system
	var presentPercent float64
	if currentHeadcount > 0 {
		presentPercent = float64(present) / float64(currentHeadcount) * 100
	}

	// Get pending approvals
	pendingLeave, _ := s.dashboardRepo.GetPendingLeaveRequests(tenantID)
	// TODO: Add timesheet and expense claims counts when those modules are available

	// Get summary stats
	newHires, _ := s.dashboardRepo.GetNewHiresThisMonth(tenantID)
	terminations, _ := s.dashboardRepo.GetTerminationsThisMonth(tenantID)

	// Payroll status (simplified - would come from payroll module)
	now := time.Now()
	currentPeriod := now.Format("January 2006")
	processingDate := time.Date(now.Year(), now.Month(), 25, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	return &models.AdminStatisticsResponse{
		View: "admin",
		Date: date.Format("2006-01-02"),
		Stats: models.AdminStats{
			TotalHeadcount: models.HeadcountStats{
				Value:         int(currentHeadcount),
				Change:        change,
				ChangePercent: math.Round(changePercent*100) / 100,
				Trend:         trend,
			},
			PresentToday: models.PresentTodayStats{
				Value:      int(present),
				Total:      int(currentHeadcount),
				Percentage: math.Round(presentPercent*100) / 100,
				Absent:     int(absent),
				OnLeave:    int(onLeave),
			},
			PendingApprovals: models.PendingApprovalsStats{
				Value: int(pendingLeave),
				Breakdown: models.PendingApprovalsBreakdown{
					LeaveRequests:        int(pendingLeave),
					TimesheetCorrections: 0,
					ExpenseClaims:        0,
				},
			},
			PayrollStatus: models.PayrollStatusStats{
				Status:         "Active",
				CurrentPeriod:  currentPeriod,
				ProcessingDate: processingDate,
				TotalAmount:    0,
				Currency:       "TZS",
			},
		},
		Summary: models.AdminSummary{
			NewHiresThisMonth:     int(newHires),
			TerminationsThisMonth: int(terminations),
			OpenPositions:         0, // Would come from recruitment module
			AvgAttendanceRate:     presentPercent,
		},
	}, nil
}

// GetEmployeeStatistics returns statistics for employee dashboard
func (s *DashboardService) GetEmployeeStatistics(userID uint, date time.Time) (*models.EmployeeStatisticsResponse, error) {
	// Get employee by user ID
	employee, err := s.dashboardRepo.GetEmployeeByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("employee not found for user")
	}

	year := date.Year()

	// Get leave balance
	balances, _ := s.dashboardRepo.GetEmployeeLeaveBalance(employee.EmployeeID, year)
	totalBalance := 0.0
	annualBalance := balances["ANNUAL"]
	sickBalance := balances["SICK"]
	personalBalance := balances["PERSONAL"]
	for _, b := range balances {
		totalBalance += b
	}

	// Get total leave used
	usedThisYear, _ := s.dashboardRepo.GetEmployeeTotalLeaveUsed(employee.EmployeeID, year)

	// Get pending leave requests
	pendingLeave, _ := s.dashboardRepo.GetEmployeePendingLeaveRequests(employee.EmployeeID)

	// Calculate next payday (last day of month typically)
	nextPayday := time.Date(year, date.Month()+1, 0, 0, 0, 0, 0, date.Location())
	if date.Day() > nextPayday.Day() {
		nextPayday = time.Date(year, date.Month()+2, 0, 0, 0, 0, 0, date.Location())
	}
	daysRemaining := int(nextPayday.Sub(date).Hours() / 24)

	employeeName := fmt.Sprintf("%s %s", employee.FirstName, employee.LastName)

	return &models.EmployeeStatisticsResponse{
		View:         "employee",
		Date:         date.Format("2006-01-02"),
		EmployeeID:   employee.ID,
		EmployeeName: employeeName,
		Stats: models.EmployeeStats{
			LeaveBalance: models.LeaveBalanceStats{
				Value: totalBalance,
				Unit:  "days",
				Breakdown: models.LeaveBalanceBreakdown{
					Annual:   annualBalance,
					Sick:     sickBalance,
					Personal: personalBalance,
				},
				UsedThisYear: usedThisYear,
			},
			HoursThisWeek: models.HoursThisWeekStats{
				Value:    0, // Would come from attendance/timesheet
				Target:   40,
				Overtime: 0,
				Status:   "on_track",
			},
			PendingRequests: models.PendingRequestsStats{
				Value: int(pendingLeave),
				Breakdown: models.PendingRequestsBreakdown{
					LeaveRequests:  int(pendingLeave),
					LetterRequests: 0,
				},
			},
			NextPayday: models.NextPaydayStats{
				Date:            nextPayday.Format("2006-01-02"),
				DaysRemaining:   daysRemaining,
				EstimatedNetPay: 0, // Would come from payroll
				Currency:        "TZS",
			},
		},
		Attendance: models.AttendanceInfo{
			TodayStatus:      "Not Clocked In", // Would come from attendance
			CheckInTime:      nil,
			CheckOutTime:     nil,
			WorkedHoursToday: 0,
		},
	}, nil
}

// ========================================
// Announcements
// ========================================

// GetAnnouncements returns paginated announcements
func (s *DashboardService) GetAnnouncements(tenantID *uint, userID uint, page, pageSize int, announcementType string, activeOnly bool) (*models.AnnouncementsListResponse, error) {
	announcements, total, err := s.announcementRepo.GetAnnouncements(tenantID, page, pageSize, announcementType, activeOnly)
	if err != nil {
		return nil, err
	}

	unreadCount, _ := s.announcementRepo.GetUnreadCount(tenantID, userID)

	var response []models.AnnouncementResponse
	for _, a := range announcements {
		creatorName, _ := s.dashboardRepo.GetUserName(a.CreatedBy)

		var attachments []models.AnnouncementAttachmentResponse
		for _, att := range a.Attachments {
			attachments = append(attachments, models.AnnouncementAttachmentResponse{
				ID:   att.ID,
				Name: att.Name,
				URL:  att.URL,
				Size: att.Size,
			})
		}

		isRead := s.announcementRepo.IsRead(a.ID, userID)

		response = append(response, models.AnnouncementResponse{
			ID:          a.ID,
			Title:       a.Title,
			Description: a.Description,
			Type:        string(a.Type),
			Priority:    string(a.Priority),
			CreatedAt:   a.CreatedAt,
			CreatedBy: models.AnnouncementCreator{
				ID:   a.CreatedBy,
				Name: creatorName,
			},
			ExpiresAt:    a.ExpiresAt,
			IsRead:       isRead,
			Attachments:  attachments,
			RelativeTime: getRelativeTime(a.CreatedAt),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.AnnouncementsListResponse{
		Announcements: response,
		Meta: models.AnnouncementsMeta{
			Page:        page,
			PageSize:    pageSize,
			Total:       total,
			TotalPages:  totalPages,
			UnreadCount: unreadCount,
		},
	}, nil
}

// MarkAnnouncementAsRead marks an announcement as read
func (s *DashboardService) MarkAnnouncementAsRead(announcementID, userID uint, tenantID *uint) (*time.Time, error) {
	// Verify announcement exists
	_, err := s.announcementRepo.GetByID(announcementID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("announcement not found")
	}

	if err := s.announcementRepo.MarkAsRead(announcementID, userID); err != nil {
		return nil, err
	}

	readAt := s.announcementRepo.GetReadAt(announcementID, userID)
	return readAt, nil
}

// ========================================
// Events
// ========================================

// GetUpcomingEvents returns upcoming events (birthdays, anniversaries)
func (s *DashboardService) GetUpcomingEvents(tenantID *uint, daysAhead int, eventType string, page, pageSize int) (*models.EventsResponse, error) {
	var events []models.UpcomingEvent
	var total int64

	now := time.Now()

	// Get birthdays
	if eventType == "" || eventType == "all" || eventType == "birthday" {
		birthdays, count, _ := s.dashboardRepo.GetUpcomingBirthdays(tenantID, daysAhead, page, pageSize)
		total += count
		for _, b := range birthdays {
			age := 0
			if b.DateOfBirth != nil {
				age = now.Year() - b.DateOfBirth.Year()
			}
			events = append(events, models.UpcomingEvent{
				ID:            b.ID,
				EmployeeID:    b.ID,
				EmployeeName:  fmt.Sprintf("%s %s", b.FirstName, b.LastName),
				EmployeePhoto: b.PhotoURL,
				Department:    b.Department,
				EventType:     "birthday",
				EventDate:     b.EventDate.Format("2006-01-02"),
				DisplayDate:   b.EventDate.Format("Jan 2"),
				DaysUntil:     calculateDaysUntil(b.EventDate, now),
				Age:           &age,
				Initials:      getInitials(b.FirstName, b.LastName),
			})
		}
	}

	// Get anniversaries
	if eventType == "" || eventType == "all" || eventType == "anniversary" {
		anniversaries, count, _ := s.dashboardRepo.GetUpcomingAnniversaries(tenantID, daysAhead, page, pageSize)
		total += count
		for _, a := range anniversaries {
			yearsOfService := 0
			if a.HireDate != nil {
				yearsOfService = now.Year() - a.HireDate.Year()
			}
			events = append(events, models.UpcomingEvent{
				ID:             a.ID,
				EmployeeID:     a.ID,
				EmployeeName:   fmt.Sprintf("%s %s", a.FirstName, a.LastName),
				EmployeePhoto:  a.PhotoURL,
				Department:     a.Department,
				EventType:      "anniversary",
				EventDate:      a.EventDate.Format("2006-01-02"),
				DisplayDate:    a.EventDate.Format("Jan 2"),
				DaysUntil:      calculateDaysUntil(a.EventDate, now),
				YearsOfService: &yearsOfService,
				Initials:       getInitials(a.FirstName, a.LastName),
			})
		}
	}

	// Get summary
	birthdaysThisMonth, _ := s.dashboardRepo.GetBirthdaysThisMonth(tenantID)
	anniversariesThisMonth, _ := s.dashboardRepo.GetAnniversariesThisMonth(tenantID)

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.EventsResponse{
		Events: events,
		Summary: models.EventsSummary{
			BirthdaysThisMonth:     int(birthdaysThisMonth),
			AnniversariesThisMonth: int(anniversariesThisMonth),
			ProbationEndsThisMonth: 0, // TODO: Implement when probation tracking is available
		},
		Meta: models.EventsMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// ========================================
// Recent Activity (Employee)
// ========================================

// GetMyActivity returns recent activities for an employee
func (s *DashboardService) GetMyActivity(userID uint, days, page, pageSize int, activityType string) (*models.MyActivityResponse, error) {
	employee, err := s.dashboardRepo.GetEmployeeByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("employee not found for user")
	}

	var activities []models.RecentActivity
	var total int64

	// Get leave requests
	if activityType == "" || activityType == "all" || activityType == "leave" {
		requests, count, _ := s.dashboardRepo.GetEmployeeRecentLeaveRequests(employee.EmployeeID, days, page, pageSize)
		total += count
		for _, r := range requests {
			leaveTypeName := ""
			if r.LeaveType != nil {
				leaveTypeName = r.LeaveType.Name
			}

			status := r.Status
			statusLabel := strings.Title(status)
			color := "blue"
			if status == "approved" {
				color = "green"
				statusLabel = "Approved"
			} else if status == "rejected" {
				color = "red"
				statusLabel = "Rejected"
			} else if status == "pending" {
				color = "orange"
				statusLabel = "Pending approval"
			}

			activities = append(activities, models.RecentActivity{
				ID:           r.ID,
				Type:         "leave_request",
				Title:        "Leave Request",
				Description:  fmt.Sprintf("%s for %s - %s", leaveTypeName, r.FromDate.Format("Jan 2"), r.ToDate.Format("Jan 2")),
				Status:       status,
				StatusLabel:  statusLabel,
				Icon:         "calendar",
				Color:        color,
				CreatedAt:    r.CreatedAt,
				RelativeTime: getRelativeTime(r.CreatedAt),
				Metadata: map[string]interface{}{
					"leaveType": leaveTypeName,
					"startDate": r.FromDate.Format("2006-01-02"),
					"endDate":   r.ToDate.Format("2006-01-02"),
					"days":      r.TotalDays,
				},
				ActionURL: fmt.Sprintf("/leave/requests/%d", r.ID),
			})
		}
	}

	// Count pending requests
	pendingLeave, _ := s.dashboardRepo.GetEmployeePendingLeaveRequests(employee.EmployeeID)

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.MyActivityResponse{
		Activities: activities,
		Summary: models.ActivitySummary{
			PendingRequests:  int(pendingLeave),
			RecentApprovals:  0, // TODO: Calculate from approved requests
			NewNotifications: 0, // TODO: Implement notifications
		},
		Meta: models.ActivityMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// ========================================
// Quick Actions
// ========================================

// GetQuickActions returns personalized quick actions based on user role
func (s *DashboardService) GetQuickActions(userID uint, isAdmin bool, pendingApprovals int) *models.QuickActionsResponse {
	if isAdmin {
		return &models.QuickActionsResponse{
			View: "admin",
			Actions: []models.QuickAction{
				{
					ID:          "add_employee",
					Title:       "Add Employee",
					Description: "Register new employee",
					Icon:        "user-plus",
					Color:       "blue",
					Href:        "/employees/add",
					Badge:       nil,
					Enabled:     true,
				},
				{
					ID:          "view_reports",
					Title:       "View Reports",
					Description: "Generate reports",
					Icon:        "chart",
					Color:       "green",
					Href:        "/reports/analytics",
					Badge:       nil,
					Enabled:     true,
				},
				{
					ID:          "approve_requests",
					Title:       "Approve Requests",
					Description: "Pending approvals",
					Icon:        "check-circle",
					Color:       "orange",
					Href:        "/leave/requests",
					Badge: func() *models.QuickActionBadge {
						if pendingApprovals > 0 {
							return &models.QuickActionBadge{
								Count: &pendingApprovals,
								Type:  "warning",
							}
						}
						return nil
					}(),
					Enabled: true,
				},
				{
					ID:          "manage_payroll",
					Title:       "Manage Payroll",
					Description: "Process payroll",
					Icon:        "currency",
					Color:       "purple",
					Href:        "/payroll/dashboard",
					Badge:       nil,
					Enabled:     true,
				},
			},
		}
	}

	// Employee view
	return &models.QuickActionsResponse{
		View: "employee",
		Actions: []models.QuickAction{
			{
				ID:          "apply_leave",
				Title:       "Apply Leave",
				Description: "Submit leave request",
				Icon:        "calendar",
				Color:       "blue",
				Href:        "/leave/apply",
				Badge:       nil,
				Enabled:     true,
			},
			{
				ID:          "my_payslips",
				Title:       "My Payslips",
				Description: "View payslips",
				Icon:        "document",
				Color:       "orange",
				Href:        "/self-service/my-payslips",
				Badge:       nil,
				Enabled:     true,
			},
			{
				ID:          "requests_letters",
				Title:       "Requests & Letters",
				Description: "Submit requests",
				Icon:        "mail",
				Color:       "green",
				Href:        "/self-service/requests",
				Badge:       nil,
				Enabled:     true,
			},
			{
				ID:          "people_directory",
				Title:       "People Directory",
				Description: "Find colleagues",
				Icon:        "users",
				Color:       "purple",
				Href:        "/self-service/directory",
				Badge:       nil,
				Enabled:     true,
			},
		},
	}
}

// ========================================
// Pending Approvals (Admin)
// ========================================

// GetPendingApprovals returns pending approvals for admin
func (s *DashboardService) GetPendingApprovals(tenantID *uint, approvalType string, page, pageSize int) (*models.PendingApprovalsResponse, error) {
	var approvals []models.PendingApproval
	var total int64

	// Get pending leave requests
	if approvalType == "" || approvalType == "all" || approvalType == "leave" {
		requests, count, err := s.dashboardRepo.GetPendingLeaveRequestsList(tenantID, page, pageSize)
		if err != nil {
			return nil, err
		}
		total += count

		for _, r := range requests {
			// Get employee details
			employee, err := s.dashboardRepo.GetEmployeeByEmployeeID(r.EmployeeID)
			employeeName := r.EmployeeID
			department := ""
			var employeePhoto *string
			var employeeIDNum uint

			if err == nil && employee != nil {
				employeeName = fmt.Sprintf("%s %s", employee.FirstName, employee.LastName)
				employeeIDNum = employee.ID
				if employee.DepartmentID != nil {
					department, _ = s.dashboardRepo.GetDepartmentName(*employee.DepartmentID)
				}
				employeePhoto = employee.PhotoURL
			}

			dates := fmt.Sprintf("%s - %s", r.FromDate.Format("Jan 2"), r.ToDate.Format("Jan 2, 2006"))
			days := r.TotalDays

			approvals = append(approvals, models.PendingApproval{
				ID:            r.ID,
				Type:          "leave",
				EmployeeID:    employeeIDNum,
				EmployeeName:  employeeName,
				EmployeePhoto: employeePhoto,
				Department:    department,
				Title:         "Leave Request",
				Description:   r.Reason,
				SubmittedAt:   r.CreatedAt,
				RelativeTime:  getRelativeTime(r.CreatedAt),
				Priority:      "normal",
				ActionURL:     fmt.Sprintf("/leave/requests/%d", r.ID),
				Dates:         &dates,
				Days:          &days,
			})
		}
	}

	// Get counts for summary
	pendingLeave, err := s.dashboardRepo.GetPendingLeaveRequests(tenantID)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.PendingApprovalsResponse{
		Summary: models.ApprovalsSummary{
			Total:                int(pendingLeave),
			LeaveRequests:        int(pendingLeave),
			TimesheetCorrections: 0,
			ExpenseClaims:        0,
		},
		Approvals: approvals,
		Meta: models.ApprovalsMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// ========================================
// Approval Counts for Multiple Roles
// ========================================

// GetApprovalCounts returns counts of pending approvals for different roles
func (s *DashboardService) GetApprovalCounts(tenantID *uint, userRole string, userID uint) (*models.ApprovalCountsResponse, error) {
	counts := &models.ApprovalCountsResponse{
		Role: userRole,
	}

	// Get pending HR requests count (Service Requests & Letters)
	if userRole == "hr" || userRole == "admin" || userRole == "super_admin" {
		hrRequestsCount, _ := s.dashboardRepo.GetPendingHRRequestsCount(tenantID)
		counts.HRRequests = int(hrRequestsCount)
	}

	// Get pending overtime approvals count
	if userRole == "hr" || userRole == "admin" || userRole == "super_admin" || userRole == "manager" {
		overtimeCount, _ := s.dashboardRepo.GetPendingOvertimeCount(tenantID)
		counts.Overtime = int(overtimeCount)
	}

	// Get pending leave approvals count
	if userRole == "hr" || userRole == "admin" || userRole == "super_admin" || userRole == "manager" {
		leaveCount, _ := s.dashboardRepo.GetPendingLeaveRequests(tenantID)
		counts.Leave = int(leaveCount)
	}

	// Get pending timesheet approvals count
	if userRole == "hr" || userRole == "admin" || userRole == "super_admin" || userRole == "manager" {
		timesheetCount, _ := s.dashboardRepo.GetPendingTimesheetCount(tenantID)
		counts.Timesheets = int(timesheetCount)
	}

	// Calculate total pending approvals
	counts.Total = counts.HRRequests + counts.Overtime + counts.Leave + counts.Timesheets

	return counts, nil
}

// ========================================
// Helper Functions
// ========================================

func getRelativeTime(t time.Time) string {
	duration := time.Since(t)

	if duration.Hours() < 1 {
		minutes := int(duration.Minutes())
		if minutes <= 1 {
			return "just now"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}

	if duration.Hours() < 24 {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	days := int(duration.Hours() / 24)
	if days == 1 {
		return "1 day ago"
	}
	if days < 7 {
		return fmt.Sprintf("%d days ago", days)
	}

	weeks := days / 7
	if weeks == 1 {
		return "1 week ago"
	}
	if weeks < 4 {
		return fmt.Sprintf("%d weeks ago", weeks)
	}

	months := days / 30
	if months == 1 {
		return "1 month ago"
	}
	return fmt.Sprintf("%d months ago", months)
}

func getInitials(firstName, lastName string) string {
	initials := ""
	if len(firstName) > 0 {
		initials += strings.ToUpper(string(firstName[0]))
	}
	if len(lastName) > 0 {
		initials += strings.ToUpper(string(lastName[0]))
	}
	return initials
}

func calculateDaysUntil(eventDate time.Time, now time.Time) int {
	// Create this year's event date
	thisYearEvent := time.Date(now.Year(), eventDate.Month(), eventDate.Day(), 0, 0, 0, 0, now.Location())
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// If the event has already passed this year, use next year
	if thisYearEvent.Before(today) {
		thisYearEvent = thisYearEvent.AddDate(1, 0, 0)
	}

	days := int(thisYearEvent.Sub(today).Hours() / 24)
	return days
}

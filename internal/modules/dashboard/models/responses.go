package models

import "time"

// ========================================
// Dashboard Statistics Types
// ========================================

// TrendDirection represents the trend direction
type TrendDirection string

const (
	TrendUp     TrendDirection = "up"
	TrendDown   TrendDirection = "down"
	TrendStable TrendDirection = "stable"
)

// HeadcountStats represents headcount statistics
type HeadcountStats struct {
	Value         int            `json:"value"`
	Change        int            `json:"change"`
	ChangePercent float64        `json:"changePercent"`
	Trend         TrendDirection `json:"trend"`
}

// PresentTodayStats represents present today statistics
type PresentTodayStats struct {
	Value      int     `json:"value"`
	Total      int     `json:"total"`
	Percentage float64 `json:"percentage"`
	Absent     int     `json:"absent"`
	OnLeave    int     `json:"onLeave"`
}

// PendingApprovalsBreakdown represents the breakdown of pending approvals
type PendingApprovalsBreakdown struct {
	LeaveRequests        int `json:"leaveRequests"`
	TimesheetCorrections int `json:"timesheetCorrections"`
	ExpenseClaims        int `json:"expenseClaims"`
}

// PendingApprovalsStats represents pending approvals statistics
type PendingApprovalsStats struct {
	Value     int                       `json:"value"`
	Breakdown PendingApprovalsBreakdown `json:"breakdown"`
}

// PayrollStatusStats represents payroll status
type PayrollStatusStats struct {
	Status         string  `json:"status"`
	CurrentPeriod  string  `json:"currentPeriod"`
	ProcessingDate string  `json:"processingDate"`
	TotalAmount    float64 `json:"totalAmount"`
	Currency       string  `json:"currency"`
}

// AdminStats represents admin dashboard statistics
type AdminStats struct {
	TotalHeadcount   HeadcountStats        `json:"totalHeadcount"`
	PresentToday     PresentTodayStats     `json:"presentToday"`
	PendingApprovals PendingApprovalsStats `json:"pendingApprovals"`
	PayrollStatus    PayrollStatusStats    `json:"payrollStatus"`
}

// AdminSummary represents admin dashboard summary
type AdminSummary struct {
	NewHiresThisMonth     int     `json:"newHiresThisMonth"`
	TerminationsThisMonth int     `json:"terminationsThisMonth"`
	OpenPositions         int     `json:"openPositions"`
	AvgAttendanceRate     float64 `json:"avgAttendanceRate"`
}

// AdminStatisticsResponse represents admin dashboard statistics response
type AdminStatisticsResponse struct {
	View    string       `json:"view"`
	Date    string       `json:"date"`
	Stats   AdminStats   `json:"stats"`
	Summary AdminSummary `json:"summary"`
}

// LeaveBalanceBreakdown represents leave balance breakdown
type LeaveBalanceBreakdown struct {
	Annual   float64 `json:"annual"`
	Sick     float64 `json:"sick"`
	Personal float64 `json:"personal"`
}

// LeaveBalanceStats represents leave balance statistics
type LeaveBalanceStats struct {
	Value        float64               `json:"value"`
	Unit         string                `json:"unit"`
	Breakdown    LeaveBalanceBreakdown `json:"breakdown"`
	UsedThisYear float64               `json:"usedThisYear"`
}

// HoursThisWeekStats represents hours worked this week
type HoursThisWeekStats struct {
	Value    float64 `json:"value"`
	Target   float64 `json:"target"`
	Overtime float64 `json:"overtime"`
	Status   string  `json:"status"` // on_track, behind, ahead
}

// PendingRequestsBreakdown represents pending requests breakdown
type PendingRequestsBreakdown struct {
	LeaveRequests  int `json:"leaveRequests"`
	LetterRequests int `json:"letterRequests"`
}

// PendingRequestsStats represents pending requests statistics
type PendingRequestsStats struct {
	Value     int                      `json:"value"`
	Breakdown PendingRequestsBreakdown `json:"breakdown"`
}

// NextPaydayStats represents next payday information
type NextPaydayStats struct {
	Date            string  `json:"date"`
	DaysRemaining   int     `json:"daysRemaining"`
	EstimatedNetPay float64 `json:"estimatedNetPay"`
	Currency        string  `json:"currency"`
}

// EmployeeStats represents employee dashboard statistics
type EmployeeStats struct {
	LeaveBalance    LeaveBalanceStats    `json:"leaveBalance"`
	HoursThisWeek   HoursThisWeekStats   `json:"hoursThisWeek"`
	PendingRequests PendingRequestsStats `json:"pendingRequests"`
	NextPayday      NextPaydayStats      `json:"nextPayday"`
}

// AttendanceInfo represents today's attendance info
type AttendanceInfo struct {
	TodayStatus      string  `json:"todayStatus"`
	CheckInTime      *string `json:"checkInTime"`
	CheckOutTime     *string `json:"checkOutTime"`
	WorkedHoursToday float64 `json:"workedHoursToday"`
}

// EmployeeStatisticsResponse represents employee dashboard statistics response
type EmployeeStatisticsResponse struct {
	View         string         `json:"view"`
	Date         string         `json:"date"`
	EmployeeID   uint           `json:"employeeId"`
	EmployeeName string         `json:"employeeName"`
	Stats        EmployeeStats  `json:"stats"`
	Attendance   AttendanceInfo `json:"attendance"`
}

// ========================================
// Announcements Types
// ========================================

// AnnouncementCreator represents the creator of an announcement
type AnnouncementCreator struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// AnnouncementAttachmentResponse represents attachment in response
type AnnouncementAttachmentResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Size string `json:"size"`
}

// AnnouncementResponse represents an announcement in the response
type AnnouncementResponse struct {
	ID           uint                             `json:"id"`
	Title        string                           `json:"title"`
	Description  string                           `json:"description"`
	Type         string                           `json:"type"`
	Priority     string                           `json:"priority"`
	CreatedAt    time.Time                        `json:"createdAt"`
	CreatedBy    AnnouncementCreator              `json:"createdBy"`
	ExpiresAt    *time.Time                       `json:"expiresAt"`
	IsRead       bool                             `json:"isRead"`
	Attachments  []AnnouncementAttachmentResponse `json:"attachments"`
	RelativeTime string                           `json:"relativeTime"`
}

// AnnouncementsMeta represents announcements pagination meta
type AnnouncementsMeta struct {
	Page        int   `json:"page"`
	PageSize    int   `json:"pageSize"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"totalPages"`
	UnreadCount int64 `json:"unreadCount"`
}

// AnnouncementsListResponse represents the announcements list response
type AnnouncementsListResponse struct {
	Announcements []AnnouncementResponse `json:"announcements"`
	Meta          AnnouncementsMeta      `json:"meta"`
}

// ========================================
// Approval Counts Types
// ========================================

// ApprovalCountsResponse represents counts of pending approvals for different roles
type ApprovalCountsResponse struct {
	Role        string `json:"role"`                    // User role (hr, admin, manager, etc.)
	Total       int    `json:"total"`                    // Total pending approvals
	HRRequests  int    `json:"hr_requests"`             // Pending HR service requests & letters
	Overtime    int    `json:"overtime"`                 // Pending overtime approvals
	Leave       int    `json:"leave"`                   // Pending leave approvals
	Timesheets  int    `json:"timesheets"`              // Pending timesheet approvals
	LastUpdated string `json:"last_updated,omitempty"`  // Timestamp of last update
}

// ========================================
// Events Types
// ========================================

// UpcomingEvent represents an upcoming event (birthday, anniversary, etc.)
type UpcomingEvent struct {
	ID              uint    `json:"id"`
	EmployeeID      uint    `json:"employeeId"`
	EmployeeName    string  `json:"employeeName"`
	EmployeePhoto   *string `json:"employeePhoto"`
	Department      string  `json:"department"`
	EventType       string  `json:"eventType"` // birthday, anniversary, probation_end
	EventDate       string  `json:"eventDate"`
	DisplayDate     string  `json:"displayDate"`
	DaysUntil       int     `json:"daysUntil"`
	Age             *int    `json:"age,omitempty"`
	YearsOfService  *int    `json:"yearsOfService,omitempty"`
	ProbationMonths *int    `json:"probationMonths,omitempty"`
	Initials        string  `json:"initials"`
}

// EventsSummary represents events summary
type EventsSummary struct {
	BirthdaysThisMonth     int `json:"birthdaysThisMonth"`
	AnniversariesThisMonth int `json:"anniversariesThisMonth"`
	ProbationEndsThisMonth int `json:"probationEndsThisMonth"`
}

// EventsMeta represents events pagination meta
type EventsMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

// EventsResponse represents the events response
type EventsResponse struct {
	Events  []UpcomingEvent `json:"events"`
	Summary EventsSummary   `json:"summary"`
	Meta    EventsMeta      `json:"meta"`
}

// ========================================
// Recent Activity Types
// ========================================

// RecentActivity represents a recent activity item
type RecentActivity struct {
	ID           uint                   `json:"id"`
	Type         string                 `json:"type"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Status       string                 `json:"status"`
	StatusLabel  string                 `json:"statusLabel"`
	Icon         string                 `json:"icon"`
	Color        string                 `json:"color"`
	CreatedAt    time.Time              `json:"createdAt"`
	RelativeTime string                 `json:"relativeTime"`
	Metadata     map[string]interface{} `json:"metadata"`
	ActionURL    string                 `json:"actionUrl"`
}

// ActivitySummary represents activity summary
type ActivitySummary struct {
	PendingRequests  int `json:"pendingRequests"`
	RecentApprovals  int `json:"recentApprovals"`
	NewNotifications int `json:"newNotifications"`
}

// ActivityMeta represents activity pagination meta
type ActivityMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

// MyActivityResponse represents the my activity response
type MyActivityResponse struct {
	Activities []RecentActivity `json:"activities"`
	Summary    ActivitySummary  `json:"summary"`
	Meta       ActivityMeta     `json:"meta"`
}

// ========================================
// Quick Actions Types
// ========================================

// QuickActionBadge represents a badge on a quick action
type QuickActionBadge struct {
	Count *int   `json:"count,omitempty"`
	Text  string `json:"text,omitempty"`
	Type  string `json:"type"` // info, warning, success, error
}

// QuickAction represents a quick action item
type QuickAction struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"`
	Color       string            `json:"color"`
	Href        string            `json:"href"`
	Badge       *QuickActionBadge `json:"badge"`
	Enabled     bool              `json:"enabled"`
}

// QuickActionsResponse represents the quick actions response
type QuickActionsResponse struct {
	View    string        `json:"view"`
	Actions []QuickAction `json:"actions"`
}

// ========================================
// Pending Approvals Types
// ========================================

// PendingApproval represents a pending approval item
type PendingApproval struct {
	ID            uint      `json:"id"`
	Type          string    `json:"type"` // leave, timesheet, expense
	EmployeeID    uint      `json:"employeeId"`
	EmployeeName  string    `json:"employeeName"`
	EmployeePhoto *string   `json:"employeePhoto"`
	Department    string    `json:"department"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	SubmittedAt   time.Time `json:"submittedAt"`
	RelativeTime  string    `json:"relativeTime"`
	Priority      string    `json:"priority"` // low, normal, high, urgent
	ActionURL     string    `json:"actionUrl"`

	// Type-specific fields
	Dates         *string  `json:"dates,omitempty"`
	Days          *float64 `json:"days,omitempty"`
	Amount        *float64 `json:"amount,omitempty"`
	Currency      *string  `json:"currency,omitempty"`
	OriginalTime  *string  `json:"originalTime,omitempty"`
	CorrectedTime *string  `json:"correctedTime,omitempty"`
}

// ApprovalsSummary represents pending approvals summary
type ApprovalsSummary struct {
	Total                int `json:"total"`
	LeaveRequests        int `json:"leaveRequests"`
	TimesheetCorrections int `json:"timesheetCorrections"`
	ExpenseClaims        int `json:"expenseClaims"`
}

// ApprovalsMeta represents approvals pagination meta
type ApprovalsMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

// PendingApprovalsResponse represents the pending approvals response
type PendingApprovalsResponse struct {
	Summary   ApprovalsSummary  `json:"summary"`
	Approvals []PendingApproval `json:"approvals"`
	Meta      ApprovalsMeta     `json:"meta"`
}

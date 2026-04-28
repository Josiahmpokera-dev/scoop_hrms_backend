package models

// AttendanceSummaryResponse represents the KPI summary
type AttendanceSummaryResponse struct {
	PresentRate         float64          `json:"presentRate"`
	AbsentRate          float64          `json:"absentRate"`
	LateRate            float64          `json:"lateRate"`
	AverageWorkingHours float64          `json:"averageWorkingHours"`
	OvertimeHours       float64          `json:"overtimeHours"`
	ExceptionRate       float64          `json:"exceptionRate"`
	OnTimeRate          float64          `json:"onTimeRate"`
	ShiftAdherence      float64          `json:"shiftAdherence"`
	TotalEmployees      int64            `json:"totalEmployees"`
	Trends              AttendanceTrends `json:"trends"`
}

type AttendanceTrends struct {
	PresentChange float64 `json:"presentChange"`
	AbsentChange  float64 `json:"absentChange"`
	LateChange    float64 `json:"lateChange"`
}

// AttendanceTrendResponse represents daily/weekly trends
type AttendanceTrendResponse struct {
	Date           string `json:"date"`
	Day            string `json:"day"`
	Present        int64  `json:"present"`
	Absent         int64  `json:"absent"`
	Late           int64  `json:"late"`
	OnLeave        int64  `json:"onLeave"`
	TotalScheduled int64  `json:"totalScheduled"`
}

// DepartmentAttendanceResponse represents department-wise stats
type DepartmentAttendanceResponse struct {
	DepartmentID         string  `json:"departmentId"`
	DepartmentName       string  `json:"departmentName"`
	Present              int64   `json:"present"`
	Absent               int64   `json:"absent"`
	Late                 int64   `json:"late"`
	OnLeave              int64   `json:"onLeave"`
	TotalEmployees       int64   `json:"totalEmployees"`
	AttendancePercentage float64 `json:"attendancePercentage"`
}

// ComplianceViolationResponse represents a compliance issue
type ComplianceViolationResponse struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	EmployeeID   string `json:"employeeId"`
	EmployeeName string `json:"employeeName"`
	Department   string `json:"department"`
	Details      string `json:"details"`
	Severity     string `json:"severity"`
	Date         string `json:"date"`
}

// OvertimeAnalysisResponse represents overtime stats
type OvertimeAnalysisResponse struct {
	TotalOvertimeHours         float64               `json:"totalOvertimeHours"`
	AverageOvertimePerEmployee float64               `json:"averageOvertimePerEmployee"`
	EstimatedCost              float64               `json:"estimatedCost"`
	Currency                   string                `json:"currency"`
	TopContributors            []OvertimeContributor `json:"topContributors"`
}

type OvertimeContributor struct {
	EmployeeID string  `json:"employeeId"`
	Name       string  `json:"name"`
	Hours      float64 `json:"hours"`
}

// ExportReportRequest represents the request body for exporting reports
type ExportReportRequest struct {
	ReportType   string `json:"reportType" binding:"required"` // daily_attendance, monthly_summary, etc.
	StartDate    string `json:"startDate" binding:"required"`
	EndDate      string `json:"endDate" binding:"required"`
	DepartmentID string `json:"departmentId"`
	Format       string `json:"format"` // pdf, csv, xlsx
}

// ExportReportResponse represents the response for export request
type ExportReportResponse struct {
	DownloadURL string `json:"downloadUrl"`
	JobID       string `json:"jobId,omitempty"`
	ExpiresAt   string `json:"expiresAt,omitempty"`
}

type AttendanceOverviewPeriod struct {
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	View          string `json:"view"`
	LateThreshold string `json:"late_threshold"`
}

type AttendanceOverviewFilters struct {
	DepartmentID *uint   `json:"department_id"`
	LocationID   *uint   `json:"location_id"`
	EmpCode      *string `json:"emp_code"`
}

type AttendanceOverviewTotals struct {
	TotalEmployees int64 `json:"total_employees"`
	Present        int64 `json:"present"`
	Absent         int64 `json:"absent"`
	OnLeave        int64 `json:"on_leave"`
	Late           int64 `json:"late"`
	Exceptions     int64 `json:"exceptions"`
}

type AttendanceOverviewUnits struct {
	Late       string `json:"late"`
	Exceptions string `json:"exceptions"`
}

type AttendanceOverviewMeta struct {
	Units AttendanceOverviewUnits `json:"units"`
}

type AttendanceOverviewData struct {
	Period  AttendanceOverviewPeriod  `json:"period"`
	Filters AttendanceOverviewFilters `json:"filters"`
	Totals  AttendanceOverviewTotals  `json:"totals"`
	Meta    AttendanceOverviewMeta    `json:"meta"`
}

// --- Comprehensive Employee Report Models ---

type ComprehensiveEmployeeReportRequest struct {
	StartDate    string `json:"startDate" binding:"required"` // YYYY-MM-DD
	EndDate      string `json:"endDate" binding:"required"`   // YYYY-MM-DD
	EmployeeID   *uint  `json:"employeeId,omitempty"`
	DepartmentID *uint  `json:"departmentId,omitempty"`
	LocationID   *uint  `json:"locationId,omitempty"`
	Page         int    `json:"page,omitempty"`
	PageSize     int    `json:"pageSize,omitempty"`
}

type ComprehensiveEmployeeReportResponse struct {
	EmployeeID       string                             `json:"employee_id"`
	EmployeeName     string                             `json:"employee_name"`
	Department       string                             `json:"department"`
	Position         string                             `json:"position,omitempty"`
	Period           EmployeeReportPeriod               `json:"period"`
	Attendance       EmployeeAttendanceStats            `json:"attendance"`
	Timesheet        EmployeeTimesheetStats             `json:"timesheet"`
	Leave            EmployeeLeaveStats                 `json:"leave"`
	Overtime         EmployeeOvertimeStats              `json:"overtime"`
	Compliance       EmployeeComplianceStats            `json:"compliance"`
}

type EmployeeReportPeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type EmployeeAttendanceStats struct {
	TotalWorkingDays     int     `json:"total_working_days"`
	PresentDays          int     `json:"present_days"`
	AbsentDays           int     `json:"absent_days"`
	LateDays             int     `json:"late_days"`
	EarlyDepartureDays   int     `json:"early_departure_days"`
	PresentPercentage    float64 `json:"present_percentage"`
	AverageWorkingHours  float64 `json:"average_working_hours"`
	TotalWorkingHours    float64 `json:"total_working_hours"`
}

type EmployeeTimesheetStats struct {
	TotalHours       float64 `json:"total_hours"`
	RegularHours     float64 `json:"regular_hours"`
	OvertimeHours    float64 `json:"overtime_hours"`
	BillableHours    float64 `json:"billable_hours"`
	NonBillableHours float64 `json:"non_billable_hours"`
}

type EmployeeLeaveStats struct {
	TotalLeaveDays   int                        `json:"total_leave_days"`
	LeaveBreakdown   []LeaveTypeBreakdown       `json:"leave_breakdown"`
	PendingRequests  int                        `json:"pending_requests"`
	ApprovedRequests int                        `json:"approved_requests"`
	RejectedRequests int                        `json:"rejected_requests"`
}

type LeaveTypeBreakdown struct {
	LeaveType string  `json:"leave_type"`
	Days      float64 `json:"days"`
}

type EmployeeOvertimeStats struct {
	TotalRequests   int     `json:"total_requests"`
	ApprovedHours   float64 `json:"approved_hours"`
	PendingHours    float64 `json:"pending_hours"`
	RejectedHours   float64 `json:"rejected_hours"`
	TotalPayout     float64 `json:"total_payout,omitempty"`
	CompOffHours    float64 `json:"comp_off_hours,omitempty"`
}

type EmployeeComplianceStats struct {
	ViolationsCount int                         `json:"violations_count"`
	Violations      []ComplianceViolationDetail `json:"violations"`
}

type ComplianceViolationDetail struct {
	Date     string `json:"date"`
	Type     string `json:"type"`
	Details  string `json:"details"`
	Severity string `json:"severity"`
}

// --- Employee Attendance Report Models ---

type EmployeeAttendanceReportResponse struct {
	Summary      EmployeeAttendanceSummary      `json:"summary"`
	Distribution EmployeeAttendanceDistribution `json:"distribution"`
	Trends       []EmployeeAttendanceTrend      `json:"trends"`
	Logs         []EmployeeAttendanceLog        `json:"logs"`
}

type EmployeeAttendanceSummary struct {
	DaysPresent        int     `json:"daysPresent"`
	TotalDays          int     `json:"totalDays"`
	AvgWorkHours       float64 `json:"avgWorkHours"`
	TotalOvertimeHours float64 `json:"totalOvertimeHours"`
	LateArrivalsCount  int     `json:"lateArrivalsCount"`
}

type EmployeeAttendanceDistribution struct {
	Present int `json:"present"`
	Absent  int `json:"absent"`
	OnLeave int `json:"onLeave"`
}

type EmployeeAttendanceTrend struct {
	Date      string  `json:"date"`
	WorkHours float64 `json:"workHours"`
}

type EmployeeAttendanceLog struct {
	Date      string  `json:"date"`
	Status    string  `json:"status"`
	CheckIn   *string `json:"checkIn"`
	CheckOut  *string `json:"checkOut"`
	WorkHours string  `json:"workHours"`
	Overtime  string  `json:"overtime"`
}

package models

// AttendanceSummaryResponse represents the KPI summary
type AttendanceSummaryResponse struct {
	PresentRate         float64            `json:"presentRate"`
	AbsentRate          float64            `json:"absentRate"`
	LateRate            float64            `json:"lateRate"`
	AverageWorkingHours float64            `json:"averageWorkingHours"`
	OvertimeHours       float64            `json:"overtimeHours"`
	ExceptionRate       float64            `json:"exceptionRate"`
	OnTimeRate          float64            `json:"onTimeRate"`
	ShiftAdherence      float64            `json:"shiftAdherence"`
	TotalEmployees      int64              `json:"totalEmployees"`
	Trends              AttendanceTrends   `json:"trends"`
}

type AttendanceTrends struct {
	PresentChange float64 `json:"presentChange"`
	AbsentChange  float64 `json:"absentChange"`
	LateChange    float64 `json:"lateChange"`
}

// AttendanceTrendResponse represents daily/weekly trends
type AttendanceTrendResponse struct {
	Date           string  `json:"date"`
	Day            string  `json:"day"`
	Present        int64   `json:"present"`
	Absent         int64   `json:"absent"`
	Late           int64   `json:"late"`
	OnLeave        int64   `json:"onLeave"`
	TotalScheduled int64   `json:"totalScheduled"`
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
}

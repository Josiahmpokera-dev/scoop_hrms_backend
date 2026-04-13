package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AttendanceReportsHandler handles attendance module reports
type AttendanceReportsHandler struct {
	db      *gorm.DB
	service *services.AttendanceReportService
}

// NewAttendanceReportsHandler creates a new reports handler
func NewAttendanceReportsHandler() *AttendanceReportsHandler {
	return &AttendanceReportsHandler{
		db:      database.GetDB(),
		service: services.NewAttendanceReportService(),
	}
}

// GetSummary returns high-level attendance KPIs
// Endpoint: GET /attendance/reports/summary
func (h *AttendanceReportsHandler) GetSummary(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	deptIDStr := c.Query("departmentId")

	if startDate == "" || endDate == "" {
		errs := map[string]string{}
		if startDate == "" {
			errs["startDate"] = "required"
		}
		if endDate == "" {
			errs["endDate"] = "required"
		}
		response.BadRequest(c, "Invalid date range", errs)
		return
	}

	var deptID *uint
	if deptIDStr != "" {
		id, err := strconv.ParseUint(deptIDStr, 10, 32)
		if err == nil {
			uid := uint(id)
			deptID = &uid
		}
	}

	summary, err := h.service.GetSummary(startDate, endDate, deptID)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch attendance summary", err)
		return
	}

	response.Success(c, "Attendance summary retrieved successfully", summary)
}

// GetTrends returns attendance trends
// Endpoint: GET /attendance/reports/trends
func (h *AttendanceReportsHandler) GetTrends(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	deptIDStr := c.Query("departmentId")
	interval := c.DefaultQuery("interval", "daily")

	if startDate == "" || endDate == "" {
		errs := map[string]string{}
		if startDate == "" {
			errs["startDate"] = "required"
		}
		if endDate == "" {
			errs["endDate"] = "required"
		}
		response.BadRequest(c, "Invalid date range", errs)
		return
	}

	var deptID *uint
	if deptIDStr != "" {
		id, err := strconv.ParseUint(deptIDStr, 10, 32)
		if err == nil {
			uid := uint(id)
			deptID = &uid
		}
	}

	trends, err := h.service.GetTrends(startDate, endDate, deptID, interval)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch attendance trends", err)
		return
	}

	response.Success(c, "Attendance trends retrieved successfully", trends)
}

// GetDepartmentStats returns department-wise attendance
// Endpoint: GET /attendance/reports/by-department
func (h *AttendanceReportsHandler) GetDepartmentStats(c *gin.Context) {
	date := c.Query("date")
	// orgIDStr := c.Query("organizationId") // Not used in repo yet, but can be added

	if date == "" {
		response.BadRequest(c, "Invalid date range", map[string]string{"date": "required"})
		return
	}

	stats, err := h.service.GetDepartmentStats(date, nil)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch department stats", err)
		return
	}

	data := make([]map[string]interface{}, 0, len(stats))
	for _, s := range stats {
		total := s.TotalEmployees
		presentRate := 0.0
		if total > 0 {
			presentRate = (float64(s.Present) / float64(total)) * 100
		}

		data = append(data, map[string]interface{}{
			"departmentId":   s.DepartmentID,
			"departmentName": s.DepartmentName,
			"present":        s.Present,
			"absent":         s.Absent,
			"late":           s.Late,
			"total":          total,
			"presentRate":    presentRate,
		})
	}

	response.Success(c, "Department attendance retrieved successfully", data)
}

// GetComplianceViolations returns compliance issues
// Endpoint: GET /attendance/reports/compliance
func (h *AttendanceReportsHandler) GetComplianceViolations(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	deptIDStr := c.Query("departmentId")
	severity := c.Query("severity")

	if startDate == "" || endDate == "" {
		errs := map[string]string{}
		if startDate == "" {
			errs["startDate"] = "required"
		}
		if endDate == "" {
			errs["endDate"] = "required"
		}
		response.BadRequest(c, "Invalid date range", errs)
		return
	}

	var deptID *uint
	if deptIDStr != "" {
		id, err := strconv.ParseUint(deptIDStr, 10, 32)
		if err == nil {
			uid := uint(id)
			deptID = &uid
		}
	}

	violations, err := h.service.GetComplianceViolations(startDate, endDate, deptID, severity)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch compliance violations", err)
		return
	}

	data := make([]map[string]interface{}, 0, len(violations))
	for _, v := range violations {
		data = append(data, map[string]interface{}{
			"id":            v.ID,
			"employeeId":    v.EmployeeID,
			"employeeName":  v.EmployeeName,
			"department":    v.Department,
			"violationType": v.Type,
			"violationDate": v.Date,
			"severity":      v.Severity,
			"details":       v.Details,
			"status":        "Open",
		})
	}

	response.Success(c, "Compliance violations retrieved successfully", data)
}

// GetOvertimeAnalysis returns overtime analysis
// Endpoint: GET /attendance/reports/overtime
func (h *AttendanceReportsHandler) GetOvertimeAnalysis(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate == "" || endDate == "" {
		errs := map[string]string{}
		if startDate == "" {
			errs["startDate"] = "required"
		}
		if endDate == "" {
			errs["endDate"] = "required"
		}
		response.BadRequest(c, "Invalid date range", errs)
		return
	}

	analysis, err := h.service.GetOvertimeAnalysisUI(startDate, endDate)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch overtime analysis", err)
		return
	}

	response.Success(c, "Overtime analysis retrieved successfully", analysis)
}

// ExportReport triggers a report export
// Endpoint: POST /attendance/reports/export
func (h *AttendanceReportsHandler) ExportReport(c *gin.Context) {
	var req models.ExportReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	result, err := h.service.ExportReport(req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "unsupported") {
			response.ValidationError(c, "Validation failed", map[string]string{"export": err.Error()})
			return
		}
		response.InternalServerError(c, "Failed to initiate report export", err)
		return
	}

	if result != nil && strings.HasPrefix(result.DownloadURL, "/storage/") {
		proto := c.GetHeader("X-Forwarded-Proto")
		if proto == "" {
			if c.Request.TLS != nil {
				proto = "https"
			} else {
				proto = "http"
			}
		}
		host := c.GetHeader("X-Forwarded-Host")
		if host == "" {
			host = c.Request.Host
		}
		result.DownloadURL = proto + "://" + host + result.DownloadURL
	}

	response.Success(c, "Report generation started", result)
}

func (h *AttendanceReportsHandler) GetOverview(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	view := c.Query("view")

	if startDate == "" || endDate == "" {
		errs := map[string]string{}
		if startDate == "" {
			errs["start_date"] = "required"
		}
		if endDate == "" {
			errs["end_date"] = "required"
		}
		response.BadRequest(c, "Invalid date range", errs)
		return
	}

	startTime, err := parseOverviewDateTime(startDate, true)
	if err != nil {
		response.BadRequest(c, "Invalid date range", map[string]string{"start_date": err.Error()})
		return
	}

	endTime, err := parseOverviewDateTime(endDate, false)
	if err != nil {
		response.BadRequest(c, "Invalid date range", map[string]string{"end_date": err.Error()})
		return
	}

	if startTime.After(endTime) {
		response.BadRequest(c, "Invalid date range", map[string]string{"start_date": "must be before end_date"})
		return
	}

	lateThreshold := c.DefaultQuery("late_threshold", "09:20:00")
	if _, err := time.Parse("15:04:05", lateThreshold); err != nil {
		response.ValidationError(c, "Validation failed", map[string]string{"late_threshold": "invalid format HH:mm:ss"})
		return
	}

	var deptID *uint
	if deptIDStr := c.Query("department_id"); deptIDStr != "" {
		id, err := strconv.ParseUint(deptIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid department_id", nil)
			return
		}
		uid := uint(id)
		deptID = &uid
	}

	var locationID *uint
	if locationIDStr := c.Query("location_id"); locationIDStr != "" {
		id, err := strconv.ParseUint(locationIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid location_id", nil)
			return
		}
		uid := uint(id)
		locationID = &uid
	}

	var empCode *string
	if emp := c.Query("emp_code"); emp != "" {
		empCode = &emp
	}

	data, err := h.service.GetOverview(startTime, endTime, view, deptID, locationID, empCode, lateThreshold)
	if err != nil {
		response.InternalServerError(c, "Unable to compute overview", err)
		return
	}

	response.Success(c, "Attendance overview retrieved successfully", data)
}

func parseOverviewDateTime(value string, isStart bool) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			if layout == "2006-01-02" {
				if isStart {
					return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
				}
				return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.UTC), nil
			}
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid format, use YYYY-MM-DD HH:MM:SS")
}

// GetTimesheetSummaryReport returns a summary of timesheets
// Query params: start_date, end_date, department_id
func (h *AttendanceReportsHandler) GetTimesheetSummaryReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		response.BadRequest(c, "start_date and end_date are required", nil)
		return
	}

	sd, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		response.BadRequest(c, "Invalid start_date format (use YYYY-MM-DD)", nil)
		return
	}
	ed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		response.BadRequest(c, "Invalid end_date format (use YYYY-MM-DD)", nil)
		return
	}

	// Summary stats
	var summary struct {
		TotalEntries     int64   `json:"total_entries"`
		TotalHours       float64 `json:"total_hours"`
		BillableHours    float64 `json:"billable_hours"`
		NonBillableHours float64 `json:"non_billable_hours"`
		UniqueEmployees  int64   `json:"unique_employees"`
		UniqueProjects   int64   `json:"unique_projects"`
	}

	h.db.Model(&models.TimesheetEntry{}).
		Where("date >= ? AND date <= ?", sd, ed).
		Select(`
			COUNT(*) as total_entries,
			COALESCE(SUM(hours), 0) as total_hours,
			COALESCE(SUM(CASE WHEN is_billable THEN hours ELSE 0 END), 0) as billable_hours,
			COALESCE(SUM(CASE WHEN NOT is_billable THEN hours ELSE 0 END), 0) as non_billable_hours,
			COUNT(DISTINCT employee_id) as unique_employees,
			COUNT(DISTINCT project_name) as unique_projects
		`).Scan(&summary)

	// Weekly submission status
	var weeklyStatus []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	h.db.Model(&models.TimesheetWeek{}).
		Where("week_start >= ? AND week_end <= ?", sd, ed).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&weeklyStatus)

	// Project distribution
	var projectDist []struct {
		ProjectName   string  `json:"project_name"`
		TotalHours    float64 `json:"total_hours"`
		BillableHours float64 `json:"billable_hours"`
		EmployeeCount int64   `json:"employee_count"`
	}
	h.db.Model(&models.TimesheetEntry{}).
		Where("date >= ? AND date <= ?", sd, ed).
		Select(`
			project_name,
			COALESCE(SUM(hours), 0) as total_hours,
			COALESCE(SUM(CASE WHEN is_billable THEN hours ELSE 0 END), 0) as billable_hours,
			COUNT(DISTINCT employee_id) as employee_count
		`).
		Group("project_name").
		Order("total_hours DESC").
		Limit(20).
		Scan(&projectDist)

	utilizationRate := 0.0
	if summary.TotalHours > 0 {
		utilizationRate = (summary.BillableHours / summary.TotalHours) * 100
	}

	response.Success(c, "Timesheet summary report", gin.H{
		"period": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
		"summary": gin.H{
			"total_entries":      summary.TotalEntries,
			"total_hours":        summary.TotalHours,
			"billable_hours":     summary.BillableHours,
			"non_billable_hours": summary.NonBillableHours,
			"utilization_rate":   utilizationRate,
			"unique_employees":   summary.UniqueEmployees,
			"unique_projects":    summary.UniqueProjects,
		},
		"weekly_submission_status": weeklyStatus,
		"project_distribution":     projectDist,
	})
}

// GetOvertimeSummaryReport returns overtime analytics
func (h *AttendanceReportsHandler) GetOvertimeSummaryReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		response.BadRequest(c, "start_date and end_date are required", nil)
		return
	}

	sd, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		response.BadRequest(c, "Invalid start_date format (use YYYY-MM-DD)", nil)
		return
	}
	ed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		response.BadRequest(c, "Invalid end_date format (use YYYY-MM-DD)", nil)
		return
	}

	// Summary
	var summary struct {
		TotalRequests   int64   `json:"total_requests"`
		ApprovedCount   int64   `json:"approved_count"`
		PendingCount    int64   `json:"pending_count"`
		RejectedCount   int64   `json:"rejected_count"`
		TotalHours      float64 `json:"total_hours"`
		ApprovedHours   float64 `json:"approved_hours"`
		TotalPayout     float64 `json:"total_payout"`
		TotalCompOff    float64 `json:"total_comp_off_hours"`
		UniqueEmployees int64   `json:"unique_employees"`
	}

	h.db.Model(&models.OvertimeRequest{}).
		Where("date >= ? AND date <= ?", sd, ed).
		Select(`
			COUNT(*) as total_requests,
			SUM(CASE WHEN status = 'approved' THEN 1 ELSE 0 END) as approved_count,
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending_count,
			SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END) as rejected_count,
			COALESCE(SUM(hours), 0) as total_hours,
			COALESCE(SUM(CASE WHEN status IN ('approved','processed') THEN hours ELSE 0 END), 0) as approved_hours,
			COALESCE(SUM(CASE WHEN status IN ('approved','processed') THEN COALESCE(payout_amount, 0) ELSE 0 END), 0) as total_payout,
			COALESCE(SUM(CASE WHEN status IN ('approved','processed') THEN COALESCE(comp_off_hours, 0) ELSE 0 END), 0) as total_comp_off,
			COUNT(DISTINCT employee_id) as unique_employees
		`).Scan(&summary)

	// By OT type
	var byType []struct {
		OvertimeType string  `json:"overtime_type"`
		Count        int64   `json:"count"`
		TotalHours   float64 `json:"total_hours"`
	}
	h.db.Model(&models.OvertimeRequest{}).
		Where("date >= ? AND date <= ?", sd, ed).
		Select("overtime_type, COUNT(*) as count, COALESCE(SUM(hours), 0) as total_hours").
		Group("overtime_type").
		Scan(&byType)

	// By compensation type
	var byCompType []struct {
		CompensationType string  `json:"compensation_type"`
		Count            int64   `json:"count"`
		TotalHours       float64 `json:"total_hours"`
	}
	h.db.Model(&models.OvertimeRequest{}).
		Where("date >= ? AND date <= ?", sd, ed).
		Select("compensation_type, COUNT(*) as count, COALESCE(SUM(hours), 0) as total_hours").
		Group("compensation_type").
		Scan(&byCompType)

	// Monthly trend
	var monthlyTrend []struct {
		Month      string  `json:"month"`
		TotalHours float64 `json:"total_hours"`
		Count      int64   `json:"count"`
	}
	h.db.Model(&models.OvertimeRequest{}).
		Where("date >= ? AND date <= ? AND status IN ?", sd, ed, []string{"approved", "processed"}).
		Select("TO_CHAR(date, 'YYYY-MM') as month, COALESCE(SUM(hours), 0) as total_hours, COUNT(*) as count").
		Group("TO_CHAR(date, 'YYYY-MM')").
		Order("month ASC").
		Scan(&monthlyTrend)

	response.Success(c, "Overtime summary report", gin.H{
		"period": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
		"summary":              summary,
		"by_overtime_type":     byType,
		"by_compensation_type": byCompType,
		"monthly_trend":        monthlyTrend,
	})
}

// GetProjectUtilizationReport returns per-project utilization metrics
func (h *AttendanceReportsHandler) GetProjectUtilizationReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		response.BadRequest(c, "start_date and end_date are required", nil)
		return
	}

	sd, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		response.BadRequest(c, "Invalid start_date format (use YYYY-MM-DD)", nil)
		return
	}
	ed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		response.BadRequest(c, "Invalid end_date format (use YYYY-MM-DD)", nil)
		return
	}

	var projects []struct {
		ProjectName      string  `json:"project_name"`
		ClientName       *string `json:"client_name"`
		TotalHours       float64 `json:"total_hours"`
		BillableHours    float64 `json:"billable_hours"`
		NonBillableHours float64 `json:"non_billable_hours"`
		TeamSize         int64   `json:"team_size"`
		UtilizationRate  float64 `json:"utilization_rate"`
	}

	h.db.Model(&models.TimesheetEntry{}).
		Where("date >= ? AND date <= ?", sd, ed).
		Select(`
			project_name,
			MAX(client_name) as client_name,
			COALESCE(SUM(hours), 0) as total_hours,
			COALESCE(SUM(CASE WHEN is_billable THEN hours ELSE 0 END), 0) as billable_hours,
			COALESCE(SUM(CASE WHEN NOT is_billable THEN hours ELSE 0 END), 0) as non_billable_hours,
			COUNT(DISTINCT employee_id) as team_size,
			CASE WHEN SUM(hours) > 0
				THEN (SUM(CASE WHEN is_billable THEN hours ELSE 0 END) / SUM(hours)) * 100
				ELSE 0
			END as utilization_rate
		`).
		Group("project_name").
		Order("total_hours DESC").
		Scan(&projects)

	response.Success(c, "Project utilization report", gin.H{
		"period": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
		"projects": projects,
	})
}

// GetEmployeeUtilizationReport returns per-employee utilization
func (h *AttendanceReportsHandler) GetEmployeeUtilizationReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		response.BadRequest(c, "start_date and end_date are required", nil)
		return
	}

	sd, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		response.BadRequest(c, "Invalid start_date format (use YYYY-MM-DD)", nil)
		return
	}
	ed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		response.BadRequest(c, "Invalid end_date format (use YYYY-MM-DD)", nil)
		return
	}

	var employees []struct {
		EmployeeID       uint    `json:"employee_id"`
		TotalHours       float64 `json:"total_hours"`
		BillableHours    float64 `json:"billable_hours"`
		NonBillableHours float64 `json:"non_billable_hours"`
		UtilizationRate  float64 `json:"utilization_rate"`
		ProjectCount     int64   `json:"project_count"`
	}

	h.db.Model(&models.TimesheetEntry{}).
		Where("date >= ? AND date <= ?", sd, ed).
		Select(`
			employee_id,
			COALESCE(SUM(hours), 0) as total_hours,
			COALESCE(SUM(CASE WHEN is_billable THEN hours ELSE 0 END), 0) as billable_hours,
			COALESCE(SUM(CASE WHEN NOT is_billable THEN hours ELSE 0 END), 0) as non_billable_hours,
			CASE WHEN SUM(hours) > 0
				THEN (SUM(CASE WHEN is_billable THEN hours ELSE 0 END) / SUM(hours)) * 100
				ELSE 0
			END as utilization_rate,
			COUNT(DISTINCT project_name) as project_count
		`).
		Group("employee_id").
		Order("total_hours DESC").
		Scan(&employees)

	response.Success(c, "Employee utilization report", gin.H{
		"period": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
		"employees": employees,
	})
}

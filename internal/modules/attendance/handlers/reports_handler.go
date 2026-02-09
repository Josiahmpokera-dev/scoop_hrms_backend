package handlers

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AttendanceReportsHandler handles attendance module reports
type AttendanceReportsHandler struct {
	db *gorm.DB
}

// NewAttendanceReportsHandler creates a new reports handler
func NewAttendanceReportsHandler() *AttendanceReportsHandler {
	return &AttendanceReportsHandler{
		db: database.GetDB(),
	}
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
		TotalEntries    int64   `json:"total_entries"`
		TotalHours      float64 `json:"total_hours"`
		BillableHours   float64 `json:"billable_hours"`
		NonBillableHours float64 `json:"non_billable_hours"`
		UniqueEmployees int64   `json:"unique_employees"`
		UniqueProjects  int64   `json:"unique_projects"`
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

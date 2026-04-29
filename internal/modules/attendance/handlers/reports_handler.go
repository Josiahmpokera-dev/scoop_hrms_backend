package handlers

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	attendanceRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/services"
	biometricRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/pkg/scheduler"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/encryption"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/storage"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// AttendanceReportsHandler handles attendance module reports
type AttendanceReportsHandler struct {
	db           *gorm.DB
	service      *services.AttendanceReportService
	emailLogRepo *attendanceRepos.EmailLogRepository
}

// NewAttendanceReportsHandler creates a new reports handler
func NewAttendanceReportsHandler() *AttendanceReportsHandler {
	return &AttendanceReportsHandler{
		db:           database.GetDB(),
		service:      services.NewAttendanceReportService(),
		emailLogRepo: attendanceRepos.NewEmailLogRepository(),
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

// GetDepartmentStatsByRange returns department-wise attendance for a date range
// Endpoint: GET /attendance/reports/department-stats
func (h *AttendanceReportsHandler) GetDepartmentStatsByRange(c *gin.Context) {
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

	stats, err := h.service.GetDepartmentAttendanceStats(startDate, endDate)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch department stats", err)
		return
	}

	data := make([]map[string]interface{}, 0, len(stats))
	for _, s := range stats {
		data = append(data, map[string]interface{}{
			"departmentId":         s.DepartmentID,
			"departmentName":       s.DepartmentName,
			"present":              s.Present,
			"absent":               s.Absent,
			"late":                 s.Late,
			"onLeave":              s.OnLeave,
			"attendancePercentage": s.AttendancePercentage,
		})
	}

	response.Success(c, "Department attendance retrieved successfully", data)
}

// GetEmployeeStats returns detailed attendance for a specific employee
// Endpoint: GET /attendance/reports/employee-stats
func (h *AttendanceReportsHandler) GetEmployeeStats(c *gin.Context) {
	employeeIDStr := c.Query("employeeId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if employeeIDStr == "" || startDate == "" || endDate == "" {
		errs := map[string]string{}
		if employeeIDStr == "" {
			errs["employeeId"] = "required"
		}
		if startDate == "" {
			errs["startDate"] = "required"
		}
		if endDate == "" {
			errs["endDate"] = "required"
		}
		response.BadRequest(c, "Invalid parameters", errs)
		return
	}

	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid employee ID", nil)
		return
	}

	data, err := h.service.GetEmployeeAttendanceStats(uint(employeeID), startDate, endDate)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch employee stats", err)
		return
	}

	response.Success(c, "Employee attendance report retrieved successfully", data)
}

// TriggerMonthlyReport manually triggers the end-of-month report generation
// Endpoint: POST /attendance/reports/trigger-monthly
func (h *AttendanceReportsHandler) TriggerMonthlyReport(c *gin.Context) {
	monthStr := c.Query("month") // format YYYY-MM
	if monthStr == "" {
		response.BadRequest(c, "Month is required (YYYY-MM)", nil)
		return
	}

	t, err := time.Parse("2006-01", monthStr)
	if err != nil {
		response.BadRequest(c, "Invalid month format (expected YYYY-MM)", nil)
		return
	}

	// Trigger via scheduler package
	if err := scheduler.ManualTrigger(t.Year(), int(t.Month())); err != nil {
		response.InternalServerError(c, "Failed to trigger report", err)
		return
	}

	response.Success(c, "Report generation triggered successfully. HR will receive an email shortly.", nil)
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

	lateThreshold := c.DefaultQuery("late_threshold", "09:00:00")
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

// DownloadEmployeeAttendanceReport generates an employee attendance report (default current month),
// uploads to storage (S3 when configured), and returns an encrypted download link.
// The same endpoint also handles token-based secure download redirect.
// Endpoint: GET /attendance/reports/employee-attendance-download
func (h *AttendanceReportsHandler) DownloadEmployeeAttendanceReport(c *gin.Context) {
	empCode := strings.TrimSpace(c.Query("emp_code"))
	if empCode == "" {
		response.BadRequest(c, "emp_code is required", map[string]string{"emp_code": "required"})
		return
	}

	employeeRepo := employeeRepos.NewEmployeeRepository()
	lookupDebug := map[string]interface{}{
		"input_emp_code":              empCode,
		"matched_by_employee_code":    false,
		"matched_by_numeric_employee_id": false,
		"matched_by_biometric_enrollment": false,
		"biometric_transactions_found_for_month": false,
	}
	employee, err := employeeRepo.FindByEmployeeID(empCode)
	if err == nil && employee != nil {
		lookupDebug["matched_by_employee_code"] = true
	}
	if err != nil || employee == nil {
		// Backward-compatible fallback:
		// some UIs send numeric employee DB id in emp_code (e.g. "10132").
		if parsedID, parseErr := strconv.ParseUint(empCode, 10, 32); parseErr == nil {
			employee, err = employeeRepo.FindByID(uint(parsedID))
			if err == nil && employee != nil {
				lookupDebug["matched_by_numeric_employee_id"] = true
			}
		}
		// Biometric-link fallback:
		// when UI sends device emp_code, resolve linked employee first.
		if err != nil || employee == nil {
			enrollmentRepo := biometricRepos.NewEnrollmentRepository()
			if enrollment, enrollErr := enrollmentRepo.FindByEmpCode(empCode); enrollErr == nil && enrollment != nil {
				employee, err = employeeRepo.FindByID(enrollment.EmployeeID)
				if err == nil && employee != nil {
					lookupDebug["matched_by_biometric_enrollment"] = true
					lookupDebug["linked_employee_id"] = enrollment.EmployeeID
				}
			}
		}
		if err != nil || employee == nil {
			// Helpful diagnostics to quickly identify why this emp_code cannot resolve.
			month := strings.TrimSpace(c.Query("month"))
			if month == "" {
				month = time.Now().UTC().Format("2006-01")
			}
			if monthStart, parseErr := time.Parse("2006-01", month); parseErr == nil {
				startTime := time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, time.UTC)
				endTime := startTime.AddDate(0, 1, 0).Add(-time.Nanosecond)
				tenantID := middleware.GetTenantID(c)
				txRepo := biometricRepos.NewBioTimeTransactionRepository()
				if exists, txErr := txRepo.HasTransactionsForEmpCode(tenantID, empCode, startTime, endTime); txErr == nil {
					lookupDebug["biometric_transactions_found_for_month"] = exists
				}
			}
			c.JSON(404, response.APIResponse{
				Success: false,
				Message: "Employee not found",
				Error:   lookupDebug,
			})
			return
		}
	}

	// Token flow: decrypt and redirect to stored report URL.
	if token := strings.TrimSpace(c.Query("token")); token != "" {
		enc, err := encryption.NewEncryptionService()
		if err != nil {
			response.InternalServerError(c, "Failed to initialize encryption service", err.Error())
			return
		}
		decryptedURL, err := enc.DecryptField(employee.ID, token)
		if err != nil || decryptedURL == "" {
			response.BadRequest(c, "Invalid or expired download token", map[string]string{"token": "invalid"})
			return
		}
		c.Redirect(302, decryptedURL)
		return
	}

	month := strings.TrimSpace(c.Query("month"))
	now := time.Now().UTC()
	var monthStart time.Time
	if month == "" {
		monthStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		month = monthStart.Format("2006-01")
	} else {
		parsed, err := time.Parse("2006-01", month)
		if err != nil {
			response.BadRequest(c, "month must be in YYYY-MM format", map[string]string{"month": "invalid_format"})
			return
		}
		monthStart = time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)

	stats, err := h.service.GetEmployeeAttendanceStats(employee.ID, monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02"))
	if err != nil {
		response.InternalServerError(c, "Failed to build employee attendance report", err.Error())
		return
	}

	fileBytes, err := buildEmployeeAttendanceExcel(employee.EmployeeID, employee.FirstName, employee.LastName, month, stats)
	if err != nil {
		response.InternalServerError(c, "Failed to generate Excel report", err.Error())
		return
	}

	store := storage.NewStorageService()
	fileName := fmt.Sprintf("attendance_%s_%s.xlsx", employee.EmployeeID, strings.ReplaceAll(month, "-", ""))
	fileURL, err := store.UploadReader(bytes.NewReader(fileBytes), fileName, "reports/attendance", employee.EmployeeID)
	if err != nil {
		response.InternalServerError(c, "Failed to upload attendance report", err.Error())
		return
	}

	enc, err := encryption.NewEncryptionService()
	if err != nil {
		response.InternalServerError(c, "Failed to initialize encryption service", err.Error())
		return
	}
	token, err := enc.EncryptField(employee.ID, fileURL)
	if err != nil {
		response.InternalServerError(c, "Failed to secure download link", err.Error())
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}
	host := c.Request.Host
	if forwardedHost := c.GetHeader("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	// UI-safe relative path (works with frontend baseURL="/api/v1")
	downloadPath := fmt.Sprintf("/attendance/reports/employee-attendance-download?emp_code=%s&token=%s",
		url.QueryEscape(employee.EmployeeID), url.QueryEscape(token))
	// Absolute link for direct browser usage
	downloadLink := fmt.Sprintf("%s://%s/api/v1%s", scheme, host, downloadPath)

	response.Success(c, "Employee attendance report generated successfully", gin.H{
		"employee_id":          employee.EmployeeID,
		"month":                month,
		"format":               "xlsx",
		"storage":              map[string]bool{"s3": store.IsS3()},
		"download_path":        downloadPath,
		"download_link":        downloadLink,
		"encrypted_token":      token,
		"report_file_url":      fileURL,
	})
}

// GetEmailLogs returns monthly report email logs with pagination.
// Endpoint: GET /attendance/reports/email-logs
func (h *AttendanceReportsHandler) GetEmailLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 15
	}

	logs, total, err := h.emailLogRepo.GetLogs(page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch email logs", err.Error())
		return
	}

	items := make([]models.EmailLogResponse, 0, len(logs))
	for _, l := range logs {
		items = append(items, models.EmailLogResponse{
			ID:          l.ID,
			Recipient:   l.Recipient,
			Subject:     l.Subject,
			ReportMonth: l.ReportMonth,
			Status:      l.Status,
			DownloadURL: l.DownloadURL,
			FileType:    l.FileType,
			FileSize:    l.FileSize,
			Date:        l.CreatedAt.Format(time.RFC3339),
		})
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	response.SuccessWithMeta(c, "Email logs retrieved successfully", items, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

func buildEmployeeAttendanceExcel(employeeID, firstName, lastName, month string, data *models.EmployeeAttendanceReportResponse) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Attendance"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Date", "Status", "Check In", "Check Out", "Work Hours", "Overtime"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	title := fmt.Sprintf("Employee: %s - %s %s | Month: %s", employeeID, firstName, lastName, month)
	f.SetCellValue(sheet, "A3", title)

	rowIdx := 5
	for _, log := range data.Logs {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIdx), log.Date)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIdx), log.Status)
		if log.CheckIn != nil {
			f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIdx), *log.CheckIn)
		}
		if log.CheckOut != nil {
			f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIdx), *log.CheckOut)
		}
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIdx), log.WorkHours)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIdx), log.Overtime)
		rowIdx++
	}

	f.SetCellValue(sheet, "H1", "Summary")
	f.SetCellValue(sheet, "H2", "Days Present")
	f.SetCellValue(sheet, "I2", data.Summary.DaysPresent)
	f.SetCellValue(sheet, "H3", "Total Days")
	f.SetCellValue(sheet, "I3", data.Summary.TotalDays)
	f.SetCellValue(sheet, "H4", "Avg Work Hours")
	f.SetCellValue(sheet, "I4", data.Summary.AvgWorkHours)
	f.SetCellValue(sheet, "H5", "Total Overtime Hours")
	f.SetCellValue(sheet, "I5", data.Summary.TotalOvertimeHours)
	f.SetCellValue(sheet, "H6", "Late Arrivals")
	f.SetCellValue(sheet, "I6", data.Summary.LateArrivalsCount)

	f.SetColWidth(sheet, "A", "I", 18)
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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

// GetComprehensiveEmployeeReport godoc
// @Summary Generate comprehensive employee report
// @Description Generate a comprehensive report for employees within a date range, combining attendance, timesheet, leave, overtime, and compliance data
// @Tags Attendance Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param employee_id query int false "Filter by employee ID"
// @Param department_id query int false "Filter by department ID"
// @Param location_id query int false "Filter by location ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} response.APIResponse{data=[]models.ComprehensiveEmployeeReportResponse,meta=response.Meta} "Comprehensive employee report retrieved successfully"
// @Failure 400 {object} response.APIResponse "Invalid date range or parameters"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 403 {object} response.APIResponse "Forbidden"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /attendance/reports/employee-comprehensive [get]
func (h *AttendanceReportsHandler) GetComprehensiveEmployeeReport(c *gin.Context) {
	// Parse and validate dates
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	// Backward compatibility with frontend camelCase params
	if startDate == "" {
		startDate = c.Query("startDate")
	}
	if endDate == "" {
		endDate = c.Query("endDate")
	}

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

	// Validate date format
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		response.BadRequest(c, "Invalid start_date format (use YYYY-MM-DD)", nil)
		return
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		response.BadRequest(c, "Invalid end_date format (use YYYY-MM-DD)", nil)
		return
	}

	// Parse optional filters
	var employeeID, departmentID, locationID *uint

	eid := c.Query("employee_id")
	if eid == "" {
		eid = c.Query("employeeId")
	}
	if eid != "" {
		if id, err := strconv.ParseUint(eid, 10, 32); err == nil {
			uid := uint(id)
			employeeID = &uid
		}
	}

	did := c.Query("department_id")
	if did == "" {
		did = c.Query("departmentId")
	}
	if did != "" {
		if id, err := strconv.ParseUint(did, 10, 32); err == nil {
			uid := uint(id)
			departmentID = &uid
		}
	}

	lid := c.Query("location_id")
	if lid == "" {
		lid = c.Query("locationId")
	}
	if lid != "" {
		if id, err := strconv.ParseUint(lid, 10, 32); err == nil {
			uid := uint(id)
			locationID = &uid
		}
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSizeStr := c.Query("page_size")
	if pageSizeStr == "" {
		pageSizeStr = c.DefaultQuery("pageSize", "20")
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Call service
	reports, total, err := h.service.GetComprehensiveEmployeeReport(startDate, endDate, employeeID, departmentID, locationID, page, pageSize)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "exceed") {
			response.BadRequest(c, "Validation failed", map[string]string{"error": err.Error()})
			return
		}
		response.InternalServerError(c, "Failed to generate comprehensive employee report", err)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response.SuccessWithMeta(c, "Comprehensive employee report retrieved successfully", reports, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

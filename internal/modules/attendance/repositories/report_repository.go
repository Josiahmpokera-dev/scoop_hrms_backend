package repositories

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type AttendanceReportRepository struct {
	db *gorm.DB
}

func NewAttendanceReportRepository() *AttendanceReportRepository {
	return &AttendanceReportRepository{
		db: database.GetDB(),
	}
}

// GetAttendanceSummary calculates high-level KPIs
func (r *AttendanceReportRepository) GetAttendanceSummary(startDate, endDate time.Time, departmentID *uint) (*models.AttendanceSummaryResponse, error) {
	var response models.AttendanceSummaryResponse

	lateThreshold := "09:20:00"

	employeesQuery := r.db.Model(&employeeModels.Employee{}).
		Where("status = ?", "active").
		Where("is_active = ?", true)
	if departmentID != nil {
		employeesQuery = employeesQuery.Where("department_id = ?", *departmentID)
	}

	var totalEmployees int64
	if err := employeesQuery.Count(&totalEmployees).Error; err != nil {
		return nil, err
	}
	response.TotalEmployees = totalEmployees

	days := int64(endDate.Sub(startDate).Hours()/24) + 1
	if days < 1 {
		days = 1
	}
	potentialEmployeeDays := totalEmployees * days

	var presentEmployeeDays int64
	presentSQL := `
SELECT COUNT(DISTINCT bt.emp_code || '|' || DATE(bt.punch_time)::text) AS present_days
FROM biotime_transactions bt
JOIN employees e ON e.employee_id = bt.emp_code
WHERE bt.punch_time BETWEEN ? AND ?
  AND e.status = 'active'
  AND e.is_active = true
`
	args := []interface{}{startDate, endDate}
	if departmentID != nil {
		presentSQL += " AND e.department_id = ?"
		args = append(args, *departmentID)
	}
	if err := r.db.Raw(presentSQL, args...).Scan(&presentEmployeeDays).Error; err != nil {
		return nil, err
	}

	var lateEmployeeDays int64
	lateSQL := `
WITH first_punch AS (
	SELECT
		e.employee_id AS employee_id,
		DATE(bt.punch_time) AS day,
		MIN(bt.punch_time) AS first_checkin
	FROM biotime_transactions bt
	JOIN employees e ON e.employee_id = bt.emp_code
	WHERE bt.punch_time BETWEEN ? AND ?
	  AND e.status = 'active'
	  AND e.is_active = true
`
	lateArgs := []interface{}{startDate, endDate}
	if departmentID != nil {
		lateSQL += " AND e.department_id = ?"
		lateArgs = append(lateArgs, *departmentID)
	}
	lateSQL += `
	GROUP BY e.employee_id, DATE(bt.punch_time)
)
SELECT COUNT(*) FROM first_punch WHERE first_checkin::time > ?::time
`
	lateArgs = append(lateArgs, lateThreshold)
	if err := r.db.Raw(lateSQL, lateArgs...).Scan(&lateEmployeeDays).Error; err != nil {
		return nil, err
	}

	var avgWorkingHours float64
	avgSQL := `
WITH daily_span AS (
	SELECT
		bt.emp_code AS emp_code,
		DATE(bt.punch_time) AS day,
		MIN(bt.punch_time) AS first_checkin,
		MAX(bt.punch_time) AS last_checkout
	FROM biotime_transactions bt
	JOIN employees e ON e.employee_id = bt.emp_code
	WHERE bt.punch_time BETWEEN ? AND ?
	  AND e.status = 'active'
	  AND e.is_active = true
`
	avgArgs := []interface{}{startDate, endDate}
	if departmentID != nil {
		avgSQL += " AND e.department_id = ?"
		avgArgs = append(avgArgs, *departmentID)
	}
	avgSQL += `
	GROUP BY bt.emp_code, DATE(bt.punch_time)
)
SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (last_checkout - first_checkin)) / 3600.0), 0)
FROM daily_span
WHERE last_checkout > first_checkin
`
	if err := r.db.Raw(avgSQL, avgArgs...).Scan(&avgWorkingHours).Error; err != nil {
		return nil, err
	}
	response.AverageWorkingHours = avgWorkingHours

	var overtimeHours float64
	otQuery := r.db.Model(&models.TimesheetEntry{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Where("entry_type = ?", models.EntryTypeOvertime)
	if departmentID != nil {
		otQuery = otQuery.Joins("JOIN employees ON employees.id = timesheet_entries.employee_id").
			Where("employees.department_id = ?", *departmentID)
	}
	if err := otQuery.Select("COALESCE(SUM(hours), 0)").Scan(&overtimeHours).Error; err != nil {
		return nil, err
	}
	response.OvertimeHours = overtimeHours

	if potentialEmployeeDays > 0 {
		response.PresentRate = (float64(presentEmployeeDays) / float64(potentialEmployeeDays)) * 100
		if response.PresentRate > 100 {
			response.PresentRate = 100
		}
		response.AbsentRate = 100 - response.PresentRate

		response.LateRate = (float64(lateEmployeeDays) / float64(potentialEmployeeDays)) * 100
		if response.LateRate < 0 {
			response.LateRate = 0
		}
		if response.LateRate > 100 {
			response.LateRate = 100
		}
	} else {
		response.PresentRate = 0
		response.AbsentRate = 0
		response.LateRate = 0
	}

	response.ExceptionRate = response.LateRate
	response.OnTimeRate = 100 - response.LateRate
	response.ShiftAdherence = response.OnTimeRate

	return &response, nil
}

// GetComplianceViolations fetches compliance issues
func (r *AttendanceReportRepository) GetComplianceViolations(startDate, endDate time.Time, departmentID *uint, severity string) ([]models.ComplianceViolationResponse, error) {
	var violations []models.ComplianceViolationResponse

	// Example Violation 1: Working more than 7 consecutive days (Simplified check)
	// For now, let's just find employees with > 60 hours in the period (assuming weekly view)

	type Result struct {
		EmployeeID     uint
		EmployeeName   string
		DepartmentName string
		TotalHours     float64
	}

	var results []Result
	query := r.db.Table("timesheet_entries").
		Select(`
			employees.id as employee_id,
			CONCAT(employees.first_name, ' ', employees.last_name) as employee_name,
			departments.name as department_name,
			SUM(hours) as total_hours
		`).
		Joins("JOIN employees ON employees.id = timesheet_entries.employee_id").
		Joins("LEFT JOIN departments ON departments.id = employees.department_id").
		Where("timesheet_entries.date BETWEEN ? AND ?", startDate, endDate).
		Group("employees.id, employees.first_name, employees.last_name, departments.name").
		Having("SUM(hours) > ?", 50) // Threshold for violation

	if departmentID != nil {
		query = query.Where("employees.department_id = ?", *departmentID)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	for _, res := range results {
		violations = append(violations, models.ComplianceViolationResponse{
			ID:           "viol_hours_exceeded", // Generate unique ID ideally
			Type:         "Weekly Hours Exceeded",
			EmployeeID:   fmt.Sprintf("emp_%d", res.EmployeeID),
			EmployeeName: res.EmployeeName,
			Department:   res.DepartmentName,
			Details:      "Worked excessive hours",
			Severity:     "High",
			Date:         endDate.Format("2006-01-02"),
		})
	}

	return violations, nil
}

// GetAttendanceTrends fetches daily trends
func (r *AttendanceReportRepository) GetAttendanceTrends(startDate, endDate time.Time, departmentID *uint) ([]models.AttendanceTrendResponse, error) {
	lateThreshold := "09:20:00"

	employeesQuery := r.db.Model(&employeeModels.Employee{}).
		Where("status = ?", "active").
		Where("is_active = ?", true)
	if departmentID != nil {
		employeesQuery = employeesQuery.Where("department_id = ?", *departmentID)
	}

	var totalEmployees int64
	if err := employeesQuery.Count(&totalEmployees).Error; err != nil {
		return nil, err
	}
	if totalEmployees == 0 {
		return []models.AttendanceTrendResponse{}, nil
	}

	sql := `
WITH days AS (
	SELECT generate_series(DATE(?), DATE(?), interval '1 day')::date AS day
),
present AS (
	SELECT DATE(bt.punch_time) AS day, COUNT(DISTINCT e.employee_id) AS present_count
	FROM biotime_transactions bt
	JOIN employees e ON e.employee_id = bt.emp_code
	WHERE bt.punch_time BETWEEN ? AND ?
	  AND e.status = 'active'
	  AND e.is_active = true
` + func() string {
		if departmentID != nil {
			return " AND e.department_id = ?"
		}
		return ""
	}() + `
	GROUP BY DATE(bt.punch_time)
),
first_punch AS (
	SELECT DATE(bt.punch_time) AS day, e.employee_id AS employee_id, MIN(bt.punch_time) AS first_checkin
	FROM biotime_transactions bt
	JOIN employees e ON e.employee_id = bt.emp_code
	WHERE bt.punch_time BETWEEN ? AND ?
	  AND e.status = 'active'
	  AND e.is_active = true
` + func() string {
		if departmentID != nil {
			return " AND e.department_id = ?"
		}
		return ""
	}() + `
	GROUP BY DATE(bt.punch_time), e.employee_id
),
late AS (
	SELECT day, COUNT(*) AS late_count
	FROM first_punch
	WHERE first_checkin::time > ?::time
	GROUP BY day
),
on_leave AS (
	SELECT d.day AS day, COUNT(DISTINCT lr.employee_id) AS on_leave_count
	FROM days d
	JOIN leave_requests lr ON lr.status IN ('approved', 'partially_approved')
		AND lr.from_date::date <= d.day AND lr.to_date::date >= d.day
	JOIN employees e ON e.employee_id = lr.employee_id
	WHERE e.status = 'active' AND e.is_active = true
` + func() string {
		if departmentID != nil {
			return " AND e.department_id = ?"
		}
		return ""
	}() + `
	GROUP BY d.day
)
SELECT
	d.day,
	COALESCE(p.present_count, 0) AS present_count,
	COALESCE(l.late_count, 0) AS late_count,
	COALESCE(ol.on_leave_count, 0) AS on_leave_count
FROM days d
LEFT JOIN present p ON p.day = d.day
LEFT JOIN late l ON l.day = d.day
LEFT JOIN on_leave ol ON ol.day = d.day
ORDER BY d.day
`

	args := []interface{}{startDate, endDate, startDate, endDate}
	if departmentID != nil {
		args = append(args, *departmentID)
	}
	args = append(args, startDate, endDate)
	if departmentID != nil {
		args = append(args, *departmentID)
	}
	args = append(args, lateThreshold)
	if departmentID != nil {
		args = append(args, *departmentID)
	}

	type row struct {
		Day     time.Time
		Present int64
		Late    int64
		OnLeave int64
	}
	var rows []row
	if err := r.db.Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	trends := make([]models.AttendanceTrendResponse, 0, len(rows))
	for _, rrow := range rows {
		presentPct := int64((float64(rrow.Present) / float64(totalEmployees) * 100) + 0.5)
		latePct := int64((float64(rrow.Late) / float64(totalEmployees) * 100) + 0.5)
		onLeavePct := int64((float64(rrow.OnLeave) / float64(totalEmployees) * 100) + 0.5)
		absentPct := int64(100 - presentPct - onLeavePct)
		if absentPct < 0 {
			absentPct = 0
		}

		trends = append(trends, models.AttendanceTrendResponse{
			Date:           rrow.Day.Format("2006-01-02"),
			Day:            rrow.Day.Weekday().String()[:3],
			Present:        presentPct,
			Absent:         absentPct,
			Late:           latePct,
			OnLeave:        onLeavePct,
			TotalScheduled: 100,
		})
	}

	return trends, nil
}

// GetDepartmentAttendance fetches department-wise stats
func (r *AttendanceReportRepository) GetDepartmentAttendance(date time.Time) ([]models.DepartmentAttendanceResponse, error) {
	lateThreshold := "09:20:00"

	type row struct {
		DepartmentID   string
		DepartmentName string
		TotalEmployees int64
		Present        int64
		Late           int64
	}

	sql := `
WITH present AS (
	SELECT e.department_id AS department_id, COUNT(DISTINCT e.employee_id) AS present
	FROM biotime_transactions bt
	JOIN employees e ON e.employee_id = bt.emp_code
	WHERE DATE(bt.punch_time) = DATE(?)
	  AND e.status = 'active'
	  AND e.is_active = true
	GROUP BY e.department_id
),
first_punch AS (
	SELECT e.department_id AS department_id, e.employee_id AS employee_id, MIN(bt.punch_time) AS first_checkin
	FROM biotime_transactions bt
	JOIN employees e ON e.employee_id = bt.emp_code
	WHERE DATE(bt.punch_time) = DATE(?)
	  AND e.status = 'active'
	  AND e.is_active = true
	GROUP BY e.department_id, e.employee_id
),
late AS (
	SELECT department_id, COUNT(*) AS late
	FROM first_punch
	WHERE first_checkin::time > ?::time
	GROUP BY department_id
)
SELECT
	d.id::text AS department_id,
	d.name AS department_name,
	(SELECT COUNT(*) FROM employees e WHERE e.department_id = d.id AND e.status = 'active' AND e.is_active = true) AS total_employees,
	COALESCE(p.present, 0) AS present,
	COALESCE(l.late, 0) AS late
FROM departments d
LEFT JOIN present p ON p.department_id = d.id
LEFT JOIN late l ON l.department_id = d.id
ORDER BY d.name
`

	var rows []row
	if err := r.db.Raw(sql, date, date, lateThreshold).Scan(&rows).Error; err != nil {
		return nil, err
	}

	results := make([]models.DepartmentAttendanceResponse, 0, len(rows))
	for _, rrow := range rows {
		absent := rrow.TotalEmployees - rrow.Present
		if absent < 0 {
			absent = 0
		}

		var pct float64
		if rrow.TotalEmployees > 0 {
			pct = (float64(rrow.Present) / float64(rrow.TotalEmployees)) * 100
		}

		results = append(results, models.DepartmentAttendanceResponse{
			DepartmentID:         rrow.DepartmentID,
			DepartmentName:       rrow.DepartmentName,
			Present:              rrow.Present,
			Absent:               absent,
			Late:                 rrow.Late,
			OnLeave:              0,
			TotalEmployees:       rrow.TotalEmployees,
			AttendancePercentage: pct,
		})
	}

	return results, nil
}

// GetOvertimeAnalysis fetches overtime stats
func (r *AttendanceReportRepository) GetOvertimeAnalysis(startDate, endDate time.Time) (*models.OvertimeAnalysisResponse, error) {
	var response models.OvertimeAnalysisResponse
	response.Currency = "TZS" // Default

	// 1. Total OT Hours
	err := r.db.Model(&models.TimesheetEntry{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Where("entry_type = ?", models.EntryTypeOvertime).
		Select("COALESCE(SUM(hours), 0)").
		Scan(&response.TotalOvertimeHours).Error
	if err != nil {
		return nil, err
	}

	// 2. Top Contributors
	err = r.db.Table("timesheet_entries").
		Select("employees.id as employee_id, CONCAT(employees.first_name, ' ', employees.last_name) as name, SUM(hours) as hours").
		Joins("JOIN employees ON employees.id = timesheet_entries.employee_id").
		Where("timesheet_entries.date BETWEEN ? AND ?", startDate, endDate).
		Where("timesheet_entries.entry_type = ?", models.EntryTypeOvertime).
		Group("employees.id, employees.first_name, employees.last_name").
		Order("hours DESC").
		Limit(5).
		Scan(&response.TopContributors).Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *AttendanceReportRepository) GetOvertimeAnalysisUI(startDate, endDate time.Time) (map[string]interface{}, error) {
	var totalOvertimeHours float64
	err := r.db.Model(&models.TimesheetEntry{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Where("entry_type = ?", models.EntryTypeOvertime).
		Select("COALESCE(SUM(hours), 0)").
		Scan(&totalOvertimeHours).Error
	if err != nil {
		return nil, err
	}

	var employeesWithOvertime int64
	err = r.db.Model(&models.TimesheetEntry{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Where("entry_type = ?", models.EntryTypeOvertime).
		Select("COUNT(DISTINCT employee_id)").
		Scan(&employeesWithOvertime).Error
	if err != nil {
		return nil, err
	}

	var overtimeCost float64
	costSQL := `
SELECT COALESCE(SUM((COALESCE(e.salary, 0) / 160.0) * 1.5 * te.hours), 0) AS overtime_cost
FROM timesheet_entries te
JOIN employees e ON e.id = te.employee_id
WHERE te.date BETWEEN ? AND ?
  AND te.entry_type = ?
`
	if err := r.db.Raw(costSQL, startDate, endDate, models.EntryTypeOvertime).Scan(&overtimeCost).Error; err != nil {
		return nil, err
	}

	type deptRow struct {
		DepartmentName string  `json:"departmentName"`
		Hours          float64 `json:"hours"`
	}
	var topDepartments []deptRow
	deptSQL := `
SELECT COALESCE(d.name, 'Unassigned') AS department_name, COALESCE(SUM(te.hours), 0) AS hours
FROM timesheet_entries te
JOIN employees e ON e.id = te.employee_id
LEFT JOIN departments d ON d.id = e.department_id
WHERE te.date BETWEEN ? AND ?
  AND te.entry_type = ?
GROUP BY COALESCE(d.name, 'Unassigned')
ORDER BY hours DESC
LIMIT 5
`
	if err := r.db.Raw(deptSQL, startDate, endDate, models.EntryTypeOvertime).Scan(&topDepartments).Error; err != nil {
		return nil, err
	}

	type trendRow struct {
		Date  string  `json:"date"`
		Hours float64 `json:"hours"`
	}
	var trends []trendRow
	trendSQL := `
SELECT DATE(te.date)::text AS date, COALESCE(SUM(te.hours), 0) AS hours
FROM timesheet_entries te
WHERE te.date BETWEEN ? AND ?
  AND te.entry_type = ?
GROUP BY DATE(te.date)
ORDER BY DATE(te.date)
`
	if err := r.db.Raw(trendSQL, startDate, endDate, models.EntryTypeOvertime).Scan(&trends).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"totalOvertimeHours":     totalOvertimeHours,
		"overtimeCost":           overtimeCost,
		"employeesWithOvertime":  employeesWithOvertime,
		"topOvertimeDepartments": topDepartments,
		"overtimeTrends":         trends,
	}, nil
}

func (r *AttendanceReportRepository) ExportDailyAttendanceRows(startTime, endTime time.Time, departmentID *uint) ([][]string, error) {
	type row struct {
		EmpCode    string
		FirstName  string
		LastName   string
		Day        time.Time
		CheckIn    *time.Time
		CheckOut   *time.Time
		PunchCount int
	}

	sql := `
SELECT
	bt.emp_code AS emp_code,
	bt.first_name AS first_name,
	bt.last_name AS last_name,
	DATE(bt.punch_time)::date AS day,
	MIN(bt.punch_time)::timestamp AS check_in,
	MAX(bt.punch_time)::timestamp AS check_out,
	COUNT(*) AS punch_count
FROM biotime_transactions bt
JOIN employees e ON e.employee_id = bt.emp_code
WHERE bt.punch_time BETWEEN ? AND ?
  AND e.status = 'active'
  AND e.is_active = true
`
	args := []interface{}{startTime, endTime}
	if departmentID != nil {
		sql += " AND e.department_id = ?"
		args = append(args, *departmentID)
	}
	sql += `
GROUP BY bt.emp_code, bt.first_name, bt.last_name, DATE(bt.punch_time)
ORDER BY day DESC, emp_code ASC
`

	var dbRows []row
	if err := r.db.Raw(sql, args...).Scan(&dbRows).Error; err != nil {
		return nil, err
	}

	out := make([][]string, 0, len(dbRows)+1)
	out = append(out, []string{"emp_code", "name", "date", "checkin", "checkout", "working_hours", "punch_count"})

	for _, rrow := range dbRows {
		name := strings.TrimSpace(strings.TrimSpace(rrow.FirstName + " " + rrow.LastName))
		dateStr := rrow.Day.Format("2006-01-02")

		checkIn := ""
		if rrow.CheckIn != nil {
			checkIn = rrow.CheckIn.Format("15:04:05")
		}
		checkOut := ""
		if rrow.CheckOut != nil {
			checkOut = rrow.CheckOut.Format("15:04:05")
		}

		working := ""
		if rrow.CheckIn != nil && rrow.CheckOut != nil && rrow.CheckOut.After(*rrow.CheckIn) {
			d := rrow.CheckOut.Sub(*rrow.CheckIn)
			h := int(d.Hours())
			m := int(d.Minutes()) % 60
			working = fmt.Sprintf("%02d:%02d", h, m)
		}

		out = append(out, []string{
			rrow.EmpCode,
			name,
			dateStr,
			checkIn,
			checkOut,
			working,
			strconv.Itoa(rrow.PunchCount),
		})
	}

	return out, nil
}

func (r *AttendanceReportRepository) GetAttendanceOverview(startTime, endTime time.Time, departmentID, locationID *uint, empCode *string, lateThreshold string) (totalEmployees, present, onLeave, late int64, err error) {
	employeesQuery := r.db.Model(&employeeModels.Employee{}).
		Where("status = ?", "active").
		Where("is_active = ?", true)

	if departmentID != nil {
		employeesQuery = employeesQuery.Where("department_id = ?", *departmentID)
	}
	if locationID != nil {
		employeesQuery = employeesQuery.Where("location_id = ?", *locationID)
	}
	if empCode != nil && strings.TrimSpace(*empCode) != "" {
		employeesQuery = employeesQuery.Where("employee_id = ?", strings.TrimSpace(*empCode))
	}

	if err := employeesQuery.Count(&totalEmployees).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	presentQuery := r.db.Table("biotime_transactions").
		Joins("JOIN employees ON employees.employee_id = biotime_transactions.emp_code").
		Where("biotime_transactions.punch_time BETWEEN ? AND ?", startTime, endTime).
		Where("employees.status = ?", "active").
		Where("employees.is_active = ?", true)

	if departmentID != nil {
		presentQuery = presentQuery.Where("employees.department_id = ?", *departmentID)
	}
	if locationID != nil {
		presentQuery = presentQuery.Where("employees.location_id = ?", *locationID)
	}
	if empCode != nil && strings.TrimSpace(*empCode) != "" {
		presentQuery = presentQuery.Where("employees.employee_id = ?", strings.TrimSpace(*empCode))
	}

	if err := presentQuery.Select("COUNT(DISTINCT employees.employee_id)").Scan(&present).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	onLeaveQuery := r.db.Table("leave_requests").
		Joins("JOIN employees ON employees.employee_id = leave_requests.employee_id").
		Where("leave_requests.status IN ?", []string{"approved", "partially_approved"}).
		Where("leave_requests.from_date <= ? AND leave_requests.to_date >= ?", endTime, startTime).
		Where("employees.status = ?", "active").
		Where("employees.is_active = ?", true)

	if departmentID != nil {
		onLeaveQuery = onLeaveQuery.Where("employees.department_id = ?", *departmentID)
	}
	if locationID != nil {
		onLeaveQuery = onLeaveQuery.Where("employees.location_id = ?", *locationID)
	}
	if empCode != nil && strings.TrimSpace(*empCode) != "" {
		onLeaveQuery = onLeaveQuery.Where("employees.employee_id = ?", strings.TrimSpace(*empCode))
	}

	if err := onLeaveQuery.Select("COUNT(DISTINCT employees.employee_id)").Scan(&onLeave).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	sql := `
WITH first_punch AS (
	SELECT
		employees.employee_id AS employee_id,
		DATE(biotime_transactions.punch_time) AS day,
		MIN(biotime_transactions.punch_time) AS first_checkin
	FROM biotime_transactions
	JOIN employees ON employees.employee_id = biotime_transactions.emp_code
	WHERE biotime_transactions.punch_time BETWEEN ? AND ?
		AND employees.status = 'active'
		AND employees.is_active = true
`
	args := []interface{}{startTime, endTime}

	if departmentID != nil {
		sql += " AND employees.department_id = ?"
		args = append(args, *departmentID)
	}
	if locationID != nil {
		sql += " AND employees.location_id = ?"
		args = append(args, *locationID)
	}
	if empCode != nil && strings.TrimSpace(*empCode) != "" {
		sql += " AND employees.employee_id = ?"
		args = append(args, strings.TrimSpace(*empCode))
	}

	sql += `
	GROUP BY employees.employee_id, DATE(biotime_transactions.punch_time)
)
SELECT COUNT(*) FROM first_punch WHERE first_checkin::time > ?::time
`
	args = append(args, lateThreshold)

	if err := r.db.Raw(sql, args...).Scan(&late).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	return totalEmployees, present, onLeave, late, nil
}

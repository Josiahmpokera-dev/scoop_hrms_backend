package repositories

import (
	"fmt"
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

	// 1. Total Employees (Active)
	var totalEmployees int64
	query := r.db.Model(&employeeModels.Employee{}).Where("status = ?", "active")
	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}
	if err := query.Count(&totalEmployees).Error; err != nil {
		return nil, err
	}
	response.TotalEmployees = totalEmployees

	// 2. Attendance Stats from Timesheet Entries
	var stats struct {
		TotalHours   float64
		TotalEntries int64
		UniqueDays   int64
	}

	tsQuery := r.db.Model(&models.TimesheetEntry{}).
		Where("date BETWEEN ? AND ?", startDate, endDate)

	if departmentID != nil {
		tsQuery = tsQuery.Joins("JOIN employees ON employees.id = timesheet_entries.employee_id").
			Where("employees.department_id = ?", *departmentID)
	}

	err := tsQuery.Select(`
		COALESCE(SUM(hours), 0) as total_hours,
		COUNT(*) as total_entries,
		COUNT(DISTINCT date) as unique_days
	`).Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	// Calculate Rates (Simplified logic)
	// Assume 8 hours per day standard
	totalWorkingDays := int64(endDate.Sub(startDate).Hours() / 24)
	if totalWorkingDays < 1 {
		totalWorkingDays = 1
	}

	// Potential total man-days = employees * days
	potentialManDays := totalEmployees * totalWorkingDays
	if potentialManDays > 0 {
		// Present rate based on entries vs potential
		// This is an approximation. Ideally we check daily schedules.
		response.PresentRate = float64(stats.TotalEntries) / float64(potentialManDays) * 100
		if response.PresentRate > 100 {
			response.PresentRate = 100
		}
		response.AbsentRate = 100 - response.PresentRate
	}

	if stats.TotalEntries > 0 {
		response.AverageWorkingHours = stats.TotalHours / float64(stats.TotalEntries)
	}

	// 3. Overtime Hours
	// Assuming Overtime is logged as 'overtime' entry type in TimesheetEntry
	// OR query from OvertimeRequest if available. Let's use TimesheetEntry for now as it has EntryType
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

	// Mocking other rates for now as they require shift data integration
	response.LateRate = 5.0 // Placeholder
	response.OnTimeRate = 100 - response.LateRate
	response.ExceptionRate = 2.5
	response.ShiftAdherence = 95.0

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
	var trends []models.AttendanceTrendResponse

	// Group by Date
	// Note: This query might need adjustment for specific SQL dialect (Postgres)
	query := r.db.Table("timesheet_entries").
		Select(`
			DATE(timesheet_entries.date) as date_str,
			COUNT(DISTINCT timesheet_entries.employee_id) as present_count
		`).
		Where("timesheet_entries.date BETWEEN ? AND ?", startDate, endDate).
		Group("DATE(timesheet_entries.date)").
		Order("DATE(timesheet_entries.date)")

	if departmentID != nil {
		query = query.Joins("JOIN employees ON employees.id = timesheet_entries.employee_id").
			Where("employees.department_id = ?", *departmentID)
	}

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dateStr string // Scan as string for compatibility
		var present int64
		if err := rows.Scan(&dateStr, &present); err != nil {
			// Try scanning as time.Time if string fails, or interface{}
			continue
		}

		// Parse date to get Day name
		date, _ := time.Parse("2006-01-02", dateStr[:10]) // Handle potential timestamp format

		trend := models.AttendanceTrendResponse{
			Date:    date.Format("2006-01-02"),
			Day:     date.Weekday().String()[:3],
			Present: present,
			// Mocking others for now
			Absent:         0, // Needs total scheduled calculation
			Late:           0,
			OnLeave:        0,
			TotalScheduled: present, // Placeholder
		}
		trends = append(trends, trend)
	}

	return trends, nil
}

// GetDepartmentAttendance fetches department-wise stats
func (r *AttendanceReportRepository) GetDepartmentAttendance(date time.Time) ([]models.DepartmentAttendanceResponse, error) {
	var results []models.DepartmentAttendanceResponse

	// Join Departments with Employees and TimesheetEntries
	// This is a complex aggregation.
	// Strategy: Get all departments, then count employees and present employees.

	err := r.db.Table("departments").
		Select(`
			departments.id as department_id,
			departments.name as department_name,
			(SELECT COUNT(*) FROM employees WHERE employees.department_id = departments.id AND employees.status = 'active') as total_employees,
			(SELECT COUNT(DISTINCT timesheet_entries.employee_id) FROM timesheet_entries 
			 JOIN employees e ON e.id = timesheet_entries.employee_id 
			 WHERE e.department_id = departments.id AND DATE(timesheet_entries.date) = ?) as present
		`, date).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Post-process to calculate rates
	for i := range results {
		if results[i].TotalEmployees > 0 {
			results[i].Absent = results[i].TotalEmployees - results[i].Present
			results[i].AttendancePercentage = (float64(results[i].Present) / float64(results[i].TotalEmployees)) * 100
		}
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

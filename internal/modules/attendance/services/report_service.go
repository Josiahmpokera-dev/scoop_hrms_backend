package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/repositories"
	"github.com/xuri/excelize/v2"
)

type AttendanceReportService struct {
	repo *repositories.AttendanceReportRepository
}

func NewAttendanceReportService() *AttendanceReportService {
	return &AttendanceReportService{
		repo: repositories.NewAttendanceReportRepository(),
	}
}

func (s *AttendanceReportService) GetSummary(startDate, endDate string, departmentID *uint) (*models.AttendanceSummaryResponse, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format")
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format")
	}

	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	current, err := s.repo.GetAttendanceSummary(start, end, departmentID)
	if err != nil {
		return nil, err
	}

	days := int(end.Sub(start).Hours()/24) + 1
	if days < 1 {
		days = 1
	}

	prevStart := start.AddDate(0, 0, -days)
	prevEnd := start.AddDate(0, 0, -1)
	prevEnd = time.Date(prevEnd.Year(), prevEnd.Month(), prevEnd.Day(), 23, 59, 59, 0, time.UTC)

	prev, err := s.repo.GetAttendanceSummary(prevStart, prevEnd, departmentID)
	if err == nil && prev != nil {
		current.Trends.PresentChange = current.PresentRate - prev.PresentRate
		current.Trends.AbsentChange = current.AbsentRate - prev.AbsentRate
		current.Trends.LateChange = current.LateRate - prev.LateRate
	}

	return current, nil
}

func (s *AttendanceReportService) GetTrends(startDate, endDate string, departmentID *uint, interval string) ([]models.AttendanceTrendResponse, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format")
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format")
	}

	_ = interval
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	return s.repo.GetAttendanceTrends(start, end, departmentID)
}

func (s *AttendanceReportService) GetDepartmentStats(dateStr string, orgID *uint) ([]models.DepartmentAttendanceResponse, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	return s.repo.GetDepartmentAttendance(date)
}

func (s *AttendanceReportService) GetDepartmentAttendanceStats(startDate, endDate string) ([]models.DepartmentAttendanceResponse, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format")
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format")
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)

	if start.After(end) {
		return nil, fmt.Errorf("start date must be before or equal to end date")
	}

	return s.repo.GetDepartmentAttendanceStats(start, end)
}

func (s *AttendanceReportService) GetEmployeeAttendanceStats(employeeID uint, startDate, endDate string) (*models.EmployeeAttendanceReportResponse, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format")
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format")
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)

	if start.After(end) {
		return nil, fmt.Errorf("start date must be before or equal to end date")
	}

	return s.repo.GetEmployeeAttendanceStats(employeeID, start, end)
}

func (s *AttendanceReportService) GetComplianceViolations(startDate, endDate string, departmentID *uint, severity string) ([]models.ComplianceViolationResponse, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format")
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format")
	}

	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	return s.repo.GetComplianceViolations(start, end, departmentID, severity)
}

func (s *AttendanceReportService) GetOvertimeAnalysis(startDate, endDate string) (*models.OvertimeAnalysisResponse, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format")
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format")
	}

	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	return s.repo.GetOvertimeAnalysis(start, end)
}

func (s *AttendanceReportService) GetOvertimeAnalysisUI(startDate, endDate string) (map[string]interface{}, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format")
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format")
	}

	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	return s.repo.GetOvertimeAnalysisUI(start, end)
}

func (s *AttendanceReportService) ExportReport(req models.ExportReportRequest) (*models.ExportReportResponse, error) {
	reportType := strings.TrimSpace(strings.ToLower(req.ReportType))
	format := strings.TrimSpace(strings.ToLower(req.Format))

	if reportType == "" {
		return nil, fmt.Errorf("reportType is required")
	}
	if format == "" {
		format = "csv"
	}
	if format != "csv" && format != "xlsx" {
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}

	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid startDate format")
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid endDate format")
	}
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
	if start.After(end) {
		return nil, fmt.Errorf("invalid date range")
	}

	var deptID *uint
	if req.DepartmentID != "" {
		if v, err := strconv.ParseUint(req.DepartmentID, 10, 32); err == nil {
			u := uint(v)
			deptID = &u
		}
	}

	var rows [][]string
	switch reportType {
	case "daily", "daily_attendance":
		rows, err = s.repo.ExportDailyAttendanceRows(start, end, deptID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported reportType: %s", reportType)
	}

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}

	reportsDir := filepath.Join(storagePath, "reports")
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to prepare reports directory")
	}

	now := time.Now().UTC()
	fileName := fmt.Sprintf("attendance_%s_%s_%d.%s", reportType, req.StartDate, now.Unix(), format)
	fullPath := filepath.Join(reportsDir, fileName)

	if format == "csv" {
		if err := writeCSV(fullPath, rows); err != nil {
			return nil, err
		}
	} else {
		if err := writeXLSX(fullPath, rows); err != nil {
			return nil, err
		}
	}

	jobID := fmt.Sprintf("job_%d", now.Unix())
	downloadURL := "/storage/reports/" + fileName
	expiresAt := now.Add(24 * time.Hour).Format(time.RFC3339)

	return &models.ExportReportResponse{
		DownloadURL: downloadURL,
		JobID:       jobID,
		ExpiresAt:   expiresAt,
	}, nil
}

func writeCSV(path string, rows [][]string) error {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	for _, r := range rows {
		if err := w.Write(r); err != nil {
			return fmt.Errorf("failed to write csv")
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("failed to write csv")
	}
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to save report")
	}
	return nil
}

func writeXLSX(path string, rows [][]string) error {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	for rIdx, row := range rows {
		for cIdx, cell := range row {
			cellName, err := excelize.CoordinatesToCellName(cIdx+1, rIdx+1)
			if err != nil {
				return fmt.Errorf("failed to build xlsx")
			}
			if err := f.SetCellValue(sheet, cellName, cell); err != nil {
				return fmt.Errorf("failed to build xlsx")
			}
		}
	}
	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("failed to save report")
	}
	return nil
}

func (s *AttendanceReportService) GetOverview(startTime, endTime time.Time, view string, departmentID, locationID *uint, empCode *string, lateThreshold string) (*models.AttendanceOverviewData, error) {
	totalEmployees, present, onLeave, late, err := s.repo.GetAttendanceOverview(startTime, endTime, departmentID, locationID, empCode, lateThreshold)
	if err != nil {
		return nil, err
	}

	absent := totalEmployees - present - onLeave
	if absent < 0 {
		absent = 0
	}

	return &models.AttendanceOverviewData{
		Period: models.AttendanceOverviewPeriod{
			StartDate:     startTime.Format("2006-01-02 15:04:05"),
			EndDate:       endTime.Format("2006-01-02 15:04:05"),
			View:          view,
			LateThreshold: lateThreshold,
		},
		Filters: models.AttendanceOverviewFilters{
			DepartmentID: departmentID,
			LocationID:   locationID,
			EmpCode:      empCode,
		},
		Totals: models.AttendanceOverviewTotals{
			TotalEmployees: totalEmployees,
			Present:        present,
			Absent:         absent,
			OnLeave:        onLeave,
			Late:           late,
			Exceptions:     late,
		},
		Meta: models.AttendanceOverviewMeta{
			Units: models.AttendanceOverviewUnits{
				Late:       "employee_days",
				Exceptions: "employee_days",
			},
		},
	}, nil
}

// GetComprehensiveEmployeeReport generates a comprehensive report for employees within a date range
func (s *AttendanceReportService) GetComprehensiveEmployeeReport(startDate, endDate string, employeeID, departmentID, locationID *uint, page, pageSize int) ([]models.ComprehensiveEmployeeReportResponse, int64, error) {
	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid start date format")
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid end date format")
	}

	// Validate date range
	if start.After(end) {
		return nil, 0, fmt.Errorf("start date must be before or equal to end date")
	}

	// Check maximum range (1 year)
	maxEnd := start.AddDate(0, 0, 365)
	if end.After(maxEnd) {
		return nil, 0, fmt.Errorf("date range cannot exceed 365 days")
	}

	// Normalize times
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)

	// Set default pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.repo.GetComprehensiveEmployeeReport(start, end, employeeID, departmentID, locationID, page, pageSize)
}

package services

import (
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/xuri/excelize/v2"
)

type ReportGenerationService struct {
	attendanceRepo *repositories.AttendanceReportRepository
	employeeRepo   *employeeRepos.EmployeeRepository
}

func NewReportGenerationService() *ReportGenerationService {
	return &ReportGenerationService{
		attendanceRepo: repositories.NewAttendanceReportRepository(),
		employeeRepo:   employeeRepos.NewEmployeeRepository(),
	}
}

func (s *ReportGenerationService) GenerateMonthlyReportExcel(year int, month int) ([]byte, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	// Fetch all active employees
	employees, _, err := s.employeeRepo.SearchEmployees(nil, nil, nil, nil, nil, nil, 1, 1000) // Assume max 1000 for now
	if err != nil {
		return nil, fmt.Errorf("failed to fetch employees: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	// Summary Sheet
	summarySheet := "Monthly Summary"
	f.SetSheetName("Sheet1", summarySheet)
	
	// Headers for summary
	headers := []string{"Employee ID", "Name", "Department", "Days Present", "Avg Work Hours", "Total OT", "Late Count"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(summarySheet, cell, h)
	}

	// Styles
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"4F81BD"}, Pattern: 1},
	})
	f.SetRowStyle(summarySheet, 1, 1, headerStyle)

	for i, emp := range employees {
		// Fetch stats for this employee
		stats, err := s.attendanceRepo.GetEmployeeAttendanceStats(emp.ID, startDate, endDate)
		if err != nil {
			continue // Skip failed employees
		}

		// Add to summary
		row := i + 2
		f.SetCellValue(summarySheet, fmt.Sprintf("A%d", row), emp.EmployeeID)
		f.SetCellValue(summarySheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s %s", emp.FirstName, emp.LastName))
		if emp.Department != nil {
			f.SetCellValue(summarySheet, fmt.Sprintf("C%d", row), emp.Department.Name)
		}
		f.SetCellValue(summarySheet, fmt.Sprintf("D%d", row), stats.Summary.DaysPresent)
		f.SetCellValue(summarySheet, fmt.Sprintf("E%d", row), fmt.Sprintf("%.2f", stats.Summary.AvgWorkHours))
		f.SetCellValue(summarySheet, fmt.Sprintf("F%d", row), fmt.Sprintf("%.2f", stats.Summary.TotalOvertimeHours))
		f.SetCellValue(summarySheet, fmt.Sprintf("G%d", row), stats.Summary.LateArrivalsCount)

		// Create individual tab
		sheetName := fmt.Sprintf("%s (%s)", emp.EmployeeID, emp.LastName)
		if len(sheetName) > 31 {
			sheetName = sheetName[:31] // Excel tab name limit
		}
		f.NewSheet(sheetName)

		// Individual sheet headers
		subHeaders := []string{"Date", "Status", "Check In", "Check Out", "Work Hours", "Overtime"}
		for j, sh := range subHeaders {
			cell, _ := excelize.CoordinatesToCellName(j+1, 1)
			f.SetCellValue(sheetName, cell, sh)
		}
		f.SetRowStyle(sheetName, 1, 1, headerStyle)

		// Individual sheet logs
		for k, log := range stats.Logs {
			lRow := k + 2
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", lRow), log.Date)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", lRow), log.Status)
			if log.CheckIn != nil {
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", lRow), *log.CheckIn)
			}
			if log.CheckOut != nil {
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", lRow), *log.CheckOut)
			}
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", lRow), log.WorkHours)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", lRow), log.Overtime)
		}
		
		f.SetColWidth(sheetName, "A", "F", 15)
	}

	f.SetColWidth(summarySheet, "A", "G", 18)
	f.SetActiveSheet(0)

	// Buffer to byte slice
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

package services

import (
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/repositories"
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

	return s.repo.GetAttendanceSummary(start, end, departmentID)
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

	return s.repo.GetAttendanceTrends(start, end, departmentID)
}

func (s *AttendanceReportService) GetDepartmentStats(dateStr string, orgID *uint) ([]models.DepartmentAttendanceResponse, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	return s.repo.GetDepartmentAttendance(date)
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

	return s.repo.GetOvertimeAnalysis(start, end)
}

func (s *AttendanceReportService) ExportReport(req models.ExportReportRequest) (*models.ExportReportResponse, error) {
	// Mock implementation for export
	// In a real scenario, this would generate a file and upload to S3/store locally
	
	// Simulate async job
	jobID := fmt.Sprintf("job_%d", time.Now().Unix())
	downloadURL := fmt.Sprintf("https://api.scoophrms.com/downloads/reports/%s_%s.%s", req.ReportType, req.StartDate, req.Format)

	return &models.ExportReportResponse{
		DownloadURL: downloadURL,
		JobID:       jobID,
	}, nil
}

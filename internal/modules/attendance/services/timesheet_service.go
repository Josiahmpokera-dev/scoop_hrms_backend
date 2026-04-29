package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
)

// TimesheetService handles timesheet business logic
type TimesheetService struct {
	timesheetRepo *repositories.TimesheetRepository
	employeeRepo  *employeeRepos.EmployeeRepository
}

// NewTimesheetService creates a new timesheet service
func NewTimesheetService() *TimesheetService {
	return &TimesheetService{
		timesheetRepo: repositories.NewTimesheetRepository(),
		employeeRepo:  employeeRepos.NewEmployeeRepository(),
	}
}

// getOrCreateWeek gets or creates the weekly timesheet for an employee
func (s *TimesheetService) getOrCreateWeek(employeeID uint, date time.Time) (*models.TimesheetWeek, error) {
	// Calculate the Monday of the week
	weekday := date.Weekday()
	daysToMonday := int(weekday - time.Monday)
	if daysToMonday < 0 {
		daysToMonday += 7
	}
	monday := date.AddDate(0, 0, -daysToMonday)
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
	sunday := monday.AddDate(0, 0, 6)

	// Try to find existing week
	week, err := s.timesheetRepo.FindWeekByEmployeeAndDate(employeeID, monday)
	if err == nil {
		return week, nil
	}

	// Create new week
	week = &models.TimesheetWeek{
		EmployeeID: employeeID,
		WeekStart:  monday,
		WeekEnd:    sunday,
		Status:     models.TimesheetStatusDraft,
	}
	if err := s.timesheetRepo.CreateWeek(week); err != nil {
		return nil, fmt.Errorf("failed to create timesheet week: %w", err)
	}

	return week, nil
}

// CreateEntry creates a new timesheet entry
func (s *TimesheetService) CreateEntry(user *userModels.User, req models.CreateTimesheetEntryRequest) (*models.TimesheetEntry, error) {
	// Resolve employee
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, errors.New("employee record not found for this user")
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	// Resolve hours (optional - default to 0)
	hours := 0.0
	if req.Hours != nil {
		hours = *req.Hours
		// Validate hours precision (0.5h) when provided
		if hours > 0 && math.Mod(hours*2, 1) != 0 {
			return nil, errors.New("hours must be in 0.5h increments")
		}
	}

	// Resolve entry type (default to "regular")
	entryType := models.EntryTypeRegular
	if req.EntryType != "" {
		entryType = models.TimesheetEntryType(req.EntryType)
	}

	// Get or create the week
	week, err := s.getOrCreateWeek(employee.ID, date)
	if err != nil {
		return nil, err
	}

	// Check week is still editable (draft or recalled)
	if week.Status != models.TimesheetStatusDraft && week.Status != models.TimesheetStatusRecalled {
		return nil, fmt.Errorf("cannot add entries to a timesheet in '%s' status", week.Status)
	}

	entry := &models.TimesheetEntry{
		TimesheetID:     week.ID,
		EmployeeID:      employee.ID,
		Date:            date,
		ProjectName:     req.ProjectName,
		ClientName:      req.ClientName,
		TaskDescription: req.TaskDescription,
		EntryType:       entryType,
		Hours:           hours,
		IsBillable:      req.IsBillable,
		Notes:           req.Notes,
	}

	if err := s.timesheetRepo.CreateEntry(entry); err != nil {
		return nil, fmt.Errorf("failed to create timesheet entry: %w", err)
	}

	// Recalculate week totals
	_ = s.timesheetRepo.RecalcWeekTotals(week.ID)

	return entry, nil
}

// BulkCreateEntries creates multiple entries for a week
func (s *TimesheetService) BulkCreateEntries(user *userModels.User, req models.BulkCreateTimesheetRequest) (*models.TimesheetWeek, error) {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, errors.New("employee record not found for this user")
	}

	weekStart, err := time.Parse("2006-01-02", req.WeekStart)
	if err != nil {
		return nil, errors.New("invalid week_start date format, use YYYY-MM-DD")
	}

	week, err := s.getOrCreateWeek(employee.ID, weekStart)
	if err != nil {
		return nil, err
	}

	if week.Status != models.TimesheetStatusDraft && week.Status != models.TimesheetStatusRecalled {
		return nil, fmt.Errorf("cannot add entries to a timesheet in '%s' status", week.Status)
	}

	var entries []models.TimesheetEntry
	for _, r := range req.Entries {
		date, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for entry: %s", r.Date)
		}
		hours := 0.0
		if r.Hours != nil {
			hours = *r.Hours
		}
		eType := models.EntryTypeRegular
		if r.EntryType != "" {
			eType = models.TimesheetEntryType(r.EntryType)
		}
		entries = append(entries, models.TimesheetEntry{
			TimesheetID:     week.ID,
			EmployeeID:      employee.ID,
			Date:            date,
			ProjectName:     r.ProjectName,
			ClientName:      r.ClientName,
			TaskDescription: r.TaskDescription,
			EntryType:       eType,
			Hours:           hours,
			IsBillable:      r.IsBillable,
			Notes:           r.Notes,
		})
	}

	if err := s.timesheetRepo.CreateEntries(entries); err != nil {
		return nil, fmt.Errorf("failed to create timesheet entries: %w", err)
	}

	_ = s.timesheetRepo.RecalcWeekTotals(week.ID)

	// Return updated week
	return s.timesheetRepo.FindWeekByID(week.ID)
}

// UpdateEntry updates a timesheet entry
func (s *TimesheetService) UpdateEntry(user *userModels.User, entryID uint, req models.CreateTimesheetEntryRequest) (*models.TimesheetEntry, error) {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, errors.New("employee record not found for this user")
	}

	entry, err := s.timesheetRepo.FindEntryByID(entryID)
	if err != nil {
		return nil, errors.New("timesheet entry not found")
	}

	if entry.EmployeeID != employee.ID {
		return nil, errors.New("you can only update your own timesheet entries")
	}

	// Check the week status
	week, err := s.timesheetRepo.FindWeekByID(entry.TimesheetID)
	if err != nil {
		return nil, errors.New("timesheet week not found")
	}
	if week.Status != models.TimesheetStatusDraft && week.Status != models.TimesheetStatusRecalled {
		return nil, fmt.Errorf("cannot update entries in a timesheet with '%s' status", week.Status)
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	hours := entry.Hours // keep existing if not provided
	if req.Hours != nil {
		hours = *req.Hours
	}
	eType := entry.EntryType // keep existing if not provided
	if req.EntryType != "" {
		eType = models.TimesheetEntryType(req.EntryType)
	}

	entry.Date = date
	entry.ProjectName = req.ProjectName
	entry.ClientName = req.ClientName
	entry.TaskDescription = req.TaskDescription
	entry.EntryType = eType
	entry.Hours = hours
	entry.IsBillable = req.IsBillable
	entry.Notes = req.Notes

	if err := s.timesheetRepo.UpdateEntry(entry); err != nil {
		return nil, fmt.Errorf("failed to update timesheet entry: %w", err)
	}

	_ = s.timesheetRepo.RecalcWeekTotals(week.ID)
	return entry, nil
}

// DeleteEntry deletes a timesheet entry
func (s *TimesheetService) DeleteEntry(user *userModels.User, entryID uint) error {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return errors.New("employee record not found for this user")
	}

	entry, err := s.timesheetRepo.FindEntryByID(entryID)
	if err != nil {
		return errors.New("timesheet entry not found")
	}

	if entry.EmployeeID != employee.ID {
		return errors.New("you can only delete your own timesheet entries")
	}

	week, err := s.timesheetRepo.FindWeekByID(entry.TimesheetID)
	if err != nil {
		return errors.New("timesheet week not found")
	}
	if week.Status != models.TimesheetStatusDraft && week.Status != models.TimesheetStatusRecalled {
		return fmt.Errorf("cannot delete entries in a timesheet with '%s' status", week.Status)
	}

	weekID := entry.TimesheetID
	if err := s.timesheetRepo.DeleteEntry(entryID); err != nil {
		return fmt.Errorf("failed to delete timesheet entry: %w", err)
	}

	_ = s.timesheetRepo.RecalcWeekTotals(weekID)
	return nil
}

// SubmitTimesheet submits a weekly timesheet for approval
func (s *TimesheetService) SubmitTimesheet(user *userModels.User, timesheetID uint) (*models.TimesheetWeek, error) {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, errors.New("employee record not found for this user")
	}

	week, err := s.timesheetRepo.FindWeekByID(timesheetID)
	if err != nil {
		return nil, errors.New("timesheet not found")
	}

	if week.EmployeeID != employee.ID {
		return nil, errors.New("you can only submit your own timesheets")
	}

	if week.Status != models.TimesheetStatusDraft && week.Status != models.TimesheetStatusRecalled {
		return nil, fmt.Errorf("cannot submit a timesheet in '%s' status", week.Status)
	}

	if len(week.Entries) == 0 {
		return nil, errors.New("cannot submit a timesheet with no entries")
	}

	now := time.Now()
	week.Status = models.TimesheetStatusSubmitted
	week.SubmittedAt = &now

	if err := s.timesheetRepo.UpdateWeek(week); err != nil {
		return nil, fmt.Errorf("failed to submit timesheet: %w", err)
	}

	return week, nil
}

// ApproveTimesheet approves or rejects a timesheet (manager/HR/Admin action)
func (s *TimesheetService) ApproveTimesheet(user *userModels.User, timesheetID uint, req models.ApproveTimesheetRequest) (*models.TimesheetWeek, error) {
	week, err := s.timesheetRepo.FindWeekByID(timesheetID)
	if err != nil {
		return nil, errors.New("timesheet not found")
	}

	if week.Status != models.TimesheetStatusSubmitted {
		return nil, fmt.Errorf("timesheet is not in 'submitted' status (current: %s)", week.Status)
	}

	// Permission check
	canApprove := false
	approverEmployee, _ := s.employeeRepo.FindByUserID(user.ID)
	var approverID uint
	if approverEmployee != nil {
		approverID = approverEmployee.ID
	}

	if user.Role == userModels.RoleAdmin || user.Role == userModels.RoleSuperAdmin || user.Role == userModels.RoleHR {
		canApprove = true
	} else if user.Role == userModels.RoleManager {
		// Check if the employee reports to this manager
		targetEmployee, _ := s.employeeRepo.FindByID(week.EmployeeID)
		if approverEmployee != nil && targetEmployee != nil && targetEmployee.ReportsToID != nil && *targetEmployee.ReportsToID == approverEmployee.ID {
			canApprove = true
		}
	}

	if !canApprove {
		return nil, errors.New("only Admin, HR, or the reporting manager can approve this timesheet")
	}

	now := time.Now()

	if req.Status == "approved" {
		week.Status = models.TimesheetStatusApproved
		week.ApprovedAt = &now
		week.ApprovedBy = &user.ID
	} else {
		week.Status = models.TimesheetStatusRejected
		week.RejectionReason = req.Comments
	}

	// Create approval record
	approval := &models.TimesheetApproval{
		TimesheetID: timesheetID,
		ApproverID:  approverID,
		Level:       1,
		Status:      req.Status,
		Comments:    req.Comments,
		ActionAt:    &now,
	}
	_ = s.timesheetRepo.CreateApproval(approval)

	if err := s.timesheetRepo.UpdateWeek(week); err != nil {
		return nil, fmt.Errorf("failed to update timesheet: %w", err)
	}

	return week, nil
}

// RecallTimesheet lets an employee recall a submitted timesheet
func (s *TimesheetService) RecallTimesheet(user *userModels.User, timesheetID uint) (*models.TimesheetWeek, error) {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, errors.New("employee record not found for this user")
	}

	week, err := s.timesheetRepo.FindWeekByID(timesheetID)
	if err != nil {
		return nil, errors.New("timesheet not found")
	}

	if week.EmployeeID != employee.ID {
		return nil, errors.New("you can only recall your own timesheets")
	}

	if week.Status != models.TimesheetStatusSubmitted {
		return nil, errors.New("can only recall submitted timesheets")
	}

	week.Status = models.TimesheetStatusRecalled
	week.SubmittedAt = nil

	if err := s.timesheetRepo.UpdateWeek(week); err != nil {
		return nil, fmt.Errorf("failed to recall timesheet: %w", err)
	}

	return week, nil
}

// GetMyTimesheets returns the employee's own timesheets
func (s *TimesheetService) GetMyTimesheets(user *userModels.User, status string, page, pageSize int) ([]models.TimesheetWeek, int64, error) {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, 0, errors.New("employee record not found for this user")
	}
	return s.timesheetRepo.ListWeeksByEmployee(employee.ID, status, 0, 0, page, pageSize)
}

// GetTimesheetByID returns a timesheet by ID (owner or admin/HR)
func (s *TimesheetService) GetTimesheetByID(user *userModels.User, timesheetID uint) (*models.TimesheetWeek, error) {
	week, err := s.timesheetRepo.FindWeekByID(timesheetID)
	if err != nil {
		return nil, errors.New("timesheet not found")
	}

	// Check permission
	if user.Role != userModels.RoleAdmin && user.Role != userModels.RoleHR {
		employee, err := s.employeeRepo.FindByUserID(user.ID)
		if err != nil || employee == nil || week.EmployeeID != employee.ID {
			return nil, errors.New("access denied")
		}
	}

	return week, nil
}

// GetPendingApprovals returns submitted timesheets awaiting approval (HR/Admin/Manager)
func (s *TimesheetService) GetPendingApprovals(user *userModels.User, page, pageSize int) ([]models.TimesheetWeek, int64, error) {
	// If HR or Admin, see everything
	if user.Role == userModels.RoleAdmin || user.Role == userModels.RoleSuperAdmin || user.Role == userModels.RoleHR {
		return s.timesheetRepo.ListPendingApprovals(nil, page, pageSize)
	}

	// If Manager, see team only
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, 0, errors.New("employee record not found")
	}

	directReports, _, err := s.employeeRepo.ListByManagerID(employee.ID, 1, 1000)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get team members: %w", err)
	}

	var teamIDs []uint
	for _, dr := range directReports {
		teamIDs = append(teamIDs, dr.ID)
	}

	if len(teamIDs) == 0 {
		return []models.TimesheetWeek{}, 0, nil
	}

	return s.timesheetRepo.ListPendingApprovals(teamIDs, page, pageSize)
}

// GetTeamTimesheets returns timesheets for a manager's team
func (s *TimesheetService) GetTeamTimesheets(user *userModels.User, weekStartStr string, status string, page, pageSize int) ([]models.TimesheetWeek, int64, error) {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, 0, errors.New("employee record not found for this user")
	}

	// Get direct reports
	directReports, _, err := s.employeeRepo.ListByManagerID(employee.ID, 1, 1000)
	if err != nil {
		return nil, 0, errors.New("failed to get team members")
	}

	var employeeIDs []uint
	for _, dr := range directReports {
		employeeIDs = append(employeeIDs, dr.ID)
	}

	if len(employeeIDs) == 0 {
		return []models.TimesheetWeek{}, 0, nil
	}

	var weekStart *time.Time
	if weekStartStr != "" {
		ws, err := time.Parse("2006-01-02", weekStartStr)
		if err == nil {
			weekStart = &ws
		}
	}

	return s.timesheetRepo.ListTeamTimesheets(employeeIDs, weekStart, status, page, pageSize)
}

// CopyLastWeek copies entries from source week to target week
func (s *TimesheetService) CopyLastWeek(user *userModels.User, req models.CopyLastWeekRequest) (*models.TimesheetWeek, error) {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, errors.New("employee record not found for this user")
	}

	sourceStart, err := time.Parse("2006-01-02", req.SourceWeekStart)
	if err != nil {
		return nil, errors.New("invalid source_week_start date format")
	}
	targetStart, err := time.Parse("2006-01-02", req.TargetWeekStart)
	if err != nil {
		return nil, errors.New("invalid target_week_start date format")
	}

	// Find source week
	sourceWeek, err := s.timesheetRepo.FindWeekByEmployeeAndDate(employee.ID, sourceStart)
	if err != nil || sourceWeek == nil {
		return nil, errors.New("source week timesheet not found")
	}

	// Get or create target week
	targetWeek, err := s.getOrCreateWeek(employee.ID, targetStart)
	if err != nil {
		return nil, err
	}

	if targetWeek.Status != models.TimesheetStatusDraft && targetWeek.Status != models.TimesheetStatusRecalled {
		return nil, fmt.Errorf("cannot copy into a timesheet with '%s' status", targetWeek.Status)
	}

	// Calculate day offset
	daysDiff := int(targetStart.Sub(sourceStart).Hours() / 24)

	// Copy entries
	var newEntries []models.TimesheetEntry
	for _, entry := range sourceWeek.Entries {
		newDate := entry.Date.AddDate(0, 0, daysDiff)
		newEntries = append(newEntries, models.TimesheetEntry{
			TimesheetID:     targetWeek.ID,
			EmployeeID:      employee.ID,
			Date:            newDate,
			ProjectName:     entry.ProjectName,
			ClientName:      entry.ClientName,
			TaskDescription: entry.TaskDescription,
			EntryType:       entry.EntryType,
			Hours:           entry.Hours,
			IsBillable:      entry.IsBillable,
			Notes:           entry.Notes,
		})
	}

	if len(newEntries) > 0 {
		if err := s.timesheetRepo.CreateEntries(newEntries); err != nil {
			return nil, fmt.Errorf("failed to copy entries: %w", err)
		}
		_ = s.timesheetRepo.RecalcWeekTotals(targetWeek.ID)
	}

	return s.timesheetRepo.FindWeekByID(targetWeek.ID)
}

// GetEmployeeStats returns timesheet statistics for an employee
func (s *TimesheetService) GetEmployeeStats(user *userModels.User, startDateStr, endDateStr string) (map[string]interface{}, error) {
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil || employee == nil {
		return nil, errors.New("employee record not found for this user")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, errors.New("invalid start_date format")
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, errors.New("invalid end_date format")
	}

	totalHours, billableHours, err := s.timesheetRepo.GetEmployeeTimesheetStats(employee.ID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	utilizationRate := 0.0
	if totalHours > 0 {
		utilizationRate = (billableHours / totalHours) * 100
	}

	projectStats, err := s.timesheetRepo.GetProjectStats(employee.ID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"employee_id":          employee.EmployeeID,
		"employee_name":        employee.FirstName + " " + employee.LastName,
		"period":               map[string]string{"start": startDateStr, "end": endDateStr},
		"total_hours":          totalHours,
		"billable_hours":       billableHours,
		"non_billable_hours":   totalHours - billableHours,
		"utilization_rate":     math.Round(utilizationRate*100) / 100,
		"project_distribution": projectStats,
	}, nil
}

// GetAllTimesheets returns all timesheets (admin/HR) with pagination
func (s *TimesheetService) GetAllTimesheets(status string, page, pageSize int) ([]models.TimesheetWeek, int64, error) {
	return s.timesheetRepo.ListWeeksByEmployee(0, status, 0, 0, page, pageSize)
}

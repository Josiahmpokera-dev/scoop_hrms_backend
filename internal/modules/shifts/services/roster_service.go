package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
)

type RosterService struct {
	repo          *repositories.RosterRepository
	shiftRepo     *repositories.ShiftRepository
	employeeRepo  *employeeRepos.EmployeeRepository
	departmentRepo *departmentRepos.DepartmentRepository
}

func NewRosterService() *RosterService {
	return &RosterService{
		repo:          repositories.NewRosterRepository(),
		shiftRepo:     repositories.NewShiftRepository(),
		employeeRepo:  employeeRepos.NewEmployeeRepository(),
		departmentRepo: departmentRepos.NewDepartmentRepository(),
	}
}

// CreateRosterAssignment creates a new roster assignment
func (s *RosterService) CreateRosterAssignment(req *models.CreateRosterAssignmentRequest, tenantID *uint, createdBy *uint) (*models.RosterAssignment, error) {
	// Validate employee exists
	_, err := s.employeeRepo.FindByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Validate shift exists
	_, err = s.shiftRepo.FindByID(req.ShiftID)
	if err != nil {
		return nil, errors.New("shift not found")
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format. Use YYYY-MM-DD")
	}

	// Check if assignment already exists
	if s.repo.Exists(req.EmployeeID, date) {
		return nil, errors.New("roster assignment already exists for this employee and date")
	}

	status := "draft"
	if req.Status != nil {
		status = *req.Status
	}

	assignment := &models.RosterAssignment{
		EmployeeID: req.EmployeeID,
		Date:       date,
		ShiftID:    req.ShiftID,
		LocationID: req.LocationID,
		Status:     status,
		CreatedBy:  createdBy,
		UpdatedBy:  createdBy,
	}

	if err := s.repo.Create(assignment); err != nil {
		return nil, fmt.Errorf("failed to create roster assignment: %w", err)
	}

	return s.repo.FindByID(assignment.ID)
}

// BulkCreateRosterAssignments creates multiple roster assignments
func (s *RosterService) BulkCreateRosterAssignments(req *models.BulkCreateRosterAssignmentRequest, tenantID *uint, createdBy *uint) ([]models.RosterAssignment, int, int, error) {
	assignments := make([]models.RosterAssignment, 0)
	created := 0
	failed := 0

	for _, reqItem := range req.Assignments {
		// Validate employee
		_, err := s.employeeRepo.FindByEmployeeID(reqItem.EmployeeID)
		if err != nil {
			failed++
			continue
		}

		// Validate shift
		_, err = s.shiftRepo.FindByID(reqItem.ShiftID)
		if err != nil {
			failed++
			continue
		}

		// Parse date
		date, err := time.Parse("2006-01-02", reqItem.Date)
		if err != nil {
			failed++
			continue
		}

		// Skip if already exists
		if s.repo.Exists(reqItem.EmployeeID, date) {
			failed++
			continue
		}

		status := "draft"
		if reqItem.Status != nil {
			status = *reqItem.Status
		}

		assignment := models.RosterAssignment{
			EmployeeID: reqItem.EmployeeID,
			Date:       date,
			ShiftID:    reqItem.ShiftID,
			LocationID: reqItem.LocationID,
			Status:     status,
			CreatedBy:  createdBy,
			UpdatedBy:  createdBy,
		}

		assignments = append(assignments, assignment)
	}

	if len(assignments) > 0 {
		if err := s.repo.BulkCreate(assignments); err != nil {
			return nil, 0, failed, fmt.Errorf("failed to bulk create assignments: %w", err)
		}
		created = len(assignments)
	}

	return assignments, created, failed, nil
}

// GetRosterAssignmentByID retrieves a roster assignment by ID
func (s *RosterService) GetRosterAssignmentByID(id uint) (*models.RosterAssignment, error) {
	return s.repo.FindByID(id)
}

// UpdateRosterAssignment updates a roster assignment
func (s *RosterService) UpdateRosterAssignment(id uint, req *models.UpdateRosterAssignmentRequest, updatedBy *uint) (*models.RosterAssignment, error) {
	assignment, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("roster assignment not found")
	}

	if req.ShiftID != nil {
		// Validate shift exists
		_, err = s.shiftRepo.FindByID(*req.ShiftID)
		if err != nil {
			return nil, errors.New("shift not found")
		}
		assignment.ShiftID = *req.ShiftID
	}

	if req.LocationID != nil {
		assignment.LocationID = req.LocationID
	}

	if req.Status != nil {
		assignment.Status = *req.Status
	}

	assignment.UpdatedBy = updatedBy

	if err := s.repo.Update(assignment); err != nil {
		return nil, fmt.Errorf("failed to update roster assignment: %w", err)
	}

	return s.repo.FindByID(id)
}

// DeleteRosterAssignment deletes a roster assignment
func (s *RosterService) DeleteRosterAssignment(id uint) error {
	return s.repo.Delete(id)
}

// ListRosterAssignments lists roster assignments with pagination and filters
func (s *RosterService) ListRosterAssignments(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.RosterAssignment, int64, error) {
	return s.repo.List(tenantID, page, pageSize, filters)
}

// PublishRosters publishes roster assignments for a date range
func (s *RosterService) PublishRosters(req *models.PublishRosterRequest, tenantID *uint) (int64, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return 0, errors.New("invalid start_date format. Use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return 0, errors.New("invalid end_date format. Use YYYY-MM-DD")
	}

	if startDate.After(endDate) {
		return 0, errors.New("start_date must be before or equal to end_date")
	}

	// Publish rosters
	count, err := s.repo.PublishRosters(tenantID, startDate, endDate)
	if err != nil {
		return 0, fmt.Errorf("failed to publish rosters: %w", err)
	}

	// TODO: Send notifications if requested
	// if req.NotifyEmployees != nil && *req.NotifyEmployees {
	//     // Implement notification logic
	// }

	return count, nil
}

// GetWeeklyRosterView retrieves roster assignments for a specific week
func (s *RosterService) GetWeeklyRosterView(req *models.WeeklyRosterRequest, tenantID *uint) (map[string]interface{}, error) {
	// Parse week (format: YYYY-WW)
	year, week, err := parseWeek(req.Week)
	if err != nil {
		return nil, errors.New("invalid week format. Use YYYY-WW")
	}

	// Calculate start and end dates for the week
	startDate := getWeekStartDate(year, week)
	endDate := startDate.AddDate(0, 0, 6)

	filters := make(map[string]interface{})
	if len(req.EmployeeIDs) > 0 {
		filters["employee_ids"] = req.EmployeeIDs
	}
	if req.DepartmentID != nil {
		filters["department_id"] = *req.DepartmentID
	}
	if req.LocationID != nil {
		filters["location_id"] = *req.LocationID
	}

	assignments, err := s.repo.FindByDateRange(tenantID, startDate, endDate, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly roster: %w", err)
	}

	// Group by employee
	employeeMap := make(map[string]map[string]interface{})
	for _, assignment := range assignments {
		if _, exists := employeeMap[assignment.EmployeeID]; !exists {
			// Get employee details
			employee, _ := s.employeeRepo.FindByEmployeeID(assignment.EmployeeID)
			employeeName := ""
			employeePhoto := ""
			department := ""
			if employee != nil {
				employeeName = employee.FullName()
				if employee.PhotoURL != nil {
					employeePhoto = *employee.PhotoURL
				}
				// Get department name
				if employee.DepartmentID != nil {
					dept, err := s.departmentRepo.FindByID(*employee.DepartmentID)
					if err == nil && dept != nil {
						department = dept.Name
					}
				}
			}
			employeeMap[assignment.EmployeeID] = map[string]interface{}{
				"employee_id":   assignment.EmployeeID,
				"employee_name": employeeName,
				"employee_photo": employeePhoto,
				"department":    department,
				"assignments":   make(map[string]interface{}),
			}
		}

		dateStr := assignment.Date.Format("2006-01-02")
		assignmentsMap := employeeMap[assignment.EmployeeID]["assignments"].(map[string]interface{})
		assignmentsMap[dateStr] = map[string]interface{}{
			"id":            assignment.ID,
			"shift_id":      assignment.ShiftID,
			"shift_name":    getShiftName(assignment.Shift),
			"shift_code":    getShiftCode(assignment.Shift),
			"shift_timing":  getShiftTiming(assignment.Shift),
			"location_id":   assignment.LocationID,
			"location_name": getLocationName(assignment.Location),
			"status":        assignment.Status,
		}
	}

	// Convert map to slice
	employees := make([]map[string]interface{}, 0)
	for _, empData := range employeeMap {
		employees = append(employees, empData)
	}

	return map[string]interface{}{
		"week":       req.Week,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"employees":  employees,
	}, nil
}

// AutoSchedule generates roster assignments automatically
func (s *RosterService) AutoSchedule(req *models.AutoScheduleRequest, tenantID *uint, createdBy *uint) ([]models.RosterAssignment, int, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, 0, errors.New("invalid start_date format. Use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, 0, errors.New("invalid end_date format. Use YYYY-MM-DD")
	}

	if startDate.After(endDate) {
		return nil, 0, errors.New("start_date must be before or equal to end_date")
	}

	// Get employee list
	employeeIDs := req.EmployeeIDs
	if len(employeeIDs) == 0 {
		// If no employees specified, get all active employees
		// This would require employee repository method to list all employees
		// For now, return error
		return nil, 0, errors.New("employee_ids are required for auto-scheduling")
	}

	// Simple auto-scheduling algorithm
	// This is a basic implementation - can be enhanced with more sophisticated rules
	assignments := make([]models.RosterAssignment, 0)
	currentDate := startDate

	// Get available shifts
	shifts, _, _ := s.shiftRepo.List(tenantID, 1, 100, map[string]interface{}{"status": "active"})
	if len(shifts) == 0 {
		return nil, 0, errors.New("no active shifts available")
	}

	// Simple round-robin assignment
	shiftIndex := 0
	for !currentDate.After(endDate) {
		for _, employeeID := range employeeIDs {
			// Skip if assignment already exists
			if s.repo.Exists(employeeID, currentDate) {
				continue
			}

			// Get preferred shift for employee
			shiftID := shifts[shiftIndex%len(shifts)].ID
			if prefs, ok := req.ShiftPreferences[employeeID]; ok && len(prefs) > 0 {
				// Use first preferred shift
				shiftID = prefs[0]
			}

			assignment := models.RosterAssignment{
				EmployeeID: employeeID,
				Date:       currentDate,
				ShiftID:    shiftID,
				LocationID: req.LocationID,
				Status:     "draft",
				CreatedBy:  createdBy,
				UpdatedBy:  createdBy,
			}

			assignments = append(assignments, assignment)
			shiftIndex++
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	if len(assignments) > 0 {
		if err := s.repo.BulkCreate(assignments); err != nil {
			return nil, 0, fmt.Errorf("failed to create auto-scheduled assignments: %w", err)
		}
	}

	return assignments, len(assignments), nil
}

// Helper functions
func parseWeek(weekStr string) (int, int, error) {
	var year, week int
	_, err := fmt.Sscanf(weekStr, "%d-%d", &year, &week)
	return year, week, err
}

func getWeekStartDate(year, week int) time.Time {
	// Simple implementation - first day of the year
	date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	// Adjust to the start of the week (Monday)
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	// Add weeks
	date = date.AddDate(0, 0, (week-1)*7)
	return date
}

func getShiftName(shift *models.Shift) string {
	if shift == nil {
		return ""
	}
	return shift.ShiftName
}

func getShiftCode(shift *models.Shift) string {
	if shift == nil {
		return ""
	}
	return shift.ShiftCode
}

func getShiftTiming(shift *models.Shift) string {
	if shift == nil {
		return ""
	}
	return fmt.Sprintf("%s - %s", shift.StartTime, shift.EndTime)
}

func getLocationName(location *models.RosterLocationRef) string {
	if location == nil {
		return ""
	}
	return location.Name
}

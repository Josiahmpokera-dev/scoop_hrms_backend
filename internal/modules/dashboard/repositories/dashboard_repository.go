package repositories

import (
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	leaveModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"gorm.io/gorm"
)

// DashboardRepository handles dashboard-related database operations
type DashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository creates a new dashboard repository
func NewDashboardRepository() *DashboardRepository {
	return &DashboardRepository{
		db: database.DB,
	}
}

// ========================================
// Employee/Headcount Statistics
// ========================================

// GetTotalHeadcount returns the total number of active employees
func (r *DashboardRepository) GetTotalHeadcount(tenantID *uint) (int64, error) {
	var count int64
	query := r.db.Model(&employeeModels.Employee{}).
		Where("status = ?", employeeModels.StatusActive)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetHeadcountLastMonth returns headcount from last month for trend calculation
func (r *DashboardRepository) GetHeadcountLastMonth(tenantID *uint) (int64, error) {
	var count int64
	lastMonth := time.Now().AddDate(0, -1, 0)

	query := r.db.Model(&employeeModels.Employee{}).
		Where("status = ?", employeeModels.StatusActive).
		Where("created_at <= ?", lastMonth)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetNewHiresThisMonth returns count of new hires this month
func (r *DashboardRepository) GetNewHiresThisMonth(tenantID *uint) (int64, error) {
	var count int64
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	query := r.db.Model(&employeeModels.Employee{}).
		Where("hire_date >= ?", startOfMonth)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetTerminationsThisMonth returns count of terminations this month
func (r *DashboardRepository) GetTerminationsThisMonth(tenantID *uint) (int64, error) {
	var count int64
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	query := r.db.Model(&employeeModels.Employee{}).
		Where("status = ?", employeeModels.StatusTerminated).
		Where("updated_at >= ?", startOfMonth)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ========================================
// Leave Statistics
// ========================================

// GetPendingLeaveRequests returns count of pending leave requests
func (r *DashboardRepository) GetPendingLeaveRequests(tenantID *uint) (int64, error) {
	var count int64
	query := r.db.Model(&leaveModels.LeaveRequest{}).
		Where("status = ?", leaveModels.LeaveRequestStatusPending)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetEmployeesOnLeaveToday returns count of employees on leave today
func (r *DashboardRepository) GetEmployeesOnLeaveToday(tenantID *uint, date time.Time) (int64, error) {
	var count int64
	query := r.db.Model(&leaveModels.LeaveRequest{}).
		Where("status = ?", leaveModels.LeaveRequestStatusApproved).
		Where("from_date <= ? AND to_date >= ?", date, date)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetEmployeePendingLeaveRequests returns count of pending leave requests for an employee
func (r *DashboardRepository) GetEmployeePendingLeaveRequests(employeeID string) (int64, error) {
	var count int64
	if err := r.db.Model(&leaveModels.LeaveRequest{}).
		Where("employee_id = ?", employeeID).
		Where("status = ?", leaveModels.LeaveRequestStatusPending).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetEmployeeLeaveBalance returns leave balance summary for an employee
func (r *DashboardRepository) GetEmployeeLeaveBalance(employeeID string, year int) (map[string]float64, error) {
	balances := make(map[string]float64)

	var results []struct {
		LeaveTypeCode string
		Balance       float64
	}

	if err := r.db.Model(&leaveModels.LeaveBalance{}).
		Select("leave_type_code, available as balance").
		Where("employee_id = ? AND year = ?", employeeID, year).
		Find(&results).Error; err != nil {
		return nil, err
	}

	for _, r := range results {
		balances[r.LeaveTypeCode] = r.Balance
	}

	return balances, nil
}

// GetEmployeeTotalLeaveUsed returns total leave used by an employee this year
func (r *DashboardRepository) GetEmployeeTotalLeaveUsed(employeeID string, year int) (float64, error) {
	var total float64
	if err := r.db.Model(&leaveModels.LeaveRequest{}).
		Select("COALESCE(SUM(total_days), 0)").
		Where("employee_id = ?", employeeID).
		Where("leave_period_year = ?", year).
		Where("status = ?", leaveModels.LeaveRequestStatusApproved).
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// ========================================
// Events (Birthdays, Anniversaries)
// ========================================

// UpcomingEventRow represents a row from upcoming events query
type UpcomingEventRow struct {
	ID           uint
	EmployeeID   string
	FirstName    string
	LastName     string
	PhotoURL     *string
	Department   string
	EventType    string
	EventDate    time.Time
	DateOfBirth  *time.Time
	HireDate     *time.Time
	ProbationEnd *time.Time
}

// GetUpcomingBirthdays returns employees with birthdays in the next N days
func (r *DashboardRepository) GetUpcomingBirthdays(tenantID *uint, daysAhead int, page, pageSize int) ([]UpcomingEventRow, int64, error) {
	var events []UpcomingEventRow
	var total int64

	now := time.Now()
	endDate := now.AddDate(0, 0, daysAhead)

	// Build query for birthdays - compare month and day
	query := r.db.Table("employees e").
		Select(`
			e.id,
			e.employee_id,
			e.first_name,
			e.last_name,
			e.photo_url,
			COALESCE(d.name, '') as department,
			'birthday' as event_type,
			e.date_of_birth as event_date,
			e.date_of_birth
		`).
		Joins("LEFT JOIN departments d ON e.department_id = d.id").
		Where("e.status = ?", employeeModels.StatusActive).
		Where("e.date_of_birth IS NOT NULL").
		Where(`(
			(EXTRACT(MONTH FROM e.date_of_birth) = EXTRACT(MONTH FROM ?::timestamp) AND EXTRACT(DAY FROM e.date_of_birth) >= EXTRACT(DAY FROM ?::timestamp))
			OR (EXTRACT(MONTH FROM e.date_of_birth) = EXTRACT(MONTH FROM ?::timestamp) AND EXTRACT(DAY FROM e.date_of_birth) <= EXTRACT(DAY FROM ?::timestamp))
			OR (EXTRACT(MONTH FROM e.date_of_birth) > EXTRACT(MONTH FROM ?::timestamp) AND EXTRACT(MONTH FROM e.date_of_birth) < EXTRACT(MONTH FROM ?::timestamp))
		)`, now, now, endDate, endDate, now, endDate)

	if tenantID != nil {
		query = query.Where("e.tenant_id = ?", *tenantID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := query.
		Order("EXTRACT(MONTH FROM e.date_of_birth), EXTRACT(DAY FROM e.date_of_birth)").
		Offset(offset).
		Limit(pageSize).
		Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetUpcomingAnniversaries returns employees with work anniversaries in the next N days
func (r *DashboardRepository) GetUpcomingAnniversaries(tenantID *uint, daysAhead int, page, pageSize int) ([]UpcomingEventRow, int64, error) {
	var events []UpcomingEventRow
	var total int64

	now := time.Now()
	endDate := now.AddDate(0, 0, daysAhead)

	query := r.db.Table("employees e").
		Select(`
			e.id,
			e.employee_id,
			e.first_name,
			e.last_name,
			e.photo_url,
			COALESCE(d.name, '') as department,
			'anniversary' as event_type,
			e.hire_date as event_date,
			e.hire_date
		`).
		Joins("LEFT JOIN departments d ON e.department_id = d.id").
		Where("e.status = ?", employeeModels.StatusActive).
		Where("e.hire_date IS NOT NULL").
		Where(`(
			(EXTRACT(MONTH FROM e.hire_date) = EXTRACT(MONTH FROM ?::timestamp) AND EXTRACT(DAY FROM e.hire_date) >= EXTRACT(DAY FROM ?::timestamp))
			OR (EXTRACT(MONTH FROM e.hire_date) = EXTRACT(MONTH FROM ?::timestamp) AND EXTRACT(DAY FROM e.hire_date) <= EXTRACT(DAY FROM ?::timestamp))
			OR (EXTRACT(MONTH FROM e.hire_date) > EXTRACT(MONTH FROM ?::timestamp) AND EXTRACT(MONTH FROM e.hire_date) < EXTRACT(MONTH FROM ?::timestamp))
		)`, now, now, endDate, endDate, now, endDate)

	if tenantID != nil {
		query = query.Where("e.tenant_id = ?", *tenantID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.
		Order("EXTRACT(MONTH FROM e.hire_date), EXTRACT(DAY FROM e.hire_date)").
		Offset(offset).
		Limit(pageSize).
		Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetBirthdaysThisMonth returns count of birthdays this month
func (r *DashboardRepository) GetBirthdaysThisMonth(tenantID *uint) (int64, error) {
	var count int64
	now := time.Now()

	query := r.db.Model(&employeeModels.Employee{}).
		Where("status = ?", employeeModels.StatusActive).
		Where("date_of_birth IS NOT NULL").
		Where("EXTRACT(MONTH FROM date_of_birth) = ?", now.Month())

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetAnniversariesThisMonth returns count of anniversaries this month
func (r *DashboardRepository) GetAnniversariesThisMonth(tenantID *uint) (int64, error) {
	var count int64
	now := time.Now()

	query := r.db.Model(&employeeModels.Employee{}).
		Where("status = ?", employeeModels.StatusActive).
		Where("hire_date IS NOT NULL").
		Where("EXTRACT(MONTH FROM hire_date) = ?", now.Month())

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ========================================
// Employee Activity
// ========================================

// GetEmployeeRecentLeaveRequests returns recent leave requests for an employee
func (r *DashboardRepository) GetEmployeeRecentLeaveRequests(employeeID string, days int, page, pageSize int) ([]leaveModels.LeaveRequest, int64, error) {
	var requests []leaveModels.LeaveRequest
	var total int64

	since := time.Now().AddDate(0, 0, -days)

	query := r.db.Model(&leaveModels.LeaveRequest{}).
		Where("employee_id = ?", employeeID).
		Where("created_at >= ?", since)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.
		Preload("LeaveType").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetPendingLeaveRequestsList returns list of pending leave requests with employee details
func (r *DashboardRepository) GetPendingLeaveRequestsList(tenantID *uint, page, pageSize int) ([]leaveModels.LeaveRequest, int64, error) {
	var requests []leaveModels.LeaveRequest
	var total int64

	query := r.db.Model(&leaveModels.LeaveRequest{}).
		Where("status = ?", leaveModels.LeaveRequestStatusPending)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.
		Preload("LeaveType").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetEmployeeByEmployeeID returns employee by employee_id string
func (r *DashboardRepository) GetEmployeeByEmployeeID(employeeID string) (*employeeModels.Employee, error) {
	var employee employeeModels.Employee
	if err := r.db.Where("employee_id = ?", employeeID).First(&employee).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetEmployeeByUserID returns employee by user_id
func (r *DashboardRepository) GetEmployeeByUserID(userID uint) (*employeeModels.Employee, error) {
	var employee employeeModels.Employee
	if err := r.db.Where("user_id = ?", userID).
		First(&employee).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetDepartmentName returns department name by ID
func (r *DashboardRepository) GetDepartmentName(departmentID uint) (string, error) {
	var name string
	if err := r.db.Table("departments").
		Select("name").
		Where("id = ?", departmentID).
		Scan(&name).Error; err != nil {
		return "", err
	}
	return name, nil
}

// GetUserName returns user's full name by ID
func (r *DashboardRepository) GetUserName(userID uint) (string, error) {
	var result struct {
		FirstName string
		LastName  string
	}
	if err := r.db.Table("users").
		Select("first_name, last_name").
		Where("id = ?", userID).
		Scan(&result).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", result.FirstName, result.LastName), nil
}

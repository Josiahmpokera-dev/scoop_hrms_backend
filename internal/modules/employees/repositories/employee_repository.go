package repositories

import (
	"fmt"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

// EmployeeRepository handles employee database operations
type EmployeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository creates a new employee repository
func NewEmployeeRepository() *EmployeeRepository {
	return &EmployeeRepository{
		db: database.GetDB(),
	}
}

// Create creates a new employee
func (r *EmployeeRepository) Create(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

// FindByID finds an employee by ID
func (r *EmployeeRepository) FindByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.Where("id = ?", id).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// FindByEmployeeID finds an employee by employee ID
func (r *EmployeeRepository) FindByEmployeeID(employeeID string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.Where("employee_id = ?", employeeID).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// FindByEmail finds an employee by email
func (r *EmployeeRepository) FindByEmail(email string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.Where("email = ?", email).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// FindByUserID finds an employee by user ID
func (r *EmployeeRepository) FindByUserID(userID uint) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.Where("user_id = ?", userID).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// Update updates an employee
func (r *EmployeeRepository) Update(employee *models.Employee) error {
	return r.db.Save(employee).Error
}

// List lists employees with pagination
func (r *EmployeeRepository) List(limit, offset int) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	// Count total
	if err := r.db.Model(&models.Employee{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&employees).Error
	return employees, total, err
}

// ListByDepartment lists employees by department
func (r *EmployeeRepository) ListByDepartment(departmentID uint) ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.Where("department_id = ?", departmentID).Find(&employees).Error
	return employees, err
}

// ListByStatus lists employees by status
func (r *EmployeeRepository) ListByStatus(status models.EmployeeStatus) ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.Where("status = ?", status).Find(&employees).Error
	return employees, err
}

// ExistsByEmployeeID checks if an employee with the employee ID exists
func (r *EmployeeRepository) ExistsByEmployeeID(employeeID string) bool {
	var count int64
	r.db.Model(&models.Employee{}).Where("employee_id = ?", employeeID).Count(&count)
	return count > 0
}

// ExistsByEmail checks if an employee with the email exists
func (r *EmployeeRepository) ExistsByEmail(email string) bool {
	var count int64
	r.db.Model(&models.Employee{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// Delete soft deletes an employee
func (r *EmployeeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Employee{}, id).Error
}

// FindByName finds an employee by full name (first_name + last_name)
// The name parameter can be "FirstName LastName" or "LastName, FirstName" or just "FirstName LastName"
func (r *EmployeeRepository) FindByName(name string) (*models.Employee, error) {
	var employee models.Employee
	
	// Check if it's in "LastName, FirstName" format
	if strings.Contains(name, ",") {
		parts := strings.Split(name, ",")
		if len(parts) == 2 {
			lastName := strings.TrimSpace(parts[0])
			firstName := strings.TrimSpace(parts[1])
			err := r.db.Where("LOWER(first_name) = LOWER(?) AND LOWER(last_name) = LOWER(?)", firstName, lastName).First(&employee).Error
			if err == nil {
				return &employee, nil
			}
		}
	}
	
	// Try "FirstName LastName" format
	parts := strings.Fields(name)
	if len(parts) >= 2 {
		firstName := parts[0]
		lastName := strings.Join(parts[1:], " ")
		err := r.db.Where("LOWER(first_name) = LOWER(?) AND LOWER(last_name) = LOWER(?)", firstName, lastName).First(&employee).Error
		if err == nil {
			return &employee, nil
		}
	}
	
	// Try matching just first name or last name if only one part
	if len(parts) == 1 {
		searchTerm := strings.ToLower(parts[0])
		err := r.db.Where("LOWER(first_name) = ? OR LOWER(last_name) = ?", searchTerm, searchTerm).First(&employee).Error
		if err == nil {
			return &employee, nil
		}
	}
	
	return nil, fmt.Errorf("employee not found with name: %s", name)
}

// ListManagers lists all active employees who can be reporting managers
func (r *EmployeeRepository) ListManagers(tenantID *uint) ([]models.Employee, error) {
	var employees []models.Employee
	query := r.db.Where("status = ? AND is_active = ?", models.StatusActive, true)
	
	// Tenant filter removed (single-tenant)
	
	err := query.Order("first_name ASC, last_name ASC").Find(&employees).Error
	return employees, err
}

// ListEmployeesWithoutUser returns active employees that don't have a linked user account yet.
// Used by the "Create User" flow to pick an employee.
func (r *EmployeeRepository) ListEmployeesWithoutUser(page, pageSize int, search string) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.Model(&models.Employee{}).
		Where("(user_id IS NULL OR user_id = 0)").
		Where("status = ? AND is_active = ?", models.StatusActive, true)

	if search != "" {
		term := "%" + search + "%"
		query = query.Where(
			"(LOWER(first_name) LIKE LOWER(?) OR LOWER(last_name) LIKE LOWER(?) OR LOWER(employee_id) LIKE LOWER(?) OR LOWER(work_email) LIKE LOWER(?))",
			term, term, term, term,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("first_name ASC, last_name ASC").Offset(offset).Limit(pageSize).Find(&employees).Error
	return employees, total, err
}

// ListByManagerID returns employees who report to the given manager
func (r *EmployeeRepository) ListByManagerID(managerID uint, page, pageSize int) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.Model(&models.Employee{}).
		Where("manager_id = ? AND status = ? AND is_active = ?", managerID, models.StatusActive, true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("first_name ASC, last_name ASC").Offset(offset).Limit(pageSize).Find(&employees).Error
	return employees, total, err
}

// SearchEmployees searches employees with filters and pagination
func (r *EmployeeRepository) SearchEmployees(tenantID *uint, search *string, departmentID, positionID, locationID *uint, status *string, page, pageSize int) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Employee{})

	// Tenant filter removed (single-tenant)

	// Search filter
	if search != nil && *search != "" {
		searchTerm := "%" + *search + "%"
		query = query.Where("(LOWER(first_name) LIKE LOWER(?) OR LOWER(last_name) LIKE LOWER(?) OR LOWER(employee_id) LIKE LOWER(?) OR LOWER(work_email) LIKE LOWER(?))",
			searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// Department filter
	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	// Position filter
	if positionID != nil {
		query = query.Where("position_id = ?", *positionID)
	}

	// Location filter
	if locationID != nil {
		query = query.Where("location_id = ?", *locationID)
	}

	// Status filter
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	} else {
		// Default to active
		query = query.Where("status = ? AND is_active = ?", models.StatusActive, true)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("first_name ASC, last_name ASC").Offset(offset).Limit(pageSize).Find(&employees).Error

	return employees, total, err
}

// PeopleDirectoryFilters holds all filter options for the People Directory search
type PeopleDirectoryFilters struct {
	Search         *string
	DepartmentID   *uint
	PositionID     *uint
	LocationID     *uint
	TeamID         *uint
	ManagerID      *uint
	EmploymentType *string
	Status         *string
	Letter         *string // First letter of first name (A-Z)
	SortBy         string  // first_name, last_name, employee_id, department, hire_date
	SortOrder      string  // asc, desc
	Page           int
	PageSize       int
}

// SearchPeopleDirectory searches employees for the People Directory with comprehensive filters
func (r *EmployeeRepository) SearchPeopleDirectory(filters PeopleDirectoryFilters) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	offset := (filters.Page - 1) * filters.PageSize
	query := r.db.Model(&models.Employee{})

	// Tenant filter removed (single-tenant)

	// Text search (name, employee_id, email, phone)
	if filters.Search != nil && *filters.Search != "" {
		searchTerm := "%" + *filters.Search + "%"
		query = query.Where(
			"(LOWER(first_name) LIKE LOWER(?) OR LOWER(last_name) LIKE LOWER(?) OR LOWER(CONCAT(first_name, ' ', last_name)) LIKE LOWER(?) OR LOWER(employee_id) LIKE LOWER(?) OR LOWER(work_email) LIKE LOWER(?) OR LOWER(phone_number) LIKE LOWER(?))",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// Alphabet letter filter
	if filters.Letter != nil && *filters.Letter != "" {
		query = query.Where("UPPER(LEFT(first_name, 1)) = UPPER(?)", *filters.Letter)
	}

	// Department filter
	if filters.DepartmentID != nil {
		query = query.Where("department_id = ?", *filters.DepartmentID)
	}

	// Position filter
	if filters.PositionID != nil {
		query = query.Where("position_id = ?", *filters.PositionID)
	}

	// Location filter
	if filters.LocationID != nil {
		query = query.Where("location_id = ?", *filters.LocationID)
	}

	// Team filter
	if filters.TeamID != nil {
		query = query.Where("team_id = ?", *filters.TeamID)
	}

	// Manager (reports_to) filter
	if filters.ManagerID != nil {
		query = query.Where("reports_to_id = ?", *filters.ManagerID)
	}

	// Employment type filter
	if filters.EmploymentType != nil && *filters.EmploymentType != "" {
		query = query.Where("employment_type = ?", *filters.EmploymentType)
	}

	// Status filter (default to active)
	if filters.Status != nil && *filters.Status != "" {
		query = query.Where("status = ?", *filters.Status)
	} else {
		query = query.Where("status = ? AND is_active = ?", models.StatusActive, true)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	orderClause := "first_name ASC, last_name ASC" // default
	validSortFields := map[string]string{
		"first_name":  "first_name",
		"last_name":   "last_name",
		"employee_id": "employee_id",
		"hire_date":   "hire_date",
		"created_at":  "created_at",
		"department":  "department_id",
	}
	if sortField, ok := validSortFields[filters.SortBy]; ok {
		direction := "ASC"
		if filters.SortOrder == "desc" {
			direction = "DESC"
		}
		orderClause = sortField + " " + direction
	}

	// Get paginated results
	err := query.Order(orderClause).Offset(offset).Limit(filters.PageSize).Find(&employees).Error

	return employees, total, err
}

// GetDepartmentEmployeeCounts returns the count of active employees per department
func (r *EmployeeRepository) GetDepartmentEmployeeCounts(tenantID *uint) ([]DepartmentEmployeeCount, error) {
	var counts []DepartmentEmployeeCount

	query := r.db.Model(&models.Employee{}).
		Select("department_id, COUNT(*) as count").
		Where("status = ? AND is_active = ? AND department_id IS NOT NULL", models.StatusActive, true).
		Group("department_id")

	// Tenant filter removed (single-tenant)

	err := query.Find(&counts).Error
	return counts, err
}

// DepartmentEmployeeCount holds department_id and employee count
type DepartmentEmployeeCount struct {
	DepartmentID uint  `json:"department_id"`
	Count        int64 `json:"count"`
}

// GetLocationEmployeeCounts returns the count of active employees per location
func (r *EmployeeRepository) GetLocationEmployeeCounts(tenantID *uint) ([]LocationEmployeeCount, error) {
	var counts []LocationEmployeeCount

	query := r.db.Model(&models.Employee{}).
		Select("location_id, COUNT(*) as count").
		Where("status = ? AND is_active = ? AND location_id IS NOT NULL", models.StatusActive, true).
		Group("location_id")

	// Tenant filter removed (single-tenant)

	err := query.Find(&counts).Error
	return counts, err
}

// LocationEmployeeCount holds location_id and employee count
type LocationEmployeeCount struct {
	LocationID uint  `json:"location_id"`
	Count      int64 `json:"count"`
}

// GetEmploymentTypeCounts returns the count of active employees per employment type
func (r *EmployeeRepository) GetEmploymentTypeCounts(tenantID *uint) ([]EmploymentTypeCount, error) {
	var counts []EmploymentTypeCount

	query := r.db.Model(&models.Employee{}).
		Select("employment_type, COUNT(*) as count").
		Where("status = ? AND is_active = ? AND employment_type IS NOT NULL", models.StatusActive, true).
		Group("employment_type")

	// Tenant filter removed (single-tenant)

	err := query.Find(&counts).Error
	return counts, err
}

// EmploymentTypeCount holds employment_type and count
type EmploymentTypeCount struct {
	EmploymentType string `json:"employment_type"`
	Count          int64  `json:"count"`
}

// GetAlphabetCounts returns the count of active employees grouped by first letter of first_name
func (r *EmployeeRepository) GetAlphabetCounts(tenantID *uint) ([]AlphabetCount, error) {
	var counts []AlphabetCount

	query := r.db.Model(&models.Employee{}).
		Select("UPPER(LEFT(first_name, 1)) as letter, COUNT(*) as count").
		Where("status = ? AND is_active = ?", models.StatusActive, true).
		Group("UPPER(LEFT(first_name, 1))").
		Order("letter ASC")

	// Tenant filter removed (single-tenant)

	err := query.Find(&counts).Error
	return counts, err
}

// AlphabetCount holds letter and count
type AlphabetCount struct {
	Letter string `json:"letter"`
	Count  int64  `json:"count"`
}

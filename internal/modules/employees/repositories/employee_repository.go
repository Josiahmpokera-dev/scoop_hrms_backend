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
	
	// Filter by tenant if provided
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	err := query.Order("first_name ASC, last_name ASC").Find(&employees).Error
	return employees, err
}

// SearchEmployees searches employees with filters and pagination
func (r *EmployeeRepository) SearchEmployees(tenantID *uint, search *string, departmentID, positionID, locationID *uint, status *string, page, pageSize int) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Employee{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

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

package services

import (
	"errors"
	"fmt"
	"time"

	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	locationRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/repositories"
	positionRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/email"
)

// EmployeeService handles employee business logic
type EmployeeService struct {
	employeeRepo   *employeeRepos.EmployeeRepository
	departmentRepo *departmentRepos.DepartmentRepository
	positionRepo   *positionRepos.JobPositionRepository
	locationRepo   *locationRepos.LocationRepository
	emailService   *email.EmailService
}

// NewEmployeeService creates a new employee service
func NewEmployeeService() *EmployeeService {
	return &EmployeeService{
		employeeRepo:   employeeRepos.NewEmployeeRepository(),
		departmentRepo: departmentRepos.NewDepartmentRepository(),
		positionRepo:   positionRepos.NewJobPositionRepository(),
		locationRepo:   locationRepos.NewLocationRepository(),
		emailService:   email.NewEmailService(),
	}
}

// OnboardEmployee creates a new employee (onboarding)
func (s *EmployeeService) OnboardEmployee(req *models.OnboardEmployeeRequest) (*models.Employee, error) {
	// Check if employee ID already exists
	if s.employeeRepo.ExistsByEmployeeID(req.EmployeeID) {
		return nil, errors.New("employee with this employee ID already exists")
	}

	// Check if email already exists
	if req.Email != "" && s.employeeRepo.ExistsByEmail(req.Email) {
		return nil, errors.New("employee with this email already exists")
	}

	// Validate department if provided and extract organization_id and organization_unit_id
	// NOTE: organization_id and organization_unit_id are automatically derived from the department
	// No need to require them in the request - they come from the department's organization_id
	var organizationID *uint
	var organizationUnitID *uint
	if req.DepartmentID != nil {
		department, err := s.departmentRepo.FindByID(*req.DepartmentID)
		if err != nil {
			return nil, errors.New("department not found")
		}
		// Automatically set organization_id from department (department holds the organization_id)
		organizationID = department.OrganizationID
		// Automatically set organization_unit_id from department if available
		organizationUnitID = department.OrganizationUnitID
	}

	// Validate position if provided
	if req.PositionID != nil {
		_, err := s.positionRepo.FindByID(*req.PositionID)
		if err != nil {
			return nil, errors.New("job position not found")
		}
	}

	// Set default hire date to today if not provided
	hireDate := req.HireDate
	if hireDate == nil {
		now := time.Now()
		hireDate = &now
	}

	// Set default employment type if not provided
	employmentType := req.EmploymentType
	if employmentType == "" {
		employmentType = "full-time"
	}

	// Set default currency if not provided
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	// Create employee
	// Map Gender to pointer
	var genderPtr *string
	if req.Gender != "" {
		genderPtr = &req.Gender
	}

	// Map EmploymentType to pointer
	var employmentTypePtr *string
	if employmentType != "" {
		employmentTypePtr = &employmentType
	}

	// Map Email to WorkEmail
	var workEmailPtr *string
	if req.Email != "" {
		workEmailPtr = &req.Email
	}

	// Map Phone to PhoneNumber
	var phoneNumberPtr *string
	if req.Phone != "" {
		phoneNumberPtr = &req.Phone
	}

	employee := &models.Employee{
		EmployeeID:               req.EmployeeID,
		FirstName:                req.FirstName,
		LastName:                 req.LastName,
		WorkEmail:                workEmailPtr,
		PhoneNumber:              phoneNumberPtr,
		DateOfBirth:              req.DateOfBirth,
		Gender:                   genderPtr,
		OrganizationID:           organizationID,     // Auto-populated from department
		OrganizationUnitID:       organizationUnitID, // Auto-populated from department
		DepartmentID:             req.DepartmentID,
		PositionID:               req.PositionID,
		LocationID:               req.LocationID,
		HireDate:                 hireDate,
		EmploymentType:           employmentTypePtr,
		Status:                   models.StatusActive,
		Salary:                   req.Salary,
		Currency:                 currency,
		EmergencyContactName:     req.EmergencyContactName,
		EmergencyContactPhone:    req.EmergencyContactPhone,
		EmergencyContactRelation: req.EmergencyContactRelation,
		Notes:                    req.Notes,
		IsActive:                 true,
	}

	if err := s.employeeRepo.Create(employee); err != nil {
		return nil, errors.New("failed to create employee")
	}

	return employee, nil
}

// SendCredentialsEmail sends login credentials to employee's personal email
func (s *EmployeeService) SendCredentialsEmail(employeeID uint) error {
	// Get employee by ID
	employee, err := s.employeeRepo.FindByID(employeeID)
	if err != nil {
		return errors.New("employee not found")
	}

	// Check if employee has personal email
	if employee.PersonalEmail == nil || *employee.PersonalEmail == "" {
		return errors.New("employee does not have a personal email address")
	}

	// Get user credentials for the employee
	userCredentials, err := s.getUserCredentialsForEmployee(employee)
	if err != nil {
		return fmt.Errorf("failed to get user credentials: %w", err)
	}

	// Send email with credentials
	return s.sendCredentialsEmail(*employee.PersonalEmail, employee, userCredentials)
}

// getUserCredentialsForEmployee retrieves user credentials for an employee
func (s *EmployeeService) getUserCredentialsForEmployee(employee *models.Employee) (*UserCredentials, error) {
	// This is a simplified implementation - in a real scenario, you would
	// query the user repository to get the actual credentials
	// For now, we'll return a mock response

	if employee.UserID == nil {
		return nil, errors.New("employee does not have a user account")
	}

	// In a real implementation, you would fetch from user repository
	// For demonstration, we'll return mock credentials
	return &UserCredentials{
		Email:    *employee.PersonalEmail,
		Username: fmt.Sprintf("%s.%s", employee.FirstName, employee.LastName),
		Password: "TemporaryPassword123!", // This should be fetched from user table
	}, nil
}

// sendCredentialsEmail sends the credentials email
func (s *EmployeeService) sendCredentialsEmail(recipientEmail string, employee *models.Employee, credentials *UserCredentials) error {
	// Build email subject and body
	subject := fmt.Sprintf("Your ScoopWorks Login Credentials - %s", employee.EmployeeID)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>ScoopWorks Login Credentials</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 5px; }
        .content { background-color: #fff; padding: 30px; border-radius: 5px; margin-top: 20px; border: 1px solid #e9ecef; }
        .credentials { background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .credential-item { margin: 10px 0; }
        .label { font-weight: bold; color: #495057; }
        .value { background-color: #fff; padding: 8px 12px; border: 1px solid #dee2e6; border-radius: 3px; font-family: monospace; }
        .footer { margin-top: 30px; text-align: center; color: #6c757d; font-size: 14px; }
        .warning { background-color: #fff3cd; border: 1px solid #ffeaa7; color: #856404; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>ScoopWorks HRMS</h1>
            <p>Your Login Credentials</p>
        </div>
        
        <div class="content">
            <h2>Hello %s %s!</h2>
            
            <p>Welcome to ScoopWorks! Your employee account has been created successfully.</p>
            
            <p>Here are your login credentials:</p>
            
            <div class="credentials">
                <div class="credential-item">
                    <span class="label">Employee ID:</span>
                    <div class="value">%s</div>
                </div>
                <div class="credential-item">
                    <span class="label">Email/Username:</span>
                    <div class="value">%s</div>
                </div>
                <div class="credential-item">
                    <span class="label">Password:</span>
                    <div class="value">%s</div>
                </div>
            </div>
            
            <div class="warning">
                <strong>Important Security Notice:</strong>
                <p>For security reasons, we recommend that you change your password immediately after your first login.</p>
            </div>
            
            <p><strong>Login URL:</strong> https://hrm.scoopworks.co.tz</p>
            
            <p>If you have any issues logging in, please contact the HR department.</p>
            
            <p>Best regards,<br>
            <strong>ScoopWorks HR Team</strong></p>
        </div>
        
        <div class="footer">
            <p>This is an automated message. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`,
		employee.FirstName,
		employee.LastName,
		employee.EmployeeID,
		credentials.Email,
		credentials.Password)

	// Send email
	return s.emailService.SendEmail(recipientEmail, subject, body)
}

// GetEmployeeByID retrieves an employee by ID
func (s *EmployeeService) GetEmployeeByID(id uint) (*models.Employee, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}
	return employee, nil
}

// GetEmployeeByEmployeeID retrieves an employee by employee ID
func (s *EmployeeService) GetEmployeeByEmployeeID(employeeID string) (*models.Employee, error) {
	employee, err := s.employeeRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}
	return employee, nil
}

// UpdateEmployee updates an employee
func (s *EmployeeService) UpdateEmployee(id uint, req *models.UpdateEmployeeRequest) (*models.Employee, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Update fields if provided
	if req.FirstName != nil {
		employee.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		employee.LastName = *req.LastName
	}
	if req.Email != nil {
		// Check if email is already taken by another employee
		if existing, _ := s.employeeRepo.FindByEmail(*req.Email); existing != nil && existing.ID != id {
			return nil, errors.New("email already in use by another employee")
		}
		employee.WorkEmail = req.Email
	}
	if req.Phone != nil {
		employee.PhoneNumber = req.Phone
	}
	if req.DateOfBirth != nil {
		employee.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != nil {
		employee.Gender = req.Gender
	}
	// Note: Address, City, State, Country, PostalCode are now in employee_addresses table
	// These fields are kept for backward compatibility but should be migrated to employee_addresses
	if req.DepartmentID != nil {
		// Validate department exists and extract organization_id and organization_unit_id
		// NOTE: When department is updated, organization_id and organization_unit_id are automatically updated
		// from the department (department holds the organization_id)
		department, err := s.departmentRepo.FindByID(*req.DepartmentID)
		if err != nil {
			return nil, errors.New("department not found")
		}
		employee.DepartmentID = req.DepartmentID
		// Automatically update organization_id from department
		employee.OrganizationID = department.OrganizationID
		// Automatically update organization_unit_id from department if available
		employee.OrganizationUnitID = department.OrganizationUnitID
	}
	if req.PositionID != nil {
		// Validate position exists
		_, err := s.positionRepo.FindByID(*req.PositionID)
		if err != nil {
			return nil, errors.New("job position not found")
		}
		employee.PositionID = req.PositionID
	}
	if req.LocationID != nil {
		// Validate location exists
		_, err := s.locationRepo.FindByID(*req.LocationID)
		if err != nil {
			return nil, errors.New("location not found")
		}
		employee.LocationID = req.LocationID
	}
	if req.EmploymentType != nil {
		employee.EmploymentType = req.EmploymentType
	}
	if req.Status != nil {
		employee.Status = models.EmployeeStatus(*req.Status)
	}
	if req.Salary != nil {
		employee.Salary = req.Salary
	}
	if req.Currency != nil {
		employee.Currency = *req.Currency
	}
	if req.EmergencyContactName != nil {
		employee.EmergencyContactName = *req.EmergencyContactName
	}
	if req.EmergencyContactPhone != nil {
		employee.EmergencyContactPhone = *req.EmergencyContactPhone
	}
	if req.EmergencyContactRelation != nil {
		employee.EmergencyContactRelation = *req.EmergencyContactRelation
	}
	if req.Notes != nil {
		employee.Notes = *req.Notes
	}
	if req.IsActive != nil {
		employee.IsActive = *req.IsActive
	}

	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, errors.New("failed to update employee")
	}

	return employee, nil
}

// ListEmployeesWithDetails lists employees with full details including related data and onboarding status
func (s *EmployeeService) ListEmployeesWithDetails(page, pageSize int) ([]models.EmployeeListResponse, int64, int, error) {
	employees, total, totalPages, err := s.ListEmployees(page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}

	// Convert to list response with related data
	responses := make([]models.EmployeeListResponse, len(employees))
	for i := range employees {
		response := s.convertToEmployeeListResponse(&employees[i])
		responses[i] = *response
	}

	return responses, total, totalPages, nil
}

// convertToEmployeeListResponse converts Employee to EmployeeListResponse with related data
func (s *EmployeeService) convertToEmployeeListResponse(emp *models.Employee) *models.EmployeeListResponse {
	response := &models.EmployeeListResponse{
		ID:            emp.ID,
		FullName:      emp.FullName(),
		EmployeeID:    emp.EmployeeID,
		Status:        string(emp.Status),
		DateOfJoining: emp.HireDate,
		IsActive:      emp.IsActive,
		CreatedAt:     emp.CreatedAt,
		UpdatedAt:     emp.UpdatedAt,
	}

	// Set email (prefer work email, fallback to personal email)
	if emp.WorkEmail != nil {
		response.Email = emp.WorkEmail
	} else if emp.PersonalEmail != nil {
		response.Email = emp.PersonalEmail
	}

	// Load department name
	if emp.DepartmentID != nil {
		response.DepartmentID = emp.DepartmentID
		department, err := s.departmentRepo.FindByID(*emp.DepartmentID)
		if err == nil && department != nil {
			response.Department = &department.Name
		}
	}

	// Load position/designation name
	if emp.PositionID != nil {
		response.PositionID = emp.PositionID
		position, err := s.positionRepo.FindByID(*emp.PositionID)
		if err == nil && position != nil {
			response.Designation = &position.Title
		}
	}

	// Load location name
	if emp.LocationID != nil {
		response.LocationID = emp.LocationID
		location, err := s.locationRepo.FindByID(*emp.LocationID)
		if err == nil && location != nil {
			response.Location = &location.Name
		}
	}

	// Load manager name
	if emp.ReportsToID != nil {
		response.ManagerID = emp.ReportsToID
		manager, err := s.employeeRepo.FindByID(*emp.ReportsToID)
		if err == nil && manager != nil {
			managerName := manager.FullName()
			response.Manager = &managerName
		}
	}

	// Check onboarding status
	onboardingStatus, progress := s.getOnboardingStatus(emp.ID)
	response.OnboardingStatus = onboardingStatus
	if progress != nil {
		response.OnboardingProgress = progress
	}

	return response
}

// ListManagers lists all active employees who can be reporting managers
func (s *EmployeeService) ListManagers(tenantID *uint, departmentID *uint) ([]map[string]interface{}, error) {
	managers, err := s.employeeRepo.ListManagers(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get managers: %w", err)
	}

	// If a departmentID filter is given, find the department head so we can flag it
	var departmentHeadID *uint
	if departmentID != nil {
		dept, err := s.departmentRepo.FindByID(*departmentID)
		if err == nil && dept != nil && dept.ManagerID != nil {
			departmentHeadID = dept.ManagerID
		}
	}

	var result []map[string]interface{}
	for _, emp := range managers {
		// If filtering by department, only include managers in that department
		if departmentID != nil {
			if emp.DepartmentID == nil || *emp.DepartmentID != *departmentID {
				continue
			}
		}

		managerData := map[string]interface{}{
			"id":          emp.ID,
			"employee_id": emp.EmployeeID,
			"full_name":   emp.FullName(),
			"first_name":  emp.FirstName,
			"last_name":   emp.LastName,
			"email":       emp.WorkEmail,
		}

		// Add department if available
		if emp.DepartmentID != nil {
			managerData["department_id"] = emp.DepartmentID
			department, err := s.departmentRepo.FindByID(*emp.DepartmentID)
			if err == nil && department != nil {
				managerData["department"] = department.Name
			}
		}

		// Add position if available
		if emp.PositionID != nil {
			managerData["position_id"] = emp.PositionID
			position, err := s.positionRepo.FindByID(*emp.PositionID)
			if err == nil && position != nil {
				managerData["position"] = position.Title
			}
		}

		// Flag the department head (suggested reporting manager)
		if departmentHeadID != nil && emp.ID == *departmentHeadID {
			managerData["is_department_head"] = true
			managerData["is_suggested"] = true
		} else {
			managerData["is_department_head"] = false
			managerData["is_suggested"] = false
		}

		result = append(result, managerData)
	}

	// If no results but we do have a department head, include the head even if not in the filtered list
	// (they might be from a different department managing this one)
	if departmentID != nil && len(result) == 0 && departmentHeadID != nil {
		headEmp, err := s.employeeRepo.FindByID(*departmentHeadID)
		if err == nil && headEmp != nil && headEmp.IsActive {
			managerData := map[string]interface{}{
				"id":                 headEmp.ID,
				"employee_id":        headEmp.EmployeeID,
				"full_name":          headEmp.FullName(),
				"first_name":         headEmp.FirstName,
				"last_name":          headEmp.LastName,
				"email":              headEmp.WorkEmail,
				"is_department_head": true,
				"is_suggested":       true,
			}
			if headEmp.DepartmentID != nil {
				managerData["department_id"] = headEmp.DepartmentID
				department, err := s.departmentRepo.FindByID(*headEmp.DepartmentID)
				if err == nil && department != nil {
					managerData["department"] = department.Name
				}
			}
			if headEmp.PositionID != nil {
				managerData["position_id"] = headEmp.PositionID
				position, err := s.positionRepo.FindByID(*headEmp.PositionID)
				if err == nil && position != nil {
					managerData["position"] = position.Title
				}
			}
			result = append(result, managerData)
		}
	}

	if result == nil {
		result = []map[string]interface{}{}
	}

	return result, nil
}

// GetDepartmentManager returns the suggested reporting manager for a department
// (the department head/manager). Returns nil if no manager is set.
func (s *EmployeeService) GetDepartmentManager(departmentID uint) (map[string]interface{}, error) {
	department, err := s.departmentRepo.FindByID(departmentID)
	if err != nil {
		return nil, fmt.Errorf("department not found")
	}

	if department.ManagerID == nil {
		return nil, nil // No manager set for this department
	}

	manager, err := s.employeeRepo.FindByID(*department.ManagerID)
	if err != nil || manager == nil {
		return nil, nil // Manager reference invalid
	}

	if manager.Status != models.StatusActive || !manager.IsActive {
		return nil, nil // Manager not active
	}

	result := map[string]interface{}{
		"id":                 manager.ID,
		"employee_id":        manager.EmployeeID,
		"full_name":          manager.FullName(),
		"first_name":         manager.FirstName,
		"last_name":          manager.LastName,
		"email":              manager.WorkEmail,
		"is_department_head": true,
		"department_id":      departmentID,
		"department":         department.Name,
	}

	if manager.PositionID != nil {
		position, err := s.positionRepo.FindByID(*manager.PositionID)
		if err == nil && position != nil {
			result["position_id"] = manager.PositionID
			result["position"] = position.Title
		}
	}

	return result, nil
}

// getOnboardingStatus checks if employee has an incomplete onboarding draft
func (s *EmployeeService) getOnboardingStatus(employeeID uint) (string, *float64) {
	// Get the employee to check their employee_id string and tenant_id
	employee, err := s.employeeRepo.FindByID(employeeID)
	if err != nil {
		return "completed", nil
	}

	draftRepo := employeeRepos.NewOnboardingDraftRepository()

	// First, check if there's a draft with employee_id_final matching this employee
	draft, err := draftRepo.FindByEmployeeIDFinal(employeeID)
	if err == nil && draft != nil {
		if !draft.IsCompleted {
			progress := draft.Progress
			return "in_progress", &progress
		}
		return "completed", nil
	}

	// Also check for drafts with employee_id string matching this employee's employee_id
	// This handles cases where draft was created with employee_id but not yet completed
	if employee.EmployeeID != "" {
		draft, err := draftRepo.FindByEmployeeIDString(employee.EmployeeID, nil)
		if err == nil && draft != nil {
			if !draft.IsCompleted {
				progress := draft.Progress
				return "in_progress", &progress
			}
			return "completed", nil
		}
	}

	// If no draft found, employee onboarding is completed (employee record exists)
	return "completed", nil
}

// ListEmployees lists employees with pagination
func (s *EmployeeService) ListEmployees(page, pageSize int) ([]models.Employee, int64, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	employees, total, err := s.employeeRepo.List(pageSize, offset)
	if err != nil {
		return nil, 0, 0, errors.New("failed to list employees")
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return employees, total, totalPages, nil
}

// ListByDepartment lists employees by department
func (s *EmployeeService) ListByDepartment(departmentID uint) ([]models.Employee, error) {
	return s.employeeRepo.ListByDepartment(departmentID)
}

// ListByStatus lists employees by status
func (s *EmployeeService) ListByStatus(status string) ([]models.Employee, error) {
	employeeStatus := models.EmployeeStatus(status)
	return s.employeeRepo.ListByStatus(employeeStatus)
}

// DeleteEmployee soft deletes an employee
func (s *EmployeeService) DeleteEmployee(id uint) error {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return errors.New("employee not found")
	}

	return s.employeeRepo.Delete(employee.ID)
}

// TerminateEmployee terminates an employee (sets status to terminated and is_active to false)
func (s *EmployeeService) TerminateEmployee(id uint, reason *string, terminationDate *time.Time, updatedBy *uint) (*models.Employee, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.Status == models.StatusTerminated {
		return nil, errors.New("employee is already terminated")
	}

	// Set status to terminated
	employee.Status = models.StatusTerminated
	employee.IsActive = false
	employee.UpdatedBy = updatedBy

	// Add termination reason to notes if provided
	if reason != nil && *reason != "" {
		if employee.Notes != "" {
			employee.Notes += "\n\nTermination Reason: " + *reason
		} else {
			employee.Notes = "Termination Reason: " + *reason
		}
	}

	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, errors.New("failed to terminate employee")
	}

	return employee, nil
}

// SuspendEmployee suspends an employee (sets status to suspended and is_active to false)
func (s *EmployeeService) SuspendEmployee(id uint, reason *string, suspensionEndDate *time.Time, updatedBy *uint) (*models.Employee, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.Status == models.StatusSuspended {
		return nil, errors.New("employee is already suspended")
	}

	if employee.Status == models.StatusTerminated {
		return nil, errors.New("cannot suspend a terminated employee")
	}

	// Set status to suspended
	employee.Status = models.StatusSuspended
	employee.IsActive = false
	employee.UpdatedBy = updatedBy

	// Add suspension reason to notes if provided
	if reason != nil && *reason != "" {
		suspensionNote := "Suspension Reason: " + *reason
		if suspensionEndDate != nil {
			suspensionNote += " (Suspension ends: " + suspensionEndDate.Format("2006-01-02") + ")"
		}
		if employee.Notes != "" {
			employee.Notes += "\n\n" + suspensionNote
		} else {
			employee.Notes = suspensionNote
		}
	}

	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, errors.New("failed to suspend employee")
	}

	return employee, nil
}

// ArchiveEmployee archives an employee (sets status to archived and is_active to false)
func (s *EmployeeService) ArchiveEmployee(id uint, reason *string, updatedBy *uint) (*models.Employee, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.Status == models.StatusArchived {
		return nil, errors.New("employee is already archived")
	}

	// Set status to archived
	employee.Status = models.StatusArchived
	employee.IsActive = false
	employee.UpdatedBy = updatedBy

	// Add archive reason to notes if provided
	if reason != nil && *reason != "" {
		if employee.Notes != "" {
			employee.Notes += "\n\nArchived Reason: " + *reason
		} else {
			employee.Notes = "Archived Reason: " + *reason
		}
	}

	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, errors.New("failed to archive employee")
	}

	return employee, nil
}

// ReactivateEmployee reactivates a suspended or archived employee
func (s *EmployeeService) ReactivateEmployee(id uint, updatedBy *uint) (*models.Employee, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.Status == models.StatusTerminated {
		return nil, errors.New("cannot reactivate a terminated employee")
	}

	if employee.Status == models.StatusActive && employee.IsActive {
		return nil, errors.New("employee is already active")
	}

	// Reactivate employee
	employee.Status = models.StatusActive
	employee.IsActive = true
	employee.UpdatedBy = updatedBy

	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, errors.New("failed to reactivate employee")
	}

	return employee, nil
}

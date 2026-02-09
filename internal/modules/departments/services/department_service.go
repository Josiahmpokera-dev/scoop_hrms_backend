package services

import (
	"errors"
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	locationRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/repositories"
	organizationRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/repositories"
)

type DepartmentService struct {
	repo            *repositories.DepartmentRepository
	organizationRepo *organizationRepos.OrganizationRepository
	locationRepo    *locationRepos.LocationRepository
	employeeRepo    *employeeRepos.EmployeeRepository
}

func NewDepartmentService() *DepartmentService {
	return &DepartmentService{
		repo:            repositories.NewDepartmentRepository(),
		organizationRepo: organizationRepos.NewOrganizationRepository(),
		locationRepo:    locationRepos.NewLocationRepository(),
		employeeRepo:    employeeRepos.NewEmployeeRepository(),
	}
}

// CreateDepartment creates a new department
func (s *DepartmentService) CreateDepartment(req *models.CreateDepartmentRequest, tenantID *uint, updatedBy *uint) (*models.Department, error) {
	// Check if code already exists
	if s.repo.ExistsByCode(req.Code) {
		return nil, errors.New("department with this code already exists")
	}

	// Validate organization if provided
	if req.OrganizationID != nil {
		_, err := s.organizationRepo.FindByID(*req.OrganizationID)
		if err != nil {
			return nil, errors.New("organization not found")
		}
	}

	// Validate parent department if provided
	var parentDepartment *models.Department
	if req.ParentDepartmentID != nil {
		parent, err := s.repo.FindByID(*req.ParentDepartmentID)
		if err != nil {
			return nil, errors.New("parent department not found")
		}
		parentDepartment = parent
	}

	// Handle location: accept either location_id (ID) or location (name)
	// Priority: location_id > location (if both provided, location_id takes precedence)
	var locationID *uint
	if req.LocationID != nil {
		// Validate location by ID
		location, err := s.locationRepo.FindByID(*req.LocationID)
		if err != nil {
			return nil, fmt.Errorf("location with ID '%d' not found", *req.LocationID)
		}
		// Verify location belongs to the tenant
		if tenantID != nil && location.TenantID != nil && *location.TenantID != *tenantID {
			return nil, fmt.Errorf("location with ID '%d' does not belong to your tenant", *req.LocationID)
		}
		locationID = &location.ID
	} else if req.Location != nil && *req.Location != "" {
		// Look up location by name if provided
		// Allow locations from any organization within the tenant (not restricted to department's organization)
		location, err := s.locationRepo.FindByName(*req.Location, tenantID, req.OrganizationID)
		if err != nil {
			// If not found in organization scope, try tenant-wide search
			if req.OrganizationID != nil {
				location, err = s.locationRepo.FindByName(*req.Location, tenantID, nil)
			}
			if err != nil {
				return nil, fmt.Errorf("location '%s' not found", *req.Location)
			}
		}
		locationID = &location.ID
	}

	// Cost center is stored as a string (not a reference)
	// No lookup needed - just use the string value directly

	// Validate manager (Head of Department) if provided
	if req.ManagerID != nil {
		_, err := s.employeeRepo.FindByID(*req.ManagerID)
		if err != nil {
			return nil, errors.New("manager (head of department) not found")
		}
	}

	// Deputy manager is stored as a string (not a reference)
	// No lookup needed - just use the string value directly

	// Determine level: if parent is provided, infer level from parent; otherwise use provided level or default to top level
	var level *string
	if req.Level != nil {
		level = req.Level
	} else if parentDepartment != nil && parentDepartment.Level != nil {
		// Infer level from parent (one level down)
		parentLevel := *parentDepartment.Level
		var inferredLevel string
		switch parentLevel {
		case "company":
			inferredLevel = "business_unit"
		case "business_unit":
			inferredLevel = "department"
		case "department":
			inferredLevel = "team"
		default:
			inferredLevel = "department" // Default fallback
		}
		level = &inferredLevel
	} else {
		// Default to top level (company) if no parent and no level provided
		defaultLevel := "company"
		level = &defaultLevel
	}

	// Set parent_department_id: if not provided, default to nil (top level)
	parentDepartmentID := req.ParentDepartmentID
	// If explicitly set to 0 or nil, it means top level (no parent)
	if req.ParentDepartmentID != nil && *req.ParentDepartmentID == 0 {
		parentDepartmentID = nil
	}

	department := &models.Department{
		TenantID:          tenantID,
		OrganizationID:    req.OrganizationID,
		OrganizationUnitID: req.OrganizationUnitID,
		Code:              req.Code,
		Name:              req.Name,
		Description:       req.Description,
		Level:             level,
		DepartmentType:    req.DepartmentType,
		ParentDepartmentID: parentDepartmentID,
		ManagerID:         req.ManagerID,
		DeputyManager:     req.DeputyManager, // Store as string directly
		BudgetAllocated:   req.BudgetAllocated,
		BudgetCurrency:   req.BudgetCurrency,
		EmployeeCapacity:  req.EmployeeCapacity,
		LocationID:       locationID,
		CostCenter:       req.CostCenter, // Store as string directly
		IsActive:         true,
		UpdatedBy:        updatedBy,
	}
	
	// Set default budget currency if not provided
	if department.BudgetCurrency == nil && department.BudgetAllocated != nil {
		defaultCurrency := "TZS"
		department.BudgetCurrency = &defaultCurrency
	}

	if req.IsActive != nil {
		department.IsActive = *req.IsActive
	}

	if err := s.repo.Create(department); err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	return department, nil
}

// GetDepartmentByID retrieves a department by ID
func (s *DepartmentService) GetDepartmentByID(id uint) (*models.Department, error) {
	return s.repo.FindByID(id)
}

// UpdateDepartment updates an existing department
func (s *DepartmentService) UpdateDepartment(id uint, req *models.UpdateDepartmentRequest, updatedBy *uint) (*models.Department, error) {
	department, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("department not found")
	}

	// Check code uniqueness if code is being updated
	if req.Code != nil && *req.Code != department.Code {
		if s.repo.ExistsByCode(*req.Code) {
			return nil, errors.New("department with this code already exists")
		}
		department.Code = *req.Code
	}

	// Validate organization if being updated
	if req.OrganizationID != nil {
		_, err := s.organizationRepo.FindByID(*req.OrganizationID)
		if err != nil {
			return nil, errors.New("organization not found")
		}
		department.OrganizationID = req.OrganizationID
	}
	if req.OrganizationUnitID != nil {
		department.OrganizationUnitID = req.OrganizationUnitID
	}

	if req.Name != nil {
		department.Name = *req.Name
	}
	if req.Description != nil {
		department.Description = req.Description
	}
	if req.DepartmentType != nil {
		department.DepartmentType = req.DepartmentType
	}
	// Deputy manager is stored as a string (not a reference)
	if req.DeputyManager != nil {
		department.DeputyManager = req.DeputyManager
	}
	if req.BudgetAllocated != nil {
		department.BudgetAllocated = req.BudgetAllocated
	}
	if req.BudgetCurrency != nil {
		department.BudgetCurrency = req.BudgetCurrency
	}
	if req.EmployeeCapacity != nil {
		department.EmployeeCapacity = req.EmployeeCapacity
	}
	// Handle location: accept either location_id (ID) or location (name)
	// Priority: location_id > location (if both provided, location_id takes precedence)
	if req.LocationID != nil && *req.LocationID != 0 {
		// Validate location by ID
		location, err := s.locationRepo.FindByID(*req.LocationID)
		if err != nil {
			return nil, fmt.Errorf("location with ID '%d' not found", *req.LocationID)
		}
		// Verify location belongs to the tenant
		if department.TenantID != nil && location.TenantID != nil && *location.TenantID != *department.TenantID {
			return nil, fmt.Errorf("location with ID '%d' does not belong to your tenant", *req.LocationID)
		}
		department.LocationID = &location.ID
	} else if req.Location != nil && *req.Location != "" {
		// Look up location by name
		// Allow locations from any organization within the tenant
		location, err := s.locationRepo.FindByName(*req.Location, department.TenantID, department.OrganizationID)
		if err != nil {
			// If not found in organization scope, try tenant-wide search
			if department.OrganizationID != nil {
				location, err = s.locationRepo.FindByName(*req.Location, department.TenantID, nil)
			}
			if err != nil {
				return nil, fmt.Errorf("location '%s' not found", *req.Location)
			}
		}
		department.LocationID = &location.ID
	} else if req.LocationID != nil && *req.LocationID == 0 {
		// Allow clearing location by setting location_id to 0
		department.LocationID = nil
	} else if req.Location != nil && *req.Location == "" {
		// Allow clearing location by setting location to empty string
		department.LocationID = nil
	}
	// Cost center is stored as a string (not a reference)
	if req.CostCenter != nil {
		department.CostCenter = req.CostCenter
	}
	if req.Level != nil {
		department.Level = req.Level
	}
	if req.ParentDepartmentID != nil {
		// Validate parent department
		if *req.ParentDepartmentID != 0 {
			parent, err := s.repo.FindByID(*req.ParentDepartmentID)
			if err != nil {
				return nil, errors.New("parent department not found")
			}
			// Prevent circular reference
			if parent.ID == id {
				return nil, errors.New("department cannot be its own parent")
			}
		}
		department.ParentDepartmentID = req.ParentDepartmentID
	}
	if req.ManagerID != nil {
		// Validate manager exists
		_, err := s.employeeRepo.FindByID(*req.ManagerID)
		if err != nil {
			return nil, errors.New("manager (head of department) not found")
		}
		department.ManagerID = req.ManagerID
	}
	if req.IsActive != nil {
		department.IsActive = *req.IsActive
	}

	department.UpdatedBy = updatedBy

	if err := s.repo.Update(department); err != nil {
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	return department, nil
}

// DeleteDepartment soft deletes a department
func (s *DepartmentService) DeleteDepartment(id uint) error {
	// Check if department has sub-departments
	department, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("department not found")
	}

	if len(department.SubDepartments) > 0 {
		return errors.New("cannot delete department with sub-departments")
	}

	return s.repo.Delete(id)
}

// ListDepartments lists departments with pagination and filters
func (s *DepartmentService) ListDepartments(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Department, int64, error) {
	return s.repo.List(tenantID, page, pageSize, filters)
}

// GetRootDepartments gets all root departments
func (s *DepartmentService) GetRootDepartments(tenantID *uint) ([]models.Department, error) {
	return s.repo.FindRootDepartments(tenantID)
}

// AssignDepartmentHead assigns an employee as the head of a department.
// This also updates the org chart by:
//   - Setting the employee's department_id to this department
//   - Setting the employee's reports_to_id to the parent department's head (if exists)
//   - Updating employees in this department who had no reporting manager to report to the new head
func (s *DepartmentService) AssignDepartmentHead(departmentID uint, employeeID uint, updatedBy *uint) (*models.Department, error) {
	department, err := s.repo.FindByID(departmentID)
	if err != nil {
		return nil, errors.New("department not found")
	}

	// Validate the employee exists and is active
	employee, err := s.employeeRepo.FindByID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}
	if !employee.IsActive {
		return nil, errors.New("employee is not active")
	}

	// Track the previous head (if any) for cascading updates
	previousHeadID := department.ManagerID

	// 1. Update the department's manager_id
	department.ManagerID = &employeeID
	department.UpdatedBy = updatedBy

	if err := s.repo.Update(department); err != nil {
		return nil, fmt.Errorf("failed to assign department head: %w", err)
	}

	// 2. Update the new head's department_id to this department (if not already)
	if employee.DepartmentID == nil || *employee.DepartmentID != departmentID {
		employee.DepartmentID = &departmentID
	}

	// 3. Set the new head's reporting manager to the parent department's head
	//    This places the new head correctly in the org chart hierarchy
	if department.ParentDepartmentID != nil {
		parentDept, err := s.repo.FindByID(*department.ParentDepartmentID)
		if err == nil && parentDept != nil && parentDept.ManagerID != nil {
			// New head reports to the parent department's head
			employee.ReportsToID = parentDept.ManagerID
		}
	} else {
		// Top-level department: the head reports to no one (or could report to CEO)
		// Clear reports_to if this is a root department
		employee.ReportsToID = nil
	}

	employee.UpdatedBy = updatedBy
	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, fmt.Errorf("failed to update department head employee record: %w", err)
	}

	// 4. Cascade: employees in this department who currently report to the old head
	//    should now report to the new head
	if previousHeadID != nil && *previousHeadID != employeeID {
		s.reassignSubordinates(*previousHeadID, employeeID, departmentID)
	}

	// 5. Cascade: employees in this department who have no reporting manager
	//    should be assigned to report to the new head
	s.assignOrphanedEmployees(employeeID, departmentID)

	// 6. Update child departments: if any child department's head reports to the old head,
	//    update them to report to the new head
	if previousHeadID != nil && *previousHeadID != employeeID {
		s.updateChildDepartmentHeads(departmentID, employeeID)
	} else {
		// Even without a previous head, ensure child department heads report to this head
		s.updateChildDepartmentHeads(departmentID, employeeID)
	}

	return department, nil
}

// RemoveDepartmentHead removes the head of department assignment.
// This updates the org chart by:
//   - Re-pointing employees who reported to the removed head to the parent department head (if exists)
//   - Clearing the department's manager_id
func (s *DepartmentService) RemoveDepartmentHead(departmentID uint, updatedBy *uint) (*models.Department, error) {
	department, err := s.repo.FindByID(departmentID)
	if err != nil {
		return nil, errors.New("department not found")
	}

	previousHeadID := department.ManagerID

	// Clear the department head
	department.ManagerID = nil
	department.UpdatedBy = updatedBy

	if err := s.repo.Update(department); err != nil {
		return nil, fmt.Errorf("failed to remove department head: %w", err)
	}

	// Cascade: reassign employees who reported to the removed head
	if previousHeadID != nil {
		// Find the fallback manager: parent department's head
		var fallbackManagerID *uint
		if department.ParentDepartmentID != nil {
			parentDept, err := s.repo.FindByID(*department.ParentDepartmentID)
			if err == nil && parentDept != nil && parentDept.ManagerID != nil {
				fallbackManagerID = parentDept.ManagerID
			}
		}

		// Reassign direct reports of the removed head within this department
		employees, _ := s.employeeRepo.ListByDepartment(departmentID)
		for _, emp := range employees {
			if emp.ReportsToID != nil && *emp.ReportsToID == *previousHeadID {
				emp.ReportsToID = fallbackManagerID
				emp.UpdatedBy = updatedBy
				_ = s.employeeRepo.Update(&emp)
			}
		}

		// Update child department heads: they should now report to the fallback manager
		if department.SubDepartments != nil {
			for _, childDept := range department.SubDepartments {
				if childDept.ManagerID != nil {
					childHead, err := s.employeeRepo.FindByID(*childDept.ManagerID)
					if err == nil && childHead != nil {
						childHead.ReportsToID = fallbackManagerID
						childHead.UpdatedBy = updatedBy
						_ = s.employeeRepo.Update(childHead)
					}
				}
			}
		}
	}

	return department, nil
}

// reassignSubordinates moves employees who reported to oldHeadID to report to newHeadID
// within the given department
func (s *DepartmentService) reassignSubordinates(oldHeadID, newHeadID uint, departmentID uint) {
	employees, err := s.employeeRepo.ListByDepartment(departmentID)
	if err != nil {
		return
	}

	for _, emp := range employees {
		if emp.ID == newHeadID {
			continue // Don't modify the new head itself
		}
		if emp.ReportsToID != nil && *emp.ReportsToID == oldHeadID {
			emp.ReportsToID = &newHeadID
			_ = s.employeeRepo.Update(&emp)
		}
	}
}

// assignOrphanedEmployees assigns employees in the department who have no reporting manager
// to report to the department head
func (s *DepartmentService) assignOrphanedEmployees(headEmployeeID uint, departmentID uint) {
	employees, err := s.employeeRepo.ListByDepartment(departmentID)
	if err != nil {
		return
	}

	for _, emp := range employees {
		if emp.ID == headEmployeeID {
			continue // Don't assign head to report to themselves
		}
		if emp.ReportsToID == nil && emp.IsActive {
			emp.ReportsToID = &headEmployeeID
			_ = s.employeeRepo.Update(&emp)
		}
	}
}

// updateChildDepartmentHeads ensures heads of child departments report to this department's head
func (s *DepartmentService) updateChildDepartmentHeads(departmentID uint, newHeadEmployeeID uint) {
	department, err := s.repo.FindByID(departmentID)
	if err != nil || department == nil {
		return
	}

	// SubDepartments are preloaded by FindByID
	for _, childDept := range department.SubDepartments {
		if childDept.ManagerID != nil && *childDept.ManagerID != newHeadEmployeeID {
			childHead, err := s.employeeRepo.FindByID(*childDept.ManagerID)
			if err == nil && childHead != nil && childHead.IsActive {
				childHead.ReportsToID = &newHeadEmployeeID
				_ = s.employeeRepo.Update(childHead)
			}
		}
	}
}

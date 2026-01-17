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

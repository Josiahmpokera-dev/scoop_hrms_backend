package services

import (
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization/models"
	positionRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/repositories"
	"gorm.io/gorm"
)

// OrgChartService handles organization chart business logic
type OrgChartService struct {
	positionRepo   *positionRepos.JobPositionRepository
	employeeRepo   *employeeRepos.EmployeeRepository
	departmentRepo *departmentRepos.DepartmentRepository
	db             *gorm.DB
}

// NewOrgChartService creates a new organization chart service
func NewOrgChartService() *OrgChartService {
	return &OrgChartService{
		positionRepo:   positionRepos.NewJobPositionRepository(),
		employeeRepo:   employeeRepos.NewEmployeeRepository(),
		departmentRepo: departmentRepos.NewDepartmentRepository(),
		db:             database.GetDB(),
	}
}

// GetOrganizationChart builds and returns the complete organization chart
func (s *OrgChartService) GetOrganizationChart(tenantID *uint) (*models.OrgChartNode, error) {
	// Find root positions (positions with no parent - ReportsToPositionID is NULL)
	var rootPositions []struct {
		ID                  uint
		Code                string
		Title               string
		DepartmentID        *uint
		ReportsToPositionID *uint
		IsActive            bool
		TenantID            *uint
	}

	query := s.db.Table("job_positions").
		Select("id, code, title, department_id, reports_to_position_id, is_active, tenant_id").
		Where("reports_to_position_id IS NULL").
		Where("is_active = ?", true).
		Where("deleted_at IS NULL")

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Find(&rootPositions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch root positions: %w", err)
	}

	if len(rootPositions) == 0 {
		// Provide more helpful error message with diagnostic info
		var totalPositions int64
		var rootPositionsCount int64
		var activeRootPositionsCount int64

		baseQuery := s.db.Table("job_positions").Where("deleted_at IS NULL")
		if tenantID != nil {
			baseQuery = baseQuery.Where("tenant_id = ?", *tenantID)
		}

		baseQuery.Count(&totalPositions)

		rootQuery := s.db.Table("job_positions").
			Where("reports_to_position_id IS NULL").
			Where("deleted_at IS NULL")
		if tenantID != nil {
			rootQuery = rootQuery.Where("tenant_id = ?", *tenantID)
		}
		rootQuery.Count(&rootPositionsCount)

		activeRootQuery := s.db.Table("job_positions").
			Where("reports_to_position_id IS NULL").
			Where("is_active = ?", true).
			Where("deleted_at IS NULL")
		if tenantID != nil {
			activeRootQuery = activeRootQuery.Where("tenant_id = ?", *tenantID)
		}
		activeRootQuery.Count(&activeRootPositionsCount)

		errorMsg := fmt.Sprintf("no root positions found. Diagnostic: total positions=%d, root positions (any status)=%d, active root positions=%d",
			totalPositions, rootPositionsCount, activeRootPositionsCount)

		if tenantID != nil {
			errorMsg += fmt.Sprintf(", tenant_id=%d", *tenantID)
		} else {
			errorMsg += ", tenant_id=null (checking all tenants)"
		}

		if rootPositionsCount > 0 && activeRootPositionsCount == 0 {
			errorMsg += ". Found root positions but they are inactive. Set is_active=true to make them available."
		} else if totalPositions > 0 && rootPositionsCount == 0 {
			errorMsg += ". Found positions but none are root positions (all have reports_to_position_id set). Create a position with reports_to_position_id=null."
		} else if totalPositions == 0 {
			errorMsg += ". No positions found in database. Create at least one position first."
		}

		return nil, fmt.Errorf("%s", errorMsg)
	}

	// For now, use the first root position as the CEO/top level
	// In a multi-organization setup, you might want to handle multiple roots
	rootPosition := rootPositions[0]

	// Build the chart starting from root
	rootNode, err := s.buildNode(rootPosition.ID, nil, 0, tenantID)
	if err != nil {
		return nil, err
	}

	return rootNode, nil
}

// buildNode recursively builds a node and its subordinates
func (s *OrgChartService) buildNode(positionID uint, managerEmpID *string, level int, tenantID *uint) (*models.OrgChartNode, error) {
	// Get position details
	position, err := s.positionRepo.FindByID(positionID)
	if err != nil {
		return nil, fmt.Errorf("position not found: %w", err)
	}

	// Find employee(s) in this position (there might be multiple, but we'll take the first active one)
	var employees []struct {
		ID          uint
		EmployeeID  string
		FirstName   string
		LastName    string
		WorkEmail   *string
		PhoneNumber *string
		PhotoURL    *string
		Status      string
		IsActive    bool
	}

	empQuery := s.db.Table("employees").
		Select("id, employee_id, first_name, last_name, work_email, phone_number, photo_url, status, is_active").
		Where("position_id = ?", positionID).
		Where("status = ?", "active").
		Where("is_active = ?", true).
		Where("deleted_at IS NULL")

	if tenantID != nil {
		empQuery = empQuery.Where("tenant_id = ?", *tenantID)
	}

	if err := empQuery.Find(&employees).Error; err != nil {
		// If error finding employees, continue with vacant position
		employees = []struct {
			ID          uint
			EmployeeID  string
			FirstName   string
			LastName    string
			WorkEmail   *string
			PhoneNumber *string
			PhotoURL    *string
			Status      string
			IsActive    bool
		}{}
	}

	// Get department name and ID
	var departmentName *string
	var departmentID *uint
	if position.DepartmentID != nil {
		dept, err := s.departmentRepo.FindByID(*position.DepartmentID)
		if err == nil && dept != nil {
			departmentName = &dept.Name
			departmentID = &dept.ID
		}
	}

	// Build node
	node := &models.OrgChartNode{
		Designation:  position.Title,
		Department:   departmentName,
		DepartmentID: departmentID,
		ManagerEmpID: managerEmpID,
		Level:        level,
		IsVacant:     len(employees) == 0,
		Subordinates: []models.OrgChartNode{},
	}

	// Fill employee data if exists
	if len(employees) > 0 {
		emp := employees[0] // Take first active employee
		empID := emp.EmployeeID
		fullName := fmt.Sprintf("%s %s", emp.FirstName, emp.LastName)

		node.ID = fmt.Sprintf("emp-%s", empID)
		node.EmpID = &empID
		node.Name = fullName
		node.Email = emp.WorkEmail
		node.Phone = emp.PhoneNumber
		if emp.PhotoURL != nil && *emp.PhotoURL != "" {
			node.Photo = emp.PhotoURL
		} else {
			defaultPhoto := "/img/avatars/default.jpg"
			node.Photo = &defaultPhoto
		}

		// Check if this employee is a department head
		// Query departments where this employee is the manager (head)
		var headDepartments []struct {
			ID   uint
			Name string
			Code string
		}
		s.db.Table("departments").
			Select("id, name, code").
			Where("manager_id = ? AND is_active = ? AND deleted_at IS NULL", emp.ID, true).
			Find(&headDepartments)

		if len(headDepartments) > 0 {
			node.IsDepartmentHead = true
			node.HeadOfDepartments = make([]models.DepartmentRef, len(headDepartments))
			for i, hd := range headDepartments {
				node.HeadOfDepartments[i] = models.DepartmentRef{
					ID:   hd.ID,
					Name: hd.Name,
					Code: hd.Code,
				}
			}
		}
	} else {
		// Vacant position
		node.ID = fmt.Sprintf("pos-%d", positionID)
		node.Name = "Open Position"
		defaultPhoto := "/img/avatars/default.jpg"
		node.Photo = &defaultPhoto
	}

	// Find child positions (positions that report to this position)
	var childPositions []struct {
		ID                  uint
		Code                string
		Title               string
		DepartmentID        *uint
		ReportsToPositionID *uint
		IsActive            bool
		TenantID            *uint
	}

	// Find child positions (positions that report to this position)
	// First, let's check all positions that report to this position (for debugging)
	var allChildPositions []struct {
		ID                  uint
		Code                string
		Title               string
		ReportsToPositionID *uint
		IsActive            bool
		TenantID            *uint
		DeletedAt           *string
	}

	// Query all positions reporting to this position (for diagnostic purposes)
	allChildQuery := s.db.Table("job_positions").
		Select("id, code, title, reports_to_position_id, is_active, tenant_id, deleted_at").
		Where("reports_to_position_id = ?", positionID)

	if tenantID != nil {
		allChildQuery = allChildQuery.Where("tenant_id = ?", *tenantID)
	}

	allChildQuery.Find(&allChildPositions)

	// Now query only active, non-deleted child positions
	childQuery := s.db.Table("job_positions").
		Select("id, code, title, department_id, reports_to_position_id, is_active, tenant_id").
		Where("reports_to_position_id = ?", positionID).
		Where("is_active = ?", true).
		Where("deleted_at IS NULL")

	if tenantID != nil {
		childQuery = childQuery.Where("tenant_id = ?", *tenantID)
	}

	if err := childQuery.Find(&childPositions).Error; err != nil {
		// If query fails, return node without children rather than failing completely
		// This allows partial chart building
		return node, nil
	}

	// If no active child positions found, but we have inactive ones, that's the issue
	if len(childPositions) == 0 && len(allChildPositions) > 0 {
		// We have child positions but they're inactive or deleted
		// Return node with empty subordinates (this is expected behavior)
		return node, nil
	}

	// If no child positions at all, this is a leaf node (normal)
	if len(childPositions) == 0 {
		return node, nil
	}

	// Recursively build child nodes
	for _, childPos := range childPositions {
		// Determine manager employee ID for child
		// If current node has an employee, use that employee's ID
		// Otherwise, pass nil (vacant position)
		childManagerEmpID := node.EmpID

		// Recursively build the child node (this will include its own subordinates)
		// This recursive call will:
		// 1. Find the child position details
		// 2. Find employees in that position
		// 3. Find child positions of this child (grandchildren)
		// 4. Recursively build those grandchildren
		// 5. Return the complete node with all nested subordinates
		childNode, err := s.buildNode(childPos.ID, childManagerEmpID, level+1, tenantID)
		if err != nil {
			// If building child fails, create a minimal node with error info
			// This ensures the hierarchy is still visible even if some nodes fail
			errorNode := models.OrgChartNode{
				ID:           fmt.Sprintf("pos-%d", childPos.ID),
				EmpID:        nil,
				Name:         "Error loading position",
				Designation:  childPos.Title,
				Department:   nil,
				Photo:        nil,
				Email:        nil,
				Phone:        nil,
				ManagerEmpID: childManagerEmpID,
				Level:        level + 1,
				IsVacant:     true,
				Subordinates: []models.OrgChartNode{},
			}
			node.Subordinates = append(node.Subordinates, errorNode)
			continue
		}

		// Append the child node (which already contains its own subordinates from the recursive call)
		// This builds the complete nested hierarchy - the childNode.Subordinates array
		// will contain all grandchildren, great-grandchildren, etc.
		node.Subordinates = append(node.Subordinates, *childNode)
	}

	return node, nil
}

package services

import (
	"errors"
	"fmt"
	"time"

	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/repositories"
)

// ProjectService handles business logic for projects
type ProjectService struct {
	projectRepo    *repositories.ProjectRepository
	taskRepo       *repositories.DailyTaskRepository
	employeeRepo   *employeeRepos.EmployeeRepository
	departmentRepo *departmentRepos.DepartmentRepository
}

// NewProjectService creates a new ProjectService
func NewProjectService() *ProjectService {
	return &ProjectService{
		projectRepo:    repositories.NewProjectRepository(),
		taskRepo:       repositories.NewDailyTaskRepository(),
		employeeRepo:   employeeRepos.NewEmployeeRepository(),
		departmentRepo: departmentRepos.NewDepartmentRepository(),
	}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(req *models.CreateProjectRequest, createdByID uint, createdByName string) (*models.Project, error) {
	// Validate department if provided
	var deptName *string
	if req.DepartmentID != nil {
		dept, err := s.departmentRepo.FindByID(*req.DepartmentID)
		if err != nil {
			return nil, errors.New("department not found")
		}
		name := dept.Name
		deptName = &name
	}

	project, err := s.buildProject(req, createdByID, createdByName, deptName)
	if err != nil {
		return nil, err
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Add initial members if provided
	s.addInitialMembers(project.ID, req.MemberIDs, createdByID)

	return s.projectRepo.FindByID(project.ID)
}

func (s *ProjectService) buildProject(req *models.CreateProjectRequest, createdByID uint, createdByName string, deptName *string) (*models.Project, error) {
	priority := models.ProjectPriorityMedium
	if req.Priority != "" {
		priority = req.Priority
	}

	// Parse date strings
	var startDate, endDate, dueDate *time.Time
	if req.StartDate != nil {
		d, err := models.ParseDate(*req.StartDate)
		if err != nil {
			return nil, err
		}
		startDate = d
	}
	if req.EndDate != nil {
		d, err := models.ParseDate(*req.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = d
	}
	if req.DueDate != nil {
		d, err := models.ParseDate(*req.DueDate)
		if err != nil {
			return nil, err
		}
		dueDate = d
	}

	return &models.Project{
		ProjectCode:    s.projectRepo.GetNextProjectCode(),
		Name:           req.Name,
		Description:    req.Description,
		Status:         models.ProjectStatusPlanning,
		Priority:       priority,
		DepartmentID:   req.DepartmentID,
		DepartmentName: deptName,
		Category:       req.Category,
		CreatedByID:    createdByID,
		CreatedByName:  createdByName,
		StartDate:      startDate,
		EndDate:        endDate,
		DueDate:        dueDate,
		EstimatedHours: req.EstimatedHours,
		Budget:         req.Budget,
		Notes:          req.Notes,
		IsActive:       true,
	}, nil
}

func (s *ProjectService) addInitialMembers(projectID uint, memberIDs []string, assignedByID uint) {
	for _, empID := range memberIDs {
		emp, err := s.employeeRepo.FindByEmployeeID(empID)
		if err != nil {
			continue
		}
		member := &models.ProjectMember{
			ProjectID:    projectID,
			EmployeeID:   empID,
			EmployeeName: emp.FirstName + " " + emp.LastName,
			Role:         models.MemberRoleMember,
			AssignedByID: &assignedByID,
			IsActive:     true,
		}
		_ = s.projectRepo.AddMember(member)
	}
}

// GetProject returns a project by ID
func (s *ProjectService) GetProject(id uint) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("project not found")
	}
	return project, nil
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(id uint, req *models.UpdateProjectRequest, updatedBy *uint) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = req.Description
	}
	if req.Status != nil {
		project.Status = *req.Status
		if *req.Status == models.ProjectStatusCompleted {
			now := time.Now()
			project.CompletedAt = &now
			project.Progress = 100
		}
	}
	if req.Priority != nil {
		project.Priority = *req.Priority
	}
	if req.DepartmentID != nil {
		dept, dErr := s.departmentRepo.FindByID(*req.DepartmentID)
		if dErr == nil {
			project.DepartmentID = req.DepartmentID
			deptName := dept.Name
			project.DepartmentName = &deptName
		}
	}
	if req.Category != nil {
		project.Category = req.Category
	}
	if req.StartDate != nil {
		d, pErr := models.ParseDate(*req.StartDate)
		if pErr != nil {
			return nil, pErr
		}
		project.StartDate = d
	}
	if req.EndDate != nil {
		d, pErr := models.ParseDate(*req.EndDate)
		if pErr != nil {
			return nil, pErr
		}
		project.EndDate = d
	}
	if req.DueDate != nil {
		d, pErr := models.ParseDate(*req.DueDate)
		if pErr != nil {
			return nil, pErr
		}
		project.DueDate = d
	}
	if req.EstimatedHours != nil {
		project.EstimatedHours = req.EstimatedHours
	}
	if req.Budget != nil {
		project.Budget = req.Budget
	}
	if req.Notes != nil {
		project.Notes = req.Notes
	}
	if req.Progress != nil {
		if *req.Progress < 0 || *req.Progress > 100 {
			return nil, errors.New("progress must be between 0 and 100")
		}
		project.Progress = *req.Progress
	}

	project.UpdatedBy = updatedBy

	if err := s.projectRepo.Update(project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return s.projectRepo.FindByID(id)
}

// DeleteProject soft-deletes a project
func (s *ProjectService) DeleteProject(id uint) error {
	_, err := s.projectRepo.FindByID(id)
	if err != nil {
		return errors.New("project not found")
	}
	return s.projectRepo.Delete(id)
}

// ListProjects lists projects with pagination and filters
func (s *ProjectService) ListProjects(page, pageSize int, filters map[string]interface{}) ([]models.Project, int64, error) {
	return s.projectRepo.List(page, pageSize, filters)
}

// ListMyProjects lists projects for a specific employee
func (s *ProjectService) ListMyProjects(employeeID string, page, pageSize int, filters map[string]interface{}) ([]models.Project, int64, error) {
	return s.projectRepo.ListByEmployeeID(employeeID, page, pageSize, filters)
}

// AddMembers adds members to a project
func (s *ProjectService) AddMembers(projectID uint, req *models.AddMembersRequest, assignedByID uint) ([]models.ProjectMember, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}
	_ = project

	var added []models.ProjectMember
	for _, entry := range req.Members {
		// Check if already a member
		if s.projectRepo.IsMember(projectID, entry.EmployeeID) {
			continue
		}

		// Verify employee exists
		emp, err := s.employeeRepo.FindByEmployeeID(entry.EmployeeID)
		if err != nil {
			continue
		}

		role := models.MemberRoleMember
		if entry.Role != "" {
			role = entry.Role
		}

		member := &models.ProjectMember{
			ProjectID:    projectID,
			EmployeeID:   entry.EmployeeID,
			EmployeeName: emp.FirstName + " " + emp.LastName,
			Role:         role,
			AssignedByID: &assignedByID,
			IsActive:     true,
		}

		if err := s.projectRepo.AddMember(member); err != nil {
			continue
		}
		added = append(added, *member)
	}

	return added, nil
}

// RemoveMember removes a member from a project
func (s *ProjectService) RemoveMember(projectID uint, employeeID string) error {
	if !s.projectRepo.IsMember(projectID, employeeID) {
		return errors.New("employee is not a member of this project")
	}
	return s.projectRepo.RemoveMember(projectID, employeeID)
}

// ListMembers lists all active members of a project
func (s *ProjectService) ListMembers(projectID uint) ([]models.ProjectMember, error) {
	_, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}
	return s.projectRepo.ListMembers(projectID)
}

// GetProjectProgress returns progress details for a project
func (s *ProjectService) GetProjectProgress(projectID uint) (map[string]interface{}, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	taskSummary, _ := s.taskRepo.GetSummaryByProject(projectID)
	totalHours, _ := s.taskRepo.GetTotalHoursForProject(projectID)

	result := map[string]interface{}{
		"project_id":      project.ID,
		"project_code":    project.ProjectCode,
		"project_name":    project.Name,
		"status":          project.Status,
		"progress":        project.Progress,
		"estimated_hours": project.EstimatedHours,
		"actual_hours":    totalHours,
		"member_count":    len(project.Members),
		"task_summary":    taskSummary,
	}

	if project.EstimatedHours != nil && *project.EstimatedHours > 0 {
		result["hours_utilization"] = (totalHours / *project.EstimatedHours) * 100
	}

	return result, nil
}

// GetProjectStatistics returns overall project statistics
func (s *ProjectService) GetProjectStatistics() (map[string]interface{}, error) {
	counts, err := s.projectRepo.CountByStatus()
	if err != nil {
		return nil, err
	}

	var totalActive int64
	for status, count := range counts {
		if status == string(models.ProjectStatusActive) || status == string(models.ProjectStatusPlanning) {
			totalActive += count
		}
	}

	return map[string]interface{}{
		"total_active":   totalActive,
		"by_status":      counts,
	}, nil
}

// RecalculateProjectProgress recalculates project progress from daily tasks
func (s *ProjectService) RecalculateProjectProgress(projectID uint) error {
	progress, err := s.taskRepo.GetTotalProgressForProject(projectID)
	if err != nil {
		return err
	}
	if progress > 100 {
		progress = 100
	}

	totalHours, _ := s.taskRepo.GetTotalHoursForProject(projectID)

	return s.projectRepo.UpdateProgress(projectID, progress, totalHours)
}

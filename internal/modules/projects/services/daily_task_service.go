package services

import (
	"errors"
	"fmt"
	"time"

	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/repositories"
)

// DailyTaskService handles business logic for daily tasks
type DailyTaskService struct {
	taskRepo     *repositories.DailyTaskRepository
	projectRepo  *repositories.ProjectRepository
	employeeRepo *employeeRepos.EmployeeRepository
}

// NewDailyTaskService creates a new DailyTaskService
func NewDailyTaskService() *DailyTaskService {
	return &DailyTaskService{
		taskRepo:     repositories.NewDailyTaskRepository(),
		projectRepo:  repositories.NewProjectRepository(),
		employeeRepo: employeeRepos.NewEmployeeRepository(),
	}
}

// GetEmployeeByUserID resolves user ID to employee record
func (s *DailyTaskService) GetEmployeeByUserID(userID uint) (string, string, error) {
	emp, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return "", "", errors.New("employee record not found for this user")
	}
	return emp.EmployeeID, emp.FirstName + " " + emp.LastName, nil
}

// CreateDailyTask creates a new daily task entry
func (s *DailyTaskService) CreateDailyTask(req *models.CreateDailyTaskRequest, employeeID string, employeeName string) (*models.DailyTask, error) {
	// Validate project exists
	project, err := s.projectRepo.FindByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	// Validate employee is a member of the project
	if !s.projectRepo.IsMember(req.ProjectID, employeeID) {
		return nil, errors.New("you are not a member of this project")
	}

	// Validate project is active
	if project.Status != models.ProjectStatusActive && project.Status != models.ProjectStatusPlanning {
		return nil, errors.New("cannot add tasks to a project that is not active or planning")
	}

	// Parse task date (default to today)
	taskDate := time.Now().Truncate(24 * time.Hour)
	if req.TaskDate != nil && *req.TaskDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.TaskDate)
		if err != nil {
			return nil, errors.New("invalid task_date format, expected YYYY-MM-DD")
		}
		taskDate = parsed
	}

	// Validate progress contribution
	progressContribution := float64(0)
	if req.ProgressContribution != nil {
		if *req.ProgressContribution < 0 || *req.ProgressContribution > 100 {
			return nil, errors.New("progress_contribution must be between 0 and 100")
		}
		progressContribution = *req.ProgressContribution

		// Check that total progress won't exceed 100%
		currentProgress, _ := s.taskRepo.GetTotalProgressForProject(req.ProjectID)
		if currentProgress+progressContribution > 100 {
			return nil, fmt.Errorf("total progress would exceed 100%% (current: %.1f%%, adding: %.1f%%)", currentProgress, progressContribution)
		}
	}

	status := models.TaskStatusCompleted
	if req.Status != nil {
		status = *req.Status
	}

	task := &models.DailyTask{
		ProjectID:            req.ProjectID,
		ProjectName:          project.Name,
		EmployeeID:           employeeID,
		EmployeeName:         employeeName,
		TaskDate:             taskDate,
		Title:                req.Title,
		Description:          req.Description,
		HoursSpent:           req.HoursSpent,
		Status:               status,
		Category:             req.Category,
		ProgressContribution: progressContribution,
		Tags:                 req.Tags,
		Notes:                req.Notes,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create daily task: %w", err)
	}

	// Recalculate project progress
	s.recalculateProgress(req.ProjectID)

	return task, nil
}

// GetDailyTask returns a daily task by ID
func (s *DailyTaskService) GetDailyTask(id uint) (*models.DailyTask, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("daily task not found")
	}
	return task, nil
}

// UpdateDailyTask updates a daily task
func (s *DailyTaskService) UpdateDailyTask(id uint, req *models.UpdateDailyTaskRequest, employeeID string) (*models.DailyTask, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("daily task not found")
	}

	// Only the task creator can update
	if task.EmployeeID != employeeID {
		return nil, errors.New("you can only update your own tasks")
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.HoursSpent != nil {
		task.HoursSpent = *req.HoursSpent
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Category != nil {
		task.Category = *req.Category
	}
	if req.ProgressContribution != nil {
		if *req.ProgressContribution < 0 || *req.ProgressContribution > 100 {
			return nil, errors.New("progress_contribution must be between 0 and 100")
		}

		// Check total progress (excluding current task's contribution)
		currentProgress, _ := s.taskRepo.GetTotalProgressForProject(task.ProjectID)
		newTotal := currentProgress - task.ProgressContribution + *req.ProgressContribution
		if newTotal > 100 {
			return nil, fmt.Errorf("total progress would exceed 100%% (current total: %.1f%%)", currentProgress)
		}
		task.ProgressContribution = *req.ProgressContribution
	}
	if req.Tags != nil {
		task.Tags = req.Tags
	}
	if req.Notes != nil {
		task.Notes = req.Notes
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update daily task: %w", err)
	}

	// Recalculate project progress
	s.recalculateProgress(task.ProjectID)

	return task, nil
}

// DeleteDailyTask deletes a daily task
func (s *DailyTaskService) DeleteDailyTask(id uint, employeeID string) error {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return errors.New("daily task not found")
	}

	// Only the task creator can delete
	if task.EmployeeID != employeeID {
		return errors.New("you can only delete your own tasks")
	}

	projectID := task.ProjectID

	if err := s.taskRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete daily task: %w", err)
	}

	// Recalculate project progress
	s.recalculateProgress(projectID)

	return nil
}

// ListDailyTasks lists daily tasks with pagination and filters
func (s *DailyTaskService) ListDailyTasks(page, pageSize int, filters map[string]interface{}) ([]models.DailyTask, int64, error) {
	return s.taskRepo.List(page, pageSize, filters)
}

// GetMyTasks returns daily tasks for an employee with optional date filters
func (s *DailyTaskService) GetMyTasks(employeeID string, page, pageSize int, filters map[string]interface{}) ([]models.DailyTask, int64, error) {
	filters["employee_id"] = employeeID
	return s.taskRepo.List(page, pageSize, filters)
}

// GetMyTodayTasks returns today's tasks for an employee
func (s *DailyTaskService) GetMyTodayTasks(employeeID string) ([]models.DailyTask, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return s.taskRepo.ListByEmployeeAndDateRange(employeeID, today, today)
}

// GetMyTaskSummary returns a task summary for an employee within a date range
func (s *DailyTaskService) GetMyTaskSummary(employeeID string, from, to time.Time) (*models.DailyTaskSummary, error) {
	return s.taskRepo.GetSummaryByEmployee(employeeID, from, to)
}

// GetProjectTaskSummary returns task summary for a project
func (s *DailyTaskService) GetProjectTaskSummary(projectID uint) (*models.DailyTaskSummary, error) {
	_, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}
	return s.taskRepo.GetSummaryByProject(projectID)
}

// GetTaskCategories returns all available task categories
func (s *DailyTaskService) GetTaskCategories() []map[string]string {
	return []map[string]string{
		{"value": "development", "label": "Development"},
		{"value": "testing", "label": "Testing / QA"},
		{"value": "design", "label": "Design / UI/UX"},
		{"value": "documentation", "label": "Documentation"},
		{"value": "meeting", "label": "Meeting"},
		{"value": "research", "label": "Research"},
		{"value": "planning", "label": "Planning"},
		{"value": "review", "label": "Code Review"},
		{"value": "support", "label": "Support / Bug Fix"},
		{"value": "other", "label": "Other"},
	}
}

// recalculateProgress recalculates the project progress from daily tasks
func (s *DailyTaskService) recalculateProgress(projectID uint) {
	progress, err := s.taskRepo.GetTotalProgressForProject(projectID)
	if err != nil {
		return
	}
	if progress > 100 {
		progress = 100
	}
	totalHours, _ := s.taskRepo.GetTotalHoursForProject(projectID)
	_ = s.projectRepo.UpdateProgress(projectID, progress, totalHours)
}

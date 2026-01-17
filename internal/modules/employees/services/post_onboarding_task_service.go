package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
)

type PostOnboardingTaskService struct {
	taskRepo     *repositories.PostOnboardingTaskRepository
	employeeRepo *repositories.EmployeeRepository
	userRepo     *userRepos.UserRepository
}

func NewPostOnboardingTaskService() *PostOnboardingTaskService {
	return &PostOnboardingTaskService{
		taskRepo:     repositories.NewPostOnboardingTaskRepository(),
		employeeRepo: repositories.NewEmployeeRepository(),
		userRepo:     userRepos.NewUserRepository(),
	}
}

// GetTaskTypeInfo returns information about a task type
func (s *PostOnboardingTaskService) GetTaskTypeInfo(taskType string) (title string, defaultAssignedTo string, description string) {
	taskInfo := map[string]struct {
		Title           string
		DefaultAssigned string
		Description     string
	}{
		string(models.TaskTypeWelcomeEmail): {
			Title:           "Send Welcome Email",
			DefaultAssigned: "HR",
			Description:     "Send welcome email to new employee",
		},
		string(models.TaskTypeCreateEmailAccount): {
			Title:           "Create Email Account",
			DefaultAssigned: "HR",
			Description:     "Create company email account for employee",
		},
		string(models.TaskTypeAllocateLaptop): {
			Title:           "Allocate Laptop",
			DefaultAssigned: "IT",
			Description:     "Allocate laptop/computer to employee",
		},
		string(models.TaskTypeDeskAllocate): {
			Title:           "Desk Allocation",
			DefaultAssigned: "HR",
			Description:     "Assign desk/workspace to employee",
		},
		string(models.TaskTypeIssueIDCard): {
			Title:           "Issue ID Card",
			DefaultAssigned: "HR",
			Description:     "Issue employee identification card",
		},
		string(models.TaskTypeTeamIntroduction): {
			Title:           "Team Introduction",
			DefaultAssigned: "HR",
			Description:     "Introduce employee to team members",
		},
		string(models.TaskTypePolicyAcknowledge): {
			Title:           "Policy Acknowledgment",
			DefaultAssigned: "HR",
			Description:     "Employee to acknowledge company policies",
		},
		string(models.TaskTypeInductionTraining): {
			Title:           "Induction Training",
			DefaultAssigned: "HR",
			Description:     "Complete induction training program",
		},
		string(models.TaskTypeAssignBuddyMentor): {
			Title:           "Assign Buddy/Mentor",
			DefaultAssigned: "HR",
			Description:     "Assign a buddy or mentor to employee",
		},
		string(models.TaskTypeSetupDevEnvironment): {
			Title:           "Setup Development Environment",
			DefaultAssigned: "IT",
			Description:     "Setup development environment for software developer",
		},
		string(models.TaskTypeFirstWeekReview): {
			Title:           "First Week Review",
			DefaultAssigned: "HR",
			Description:     "Conduct first week review with employee",
		},
		string(models.TaskTypeFirstMonthReview): {
			Title:           "First Month Review",
			DefaultAssigned: "HR",
			Description:     "Conduct first month review with employee",
		},
		string(models.TaskTypeMidProbationReview): {
			Title:           "Mid-Probation Review",
			DefaultAssigned: "HR",
			Description:     "Conduct mid-probation period review",
		},
		string(models.TaskTypeProbationConfirmation): {
			Title:           "Probation Confirmation Assessment",
			DefaultAssigned: "HR",
			Description:     "Final probation confirmation assessment",
		},
	}

	if info, ok := taskInfo[strings.ToLower(taskType)]; ok {
		return info.Title, info.DefaultAssigned, info.Description
	}

	return taskType, "HR", "Post-onboarding task"
}

// CreateTask creates a new post-onboarding task
func (s *PostOnboardingTaskService) CreateTask(req *models.CreatePostOnboardingTaskRequest, tenantID *uint, createdBy *uint) (*models.PostOnboardingTask, error) {
	// Validate employee exists
	_, err := s.employeeRepo.FindByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Validate task type
	validTaskTypes := []string{
		string(models.TaskTypeWelcomeEmail),
		string(models.TaskTypeCreateEmailAccount),
		string(models.TaskTypeAllocateLaptop),
		string(models.TaskTypeDeskAllocate),
		string(models.TaskTypeIssueIDCard),
		string(models.TaskTypeTeamIntroduction),
		string(models.TaskTypePolicyAcknowledge),
		string(models.TaskTypeInductionTraining),
		string(models.TaskTypeAssignBuddyMentor),
		string(models.TaskTypeSetupDevEnvironment),
		string(models.TaskTypeFirstWeekReview),
		string(models.TaskTypeFirstMonthReview),
		string(models.TaskTypeMidProbationReview),
		string(models.TaskTypeProbationConfirmation),
	}

	isValidType := false
	for _, t := range validTaskTypes {
		if strings.ToLower(req.TaskType) == strings.ToLower(t) {
			isValidType = true
			req.TaskType = t // Normalize to correct case
			break
		}
	}
	if !isValidType {
		return nil, errors.New("invalid task type")
	}

	// Get task info
	title, defaultAssigned, _ := s.GetTaskTypeInfo(req.TaskType)

	// Use provided title or default
	if req.Title != nil && *req.Title != "" {
		title = *req.Title
	}

	// Use provided assigned_to or default
	assignedTo := defaultAssigned
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		assignedTo = *req.AssignedTo
	}

	// Set priority
	priority := string(models.TaskPriorityMedium)
	if req.Priority != nil && *req.Priority != "" {
		validPriorities := []string{
			string(models.TaskPriorityLow),
			string(models.TaskPriorityMedium),
			string(models.TaskPriorityHigh),
			string(models.TaskPriorityCritical),
		}
		for _, p := range validPriorities {
			if strings.ToLower(*req.Priority) == strings.ToLower(p) {
				priority = p
				break
			}
		}
	}

	// Parse due date
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, errors.New("invalid due_date format, expected YYYY-MM-DD")
		}
		dueDate = &parsed
	}

	// Validate assigned user if provided
	if req.AssignedToUserID != nil {
		_, err := s.userRepo.FindByID(*req.AssignedToUserID)
		if err != nil {
			return nil, errors.New("assigned user not found")
		}
	}

	// Create task
	task := &models.PostOnboardingTask{
		TenantID:         tenantID,
		EmployeeID:       req.EmployeeID,
		TaskType:         req.TaskType,
		Title:            title,
		Description:      req.Description,
		Status:           string(models.TaskStatusPending),
		Priority:         priority,
		AssignedTo:       &assignedTo,
		AssignedToUserID: req.AssignedToUserID,
		DueDate:          dueDate,
		Notes:            req.Notes,
		Metadata:         req.Metadata,
		CreatedBy:        createdBy,
		UpdatedBy:        createdBy,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// BulkCreateTasks creates multiple tasks for an employee
func (s *PostOnboardingTaskService) BulkCreateTasks(req *models.BulkCreatePostOnboardingTasksRequest, tenantID *uint, createdBy *uint) ([]models.PostOnboardingTask, error) {
	// Validate employee exists
	_, err := s.employeeRepo.FindByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	var tasks []models.PostOnboardingTask
	for _, taskType := range req.TaskTypes {
		createReq := &models.CreatePostOnboardingTaskRequest{
			EmployeeID: req.EmployeeID,
			TaskType:   taskType,
			Priority:   req.Priority,
			DueDate:    req.DueDate,
		}

		task, err := s.CreateTask(createReq, tenantID, createdBy)
		if err != nil {
			return nil, fmt.Errorf("failed to create task %s: %w", taskType, err)
		}

		tasks = append(tasks, *task)
	}

	return tasks, nil
}

// UpdateTask updates a post-onboarding task
func (s *PostOnboardingTaskService) UpdateTask(req *models.UpdatePostOnboardingTaskRequest, tenantID *uint, updatedBy *uint) (*models.PostOnboardingTask, error) {
	task, err := s.taskRepo.FindByID(req.ID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	// Verify tenant ownership
	if tenantID != nil && task.TenantID != nil && *task.TenantID != *tenantID {
		return nil, errors.New("task does not belong to your tenant")
	}

	// Update fields
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.Status != nil {
		validStatuses := []string{
			string(models.TaskStatusPending),
			string(models.TaskStatusInProgress),
			string(models.TaskStatusCompleted),
			string(models.TaskStatusSkipped),
		}
		isValidStatus := false
		for _, s := range validStatuses {
			if strings.ToLower(*req.Status) == strings.ToLower(s) {
				task.Status = s
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			return nil, errors.New("invalid status, must be one of: pending, in_progress, completed, skipped")
		}

		// If marking as completed, set completed_at and completed_by
		if task.Status == string(models.TaskStatusCompleted) && task.CompletedAt == nil {
			now := time.Now()
			task.CompletedAt = &now
			task.CompletedByUserID = updatedBy
		}
	}
	if req.Priority != nil {
		validPriorities := []string{
			string(models.TaskPriorityLow),
			string(models.TaskPriorityMedium),
			string(models.TaskPriorityHigh),
			string(models.TaskPriorityCritical),
		}
		for _, p := range validPriorities {
			if strings.ToLower(*req.Priority) == strings.ToLower(p) {
				task.Priority = p
				break
			}
		}
	}
	if req.AssignedTo != nil {
		task.AssignedTo = req.AssignedTo
	}
	if req.AssignedToUserID != nil {
		// Validate user exists
		_, err := s.userRepo.FindByID(*req.AssignedToUserID)
		if err != nil {
			return nil, errors.New("assigned user not found")
		}
		task.AssignedToUserID = req.AssignedToUserID
	}
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, errors.New("invalid due_date format, expected YYYY-MM-DD")
		}
		task.DueDate = &parsed
	}
	if req.Notes != nil {
		task.Notes = req.Notes
	}
	if req.Metadata != nil {
		task.Metadata = req.Metadata
	}

	task.UpdatedBy = updatedBy
	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// CompleteTask marks a task as completed
func (s *PostOnboardingTaskService) CompleteTask(req *models.CompletePostOnboardingTaskRequest, tenantID *uint, completedBy *uint) (*models.PostOnboardingTask, error) {
	task, err := s.taskRepo.FindByID(req.ID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	// Verify tenant ownership
	if tenantID != nil && task.TenantID != nil && *task.TenantID != *tenantID {
		return nil, errors.New("task does not belong to your tenant")
	}

	// Check if task can be completed
	if !task.CanBeCompleted() {
		return nil, errors.New("task cannot be completed (already completed or skipped)")
	}

	// Mark as completed
	now := time.Now()
	task.Status = string(models.TaskStatusCompleted)
	task.CompletedAt = &now
	task.CompletedByUserID = completedBy
	if req.Notes != nil {
		task.Notes = req.Notes
	}
	task.UpdatedBy = completedBy
	task.UpdatedAt = now

	if err := s.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to complete task: %w", err)
	}

	return task, nil
}

// GetTask retrieves a task by ID
func (s *PostOnboardingTaskService) GetTask(id uint, tenantID *uint) (*models.PostOnboardingTask, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("task not found")
	}

	// Verify tenant ownership
	if tenantID != nil && task.TenantID != nil && *task.TenantID != *tenantID {
		return nil, errors.New("task does not belong to your tenant")
	}

	return task, nil
}

// GetTasksByEmployeeID retrieves all tasks for an employee
func (s *PostOnboardingTaskService) GetTasksByEmployeeID(employeeID string, tenantID *uint) ([]models.PostOnboardingTask, error) {
	// Validate employee exists
	employee, err := s.employeeRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Verify tenant ownership
	if tenantID != nil && employee.TenantID != nil && *employee.TenantID != *tenantID {
		return nil, errors.New("employee does not belong to your tenant")
	}

	tasks, err := s.taskRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tasks: %w", err)
	}

	return tasks, nil
}

// ListTasks lists tasks with pagination and filters
func (s *PostOnboardingTaskService) ListTasks(req *models.GetPostOnboardingTasksRequest, tenantID *uint) ([]models.PostOnboardingTask, int64, error) {
	// Set defaults
	page := 1
	if req.Page > 0 {
		page = req.Page
	}

	pageSize := 20
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.EmployeeID != nil && *req.EmployeeID != "" {
		filters["employee_id"] = *req.EmployeeID
	}
	if req.Status != nil && *req.Status != "" {
		filters["status"] = *req.Status
	}
	if req.TaskType != nil && *req.TaskType != "" {
		filters["task_type"] = *req.TaskType
	}
	if req.Priority != nil && *req.Priority != "" {
		filters["priority"] = *req.Priority
	}
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		filters["assigned_to"] = *req.AssignedTo
	}

	return s.taskRepo.List(tenantID, page, pageSize, filters)
}

// GetTaskTypes returns available task types
func (s *PostOnboardingTaskService) GetTaskTypes() []models.GetTaskTypesResponse {
	return []models.GetTaskTypesResponse{
		{Type: string(models.TaskTypeWelcomeEmail), Name: "Send Welcome Email", Description: "Send welcome email to new employee", DefaultAssignedTo: "HR", IsConditional: false},
		{Type: string(models.TaskTypeCreateEmailAccount), Name: "Create Email Account", Description: "Create company email account for employee", DefaultAssignedTo: "HR", IsConditional: false},
		{Type: string(models.TaskTypeAllocateLaptop), Name: "Allocate Laptop", Description: "Allocate laptop/computer to employee", DefaultAssignedTo: "IT", IsConditional: false},
		{Type: string(models.TaskTypeDeskAllocate), Name: "Desk Allocation", Description: "Assign desk/workspace to employee", DefaultAssignedTo: "HR", IsConditional: false},
		{Type: string(models.TaskTypeIssueIDCard), Name: "Issue ID Card", Description: "Issue employee identification card", DefaultAssignedTo: "HR", IsConditional: false},
		{Type: string(models.TaskTypeTeamIntroduction), Name: "Team Introduction", Description: "Introduce employee to team members", DefaultAssignedTo: "HR", IsConditional: false},
		{Type: string(models.TaskTypePolicyAcknowledge), Name: "Policy Acknowledgment", Description: "Employee to acknowledge company policies", DefaultAssignedTo: "HR", IsConditional: false},
		{Type: string(models.TaskTypeInductionTraining), Name: "Induction Training", Description: "Complete induction training program", DefaultAssignedTo: "HR", IsConditional: false},
		{Type: string(models.TaskTypeAssignBuddyMentor), Name: "Assign Buddy/Mentor", Description: "Assign a buddy or mentor to employee", DefaultAssignedTo: "HR", IsConditional: false},
		{Type: string(models.TaskTypeSetupDevEnvironment), Name: "Setup Development Environment", Description: "Setup development environment for software developer", DefaultAssignedTo: "IT", IsConditional: true},
		{Type: string(models.TaskTypeFirstWeekReview), Name: "First Week Review", Description: "Conduct first week review with employee", DefaultAssignedTo: "HR", IsConditional: false},
		{Type: string(models.TaskTypeFirstMonthReview), Name: "First Month Review", Description: "Conduct first month review with employee", DefaultAssignedTo: "HR", IsConditional: false},
		{Type: string(models.TaskTypeMidProbationReview), Name: "Mid-Probation Review", Description: "Conduct mid-probation period review", DefaultAssignedTo: "HR", IsConditional: false},
		{Type: string(models.TaskTypeProbationConfirmation), Name: "Probation Confirmation Assessment", Description: "Final probation confirmation assessment", DefaultAssignedTo: "HR", IsConditional: false},
	}
}

// GetTaskSummary returns a summary of tasks for an employee
func (s *PostOnboardingTaskService) GetTaskSummary(employeeID string, tenantID *uint) (map[string]interface{}, error) {
	// Validate employee exists
	employee, err := s.employeeRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Verify tenant ownership
	if tenantID != nil && employee.TenantID != nil && *employee.TenantID != *tenantID {
		return nil, errors.New("employee does not belong to your tenant")
	}

	// Get task counts
	counts, err := s.taskRepo.CountByEmployeeID(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task counts: %w", err)
	}

	// Get all tasks
	tasks, err := s.taskRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	// Calculate completion percentage
	totalTasks := int64(len(tasks))
	completedTasks := counts[string(models.TaskStatusCompleted)]
	var completionPercentage float64
	if totalTasks > 0 {
		completionPercentage = float64(completedTasks) / float64(totalTasks) * 100
	}

	return map[string]interface{}{
		"employee_id":           employeeID,
		"total_tasks":           totalTasks,
		"pending":               counts[string(models.TaskStatusPending)],
		"in_progress":           counts[string(models.TaskStatusInProgress)],
		"completed":             completedTasks,
		"skipped":               counts[string(models.TaskStatusSkipped)],
		"completion_percentage": completionPercentage,
		"tasks":                 tasks,
	}, nil
}

// DeleteTask deletes a task
func (s *PostOnboardingTaskService) DeleteTask(id uint, tenantID *uint) error {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return errors.New("task not found")
	}

	// Verify tenant ownership
	if tenantID != nil && task.TenantID != nil && *task.TenantID != *tenantID {
		return errors.New("task does not belong to your tenant")
	}

	if err := s.taskRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// ListEmployeesWithTaskCompletion lists all employees with their post-onboarding task completion percentages
func (s *PostOnboardingTaskService) ListEmployeesWithTaskCompletion(tenantID *uint, page, pageSize int) ([]map[string]interface{}, int64, error) {
	// Get employees with pagination
	offset := (page - 1) * pageSize
	employees, total, err := s.employeeRepo.List(pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get employees: %w", err)
	}

	// Filter by tenant if needed
	filteredEmployees := make([]models.Employee, 0)
	for _, emp := range employees {
		if tenantID == nil || (emp.TenantID != nil && *emp.TenantID == *tenantID) {
			filteredEmployees = append(filteredEmployees, emp)
		}
	}

	// Get employee IDs
	employeeIDs := make([]string, len(filteredEmployees))
	for i, emp := range filteredEmployees {
		employeeIDs[i] = emp.EmployeeID
	}

	// Get task completion data for all employees in batch
	taskDataMap := make(map[string]struct {
		TotalTasks        int64
		CompletedTasks    int64
		PendingTasks      int64
		InProgressTasks   int64
		SkippedTasks      int64
		CompletionPercent float64
		IsCompleted       bool
	})

	if len(employeeIDs) > 0 {
		// Get task counts for all employees at once
		for _, empID := range employeeIDs {
			counts, err := s.taskRepo.CountByEmployeeID(empID)
			if err != nil {
				continue // Skip if error, will default to 0
			}

			totalTasks := int64(0)
			completedTasks := int64(0)
			pendingTasks := int64(0)
			inProgressTasks := int64(0)
			skippedTasks := int64(0)

			for status, count := range counts {
				totalTasks += count
				switch status {
				case string(models.TaskStatusCompleted):
					completedTasks = count
				case string(models.TaskStatusPending):
					pendingTasks = count
				case string(models.TaskStatusInProgress):
					inProgressTasks = count
				case string(models.TaskStatusSkipped):
					skippedTasks = count
				}
			}

			var completionPercent float64
			var isCompleted bool
			if totalTasks > 0 {
				completionPercent = float64(completedTasks) / float64(totalTasks) * 100
				isCompleted = completedTasks == totalTasks
			}

			taskDataMap[empID] = struct {
				TotalTasks        int64
				CompletedTasks    int64
				PendingTasks      int64
				InProgressTasks   int64
				SkippedTasks      int64
				CompletionPercent float64
				IsCompleted       bool
			}{
				TotalTasks:        totalTasks,
				CompletedTasks:    completedTasks,
				PendingTasks:      pendingTasks,
				InProgressTasks:   inProgressTasks,
				SkippedTasks:      skippedTasks,
				CompletionPercent: completionPercent,
				IsCompleted:       isCompleted,
			}
		}
	}

	// Build response with employee details and task completion
	result := make([]map[string]interface{}, len(filteredEmployees))
	for i, emp := range filteredEmployees {
		employeeData := map[string]interface{}{
			"employee_id":    emp.EmployeeID,
			"first_name":     emp.FirstName,
			"last_name":      emp.LastName,
			"email":          emp.WorkEmail,
			"department_id":  emp.DepartmentID,
			"position_id":    emp.PositionID,
			"status":         emp.Status,
			"hire_date":      emp.HireDate,
		}

		// Add task completion data if available
		if taskInfo, ok := taskDataMap[emp.EmployeeID]; ok {
			employeeData["total_tasks"] = taskInfo.TotalTasks
			employeeData["completed_tasks"] = taskInfo.CompletedTasks
			employeeData["pending_tasks"] = taskInfo.PendingTasks
			employeeData["in_progress_tasks"] = taskInfo.InProgressTasks
			employeeData["skipped_tasks"] = taskInfo.SkippedTasks
			employeeData["completion_percentage"] = taskInfo.CompletionPercent
			employeeData["is_completed"] = taskInfo.IsCompleted
		} else {
			// Employee has no tasks
			employeeData["total_tasks"] = 0
			employeeData["completed_tasks"] = 0
			employeeData["pending_tasks"] = 0
			employeeData["in_progress_tasks"] = 0
			employeeData["skipped_tasks"] = 0
			employeeData["completion_percentage"] = 0.0
			employeeData["is_completed"] = false
		}

		result[i] = employeeData
	}

	return result, total, nil
}

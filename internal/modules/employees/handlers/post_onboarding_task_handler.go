package handlers

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type PostOnboardingTaskHandler struct {
	service  *services.PostOnboardingTaskService
	userRepo *userRepos.UserRepository
}

func NewPostOnboardingTaskHandler() *PostOnboardingTaskHandler {
	return &PostOnboardingTaskHandler{
		service:  services.NewPostOnboardingTaskService(),
		userRepo: userRepos.NewUserRepository(),
	}
}

// toTaskResponse converts PostOnboardingTask model to PostOnboardingTaskResponse
func (h *PostOnboardingTaskHandler) toTaskResponse(task *models.PostOnboardingTask) *models.PostOnboardingTaskResponse {
	resp := &models.PostOnboardingTaskResponse{
		ID:                task.ID,
		EmployeeID:        task.EmployeeID,
		TaskType:          task.TaskType,
		Title:             task.Title,
		Description:       task.Description,
		Status:            task.Status,
		Priority:          task.Priority,
		AssignedTo:        task.AssignedTo,
		AssignedToUserID:  task.AssignedToUserID,
		CompletedByUserID: task.CompletedByUserID,
		CompletedAt:       task.CompletedAt,
		DueDate:           task.DueDate,
		Notes:             task.Notes,
		Metadata:          task.Metadata,
		CreatedBy:         task.CreatedBy,
		UpdatedBy:         task.UpdatedBy,
		CreatedAt:         task.CreatedAt,
		UpdatedAt:         task.UpdatedAt,
	}

	// Get employee name if needed (can be loaded separately)
	// For now, we'll leave it empty and can be populated in a separate call if needed

	// Get assigned user name
	if task.AssignedToUserID != nil {
		user, err := h.userRepo.FindByID(*task.AssignedToUserID)
		if err == nil && user != nil {
			fullName := user.FirstName + " " + user.LastName
			resp.AssignedToUserName = &fullName
		}
	}

	// Get completed by user name
	if task.CompletedByUserID != nil {
		user, err := h.userRepo.FindByID(*task.CompletedByUserID)
		if err == nil && user != nil {
			fullName := user.FirstName + " " + user.LastName
			resp.CompletedByUserName = &fullName
		}
	}

	return resp
}

// CreateTask handles creating a new post-onboarding task
func (h *PostOnboardingTaskHandler) CreateTask(c *gin.Context) {
	var req models.CreatePostOnboardingTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var createdBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		createdBy = &userObj.ID
	}

	task, err := h.service.CreateTask(&req, tenantID, createdBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Post-onboarding task created successfully", h.toTaskResponse(task))
}

// BulkCreateTasks handles creating multiple tasks for an employee
func (h *PostOnboardingTaskHandler) BulkCreateTasks(c *gin.Context) {
	var req models.BulkCreatePostOnboardingTasksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var createdBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		createdBy = &userObj.ID
	}

	tasks, err := h.service.BulkCreateTasks(&req, tenantID, createdBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert to response format
	taskResponses := make([]models.PostOnboardingTaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = *h.toTaskResponse(&task)
	}

	response.Created(c, "Post-onboarding tasks created successfully", taskResponses)
}

// UpdateTask handles updating a post-onboarding task
func (h *PostOnboardingTaskHandler) UpdateTask(c *gin.Context) {
	var req models.UpdatePostOnboardingTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	task, err := h.service.UpdateTask(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Task updated successfully", h.toTaskResponse(task))
}

// CompleteTask handles completing a task
func (h *PostOnboardingTaskHandler) CompleteTask(c *gin.Context) {
	var req models.CompletePostOnboardingTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var completedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		completedBy = &userObj.ID
	}

	task, err := h.service.CompleteTask(&req, tenantID, completedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Task completed successfully", h.toTaskResponse(task))
}

// GetTask handles getting a task by ID
func (h *PostOnboardingTaskHandler) GetTask(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", "id is required in request body")
		return
	}

	tenantID := middleware.GetTenantID(c)
	task, err := h.service.GetTask(req.ID, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Task retrieved successfully", h.toTaskResponse(task))
}

// GetTasksByEmployeeID handles getting all tasks for an employee
func (h *PostOnboardingTaskHandler) GetTasksByEmployeeID(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		response.ValidationError(c, "Validation failed", "employee_id is required")
		return
	}

	tenantID := middleware.GetTenantID(c)
	tasks, err := h.service.GetTasksByEmployeeID(employeeID, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert to response format
	taskResponses := make([]models.PostOnboardingTaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = *h.toTaskResponse(&task)
	}

	response.Success(c, "Tasks retrieved successfully", taskResponses)
}

// GetTaskSummary handles getting task summary for an employee
func (h *PostOnboardingTaskHandler) GetTaskSummary(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		response.ValidationError(c, "Validation failed", "employee_id is required")
		return
	}

	tenantID := middleware.GetTenantID(c)
	summary, err := h.service.GetTaskSummary(employeeID, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Task summary retrieved successfully", summary)
}

// ListTasks handles listing tasks with pagination and filters
func (h *PostOnboardingTaskHandler) ListTasks(c *gin.Context) {
	var req models.GetPostOnboardingTasksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	tasks, total, err := h.service.ListTasks(&req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert to response format
	taskResponses := make([]models.PostOnboardingTaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = *h.toTaskResponse(&task)
	}

	// Calculate pagination metadata
	page := req.Page
	if page == 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 20
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response.SuccessWithMeta(c, "Tasks retrieved successfully", taskResponses, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetTaskTypes handles getting available task types
func (h *PostOnboardingTaskHandler) GetTaskTypes(c *gin.Context) {
	taskTypes := h.service.GetTaskTypes()
	response.Success(c, "Task types retrieved successfully", taskTypes)
}

// DeleteTask handles deleting a task
func (h *PostOnboardingTaskHandler) DeleteTask(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", "id is required in request body")
		return
	}

	tenantID := middleware.GetTenantID(c)
	if err := h.service.DeleteTask(req.ID, tenantID); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Task deleted successfully", nil)
}

// ListEmployeesWithTaskCompletion handles listing all employees with their post-onboarding task completion percentages
func (h *PostOnboardingTaskHandler) ListEmployeesWithTaskCompletion(c *gin.Context) {
	var req struct {
		Page     int `form:"page" binding:"omitempty,min=1"`
		PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	// Set defaults
	page := 1
	if req.Page > 0 {
		page = req.Page
	}
	pageSize := 20
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	tenantID := middleware.GetTenantID(c)
	employees, total, err := h.service.ListEmployeesWithTaskCompletion(tenantID, page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Calculate pagination metadata
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response.SuccessWithMeta(c, "Employees with task completion retrieved successfully", employees, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

package handlers

import (
	"strconv"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// DailyTaskHandler handles HTTP requests for daily tasks
type DailyTaskHandler struct {
	service        *services.DailyTaskService
	projectService *services.ProjectService
}

// NewDailyTaskHandler creates a new DailyTaskHandler
func NewDailyTaskHandler() *DailyTaskHandler {
	return &DailyTaskHandler{
		service:        services.NewDailyTaskService(),
		projectService: services.NewProjectService(),
	}
}

// ────────────────────────── Daily Task Endpoints ──────────────────────────

// CreateDailyTask creates a new daily task entry
// @Summary Log a daily task
// @Tags Daily Tasks
// @Accept json
// @Produce json
// @Param body body models.CreateDailyTaskRequest true "Daily task data"
// @Success 201 {object} response.APIResponse
// @Router /api/v1/daily-tasks [post]
func (h *DailyTaskHandler) CreateDailyTask(c *gin.Context) {
	var req models.CreateDailyTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	employeeID, employeeName, ok := getEmployeeFromUser(c)
	if !ok || employeeID == "" {
		response.BadRequest(c, "No employee profile linked to your account", nil)
		return
	}

	task, err := h.service.CreateDailyTask(&req, employeeID, employeeName)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Daily task logged successfully", task)
}

// GetDailyTask returns a daily task by ID
// @Summary Get daily task details
// @Tags Daily Tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/daily-tasks/{id} [get]
func (h *DailyTaskHandler) GetDailyTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid task ID", nil)
		return
	}

	task, err := h.service.GetDailyTask(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Daily task retrieved successfully", task)
}

// UpdateDailyTask updates a daily task
// @Summary Update a daily task
// @Tags Daily Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param body body models.UpdateDailyTaskRequest true "Updated task data"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/daily-tasks/{id} [put]
func (h *DailyTaskHandler) UpdateDailyTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid task ID", nil)
		return
	}

	var req models.UpdateDailyTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	employeeID, _, ok := getEmployeeFromUser(c)
	if !ok || employeeID == "" {
		response.BadRequest(c, "No employee profile linked to your account", nil)
		return
	}

	task, err := h.service.UpdateDailyTask(uint(id), &req, employeeID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Daily task updated successfully", task)
}

// DeleteDailyTask deletes a daily task
// @Summary Delete a daily task
// @Tags Daily Tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/daily-tasks/{id} [delete]
func (h *DailyTaskHandler) DeleteDailyTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid task ID", nil)
		return
	}

	employeeID, _, ok := getEmployeeFromUser(c)
	if !ok || employeeID == "" {
		response.BadRequest(c, "No employee profile linked to your account", nil)
		return
	}

	if err := h.service.DeleteDailyTask(uint(id), employeeID); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Daily task deleted successfully", nil)
}

// ListDailyTasks lists daily tasks with filters (HR/Admin view)
// @Summary List daily tasks (Admin)
// @Tags Daily Tasks
// @Produce json
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param project_id query int false "Filter by project"
// @Param employee_id query string false "Filter by employee"
// @Param status query string false "Filter by status"
// @Param category query string false "Filter by category"
// @Param date_from query string false "From date (YYYY-MM-DD)"
// @Param date_to query string false "To date (YYYY-MM-DD)"
// @Param search query string false "Search"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/daily-tasks [get]
func (h *DailyTaskHandler) ListDailyTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}

	filters := buildTaskFilters(c)

	tasks, total, err := h.service.ListDailyTasks(page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list daily tasks", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.SuccessWithMeta(c, "Daily tasks retrieved successfully", tasks, meta)
}

// GetTaskCategories returns available task categories
// @Summary Get task categories
// @Tags Daily Tasks
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/daily-tasks/categories [get]
func (h *DailyTaskHandler) GetTaskCategories(c *gin.Context) {
	categories := h.service.GetTaskCategories()
	response.Success(c, "Task categories retrieved successfully", categories)
}

// GetProjectTaskSummary returns task summary for a specific project
// @Summary Get project task summary
// @Tags Daily Tasks
// @Produce json
// @Param project_id path int true "Project ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/daily-tasks/project/{project_id}/summary [get]
func (h *DailyTaskHandler) GetProjectTaskSummary(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("project_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid project ID", nil)
		return
	}

	summary, err := h.service.GetProjectTaskSummary(uint(projectID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Project task summary retrieved successfully", summary)
}

// ────────────────────────── Self-Service Endpoints ──────────────────────────

// GetMyTasks returns daily tasks for the current employee
// @Summary Get my daily tasks
// @Tags Daily Tasks (Self-Service)
// @Produce json
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param project_id query int false "Filter by project"
// @Param date_from query string false "From date (YYYY-MM-DD)"
// @Param date_to query string false "To date (YYYY-MM-DD)"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/daily-tasks [get]
func (h *DailyTaskHandler) GetMyTasks(c *gin.Context) {
	employeeID, _, ok := getEmployeeFromUser(c)
	if !ok || employeeID == "" {
		response.BadRequest(c, "No employee profile linked to your account", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}

	filters := buildTaskFilters(c)

	tasks, total, err := h.service.GetMyTasks(employeeID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list tasks", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.SuccessWithMeta(c, "My tasks retrieved successfully", tasks, meta)
}

// GetMyTodayTasks returns today's tasks for the current employee
// @Summary Get my today's tasks
// @Tags Daily Tasks (Self-Service)
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/daily-tasks/today [get]
func (h *DailyTaskHandler) GetMyTodayTasks(c *gin.Context) {
	employeeID, _, ok := getEmployeeFromUser(c)
	if !ok || employeeID == "" {
		response.BadRequest(c, "No employee profile linked to your account", nil)
		return
	}

	tasks, err := h.service.GetMyTodayTasks(employeeID)
	if err != nil {
		response.InternalServerError(c, "Failed to get today's tasks", err.Error())
		return
	}

	response.Success(c, "Today's tasks retrieved successfully", tasks)
}

// GetMyTaskSummary returns task summary for the current employee
// @Summary Get my task summary
// @Tags Daily Tasks (Self-Service)
// @Produce json
// @Param date_from query string false "From date (YYYY-MM-DD), defaults to start of current month"
// @Param date_to query string false "To date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/self-service/daily-tasks/summary [get]
func (h *DailyTaskHandler) GetMyTaskSummary(c *gin.Context) {
	employeeID, _, ok := getEmployeeFromUser(c)
	if !ok || employeeID == "" {
		response.BadRequest(c, "No employee profile linked to your account", nil)
		return
	}

	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	to := now

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			from = parsed
		}
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		if parsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			to = parsed
		}
	}

	summary, err := h.service.GetMyTaskSummary(employeeID, from, to)
	if err != nil {
		response.InternalServerError(c, "Failed to get task summary", err.Error())
		return
	}

	response.Success(c, "Task summary retrieved successfully", summary)
}

// ────────────────────────── Helper Functions ──────────────────────────

func buildTaskFilters(c *gin.Context) map[string]interface{} {
	filters := make(map[string]interface{})

	if projectID := c.Query("project_id"); projectID != "" {
		if id, err := strconv.ParseUint(projectID, 10, 32); err == nil {
			filters["project_id"] = uint(id)
		}
	}
	if employeeID := c.Query("employee_id"); employeeID != "" {
		filters["employee_id"] = employeeID
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters["date_from"] = parsed
		}
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		if parsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters["date_to"] = parsed
		}
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	return filters
}

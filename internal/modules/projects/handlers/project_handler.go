package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// ProjectHandler handles HTTP requests for projects
type ProjectHandler struct {
	service *services.ProjectService
}

// NewProjectHandler creates a new ProjectHandler
func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		service: services.NewProjectService(),
	}
}

// ────────────────────────── Helper Functions ──────────────────────────

func getUserID(c *gin.Context) uint {
	user, exists := c.Get("user")
	if !exists {
		return 0
	}
	if userObj, ok := user.(*userModels.User); ok {
		return userObj.ID
	}
	return 0
}

func getUserName(c *gin.Context) string {
	user, exists := c.Get("user")
	if !exists {
		return ""
	}
	if userObj, ok := user.(*userModels.User); ok {
		return userObj.Username
	}
	return ""
}

func getEmployeeFromUser(c *gin.Context) (string, string, bool) {
	user, exists := c.Get("user")
	if !exists {
		return "", "", false
	}
	userObj, ok := user.(*userModels.User)
	if !ok {
		return "", "", false
	}
	// Look up employee by user ID using the service
	svc := services.NewDailyTaskService()
	empID, empName, err := svc.GetEmployeeByUserID(userObj.ID)
	if err != nil {
		return "", "", false
	}
	return empID, empName, true
}

// ────────────────────────── Project Endpoints ──────────────────────────

// CreateProject creates a new project
// @Summary Create a new project
// @Tags Projects
// @Accept json
// @Produce json
// @Param body body models.CreateProjectRequest true "Project data"
// @Success 201 {object} response.APIResponse
// @Router/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	userID := getUserID(c)
	userName := getUserName(c)

	project, err := h.service.CreateProject(&req, userID, userName)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Project created successfully", project)
}

// GetProject returns a project by ID
// @Summary Get project details
// @Tags Projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} response.APIResponse
// @Router/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid project ID", nil)
		return
	}

	project, err := h.service.GetProject(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Project retrieved successfully", project)
}

// UpdateProject updates a project
// @Summary Update a project
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param body body models.UpdateProjectRequest true "Updated project data"
// @Success 200 {object} response.APIResponse
// @Router/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid project ID", nil)
		return
	}

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	userID := getUserID(c)
	project, err := h.service.UpdateProject(uint(id), &req, &userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Project updated successfully", project)
}

// DeleteProject deletes a project
// @Summary Delete a project
// @Tags Projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} response.APIResponse
// @Router/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid project ID", nil)
		return
	}

	if err := h.service.DeleteProject(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Project deleted successfully", nil)
}

// ListProjects lists all projects with pagination and filters
// @Summary List projects
// @Tags Projects
// @Produce json
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param status query string false "Filter by status"
// @Param priority query string false "Filter by priority"
// @Param department_id query int false "Filter by department"
// @Param category query string false "Filter by category"
// @Param search query string false "Search by name/code"
// @Success 200 {object} response.APIResponse
// @Router/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if priority := c.Query("priority"); priority != "" {
		filters["priority"] = priority
	}
	if deptID := c.Query("department_id"); deptID != "" {
		if id, err := strconv.ParseUint(deptID, 10, 32); err == nil {
			filters["department_id"] = uint(id)
		}
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	projects, total, err := h.service.ListProjects(page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list projects", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.SuccessWithMeta(c, "Projects retrieved successfully", projects, meta)
}

// ────────────────────────── Member Endpoints ──────────────────────────

// AddMembers adds members to a project
// @Summary Add members to a project
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param body body models.AddMembersRequest true "Members to add"
// @Success 201 {object} response.APIResponse
// @Router/projects/{id}/members [post]
func (h *ProjectHandler) AddMembers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid project ID", nil)
		return
	}

	var req models.AddMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	userID := getUserID(c)
	added, err := h.service.AddMembers(uint(id), &req, userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Members added successfully", added)
}

// RemoveMember removes a member from a project
// @Summary Remove a member from a project
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param body body models.RemoveMemberRequest true "Member to remove"
// @Success 200 {object} response.APIResponse
// @Router/projects/{id}/members/remove [post]
func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid project ID", nil)
		return
	}

	var req models.RemoveMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	if err := h.service.RemoveMember(uint(id), req.EmployeeID); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Member removed successfully", nil)
}

// ListMembers lists all active members of a project
// @Summary List project members
// @Tags Projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} response.APIResponse
// @Router/projects/{id}/members [get]
func (h *ProjectHandler) ListMembers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid project ID", nil)
		return
	}

	members, err := h.service.ListMembers(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Members retrieved successfully", members)
}

// GetProjectProgress returns project progress details
// @Summary Get project progress
// @Tags Projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} response.APIResponse
// @Router/projects/{id}/progress [get]
func (h *ProjectHandler) GetProjectProgress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid project ID", nil)
		return
	}

	progress, err := h.service.GetProjectProgress(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Project progress retrieved successfully", progress)
}

// GetStatistics returns overall project statistics
// @Summary Get project statistics
// @Tags Projects
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router/projects/statistics [get]
func (h *ProjectHandler) GetStatistics(c *gin.Context) {
	_ = middleware.GetTenantID(c) // backward compatibility

	stats, err := h.service.GetProjectStatistics()
	if err != nil {
		response.InternalServerError(c, "Failed to get statistics", err.Error())
		return
	}

	response.Success(c, "Project statistics retrieved successfully", stats)
}

// ────────────────────────── Self-Service Endpoints ──────────────────────────

// ListMyProjects lists projects the current user is a member of
// @Summary List my projects
// @Tags Projects (Self-Service)
// @Produce json
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param status query string false "Filter by status"
// @Success 200 {object} response.APIResponse
// @Router/self-service/projects [get]
func (h *ProjectHandler) ListMyProjects(c *gin.Context) {
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

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	projects, total, err := h.service.ListMyProjects(employeeID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list projects", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.SuccessWithMeta(c, "My projects retrieved successfully", projects, meta)
}

// ────────────────────────── Assign / Unassign Endpoints ──────────────────────────

// AssignProject assigns a project to one or more employees
// @Summary Assign project to employees
// @Description Assign a project to one or more employees by their employee IDs. Employees become active project members.
// @Tags Projects
// @Accept json
// @Produce json
// @Param body body models.AssignProjectRequest true "Assignment data"
// @Success 201 {object} response.APIResponse
// @Router/projects/assign [post]
func (h *ProjectHandler) AssignProject(c *gin.Context) {
	var req models.AssignProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	userID := getUserID(c)

	// Build AddMembersRequest from AssignProjectRequest
	role := models.MemberRoleMember
	if req.Role != "" {
		role = req.Role
	}

	members := make([]models.AddMemberEntry, len(req.EmployeeIDs))
	for i, empID := range req.EmployeeIDs {
		members[i] = models.AddMemberEntry{
			EmployeeID: empID,
			Role:       role,
		}
	}

	addReq := &models.AddMembersRequest{Members: members}
	added, err := h.service.AddMembers(req.ProjectID, addReq, userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Fetch updated project with members
	project, _ := h.service.GetProject(req.ProjectID)

	response.Created(c, "Project assigned to employees successfully", gin.H{
		"project":       project,
		"members_added": added,
	})
}

// UnassignProject removes an employee from a project
// @Summary Unassign employee from project
// @Description Remove an employee from a project assignment.
// @Tags Projects
// @Accept json
// @Produce json
// @Param body body models.UnassignProjectRequest true "Unassignment data"
// @Success 200 {object} response.APIResponse
// @Router/projects/unassign [post]
func (h *ProjectHandler) UnassignProject(c *gin.Context) {
	var req models.UnassignProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	if err := h.service.RemoveMember(req.ProjectID, req.EmployeeID); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Employee unassigned from project successfully", nil)
}

// GetMyProjectDetails returns project details for a project the employee is assigned to
// @Summary Get assigned project details
// @Tags Projects (Self-Service)
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} response.APIResponse
// @Router/self-service/projects/{id} [get]
func (h *ProjectHandler) GetMyProjectDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid project ID", nil)
		return
	}

	employeeID, _, ok := getEmployeeFromUser(c)
	if !ok || employeeID == "" {
		response.BadRequest(c, "No employee profile linked to your account", nil)
		return
	}

	project, err := h.service.GetProject(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	// Verify employee is a member of the project
	isMember := false
	for _, m := range project.Members {
		if m.EmployeeID == employeeID && m.IsActive {
			isMember = true
			break
		}
	}
	if !isMember {
		response.Forbidden(c, "You are not assigned to this project")
		return
	}

	// Get progress info
	progress, _ := h.service.GetProjectProgress(uint(id))

	response.Success(c, "Project details retrieved successfully", gin.H{
		"project":  project,
		"progress": progress,
	})
}

package handlers

import (
	"strconv"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// EmployeeHandler handles employee HTTP requests
type EmployeeHandler struct {
	employeeService *services.EmployeeService
}

// NewEmployeeHandler creates a new employee handler
func NewEmployeeHandler() *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: services.NewEmployeeService(),
	}
}

// OnboardEmployee handles employee onboarding
// @Summary Onboard a new employee
// @Description Create a new employee record (onboarding)
// @Tags Employees
// @Accept json
// @Produce json
// @Param request body models.OnboardEmployeeRequest true "Employee data"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 422 {object} response.APIResponse
// @Router /api/v1/employees [post]
func (h *EmployeeHandler) OnboardEmployee(c *gin.Context) {
	var req models.OnboardEmployeeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	employee, err := h.employeeService.OnboardEmployee(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Employee onboarded successfully", employee)
}

// GetEmployee handles getting an employee by ID
// @Summary Get employee by ID
// @Description Get employee details by ID
// @Tags Employees
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/employees/{id} [get]
func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid employee ID", nil)
		return
	}

	employee, err := h.employeeService.GetEmployeeByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Employee retrieved successfully", employee)
}

// GetEmployeeByEmployeeID handles getting an employee by employee ID
// @Summary Get employee by employee ID
// @Description Get employee details by employee ID (employee number)
// @Tags Employees
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/employees/employee-id/{employee_id} [get]
func (h *EmployeeHandler) GetEmployeeByEmployeeID(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		response.BadRequest(c, "Employee ID is required", nil)
		return
	}

	employee, err := h.employeeService.GetEmployeeByEmployeeID(employeeID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Employee retrieved successfully", employee)
}

// UpdateEmployee handles updating an employee
// @Summary Update employee
// @Description Update employee information
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param request body models.UpdateEmployeeRequest true "Update data"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid employee ID", nil)
		return
	}

	var req models.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	employee, err := h.employeeService.UpdateEmployee(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Employee updated successfully", employee)
}

// ListEmployees handles listing employees with pagination
// @Summary List employees
// @Description Get paginated list of employees
// @Tags Employees
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees [get]
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Use the enhanced list method with all details
	employees, total, totalPages, err := h.employeeService.ListEmployeesWithDetails(page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list employees", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	response.SuccessWithMeta(c, "Employees retrieved successfully", employees, meta)
}

// ListByDepartment handles listing employees by department
// @Summary List employees by department
// @Description Get employees in a specific department
// @Tags Employees
// @Produce json
// @Param department_id path int true "Department ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/department/{department_id} [get]
func (h *EmployeeHandler) ListByDepartment(c *gin.Context) {
	idParam := c.Param("department_id")
	departmentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid department ID", nil)
		return
	}

	employees, err := h.employeeService.ListByDepartment(uint(departmentID))
	if err != nil {
		response.InternalServerError(c, err.Error(), nil)
		return
	}

	response.Success(c, "Employees retrieved successfully", employees)
}

// ListByStatus handles listing employees by status
// @Summary List employees by status
// @Description Get employees with a specific status
// @Tags Employees
// @Produce json
// @Param status path string true "Employee status" Enums(active, inactive, on_leave, terminated)
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/status/{status} [get]
func (h *EmployeeHandler) ListByStatus(c *gin.Context) {
	status := c.Param("status")
	if status == "" {
		response.BadRequest(c, "Status is required", nil)
		return
	}

	employees, err := h.employeeService.ListByStatus(status)
	if err != nil {
		response.InternalServerError(c, err.Error(), nil)
		return
	}

	response.Success(c, "Employees retrieved successfully", employees)
}

// DeleteEmployee handles deleting an employee
// @Summary Delete employee
// @Description Soft delete an employee
// @Tags Employees
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid employee ID", nil)
		return
	}

	if err := h.employeeService.DeleteEmployee(uint(id)); err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Employee deleted successfully", nil)
}

// TerminateEmployee handles terminating an employee (POST with ID in body)
func (h *EmployeeHandler) TerminateEmployee(c *gin.Context) {
	var req struct {
		ID              uint       `json:"id" binding:"required"`
		Reason          *string    `json:"reason,omitempty"`
		TerminationDate *time.Time `json:"termination_date,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	employee, err := h.employeeService.TerminateEmployee(req.ID, req.Reason, req.TerminationDate, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Employee terminated successfully", employee)
}

// SuspendEmployee handles suspending an employee (POST with ID in body)
func (h *EmployeeHandler) SuspendEmployee(c *gin.Context) {
	var req struct {
		ID                uint       `json:"id" binding:"required"`
		Reason            *string    `json:"reason,omitempty"`
		SuspensionEndDate *time.Time `json:"suspension_end_date,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	employee, err := h.employeeService.SuspendEmployee(req.ID, req.Reason, req.SuspensionEndDate, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Employee suspended successfully", employee)
}

// ArchiveEmployee handles archiving an employee (POST with ID in body)
func (h *EmployeeHandler) ArchiveEmployee(c *gin.Context) {
	var req struct {
		ID     uint    `json:"id" binding:"required"`
		Reason *string `json:"reason,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	employee, err := h.employeeService.ArchiveEmployee(req.ID, req.Reason, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Employee archived successfully", employee)
}

// ReactivateEmployee handles reactivating a suspended or archived employee (POST with ID in body)
func (h *EmployeeHandler) ReactivateEmployee(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	employee, err := h.employeeService.ReactivateEmployee(req.ID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Employee reactivated successfully", employee)
}

// ListManagers handles listing all potential reporting managers
// @Summary List reporting managers
// @Description Get list of all active employees who can be reporting managers.
//
//	Optionally filter by department_id — when provided, only managers in that department
//	are returned and the department head is flagged with is_department_head=true and is_suggested=true.
//
// @Tags Employees
// @Produce json
// @Param department_id query int false "Filter by department ID (also flags department head as suggested)"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/managers [get]
func (h *EmployeeHandler) ListManagers(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	// Optional department_id filter
	var departmentID *uint
	if deptStr := c.Query("department_id"); deptStr != "" {
		if deptID, err := strconv.ParseUint(deptStr, 10, 32); err == nil {
			dID := uint(deptID)
			departmentID = &dID
		}
	}

	managers, err := h.employeeService.ListManagers(tenantID, departmentID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Managers retrieved successfully", managers)
}

// GetDepartmentManager returns the suggested reporting manager (department head) for a department.
// @Summary Get suggested manager for a department
// @Description Returns the department head/manager who should be the default reporting manager
//
//	for new employees in this department. Returns null data if no manager is set.
//
// @Tags Employees
// @Produce json
// @Param department_id path int true "Department ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/department-manager/{department_id} [get]
func (h *EmployeeHandler) GetDepartmentManager(c *gin.Context) {
	departmentID, err := strconv.ParseUint(c.Param("department_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid department ID", nil)
		return
	}

	manager, err := h.employeeService.GetDepartmentManager(uint(departmentID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	if manager == nil {
		response.Success(c, "No manager is assigned to this department", nil)
		return
	}

	response.Success(c, "Department manager retrieved successfully", manager)
}

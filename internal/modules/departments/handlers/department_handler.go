package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/services"
	locationRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/types"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// DepartmentResponse represents the department response with location name instead of location_id
type DepartmentResponse struct {
	ID                uint           `json:"id"`
	TenantID          *uint          `json:"tenant_id,omitempty"`
	OrganizationID    *uint          `json:"organization_id,omitempty"`
	OrganizationUnitID *uint         `json:"organization_unit_id,omitempty"`
	Code              string         `json:"code"`
	Name              string         `json:"name"`
	Description       *string        `json:"description,omitempty"`
	Level             *string        `json:"level,omitempty"`
	DepartmentType    *string        `json:"department_type,omitempty"`
	ParentDepartmentID *uint         `json:"parent_department_id,omitempty"`
	ManagerID         *uint          `json:"manager_id,omitempty"`
	DeputyManager     *string        `json:"deputy_manager,omitempty"`
	BudgetAllocated   *float64       `json:"budget_allocated,omitempty"`
	BudgetCurrency    *string        `json:"budget_currency,omitempty"`
	EmployeeCapacity  *int           `json:"employee_capacity,omitempty"`
	Location          *string        `json:"location,omitempty"` // Location name instead of location_id
	CostCenter        *string        `json:"cost_center,omitempty"`
	IsActive          bool           `json:"is_active"`
	CreatedAt         string         `json:"created_at"`
	UpdatedAt         string         `json:"updated_at"`
	UpdatedBy         *uint          `json:"updated_by,omitempty"`
}

type DepartmentHandler struct {
	service     *services.DepartmentService
	locationRepo *locationRepos.LocationRepository
}

func NewDepartmentHandler() *DepartmentHandler {
	return &DepartmentHandler{
		service:     services.NewDepartmentService(),
		locationRepo: locationRepos.NewLocationRepository(),
	}
}

// toDepartmentResponse converts Department model to DepartmentResponse with location name
func (h *DepartmentHandler) toDepartmentResponse(dept *models.Department) *DepartmentResponse {
	resp := &DepartmentResponse{
		ID:                dept.ID,
		TenantID:          dept.TenantID,
		OrganizationID:    dept.OrganizationID,
		OrganizationUnitID: dept.OrganizationUnitID,
		Code:              dept.Code,
		Name:              dept.Name,
		Description:       dept.Description,
		Level:             dept.Level,
		DepartmentType:    dept.DepartmentType,
		ParentDepartmentID: dept.ParentDepartmentID,
		ManagerID:         dept.ManagerID,
		DeputyManager:     dept.DeputyManager,
		BudgetAllocated:   dept.BudgetAllocated,
		BudgetCurrency:    dept.BudgetCurrency,
		EmployeeCapacity:  dept.EmployeeCapacity,
		CostCenter:        dept.CostCenter,
		IsActive:          dept.IsActive,
		CreatedAt:         dept.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         dept.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedBy:         dept.UpdatedBy,
	}

	// Get location name if location_id exists
	if dept.LocationID != nil {
		location, err := h.locationRepo.FindByID(*dept.LocationID)
		if err == nil && location != nil {
			resp.Location = &location.Name
		}
	}

	return resp
}

// CreateDepartment handles department creation
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var req models.CreateDepartmentRequest
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

	department, err := h.service.CreateDepartment(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert to response format with location name
	departmentResponse := h.toDepartmentResponse(department)
	response.Created(c, "Department created successfully", departmentResponse)
}

// GetDepartment handles getting a department by ID
func (h *DepartmentHandler) GetDepartment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid department ID", nil)
		return
	}

	department, err := h.service.GetDepartmentByID(uint(id))
	if err != nil {
		response.NotFound(c, "Department not found")
		return
	}

	// Convert to response format with location name
	departmentResponse := h.toDepartmentResponse(department)
	response.Success(c, "Department retrieved successfully", departmentResponse)
}

// UpdateDepartment handles department updates
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid department ID", nil)
		return
	}

	var req models.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	department, err := h.service.UpdateDepartment(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert to response format with location name
	departmentResponse := h.toDepartmentResponse(department)
	response.Success(c, "Department updated successfully", departmentResponse)
}

// DeleteDepartment handles department deletion
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid department ID", nil)
		return
	}

	if err := h.service.DeleteDepartment(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Department deleted successfully", nil)
}

// ListDepartments handles listing departments
func (h *DepartmentHandler) ListDepartments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)
	
	// Build filters from query parameters
	filters := make(map[string]interface{})
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"] = isActive == "true"
	}
	if parentID := c.Query("parent_department_id"); parentID != "" {
		if parentID == "null" {
			filters["parent_department_id"] = nil
		} else {
			if pid, err := strconv.ParseUint(parentID, 10, 32); err == nil {
				filters["parent_department_id"] = uint(pid)
			}
		}
	}

	departments, total, err := h.service.ListDepartments(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list departments", err.Error())
		return
	}

	// Convert to response format with location names
	departmentResponses := make([]*DepartmentResponse, len(departments))
	for i, dept := range departments {
		departmentResponses[i] = h.toDepartmentResponse(&dept)
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.SuccessWithMeta(c, "Departments retrieved successfully", departmentResponses, meta)
}

// GetRootDepartments handles getting root departments
func (h *DepartmentHandler) GetRootDepartments(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	
	departments, err := h.service.GetRootDepartments(tenantID)
	if err != nil {
		response.InternalServerError(c, "Failed to get root departments", err.Error())
		return
	}

	// Convert to response format with location names
	departmentResponses := make([]*DepartmentResponse, len(departments))
	for i, dept := range departments {
		departmentResponses[i] = h.toDepartmentResponse(&dept)
	}

	response.Success(c, "Root departments retrieved successfully", departmentResponses)
}

// HandleAction handles POST-only action-based requests
func (h *DepartmentHandler) HandleAction(c *gin.Context) {
	var req types.APIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	switch req.Action {
	case "create":
		var createReq models.CreateDepartmentRequest
		if err := mapToStruct(req.Data, &createReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		department, err := h.service.CreateDepartment(&createReq, tenantID, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		// Convert to response format with location name
		departmentResponse := h.toDepartmentResponse(department)
		response.Created(c, "Department created successfully", departmentResponse)

	case "read":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for read action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid department ID", nil)
			return
		}
		department, err := h.service.GetDepartmentByID(uint(id))
		if err != nil {
			response.NotFound(c, "Department not found")
			return
		}
		// Convert to response format with location name
		departmentResponse := h.toDepartmentResponse(department)
		response.Success(c, "Department retrieved successfully", departmentResponse)

	case "update":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for update action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid department ID", nil)
			return
		}
		var updateReq models.UpdateDepartmentRequest
		if err := mapToStruct(req.Data, &updateReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		department, err := h.service.UpdateDepartment(uint(id), &updateReq, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		// Convert to response format with location name
		departmentResponse := h.toDepartmentResponse(department)
		response.Success(c, "Department updated successfully", departmentResponse)

	case "delete":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for delete action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid department ID", nil)
			return
		}
		if err := h.service.DeleteDepartment(uint(id)); err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Department deleted successfully", nil)

	case "list":
		page := 1
		pageSize := 20
		if req.Pagination != nil {
			page = req.Pagination.GetPage()
			pageSize = req.Pagination.GetPageSize()
		}

		filters := req.Filters
		if filters == nil {
			filters = make(map[string]interface{})
		}

		departments, total, err := h.service.ListDepartments(tenantID, page, pageSize, filters)
		if err != nil {
			response.InternalServerError(c, "Failed to list departments", err.Error())
			return
		}

		// Convert to response format with location names
		departmentResponses := make([]*DepartmentResponse, len(departments))
		for i, dept := range departments {
			departmentResponses[i] = h.toDepartmentResponse(&dept)
		}

		meta := &response.Meta{
			Page:       page,
			PerPage:    pageSize,
			Total:      total,
			TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		}

		response.SuccessWithMeta(c, "Departments retrieved successfully", departmentResponses, meta)

	default:
		response.BadRequest(c, "Invalid action. Must be one of: create, read, update, delete, list", nil)
	}
}

// Helper function to map map[string]interface{} to struct
func mapToStruct(data map[string]interface{}, target interface{}) error {
	// Simple implementation - in production, use a proper map-to-struct library
	// For now, we'll use JSON marshaling/unmarshaling
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, target)
}

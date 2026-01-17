package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/cost_centers/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/cost_centers/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/types"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type CostCenterHandler struct {
	service *services.CostCenterService
}

func NewCostCenterHandler() *CostCenterHandler {
	return &CostCenterHandler{
		service: services.NewCostCenterService(),
	}
}

// CreateCostCenter handles cost center creation
func (h *CostCenterHandler) CreateCostCenter(c *gin.Context) {
	var req models.CreateCostCenterRequest
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

	costCenter, err := h.service.CreateCostCenter(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Cost center created successfully", costCenter)
}

// GetCostCenter handles getting a cost center by ID
func (h *CostCenterHandler) GetCostCenter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid cost center ID", nil)
		return
	}

	costCenter, err := h.service.GetCostCenterByID(uint(id))
	if err != nil {
		response.NotFound(c, "Cost center not found")
		return
	}

	response.Success(c, "Cost center retrieved successfully", costCenter)
}

// UpdateCostCenter handles cost center updates
func (h *CostCenterHandler) UpdateCostCenter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid cost center ID", nil)
		return
	}

	var req models.UpdateCostCenterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	costCenter, err := h.service.UpdateCostCenter(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Cost center updated successfully", costCenter)
}

// DeleteCostCenter handles cost center deletion
func (h *CostCenterHandler) DeleteCostCenter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid cost center ID", nil)
		return
	}

	if err := h.service.DeleteCostCenter(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Cost center deleted successfully", nil)
}

// ListCostCenters handles listing cost centers
func (h *CostCenterHandler) ListCostCenters(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)
	var organizationID *uint
	if orgIDStr := c.Query("organization_id"); orgIDStr != "" {
		if orgID, err := strconv.ParseUint(orgIDStr, 10, 32); err == nil {
			orgIDUint := uint(orgID)
			organizationID = &orgIDUint
		}
	}
	
	// Build filters from query parameters
	filters := make(map[string]interface{})
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"] = isActive == "true"
	}
	if costCenterType := c.Query("cost_center_type"); costCenterType != "" {
		filters["cost_center_type"] = costCenterType
	}
	if departmentID := c.Query("department_id"); departmentID != "" {
		if deptID, err := strconv.ParseUint(departmentID, 10, 32); err == nil {
			filters["department_id"] = uint(deptID)
		}
	}
	if parentCostCenterID := c.Query("parent_cost_center_id"); parentCostCenterID != "" {
		if parentCostCenterID == "null" {
			filters["parent_cost_center_id"] = nil
		} else {
			if pid, err := strconv.ParseUint(parentCostCenterID, 10, 32); err == nil {
				filters["parent_cost_center_id"] = uint(pid)
			}
		}
	}

	costCenters, total, err := h.service.ListCostCenters(tenantID, organizationID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list cost centers", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.SuccessWithMeta(c, "Cost centers retrieved successfully", costCenters, meta)
}

// GetRootCostCenters handles getting root cost centers
func (h *CostCenterHandler) GetRootCostCenters(c *gin.Context) {
	var organizationID *uint
	if orgIDStr := c.Query("organization_id"); orgIDStr != "" {
		if orgID, err := strconv.ParseUint(orgIDStr, 10, 32); err == nil {
			orgIDUint := uint(orgID)
			organizationID = &orgIDUint
		}
	}
	
	costCenters, err := h.service.GetRootCostCenters(organizationID)
	if err != nil {
		response.InternalServerError(c, "Failed to get root cost centers", err.Error())
		return
	}

	response.Success(c, "Root cost centers retrieved successfully", costCenters)
}

// HandleAction handles POST-only action-based requests
func (h *CostCenterHandler) HandleAction(c *gin.Context) {
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

	var organizationID *uint
	if orgID, ok := req.Filters["organization_id"]; ok {
		if orgIDFloat, ok := orgID.(float64); ok {
			orgIDUint := uint(orgIDFloat)
			organizationID = &orgIDUint
		}
	}

	switch req.Action {
	case "create":
		var createReq models.CreateCostCenterRequest
		if err := mapToStruct(req.Data, &createReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		costCenter, err := h.service.CreateCostCenter(&createReq, tenantID, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Created(c, "Cost center created successfully", costCenter)

	case "read":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for read action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid cost center ID", nil)
			return
		}
		costCenter, err := h.service.GetCostCenterByID(uint(id))
		if err != nil {
			response.NotFound(c, "Cost center not found")
			return
		}
		response.Success(c, "Cost center retrieved successfully", costCenter)

	case "update":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for update action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid cost center ID", nil)
			return
		}
		var updateReq models.UpdateCostCenterRequest
		if err := mapToStruct(req.Data, &updateReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		costCenter, err := h.service.UpdateCostCenter(uint(id), &updateReq, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Cost center updated successfully", costCenter)

	case "delete":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for delete action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid cost center ID", nil)
			return
		}
		if err := h.service.DeleteCostCenter(uint(id)); err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Cost center deleted successfully", nil)

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

		costCenters, total, err := h.service.ListCostCenters(tenantID, organizationID, page, pageSize, filters)
		if err != nil {
			response.InternalServerError(c, "Failed to list cost centers", err.Error())
			return
		}

		meta := &response.Meta{
			Page:       page,
			PerPage:    pageSize,
			Total:      total,
			TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		}

		response.SuccessWithMeta(c, "Cost centers retrieved successfully", costCenters, meta)

	default:
		response.BadRequest(c, "Invalid action. Must be one of: create, read, update, delete, list", nil)
	}
}

// Helper function to map map[string]interface{} to struct
func mapToStruct(data map[string]interface{}, target interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, target)
}

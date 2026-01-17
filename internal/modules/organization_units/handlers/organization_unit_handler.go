package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization_units/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization_units/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/types"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type OrganizationUnitHandler struct {
	service *services.OrganizationUnitService
}

func NewOrganizationUnitHandler() *OrganizationUnitHandler {
	return &OrganizationUnitHandler{
		service: services.NewOrganizationUnitService(),
	}
}

// CreateOrganizationUnit handles organization unit creation
func (h *OrganizationUnitHandler) CreateOrganizationUnit(c *gin.Context) {
	var req models.CreateOrganizationUnitRequest
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

	unit, err := h.service.CreateOrganizationUnit(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Organization unit created successfully", unit)
}

// GetOrganizationUnit handles getting an organization unit by ID
func (h *OrganizationUnitHandler) GetOrganizationUnit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid organization unit ID", nil)
		return
	}

	unit, err := h.service.GetOrganizationUnitByID(uint(id))
	if err != nil {
		response.NotFound(c, "Organization unit not found")
		return
	}

	response.Success(c, "Organization unit retrieved successfully", unit)
}

// UpdateOrganizationUnit handles organization unit updates
func (h *OrganizationUnitHandler) UpdateOrganizationUnit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid organization unit ID", nil)
		return
	}

	var req models.UpdateOrganizationUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	unit, err := h.service.UpdateOrganizationUnit(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Organization unit updated successfully", unit)
}

// DeleteOrganizationUnit handles organization unit deletion
func (h *OrganizationUnitHandler) DeleteOrganizationUnit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid organization unit ID", nil)
		return
	}

	if err := h.service.DeleteOrganizationUnit(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Organization unit deleted successfully", nil)
}

// ListOrganizationUnits handles listing organization units
func (h *OrganizationUnitHandler) ListOrganizationUnits(c *gin.Context) {
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
	if unitType := c.Query("unit_type"); unitType != "" {
		filters["unit_type"] = unitType
	}
	if parentUnitID := c.Query("parent_unit_id"); parentUnitID != "" {
		if parentUnitID == "null" {
			filters["parent_unit_id"] = nil
		} else {
			if pid, err := strconv.ParseUint(parentUnitID, 10, 32); err == nil {
				filters["parent_unit_id"] = uint(pid)
			}
		}
	}

	units, total, err := h.service.ListOrganizationUnits(tenantID, organizationID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list organization units", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.SuccessWithMeta(c, "Organization units retrieved successfully", units, meta)
}

// GetRootOrganizationUnits handles getting root organization units
func (h *OrganizationUnitHandler) GetRootOrganizationUnits(c *gin.Context) {
	var organizationID *uint
	if orgIDStr := c.Query("organization_id"); orgIDStr != "" {
		if orgID, err := strconv.ParseUint(orgIDStr, 10, 32); err == nil {
			orgIDUint := uint(orgID)
			organizationID = &orgIDUint
		}
	}
	
	units, err := h.service.GetRootOrganizationUnits(organizationID)
	if err != nil {
		response.InternalServerError(c, "Failed to get root organization units", err.Error())
		return
	}

	response.Success(c, "Root organization units retrieved successfully", units)
}

// HandleAction handles POST-only action-based requests
func (h *OrganizationUnitHandler) HandleAction(c *gin.Context) {
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
		var createReq models.CreateOrganizationUnitRequest
		if err := mapToStruct(req.Data, &createReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		unit, err := h.service.CreateOrganizationUnit(&createReq, tenantID, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Created(c, "Organization unit created successfully", unit)

	case "read":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for read action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid organization unit ID", nil)
			return
		}
		unit, err := h.service.GetOrganizationUnitByID(uint(id))
		if err != nil {
			response.NotFound(c, "Organization unit not found")
			return
		}
		response.Success(c, "Organization unit retrieved successfully", unit)

	case "update":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for update action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid organization unit ID", nil)
			return
		}
		var updateReq models.UpdateOrganizationUnitRequest
		if err := mapToStruct(req.Data, &updateReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		unit, err := h.service.UpdateOrganizationUnit(uint(id), &updateReq, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Organization unit updated successfully", unit)

	case "delete":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for delete action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid organization unit ID", nil)
			return
		}
		if err := h.service.DeleteOrganizationUnit(uint(id)); err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Organization unit deleted successfully", nil)

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

		units, total, err := h.service.ListOrganizationUnits(tenantID, organizationID, page, pageSize, filters)
		if err != nil {
			response.InternalServerError(c, "Failed to list organization units", err.Error())
			return
		}

		meta := &response.Meta{
			Page:       page,
			PerPage:    pageSize,
			Total:      total,
			TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		}

		response.SuccessWithMeta(c, "Organization units retrieved successfully", units, meta)

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

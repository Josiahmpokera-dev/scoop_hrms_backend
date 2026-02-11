package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/types"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	service *services.OrganizationService
}

func NewOrganizationHandler() *OrganizationHandler {
	return &OrganizationHandler{
		service: services.NewOrganizationService(),
	}
}

// CreateOrganization handles organization creation
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req models.CreateOrganizationRequest
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

	organization, err := h.service.CreateOrganization(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Organization created successfully", organization)
}

// GetOrganization handles getting an organization by ID
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	idStr := c.Param("id")
	
	// Handle "me" case - redirect to GetMyOrganization
	if idStr == "me" {
		h.GetMyOrganization(c)
		return
	}
	
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID", nil)
		return
	}

	organization, err := h.service.GetOrganizationByID(uint(id))
	if err != nil {
		response.NotFound(c, "Organization not found")
		return
	}

	response.Success(c, "Organization retrieved successfully", organization)
}

// UpdateOrganization handles organization updates
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID", nil)
		return
	}

	var req models.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	organization, err := h.service.UpdateOrganization(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Organization updated successfully", organization)
}

// DeleteOrganization handles organization deletion
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID", nil)
		return
	}

	if err := h.service.DeleteOrganization(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Organization deleted successfully", nil)
}

// GetMyOrganization handles getting the current user's organization
// This is useful after onboarding to get the organization that was created
func (h *OrganizationHandler) GetMyOrganization(c *gin.Context) {
	// Get tenant_id from context
	tenantID := middleware.GetTenantID(c)

	if tenantID == nil {
		response.BadRequest(c, "User is not associated with a tenant. Please complete onboarding first.", nil)
		return
	}

	// Get the first organization for this tenant (usually the one created during onboarding)
	organizations, total, err := h.service.ListOrganizations(tenantID, 1, 1, nil)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve organization", err.Error())
		return
	}

	if total == 0 || len(organizations) == 0 {
		response.NotFound(c, "No organization found for your account. Please create one first.")
		return
	}

	response.Success(c, "Organization retrieved successfully", organizations[0])
}

// ListOrganizations handles listing organizations
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
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
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if country := c.Query("country"); country != "" {
		filters["country"] = country
	}

	organizations, total, err := h.service.ListOrganizations(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list organizations", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.SuccessWithMeta(c, "Organizations retrieved successfully", organizations, meta)
}

// HandleAction handles POST-only action-based requests
func (h *OrganizationHandler) HandleAction(c *gin.Context) {
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
		var createReq models.CreateOrganizationRequest
		if err := mapToStruct(req.Data, &createReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		organization, err := h.service.CreateOrganization(&createReq, tenantID, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Created(c, "Organization created successfully", organization)

	case "read":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for read action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid organization ID", nil)
			return
		}
		organization, err := h.service.GetOrganizationByID(uint(id))
		if err != nil {
			response.NotFound(c, "Organization not found")
			return
		}
		response.Success(c, "Organization retrieved successfully", organization)

	case "update":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for update action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid organization ID", nil)
			return
		}
		var updateReq models.UpdateOrganizationRequest
		if err := mapToStruct(req.Data, &updateReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		organization, err := h.service.UpdateOrganization(uint(id), &updateReq, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Organization updated successfully", organization)

	case "delete":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for delete action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid organization ID", nil)
			return
		}
		if err := h.service.DeleteOrganization(uint(id)); err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Organization deleted successfully", nil)

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

		organizations, total, err := h.service.ListOrganizations(tenantID, page, pageSize, filters)
		if err != nil {
			response.InternalServerError(c, "Failed to list organizations", err.Error())
			return
		}

		meta := &response.Meta{
			Page:       page,
			PerPage:    pageSize,
			Total:      total,
			TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		}

		response.SuccessWithMeta(c, "Organizations retrieved successfully", organizations, meta)

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

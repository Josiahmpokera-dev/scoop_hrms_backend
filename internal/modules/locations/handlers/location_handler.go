package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/services"
	organizationServices "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/types"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	service *services.LocationService
}

func NewLocationHandler() *LocationHandler {
	return &LocationHandler{
		service: services.NewLocationService(),
	}
}

// CreateLocation handles location creation
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var req models.CreateLocationRequest
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

	location, err := h.service.CreateLocation(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Location created successfully", location)
}

// GetLocation handles getting a location by ID
func (h *LocationHandler) GetLocation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid location ID", nil)
		return
	}

	location, err := h.service.GetLocationByID(uint(id))
	if err != nil {
		response.NotFound(c, "Location not found")
		return
	}

	response.Success(c, "Location retrieved successfully", location)
}

// UpdateLocation handles location updates
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid location ID", nil)
		return
	}

	var req models.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	location, err := h.service.UpdateLocation(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Location updated successfully", location)
}

// DeleteLocation handles location deletion
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid location ID", nil)
		return
	}

	if err := h.service.DeleteLocation(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Location deleted successfully", nil)
}

// ListLocations handles listing locations
func (h *LocationHandler) ListLocations(c *gin.Context) {
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
	if isHeadOffice := c.Query("is_head_office"); isHeadOffice != "" {
		filters["is_head_office"] = isHeadOffice == "true"
	}
	if locationType := c.Query("location_type"); locationType != "" {
		filters["location_type"] = locationType
	}
	if country := c.Query("country"); country != "" {
		filters["country"] = country
	}
	if city := c.Query("city"); city != "" {
		filters["city"] = city
	}

	locations, total, err := h.service.ListLocations(tenantID, organizationID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list locations", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.SuccessWithMeta(c, "Locations retrieved successfully", locations, meta)
}

// GetHeadOffice handles getting the head office location
func (h *LocationHandler) GetHeadOffice(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	
	location, err := h.service.GetHeadOffice(tenantID)
	if err != nil {
		response.NotFound(c, "Head office location not found")
		return
	}

	response.Success(c, "Head office location retrieved successfully", location)
}

// GetMyOrganizationLocations handles getting locations for the current user's organization
func (h *LocationHandler) GetMyOrganizationLocations(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	
	// Get user's organization (first organization for the tenant)
	orgService := organizationServices.NewOrganizationService()
	organizations, total, err := orgService.ListOrganizations(tenantID, 1, 1, nil)
	if err != nil || total == 0 || len(organizations) == 0 {
		response.NotFound(c, "No organization found for your account. Please create one first.")
		return
	}

	// Get organization ID (use the first organization)
	organizationID := organizations[0].ID

	// Get locations for this organization
	locations, _, err := h.service.ListLocations(tenantID, &organizationID, 1, 100, map[string]interface{}{})
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve locations", err.Error())
		return
	}

	response.Success(c, "Locations retrieved successfully", locations)
}

// HandleAction handles POST-only action-based requests
func (h *LocationHandler) HandleAction(c *gin.Context) {
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
		var createReq models.CreateLocationRequest
		if err := mapToStruct(req.Data, &createReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		location, err := h.service.CreateLocation(&createReq, tenantID, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Created(c, "Location created successfully", location)

	case "read":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for read action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid location ID", nil)
			return
		}
		location, err := h.service.GetLocationByID(uint(id))
		if err != nil {
			response.NotFound(c, "Location not found")
			return
		}
		response.Success(c, "Location retrieved successfully", location)

	case "update":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for update action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid location ID", nil)
			return
		}
		var updateReq models.UpdateLocationRequest
		if err := mapToStruct(req.Data, &updateReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		location, err := h.service.UpdateLocation(uint(id), &updateReq, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Location updated successfully", location)

	case "delete":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for delete action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid location ID", nil)
			return
		}
		if err := h.service.DeleteLocation(uint(id)); err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Location deleted successfully", nil)

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

		var organizationID *uint
		if orgID, ok := filters["organization_id"]; ok {
			if orgIDFloat, ok := orgID.(float64); ok {
				orgIDUint := uint(orgIDFloat)
				organizationID = &orgIDUint
			}
		}

		locations, total, err := h.service.ListLocations(tenantID, organizationID, page, pageSize, filters)
		if err != nil {
			response.InternalServerError(c, "Failed to list locations", err.Error())
			return
		}

		meta := &response.Meta{
			Page:       page,
			PerPage:    pageSize,
			Total:      total,
			TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		}

		response.SuccessWithMeta(c, "Locations retrieved successfully", locations, meta)

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

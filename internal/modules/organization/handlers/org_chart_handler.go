package handlers

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// OrgChartHandler handles organization chart HTTP requests
type OrgChartHandler struct {
	service *services.OrgChartService
}

// NewOrgChartHandler creates a new organization chart handler
func NewOrgChartHandler() *OrgChartHandler {
	return &OrgChartHandler{
		service: services.NewOrgChartService(),
	}
}

// GetOrganizationChart handles getting the organization chart
func (h *OrgChartHandler) GetOrganizationChart(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	chart, err := h.service.GetOrganizationChart(tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Organization chart retrieved successfully", chart)
}

// GetOrganizationChartDebug provides diagnostic information about the organization chart
func (h *OrgChartHandler) GetOrganizationChartDebug(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	db := database.GetDB()

	// Get root positions
	var rootPositions []struct {
		ID                  uint
		Code                string
		Title               string
		ReportsToPositionID *uint
		IsActive            bool
		TenantID            *uint
	}

	rootQuery := db.Table("job_positions").
		Select("id, code, title, reports_to_position_id, is_active, tenant_id").
		Where("reports_to_position_id IS NULL").
		Where("deleted_at IS NULL")

	if tenantID != nil {
		rootQuery = rootQuery.Where("tenant_id = ?", *tenantID)
	}

	rootQuery.Find(&rootPositions)

	// For each root position, get its children
	diagnostics := make([]map[string]interface{}, 0)

	for _, root := range rootPositions {
		// Get all child positions (any status)
		var allChildren []struct {
			ID                  uint
			Code                string
			Title               string
			ReportsToPositionID *uint
			IsActive            bool
			TenantID            *uint
			DeletedAt           *string
		}

		childQuery := db.Table("job_positions").
			Select("id, code, title, reports_to_position_id, is_active, tenant_id, deleted_at").
			Where("reports_to_position_id = ?", root.ID)

		if tenantID != nil {
			childQuery = childQuery.Where("tenant_id = ?", *tenantID)
		}

		childQuery.Find(&allChildren)

		// Count active children
		activeChildren := 0
		for _, child := range allChildren {
			if child.IsActive && (child.DeletedAt == nil || *child.DeletedAt == "") {
				activeChildren++
			}
		}

		// For each child, check if it has its own children (grandchildren of root)
		childrenWithDetails := make([]map[string]interface{}, 0)
		for _, child := range allChildren {
			// Check if this child has its own children
			var grandChildren []struct {
				ID                  uint
				Code                string
				Title               string
				ReportsToPositionID *uint
				IsActive            bool
				TenantID            *uint
				DeletedAt           *string
			}

			grandChildQuery := db.Table("job_positions").
				Select("id, code, title, reports_to_position_id, is_active, tenant_id, deleted_at").
				Where("reports_to_position_id = ?", child.ID)

			if tenantID != nil {
				grandChildQuery = grandChildQuery.Where("tenant_id = ?", *tenantID)
			}

			grandChildQuery.Find(&grandChildren)

			activeGrandChildren := 0
			for _, gc := range grandChildren {
				if gc.IsActive && (gc.DeletedAt == nil || *gc.DeletedAt == "") {
					activeGrandChildren++
				}
			}

			childDetail := map[string]interface{}{
				"id":                     child.ID,
				"code":                   child.Code,
				"title":                  child.Title,
				"reports_to_position_id": child.ReportsToPositionID,
				"is_active":              child.IsActive,
				"total_children":         len(grandChildren),
				"active_children":        activeGrandChildren,
				"inactive_children":      len(grandChildren) - activeGrandChildren,
			}

			if len(grandChildren) > 0 {
				childDetail["children"] = grandChildren
			}

			childrenWithDetails = append(childrenWithDetails, childDetail)
		}

		diagnostics = append(diagnostics, map[string]interface{}{
			"root_position": map[string]interface{}{
				"id":                     root.ID,
				"code":                   root.Code,
				"title":                  root.Title,
				"is_active":              root.IsActive,
				"reports_to_position_id": root.ReportsToPositionID,
			},
			"total_children":    len(allChildren),
			"active_children":   activeChildren,
			"inactive_children": len(allChildren) - activeChildren,
			"children":         childrenWithDetails, // Now includes grandchildren info
		})
	}

	response.Success(c, "Organization chart diagnostics", map[string]interface{}{
		"tenant_id":      tenantID,
		"root_positions": len(rootPositions),
		"diagnostics":    diagnostics,
		"help": map[string]string{
			"issue": "If active_children is 0 but total_children > 0, your child positions are inactive. Set is_active=true.",
			"fix":   "Update child positions with: POST /api/v1/job-positions/update with {\"id\": X, \"is_active\": true}",
		},
	})
}

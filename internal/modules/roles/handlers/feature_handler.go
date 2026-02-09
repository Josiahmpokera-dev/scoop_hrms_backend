package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// FeatureHandler handles feature access HTTP requests
type FeatureHandler struct {
	featureService *services.FeatureService
}

// NewFeatureHandler creates a new feature handler
func NewFeatureHandler() *FeatureHandler {
	return &FeatureHandler{
		featureService: services.NewFeatureService(),
	}
}

// GetAllFeatures returns the complete feature catalog with all modules and permissions.
// This is for admin UI to display all available features when configuring roles.
//
// @Summary     List all features
// @Description Get the complete feature catalog with all modules and permissions
// @Tags        Features
// @Produce     json
// @Success     200 {object} response.APIResponse
// @Router      /api/v1/features [get]
func (h *FeatureHandler) GetAllFeatures(c *gin.Context) {
	features := h.featureService.GetAllFeatures()

	// Group by category for easier frontend rendering
	categories := make(map[string][]interface{})
	for _, f := range features {
		categories[f.Category] = append(categories[f.Category], f)
	}

	response.Success(c, "Feature catalog retrieved successfully", gin.H{
		"features":   features,
		"categories": categories,
		"total":      len(features),
	})
}

// GetMyFeatureAccess returns the feature access map for the currently authenticated user.
// The frontend uses this to determine which menus, pages, and actions to show/hide.
//
// @Summary     Get my feature access
// @Description Get feature access permissions for the currently authenticated user
// @Tags        Features
// @Produce     json
// @Success     200 {object} response.APIResponse
// @Router      /api/v1/features/my-access [get]
func (h *FeatureHandler) GetMyFeatureAccess(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "Authentication required")
		return
	}

	userObj, ok := user.(*userModels.User)
	if !ok {
		response.Unauthorized(c, "Invalid user context")
		return
	}

	featureAccess, err := h.featureService.GetMyFeatureAccess(userObj.ID)
	if err != nil {
		response.InternalServerError(c, "Failed to get feature access", err.Error())
		return
	}

	// Build a flat permission map for quick lookups on the frontend
	permissionMap := make(map[string]bool)
	accessibleFeatures := make([]string, 0)
	for _, fa := range featureAccess {
		if fa.HasAccess {
			accessibleFeatures = append(accessibleFeatures, fa.Code)
		}
		for _, pa := range fa.Permissions {
			permissionMap[pa.Code] = pa.Granted
		}
	}

	response.Success(c, "Feature access retrieved successfully", gin.H{
		"features":            featureAccess,
		"permissions":         permissionMap,
		"accessible_features": accessibleFeatures,
		"user_id":             userObj.ID,
		"user_role":           string(userObj.Role),
	})
}

// GetRoleFeatureAccess returns the feature access map for a specific role.
// Used by admins to view/configure what a role can access.
//
// @Summary     Get role feature access
// @Description Get feature access permissions for a specific role
// @Tags        Features
// @Produce     json
// @Param       id path int true "Role ID"
// @Success     200 {object} response.APIResponse
// @Router      /api/v1/features/roles/{id} [get]
func (h *FeatureHandler) GetRoleFeatureAccess(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid role ID", nil)
		return
	}

	role, err := h.featureService.GetRoleByID(uint(roleID))
	if err != nil {
		response.NotFound(c, "Role not found")
		return
	}

	featureAccess, err := h.featureService.GetRoleFeatureAccess(uint(roleID))
	if err != nil {
		response.InternalServerError(c, "Failed to get role feature access", err.Error())
		return
	}

	// Build a flat permission map
	permissionMap := make(map[string]bool)
	accessibleFeatures := make([]string, 0)
	for _, fa := range featureAccess {
		if fa.HasAccess {
			accessibleFeatures = append(accessibleFeatures, fa.Code)
		}
		for _, pa := range fa.Permissions {
			permissionMap[pa.Code] = pa.Granted
		}
	}

	response.Success(c, "Role feature access retrieved successfully", gin.H{
		"role": gin.H{
			"id":          role.ID,
			"code":        role.Code,
			"name":        role.Name,
			"description": role.Description,
		},
		"features":            featureAccess,
		"permissions":         permissionMap,
		"accessible_features": accessibleFeatures,
	})
}

// UpdateRolePermissions updates the permissions assigned to a role.
// Accepts a list of permission codes to replace the current permissions.
//
// @Summary     Update role permissions
// @Description Update the permissions assigned to a specific role
// @Tags        Features
// @Accept      json
// @Produce     json
// @Param       id   path int                      true "Role ID"
// @Param       body body UpdateRolePermissionsReq  true "Permission codes"
// @Success     200 {object} response.APIResponse
// @Router      /api/v1/features/roles/{id} [put]
func (h *FeatureHandler) UpdateRolePermissions(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid role ID", nil)
		return
	}

	var req UpdateRolePermissionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	if len(req.Permissions) == 0 {
		response.ValidationError(c, "Validation failed", "At least one permission is required")
		return
	}

	if err := h.featureService.UpdateRolePermissions(uint(roleID), req.Permissions); err != nil {
		response.BadRequest(c, "Failed to update role permissions", err.Error())
		return
	}

	// Return updated feature access
	featureAccess, err := h.featureService.GetRoleFeatureAccess(uint(roleID))
	if err != nil {
		response.Success(c, "Role permissions updated successfully", nil)
		return
	}

	response.Success(c, "Role permissions updated successfully", gin.H{
		"features":          featureAccess,
		"permissions_count": len(req.Permissions),
	})
}

// UpdateRolePermissionsReq is the request body for updating role permissions
type UpdateRolePermissionsReq struct {
	Permissions []string `json:"permissions" binding:"required"`
}

package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// RoleHandler handles role HTTP requests.
type RoleHandler struct {
	roleService *services.RoleService
}

// NewRoleHandler creates a new role handler.
func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: services.NewRoleService(),
	}
}

// ListRoles returns a paginated list of RBAC roles.
//
// @Summary List roles
// @Description Get paginated list of roles (RBAC)
// @Tags Roles
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20) maximum(100)
// @Success 200 {object} response.APIResponse
// @Router/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	roles, total, err := h.roleService.ListRoles(tenantID, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list roles", err.Error())
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Roles retrieved successfully", roles, meta)
}

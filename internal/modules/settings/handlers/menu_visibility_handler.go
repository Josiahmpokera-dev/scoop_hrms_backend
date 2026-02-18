package handlers

import (
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/settings/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/settings/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type MenuVisibilityHandler struct {
	service *services.MenuVisibilityService
}

func NewMenuVisibilityHandler() *MenuVisibilityHandler {
	return &MenuVisibilityHandler{service: services.NewMenuVisibilityService()}
}

// GetAll returns the full menu tree with visibility flags + summary.
func (h *MenuVisibilityHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetFullTree()
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve menu visibility", err.Error())
		return
	}
	response.Success(c, "Menu visibility retrieved successfully", data)
}

// BulkUpdate applies visibility changes for multiple items.
func (h *MenuVisibilityHandler) BulkUpdate(c *gin.Context) {
	var req models.BulkUpdateVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	callerID := userID.(uint)

	items, err := h.service.BulkUpdate(req.Items, callerID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	msg := fmt.Sprintf("Menu visibility updated for %d items", len(items))
	response.Success(c, msg, map[string]interface{}{
		"updated_count": len(items),
		"items":         items,
	})
}

// ToggleSingle toggles visibility of a single menu item.
func (h *MenuVisibilityHandler) ToggleSingle(c *gin.Context) {
	menuKey := c.Param("menuKey")

	var req models.ToggleVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	callerID := userID.(uint)

	cascaded, err := h.service.ToggleSingle(menuKey, req.Visible, callerID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	state := "visible"
	if !req.Visible {
		state = "hidden"
	}
	extra := ""
	if len(cascaded) > 1 {
		extra = fmt.Sprintf(" and %d children", len(cascaded)-1)
	}
	msg := fmt.Sprintf("Menu item '%s'%s is now %s", menuKey, extra, state)

	response.Success(c, msg, map[string]interface{}{
		"menu_key": menuKey,
		"visible":  req.Visible,
		"cascaded": cascaded,
	})
}

// Reset resets all menu items to visible.
func (h *MenuVisibilityHandler) Reset(c *gin.Context) {
	count, err := h.service.ResetAll()
	if err != nil {
		response.InternalServerError(c, "Failed to reset menu visibility", err.Error())
		return
	}
	msg := fmt.Sprintf("Menu visibility reset to defaults. All %d items are now visible.", count)
	response.Success(c, msg, map[string]interface{}{
		"reset_count": count,
	})
}

// GetActive returns the list of hidden menu keys (for frontend filtering).
func (h *MenuVisibilityHandler) GetActive(c *gin.Context) {
	keys, err := h.service.GetHiddenKeys()
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve hidden keys", err.Error())
		return
	}
	if keys == nil {
		keys = []string{}
	}
	response.Success(c, "Hidden menu keys retrieved successfully", map[string]interface{}{
		"hidden_keys": keys,
	})
}

package handlers

import (
	"strconv"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user management HTTP requests.
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler.
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// ListUsers returns a paginated list of users for transfer-role and other admin/HR use.
// Optional filters: role (admin, hr, it, user), search (email, first_name, last_name, username).
//
// @Summary List users
// @Description Get paginated list of users. Use role=user to get normal employees for transfer-role target.
// @Tags Users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20) maximum(100)
// @Param role query string false "Filter by role: admin, hr, it, user"
// @Param search query string false "Search by email, first name, last name, username"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
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
	role := strings.TrimSpace(c.Query("role"))
	search := strings.TrimSpace(c.Query("search"))

	users, total, err := h.userService.ListUsers(tenantID, page, pageSize, role, search)
	if err != nil {
		response.InternalServerError(c, "Failed to list users", err.Error())
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Users retrieved successfully", users, meta)
}

// ListSpecialRoleUsers returns users who have role admin, hr, or it (for transfer-role "from" picker).
//
// @Summary List IT, HR and Admin users
// @Description Get paginated list of users with role admin, hr, or it
// @Tags Users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20) maximum(100)
// @Param search query string false "Search by email, first name, last name, username"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/users/special-roles [get]
func (h *UserHandler) ListSpecialRoleUsers(c *gin.Context) {
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
	search := strings.TrimSpace(c.Query("search"))

	users, total, err := h.userService.ListSpecialRoleUsers(tenantID, page, pageSize, search)
	if err != nil {
		response.InternalServerError(c, "Failed to list special role users", err.Error())
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Special role users (IT, HR, Admin) retrieved successfully", users, meta)
}

// TransferRole transfers a special role (admin, hr, it) to another user.
// Admin can transfer any role (must provide from_user_id). Role holder can transfer their own role to another user.
//
// @Summary Transfer role to another user
// @Description Transfer admin, hr, or it role from one user to another. Admin must provide from_user_id; role holder transfers their own role.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.TransferRoleRequest true "Transfer role request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /api/v1/users/transfer-role [post]
func (h *UserHandler) TransferRole(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)

	var req models.TransferRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	// When caller is not admin, they can only transfer their own role; from_user_id is ignored
	result, err := h.userService.TransferRole(callerID, req.Role, req.TargetUserID, req.FromUserID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Role transferred successfully", gin.H{
		"from_user":   result.FromUser,
		"target_user": result.TargetUser,
		"role":        result.Role,
	})
}

// AssignRole assigns a special role (admin, hr, or it) to a normal user.
// Only an Admin can call this. The target user must have role "user".
//
// @Summary Assign role to a user
// @Description Assign admin, hr, or it role to a normal user. Admin only. Target must have role "user".
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.AssignRoleRequest true "Assign role request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /api/v1/users/assign-role [post]
func (h *UserHandler) AssignRole(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)

	var req models.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	result, err := h.userService.AssignRole(callerID, req.UserID, req.Role, req.PositionID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	resp := gin.H{"user": result.User, "role": result.Role}
	if result.PositionUpdated {
		resp["position_updated"] = true
	}
	response.Success(c, "Role assigned successfully", resp)
}

// parseUserIDParam returns the user ID from the URL param "id". Returns 0 and false if invalid.
func parseUserIDParam(c *gin.Context) (uint, bool) {
	idStr := c.Param("id")
	if idStr == "" {
		return 0, false
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint(id), true
}

// SuspendUser suspends a user; they cannot login until unblocked. HR/Admin only.
//
// @Summary Suspend user
// @Description Set user status to suspended. User cannot login. HR/Admin only.
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/users/{id}/suspend [post]
func (h *UserHandler) SuspendUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)
	targetID, ok := parseUserIDParam(c)
	if !ok {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}
	user, err := h.userService.SuspendUser(callerID, targetID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "User suspended successfully", user)
}

// BlockUser blocks a user; they cannot login until unblocked. HR/Admin only.
//
// @Summary Block user
// @Description Set user status to blocked. User cannot login. HR/Admin only.
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/users/{id}/block [post]
func (h *UserHandler) BlockUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)
	targetID, ok := parseUserIDParam(c)
	if !ok {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}
	user, err := h.userService.BlockUser(callerID, targetID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "User blocked successfully", user)
}

// UnblockUser restores user status to active so they can login again. HR/Admin only.
//
// @Summary Unblock user
// @Description Set user status to active (restores login). Use for suspended or blocked users. HR/Admin only.
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/users/{id}/unblock [post]
func (h *UserHandler) UnblockUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)
	targetID, ok := parseUserIDParam(c)
	if !ok {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}
	user, err := h.userService.UnblockUser(callerID, targetID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "User unblocked successfully", user)
}

// UnsuspendUser restores user status to active (for suspended users). HR/Admin only.
//
// @Summary Unsuspend user
// @Description Set user status to active (restores login). Use for suspended users. HR/Admin only.
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/users/{id}/unsuspend [post]
func (h *UserHandler) UnsuspendUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)
	targetID, ok := parseUserIDParam(c)
	if !ok {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}
	user, err := h.userService.UnsuspendUser(callerID, targetID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "User unsuspended successfully", user)
}

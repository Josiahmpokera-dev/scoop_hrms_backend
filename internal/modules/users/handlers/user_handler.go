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

// CreateUser creates a new system user account linked to an existing employee.
// The user's name and email are taken from the employee record automatically.
// The employee must exist and must not already have a user account.
//
// @Summary Create user from employee
// @Description Create a login account for an existing employee with a specified role. Admin or HR only.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "Create user request"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	result, err := h.userService.CreateUser(callerID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "User created successfully", gin.H{
		"user":        result.User,
		"employee_id": result.EmployeeID,
		"roles":       result.Roles,
		"credentials": gin.H{
			"email":    result.Email,
			"password": result.Password,
		},
	})
}

// ListEmployeesWithoutUser returns employees that don't have a login account yet.
// Use this to populate the employee picker when creating a new user.
//
// @Summary List employees without user accounts
// @Description Get employees that can have a user account created. Admin or HR only.
// @Tags Users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20) maximum(100)
// @Param search query string false "Search by name, employee ID, or email"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/users/available-employees [get]
func (h *UserHandler) ListEmployeesWithoutUser(c *gin.Context) {
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

	employees, total, err := h.userService.ListEmployeesWithoutUser(page, pageSize, search)
	if err != nil {
		response.InternalServerError(c, "Failed to list employees", err.Error())
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Employees without user accounts retrieved successfully", employees, meta)
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

// AddRole adds an additional role to a user without removing existing roles.
// This allows users to have multiple roles like ["admin", "employee"] or ["hr", "it"].
// Only an Admin can call this.
//
// @Summary Add role to user
// @Description Add an additional role to a user (supports multiple roles per user). Admin only.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.AddRoleRequest true "Add role request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /api/v1/users/add-role [post]
func (h *UserHandler) AddRole(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)

	var req models.AddRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	roles, err := h.userService.AddRoleToUser(callerID, req.UserID, req.Role)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Role added successfully", gin.H{"user_id": req.UserID, "roles": roles})
}

// RemoveRole removes a role from a user.
// Cannot remove the last role - user must have at least one role.
// Only an Admin can call this.
//
// @Summary Remove role from user
// @Description Remove a role from a user. Cannot remove the last role. Admin only.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.RemoveRoleRequest true "Remove role request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /api/v1/users/remove-role [post]
func (h *UserHandler) RemoveRole(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)

	var req models.RemoveRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	roles, err := h.userService.RemoveRoleFromUser(callerID, req.UserID, req.Role)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Role removed successfully", gin.H{"user_id": req.UserID, "roles": roles})
}

// SetRoles replaces all roles for a user with the specified roles.
// At least one role must be provided.
// Only an Admin can call this.
//
// @Summary Set user roles
// @Description Replace all roles for a user with new role array. Admin only.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.SetRolesRequest true "Set roles request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /api/v1/users/set-roles [post]
func (h *UserHandler) SetRoles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)

	var req models.SetRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	roles, err := h.userService.SetUserRoles(callerID, req.UserID, req.Roles)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Roles updated successfully", gin.H{"user_id": req.UserID, "roles": roles})
}

// ChangeRoles applies a checkbox-style role change.
// Send the full set of desired (checked) roles. The backend computes the diff: what was added, removed, unchanged.
// Only Admin can call this.
//
// @Summary Change user roles (checkbox style)
// @Description Apply a checkbox-style role change. Send all checked roles; backend computes the diff. Admin only.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.ChangeRolesRequest true "Change roles request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /api/v1/users/change-roles [put]
func (h *UserHandler) ChangeRoles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)

	var req models.ChangeRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	result, err := h.userService.ChangeUserRoles(callerID, req.UserID, req.Roles)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Roles changed successfully", result)
}

// GetUser returns a single user's profile with their roles and account information.
// Does NOT include employee-related data — only user/account data.
//
// @Summary Get user detail
// @Description Get a user's profile, roles, and account info. HR or Admin only.
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	targetID, ok := parseUserIDParam(c)
	if !ok {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	detail, err := h.userService.GetUserDetail(targetID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "User retrieved successfully", detail)
}

// GetUserRoles returns all roles for a user (combined legacy + RBAC, deduplicated).
//
// @Summary Get user roles
// @Description Get all roles for a user
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/users/{id}/roles [get]
func (h *UserHandler) GetUserRoles(c *gin.Context) {
	targetID, ok := parseUserIDParam(c)
	if !ok {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	roles, err := h.userService.GetUserRoles(targetID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "User roles retrieved successfully", gin.H{"user_id": targetID, "roles": roles})
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

// ListBlockedUsers returns users blocked from login (rate-limit or manual). Admin only.
//
// @Summary List blocked users
// @Description Get users blocked from login (failed attempts or manually blocked). Admin only.
// @Tags Security
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20) maximum(100)
// @Param include_suspended query bool false "Include suspended users" default(false)
// @Success 200 {object} response.APIResponse
// @Router /api/v1/security/blocked-users [get]
func (h *UserHandler) ListBlockedUsers(c *gin.Context) {
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
	includeSuspended := false
	if c.Query("include_suspended") == "true" || c.Query("include_suspended") == "1" {
		includeSuspended = true
	}
	users, total, err := h.userService.ListBlockedUsers(tenantID, page, pageSize, includeSuspended)
	if err != nil {
		response.InternalServerError(c, "Failed to list blocked users", err.Error())
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{Page: page, PerPage: pageSize, Total: total, TotalPages: totalPages}
	list := make([]map[string]interface{}, 0, len(users))
	for _, u := range users {
		list = append(list, map[string]interface{}{
			"id":               u.ID,
			"email":            u.Email,
			"firstName":        u.FirstName,
			"lastName":         u.LastName,
			"username":         u.Username,
			"status":           u.Status,
			"failedLoginCount": u.FailedLoginCount,
			"lockedUntil":      u.LockedUntil,
			"updatedAt":        u.UpdatedAt,
		})
	}
	response.SuccessWithMeta(c, "Blocked users retrieved successfully", list, meta)
}

// GetRateLimitStatus returns login rate-limit status for an email. Admin only.
//
// @Summary Get rate limit status
// @Description Get login rate-limit status for a user by email (blocked, attempts remaining, etc.). Admin only.
// @Tags Security
// @Produce json
// @Param email query string true "User email"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/security/rate-limit/status [get]
func (h *UserHandler) GetRateLimitStatus(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.BadRequest(c, "email is required", nil)
		return
	}
	status, err := h.userService.GetRateLimitStatus(email)
	if err != nil {
		response.InternalServerError(c, "Failed to get rate limit status", err.Error())
		return
	}
	response.Success(c, "Rate limit status retrieved", status)
}

// ResetPassword resets a user's password to a specified value or the default.
// Also unblocks the account and clears failed login attempts.
//
// @Summary Reset user password
// @Description Reset a user's password. If no password is provided, defaults to "GreenTelecom@2026". Also unblocks account and clears failed login count.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /api/v1/users/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	callerID := userID.(uint)

	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	newPassword := ""
	if req.Password != nil {
		newPassword = *req.Password
	}

	result, err := h.userService.ResetUserPassword(callerID, req.UserID, newPassword)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Password reset successfully", result)
}

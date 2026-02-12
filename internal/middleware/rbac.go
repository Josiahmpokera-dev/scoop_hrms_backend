package middleware

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// PermissionMiddleware checks if the user has the required permission
func PermissionMiddleware(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		userObj, ok := user.(*models.User)
		if !ok {
			response.Unauthorized(c, "Invalid user context")
			c.Abort()
			return
		}

		// Check if user has permission through roles
		roleRepo := repositories.NewRoleRepository()
		userRoles, err := roleRepo.GetUserRoles(userObj.ID)
		if err != nil {
			response.Forbidden(c, "Failed to check permissions")
			c.Abort()
			return
		}

		hasPermission := false
		for _, role := range userRoles {
			for _, perm := range role.Permissions {
				if perm.Code == permissionCode {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		// Also check legacy role-based access (for backward compatibility)
		if !hasPermission {
			hasPermission = checkLegacyRolePermission(userObj.Role, permissionCode)
		}

		if !hasPermission {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkLegacyRolePermission checks permissions based on legacy role system.
// Admin gets all permissions; other roles get a predefined subset for backward compatibility.
func checkLegacyRolePermission(role models.UserRole, permissionCode string) bool {
	// Admin and Super Admin have full access to everything
	if role == models.RoleAdmin || role == models.RoleSuperAdmin {
		return true
	}

	// Map legacy roles to permissions for backward compatibility
	rolePermissions := map[models.UserRole][]string{
		models.RoleHR: {
			"dashboard:view", "dashboard:manage_announcements",
			"employee:read", "employee:create", "employee:update", "employee:delete",
			"employee:onboard", "employee:offboard",
			"leave:read", "leave:create", "leave:approve",
			"leave:manage_types", "leave:manage_policies", "leave:manage_holidays",
			"attendance:read", "attendance:approve",
			"shift:read", "shift:create", "shift:update", "shift:delete",
			"shift:manage_rosters", "shift:approve_swaps",
			"department:read", "department:create", "department:update", "department:delete",
			"team:read", "team:create", "team:update", "team:delete",
			"position:read", "position:create", "position:update", "position:delete",
			"organization:read", "organization:update",
			"organization_unit:read", "organization_unit:create", "organization_unit:update", "organization_unit:delete",
			"location:read", "location:create", "location:update", "location:delete",
			"cost_center:read", "cost_center:create", "cost_center:update", "cost_center:delete",
			"payroll:read", "payroll:run", "payroll:manage_structures",
			"payroll:manage_loans", "payroll:manage_reports",
			"asset:read", "asset:create", "asset:update", "asset:assign",
			"helpdesk:read", "helpdesk:create", "helpdesk:manage", "helpdesk:manage_kb",
			"user:read",
			"report:view", "report:export",
		},
		models.RoleIT: {
			"dashboard:view",
			"employee:read",
			"asset:read", "asset:create", "asset:update", "asset:assign",
			"helpdesk:read", "helpdesk:create", "helpdesk:manage", "helpdesk:manage_kb",
			"user:read",
			"attendance:read", "attendance:manage_config",
		},
		models.RoleManager: {
			"dashboard:view",
			"employee:read", "employee:update",
			"leave:read", "leave:create", "leave:approve",
			"attendance:read", "attendance:approve",
			"shift:read", "shift:approve_swaps",
			"department:read", "team:read", "position:read",
			"organization:read", "location:read",
			"helpdesk:read", "helpdesk:create",
			"report:view",
			"payroll:read", "asset:read",
		},
		models.RoleEmployee: {
			"dashboard:view",
			"employee:read",
			"leave:read", "leave:create",
			"shift:read",
			"department:read", "team:read", "position:read",
			"organization:read", "location:read",
			"helpdesk:read", "helpdesk:create",
			"payroll:read", "asset:read", "attendance:read",
		},
		models.RoleUser: {
			"dashboard:view",
			"employee:read",
			"leave:read", "leave:create",
			"shift:read",
			"department:read", "team:read", "position:read",
			"organization:read", "location:read",
			"helpdesk:read", "helpdesk:create",
			"payroll:read", "asset:read", "attendance:read",
		},
	}

	permissions, exists := rolePermissions[role]
	if !exists {
		return false
	}

	for _, perm := range permissions {
		if perm == permissionCode {
			return true
		}
	}

	return false
}

// RequireAnyPermission checks if user has any of the specified permissions
func RequireAnyPermission(permissionCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		userObj, ok := user.(*models.User)
		if !ok {
			response.Unauthorized(c, "Invalid user context")
			c.Abort()
			return
		}

		roleRepo := repositories.NewRoleRepository()
		userRoles, err := roleRepo.GetUserRoles(userObj.ID)
		if err != nil {
			response.Forbidden(c, "Failed to check permissions")
			c.Abort()
			return
		}

		hasPermission := false
		for _, role := range userRoles {
			for _, perm := range role.Permissions {
				for _, requiredPerm := range permissionCodes {
					if perm.Code == requiredPerm {
						hasPermission = true
						break
					}
				}
				if hasPermission {
					break
				}
			}
			if hasPermission {
				break
			}
		}

		// Check legacy role permissions
		if !hasPermission {
			for _, requiredPerm := range permissionCodes {
				if checkLegacyRolePermission(userObj.Role, requiredPerm) {
					hasPermission = true
					break
				}
			}
		}

		if !hasPermission {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

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

// checkLegacyRolePermission checks permissions based on legacy role system
func checkLegacyRolePermission(role models.UserRole, permissionCode string) bool {
	// Map legacy roles to permissions for backward compatibility
	rolePermissions := map[models.UserRole][]string{
		models.RoleAdmin: {
			"employee:read", "employee:create", "employee:update", "employee:delete",
			"department:read", "department:create", "department:update", "department:delete",
			"payroll:run", "attendance:approve",
			"user:read", "user:create", "user:update", "user:delete",
		},
		models.RoleHR: {
			"employee:read", "employee:create", "employee:update",
			"department:read", "payroll:run", "attendance:approve",
			"user:read",
		},
		models.RoleUser: {
			"employee:read",
		},
		models.RoleIT: {
			"employee:read",
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

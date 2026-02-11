package middleware

import (
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/auth/services"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	roleRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// buildUserRolesSlice returns a deduplicated slice of role strings for middleware (legacy user_type + assigned role codes).
// Always includes "user" as a base role since all system users are treated as users.
func buildUserRolesSlice(legacyRole string, roleCodes []string) []string {
	seen := make(map[string]bool)
	var roles []string

	// Add legacy role first
	legacy := strings.ToLower(strings.TrimSpace(legacyRole))
	if legacy != "" && !seen[legacy] {
		seen[legacy] = true
		roles = append(roles, legacy)
	}

	// Add RBAC roles from the database
	for _, code := range roleCodes {
		c := strings.ToLower(strings.TrimSpace(code))
		if c != "" && !seen[c] {
			seen[c] = true
			roles = append(roles, c)
		}
	}

	// Always ensure "user" is in the roles array - all system users are treated as users
	if !seen["user"] {
		roles = append(roles, "user")
	}

	return roles
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		userID, role, err := authService.ValidateToken(token)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Load user to get tenant_id and build roles from DB
		userRepo := userRepos.NewUserRepository()
		roleRepo := roleRepos.NewRoleRepository()
		user, err := userRepo.FindByID(userID)
		if err == nil && user != nil {
			// Set full user object in context
			c.Set("user", user)
			// Build roles array: legacy user_type + assigned role codes
			codes, _ := roleRepo.GetUserRoleCodes(user.ID)
			userRoles := buildUserRolesSlice(string(user.Role), codes)
			c.Set("user_roles", userRoles)
			// Primary role for backward compatibility
			if len(userRoles) > 0 {
				c.Set("user_role", userRoles[0])
			} else {
				c.Set("user_role", string(user.Role))
			}
		} else {
			// Token valid but user not found: still set user_id and user_role from token
			c.Set("user_role", role)
			c.Set("user_roles", []string{role})
		}

		// Set user information in context
		c.Set("user_id", userID)

		c.Next()
	}
}

// containsRole returns true if the user has the given role (checks both user_roles array and legacy user_role).
func containsRole(c *gin.Context, want string) bool {
	if roles, exists := c.Get("user_roles"); exists {
		if arr, ok := roles.([]string); ok {
			for _, r := range arr {
				if r == want {
					return true
				}
			}
		}
	}
	if role, exists := c.Get("user_role"); exists {
		if r, ok := role.(string); ok && r == want {
			return true
		}
	}
	return false
}

// AdminMiddleware ensures the user is an admin (user has "admin" in roles array or as primary role)
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user_id"); !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}
		if !containsRole(c, "admin") {
			response.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// HRMiddleware ensures the user is HR or Admin (user has "admin" or "hr" in roles array or as primary role)
func HRMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user_id"); !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}
		if !containsRole(c, "admin") && !containsRole(c, "hr") {
			response.Forbidden(c, "HR or Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

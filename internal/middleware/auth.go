package middleware

import (
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/auth/services"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

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

		// Load user to get tenant_id
		userRepo := userRepos.NewUserRepository()
		user, err := userRepo.FindByID(userID)
		if err == nil && user != nil {
			// Set full user object in context
			c.Set("user", user)
			// Set tenant_id in context for easy access
			if user.TenantID != nil {
				c.Set(TenantIDKey, *user.TenantID)
			}
		}

		// Set user information in context
		c.Set("user_id", userID)
		c.Set("user_role", role)

		c.Next()
	}
}

// AdminMiddleware ensures the user is an admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		if role != "admin" {
			response.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// HRMiddleware ensures the user is HR or Admin
func HRMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		if role != "admin" && role != "hr" {
			response.Forbidden(c, "HR or Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

package middleware

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/tenants/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

const TenantIDKey = "tenant_id"

// TenantMiddleware extracts tenant ID from header or context
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get tenant ID from header
		tenantIDStr := c.GetHeader("X-Tenant-ID")
		
		// If not in header, try to get from authenticated user
		if tenantIDStr == "" {
			user, exists := c.Get("user")
			if exists {
				// Assuming user has TenantID field
				// This will be set by AuthMiddleware
				if userObj, ok := user.(interface{ GetTenantID() *uint }); ok {
					if tenantID := userObj.GetTenantID(); tenantID != nil {
						c.Set(TenantIDKey, *tenantID)
						c.Next()
						return
					}
				}
			}
		} else {
			// Parse tenant ID from header
			tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
			if err == nil {
				c.Set(TenantIDKey, uint(tenantID))
				c.Next()
				return
			}
		}

		// For now, allow requests without tenant (backward compatibility)
		// In production, you might want to require tenant ID
		c.Next()
	}
}

// RequireTenantMiddleware requires tenant ID to be present
func RequireTenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, exists := c.Get(TenantIDKey)
		if !exists {
			response.Unauthorized(c, "Tenant ID is required")
			c.Abort()
			return
		}

		// Verify tenant exists and is active
		tenantRepo := repositories.NewTenantRepository()
		tenant, err := tenantRepo.FindByID(tenantID.(uint))
		if err != nil || !tenant.IsActive() {
			response.Forbidden(c, "Invalid or inactive tenant")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetTenantID retrieves tenant ID from context
func GetTenantID(c *gin.Context) *uint {
	tenantID, exists := c.Get(TenantIDKey)
	if !exists {
		return nil
	}
	
	if id, ok := tenantID.(uint); ok {
		return &id
	}
	return nil
}

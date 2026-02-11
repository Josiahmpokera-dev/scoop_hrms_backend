package middleware

import "github.com/gin-gonic/gin"

// TenantIDKey is kept for backward compatibility (single-tenant mode).
const TenantIDKey = "tenant_id"

// GetTenantID always returns nil in single-tenant mode.
func GetTenantID(c *gin.Context) *uint {
	return nil
}

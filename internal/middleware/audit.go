package middleware

import (
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/repositories"
	"github.com/gin-gonic/gin"
)

const userIDKey = "user_id"

// deriveResourceAndAction extracts resource and action from path for audit.
// e.g. /api/v1/users/assign-role -> resource=users, action=assign-role
// e.g. /api/v1/employees -> resource=employees, action=list
func deriveResourceAndAction(path, method string) (resource, action string) {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	// Expect /api/v1/<resource>/... or /api/v1/<resource>
	if len(parts) >= 3 {
		resource = parts[2]
		if len(parts) > 3 {
			action = parts[3]
		} else {
			action = method
		}
	} else {
		resource = "api"
		action = method
	}
	return resource, action
}

// AuditMiddleware logs every request to the audit_logs table (who, what, when, where, outcome).
// Should be registered on routes that need auditing (e.g. v1 API group). Runs after handler to capture status code.
func AuditMiddleware() gin.HandlerFunc {
	repo := repositories.NewAuditRepository()

	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()
		var userID, tenantID *uint
		if uid, exists := c.Get(userIDKey); exists {
			if u, ok := uid.(uint); ok {
				userID = &u
			}
		}
		if tid, exists := c.Get(TenantIDKey); exists {
			if t, ok := tid.(uint); ok {
				tenantID = &t
			}
		}

		c.Next()

		statusCode := c.Writer.Status()
		resource, action := deriveResourceAndAction(path, method)

		entry := &models.AuditLog{
			TenantID:   tenantID,
			UserID:     userID,
			Action:     action,
			Resource:   resource,
			Method:     method,
			Path:       path,
			StatusCode: statusCode,
			IP:         ip,
			UserAgent:  userAgent,
		}
		_ = repo.Create(entry)
	}
}

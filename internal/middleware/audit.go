package middleware

import (
	"log"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/repositories"
	"github.com/gin-gonic/gin"
)

const userIDKey = "user_id"

// getClientIPForAudit returns the client IP for audit logging.
// Prefers X-Forwarded-For (first IP), X-Real-IP, CF-Connecting-IP when present (e.g. behind reverse proxy),
// then falls back to c.ClientIP(). When client and server are on the same machine (localhost), result is ::1 or 127.0.0.1.
func getClientIPForAudit(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// First IP is the client, rest are proxies
		if idx := strings.Index(xff, ","); idx > 0 {
			xff = strings.TrimSpace(xff[:idx])
		} else {
			xff = strings.TrimSpace(xff)
		}
		if xff != "" {
			return xff
		}
	}
	if xri := strings.TrimSpace(c.GetHeader("X-Real-IP")); xri != "" {
		return xri
	}
	if cf := strings.TrimSpace(c.GetHeader("CF-Connecting-IP")); cf != "" {
		return cf
	}
	return c.ClientIP()
}

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
		ip := getClientIPForAudit(c)
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
		if err := repo.Create(entry); err != nil {
			log.Printf("[audit] failed to write audit log: %v (path=%s method=%s)", err, path, method)
		}
	}
}

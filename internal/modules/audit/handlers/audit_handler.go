package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// AuditHandler handles audit log API requests.
type AuditHandler struct {
	auditService *services.AuditService
}

// NewAuditHandler returns a new audit handler.
func NewAuditHandler() *AuditHandler {
	return &AuditHandler{auditService: services.NewAuditService()}
}

// ListAuditLogs returns paginated audit logs with optional filters and search.
//
// @Summary List audit logs
// @Description Get paginated audit logs (Security & Audit). Admin only. Supports pagination, search, and filters.
// @Tags Security & Audit
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20) maximum(100)
// @Param user_id query int false "Filter by user ID (who performed the action)"
// @Param action query string false "Filter by action (partial match)"
// @Param resource query string false "Filter by resource (partial match)"
// @Param method query string false "Filter by HTTP method (GET, POST, PUT, PATCH, DELETE)"
// @Param date_from query string false "Filter from date (RFC3339 or 2006-01-02)"
// @Param date_to query string false "Filter to date (RFC3339 or 2006-01-02)"
// @Param status_code query int false "Filter by HTTP response status code"
// @Param search query string false "Search in action, resource, path, details, reason"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/security/audit [get]
func (h *AuditHandler) ListAuditLogs(c *gin.Context) {
	// Do not filter by tenant so admins see all audit logs (including unauthenticated requests with tenant_id = nil)
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

	filter := services.ListFilter{
		TenantID: nil, // Default: show all logs; optional ?tenant_id= for scoping
		Action:   strings.TrimSpace(c.Query("action")),
		Resource: strings.TrimSpace(c.Query("resource")),
		Method:   strings.TrimSpace(c.Query("method")),
		Search:   strings.TrimSpace(c.Query("search")),
	}

	if tid := c.Query("tenant_id"); tid != "" {
		if parsed, err := strconv.ParseUint(tid, 10, 32); err == nil {
			tidUint := uint(parsed)
			filter.TenantID = &tidUint
		}
	}
	if u := c.Query("user_id"); u != "" {
		if parsed, err := strconv.ParseUint(u, 10, 32); err == nil {
			uid := uint(parsed)
			filter.UserID = &uid
		}
	}
	if s := c.Query("status_code"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil {
			filter.StatusCode = &parsed
		}
	}
	if df := c.Query("date_from"); df != "" {
		for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05Z07:00", "2006-01-02"} {
			if t, err := time.Parse(layout, df); err == nil {
				filter.DateFrom = &t
				break
			}
		}
	}
	if dt := c.Query("date_to"); dt != "" {
		for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05Z07:00", "2006-01-02"} {
			if t, err := time.Parse(layout, dt); err == nil {
				filter.DateTo = &t
				break
			}
		}
	}

	logs, total, err := h.auditService.List(filter, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list audit logs", err.Error())
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Audit logs retrieved successfully", logs, meta)
}

package services

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/repositories"
)

// AuditService handles audit logging and listing.
type AuditService struct {
	repo *repositories.AuditRepository
}

// NewAuditService returns a new audit service.
func NewAuditService() *AuditService {
	return &AuditService{repo: repositories.NewAuditRepository()}
}

// Log creates an audit log entry. Safe to call from middleware or handlers.
func (s *AuditService) Log(entry *models.AuditLog) error {
	return s.repo.Create(entry)
}

// ListFilter is the filter for listing audit logs (exposed to handlers).
type ListFilter struct {
	UserID     *uint
	Action     string
	Resource   string
	Method     string
	DateFrom   *time.Time
	DateTo     *time.Time
	Search     string
	StatusCode *int
}

// List returns audit logs with pagination and filters.
func (s *AuditService) List(filter ListFilter, page, pageSize int) ([]models.AuditLog, int64, error) {
	rf := repositories.ListFilter{
		UserID:     filter.UserID,
		Action:     filter.Action,
		Resource:   filter.Resource,
		Method:     filter.Method,
		DateFrom:   filter.DateFrom,
		DateTo:     filter.DateTo,
		Search:     filter.Search,
		StatusCode: filter.StatusCode,
	}
	return s.repo.List(rf, page, pageSize)
}

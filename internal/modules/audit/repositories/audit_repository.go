package repositories

import (
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/models"
	"gorm.io/gorm"
)

// AuditRepository handles audit log persistence.
type AuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository returns a new audit repository.
func NewAuditRepository() *AuditRepository {
	return &AuditRepository{db: database.GetDB()}
}

// Create inserts an audit log entry.
func (r *AuditRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

// ListFilter holds filters for listing audit logs.
type ListFilter struct {
	UserID     *uint
	Action     string   // exact or partial
	Resource   string   // exact or partial
	Method     string   // GET, POST, etc.
	DateFrom   *time.Time
	DateTo     *time.Time
	Search     string   // search in action, resource, path, details
	StatusCode *int     // filter by response status
}

// List returns audit logs with pagination and filters.
func (r *AuditRepository) List(filter ListFilter, page, pageSize int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := r.db.Model(&models.AuditLog{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Action != "" {
		query = query.Where("action ILIKE ?", "%"+strings.TrimSpace(filter.Action)+"%")
	}
	if filter.Resource != "" {
		query = query.Where("resource ILIKE ?", "%"+strings.TrimSpace(filter.Resource)+"%")
	}
	if filter.Method != "" {
		query = query.Where("method = ?", strings.ToUpper(strings.TrimSpace(filter.Method)))
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", filter.DateTo)
	}
	if filter.StatusCode != nil {
		query = query.Where("status_code = ?", *filter.StatusCode)
	}
	if filter.Search != "" {
		s := "%" + strings.TrimSpace(filter.Search) + "%"
		query = query.Where(
			"action ILIKE ? OR resource ILIKE ? OR path ILIKE ? OR details ILIKE ? OR reason ILIKE ?",
			s, s, s, s, s,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog represents a system audit log entry for security and compliance.
// Captures who did what, when, where, and outcome for every action.
type AuditLog struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	TenantID   *uint          `json:"tenant_id,omitempty" gorm:"index"`
	UserID     *uint          `json:"user_id,omitempty" gorm:"index"`           // Who performed the action (nil if unauthenticated)
	Action     string         `json:"action" gorm:"size:100;index"`             // e.g. assign_role, create_employee, login
	Resource   string         `json:"resource" gorm:"size:80;index"`           // e.g. users, employees, leave
	ResourceID *string        `json:"resource_id,omitempty" gorm:"size:50"`    // Target entity ID if applicable
	Method     string         `json:"method" gorm:"size:10;index"`              // HTTP method: GET, POST, PUT, PATCH, DELETE
	Path       string         `json:"path" gorm:"size:512;index"`               // Request path
	StatusCode int            `json:"status_code" gorm:"index"`                 // HTTP response status
	IP         string         `json:"ip" gorm:"size:45"`                        // Client IP
	UserAgent  string         `json:"user_agent" gorm:"size:512"`              // Client user agent
	Details    string         `json:"details,omitempty" gorm:"type:text"`       // Optional JSON or text details
	Reason     string         `json:"reason,omitempty" gorm:"size:255"`        // Optional reason/notes
	CreatedAt  time.Time      `json:"created_at" gorm:"index"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for AuditLog.
func (AuditLog) TableName() string {
	return "audit_logs"
}

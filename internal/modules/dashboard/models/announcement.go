package models

import (
	"time"

	"gorm.io/gorm"
)

// AnnouncementType represents the type of announcement
type AnnouncementType string

const (
	AnnouncementTypeInfo    AnnouncementType = "info"
	AnnouncementTypeWarning AnnouncementType = "warning"
	AnnouncementTypeSuccess AnnouncementType = "success"
	AnnouncementTypeUrgent  AnnouncementType = "urgent"
)

// AnnouncementPriority represents the priority of announcement
type AnnouncementPriority string

const (
	AnnouncementPriorityLow      AnnouncementPriority = "low"
	AnnouncementPriorityNormal   AnnouncementPriority = "normal"
	AnnouncementPriorityHigh     AnnouncementPriority = "high"
	AnnouncementPriorityCritical AnnouncementPriority = "critical"
)

// Announcement represents a company announcement
type Announcement struct {
	ID          uint                 `json:"id" gorm:"primaryKey"`
	TenantID    *uint                `json:"tenant_id,omitempty" gorm:"index"`
	Title       string               `json:"title" gorm:"not null;size:255"`
	Description string               `json:"description" gorm:"type:text;not null"`
	Type        AnnouncementType     `json:"type" gorm:"type:varchar(20);default:'info'"`
	Priority    AnnouncementPriority `json:"priority" gorm:"type:varchar(20);default:'normal'"`
	IsActive    bool                 `json:"is_active" gorm:"default:true"`
	ExpiresAt   *time.Time           `json:"expires_at,omitempty"`
	CreatedBy   uint                 `json:"created_by" gorm:"not null;index"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	DeletedAt   gorm.DeletedAt       `json:"-" gorm:"index"`

	// Relationships
	Attachments []AnnouncementAttachment `json:"attachments,omitempty" gorm:"foreignKey:AnnouncementID"`
	ReadRecords []AnnouncementRead       `json:"read_records,omitempty" gorm:"foreignKey:AnnouncementID"`
}

// TableName specifies the table name
func (Announcement) TableName() string {
	return "announcements"
}

// AnnouncementAttachment represents an attachment to an announcement
type AnnouncementAttachment struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	AnnouncementID uint           `json:"announcement_id" gorm:"not null;index"`
	Name           string         `json:"name" gorm:"not null;size:255"`
	URL            string         `json:"url" gorm:"not null;size:500"`
	Size           string         `json:"size" gorm:"size:50"`
	MimeType       *string        `json:"mime_type,omitempty" gorm:"size:100"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (AnnouncementAttachment) TableName() string {
	return "announcement_attachments"
}

// AnnouncementRead tracks which users have read which announcements
type AnnouncementRead struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	AnnouncementID uint      `json:"announcement_id" gorm:"not null;index;uniqueIndex:idx_announcement_user"`
	UserID         uint      `json:"user_id" gorm:"not null;index;uniqueIndex:idx_announcement_user"`
	ReadAt         time.Time `json:"read_at" gorm:"not null"`
}

// TableName specifies the table name
func (AnnouncementRead) TableName() string {
	return "announcement_reads"
}

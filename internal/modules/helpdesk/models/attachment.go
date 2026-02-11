package models

import (
	"time"

	"gorm.io/gorm"
)

// Attachment represents an attachment on a helpdesk ticket
type Attachment struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TicketID    uint           `json:"ticket_id" gorm:"index;not null"`
	Name        string         `json:"name" gorm:"not null;size:255"`
	URL         string         `json:"url" gorm:"not null;size:500"`
	Size        int64          `json:"size" gorm:"not null"` // Size in bytes
	Type        string          `json:"type" gorm:"size:100"` // MIME type
	UploadedBy  uint           `json:"uploaded_by" gorm:"index;not null"` // User ID
	UploadedAt  time.Time      `json:"uploaded_at" gorm:"not null"`
	
	// Timestamps
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (Attachment) TableName() string {
	return "helpdesk_attachments"
}

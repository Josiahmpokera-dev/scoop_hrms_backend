package models

import (
	"time"

	"gorm.io/gorm"
)

// CommentAuthorType represents the type of comment author
type CommentAuthorType string

const (
	CommentAuthorTypeRequester CommentAuthorType = "Requester"
	CommentAuthorTypeAgent     CommentAuthorType = "Agent"
	CommentAuthorTypeAdmin     CommentAuthorType = "Admin"
)

// Comment represents a comment on a helpdesk ticket
type Comment struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	TenantID    *uint            `json:"tenant_id,omitempty" gorm:"index"`
	TicketID    uint             `json:"ticket_id" gorm:"index;not null"`
	AuthorID    uint             `json:"author_id" gorm:"index;not null"` // User ID
	AuthorName  string           `json:"author_name" gorm:"not null;size:255"`
	AuthorType  CommentAuthorType `json:"author_type" gorm:"type:varchar(50);not null"`
	Text        string           `json:"text" gorm:"type:text;not null"`
	IsPrivate   bool             `json:"is_private" gorm:"default:false"` // Private comments visible only to agents/admin
	Timestamp   time.Time        `json:"timestamp" gorm:"not null"`
	
	// Timestamps
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (Comment) TableName() string {
	return "helpdesk_comments"
}

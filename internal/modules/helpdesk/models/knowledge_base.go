package models

import (
	"time"

	"gorm.io/gorm"
)

// KBArticleStatus represents the status of a knowledge base article
type KBArticleStatus string

const (
	KBArticleStatusPublished KBArticleStatus = "Published"
	KBArticleStatusDraft     KBArticleStatus = "Draft"
	KBArticleStatusArchived  KBArticleStatus = "Archived"
)

// KnowledgeBaseArticle represents a knowledge base article
type KnowledgeBaseArticle struct {
	ID              uint             `json:"id" gorm:"primaryKey"`
	TenantID        *uint            `json:"tenant_id,omitempty" gorm:"index"`
	ArticleNumber   string           `json:"article_number" gorm:"uniqueIndex;not null;size:50"` // e.g., KB-001
	Title           string           `json:"title" gorm:"not null;size:255"`
	Content         string           `json:"content" gorm:"type:text;not null"`
	ContentPreview  *string          `json:"content_preview,omitempty" gorm:"type:text"` // Preview for list views
	Category        string           `json:"category" gorm:"not null;size:100"` // IT, Payroll, HR, Leave, etc.
	
	// Tags (stored as JSON array)
	TagsJSON        *string          `json:"-" gorm:"type:text"`
	
	// Statistics
	Views           int              `json:"views" gorm:"default:0"`
	Helpful         int              `json:"helpful" gorm:"default:0"`
	NotHelpful      int              `json:"not_helpful" gorm:"default:0"`
	
	// Metadata
	Author          string           `json:"author" gorm:"not null;size:255"` // Author name or team
	AuthorID        *uint            `json:"author_id,omitempty" gorm:"index"` // User ID
	Status          KBArticleStatus  `json:"status" gorm:"type:varchar(50);default:'Draft'"`
	Featured        bool             `json:"featured" gorm:"default:false"`
	
	// Timestamps
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	UpdatedBy       *uint            `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt       gorm.DeletedAt   `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (KnowledgeBaseArticle) TableName() string {
	return "helpdesk_kb_articles"
}

// KBArticleFeedback represents feedback on a knowledge base article
type KBArticleFeedback struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TenantID    *uint          `json:"tenant_id,omitempty" gorm:"index"`
	ArticleID   uint           `json:"article_id" gorm:"index;not null"`
	UserID      uint           `json:"user_id" gorm:"index;not null"`
	IsHelpful   bool           `json:"is_helpful" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (KBArticleFeedback) TableName() string {
	return "helpdesk_kb_feedback"
}

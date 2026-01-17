package models

import (
	"time"

	"gorm.io/gorm"
)

// Permission represents a permission in the RBAC system
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null;size:150"`
	Name        string         `json:"name" gorm:"not null;size:255"`
	Resource    string         `json:"resource" gorm:"not null;size:100"`
	Action      string         `json:"action" gorm:"not null;size:50"`
	Description *string        `json:"description,omitempty" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UpdatedBy   *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for Permission model
func (Permission) TableName() string {
	return "permissions"
}

// GetFullCode returns the full permission code (resource:action)
func (p *Permission) GetFullCode() string {
	return p.Resource + ":" + p.Action
}

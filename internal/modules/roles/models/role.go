package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a role in the RBAC system
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TenantID    *uint          `json:"tenant_id,omitempty" gorm:"index"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null;size:100"`
	Name        string         `json:"name" gorm:"not null;size:150"`
	Description *string        `json:"description,omitempty" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UpdatedBy   *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

// TableName specifies the table name for Role model
func (Role) TableName() string {
	return "roles"
}

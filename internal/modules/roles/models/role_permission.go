package models

import (
	"time"

	"gorm.io/gorm"
)

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	RoleID       uint           `json:"role_id" gorm:"primaryKey"`
	PermissionID uint           `json:"permission_id" gorm:"primaryKey"`
	CreatedAt    time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `json:"updated_at"`
	UpdatedBy    *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for RolePermission model
func (RolePermission) TableName() string {
	return "role_permissions"
}

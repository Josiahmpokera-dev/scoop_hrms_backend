package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID     uint           `json:"user_id" gorm:"primaryKey"`
	RoleID     uint           `json:"role_id" gorm:"primaryKey"`
	AssignedBy *uint          `json:"assigned_by,omitempty" gorm:"index"`
	AssignedAt time.Time      `json:"assigned_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `json:"updated_at"`
	UpdatedBy  *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for UserRole model
func (UserRole) TableName() string {
	return "user_roles"
}

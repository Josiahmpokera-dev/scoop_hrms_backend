package models

import (
	"time"

	"gorm.io/gorm"
)

// LeaveType represents a type of leave (Annual, Sick, Emergency, etc.)
type LeaveType struct {
	ID                     uint           `json:"id" gorm:"primaryKey"`
	TenantID               *uint          `json:"tenant_id,omitempty" gorm:"index"`
	Code                   string         `json:"code" gorm:"uniqueIndex;not null;size:10"` // AL, SL, EL, etc.
	Name                   string         `json:"name" gorm:"not null;size:100"`
	Icon                   *string        `json:"icon,omitempty" gorm:"size:10"` // Emoji icon
	Category               string         `json:"category" gorm:"not null;size:50"` // Annual, Emergency, Sick, Other
	PaidLeave              bool           `json:"paid_leave" gorm:"default:true"`
	RequiresDocumentation   bool           `json:"requires_documentation" gorm:"default:false"`
	IsActive                bool           `json:"is_active" gorm:"default:true"`
	Description             *string        `json:"description,omitempty" gorm:"type:text"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	CreatedBy               *uint          `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy               *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt               gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (LeaveType) TableName() string {
	return "leave_types"
}

// CreateLeaveTypeRequest represents the request to create a leave type
type CreateLeaveTypeRequest struct {
	Code                 string  `json:"code" binding:"required"`
	Name                 string  `json:"name" binding:"required"`
	Icon                 *string `json:"icon,omitempty"`
	Category             string  `json:"category" binding:"required"`
	PaidLeave            *bool   `json:"paid_leave"`
	RequiresDocumentation *bool   `json:"requires_documentation"`
	IsActive             *bool   `json:"is_active"`
	Description          *string `json:"description"`
}

// UpdateLeaveTypeRequest represents the request to update a leave type
type UpdateLeaveTypeRequest struct {
	Name                 *string `json:"name"`
	Icon                 *string `json:"icon"`
	Category             *string `json:"category"`
	PaidLeave            *bool   `json:"paid_leave"`
	RequiresDocumentation *bool   `json:"requires_documentation"`
	IsActive             *bool   `json:"is_active"`
	Description          *string `json:"description"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// TenantStatus represents tenant status
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusInactive  TenantStatus = "inactive"
)

// Tenant represents a tenant in the multi-tenant system
type Tenant struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	Name           string       `json:"name" gorm:"not null;size:255"`
	Domain         *string      `json:"domain,omitempty" gorm:"uniqueIndex;size:255"`
	Status         TenantStatus `json:"status" gorm:"type:varchar(20);default:'active';not null"`
	SubscriptionPlan *string    `json:"subscription_plan,omitempty" gorm:"size:100"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	UpdatedBy      *uint        `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for Tenant model
func (Tenant) TableName() string {
	return "tenants"
}

// IsActive checks if tenant is active
func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive
}

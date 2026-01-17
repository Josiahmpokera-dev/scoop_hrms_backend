package models

import (
	"time"

	"gorm.io/gorm"
)

// AddressType represents the type of address
type AddressType string

const (
	AddressTypeCurrent    AddressType = "current"
	AddressTypePermanent  AddressType = "permanent"
)

// EmployeeAddress represents an employee's address (current or permanent)
type EmployeeAddress struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	EmployeeID  *uint          `json:"employee_id,omitempty" gorm:"index"` // Nullable for drafts
	EmployeeIDString *string    `json:"employee_id_string,omitempty" gorm:"size:50;index"` // Employee ID string (e.g., "EMP466088") - for tracking during draft
	DraftID     *uint          `json:"draft_id,omitempty" gorm:"index"`   // References employee_onboarding_drafts(id)
	AddressType AddressType    `json:"address_type" gorm:"type:varchar(20);not null"` // current or permanent
	AddressLine1 *string       `json:"address_line1,omitempty" gorm:"type:text"`
	AddressLine2 *string       `json:"address_line2,omitempty" gorm:"type:text"`
	City        *string        `json:"city,omitempty" gorm:"size:100"`
	State       *string        `json:"state,omitempty" gorm:"size:100"`
	PostalCode  *string        `json:"postal_code,omitempty" gorm:"size:20"`
	Country     *string        `json:"country,omitempty" gorm:"size:100"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (EmployeeAddress) TableName() string {
	return "employee_addresses"
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeeEmergencyContact represents an emergency contact for an employee
type EmployeeEmergencyContact struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	EmployeeID     *uint          `json:"employee_id,omitempty" gorm:"index"` // Nullable for drafts
	EmployeeIDString *string       `json:"employee_id_string,omitempty" gorm:"size:50;index"` // Employee ID string (e.g., "EMP466088") - for tracking during draft
	DraftID        *uint          `json:"draft_id,omitempty" gorm:"index"`    // References employee_onboarding_drafts(id)
	ContactName    string         `json:"contact_name" gorm:"not null;size:100"`
	Relationship   *string        `json:"relationship,omitempty" gorm:"size:50"` // spouse, parent, sibling, friend, etc.
	PhoneNumber    string         `json:"phone_number" gorm:"not null;size:50"`
	AlternatePhone *string        `json:"alternate_phone,omitempty" gorm:"size:50"`
	Email          *string        `json:"email,omitempty" gorm:"size:191"`
	Address        *string        `json:"address,omitempty" gorm:"type:text"`
	IsPrimary      bool           `json:"is_primary" gorm:"default:false"` // Primary emergency contact
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (EmployeeEmergencyContact) TableName() string {
	return "employee_emergency_contacts"
}

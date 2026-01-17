package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeeAsset represents company assets assigned to an employee during onboarding
type EmployeeAsset struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	EmployeeID *uint          `json:"employee_id,omitempty" gorm:"index"` // References employees(id) - NULL during draft
	EmployeeIDString *string  `json:"employee_id_string,omitempty" gorm:"size:50;index"` // Employee ID string (e.g., "EMP466088") - for tracking during draft
	DraftID   *uint          `json:"draft_id,omitempty" gorm:"index"`     // References employee_onboarding_drafts(id) - NULL after completion
	
	AssetType        string     `json:"asset_type" gorm:"not null;size:100"` // laptop, phone, vehicle, etc.
	AssetName        string     `json:"asset_name" gorm:"not null;size:200"` // e.g., "MacBook Pro 16"
	SerialNumber     *string    `json:"serial_number,omitempty" gorm:"size:100"`
	AssetTag         *string    `json:"asset_tag,omitempty" gorm:"size:50"`
	AssignedDate     *time.Time `json:"assigned_date,omitempty"`
	ExpectedReturnDate *time.Time `json:"expected_return_date,omitempty"`
	Condition        *string    `json:"condition,omitempty" gorm:"size:50"` // new, good, fair, poor
	Notes            *string    `json:"notes,omitempty" gorm:"type:text"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (EmployeeAsset) TableName() string {
	return "employee_assets"
}

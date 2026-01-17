package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeeBasicInformation represents personal/basic information for an employee during onboarding
type EmployeeBasicInformation struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	EmployeeID     *uint           `json:"employee_id,omitempty" gorm:"index"` // References employees(id) - NULL during draft
	EmployeeIDString *string       `json:"employee_id_string,omitempty" gorm:"size:50;index"` // Employee ID string (e.g., "EMP466088") - for tracking during draft
	DraftID       *uint           `json:"draft_id,omitempty" gorm:"index"`     // References employee_onboarding_drafts(id) - NULL after completion
	
	// Profile Picture
	PhotoURL *string `json:"photo_url,omitempty" gorm:"size:500"`
	
	// Basic Information
	FirstName    string     `json:"first_name" gorm:"not null;size:120"`
	MiddleName   *string    `json:"middle_name,omitempty" gorm:"size:120"`
	LastName     string     `json:"last_name" gorm:"not null;size:120"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty"`
	Gender       *string    `json:"gender,omitempty" gorm:"size:20"` // male, female, other
	MaritalStatus *string   `json:"marital_status,omitempty" gorm:"size:20"` // single, married, divorced, widowed
	BloodGroup   *string    `json:"blood_group,omitempty" gorm:"size:10"` // A+, A-, B+, B-, AB+, AB-, O+, O-
	Nationality  *string    `json:"nationality,omitempty" gorm:"size:100"`
	
	// Contact Information
	PersonalEmail    *string `json:"personal_email,omitempty" gorm:"size:191"`
	MobileNumber     *string `json:"mobile_number,omitempty" gorm:"size:50"`
	AlternateNumber  *string `json:"alternate_number,omitempty" gorm:"size:50"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (EmployeeBasicInformation) TableName() string {
	return "employee_basic_information"
}

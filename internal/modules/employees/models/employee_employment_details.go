package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeeEmploymentDetails represents employment details for an employee during onboarding
type EmployeeEmploymentDetails struct {
	ID         uint  `json:"id" gorm:"primaryKey"`
	EmployeeID *uint `json:"employee_id,omitempty" gorm:"index"` // References employees(id) - NULL during draft
	DraftID    *uint `json:"draft_id,omitempty" gorm:"index"`    // References employee_onboarding_drafts(id) - NULL after completion

	// Employment ID & Email
	EmployeeIDString *string    `json:"employee_id_string,omitempty" gorm:"size:50"` // Employee ID (if assigned during draft)
	OfficialEmail    *string    `json:"official_email,omitempty" gorm:"size:191"`
	DateOfJoining    *time.Time `json:"date_of_joining,omitempty"`

	// Organization Structure
	DepartmentID       *uint   `json:"department_id,omitempty" gorm:"index"`
	PositionID         *uint   `json:"position_id,omitempty" gorm:"index"` // Designation
	Grade              *string `json:"grade,omitempty" gorm:"size:20"`
	ReportingManagerID *uint   `json:"reporting_manager_id,omitempty" gorm:"index"` // Manager ID
	EmploymentType     *string `json:"employment_type,omitempty" gorm:"size:50"`    // full_time, part_time, contract, intern

	// Work Location & Shift
	LocationID *uint   `json:"location_id,omitempty" gorm:"index"`
	Shift      *string `json:"shift,omitempty" gorm:"size:50"` // Day, Night, Rotating, etc.
	WorkPhone  *string `json:"work_phone,omitempty" gorm:"size:50"`

	// Probation Period
	ProbationPeriodDays      *int       `json:"probation_period_days,omitempty"` // Minimum 90 days
	ExpectedConfirmationDate *time.Time `json:"expected_confirmation_date,omitempty"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (EmployeeEmploymentDetails) TableName() string {
	return "employee_employment_details"
}

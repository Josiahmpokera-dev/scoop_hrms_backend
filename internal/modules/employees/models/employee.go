package models

import (
	"time"

	"gorm.io/gorm"
)

// NOTE: This model includes common employee fields.
// If your database table has different column names, adjust the gorm tags accordingly.
// For example, if your column is "emp_id" instead of "employee_id", use: gorm:"column:emp_id"

// EmployeeStatus represents employee employment status
type EmployeeStatus string

const (
	StatusActive     EmployeeStatus = "active"
	StatusInactive   EmployeeStatus = "inactive"
	StatusOnLeave    EmployeeStatus = "on_leave"
	StatusSuspended  EmployeeStatus = "suspended"
	StatusTerminated EmployeeStatus = "terminated"
	StatusArchived   EmployeeStatus = "archived"
)

// Employee represents an employee in the system
type Employee struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	TenantID       *uint      `json:"tenant_id,omitempty" gorm:"index"`                // Multi-tenant support
	UserID         *uint      `json:"user_id,omitempty" gorm:"index"`                  // References users(id)
	OrganizationID *uint      `json:"organization_id,omitempty" gorm:"index"`          // References organizations(id)
	EmployeeID     string     `json:"employee_id" gorm:"uniqueIndex;not null;size:50"` // Employee number/ID
	FirstName      string     `json:"first_name" gorm:"not null;size:120"`
	MiddleName     *string    `json:"middle_name,omitempty" gorm:"size:120"`
	LastName       string     `json:"last_name" gorm:"not null;size:120"`
	PreferredName  *string    `json:"preferred_name,omitempty" gorm:"size:120"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty"`
	Gender         *string    `json:"gender,omitempty" gorm:"size:20"`
	MaritalStatus  *string    `json:"marital_status,omitempty" gorm:"size:20"`
	BloodGroup     *string    `json:"blood_group,omitempty" gorm:"size:10"` // A+, B+, O+, AB+, etc.
	Nationality    *string    `json:"nationality,omitempty" gorm:"size:100"`
	WorkEmail      *string    `json:"work_email,omitempty" gorm:"uniqueIndex;size:191"`
	PersonalEmail  *string    `json:"personal_email,omitempty" gorm:"size:191"`
	PhoneNumber    *string    `json:"phone_number,omitempty" gorm:"size:50"`
	AlternatePhone *string    `json:"alternate_phone,omitempty" gorm:"size:50"`

	// Legacy fields for backward compatibility
	Email string `json:"email,omitempty" gorm:"-"` // Computed: work_email or personal_email
	Phone string `json:"phone,omitempty" gorm:"-"` // Computed: phone_number

	// Employment Details
	DepartmentID             *uint          `json:"department_id,omitempty" gorm:"index"`
	TeamID                   *uint          `json:"team_id,omitempty" gorm:"index"`              // References teams(id)
	OrganizationUnitID       *uint          `json:"organization_unit_id,omitempty" gorm:"index"` // References organization_units(id)
	PositionID               *uint          `json:"position_id,omitempty" gorm:"index"`          // References job_positions(id)
	LocationID               *uint          `json:"location_id,omitempty" gorm:"index"`          // References locations(id)
	CostCenterID             *uint          `json:"cost_center_id,omitempty" gorm:"index"`       // References cost_centers(id)
	HireDate                 *time.Time     `json:"hire_date,omitempty"`                         // Date of Joining
	EmploymentType           *string        `json:"employment_type,omitempty" gorm:"size:50"`    // full_time, part_time, contract, intern
	Status                   EmployeeStatus `json:"status" gorm:"type:varchar(50);default:'active'"`
	ManagerID                *uint          `json:"manager_id,omitempty" gorm:"index"`    // References employees(id)
	ReportsToID              *uint          `json:"reports_to_id,omitempty" gorm:"index"` // Direct manager / Reporting Manager
	DottedLineManagerID      *uint          `json:"dotted_line_manager_id,omitempty" gorm:"index"`
	MatrixManagerID          *uint          `json:"matrix_manager_id,omitempty" gorm:"index"`
	Grade                    *string        `json:"grade,omitempty" gorm:"size:20"`
	Level                    *int           `json:"level,omitempty"`
	PhotoURL                 *string        `json:"photo_url,omitempty" gorm:"size:500"`
	Shift                    *string        `json:"shift,omitempty" gorm:"size:50"` // Day, Night, Rotating, etc.
	WorkPhone                *string        `json:"work_phone,omitempty" gorm:"size:50"`
	ProbationPeriodDays      *int           `json:"probation_period_days,omitempty"`      // Probation period in days (min 90)
	ExpectedConfirmationDate *time.Time     `json:"expected_confirmation_date,omitempty"` // Expected confirmation date after probation

	// Salary & Benefits (keeping for backward compatibility)
	Salary   *float64 `json:"salary,omitempty" gorm:"type:decimal(10,2)"`
	Currency string   `json:"currency" gorm:"size:10;default:'USD'"`

	// Emergency Contact (keeping for backward compatibility)
	EmergencyContactName     string `json:"emergency_contact_name" gorm:"size:100"`
	EmergencyContactPhone    string `json:"emergency_contact_phone" gorm:"size:20"`
	EmergencyContactRelation string `json:"emergency_contact_relation" gorm:"size:50"`

	// Additional Info
	Notes    string `json:"notes" gorm:"type:text"`
	IsActive bool   `json:"is_active" gorm:"default:true"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for Employee model
func (Employee) TableName() string {
	return "employees"
}

// FullName returns the full name of the employee
func (e *Employee) FullName() string {
	return e.FirstName + " " + e.LastName
}

// IsActiveStatus checks if employee is active
func (e *Employee) IsActiveStatus() bool {
	return e.Status == StatusActive && e.IsActive
}

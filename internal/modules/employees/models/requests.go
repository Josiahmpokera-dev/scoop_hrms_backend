package models

import "time"

// OnboardEmployeeRequest represents the request to onboard a new employee
type OnboardEmployeeRequest struct {
	EmployeeID    string     `json:"employee_id" binding:"required,min=3"`
	FirstName     string     `json:"first_name" binding:"required,min=2"`
	LastName      string     `json:"last_name" binding:"required,min=2"`
	Email         string     `json:"email" binding:"required,email"`
	Phone         string     `json:"phone,omitempty"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`
	Gender        string     `json:"gender,omitempty"`
	Address       string     `json:"address,omitempty"`
	City          string     `json:"city,omitempty"`
	State         string     `json:"state,omitempty"`
	Country       string     `json:"country,omitempty"`
	PostalCode    string     `json:"postal_code,omitempty"`
	
	// Employment Details
	DepartmentID  *uint       `json:"department_id,omitempty"`
	PositionID    *uint       `json:"position_id,omitempty"`
	LocationID    *uint       `json:"location_id,omitempty"`
	HireDate      *time.Time  `json:"hire_date,omitempty"`
	EmploymentType string     `json:"employment_type,omitempty"`
	
	// Salary
	Salary        *float64    `json:"salary,omitempty"`
	Currency      string      `json:"currency,omitempty"`
	
	// Emergency Contact
	EmergencyContactName     string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone    string `json:"emergency_contact_phone,omitempty"`
	EmergencyContactRelation string `json:"emergency_contact_relation,omitempty"`
	
	// Additional
	Notes         string      `json:"notes,omitempty"`
}

// UpdateEmployeeRequest represents the request to update an employee
type UpdateEmployeeRequest struct {
	FirstName     *string     `json:"first_name,omitempty"`
	LastName      *string     `json:"last_name,omitempty"`
	Email         *string     `json:"email,omitempty" binding:"omitempty,email"`
	Phone         *string     `json:"phone,omitempty"`
	DateOfBirth   *time.Time  `json:"date_of_birth,omitempty"`
	Gender        *string     `json:"gender,omitempty"`
	Address       *string     `json:"address,omitempty"`
	City          *string     `json:"city,omitempty"`
	State         *string     `json:"state,omitempty"`
	Country       *string     `json:"country,omitempty"`
	PostalCode    *string     `json:"postal_code,omitempty"`
	
	DepartmentID  *uint       `json:"department_id,omitempty"`
	PositionID    *uint       `json:"position_id,omitempty"`
	LocationID    *uint       `json:"location_id,omitempty"`
	EmploymentType *string    `json:"employment_type,omitempty"`
	Status        *string     `json:"status,omitempty"`
	Salary        *float64    `json:"salary,omitempty"`
	Currency      *string     `json:"currency,omitempty"`
	
	EmergencyContactName     *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone    *string `json:"emergency_contact_phone,omitempty"`
	EmergencyContactRelation *string `json:"emergency_contact_relation,omitempty"`
	
	Notes         *string     `json:"notes,omitempty"`
	IsActive      *bool       `json:"is_active,omitempty"`
}

// TerminateEmployeeRequest represents the request to terminate an employee
type TerminateEmployeeRequest struct {
	ID             uint       `json:"id" binding:"required"`
	Reason         *string    `json:"reason,omitempty"`           // Reason for termination
	TerminationDate *time.Time `json:"termination_date,omitempty"` // Date of termination (defaults to today)
}

// SuspendEmployeeRequest represents the request to suspend an employee
type SuspendEmployeeRequest struct {
	ID              uint       `json:"id" binding:"required"`
	Reason          *string    `json:"reason,omitempty"`           // Reason for suspension
	SuspensionEndDate *time.Time `json:"suspension_end_date,omitempty"` // Expected end date of suspension
}

// ArchiveEmployeeRequest represents the request to archive an employee
type ArchiveEmployeeRequest struct {
	ID     uint    `json:"id" binding:"required"`
	Reason *string `json:"reason,omitempty"` // Reason for archiving
}

// ReactivateEmployeeRequest represents the request to reactivate an employee
type ReactivateEmployeeRequest struct {
	ID uint `json:"id" binding:"required"`
}

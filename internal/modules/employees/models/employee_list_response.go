package models

import "time"

// EmployeeListResponse represents the employee data for list view
type EmployeeListResponse struct {
	ID                uint       `json:"id"`
	FullName          string     `json:"full_name"`           // First Name + Last Name
	EmployeeID       string     `json:"employee_id"`         // Employee ID or code
	Email             *string    `json:"email,omitempty"`     // Work email or personal email
	Department        *string    `json:"department,omitempty"` // Department name
	DepartmentID      *uint      `json:"department_id,omitempty"`
	Designation       *string    `json:"designation,omitempty"` // Position/Job title
	PositionID        *uint      `json:"position_id,omitempty"`
	Status            string     `json:"status"`              // active, inactive, on_leave, terminated
	Location          *string    `json:"location,omitempty"`  // Location name
	LocationID        *uint      `json:"location_id,omitempty"`
	Manager           *string    `json:"manager,omitempty"`   // Manager full name
	ManagerID         *uint      `json:"manager_id,omitempty"`
	DateOfJoining     *time.Time `json:"date_of_joining,omitempty"` // Hire date
	OnboardingStatus  string     `json:"onboarding_status"`   // completed, in_progress, not_started
	OnboardingProgress *float64  `json:"onboarding_progress,omitempty"` // 0-100 if in_progress
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

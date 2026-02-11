package models

import (
	"time"

	"gorm.io/gorm"
)

// RosterStatus represents the status of a roster assignment
type RosterStatus string

const (
	RosterStatusDraft     RosterStatus = "draft"
	RosterStatusPublished RosterStatus = "published"
	RosterStatusSwapped   RosterStatus = "swapped"
)

// RosterAssignment represents a roster assignment for an employee
type RosterAssignment struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	EmployeeID string         `json:"employee_id" gorm:"not null;size:50;index"` // References employees(employee_id)
	Date        time.Time      `json:"date" gorm:"type:date;not null;index"`
	ShiftID     uint           `json:"shift_id" gorm:"not null;index"` // References shifts(id)
	LocationID *uint          `json:"location_id,omitempty" gorm:"index"` // References locations(id)
	Status      string         `json:"status" gorm:"not null;size:20;default:'draft'"` // draft, published, swapped
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedBy   *uint          `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy   *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Shift    *Shift            `json:"shift,omitempty" gorm:"foreignKey:ShiftID"`
	Location *RosterLocationRef `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// RosterLocationRef is a simplified location reference for roster assignments
type RosterLocationRef struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

// TableName specifies the table name
func (RosterLocationRef) TableName() string {
	return "locations"
}

// TableName specifies the table name
func (RosterAssignment) TableName() string {
	return "roster_assignments"
}

// CreateRosterAssignmentRequest represents the request to create a roster assignment
type CreateRosterAssignmentRequest struct {
	EmployeeID string  `json:"employee_id" binding:"required"`
	Date       string  `json:"date" binding:"required"` // YYYY-MM-DD format
	ShiftID    uint    `json:"shift_id" binding:"required"`
	LocationID *uint   `json:"location_id"`
	Status     *string `json:"status"` // draft, published, swapped
}

// BulkCreateRosterAssignmentRequest represents the request to bulk create roster assignments
type BulkCreateRosterAssignmentRequest struct {
	Assignments []CreateRosterAssignmentRequest `json:"assignments" binding:"required"`
}

// UpdateRosterAssignmentRequest represents the request to update a roster assignment
type UpdateRosterAssignmentRequest struct {
	ShiftID    *uint   `json:"shift_id"`
	LocationID *uint   `json:"location_id"`
	Status     *string `json:"status"`
}

// PublishRosterRequest represents the request to publish rosters
type PublishRosterRequest struct {
	StartDate            string `json:"start_date" binding:"required"` // YYYY-MM-DD format
	EndDate              string `json:"end_date" binding:"required"`   // YYYY-MM-DD format
	NotifyEmployees      *bool  `json:"notify_employees"`
	NotificationMessage  *string `json:"notification_message"`
}

// AutoScheduleRequest represents the request for auto-scheduling
type AutoScheduleRequest struct {
	StartDate        string                `json:"start_date" binding:"required"` // YYYY-MM-DD format
	EndDate          string                `json:"end_date" binding:"required"`   // YYYY-MM-DD format
	EmployeeIDs      []string              `json:"employee_ids"`
	DepartmentID     *uint                 `json:"department_id"`
	LocationID       *uint                 `json:"location_id"`
	ShiftPreferences map[string][]uint     `json:"shift_preferences"` // employee_id -> []shift_ids
	Rules            *AutoScheduleRules    `json:"rules"`
}

// AutoScheduleRules represents rules for auto-scheduling
type AutoScheduleRules struct {
	MaxConsecutiveDays   *int  `json:"max_consecutive_days"`
	MinRestDaysPerWeek   *int  `json:"min_rest_days_per_week"`
	PreferWeekendOff     *bool `json:"prefer_weekend_off"`
}

// WeeklyRosterRequest represents the request for weekly roster view
type WeeklyRosterRequest struct {
	Week         string   `json:"week" binding:"required"` // YYYY-WW format
	EmployeeIDs  []string `json:"employee_ids"`
	DepartmentID *uint    `json:"department_id"`
	LocationID   *uint    `json:"location_id"`
}

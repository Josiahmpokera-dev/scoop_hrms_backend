package models

import (
	"time"

	"gorm.io/gorm"
)

// RosterChangeRequestStatus represents the status of a roster change request
type RosterChangeRequestStatus string

const (
	RosterChangeRequestStatusPending  RosterChangeRequestStatus = "pending"
	RosterChangeRequestStatusApproved RosterChangeRequestStatus = "approved"
	RosterChangeRequestStatusRejected RosterChangeRequestStatus = "rejected"
)

// RosterChangeRequest represents a request by an employee to change their roster assignment
type RosterChangeRequest struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	TenantID        *uint          `json:"tenant_id,omitempty" gorm:"index"`
	RequestedBy     string         `json:"requested_by" gorm:"not null;size:50;index"` // References employees(employee_id)
	AssignmentID    uint           `json:"assignment_id" gorm:"not null;index"` // References roster_assignments(id)
	RequestedShiftID *uint         `json:"requested_shift_id,omitempty" gorm:"index"` // New shift (optional - if not provided, keeps current)
	RequestedDate   *time.Time     `json:"requested_date,omitempty"` // New date (optional - if not provided, keeps current)
	RequestedLocationID *uint       `json:"requested_location_id,omitempty" gorm:"index"` // New location (optional)
	Reason          string         `json:"reason" gorm:"type:text"`
	Status          string         `json:"status" gorm:"not null;size:20;default:'pending'"` // pending, approved, rejected
	RejectionReason *string        `json:"rejection_reason,omitempty" gorm:"type:text"`
	RequestedAt     time.Time      `json:"requested_at" gorm:"default:CURRENT_TIMESTAMP"`
	ReviewedAt      *time.Time     `json:"reviewed_at,omitempty"`
	ReviewedBy      *uint          `json:"reviewed_by,omitempty" gorm:"index"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Assignment      *RosterAssignment `json:"assignment,omitempty" gorm:"foreignKey:AssignmentID"`
	RequestedShift  *Shift            `json:"requested_shift,omitempty" gorm:"foreignKey:RequestedShiftID"`
}

// TableName specifies the table name
func (RosterChangeRequest) TableName() string {
	return "roster_change_requests"
}

// CreateRosterChangeRequestRequest represents the request to create a roster change request
type CreateRosterChangeRequestRequest struct {
	AssignmentID       uint   `json:"assignment_id" binding:"required"`
	RequestedShiftID   *uint  `json:"requested_shift_id"` // Optional - new shift
	RequestedDate     *string `json:"requested_date"`     // Optional - new date (YYYY-MM-DD format)
	RequestedLocationID *uint  `json:"requested_location_id"` // Optional - new location
	Reason            string `json:"reason"`
}

// ApproveRosterChangeRequestRequest represents the request to approve a roster change request
type ApproveRosterChangeRequestRequest struct {
	NotifyEmployee *bool   `json:"notify_employee"`
	Notes          *string `json:"notes"`
}

// RejectRosterChangeRequestRequest represents the request to reject a roster change request
type RejectRosterChangeRequestRequest struct {
	RejectionReason *string `json:"rejection_reason" binding:"required"`
	NotifyEmployee  *bool   `json:"notify_employee"`
}

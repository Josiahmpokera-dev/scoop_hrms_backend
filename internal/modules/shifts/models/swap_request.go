package models

import (
	"time"

	"gorm.io/gorm"
)

// SwapRequestStatus represents the status of a swap request
type SwapRequestStatus string

const (
	SwapRequestStatusPending  SwapRequestStatus = "pending"
	SwapRequestStatusApproved SwapRequestStatus = "approved"
	SwapRequestStatusRejected SwapRequestStatus = "rejected"
)

// SwapRequest represents a shift swap request between two employees
type SwapRequest struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	RequestedBy         string         `json:"requested_by" gorm:"not null;size:50;index"` // References employees(employee_id)
	RequestedWith       string         `json:"requested_with" gorm:"not null;size:50;index"` // References employees(employee_id)
	AssignmentID        uint           `json:"assignment_id" gorm:"not null;index"` // References roster_assignments(id)
	SwapAssignmentID    uint           `json:"swap_assignment_id" gorm:"not null;index"` // References roster_assignments(id)
	Reason              string         `json:"reason" gorm:"type:text"`
	Status              string         `json:"status" gorm:"not null;size:20;default:'pending'"` // pending, approved, rejected
	RejectionReason     *string        `json:"rejection_reason,omitempty" gorm:"type:text"`
	RequestedAt         time.Time      `json:"requested_at" gorm:"default:CURRENT_TIMESTAMP"`
	ReviewedAt          *time.Time     `json:"reviewed_at,omitempty"`
	ReviewedBy          *uint          `json:"reviewed_by,omitempty" gorm:"index"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Assignment     *RosterAssignment `json:"assignment,omitempty" gorm:"foreignKey:AssignmentID"`
	SwapAssignment *RosterAssignment `json:"swap_assignment,omitempty" gorm:"foreignKey:SwapAssignmentID"`
}

// TableName specifies the table name
func (SwapRequest) TableName() string {
	return "swap_requests"
}

// CreateSwapRequestRequest represents the request to create a swap request
type CreateSwapRequestRequest struct {
	AssignmentID     uint   `json:"assignment_id" binding:"required"`
	SwapAssignmentID uint   `json:"swap_assignment_id" binding:"required"`
	Reason           string `json:"reason"`
}

// ApproveSwapRequestRequest represents the request to approve a swap request
type ApproveSwapRequestRequest struct {
	NotifyEmployees *bool   `json:"notify_employees"`
	Notes           *string `json:"notes"`
}

// RejectSwapRequestRequest represents the request to reject a swap request
type RejectSwapRequestRequest struct {
	RejectionReason *string `json:"rejection_reason" binding:"required"`
	NotifyEmployees *bool   `json:"notify_employees"`
}

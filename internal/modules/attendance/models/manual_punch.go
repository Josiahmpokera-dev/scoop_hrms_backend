package models

import (
	"time"

	"gorm.io/gorm"
)

// ManualPunchStatus represents the status of a manual punch request
type ManualPunchStatus string

const (
	ManualPunchStatusPending   ManualPunchStatus = "pending"
	ManualPunchStatusApproved  ManualPunchStatus = "approved"
	ManualPunchStatusRejected  ManualPunchStatus = "rejected"
	ManualPunchStatusCancelled ManualPunchStatus = "cancelled"
)

// ManualPunchType represents the type of manual punch (in/out)
type ManualPunchType string

const (
	ManualPunchTypeIn  ManualPunchType = "in"
	ManualPunchTypeOut ManualPunchType = "out"
)

// ManualPunch represents a manual attendance punch request
// This allows employees to request manual time entries when biometric/fingerprint fails
type ManualPunch struct {
	ID                 uint              `json:"id" gorm:"primaryKey"`
	EmployeeID         uint              `json:"employee_id" gorm:"index;not null"`           // References employees(id)
	PunchDate          time.Time         `json:"punch_date" gorm:"type:date;not null;index"`  // Date of the punch
	PunchTime          time.Time         `json:"punch_time" gorm:"not null"`                  // Actual time of punch
	PunchType          ManualPunchType   `json:"punch_type" gorm:"type:varchar(10);not null"` // in or out
	Status             ManualPunchStatus `json:"status" gorm:"type:varchar(20);default:'pending';not null;index"`
	Reason             string            `json:"reason" gorm:"type:text;not null"`              // Reason for manual punch
	SupportingDocument *string           `json:"supporting_document,omitempty" gorm:"size:255"` // Optional document URL
	ApproverID         *uint             `json:"approver_id,omitempty" gorm:"index"`            // References employees(id) - who approved/rejected
	ApprovedAt         *time.Time        `json:"approved_at,omitempty"`
	RejectionReason    *string           `json:"rejection_reason,omitempty" gorm:"type:text"` // Reason for rejection if applicable

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Virtual fields (not stored in database)
	EmployeeName string `json:"employee_name,omitempty" gorm:"-"`
	ApproverName string `json:"approver_name,omitempty" gorm:"-"`
}

// ManualPunchRequest represents the request payload for creating a manual punch
type ManualPunchRequest struct {
	PunchType   ManualPunchType `json:"punch_type" binding:"required,oneof=in out"`
	Reason      string          `json:"reason" binding:"required,min=10,max=500"`
	SupportingDocument *string  `json:"supporting_document,omitempty"`
}

// ManualPunchApprovalRequest represents the request payload for approving/rejecting a manual punch
type ManualPunchApprovalRequest struct {
	Status          ManualPunchStatus `json:"status" binding:"required,oneof=approved rejected"`
	RejectionReason *string           `json:"rejection_reason,omitempty"` // Required if status is rejected
}

// ManualPunchFilter represents filter criteria for querying manual punches
type ManualPunchFilter struct {
	EmployeeID *uint             `json:"employee_id,omitempty"`
	Status     ManualPunchStatus `json:"status,omitempty"`
	PunchType  ManualPunchType   `json:"punch_type,omitempty"`
	StartDate  *time.Time        `json:"start_date,omitempty"`
	EndDate    *time.Time        `json:"end_date,omitempty"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

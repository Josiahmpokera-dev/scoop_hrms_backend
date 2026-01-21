package models

import (
	"time"

	"gorm.io/gorm"
)

// ProfileUpdateRequestStatus represents the status of a profile update request
type ProfileUpdateRequestStatus string

const (
	ProfileUpdateStatusPending  ProfileUpdateRequestStatus = "pending_approval"
	ProfileUpdateStatusApproved ProfileUpdateRequestStatus = "approved"
	ProfileUpdateStatusRejected ProfileUpdateRequestStatus = "rejected"
)

// ProfileUpdateSection represents the section of profile being updated
type ProfileUpdateSection string

const (
	ProfileSectionPersonal         ProfileUpdateSection = "personal"
	ProfileSectionEmergencyContacts ProfileUpdateSection = "emergency_contacts"
	ProfileSectionAddress          ProfileUpdateSection = "address"
)

// ProfileUpdateRequest represents a profile update request that requires manager approval
type ProfileUpdateRequest struct {
	ID                uint                      `json:"id" gorm:"primaryKey"`
	TenantID          *uint                     `json:"tenant_id,omitempty" gorm:"index"`
	EmployeeID        uint                      `json:"employee_id" gorm:"index;not null"`
	UpdateRequestID   string                    `json:"update_request_id" gorm:"uniqueIndex;not null;size:50"` // e.g., PRU-2026-001
	Section           ProfileUpdateSection      `json:"section" gorm:"type:varchar(50);not null"`              // personal, emergency_contacts, address
	Status            ProfileUpdateRequestStatus `json:"status" gorm:"type:varchar(50);default:'pending_approval'"`
	UpdatesJSON       string                    `json:"-" gorm:"type:text;not null"` // JSON object with the updates
	Reason            *string                  `json:"reason,omitempty" gorm:"type:text"`
	
	// Approval
	ApproverID        *uint      `json:"approver_id,omitempty" gorm:"index"` // Manager/HR who approves
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	RejectedAt        *time.Time `json:"rejected_at,omitempty"`
	RejectionReason   *string    `json:"rejection_reason,omitempty" gorm:"type:text"`
	
	// Timestamps
	SubmittedAt       time.Time      `json:"submitted_at" gorm:"not null"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	UpdatedBy         *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (ProfileUpdateRequest) TableName() string {
	return "profile_update_requests"
}

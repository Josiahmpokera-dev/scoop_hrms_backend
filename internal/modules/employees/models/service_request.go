package models

import (
	"time"

	"gorm.io/gorm"
)

// ServiceRequestType represents the type of service request
type ServiceRequestType string

const (
	ServiceRequestTypeHRLetter   ServiceRequestType = "hr_letter"
	ServiceRequestTypeITRequest  ServiceRequestType = "it_request"
	ServiceRequestTypeFacilities ServiceRequestType = "facilities"
)

// ServiceRequestStatus represents the status of a service request
type ServiceRequestStatus string

const (
	ServiceRequestStatusDraft       ServiceRequestStatus = "draft"
	ServiceRequestStatusSubmitted   ServiceRequestStatus = "submitted"
	ServiceRequestStatusInProgress  ServiceRequestStatus = "in_progress"
	ServiceRequestStatusApproved    ServiceRequestStatus = "approved"
	ServiceRequestStatusRejected    ServiceRequestStatus = "rejected"
	ServiceRequestStatusCompleted   ServiceRequestStatus = "completed"
	ServiceRequestStatusCancelled  ServiceRequestStatus = "cancelled"
)

// ServiceRequestPriority represents the priority of a service request
type ServiceRequestPriority string

const (
	ServiceRequestPriorityLow    ServiceRequestPriority = "low"
	ServiceRequestPriorityMedium ServiceRequestPriority = "medium"
	ServiceRequestPriorityHigh  ServiceRequestPriority = "high"
	ServiceRequestPriorityUrgent ServiceRequestPriority = "urgent"
)

// ServiceRequest represents a service request (HR Letter, IT Request, Facilities)
type ServiceRequest struct {
	ID            uint                  `json:"id" gorm:"primaryKey"`
	EmployeeID    uint                  `json:"employee_id" gorm:"index;not null"`
	RequestNumber string                `json:"request_number" gorm:"uniqueIndex;not null;size:50"` // e.g., SR-2026-001
	Type          ServiceRequestType    `json:"type" gorm:"type:varchar(50);not null"`              // hr_letter, it_request, facilities
	Category      string                `json:"category" gorm:"size:100"`                            // Employment Verification, Software Access, etc.
	Subject       string                `json:"subject" gorm:"size:255;not null"`
	Description   *string               `json:"description,omitempty" gorm:"type:text"`
	Status        ServiceRequestStatus   `json:"status" gorm:"type:varchar(50);default:'submitted'"`
	Priority      ServiceRequestPriority `json:"priority" gorm:"type:varchar(20);default:'medium'"`
	
	// HR Letter specific fields
	LetterType    *string `json:"letter_type,omitempty" gorm:"size:100"`    // Employment Verification, Salary Certificate, etc.
	Purpose       *string `json:"purpose,omitempty" gorm:"type:text"`      // Purpose of the letter
	AddressedTo   *string `json:"addressed_to,omitempty" gorm:"size:255"`   // To Whom It May Concern, etc.
	AdditionalNotes *string `json:"additional_notes,omitempty" gorm:"type:text"`
	
	// IT Request / Facilities specific fields
	RequestedItemsJSON *string `json:"-" gorm:"type:text"` // JSON array of requested items
	
	// Assignment
	AssignedToID  *uint   `json:"assigned_to_id,omitempty" gorm:"index"` // User/Team ID
	AssignedTo    *string `json:"assigned_to,omitempty" gorm:"size:255"` // Team name or user name
	
	// SLA
	SLAHours      *int       `json:"sla_hours,omitempty"` // Expected processing time in hours
	EstimatedCompletionDate *time.Time `json:"estimated_completion_date,omitempty"`
	
	// Document
	DocumentID    *uint   `json:"document_id,omitempty" gorm:"index"` // Reference to generated document
	DocumentURL   *string `json:"document_url,omitempty" gorm:"size:500"`
	FileName      *string `json:"file_name,omitempty" gorm:"size:255"`
	
	// Timeline
	RequestedDate time.Time  `json:"requested_date" gorm:"not null"`
	CompletedDate *time.Time `json:"completed_date,omitempty"`
	CancelledDate *time.Time `json:"cancelled_date,omitempty"`
	CancelledReason *string  `json:"cancelled_reason,omitempty" gorm:"type:text"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (ServiceRequest) TableName() string {
	return "service_requests"
}

// RequestedItem represents an item requested in IT/Facilities requests
type RequestedItem struct {
	ItemType      string `json:"item_type"`      // software, hardware, access_card, etc.
	ItemName      string `json:"item_name"`      // Adobe Creative Suite, Laptop, etc.
	Justification string `json:"justification"`  // Why this item is needed
}

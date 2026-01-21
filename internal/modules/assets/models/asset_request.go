package models

import (
	"time"

	"gorm.io/gorm"
)

// AssetRequestType represents the type of asset request
type AssetRequestType string

const (
	AssetRequestTypeNew      AssetRequestType = "new"      // Request for new asset
	AssetRequestTypeReplacement AssetRequestType = "replacement" // Request to replace existing asset
	AssetRequestTypeAdditional AssetRequestType = "additional" // Request for additional asset
)

// AssetRequestStatus represents the status of an asset request
type AssetRequestStatus string

const (
	AssetRequestStatusPending   AssetRequestStatus = "pending"   // Awaiting HR review
	AssetRequestStatusApproved  AssetRequestStatus = "approved"  // Approved by HR
	AssetRequestStatusRejected  AssetRequestStatus = "rejected"  // Rejected by HR
	AssetRequestStatusFulfilled AssetRequestStatus = "fulfilled" // Asset assigned
	AssetRequestStatusCancelled AssetRequestStatus = "cancelled" // Cancelled by employee
)

// AssetRequestPriority represents the priority of an asset request
type AssetRequestPriority string

const (
	AssetRequestPriorityLow    AssetRequestPriority = "low"
	AssetRequestPriorityMedium AssetRequestPriority = "medium"
	AssetRequestPriorityHigh   AssetRequestPriority = "high"
	AssetRequestPriorityUrgent AssetRequestPriority = "urgent"
)

// AssetRequest represents a request from an employee to HR for an asset
type AssetRequest struct {
	ID            uint                 `json:"id" gorm:"primaryKey"`
	TenantID      *uint                `json:"tenant_id,omitempty" gorm:"index"`
	EmployeeID    string               `json:"employee_id" gorm:"index;not null;size:50"` // Employee ID (e.g., EMP001)
	RequestNumber string               `json:"request_number" gorm:"uniqueIndex;not null;size:50"` // e.g., AR-2026-001
	RequestType   AssetRequestType     `json:"request_type" gorm:"type:varchar(50);not null"` // new, replacement, additional
	AssetType     string               `json:"asset_type" gorm:"not null;size:50"` // laptop, mobile_phone, etc.
	Brand         *string              `json:"brand,omitempty" gorm:"size:100"` // Preferred brand (optional)
	Model         *string              `json:"model,omitempty" gorm:"size:100"` // Preferred model (optional)
	Priority      AssetRequestPriority `json:"priority" gorm:"type:varchar(20);default:'medium'"`
	Status        AssetRequestStatus   `json:"status" gorm:"type:varchar(50);default:'pending'"`
	Justification string               `json:"justification" gorm:"type:text;not null"` // Why the asset is needed
	Notes         *string             `json:"notes,omitempty" gorm:"type:text"`
	
	// Replacement/Additional specific
	ReplacingAssetID *uint   `json:"replacing_asset_id,omitempty" gorm:"index"` // If replacing existing asset
	ReplacingAssetCode *string `json:"replacing_asset_code,omitempty" gorm:"size:50"` // Asset code being replaced
	
	// Approval
	ApprovedBy    *uint      `json:"approved_by,omitempty" gorm:"index"` // HR user ID
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	RejectedAt    *time.Time `json:"rejected_at,omitempty"`
	RejectionReason *string  `json:"rejection_reason,omitempty" gorm:"type:text"`
	
	// Fulfillment
	FulfilledAt   *time.Time `json:"fulfilled_at,omitempty"`
	AssignedAssetID *uint     `json:"assigned_asset_id,omitempty" gorm:"index"` // Asset that was assigned
	AssignedAssetCode *string  `json:"assigned_asset_code,omitempty" gorm:"size:50"`
	
	// Timeline
	RequestedDate time.Time  `json:"requested_date" gorm:"not null"`
	CancelledDate *time.Time `json:"cancelled_date,omitempty"`
	CancelledReason *string  `json:"cancelled_reason,omitempty" gorm:"type:text"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (AssetRequest) TableName() string {
	return "asset_requests"
}

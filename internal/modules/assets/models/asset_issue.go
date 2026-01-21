package models

import (
	"time"

	"gorm.io/gorm"
)

// AssetIssueType represents the type of asset issue
type AssetIssueType string

const (
	AssetIssueTypeMalfunction AssetIssueType = "malfunction" // Asset not working properly
	AssetIssueTypeDamage      AssetIssueType = "damage"      // Asset is damaged
	AssetIssueTypeLost        AssetIssueType = "lost"        // Asset is lost
	AssetIssueTypeStolen      AssetIssueType = "stolen"      // Asset is stolen
	AssetIssueTypeOther       AssetIssueType = "other"        // Other issues
)

// AssetIssueStatus represents the status of an asset issue
type AssetIssueStatus string

const (
	AssetIssueStatusReported  AssetIssueStatus = "reported"  // Issue reported by employee
	AssetIssueStatusUnderReview AssetIssueStatus = "under_review" // HR/IT reviewing
	AssetIssueStatusInProgress AssetIssueStatus = "in_progress" // Being resolved
	AssetIssueStatusResolved  AssetIssueStatus = "resolved"  // Issue resolved
	AssetIssueStatusClosed    AssetIssueStatus = "closed"    // Issue closed
)

// AssetIssuePriority represents the priority of an asset issue
type AssetIssuePriority string

const (
	AssetIssuePriorityLow    AssetIssuePriority = "low"
	AssetIssuePriorityMedium AssetIssuePriority = "medium"
	AssetIssuePriorityHigh   AssetIssuePriority = "high"
	AssetIssuePriorityUrgent AssetIssuePriority = "urgent"
)

// AssetIssue represents an issue reported by an employee regarding an assigned asset
type AssetIssue struct {
	ID            uint               `json:"id" gorm:"primaryKey"`
	TenantID      *uint              `json:"tenant_id,omitempty" gorm:"index"`
	EmployeeID    string             `json:"employee_id" gorm:"index;not null;size:50"` // Employee ID (e.g., EMP001)
	IssueNumber   string             `json:"issue_number" gorm:"uniqueIndex;not null;size:50"` // e.g., AI-2026-001
	AssetID       uint               `json:"asset_id" gorm:"index;not null"` // References assets(id)
	AssetCode     string             `json:"asset_code" gorm:"not null;size:50"` // Asset code for quick reference
	IssueType     AssetIssueType     `json:"issue_type" gorm:"type:varchar(50);not null"` // malfunction, damage, lost, stolen, other
	Priority      AssetIssuePriority `json:"priority" gorm:"type:varchar(20);default:'medium'"`
	Status        AssetIssueStatus   `json:"status" gorm:"type:varchar(50);default:'reported'"`
	Title         string             `json:"title" gorm:"not null;size:255"` // Brief title of the issue
	Description   string             `json:"description" gorm:"type:text;not null"` // Detailed description
	ReportedDate  time.Time          `json:"reported_date" gorm:"not null"`
	
	// Resolution
	ResolvedAt    *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy    *uint      `json:"resolved_by,omitempty" gorm:"index"` // HR/IT user ID
	ResolutionNotes *string  `json:"resolution_notes,omitempty" gorm:"type:text"`
	ActionTaken   *string    `json:"action_taken,omitempty" gorm:"type:text"` // What action was taken
	
	// For lost/stolen assets
	IncidentDate  *time.Time `json:"incident_date,omitempty"` // When the incident occurred
	IncidentLocation *string `json:"incident_location,omitempty" gorm:"type:text"` // Where it happened
	PoliceReportNumber *string `json:"police_report_number,omitempty" gorm:"size:100"` // If police report was filed
	
	// Attachments (stored as JSON array of file URLs)
	AttachmentURLsJSON *string `json:"-" gorm:"type:text"` // JSON array of attachment URLs
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (AssetIssue) TableName() string {
	return "asset_issues"
}

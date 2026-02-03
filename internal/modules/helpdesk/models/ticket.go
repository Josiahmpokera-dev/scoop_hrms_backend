package models

import (
	"time"

	"gorm.io/gorm"
)

// TicketStatus represents the status of a helpdesk ticket
type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "Open"
	TicketStatusInProgress TicketStatus = "In Progress"
	TicketStatusPending    TicketStatus = "Pending"
	TicketStatusResolved   TicketStatus = "Resolved"
	TicketStatusClosed     TicketStatus = "Closed"
	TicketStatusCancelled  TicketStatus = "Cancelled"
)

// TicketPriority represents the priority of a helpdesk ticket
type TicketPriority string

const (
	TicketPriorityCritical TicketPriority = "Critical"
	TicketPriorityHigh     TicketPriority = "High"
	TicketPriorityMedium   TicketPriority = "Medium"
	TicketPriorityLow      TicketPriority = "Low"
)

// TicketChannel represents the channel through which the ticket was created
type TicketChannel string

const (
	TicketChannelPortal   TicketChannel = "Portal"
	TicketChannelEmail    TicketChannel = "Email"
	TicketChannelWhatsApp TicketChannel = "WhatsApp"
	TicketChannelPhone    TicketChannel = "Phone"
	TicketChannelChat     TicketChannel = "Chat"
)

// Ticket represents a helpdesk ticket
type Ticket struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	TenantID     *uint          `json:"tenant_id,omitempty" gorm:"index"`
	TicketNumber string         `json:"ticket_number" gorm:"uniqueIndex;not null;size:50"` // e.g., HD-2026-001
	Title        string         `json:"title" gorm:"not null;size:255"`
	Description  string         `json:"description" gorm:"type:text;not null"`
	Category     string         `json:"category" gorm:"not null;size:100"` // IT, HR, Payroll, Facilities, Leave, General
	SubCategory  *string        `json:"sub_category,omitempty" gorm:"size:100"`
	Priority     TicketPriority `json:"priority" gorm:"type:varchar(20);default:'Medium'"`
	Status       TicketStatus   `json:"status" gorm:"type:varchar(50);default:'Open'"`
	Channel      TicketChannel  `json:"channel" gorm:"type:varchar(50);default:'Portal'"`

	// Requester: either Employee (RequesterID + RequesterEmployeeID) or User-only e.g. Admin (RequesterUserID + RequesterName/Email)
	RequesterID         uint    `json:"requester_id" gorm:"index"`                      // Employee ID (0 when requester is user-only)
	RequesterEmployeeID *string `json:"requester_employee_id,omitempty" gorm:"size:50"` // Employee ID string (e.g., EMP001)
	RequesterUserID     *uint   `json:"requester_user_id,omitempty" gorm:"index"`       // User ID when requester has no employee (e.g. Admin)
	RequesterName       *string `json:"requester_name,omitempty" gorm:"size:255"`       // Display name when user-only requester
	RequesterEmail      *string `json:"requester_email,omitempty" gorm:"size:255"`      // Email when user-only requester

	// Assignment
	AssignedToID    *uint   `json:"assigned_to_id,omitempty" gorm:"index"` // User ID of agent
	AssignedToName  *string `json:"assigned_to_name,omitempty" gorm:"size:255"`
	AssignedToEmail *string `json:"assigned_to_email,omitempty" gorm:"size:255"`
	AssignedToTeam  *string `json:"assigned_to_team,omitempty" gorm:"size:100"` // HR Support, IT Support, etc.
	Queue           *string `json:"queue,omitempty" gorm:"size:100"`            // Queue name

	// SLA
	DueDate           *time.Time `json:"due_date,omitempty"`
	SLAStatus         *string    `json:"sla_status,omitempty" gorm:"size:50"` // Within SLA, Breached, At Risk
	FirstResponseTime *time.Time `json:"first_response_time,omitempty"`       // When first response was given
	ResolutionTime    *time.Time `json:"resolution_time,omitempty"`           // When ticket was resolved

	// Resolution
	ResolvedAt        *time.Time `json:"resolved_at,omitempty"`
	ResolvedByID      *uint      `json:"resolved_by_id,omitempty" gorm:"index"`
	ResolutionSummary *string    `json:"resolution_summary,omitempty" gorm:"type:text"`
	InternalNote      *string    `json:"internal_note,omitempty" gorm:"type:text"` // Private note for agents/admin

	// Closure
	ClosedAt   *time.Time `json:"closed_at,omitempty"`
	ClosedByID *uint      `json:"closed_by_id,omitempty" gorm:"index"`

	// CSAT
	CSATRating  *int    `json:"csat_rating,omitempty"` // 1-5
	CSATComment *string `json:"csat_comment,omitempty" gorm:"type:text"`

	// Tags (stored as JSON array)
	TagsJSON *string `json:"-" gorm:"type:text"`

	// Watchers (stored as JSON array of user IDs)
	WatchersJSON *string `json:"-" gorm:"type:text"`

	// Related tickets (stored as JSON array of ticket IDs)
	RelatedTicketsJSON *string `json:"-" gorm:"type:text"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (Ticket) TableName() string {
	return "helpdesk_tickets"
}

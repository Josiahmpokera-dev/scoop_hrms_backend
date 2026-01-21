package models

import (
	"time"

	"gorm.io/gorm"
)

// TicketCategory represents a ticket category configuration
type TicketCategory struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	TenantID          *uint          `json:"tenant_id,omitempty" gorm:"index"`
	Name              string         `json:"name" gorm:"uniqueIndex;not null;size:100"` // IT, HR, Payroll, etc.
	Description       *string        `json:"description,omitempty" gorm:"type:text"`
	
	// SLA Configuration
	SLAFirstResponseMinutes *int      `json:"sla_first_response_minutes,omitempty"` // Minutes for first response
	SLAResolutionHours      *int      `json:"sla_resolution_hours,omitempty"` // Hours for resolution
	
	// Default settings
	DefaultPriority         *string   `json:"default_priority,omitempty" gorm:"size:20"` // Low, Medium, High, Critical
	AutoAssignTo            *string   `json:"auto_assign_to,omitempty" gorm:"size:100"` // Queue or team name
	
	// Required fields (stored as JSON array)
	RequiredFieldsJSON      *string   `json:"-" gorm:"type:text"`
	
	// Templates (stored as JSON array)
	TemplatesJSON           *string   `json:"-" gorm:"type:text"`
	
	// Timestamps
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	UpdatedBy              *uint     `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt              gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (TicketCategory) TableName() string {
	return "helpdesk_ticket_categories"
}

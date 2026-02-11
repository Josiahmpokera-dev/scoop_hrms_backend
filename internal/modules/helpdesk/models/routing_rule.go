package models

import (
	"time"

	"gorm.io/gorm"
)

// RoutingRuleType represents the type of routing target
type RoutingRuleType string

const (
	RoutingRuleTypeQueue RoutingRuleType = "queue"
	RoutingRuleTypeUser  RoutingRuleType = "user"
)

// RoutingRule represents a routing rule for automatic ticket assignment
type RoutingRule struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	Name              string         `json:"name" gorm:"not null;size:255"`
	IsActive          bool           `json:"is_active" gorm:"default:true"`
	
	// Match conditions (stored as JSON)
	MatchConditionsJSON string       `json:"-" gorm:"type:text;not null"` // JSON: {category, subCategory, priority_in}
	
	// Route to (stored as JSON)
	RouteToJSON        string        `json:"-" gorm:"type:text;not null"` // JSON: {type: "queue"|"user", queue?: string, assignee_user_id?: uint}
	
	// Fallback (stored as JSON)
	FallbackJSON       *string       `json:"-" gorm:"type:text"` // JSON: {type: "user", assignee_user_id: uint}
	
	// Escalation (stored as JSON)
	EscalationJSON     *string       `json:"-" gorm:"type:text"` // JSON: {after_minutes_without_first_response: int, escalate_to_user_id: uint}
	
	// Timestamps
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	UpdatedBy          *uint         `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (RoutingRule) TableName() string {
	return "helpdesk_routing_rules"
}

// MatchConditions represents the match conditions for routing
type MatchConditions struct {
	Category    *string   `json:"category,omitempty"`
	SubCategory *string   `json:"subCategory,omitempty"`
	PriorityIn  []string  `json:"priority_in,omitempty"`
}

// RouteTo represents where to route the ticket
type RouteTo struct {
	Type           string  `json:"type"` // "queue" or "user"
	Queue          *string `json:"queue,omitempty"`
	AssigneeUserID *uint   `json:"assignee_user_id,omitempty"`
}

// Fallback represents fallback routing
type Fallback struct {
	Type           string `json:"type"` // "user"
	AssigneeUserID uint   `json:"assignee_user_id"`
}

// Escalation represents escalation rules
type Escalation struct {
	AfterMinutesWithoutFirstResponse int  `json:"after_minutes_without_first_response"`
	EscalateToUserID                uint `json:"escalate_to_user_id"`
}

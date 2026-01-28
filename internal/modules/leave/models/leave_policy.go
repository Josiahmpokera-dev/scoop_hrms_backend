package models

import (
	"time"

	"gorm.io/gorm"
)

// LeavePolicy represents a leave policy configuration
type LeavePolicy struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	TenantID              *uint          `json:"tenant_id,omitempty" gorm:"index"`
	PolicyName            string         `json:"policy_name" gorm:"not null;size:200"`
	Country               string         `json:"country" gorm:"not null;size:100"`
	LeaveTypeCode         string         `json:"leave_type_code" gorm:"not null;size:10;index"` // References leave_types(code)
	Entitlement           int            `json:"entitlement" gorm:"not null"` // Days per year
	AccrualFrequency      string         `json:"accrual_frequency" gorm:"not null;size:50"` // Annual, Monthly, Quarterly
	ProrationOnJoin       bool           `json:"proration_on_join" gorm:"default:true"`
	ProrationOnExit       bool           `json:"proration_on_exit" gorm:"default:true"`
	CarryForward          bool           `json:"carry_forward" gorm:"default:false"`
	CarryForwardLimit     *int            `json:"carry_forward_limit,omitempty"` // Max days that can be carried forward
	CarryForwardExpiry    *time.Time      `json:"carry_forward_expiry,omitempty"` // Expiry date for carried forward leaves
	EncashmentAllowed     bool           `json:"encashment_allowed" gorm:"default:false"`
	EncashmentLimit       *int            `json:"encashment_limit,omitempty"` // Max days that can be encashed
	NegativeBalanceAllowed bool          `json:"negative_balance_allowed" gorm:"default:false"`
	SandwichRules         bool           `json:"sandwich_rules" gorm:"default:false"` // Allow sandwiching weekends/holidays
	HalfDayAllowed        bool           `json:"half_day_allowed" gorm:"default:false"`
	MinimumNoticeDays     int            `json:"minimum_notice_days" gorm:"default:0"` // Minimum days notice required
	MaximumDaysPerRequest *int           `json:"maximum_days_per_request,omitempty"` // Max days per single request
	IsActive              bool           `json:"is_active" gorm:"default:true"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	CreatedBy             *uint          `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy             *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	LeaveType *LeaveType `json:"leave_type,omitempty" gorm:"foreignKey:LeaveTypeCode;references:Code"`
}

// TableName specifies the table name
func (LeavePolicy) TableName() string {
	return "leave_policies"
}

// CreateLeavePolicyRequest represents the request to create a leave policy
type CreateLeavePolicyRequest struct {
	PolicyName            string     `json:"policy_name" binding:"required"`
	Country               string     `json:"country" binding:"required"`
	LeaveTypeCode         string     `json:"leave_type_code" binding:"required"`
	Entitlement           int        `json:"entitlement" binding:"required"`
	AccrualFrequency      string     `json:"accrual_frequency" binding:"required"`
	ProrationOnJoin       *bool      `json:"proration_on_join"`
	ProrationOnExit       *bool      `json:"proration_on_exit"`
	CarryForward          *bool      `json:"carry_forward"`
	CarryForwardLimit     *int       `json:"carry_forward_limit"`
	CarryForwardExpiry    *time.Time `json:"carry_forward_expiry"`
	EncashmentAllowed     *bool      `json:"encashment_allowed"`
	EncashmentLimit       *int       `json:"encashment_limit"`
	NegativeBalanceAllowed *bool     `json:"negative_balance_allowed"`
	SandwichRules         *bool      `json:"sandwich_rules"`
	HalfDayAllowed        *bool      `json:"half_day_allowed"`
	MinimumNoticeDays     *int       `json:"minimum_notice_days"`
	MaximumDaysPerRequest *int       `json:"maximum_days_per_request"`
	IsActive              *bool      `json:"is_active"`
}

// UpdateLeavePolicyRequest represents the request to update a leave policy
type UpdateLeavePolicyRequest struct {
	PolicyName            *string    `json:"policy_name"`
	Entitlement           *int       `json:"entitlement"`
	AccrualFrequency      *string    `json:"accrual_frequency"`
	ProrationOnJoin       *bool      `json:"proration_on_join"`
	ProrationOnExit       *bool      `json:"proration_on_exit"`
	CarryForward          *bool      `json:"carry_forward"`
	CarryForwardLimit     *int       `json:"carry_forward_limit"`
	CarryForwardExpiry    *time.Time `json:"carry_forward_expiry"`
	EncashmentAllowed     *bool      `json:"encashment_allowed"`
	EncashmentLimit       *int       `json:"encashment_limit"`
	NegativeBalanceAllowed *bool     `json:"negative_balance_allowed"`
	SandwichRules         *bool      `json:"sandwich_rules"`
	HalfDayAllowed        *bool      `json:"half_day_allowed"`
	MinimumNoticeDays     *int       `json:"minimum_notice_days"`
	MaximumDaysPerRequest *int       `json:"maximum_days_per_request"`
	IsActive              *bool      `json:"is_active"`
}

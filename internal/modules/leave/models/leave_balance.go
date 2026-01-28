package models

import (
	"time"

	"gorm.io/gorm"
)

// LeaveBalance represents leave balance for an employee
type LeaveBalance struct {
	ID                        uint           `json:"id" gorm:"primaryKey"`
	TenantID                  *uint          `json:"tenant_id,omitempty" gorm:"index"`
	EmployeeID                string         `json:"employee_id" gorm:"not null;size:50;index"` // References employees(employee_id)
	LeaveTypeCode             string         `json:"leave_type_code" gorm:"not null;size:10;index"` // References leave_types(code)
	Year                      int            `json:"year" gorm:"not null;index"` // Leave period year
	Entitlement               float64        `json:"entitlement" gorm:"type:decimal(5,2);not null"` // Total entitlement for the year
	Used                      float64        `json:"used" gorm:"type:decimal(5,2);default:0"` // Days used
	Pending                   float64        `json:"pending" gorm:"type:decimal(5,2);default:0"` // Days in pending requests
	Available                 float64        `json:"available" gorm:"type:decimal(5,2);default:0"` // Available balance
	CarriedForward            float64        `json:"carried_forward" gorm:"type:decimal(5,2);default:0"` // From previous year
	ExpiresOn                 *time.Time     `json:"expires_on,omitempty"` // Expiry date for carried forward balance
	LeaveTakenFromPreviousYear float64       `json:"leave_taken_from_previous_year" gorm:"type:decimal(5,2);default:0"`
	LeaveBalanceFromPreviousYear float64      `json:"leave_balance_from_previous_year" gorm:"type:decimal(5,2);default:0"`
	PreviousDaysUsed          float64        `json:"previous_days_used" gorm:"type:decimal(5,2);default:0"`
	LastAccrualDate           *time.Time     `json:"last_accrual_date,omitempty"`
	CreatedAt                 time.Time      `json:"created_at"`
	UpdatedAt                 time.Time      `json:"updated_at"`
	DeletedAt                 gorm.DeletedAt `json:"-" gorm:"index"`

	// Unique constraint on employee_id + leave_type_code + year
	// This is handled at application level or via unique index

	// Relationships
	LeaveType *LeaveType `json:"leave_type,omitempty" gorm:"foreignKey:LeaveTypeCode;references:Code"`
}

// TableName specifies the table name
func (LeaveBalance) TableName() string {
	return "leave_balances"
}

// Ensure unique constraint on employee_id + leave_type_code + year
// This should be added via migration or unique index

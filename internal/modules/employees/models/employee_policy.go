package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeePolicy represents leave and attendance policy assignments
type EmployeePolicy struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	EmployeeID       *uint          `json:"employee_id,omitempty" gorm:"index"` // Nullable for drafts
	EmployeeIDString *string        `json:"employee_id_string,omitempty" gorm:"size:50;index"` // Employee ID string (e.g., "EMP466088") - for tracking during draft
	DraftID          *uint          `json:"draft_id,omitempty" gorm:"index"`   // References employee_onboarding_drafts(id)
	LeavePolicyID    *uint          `json:"leave_policy_id,omitempty" gorm:"index"` // References leave_policies(id) - to be created
	AttendancePolicyID *uint        `json:"attendance_policy_id,omitempty" gorm:"index"` // References attendance_policies(id) - to be created
	WeeklyOffDays    *string        `json:"weekly_off_days,omitempty" gorm:"size:50"` // e.g., "Saturday,Sunday" or "Friday"
	EffectiveDate    *time.Time     `json:"effective_date,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (EmployeePolicy) TableName() string {
	return "employee_policies"
}

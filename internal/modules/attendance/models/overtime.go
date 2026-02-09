package models

import (
	"time"

	"gorm.io/gorm"
)

// OvertimeRequestStatus represents the status of an overtime request
type OvertimeRequestStatus string

const (
	OTStatusPending   OvertimeRequestStatus = "pending"
	OTStatusApproved  OvertimeRequestStatus = "approved"
	OTStatusRejected  OvertimeRequestStatus = "rejected"
	OTStatusProcessed OvertimeRequestStatus = "processed" // Payout or comp-off applied
	OTStatusCancelled OvertimeRequestStatus = "cancelled"
)

// OvertimeType represents the type of overtime
type OvertimeType string

const (
	OTTypeWeekday OvertimeType = "weekday"
	OTTypeWeekend OvertimeType = "weekend"
	OTTypeHoliday OvertimeType = "holiday"
)

// CompensationType represents how overtime is compensated
type CompensationType string

const (
	CompTypePayout  CompensationType = "payout"
	CompTypeCompOff CompensationType = "comp_off"
	CompTypeBoth    CompensationType = "both"
)

// OvertimePolicy defines overtime rules for a department or role
type OvertimePolicy struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	Name         string  `json:"name" gorm:"size:255;not null"`
	DepartmentID *uint   `json:"department_id,omitempty" gorm:"index"` // NULL = applies to all departments
	IsActive     bool    `json:"is_active" gorm:"default:true"`

	// Thresholds (hours beyond which OT kicks in)
	DailyThresholdHours  float64 `json:"daily_threshold_hours" gorm:"type:decimal(5,2);default:8"`  // e.g., 8 hours
	WeeklyThresholdHours float64 `json:"weekly_threshold_hours" gorm:"type:decimal(5,2);default:40"` // e.g., 40 hours

	// Multiplier rates
	WeekdayMultiplier float64 `json:"weekday_multiplier" gorm:"type:decimal(4,2);default:1.5"`  // e.g., 1.5x
	WeekendMultiplier float64 `json:"weekend_multiplier" gorm:"type:decimal(4,2);default:2.0"`  // e.g., 2.0x
	HolidayMultiplier float64 `json:"holiday_multiplier" gorm:"type:decimal(4,2);default:2.5"`  // e.g., 2.5x

	// Limits
	MaxDailyOTHours   float64 `json:"max_daily_ot_hours" gorm:"type:decimal(5,2);default:4"`    // Max 4 hrs/day
	MaxWeeklyOTHours  float64 `json:"max_weekly_ot_hours" gorm:"type:decimal(5,2);default:20"`  // Max 20 hrs/week
	MaxMonthlyOTHours float64 `json:"max_monthly_ot_hours" gorm:"type:decimal(5,2);default:60"` // Max 60 hrs/month

	// Settings
	RequiresPreApproval bool `json:"requires_pre_approval" gorm:"default:true"`
	AllowCompOff        bool `json:"allow_comp_off" gorm:"default:true"`
	CompOffRatio        float64 `json:"comp_off_ratio" gorm:"type:decimal(4,2);default:1.0"` // 1 OT hour = 1 comp-off hour
	CompOffExpiryDays   int     `json:"comp_off_expiry_days" gorm:"default:90"` // Comp-off expires after N days

	// Budget
	MonthlyBudgetCap *float64 `json:"monthly_budget_cap,omitempty" gorm:"type:decimal(12,2)"` // Budget limit per department/month

	// Audit
	CreatedBy *uint `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy *uint `json:"updated_by,omitempty" gorm:"index"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (OvertimePolicy) TableName() string {
	return "overtime_policies"
}

// OvertimeRequest represents an employee's overtime request
type OvertimeRequest struct {
	ID              uint                  `json:"id" gorm:"primaryKey"`
	EmployeeID      uint                  `json:"employee_id" gorm:"index;not null"`
	PolicyID        *uint                 `json:"policy_id,omitempty" gorm:"index"` // References overtime_policies(id)
	Date            time.Time             `json:"date" gorm:"type:date;not null;index"`
	OvertimeType    OvertimeType          `json:"overtime_type" gorm:"type:varchar(20);not null"` // weekday, weekend, holiday
	Hours           float64               `json:"hours" gorm:"type:decimal(5,2);not null"`
	Reason          string                `json:"reason" gorm:"type:text;not null"`
	CompensationType CompensationType     `json:"compensation_type" gorm:"type:varchar(20);not null"` // payout, comp_off, both
	Status          OvertimeRequestStatus `json:"status" gorm:"type:varchar(20);default:'pending';not null"`

	// Computed fields
	Multiplier    float64  `json:"multiplier" gorm:"type:decimal(4,2);default:1.0"`
	PayoutAmount  *float64 `json:"payout_amount,omitempty" gorm:"type:decimal(12,2)"` // Calculated payout
	CompOffHours  *float64 `json:"comp_off_hours,omitempty" gorm:"type:decimal(5,2)"` // Comp-off hours credited
	CompOffExpiry *time.Time `json:"comp_off_expiry,omitempty"` // When comp-off expires

	// Approval
	ApprovedBy      *uint   `json:"approved_by,omitempty" gorm:"index"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	RejectionReason *string `json:"rejection_reason,omitempty" gorm:"type:text"`

	// Processing
	ProcessedAt  *time.Time `json:"processed_at,omitempty"` // When payout/comp-off was applied
	PayrollRunID *uint      `json:"payroll_run_id,omitempty" gorm:"index"` // If linked to payroll

	// Auto-detected (from attendance)
	IsAutoDetected  bool     `json:"is_auto_detected" gorm:"default:false"` // Detected from attendance data
	AttendanceHours *float64 `json:"attendance_hours,omitempty" gorm:"type:decimal(5,2)"` // Actual hours from attendance
	ShiftHours      *float64 `json:"shift_hours,omitempty" gorm:"type:decimal(5,2)"` // Scheduled shift hours

	// Audit
	CreatedBy *uint `json:"created_by,omitempty" gorm:"index"`

	// Relationships
	Approvals []OvertimeApproval `json:"approvals,omitempty" gorm:"foreignKey:OvertimeRequestID"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (OvertimeRequest) TableName() string {
	return "overtime_requests"
}

// OvertimeApproval represents an approval action for an OT request
type OvertimeApproval struct {
	ID                uint    `json:"id" gorm:"primaryKey"`
	OvertimeRequestID uint    `json:"overtime_request_id" gorm:"index;not null"`
	ApproverID        uint    `json:"approver_id" gorm:"index;not null"`
	Level             int     `json:"level" gorm:"not null;default:1"` // 1=Manager, 2=HR
	Status            string  `json:"status" gorm:"type:varchar(20);default:'pending';not null"`
	Comments          *string `json:"comments,omitempty" gorm:"type:text"`
	ActionAt          *time.Time `json:"action_at,omitempty"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (OvertimeApproval) TableName() string {
	return "overtime_approvals"
}

// --- Request DTOs ---

// CreateOvertimePolicyRequest is for creating/updating an OT policy
type CreateOvertimePolicyRequest struct {
	Name                 string   `json:"name" binding:"required,min=2"`
	DepartmentID         *uint    `json:"department_id,omitempty"`
	DailyThresholdHours  float64  `json:"daily_threshold_hours" binding:"required,gt=0"`
	WeeklyThresholdHours float64  `json:"weekly_threshold_hours" binding:"required,gt=0"`
	WeekdayMultiplier    float64  `json:"weekday_multiplier" binding:"required,gt=0"`
	WeekendMultiplier    float64  `json:"weekend_multiplier" binding:"required,gt=0"`
	HolidayMultiplier    float64  `json:"holiday_multiplier" binding:"required,gt=0"`
	MaxDailyOTHours      float64  `json:"max_daily_ot_hours" binding:"required,gt=0"`
	MaxWeeklyOTHours     float64  `json:"max_weekly_ot_hours" binding:"required,gt=0"`
	MaxMonthlyOTHours    float64  `json:"max_monthly_ot_hours" binding:"required,gt=0"`
	RequiresPreApproval  bool     `json:"requires_pre_approval"`
	AllowCompOff         bool     `json:"allow_comp_off"`
	CompOffRatio         float64  `json:"comp_off_ratio"`
	CompOffExpiryDays    int      `json:"comp_off_expiry_days"`
	MonthlyBudgetCap     *float64 `json:"monthly_budget_cap,omitempty"`
}

// CreateOvertimeRequestDTO is for submitting an OT request
type CreateOvertimeRequestDTO struct {
	Date             string `json:"date" binding:"required"`                                      // YYYY-MM-DD
	Hours            float64 `json:"hours" binding:"required,gt=0,max=24"`
	OvertimeType     string `json:"overtime_type" binding:"required,oneof=weekday weekend holiday"`
	Reason           string `json:"reason" binding:"required,min=3"`
	CompensationType string `json:"compensation_type" binding:"required,oneof=payout comp_off both"`
}

// ApproveOvertimeRequest is for approving/rejecting an OT request
type ApproveOvertimeRequestDTO struct {
	Status   string  `json:"status" binding:"required,oneof=approved rejected"`
	Comments *string `json:"comments,omitempty"`
}

package models

import (
	"time"

	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

// TimesheetStatus represents the status of a timesheet submission
type TimesheetStatus string

const (
	TimesheetStatusDraft     TimesheetStatus = "draft"
	TimesheetStatusSubmitted TimesheetStatus = "submitted"
	TimesheetStatusApproved  TimesheetStatus = "approved"
	TimesheetStatusRejected  TimesheetStatus = "rejected"
	TimesheetStatusRecalled  TimesheetStatus = "recalled"
)

// TimesheetEntryType represents the type/category of a timesheet entry
type TimesheetEntryType string

const (
	EntryTypeRegular  TimesheetEntryType = "regular"
	EntryTypeMeeting  TimesheetEntryType = "meeting"
	EntryTypeTraining TimesheetEntryType = "training"
	EntryTypeTravel   TimesheetEntryType = "travel"
	EntryTypeSupport  TimesheetEntryType = "support"
	EntryTypeOvertime TimesheetEntryType = "overtime"
	EntryTypeBreak    TimesheetEntryType = "break"
	EntryTypeOther    TimesheetEntryType = "other"
)

// TimesheetEntry represents a single timesheet entry (hours worked on a project/task for a date)
type TimesheetEntry struct {
	ID              uint               `json:"id" gorm:"primaryKey"`
	TimesheetID     uint               `json:"timesheet_id" gorm:"index;not null"` // References timesheet_weeks(id)
	EmployeeID      uint               `json:"employee_id" gorm:"index;not null"`  // References employees(id)
	Date            time.Time          `json:"date" gorm:"type:date;not null;index"`
	ProjectID       *uint              `json:"project_id,omitempty" gorm:"index"`     // References projects (if you have them)
	ProjectName     string             `json:"project_name" gorm:"size:255;not null"` // Project name (freeform or from reference)
	ClientName      *string            `json:"client_name,omitempty" gorm:"size:255"`
	TaskDescription string             `json:"task_description" gorm:"type:text;not null"`
	EntryType       TimesheetEntryType `json:"entry_type" gorm:"type:varchar(30);default:'regular';not null"` // Type of work
	Hours           float64            `json:"hours" gorm:"type:decimal(5,2);default:0"`                      // 0.5h precision, optional (0 = not yet filled)
	IsBillable      bool               `json:"is_billable" gorm:"default:false"`
	Notes           *string            `json:"notes,omitempty" gorm:"type:text"`

	// Relationships
	Employee *employeeModels.Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID;references:ID"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (TimesheetEntry) TableName() string {
	return "timesheet_entries"
}

// TimesheetWeek represents a weekly timesheet submission by an employee
type TimesheetWeek struct {
	ID              uint            `json:"id" gorm:"primaryKey"`
	EmployeeID      uint            `json:"employee_id" gorm:"index;not null"`          // References employees(id)
	WeekStart       time.Time       `json:"week_start" gorm:"type:date;not null;index"` // Monday of the week
	WeekEnd         time.Time       `json:"week_end" gorm:"type:date;not null"`         // Sunday of the week
	Status          TimesheetStatus `json:"status" gorm:"type:varchar(20);default:'draft';not null"`
	TotalHours      float64         `json:"total_hours" gorm:"type:decimal(6,2);default:0"`
	BillableHours   float64         `json:"billable_hours" gorm:"type:decimal(6,2);default:0"`
	SubmittedAt     *time.Time      `json:"submitted_at,omitempty"`
	ApprovedAt      *time.Time      `json:"approved_at,omitempty"`
	ApprovedBy      *uint           `json:"approved_by,omitempty" gorm:"index"`
	RejectionReason *string         `json:"rejection_reason,omitempty" gorm:"type:text"`

	// Relationships
	Employee  *employeeModels.Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID;references:ID"`
	Entries   []TimesheetEntry         `json:"entries,omitempty" gorm:"foreignKey:TimesheetID"`
	Approvals []TimesheetApproval      `json:"approvals,omitempty" gorm:"foreignKey:TimesheetID"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (TimesheetWeek) TableName() string {
	return "timesheet_weeks"
}

// TimesheetApproval represents an approval action in the approval chain
type TimesheetApproval struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	TimesheetID uint       `json:"timesheet_id" gorm:"index;not null"`                        // References timesheet_weeks(id)
	ApproverID  uint       `json:"approver_id" gorm:"index;not null"`                         // References employees(id)
	Level       int        `json:"level" gorm:"not null;default:1"`                           // 1=Manager, 2=HR/Finance
	Status      string     `json:"status" gorm:"type:varchar(20);default:'pending';not null"` // pending, approved, rejected
	Comments    *string    `json:"comments,omitempty" gorm:"type:text"`
	ActionAt    *time.Time `json:"action_at,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"` // SLA due date

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (TimesheetApproval) TableName() string {
	return "timesheet_approvals"
}

// --- Request DTOs ---

// CreateTimesheetEntryRequest is the request body for creating a timesheet entry
type CreateTimesheetEntryRequest struct {
	Date            string   `json:"date" binding:"required"` // YYYY-MM-DD
	ProjectName     string   `json:"project_name" binding:"required,min=1"`
	ClientName      *string  `json:"client_name,omitempty"`
	TaskDescription string   `json:"task_description" binding:"required,min=1"`
	EntryType       string   `json:"entry_type" binding:"omitempty,oneof=regular meeting training travel support overtime break other"` // Type of work (default: regular)
	Hours           *float64 `json:"hours,omitempty" binding:"omitempty,min=0,max=24"`                                                  // Hours worked (optional, 0.5h precision)
	IsBillable      bool     `json:"is_billable"`
	Notes           *string  `json:"notes,omitempty"`
}

// BulkCreateTimesheetRequest is for creating multiple entries at once
type BulkCreateTimesheetRequest struct {
	WeekStart string                        `json:"week_start" binding:"required"` // YYYY-MM-DD (Monday)
	Entries   []CreateTimesheetEntryRequest `json:"entries" binding:"required,min=1"`
}

// SubmitTimesheetRequest is for submitting a weekly timesheet for approval
type SubmitTimesheetRequest struct {
	TimesheetID uint `json:"timesheet_id" binding:"required"`
}

// ApproveTimesheetRequest is for approving/rejecting a timesheet
type ApproveTimesheetRequest struct {
	Status   string  `json:"status" binding:"required,oneof=approved rejected"`
	Comments *string `json:"comments,omitempty"`
}

// CopyLastWeekRequest copies entries from the previous week
type CopyLastWeekRequest struct {
	SourceWeekStart string `json:"source_week_start" binding:"required"` // YYYY-MM-DD
	TargetWeekStart string `json:"target_week_start" binding:"required"` // YYYY-MM-DD
}

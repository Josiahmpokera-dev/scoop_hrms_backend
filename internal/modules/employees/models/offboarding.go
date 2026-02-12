package models

import (
	"time"

	"gorm.io/gorm"
)

// SeparationType represents the type of employee separation
type SeparationType string

const (
	SeparationTypeResignation     SeparationType = "Resignation"
	SeparationTypeTermination     SeparationType = "Termination"
	SeparationTypeRetirement      SeparationType = "Retirement"
	SeparationTypeMutualSeparation SeparationType = "Mutual Separation"
	SeparationTypeEndOfContract   SeparationType = "End of Contract"
	SeparationTypeLayoff          SeparationType = "Layoff"
)

// OffboardingStatus represents the status of an offboarding workflow
type OffboardingStatus string

const (
	OffboardingStatusInProgress OffboardingStatus = "In Progress"
	OffboardingStatusCompleted  OffboardingStatus = "Completed"
	OffboardingStatusOnHold      OffboardingStatus = "On Hold"
)

// ClearanceStatus represents the status of a clearance
type ClearanceStatus string

const (
	ClearanceStatusPending  ClearanceStatus = "Pending"
	ClearanceStatusInProgress ClearanceStatus = "In Progress"
	ClearanceStatusCleared  ClearanceStatus = "Cleared"
	ClearanceStatusIssues   ClearanceStatus = "Issues"
)

// ClearanceDepartment represents the department for clearance
type ClearanceDepartment string

const (
	ClearanceDepartmentManager ClearanceDepartment = "Manager"
	ClearanceDepartmentIT      ClearanceDepartment = "IT"
	ClearanceDepartmentHR      ClearanceDepartment = "HR"
	ClearanceDepartmentFinance ClearanceDepartment = "Finance"
	ClearanceDepartmentAssets  ClearanceDepartment = "Assets"
)

// ExitInterviewStatus represents the status of exit interview
type ExitInterviewStatus string

const (
	ExitInterviewStatusNotScheduled ExitInterviewStatus = "Not Scheduled"
	ExitInterviewStatusScheduled    ExitInterviewStatus = "Scheduled"
	ExitInterviewStatusCompleted    ExitInterviewStatus = "Completed"
	ExitInterviewStatusCancelled    ExitInterviewStatus = "Cancelled"
)

// SettlementStatus represents the status of final settlement
type SettlementStatus string

const (
	SettlementStatusPending    SettlementStatus = "Pending"
	SettlementStatusCalculated SettlementStatus = "Calculated"
	SettlementStatusApproved   SettlementStatus = "Approved"
	SettlementStatusPaid       SettlementStatus = "Paid"
)

// PaymentMethod represents payment method for settlement
type PaymentMethod string

const (
	PaymentMethodBankTransfer PaymentMethod = "Bank Transfer"
	PaymentMethodCheque       PaymentMethod = "Cheque"
	PaymentMethodCash         PaymentMethod = "Cash"
	PaymentMethodOther        PaymentMethod = "Other"
)

// ReasonCode represents predefined reason codes for separation
type ReasonCode string

const (
	ReasonCodeBetterOpportunity ReasonCode = "Better Opportunity"
	ReasonCodeHigherStudies     ReasonCode = "Higher Studies"
	ReasonCodeRelocation        ReasonCode = "Relocation"
	ReasonCodePersonalReasons   ReasonCode = "Personal Reasons"
	ReasonCodeHealthIssues      ReasonCode = "Health Issues"
	ReasonCodePerformanceIssues ReasonCode = "Performance Issues"
	ReasonCodePolicyViolation   ReasonCode = "Policy Violation"
	ReasonCodeOther             ReasonCode = "Other"
)

// OffboardingWorkflow represents an employee offboarding workflow
type OffboardingWorkflow struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	OffboardingID         string         `json:"offboarding_id" gorm:"uniqueIndex;not null;size:50"` // e.g., "off-001"
	EmployeeID            string         `json:"employee_id" gorm:"not null;size:50;index"` // Employee ID string (e.g., "EMP013")
	EmployeeDBID          *uint          `json:"employee_db_id,omitempty" gorm:"index"`      // References employees(id)
	
	// Separation Details
	SeparationType        string         `json:"separation_type" gorm:"not null;size:50"` // Resignation, Termination, etc.
	ResignationDate       *time.Time     `json:"resignation_date,omitempty"`              // Date when resignation/notice was submitted
	LastWorkingDate       *time.Time     `json:"last_working_date,omitempty"`             // Employee's last working day
	NoticePeriodDays      *int           `json:"notice_period_days,omitempty"`            // Standard notice period in days
	ServedNoticePeriod    *int           `json:"served_notice_period,omitempty"`          // Calculated served notice period
	BuyoutAmount          *float64       `json:"buyout_amount,omitempty" gorm:"type:decimal(10,2)"` // Notice period buyout amount
	Reason                *string        `json:"reason,omitempty" gorm:"type:text"`       // Detailed reason for separation
	ReasonCode            *string        `json:"reason_code,omitempty" gorm:"size:50"`    // Predefined reason code
	AdditionalNotes       *string        `json:"additional_notes,omitempty" gorm:"type:text"`
	
	// Exit Interview
	ExitInterviewRequired *bool          `json:"exit_interview_required,omitempty" gorm:"default:true"`
	ExitInterviewStatus   string         `json:"exit_interview_status" gorm:"size:50;default:'Not Scheduled'"` // Not Scheduled, Scheduled, Completed, Cancelled
	ExitInterviewDate     *time.Time     `json:"exit_interview_date,omitempty"`
	ExitInterviewLocation *string        `json:"exit_interview_location,omitempty" gorm:"size:255"`
	ExitInterviewerID     *string        `json:"exit_interviewer_id,omitempty" gorm:"size:50"` // Employee ID of interviewer
	ExitInterviewNotes    *string        `json:"exit_interview_notes,omitempty" gorm:"type:text"`
	ExitInterviewRating   *int           `json:"exit_interview_rating,omitempty"` // 1-5 rating
	WouldRecommend        *bool          `json:"would_recommend,omitempty"`
	ExitInterviewConductedByID *string    `json:"exit_interview_conducted_by_id,omitempty" gorm:"size:50"`
	ExitInterviewCompletedAt   *time.Time `json:"exit_interview_completed_at,omitempty"`
	
	// Clearance Status
	ClearanceStatus       string         `json:"clearance_status" gorm:"size:50;default:'Pending'"` // Pending, In Progress, Completed
	ClearanceProgress     *float64       `json:"clearance_progress,omitempty" gorm:"type:decimal(5,2);default:0"` // Percentage 0-100
	CompletedClearances   *int           `json:"completed_clearances,omitempty" gorm:"default:0"`
	TotalClearances       *int           `json:"total_clearances,omitempty" gorm:"default:5"`
	
	// Final Settlement
	FinalSettlementStatus string         `json:"final_settlement" gorm:"size:50;default:'Pending'"` // Pending, Calculated, Approved, Paid
	
	// Workflow Status
	Status                string         `json:"status" gorm:"size:50;default:'In Progress'"` // In Progress, Completed, On Hold
	
	// Access & Completion
	AccessRevokedDate     *time.Time     `json:"access_revoked_date,omitempty"`
	CompletedAt           *time.Time     `json:"completed_at,omitempty"`
	CompletedByID         *string         `json:"completed_by_id,omitempty" gorm:"size:50"`
	
	// Audit
	InitiatedByID         *string        `json:"initiated_by_id,omitempty" gorm:"size:50"` // Employee ID of person initiating
	CreatedBy             *uint          `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy             *uint          `json:"updated_by,omitempty" gorm:"index"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (OffboardingWorkflow) TableName() string {
	return "offboarding_workflows"
}

// OffboardingClearance represents a clearance task for an offboarding workflow
type OffboardingClearance struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	ClearanceID   string         `json:"clearance_id" gorm:"uniqueIndex;not null;size:50"` // e.g., "clear-001"
	OffboardingID string         `json:"offboarding_id" gorm:"not null;size:50;index"`      // References offboarding_workflows(offboarding_id)
	
	// Clearance Details
	Department    string         `json:"department" gorm:"not null;size:50"` // Manager, IT, HR, Finance, Assets
	Status        string         `json:"status" gorm:"size:50;default:'Pending'"` // Pending, Cleared, Issues
	ClearedByID   *string        `json:"cleared_by_id,omitempty" gorm:"size:50"` // Employee ID of person clearing
	ClearanceDate *time.Time     `json:"clearance_date,omitempty"`
	Notes         *string        `json:"notes,omitempty" gorm:"type:text"`
	Issues        *string        `json:"issues,omitempty" gorm:"type:text"` // Description of issues preventing clearance
	
	// Timestamps
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (OffboardingClearance) TableName() string {
	return "offboarding_clearances"
}

// OffboardingAssetReturn represents asset return tracking for offboarding
type OffboardingAssetReturn struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	OffboardingID string         `json:"offboarding_id" gorm:"not null;size:50;index"` // References offboarding_workflows(offboarding_id)
	AssetID       uint           `json:"asset_id" gorm:"not null;index"`               // References assets(id)
	
	// Return Details
	ReturnStatus  string         `json:"return_status" gorm:"size:50;default:'Pending'"` // Pending, Returned, Issue
	ReturnDate    *time.Time     `json:"return_date,omitempty"`
	ReturnCondition *string      `json:"return_condition,omitempty" gorm:"size:50"` // Excellent, Good, Fair, Poor, Damaged
	ReturnNotes   *string        `json:"return_notes,omitempty" gorm:"type:text"`
	
	// Issue Details (if return_status is "Issue")
	IssueType     *string        `json:"issue_type,omitempty" gorm:"size:50"` // Not Returned, Damaged, Missing, Stolen, Other
	IssueDescription *string     `json:"issue_description,omitempty" gorm:"type:text"`
	ReportedByID  *string        `json:"reported_by_id,omitempty" gorm:"size:50"`
	
	// Returned By
	ReturnedByID  *string        `json:"returned_by_id,omitempty" gorm:"size:50"` // Employee ID of person receiving asset
	
	// Expected Return
	ExpectedReturnDate *time.Time `json:"expected_return_date,omitempty"`
	
	// Timestamps
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (OffboardingAssetReturn) TableName() string {
	return "offboarding_asset_returns"
}

// FinalSettlement represents the final settlement for an offboarding workflow
type FinalSettlement struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	SettlementID  string         `json:"settlement_id" gorm:"uniqueIndex;not null;size:50"` // e.g., "sett-001"
	OffboardingID string         `json:"offboarding_id" gorm:"not null;size:50;index"`      // References offboarding_workflows(offboarding_id)
	EmployeeID    string         `json:"employee_id" gorm:"not null;size:50;index"`         // Employee ID string
	
	// Settlement Status
	Status        string         `json:"status" gorm:"size:50;default:'Pending'"` // Pending, Calculated, Approved, Paid
	
	// Calculation
	CalculatedDate *time.Time    `json:"calculated_date,omitempty"`
	CalculatedByID *string       `json:"calculated_by_id,omitempty" gorm:"size:50"`
	
	// Approval
	ApprovedDate  *time.Time     `json:"approved_date,omitempty"`
	ApprovedByID  *string       `json:"approved_by_id,omitempty" gorm:"size:50"`
	ApprovalNotes *string       `json:"approval_notes,omitempty" gorm:"type:text"`
	
	// Payment
	PaidDate      *time.Time     `json:"paid_date,omitempty"`
	PaidByID      *string        `json:"paid_by_id,omitempty" gorm:"size:50"`
	PaymentMethod *string        `json:"payment_method,omitempty" gorm:"size:50"` // Bank Transfer, Cheque, Cash, Other
	PaymentReference *string     `json:"payment_reference,omitempty" gorm:"size:255"`
	PaymentNotes  *string        `json:"payment_notes,omitempty" gorm:"type:text"`
	
	// Settlement Breakdown (stored as JSON)
	SettlementBreakdown string    `json:"settlement_breakdown,omitempty" gorm:"type:text"` // JSON string with earnings, deductions, net_settlement
	
	// Notes
	Notes         *string        `json:"notes,omitempty" gorm:"type:text"`
	
	// Timestamps
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (FinalSettlement) TableName() string {
	return "final_settlements"
}

package models

import (
	"time"
)

// InitiateSeparationRequest represents the request to initiate employee separation
type InitiateSeparationRequest struct {
	EmployeeID            string     `json:"employee_id" binding:"required"`
	SeparationType        string     `json:"separation_type" binding:"required,oneof=Resignation Termination Retirement 'Mutual Separation' 'End of Contract' Layoff"`
	ResignationDate       string     `json:"resignation_date" binding:"required"` // ISO 8601 format
	LastWorkingDate       string     `json:"last_working_date" binding:"required"` // ISO 8601 format
	NoticePeriodDays      int        `json:"notice_period_days" binding:"required,min=1"`
	Reason                string     `json:"reason" binding:"required"`
	ReasonCode            *string    `json:"reason_code,omitempty"` // Better Opportunity, Higher Studies, etc.
	ExitInterviewRequired *bool      `json:"exit_interview_required,omitempty"`
	RequireExecutiveApproval *bool   `json:"require_executive_approval,omitempty"` // true => CEO/MD/Director final approval required
	AdditionalNotes       *string    `json:"additional_notes,omitempty"`
	InitiatedBy           *string    `json:"initiated_by,omitempty"` // Employee ID of person initiating
}

// UpdateOffboardingWorkflowRequest represents the request to update an offboarding workflow
type UpdateOffboardingWorkflowRequest struct {
	LastWorkingDate *string `json:"last_working_date,omitempty"` // ISO 8601 format
	Reason          *string `json:"reason,omitempty"`
	AdditionalNotes *string `json:"additional_notes,omitempty"`
	Status          *string `json:"status,omitempty" binding:"omitempty,oneof='In Progress' Completed 'On Hold'"`
}

// UpdateClearanceRequest represents the request to update a clearance status
type UpdateClearanceRequest struct {
	Status      string  `json:"status" binding:"required,oneof=Pending Cleared Issues"`
	ClearedByID *string `json:"cleared_by_id,omitempty"` // Required if status is "Cleared"
	Notes       *string `json:"notes,omitempty"`
	Issues      *string `json:"issues,omitempty"` // Required if status is "Issues"
}

// RecordAssetReturnRequest represents the request to record asset return
type RecordAssetReturnRequest struct {
	ReturnDate     string  `json:"return_date" binding:"required"` // ISO 8601 format
	ReturnCondition string  `json:"return_condition" binding:"required,oneof=Excellent Good Fair Poor Damaged"`
	ReturnNotes    *string `json:"return_notes,omitempty"`
	ReturnedByID   string  `json:"returned_by_id" binding:"required"` // Employee ID
}

// RecordAssetIssueRequest represents the request to record asset issue
type RecordAssetIssueRequest struct {
	IssueType        string `json:"issue_type" binding:"required,oneof='Not Returned' Damaged Missing Stolen Other"`
	IssueDescription string `json:"issue_description" binding:"required"`
	ReportedByID     string `json:"reported_by_id" binding:"required"` // Employee ID
}

// ManageAssetClearanceRequest manages asset clearance status in a single API.
type ManageAssetClearanceRequest struct {
	Status            string  `json:"status" binding:"required,oneof=cleared issues pending"` // cleared, issues, pending
	ActedByEmployeeID *string `json:"acted_by_employee_id,omitempty"`                          // optional employee_id of actor
	ReturnDate        *string `json:"return_date,omitempty"`                                   // YYYY-MM-DD (for cleared)
	ReturnCondition   *string `json:"return_condition,omitempty"`                              // Excellent, Good, Fair, Poor, Damaged
	IssueType         *string `json:"issue_type,omitempty"`                                    // Not Returned, Damaged, Missing, Stolen, Other
	IssueDescription  *string `json:"issue_description,omitempty"`                             // required when status=issues
	Notes             *string `json:"notes,omitempty"`
}

// ScheduleExitInterviewRequest represents the request to schedule exit interview
type ScheduleExitInterviewRequest struct {
	InterviewDate string  `json:"interview_date" binding:"required"` // ISO 8601 format
	InterviewerID string  `json:"interviewer_id" binding:"required"` // Employee ID
	Location      *string `json:"location,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

// CompleteExitInterviewRequest represents the request to complete exit interview
type CompleteExitInterviewRequest struct {
	InterviewNotes string  `json:"interview_notes" binding:"required"`
	FeedbackRating *int     `json:"feedback_rating,omitempty" binding:"omitempty,min=1,max=5"`
	WouldRecommend *bool    `json:"would_recommend,omitempty"`
	ConductedByID  string   `json:"conducted_by_id" binding:"required"` // Employee ID
}

// CancelExitInterviewRequest represents the request to cancel exit interview
type CancelExitInterviewRequest struct {
	CancellationReason string `json:"cancellation_reason" binding:"required"`
	CancelledByID      string `json:"cancelled_by_id" binding:"required"` // Employee ID
}

// CalculateSettlementRequest represents the request to calculate final settlement
type CalculateSettlementRequest struct {
	// Optional: Can be calculated automatically based on employee data
	// Or can provide custom values
	OutstandingSalary *float64 `json:"outstanding_salary,omitempty"`
	LeaveEncashment   *float64 `json:"leave_encashment,omitempty"`
	Bonus             *float64 `json:"bonus,omitempty"`
	Incentives        *float64 `json:"incentives,omitempty"`
	OtherEarnings     *float64 `json:"other_earnings,omitempty"`
	OutstandingLoans  *float64 `json:"outstanding_loans,omitempty"`
	Advances          *float64 `json:"advances,omitempty"`
	AssetDamages      *float64 `json:"asset_damages,omitempty"`
	TaxDeductions     *float64 `json:"tax_deductions,omitempty"`
	OtherDeductions   *float64 `json:"other_deductions,omitempty"`
	Currency          *string  `json:"currency,omitempty"` // Default: TZS
	Notes             *string  `json:"notes,omitempty"`
}

// ApproveSettlementRequest represents the request to approve final settlement
type ApproveSettlementRequest struct {
	ApprovedByID   string  `json:"approved_by_id" binding:"required"` // Employee ID
	ApprovalNotes  *string `json:"approval_notes,omitempty"`
}

// PaySettlementRequest represents the request to mark settlement as paid
type PaySettlementRequest struct {
	PaidDate        string  `json:"paid_date" binding:"required"` // ISO 8601 format
	PaymentMethod   string  `json:"payment_method" binding:"required,oneof='Bank Transfer' Cheque Cash Other"`
	PaymentReference *string `json:"payment_reference,omitempty"`
	PaidByID        string  `json:"paid_by_id" binding:"required"` // Employee ID
	PaymentNotes    *string `json:"payment_notes,omitempty"`
}

// CompleteOffboardingRequest represents the request to complete offboarding
type CompleteOffboardingRequest struct {
	CompletedByID *string `json:"completed_by_id,omitempty"` // Employee ID
}

// ApproveOffboardingLevelOneRequest represents manager/hr level-one approval.
type ApproveOffboardingLevelOneRequest struct {
	Department *string `json:"department,omitempty"` // Optional when approver is admin; allowed: Manager, HR
	Notes      *string `json:"notes,omitempty"`
}

// ApproveOffboardingFinalRequest represents CEO/Director final approval.
type ApproveOffboardingFinalRequest struct {
	Notes *string `json:"notes,omitempty"`
}

// CancelOffboardingRequest represents cancellation/hold of offboarding workflow.
type CancelOffboardingRequest struct {
	Reason string  `json:"reason" binding:"required"`
	Notes  *string `json:"notes,omitempty"`
}

// ResumeOffboardingRequest represents request to resume an on-hold workflow.
type ResumeOffboardingRequest struct {
	Reason string  `json:"reason" binding:"required"`
	Notes  *string `json:"notes,omitempty"`
}

// SettlementBreakdown represents the breakdown of final settlement
type SettlementBreakdown struct {
	Earnings struct {
		OutstandingSalary float64 `json:"outstanding_salary"`
		LeaveEncashment   float64 `json:"leave_encashment"`
		Bonus             float64 `json:"bonus"`
		Incentives        float64 `json:"incentives"`
		OtherEarnings     float64 `json:"other_earnings"`
		TotalEarnings     float64 `json:"total_earnings"`
	} `json:"earnings"`
	Deductions struct {
		OutstandingLoans float64 `json:"outstanding_loans"`
		Advances         float64 `json:"advances"`
		AssetDamages     float64 `json:"asset_damages"`
		TaxDeductions    float64 `json:"tax_deductions"`
		OtherDeductions  float64 `json:"other_deductions"`
		TotalDeductions  float64 `json:"total_deductions"`
	} `json:"deductions"`
	NetSettlement float64 `json:"net_settlement"`
	Currency      string  `json:"currency"`
}

// ListOffboardingWorkflowsRequest represents the request to list offboarding workflows
type ListOffboardingWorkflowsRequest struct {
	Page            int     `form:"page" binding:"omitempty,min=1"`
	PageSize        int     `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status          *string `form:"status"`
	SeparationType  *string `form:"separation_type"`
	ClearanceStatus *string `form:"clearance_status"`
	FinalSettlement *string `form:"final_settlement"`
	Search          *string `form:"search"`
	Department      *string `form:"department"`
	DateFrom        *string `form:"date_from"` // ISO 8601 format
	DateTo          *string `form:"date_to"`   // ISO 8601 format
}

// OffboardingWorkflowResponse represents the response for an offboarding workflow
type OffboardingWorkflowResponse struct {
	OffboardingID         string                      `json:"offboarding_id"`
	EmployeeID            string                      `json:"employee_id"`
	EmployeeDBID          *uint                        `json:"employee_db_id,omitempty"`
	EmployeeName          *string                      `json:"employee_name,omitempty"`
	EmpID                 string                      `json:"emp_id"`
	PhotoURL              *string                      `json:"photo_url,omitempty"`
	Designation           *string                     `json:"designation,omitempty"`
	Department            *string                     `json:"department,omitempty"`
	SeparationType        string                      `json:"separation_type"`
	ResignationDate       *time.Time                  `json:"resignation_date,omitempty"`
	LastWorkingDate       *time.Time                  `json:"last_working_date,omitempty"`
	NoticePeriod          *int                        `json:"notice_period,omitempty"`
	ServedNoticePeriod    *int                        `json:"served_notice_period,omitempty"`
	BuyoutAmount          *float64                    `json:"buyout_amount,omitempty"`
	Reason                *string                     `json:"reason,omitempty"`
	ReasonCode            *string                     `json:"reason_code,omitempty"`
	ExitInterviewStatus   string                      `json:"exit_interview_status"`
	ExitInterviewDate     *time.Time                  `json:"exit_interview_date,omitempty"`
	ExitInterviewConductedBy *string                  `json:"exit_interview_conducted_by,omitempty"`
	ExitInterviewNotes    *string                     `json:"exit_interview_notes,omitempty"`
	ClearanceStatus       string                      `json:"clearance_status"`
	FinalSettlement       string                      `json:"final_settlement"`
	Status                string                      `json:"status"`
	ClearanceProgress     *float64                    `json:"clearance_progress,omitempty"`
	CompletedClearances   *int                        `json:"completed_clearances,omitempty"`
	TotalClearances       *int                        `json:"total_clearances,omitempty"`
	AccessRevokedDate     *time.Time                  `json:"access_revoked_date,omitempty"`
	Clearances            []OffboardingClearanceResponse `json:"clearances,omitempty"`
	Assets                []OffboardingAssetResponse    `json:"assets,omitempty"`
	FinalSettlementData   *FinalSettlementResponse      `json:"final_settlement_data,omitempty"`
	CreatedAt             time.Time                     `json:"created_at"`
	UpdatedAt             time.Time                     `json:"updated_at"`
}

// OffboardingClearanceResponse represents the response for a clearance
type OffboardingClearanceResponse struct {
	ID            string     `json:"id"`
	OffboardingID string     `json:"offboarding_id"`
	Department    string     `json:"department"`
	ClearedBy     *string    `json:"cleared_by,omitempty"`
	ClearedByID   *string    `json:"cleared_by_id,omitempty"`
	ClearanceDate *time.Time `json:"clearance_date,omitempty"`
	Status        string     `json:"status"`
	Notes         *string    `json:"notes,omitempty"`
	Issues        *string    `json:"issues,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// OffboardingAssetResponse represents the response for an asset in offboarding
type OffboardingAssetResponse struct {
	AssetID            uint       `json:"asset_id"`
	AssetCode          string     `json:"asset_code"`
	AssetType          string     `json:"asset_type"`
	AssetName          *string    `json:"asset_name,omitempty"`
	Brand              *string    `json:"brand,omitempty"`
	Model              *string    `json:"model,omitempty"`
	SerialNumber       *string    `json:"serial_number,omitempty"`
	AssetTag           *string    `json:"asset_tag,omitempty"`
	AssignedDate       *time.Time `json:"assigned_date,omitempty"`
	ReturnStatus       string     `json:"return_status"`
	ReturnDate         *time.Time `json:"return_date,omitempty"`
	ReturnCondition    *string   `json:"return_condition,omitempty"`
	ReturnNotes        *string   `json:"return_notes,omitempty"`
	ExpectedReturnDate *time.Time `json:"expected_return_date,omitempty"`
	IssueType          *string   `json:"issue_type,omitempty"`
	IssueDescription  *string   `json:"issue_description,omitempty"`
}

// FinalSettlementResponse represents the response for final settlement
type FinalSettlementResponse struct {
	SettlementID      string             `json:"settlement_id"`
	OffboardingID    string             `json:"offboarding_id"`
	EmployeeID       string             `json:"employee_id"`
	Status           string             `json:"status"`
	CalculatedDate   *time.Time         `json:"calculated_date,omitempty"`
	CalculatedBy     *string            `json:"calculated_by,omitempty"`
	ApprovedDate     *time.Time         `json:"approved_date,omitempty"`
	ApprovedBy       *string            `json:"approved_by,omitempty"`
	ApprovedByID     *string            `json:"approved_by_id,omitempty"`
	PaidDate         *time.Time         `json:"paid_date,omitempty"`
	PaidBy           *string            `json:"paid_by,omitempty"`
	PaidByID         *string            `json:"paid_by_id,omitempty"`
	TotalAmount      *float64           `json:"total_amount,omitempty"`
	Breakdown        *SettlementBreakdown `json:"breakdown,omitempty"`
	PaymentMethod    *string            `json:"payment_method,omitempty"`
	PaymentReference *string            `json:"payment_reference,omitempty"`
	Notes            *string            `json:"notes,omitempty"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

// OffboardingStatisticsResponse represents the response for offboarding statistics
type OffboardingStatisticsResponse struct {
	ActiveSeparations            int                    `json:"active_separations"`
	TotalClearances              int                    `json:"total_clearances"`
	CompletedClearances          int                    `json:"completed_clearances"`
	ClearanceCompletionPercentage float64               `json:"clearance_completion_percentage"`
	ScheduledExitInterviews      int                    `json:"scheduled_exit_interviews"`
	CompletedExitInterviews      int                    `json:"completed_exit_interviews"`
	PendingSettlements           int                    `json:"pending_settlements"`
	SettlementsByStatus          map[string]int        `json:"settlements_by_status"`
	SeparationsByType            map[string]int        `json:"separations_by_type"`
	WorkflowsByStatus            map[string]int        `json:"workflows_by_status"`
}

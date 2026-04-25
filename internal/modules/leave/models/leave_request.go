package models

import (
	"time"

	"gorm.io/gorm"
)

// LeaveRequestStatus represents the status of a leave request
type LeaveRequestStatus string

const (
	LeaveRequestStatusDraft           LeaveRequestStatus = "draft"
	LeaveRequestStatusPending         LeaveRequestStatus = "pending"
	LeaveRequestStatusReturnedForInfo LeaveRequestStatus = "returned_for_info"
	LeaveRequestStatusApproved        LeaveRequestStatus = "approved"
	LeaveRequestStatusPartiallyApproved LeaveRequestStatus = "partially_approved"
	LeaveRequestStatusRejected        LeaveRequestStatus = "rejected"
	LeaveRequestStatusCancelled       LeaveRequestStatus = "cancelled"
)

// LeaveRequest represents a leave application/request
type LeaveRequest struct {
	ID                        uint           `json:"id" gorm:"primaryKey"`
	ApplicationNumber         string         `json:"application_number" gorm:"uniqueIndex;not null;size:50"` // LV-2026-00123
	DocumentNumber            string         `json:"document_number" gorm:"default:'HR.FO.04.00';size:50"`
	EmployeeID                string         `json:"employee_id" gorm:"not null;size:50;index"` // References employees(employee_id)
	LeaveTypeCode             string         `json:"leave_type_code" gorm:"not null;size:10;index"` // References leave_types(code)
	FromDate                  time.Time      `json:"from_date" gorm:"not null"`
	ToDate                    time.Time      `json:"to_date" gorm:"not null"`
	TotalDays                 float64        `json:"total_days" gorm:"type:decimal(5,2);not null"`
	HalfDay                   bool           `json:"half_day" gorm:"default:false"`
	Reason                    string         `json:"reason" gorm:"type:text"`
	WorkDelegatedTo           *string        `json:"work_delegated_to,omitempty" gorm:"size:50"` // Employee ID
	HandoverNotes             *string        `json:"handover_notes,omitempty" gorm:"type:text"`
	ApplicationDate           time.Time      `json:"application_date" gorm:"not null"`
	LeavePeriodYear           int            `json:"leave_period_year" gorm:"not null"`
	ReportingBackDate         time.Time      `json:"reporting_back_date" gorm:"not null"`
	EmployeeContactNumber     string         `json:"employee_contact_number" gorm:"not null;size:50"`
	EmergencyContactPerson    string         `json:"emergency_contact_person" gorm:"not null;size:200"`
	EmergencyContactNumber    string         `json:"emergency_contact_number" gorm:"not null;size:50"`
	EmployeeSignature         *string        `json:"employee_signature,omitempty" gorm:"type:text"` // Base64 encoded signature
	Status                    string         `json:"status" gorm:"not null;size:50;default:'draft';index"` // draft, pending, approved, rejected, cancelled, returned_for_info
	InformationRequired       *string        `json:"information_required,omitempty" gorm:"type:text"` // For returned_for_info status
	CancellationReason        *string        `json:"cancellation_reason,omitempty" gorm:"type:text"`
	ApprovedDays             *float64        `json:"approved_days,omitempty" gorm:"type:decimal(5,2)"` // For partial approval
	UnpaidDays               *float64        `json:"unpaid_days,omitempty" gorm:"type:decimal(5,2)"` // For partial approval
	WillResultInLOP          bool           `json:"will_result_in_lop" gorm:"default:false"` // Loss of Pay
	LOPDays                  float64        `json:"lop_days" gorm:"type:decimal(5,2);default:0"`
	BalanceAfterApproval     *float64        `json:"balance_after_approval,omitempty" gorm:"type:decimal(5,2)"`
	CreatedAt                 time.Time      `json:"created_at"`
	UpdatedAt                 time.Time      `json:"updated_at"`
	CreatedBy                 *uint          `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy                 *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt                 gorm.DeletedAt `json:"-" gorm:"index"`

	// Leave Tracker Information (stored as JSON or separate fields)
	PreviousDaysUsed              float64 `json:"previous_days_used" gorm:"type:decimal(5,2);default:0"`
	DaysRemainingAfterRequest     float64 `json:"days_remaining_after_request" gorm:"type:decimal(5,2);default:0"`
	LeaveTakenFromPreviousYear    float64 `json:"leave_taken_from_previous_year" gorm:"type:decimal(5,2);default:0"`
	LeaveBalanceFromPreviousYear  float64 `json:"leave_balance_from_previous_year" gorm:"type:decimal(5,2);default:0"`

	// Relationships
	LeaveType  *LeaveType       `json:"leave_type,omitempty" gorm:"foreignKey:LeaveTypeCode;references:Code"`
	Documents  []LeaveDocument  `json:"documents,omitempty" gorm:"foreignKey:LeaveRequestID"`
	Approvals  []LeaveApproval  `json:"approvals,omitempty" gorm:"foreignKey:LeaveRequestID"`
}

// TableName specifies the table name
func (LeaveRequest) TableName() string {
	return "leave_requests"
}

// LeaveApproval represents approval workflow for leave requests
type LeaveApproval struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	LeaveRequestID  uint           `json:"leave_request_id" gorm:"not null;index"`
	Level           int            `json:"level" gorm:"not null"` // 1, 2, 3 (Head of Dept, HR, Director/CEO)
	ApproverType    string         `json:"approver_type" gorm:"not null;size:50"` // head_of_department, hr_department, director_ceo
	ApproverID      *string        `json:"approver_id,omitempty" gorm:"size:50"` // Employee ID
	ApproverName    *string        `json:"approver_name,omitempty" gorm:"size:200"`
	ApproverSignature *string      `json:"approver_signature,omitempty" gorm:"type:text"` // Base64 encoded
	Status          string         `json:"status" gorm:"not null;size:50;default:'pending'"` // pending, approved, rejected
	Remarks         *string        `json:"remarks,omitempty" gorm:"type:text"`
	ApprovedAt      *time.Time     `json:"approved_at,omitempty"`
	IsApplicable    bool           `json:"is_applicable" gorm:"default:true"` // Some levels may not be applicable
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	LeaveRequest *LeaveRequest `json:"leave_request,omitempty" gorm:"foreignKey:LeaveRequestID"`
}

// TableName specifies the table name
func (LeaveApproval) TableName() string {
	return "leave_approvals"
}

// LeaveDocument represents documents attached to leave requests
type LeaveDocument struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	LeaveRequestID uint          `json:"leave_request_id" gorm:"not null;index"`
	FileName      string         `json:"file_name" gorm:"not null;size:255"`
	FileURL       string         `json:"file_url" gorm:"not null;size:500"`
	FileType      string         `json:"file_type" gorm:"size:100"`
	FileSize      int64          `json:"file_size"` // in bytes
	DocumentType  *string        `json:"document_type,omitempty" gorm:"size:100"` // medical_certificate, travel_document, etc.
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	LeaveRequest *LeaveRequest `json:"leave_request,omitempty" gorm:"foreignKey:LeaveRequestID"`
}

// TableName specifies the table name
func (LeaveDocument) TableName() string {
	return "leave_documents"
}

// CreateLeaveRequestRequest represents the request to create a leave request
type CreateLeaveRequestRequest struct {
	LeaveTypeCode                string     `json:"leave_type_code" binding:"required"`
	FromDate                    string     `json:"from_date" binding:"required"`
	ToDate                      string     `json:"to_date" binding:"required"`
	HalfDay                     *bool      `json:"half_day"`
	Reason                      string     `json:"reason" binding:"required"`
	WorkDelegatedTo             *string    `json:"work_delegated_to"`
	HandoverNotes               *string    `json:"handover_notes"`
	SaveAsDraft                 *bool      `json:"save_as_draft"`
	ApplicationDate             string     `json:"application_date" binding:"required"`
	LeavePeriodYear             int        `json:"leave_period_year" binding:"required"`
	ReportingBackDate           string     `json:"reporting_back_date" binding:"required"`
	EmployeeContactNumber       string     `json:"employee_contact_number" binding:"required"`
	EmergencyContactPerson      string     `json:"emergency_contact_person" binding:"required"`
	EmergencyContactNumber      string     `json:"emergency_contact_number" binding:"required"`
	EmployeeSignature           *string    `json:"employee_signature"`
	Documents                   []LeaveDocumentInput `json:"documents,omitempty"`
	LeaveTracker                LeaveTrackerInput     `json:"leave_tracker" binding:"required"`
}

// LeaveDocumentInput represents document input for leave request
type LeaveDocumentInput struct {
	FileName     string  `json:"file_name" binding:"required"`
	FileURL      string  `json:"file_url" binding:"required"`
	FileType     string  `json:"file_type" binding:"required"`
	FileSize     int64   `json:"file_size"`
	DocumentType *string `json:"document_type"`
}

// LeaveTrackerInput represents leave tracker information
type LeaveTrackerInput struct {
	StartDate                   string  `json:"start_date" binding:"required"`
	EndDate                     string  `json:"end_date" binding:"required"`
	NumberOfDays                float64 `json:"number_of_days" binding:"required"`
	PreviousDaysUsed            float64 `json:"previous_days_used"`
	DaysRemainingAfterRequest   float64 `json:"days_remaining_after_request"`
	LeaveTakenFromPreviousYear  float64 `json:"leave_taken_from_previous_year"`
	LeaveBalanceFromPreviousYear float64 `json:"leave_balance_from_previous_year"`
}

// UpdateLeaveRequestRequest represents the request to update a leave request
type UpdateLeaveRequestRequest struct {
	FromDate                    *string   `json:"from_date"`
	ToDate                      *string   `json:"to_date"`
	HalfDay                     *bool      `json:"half_day"`
	Reason                      *string   `json:"reason"`
	WorkDelegatedTo             *string   `json:"work_delegated_to"`
	HandoverNotes               *string   `json:"handover_notes"`
	ReportingBackDate           *string   `json:"reporting_back_date"`
	EmployeeContactNumber       *string   `json:"employee_contact_number"`
	EmergencyContactPerson      *string   `json:"emergency_contact_person"`
	EmergencyContactNumber      *string   `json:"emergency_contact_number"`
	LeaveTracker                *LeaveTrackerInput `json:"leave_tracker"`
}

// ApproveLeaveRequestRequest represents the request to approve a leave request
type ApproveLeaveRequestRequest struct {
	Action          string  `json:"action" binding:"required"` // "approve"
	ApproverType    string  `json:"approver_type" binding:"required"`
	ApproverName    string  `json:"approver_name" binding:"required"`
	ApproverSignature *string `json:"approver_signature" binding:"omitempty"`
	Remarks         *string `json:"remarks"`
	PartialApproval *bool   `json:"partial_approval"`
	PartialDays     *float64 `json:"partial_days"`
	ConvertToUnpaid *bool   `json:"convert_to_unpaid"`
	NotifyEmployee  *bool   `json:"notify_employee"`
}

// RejectLeaveRequestRequest represents the request to reject a leave request
type RejectLeaveRequestRequest struct {
	ApproverType     string  `json:"approver_type" binding:"required"`
	ApproverName     string  `json:"approver_name" binding:"required"`
	ApproverSignature *string `json:"approver_signature" binding:"omitempty"`
	RejectionReason  string  `json:"rejection_reason" binding:"required"`
	Remarks          *string `json:"remarks"`
	NotifyEmployee   *bool   `json:"notify_employee"`
}

// ReturnForInfoRequest represents the request to return a leave request for information
type ReturnForInfoRequest struct {
	InformationRequired string `json:"information_required" binding:"required"`
	NotifyEmployee      *bool  `json:"notify_employee"`
}

// CancelLeaveRequestRequest represents the request to cancel a leave request
type CancelLeaveRequestRequest struct {
	CancellationReason string `json:"cancellation_reason" binding:"required"`
	NotifyApprover     *bool  `json:"notify_approver"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type JobRequisition struct {
	ID              string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	RequisitionNo   string         `json:"requisitionNumber" gorm:"unique;not null"`
	JobTitle        string         `json:"jobTitle" gorm:"not null"`
	Department      string         `json:"department" gorm:"not null"`
	Location        string         `json:"location"`
	Grade           string         `json:"grade"`
	Headcount       int            `json:"headcount" gorm:"default:1"`
	EmploymentType  string         `json:"employmentType"`
	HiringManager   string         `json:"hiringManager"` // Could be a User ID in real implementation
	RequestedBy     string         `json:"requestedBy"`
	TargetStartDate time.Time      `json:"targetStartDate"`
	EstimatedBudget float64        `json:"estimatedBudget"`
	Currency        string         `json:"currency" gorm:"default:'TZS'"`
	Justification   string         `json:"justification"`
	Status          string         `json:"status" gorm:"default:'Pending Approval'"` // Draft, Pending Approval, Approved, Rejected, Closed
	ApprovalFlow    []ApprovalStep `json:"approvalFlow" gorm:"foreignKey:RequisitionID"`
	CreatedAt       time.Time      `json:"createdDate"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type ApprovalStep struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	RequisitionID string    `json:"requisitionId" gorm:"index"`
	Level         int       `json:"level"`
	Approver      string    `json:"approver"`
	Role          string    `json:"role"`
	Status        string    `json:"status"` // Pending, Approved, Rejected
	Comments      string    `json:"comments"`
	Timestamp     time.Time `json:"timestamp"`
}

type CreateRequisitionRequest struct {
	JobTitle        string  `json:"jobTitle" binding:"required"`
	Department      string  `json:"department" binding:"required"`
	Location        string  `json:"location"`
	Grade           string  `json:"grade"`
	Headcount       int     `json:"headcount" binding:"required"`
	EmploymentType  string  `json:"employmentType" binding:"required"`
	HiringManager   string  `json:"hiringManager" binding:"required"`
	TargetStartDate string  `json:"targetStartDate" binding:"required"` // Parse to time.Time
	EstimatedBudget float64 `json:"estimatedBudget"`
	Currency        string  `json:"currency"`
	Justification   string  `json:"justification" binding:"required"`
}

type UpdateRequisitionRequest struct {
	Headcount     *int    `json:"headcount"`
	Justification *string `json:"justification"`
	Status        *string `json:"status"`
}

type ApprovalRequest struct {
	Action   string `json:"action" binding:"required"` // Approve, Reject
	Comments string `json:"comments"`
}

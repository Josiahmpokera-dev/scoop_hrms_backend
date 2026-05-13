package models

import (
	"encoding/json"
	"strconv"
	"strings"
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
	Priority        string         `json:"priority" gorm:"size:50"`                  // e.g. Low, Medium, High
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
	RequisitionNo   string  `json:"requisitionNo"` // optional; server generates if empty
	JobTitle        string  `json:"jobTitle" binding:"required"`
	Department      string  `json:"department" binding:"required"`
	Location        string  `json:"location"`
	Grade           string  `json:"grade"`
	Headcount       int     `json:"headcount" binding:"required,gt=0"`
	EmploymentType  string  `json:"employmentType"` // optional; default full_time (accepts Permanent, Contract, etc.)
	HiringManager   string  `json:"hiringManager"`
	TargetStartDate string  `json:"targetStartDate" binding:"required"` // YYYY-MM-DD
	EstimatedBudget float64 `json:"estimatedBudget"`
	Currency        string  `json:"currency"`
	Justification   string  `json:"justification" binding:"required"`
	Priority        string  `json:"priority"`
	Status          string  `json:"status"` // optional; e.g. Draft
}

// UnmarshalJSON accepts camelCase/snake_case, requisitionNumber as alias of requisitionNo,
// and headcount as JSON number or string (clients often stringify numbers).
func (req *CreateRequisitionRequest) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	getStr := func(keys ...string) string {
		for _, k := range keys {
			v, ok := raw[k]
			if !ok || string(v) == "null" {
				continue
			}
			var s string
			if err := json.Unmarshal(v, &s); err == nil {
				return strings.TrimSpace(s)
			}
		}
		return ""
	}
	getIntFlexible := func(keys ...string) int {
		for _, k := range keys {
			v, ok := raw[k]
			if !ok || string(v) == "null" {
				continue
			}
			var n int
			if err := json.Unmarshal(v, &n); err == nil {
				return n
			}
			var f float64
			if err := json.Unmarshal(v, &f); err == nil {
				return int(f)
			}
			var s string
			if err := json.Unmarshal(v, &s); err == nil {
				i, err := strconv.Atoi(strings.TrimSpace(s))
				if err == nil {
					return i
				}
			}
		}
		return 0
	}
	getFloatFlexible := func(keys ...string) float64 {
		for _, k := range keys {
			v, ok := raw[k]
			if !ok || string(v) == "null" {
				continue
			}
			var f float64
			if err := json.Unmarshal(v, &f); err == nil {
				return f
			}
			var s string
			if err := json.Unmarshal(v, &s); err == nil {
				f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
				if err == nil {
					return f
				}
			}
		}
		return 0
	}

	req.RequisitionNo = getStr(
		"requisitionNo", "requisition_no",
		"requisitionNumber", "requisition_number",
	)
	req.JobTitle = getStr("jobTitle", "job_title")
	req.Department = getStr("department")
	req.Location = getStr("location")
	req.Grade = getStr("grade")
	req.Headcount = getIntFlexible("headcount", "head_count")
	req.EmploymentType = getStr("employmentType", "employment_type")
	req.HiringManager = getStr("hiringManager", "hiring_manager")
	req.TargetStartDate = getStr("targetStartDate", "target_start_date")
	req.Currency = getStr("currency")
	req.Justification = getStr("justification")
	req.Priority = getStr("priority")
	req.Status = getStr("status")
	req.EstimatedBudget = getFloatFlexible("estimatedBudget", "estimated_budget")
	return nil
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

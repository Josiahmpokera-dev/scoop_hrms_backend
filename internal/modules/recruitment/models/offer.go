package models

import (
	"time"

	"gorm.io/gorm"
)

type Offer struct {
	ID              string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ApplicationID   string         `json:"applicationId" gorm:"not null;index"`
	CandidateID     string         `json:"candidateId" gorm:"not null"`
	JobOpeningID    string         `json:"jobId" gorm:"not null"`
	JoiningDate     time.Time      `json:"joiningDate"`
	ProbationPeriod int            `json:"probationPeriod"` // Months
	AnnualCTC       float64        `json:"annualCTC"`
	Currency        string         `json:"currency" gorm:"default:'TZS'"`
	Components      []SalaryComp   `json:"components" gorm:"serializer:json"`
	Benefits        []string       `json:"benefits" gorm:"serializer:json"`
	ExpiryDate      time.Time      `json:"expiryDate"`
	Status          string         `json:"status" gorm:"default:'Draft'"` // Draft, Sent, Accepted, Rejected
	SentAt          *time.Time     `json:"sentAt"`
	AcceptedAt      *time.Time     `json:"acceptedAt"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type SalaryComp struct {
	Component string  `json:"component"`
	Amount    float64 `json:"amount"`
}

type CreateOfferRequest struct {
	CandidateID     string       `json:"candidateId" binding:"required"`
	JobID           string       `json:"jobId" binding:"required"`
	JoiningDate     string       `json:"joiningDate" binding:"required"`
	ProbationPeriod int          `json:"probationPeriod"`
	Compensation    Compensation `json:"compensation"`
	Benefits        []string     `json:"benefits"`
	ExpiryDate      string       `json:"expiryDate"`
}

type Compensation struct {
	AnnualCTC  float64      `json:"annualCTC"`
	Currency   string       `json:"currency"`
	Components []SalaryComp `json:"components"`
}

type OfferApprovalRequest struct {
	Action   string `json:"action" binding:"required"` // Approve, Reject
	Comments string `json:"comments"`
}

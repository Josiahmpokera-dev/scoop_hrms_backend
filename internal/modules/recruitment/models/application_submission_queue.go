package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	ApplicationQueueStatusPending    = "pending"
	ApplicationQueueStatusProcessing = "processing"
	ApplicationQueueStatusCompleted  = "completed"
	ApplicationQueueStatusFailed     = "failed"
)

// ApplicationSubmissionQueue stores public apply payloads for durable processing/retry.
type ApplicationSubmissionQueue struct {
	ID            string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	JobOpeningID  string         `json:"jobOpeningId" gorm:"index"`
	CandidateEmail string        `json:"candidateEmail" gorm:"index"`
	Payload       string         `json:"payload" gorm:"type:text"`
	Status        string         `json:"status" gorm:"size:32;index;default:'pending'"`
	RetryCount    int            `json:"retryCount" gorm:"default:0"`
	LastError     string         `json:"lastError" gorm:"type:text"`
	ApplicationID string         `json:"applicationId" gorm:"index"`
	ProcessedAt   *time.Time     `json:"processedAt"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}


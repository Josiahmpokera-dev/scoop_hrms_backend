package models

import (
	"time"

	"gorm.io/gorm"
)

type Candidate struct {
	ID           string           `json:"id" gorm:"primaryKey;type:varchar(36)"`
	FirstName    string           `json:"firstName" gorm:"not null"`
	LastName     string           `json:"lastName" gorm:"not null"`
	Email        string           `json:"email" gorm:"unique;not null"`
	Phone        string           `json:"phone"`
	ResumeURL    string           `json:"resumeUrl"`
	LinkedInURL  string           `json:"linkedInUrl"`
	PortfolioURL string           `json:"portfolioUrl"`
	CoverLetter  string           `json:"coverLetter"`
	Skills       []string         `json:"skills" gorm:"serializer:json"`
	Applications []JobApplication `json:"applications" gorm:"foreignKey:CandidateID"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt   `json:"-" gorm:"index"`
}

type JobApplication struct {
	ID             string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	CandidateID    string         `json:"candidateId" gorm:"not null;index"`
	JobOpeningID   string         `json:"jobOpeningId" gorm:"not null;index"`
	JobOpening     JobOpening     `json:"jobOpening" gorm:"foreignKey:JobOpeningID"`
	Stage          string         `json:"stage" gorm:"default:'Applied'"` // Applied, Screening, Interview, Offer, Hired, Rejected
	Score          float64        `json:"score"`
	Notes          string         `json:"notes"`
	Interviews     []Interview    `json:"interviews" gorm:"foreignKey:ApplicationID"`
	Offer          *Offer         `json:"offer,omitempty" gorm:"foreignKey:ApplicationID"`
	AppliedDate    time.Time      `json:"appliedDate"`
	LastStatusDate time.Time      `json:"lastStatusDate"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type ApplyJobRequest struct {
	JobID        string `json:"jobId" binding:"required"`
	FirstName    string `json:"firstName" binding:"required"`
	LastName     string `json:"lastName" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Phone        string `json:"phone"`
	Resume       string `json:"resume"` // Base64 or URL
	LinkedInURL  string `json:"linkedInUrl"`
	PortfolioURL string `json:"portfolioUrl"`
	CoverLetter  string `json:"coverLetter"`
}

type UpdateStageRequest struct {
	JobID string `json:"jobId"` // Optional verification
	Stage string `json:"stage" binding:"required"`
	Notes string `json:"notes"`
}

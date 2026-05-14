package models

import (
	"time"

	"gorm.io/gorm"
)

type TalentPoolCandidate struct {
	ID                  string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name                string         `json:"name" gorm:"not null"`
	Email               string         `json:"email" gorm:"unique;not null"`
	Skills              []string       `json:"skills" gorm:"serializer:json"`
	PreferredRole       []string       `json:"preferredRole" gorm:"serializer:json"`
	ResumeURL           string         `json:"resumeUrl"`
	ExperienceMin       int            `json:"experienceMin"` // Years
	Location            string         `json:"location"`
	SourceJobOpeningID  string         `json:"sourceJobOpeningId" gorm:"size:36;index"`  // last known job context
	SourceApplicationID string         `json:"sourceApplicationId" gorm:"size:36;index"` // optional link back
	InternalNotes       string         `json:"internalNotes" gorm:"type:text"`           // HR notes for outreach
	LastContactedAt     *time.Time     `json:"lastContactedAt"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}

type AddToTalentPoolRequest struct {
	Name          string   `json:"name" binding:"required"`
	Email         string   `json:"email" binding:"required,email"`
	Skills        []string `json:"skills"`
	PreferredRole []string `json:"preferredRole"`
	ResumeURL     string   `json:"resumeUrl"`
}

// MoveApplicationToTalentPoolRequest moves an application out of the active pipeline into TalentPool stage.
type MoveApplicationToTalentPoolRequest struct {
	Notes             string `json:"notes"`             // stored on JobApplication
	SyncTalentPoolRow *bool  `json:"syncTalentPoolRow"` // nil/true: upsert talent_pool_candidates for outreach; false: stage only
	InternalPoolNotes string `json:"internalPoolNotes"` // stored on talent pool row when syncing
}

// ContactTalentPoolRequest sends an email to the talent pool candidate (SMTP must be configured).
type ContactTalentPoolRequest struct {
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// TalentPoolPipelineApplicationBrief is one application currently in TalentPool stage.
type TalentPoolPipelineApplicationBrief struct {
	ApplicationID   string    `json:"applicationId"`
	JobOpeningID      string    `json:"jobOpeningId"`
	JobTitle          string    `json:"jobTitle"`
	Department        string    `json:"department,omitempty"`
	Location          string    `json:"location,omitempty"`
	AppliedDate       time.Time `json:"appliedDate"`
	LastStatusDate    time.Time `json:"lastStatusDate"`
	Notes             string    `json:"notes,omitempty"`
}

// TalentPoolPipelineCandidateSummary is one row for list/search/filter APIs (candidates with TalentPool applications).
type TalentPoolPipelineCandidateSummary struct {
	CandidateID                 string                               `json:"candidateId"`
	FirstName                   string                               `json:"firstName"`
	LastName                    string                               `json:"lastName"`
	FullName                    string                               `json:"fullName"`
	Email                       string                               `json:"email"`
	Phone                       string                               `json:"phone,omitempty"`
	ResumeURL                   string                               `json:"resumeUrl,omitempty"`
	Skills                      []string                             `json:"skills,omitempty"`
	LastTalentPoolApplicationAt time.Time                            `json:"lastTalentPoolApplicationAt"`
	TalentPoolApplications      []TalentPoolPipelineApplicationBrief `json:"talentPoolApplications"`
}

// TalentPoolPipelineCandidateDetail extends the summary for GET detail (full profile + pool registry if any).
type TalentPoolPipelineCandidateDetail struct {
	TalentPoolPipelineCandidateSummary
	LinkedInURL   string                 `json:"linkedInUrl,omitempty"`
	PortfolioURL  string                 `json:"portfolioUrl,omitempty"`
	CoverLetter   string                 `json:"coverLetter,omitempty"`
	RegistryEntry *TalentPoolCandidate   `json:"registryEntry,omitempty"` // talent_pool_candidates row matched by email, if any
}

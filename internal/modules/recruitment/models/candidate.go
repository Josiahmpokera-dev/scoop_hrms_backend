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
	Candidate      *Candidate     `json:"candidate,omitempty" gorm:"foreignKey:CandidateID"`
	JobOpeningID   string         `json:"jobOpeningId" gorm:"not null;index"`
	JobOpening     JobOpening     `json:"jobOpening" gorm:"foreignKey:JobOpeningID"`
	Stage          string         `json:"stage" gorm:"default:'Applied'"` // see recruitment_pipeline.go constants
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
	JobID             string `json:"jobId" binding:"required_without=ApplyToken"`
	ApplyToken        string `json:"applyToken" binding:"required_without=JobID"`
	FirstName         string `json:"firstName" binding:"required"`
	LastName          string `json:"lastName" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	Phone             string `json:"phone"`
	Resume            string `json:"resume"` // Base64 or URL
	ResumeName        string `json:"resumeName"`
	LinkedInURL       string `json:"linkedInUrl"`
	PortfolioURL      string `json:"portfolioUrl"`
	CoverLetter       string `json:"coverLetter"`
	ApplicationLetter string `json:"applicationLetter"` // optional; used when coverLetter is empty
}

type UpdateStageRequest struct {
	JobID string `json:"jobId"` // Optional verification
	Stage string `json:"stage" binding:"required"`
	Notes string `json:"notes"`
}

const (
	ApplicationActionAccept                     = "accept"
	ApplicationActionReject                     = "reject"
	ApplicationActionMoveToInterviewStageOne    = "move_to_interview_stage_one"
)

// ApplicationActionRequest provides explicit actions for application progression.
type ApplicationActionRequest struct {
	Action string `json:"action" binding:"required"` // accept | reject | move_to_interview_stage_one
	Notes  string `json:"notes"`

	// Used only when action=move_to_interview_stage_one
	InterviewType          string   `json:"interviewType"` // default: "Stage 1"
	ScheduledDate          string   `json:"scheduledDate"` // required for stage-one interview action
	ScheduledTime          string   `json:"scheduledTime"` // required for stage-one interview action
	Duration               int      `json:"duration"`
	Mode                   string   `json:"mode"`
	InterviewerEmployeeIDs []string `json:"interviewerEmployeeIds"`
	Interviewers           []string `json:"interviewers"`
}

// RecentApplicantView is a compact row for GET /applications/recent (latest applications first).
type RecentApplicantView struct {
	ApplicationID  string    `json:"applicationId"`
	CandidateID    string    `json:"candidateId"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	FullName       string    `json:"fullName"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone,omitempty"`
	JobOpeningID   string    `json:"jobOpeningId"`
	JobTitle       string    `json:"jobTitle"`
	Department     string    `json:"department,omitempty"`
	Location       string    `json:"location,omitempty"`
	AppliedAt      time.Time `json:"appliedAt"`
	Stage          string    `json:"stage"`
	Score          float64   `json:"score,omitempty"`
}

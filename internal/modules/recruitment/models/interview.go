package models

import (
	"time"

	"gorm.io/gorm"
)

type Interview struct {
	ID            string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ApplicationID string         `json:"applicationId" gorm:"not null;index"`
	CandidateID   string         `json:"candidateId" gorm:"not null"`
	JobOpeningID  string         `json:"jobId" gorm:"not null"`
	InterviewType string         `json:"interviewType" binding:"required"` // Technical, HR, etc
	Round         int            `json:"round"`
	ScheduledDate time.Time      `json:"scheduledDate"`
	ScheduledTime string         `json:"scheduledTime"` // Keep as string for now or combine with date
	DurationMin   int            `json:"duration"`
	Mode          string         `json:"mode"`         // Video Call, In-Person
	Interviewers  []string       `json:"interviewers" gorm:"serializer:json"` // Emails or IDs
	Status        string         `json:"status" gorm:"default:'Scheduled'"`   // Scheduled, Completed, Cancelled
	Feedback      *Feedback      `json:"feedback" gorm:"foreignKey:InterviewID"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type Feedback struct {
	ID               uint     `json:"id" gorm:"primaryKey"`
	InterviewID      string   `json:"interviewId" gorm:"uniqueIndex"`
	InterviewerEmail string   `json:"interviewerEmail"`
	Rating           int      `json:"rating"` // 1-5
	Strengths        []string `json:"strengths" gorm:"serializer:json"`
	Concerns         []string `json:"concerns" gorm:"serializer:json"`
	Recommendation   string   `json:"recommendation"` // Hire, No Hire
	Comments         string   `json:"comments"`
	CreatedAt        time.Time
}

type ScheduleInterviewRequest struct {
	CandidateID   string   `json:"candidateId" binding:"required"`
	JobID         string   `json:"jobId" binding:"required"`
	InterviewType string   `json:"interviewType" binding:"required"`
	Round         int      `json:"round"`
	ScheduledDate string   `json:"scheduledDate" binding:"required"`
	ScheduledTime string   `json:"scheduledTime" binding:"required"`
	Duration      int      `json:"duration"`
	Mode          string   `json:"mode"`
	Interviewers  []string `json:"interviewers"`
}

type SubmitFeedbackRequest struct {
	InterviewerEmail string   `json:"interviewerEmail" binding:"required"`
	Rating           int      `json:"rating" binding:"required"`
	Strengths        []string `json:"strengths"`
	Concerns         []string `json:"concerns"`
	Recommendation   string   `json:"recommendation" binding:"required"`
	Comments         string   `json:"comments"`
}

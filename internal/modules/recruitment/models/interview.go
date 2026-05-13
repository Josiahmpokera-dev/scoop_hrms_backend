package models

import (
	"time"

	"gorm.io/gorm"
)

// Interview rounds are capped at three stages (same candidate/application).
const InterviewMaxRound = 3

// Interview lifecycle statuses.
const (
	InterviewStatusScheduled               = "Scheduled"
	InterviewStatusPanelCompleted          = "PanelCompleted"
	InterviewStatusLegacyCompleted         = "Completed" // legacy: panel submitted feedback before HR gate existed
	InterviewStatusHrApprovedNextRound     = "HrApprovedNextRound"
	InterviewStatusHrRequestedReinterview = "HrRequestedReinterview"
	InterviewStatusHrRejected              = "HrRejected"
	InterviewStatusCancelled               = "Cancelled"
)

// HRInterviewDecision values for POST .../interviews/:id/hr-decision.
const (
	HrInterviewDecisionProceedNextRound    = "proceed_next_round"
	HrInterviewDecisionScheduleReinterview = "schedule_reinterview"
	HrInterviewDecisionRejectCandidate     = "reject_candidate"
)

// InterviewerAssignment stores employee-backed panel members at schedule time (snapshot).
type InterviewerAssignment struct {
	EmployeeID   string `json:"employeeId"`
	DisplayName  string `json:"displayName"`
	WorkEmail    string `json:"workEmail,omitempty"`
	DepartmentID *uint  `json:"departmentId,omitempty"`
	PositionID   *uint  `json:"positionId,omitempty"`
}

type Interview struct {
	ID            string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ApplicationID string         `json:"applicationId" gorm:"not null;index"`
	CandidateID   string         `json:"candidateId" gorm:"not null"`
	JobOpeningID  string         `json:"jobId" gorm:"not null"`
	InterviewType string         `json:"interviewType" binding:"required"` // Technical, HR, etc
	Round         int            `json:"round"`
	Attempt       int            `json:"attempt" gorm:"default:1;index"` // increases on same-round re-interview
	ScheduledDate time.Time      `json:"scheduledDate"`
	ScheduledTime string         `json:"scheduledTime"` // Keep as string for now or combine with date
	DurationMin   int            `json:"duration"`
	Mode          string         `json:"mode"` // Video Call, In-Person
	Interviewers  []string       `json:"interviewers" gorm:"serializer:json"`
	// InterviewerPanel is populated from interviewerEmployeeIds (preferred); Interviewers holds emails for notifications / legacy clients.
	InterviewerPanel []InterviewerAssignment `json:"interviewerPanel" gorm:"serializer:json"`
	Status           string                  `json:"status" gorm:"default:'Scheduled'"`

	HrDecision           string     `json:"hrDecision,omitempty"`           // proceed_next_round | schedule_reinterview | reject_candidate
	HrRemarks            string     `json:"hrRemarks,omitempty"`
	HrDecidedAt          *time.Time `json:"hrDecidedAt,omitempty"`
	HrDeciderEmployeeID  string     `json:"hrDeciderEmployeeId,omitempty"`
	PanelCompletedAt     *time.Time `json:"panelCompletedAt,omitempty"`

	Feedback  *Feedback      `json:"feedback" gorm:"foreignKey:InterviewID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
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
	CandidateID            string   `json:"candidateId" binding:"required"`
	JobID                  string   `json:"jobId" binding:"required"`
	InterviewType          string   `json:"interviewType" binding:"required"`
	Round                  int      `json:"round" binding:"required,min=1,max=3"`
	ScheduledDate          string   `json:"scheduledDate" binding:"required"`
	ScheduledTime          string   `json:"scheduledTime" binding:"required"`
	Duration               int      `json:"duration"`
	Mode                   string   `json:"mode"`
	Interviewers           []string `json:"interviewers"`
	InterviewerEmployeeIDs []string `json:"interviewerEmployeeIds"`
	IsReinterview          bool     `json:"isReinterview"`
}

// HRInterviewDecisionRequest records HR approval after panel feedback before the next round or re-interview.
type HRInterviewDecisionRequest struct {
	Decision            string `json:"decision" binding:"required"` // proceed_next_round | schedule_reinterview | reject_candidate
	Remarks             string `json:"remarks" binding:"required"`
	DecidedByEmployeeID string `json:"decidedByEmployeeId"` // optional override if JWT user has no employee row
}

type SubmitFeedbackRequest struct {
	InterviewerEmail string   `json:"interviewerEmail" binding:"required"`
	Rating           int      `json:"rating" binding:"required"`
	Strengths        []string `json:"strengths"`
	Concerns         []string `json:"concerns"`
	Recommendation   string   `json:"recommendation" binding:"required"`
	Comments         string   `json:"comments"`
}

// ManualScheduleInterviewRequest schedules an interview for someone not yet in the system:
// creates Candidate + JobApplication (stage Interview) + Interview in one flow.
type ManualScheduleInterviewRequest struct {
	JobOpeningID  string   `json:"jobOpeningId" binding:"required"`
	FirstName     string   `json:"firstName" binding:"required"`
	LastName      string   `json:"lastName" binding:"required"`
	Email         string   `json:"email" binding:"required,email"`
	Phone         string   `json:"phone"`
	ResumeURL     string   `json:"resumeUrl"`
	Notes         string   `json:"notes"` // optional HR notes on application
	InterviewType string   `json:"interviewType" binding:"required"`
	Round         int      `json:"round" binding:"required,min=1,max=3"`
	ScheduledDate string   `json:"scheduledDate" binding:"required"`
	ScheduledTime string   `json:"scheduledTime" binding:"required"`
	Duration      int      `json:"duration"`
	Mode          string   `json:"mode"`
	Interviewers  []string `json:"interviewers"`
	InterviewerEmployeeIDs []string `json:"interviewerEmployeeIds"`
	IsReinterview          bool     `json:"isReinterview"`
}

// BatchScheduleInterviewRequest creates interviews for existing applications (pipeline).
type BatchScheduleInterviewRequest struct {
	JobOpeningID   string   `json:"jobOpeningId" binding:"required"`
	ApplicationIDs []string `json:"applicationIds" binding:"required"`
	InterviewType  string   `json:"interviewType" binding:"required"`
	Round          int      `json:"round" binding:"required,min=1,max=3"`
	ScheduledDate  string   `json:"scheduledDate" binding:"required"` // YYYY-MM-DD
	ScheduledTime  string   `json:"scheduledTime" binding:"required"`
	Duration       int      `json:"duration"`
	Mode           string   `json:"mode"`
	Interviewers   []string `json:"interviewers"`
	InterviewerEmployeeIDs []string `json:"interviewerEmployeeIds"`
	IsReinterview          bool     `json:"isReinterview"`
	UpdateStage    *bool    `json:"updateStage"` // omit or true: set application stage to Interview; false skips
}

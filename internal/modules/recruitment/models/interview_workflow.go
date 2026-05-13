package models

import (
	"time"

	"gorm.io/gorm"
)

// --- Stage attempt lifecycle (UI cards) ---

const (
	IWAttemptStatusScheduled   = "Scheduled"
	IWAttemptStatusInProgress  = "InProgress"
	IWAttemptStatusCompleted   = "Completed"
	IWAttemptStatusMissed       = "Missed"
	IWAttemptStatusSuperseded  = "Superseded" // replaced by repeat or reschedule
)

// --- Attendance ---

const (
	IWAttendancePending = "pending"
	IWAttendancePresent = "present"
	IWAttendanceAbsent  = "absent"
	IWAttendanceExcused = "excused"
)

// --- Per-stage outcome (after evaluation) ---

const (
	IWStageOutcomeProceed     = "proceed"
	IWStageOutcomeRepeatStage = "repeat_stage"
	IWStageOutcomeFail        = "fail"
)

// --- Process aggregate status ---

const (
	IWProcessInProgress          = "in_progress"
	IWProcessPendingApprovals    = "pending_final_approvals"
	IWProcessApproved            = "approved"
	IWProcessRejected            = "rejected"
	IWProcessFailed              = "failed"
)

// --- Final approval row ---

const (
	IWApprovalPending   = "pending"
	IWApprovalApproved  = "approved"
	IWApprovalRejected  = "rejected"
	IWApprovalSkipped   = "skipped" // e.g. CEO not in chain
)

const (
	IWApprovalRoleHR      = "hr"
	IWApprovalRoleManager = "manager"
	IWApprovalRoleCEO     = "ceo"
)

// --- Audit timeline (append-only) ---

const (
	IWAuditProcessStarted        = "process_started"
	IWAuditStageScheduled        = "stage_scheduled"
	IWAuditAttendanceRecorded    = "attendance_recorded"
	IWAuditStageStarted          = "stage_started"
	IWAuditEvaluationSubmitted   = "evaluation_submitted"
	IWAuditAbsenceFail           = "absence_fail"
	IWAuditAbsenceReschedule     = "absence_reschedule"
	IWAuditFinalApprovalDecided  = "final_approval_decided"
)

// ScoringDimension defines structured rubric for one score axis (stored as JSON on definition stage / snapshot).
type ScoringDimension struct {
	Key      string  `json:"key" binding:"required"`
	Label    string  `json:"label" binding:"required"`
	Min      float64 `json:"min"`
	Max      float64 `json:"max" binding:"required"`
	Required bool    `json:"required"`
}

// InterviewWorkflowDefinition is a reusable template bound to one job opening (HR configures stages + rubric + CEO flag).
type InterviewWorkflowDefinition struct {
	ID               string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	JobOpeningID     string         `json:"jobOpeningId" gorm:"not null;index;type:varchar(36)"`
	Name             string         `json:"name" gorm:"not null;size:200"`
	IncludeCEOApproval bool         `json:"includeCeoApproval" gorm:"default:false"`
	Stages           []InterviewWorkflowDefinitionStage `json:"stages,omitempty" gorm:"foreignKey:DefinitionID;references:ID"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// InterviewWorkflowDefinitionStage is one ordered step in the definition (sequence starts at 1).
type InterviewWorkflowDefinitionStage struct {
	ID             string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	DefinitionID   string         `json:"definitionId" gorm:"not null;index;type:varchar(36)"`
	Sequence       int            `json:"sequence" gorm:"not null;index"`
	Title          string         `json:"title" gorm:"not null;size:200"`
	Description    string         `json:"description" gorm:"type:text"`
	// ScoringDimensions JSON array of ScoringDimension
	ScoringDimensionsJSON string `json:"scoringDimensionsJson" gorm:"column:scoring_dimensions_json;type:text"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// InterviewWorkflowProcessStageSnapshot freezes rubric/titles per running process (definition edits do not change in-flight processes).
type InterviewWorkflowProcessStageSnapshot struct {
	ID                    string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ProcessID             string         `json:"processId" gorm:"not null;index;type:varchar(36)"`
	Sequence              int            `json:"sequence" gorm:"not null;index"`
	Title                 string         `json:"title" gorm:"not null;size:200"`
	Description           string         `json:"description" gorm:"type:text"`
	ScoringDimensionsJSON string         `json:"scoringDimensionsJson" gorm:"column:scoring_dimensions_json;type:text"`
	CreatedAt             time.Time      `json:"createdAt"`
}

// InterviewWorkflowProcess is one candidate/application running through a definition + final approvals.
type InterviewWorkflowProcess struct {
	ID            string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ApplicationID string `json:"applicationId" gorm:"not null;uniqueIndex;type:varchar(36)"`
	DefinitionID  string `json:"definitionId" gorm:"not null;type:varchar(36)"`
	JobOpeningID  string `json:"jobOpeningId" gorm:"not null;index;type:varchar(36)"`
	Status        string `json:"status" gorm:"not null;default:'in_progress';index"`
	// IncludeCEOApproval is copied from the definition at process start (approvals chain).
	IncludeCEOApproval bool `json:"includeCeoApproval" gorm:"default:false"`

	// Frozen approvers (set at process start; CEO row only used when definition.IncludeCEOApproval)
	ApprovalHREmployeeID      string `json:"approvalHrEmployeeId" gorm:"not null;size:50"`
	ApprovalManagerEmployeeID string `json:"approvalManagerEmployeeId" gorm:"not null;size:50"`
	ApprovalCEOEmployeeID   string `json:"approvalCeoEmployeeId,omitempty" gorm:"size:50"`

	Snapshots []InterviewWorkflowProcessStageSnapshot `json:"snapshots,omitempty" gorm:"foreignKey:ProcessID;references:ID"`
	Attempts  []InterviewWorkflowStageAttempt         `json:"attempts,omitempty" gorm:"foreignKey:ProcessID;references:ID"`
	Approvals []InterviewWorkflowFinalApproval        `json:"approvals,omitempty" gorm:"foreignKey:ProcessID;references:ID"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// InterviewWorkflowStageAttempt is one schedulable unit (exclusive assignee). Repeats use a new row + higher AttemptNumber.
type InterviewWorkflowStageAttempt struct {
	ID                 string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ProcessID          string    `json:"processId" gorm:"not null;index;type:varchar(36)"`
	StageSequence      int       `json:"stageSequence" gorm:"not null;index"`
	AttemptNumber      int       `json:"attemptNumber" gorm:"not null;default:1;index"`
	AssignedEmployeeID string    `json:"assignedEmployeeId" gorm:"not null;size:50;index"`
	ScheduledDate      time.Time `json:"scheduledDate"`
	ScheduledTime      string    `json:"scheduledTime" gorm:"size:32"`
	DurationMin        int       `json:"durationMin"`
	Mode               string    `json:"mode" gorm:"size:64"`
	Status             string    `json:"status" gorm:"not null;default:'Scheduled';index"`

	AttendanceStatus    string     `json:"attendanceStatus" gorm:"not null;default:'pending';size:32"`
	AttendanceAt        *time.Time `json:"attendanceAt,omitempty"`
	AttendanceBy        string     `json:"attendanceByEmployeeId,omitempty" gorm:"size:50"`
	StartedAt           *time.Time `json:"startedAt,omitempty"`
	StartedByEmployeeID string     `json:"startedByEmployeeId,omitempty" gorm:"size:50"`

	// Immutable after evaluation submitted (service enforces)
	EvaluationScoresJSON string     `json:"evaluationScores,omitempty" gorm:"column:evaluation_scores_json;type:text"`
	EvaluationRemarks    string     `json:"evaluationRemarks,omitempty" gorm:"type:text"`
	StageOutcome         string     `json:"stageOutcome,omitempty" gorm:"size:32"` // proceed | repeat_stage | fail
	EvaluatedAt          *time.Time `json:"evaluatedAt,omitempty"`
	EvaluatedByEmployeeID string    `json:"evaluatedByEmployeeId,omitempty" gorm:"size:50"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// InterviewWorkflowFinalApproval is one step in the post-interview chain (sequential: HR then Manager then optional CEO).
type InterviewWorkflowFinalApproval struct {
	ID                 string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ProcessID          string     `json:"processId" gorm:"not null;index;type:varchar(36)"`
	StepOrder          int        `json:"stepOrder" gorm:"not null;index"` // 1,2,3
	RoleCode           string     `json:"roleCode" gorm:"not null;size:32"` // hr | manager | ceo
	AssignedEmployeeID string     `json:"assignedEmployeeId" gorm:"not null;size:50;index"`
	Status             string     `json:"status" gorm:"not null;default:'pending';index"`
	Remarks            string     `json:"remarks,omitempty" gorm:"type:text"`
	DecidedAt          *time.Time `json:"decidedAt,omitempty"`
	DecidedByEmployeeID string    `json:"decidedByEmployeeId,omitempty" gorm:"size:50"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// InterviewWorkflowAuditEvent is append-only timeline / audit trail.
type InterviewWorkflowAuditEvent struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	ProcessID          string         `json:"processId" gorm:"not null;index;type:varchar(36)"`
	AttemptID          string         `json:"attemptId,omitempty" gorm:"index;type:varchar(36)"`
	EventType          string         `json:"eventType" gorm:"not null;size:64;index"`
	ActorEmployeeID    string         `json:"actorEmployeeId,omitempty" gorm:"size:50"`
	Notes              string         `json:"notes,omitempty" gorm:"type:text"`
	PayloadJSON        string         `json:"payloadJson,omitempty" gorm:"column:payload_json;type:text"`
	CreatedAt          time.Time      `json:"createdAt"`
}

// --- Request DTOs ---

type CreateInterviewWorkflowDefinitionRequest struct {
	JobOpeningID         string                              `json:"jobOpeningId" binding:"required"`
	Name                 string                              `json:"name" binding:"required"`
	IncludeCEOApproval   bool                                `json:"includeCeoApproval"`
	Stages               []InterviewWorkflowDefinitionStageDTO `json:"stages" binding:"required,min=1,dive"`
}

type InterviewWorkflowDefinitionStageDTO struct {
	Sequence              int                `json:"sequence" binding:"required,min=1"`
	Title                 string             `json:"title" binding:"required"`
	Description           string             `json:"description"`
	ScoringDimensions     []ScoringDimension `json:"scoringDimensions" binding:"required,min=1,dive"`
}

type StartInterviewWorkflowProcessRequest struct {
	ApplicationID string `json:"applicationId" binding:"required"`
	DefinitionID  string `json:"definitionId" binding:"required"`
	FirstStage    struct {
		AssignedEmployeeID string `json:"assignedEmployeeId" binding:"required"`
		ScheduledDate      string `json:"scheduledDate" binding:"required"`
		ScheduledTime      string `json:"scheduledTime" binding:"required"`
		DurationMin        int    `json:"durationMin"`
		Mode               string `json:"mode"`
	} `json:"firstStage" binding:"required"`
	FinalApprovals struct {
		HREmployeeID      string `json:"hrEmployeeId" binding:"required"`
		ManagerEmployeeID string `json:"managerEmployeeId" binding:"required"`
		CEOEmployeeID     string `json:"ceoEmployeeId"` // required when definition has IncludeCEOApproval
	} `json:"finalApprovals" binding:"required"`
}

type ScheduleInterviewWorkflowStageRequest struct {
	StageSequence      int    `json:"stageSequence" binding:"required,min=1"`
	AssignedEmployeeID string `json:"assignedEmployeeId" binding:"required"`
	ScheduledDate      string `json:"scheduledDate" binding:"required"`
	ScheduledTime      string `json:"scheduledTime" binding:"required"`
	DurationMin        int    `json:"durationMin"`
	Mode               string `json:"mode"`
}

type RecordInterviewWorkflowAttendanceRequest struct {
	Status string `json:"status" binding:"required"` // present | absent | excused
	Notes  string `json:"notes"`
}

type InterviewWorkflowAbsenceActionRequest struct {
	Action string `json:"action" binding:"required"` // fail_candidate | reschedule
	// When action=reschedule
	NextScheduledDate string `json:"nextScheduledDate"`
	NextScheduledTime string `json:"nextScheduledTime"`
	DurationMin       int    `json:"durationMin"`
	Mode              string `json:"mode"`
	Notes             string `json:"notes" binding:"required"`
}

type SubmitInterviewWorkflowEvaluationRequest struct {
	Scores  map[string]float64 `json:"scores" binding:"required"`
	Remarks string             `json:"remarks" binding:"required"`
	Outcome string             `json:"outcome" binding:"required"` // proceed | repeat_stage | fail
	// When outcome=repeat_stage
	NextScheduledDate string `json:"nextScheduledDate"`
	NextScheduledTime string `json:"nextScheduledTime"`
	DurationMin       int    `json:"durationMin"`
	Mode              string `json:"mode"`
}

type InterviewWorkflowFinalApprovalDecisionRequest struct {
	Decision string `json:"decision" binding:"required"` // approve | reject
	Remarks  string `json:"remarks" binding:"required"`
}

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Goal - Performance goal (OKR, KPI, Project, Development)
type Goal struct {
	ID                   uint           `json:"id" gorm:"primaryKey"`
	GoalCode             string         `json:"goal_code" gorm:"uniqueIndex;not null;size:20"` // G-001, G-002...
	Title                string         `json:"title" gorm:"not null;size:300"`
	Description          *string        `json:"description,omitempty" gorm:"type:text"`
	Type                 string         `json:"type" gorm:"not null;size:20;index"`   // OKR, KPI, Project, Development
	Level                string         `json:"level" gorm:"not null;size:20;index"`   // Individual, Team, Department, Company
	OwnerID              uint           `json:"owner_id" gorm:"not null;index"`       // User ID
	OwnerName            string         `json:"owner_name" gorm:"size:200"`
	Department           *string        `json:"department,omitempty" gorm:"size:100;index"`
	ParentGoalID         *uint          `json:"parent_goal_id,omitempty" gorm:"index"`
	Weight               float64        `json:"weight" gorm:"type:decimal(5,2);default:100"`
	DueDate              *time.Time     `json:"due_date,omitempty"`
	Status               string         `json:"status" gorm:"type:varchar(30);default:'Not Started';index"`       // Not Started, In Progress, At Risk, Completed, Cancelled
	Progress             float64        `json:"progress" gorm:"type:decimal(5,2);default:0"`
	Visibility           string         `json:"visibility" gorm:"type:varchar(20);default:'Manager'"`            // Public, Manager, Private
	ApprovalStatus       string         `json:"approval_status" gorm:"type:varchar(30);default:'Draft'"`          // Draft, Pending Approval, Approved, Rejected
	ApprovedByID         *uint          `json:"approved_by_id,omitempty"`
	ApprovedByName       *string        `json:"approved_by_name,omitempty" gorm:"size:200"`
	ApprovedDate         *time.Time     `json:"approved_date,omitempty"`
	CompletionStatus     *string        `json:"completion_status,omitempty" gorm:"type:varchar(30)"`              // null, Pending Verification, Verified, Rejected
	CompletionEvidence   *string        `json:"completion_evidence,omitempty" gorm:"type:text"`
	CompletionEvidenceURL *string       `json:"completion_evidence_url,omitempty" gorm:"size:500"`
	FinalRating          *float64       `json:"final_rating,omitempty" gorm:"type:decimal(3,1)"`
	VerifiedByID         *uint          `json:"verified_by_id,omitempty"`
	VerifiedByName       *string        `json:"verified_by_name,omitempty" gorm:"size:200"`
	VerifiedDate         *time.Time     `json:"verified_date,omitempty"`
	Tags                 *string        `json:"tags,omitempty" gorm:"type:text"`       // JSON array stored as text
	LinkedProjectID      *uint          `json:"linked_project_id,omitempty" gorm:"index"`
	LinkedProjectCode    *string        `json:"linked_project_code,omitempty" gorm:"size:50"`
	LinkedProjectName    *string        `json:"linked_project_name,omitempty" gorm:"size:200"`
	AssignedToID         *uint          `json:"assigned_to_id,omitempty" gorm:"index"` // For manager-assigned goals
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`
	// Relations
	KeyResults []KeyResult   `json:"key_results,omitempty" gorm:"foreignKey:GoalID"`
	CheckIns   []GoalCheckIn `json:"check_ins,omitempty" gorm:"foreignKey:GoalID"`
}

func (Goal) TableName() string { return "performance_goals" }

// KeyResult - Measurable key result for an OKR goal
type KeyResult struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	KRCode       string         `json:"kr_code" gorm:"uniqueIndex;not null;size:20"` // KR-001
	GoalID       uint           `json:"goal_id" gorm:"not null;index"`
	Title        string         `json:"title" gorm:"not null;size:300"`
	TargetValue  float64        `json:"target_value" gorm:"type:decimal(15,2);not null"`
	CurrentValue float64        `json:"current_value" gorm:"type:decimal(15,2);default:0"`
	Unit         string         `json:"unit" gorm:"not null;size:50"`
	DueDate      *time.Time     `json:"due_date,omitempty"`
	Status       string         `json:"status" gorm:"type:varchar(20);default:'On Track'"` // On Track, At Risk, Behind, Completed
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (KeyResult) TableName() string { return "performance_key_results" }

// GoalCheckIn - Progress check-in for a goal
type GoalCheckIn struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	GoalID        uint           `json:"goal_id" gorm:"not null;index"`
	Progress      float64        `json:"progress" gorm:"type:decimal(5,2);not null"`
	Status        string         `json:"status" gorm:"type:varchar(20);not null"` // On Track, At Risk, Behind
	Notes         string         `json:"notes" gorm:"type:text;not null"`
	EvidenceURL   *string        `json:"evidence_url,omitempty" gorm:"size:500"`
	UpdatedByID   uint           `json:"updated_by_id" gorm:"not null"`
	UpdatedByName string         `json:"updated_by_name" gorm:"size:200"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (GoalCheckIn) TableName() string { return "performance_goal_check_ins" }

// --- Request DTOs ---

// CreateGoalRequest - Create a new goal
type CreateGoalRequest struct {
	Title        string     `json:"title" binding:"required,max=300"`
	Description  *string    `json:"description,omitempty" binding:"omitempty"`
	Type         string     `json:"type" binding:"required,oneof=OKR KPI Project Development"`
	Level        string     `json:"level" binding:"required,oneof=Individual Team Department Company"`
	OwnerID      uint       `json:"owner_id" binding:"required"`
	Department   *string    `json:"department,omitempty" binding:"omitempty,max=100"`
	ParentGoalID *uint      `json:"parent_goal_id,omitempty" binding:"omitempty"`
	Weight       float64    `json:"weight" binding:"omitempty"` // default 100
	DueDate      *time.Time `json:"due_date,omitempty" binding:"omitempty"`
	Visibility   string     `json:"visibility" binding:"omitempty,oneof=Public Manager Private"`
}

// UpdateGoalRequest - Update an existing goal
type UpdateGoalRequest struct {
	Title        *string    `json:"title,omitempty" binding:"omitempty,max=300"`
	Description  *string    `json:"description,omitempty" binding:"omitempty"`
	Type         *string    `json:"type,omitempty" binding:"omitempty,oneof=OKR KPI Project Development"`
	Level        *string    `json:"level,omitempty" binding:"omitempty,oneof=Individual Team Department Company"`
	Department   *string    `json:"department,omitempty" binding:"omitempty,max=100"`
	ParentGoalID *uint      `json:"parent_goal_id,omitempty" binding:"omitempty"`
	Weight       *float64   `json:"weight,omitempty" binding:"omitempty"`
	DueDate      *time.Time `json:"due_date,omitempty" binding:"omitempty"`
	Status       *string    `json:"status,omitempty" binding:"omitempty"`
	Progress     *float64   `json:"progress,omitempty" binding:"omitempty"`
	Visibility   *string    `json:"visibility,omitempty" binding:"omitempty,oneof=Public Manager Private"`
}

// CreateKeyResultRequest - Create a key result for a goal
type CreateKeyResultRequest struct {
	GoalID      uint       `json:"goal_id" binding:"required"`
	Title       string     `json:"title" binding:"required,max=300"`
	TargetValue float64    `json:"target_value" binding:"required"`
	Unit        string     `json:"unit" binding:"required,max=50"`
	DueDate     *time.Time `json:"due_date,omitempty" binding:"omitempty"`
}

// UpdateKeyResultRequest - Update a key result
type UpdateKeyResultRequest struct {
	Title        *string    `json:"title,omitempty" binding:"omitempty,max=300"`
	TargetValue  *float64   `json:"target_value,omitempty" binding:"omitempty"`
	CurrentValue *float64   `json:"current_value,omitempty" binding:"omitempty"`
	Unit         *string    `json:"unit,omitempty" binding:"omitempty,max=50"`
	DueDate      *time.Time `json:"due_date,omitempty" binding:"omitempty"`
	Status       *string    `json:"status,omitempty" binding:"omitempty,oneof=On Track At Risk Behind Completed"`
}

// CreateCheckInRequest - Create a goal check-in
type CreateCheckInRequest struct {
	GoalID        uint    `json:"goal_id" binding:"required"`
	Progress      float64 `json:"progress" binding:"required"`
	Status        string  `json:"status" binding:"required,oneof=On Track At Risk Behind"`
	Notes         string  `json:"notes" binding:"required"`
	EvidenceURL   *string `json:"evidence_url,omitempty" binding:"omitempty,max=500"`
	UpdatedByID   uint    `json:"updated_by_id" binding:"required"`
	UpdatedByName string  `json:"updated_by_name" binding:"required,max=200"`
}

// SubmitForApprovalRequest - Submit goal for approval
type SubmitForApprovalRequest struct {
	GoalID uint `json:"goal_id" binding:"required"`
}

// ApproveGoalRequest - Approve or reject a goal
type ApproveGoalRequest struct {
	GoalID   uint   `json:"goal_id" binding:"required"`
	Approved bool   `json:"approved"` // true = approve, false = reject
	Comments *string `json:"comments,omitempty" binding:"omitempty"`
}

// RequestCompletionRequest - Employee requests completion verification
type RequestCompletionRequest struct {
	GoalID              uint    `json:"goal_id" binding:"required"`
	CompletionEvidence  *string `json:"completion_evidence,omitempty" binding:"omitempty"`
	CompletionEvidenceURL *string `json:"completion_evidence_url,omitempty" binding:"omitempty,max=500"`
}

// VerifyCompletionRequest - Manager verifies or rejects completion
type VerifyCompletionRequest struct {
	GoalID     uint    `json:"goal_id" binding:"required"`
	Verified   bool    `json:"verified"` // true = verified, false = rejected
	FinalRating *float64 `json:"final_rating,omitempty" binding:"omitempty"`
	Comments   *string `json:"comments,omitempty" binding:"omitempty"`
}

// KeyResultInput - Key result payload for assign goal (no goal_id; server sets it)
type KeyResultInput struct {
	Title       string     `json:"title" binding:"required,max=300"`
	TargetValue float64    `json:"target_value" binding:"required"`
	Unit        string     `json:"unit" binding:"required,max=50"`
	DueDate     *time.Time `json:"due_date,omitempty" binding:"omitempty"`
}

// AssignGoalRequest - Assign goal to employee (includes key_results)
type AssignGoalRequest struct {
	GoalID      uint             `json:"goal_id" binding:"required"`
	AssignedTo  uint             `json:"assigned_to" binding:"required"`
	KeyResults  []KeyResultInput  `json:"key_results,omitempty" binding:"omitempty"`
}

// LinkProjectRequest - Link goal to a project
type LinkProjectRequest struct {
	GoalID             uint    `json:"goal_id" binding:"required"`
	ProjectID          uint    `json:"project_id" binding:"required"`
	ProjectCode        string  `json:"project_code" binding:"required,max=50"`
	ProjectName        string  `json:"project_name" binding:"required,max=200"`
	ContributionNotes  *string `json:"contribution_notes,omitempty" binding:"omitempty"`
}

// NextGoalCode returns the next goal code in sequence (e.g. G-001, G-002).
// Pass the next ID or the current max ID + 1 from the database.
func NextGoalCode(id uint) string {
	return fmt.Sprintf("G-%03d", id)
}

// NextKRCode returns the next key result code in sequence (e.g. KR-001, KR-002).
// Pass the next ID or the current max ID + 1 from the database.
func NextKRCode(id uint) string {
	return fmt.Sprintf("KR-%03d", id)
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// AppraisalCycle represents an appraisal/performance review cycle.
type AppraisalCycle struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null;size:20"`
	Name        string         `json:"name" gorm:"not null;size:200"`
	CycleType   string         `json:"cycle_type" gorm:"not null;size:50;index"` // Annual, Mid-Year, Quarterly
	Year        int            `json:"year" gorm:"not null;index"`
	StartDate   time.Time      `json:"start_date" gorm:"not null"`
	EndDate     time.Time      `json:"end_date" gorm:"not null"`
	Status      string         `json:"status" gorm:"type:varchar(30);default:'Draft';index"` // Draft, Active, Closed
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	WorkflowSteps []AppraisalWorkflowStep `json:"workflow_steps,omitempty" gorm:"foreignKey:AppraisalCycleID"`
}

func (AppraisalCycle) TableName() string { return "performance_appraisal_cycles" }

// AppraisalWorkflowStep defines a step in the appraisal workflow.
type AppraisalWorkflowStep struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	AppraisalCycleID uint           `json:"appraisal_cycle_id" gorm:"not null;index"`
	StepOrder        int            `json:"step_order" gorm:"not null"`
	Name             string         `json:"name" gorm:"not null;size:100"`
	Description      *string        `json:"description,omitempty" gorm:"type:text"`
	DueDate          *time.Time     `json:"due_date,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

func (AppraisalWorkflowStep) TableName() string { return "performance_appraisal_workflow_steps" }

// Appraisal represents a single employee's appraisal in a cycle.
type Appraisal struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Code             string         `json:"code" gorm:"uniqueIndex;not null;size:20"`
	AppraisalCycleID uint           `json:"appraisal_cycle_id" gorm:"not null;index"`
	EmployeeID       uint           `json:"employee_id" gorm:"not null;index"`
	Department       *string        `json:"department,omitempty" gorm:"size:100;index"`
	Status           string         `json:"status" gorm:"type:varchar(30);default:'Not Started';index"`
	CurrentStep      *int           `json:"current_step,omitempty"`
	SelfRating       *float64       `json:"self_rating,omitempty" gorm:"type:decimal(3,2)"`
	ManagerRating    *float64       `json:"manager_rating,omitempty" gorm:"type:decimal(3,2)"`
	OverallRating    *float64       `json:"overall_rating,omitempty" gorm:"type:decimal(3,2)"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Appraisal) TableName() string { return "performance_appraisals" }

// --- Request DTOs ---

// CreateAppraisalCycleRequest - Create an appraisal cycle
type CreateAppraisalCycleRequest struct {
	Name      string     `json:"name" binding:"required,max=200"`
	CycleType string     `json:"cycle_type" binding:"required,oneof=Annual Mid-Year Quarterly"`
	Year      int        `json:"year" binding:"required"`
	StartDate time.Time  `json:"start_date" binding:"required"`
	EndDate   time.Time  `json:"end_date" binding:"required"`
	Steps     []WorkflowStepInput `json:"steps,omitempty" binding:"omitempty"`
}

// WorkflowStepInput - Workflow step for create
type WorkflowStepInput struct {
	StepOrder   int        `json:"step_order" binding:"required"`
	Name        string     `json:"name" binding:"required,max=100"`
	Description *string    `json:"description,omitempty" binding:"omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty" binding:"omitempty"`
}

// UpdateAppraisalCycleRequest - Update an appraisal cycle
type UpdateAppraisalCycleRequest struct {
	Name      *string    `json:"name,omitempty" binding:"omitempty,max=200"`
	CycleType *string    `json:"cycle_type,omitempty" binding:"omitempty,oneof=Annual Mid-Year Quarterly"`
	Year      *int       `json:"year,omitempty" binding:"omitempty"`
	StartDate *time.Time `json:"start_date,omitempty" binding:"omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty" binding:"omitempty"`
}

// SubmitSelfReviewRequest - Submit self review
type SubmitSelfReviewRequest struct {
	AppraisalID uint     `json:"appraisal_id" binding:"required"`
	SelfRating  float64  `json:"self_rating" binding:"required"`
	Comments    *string  `json:"comments,omitempty" binding:"omitempty"`
}

// SubmitManagerReviewRequest - Submit manager review
type SubmitManagerReviewRequest struct {
	AppraisalID    uint    `json:"appraisal_id" binding:"required"`
	ManagerRating  float64 `json:"manager_rating" binding:"required"`
	Comments       *string `json:"comments,omitempty" binding:"omitempty"`
}

// SubmitCalibrationRequest - Submit calibration
type SubmitCalibrationRequest struct {
	AppraisalID    uint     `json:"appraisal_id" binding:"required"`
	OverallRating  float64  `json:"overall_rating" binding:"required"`
	Comments       *string  `json:"comments,omitempty" binding:"omitempty"`
}

// SendBackRequest - Send appraisal back to previous step
type SendBackRequest struct {
	AppraisalID uint    `json:"appraisal_id" binding:"required"`
	ToStep      string  `json:"to_step" binding:"required"`
	Comments    *string `json:"comments,omitempty" binding:"omitempty"`
}

// FinalizeRequest - Finalize appraisal with final rating and outcomes
type FinalizeRequest struct {
	AppraisalID    uint     `json:"appraisal_id" binding:"required"`
	OverallRating  float64  `json:"overall_rating" binding:"required"`
	Outcomes       *string  `json:"outcomes,omitempty" binding:"omitempty"` // JSON or text
}

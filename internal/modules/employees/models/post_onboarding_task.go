package models

import (
	"time"

	"gorm.io/gorm"
)

// PostOnboardingTaskStatus represents the status of a post-onboarding task
type PostOnboardingTaskStatus string

const (
	TaskStatusPending    PostOnboardingTaskStatus = "pending"
	TaskStatusInProgress PostOnboardingTaskStatus = "in_progress"
	TaskStatusCompleted  PostOnboardingTaskStatus = "completed"
	TaskStatusSkipped    PostOnboardingTaskStatus = "skipped"
)

// PostOnboardingTaskPriority represents the priority of a task
type PostOnboardingTaskPriority string

const (
	TaskPriorityLow      PostOnboardingTaskPriority = "low"
	TaskPriorityMedium   PostOnboardingTaskPriority = "medium"
	TaskPriorityHigh     PostOnboardingTaskPriority = "high"
	TaskPriorityCritical PostOnboardingTaskPriority = "critical"
)

// PostOnboardingTaskType represents the type/category of task
type PostOnboardingTaskType string

const (
	TaskTypeWelcomeEmail              PostOnboardingTaskType = "send_welcome_email"
	TaskTypeCreateEmailAccount        PostOnboardingTaskType = "create_email_account"
	TaskTypeAllocateLaptop            PostOnboardingTaskType = "allocate_laptop"
	TaskTypeDeskAllocate              PostOnboardingTaskType = "desk_allocate"
	TaskTypeIssueIDCard               PostOnboardingTaskType = "issue_id_card"
	TaskTypeTeamIntroduction          PostOnboardingTaskType = "team_introduction"
	TaskTypePolicyAcknowledge         PostOnboardingTaskType = "policy_acknowledge"
	TaskTypeInductionTraining         PostOnboardingTaskType = "induction_training"
	TaskTypeAssignBuddyMentor         PostOnboardingTaskType = "assign_buddy_mentor"
	TaskTypeSetupDevEnvironment       PostOnboardingTaskType = "setup_dev_environment"
	TaskTypeFirstWeekReview           PostOnboardingTaskType = "first_week_review"
	TaskTypeFirstMonthReview          PostOnboardingTaskType = "first_month_review"
	TaskTypeMidProbationReview        PostOnboardingTaskType = "mid_probation_review"
	TaskTypeProbationConfirmation     PostOnboardingTaskType = "probation_confirmation"
)

// PostOnboardingTask represents a post-onboarding task/checklist item
type PostOnboardingTask struct {
	ID                uint                        `json:"id" gorm:"primaryKey"`
	EmployeeID        string                      `json:"employee_id" gorm:"not null;size:50;index"` // Employee ID (e.g., EMP001)
	TaskType          string                      `json:"task_type" gorm:"not null;size:100;index"` // Task type/category
	Title             string                      `json:"title" gorm:"not null;size:255"`           // Task title
	Description       *string                     `json:"description,omitempty" gorm:"type:text"`   // Task description
	Status            string                      `json:"status" gorm:"default:'pending';size:50"`  // pending, in_progress, completed, skipped
	Priority          string                      `json:"priority" gorm:"default:'medium';size:50"`  // low, medium, high, critical
	AssignedTo        *string                     `json:"assigned_to,omitempty" gorm:"size:100"`    // Assigned to role/department (HR, IT, etc.)
	AssignedToUserID  *uint                       `json:"assigned_to_user_id,omitempty" gorm:"index"` // Specific user assigned to task
	CompletedByUserID *uint                       `json:"completed_by_user_id,omitempty" gorm:"index"` // User who completed the task
	CompletedAt       *time.Time                  `json:"completed_at,omitempty"`                    // When task was completed
	DueDate           *time.Time                  `json:"due_date,omitempty"`                       // Task due date
	Notes             *string                     `json:"notes,omitempty" gorm:"type:text"`         // Additional notes
	Metadata          *string                     `json:"metadata,omitempty" gorm:"type:text"`        // JSON metadata for task-specific data
	CreatedBy         *uint                       `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy         *uint                       `json:"updated_by,omitempty" gorm:"index"`
	CreatedAt         time.Time                   `json:"created_at"`
	UpdatedAt         time.Time                   `json:"updated_at"`
	DeletedAt         gorm.DeletedAt              `json:"-" gorm:"index"`

	// Relationships (will be loaded separately if needed)
	// Employee         *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID;references:EmployeeID"`
	// AssignedToUser   *User     `json:"assigned_to_user,omitempty" gorm:"foreignKey:AssignedToUserID"`
	// CompletedByUser  *User     `json:"completed_by_user,omitempty" gorm:"foreignKey:CompletedByUserID"`
}

// TableName specifies the table name
func (PostOnboardingTask) TableName() string {
	return "post_onboarding_tasks"
}

// IsCompleted checks if task is completed
func (t *PostOnboardingTask) IsCompleted() bool {
	return t.Status == string(TaskStatusCompleted)
}

// IsPending checks if task is pending
func (t *PostOnboardingTask) IsPending() bool {
	return t.Status == string(TaskStatusPending)
}

// CanBeCompleted checks if task can be marked as completed
func (t *PostOnboardingTask) CanBeCompleted() bool {
	return t.Status != string(TaskStatusCompleted) && t.Status != string(TaskStatusSkipped)
}

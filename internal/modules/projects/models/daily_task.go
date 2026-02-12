package models

import (
	"time"

	"gorm.io/gorm"
)

// TaskStatus represents the status of a daily task
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusBlocked    TaskStatus = "blocked"
)

// TaskCategory represents the category of work done
type TaskCategory string

const (
	TaskCategoryDevelopment   TaskCategory = "development"
	TaskCategoryTesting       TaskCategory = "testing"
	TaskCategoryDesign        TaskCategory = "design"
	TaskCategoryDocumentation TaskCategory = "documentation"
	TaskCategoryMeeting       TaskCategory = "meeting"
	TaskCategoryResearch      TaskCategory = "research"
	TaskCategoryPlanning      TaskCategory = "planning"
	TaskCategoryReview        TaskCategory = "review"
	TaskCategorySupport       TaskCategory = "support"
	TaskCategoryOther         TaskCategory = "other"
)

// DailyTask represents a daily task entry logged by an employee
type DailyTask struct {
	ID                   uint           `json:"id" gorm:"primaryKey"`
	ProjectID            uint           `json:"project_id" gorm:"not null;index"`
	ProjectName          string         `json:"project_name" gorm:"size:200"`
	EmployeeID           string         `json:"employee_id" gorm:"not null;size:50;index"`
	EmployeeName         string         `json:"employee_name" gorm:"size:200"`
	TaskDate             time.Time      `json:"task_date" gorm:"type:date;not null;index"`
	Title                string         `json:"title" gorm:"not null;size:300"`
	Description          *string        `json:"description,omitempty" gorm:"type:text"`
	HoursSpent           float64        `json:"hours_spent" gorm:"type:decimal(5,2);not null"`
	Status               TaskStatus     `json:"status" gorm:"type:varchar(20);default:'completed';index"`
	Category             TaskCategory   `json:"category" gorm:"type:varchar(30);default:'other';index"`
	ProgressContribution float64        `json:"progress_contribution" gorm:"type:decimal(5,2);default:0"`
	Tags                 *string        `json:"tags,omitempty" gorm:"size:500"`
	Notes                *string        `json:"notes,omitempty" gorm:"type:text"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for DailyTask
func (DailyTask) TableName() string {
	return "daily_tasks"
}

// ──────────────────────────── Request DTOs ────────────────────────────

// CreateDailyTaskRequest is the request body for creating a daily task
type CreateDailyTaskRequest struct {
	ProjectID            uint         `json:"project_id" binding:"required"`
	TaskDate             *string      `json:"task_date,omitempty"`   // YYYY-MM-DD, defaults to today
	Title                string       `json:"title" binding:"required,min=2,max=300"`
	Description          *string      `json:"description,omitempty"`
	HoursSpent           float64      `json:"hours_spent" binding:"required,gt=0,lte=24"`
	Status               *TaskStatus  `json:"status,omitempty"`
	Category             TaskCategory `json:"category" binding:"required"`
	ProgressContribution *float64     `json:"progress_contribution,omitempty"` // 0-100
	Tags                 *string      `json:"tags,omitempty"`
	Notes                *string      `json:"notes,omitempty"`
}

// UpdateDailyTaskRequest is the request body for updating a daily task
type UpdateDailyTaskRequest struct {
	Title                *string       `json:"title,omitempty" binding:"omitempty,min=2,max=300"`
	Description          *string       `json:"description,omitempty"`
	HoursSpent           *float64      `json:"hours_spent,omitempty" binding:"omitempty,gt=0,lte=24"`
	Status               *TaskStatus   `json:"status,omitempty"`
	Category             *TaskCategory `json:"category,omitempty"`
	ProgressContribution *float64      `json:"progress_contribution,omitempty"`
	Tags                 *string       `json:"tags,omitempty"`
	Notes                *string       `json:"notes,omitempty"`
}

// DailyTaskSummary represents a summary of daily tasks
type DailyTaskSummary struct {
	TotalTasks      int64              `json:"total_tasks"`
	TotalHours      float64            `json:"total_hours"`
	CompletedTasks  int64              `json:"completed_tasks"`
	PendingTasks    int64              `json:"pending_tasks"`
	InProgressTasks int64              `json:"in_progress_tasks"`
	BlockedTasks    int64              `json:"blocked_tasks"`
	ByCategory      []CategorySummary  `json:"by_category"`
	ByProject       []ProjectTaskSummary `json:"by_project"`
}

// CategorySummary summarizes tasks by category
type CategorySummary struct {
	Category   string  `json:"category"`
	TaskCount  int64   `json:"task_count"`
	TotalHours float64 `json:"total_hours"`
}

// ProjectTaskSummary summarizes tasks by project
type ProjectTaskSummary struct {
	ProjectID   uint    `json:"project_id"`
	ProjectName string  `json:"project_name"`
	TaskCount   int64   `json:"task_count"`
	TotalHours  float64 `json:"total_hours"`
}

package models

import "time"

// CreatePostOnboardingTaskRequest represents the request to create a post-onboarding task
type CreatePostOnboardingTaskRequest struct {
	EmployeeID       string     `json:"employee_id" binding:"required"` // Employee ID (e.g., EMP001)
	TaskType         string     `json:"task_type" binding:"required"`   // Task type
	Title            *string    `json:"title,omitempty"`                // Optional custom title (auto-generated if not provided)
	Description      *string    `json:"description,omitempty"`
	Priority         *string    `json:"priority,omitempty"` // low, medium, high, critical (default: medium)
	AssignedTo       *string    `json:"assigned_to,omitempty"` // HR, IT, etc.
	AssignedToUserID *uint      `json:"assigned_to_user_id,omitempty"` // Specific user ID
	DueDate          *string    `json:"due_date,omitempty"` // YYYY-MM-DD format
	Notes            *string    `json:"notes,omitempty"`
	Metadata         *string    `json:"metadata,omitempty"` // JSON string for task-specific data
}

// UpdatePostOnboardingTaskRequest represents the request to update a post-onboarding task
type UpdatePostOnboardingTaskRequest struct {
	ID               uint       `json:"id" binding:"required"`
	Title            *string    `json:"title,omitempty"`
	Description      *string    `json:"description,omitempty"`
	Status           *string    `json:"status,omitempty"` // pending, in_progress, completed, skipped
	Priority         *string    `json:"priority,omitempty"`
	AssignedTo       *string    `json:"assigned_to,omitempty"`
	AssignedToUserID *uint      `json:"assigned_to_user_id,omitempty"`
	DueDate          *string    `json:"due_date,omitempty"` // YYYY-MM-DD format
	Notes            *string    `json:"notes,omitempty"`
	Metadata         *string    `json:"metadata,omitempty"`
}

// CompletePostOnboardingTaskRequest represents the request to complete a task
type CompletePostOnboardingTaskRequest struct {
	ID      uint    `json:"id" binding:"required"`
	Notes   *string `json:"notes,omitempty"` // Completion notes
}

// GetPostOnboardingTasksRequest represents the request to get tasks with filters
type GetPostOnboardingTasksRequest struct {
	EmployeeID *string `form:"employee_id"` // Filter by employee ID
	Status     *string `form:"status"`      // Filter by status
	TaskType   *string `form:"task_type"`   // Filter by task type
	Priority   *string `form:"priority"`    // Filter by priority
	AssignedTo *string `form:"assigned_to"` // Filter by assigned to
	Page       int     `form:"page" binding:"omitempty,min=1"`
	PageSize   int     `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// PostOnboardingTaskResponse represents the response for a post-onboarding task
type PostOnboardingTaskResponse struct {
	ID                uint       `json:"id"`
	EmployeeID        string     `json:"employee_id"`
	EmployeeName      *string    `json:"employee_name,omitempty"`
	TaskType          string     `json:"task_type"`
	Title             string     `json:"title"`
	Description       *string    `json:"description,omitempty"`
	Status            string     `json:"status"`
	Priority          string     `json:"priority"`
	AssignedTo        *string    `json:"assigned_to,omitempty"`
	AssignedToUserID  *uint      `json:"assigned_to_user_id,omitempty"`
	AssignedToUserName *string   `json:"assigned_to_user_name,omitempty"`
	CompletedByUserID *uint      `json:"completed_by_user_id,omitempty"`
	CompletedByUserName *string  `json:"completed_by_user_name,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	DueDate           *time.Time `json:"due_date,omitempty"`
	Notes             *string    `json:"notes,omitempty"`
	Metadata          *string    `json:"metadata,omitempty"`
	CreatedBy         *uint      `json:"created_by,omitempty"`
	UpdatedBy         *uint      `json:"updated_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// BulkCreatePostOnboardingTasksRequest represents the request to create multiple tasks for an employee
type BulkCreatePostOnboardingTasksRequest struct {
	EmployeeID string   `json:"employee_id" binding:"required"`
	TaskTypes  []string `json:"task_types" binding:"required,min=1"` // List of task types to create
	Priority   *string  `json:"priority,omitempty"` // Default priority for all tasks
	DueDate     *string  `json:"due_date,omitempty"` // Default due date for all tasks
}

// GetTaskTypesResponse represents available task types
type GetTaskTypesResponse struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DefaultAssignedTo string `json:"default_assigned_to"` // Default department/role
	IsConditional bool   `json:"is_conditional"` // Whether task is conditional (e.g., dev environment only for developers)
}

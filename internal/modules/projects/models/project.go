package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ProjectStatus represents the current status of a project
type ProjectStatus string

const (
	ProjectStatusPlanning   ProjectStatus = "planning"
	ProjectStatusActive     ProjectStatus = "active"
	ProjectStatusOnHold     ProjectStatus = "on_hold"
	ProjectStatusCompleted  ProjectStatus = "completed"
	ProjectStatusCancelled  ProjectStatus = "cancelled"
)

// ProjectPriority represents the priority level of a project
type ProjectPriority string

const (
	ProjectPriorityLow      ProjectPriority = "low"
	ProjectPriorityMedium   ProjectPriority = "medium"
	ProjectPriorityHigh     ProjectPriority = "high"
	ProjectPriorityCritical ProjectPriority = "critical"
)

// Project represents a project in the system
type Project struct {
	ID              uint            `json:"id" gorm:"primaryKey"`
	ProjectCode     string          `json:"project_code" gorm:"uniqueIndex;not null;size:50"`
	Name            string          `json:"name" gorm:"not null;size:200"`
	Description     *string         `json:"description,omitempty" gorm:"type:text"`
	Status          ProjectStatus   `json:"status" gorm:"type:varchar(20);default:'planning';index"`
	Priority        ProjectPriority `json:"priority" gorm:"type:varchar(20);default:'medium'"`
	DepartmentID    *uint           `json:"department_id,omitempty" gorm:"index"`
	DepartmentName  *string         `json:"department_name,omitempty" gorm:"size:100"`
	Category        *string         `json:"category,omitempty" gorm:"size:100;index"`
	CreatedByID     uint            `json:"created_by_id" gorm:"not null;index"`
	CreatedByName   string          `json:"created_by_name" gorm:"size:200"`
	StartDate       *time.Time      `json:"start_date,omitempty"`
	EndDate         *time.Time      `json:"end_date,omitempty"`
	DueDate         *time.Time      `json:"due_date,omitempty"`
	Progress        float64         `json:"progress" gorm:"type:decimal(5,2);default:0"`
	EstimatedHours  *float64        `json:"estimated_hours,omitempty" gorm:"type:decimal(10,2)"`
	ActualHours     float64         `json:"actual_hours" gorm:"type:decimal(10,2);default:0"`
	Budget          *float64        `json:"budget,omitempty" gorm:"type:decimal(15,2)"`
	Notes           *string         `json:"notes,omitempty" gorm:"type:text"`
	IsActive        bool            `json:"is_active" gorm:"default:true"`
	CompletedAt     *time.Time      `json:"completed_at,omitempty"`
	Members         []ProjectMember `json:"members,omitempty" gorm:"foreignKey:ProjectID"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	UpdatedBy       *uint           `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt       gorm.DeletedAt  `json:"-" gorm:"index"`
}

// TableName specifies the table name for Project
func (Project) TableName() string {
	return "projects"
}

// MemberRole represents the role of a member in a project
type MemberRole string

const (
	MemberRoleLead     MemberRole = "lead"
	MemberRoleMember   MemberRole = "member"
	MemberRoleReviewer MemberRole = "reviewer"
)

// ProjectMember represents an employee assigned to a project
type ProjectMember struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	ProjectID    uint           `json:"project_id" gorm:"not null;index;uniqueIndex:idx_project_employee"`
	EmployeeID   string         `json:"employee_id" gorm:"not null;size:50;index;uniqueIndex:idx_project_employee"`
	EmployeeName string         `json:"employee_name" gorm:"size:200"`
	Role         MemberRole     `json:"role" gorm:"type:varchar(20);default:'member'"`
	AssignedByID *uint          `json:"assigned_by_id,omitempty"`
	JoinedAt     time.Time      `json:"joined_at" gorm:"autoCreateTime"`
	LeftAt       *time.Time     `json:"left_at,omitempty"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for ProjectMember
func (ProjectMember) TableName() string {
	return "project_members"
}

// ──────────────────────────── Request DTOs ────────────────────────────

// CreateProjectRequest is the request body for creating a project
type CreateProjectRequest struct {
	Name           string          `json:"name" binding:"required,min=2,max=200"`
	Description    *string         `json:"description,omitempty"`
	Priority       ProjectPriority `json:"priority,omitempty"`
	DepartmentID   *uint           `json:"department_id,omitempty"`
	Category       *string         `json:"category,omitempty"`
	StartDate      *string         `json:"start_date,omitempty"`  // YYYY-MM-DD or ISO 8601
	EndDate        *string         `json:"end_date,omitempty"`    // YYYY-MM-DD or ISO 8601
	DueDate        *string         `json:"due_date,omitempty"`    // YYYY-MM-DD or ISO 8601
	EstimatedHours *float64        `json:"estimated_hours,omitempty"`
	Budget         *float64        `json:"budget,omitempty"`
	Notes          *string         `json:"notes,omitempty"`
	MemberIDs      []string        `json:"member_ids,omitempty"` // Employee IDs to assign initially
}

// UpdateProjectRequest is the request body for updating a project
type UpdateProjectRequest struct {
	Name           *string          `json:"name,omitempty" binding:"omitempty,min=2,max=200"`
	Description    *string          `json:"description,omitempty"`
	Status         *ProjectStatus   `json:"status,omitempty"`
	Priority       *ProjectPriority `json:"priority,omitempty"`
	DepartmentID   *uint            `json:"department_id,omitempty"`
	Category       *string          `json:"category,omitempty"`
	StartDate      *string          `json:"start_date,omitempty"`  // YYYY-MM-DD or ISO 8601
	EndDate        *string          `json:"end_date,omitempty"`    // YYYY-MM-DD or ISO 8601
	DueDate        *string          `json:"due_date,omitempty"`    // YYYY-MM-DD or ISO 8601
	EstimatedHours *float64         `json:"estimated_hours,omitempty"`
	Budget         *float64         `json:"budget,omitempty"`
	Notes          *string          `json:"notes,omitempty"`
	Progress       *float64         `json:"progress,omitempty"` // Manual progress override
}

// ParseDate parses a date string in YYYY-MM-DD or ISO 8601 format
func ParseDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	// Try YYYY-MM-DD first
	t, err := time.Parse("2006-01-02", s)
	if err == nil {
		return &t, nil
	}
	// Try full ISO 8601
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		return &t, nil
	}
	return nil, fmt.Errorf("invalid date format: %s (expected YYYY-MM-DD or ISO 8601)", s)
}

// AddMembersRequest is the request body for adding members to a project
type AddMembersRequest struct {
	Members []AddMemberEntry `json:"members" binding:"required,min=1"`
}

// AddMemberEntry represents a single member to add
type AddMemberEntry struct {
	EmployeeID string     `json:"employee_id" binding:"required"`
	Role       MemberRole `json:"role,omitempty"`
}

// RemoveMemberRequest is the request body for removing a member
type RemoveMemberRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
}

// AssignProjectRequest is the request body for assigning a project to employees
type AssignProjectRequest struct {
	ProjectID   uint       `json:"project_id" binding:"required"`
	EmployeeIDs []string   `json:"employee_ids" binding:"required,min=1"` // One or more employee IDs
	Role        MemberRole `json:"role,omitempty"`                        // Role for all (default: member)
}

// UnassignProjectRequest is the request body for removing an employee from a project
type UnassignProjectRequest struct {
	ProjectID  uint   `json:"project_id" binding:"required"`
	EmployeeID string `json:"employee_id" binding:"required"`
}

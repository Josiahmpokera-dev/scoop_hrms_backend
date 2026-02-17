package models

import (
	"time"

	"gorm.io/gorm"
)

// DepartmentTarget represents a department-level KPI/target.
type DepartmentTarget struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null;size:20"`
	Title       string         `json:"title" gorm:"not null;size:300"`
	Description *string        `json:"description,omitempty" gorm:"type:text"`
	Department  string         `json:"department" gorm:"not null;size:100;index"`
	Category    string         `json:"category" gorm:"not null;size:50;index"`
	Metric      string         `json:"metric" gorm:"not null;size:200"`
	TargetValue float64        `json:"target_value" gorm:"type:decimal(15,2);not null"`
	CurrentValue float64       `json:"current_value" gorm:"type:decimal(15,2);default:0"`
	Unit        string         `json:"unit" gorm:"not null;size:50"`
	Period      string         `json:"period" gorm:"not null;size:50;index"`
	StartDate   *time.Time     `json:"start_date,omitempty"`
	EndDate     *time.Time     `json:"end_date,omitempty"`
	Status      string         `json:"status" gorm:"type:varchar(30);default:'Not Started';index"`
	Priority    string         `json:"priority" gorm:"size:20"`
	CreatedByID     uint           `json:"created_by_id" gorm:"index"`
	LinkedGoalIDs   *string        `json:"linked_goal_ids,omitempty" gorm:"type:text"` // JSON array of goal codes
	LinkedProjectID *uint          `json:"linked_project_id,omitempty" gorm:"index"`
	LinkedProjectCode *string      `json:"linked_project_code,omitempty" gorm:"size:50"`
	LinkedProjectName *string      `json:"linked_project_name,omitempty" gorm:"size:200"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	Milestones []DepartmentTargetMilestone `json:"milestones,omitempty" gorm:"foreignKey:DepartmentTargetID"`
}

func (DepartmentTarget) TableName() string { return "performance_department_targets" }

// DepartmentTargetMilestone is a milestone for a department target.
type DepartmentTargetMilestone struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	DepartmentTargetID uint           `json:"department_target_id" gorm:"not null;index"`
	Title              string         `json:"title" gorm:"not null;size:300"`
	DueDate            *time.Time     `json:"due_date,omitempty"`
	Status             string         `json:"status" gorm:"type:varchar(30);default:'Not Started'"`
	CompletedDate      *time.Time     `json:"completed_date,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
}

func (DepartmentTargetMilestone) TableName() string { return "performance_department_target_milestones" }

// --- Request DTOs ---

// CreateDepartmentTargetRequest - Create a department target
type CreateDepartmentTargetRequest struct {
	Title        string      `json:"title" binding:"required,max=300"`
	Description  *string     `json:"description,omitempty" binding:"omitempty"`
	Department   string      `json:"department" binding:"required,max=100"`
	Category     string      `json:"category" binding:"required,max=50"`
	Metric       string      `json:"metric" binding:"required,max=200"`
	TargetValue  float64     `json:"target_value" binding:"required"`
	Unit         string      `json:"unit" binding:"required,max=50"`
	Period       string      `json:"period" binding:"required,max=50"`
	StartDate    *time.Time  `json:"start_date,omitempty" binding:"omitempty"`
	EndDate      *time.Time  `json:"end_date,omitempty" binding:"omitempty"`
	Priority     string      `json:"priority,omitempty" binding:"omitempty,max=20"`
	Milestones   []MilestoneInput `json:"milestones,omitempty" binding:"omitempty"`
}

// MilestoneInput - Milestone payload for create
type MilestoneInput struct {
	Title   string     `json:"title" binding:"required,max=300"`
	DueDate *time.Time `json:"due_date,omitempty" binding:"omitempty"`
}

// UpdateDepartmentTargetRequest - Update a department target
type UpdateDepartmentTargetRequest struct {
	Title       *string    `json:"title,omitempty" binding:"omitempty,max=300"`
	Description *string    `json:"description,omitempty" binding:"omitempty"`
	Category    *string    `json:"category,omitempty" binding:"omitempty,max=50"`
	Metric      *string    `json:"metric,omitempty" binding:"omitempty,max=200"`
	TargetValue *float64   `json:"target_value,omitempty" binding:"omitempty"`
	Unit        *string    `json:"unit,omitempty" binding:"omitempty,max=50"`
	Period      *string    `json:"period,omitempty" binding:"omitempty,max=50"`
	StartDate   *time.Time `json:"start_date,omitempty" binding:"omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty" binding:"omitempty"`
	Status      *string    `json:"status,omitempty" binding:"omitempty"`
	Priority    *string    `json:"priority,omitempty" binding:"omitempty,max=20"`
}

// UpdateDepartmentTargetProgressRequest - Update progress on a department target
type UpdateDepartmentTargetProgressRequest struct {
	CurrentValue float64 `json:"current_value" binding:"required"`
}

// LinkDepartmentTargetGoalRequest - Link department target to a goal
type LinkDepartmentTargetGoalRequest struct {
	GoalID uint `json:"goal_id" binding:"required"`
}

// LinkDepartmentTargetProjectRequest - Link department target to a project
type LinkDepartmentTargetProjectRequest struct {
	ProjectID   uint   `json:"project_id" binding:"required"`
	ProjectCode string `json:"project_code" binding:"required,max=50"`
	ProjectName string `json:"project_name" binding:"required,max=200"`
}

// LinkProjectRequestDepartment - Link department target to a project (alias)
type LinkProjectRequestDepartment struct {
	ProjectID   uint   `json:"project_id" binding:"required"`
	ProjectCode string `json:"project_code" binding:"required,max=50"`
	ProjectName string `json:"project_name" binding:"required,max=200"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeeTarget represents an employee-level target (linked to department target or standalone).
type EmployeeTarget struct {
	ID                   uint           `json:"id" gorm:"primaryKey"`
	Code                 string         `json:"code" gorm:"uniqueIndex;not null;size:20"`
	EmployeeID           string         `json:"employee_id" gorm:"not null;size:50;index"`
	DepartmentTargetID   *uint          `json:"department_target_id,omitempty" gorm:"index"`
	Title                string         `json:"title" gorm:"not null;size:300"`
	Description          *string        `json:"description,omitempty" gorm:"type:text"`
	TargetValue          float64        `json:"target_value" gorm:"type:decimal(15,2);not null"`
	CurrentValue         float64        `json:"current_value" gorm:"type:decimal(15,2);default:0"`
	Unit                 string         `json:"unit" gorm:"not null;size:50"`
	Period               string         `json:"period" gorm:"not null;size:50;index"`
	Department           *string        `json:"department,omitempty" gorm:"size:100;index"`
	Status               string         `json:"status" gorm:"type:varchar(30);default:'Not Started';index"`
	AssignedByID         uint           `json:"assigned_by_id" gorm:"not null;index"`
	DueDate              *time.Time     `json:"due_date,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`
}

func (EmployeeTarget) TableName() string { return "performance_employee_targets" }

// --- Request DTOs ---

// CreateEmployeeTargetRequest - Create an employee target
type CreateEmployeeTargetRequest struct {
	EmployeeID         string     `json:"employee_id" binding:"required,max=50"`
	DepartmentTargetID *uint      `json:"department_target_id,omitempty" binding:"omitempty"`
	Title              string     `json:"title" binding:"required,max=300"`
	Description        *string    `json:"description,omitempty" binding:"omitempty"`
	TargetValue        float64    `json:"target_value" binding:"required"`
	Unit               string     `json:"unit" binding:"required,max=50"`
	Period             string     `json:"period" binding:"required,max=50"`
	Department         *string    `json:"department,omitempty" binding:"omitempty,max=100"`
	DueDate            *time.Time `json:"due_date,omitempty" binding:"omitempty"`
}

// UpdateEmployeeTargetRequest - Update an employee target
type UpdateEmployeeTargetRequest struct {
	Title       *string    `json:"title,omitempty" binding:"omitempty,max=300"`
	Description *string    `json:"description,omitempty" binding:"omitempty"`
	TargetValue *float64   `json:"target_value,omitempty" binding:"omitempty"`
	Unit        *string    `json:"unit,omitempty" binding:"omitempty,max=50"`
	Period      *string    `json:"period,omitempty" binding:"omitempty,max=50"`
	Department  *string    `json:"department,omitempty" binding:"omitempty,max=100"`
	Status      *string    `json:"status,omitempty" binding:"omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty" binding:"omitempty"`
}

// UpdateEmployeeTargetProgressRequest - Update progress on an employee target
type UpdateEmployeeTargetProgressRequest struct {
	CurrentValue float64 `json:"current_value" binding:"required"`
}

// BulkAssignTargetRequest - Bulk assign one target template to multiple employees
type BulkAssignTargetRequest struct {
	EmployeeIDs []string   `json:"employee_ids" binding:"required"`
	Title       string     `json:"title" binding:"required,max=300"`
	Description *string    `json:"description,omitempty" binding:"omitempty"`
	TargetValue float64   `json:"target_value" binding:"required"`
	Unit        string    `json:"unit" binding:"required,max=50"`
	Period      string    `json:"period" binding:"required,max=50"`
	Department  *string   `json:"department,omitempty" binding:"omitempty,max=100"`
	DueDate     *time.Time `json:"due_date,omitempty" binding:"omitempty"`
}

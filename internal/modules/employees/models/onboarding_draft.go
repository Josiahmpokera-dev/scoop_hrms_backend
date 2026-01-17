package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// OnboardingStep represents the step number in the onboarding process
type OnboardingStep int

const (
	StepPersonalInfo OnboardingStep = 1  // Personal Information
	StepEmployment   OnboardingStep = 2  // Employment Details
	StepSalary       OnboardingStep = 3  // Salary & CTC
	StepBank         OnboardingStep = 4  // Bank Account
	StepStatutory    OnboardingStep = 5  // Statutory Requirements
	StepDocuments    OnboardingStep = 6  // Documents
	StepAssets       OnboardingStep = 7  // Company Assets
	StepPolicies     OnboardingStep = 8  // Leave & Attendance Policies
	StepEmergency    OnboardingStep = 9  // Emergency Contacts
	StepNotes        OnboardingStep = 10 // Internal Notes
)

// EmployeeOnboardingDraft represents a draft employee onboarding record
type EmployeeOnboardingDraft struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	TenantID        *uint          `json:"tenant_id,omitempty" gorm:"index"`
	EmployeeID      *string        `json:"employee_id,omitempty" gorm:"size:50"`                              // Employee ID (if assigned)
	StepData        string         `json:"step_data,omitempty" gorm:"column:step_data;type:text"`             // JSON object storing step 1 and step 2 data
	CompletedSteps  string         `json:"completed_steps,omitempty" gorm:"column:completed_steps;type:text"` // JSON array of completed step numbers (e.g., "[1,2,3]")
	Progress        float64        `json:"progress" gorm:"column:progress;type:decimal(5,2);default:0"`       // Completion percentage (0-100)
	IsCompleted     bool           `json:"is_completed" gorm:"default:false"`                                 // Whether onboarding is fully completed
	EmployeeIDFinal *uint          `json:"employee_id_final,omitempty" gorm:"index"`                          // Final employee ID after completion
	CreatedBy       *uint          `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy       *uint          `json:"updated_by,omitempty" gorm:"index"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships - will be loaded separately
	Employee *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeIDFinal"`
}

// TableName specifies the table name
func (EmployeeOnboardingDraft) TableName() string {
	return "employee_onboarding_drafts"
}

// GetCompletedStepsList returns the completed steps as a slice of integers
func (d *EmployeeOnboardingDraft) GetCompletedStepsList() []int {
	if d.CompletedSteps == "" {
		return []int{}
	}
	var steps []int
	json.Unmarshal([]byte(d.CompletedSteps), &steps)
	return steps
}

// SetCompletedStepsList sets the completed steps from a slice
func (d *EmployeeOnboardingDraft) SetCompletedStepsList(steps []int) {
	if len(steps) == 0 {
		d.CompletedSteps = "[]"
		return
	}
	data, _ := json.Marshal(steps)
	d.CompletedSteps = string(data)
}

// HasStepCompleted checks if a specific step is completed
func (d *EmployeeOnboardingDraft) HasStepCompleted(step OnboardingStep) bool {
	steps := d.GetCompletedStepsList()
	stepInt := int(step)
	for _, s := range steps {
		if s == stepInt {
			return true
		}
	}
	return false
}

// AddCompletedStep adds a step to completed steps if not already present
func (d *EmployeeOnboardingDraft) AddCompletedStep(step OnboardingStep) {
	steps := d.GetCompletedStepsList()
	stepInt := int(step)
	for _, s := range steps {
		if s == stepInt {
			return // Already completed
		}
	}
	steps = append(steps, stepInt)
	d.SetCompletedStepsList(steps)
}

// CalculateProgress calculates the completion percentage
func (d *EmployeeOnboardingDraft) CalculateProgress() float64 {
	// Total steps = 10
	totalSteps := 10.0
	steps := d.GetCompletedStepsList()
	completedCount := float64(len(steps))
	return (completedCount / totalSteps) * 100.0
}

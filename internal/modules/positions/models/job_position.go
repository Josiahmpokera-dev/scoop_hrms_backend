package models

import (
	"time"

	"gorm.io/gorm"
)

// JobPosition represents a job position/title in the system
type JobPosition struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	OrganizationID      *uint          `json:"organization_id,omitempty" gorm:"index"`        // References organizations(id)
	Code                string         `json:"code" gorm:"uniqueIndex;not null;size:20"`      // Position Code *
	Title               string         `json:"title" gorm:"not null;size:100"`                // Position Title *
	Grade               *string        `json:"grade,omitempty" gorm:"size:20"`                // Grade Level *
	DepartmentID        *uint          `json:"department_id,omitempty" gorm:"index"`          // Department * - References departments(id)
	ReportsToPositionID *uint          `json:"reports_to_position_id,omitempty" gorm:"index"` // Reports To - References job_positions(id)
	BudgetedHeadcount   *int           `json:"budgeted_headcount,omitempty"`                  // Budgeted Headcount
	CurrentHeadcount    *int           `json:"current_headcount,omitempty"`                   // Current Headcount
	EmploymentType      *string        `json:"employment_type,omitempty" gorm:"size:50"`      // Employment Type: full_time, part_time, contract, intern
	KeyCompetencies     *string        `json:"key_competencies,omitempty" gorm:"type:text"`   // Key Competencies (comma separated)
	IsActive            bool           `json:"is_active" gorm:"default:true"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	UpdatedBy           *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for JobPosition model
func (JobPosition) TableName() string {
	return "job_positions"
}

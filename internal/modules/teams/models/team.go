package models

import (
	"time"

	"gorm.io/gorm"
)

// TeamType represents the type of team
type TeamType string

const (
	TeamTypeProject         TeamType = "project"
	TeamTypeFunctional      TeamType = "functional"
	TeamTypeCrossFunctional TeamType = "cross-functional"
	TeamTypeVirtual         TeamType = "virtual"
)

// Team represents a team within a department
type Team struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	DepartmentID *uint         `json:"department_id,omitempty" gorm:"index"` // References departments(id)
	Code        string         `json:"code" gorm:"uniqueIndex;not null;size:20"`
	Name        string         `json:"name" gorm:"not null;size:100"`
	Description *string        `json:"description,omitempty" gorm:"type:text"`
	TeamLeadID  *uint          `json:"team_lead_id,omitempty" gorm:"index"` // References employees(id)
	TeamType    *string        `json:"team_type,omitempty" gorm:"size:50"` // project, functional, cross-functional, virtual
	ProjectID   *uint          `json:"project_id,omitempty" gorm:"index"` // For project-based teams
	MaxMembers  *int           `json:"max_members,omitempty" gorm:"default:10"`
	LocationID  *uint          `json:"location_id,omitempty" gorm:"index"` // References locations(id)
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UpdatedBy   *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for Team model
func (Team) TableName() string {
	return "teams"
}

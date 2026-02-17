package models

import (
	"time"

	"gorm.io/gorm"
)

// Feedback360Campaign represents a 360° feedback campaign.
type Feedback360Campaign struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null;size:20"`
	Name        string         `json:"name" gorm:"not null;size:200"`
	Description *string        `json:"description,omitempty" gorm:"type:text"`
	Department  *string        `json:"department,omitempty" gorm:"size:100;index"`
	Status      string         `json:"status" gorm:"type:varchar(30);default:'Draft';index"`
	StartDate   *time.Time     `json:"start_date,omitempty"`
	EndDate     *time.Time     `json:"end_date,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	RaterGroups []Feedback360RaterGroup `json:"rater_groups,omitempty" gorm:"foreignKey:CampaignID"`
}

func (Feedback360Campaign) TableName() string { return "performance_feedback360_campaigns" }

// Feedback360RaterGroup defines a group of raters (e.g. Peers, Direct Reports).
type Feedback360RaterGroup struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CampaignID uint           `json:"campaign_id" gorm:"not null;index"`
	Name       string         `json:"name" gorm:"not null;size:100"`
	MinRaters  int            `json:"min_raters" gorm:"default:1"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Feedback360RaterGroup) TableName() string { return "performance_feedback360_rater_groups" }

// --- Request DTOs ---

// LaunchFeedback360Request - Launch a 360 feedback campaign
type LaunchFeedback360Request struct {
	Name        string     `json:"name" binding:"required,max=200"`
	Description *string    `json:"description,omitempty" binding:"omitempty"`
	Department  *string    `json:"department,omitempty" binding:"omitempty,max=100"`
	StartDate   *time.Time `json:"start_date,omitempty" binding:"omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty" binding:"omitempty"`
	RaterGroups []RaterGroupInput `json:"rater_groups,omitempty" binding:"omitempty"`
}

// RaterGroupInput - Rater group for launch
type RaterGroupInput struct {
	Name      string `json:"name" binding:"required,max=100"`
	MinRaters int    `json:"min_raters" binding:"omitempty"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// TalentReview - Employee talent review entry (9-box)
type TalentReview struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	EmployeeID          uint           `json:"employee_id" gorm:"not null;uniqueIndex"`
	EmployeeName        string         `json:"employee_name" gorm:"size:200"`
	Designation         *string        `json:"designation,omitempty" gorm:"size:200"`
	Department          *string        `json:"department,omitempty" gorm:"size:100;index"`
	Photo               *string        `json:"photo,omitempty" gorm:"size:500"`
	PerformanceRating   float64        `json:"performance_rating" gorm:"type:decimal(3,1);default:0"`
	PotentialRating     float64        `json:"potential_rating" gorm:"type:decimal(3,1);default:0"`
	Box                 int            `json:"box" gorm:"default:5;index"` // 1-9
	BoxLabel            string         `json:"box_label" gorm:"size:50"`
	CriticalRole        bool           `json:"critical_role" gorm:"default:false"`
	Tenure              *string        `json:"tenure,omitempty" gorm:"size:50"`
	RiskOfLoss          string         `json:"risk_of_loss" gorm:"type:varchar(20);default:'Low'"` // Low, Medium, High
	Readiness           *string        `json:"readiness,omitempty" gorm:"size:30"` // Now, 6-12 Months, 1-2 Years
	DevelopmentPriority *string        `json:"development_priority,omitempty" gorm:"size:20"` // Low, Medium, High
	SuccessorFor        *string        `json:"successor_for,omitempty" gorm:"type:text"` // JSON
	Notes               *string        `json:"notes,omitempty" gorm:"type:text"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}

func (TalentReview) TableName() string { return "performance_talent_reviews" }

// CalibrationSession - Calibration session for talent review
type CalibrationSession struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"not null;size:200"`
	Department     *string        `json:"department,omitempty" gorm:"size:100"`
	ParticipantIDs *string        `json:"participant_ids,omitempty" gorm:"type:text"` // JSON
	CreatedByID    uint           `json:"created_by_id" gorm:"not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (CalibrationSession) TableName() string { return "performance_calibration_sessions" }

// SuccessionPlan - Succession plan for a role
type SuccessionPlan struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	EmployeeID      uint           `json:"employee_id" gorm:"not null;index"`
	Role            string         `json:"role" gorm:"not null;size:200"`
	Readiness       string         `json:"readiness" gorm:"size:30"`
	DevelopmentPlan *string        `json:"development_plan,omitempty" gorm:"type:text"`
	CreatedByID     uint           `json:"created_by_id" gorm:"not null"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (SuccessionPlan) TableName() string { return "performance_succession_plans" }

// --- Request DTOs ---

// UpdateTalentReviewRequest - Update a talent review entry
type UpdateTalentReviewRequest struct {
	PerformanceRating   *float64 `json:"performance_rating,omitempty" binding:"omitempty"`
	PotentialRating     *float64 `json:"potential_rating,omitempty" binding:"omitempty"`
	Box                 *int     `json:"box,omitempty" binding:"omitempty,min=1,max=9"`
	BoxLabel            *string  `json:"box_label,omitempty" binding:"omitempty,max=50"`
	CriticalRole        *bool    `json:"critical_role,omitempty" binding:"omitempty"`
	Tenure              *string  `json:"tenure,omitempty" binding:"omitempty,max=50"`
	RiskOfLoss          *string  `json:"risk_of_loss,omitempty" binding:"omitempty,oneof=Low Medium High"`
	Readiness           *string  `json:"readiness,omitempty" binding:"omitempty,max=30"`
	DevelopmentPriority *string  `json:"development_priority,omitempty" binding:"omitempty,oneof=Low Medium High"`
	SuccessorFor        *string  `json:"successor_for,omitempty" binding:"omitempty"`
	Notes               *string  `json:"notes,omitempty" binding:"omitempty"`
}

// CreateCalibrationSessionRequest - Create a calibration session
type CreateCalibrationSessionRequest struct {
	Name           string  `json:"name" binding:"required,max=200"`
	Department     *string `json:"department,omitempty" binding:"omitempty,max=100"`
	ParticipantIDs *string `json:"participant_ids,omitempty" binding:"omitempty"` // JSON array
	CreatedByID    uint    `json:"created_by_id" binding:"required"`
}

// AddSuccessionPlanRequest - Add a succession plan
type AddSuccessionPlanRequest struct {
	EmployeeID      uint    `json:"employee_id" binding:"required"`
	Role            string  `json:"role" binding:"required,max=200"`
	Readiness       string  `json:"readiness" binding:"required,max=30"`
	DevelopmentPlan *string `json:"development_plan,omitempty" binding:"omitempty"`
	CreatedByID     uint    `json:"created_by_id" binding:"required"`
}

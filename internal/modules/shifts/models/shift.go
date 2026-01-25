package models

import (
	"time"

	locationModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/models"
	"gorm.io/gorm"
)

// ShiftType represents the type of shift
type ShiftType string

const (
	ShiftTypeFixed     ShiftType = "fixed"
	ShiftTypeSplit     ShiftType = "split"
	ShiftTypeRotating  ShiftType = "rotating"
	ShiftTypeNight     ShiftType = "night"
	ShiftTypeFlexible  ShiftType = "flexible"
)

// Shift represents a work shift configuration
type Shift struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	TenantID           *uint          `json:"tenant_id,omitempty" gorm:"index"`
	ShiftName          string         `json:"shift_name" gorm:"not null;size:100"`
	ShiftCode          string         `json:"shift_code" gorm:"uniqueIndex;not null;size:20"`
	ShiftType          string         `json:"shift_type" gorm:"not null;size:50"` // fixed, split, rotating, night, flexible
	StartTime          string         `json:"start_time" gorm:"not null;size:8"` // HH:MM:SS format
	EndTime            string         `json:"end_time" gorm:"not null;size:8"`   // HH:MM:SS format
	WorkingHours       float64        `json:"working_hours" gorm:"type:decimal(5,2)"`
	BreakDuration      int            `json:"break_duration" gorm:"default:0"` // in minutes
	GraceMinutes       int            `json:"grace_minutes" gorm:"default:0"` // grace period before marking late
	LateMarkAfter      int            `json:"late_mark_after" gorm:"default:0"` // minutes after start time to mark as late
	EarlyGoingMinutes  int            `json:"early_going_minutes" gorm:"default:0"` // minutes before end time allowed for early leaving
	HalfDayHours       float64        `json:"half_day_hours" gorm:"type:decimal(5,2)"`
	MinimumHours       float64        `json:"minimum_hours" gorm:"type:decimal(5,2)"`
	CrossDay           bool           `json:"cross_day" gorm:"default:false"` // shift crosses midnight
	NightShift         bool           `json:"night_shift" gorm:"default:false"`
	WeeklyOff          string         `json:"weekly_off" gorm:"type:jsonb"` // JSON array of day names
	IsActive           bool           `json:"is_active" gorm:"default:true"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	CreatedBy          *uint          `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy          *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`

	// Many-to-many relationship with locations
	// Uses the actual Location model from locations module
	// Join table: shift_locations with columns: shift_id, location_id
	Locations []locationModels.Location `json:"locations,omitempty" gorm:"many2many:shift_locations;"`
}

// TableName specifies the table name
func (Shift) TableName() string {
	return "shifts"
}

// ShiftLocation represents the join table for shifts and locations
type ShiftLocation struct {
	ShiftID    uint `gorm:"primaryKey;column:shift_id"`
	LocationID uint `gorm:"primaryKey;column:location_id"`
}

// TableName specifies the table name
func (ShiftLocation) TableName() string {
	return "shift_locations"
}

// CreateShiftRequest represents the request to create a shift
type CreateShiftRequest struct {
	ShiftName         string   `json:"shift_name" binding:"required"`
	ShiftCode         string   `json:"shift_code" binding:"required"`
	ShiftType         string   `json:"shift_type" binding:"required,oneof=fixed split rotating night flexible"`
	StartTime         string   `json:"start_time" binding:"required"`
	EndTime           string   `json:"end_time" binding:"required"`
	BreakDuration     *int    `json:"break_duration"`
	GraceMinutes      *int    `json:"grace_minutes"`
	LateMarkAfter     *int    `json:"late_mark_after"`
	EarlyGoingMinutes *int    `json:"early_going_minutes"`
	HalfDayHours      *float64 `json:"half_day_hours"`
	MinimumHours      *float64 `json:"minimum_hours"`
	CrossDay          *bool   `json:"cross_day"`
	NightShift        *bool   `json:"night_shift"`
	WeeklyOff         []string `json:"weekly_off"` // Array of day names
	LocationIDs       []uint   `json:"location_ids"`
	IsActive          *bool   `json:"is_active"`
}

// UpdateShiftRequest represents the request to update a shift
type UpdateShiftRequest struct {
	ShiftName         *string   `json:"shift_name"`
	ShiftCode         *string   `json:"shift_code"`
	ShiftType         *string   `json:"shift_type"`
	StartTime         *string   `json:"start_time"`
	EndTime           *string   `json:"end_time"`
	BreakDuration     *int     `json:"break_duration"`
	GraceMinutes      *int     `json:"grace_minutes"`
	LateMarkAfter     *int     `json:"late_mark_after"`
	EarlyGoingMinutes *int     `json:"early_going_minutes"`
	HalfDayHours      *float64 `json:"half_day_hours"`
	MinimumHours      *float64 `json:"minimum_hours"`
	CrossDay          *bool    `json:"cross_day"`
	NightShift        *bool    `json:"night_shift"`
	WeeklyOff         []string `json:"weekly_off"`
	LocationIDs       []uint   `json:"location_ids"`
	IsActive          *bool    `json:"is_active"`
}

// DuplicateShiftRequest represents the request to duplicate a shift
type DuplicateShiftRequest struct {
	ShiftName *string `json:"shift_name"`
	ShiftCode *string `json:"shift_code"`
}

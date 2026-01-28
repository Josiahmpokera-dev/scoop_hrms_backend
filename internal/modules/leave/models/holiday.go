package models

import (
	"time"

	"gorm.io/gorm"
)

// HolidayType represents the type of holiday
type HolidayType string

const (
	HolidayTypePublic    HolidayType = "Public"
	HolidayTypeCompany   HolidayType = "Company"
	HolidayTypeOptional  HolidayType = "Optional"
	HolidayTypeRestricted HolidayType = "Restricted"
)

// Holiday represents a holiday
type Holiday struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	TenantID        *uint          `json:"tenant_id,omitempty" gorm:"index"`
	Name            string         `json:"name" gorm:"not null;size:200"`
	Date            time.Time      `json:"date" gorm:"not null;index"` // Date of holiday
	Type            string         `json:"type" gorm:"not null;size:50"` // Public, Company, Optional, Restricted
	IsFloater       bool           `json:"is_floater" gorm:"default:false"` // Floating holiday
	Location        string         `json:"location" gorm:"type:jsonb"` // JSON array of location names or "All Locations"
	ApplicableFor   string         `json:"applicable_for" gorm:"type:jsonb"` // JSON array of employee groups or "All Employees"
	Description     *string        `json:"description,omitempty" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	CreatedBy       *uint          `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy       *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (Holiday) TableName() string {
	return "holidays"
}

// CreateHolidayRequest represents the request to create a holiday
type CreateHolidayRequest struct {
	Name          string   `json:"name" binding:"required"`
	Date          string   `json:"date" binding:"required"` // YYYY-MM-DD
	Type          string   `json:"type" binding:"required"`
	IsFloater     *bool    `json:"is_floater"`
	Location      []string `json:"location" binding:"required"` // Array of location names
	ApplicableFor []string `json:"applicable_for" binding:"required"` // Array of employee groups
	Description   *string  `json:"description"`
}

// UpdateHolidayRequest represents the request to update a holiday
type UpdateHolidayRequest struct {
	Name          *string   `json:"name"`
	Date          *string   `json:"date"`
	Type          *string   `json:"type"`
	IsFloater     *bool     `json:"is_floater"`
	Location      []string  `json:"location"`
	ApplicableFor []string  `json:"applicable_for"`
	Description   *string   `json:"description"`
}

// BulkUploadHolidaysRequest represents the request for bulk upload
type BulkUploadHolidaysRequest struct {
	Year              int  `json:"year" binding:"required"`
	OverwriteExisting *bool `json:"overwrite_existing"`
}

// ImportHolidaysRequest represents the request to import holidays from API
type ImportHolidaysRequest struct {
	Country           string `json:"country" binding:"required"`
	Year              int    `json:"year" binding:"required"`
	OverwriteExisting *bool  `json:"overwrite_existing"`
}

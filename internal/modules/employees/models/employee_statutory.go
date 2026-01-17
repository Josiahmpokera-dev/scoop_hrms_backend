package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeeStatutoryInfo represents statutory and work authorization information
type EmployeeStatutoryInfo struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	EmployeeID            *uint          `json:"employee_id,omitempty" gorm:"index"` // Nullable for drafts
	EmployeeIDString      *string        `json:"employee_id_string,omitempty" gorm:"size:50;index"` // Employee ID string (e.g., "EMP466088") - for tracking during draft
	DraftID               *uint          `json:"draft_id,omitempty" gorm:"index"`   // References employee_onboarding_drafts(id)
	
	// Tanzania Statutory Requirements
	TINNumber             *string        `json:"tin_number,omitempty" gorm:"size:50"` // Tax Identification Number
	NSSFNumber            *string        `json:"nssf_number,omitempty" gorm:"size:50"` // National Social Security Fund
	NHIFNumber            *string        `json:"nhif_number,omitempty" gorm:"size:50"` // National Health Insurance Fund
	WCFNumber             *string        `json:"wcf_number,omitempty" gorm:"size:50"` // Workers Compensation Fund
	SDLNumber             *string        `json:"sdl_number,omitempty" gorm:"size:50"` // Skills Development Levy

	// Work Authorization (For Expatriates)
	PassportNumber        *string        `json:"passport_number,omitempty" gorm:"size:50"`
	PassportExpiryDate    *time.Time     `json:"passport_expiry_date,omitempty"`
	WorkPermitNumber      *string        `json:"work_permit_number,omitempty" gorm:"size:50"`
	WorkPermitExpiryDate  *time.Time     `json:"work_permit_expiry_date,omitempty"`

	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (EmployeeStatutoryInfo) TableName() string {
	return "employee_statutory_info"
}

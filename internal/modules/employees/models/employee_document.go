package models

import (
	"time"

	"gorm.io/gorm"
)

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypeIdentity       DocumentType = "identity"
	DocumentTypeWorkPermit     DocumentType = "work_permit"
	DocumentTypeEducation      DocumentType = "education"
	DocumentTypeContract       DocumentType = "contract"
	DocumentTypeTaxStatutory   DocumentType = "tax_statutory"
	DocumentTypeOther          DocumentType = "other"
)

// EmployeeDocument represents a document uploaded for an employee
type EmployeeDocument struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	EmployeeID  *uint          `json:"employee_id,omitempty" gorm:"index"` // Nullable for drafts
	EmployeeIDString *string    `json:"employee_id_string,omitempty" gorm:"size:50;index"` // Employee ID string (e.g., "EMP466088") - for tracking during draft
	DraftID     *uint          `json:"draft_id,omitempty" gorm:"index"`   // References employee_onboarding_drafts(id)
	DocumentType DocumentType  `json:"document_type" gorm:"type:varchar(50);not null"`
	FileName    string         `json:"file_name" gorm:"not null;size:255"`
	FileURL     string         `json:"file_url" gorm:"not null;size:500"` // URL to the stored file
	FileSize    *int64         `json:"file_size,omitempty"` // File size in bytes
	MimeType    *string        `json:"mime_type,omitempty" gorm:"size:100"`
	Description *string        `json:"description,omitempty" gorm:"type:text"`
	UploadedBy  *uint          `json:"uploaded_by,omitempty" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (EmployeeDocument) TableName() string {
	return "employee_documents"
}

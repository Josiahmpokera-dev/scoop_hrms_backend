package models

import (
	"time"

	"gorm.io/gorm"
)

// BioTimeTransaction represents a BioTime transaction stored in our database
type BioTimeTransaction struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	TenantID           *uint          `json:"tenant_id,omitempty" gorm:"index"`

	// BioTime Transaction Data
	BioTimeTransactionID int          `json:"biotime_transaction_id" gorm:"column:biotime_transaction_id;type:integer;not null;index"` // ID from BioTime API
	EmpCode             string        `json:"emp_code" gorm:"size:100;index"`
	FirstName           string        `json:"first_name" gorm:"size:255"`
	LastName            string        `json:"last_name" gorm:"size:255"`
	Department          string        `json:"department" gorm:"size:255"`
	Position            string        `json:"position" gorm:"size:255"`
	PunchTime           time.Time     `json:"punch_time" gorm:"index"`
	PunchState          string        `json:"punch_state" gorm:"size:50"`
	PunchStateDisplay   string        `json:"punch_state_display" gorm:"size:100"`
	VerifyType          int           `json:"verify_type"`
	VerifyTypeDisplay   string        `json:"verify_type_display" gorm:"size:100"`
	WorkCode            string        `json:"work_code" gorm:"size:100"`
	GPSLocation         string        `json:"gps_location" gorm:"type:text"`
	AreaAlias           *string       `json:"area_alias" gorm:"size:255"`
	TerminalSN          string        `json:"terminal_sn" gorm:"size:100;index"`
	Temperature         float64       `json:"temperature"`
	TerminalAlias       *string       `json:"terminal_alias" gorm:"size:255"`
	UploadTime          *time.Time    `json:"upload_time"`
	
	// Metadata
	SyncedAt            time.Time     `json:"synced_at" gorm:"default:CURRENT_TIMESTAMP"`
	
	// Timestamps
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (BioTimeTransaction) TableName() string {
	return "biotime_transactions"
}

// Unique constraint: (tenant_id, biotime_transaction_id, punch_time) to prevent duplicates
// This will be handled at application level

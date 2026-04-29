package models

import (
	"time"

	"gorm.io/gorm"
)

// EmailLog tracks automated emails sent by the system (like monthly reports)
type EmailLog struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Recipient    string         `gorm:"size:255;not null" json:"recipient"`
	Subject      string         `gorm:"size:255;not null" json:"subject"`
	ReportMonth  string         `gorm:"size:50" json:"report_month"` // e.g., "2026-04"
	Status       string         `gorm:"size:50;not null" json:"status"` // "Delivered", "Failed"
	ErrorMessage string         `gorm:"type:text" json:"error_message,omitempty"`
	DownloadURL  string         `gorm:"size:255" json:"download_url,omitempty"`
	FileType     string         `gorm:"size:50" json:"file_type,omitempty"` // e.g., "xlsx", "pdf"
	FileSize     int64          `json:"file_size,omitempty"`
	CreatedAt    time.Time      `json:"sent_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// EmailLogResponse is the data returned to the API
type EmailLogResponse struct {
	ID          uint   `json:"id"`
	Recipient   string `json:"recipient"`
	Subject     string `json:"subject"`
	ReportMonth string `json:"report_month"`
	Status      string `json:"status"`
	DownloadURL string `json:"download_url,omitempty"`
	FileType    string `json:"file_type,omitempty"`
	FileSize    int64  `json:"file_size,omitempty"`
	Date        string `json:"date"` // formatted CreatedAt
}

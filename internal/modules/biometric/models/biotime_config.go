package models

import (
	"time"

	"gorm.io/gorm"
)

// BioTimeConfig represents BioTime API configuration stored in database
type BioTimeConfig struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	
	// API Configuration
	BaseURL     string         `json:"base_url" gorm:"not null;size:255"`
	Username    string         `json:"username" gorm:"not null;size:100"`
	Password    string         `json:"password" gorm:"not null;size:255"` // Encrypted in production
	
	// Token Management
	Token       *string        `json:"token,omitempty" gorm:"type:text"` // JWT token from BioTime
	TokenExpiry *time.Time     `json:"token_expiry,omitempty"` // When token expires
	
	// Status
	Enabled     bool           `json:"enabled" gorm:"default:true"`
	LastSync    *time.Time     `json:"last_sync,omitempty"` // Last successful sync/connection
	LastError   *string        `json:"last_error,omitempty" gorm:"type:text"` // Last error message
	
	// Timestamps
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UpdatedBy   *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (BioTimeConfig) TableName() string {
	return "biotime_configs"
}

// IsTokenValid checks if the stored token is still valid
func (b *BioTimeConfig) IsTokenValid() bool {
	if b.Token == nil || *b.Token == "" {
		return false
	}
	if b.TokenExpiry == nil {
		return false
	}
	// Consider token valid if it expires in more than 5 minutes (buffer)
	return time.Now().Add(5 * time.Minute).Before(*b.TokenExpiry)
}

// SetToken sets the token and calculates expiry
func (b *BioTimeConfig) SetToken(token string, expiresIn time.Duration) {
	b.Token = &token
	expiry := time.Now().Add(expiresIn)
	b.TokenExpiry = &expiry
}

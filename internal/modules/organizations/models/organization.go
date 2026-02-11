package models

import (
	"time"

	"gorm.io/gorm"
)

// OrganizationStatus represents organization status
type OrganizationStatus string

const (
	OrganizationStatusActive    OrganizationStatus = "active"
	OrganizationStatusSuspended OrganizationStatus = "suspended"
	OrganizationStatusInactive  OrganizationStatus = "inactive"
)

// Organization represents an organization/company in the system
type Organization struct {
	ID                uint               `json:"id" gorm:"primaryKey"`
	Code              string             `json:"code" gorm:"uniqueIndex;not null;size:50"`
	Name              string             `json:"name" gorm:"not null;size:255"`
	LegalName         *string            `json:"legal_name,omitempty" gorm:"size:255"`
	RegistrationNumber *string           `json:"registration_number,omitempty" gorm:"size:100"`
	TaxID             *string            `json:"tax_id,omitempty" gorm:"size:100"`
	Industry          *string            `json:"industry,omitempty" gorm:"size:100"`
	CompanySize       *string            `json:"company_size,omitempty" gorm:"size:50"` // startup, small, medium, large, enterprise
	Description       *string            `json:"description,omitempty" gorm:"type:text"`
	Website           *string            `json:"website,omitempty" gorm:"size:255"`
	Domain            *string            `json:"domain" gorm:"size:255"` // Organization domain (e.g., company.com)
	LogoURL           *string            `json:"logo_url,omitempty" gorm:"size:500"`
	AddressLine1      *string            `json:"address_line1,omitempty" gorm:"type:text"`
	AddressLine2      *string            `json:"address_line2,omitempty" gorm:"type:text"`
	City              *string            `json:"city,omitempty" gorm:"size:100"`
	State             *string            `json:"state,omitempty" gorm:"size:100"`
	Country           *string            `json:"country,omitempty" gorm:"size:100"`
	PostalCode        *string            `json:"postal_code,omitempty" gorm:"size:20"`
	PhoneNumber       *string            `json:"phone_number,omitempty" gorm:"size:50"`
	Email             *string            `json:"email,omitempty" gorm:"size:191"`
	FoundedDate       *time.Time         `json:"founded_date,omitempty"`
	FiscalYearStart   *int               `json:"fiscal_year_start,omitempty" gorm:"default:1"` // 1-12 (January to December)
	Timezone          *string            `json:"timezone,omitempty" gorm:"size:100;default:'Africa/Dar_es_Salaam'"`
	CurrencyCode      *string            `json:"currency_code,omitempty" gorm:"size:3;default:'TZS'"`
	Status            OrganizationStatus `json:"status" gorm:"type:varchar(50);default:'active';not null"`
	IsActive          bool               `json:"is_active" gorm:"default:true"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	UpdatedBy         *uint              `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt         gorm.DeletedAt     `json:"-" gorm:"index"`
}

// TableName specifies the table name for Organization model
func (Organization) TableName() string {
	return "organizations"
}

// IsActiveStatus checks if organization is active
func (o *Organization) IsActiveStatus() bool {
	return o.Status == OrganizationStatusActive && o.IsActive
}

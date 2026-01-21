package models

import "time"

// CreateOrganizationRequest represents the request to create an organization
type CreateOrganizationRequest struct {
	Code              string     `json:"code" binding:"required,min=2,max=50"`
	Name              string     `json:"name" binding:"required,min=2,max=255"`
	LegalName         *string    `json:"legal_name,omitempty" binding:"omitempty,max=255"`
	RegistrationNumber *string   `json:"registration_number,omitempty" binding:"omitempty,max=100"`
	TaxID             *string    `json:"tax_id,omitempty" binding:"omitempty,max=100"`
	Industry          *string    `json:"industry,omitempty" binding:"omitempty,max=100"`
	CompanySize       *string    `json:"company_size,omitempty" binding:"omitempty,oneof=startup small medium large enterprise"`
	Description       *string    `json:"description,omitempty"`
	Website           *string    `json:"website,omitempty" binding:"omitempty,max=255"`
	Domain            *string    `json:"domain,omitempty" binding:"omitempty,max=255"` // Organization domain (e.g., company.com)
	LogoURL           *string    `json:"logo_url,omitempty" binding:"omitempty,max=500"`
	AddressLine1      *string    `json:"address_line1,omitempty"`
	AddressLine2      *string    `json:"address_line2,omitempty"`
	City              *string    `json:"city,omitempty" binding:"omitempty,max=100"`
	State             *string    `json:"state,omitempty" binding:"omitempty,max=100"`
	Country           *string    `json:"country,omitempty" binding:"omitempty,max=100"`
	PostalCode        *string    `json:"postal_code,omitempty" binding:"omitempty,max=20"`
	PhoneNumber       *string    `json:"phone_number,omitempty" binding:"omitempty,max=50"`
	Email             *string    `json:"email,omitempty" binding:"omitempty,email"`
	FoundedDate       *time.Time `json:"founded_date,omitempty"`
	FiscalYearStart   *int       `json:"fiscal_year_start,omitempty" binding:"omitempty,min=1,max=12"`
	Timezone          *string    `json:"timezone,omitempty" binding:"omitempty,max=100"`
	CurrencyCode      *string    `json:"currency_code,omitempty" binding:"omitempty,len=3"`
	Status            *string    `json:"status,omitempty"`
	IsActive          *bool      `json:"is_active,omitempty"`
}

// UpdateOrganizationRequest represents the request to update an organization
type UpdateOrganizationRequest struct {
	Code              *string    `json:"code,omitempty" binding:"omitempty,min=2,max=50"`
	Name              *string    `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	LegalName         *string    `json:"legal_name,omitempty" binding:"omitempty,max=255"`
	RegistrationNumber *string   `json:"registration_number,omitempty" binding:"omitempty,max=100"`
	TaxID             *string    `json:"tax_id,omitempty" binding:"omitempty,max=100"`
	Industry          *string    `json:"industry,omitempty" binding:"omitempty,max=100"`
	CompanySize       *string    `json:"company_size,omitempty" binding:"omitempty,oneof=startup small medium large enterprise"`
	Description       *string    `json:"description,omitempty"`
	Website           *string    `json:"website,omitempty" binding:"omitempty,max=255"`
	Domain            *string    `json:"domain,omitempty" binding:"omitempty,max=255"` // Organization domain (e.g., company.com)
	LogoURL           *string    `json:"logo_url,omitempty" binding:"omitempty,max=500"`
	AddressLine1      *string    `json:"address_line1,omitempty"`
	AddressLine2      *string    `json:"address_line2,omitempty"`
	City              *string    `json:"city,omitempty" binding:"omitempty,max=100"`
	State             *string    `json:"state,omitempty" binding:"omitempty,max=100"`
	Country           *string    `json:"country,omitempty" binding:"omitempty,max=100"`
	PostalCode        *string    `json:"postal_code,omitempty" binding:"omitempty,max=20"`
	PhoneNumber       *string    `json:"phone_number,omitempty" binding:"omitempty,max=50"`
	Email             *string    `json:"email,omitempty" binding:"omitempty,email"`
	FoundedDate       *time.Time `json:"founded_date,omitempty"`
	FiscalYearStart   *int       `json:"fiscal_year_start,omitempty" binding:"omitempty,min=1,max=12"`
	Timezone          *string    `json:"timezone,omitempty" binding:"omitempty,max=100"`
	CurrencyCode      *string    `json:"currency_code,omitempty" binding:"omitempty,len=3"`
	Status            *string    `json:"status,omitempty"`
	IsActive          *bool      `json:"is_active,omitempty"`
}

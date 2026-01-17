package models

// CreateLocationRequest represents the request to create a location
type CreateLocationRequest struct {
	OrganizationID *uint    `json:"organization_id,omitempty"`
	Name           string   `json:"name" binding:"required,min=2,max=100"`
	LocationType   *string  `json:"location_type,omitempty" binding:"omitempty,oneof=head_office branch remote satellite"`
	AddressLine1   *string  `json:"address_line1,omitempty"`
	AddressLine2   *string  `json:"address_line2,omitempty"`
	City           *string  `json:"city,omitempty" binding:"omitempty,max=100"`
	State          *string  `json:"state,omitempty" binding:"omitempty,max=100"`
	Country        *string  `json:"country,omitempty" binding:"omitempty,max=100"`
	PostalCode     *string  `json:"postal_code,omitempty" binding:"omitempty,max=20"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	Timezone       *string  `json:"timezone,omitempty" binding:"omitempty,max=100"`
	PhoneNumber    *string  `json:"phone_number,omitempty" binding:"omitempty,max=50"`
	Capacity       *int     `json:"capacity,omitempty"`
	Facilities     *string  `json:"facilities,omitempty"` // JSON array string
	IsHeadOffice   *bool    `json:"is_head_office,omitempty"`
	IsActive       *bool    `json:"is_active,omitempty"`
}

// UpdateLocationRequest represents the request to update a location
type UpdateLocationRequest struct {
	OrganizationID *uint    `json:"organization_id,omitempty"`
	Name           *string  `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	LocationType   *string  `json:"location_type,omitempty" binding:"omitempty,oneof=head_office branch remote satellite"`
	AddressLine1   *string  `json:"address_line1,omitempty"`
	AddressLine2   *string  `json:"address_line2,omitempty"`
	City           *string  `json:"city,omitempty" binding:"omitempty,max=100"`
	State          *string  `json:"state,omitempty" binding:"omitempty,max=100"`
	Country        *string  `json:"country,omitempty" binding:"omitempty,max=100"`
	PostalCode     *string  `json:"postal_code,omitempty" binding:"omitempty,max=20"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	Timezone       *string  `json:"timezone,omitempty" binding:"omitempty,max=100"`
	PhoneNumber    *string  `json:"phone_number,omitempty" binding:"omitempty,max=50"`
	Capacity       *int     `json:"capacity,omitempty"`
	Facilities     *string  `json:"facilities,omitempty"`
	IsHeadOffice   *bool    `json:"is_head_office,omitempty"`
	IsActive       *bool    `json:"is_active,omitempty"`
}

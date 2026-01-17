package models

import (
	"time"

	"gorm.io/gorm"
)

// Location represents a work location in the system
type Location struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TenantID    *uint          `json:"tenant_id,omitempty" gorm:"index"`
	OrganizationID *uint       `json:"organization_id,omitempty" gorm:"index"` // References organizations(id)
	Name        string         `json:"name" gorm:"not null;size:100"`
	LocationType *string       `json:"location_type,omitempty" gorm:"size:50"` // head_office, branch, remote, satellite
	AddressLine1 *string       `json:"address_line1,omitempty" gorm:"type:text"`
	AddressLine2 *string       `json:"address_line2,omitempty" gorm:"type:text"`
	City        *string        `json:"city,omitempty" gorm:"size:100"`
	State       *string        `json:"state,omitempty" gorm:"size:100"`
	Country     *string        `json:"country,omitempty" gorm:"size:100"`
	PostalCode  *string        `json:"postal_code,omitempty" gorm:"size:20"`
	Latitude    *float64       `json:"latitude,omitempty" gorm:"type:decimal(10,8)"`
	Longitude   *float64       `json:"longitude,omitempty" gorm:"type:decimal(11,8)"`
	Timezone    *string        `json:"timezone,omitempty" gorm:"size:100;default:'Africa/Dar_es_Salaam'"`
	PhoneNumber *string        `json:"phone_number,omitempty" gorm:"size:50"`
	Capacity    *int           `json:"capacity,omitempty"` // Max employees
	Facilities  *string        `json:"facilities,omitempty" gorm:"type:text"` // JSON array of available facilities
	IsHeadOffice bool          `json:"is_head_office" gorm:"default:false"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	UpdatedBy    *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for Location model
func (Location) TableName() string {
	return "locations"
}

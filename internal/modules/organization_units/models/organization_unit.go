package models

import (
	"time"

	"gorm.io/gorm"
)

// OrganizationUnitType represents the type of organization unit
type OrganizationUnitType string

const (
	UnitTypeDivision     OrganizationUnitType = "division"
	UnitTypeBusinessUnit OrganizationUnitType = "business_unit"
	UnitTypeSubsidiary   OrganizationUnitType = "subsidiary"
	UnitTypeBranch       OrganizationUnitType = "branch"
)

// OrganizationUnit represents a business unit/division within an organization
type OrganizationUnit struct {
	ID              uint                 `json:"id" gorm:"primaryKey"`
	TenantID        *uint                `json:"tenant_id,omitempty" gorm:"index"`
	OrganizationID  *uint                `json:"organization_id,omitempty" gorm:"index"` // References organizations(id)
	Code            string               `json:"code" gorm:"uniqueIndex;not null;size:20"`
	Name            string               `json:"name" gorm:"not null;size:100"`
	Description     *string               `json:"description,omitempty" gorm:"type:text"`
	UnitType        *string              `json:"unit_type,omitempty" gorm:"size:50"` // division, business_unit, subsidiary, branch
	ParentUnitID    *uint                `json:"parent_unit_id,omitempty" gorm:"index"` // Self-referential
	HeadID          *uint                `json:"head_id,omitempty" gorm:"index"` // References employees(id)
	BudgetAllocated *float64             `json:"budget_allocated,omitempty" gorm:"type:decimal(15,2)"`
	BudgetCurrency  *string              `json:"budget_currency,omitempty" gorm:"size:3;default:'TZS'"`
	LocationID      *uint                `json:"location_id,omitempty" gorm:"index"` // References locations(id)
	IsActive        bool                 `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	UpdatedBy       *uint                `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt       gorm.DeletedAt       `json:"-" gorm:"index"`
	
	// Self-referential relationship
	ParentUnit   *OrganizationUnit `json:"parent_unit,omitempty" gorm:"foreignKey:ParentUnitID"`
	SubUnits     []OrganizationUnit `json:"sub_units,omitempty" gorm:"foreignKey:ParentUnitID"`
}

// TableName specifies the table name for OrganizationUnit model
func (OrganizationUnit) TableName() string {
	return "organization_units"
}

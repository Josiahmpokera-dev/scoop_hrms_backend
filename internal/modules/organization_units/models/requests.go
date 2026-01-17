package models

// CreateOrganizationUnitRequest represents the request to create an organization unit
type CreateOrganizationUnitRequest struct {
	OrganizationID  *uint    `json:"organization_id,omitempty"`
	Code            string   `json:"code" binding:"required,min=2,max=20"`
	Name            string   `json:"name" binding:"required,min=2,max=100"`
	Description     *string  `json:"description,omitempty"`
	UnitType        *string  `json:"unit_type,omitempty" binding:"omitempty,oneof=division business_unit subsidiary branch"`
	ParentUnitID    *uint    `json:"parent_unit_id,omitempty"`
	HeadID          *uint    `json:"head_id,omitempty"`
	BudgetAllocated *float64 `json:"budget_allocated,omitempty"`
	BudgetCurrency  *string  `json:"budget_currency,omitempty" binding:"omitempty,len=3"`
	LocationID      *uint    `json:"location_id,omitempty"`
	IsActive        *bool    `json:"is_active,omitempty"`
}

// UpdateOrganizationUnitRequest represents the request to update an organization unit
type UpdateOrganizationUnitRequest struct {
	OrganizationID  *uint    `json:"organization_id,omitempty"`
	Code            *string  `json:"code,omitempty" binding:"omitempty,min=2,max=20"`
	Name            *string  `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Description     *string  `json:"description,omitempty"`
	UnitType        *string  `json:"unit_type,omitempty" binding:"omitempty,oneof=division business_unit subsidiary branch"`
	ParentUnitID    *uint    `json:"parent_unit_id,omitempty"`
	HeadID          *uint    `json:"head_id,omitempty"`
	BudgetAllocated *float64 `json:"budget_allocated,omitempty"`
	BudgetCurrency  *string  `json:"budget_currency,omitempty" binding:"omitempty,len=3"`
	LocationID      *uint    `json:"location_id,omitempty"`
	IsActive        *bool    `json:"is_active,omitempty"`
}

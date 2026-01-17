package models

// CreateCostCenterRequest represents the request to create a cost center
type CreateCostCenterRequest struct {
	OrganizationID     *uint    `json:"organization_id,omitempty"`
	Code               string   `json:"code" binding:"required,min=2,max=20"`
	Name               string   `json:"name" binding:"required,min=2,max=100"`
	Description        *string  `json:"description,omitempty"`
	CostCenterType     *string  `json:"cost_center_type,omitempty" binding:"omitempty,oneof=department project overhead"`
	ParentCostCenterID *uint    `json:"parent_cost_center_id,omitempty"`
	BudgetAllocated    *float64 `json:"budget_allocated,omitempty"`
	BudgetCurrency     *string  `json:"budget_currency,omitempty" binding:"omitempty,len=3"`
	ManagerID          *uint    `json:"manager_id,omitempty"`
	DepartmentID       *uint    `json:"department_id,omitempty"`
	IsActive           *bool   `json:"is_active,omitempty"`
}

// UpdateCostCenterRequest represents the request to update a cost center
type UpdateCostCenterRequest struct {
	OrganizationID     *uint    `json:"organization_id,omitempty"`
	Code               *string  `json:"code,omitempty" binding:"omitempty,min=2,max=20"`
	Name               *string  `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Description        *string  `json:"description,omitempty"`
	CostCenterType     *string  `json:"cost_center_type,omitempty" binding:"omitempty,oneof=department project overhead"`
	ParentCostCenterID *uint    `json:"parent_cost_center_id,omitempty"`
	BudgetAllocated    *float64 `json:"budget_allocated,omitempty"`
	BudgetCurrency     *string  `json:"budget_currency,omitempty" binding:"omitempty,len=3"`
	ManagerID          *uint    `json:"manager_id,omitempty"`
	DepartmentID       *uint    `json:"department_id,omitempty"`
	IsActive           *bool    `json:"is_active,omitempty"`
}

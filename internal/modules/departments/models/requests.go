package models

import (
	"encoding/json"
	"fmt"
)

// CreateDepartmentRequest represents the request to create a department
type CreateDepartmentRequest struct {
	OrganizationID     *uint    `json:"organization_id,omitempty"`                                                       // Organization this department belongs to
	OrganizationUnitID *uint    `json:"organization_unit_id,omitempty"`                                                  // Business unit/division
	Code               string   `json:"code" binding:"required,min=2,max=20"`                                            // Department code (required)
	Level              *string  `json:"level,omitempty" binding:"omitempty,oneof=company business_unit department team"` // Level: company, business_unit, department, team
	Name               string   `json:"name" binding:"required,min=2,max=100"`
	Description        *string  `json:"description,omitempty"`
	DepartmentType     *string  `json:"department_type,omitempty" binding:"omitempty,oneof=core support operational strategic"`
	ParentDepartmentID *uint    `json:"parent_department_id,omitempty"` // Parent Department (none for top level, default if not provided)
	ManagerID          *uint    `json:"manager_id,omitempty"`           // Head of Department
	DeputyManager      *string  `json:"deputy_manager,omitempty"`       // Deputy Manager name (e.g., "John Doe")
	LocationID         *uint    `json:"location_id,omitempty"`          // Location ID (references locations table)
	Location           *string  `json:"location,omitempty"`             // Location name (alternative to location_id)
	CostCenter         *string  `json:"cost_center,omitempty"`          // Cost Center name
	BudgetAllocated    *float64 `json:"budget_allocated,omitempty"`     // Budget
	BudgetCurrency     *string  `json:"budget_currency,omitempty" binding:"omitempty,len=3"`
	EmployeeCapacity   *int     `json:"employee_capacity,omitempty"`
	ShiftIDs           []uint   `json:"shift_ids,omitempty"` // List of shift IDs to associate with this department
	IsActive           *bool    `json:"is_active,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle location as both number (ID) and string (name)
func (r *CreateDepartmentRequest) UnmarshalJSON(data []byte) error {
	// Define a temporary struct with location as interface{} to handle both types
	type Alias CreateDepartmentRequest
	aux := &struct {
		Location interface{} `json:"location,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle location field - can be number (ID) or string (name)
	if aux.Location != nil {
		switch v := aux.Location.(type) {
		case float64:
			// If location is a number, convert it to location_id
			locationID := uint(v)
			r.LocationID = &locationID
			r.Location = nil
		case string:
			// If location is a string, use it as location name
			r.Location = &v
		case nil:
			// If location is null, leave both fields nil
			r.Location = nil
		default:
			return fmt.Errorf("location must be either a number (ID) or string (name), got %T", v)
		}
	}

	return nil
}

// UpdateDepartmentRequest represents the request to update a department
type UpdateDepartmentRequest struct {
	OrganizationID     *uint    `json:"organization_id,omitempty"`
	OrganizationUnitID *uint    `json:"organization_unit_id,omitempty"`
	Code               *string  `json:"code,omitempty" binding:"omitempty,min=2,max=20"`
	Level              *string  `json:"level,omitempty" binding:"omitempty,oneof=company business_unit department team"`
	Name               *string  `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Description        *string  `json:"description,omitempty"`
	DepartmentType     *string  `json:"department_type,omitempty" binding:"omitempty,oneof=core support operational strategic"`
	ParentDepartmentID *uint    `json:"parent_department_id,omitempty"`
	ManagerID          *uint    `json:"manager_id,omitempty"`       // Head of Department
	DeputyManager      *string  `json:"deputy_manager,omitempty"`   // Deputy Manager name (e.g., "John Doe")
	LocationID         *uint    `json:"location_id,omitempty"`      // Location ID (references locations table)
	Location           *string  `json:"location,omitempty"`         // Location name (alternative to location_id)
	CostCenter         *string  `json:"cost_center,omitempty"`      // Cost Center name
	BudgetAllocated    *float64 `json:"budget_allocated,omitempty"` // Budget
	BudgetCurrency     *string  `json:"budget_currency,omitempty" binding:"omitempty,len=3"`
	EmployeeCapacity   *int     `json:"employee_capacity,omitempty"`
	ShiftIDs           []uint   `json:"shift_ids,omitempty"` // List of shift IDs to update
	IsActive           *bool    `json:"is_active,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle location as both number (ID) and string (name)
func (r *UpdateDepartmentRequest) UnmarshalJSON(data []byte) error {
	// Define a temporary struct with location as interface{} to handle both types
	type Alias UpdateDepartmentRequest
	aux := &struct {
		Location interface{} `json:"location,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle location field - can be number (ID) or string (name)
	if aux.Location != nil {
		switch v := aux.Location.(type) {
		case float64:
			// If location is a number, convert it to location_id
			locationID := uint(v)
			r.LocationID = &locationID
			r.Location = nil
		case string:
			// If location is a string, use it as location name
			r.Location = &v
		case nil:
			// If location is null, leave both fields nil
			r.Location = nil
		default:
			return fmt.Errorf("location must be either a number (ID) or string (name), got %T", v)
		}
	}

	return nil
}

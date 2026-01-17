package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// FlexibleEmployeeID can accept either a number (database ID) or string (employee_id)
type FlexibleEmployeeID struct {
	IsNumeric bool
	NumericID uint
	StringID  string
}

// UnmarshalJSON implements json.Unmarshaler to accept both number and string
func (f *FlexibleEmployeeID) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as number first
	var numID uint
	if err := json.Unmarshal(data, &numID); err == nil {
		f.IsNumeric = true
		f.NumericID = numID
		return nil
	}

	// Try to unmarshal as string
	var strID string
	if err := json.Unmarshal(data, &strID); err == nil {
		f.IsNumeric = false
		f.StringID = strID
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into FlexibleEmployeeID (expected number or string)", string(data))
}

// String returns the employee identifier as string
func (f *FlexibleEmployeeID) String() string {
	if f.IsNumeric {
		return strconv.FormatUint(uint64(f.NumericID), 10)
	}
	return f.StringID
}

// CreateAssetRequest represents the request to create a new asset
type CreateAssetRequest struct {
	AssetType            string             `json:"asset_type" binding:"required"` // laptop, mobile_phone, etc.
	Brand                string             `json:"brand" binding:"required,max=100"`
	Model                string             `json:"model" binding:"required,max=100"`
	SerialNumber         string             `json:"serial_number" binding:"required,max=100"`
	AssetCode            *string            `json:"asset_code,omitempty"` // Optional - auto-generated if not provided
	PurchaseDate         *string            `json:"purchase_date,omitempty"` // YYYY-MM-DD format
	Value                 *float64           `json:"value,omitempty"`
	WarrantyExpiry        *string            `json:"warranty_expiry,omitempty"` // YYYY-MM-DD format
	Condition             *string            `json:"condition,omitempty"` // excellent, good, fair, poor
	AssignedToEmployeeID *FlexibleEmployeeID `json:"assigned_to_employee_id,omitempty"` // Employee ID (number or string) if assigning immediately
	Notes                 *string            `json:"notes,omitempty"`
}

// UpdateAssetRequest represents the request to update an asset
type UpdateAssetRequest struct {
	ID            uint     `json:"id" binding:"required"`
	AssetType     *string  `json:"asset_type,omitempty"`
	Brand         *string  `json:"brand,omitempty" binding:"omitempty,max=100"`
	Model         *string  `json:"model,omitempty" binding:"omitempty,max=100"`
	SerialNumber  *string  `json:"serial_number,omitempty" binding:"omitempty,max=100"`
	Condition     *string  `json:"condition,omitempty"`
	Value         *float64 `json:"value,omitempty"`
	PurchaseDate  *string  `json:"purchase_date,omitempty"` // YYYY-MM-DD format
	WarrantyExpiry *string  `json:"warranty_expiry,omitempty"` // YYYY-MM-DD format
	Notes          *string  `json:"notes,omitempty"`
}

// AssignAssetRequest represents the request to assign an asset to an employee
type AssignAssetRequest struct {
	AssetID      uint     `json:"asset_id" binding:"required"`
	EmployeeID   string   `json:"employee_id" binding:"required"`
	AssignedDate *string  `json:"assigned_date,omitempty"` // YYYY-MM-DD format, defaults to today
	Notes        *string  `json:"notes,omitempty"`
}

// ReturnAssetRequest represents the request to return an asset
type ReturnAssetRequest struct {
	AssetID    uint     `json:"asset_id" binding:"required"`
	ReturnDate *string  `json:"return_date,omitempty"` // YYYY-MM-DD format, defaults to today
	Condition  *string  `json:"condition,omitempty"` // Update condition on return
	Notes      *string  `json:"notes,omitempty"`
}

// ReassignAssetRequest represents the request to reassign an asset to another employee
type ReassignAssetRequest struct {
	AssetID      uint     `json:"asset_id" binding:"required"`
	NewEmployeeID string  `json:"new_employee_id" binding:"required"` // Employee ID to reassign to
	ReassignDate *string  `json:"reassign_date,omitempty"` // YYYY-MM-DD format, defaults to today
	ReturnDate   *string  `json:"return_date,omitempty"` // YYYY-MM-DD format for previous employee return (defaults to reassign_date)
	Condition    *string  `json:"condition,omitempty"` // Update condition on reassignment
	Notes         *string  `json:"notes,omitempty"` // Notes about the reassignment
}

// MarkForRepairRequest represents the request to mark an asset for repair
type MarkForRepairRequest struct {
	AssetID      uint    `json:"asset_id" binding:"required"`
	RepairReason *string `json:"repair_reason,omitempty"`
	RepairNotes  *string `json:"repair_notes,omitempty"`
}

// CompleteRepairRequest represents the request to complete asset repair
type CompleteRepairRequest struct {
	AssetID    uint     `json:"asset_id" binding:"required"`
	RepairCost *float64 `json:"repair_cost,omitempty"`
	Condition  *string  `json:"condition,omitempty"` // Update condition after repair
	Notes      *string  `json:"notes,omitempty"`
}

// RetireAssetRequest represents the request to retire an asset
type RetireAssetRequest struct {
	AssetID          uint    `json:"asset_id" binding:"required"`
	RetirementReason *string `json:"retirement_reason,omitempty"`
	RetirementDate   *string `json:"retirement_date,omitempty"` // YYYY-MM-DD format, defaults to today
	Notes            *string `json:"notes,omitempty"`
}

// DeleteAssetRequest represents the request to delete an asset
type DeleteAssetRequest struct {
	ID uint `json:"id" binding:"required"`
}

// GetAssetRequest represents the request to get an asset by ID
type GetAssetRequest struct {
	ID uint `json:"id" binding:"required"`
}

// ListAssetsRequest represents the request to list assets with filters
type ListAssetsRequest struct {
	Page       int     `form:"page" binding:"omitempty,min=1"`
	PageSize   int     `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search     *string `form:"search"`
	Status     *string `form:"status"` // in_use, available, under_repair, retired
	AssetType  *string `form:"asset_type"`
	AssignedTo *string `form:"assigned_to"` // Employee ID
	Department *string `form:"department"`
}

// AssetResponse represents the response for an asset
type AssetResponse struct {
	ID                  uint       `json:"id"`
	AssetCode           string     `json:"asset_code"`
	AssetType           string     `json:"asset_type"`
	Brand               string     `json:"brand"`
	Model               string     `json:"model"`
	SerialNumber        string     `json:"serial_number"`
	AssignedTo          *string    `json:"assigned_to,omitempty"` // Employee name
	EmployeeID          *string    `json:"employee_id,omitempty"`
	EmployeePhoto       *string    `json:"employee_photo,omitempty"`
	Department          *string    `json:"department,omitempty"`
	AssignedDate        *time.Time `json:"assigned_date,omitempty"`
	ReturnDate          *time.Time `json:"return_date,omitempty"`
	Status              string     `json:"status"`
	Condition           string     `json:"condition"`
	PurchaseDate        *time.Time `json:"purchase_date,omitempty"`
	WarrantyExpiry      *time.Time `json:"warranty_expiry,omitempty"`
	Value               *float64   `json:"value,omitempty"`
	Notes               *string    `json:"notes,omitempty"`
	RepairReason        *string    `json:"repair_reason,omitempty"`
	RepairNotes         *string    `json:"repair_notes,omitempty"`
	RepairDate          *time.Time `json:"repair_date,omitempty"`
	RepairCost          *float64   `json:"repair_cost,omitempty"`
	RepairCompletedDate *time.Time `json:"repair_completed_date,omitempty"`
	RetirementReason    *string    `json:"retirement_reason,omitempty"`
	RetirementDate      *time.Time `json:"retirement_date,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// AssetTypeResponse represents an asset type
type AssetTypeResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

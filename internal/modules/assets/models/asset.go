package models

import (
	"time"

	"gorm.io/gorm"
)

// AssetStatus represents the status of an asset
type AssetStatus string

const (
	AssetStatusAvailable   AssetStatus = "available"
	AssetStatusInUse       AssetStatus = "in_use"
	AssetStatusUnderRepair AssetStatus = "under_repair"
	AssetStatusRetired     AssetStatus = "retired"
)

// AssetType represents the type of asset
type AssetType string

const (
	AssetTypeLaptop         AssetType = "laptop"
	AssetTypeMobilePhone     AssetType = "mobile_phone"
	AssetTypeDesktop         AssetType = "desktop"
	AssetTypeTablet          AssetType = "tablet"
	AssetTypeIDCard          AssetType = "id_card"
	AssetTypeAccessCard      AssetType = "access_card"
	AssetTypeVehicle         AssetType = "vehicle"
	AssetTypeOfficeEquipment AssetType = "office_equipment"
)

// AssetCondition represents the condition of an asset
type AssetCondition string

const (
	AssetConditionExcellent AssetCondition = "excellent"
	AssetConditionGood      AssetCondition = "good"
	AssetConditionFair      AssetCondition = "fair"
	AssetConditionPoor      AssetCondition = "poor"
)

// Asset represents an asset/equipment in the system
type Asset struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	AssetCode         string         `json:"asset_code" gorm:"uniqueIndex;not null;size:50"`
	AssetType         string         `json:"asset_type" gorm:"not null;size:50"`
	Brand             string         `json:"brand" gorm:"not null;size:100"`
	Model             string         `json:"model" gorm:"not null;size:100"`
	SerialNumber      string         `json:"serial_number" gorm:"uniqueIndex;not null;size:100"`
	AssignedTo        *string        `json:"assigned_to,omitempty" gorm:"size:255"` // Employee name
	EmployeeID        *string        `json:"employee_id,omitempty" gorm:"size:50"`  // Employee ID (e.g., EMP001)
	Department        *string        `json:"department,omitempty" gorm:"size:100"`
	AssignedDate      *time.Time     `json:"assigned_date,omitempty"`
	ReturnDate        *time.Time     `json:"return_date,omitempty"`
	Status            string         `json:"status" gorm:"default:'available';size:20"`
	Condition         string         `json:"condition" gorm:"default:'excellent';size:20"`
	PurchaseDate      *time.Time     `json:"purchase_date,omitempty"`
	WarrantyExpiry    *time.Time     `json:"warranty_expiry,omitempty"`
	Value             *float64       `json:"value,omitempty" gorm:"type:decimal(10,2)"`
	Notes             *string        `json:"notes,omitempty" gorm:"type:text"`
	RepairReason      *string        `json:"repair_reason,omitempty" gorm:"type:text"`
	RepairNotes       *string        `json:"repair_notes,omitempty" gorm:"type:text"`
	RepairDate        *time.Time     `json:"repair_date,omitempty"`
	RepairCost        *float64       `json:"repair_cost,omitempty" gorm:"type:decimal(10,2)"`
	RepairCompletedDate *time.Time   `json:"repair_completed_date,omitempty"`
	RetirementReason  *string        `json:"retirement_reason,omitempty" gorm:"type:text"`
	RetirementDate    *time.Time     `json:"retirement_date,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	UpdatedBy         *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for Asset model
func (Asset) TableName() string {
	return "assets"
}

// IsAvailable checks if asset is available for assignment
func (a *Asset) IsAvailable() bool {
	return a.Status == string(AssetStatusAvailable)
}

// IsAssigned checks if asset is currently assigned
func (a *Asset) IsAssigned() bool {
	return a.Status == string(AssetStatusInUse) && a.EmployeeID != nil
}

// CanBeAssigned checks if asset can be assigned
func (a *Asset) CanBeAssigned() bool {
	return a.Status == string(AssetStatusAvailable) && a.DeletedAt.Time.IsZero()
}

// CanBeReturned checks if asset can be returned
func (a *Asset) CanBeReturned() bool {
	return a.Status == string(AssetStatusInUse) && a.EmployeeID != nil
}

// CanBeRepaired checks if asset can be marked for repair
func (a *Asset) CanBeRepaired() bool {
	return a.Status != string(AssetStatusUnderRepair) && a.Status != string(AssetStatusRetired)
}

// CanBeRetired checks if asset can be retired
func (a *Asset) CanBeRetired() bool {
	return a.Status != string(AssetStatusRetired)
}

// CanBeDeleted checks if asset can be deleted
func (a *Asset) CanBeDeleted() bool {
	return a.Status == string(AssetStatusRetired)
}

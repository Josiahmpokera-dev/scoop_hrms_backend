package models

import (
	"time"

	"gorm.io/gorm"
)

// CostCenterType represents the type of cost center
type CostCenterType string

const (
	CostCenterTypeDepartment CostCenterType = "department"
	CostCenterTypeProject    CostCenterType = "project"
	CostCenterTypeOverhead   CostCenterType = "overhead"
)

// CostCenter represents a cost center in the system
type CostCenter struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	TenantID          *uint          `json:"tenant_id,omitempty" gorm:"index"`
	OrganizationID    *uint          `json:"organization_id,omitempty" gorm:"index"` // References organizations(id)
	Code              string         `json:"code" gorm:"uniqueIndex;not null;size:20"`
	Name              string         `json:"name" gorm:"not null;size:100"`
	Description       *string        `json:"description,omitempty" gorm:"type:text"`
	CostCenterType    *string        `json:"cost_center_type,omitempty" gorm:"size:50"` // department, project, overhead, etc.
	ParentCostCenterID *uint         `json:"parent_cost_center_id,omitempty" gorm:"index"` // Self-referential
	BudgetAllocated   *float64       `json:"budget_allocated,omitempty" gorm:"type:decimal(15,2)"`
	BudgetCurrency    *string        `json:"budget_currency,omitempty" gorm:"size:3;default:'TZS'"`
	ManagerID         *uint          `json:"manager_id,omitempty" gorm:"index"` // References employees(id)
	DepartmentID      *uint          `json:"department_id,omitempty" gorm:"index"` // References departments(id)
	IsActive          bool           `json:"is_active" gorm:"default:true"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	UpdatedBy         *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Self-referential relationship
	ParentCostCenter *CostCenter `json:"parent_cost_center,omitempty" gorm:"foreignKey:ParentCostCenterID"`
	SubCostCenters  []CostCenter `json:"sub_cost_centers,omitempty" gorm:"foreignKey:ParentCostCenterID"`
}

// TableName specifies the table name for CostCenter model
func (CostCenter) TableName() string {
	return "cost_centers"
}

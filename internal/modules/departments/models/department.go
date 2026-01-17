package models

import (
	"time"

	"gorm.io/gorm"
)

// DepartmentLevel represents the hierarchical level of a department
type DepartmentLevel string

const (
	DepartmentLevelCompany      DepartmentLevel = "company"
	DepartmentLevelBusinessUnit  DepartmentLevel = "business_unit"
	DepartmentLevelDepartment    DepartmentLevel = "department"
	DepartmentLevelTeam          DepartmentLevel = "team"
)

// Department represents a department in the organizational hierarchy
type Department struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	TenantID          *uint          `json:"tenant_id,omitempty" gorm:"index"`
	OrganizationID    *uint          `json:"organization_id,omitempty" gorm:"index"` // References organizations(id)
	OrganizationUnitID *uint         `json:"organization_unit_id,omitempty" gorm:"index"` // References organization_units(id)
	Code              string         `json:"code" gorm:"uniqueIndex;not null;size:20"`
	Name              string         `json:"name" gorm:"not null;size:100"`
	Description       *string        `json:"description,omitempty" gorm:"type:text"`
	Level             *string        `json:"level,omitempty" gorm:"size:50"` // company, business_unit, department, team
	DepartmentType    *string        `json:"department_type,omitempty" gorm:"size:50"` // core, support, operational, strategic
	ParentDepartmentID *uint         `json:"parent_department_id,omitempty" gorm:"index"`
	ManagerID         *uint          `json:"manager_id,omitempty" gorm:"index"` // Head of Department - References employees(id)
	DeputyManager     *string        `json:"deputy_manager,omitempty" gorm:"size:100"` // Deputy Manager name as string (not a reference)
	BudgetAllocated   *float64       `json:"budget_allocated,omitempty" gorm:"type:decimal(15,2)"`
	BudgetCurrency    *string        `json:"budget_currency,omitempty" gorm:"size:3;default:'TZS'"`
	EmployeeCapacity  *int           `json:"employee_capacity,omitempty"`
	LocationID        *uint          `json:"location_id,omitempty" gorm:"index"` // References locations(id)
	CostCenter        *string        `json:"cost_center,omitempty" gorm:"size:100"` // Cost center as string (not a reference)
	IsActive          bool           `json:"is_active" gorm:"default:true"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	UpdatedBy         *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Self-referential relationship for parent department
	ParentDepartment *Department `json:"parent_department,omitempty" gorm:"foreignKey:ParentDepartmentID"`
	// Sub-departments (children)
	SubDepartments   []Department `json:"sub_departments,omitempty" gorm:"foreignKey:ParentDepartmentID"`
}

// TableName specifies the table name for Department model
func (Department) TableName() string {
	return "departments"
}

// IsRoot checks if department is a root department (no parent)
func (d *Department) IsRoot() bool {
	return d.ParentDepartmentID == nil
}

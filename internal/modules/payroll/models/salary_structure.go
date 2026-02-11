package models

import (
	"time"

	"gorm.io/gorm"
)

// SalaryStructure represents a salary template
type SalaryStructure struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	TemplateName          string         `json:"templateName" gorm:"not null;size:255"`
	Grade                 string         `json:"grade" gorm:"size:50"` // e.g., G1, G2, G3
	Location              string         `json:"location" gorm:"size:100"`
	Country               string         `json:"country" gorm:"size:100;default:'Tanzania'"`
	CTC                   float64        `json:"ctc" gorm:"type:decimal(15,2)"` // Cost to Company
	Basic                 float64        `json:"basic" gorm:"type:decimal(15,2)"`
	HRA                   float64        `json:"hra" gorm:"type:decimal(15,2)"`        // House Rent Allowance
	Transport             float64        `json:"transport" gorm:"type:decimal(15,2)"` // Transport Allowance
	Medical               float64        `json:"medical" gorm:"type:decimal(15,2)"`
	OtherAllowances       float64        `json:"otherAllowances" gorm:"type:decimal(15,2)"`
	GrossSalary           float64        `json:"grossSalary" gorm:"type:decimal(15,2)"` // Calculated
	// Deductions (for display)
	PAYEDeduction         float64        `json:"payeDeduction" gorm:"type:decimal(15,2)"`
	NSSFEmployee          float64        `json:"nssfEmployee" gorm:"type:decimal(15,2)"`
	NHIFEmployee          float64        `json:"nhifEmployee" gorm:"type:decimal(15,2)"`
	TotalDeductions       float64        `json:"totalDeductions" gorm:"type:decimal(15,2)"` // Calculated
	// Employer contributions (for display)
	NSSFEmployer          float64        `json:"nssfEmployer" gorm:"type:decimal(15,2)"`
	NHIFEmployer          float64        `json:"nhifEmployer" gorm:"type:decimal(15,2)"`
	SDL                   float64        `json:"sdl" gorm:"type:decimal(15,2)"`
	WCF                   float64        `json:"wcf" gorm:"type:decimal(15,2)"`
	TotalEmployerContribs float64        `json:"totalEmployerContributions" gorm:"type:decimal(15,2)"`
	NetPay                float64        `json:"netPay" gorm:"type:decimal(15,2)"` // Calculated
	IsActive              bool           `json:"isActive" gorm:"default:true"`
	EmployeesCount        int            `json:"employeesCount" gorm:"-"` // Virtual field
	// Timestamps
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             time.Time      `json:"updatedAt"`
	UpdatedByID           *uint          `json:"updatedById,omitempty"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (SalaryStructure) TableName() string {
	return "salary_structures"
}

// ComponentType represents the type of salary component
type ComponentType string

const (
	ComponentTypeEarning              ComponentType = "Earning"
	ComponentTypeDeduction            ComponentType = "Deduction"
	ComponentTypeEmployerContribution ComponentType = "Employer Contribution"
)

// CalculationType represents how a component is calculated
type CalculationType string

const (
	CalculationTypeFixed            CalculationType = "Fixed"
	CalculationTypePercentageBasic  CalculationType = "Percentage of Basic"
	CalculationTypePercentageGross  CalculationType = "Percentage of Gross"
	CalculationTypeFormula          CalculationType = "Formula"
)

// SalaryComponent represents a reusable salary component
type SalaryComponent struct {
	ID                uint            `json:"id" gorm:"primaryKey"`
	ComponentName     string          `json:"componentName" gorm:"not null;size:255"`
	ComponentCode     string          `json:"componentCode" gorm:"not null;size:50;uniqueIndex"` // e.g., BASIC, HRA, PAYE
	ComponentType     ComponentType   `json:"componentType" gorm:"type:varchar(30);not null"`
	CalculationType   CalculationType `json:"calculationType" gorm:"type:varchar(30);not null"`
	DefaultAmount     *float64        `json:"defaultAmount,omitempty" gorm:"type:decimal(15,2)"`
	DefaultPercentage *float64        `json:"defaultPercentage,omitempty" gorm:"type:decimal(5,2)"`
	Formula           *string         `json:"formula,omitempty" gorm:"size:500"` // e.g., PAYE_SLAB_CALCULATION
	IsTaxable         bool            `json:"isTaxable" gorm:"default:true"`
	IsStatutory       bool            `json:"isStatutory" gorm:"default:false"`
	IsRecurring       bool            `json:"isRecurring" gorm:"default:true"`
	ApplicableFor     string          `json:"applicableFor" gorm:"size:500;default:'[\"All\"]'"` // JSON array
	DisplayInPayslip  bool            `json:"displayInPayslip" gorm:"default:true"`
	Country           string          `json:"country" gorm:"size:100;default:'Tanzania'"`
	IsActive          bool            `json:"isActive" gorm:"default:true"`
	SortOrder         int             `json:"sortOrder" gorm:"default:0"`
	// Timestamps
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt  `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (SalaryComponent) TableName() string {
	return "salary_components"
}

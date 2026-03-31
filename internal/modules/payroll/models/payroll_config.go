package models

import (
	"time"
)

// PayrollConfigCategory represents the category of payroll configuration
type PayrollConfigCategory string

const (
	ConfigCategoryNSSF    PayrollConfigCategory = "NSSF"
	ConfigCategoryPAYE    PayrollConfigCategory = "PAYE"
	ConfigCategorySDL     PayrollConfigCategory = "SDL"
	ConfigCategoryWCF     PayrollConfigCategory = "WCF"
	ConfigCategoryHESLB   PayrollConfigCategory = "HESLB"
	ConfigCategoryFormula PayrollConfigCategory = "FORMULA"
	ConfigCategoryGeneral PayrollConfigCategory = "GENERAL"
)

// PayrollConfig represents payroll configuration settings
type PayrollConfig struct {
	ID          uint                  `json:"id" gorm:"primaryKey"`
	ConfigKey   string                `json:"config_key" gorm:"uniqueIndex;not null;size:100"`
	Value       string                `json:"value" gorm:"not null;type:text"`
	Label       string                `json:"label" gorm:"not null;size:255"`
	Description string                `json:"description" gorm:"type:text"`
	Category    PayrollConfigCategory `json:"category" gorm:"type:varchar(30);not null"`
	Editable    bool                  `json:"editable" gorm:"default:false"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// TableName specifies the table name
func (PayrollConfig) TableName() string {
	return "payroll_config"
}

// PAYEBand represents a PAYE tax band
type PAYEBand struct {
	Band    int      `json:"band"`
	Min     float64  `json:"min"`
	Max     *float64 `json:"max,omitempty"`
	Rate    float64  `json:"rate"`
	BaseTax float64  `json:"base_tax"`
	Label   string   `json:"label"`
}

// PAYECalculationResult represents PAYE calculation result
type PAYECalculationResult struct {
	TaxablePay  float64       `json:"taxable_pay"`
	BandApplied int           `json:"band_applied"`
	BandLabel   string        `json:"band_label"`
	PAYETax     float64       `json:"paye_tax"`
	Breakdown   PAYEBreakdown `json:"breakdown"`
}

// PAYEBreakdown represents PAYE calculation breakdown
type PAYEBreakdown struct {
	BaseTax      float64 `json:"base_tax"`
	Rate         float64 `json:"rate"`
	ExcessAmount float64 `json:"excess_amount"`
	RateTax      float64 `json:"rate_tax"`
}

// NSSFCalculationResult represents NSSF calculation result
type NSSFCalculationResult struct {
	GrossSalary  float64 `json:"gross_salary"`
	NSSEmployee  float64 `json:"nssf_employee"`
	NSSEmployer  float64 `json:"nssf_employer"`
	NSSTotal     float64 `json:"nssf_total"`
	RateEmployee float64 `json:"rate_employee"`
	RateEmployer float64 `json:"rate_employer"`
}

// CTCCalculationResult represents Cost to Company calculation result
type CTCCalculationResult struct {
	GrossSalary float64 `json:"gross_salary"`
	NSSEmployer float64 `json:"nssf_employer"`
	WCF         float64 `json:"wcf"`
	SDL         float64 `json:"sdl"`
	TotalCTC    float64 `json:"total_ctc"`
}

package services

import (
	"math"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/repositories"
)

// TaxCalculator handles all statutory calculations for Tanzania
type TaxCalculator struct {
	complianceRepo *repositories.ComplianceRepository
}

// NewTaxCalculator creates a new tax calculator instance
func NewTaxCalculator() *TaxCalculator {
	return &TaxCalculator{
		complianceRepo: repositories.NewComplianceRepository(),
	}
}

// TaxCalculationResult holds the result of tax calculations
type TaxCalculationResult struct {
	TaxableIncome        float64 `json:"taxableIncome"`
	PAYE                 float64 `json:"paye"`
	NSSFEmployee         float64 `json:"nssfEmployee"`
	NSSFEmployer         float64 `json:"nssfEmployer"`
	NHIFEmployee         float64 `json:"nhifEmployee"`
	NHIFEmployer         float64 `json:"nhifEmployer"`
	SDL                  float64 `json:"sdl"`
	WCF                  float64 `json:"wcf"`
	TotalEmployeeDeductions float64 `json:"totalEmployeeDeductions"`
	TotalEmployerContributions float64 `json:"totalEmployerContributions"`
	NetPay               float64 `json:"netPay"`
}

// CalculatePAYE calculates Tanzania PAYE tax based on gross salary
func (tc *TaxCalculator) CalculatePAYE(grossSalary float64, taxYear int, tenantID *uint) (float64, error) {
	slabs, err := tc.complianceRepo.GetTaxSlabs("Tanzania", taxYear, tenantID)
	if err != nil || len(slabs) == 0 {
		// Fallback to hardcoded 2024 rates if no slabs found
		return tc.calculatePAYEFallback(grossSalary), nil
	}

	return tc.calculatePAYEFromSlabs(grossSalary, slabs), nil
}

// calculatePAYEFromSlabs calculates PAYE using database slabs
func (tc *TaxCalculator) calculatePAYEFromSlabs(grossSalary float64, slabs []models.TaxSlab) float64 {
	if grossSalary <= 0 {
		return 0
	}

	for _, slab := range slabs {
		if grossSalary >= slab.SlabFrom && grossSalary <= slab.SlabTo {
			taxableAmount := grossSalary - slab.SlabFrom + 1
			tax := slab.FixedAmount + (taxableAmount * slab.TaxRate / 100)
			return math.Round(tax*100) / 100
		}
	}

	// If above all slabs, use highest slab
	if len(slabs) > 0 {
		highestSlab := slabs[len(slabs)-1]
		taxableAmount := grossSalary - highestSlab.SlabFrom + 1
		tax := highestSlab.FixedAmount + (taxableAmount * highestSlab.TaxRate / 100)
		return math.Round(tax*100) / 100
	}

	return 0
}

// calculatePAYEFallback calculates PAYE using hardcoded 2024 Tanzania rates
func (tc *TaxCalculator) calculatePAYEFallback(grossSalary float64) float64 {
	if grossSalary <= 0 {
		return 0
	}

	var tax float64

	switch {
	case grossSalary <= 270000:
		tax = 0
	case grossSalary <= 520000:
		tax = (grossSalary - 270000) * 0.09
	case grossSalary <= 760000:
		tax = 22500 + (grossSalary-520000)*0.20
	case grossSalary <= 1000000:
		tax = 70500 + (grossSalary-760000)*0.25
	default:
		tax = 130500 + (grossSalary-1000000)*0.30
	}

	return math.Round(tax*100) / 100
}

// CalculateNSSF calculates NSSF contributions (employee and employer)
func (tc *TaxCalculator) CalculateNSSF(grossSalary float64, tenantID *uint) (employeeContrib, employerContrib float64) {
	rule, err := tc.complianceRepo.GetStatutoryRuleByCode("NSSF", "Tanzania", tenantID)
	
	var employeeRate, employerRate float64 = 10.0, 10.0 // Default rates
	if err == nil && rule != nil {
		if rule.EmployeeRate != nil {
			employeeRate = *rule.EmployeeRate
		}
		if rule.EmployerRate != nil {
			employerRate = *rule.EmployerRate
		}
	}

	employeeContrib = math.Round(grossSalary*employeeRate) / 100
	employerContrib = math.Round(grossSalary*employerRate) / 100

	return employeeContrib, employerContrib
}

// CalculateNHIF calculates NHIF contributions based on schedule
func (tc *TaxCalculator) CalculateNHIF(grossSalary float64, tenantID *uint) (employeeContrib, employerContrib float64) {
	schedule, err := tc.complianceRepo.GetNHIFSchedule("Tanzania", tenantID)
	if err != nil || len(schedule) == 0 {
		// Fallback to hardcoded rates
		return tc.calculateNHIFFallback(grossSalary)
	}

	for _, s := range schedule {
		if grossSalary >= s.SalaryFrom && grossSalary <= s.SalaryTo {
			return s.EmployeeAmount, s.EmployerAmount
		}
	}

	// If above all bands, use highest band
	if len(schedule) > 0 {
		highest := schedule[len(schedule)-1]
		return highest.EmployeeAmount, highest.EmployerAmount
	}

	return tc.calculateNHIFFallback(grossSalary)
}

// calculateNHIFFallback calculates NHIF using hardcoded schedule
func (tc *TaxCalculator) calculateNHIFFallback(grossSalary float64) (employeeContrib, employerContrib float64) {
	switch {
	case grossSalary <= 150000:
		return 2500, 2500
	case grossSalary <= 250000:
		return 5000, 5000
	case grossSalary <= 400000:
		return 10000, 10000
	case grossSalary <= 600000:
		return 15000, 15000
	default:
		return 20000, 20000
	}
}

// CalculateSDL calculates Skills Development Levy (employer only)
func (tc *TaxCalculator) CalculateSDL(totalPayroll float64, tenantID *uint) float64 {
	rule, err := tc.complianceRepo.GetStatutoryRuleByCode("SDL", "Tanzania", tenantID)
	
	var rate float64 = 5.0 // Default 5%
	if err == nil && rule != nil && rule.EmployerRate != nil {
		rate = *rule.EmployerRate
	}

	return math.Round(totalPayroll*rate) / 100
}

// CalculateWCF calculates Workers Compensation Fund (employer only)
func (tc *TaxCalculator) CalculateWCF(totalPayroll float64, tenantID *uint) float64 {
	rule, err := tc.complianceRepo.GetStatutoryRuleByCode("WCF", "Tanzania", tenantID)
	
	var rate float64 = 1.0 // Default 1%
	if err == nil && rule != nil && rule.EmployerRate != nil {
		rate = *rule.EmployerRate
	}

	return math.Round(totalPayroll*rate) / 100
}

// CalculateAllDeductions calculates all statutory deductions and contributions
func (tc *TaxCalculator) CalculateAllDeductions(grossSalary float64, taxYear int, tenantID *uint) *TaxCalculationResult {
	result := &TaxCalculationResult{}

	// NSSF (calculated before PAYE as it's a deduction from taxable income)
	result.NSSFEmployee, result.NSSFEmployer = tc.CalculateNSSF(grossSalary, tenantID)

	// Taxable income for PAYE (gross minus NSSF employee contribution)
	result.TaxableIncome = grossSalary - result.NSSFEmployee

	// PAYE
	paye, _ := tc.CalculatePAYE(result.TaxableIncome, taxYear, tenantID)
	result.PAYE = paye

	// NHIF
	result.NHIFEmployee, result.NHIFEmployer = tc.CalculateNHIF(grossSalary, tenantID)

	// SDL (employer only)
	result.SDL = tc.CalculateSDL(grossSalary, tenantID)

	// WCF (employer only)
	result.WCF = tc.CalculateWCF(grossSalary, tenantID)

	// Calculate totals
	result.TotalEmployeeDeductions = result.PAYE + result.NSSFEmployee + result.NHIFEmployee
	result.TotalEmployerContributions = result.NSSFEmployer + result.NHIFEmployer + result.SDL + result.WCF
	result.NetPay = grossSalary - result.TotalEmployeeDeductions

	return result
}

// SimulateTax provides a tax simulation for a given CTC
func (tc *TaxCalculator) SimulateTax(ctc float64, taxYear int, tenantID *uint) map[string]interface{} {
	// Assume monthly CTC for simulation
	grossSalary := ctc

	result := tc.CalculateAllDeductions(grossSalary, taxYear, tenantID)

	return map[string]interface{}{
		"grossSalary":              grossSalary,
		"taxableIncome":            result.TaxableIncome,
		"paye":                     result.PAYE,
		"nssfEmployee":             result.NSSFEmployee,
		"nssfEmployer":             result.NSSFEmployer,
		"nhifEmployee":             result.NHIFEmployee,
		"nhifEmployer":             result.NHIFEmployer,
		"sdl":                      result.SDL,
		"wcf":                      result.WCF,
		"totalEmployeeDeductions":  result.TotalEmployeeDeductions,
		"totalEmployerContributions": result.TotalEmployerContributions,
		"netPay":                   result.NetPay,
		"effectiveTaxRate":         math.Round(result.PAYE/grossSalary*10000) / 100,
	}
}

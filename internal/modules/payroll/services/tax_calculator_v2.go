package services

import (
	"fmt"
	"math"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/repositories"
)

// TaxCalculatorV2 handles all statutory calculations using payroll configuration
type TaxCalculatorV2 struct {
	configService *PayrollConfigService
}

// NewTaxCalculatorV2 creates a new tax calculator instance using config service
func NewTaxCalculatorV2() *TaxCalculatorV2 {
	return &TaxCalculatorV2{
		configService: NewPayrollConfigService(),
	}
}

// TaxCalculationResult holds the result of tax calculations
type TaxCalculationResultV2 struct {
	TaxableIncome              float64 `json:"taxableIncome"`
	PAYE                       float64 `json:"paye"`
	NSSFEmployee               float64 `json:"nssfEmployee"`
	NSSFEmployer               float64 `json:"nssfEmployer"`
	NHIFEmployee               float64 `json:"nhifEmployee"`
	NHIFEmployer               float64 `json:"nhifEmployer"`
	SDL                        float64 `json:"sdl"`
	WCF                        float64 `json:"wcf"`
	TotalEmployeeDeductions    float64 `json:"totalEmployeeDeductions"`
	TotalEmployerContributions float64 `json:"totalEmployerContributions"`
	NetPay                     float64 `json:"netPay"`
}

// CalculatePAYE calculates Tanzania PAYE tax using payroll configuration
func (tc *TaxCalculatorV2) CalculatePAYE(taxablePay float64) (float64, error) {
	result, err := tc.configService.CalculatePAYE(taxablePay)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate PAYE: %w", err)
	}
	return result.PAYETax, nil
}

// CalculateNSSF calculates NSSF contributions using payroll configuration
func (tc *TaxCalculatorV2) CalculateNSSF(grossSalary float64) (employeeContrib float64, employerContrib float64, err error) {
	result, err := tc.configService.CalculateNSSF(grossSalary)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to calculate NSSF: %w", err)
	}
	return result.NSSEmployee, result.NSSEmployer, nil
}

// CalculateNHIF calculates NHIF contributions (fallback to database for now)
func (tc *TaxCalculatorV2) CalculateNHIF(grossSalary float64, tenantID *uint) (employeeContrib, employerContrib float64) {
	// For now, use the original method as NHIF is not in config yet
	complianceRepo := repositories.NewComplianceRepository()
	schedule, err := complianceRepo.GetNHIFSchedule("Tanzania", tenantID)
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
func (tc *TaxCalculatorV2) calculateNHIFFallback(grossSalary float64) (employeeContrib, employerContrib float64) {
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

// CalculateSDL calculates Skills Development Levy using payroll configuration
func (tc *TaxCalculatorV2) CalculateSDL(totalPayroll float64) (float64, error) {
	sdlRate, err := tc.configService.GetConfigurationValue("SDL_EMPLOYER_RATE")
	if err != nil {
		return 0, fmt.Errorf("failed to get SDL rate: %w", err)
	}

	rate, ok := sdlRate.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid SDL rate format")
	}

	sdl := totalPayroll * rate
	return float64(int(sdl + 0.5)), nil // Round to nearest integer (TZS)
}

// CalculateWCF calculates Workers Compensation Fund using payroll configuration
func (tc *TaxCalculatorV2) CalculateWCF(totalPayroll float64) (float64, error) {
	wcfRate, err := tc.configService.GetConfigurationValue("WCF_EMPLOYER_RATE")
	if err != nil {
		return 0, fmt.Errorf("failed to get WCF rate: %w", err)
	}

	rate, ok := wcfRate.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid WCF rate format")
	}

	wcf := totalPayroll * rate
	return float64(int(wcf + 0.5)), nil // Round to nearest integer (TZS)
}

// CalculateAllDeductions calculates all statutory deductions and contributions using configuration
func (tc *TaxCalculatorV2) CalculateAllDeductions(grossSalary float64, taxYear int, tenantID *uint) (*TaxCalculationResultV2, error) {
	result := &TaxCalculationResultV2{}

	// NSSF (calculated before PAYE as it's a deduction from taxable income)
	var err error
	result.NSSFEmployee, result.NSSFEmployer, err = tc.CalculateNSSF(grossSalary)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate NSSF: %w", err)
	}

	// Taxable income for PAYE (gross minus NSSF employee contribution)
	result.TaxableIncome = grossSalary - result.NSSFEmployee

	// PAYE
	result.PAYE, err = tc.CalculatePAYE(result.TaxableIncome)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate PAYE: %w", err)
	}

	// NHIF
	result.NHIFEmployee, result.NHIFEmployer = tc.CalculateNHIF(grossSalary, tenantID)

	// SDL (employer only)
	result.SDL, err = tc.CalculateSDL(grossSalary)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate SDL: %w", err)
	}

	// WCF (employer only)
	result.WCF, err = tc.CalculateWCF(grossSalary)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate WCF: %w", err)
	}

	// Calculate totals
	result.TotalEmployeeDeductions = result.PAYE + result.NSSFEmployee + result.NHIFEmployee
	result.TotalEmployerContributions = result.NSSFEmployer + result.NHIFEmployer + result.SDL + result.WCF
	result.NetPay = grossSalary - result.TotalEmployeeDeductions

	return result, nil
}

// SimulateTax provides a tax simulation for a given CTC using configuration
func (tc *TaxCalculatorV2) SimulateTax(ctc float64) (map[string]interface{}, error) {
	result, err := tc.CalculateAllDeductions(ctc, 2024, nil)
	if err != nil {
		return nil, err
	}

	effectiveTaxRate := 0.0
	if ctc > 0 {
		effectiveTaxRate = (result.PAYE / ctc) * 100
	}

	return map[string]interface{}{
		"grossSalary":                ctc,
		"taxableIncome":              result.TaxableIncome,
		"paye":                       result.PAYE,
		"nssfEmployee":               result.NSSFEmployee,
		"nssfEmployer":               result.NSSFEmployer,
		"nhifEmployee":               result.NHIFEmployee,
		"nhifEmployer":               result.NHIFEmployer,
		"sdl":                        result.SDL,
		"wcf":                        result.WCF,
		"totalEmployeeDeductions":    result.TotalEmployeeDeductions,
		"totalEmployerContributions": result.TotalEmployerContributions,
		"netPay":                     result.NetPay,
		"effectiveTaxRate":           math.Round(effectiveTaxRate*100) / 100,
	}, nil
}

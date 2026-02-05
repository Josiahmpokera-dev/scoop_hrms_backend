package services

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/repositories"
)

// ComplianceService handles business logic for compliance
type ComplianceService struct {
	repo          *repositories.ComplianceRepository
	taxCalculator *TaxCalculator
}

// NewComplianceService creates a new service instance
func NewComplianceService() *ComplianceService {
	return &ComplianceService{
		repo:          repositories.NewComplianceRepository(),
		taxCalculator: NewTaxCalculator(),
	}
}

// GetTaxSlabs retrieves tax slabs for a country and year
func (s *ComplianceService) GetTaxSlabs(country string, taxYear int, tenantID *uint) ([]models.TaxSlab, error) {
	return s.repo.GetTaxSlabs(country, taxYear, tenantID)
}

// GetStatutoryRules retrieves statutory rules for a country
func (s *ComplianceService) GetStatutoryRules(country string, tenantID *uint) ([]models.StatutoryRule, error) {
	return s.repo.GetStatutoryRules(country, tenantID)
}

// GetNHIFSchedule retrieves NHIF schedule for a country
func (s *ComplianceService) GetNHIFSchedule(country string, tenantID *uint) ([]models.NHIFSchedule, error) {
	return s.repo.GetNHIFSchedule(country, tenantID)
}

// GetComplianceSummary retrieves compliance summary for a period
func (s *ComplianceService) GetComplianceSummary(tenantID *uint, payMonth, payYear int) (map[string]interface{}, error) {
	return s.repo.GetComplianceSummary(tenantID, payMonth, payYear)
}

// SeedTanzaniaCompliance seeds Tanzania compliance data
func (s *ComplianceService) SeedTanzaniaCompliance() error {
	return s.repo.SeedTanzaniaCompliance()
}

// SimulateTax performs tax simulation for a given salary
func (s *ComplianceService) SimulateTax(grossSalary float64, taxYear int, tenantID *uint) map[string]interface{} {
	return s.taxCalculator.SimulateTax(grossSalary, taxYear, tenantID)
}

// GetCompliancePayments retrieves compliance payments for a period
func (s *ComplianceService) GetCompliancePayments(tenantID *uint, payMonth, payYear int) ([]models.CompliancePayment, error) {
	return s.repo.GetCompliancePayments(tenantID, payMonth, payYear)
}

// CreateCompliancePayment creates a compliance payment record
func (s *ComplianceService) CreateCompliancePayment(payment *models.CompliancePayment) error {
	return s.repo.CreateCompliancePayment(payment)
}

// UpdateCompliancePayment updates a compliance payment
func (s *ComplianceService) UpdateCompliancePayment(payment *models.CompliancePayment) error {
	return s.repo.UpdateCompliancePayment(payment)
}

// CalculateMonthlyStatutory calculates all statutory contributions for a month
func (s *ComplianceService) CalculateMonthlyStatutory(tenantID *uint, payMonth, payYear int, totalGross float64, employeeCount int) (map[string]interface{}, error) {
	taxResult := s.taxCalculator.CalculateAllDeductions(totalGross, payYear, tenantID)

	// Due date is 7th of next month
	dueDate := time.Date(payYear, time.Month(payMonth)+1, 7, 0, 0, 0, 0, time.UTC)
	if payMonth == 12 {
		dueDate = time.Date(payYear+1, 1, 7, 0, 0, 0, 0, time.UTC)
	}

	return map[string]interface{}{
		"payMonth":       payMonth,
		"payYear":        payYear,
		"employeeCount":  employeeCount,
		"totalGross":     totalGross,
		"dueDate":        dueDate.Format("2006-01-02"),
		"daysUntilDue":   int(time.Until(dueDate).Hours() / 24),
		"paye": map[string]interface{}{
			"amount":  taxResult.PAYE,
			"dueDate": dueDate.Format("2006-01-02"),
		},
		"nssf": map[string]interface{}{
			"employeeAmount": taxResult.NSSFEmployee,
			"employerAmount": taxResult.NSSFEmployer,
			"totalAmount":    taxResult.NSSFEmployee + taxResult.NSSFEmployer,
			"dueDate":        dueDate.Format("2006-01-02"),
		},
		"nhif": map[string]interface{}{
			"employeeAmount": taxResult.NHIFEmployee,
			"employerAmount": taxResult.NHIFEmployer,
			"totalAmount":    taxResult.NHIFEmployee + taxResult.NHIFEmployer,
			"dueDate":        dueDate.Format("2006-01-02"),
		},
		"sdl": map[string]interface{}{
			"amount":  taxResult.SDL,
			"dueDate": dueDate.Format("2006-01-02"),
		},
		"wcf": map[string]interface{}{
			"amount":  taxResult.WCF,
			"dueDate": dueDate.Format("2006-01-02"),
		},
		"totalStatutory": taxResult.TotalEmployeeDeductions + taxResult.TotalEmployerContributions,
	}, nil
}

// GenerateComplianceReturn generates compliance return data for filing
func (s *ComplianceService) GenerateComplianceReturn(tenantID *uint, payMonth, payYear int, returnType string) (map[string]interface{}, error) {
	summary, err := s.GetComplianceSummary(tenantID, payMonth, payYear)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"returnType":     returnType,
		"payMonth":       payMonth,
		"payYear":        payYear,
		"payPeriod":      FormatPayslipPeriod(payMonth, payYear),
		"generatedAt":    time.Now().Format(time.RFC3339),
		"summary":        summary,
		"status":         "Generated",
		"filingDeadline": time.Date(payYear, time.Month(payMonth)+1, 7, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
	}, nil
}

// GetAllComplianceRules retrieves all compliance rules for display
func (s *ComplianceService) GetAllComplianceRules(country string, tenantID *uint) (map[string]interface{}, error) {
	taxSlabs, err := s.repo.GetTaxSlabs(country, time.Now().Year(), tenantID)
	if err != nil {
		taxSlabs = []models.TaxSlab{}
	}

	rules, err := s.repo.GetStatutoryRules(country, tenantID)
	if err != nil {
		rules = []models.StatutoryRule{}
	}

	nhifSchedule, err := s.repo.GetNHIFSchedule(country, tenantID)
	if err != nil {
		nhifSchedule = []models.NHIFSchedule{}
	}

	// Format tax slabs
	var formattedSlabs []map[string]interface{}
	for _, slab := range taxSlabs {
		formattedSlabs = append(formattedSlabs, map[string]interface{}{
			"from":        slab.SlabFrom,
			"to":          slab.SlabTo,
			"rate":        slab.TaxRate,
			"fixedAmount": slab.FixedAmount,
			"description": slab.Description,
		})
	}

	// Format rules
	formattedRules := make(map[string]interface{})
	for _, rule := range rules {
		ruleData := map[string]interface{}{
			"name":        rule.Name,
			"basis":       rule.Basis,
			"description": rule.Description,
			"authority":   rule.Authority,
			"dueDay":      rule.DueDay,
		}
		if rule.EmployeeRate != nil {
			ruleData["employeeRate"] = *rule.EmployeeRate
		}
		if rule.EmployerRate != nil {
			ruleData["employerRate"] = *rule.EmployerRate
		}
		if rule.Cap != nil {
			ruleData["cap"] = *rule.Cap
		}
		formattedRules[rule.Code] = ruleData
	}

	// Format NHIF schedule
	var formattedNHIF []map[string]interface{}
	for _, sch := range nhifSchedule {
		formattedNHIF = append(formattedNHIF, map[string]interface{}{
			"from":     sch.SalaryFrom,
			"to":       sch.SalaryTo,
			"employee": sch.EmployeeAmount,
			"employer": sch.EmployerAmount,
		})
	}

	return map[string]interface{}{
		"country":        country,
		"taxYear":        time.Now().Year(),
		"payeSlabs":      formattedSlabs,
		"statutoryRules": formattedRules,
		"nhifSchedule":   formattedNHIF,
	}, nil
}

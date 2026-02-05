package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"gorm.io/gorm"
)

// ComplianceRepository handles database operations for compliance
type ComplianceRepository struct {
	db *gorm.DB
}

// NewComplianceRepository creates a new repository instance
func NewComplianceRepository() *ComplianceRepository {
	return &ComplianceRepository{
		db: database.DB,
	}
}

// GetTaxSlabs retrieves tax slabs for a country and year
func (r *ComplianceRepository) GetTaxSlabs(country string, taxYear int, tenantID *uint) ([]models.TaxSlab, error) {
	var slabs []models.TaxSlab
	query := r.db.Where("country = ? AND tax_year = ? AND is_active = ?", country, taxYear, true)
	if tenantID != nil {
		query = query.Where("tenant_id = ? OR tenant_id IS NULL", *tenantID)
	}
	if err := query.Order("sort_order ASC, slab_from ASC").Find(&slabs).Error; err != nil {
		return nil, err
	}
	return slabs, nil
}

// CreateTaxSlab creates a tax slab
func (r *ComplianceRepository) CreateTaxSlab(slab *models.TaxSlab) error {
	return r.db.Create(slab).Error
}

// UpdateTaxSlab updates a tax slab
func (r *ComplianceRepository) UpdateTaxSlab(slab *models.TaxSlab) error {
	return r.db.Save(slab).Error
}

// GetStatutoryRules retrieves statutory rules for a country
func (r *ComplianceRepository) GetStatutoryRules(country string, tenantID *uint) ([]models.StatutoryRule, error) {
	var rules []models.StatutoryRule
	query := r.db.Where("country = ? AND is_active = ?", country, true)
	if tenantID != nil {
		query = query.Where("tenant_id = ? OR tenant_id IS NULL", *tenantID)
	}
	if err := query.Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// GetStatutoryRuleByCode retrieves a statutory rule by code
func (r *ComplianceRepository) GetStatutoryRuleByCode(code, country string, tenantID *uint) (*models.StatutoryRule, error) {
	var rule models.StatutoryRule
	query := r.db.Where("code = ? AND country = ? AND is_active = ?", code, country, true)
	if tenantID != nil {
		query = query.Where("tenant_id = ? OR tenant_id IS NULL", *tenantID)
	}
	if err := query.First(&rule).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// CreateStatutoryRule creates a statutory rule
func (r *ComplianceRepository) CreateStatutoryRule(rule *models.StatutoryRule) error {
	return r.db.Create(rule).Error
}

// UpdateStatutoryRule updates a statutory rule
func (r *ComplianceRepository) UpdateStatutoryRule(rule *models.StatutoryRule) error {
	return r.db.Save(rule).Error
}

// GetNHIFSchedule retrieves NHIF schedule for a country
func (r *ComplianceRepository) GetNHIFSchedule(country string, tenantID *uint) ([]models.NHIFSchedule, error) {
	var schedule []models.NHIFSchedule
	query := r.db.Where("country = ? AND is_active = ?", country, true)
	if tenantID != nil {
		query = query.Where("tenant_id = ? OR tenant_id IS NULL", *tenantID)
	}
	if err := query.Order("sort_order ASC, salary_from ASC").Find(&schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}

// CreateNHIFSchedule creates an NHIF schedule entry
func (r *ComplianceRepository) CreateNHIFSchedule(schedule *models.NHIFSchedule) error {
	return r.db.Create(schedule).Error
}

// GetCompliancePayments retrieves compliance payments for a period
func (r *ComplianceRepository) GetCompliancePayments(tenantID *uint, payMonth, payYear int) ([]models.CompliancePayment, error) {
	var payments []models.CompliancePayment
	query := r.db.Where("pay_month = ? AND pay_year = ?", payMonth, payYear)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// GetCompliancePaymentByCode retrieves a compliance payment by statutory code
func (r *ComplianceRepository) GetCompliancePaymentByCode(tenantID *uint, payMonth, payYear int, code string) (*models.CompliancePayment, error) {
	var payment models.CompliancePayment
	query := r.db.Where("pay_month = ? AND pay_year = ? AND statutory_code = ?", payMonth, payYear, code)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// CreateCompliancePayment creates a compliance payment record
func (r *ComplianceRepository) CreateCompliancePayment(payment *models.CompliancePayment) error {
	return r.db.Create(payment).Error
}

// UpdateCompliancePayment updates a compliance payment
func (r *ComplianceRepository) UpdateCompliancePayment(payment *models.CompliancePayment) error {
	return r.db.Save(payment).Error
}

// GetComplianceSummary retrieves compliance summary for a period
func (r *ComplianceRepository) GetComplianceSummary(tenantID *uint, payMonth, payYear int) (map[string]interface{}, error) {
	payments, err := r.GetCompliancePayments(tenantID, payMonth, payYear)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Calculate due date (7th of next month)
	nextMonth := time.Date(payYear, time.Month(payMonth)+1, 7, 0, 0, 0, 0, time.UTC)
	if payMonth == 12 {
		nextMonth = time.Date(payYear+1, 1, 7, 0, 0, 0, 0, time.UTC)
	}
	daysUntilDue := int(time.Until(nextMonth).Hours() / 24)
	if daysUntilDue < 0 {
		daysUntilDue = 0
	}

	summary := map[string]interface{}{
		"payMonth":      payMonth,
		"payYear":       payYear,
		"nextDueDate":   nextMonth.Format("2006-01-02"),
		"daysUntilDue":  daysUntilDue,
	}

	var totalStatutory float64
	for _, p := range payments {
		switch p.StatutoryCode {
		case "PAYE":
			summary["paye"] = map[string]interface{}{
				"amount":         p.TotalAmount,
				"employeesCount": p.EmployeesCount,
				"dueDate":        p.DueDate.Format("2006-01-02"),
				"status":         p.Status,
			}
		case "NSSF":
			summary["nssf"] = map[string]interface{}{
				"employeeAmount":  p.EmployeeAmount,
				"employerAmount":  p.EmployerAmount,
				"totalAmount":     p.TotalAmount,
				"employeesCount":  p.EmployeesCount,
				"dueDate":         p.DueDate.Format("2006-01-02"),
				"status":          p.Status,
			}
		case "NHIF":
			summary["nhif"] = map[string]interface{}{
				"employeeAmount":  p.EmployeeAmount,
				"employerAmount":  p.EmployerAmount,
				"totalAmount":     p.TotalAmount,
				"employeesCount":  p.EmployeesCount,
				"dueDate":         p.DueDate.Format("2006-01-02"),
				"status":          p.Status,
			}
		case "SDL":
			summary["sdl"] = map[string]interface{}{
				"amount":  p.TotalAmount,
				"dueDate": p.DueDate.Format("2006-01-02"),
				"status":  p.Status,
			}
		case "WCF":
			summary["wcf"] = map[string]interface{}{
				"amount":  p.TotalAmount,
				"dueDate": p.DueDate.Format("2006-01-02"),
				"status":  p.Status,
			}
		}
		totalStatutory += p.TotalAmount
	}
	summary["totalStatutory"] = totalStatutory

	return summary, nil
}

// SeedTanzaniaCompliance seeds Tanzania compliance data
func (r *ComplianceRepository) SeedTanzaniaCompliance() error {
	// Check if already seeded
	var count int64
	r.db.Model(&models.TaxSlab{}).Where("country = ? AND tax_year = ?", "Tanzania", 2024).Count(&count)
	if count > 0 {
		return nil // Already seeded
	}

	// Tanzania PAYE Tax Slabs 2024 (Monthly)
	taxSlabs := []models.TaxSlab{
		{Country: "Tanzania", TaxYear: 2024, Period: "Monthly", SlabFrom: 0, SlabTo: 270000, TaxRate: 0, FixedAmount: 0, Description: "Tax-free threshold", SortOrder: 1, IsActive: true},
		{Country: "Tanzania", TaxYear: 2024, Period: "Monthly", SlabFrom: 270001, SlabTo: 520000, TaxRate: 9, FixedAmount: 0, Description: "9% on amount above 270,000", SortOrder: 2, IsActive: true},
		{Country: "Tanzania", TaxYear: 2024, Period: "Monthly", SlabFrom: 520001, SlabTo: 760000, TaxRate: 20, FixedAmount: 22500, Description: "TZS 22,500 + 20% on amount above 520,000", SortOrder: 3, IsActive: true},
		{Country: "Tanzania", TaxYear: 2024, Period: "Monthly", SlabFrom: 760001, SlabTo: 1000000, TaxRate: 25, FixedAmount: 70500, Description: "TZS 70,500 + 25% on amount above 760,000", SortOrder: 4, IsActive: true},
		{Country: "Tanzania", TaxYear: 2024, Period: "Monthly", SlabFrom: 1000001, SlabTo: 999999999, TaxRate: 30, FixedAmount: 130500, Description: "TZS 130,500 + 30% on amount above 1,000,000", SortOrder: 5, IsActive: true},
	}
	for _, slab := range taxSlabs {
		if err := r.db.Create(&slab).Error; err != nil {
			return err
		}
	}

	// Statutory Rules
	nssfEmployeeRate := 10.0
	nssfEmployerRate := 10.0
	sdlEmployerRate := 5.0
	wcfEmployerRate := 1.0

	statutoryRules := []models.StatutoryRule{
		{Country: "Tanzania", Code: "NSSF", Name: "National Social Security Fund", EmployeeRate: &nssfEmployeeRate, EmployerRate: &nssfEmployerRate, Basis: "Gross salary", Description: "Employee: 10% of gross salary. Employer: 10% of gross salary.", Authority: "NSSF Tanzania", DueDay: 7, IsActive: true},
		{Country: "Tanzania", Code: "NHIF", Name: "National Health Insurance Fund", Basis: "Schedule", Description: "Contribution as per NHIF schedule based on salary bands.", Authority: "NHIF Tanzania", DueDay: 7, IsActive: true},
		{Country: "Tanzania", Code: "SDL", Name: "Skills Development Levy", EmployerRate: &sdlEmployerRate, Basis: "Total payroll", Description: "Employer only: 5% of total gross payroll.", Authority: "TRA", DueDay: 7, IsActive: true},
		{Country: "Tanzania", Code: "WCF", Name: "Workers Compensation Fund", EmployerRate: &wcfEmployerRate, Basis: "Total payroll", Description: "Employer only: 1% of total gross payroll.", Authority: "WCF Tanzania", DueDay: 7, IsActive: true},
	}
	for _, rule := range statutoryRules {
		if err := r.db.Create(&rule).Error; err != nil {
			return err
		}
	}

	// NHIF Schedule
	nhifSchedule := []models.NHIFSchedule{
		{Country: "Tanzania", SalaryFrom: 0, SalaryTo: 150000, EmployeeAmount: 2500, EmployerAmount: 2500, SortOrder: 1, IsActive: true},
		{Country: "Tanzania", SalaryFrom: 150001, SalaryTo: 250000, EmployeeAmount: 5000, EmployerAmount: 5000, SortOrder: 2, IsActive: true},
		{Country: "Tanzania", SalaryFrom: 250001, SalaryTo: 400000, EmployeeAmount: 10000, EmployerAmount: 10000, SortOrder: 3, IsActive: true},
		{Country: "Tanzania", SalaryFrom: 400001, SalaryTo: 600000, EmployeeAmount: 15000, EmployerAmount: 15000, SortOrder: 4, IsActive: true},
		{Country: "Tanzania", SalaryFrom: 600001, SalaryTo: 999999999, EmployeeAmount: 20000, EmployerAmount: 20000, SortOrder: 5, IsActive: true},
	}
	for _, schedule := range nhifSchedule {
		if err := r.db.Create(&schedule).Error; err != nil {
			return err
		}
	}

	return nil
}

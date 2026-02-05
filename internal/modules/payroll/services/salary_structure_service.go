package services

import (
	"errors"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/repositories"
)

// SalaryStructureService handles business logic for salary structures
type SalaryStructureService struct {
	repo          *repositories.SalaryStructureRepository
	componentRepo *repositories.SalaryComponentRepository
	taxCalculator *TaxCalculator
}

// NewSalaryStructureService creates a new service instance
func NewSalaryStructureService() *SalaryStructureService {
	return &SalaryStructureService{
		repo:          repositories.NewSalaryStructureRepository(),
		componentRepo: repositories.NewSalaryComponentRepository(),
		taxCalculator: NewTaxCalculator(),
	}
}

// CreateSalaryStructure creates a new salary structure template
func (s *SalaryStructureService) CreateSalaryStructure(structure *models.SalaryStructure) error {
	// Calculate gross salary from components
	structure.GrossSalary = structure.Basic + structure.HRA + structure.Transport + structure.Medical + structure.OtherAllowances

	// Calculate statutory deductions using Tanzania rates
	taxYear := 2024 // Default to current year
	taxResult := s.taxCalculator.CalculateAllDeductions(structure.GrossSalary, taxYear, structure.TenantID)

	structure.PAYEDeduction = taxResult.PAYE
	structure.NSSFEmployee = taxResult.NSSFEmployee
	structure.NHIFEmployee = taxResult.NHIFEmployee
	structure.TotalDeductions = taxResult.TotalEmployeeDeductions

	structure.NSSFEmployer = taxResult.NSSFEmployer
	structure.NHIFEmployer = taxResult.NHIFEmployer
	structure.SDL = taxResult.SDL
	structure.WCF = taxResult.WCF
	structure.TotalEmployerContribs = taxResult.TotalEmployerContributions

	structure.NetPay = taxResult.NetPay

	return s.repo.Create(structure)
}

// GetSalaryStructure retrieves a salary structure by ID
func (s *SalaryStructureService) GetSalaryStructure(id uint, tenantID *uint) (*models.SalaryStructure, error) {
	return s.repo.GetByID(id, tenantID)
}

// UpdateSalaryStructure updates a salary structure
func (s *SalaryStructureService) UpdateSalaryStructure(structure *models.SalaryStructure) error {
	// Recalculate gross salary and deductions
	structure.GrossSalary = structure.Basic + structure.HRA + structure.Transport + structure.Medical + structure.OtherAllowances

	taxYear := 2024
	taxResult := s.taxCalculator.CalculateAllDeductions(structure.GrossSalary, taxYear, structure.TenantID)

	structure.PAYEDeduction = taxResult.PAYE
	structure.NSSFEmployee = taxResult.NSSFEmployee
	structure.NHIFEmployee = taxResult.NHIFEmployee
	structure.TotalDeductions = taxResult.TotalEmployeeDeductions

	structure.NSSFEmployer = taxResult.NSSFEmployer
	structure.NHIFEmployer = taxResult.NHIFEmployer
	structure.SDL = taxResult.SDL
	structure.WCF = taxResult.WCF
	structure.TotalEmployerContribs = taxResult.TotalEmployerContributions

	structure.NetPay = taxResult.NetPay

	return s.repo.Update(structure)
}

// DeleteSalaryStructure deletes a salary structure
func (s *SalaryStructureService) DeleteSalaryStructure(id uint, tenantID *uint) error {
	return s.repo.Delete(id, tenantID)
}

// ListSalaryStructures lists salary structures with filters
func (s *SalaryStructureService) ListSalaryStructures(tenantID *uint, isActive *bool, grade, location, country string, page, pageSize int) ([]models.SalaryStructure, int64, error) {
	return s.repo.List(tenantID, isActive, grade, location, country, page, pageSize)
}

// SimulateSalary simulates salary calculation based on CTC
func (s *SalaryStructureService) SimulateSalary(ctc float64, tenantID *uint) map[string]interface{} {
	// Default breakdown percentages
	basicPercentage := 50.0
	hraPercentage := 20.0
	transportPercentage := 10.0
	medicalPercentage := 10.0
	otherPercentage := 10.0

	basic := ctc * basicPercentage / 100
	hra := ctc * hraPercentage / 100
	transport := ctc * transportPercentage / 100
	medical := ctc * medicalPercentage / 100
	other := ctc * otherPercentage / 100
	grossSalary := basic + hra + transport + medical + other

	taxYear := 2024
	taxResult := s.taxCalculator.CalculateAllDeductions(grossSalary, taxYear, tenantID)

	return map[string]interface{}{
		"ctc": ctc,
		"earnings": map[string]interface{}{
			"basic":           basic,
			"hra":             hra,
			"transport":       transport,
			"medical":         medical,
			"otherAllowances": other,
			"grossSalary":     grossSalary,
		},
		"deductions": map[string]interface{}{
			"paye":         taxResult.PAYE,
			"nssfEmployee": taxResult.NSSFEmployee,
			"nhifEmployee": taxResult.NHIFEmployee,
			"total":        taxResult.TotalEmployeeDeductions,
		},
		"employerContributions": map[string]interface{}{
			"nssfEmployer": taxResult.NSSFEmployer,
			"nhifEmployer": taxResult.NHIFEmployer,
			"sdl":          taxResult.SDL,
			"wcf":          taxResult.WCF,
			"total":        taxResult.TotalEmployerContributions,
		},
		"netPay": taxResult.NetPay,
	}
}

// CreateSalaryComponent creates a new salary component
func (s *SalaryStructureService) CreateSalaryComponent(component *models.SalaryComponent) error {
	// Check if component code already exists
	existing, err := s.componentRepo.GetByCode(component.ComponentCode, component.TenantID)
	if err == nil && existing != nil {
		return errors.New("component code already exists")
	}

	return s.componentRepo.Create(component)
}

// GetSalaryComponent retrieves a salary component by ID
func (s *SalaryStructureService) GetSalaryComponent(id uint, tenantID *uint) (*models.SalaryComponent, error) {
	return s.componentRepo.GetByID(id, tenantID)
}

// UpdateSalaryComponent updates a salary component
func (s *SalaryStructureService) UpdateSalaryComponent(component *models.SalaryComponent) error {
	return s.componentRepo.Update(component)
}

// DeleteSalaryComponent deletes a salary component
func (s *SalaryStructureService) DeleteSalaryComponent(id uint, tenantID *uint) error {
	return s.componentRepo.Delete(id, tenantID)
}

// ListSalaryComponents lists salary components with filters
func (s *SalaryStructureService) ListSalaryComponents(tenantID *uint, componentType string, isActive, isStatutory *bool, country string, page, pageSize int) ([]models.SalaryComponent, int64, error) {
	return s.componentRepo.List(tenantID, componentType, isActive, isStatutory, country, page, pageSize)
}

// GetAllActiveComponents retrieves all active salary components
func (s *SalaryStructureService) GetAllActiveComponents(tenantID *uint, country string) ([]models.SalaryComponent, error) {
	return s.componentRepo.GetAllActive(tenantID, country)
}

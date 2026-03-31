package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/repositories"
)

// PayrollConfigService handles payroll configuration business logic
type PayrollConfigService struct {
	repo  *repositories.PayrollConfigRepository
	cache *ConfigCache
}

// NewPayrollConfigService creates a new payroll config service
func NewPayrollConfigService() *PayrollConfigService {
	service := &PayrollConfigService{
		repo: repositories.NewPayrollConfigRepository(),
	}
	// Initialize cache with the service instance to avoid circular dependency
	service.cache = NewConfigCache(service)
	return service
}

// GetAllConfigurations returns all payroll configurations
func (s *PayrollConfigService) GetAllConfigurations() (map[string]interface{}, error) {
	return s.cache.GetAllConfigurations()
}

// GetConfiguration returns a single configuration by key
func (s *PayrollConfigService) GetConfiguration(configKey string) (*models.PayrollConfig, error) {
	return s.cache.GetConfiguration(configKey)
}

// GetConfigurationValue returns a single configuration value by key
func (s *PayrollConfigService) GetConfigurationValue(configKey string) (interface{}, error) {
	return s.cache.GetConfigurationValue(configKey)
}

// GetPAYEBands returns all PAYE tax bands
func (s *PayrollConfigService) GetPAYEBands() ([]models.PAYEBand, error) {
	payeConfigs, err := s.repo.GetPAYEBands()
	if err != nil {
		return nil, fmt.Errorf("failed to get PAYE bands: %w", err)
	}

	bandMap := make(map[int]*models.PAYEBand)

	for _, config := range payeConfigs {
		bandNum := extractBandNumber(config.ConfigKey)
		if bandNum <= 0 {
			continue
		}

		// Get or create band
		band, exists := bandMap[bandNum]
		if !exists {
			band = &models.PAYEBand{Band: bandNum}
			bandMap[bandNum] = band
		}

		// Set band label if not set
		if band.Label == "" {
			band.Label = s.getPAYEBandLabel(bandNum)
		}

		// Set values based on config key suffix
		switch {
		case strings.HasSuffix(config.ConfigKey, "_MIN"):
			value, _ := strconv.ParseFloat(config.Value, 64)
			band.Min = value
		case strings.HasSuffix(config.ConfigKey, "_MAX"):
			if config.Value != "null" && config.Value != "" {
				value, _ := strconv.ParseFloat(config.Value, 64)
				band.Max = &value
			}
		case strings.HasSuffix(config.ConfigKey, "_RATE"):
			value, _ := strconv.ParseFloat(config.Value, 64)
			band.Rate = value
		case strings.HasSuffix(config.ConfigKey, "_BASE_TAX"):
			value, _ := strconv.ParseFloat(config.Value, 64)
			band.BaseTax = value
		}
	}

	// Convert map to slice
	bands := make([]models.PAYEBand, 0, len(bandMap))
	for i := 1; i <= len(bandMap); i++ {
		if band, exists := bandMap[i]; exists {
			bands = append(bands, *band)
		}
	}

	return bands, nil
}

// CalculatePAYE calculates PAYE tax for given taxable pay
func (s *PayrollConfigService) CalculatePAYE(taxablePay float64) (*models.PAYECalculationResult, error) {
	return s.cache.CalculatePAYE(taxablePay)
}

// CalculateNSSF calculates NSSF contributions
func (s *PayrollConfigService) CalculateNSSF(grossSalary float64) (*models.NSSFCalculationResult, error) {
	employeeRate, err := s.getFloatValue("NSSF_EMPLOYEE_RATE")
	if err != nil {
		return nil, err
	}

	employerRate, err := s.getFloatValue("NSSF_EMPLOYER_RATE")
	if err != nil {
		return nil, err
	}

	nssfEmployee := grossSalary * employeeRate
	nssfEmployer := grossSalary * employerRate
	nssfTotal := nssfEmployee + nssfEmployer

	// Round to nearest integer (TZS)
	nssfEmployee = float64(int(nssfEmployee + 0.5))
	nssfEmployer = float64(int(nssfEmployer + 0.5))
	nssfTotal = float64(int(nssfTotal + 0.5))

	return &models.NSSFCalculationResult{
		GrossSalary:  grossSalary,
		NSSEmployee:  nssfEmployee,
		NSSEmployer:  nssfEmployer,
		NSSTotal:     nssfTotal,
		RateEmployee: employeeRate,
		RateEmployer: employerRate,
	}, nil
}

// CalculateCTC calculates Cost to Company
func (s *PayrollConfigService) CalculateCTC(grossSalary float64) (*models.CTCCalculationResult, error) {
	nssfEmployerRate, err := s.getFloatValue("NSSF_EMPLOYER_RATE")
	if err != nil {
		return nil, err
	}

	wcfRate, err := s.getFloatValue("WCF_EMPLOYER_RATE")
	if err != nil {
		return nil, err
	}

	sdlRate, err := s.getFloatValue("SDL_EMPLOYER_RATE")
	if err != nil {
		return nil, err
	}

	nssfEmployer := grossSalary * nssfEmployerRate
	wcf := grossSalary * wcfRate
	sdl := grossSalary * sdlRate
	totalCTC := grossSalary + nssfEmployer + wcf + sdl

	// Round to nearest integer (TZS)
	nssfEmployer = float64(int(nssfEmployer + 0.5))
	wcf = float64(int(wcf + 0.5))
	sdl = float64(int(sdl + 0.5))
	totalCTC = float64(int(totalCTC + 0.5))

	return &models.CTCCalculationResult{
		GrossSalary: grossSalary,
		NSSEmployer: nssfEmployer,
		WCF:         wcf,
		SDL:         sdl,
		TotalCTC:    totalCTC,
	}, nil
}

// GetNSSFRates returns NSSF rates
func (s *PayrollConfigService) GetNSSFRates() (map[string]interface{}, error) {
	configs, err := s.repo.GetNSSFRates()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, config := range configs {
		value, err := strconv.ParseFloat(config.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", config.ConfigKey, err)
		}
		result[config.ConfigKey] = value
	}

	return result, nil
}

// GetEmployerContributions returns employer contribution rates
func (s *PayrollConfigService) GetEmployerContributions() (map[string]interface{}, error) {
	configs, err := s.repo.GetEmployerContributions()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, config := range configs {
		value, err := strconv.ParseFloat(config.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", config.ConfigKey, err)
		}
		result[config.ConfigKey] = value
	}

	return result, nil
}

// GetFormulaConfigurations returns formula configurations
func (s *PayrollConfigService) GetFormulaConfigurations() (map[string]string, error) {
	configs, err := s.repo.GetFormulaConfigurations()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, config := range configs {
		result[config.ConfigKey] = config.Value
	}

	return result, nil
}

// Helper methods

func (s *PayrollConfigService) getFloatValue(configKey string) (float64, error) {
	value, err := s.repo.GetValue(configKey)
	if err != nil {
		return 0, err
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s value: %w", configKey, err)
	}

	return floatValue, nil
}

func (s *PayrollConfigService) parseConfigValue(configKey string, value string) (interface{}, error) {
	// Determine value type based on key
	if strings.Contains(configKey, "_RATE") || strings.Contains(configKey, "_TAX") {
		return strconv.ParseFloat(value, 64)
	}
	if strings.Contains(configKey, "_MAX") || strings.Contains(configKey, "_MIN") {
		return strconv.ParseFloat(value, 64)
	}
	if strings.Contains(configKey, "_FORMULA") {
		return value, nil // Keep as string for formulas
	}
	if configKey == "PAYROLL_CURRENCY" || configKey == "PAYROLL_FREQUENCY" {
		return value, nil // Keep as string
	}
	return value, nil // Default to string
}

func (s *PayrollConfigService) buildPAYEBand(bands []models.PAYEBand, band int, config models.PayrollConfig) []models.PAYEBand {
	// Find or create band
	var currentBand *models.PAYEBand
	for i := range bands {
		if bands[i].Band == band {
			currentBand = &bands[i]
			break
		}
	}

	if currentBand == nil {
		bands = append(bands, models.PAYEBand{Band: band})
		currentBand = &bands[len(bands)-1]
	}

	// Set values based on config key
	switch {
	case strings.HasSuffix(config.ConfigKey, "_MIN"):
		value, _ := strconv.ParseFloat(config.Value, 64)
		currentBand.Min = value
	case strings.HasSuffix(config.ConfigKey, "_MAX"):
		if config.Value != "null" {
			value, _ := strconv.ParseFloat(config.Value, 64)
			currentBand.Max = &value
		}
	case strings.HasSuffix(config.ConfigKey, "_RATE"):
		value, _ := strconv.ParseFloat(config.Value, 64)
		currentBand.Rate = value
	case strings.HasSuffix(config.ConfigKey, "_BASE_TAX"):
		value, _ := strconv.ParseFloat(config.Value, 64)
		currentBand.BaseTax = value
	}

	// Set label
	if currentBand.Label == "" {
		currentBand.Label = s.getPAYEBandLabel(band)
	}

	return bands
}

func (s *PayrollConfigService) getPAYEBandLabel(band int) string {
	switch band {
	case 1:
		return "Up to TZS 270,000 — 0%"
	case 2:
		return "TZS 270,001 - 520,000 — 8%"
	case 3:
		return "TZS 520,001 - 760,000 — 20%"
	case 4:
		return "TZS 760,001 - 1,000,000 — 25%"
	case 5:
		return "Above TZS 1,000,000 — 30%"
	default:
		return fmt.Sprintf("Band %d", band)
	}
}

func extractBandNumber(configKey string) int {
	// Extract band number from PAYE_BAND_X_MIN
	parts := strings.Split(configKey, "_")
	if len(parts) >= 3 {
		var band int
		fmt.Sscanf(parts[2], "%d", &band)
		return band
	}
	return 0
}

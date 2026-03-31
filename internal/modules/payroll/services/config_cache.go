package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
)

// ConfigCache provides thread-safe caching for payroll configuration values
type ConfigCache struct {
	cache    map[string]interface{}
	expiry   map[string]time.Time
	mu       sync.RWMutex
	duration time.Duration
	service  *PayrollConfigService
}

// NewConfigCache creates a new configuration cache
func NewConfigCache(service *PayrollConfigService) *ConfigCache {
	return &ConfigCache{
		cache:    make(map[string]interface{}),
		expiry:   make(map[string]time.Time),
		duration: 1 * time.Hour, // 1 hour cache duration as specified
		service:  service,
	}
}

// GetAllConfigurations returns all configurations from cache or fetches from service
func (c *ConfigCache) GetAllConfigurations() (map[string]interface{}, error) {
	c.mu.RLock()
	if expiry, exists := c.expiry["all_configs"]; exists && time.Now().Before(expiry) {
		if configs, ok := c.cache["all_configs"].(map[string]interface{}); ok {
			c.mu.RUnlock()
			return configs, nil
		}
	}
	c.mu.RUnlock()

	// Fetch from service
	configs, err := c.service.GetAllConfigurations()
	if err != nil {
		return nil, err
	}

	// Update cache
	c.mu.Lock()
	c.cache["all_configs"] = configs
	c.expiry["all_configs"] = time.Now().Add(c.duration)
	c.mu.Unlock()

	return configs, nil
}

// GetConfiguration returns a single configuration from cache or fetches from service
func (c *ConfigCache) GetConfiguration(configKey string) (*models.PayrollConfig, error) {
	cacheKey := fmt.Sprintf("config_%s", configKey)

	c.mu.RLock()
	if expiry, exists := c.expiry[cacheKey]; exists && time.Now().Before(expiry) {
		if config, ok := c.cache[cacheKey].(*models.PayrollConfig); ok {
			c.mu.RUnlock()
			return config, nil
		}
	}
	c.mu.RUnlock()

	// Fetch from service
	config, err := c.service.GetConfiguration(configKey)
	if err != nil {
		return nil, err
	}

	// Update cache
	c.mu.Lock()
	c.cache[cacheKey] = config
	c.expiry[cacheKey] = time.Now().Add(c.duration)
	c.mu.Unlock()

	return config, nil
}

// GetConfigurationValue returns a configuration value from cache or fetches from service
func (c *ConfigCache) GetConfigurationValue(configKey string) (interface{}, error) {
	cacheKey := fmt.Sprintf("value_%s", configKey)

	c.mu.RLock()
	if expiry, exists := c.expiry[cacheKey]; exists && time.Now().Before(expiry) {
		if value, ok := c.cache[cacheKey]; ok {
			c.mu.RUnlock()
			return value, nil
		}
	}
	c.mu.RUnlock()

	// Fetch from service
	value, err := c.service.GetConfigurationValue(configKey)
	if err != nil {
		return nil, err
	}

	// Update cache
	c.mu.Lock()
	c.cache[cacheKey] = value
	c.expiry[cacheKey] = time.Now().Add(c.duration)
	c.mu.Unlock()

	return value, nil
}

// CalculatePAYE calculates PAYE using cache
func (c *ConfigCache) CalculatePAYE(taxablePay float64) (*models.PAYECalculationResult, error) {
	// PAYE calculation is always fresh as it depends on taxable pay amount
	return c.service.CalculatePAYE(taxablePay)
}

// CalculateNSSF calculates NSSF using cache for rates
func (c *ConfigCache) CalculateNSSF(grossSalary float64) (*models.NSSFCalculationResult, error) {
	// Get rates from cache
	employeeRate, err := c.getFloatValueCached("NSSF_EMPLOYEE_RATE")
	if err != nil {
		return nil, err
	}

	employerRate, err := c.getFloatValueCached("NSSF_EMPLOYER_RATE")
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

// CalculateCTC calculates CTC using cache for rates
func (c *ConfigCache) CalculateCTC(grossSalary float64) (*models.CTCCalculationResult, error) {
	// Get rates from cache
	nssfEmployerRate, err := c.getFloatValueCached("NSSF_EMPLOYER_RATE")
	if err != nil {
		return nil, err
	}

	wcfRate, err := c.getFloatValueCached("WCF_EMPLOYER_RATE")
	if err != nil {
		return nil, err
	}

	sdlRate, err := c.getFloatValueCached("SDL_EMPLOYER_RATE")
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

// getFloatValueCached gets a float value from cache
func (c *ConfigCache) getFloatValueCached(configKey string) (float64, error) {
	value, err := c.GetConfigurationValue(configKey)
	if err != nil {
		return 0, err
	}

	floatValue, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid value type for %s", configKey)
	}

	return floatValue, nil
}

// ClearCache clears all cached values
func (c *ConfigCache) ClearCache() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]interface{})
	c.expiry = make(map[string]time.Time)
}

// GetCacheStats returns cache statistics
func (c *ConfigCache) GetCacheStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	expiredCount := 0
	activeCount := 0
	now := time.Now()

	for _, expiry := range c.expiry {
		if now.After(expiry) {
			expiredCount++
		} else {
			activeCount++
		}
	}

	return map[string]interface{}{
		"total_keys":      len(c.cache),
		"active_entries":  activeCount,
		"expired_entries": expiredCount,
		"cache_duration":  c.duration.String(),
	}
}

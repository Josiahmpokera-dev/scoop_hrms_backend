package repositories

import (
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"gorm.io/gorm"
)

// PayrollConfigRepository handles payroll configuration database operations
type PayrollConfigRepository struct {
	db *gorm.DB
}

// NewPayrollConfigRepository creates a new payroll config repository
func NewPayrollConfigRepository() *PayrollConfigRepository {
	return &PayrollConfigRepository{
		db: database.GetDB(),
	}
}

// GetByKey retrieves a configuration by key
func (r *PayrollConfigRepository) GetByKey(configKey string) (*models.PayrollConfig, error) {
	var config models.PayrollConfig
	err := r.db.Where("config_key = ?", configKey).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("configuration not found: %s", configKey)
		}
		return nil, err
	}
	return &config, nil
}

// GetByKeys retrieves multiple configurations by keys
func (r *PayrollConfigRepository) GetByKeys(configKeys []string) ([]models.PayrollConfig, error) {
	var configs []models.PayrollConfig
	err := r.db.Where("config_key IN ?", configKeys).Find(&configs).Error
	return configs, err
}

// GetByCategory retrieves configurations by category
func (r *PayrollConfigRepository) GetByCategory(category models.PayrollConfigCategory) ([]models.PayrollConfig, error) {
	var configs []models.PayrollConfig
	err := r.db.Where("category = ?", category).Find(&configs).Error
	return configs, err
}

// GetAll retrieves all configurations
func (r *PayrollConfigRepository) GetAll() ([]models.PayrollConfig, error) {
	var configs []models.PayrollConfig
	err := r.db.Find(&configs).Error
	return configs, err
}

// GetAllNonEditable retrieves all non-editable configurations
func (r *PayrollConfigRepository) GetAllNonEditable() ([]models.PayrollConfig, error) {
	var configs []models.PayrollConfig
	err := r.db.Where("editable = ?", false).Find(&configs).Error
	return configs, err
}

// GetValue retrieves configuration value by key
func (r *PayrollConfigRepository) GetValue(configKey string) (string, error) {
	config, err := r.GetByKey(configKey)
	if err != nil {
		return "", err
	}
	return config.Value, nil
}

// GetPAYEBands retrieves all PAYE tax band configurations
func (r *PayrollConfigRepository) GetPAYEBands() ([]models.PayrollConfig, error) {
	var configs []models.PayrollConfig
	err := r.db.Where("category = ? AND config_key LIKE ?", models.ConfigCategoryPAYE, "PAYE_BAND_%").
		Order("config_key ASC").Find(&configs).Error
	return configs, err
}

// GetNSSFRates retrieves all NSSF rate configurations
func (r *PayrollConfigRepository) GetNSSFRates() ([]models.PayrollConfig, error) {
	var configs []models.PayrollConfig
	err := r.db.Where("category = ? AND config_key LIKE ?", models.ConfigCategoryNSSF, "NSSF_%").
		Find(&configs).Error
	return configs, err
}

// GetEmployerContributions retrieves all employer contribution rates
func (r *PayrollConfigRepository) GetEmployerContributions() ([]models.PayrollConfig, error) {
	var configs []models.PayrollConfig
	err := r.db.Where("config_key IN ?", []string{
		"NSSF_EMPLOYER_RATE",
		"WCF_EMPLOYER_RATE", 
		"SDL_EMPLOYER_RATE",
	}).Find(&configs).Error
	return configs, err
}

// GetFormulaConfigurations retrieves all formula configurations
func (r *PayrollConfigRepository) GetFormulaConfigurations() ([]models.PayrollConfig, error) {
	var configs []models.PayrollConfig
	err := r.db.Where("category = ?", models.ConfigCategoryFormula).Find(&configs).Error
	return configs, err
}

// Count returns the total number of configurations
func (r *PayrollConfigRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.PayrollConfig{}).Count(&count).Error
	return count, err
}

// CountByCategory returns the number of configurations by category
func (r *PayrollConfigRepository) CountByCategory(category models.PayrollConfigCategory) (int64, error) {
	var count int64
	err := r.db.Model(&models.PayrollConfig{}).Where("category = ?", category).Count(&count).Error
	return count, err
}

// Note: No Create, Update, or Delete methods are provided as these configurations
// are read-only and can only be modified through database migrations
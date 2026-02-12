package repositories

import (
	"errors"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/models"
	"gorm.io/gorm"
)

// BioTimeConfigRepository handles BioTime config database operations
type BioTimeConfigRepository struct {
	db *gorm.DB
}

// NewBioTimeConfigRepository creates a new BioTime config repository
func NewBioTimeConfigRepository() *BioTimeConfigRepository {
	return &BioTimeConfigRepository{
		db: database.GetDB(),
	}
}

// FindByTenantID finds BioTime config by tenant ID
func (r *BioTimeConfigRepository) FindByTenantID(tenantID *uint) (*models.BioTimeConfig, error) {
	var config models.BioTimeConfig
	query := r.db.Model(&models.BioTimeConfig{})
	
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}
	
	err := query.First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// FindOrCreate finds existing config or creates a new one
func (r *BioTimeConfigRepository) FindOrCreate(tenantID *uint, baseURL, username, password string) (*models.BioTimeConfig, error) {
	config, err := r.FindByTenantID(tenantID)
	if err == nil {
		// Update if values changed
		if config.BaseURL != baseURL || config.Username != username || config.Password != password {
			config.BaseURL = baseURL
			config.Username = username
			config.Password = password
			if err := r.Update(config); err != nil {
				return nil, err
			}
		}
		return config, nil
	}
	
	// Create new config if not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		config = &models.BioTimeConfig{
			BaseURL:  baseURL,
			Username: username,
			Password: password,
			Enabled:  true,
		}
		if err := r.Create(config); err != nil {
			return nil, err
		}
		return config, nil
	}
	
	return nil, err
}

// Create creates a new BioTime config
func (r *BioTimeConfigRepository) Create(config *models.BioTimeConfig) error {
	return r.db.Create(config).Error
}

// Update updates a BioTime config
func (r *BioTimeConfigRepository) Update(config *models.BioTimeConfig) error {
	return r.db.Save(config).Error
}

// UpdateToken updates only the token and expiry
func (r *BioTimeConfigRepository) UpdateToken(tenantID *uint, token string, expiry time.Time) error {
	query := r.db.Model(&models.BioTimeConfig{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}
	
	return query.Updates(map[string]interface{}{
		"token":        token,
		"token_expiry": expiry,
		"last_sync":    time.Now(),
		"last_error":   nil,
	}).Error
}

// UpdateLastError updates the last error message
func (r *BioTimeConfigRepository) UpdateLastError(tenantID *uint, errorMsg string) error {
	query := r.db.Model(&models.BioTimeConfig{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}
	
	return query.Update("last_error", errorMsg).Error
}

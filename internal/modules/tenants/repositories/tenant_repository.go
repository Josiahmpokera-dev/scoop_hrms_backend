package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/tenants/models"
	"gorm.io/gorm"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository() *TenantRepository {
	return &TenantRepository{
		db: database.GetDB(),
	}
}

// Create creates a new tenant
func (r *TenantRepository) Create(tenant *models.Tenant) error {
	return r.db.Create(tenant).Error
}

// FindByID finds a tenant by ID
func (r *TenantRepository) FindByID(id uint) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.First(&tenant, id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// FindByDomain finds a tenant by domain
func (r *TenantRepository) FindByDomain(domain string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.Where("domain = ?", domain).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// Update updates a tenant
func (r *TenantRepository) Update(tenant *models.Tenant) error {
	return r.db.Save(tenant).Error
}

// Delete soft deletes a tenant
func (r *TenantRepository) Delete(id uint) error {
	return r.db.Delete(&models.Tenant{}, id).Error
}

// List returns all tenants with pagination
func (r *TenantRepository) List(page, pageSize int) ([]models.Tenant, int64, error) {
	var tenants []models.Tenant
	var total int64

	offset := (page - 1) * pageSize

	// Count total
	if err := r.db.Model(&models.Tenant{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Offset(offset).Limit(pageSize).Find(&tenants).Error
	return tenants, total, err
}

// ExistsByDomain checks if a tenant with the given domain exists
func (r *TenantRepository) ExistsByDomain(domain string) bool {
	var count int64
	r.db.Model(&models.Tenant{}).Where("domain = ?", domain).Count(&count)
	return count > 0
}

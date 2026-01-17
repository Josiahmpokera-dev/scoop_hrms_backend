package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/models"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{
		db: database.GetDB(),
	}
}

// Create creates a new organization
func (r *OrganizationRepository) Create(organization *models.Organization) error {
	return r.db.Create(organization).Error
}

// FindByID finds an organization by ID
func (r *OrganizationRepository) FindByID(id uint) (*models.Organization, error) {
	var organization models.Organization
	err := r.db.First(&organization, id).Error
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

// FindByCode finds an organization by code
func (r *OrganizationRepository) FindByCode(code string) (*models.Organization, error) {
	var organization models.Organization
	err := r.db.Where("code = ?", code).First(&organization).Error
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

// Update updates an organization
func (r *OrganizationRepository) Update(organization *models.Organization) error {
	return r.db.Save(organization).Error
}

// Delete soft deletes an organization
func (r *OrganizationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Organization{}, id).Error
}

// List returns all organizations with pagination
func (r *OrganizationRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Organization, int64, error) {
	var organizations []models.Organization
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Organization{})

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply additional filters
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if country, ok := filters["country"].(string); ok && country != "" {
		query = query.Where("country = ?", country)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(pageSize).Find(&organizations).Error
	return organizations, total, err
}

// ExistsByCode checks if an organization with the given code exists
func (r *OrganizationRepository) ExistsByCode(code string) bool {
	var count int64
	r.db.Model(&models.Organization{}).Where("code = ?", code).Count(&count)
	return count > 0
}

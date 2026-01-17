package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/models"
	"gorm.io/gorm"
)

type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository() *LocationRepository {
	return &LocationRepository{
		db: database.GetDB(),
	}
}

// Create creates a new location
func (r *LocationRepository) Create(location *models.Location) error {
	return r.db.Create(location).Error
}

// FindByID finds a location by ID
func (r *LocationRepository) FindByID(id uint) (*models.Location, error) {
	var location models.Location
	err := r.db.First(&location, id).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

// Update updates a location
func (r *LocationRepository) Update(location *models.Location) error {
	return r.db.Save(location).Error
}

// Delete soft deletes a location
func (r *LocationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Location{}, id).Error
}

// List returns all locations with pagination
func (r *LocationRepository) List(tenantID, organizationID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Location, int64, error) {
	var locations []models.Location
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Location{})

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply organization filter
	if organizationID != nil {
		query = query.Where("organization_id = ?", *organizationID)
	}

	// Apply additional filters
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	if isHeadOffice, ok := filters["is_head_office"].(bool); ok {
		query = query.Where("is_head_office = ?", isHeadOffice)
	}
	if locationType, ok := filters["location_type"].(string); ok && locationType != "" {
		query = query.Where("location_type = ?", locationType)
	}
	if country, ok := filters["country"].(string); ok && country != "" {
		query = query.Where("country = ?", country)
	}
	if city, ok := filters["city"].(string); ok && city != "" {
		query = query.Where("city = ?", city)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(pageSize).Find(&locations).Error
	return locations, total, err
}

// FindHeadOffice finds the head office location
func (r *LocationRepository) FindHeadOffice(tenantID *uint) (*models.Location, error) {
	var location models.Location
	query := r.db.Where("is_head_office = ?", true)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

// FindByName finds a location by name
// Searches within tenant scope, and optionally within organization scope
func (r *LocationRepository) FindByName(name string, tenantID *uint, organizationID *uint) (*models.Location, error) {
	var location models.Location
	query := r.db.Where("LOWER(name) = LOWER(?)", name) // Case-insensitive search

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// If organizationID is provided, search within that organization
	// Otherwise, search across all organizations in the tenant
	if organizationID != nil {
		query = query.Where("organization_id = ?", *organizationID)
	}

	err := query.First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

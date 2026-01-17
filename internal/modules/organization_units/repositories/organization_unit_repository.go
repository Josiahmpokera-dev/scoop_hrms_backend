package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization_units/models"
	"gorm.io/gorm"
)

type OrganizationUnitRepository struct {
	db *gorm.DB
}

func NewOrganizationUnitRepository() *OrganizationUnitRepository {
	return &OrganizationUnitRepository{
		db: database.GetDB(),
	}
}

// Create creates a new organization unit
func (r *OrganizationUnitRepository) Create(unit *models.OrganizationUnit) error {
	return r.db.Create(unit).Error
}

// FindByID finds an organization unit by ID
func (r *OrganizationUnitRepository) FindByID(id uint) (*models.OrganizationUnit, error) {
	var unit models.OrganizationUnit
	err := r.db.Preload("ParentUnit").Preload("SubUnits").First(&unit, id).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// FindByCode finds an organization unit by code
func (r *OrganizationUnitRepository) FindByCode(code string) (*models.OrganizationUnit, error) {
	var unit models.OrganizationUnit
	err := r.db.Where("code = ?", code).First(&unit).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// Update updates an organization unit
func (r *OrganizationUnitRepository) Update(unit *models.OrganizationUnit) error {
	return r.db.Save(unit).Error
}

// Delete soft deletes an organization unit
func (r *OrganizationUnitRepository) Delete(id uint) error {
	return r.db.Delete(&models.OrganizationUnit{}, id).Error
}

// List returns all organization units with pagination
func (r *OrganizationUnitRepository) List(tenantID, organizationID *uint, page, pageSize int, filters map[string]interface{}) ([]models.OrganizationUnit, int64, error) {
	var units []models.OrganizationUnit
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.OrganizationUnit{})

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
	if unitType, ok := filters["unit_type"].(string); ok && unitType != "" {
		query = query.Where("unit_type = ?", unitType)
	}
	if parentUnitID, ok := filters["parent_unit_id"]; ok {
		if parentUnitID == nil {
			query = query.Where("parent_unit_id IS NULL")
		} else {
			query = query.Where("parent_unit_id = ?", parentUnitID)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Preload("ParentUnit").Offset(offset).Limit(pageSize).Find(&units).Error
	return units, total, err
}

// FindRootUnits finds all root organization units (no parent)
func (r *OrganizationUnitRepository) FindRootUnits(organizationID *uint) ([]models.OrganizationUnit, error) {
	var units []models.OrganizationUnit
	query := r.db.Where("parent_unit_id IS NULL")

	if organizationID != nil {
		query = query.Where("organization_id = ?", *organizationID)
	}

	err := query.Preload("SubUnits").Find(&units).Error
	return units, err
}

// ExistsByCode checks if an organization unit with the given code exists
func (r *OrganizationUnitRepository) ExistsByCode(code string) bool {
	var count int64
	r.db.Model(&models.OrganizationUnit{}).Where("code = ?", code).Count(&count)
	return count > 0
}

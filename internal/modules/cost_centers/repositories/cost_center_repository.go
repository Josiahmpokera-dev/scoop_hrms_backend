package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/cost_centers/models"
	"gorm.io/gorm"
)

type CostCenterRepository struct {
	db *gorm.DB
}

func NewCostCenterRepository() *CostCenterRepository {
	return &CostCenterRepository{
		db: database.GetDB(),
	}
}

// Create creates a new cost center
func (r *CostCenterRepository) Create(costCenter *models.CostCenter) error {
	return r.db.Create(costCenter).Error
}

// FindByID finds a cost center by ID
func (r *CostCenterRepository) FindByID(id uint) (*models.CostCenter, error) {
	var costCenter models.CostCenter
	err := r.db.Preload("ParentCostCenter").Preload("SubCostCenters").First(&costCenter, id).Error
	if err != nil {
		return nil, err
	}
	return &costCenter, nil
}

// FindByCode finds a cost center by code
func (r *CostCenterRepository) FindByCode(code string) (*models.CostCenter, error) {
	var costCenter models.CostCenter
	err := r.db.Where("code = ?", code).First(&costCenter).Error
	if err != nil {
		return nil, err
	}
	return &costCenter, nil
}

// Update updates a cost center
func (r *CostCenterRepository) Update(costCenter *models.CostCenter) error {
	return r.db.Save(costCenter).Error
}

// Delete soft deletes a cost center
func (r *CostCenterRepository) Delete(id uint) error {
	return r.db.Delete(&models.CostCenter{}, id).Error
}

// List returns all cost centers with pagination
func (r *CostCenterRepository) List(tenantID, organizationID *uint, page, pageSize int, filters map[string]interface{}) ([]models.CostCenter, int64, error) {
	var costCenters []models.CostCenter
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.CostCenter{})

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
	if costCenterType, ok := filters["cost_center_type"].(string); ok && costCenterType != "" {
		query = query.Where("cost_center_type = ?", costCenterType)
	}
	if departmentID, ok := filters["department_id"]; ok {
		query = query.Where("department_id = ?", departmentID)
	}
	if parentCostCenterID, ok := filters["parent_cost_center_id"]; ok {
		if parentCostCenterID == nil {
			query = query.Where("parent_cost_center_id IS NULL")
		} else {
			query = query.Where("parent_cost_center_id = ?", parentCostCenterID)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Preload("ParentCostCenter").Offset(offset).Limit(pageSize).Find(&costCenters).Error
	return costCenters, total, err
}

// FindRootCostCenters finds all root cost centers (no parent)
func (r *CostCenterRepository) FindRootCostCenters(organizationID *uint) ([]models.CostCenter, error) {
	var costCenters []models.CostCenter
	query := r.db.Where("parent_cost_center_id IS NULL")

	if organizationID != nil {
		query = query.Where("organization_id = ?", *organizationID)
	}

	err := query.Preload("SubCostCenters").Find(&costCenters).Error
	return costCenters, err
}

// ExistsByCode checks if a cost center with the given code exists
func (r *CostCenterRepository) ExistsByCode(code string) bool {
	var count int64
	r.db.Model(&models.CostCenter{}).Where("code = ?", code).Count(&count)
	return count > 0
}

// FindByName finds a cost center by name
func (r *CostCenterRepository) FindByName(name string, tenantID *uint) (*models.CostCenter, error) {
	var costCenter models.CostCenter
	query := r.db.Where("name = ?", name)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.First(&costCenter).Error
	if err != nil {
		return nil, err
	}
	return &costCenter, nil
}

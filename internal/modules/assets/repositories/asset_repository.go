package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	"gorm.io/gorm"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository() *AssetRepository {
	return &AssetRepository{
		db: database.GetDB(),
	}
}

// Create creates a new asset
func (r *AssetRepository) Create(asset *models.Asset) error {
	return r.db.Create(asset).Error
}

// FindByID finds an asset by ID
func (r *AssetRepository) FindByID(id uint) (*models.Asset, error) {
	var asset models.Asset
	err := r.db.First(&asset, id).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// FindByAssetCode finds an asset by asset code
func (r *AssetRepository) FindByAssetCode(assetCode string) (*models.Asset, error) {
	var asset models.Asset
	err := r.db.Where("asset_code = ?", assetCode).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// FindBySerialNumber finds an asset by serial number
func (r *AssetRepository) FindBySerialNumber(serialNumber string) (*models.Asset, error) {
	var asset models.Asset
	err := r.db.Where("serial_number = ?", serialNumber).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// ExistsByAssetCode checks if an asset with the given code exists
func (r *AssetRepository) ExistsByAssetCode(assetCode string) bool {
	var count int64
	r.db.Model(&models.Asset{}).Where("asset_code = ?", assetCode).Count(&count)
	return count > 0
}

// ExistsBySerialNumber checks if an asset with the given serial number exists
func (r *AssetRepository) ExistsBySerialNumber(serialNumber string) bool {
	var count int64
	r.db.Model(&models.Asset{}).Where("serial_number = ?", serialNumber).Count(&count)
	return count > 0
}

// Update updates an asset
func (r *AssetRepository) Update(asset *models.Asset) error {
	return r.db.Save(asset).Error
}

// Delete permanently deletes an asset (only for retired assets)
func (r *AssetRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&models.Asset{}, id).Error
}

// List returns all assets with pagination and filters
func (r *AssetRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Asset, int64, error) {
	var assets []models.Asset
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Asset{})

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply search filter
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where(
			"asset_code ILIKE ? OR serial_number ILIKE ? OR brand ILIKE ? OR model ILIKE ? OR assigned_to ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}

	// Apply status filter
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	// Apply asset type filter
	if assetType, ok := filters["asset_type"].(string); ok && assetType != "" {
		query = query.Where("asset_type = ?", assetType)
	}

	// Apply assigned_to filter (employee ID)
	if assignedTo, ok := filters["assigned_to"].(string); ok && assignedTo != "" {
		query = query.Where("employee_id = ?", assignedTo)
	}

	// Apply department filter
	if department, ok := filters["department"].(string); ok && department != "" {
		query = query.Where("department = ?", department)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}

// GetNextAssetCodeNumber gets the next sequential number for an asset type prefix
func (r *AssetRepository) GetNextAssetCodeNumber(prefix string, tenantID *uint) (int, error) {
	var maxNumber int
	query := r.db.Model(&models.Asset{}).
		Where("asset_code LIKE ?", prefix+"-%")

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.Select("COALESCE(MAX(CAST(SUBSTRING(asset_code FROM '[0-9]+$') AS INTEGER)), 0)").
		Scan(&maxNumber).Error

	if err != nil {
		return 1, err
	}

	return maxNumber + 1, nil
}

// FindByEmployeeID finds assets assigned to an employee
func (r *AssetRepository) FindByEmployeeID(employeeID string, tenantID *uint) ([]models.Asset, error) {
	var assets []models.Asset
	query := r.db.Where("employee_id = ? AND status = ?", employeeID, "in_use")
	
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	err := query.Order("assigned_date DESC").Find(&assets).Error
	return assets, err
}

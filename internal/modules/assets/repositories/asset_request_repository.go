package repositories

import (
	"fmt"
	"time"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	"gorm.io/gorm"
)

// AssetRequestRepository handles asset request database operations
type AssetRequestRepository struct {
	db *gorm.DB
}

// NewAssetRequestRepository creates a new asset request repository
func NewAssetRequestRepository() *AssetRequestRepository {
	return &AssetRequestRepository{
		db: database.GetDB(),
	}
}

// Create creates a new asset request
func (r *AssetRequestRepository) Create(request *models.AssetRequest) error {
	return r.db.Create(request).Error
}

// FindByID finds an asset request by ID
func (r *AssetRequestRepository) FindByID(id uint) (*models.AssetRequest, error) {
	var request models.AssetRequest
	err := r.db.Where("id = ?", id).First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// FindByEmployeeID finds asset requests by employee ID
func (r *AssetRequestRepository) FindByEmployeeID(employeeID string, tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.AssetRequest, int64, error) {
	var requests []models.AssetRequest
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.AssetRequest{}).Where("employee_id = ?", employeeID)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if requestType, ok := filters["request_type"].(string); ok && requestType != "" {
		query = query.Where("request_type = ?", requestType)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if assetType, ok := filters["asset_type"].(string); ok && assetType != "" {
		query = query.Where("asset_type = ?", assetType)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error
	return requests, total, err
}

// Update updates an asset request
func (r *AssetRequestRepository) Update(request *models.AssetRequest) error {
	return r.db.Save(request).Error
}

// GenerateRequestNumber generates a unique request number
func (r *AssetRequestRepository) GenerateRequestNumber(tenantID *uint) (string, error) {
	var count int64
	query := r.db.Model(&models.AssetRequest{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Count(&count).Error; err != nil {
		return "", err
	}
	
	// Format: AR-YYYY-XXX (e.g., AR-2026-001)
	year := time.Now().Format("2006")
	sequence := int(count) + 1
	return "AR-" + year + "-" + formatAssetRequestSequence(sequence), nil
}

// ListAll lists all asset requests with pagination and filters (for HR/Admin)
func (r *AssetRequestRepository) ListAll(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.AssetRequest, int64, error) {
	var requests []models.AssetRequest
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.AssetRequest{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if employeeID, ok := filters["employee_id"].(string); ok && employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}
	if requestType, ok := filters["request_type"].(string); ok && requestType != "" {
		query = query.Where("request_type = ?", requestType)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if assetType, ok := filters["asset_type"].(string); ok && assetType != "" {
		query = query.Where("asset_type = ?", assetType)
	}
	if priority, ok := filters["priority"].(string); ok && priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error
	return requests, total, err
}

// formatAssetRequestSequence formats sequence number with leading zeros
func formatAssetRequestSequence(seq int) string {
	if seq < 10 {
		return "00" + fmt.Sprintf("%d", seq)
	} else if seq < 100 {
		return "0" + fmt.Sprintf("%d", seq)
	}
	return fmt.Sprintf("%d", seq)
}

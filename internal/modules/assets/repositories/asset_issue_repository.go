package repositories

import (
	"fmt"
	"time"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	"gorm.io/gorm"
)

// AssetIssueRepository handles asset issue database operations
type AssetIssueRepository struct {
	db *gorm.DB
}

// NewAssetIssueRepository creates a new asset issue repository
func NewAssetIssueRepository() *AssetIssueRepository {
	return &AssetIssueRepository{
		db: database.GetDB(),
	}
}

// Create creates a new asset issue
func (r *AssetIssueRepository) Create(issue *models.AssetIssue) error {
	return r.db.Create(issue).Error
}

// FindByID finds an asset issue by ID
func (r *AssetIssueRepository) FindByID(id uint) (*models.AssetIssue, error) {
	var issue models.AssetIssue
	err := r.db.Where("id = ?", id).First(&issue).Error
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// FindByEmployeeID finds asset issues by employee ID
func (r *AssetIssueRepository) FindByEmployeeID(employeeID string, tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.AssetIssue, int64, error) {
	var issues []models.AssetIssue
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.AssetIssue{}).Where("employee_id = ?", employeeID)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if issueType, ok := filters["issue_type"].(string); ok && issueType != "" {
		query = query.Where("issue_type = ?", issueType)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if assetID, ok := filters["asset_id"].(uint); ok && assetID > 0 {
		query = query.Where("asset_id = ?", assetID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&issues).Error
	return issues, total, err
}

// Update updates an asset issue
func (r *AssetIssueRepository) Update(issue *models.AssetIssue) error {
	return r.db.Save(issue).Error
}

// GenerateIssueNumber generates a unique issue number
func (r *AssetIssueRepository) GenerateIssueNumber(tenantID *uint) (string, error) {
	var count int64
	query := r.db.Model(&models.AssetIssue{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Count(&count).Error; err != nil {
		return "", err
	}
	
	// Format: AI-YYYY-XXX (e.g., AI-2026-001)
	year := time.Now().Format("2006")
	sequence := int(count) + 1
	return "AI-" + year + "-" + formatAssetIssueSequence(sequence), nil
}

// formatAssetIssueSequence formats sequence number with leading zeros
func formatAssetIssueSequence(seq int) string {
	if seq < 10 {
		return "00" + fmt.Sprintf("%d", seq)
	} else if seq < 100 {
		return "0" + fmt.Sprintf("%d", seq)
	}
	return fmt.Sprintf("%d", seq)
}

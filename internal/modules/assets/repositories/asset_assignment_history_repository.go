package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	"gorm.io/gorm"
)

type AssetAssignmentHistoryRepository struct {
	db *gorm.DB
}

func NewAssetAssignmentHistoryRepository() *AssetAssignmentHistoryRepository {
	return &AssetAssignmentHistoryRepository{db: database.GetDB()}
}

func (r *AssetAssignmentHistoryRepository) Create(entry *models.AssetAssignmentHistory) error {
	return r.db.Create(entry).Error
}

func (r *AssetAssignmentHistoryRepository) List(
	tenantID *uint,
	page, pageSize int,
	assetID *uint,
	employeeID *string,
	action *string,
	startDate, endDate *time.Time,
) ([]models.AssetAssignmentHistory, int64, error) {
	var rows []models.AssetAssignmentHistory
	var total int64

	query := r.db.Model(&models.AssetAssignmentHistory{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}
	if assetID != nil && *assetID > 0 {
		query = query.Where("asset_id = ?", *assetID)
	}
	if employeeID != nil && *employeeID != "" {
		query = query.Where("from_employee_id = ? OR to_employee_id = ?", *employeeID, *employeeID)
	}
	if action != nil && *action != "" {
		query = query.Where("action = ?", *action)
	}
	if startDate != nil {
		query = query.Where("occurred_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("occurred_at <= ?", *endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("occurred_at DESC, id DESC").Offset(offset).Limit(pageSize).Find(&rows).Error
	return rows, total, err
}


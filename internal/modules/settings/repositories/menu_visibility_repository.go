package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/settings/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MenuVisibilityRepository struct {
	db *gorm.DB
}

func NewMenuVisibilityRepository() *MenuVisibilityRepository {
	return &MenuVisibilityRepository{db: database.GetDB()}
}

// GetAll returns every stored setting.
func (r *MenuVisibilityRepository) GetAll() ([]models.MenuVisibilitySetting, error) {
	var items []models.MenuVisibilitySetting
	if err := r.db.Order("menu_key ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetHiddenKeys returns only the menu_key values where visible = false.
func (r *MenuVisibilityRepository) GetHiddenKeys() ([]string, error) {
	var keys []string
	if err := r.db.Model(&models.MenuVisibilitySetting{}).
		Where("visible = ?", false).
		Pluck("menu_key", &keys).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

// UpsertMany inserts or updates a batch of visibility settings.
func (r *MenuVisibilityRepository) UpsertMany(items []models.MenuVisibilitySetting) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "menu_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"visible", "updated_by", "updated_at"}),
	}).Create(&items).Error
}

// UpsertOne inserts or updates a single setting.
func (r *MenuVisibilityRepository) UpsertOne(item *models.MenuVisibilitySetting) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "menu_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"visible", "updated_by", "updated_at"}),
	}).Create(item).Error
}

// ResetAll sets every stored row to visible = true.
func (r *MenuVisibilityRepository) ResetAll() (int64, error) {
	result := r.db.Model(&models.MenuVisibilitySetting{}).
		Where("visible = ?", false).
		Update("visible", true)
	return result.RowsAffected, result.Error
}

// CountAll returns the total number of stored settings.
func (r *MenuVisibilityRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.MenuVisibilitySetting{}).Count(&count).Error
	return count, err
}

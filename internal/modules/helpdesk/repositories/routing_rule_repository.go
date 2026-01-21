package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	"gorm.io/gorm"
)

// RoutingRuleRepository handles routing rule database operations
type RoutingRuleRepository struct {
	db *gorm.DB
}

// NewRoutingRuleRepository creates a new routing rule repository
func NewRoutingRuleRepository() *RoutingRuleRepository {
	return &RoutingRuleRepository{
		db: database.GetDB(),
	}
}

// Create creates a new routing rule
func (r *RoutingRuleRepository) Create(rule *models.RoutingRule) error {
	return r.db.Create(rule).Error
}

// FindByID finds a routing rule by ID
func (r *RoutingRuleRepository) FindByID(id uint) (*models.RoutingRule, error) {
	var rule models.RoutingRule
	err := r.db.Where("id = ?", id).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// ListAll lists all routing rules
func (r *RoutingRuleRepository) ListAll(tenantID *uint) ([]models.RoutingRule, error) {
	var rules []models.RoutingRule
	query := r.db.Model(&models.RoutingRule{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.Order("created_at DESC").Find(&rules).Error
	return rules, err
}

// FindActiveRules finds active routing rules
func (r *RoutingRuleRepository) FindActiveRules(tenantID *uint) ([]models.RoutingRule, error) {
	var rules []models.RoutingRule
	query := r.db.Where("is_active = ?", true)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.Order("created_at DESC").Find(&rules).Error
	return rules, err
}

// Update updates a routing rule
func (r *RoutingRuleRepository) Update(rule *models.RoutingRule) error {
	return r.db.Save(rule).Error
}

// Delete deletes a routing rule
func (r *RoutingRuleRepository) Delete(id uint) error {
	return r.db.Delete(&models.RoutingRule{}, id).Error
}

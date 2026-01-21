package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	"gorm.io/gorm"
)

// TicketCategoryRepository handles ticket category database operations
type TicketCategoryRepository struct {
	db *gorm.DB
}

// NewTicketCategoryRepository creates a new ticket category repository
func NewTicketCategoryRepository() *TicketCategoryRepository {
	return &TicketCategoryRepository{
		db: database.GetDB(),
	}
}

// Create creates a new ticket category
func (r *TicketCategoryRepository) Create(category *models.TicketCategory) error {
	return r.db.Create(category).Error
}

// FindByID finds a ticket category by ID
func (r *TicketCategoryRepository) FindByID(id uint) (*models.TicketCategory, error) {
	var category models.TicketCategory
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindByName finds a ticket category by name
func (r *TicketCategoryRepository) FindByName(name string, tenantID *uint) (*models.TicketCategory, error) {
	var category models.TicketCategory
	query := r.db.Where("name = ?", name)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// ListAll lists all ticket categories
func (r *TicketCategoryRepository) ListAll(tenantID *uint) ([]models.TicketCategory, error) {
	var categories []models.TicketCategory
	query := r.db.Model(&models.TicketCategory{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.Order("name ASC").Find(&categories).Error
	return categories, err
}

// Update updates a ticket category
func (r *TicketCategoryRepository) Update(category *models.TicketCategory) error {
	return r.db.Save(category).Error
}

// Delete deletes a ticket category
func (r *TicketCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.TicketCategory{}, id).Error
}

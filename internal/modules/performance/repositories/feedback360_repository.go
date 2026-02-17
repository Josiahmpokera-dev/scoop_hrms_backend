package repositories

import (
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"gorm.io/gorm"
)

// Feedback360Repository handles database operations for 360° feedback campaigns and rater groups.
type Feedback360Repository struct {
	db *gorm.DB
}

// NewFeedback360Repository creates a new Feedback360Repository.
func NewFeedback360Repository() *Feedback360Repository {
	return &Feedback360Repository{db: database.GetDB()}
}

// Create creates a new 360° feedback campaign.
func (r *Feedback360Repository) Create(campaign *models.Feedback360Campaign) error {
	return r.db.Create(campaign).Error
}

// FindByID finds a campaign by ID with RaterGroups preloaded.
func (r *Feedback360Repository) FindByID(id uint) (*models.Feedback360Campaign, error) {
	var campaign models.Feedback360Campaign
	err := r.db.Preload("RaterGroups").First(&campaign, id).Error
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

// Update updates a 360° feedback campaign.
func (r *Feedback360Repository) Update(campaign *models.Feedback360Campaign) error {
	return r.db.Save(campaign).Error
}

// List returns campaigns with pagination and optional filters; preloads RaterGroups.
func (r *Feedback360Repository) List(status, department string, employeeID *uint, page, pageSize int) ([]models.Feedback360Campaign, int64, error) {
	var campaigns []models.Feedback360Campaign
	var total int64

	query := r.db.Model(&models.Feedback360Campaign{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if department != "" {
		query = query.Where("department = ?", department)
	}
	// employeeID filter: if your schema has a participant/employee link table, join here; otherwise omit
	if employeeID != nil {
		// Optional: query = query.Where("id IN (SELECT campaign_id FROM performance_feedback360_participants WHERE employee_id = ?)", *employeeID)
		// Leaving as no filter when employeeID is set unless we have a participants table
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("RaterGroups").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&campaigns).Error

	return campaigns, total, err
}

// GetNextCode returns the next campaign code (e.g. F360-001).
func (r *Feedback360Repository) GetNextCode() string {
	var count int64
	r.db.Model(&models.Feedback360Campaign{}).Count(&count)
	return fmt.Sprintf("F360-%03d", count+1)
}

// CreateRaterGroup creates a new rater group for a campaign.
func (r *Feedback360Repository) CreateRaterGroup(group *models.Feedback360RaterGroup) error {
	return r.db.Create(group).Error
}

// UpdateRaterGroup updates a rater group.
func (r *Feedback360Repository) UpdateRaterGroup(group *models.Feedback360RaterGroup) error {
	return r.db.Save(group).Error
}

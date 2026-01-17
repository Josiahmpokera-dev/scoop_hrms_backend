package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/models"
	"gorm.io/gorm"
)

type JobPositionRepository struct {
	db *gorm.DB
}

func NewJobPositionRepository() *JobPositionRepository {
	return &JobPositionRepository{
		db: database.GetDB(),
	}
}

// Create creates a new job position
func (r *JobPositionRepository) Create(position *models.JobPosition) error {
	return r.db.Create(position).Error
}

// FindByID finds a job position by ID
func (r *JobPositionRepository) FindByID(id uint) (*models.JobPosition, error) {
	var position models.JobPosition
	err := r.db.First(&position, id).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// FindByCode finds a job position by code
func (r *JobPositionRepository) FindByCode(code string) (*models.JobPosition, error) {
	var position models.JobPosition
	err := r.db.Where("code = ?", code).First(&position).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// Update updates a job position
func (r *JobPositionRepository) Update(position *models.JobPosition) error {
	return r.db.Save(position).Error
}

// Delete soft deletes a job position
func (r *JobPositionRepository) Delete(id uint) error {
	return r.db.Delete(&models.JobPosition{}, id).Error
}

// List returns all job positions with pagination
func (r *JobPositionRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.JobPosition, int64, error) {
	var positions []models.JobPosition
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.JobPosition{})

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply additional filters
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	if grade, ok := filters["grade"].(string); ok && grade != "" {
		query = query.Where("grade = ?", grade)
	}
	if level, ok := filters["level"].(int); ok {
		query = query.Where("level = ?", level)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(pageSize).Find(&positions).Error
	return positions, total, err
}

// ExistsByCode checks if a job position with the given code exists
func (r *JobPositionRepository) ExistsByCode(code string) bool {
	var count int64
	r.db.Model(&models.JobPosition{}).Where("code = ?", code).Count(&count)
	return count > 0
}

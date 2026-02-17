package repositories

import (
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"gorm.io/gorm"
)

// DepartmentTargetRepository handles database operations for department targets and milestones.
type DepartmentTargetRepository struct {
	db *gorm.DB
}

// NewDepartmentTargetRepository creates a new DepartmentTargetRepository.
func NewDepartmentTargetRepository() *DepartmentTargetRepository {
	return &DepartmentTargetRepository{db: database.GetDB()}
}

// Create creates a new department target.
func (r *DepartmentTargetRepository) Create(target *models.DepartmentTarget) error {
	return r.db.Create(target).Error
}

// FindByID finds a department target by ID with Milestones preloaded.
func (r *DepartmentTargetRepository) FindByID(id uint) (*models.DepartmentTarget, error) {
	var target models.DepartmentTarget
	err := r.db.Preload("Milestones").First(&target, id).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

// Update updates a department target.
func (r *DepartmentTargetRepository) Update(target *models.DepartmentTarget) error {
	return r.db.Save(target).Error
}

// Delete soft deletes a department target.
func (r *DepartmentTargetRepository) Delete(id uint) error {
	return r.db.Delete(&models.DepartmentTarget{}, id).Error
}

// List returns department targets with pagination and optional filters.
func (r *DepartmentTargetRepository) List(department, category, status, period, search string, page, pageSize int) ([]models.DepartmentTarget, int64, error) {
	var targets []models.DepartmentTarget
	var total int64

	query := r.db.Model(&models.DepartmentTarget{})

	if department != "" {
		query = query.Where("department = ?", department)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if period != "" {
		query = query.Where("period = ?", period)
	}
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ? OR metric ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&targets).Error

	return targets, total, err
}

// GetNextCode returns the next department target code (e.g. DT-001).
func (r *DepartmentTargetRepository) GetNextCode() string {
	var count int64
	r.db.Model(&models.DepartmentTarget{}).Count(&count)
	return fmt.Sprintf("DT-%03d", count+1)
}

// FindMilestoneByID finds a department target milestone by ID.
func (r *DepartmentTargetRepository) FindMilestoneByID(id uint) (*models.DepartmentTargetMilestone, error) {
	var m models.DepartmentTargetMilestone
	err := r.db.First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// UpdateMilestone updates a department target milestone.
func (r *DepartmentTargetRepository) UpdateMilestone(m *models.DepartmentTargetMilestone) error {
	return r.db.Save(m).Error
}

// CreateMilestone creates a new department target milestone.
func (r *DepartmentTargetRepository) CreateMilestone(m *models.DepartmentTargetMilestone) error {
	return r.db.Create(m).Error
}

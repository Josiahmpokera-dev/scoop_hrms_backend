package repositories

import (
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"gorm.io/gorm"
)

// GoalRepository handles database operations for goals, key results, and check-ins.
type GoalRepository struct {
	db *gorm.DB
}

// NewGoalRepository creates a new GoalRepository.
func NewGoalRepository() *GoalRepository {
	return &GoalRepository{db: database.GetDB()}
}

// Create creates a new goal.
func (r *GoalRepository) Create(goal *models.Goal) error {
	return r.db.Create(goal).Error
}

// FindByID finds a goal by ID with KeyResults and CheckIns preloaded.
func (r *GoalRepository) FindByID(id uint) (*models.Goal, error) {
	var goal models.Goal
	err := r.db.Preload("KeyResults").Preload("CheckIns").First(&goal, id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// Update updates a goal.
func (r *GoalRepository) Update(goal *models.Goal) error {
	return r.db.Save(goal).Error
}

// Delete soft deletes a goal.
func (r *GoalRepository) Delete(id uint) error {
	return r.db.Delete(&models.Goal{}, id).Error
}

// List returns goals with pagination and optional filters.
func (r *GoalRepository) List(ownerID *uint, level, status, department, search string, page, pageSize int) ([]models.Goal, int64, error) {
	var goals []models.Goal
	var total int64

	query := r.db.Model(&models.Goal{})

	if ownerID != nil {
		query = query.Where("owner_id = ?", *ownerID)
	}
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if department != "" {
		query = query.Where("department = ?", department)
	}
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title ILIKE ? OR (description IS NOT NULL AND description::text ILIKE ?)", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("KeyResults").Preload("CheckIns").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&goals).Error

	return goals, total, err
}

// CountByStatus returns counts grouped by status for the given owner (or all if nil).
func (r *GoalRepository) CountByStatus(ownerID *uint) (map[string]int64, error) {
	type result struct {
		Status string
		Count  int64
	}
	var results []result
	query := r.db.Model(&models.Goal{}).Select("status, count(*) as count").Group("status")
	if ownerID != nil {
		query = query.Where("owner_id = ?", *ownerID)
	}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}
	counts := make(map[string]int64)
	for _, row := range results {
		counts[row.Status] = row.Count
	}
	return counts, nil
}

// GetNextGoalCode returns the next goal code (e.g. G-001).
func (r *GoalRepository) GetNextGoalCode() string {
	var count int64
	r.db.Model(&models.Goal{}).Count(&count)
	return fmt.Sprintf("G-%03d", count+1)
}

// GetNextKRCode returns the next key result code (e.g. KR-001).
func (r *GoalRepository) GetNextKRCode() string {
	var count int64
	r.db.Model(&models.KeyResult{}).Count(&count)
	return fmt.Sprintf("KR-%03d", count+1)
}

// CreateKeyResult creates a new key result.
func (r *GoalRepository) CreateKeyResult(kr *models.KeyResult) error {
	return r.db.Create(kr).Error
}

// FindKeyResultByID finds a key result by ID.
func (r *GoalRepository) FindKeyResultByID(id uint) (*models.KeyResult, error) {
	var kr models.KeyResult
	err := r.db.First(&kr, id).Error
	if err != nil {
		return nil, err
	}
	return &kr, nil
}

// UpdateKeyResult updates a key result.
func (r *GoalRepository) UpdateKeyResult(kr *models.KeyResult) error {
	return r.db.Save(kr).Error
}

// DeleteKeyResult soft deletes a key result.
func (r *GoalRepository) DeleteKeyResult(id uint) error {
	return r.db.Delete(&models.KeyResult{}, id).Error
}

// CreateCheckIn creates a new goal check-in.
func (r *GoalRepository) CreateCheckIn(checkIn *models.GoalCheckIn) error {
	return r.db.Create(checkIn).Error
}

// ListAllForAlignment returns all goals with minimal fields for alignment map (id, title, parent_goal_id, owner_id).
func (r *GoalRepository) ListAllForAlignment() ([]models.Goal, error) {
	var goals []models.Goal
	err := r.db.Select("id", "goal_code", "title", "parent_goal_id", "owner_id", "level").Find(&goals).Error
	return goals, err
}

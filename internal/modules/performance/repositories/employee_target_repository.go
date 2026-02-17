package repositories

import (
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"gorm.io/gorm"
)

// EmployeeTargetRepository handles database operations for employee targets.
type EmployeeTargetRepository struct {
	db *gorm.DB
}

// NewEmployeeTargetRepository creates a new EmployeeTargetRepository.
func NewEmployeeTargetRepository() *EmployeeTargetRepository {
	return &EmployeeTargetRepository{db: database.GetDB()}
}

// Create creates a new employee target.
func (r *EmployeeTargetRepository) Create(target *models.EmployeeTarget) error {
	return r.db.Create(target).Error
}

// FindByID finds an employee target by ID.
func (r *EmployeeTargetRepository) FindByID(id uint) (*models.EmployeeTarget, error) {
	var target models.EmployeeTarget
	err := r.db.First(&target, id).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

// Update updates an employee target.
func (r *EmployeeTargetRepository) Update(target *models.EmployeeTarget) error {
	return r.db.Save(target).Error
}

// Delete soft deletes an employee target.
func (r *EmployeeTargetRepository) Delete(id uint) error {
	return r.db.Delete(&models.EmployeeTarget{}, id).Error
}

// List returns employee targets with pagination and optional filters.
// employeeID, department, departmentTargetID, status, period, assignedBy are filter strings (e.g. ID as string).
func (r *EmployeeTargetRepository) List(employeeID, department, departmentTargetID, status, period, assignedBy string, page, pageSize int) ([]models.EmployeeTarget, int64, error) {
	var targets []models.EmployeeTarget
	var total int64

	query := r.db.Model(&models.EmployeeTarget{})

	if employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}
	if department != "" {
		query = query.Where("department = ?", department)
	}
	if departmentTargetID != "" {
		query = query.Where("department_target_id = ?", departmentTargetID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if period != "" {
		query = query.Where("period = ?", period)
	}
	if assignedBy != "" {
		query = query.Where("assigned_by_id = ?", assignedBy)
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

// GetNextCode returns the next employee target code (e.g. ET-001).
func (r *EmployeeTargetRepository) GetNextCode() string {
	var count int64
	r.db.Model(&models.EmployeeTarget{}).Count(&count)
	return fmt.Sprintf("ET-%03d", count+1)
}

// CreateBatch creates multiple employee targets in one call.
func (r *EmployeeTargetRepository) CreateBatch(targets []*models.EmployeeTarget) error {
	if len(targets) == 0 {
		return nil
	}
	return r.db.Create(&targets).Error
}

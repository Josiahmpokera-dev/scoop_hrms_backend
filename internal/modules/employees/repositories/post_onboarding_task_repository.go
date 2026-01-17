package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type PostOnboardingTaskRepository struct {
	db *gorm.DB
}

func NewPostOnboardingTaskRepository() *PostOnboardingTaskRepository {
	return &PostOnboardingTaskRepository{
		db: database.GetDB(),
	}
}

// Create creates a new post-onboarding task
func (r *PostOnboardingTaskRepository) Create(task *models.PostOnboardingTask) error {
	return r.db.Create(task).Error
}

// FindByID finds a task by ID
func (r *PostOnboardingTaskRepository) FindByID(id uint) (*models.PostOnboardingTask, error) {
	var task models.PostOnboardingTask
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// FindByEmployeeID finds all tasks for an employee
func (r *PostOnboardingTaskRepository) FindByEmployeeID(employeeID string) ([]models.PostOnboardingTask, error) {
	var tasks []models.PostOnboardingTask
	err := r.db.Where("employee_id = ?", employeeID).Order("created_at ASC").Find(&tasks).Error
	return tasks, err
}

// Update updates a task
func (r *PostOnboardingTaskRepository) Update(task *models.PostOnboardingTask) error {
	return r.db.Save(task).Error
}

// Delete soft deletes a task
func (r *PostOnboardingTaskRepository) Delete(id uint) error {
	return r.db.Delete(&models.PostOnboardingTask{}, id).Error
}

// List returns tasks with pagination and filters
func (r *PostOnboardingTaskRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.PostOnboardingTask, int64, error) {
	var tasks []models.PostOnboardingTask
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.PostOnboardingTask{})

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if employeeID, ok := filters["employee_id"].(string); ok && employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if taskType, ok := filters["task_type"].(string); ok && taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}
	if priority, ok := filters["priority"].(string); ok && priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if assignedTo, ok := filters["assigned_to"].(string); ok && assignedTo != "" {
		query = query.Where("assigned_to = ?", assignedTo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if err := query.Order("priority DESC, created_at ASC").Offset(offset).Limit(pageSize).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// CountByEmployeeID counts tasks for an employee by status
func (r *PostOnboardingTaskRepository) CountByEmployeeID(employeeID string) (map[string]int64, error) {
	var counts []struct {
		Status string
		Count  int64
	}

	err := r.db.Model(&models.PostOnboardingTask{}).
		Select("status, COUNT(*) as count").
		Where("employee_id = ?", employeeID).
		Group("status").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int64)
	for _, c := range counts {
		result[c.Status] = c.Count
	}

	return result, nil
}

package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/models"
	"gorm.io/gorm"
)

// DailyTaskRepository handles database operations for daily tasks
type DailyTaskRepository struct {
	db *gorm.DB
}

// NewDailyTaskRepository creates a new DailyTaskRepository
func NewDailyTaskRepository() *DailyTaskRepository {
	return &DailyTaskRepository{
		db: database.GetDB(),
	}
}

// Create creates a new daily task
func (r *DailyTaskRepository) Create(task *models.DailyTask) error {
	return r.db.Create(task).Error
}

// FindByID finds a daily task by ID
func (r *DailyTaskRepository) FindByID(id uint) (*models.DailyTask, error) {
	var task models.DailyTask
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Update updates a daily task
func (r *DailyTaskRepository) Update(task *models.DailyTask) error {
	return r.db.Save(task).Error
}

// Delete soft deletes a daily task
func (r *DailyTaskRepository) Delete(id uint) error {
	return r.db.Delete(&models.DailyTask{}, id).Error
}

// List returns daily tasks with pagination and filters
func (r *DailyTaskRepository) List(page, pageSize int, filters map[string]interface{}) ([]models.DailyTask, int64, error) {
	var tasks []models.DailyTask
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.DailyTask{})

	// Apply filters
	if projectID, ok := filters["project_id"].(uint); ok && projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if employeeID, ok := filters["employee_id"].(string); ok && employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if dateFrom, ok := filters["date_from"].(time.Time); ok {
		query = query.Where("task_date >= ?", dateFrom)
	}
	if dateTo, ok := filters["date_to"].(time.Time); ok {
		query = query.Where("task_date <= ?", dateTo)
	}
	if taskDate, ok := filters["task_date"].(time.Time); ok {
		query = query.Where("task_date = ?", taskDate)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("task_date DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&tasks).Error

	return tasks, total, err
}

// ListByEmployeeAndDateRange returns tasks for an employee in a date range
func (r *DailyTaskRepository) ListByEmployeeAndDateRange(employeeID string, from, to time.Time) ([]models.DailyTask, error) {
	var tasks []models.DailyTask
	err := r.db.Where("employee_id = ? AND task_date >= ? AND task_date <= ?", employeeID, from, to).
		Order("task_date DESC, created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

// GetTotalProgressForProject calculates the total progress contribution for a project
func (r *DailyTaskRepository) GetTotalProgressForProject(projectID uint) (float64, error) {
	var total float64
	err := r.db.Model(&models.DailyTask{}).
		Where("project_id = ? AND status = ?", projectID, models.TaskStatusCompleted).
		Select("COALESCE(SUM(progress_contribution), 0)").
		Scan(&total).Error
	return total, err
}

// GetTotalHoursForProject calculates the total hours spent on a project
func (r *DailyTaskRepository) GetTotalHoursForProject(projectID uint) (float64, error) {
	var total float64
	err := r.db.Model(&models.DailyTask{}).
		Where("project_id = ?", projectID).
		Select("COALESCE(SUM(hours_spent), 0)").
		Scan(&total).Error
	return total, err
}

// GetSummaryByEmployee returns task summary for an employee within a date range
func (r *DailyTaskRepository) GetSummaryByEmployee(employeeID string, from, to time.Time) (*models.DailyTaskSummary, error) {
	summary := &models.DailyTaskSummary{}

	baseQuery := r.db.Model(&models.DailyTask{}).
		Where("employee_id = ? AND task_date >= ? AND task_date <= ?", employeeID, from, to)

	// Total tasks and hours
	baseQuery.Count(&summary.TotalTasks)
	baseQuery.Select("COALESCE(SUM(hours_spent), 0)").Scan(&summary.TotalHours)

	// Count by status
	r.db.Model(&models.DailyTask{}).
		Where("employee_id = ? AND task_date >= ? AND task_date <= ? AND status = ?", employeeID, from, to, "completed").
		Count(&summary.CompletedTasks)
	r.db.Model(&models.DailyTask{}).
		Where("employee_id = ? AND task_date >= ? AND task_date <= ? AND status = ?", employeeID, from, to, "pending").
		Count(&summary.PendingTasks)
	r.db.Model(&models.DailyTask{}).
		Where("employee_id = ? AND task_date >= ? AND task_date <= ? AND status = ?", employeeID, from, to, "in_progress").
		Count(&summary.InProgressTasks)
	r.db.Model(&models.DailyTask{}).
		Where("employee_id = ? AND task_date >= ? AND task_date <= ? AND status = ?", employeeID, from, to, "blocked").
		Count(&summary.BlockedTasks)

	// By category
	r.db.Model(&models.DailyTask{}).
		Select("category, count(*) as task_count, COALESCE(SUM(hours_spent), 0) as total_hours").
		Where("employee_id = ? AND task_date >= ? AND task_date <= ?", employeeID, from, to).
		Group("category").
		Scan(&summary.ByCategory)

	// By project
	r.db.Model(&models.DailyTask{}).
		Select("project_id, project_name, count(*) as task_count, COALESCE(SUM(hours_spent), 0) as total_hours").
		Where("employee_id = ? AND task_date >= ? AND task_date <= ?", employeeID, from, to).
		Group("project_id, project_name").
		Scan(&summary.ByProject)

	return summary, nil
}

// GetSummaryByProject returns task summary for a project
func (r *DailyTaskRepository) GetSummaryByProject(projectID uint) (*models.DailyTaskSummary, error) {
	summary := &models.DailyTaskSummary{}

	baseQuery := r.db.Model(&models.DailyTask{}).Where("project_id = ?", projectID)

	baseQuery.Count(&summary.TotalTasks)
	r.db.Model(&models.DailyTask{}).Where("project_id = ?", projectID).
		Select("COALESCE(SUM(hours_spent), 0)").Scan(&summary.TotalHours)

	r.db.Model(&models.DailyTask{}).
		Where("project_id = ? AND status = ?", projectID, "completed").Count(&summary.CompletedTasks)
	r.db.Model(&models.DailyTask{}).
		Where("project_id = ? AND status = ?", projectID, "pending").Count(&summary.PendingTasks)
	r.db.Model(&models.DailyTask{}).
		Where("project_id = ? AND status = ?", projectID, "in_progress").Count(&summary.InProgressTasks)
	r.db.Model(&models.DailyTask{}).
		Where("project_id = ? AND status = ?", projectID, "blocked").Count(&summary.BlockedTasks)

	// By category
	r.db.Model(&models.DailyTask{}).
		Select("category, count(*) as task_count, COALESCE(SUM(hours_spent), 0) as total_hours").
		Where("project_id = ?", projectID).
		Group("category").
		Scan(&summary.ByCategory)

	return summary, nil
}

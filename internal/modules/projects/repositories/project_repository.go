package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/models"
	"gorm.io/gorm"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new ProjectRepository
func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{
		db: database.GetDB(),
	}
}

// Create creates a new project
func (r *ProjectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

// FindByID finds a project by ID, preloading active members
func (r *ProjectRepository) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Members", "is_active = ?", true).First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// FindByCode finds a project by project code
func (r *ProjectRepository) FindByCode(code string) (*models.Project, error) {
	var project models.Project
	err := r.db.Where("project_code = ?", code).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// ExistsByCode checks if a project with the given code exists
func (r *ProjectRepository) ExistsByCode(code string) bool {
	var count int64
	r.db.Model(&models.Project{}).Where("project_code = ?", code).Count(&count)
	return count > 0
}

// Update updates a project
func (r *ProjectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

// Delete soft deletes a project
func (r *ProjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}

// List returns projects with pagination and filters
func (r *ProjectRepository) List(page, pageSize int, filters map[string]interface{}) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Project{})

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if priority, ok := filters["priority"].(string); ok && priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if departmentID, ok := filters["department_id"].(uint); ok && departmentID > 0 {
		query = query.Where("department_id = ?", departmentID)
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if createdByID, ok := filters["created_by_id"].(uint); ok && createdByID > 0 {
		query = query.Where("created_by_id = ?", createdByID)
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name ILIKE ? OR project_code ILIKE ? OR description ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Members", "is_active = ?", true).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&projects).Error

	return projects, total, err
}

// ListByEmployeeID returns projects where the given employee is a member
func (r *ProjectRepository) ListByEmployeeID(employeeID string, page, pageSize int, filters map[string]interface{}) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	offset := (page - 1) * pageSize

	// Subquery to find project IDs the employee belongs to
	subQuery := r.db.Model(&models.ProjectMember{}).
		Select("project_id").
		Where("employee_id = ? AND is_active = ?", employeeID, true)

	query := r.db.Model(&models.Project{}).Where("id IN (?)", subQuery)

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Members", "is_active = ?", true).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&projects).Error

	return projects, total, err
}

// GetNextProjectCode generates the next sequential project code
func (r *ProjectRepository) GetNextProjectCode() string {
	var count int64
	r.db.Model(&models.Project{}).Unscoped().Count(&count)

	for i := count + 1; ; i++ {
		code := "PRJ-" + padNumber(i)
		if !r.ExistsByCode(code) {
			return code
		}
	}
}

// UpdateProgress updates only the progress field
func (r *ProjectRepository) UpdateProgress(id uint, progress float64, actualHours float64) error {
	return r.db.Model(&models.Project{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"progress":     progress,
			"actual_hours": actualHours,
		}).Error
}

// CountByStatus counts projects by status
func (r *ProjectRepository) CountByStatus() (map[string]int64, error) {
	type result struct {
		Status string
		Count  int64
	}
	var results []result
	err := r.db.Model(&models.Project{}).
		Select("status, count(*) as count").
		Group("status").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

// ─────────────── Project Member Methods ───────────────

// AddMember adds a member to a project
func (r *ProjectRepository) AddMember(member *models.ProjectMember) error {
	return r.db.Create(member).Error
}

// RemoveMember deactivates a member from a project
func (r *ProjectRepository) RemoveMember(projectID uint, employeeID string) error {
	return r.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND employee_id = ? AND is_active = ?", projectID, employeeID, true).
		Updates(map[string]interface{}{"is_active": false}).Error
}

// FindMember finds an active member of a project
func (r *ProjectRepository) FindMember(projectID uint, employeeID string) (*models.ProjectMember, error) {
	var member models.ProjectMember
	err := r.db.Where("project_id = ? AND employee_id = ? AND is_active = ?", projectID, employeeID, true).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// ListMembers lists all active members of a project
func (r *ProjectRepository) ListMembers(projectID uint) ([]models.ProjectMember, error) {
	var members []models.ProjectMember
	err := r.db.Where("project_id = ? AND is_active = ?", projectID, true).
		Order("role ASC, joined_at ASC").
		Find(&members).Error
	return members, err
}

// IsMember checks if an employee is an active member of a project
func (r *ProjectRepository) IsMember(projectID uint, employeeID string) bool {
	var count int64
	r.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND employee_id = ? AND is_active = ?", projectID, employeeID, true).
		Count(&count)
	return count > 0
}

// padNumber pads a number with leading zeros to 4 digits
func padNumber(n int64) string {
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	for len(s) < 4 {
		s = "0" + s
	}
	return s
}

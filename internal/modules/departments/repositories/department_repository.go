package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/models"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository() *DepartmentRepository {
	return &DepartmentRepository{
		db: database.GetDB(),
	}
}

// Create creates a new department
func (r *DepartmentRepository) Create(department *models.Department) error {
	return r.db.Create(department).Error
}

// FindByID finds a department by ID
func (r *DepartmentRepository) FindByID(id uint) (*models.Department, error) {
	var department models.Department
	err := r.db.Preload("ParentDepartment").Preload("SubDepartments").Preload("Shifts").First(&department, id).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

// FindByCode finds a department by code
func (r *DepartmentRepository) FindByCode(code string) (*models.Department, error) {
	var department models.Department
	err := r.db.Where("code = ?", code).First(&department).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

// Update updates a department
func (r *DepartmentRepository) Update(department *models.Department) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(department).Error
}

// Delete soft deletes a department
func (r *DepartmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Department{}, id).Error
}

// List returns all departments with pagination
func (r *DepartmentRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Department, int64, error) {
	var departments []models.Department
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Department{})

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply additional filters
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	if parentID, ok := filters["parent_department_id"]; ok {
		if parentID == nil {
			query = query.Where("parent_department_id IS NULL")
		} else {
			query = query.Where("parent_department_id = ?", parentID)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Preload("ParentDepartment").Preload("Shifts").Offset(offset).Limit(pageSize).Find(&departments).Error
	return departments, total, err
}

// FindRootDepartments finds all root departments (no parent)
func (r *DepartmentRepository) FindRootDepartments(tenantID *uint) ([]models.Department, error) {
	var departments []models.Department
	query := r.db.Where("parent_department_id IS NULL")

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.Preload("SubDepartments").Find(&departments).Error
	return departments, err
}

// FindByManager finds departments managed by a specific employee
func (r *DepartmentRepository) FindByManager(managerID uint) ([]models.Department, error) {
	var departments []models.Department
	err := r.db.Where("manager_id = ?", managerID).Find(&departments).Error
	return departments, err
}

// ExistsByCode checks if a department with the given code exists
func (r *DepartmentRepository) ExistsByCode(code string) bool {
	var count int64
	r.db.Model(&models.Department{}).Where("code = ?", code).Count(&count)
	return count > 0
}

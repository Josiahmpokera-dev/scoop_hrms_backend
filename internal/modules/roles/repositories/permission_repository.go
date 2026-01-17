package repositories

import (
	"errors"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/models"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{
		db: database.GetDB(),
	}
}

// Create creates a new permission
func (r *PermissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

// FindByID finds a permission by ID
func (r *PermissionRepository) FindByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByCode finds a permission by code
func (r *PermissionRepository) FindByCode(code string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByResourceAndAction finds a permission by resource and action
func (r *PermissionRepository) FindByResourceAndAction(resource, action string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("resource = ? AND action = ?", resource, action).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// List returns all permissions with pagination
func (r *PermissionRepository) List(page, pageSize int) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	offset := (page - 1) * pageSize

	// Count total
	if err := r.db.Model(&models.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Offset(offset).Limit(pageSize).Find(&permissions).Error
	return permissions, total, err
}

// SeedDefaultPermissions seeds default permissions for the system
func (r *PermissionRepository) SeedDefaultPermissions() error {
	defaultPermissions := []models.Permission{
		{Code: "employee:read", Name: "Read Employees", Resource: "employee", Action: "read"},
		{Code: "employee:create", Name: "Create Employees", Resource: "employee", Action: "create"},
		{Code: "employee:update", Name: "Update Employees", Resource: "employee", Action: "update"},
		{Code: "employee:delete", Name: "Delete Employees", Resource: "employee", Action: "delete"},
		{Code: "department:read", Name: "Read Departments", Resource: "department", Action: "read"},
		{Code: "department:create", Name: "Create Departments", Resource: "department", Action: "create"},
		{Code: "department:update", Name: "Update Departments", Resource: "department", Action: "update"},
		{Code: "department:delete", Name: "Delete Departments", Resource: "department", Action: "delete"},
		{Code: "payroll:run", Name: "Run Payroll", Resource: "payroll", Action: "run"},
		{Code: "attendance:approve", Name: "Approve Attendance", Resource: "attendance", Action: "approve"},
		{Code: "user:read", Name: "Read Users", Resource: "user", Action: "read"},
		{Code: "user:create", Name: "Create Users", Resource: "user", Action: "create"},
		{Code: "user:update", Name: "Update Users", Resource: "user", Action: "update"},
		{Code: "user:delete", Name: "Delete Users", Resource: "user", Action: "delete"},
	}

	for _, perm := range defaultPermissions {
		var existing models.Permission
		// Use a silent query to avoid logging "record not found" errors during seeding
		err := r.db.Where("code = ?", perm.Code).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Permission doesn't exist, create it
				if err := r.db.Create(&perm).Error; err != nil {
					return err
				}
			} else {
				// Some other error occurred
				return err
			}
		}
		// Permission already exists, skip creation
	}

	return nil
}

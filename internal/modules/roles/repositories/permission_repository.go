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

// SeedDefaultPermissions seeds default permissions for the system.
// Permissions are derived from the feature catalog to keep a single source of truth.
func (r *PermissionRepository) SeedDefaultPermissions() error {
	// Build permissions from the feature catalog
	features := models.GetAllFeatures()
	for _, feature := range features {
		for _, fp := range feature.Permissions {
			// Parse resource and action from permission code (e.g., "employee:read")
			parts := splitPermCode(fp.Code)
			resource := feature.Code
			action := fp.Code
			if len(parts) == 2 {
				resource = parts[0]
				action = parts[1]
			}

			perm := models.Permission{
				Code:     fp.Code,
				Name:     fp.Name,
				Resource: resource,
				Action:   action,
			}

			var existing models.Permission
			err := r.db.Where("code = ?", perm.Code).First(&existing).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					if err := r.db.Create(&perm).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}
		}
	}

	return nil
}

// splitPermCode splits a permission code like "employee:read" into ["employee", "read"]
func splitPermCode(code string) []string {
	for i, ch := range code {
		if ch == ':' {
			return []string{code[:i], code[i+1:]}
		}
	}
	return []string{code}
}

// FindAll returns all permissions (no pagination)
func (r *PermissionRepository) FindAll() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Order("resource, action").Find(&permissions).Error
	return permissions, err
}

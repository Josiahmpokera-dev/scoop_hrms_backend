package repositories

import (
	"errors"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		db: database.GetDB(),
	}
}

// Create creates a new role
func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// FindByID finds a role by ID
func (r *RoleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByCode finds a role by code
func (r *RoleRepository) FindByCode(code string) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").Where("code = ?", code).First(&role).Error
	if err != nil {
		// Check if it's a "record not found" error - this is expected and not a real error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &role, nil
}

// Update updates a role
func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete soft deletes a role
func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// List returns all roles with pagination
func (r *RoleRepository) List(tenantID *uint, page, pageSize int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Role{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Preload("Permissions").Offset(offset).Limit(pageSize).Find(&roles).Error
	return roles, total, err
}

// AssignPermissions assigns permissions to a role
func (r *RoleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	var permissions []models.Permission
	if err := r.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}
	
	role := models.Role{ID: roleID}
	return r.db.Model(&role).Association("Permissions").Replace(permissions)
}

// GetUserRoles gets all roles for a user
func (r *RoleRepository) GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Table("roles").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND user_roles.deleted_at IS NULL", userID).
		Preload("Permissions").
		Find(&roles).Error
	return roles, err
}

// HasPermission checks if a role has a specific permission
func (r *RoleRepository) HasPermission(roleID uint, permissionCode string) bool {
	var count int64
	r.db.Table("role_permissions").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.code = ?", roleID, permissionCode).
		Count(&count)
	return count > 0
}

// AssignUserRole assigns a role to a user
func (r *RoleRepository) AssignUserRole(userID, roleID uint, assignedBy *uint) error {
	userRole := &models.UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: assignedBy,
		AssignedAt: time.Now(),
	}
	return r.db.Create(userRole).Error
}

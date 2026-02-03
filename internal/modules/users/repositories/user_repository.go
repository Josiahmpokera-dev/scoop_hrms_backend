package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"gorm.io/gorm"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(userID uint) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", now).Error
}

// ExistsByEmail checks if a user with the email exists
func (r *UserRepository) ExistsByEmail(email string) bool {
	var count int64
	r.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// FindByRole finds users by role (user_type column)
func (r *UserRepository) FindByRole(role models.UserRole) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("user_type = ?", role).Find(&users).Error
	return users, err
}

// List returns users with pagination and optional filters (tenant, role, search).
func (r *UserRepository) List(tenantID *uint, page, pageSize int, role string, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if role != "" {
		query = query.Where("user_type = ?", role)
	}
	if search != "" {
		term := "%" + search + "%"
		query = query.Where(
			"email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ? OR username ILIKE ?",
			term, term, term, term,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListByRoles returns users whose user_type is in the given roles (e.g. admin, hr, it).
func (r *UserRepository) ListByRoles(tenantID *uint, page, pageSize int, roles []string, search string) ([]models.User, int64, error) {
	if len(roles) == 0 {
		return []models.User{}, 0, nil
	}
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{}).Where("user_type IN ?", roles)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if search != "" {
		term := "%" + search + "%"
		query = query.Where(
			"email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ? OR username ILIKE ?",
			term, term, term, term,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListBlockedUsers returns users with status blocked (and optionally suspended) for rate-limit / security admin.
func (r *UserRepository) ListBlockedUsers(tenantID *uint, page, pageSize int, includeSuspended bool) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	query := r.db.Model(&models.User{})
	if includeSuspended {
		query = query.Where("status IN ?", []string{models.UserStatusBlocked, models.UserStatusSuspended})
	} else {
		query = query.Where("status = ?", models.UserStatusBlocked)
	}
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// ListNonEmployeeUsers returns users who are not yet linked to any employee (for onboarding existing users).
// Supports tenant filter, pagination, and search by email, first_name, last_name, username.
func (r *UserRepository) ListNonEmployeeUsers(tenantID *uint, page, pageSize int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Use Model so GORM generates correct column names; LEFT JOIN employees and keep only users with no employee
	query := r.db.Model(&models.User{}).
		Joins("LEFT JOIN employees ON users.id = employees.user_id AND employees.deleted_at IS NULL").
		Where("employees.id IS NULL")

	if tenantID != nil {
		query = query.Where("users.tenant_id = ?", *tenantID)
	}
	if search != "" {
		term := "%" + search + "%"
		query = query.Where(
			"users.email ILIKE ? OR users.first_name ILIKE ? OR users.last_name ILIKE ? OR users.username ILIKE ?",
			term, term, term, term,
		)
	}

	query = query.Where("users.is_active = ?", true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("users.created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/models"
	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{
		db: database.GetDB(),
	}
}

// Create creates a new team
func (r *TeamRepository) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

// FindByID finds a team by ID
func (r *TeamRepository) FindByID(id uint) (*models.Team, error) {
	var team models.Team
	err := r.db.First(&team, id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// FindByCode finds a team by code
func (r *TeamRepository) FindByCode(code string) (*models.Team, error) {
	var team models.Team
	err := r.db.Where("code = ?", code).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// Update updates a team
func (r *TeamRepository) Update(team *models.Team) error {
	return r.db.Save(team).Error
}

// Delete soft deletes a team
func (r *TeamRepository) Delete(id uint) error {
	return r.db.Delete(&models.Team{}, id).Error
}

// List returns all teams with pagination
func (r *TeamRepository) List(tenantID, departmentID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Team, int64, error) {
	var teams []models.Team
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Team{})

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply department filter
	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	// Apply additional filters
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	if teamType, ok := filters["team_type"].(string); ok && teamType != "" {
		query = query.Where("team_type = ?", teamType)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(pageSize).Find(&teams).Error
	return teams, total, err
}

// FindByDepartment finds teams by department
func (r *TeamRepository) FindByDepartment(departmentID uint) ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Where("department_id = ?", departmentID).Find(&teams).Error
	return teams, err
}

// ExistsByCode checks if a team with the given code exists
func (r *TeamRepository) ExistsByCode(code string) bool {
	var count int64
	r.db.Model(&models.Team{}).Where("code = ?", code).Count(&count)
	return count > 0
}

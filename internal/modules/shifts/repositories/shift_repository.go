package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	"gorm.io/gorm"
)

type ShiftRepository struct {
	db *gorm.DB
}

func NewShiftRepository() *ShiftRepository {
	return &ShiftRepository{
		db: database.GetDB(),
	}
}

// Create creates a new shift
func (r *ShiftRepository) Create(shift *models.Shift) error {
	return r.db.Create(shift).Error
}

// FindByID finds a shift by ID
func (r *ShiftRepository) FindByID(id uint) (*models.Shift, error) {
	var shift models.Shift
	err := r.db.Preload("Locations").First(&shift, id).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

// FindByCode finds a shift by code
func (r *ShiftRepository) FindByCode(code string) (*models.Shift, error) {
	var shift models.Shift
	err := r.db.Where("shift_code = ?", code).First(&shift).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

// Update updates a shift
func (r *ShiftRepository) Update(shift *models.Shift) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(shift).Error
}

// Delete soft deletes a shift
func (r *ShiftRepository) Delete(id uint) error {
	return r.db.Delete(&models.Shift{}, id).Error
}

// List returns all shifts with pagination and filters
func (r *ShiftRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Shift, int64, error) {
	var shifts []models.Shift
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Shift{}).Preload("Locations")

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		if status == "active" {
			query = query.Where("is_active = ?", true)
		} else if status == "inactive" {
			query = query.Where("is_active = ?", false)
		}
		// "all" means no filter
	}

	if shiftType, ok := filters["shift_type"].(string); ok && shiftType != "" {
		query = query.Where("shift_type = ?", shiftType)
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("shift_name ILIKE ? OR shift_code ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&shifts).Error
	return shifts, total, err
}

// ExistsByCode checks if a shift with the given code exists
func (r *ShiftRepository) ExistsByCode(code string) bool {
	var count int64
	r.db.Model(&models.Shift{}).Where("shift_code = ?", code).Count(&count)
	return count > 0
}

// HasActiveRosters checks if shift is assigned to any active (published) rosters
func (r *ShiftRepository) HasActiveRosters(shiftID uint) bool {
	var count int64
	r.db.Model(&models.RosterAssignment{}).
		Where("shift_id = ? AND status = ?", shiftID, "published").
		Count(&count)
	return count > 0
}

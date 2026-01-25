package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	"gorm.io/gorm"
)

type RosterRepository struct {
	db *gorm.DB
}

func NewRosterRepository() *RosterRepository {
	return &RosterRepository{
		db: database.GetDB(),
	}
}

// Create creates a new roster assignment
func (r *RosterRepository) Create(assignment *models.RosterAssignment) error {
	return r.db.Create(assignment).Error
}

// BulkCreate creates multiple roster assignments
func (r *RosterRepository) BulkCreate(assignments []models.RosterAssignment) error {
	if len(assignments) == 0 {
		return nil
	}
	return r.db.CreateInBatches(assignments, 100).Error
}

// FindByID finds a roster assignment by ID
func (r *RosterRepository) FindByID(id uint) (*models.RosterAssignment, error) {
	var assignment models.RosterAssignment
	err := r.db.Preload("Shift").Preload("Location").First(&assignment, id).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// Update updates a roster assignment
func (r *RosterRepository) Update(assignment *models.RosterAssignment) error {
	return r.db.Save(assignment).Error
}

// Delete soft deletes a roster assignment
func (r *RosterRepository) Delete(id uint) error {
	return r.db.Delete(&models.RosterAssignment{}, id).Error
}

// List returns roster assignments with pagination and filters
func (r *RosterRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.RosterAssignment, int64, error) {
	var assignments []models.RosterAssignment
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.RosterAssignment{}).Preload("Shift").Preload("Location")

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("date >= ?", t)
		}
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("date <= ?", t)
		}
	}

	if employeeID, ok := filters["employee_id"].(string); ok && employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}

	if shiftID, ok := filters["shift_id"].(uint); ok && shiftID > 0 {
		query = query.Where("shift_id = ?", shiftID)
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if locationID, ok := filters["location_id"].(uint); ok && locationID > 0 {
		query = query.Where("location_id = ?", locationID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("date DESC, employee_id ASC").Offset(offset).Limit(pageSize).Find(&assignments).Error
	return assignments, total, err
}

// FindByDateRange finds roster assignments within a date range
func (r *RosterRepository) FindByDateRange(tenantID *uint, startDate, endDate time.Time, filters map[string]interface{}) ([]models.RosterAssignment, error) {
	var assignments []models.RosterAssignment
	query := r.db.Model(&models.RosterAssignment{}).Preload("Shift").Preload("Location").
		Where("date >= ? AND date <= ?", startDate, endDate)

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply additional filters
	if employeeIDs, ok := filters["employee_ids"].([]string); ok && len(employeeIDs) > 0 {
		query = query.Where("employee_id IN ?", employeeIDs)
	}

	if departmentID, ok := filters["department_id"].(uint); ok && departmentID > 0 {
		// This would require joining with employees table - simplified for now
		// You may need to adjust based on your employee model structure
	}

	if locationID, ok := filters["location_id"].(uint); ok && locationID > 0 {
		query = query.Where("location_id = ?", locationID)
	}

	err := query.Order("date ASC, employee_id ASC").Find(&assignments).Error
	return assignments, err
}

// PublishRosters updates roster status to published for a date range
func (r *RosterRepository) PublishRosters(tenantID *uint, startDate, endDate time.Time) (int64, error) {
	result := r.db.Model(&models.RosterAssignment{}).
		Where("tenant_id = ? AND date >= ? AND date <= ? AND status = ?", tenantID, startDate, endDate, "draft").
		Update("status", "published")
	return result.RowsAffected, result.Error
}

// FindByEmployeeAndDate finds roster assignment for a specific employee and date
func (r *RosterRepository) FindByEmployeeAndDate(employeeID string, date time.Time) (*models.RosterAssignment, error) {
	var assignment models.RosterAssignment
	err := r.db.Where("employee_id = ? AND date = ?", employeeID, date).
		Preload("Shift").Preload("Location").
		First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// Exists checks if a roster assignment exists for employee and date
func (r *RosterRepository) Exists(employeeID string, date time.Time) bool {
	var count int64
	r.db.Model(&models.RosterAssignment{}).
		Where("employee_id = ? AND date = ?", employeeID, date).
		Count(&count)
	return count > 0
}

package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"gorm.io/gorm"
)

type HolidayRepository struct {
	db *gorm.DB
}

func NewHolidayRepository() *HolidayRepository {
	return &HolidayRepository{
		db: database.GetDB(),
	}
}

// Create creates a new holiday
func (r *HolidayRepository) Create(holiday *models.Holiday) error {
	return r.db.Create(holiday).Error
}

// FindByID finds a holiday by ID
func (r *HolidayRepository) FindByID(id uint) (*models.Holiday, error) {
	var holiday models.Holiday
	err := r.db.First(&holiday, id).Error
	if err != nil {
		return nil, err
	}
	return &holiday, nil
}

// Update updates a holiday
func (r *HolidayRepository) Update(holiday *models.Holiday) error {
	return r.db.Save(holiday).Error
}

// Delete soft deletes a holiday
func (r *HolidayRepository) Delete(id uint) error {
	return r.db.Delete(&models.Holiday{}, id).Error
}

// List returns holidays with pagination and filters
func (r *HolidayRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Holiday, int64, error) {
	var holidays []models.Holiday
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Holiday{})

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if year, ok := filters["year"].(int); ok && year > 0 {
		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year, 12, 31, 23, 59, 59, 999999999, time.UTC)
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	if holidayType, ok := filters["type"].(string); ok && holidayType != "" {
		query = query.Where("type = ?", holidayType)
	}

	if isFloater, ok := filters["is_floater"].(bool); ok {
		query = query.Where("is_floater = ?", isFloater)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("date ASC").Offset(offset).Limit(pageSize).Find(&holidays).Error
	return holidays, total, err
}

// FindByDateRange finds holidays within a date range
func (r *HolidayRepository) FindByDateRange(startDate, endDate time.Time, tenantID *uint) ([]models.Holiday, error) {
	var holidays []models.Holiday
	query := r.db.Where("date >= ? AND date <= ?", startDate, endDate)
	
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	err := query.Order("date ASC").Find(&holidays).Error
	return holidays, err
}

// FindByDate finds holidays on a specific date
func (r *HolidayRepository) FindByDate(date time.Time, tenantID *uint) ([]models.Holiday, error) {
	var holidays []models.Holiday
	query := r.db.Where("DATE(date) = DATE(?)", date)
	
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	err := query.Find(&holidays).Error
	return holidays, err
}

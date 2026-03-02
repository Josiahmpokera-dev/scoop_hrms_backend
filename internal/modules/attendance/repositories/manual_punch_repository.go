package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"gorm.io/gorm"
)

// ManualPunchRepository handles manual punch database operations
type ManualPunchRepository struct {
	db *gorm.DB
}

// NewManualPunchRepository creates a new manual punch repository
func NewManualPunchRepository() *ManualPunchRepository {
	return &ManualPunchRepository{
		db: database.GetDB(),
	}
}

// Create creates a new manual punch record
func (r *ManualPunchRepository) Create(punch *models.ManualPunch) error {
	return r.db.Create(punch).Error
}

// FindByID finds a manual punch by ID
func (r *ManualPunchRepository) FindByID(id uint) (*models.ManualPunch, error) {
	var punch models.ManualPunch
	err := r.db.Where("id = ?", id).First(&punch).Error
	if err != nil {
		return nil, err
	}
	return &punch, nil
}

// FindByEmployeeAndDate finds manual punches for an employee on a specific date
func (r *ManualPunchRepository) FindByEmployeeAndDate(employeeID uint, date time.Time) ([]models.ManualPunch, error) {
	var punches []models.ManualPunch
	err := r.db.Where("employee_id = ? AND punch_date = ?", employeeID, date).Find(&punches).Error
	return punches, err
}

// FindByEmployee finds all manual punches for an employee with optional filters
func (r *ManualPunchRepository) FindByEmployee(employeeID uint, status models.ManualPunchStatus, startDate, endDate *time.Time) ([]models.ManualPunch, error) {
	query := r.db.Where("employee_id = ?", employeeID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if startDate != nil && endDate != nil {
		query = query.Where("punch_date BETWEEN ? AND ?", startDate, endDate)
	} else if startDate != nil {
		query = query.Where("punch_date >= ?", startDate)
	} else if endDate != nil {
		query = query.Where("punch_date <= ?", endDate)
	}

	var punches []models.ManualPunch
	err := query.Order("punch_date DESC, punch_time DESC").Find(&punches).Error
	return punches, err
}

// FindAll finds all manual punches with optional filters and pagination
func (r *ManualPunchRepository) FindAll(filter models.ManualPunchFilter) ([]models.ManualPunch, int64, error) {
	query := r.db.Model(&models.ManualPunch{})

	if filter.EmployeeID != nil {
		query = query.Where("employee_id = ?", *filter.EmployeeID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.PunchType != "" {
		query = query.Where("punch_type = ?", filter.PunchType)
	}

	if filter.StartDate != nil && filter.EndDate != nil {
		query = query.Where("punch_date BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	} else if filter.StartDate != nil {
		query = query.Where("punch_date >= ?", filter.StartDate)
	} else if filter.EndDate != nil {
		query = query.Where("punch_date <= ?", filter.EndDate)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit)

	var punches []models.ManualPunch
	err := query.Order("punch_date DESC, punch_time DESC").Find(&punches).Error
	return punches, total, err
}

// Update updates a manual punch record
func (r *ManualPunchRepository) Update(punch *models.ManualPunch) error {
	return r.db.Save(punch).Error
}

// UpdateStatus updates the status of a manual punch
func (r *ManualPunchRepository) UpdateStatus(id uint, status models.ManualPunchStatus, approverID *uint, rejectionReason *string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == models.ManualPunchStatusApproved || status == models.ManualPunchStatusRejected {
		updates["approver_id"] = approverID
		updates["approved_at"] = time.Now()
	}

	if status == models.ManualPunchStatusRejected && rejectionReason != nil {
		updates["rejection_reason"] = *rejectionReason
	}

	return r.db.Model(&models.ManualPunch{}).Where("id = ?", id).Updates(updates).Error
}

// Delete soft deletes a manual punch record
func (r *ManualPunchRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&models.ManualPunch{}).Error
}

// CountPendingByEmployee counts pending manual punch requests for an employee
func (r *ManualPunchRepository) CountPendingByEmployee(employeeID uint) (int64, error) {
	var count int64
	result := r.db.Model(&models.ManualPunch{}).
		Where("employee_id = ? AND status = ?", employeeID, models.ManualPunchStatusPending).
		Count(&count)
	return count, result.Error
}

// GetPendingRequests gets all pending manual punch requests for approval
func (r *ManualPunchRepository) GetPendingRequests() ([]models.ManualPunch, error) {
	var punches []models.ManualPunch
	result := r.db.Where("status = ?", models.ManualPunchStatusPending).
		Order("created_at ASC").
		Find(&punches)
	return punches, result.Error
}

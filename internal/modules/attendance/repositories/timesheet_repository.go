package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"gorm.io/gorm"
)

// TimesheetRepository handles timesheet database operations
type TimesheetRepository struct {
	db *gorm.DB
}

// NewTimesheetRepository creates a new timesheet repository
func NewTimesheetRepository() *TimesheetRepository {
	return &TimesheetRepository{
		db: database.GetDB(),
	}
}

// --- TimesheetWeek ---

// CreateWeek creates a new timesheet week
func (r *TimesheetRepository) CreateWeek(week *models.TimesheetWeek) error {
	return r.db.Create(week).Error
}

// FindWeekByID finds a timesheet week by ID with entries and approvals
func (r *TimesheetRepository) FindWeekByID(id uint) (*models.TimesheetWeek, error) {
	var week models.TimesheetWeek
	err := r.db.Preload("Entries", func(db *gorm.DB) *gorm.DB {
		return db.Order("date ASC")
	}).Preload("Approvals").Where("id = ?", id).First(&week).Error
	if err != nil {
		return nil, err
	}
	return &week, nil
}

// FindWeekByEmployeeAndDate finds a timesheet week for an employee and week start date
func (r *TimesheetRepository) FindWeekByEmployeeAndDate(employeeID uint, weekStart time.Time) (*models.TimesheetWeek, error) {
	var week models.TimesheetWeek
	err := r.db.Preload("Entries", func(db *gorm.DB) *gorm.DB {
		return db.Order("date ASC")
	}).Preload("Approvals").
		Where("employee_id = ? AND week_start = ?", employeeID, weekStart).
		First(&week).Error
	if err != nil {
		return nil, err
	}
	return &week, nil
}

// UpdateWeek updates a timesheet week
func (r *TimesheetRepository) UpdateWeek(week *models.TimesheetWeek) error {
	return r.db.Save(week).Error
}

// ListWeeksByEmployee returns all timesheet weeks for an employee with pagination
func (r *TimesheetRepository) ListWeeksByEmployee(employeeID uint, status string, page, pageSize int) ([]models.TimesheetWeek, int64, error) {
	var weeks []models.TimesheetWeek
	var total int64

	query := r.db.Model(&models.TimesheetWeek{})
	if employeeID > 0 {
		query = query.Where("employee_id = ?", employeeID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Entries", func(db *gorm.DB) *gorm.DB {
		return db.Order("date ASC")
	}).Order("week_start DESC").Offset(offset).Limit(pageSize).Find(&weeks).Error
	return weeks, total, err
}

// ListPendingApprovals lists timesheets pending approval for a manager
func (r *TimesheetRepository) ListPendingApprovals(page, pageSize int) ([]models.TimesheetWeek, int64, error) {
	var weeks []models.TimesheetWeek
	var total int64

	query := r.db.Model(&models.TimesheetWeek{}).Where("status = ?", models.TimesheetStatusSubmitted)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Entries", func(db *gorm.DB) *gorm.DB {
		return db.Order("date ASC")
	}).Preload("Approvals").Order("submitted_at ASC").Offset(offset).Limit(pageSize).Find(&weeks).Error
	return weeks, total, err
}

// ListTeamTimesheets lists timesheets for a list of employee IDs (manager view)
func (r *TimesheetRepository) ListTeamTimesheets(employeeIDs []uint, weekStart *time.Time, status string, page, pageSize int) ([]models.TimesheetWeek, int64, error) {
	var weeks []models.TimesheetWeek
	var total int64

	query := r.db.Model(&models.TimesheetWeek{}).Where("employee_id IN ?", employeeIDs)
	if weekStart != nil {
		query = query.Where("week_start = ?", *weekStart)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Entries", func(db *gorm.DB) *gorm.DB {
		return db.Order("date ASC")
	}).Order("week_start DESC").Offset(offset).Limit(pageSize).Find(&weeks).Error
	return weeks, total, err
}

// --- TimesheetEntry ---

// CreateEntry creates a new timesheet entry
func (r *TimesheetRepository) CreateEntry(entry *models.TimesheetEntry) error {
	return r.db.Create(entry).Error
}

// CreateEntries creates multiple entries
func (r *TimesheetRepository) CreateEntries(entries []models.TimesheetEntry) error {
	return r.db.Create(&entries).Error
}

// FindEntryByID finds a timesheet entry by ID
func (r *TimesheetRepository) FindEntryByID(id uint) (*models.TimesheetEntry, error) {
	var entry models.TimesheetEntry
	err := r.db.Where("id = ?", id).First(&entry).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// UpdateEntry updates a timesheet entry
func (r *TimesheetRepository) UpdateEntry(entry *models.TimesheetEntry) error {
	return r.db.Save(entry).Error
}

// DeleteEntry soft deletes a timesheet entry
func (r *TimesheetRepository) DeleteEntry(id uint) error {
	return r.db.Delete(&models.TimesheetEntry{}, id).Error
}

// ListEntriesByWeek returns all entries for a timesheet week
func (r *TimesheetRepository) ListEntriesByWeek(timesheetID uint) ([]models.TimesheetEntry, error) {
	var entries []models.TimesheetEntry
	err := r.db.Where("timesheet_id = ?", timesheetID).Order("date ASC").Find(&entries).Error
	return entries, err
}

// ListEntriesByEmployeeAndDateRange returns entries for an employee in a date range
func (r *TimesheetRepository) ListEntriesByEmployeeAndDateRange(employeeID uint, startDate, endDate time.Time) ([]models.TimesheetEntry, error) {
	var entries []models.TimesheetEntry
	err := r.db.Where("employee_id = ? AND date >= ? AND date <= ?", employeeID, startDate, endDate).
		Order("date ASC").Find(&entries).Error
	return entries, err
}

// --- TimesheetApproval ---

// CreateApproval creates a new approval record
func (r *TimesheetRepository) CreateApproval(approval *models.TimesheetApproval) error {
	return r.db.Create(approval).Error
}

// --- Aggregation / Reports ---

// GetEmployeeTimesheetStats returns total/billable hours for an employee in a date range
func (r *TimesheetRepository) GetEmployeeTimesheetStats(employeeID uint, startDate, endDate time.Time) (totalHours float64, billableHours float64, err error) {
	var result struct {
		TotalHours    float64
		BillableHours float64
	}
	err = r.db.Model(&models.TimesheetEntry{}).
		Select("COALESCE(SUM(hours), 0) as total_hours, COALESCE(SUM(CASE WHEN is_billable THEN hours ELSE 0 END), 0) as billable_hours").
		Where("employee_id = ? AND date >= ? AND date <= ?", employeeID, startDate, endDate).
		Scan(&result).Error
	return result.TotalHours, result.BillableHours, err
}

// GetProjectStats returns hours by project for an employee
func (r *TimesheetRepository) GetProjectStats(employeeID uint, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Model(&models.TimesheetEntry{}).
		Select("project_name, COALESCE(SUM(hours), 0) as total_hours, COALESCE(SUM(CASE WHEN is_billable THEN hours ELSE 0 END), 0) as billable_hours, COUNT(*) as entry_count").
		Where("employee_id = ? AND date >= ? AND date <= ?", employeeID, startDate, endDate).
		Group("project_name").
		Order("total_hours DESC").
		Find(&results).Error
	return results, err
}

// RecalcWeekTotals recalculates total and billable hours for a timesheet week
func (r *TimesheetRepository) RecalcWeekTotals(weekID uint) error {
	var result struct {
		TotalHours    float64
		BillableHours float64
	}
	err := r.db.Model(&models.TimesheetEntry{}).
		Select("COALESCE(SUM(hours), 0) as total_hours, COALESCE(SUM(CASE WHEN is_billable THEN hours ELSE 0 END), 0) as billable_hours").
		Where("timesheet_id = ?", weekID).
		Scan(&result).Error
	if err != nil {
		return err
	}
	return r.db.Model(&models.TimesheetWeek{}).Where("id = ?", weekID).
		Updates(map[string]interface{}{
			"total_hours":    result.TotalHours,
			"billable_hours": result.BillableHours,
		}).Error
}

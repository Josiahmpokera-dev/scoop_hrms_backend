package repositories

import (
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/models"
	"gorm.io/gorm"
)

// BioTimeTransactionRepository handles BioTime transaction database operations
type BioTimeTransactionRepository struct {
	db *gorm.DB
}

// NewBioTimeTransactionRepository creates a new BioTime transaction repository
func NewBioTimeTransactionRepository() *BioTimeTransactionRepository {
	return &BioTimeTransactionRepository{
		db: database.GetDB(),
	}
}

// Create creates a new BioTime transaction
func (r *BioTimeTransactionRepository) Create(transaction *models.BioTimeTransaction) error {
	return r.db.Create(transaction).Error
}

// FindByBioTimeIDAndTenant finds a transaction by BioTime transaction ID and tenant
func (r *BioTimeTransactionRepository) FindByBioTimeIDAndTenant(bioTimeID int, tenantID *uint) (*models.BioTimeTransaction, error) {
	var transaction models.BioTimeTransaction
	query := r.db.Where("biotime_transaction_id = ?", bioTimeID)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	err := query.First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// Exists checks if a transaction exists by BioTime ID, tenant, and punch time
func (r *BioTimeTransactionRepository) Exists(bioTimeID int, tenantID *uint, punchTime time.Time) (bool, error) {
	var count int64
	query := r.db.Model(&models.BioTimeTransaction{}).
		Where("biotime_transaction_id = ? AND punch_time = ?", bioTimeID, punchTime)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// BulkCreate creates multiple transactions in a batch
func (r *BioTimeTransactionRepository) BulkCreate(transactions []models.BioTimeTransaction) error {
	if len(transactions) == 0 {
		return nil
	}
	return r.db.CreateInBatches(transactions, 100).Error
}

// List lists transactions with pagination
func (r *BioTimeTransactionRepository) List(tenantID *uint, page, pageSize int) ([]models.BioTimeTransaction, int64, error) {
	var transactions []models.BioTimeTransaction
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.BioTimeTransaction{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("punch_time DESC").Offset(offset).Limit(pageSize).Find(&transactions).Error
	return transactions, total, err
}

// DailyAttendanceRecord represents a daily attendance record for an employee
type DailyAttendanceRecord struct {
	EmpCode    string     `json:"emp_code"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Department string     `json:"department"`
	Date       time.Time  `json:"date"`
	CheckIn    *time.Time `json:"check_in"`
	CheckOut   *time.Time `json:"check_out"`
	PunchCount int        `json:"punch_count"`
}

// CalendarAttendanceRecord represents one employee-day attendance aggregation.
type CalendarAttendanceRecord struct {
	Date       time.Time  `json:"date"`
	CheckIn    *time.Time `json:"check_in"`
	CheckOut   *time.Time `json:"check_out"`
	PunchCount int        `json:"punch_count"`
}

// GetDailyAttendance gets daily attendance records grouped by employee and date
// Filters by datetime range but groups by date for daily summaries
func (r *BioTimeTransactionRepository) GetDailyAttendance(tenantID *uint, startTime, endTime time.Time, empCode *string, page, pageSize int) ([]DailyAttendanceRecord, int64, error) {
	var records []DailyAttendanceRecord
	var total int64

	// Build base query - filter by datetime range, but group by date
	query := r.db.Model(&models.BioTimeTransaction{}).
		Select(`
			emp_code,
			first_name,
			last_name,
			department,
			DATE(punch_time)::date as date,
			MIN(punch_time)::timestamp as check_in,
			MAX(punch_time)::timestamp as check_out,
			COUNT(*) as punch_count
		`).
		Where("punch_time >= ? AND punch_time <= ?", startTime, endTime).
		Group("emp_code, first_name, last_name, department, DATE(punch_time)")

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	// Apply employee code filter if provided
	if empCode != nil && *empCode != "" {
		query = query.Where("emp_code = ?", *empCode)
	}

	// Count total records - using subquery for distinct count
	var countResult struct {
		Count int64
	}

	// Build count query with all filters - filter by datetime but count distinct date groups
	countSQL := `
		SELECT COUNT(*) as count FROM (
			SELECT emp_code, DATE(punch_time)
			FROM biotime_transactions
			WHERE punch_time >= ? AND punch_time <= ?
	`
	countArgs := []interface{}{startTime, endTime}

	if tenantID != nil {
		countSQL += " AND tenant_id = ?"
		countArgs = append(countArgs, *tenantID)
	} else {
		countSQL += " AND tenant_id IS NULL"
	}

	if empCode != nil && *empCode != "" {
		countSQL += " AND emp_code = ?"
		countArgs = append(countArgs, *empCode)
	}

	countSQL += " GROUP BY emp_code, DATE(punch_time)"
	countSQL += ") as distinct_records"

	if err := r.db.Raw(countSQL, countArgs...).Scan(&countResult).Error; err != nil {
		return nil, 0, err
	}
	total = countResult.Count

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Order("date DESC, emp_code ASC").Offset(offset).Limit(pageSize).Scan(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetEmployeeDailyAttendanceForMonth returns one-day aggregates for a single employee in a date range.
func (r *BioTimeTransactionRepository) GetEmployeeDailyAttendanceForMonth(tenantID *uint, empCode string, startTime, endTime time.Time) ([]CalendarAttendanceRecord, error) {
	var records []CalendarAttendanceRecord

	query := r.db.Model(&models.BioTimeTransaction{}).
		Select(`
			DATE(punch_time)::date as date,
			MIN(punch_time)::timestamp as check_in,
			MAX(punch_time)::timestamp as check_out,
			COUNT(*) as punch_count
		`).
		Where("emp_code = ?", empCode).
		Where("punch_time >= ? AND punch_time <= ?", startTime, endTime).
		Group("DATE(punch_time)").
		Order("DATE(punch_time) ASC")

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	if err := query.Scan(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

// HasTransactionsForEmpCode checks whether biometric transactions exist for the emp_code in the range.
func (r *BioTimeTransactionRepository) HasTransactionsForEmpCode(tenantID *uint, empCode string, startTime, endTime time.Time) (bool, error) {
	var count int64
	query := r.db.Model(&models.BioTimeTransaction{}).
		Where("emp_code = ?", empCode).
		Where("punch_time >= ? AND punch_time <= ?", startTime, endTime)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// LateArrivalRecord represents a late arrival record for an employee
type LateArrivalRecord struct {
	EmpCode     string     `json:"emp_code"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Department  string     `json:"department"`
	Date        time.Time  `json:"date"`
	CheckIn     *time.Time `json:"check_in"`
	MinutesLate int        `json:"minutes_late"`
}

// GetLateArrivals gets employees who checked in after 09:00
// Filters by datetime range but groups by date for daily summaries
func (r *BioTimeTransactionRepository) GetLateArrivals(tenantID *uint, startTime, endTime time.Time, empCode *string, page, pageSize int) ([]LateArrivalRecord, int64, error) {
	var records []LateArrivalRecord
	var total int64

	// Build base SQL query - filter by datetime range, group by date, and filter for late arrivals
	// Late arrival = check-in time (MIN punch_time) is after 09:00
	// minutes_late = (check_in_time - 09:00) in minutes
	// Use CTE to first get the min punch_time, then calculate minutes_late
	baseSQL := `
		WITH daily_checkins AS (
			SELECT 
				emp_code,
				first_name,
				last_name,
				department,
				DATE(punch_time)::date as date,
				MIN(punch_time)::timestamp as check_in
			FROM biotime_transactions
			WHERE punch_time >= ? AND punch_time <= ?
	`
	args := []interface{}{startTime, endTime}

	// Apply tenant filter
	if tenantID != nil {
		baseSQL += " AND tenant_id = ?"
		args = append(args, *tenantID)
	} else {
		baseSQL += " AND tenant_id IS NULL"
	}

	// Apply employee code filter if provided
	if empCode != nil && *empCode != "" {
		baseSQL += " AND emp_code = ?"
		args = append(args, *empCode)
	}

	baseSQL += `
			GROUP BY emp_code, first_name, last_name, department, DATE(punch_time)
			HAVING MIN(punch_time)::time > TIME '09:00:00'
		)
		SELECT 
			emp_code,
			first_name,
			last_name,
			department,
			date,
			check_in,
			(EXTRACT(EPOCH FROM (check_in - (date::timestamp + INTERVAL '9 hours'))) / 60)::integer as minutes_late
		FROM daily_checkins
	`

	// Count query
	countSQL := `
		SELECT COUNT(*) as count FROM (
			SELECT emp_code, DATE(punch_time) as date
			FROM biotime_transactions
			WHERE punch_time >= ? AND punch_time <= ?
	`
	countArgs := []interface{}{startTime, endTime}

	if tenantID != nil {
		countSQL += " AND tenant_id = ?"
		countArgs = append(countArgs, *tenantID)
	} else {
		countSQL += " AND tenant_id IS NULL"
	}

	if empCode != nil && *empCode != "" {
		countSQL += " AND emp_code = ?"
		countArgs = append(countArgs, *empCode)
	}

	countSQL += `
			GROUP BY emp_code, DATE(punch_time)
			HAVING MIN(punch_time)::time > TIME '09:00:00'
		) as late_arrivals
	`

	var countResult struct {
		Count int64
	}
	if err := r.db.Raw(countSQL, countArgs...).Scan(&countResult).Error; err != nil {
		return nil, 0, err
	}
	total = countResult.Count

	// Apply pagination and ordering
	offset := (page - 1) * pageSize
	baseSQL += " ORDER BY date DESC, minutes_late DESC, emp_code ASC"
	baseSQL += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

	if err := r.db.Raw(baseSQL, args...).Scan(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

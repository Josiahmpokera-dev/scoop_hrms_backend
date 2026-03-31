package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"gorm.io/gorm"
)

// PayrollRunRepository handles database operations for payroll runs
type PayrollRunRepository struct {
	db *gorm.DB
}

// NewPayrollRunRepository creates a new repository instance
func NewPayrollRunRepository() *PayrollRunRepository {
	return &PayrollRunRepository{
		db: database.DB,
	}
}

// Create creates a new payroll run
func (r *PayrollRunRepository) Create(run *models.PayrollRun) error {
	return r.db.Create(run).Error
}

// GetByID retrieves a payroll run by ID
func (r *PayrollRunRepository) GetByID(id uint, tenantID *uint) (*models.PayrollRun, error) {
	var run models.PayrollRun
	query := r.db.Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.First(&run).Error; err != nil {
		return nil, err
	}
	return &run, nil
}

// Update updates a payroll run
func (r *PayrollRunRepository) Update(run *models.PayrollRun) error {
	return r.db.Save(run).Error
}

// Delete soft-deletes a payroll run
func (r *PayrollRunRepository) Delete(id uint, tenantID *uint) error {
	query := r.db.Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	return query.Delete(&models.PayrollRun{}).Error
}

// List retrieves payroll runs with pagination and filters
func (r *PayrollRunRepository) List(tenantID *uint, status string, payYear, payMonth, page, pageSize int) ([]models.PayrollRun, int64, error) {
	var runs []models.PayrollRun
	var total int64

	query := r.db.Model(&models.PayrollRun{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if payYear > 0 {
		query = query.Where("pay_year = ?", payYear)
	}
	if payMonth > 0 {
		query = query.Where("pay_month = ?", payMonth)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&runs).Error; err != nil {
		return nil, 0, err
	}

	return runs, total, nil
}

// GetCurrentRun retrieves the current active payroll run (Draft or in review stages)
func (r *PayrollRunRepository) GetCurrentRun(tenantID *uint) (*models.PayrollRun, error) {
	var run models.PayrollRun
	query := r.db.Where("status IN ?", []string{
		string(models.PayrollRunStatusDraft),
		string(models.PayrollRunStatusPendingHR),
		string(models.PayrollRunStatusHRReviewed),
		string(models.PayrollRunStatusPendingFinance),
		string(models.PayrollRunStatusFinanceReviewed),
		string(models.PayrollRunStatusPendingManagement),
	})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Order("created_at DESC").First(&run).Error; err != nil {
		return nil, err
	}
	return &run, nil
}

// GetLastFinalized retrieves the last finalized payroll run
func (r *PayrollRunRepository) GetLastFinalized(tenantID *uint) (*models.PayrollRun, error) {
	var run models.PayrollRun
	query := r.db.Where("status IN ?", []string{string(models.PayrollRunStatusFinalized), string(models.PayrollRunStatusDisbursed), string(models.PayrollRunStatusClosed)})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Order("pay_year DESC, pay_month DESC").First(&run).Error; err != nil {
		return nil, err
	}
	return &run, nil
}

// GetRecentRuns retrieves recent payroll runs
func (r *PayrollRunRepository) GetRecentRuns(tenantID *uint, limit int) ([]models.PayrollRun, error) {
	var runs []models.PayrollRun
	query := r.db.Model(&models.PayrollRun{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Order("created_at DESC").Limit(limit).Find(&runs).Error; err != nil {
		return nil, err
	}
	return runs, nil
}

// GetByPeriod retrieves a payroll run by month and year
func (r *PayrollRunRepository) GetByPeriod(tenantID *uint, payMonth, payYear int) (*models.PayrollRun, error) {
	var run models.PayrollRun
	query := r.db.Where("pay_month = ? AND pay_year = ?", payMonth, payYear)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.First(&run).Error; err != nil {
		return nil, err
	}
	return &run, nil
}

// PayrollRunEmployeeRepository handles payroll run employee operations
type PayrollRunEmployeeRepository struct {
	db *gorm.DB
}

// NewPayrollRunEmployeeRepository creates a new repository instance
func NewPayrollRunEmployeeRepository() *PayrollRunEmployeeRepository {
	return &PayrollRunEmployeeRepository{
		db: database.DB,
	}
}

// Create creates a payroll run employee record
func (r *PayrollRunEmployeeRepository) Create(emp *models.PayrollRunEmployee) error {
	return r.db.Create(emp).Error
}

// CreateBatch creates multiple payroll run employee records
func (r *PayrollRunEmployeeRepository) CreateBatch(emps []models.PayrollRunEmployee) error {
	return r.db.CreateInBatches(emps, 100).Error
}

// GetByRunID retrieves employees for a payroll run with pagination
func (r *PayrollRunEmployeeRepository) GetByRunID(runID uint, search, departmentID string, hasChanges *bool, page, pageSize int) ([]models.PayrollRunEmployee, int64, error) {
	var employees []models.PayrollRunEmployee
	var total int64

	query := r.db.Model(&models.PayrollRunEmployee{}).Where("payroll_run_id = ?", runID)
	
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("(employee_name ILIKE ? OR employee_code ILIKE ?)", searchPattern, searchPattern)
	}
	if departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	}
	if hasChanges != nil {
		query = query.Where("has_changes = ?", *hasChanges)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("employee_name ASC").Offset(offset).Limit(pageSize).Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

// DeleteByRunID deletes all employee records for a payroll run
func (r *PayrollRunEmployeeRepository) DeleteByRunID(runID uint) error {
	return r.db.Where("payroll_run_id = ?", runID).Delete(&models.PayrollRunEmployee{}).Error
}

// UpdateEmployeeTotals updates the payroll run totals after calculation
func (r *PayrollRunRepository) UpdateEmployeeTotals(runID uint, totalEmployees int, totalGross, totalDeductions, totalNet, totalEmployerContributions float64) error {
	return r.db.Model(&models.PayrollRun{}).Where("id = ?", runID).Updates(map[string]interface{}{
		"total_employees":             totalEmployees,
		"total_gross":                 totalGross,
		"total_deductions":            totalDeductions,
		"total_net":                   totalNet,
		"total_employer_contributions": totalEmployerContributions,
	}).Error
}

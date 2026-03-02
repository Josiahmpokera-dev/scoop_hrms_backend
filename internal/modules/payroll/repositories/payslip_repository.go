package repositories

import (
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"gorm.io/gorm"
)

// PayrollEncryptionService defines the interface for payroll data encryption
type PayrollEncryptionService interface {
	EncryptPayslip(payslip *models.Payslip) error
	DecryptPayslip(payslip *models.Payslip) error
	DecryptAllPayslips(payslips []models.Payslip) error
}

// PayslipRepository handles database operations for payslips
type PayslipRepository struct {
	db                *gorm.DB
	encryptionService PayrollEncryptionService
}

// NewPayslipRepository creates a new repository instance
func NewPayslipRepository() *PayslipRepository {
	return &PayslipRepository{
		db: database.DB,
	}
}

// NewPayslipRepositoryWithEncryption creates a new repository instance with encryption service
func NewPayslipRepositoryWithEncryption(encryptionService PayrollEncryptionService) *PayslipRepository {
	return &PayslipRepository{
		db:                database.DB,
		encryptionService: encryptionService,
	}
}

// Create creates a new payslip
func (r *PayslipRepository) Create(payslip *models.Payslip) error {
	// Encrypt sensitive data before saving
	if r.encryptionService != nil {
		if err := r.encryptionService.EncryptPayslip(payslip); err != nil {
			return err
		}
	}

	return r.db.Create(payslip).Error
}

// CreateBatch creates multiple payslips
func (r *PayslipRepository) CreateBatch(payslips []models.Payslip) error {
	return r.db.CreateInBatches(payslips, 100).Error
}

// GetByID retrieves a payslip by ID
func (r *PayslipRepository) GetByID(id uint, tenantID *uint) (*models.Payslip, error) {
	var payslip models.Payslip
	query := r.db.Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.First(&payslip).Error; err != nil {
		return nil, err
	}

	// Decrypt sensitive data after loading
	if r.encryptionService != nil {
		if err := r.encryptionService.DecryptPayslip(&payslip); err != nil {
			// Log error but don't fail the operation
			log.Printf("Warning: Failed to decrypt payslip %d: %v", id, err)
		}
	}

	return &payslip, nil
}

// GetByEmployeeAndPeriod retrieves a payslip by employee and period
func (r *PayslipRepository) GetByEmployeeAndPeriod(employeeID uint, payMonth, payYear int, tenantID *uint) (*models.Payslip, error) {
	var payslip models.Payslip
	query := r.db.Where("employee_id = ? AND pay_month = ? AND pay_year = ?", employeeID, payMonth, payYear)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.First(&payslip).Error; err != nil {
		return nil, err
	}

	// Decrypt sensitive data after loading
	if r.encryptionService != nil {
		if err := r.encryptionService.DecryptPayslip(&payslip); err != nil {
			// Log error but don't fail the operation
			log.Printf("Warning: Failed to decrypt payslip for employee %d: %v", employeeID, err)
		}
	}

	return &payslip, nil
}

// Update updates a payslip
func (r *PayslipRepository) Update(payslip *models.Payslip) error {
	// Encrypt sensitive data before saving
	if r.encryptionService != nil {
		if err := r.encryptionService.EncryptPayslip(payslip); err != nil {
			return err
		}
	}

	return r.db.Save(payslip).Error
}

// List retrieves payslips with pagination and filters
func (r *PayslipRepository) List(tenantID *uint, payMonth, payYear int, employeeID, departmentID *uint, status string, page, pageSize int) ([]models.Payslip, int64, error) {
	var payslips []models.Payslip
	var total int64

	query := r.db.Model(&models.Payslip{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if payMonth > 0 {
		query = query.Where("pay_month = ?", payMonth)
	}
	if payYear > 0 {
		query = query.Where("pay_year = ?", payYear)
	}
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("pay_year DESC, pay_month DESC, employee_name ASC").Offset(offset).Limit(pageSize).Find(&payslips).Error; err != nil {
		return nil, 0, err
	}

	// Decrypt sensitive data for all payslips
	if r.encryptionService != nil {
		if err := r.encryptionService.DecryptAllPayslips(payslips); err != nil {
			// Log error but don't fail the operation
			log.Printf("Warning: Failed to decrypt some payslips: %v", err)
		}
	}

	return payslips, total, nil
}

// GetItemsByPayslipID retrieves all items for a payslip
func (r *PayslipRepository) GetItemsByPayslipID(payslipID uint) ([]models.PayslipItem, error) {
	var items []models.PayslipItem
	if err := r.db.Where("payslip_id = ?", payslipID).Order("sort_order ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CreateItem creates a payslip item
func (r *PayslipRepository) CreateItem(item *models.PayslipItem) error {
	return r.db.Create(item).Error
}

// CreateItemsBatch creates multiple payslip items
func (r *PayslipRepository) CreateItemsBatch(items []models.PayslipItem) error {
	return r.db.CreateInBatches(items, 500).Error
}

// DeleteItemsByPayslipID deletes all items for a payslip
func (r *PayslipRepository) DeleteItemsByPayslipID(payslipID uint) error {
	return r.db.Where("payslip_id = ?", payslipID).Delete(&models.PayslipItem{}).Error
}

// GetByPayrollRunID retrieves all payslips for a payroll run
func (r *PayslipRepository) GetByPayrollRunID(runID uint) ([]models.Payslip, error) {
	var payslips []models.Payslip
	if err := r.db.Where("payroll_run_id = ?", runID).Find(&payslips).Error; err != nil {
		return nil, err
	}
	return payslips, nil
}

// GetYTDTotals retrieves year-to-date totals for an employee
func (r *PayslipRepository) GetYTDTotals(employeeID uint, year int, upToMonth int) (grossTotal, taxTotal, nssfTotal, nhifTotal, netTotal float64, err error) {
	var result struct {
		GrossTotal float64
		TaxTotal   float64
		NSSFTotal  float64
		NHIFTotal  float64
		NetTotal   float64
	}
	err = r.db.Model(&models.Payslip{}).
		Select(`
			COALESCE(SUM(gross_salary), 0) as gross_total,
			COALESCE(SUM(ytd_tax), 0) as tax_total,
			COALESCE(SUM(ytd_nssf), 0) as nssf_total,
			COALESCE(SUM(ytd_nhif), 0) as nhif_total,
			COALESCE(SUM(net_pay), 0) as net_total
		`).
		Where("employee_id = ? AND pay_year = ? AND pay_month <= ?", employeeID, year, upToMonth).
		Scan(&result).Error
	return result.GrossTotal, result.TaxTotal, result.NSSFTotal, result.NHIFTotal, result.NetTotal, err
}

// GetSummary retrieves payslip summary for a period
func (r *PayslipRepository) GetSummary(tenantID *uint, employeeID *uint, payMonth, payYear int) (map[string]interface{}, error) {
	var result struct {
		NetPay          float64
		GrossSalary     float64
		TotalDeductions float64
		YTDGross        float64
		YTDTax          float64
		YTDNet          float64
	}

	query := r.db.Model(&models.Payslip{}).
		Select(`
			COALESCE(SUM(net_pay), 0) as net_pay,
			COALESCE(SUM(gross_salary), 0) as gross_salary,
			COALESCE(SUM(total_deductions), 0) as total_deductions,
			COALESCE(SUM(ytd_gross), 0) as ytd_gross,
			COALESCE(SUM(ytd_tax), 0) as ytd_tax,
			COALESCE(SUM(ytd_net), 0) as ytd_net
		`).
		Where("pay_month = ? AND pay_year = ?", payMonth, payYear)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"netPay":          result.NetPay,
		"grossSalary":     result.GrossSalary,
		"totalDeductions": result.TotalDeductions,
		"ytdGross":        result.YTDGross,
		"ytdTax":          result.YTDTax,
		"ytdNet":          result.YTDNet,
	}, nil
}

// UpdateStatus updates payslip status
func (r *PayslipRepository) UpdateStatus(payslipID uint, status models.PayslipStatus) error {
	return r.db.Model(&models.Payslip{}).Where("id = ?", payslipID).Update("status", status).Error
}

// UpdatePaymentStatus updates payslip payment status
func (r *PayslipRepository) UpdatePaymentStatus(payslipID uint, status models.PayslipStatus, paymentRef string) error {
	return r.db.Model(&models.Payslip{}).Where("id = ?", payslipID).Updates(map[string]interface{}{
		"status":            status,
		"payment_reference": paymentRef,
	}).Error
}

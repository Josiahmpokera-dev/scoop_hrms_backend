package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"gorm.io/gorm"
)

// LoanRepository handles database operations for loans
type LoanRepository struct {
	db *gorm.DB
}

// NewLoanRepository creates a new repository instance
func NewLoanRepository() *LoanRepository {
	return &LoanRepository{
		db: database.DB,
	}
}

// Create creates a new loan
func (r *LoanRepository) Create(loan *models.Loan) error {
	return r.db.Create(loan).Error
}

// GetByID retrieves a loan by ID
func (r *LoanRepository) GetByID(id uint, tenantID *uint) (*models.Loan, error) {
	var loan models.Loan
	query := r.db.Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.First(&loan).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

// Update updates a loan
func (r *LoanRepository) Update(loan *models.Loan) error {
	return r.db.Save(loan).Error
}

// List retrieves loans with pagination and filters
func (r *LoanRepository) List(tenantID *uint, employeeID *uint, status, loanType string, page, pageSize int) ([]models.Loan, int64, error) {
	var loans []models.Loan
	var total int64

	query := r.db.Model(&models.Loan{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if loanType != "" {
		query = query.Where("loan_type = ?", loanType)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&loans).Error; err != nil {
		return nil, 0, err
	}

	return loans, total, nil
}

// GetRepaymentsByLoanID retrieves all repayments for a loan
func (r *LoanRepository) GetRepaymentsByLoanID(loanID uint) ([]models.LoanRepayment, error) {
	var repayments []models.LoanRepayment
	if err := r.db.Where("loan_id = ?", loanID).Order("installment ASC").Find(&repayments).Error; err != nil {
		return nil, err
	}
	return repayments, nil
}

// CreateRepayment creates a loan repayment record
func (r *LoanRepository) CreateRepayment(repayment *models.LoanRepayment) error {
	return r.db.Create(repayment).Error
}

// CreateRepaymentsBatch creates multiple loan repayment records
func (r *LoanRepository) CreateRepaymentsBatch(repayments []models.LoanRepayment) error {
	return r.db.CreateInBatches(repayments, 100).Error
}

// UpdateRepayment updates a loan repayment
func (r *LoanRepository) UpdateRepayment(repayment *models.LoanRepayment) error {
	return r.db.Save(repayment).Error
}

// GetPendingRepayments retrieves pending repayments for EMI deduction
func (r *LoanRepository) GetPendingRepayments(employeeID uint, dueDate time.Time) ([]models.LoanRepayment, error) {
	var repayments []models.LoanRepayment
	if err := r.db.
		Joins("JOIN loans ON loan_repayments.loan_id = loans.id").
		Where("loans.employee_id = ? AND loans.status = ? AND loan_repayments.status = ? AND loan_repayments.due_date <= ?",
			employeeID, models.LoanStatusActive, models.LoanRepaymentStatusPending, dueDate).
		Find(&repayments).Error; err != nil {
		return nil, err
	}
	return repayments, nil
}

// GetSummary retrieves loan summary
func (r *LoanRepository) GetSummary(tenantID *uint, employeeID *uint) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Active loans count
	var activeCount int64
	query := r.db.Model(&models.Loan{}).Where("status = ?", models.LoanStatusActive)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	query.Count(&activeCount)
	result["activeLoansCount"] = activeCount

	// Total disbursed and outstanding
	var totals struct {
		TotalDisbursed     float64
		TotalOutstanding   float64
	}
	query = r.db.Model(&models.Loan{}).
		Select("COALESCE(SUM(amount), 0) as total_disbursed, COALESCE(SUM(outstanding_balance), 0) as total_outstanding").
		Where("status = ?", models.LoanStatusActive)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	query.Scan(&totals)
	result["totalDisbursed"] = totals.TotalDisbursed
	result["totalOutstanding"] = totals.TotalOutstanding

	// This month's EMI total
	now := time.Now()
	var thisMonthEMI float64
	r.db.Model(&models.LoanRepayment{}).
		Joins("JOIN loans ON loan_repayments.loan_id = loans.id").
		Where("loans.status = ? AND EXTRACT(MONTH FROM loan_repayments.due_date) = ? AND EXTRACT(YEAR FROM loan_repayments.due_date) = ?",
			models.LoanStatusActive, now.Month(), now.Year()).
		Select("COALESCE(SUM(loan_repayments.emi_amount), 0)").
		Scan(&thisMonthEMI)
	result["thisMonthEmiTotal"] = thisMonthEMI

	// Pending approvals
	var pendingCount int64
	query = r.db.Model(&models.Loan{}).Where("status = ?", models.LoanStatusPendingApproval)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	query.Count(&pendingCount)
	result["pendingApprovals"] = pendingCount

	// By type
	byType := make(map[string]map[string]interface{})
	var typeStats []struct {
		LoanType string
		Count    int64
		Amount   float64
	}
	query = r.db.Model(&models.Loan{}).
		Select("loan_type, COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Where("status = ?", models.LoanStatusActive).
		Group("loan_type")
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	query.Scan(&typeStats)
	for _, ts := range typeStats {
		byType[ts.LoanType] = map[string]interface{}{
			"count":  ts.Count,
			"amount": ts.Amount,
		}
	}
	result["byType"] = byType

	return result, nil
}

// GetActiveLoansForEmployee retrieves active loans for an employee
func (r *LoanRepository) GetActiveLoansForEmployee(employeeID uint) ([]models.Loan, error) {
	var loans []models.Loan
	if err := r.db.Where("employee_id = ? AND status = ?", employeeID, models.LoanStatusActive).Find(&loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

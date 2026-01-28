package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"gorm.io/gorm"
)

type LeaveBalanceRepository struct {
	db *gorm.DB
}

func NewLeaveBalanceRepository() *LeaveBalanceRepository {
	return &LeaveBalanceRepository{
		db: database.GetDB(),
	}
}

// CreateOrUpdate creates or updates a leave balance
func (r *LeaveBalanceRepository) CreateOrUpdate(balance *models.LeaveBalance) error {
	var existing models.LeaveBalance
	err := r.db.Where("employee_id = ? AND leave_type_code = ? AND year = ?", 
		balance.EmployeeID, balance.LeaveTypeCode, balance.Year).
		First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new
		return r.db.Create(balance).Error
	} else if err != nil {
		return err
	}
	
	// Update existing
	balance.ID = existing.ID
	return r.db.Save(balance).Error
}

// FindByEmployeeAndTypeAndYear finds balance for specific employee, type, and year
func (r *LeaveBalanceRepository) FindByEmployeeAndTypeAndYear(employeeID, leaveTypeCode string, year int) (*models.LeaveBalance, error) {
	var balance models.LeaveBalance
	err := r.db.Preload("LeaveType").
		Where("employee_id = ? AND leave_type_code = ? AND year = ?", employeeID, leaveTypeCode, year).
		First(&balance).Error
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

// FindByEmployeeAndYear finds all balances for an employee in a year
func (r *LeaveBalanceRepository) FindByEmployeeAndYear(employeeID string, year int, tenantID *uint) ([]models.LeaveBalance, error) {
	var balances []models.LeaveBalance
	query := r.db.Preload("LeaveType").
		Where("employee_id = ? AND year = ?", employeeID, year)
	
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	err := query.Find(&balances).Error
	return balances, err
}

// Update updates a leave balance
func (r *LeaveBalanceRepository) Update(balance *models.LeaveBalance) error {
	return r.db.Save(balance).Error
}

// IncrementUsed increments the used days
func (r *LeaveBalanceRepository) IncrementUsed(employeeID, leaveTypeCode string, year int, days float64) error {
	return r.db.Model(&models.LeaveBalance{}).
		Where("employee_id = ? AND leave_type_code = ? AND year = ?", employeeID, leaveTypeCode, year).
		UpdateColumn("used", gorm.Expr("used + ?", days)).Error
}

// IncrementPending increments the pending days
func (r *LeaveBalanceRepository) IncrementPending(employeeID, leaveTypeCode string, year int, days float64) error {
	return r.db.Model(&models.LeaveBalance{}).
		Where("employee_id = ? AND leave_type_code = ? AND year = ?", employeeID, leaveTypeCode, year).
		UpdateColumn("pending", gorm.Expr("pending + ?", days)).Error
}

// DecrementPending decrements the pending days
func (r *LeaveBalanceRepository) DecrementPending(employeeID, leaveTypeCode string, year int, days float64) error {
	return r.db.Model(&models.LeaveBalance{}).
		Where("employee_id = ? AND leave_type_code = ? AND year = ?", employeeID, leaveTypeCode, year).
		UpdateColumn("pending", gorm.Expr("pending - ?", days)).Error
}

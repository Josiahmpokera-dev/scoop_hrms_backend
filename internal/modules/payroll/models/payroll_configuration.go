package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeePayrollConfiguration extends existing employee data with payroll-specific fields
// This links to your existing employee_bank_accounts and employee_salary tables
// and adds Tanzania-specific payroll configuration

type EmployeePayrollConfiguration struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	EmployeeID        uint      `json:"employee_id" gorm:"uniqueIndex;not null"`
	GrossSalary       float64   `json:"gross_salary" gorm:"type:decimal(15,2);not null"` // Monthly gross in TZS
	SalaryGrade       string    `json:"salary_grade" gorm:"size:20;not null"`
	TINNumber         *string   `json:"tin_number" gorm:"size:50"`         // Tanzania Revenue Authority
	NSSFNumber        *string   `json:"nssf_number" gorm:"size:50"`        // NSSF membership number
	HESLBFlag         bool      `json:"heslb_flag" gorm:"default:false"`   // Student loan repayment flag
	BankName          string    `json:"bank_name" gorm:"size:50;not null"` // CRDB or DTB
	BankAccountNumber string    `json:"bank_account_number" gorm:"size:50;not null"`
	EffectiveDate     time.Time `json:"effective_date" gorm:"not null"`
	IsActive          bool      `json:"is_active" gorm:"default:true"`

	// Links to existing employee data (for reference)
	EmployeeCode string `json:"employee_code" gorm:"size:50"`
	EmployeeName string `json:"employee_name" gorm:"size:255"`
	Department   string `json:"department" gorm:"size:100"`
	Designation  string `json:"designation" gorm:"size:100"`

	// Audit fields
	CreatedByID *uint          `json:"created_by_id" gorm:"index"`
	UpdatedByID *uint          `json:"updated_by_id" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (EmployeePayrollConfiguration) TableName() string {
	return "employee_payroll_configurations"
}

// PayrollMonthlyData captures monthly variables for payroll calculation
type PayrollMonthlyData struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	PayrollRunID     uint      `json:"payroll_run_id" gorm:"index;not null"`
	EmployeeID       uint      `json:"employee_id" gorm:"index;not null"`
	Month            time.Time `json:"month" gorm:"not null"`
	WorkingDays      int       `json:"working_days" gorm:"not null"`       // Total working days in month
	DaysAttended     int       `json:"days_attended" gorm:"not null"`      // Actual days worked
	PaidLeaveDays    int       `json:"paid_leave_days" gorm:"default:0"`   // Paid leave days
	UnpaidLeaveDays  int       `json:"unpaid_leave_days" gorm:"default:0"` // Unpaid leave days
	OvertimeHours    float64   `json:"overtime_hours" gorm:"type:decimal(8,2);default:0"`
	OvertimeRate     float64   `json:"overtime_rate" gorm:"type:decimal(8,2);default:0"` // Hourly rate
	BonusAmount      float64   `json:"bonus_amount" gorm:"type:decimal(15,2);default:0"`
	ManualAdjustment float64   `json:"manual_adjustment" gorm:"type:decimal(15,2);default:0"`
	AdjustmentReason *string   `json:"adjustment_reason" gorm:"size:500"`
	LoanDeduction    float64   `json:"loan_deduction" gorm:"type:decimal(15,2);default:0"`
	HESLBDeduction   float64   `json:"heslb_deduction" gorm:"type:decimal(15,2);default:0"`

	// Calculated fields
	ProRatedSalary float64 `json:"pro_rated_salary" gorm:"type:decimal(15,2);default:0"`
	TotalEarnings  float64 `json:"total_earnings" gorm:"type:decimal(15,2);default:0"`

	// Audit fields
	CreatedByID *uint          `json:"created_by_id" gorm:"index"`
	UpdatedByID *uint          `json:"updated_by_id" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (PayrollMonthlyData) TableName() string {
	return "payroll_monthly_data"
}

// PayrollApproval tracks the approval workflow
type PayrollApproval struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	PayrollRunID      uint       `json:"payroll_run_id" gorm:"uniqueIndex;not null"`
	HRManagerID       uint       `json:"hr_manager_id" gorm:"not null"`
	HRApprovedAt      *time.Time `json:"hr_approved_at"`
	FinanceManagerID  *uint      `json:"finance_manager_id"`
	FinanceApprovedAt *time.Time `json:"finance_approved_at"`
	FinanceComments   *string    `json:"finance_comments" gorm:"size:1000"`
	IsLocked          bool       `json:"is_locked" gorm:"default:false"`
	LockedAt          *time.Time `json:"locked_at"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (PayrollApproval) TableName() string {
	return "payroll_approvals"
}

// BankSettings stores company bank account details for payroll exports
type BankSettings struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	CompanyCRDBAccount string    `json:"company_crdb_account" gorm:"size:50;not null"` // GTL company CRDB account number
	CompanyDTBAccount  string    `json:"company_dtb_account" gorm:"size:50;not null"`  // GTL company DTB account number
	DefaultPaymentDate time.Time `json:"default_payment_date" gorm:"not null"`         // Default payment date for exports
	CompanyName        string    `json:"company_name" gorm:"size:100;default:'GTL'"`   // Company name for narration

	// Audit fields
	CreatedByID *uint          `json:"created_by_id" gorm:"index"`
	UpdatedByID *uint          `json:"updated_by_id" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (BankSettings) TableName() string {
	return "bank_settings"
}

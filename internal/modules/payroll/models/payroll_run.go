package models

import (
	"time"

	"gorm.io/gorm"
)

// PayrollRunStatus represents the status of a payroll run
type PayrollRunStatus string

const (
	PayrollRunStatusDraft     PayrollRunStatus = "Draft"
	PayrollRunStatusInReview  PayrollRunStatus = "In Review"
	PayrollRunStatusApproved  PayrollRunStatus = "Approved"
	PayrollRunStatusFinalized PayrollRunStatus = "Finalized"
	PayrollRunStatusDisbursed PayrollRunStatus = "Disbursed"
	PayrollRunStatusClosed    PayrollRunStatus = "Closed"
)

// PayFrequency represents pay frequency
type PayFrequency string

const (
	PayFrequencyMonthly  PayFrequency = "Monthly"
	PayFrequencyBiWeekly PayFrequency = "Bi-Weekly"
	PayFrequencyWeekly   PayFrequency = "Weekly"
)

// PayrollRun represents a payroll processing run
type PayrollRun struct {
	ID                        uint             `json:"id" gorm:"primaryKey"`
	TenantID                  *uint            `json:"tenantId,omitempty" gorm:"index"`
	RunName                   string           `json:"runName" gorm:"not null;size:255"`
	PayPeriod                 string           `json:"payPeriod" gorm:"size:100"` // e.g., "01 Nov - 30 Nov 2024"
	PayMonth                  int              `json:"payMonth" gorm:"not null"`  // 1-12
	PayYear                   int              `json:"payYear" gorm:"not null"`
	PayFrequency              PayFrequency     `json:"payFrequency" gorm:"type:varchar(20);default:'Monthly'"`
	Status                    PayrollRunStatus `json:"status" gorm:"type:varchar(30);default:'Draft'"`
	CurrentStep               int              `json:"currentStep" gorm:"default:0"` // 0-5
	TotalEmployees            int              `json:"totalEmployees" gorm:"default:0"`
	TotalGross                float64          `json:"totalGross" gorm:"type:decimal(15,2);default:0"`
	TotalDeductions           float64          `json:"totalDeductions" gorm:"type:decimal(15,2);default:0"`
	TotalNet                  float64          `json:"totalNet" gorm:"type:decimal(15,2);default:0"`
	TotalEmployerContributions float64         `json:"totalEmployerContributions" gorm:"type:decimal(15,2);default:0"`
	// Breakdown
	TotalBasic          float64 `json:"totalBasic" gorm:"type:decimal(15,2);default:0"`
	TotalAllowances     float64 `json:"totalAllowances" gorm:"type:decimal(15,2);default:0"`
	TotalPAYE           float64 `json:"totalPaye" gorm:"type:decimal(15,2);default:0"`
	TotalNSSFEmployee   float64 `json:"totalNssfEmployee" gorm:"type:decimal(15,2);default:0"`
	TotalNHIFEmployee   float64 `json:"totalNhifEmployee" gorm:"type:decimal(15,2);default:0"`
	TotalLoanDeductions float64 `json:"totalLoanDeductions" gorm:"type:decimal(15,2);default:0"`
	TotalNSSFEmployer   float64 `json:"totalNssfEmployer" gorm:"type:decimal(15,2);default:0"`
	TotalNHIFEmployer   float64 `json:"totalNhifEmployer" gorm:"type:decimal(15,2);default:0"`
	TotalSDL            float64 `json:"totalSdl" gorm:"type:decimal(15,2);default:0"`
	TotalWCF            float64 `json:"totalWcf" gorm:"type:decimal(15,2);default:0"`
	// Dates
	CutoffDate        *time.Time `json:"cutoffDate,omitempty"`
	DisbursementDate  *time.Time `json:"disbursementDate,omitempty"`
	// Audit
	CreatedByID   *uint      `json:"createdById,omitempty" gorm:"index"`
	CreatedByName string     `json:"createdBy" gorm:"size:255"`
	ApprovedByID  *uint      `json:"approvedById,omitempty"`
	ApprovedByName *string   `json:"approvedBy,omitempty" gorm:"size:255"`
	ApprovedAt    *time.Time `json:"approvedAt,omitempty"`
	FinalizedByID *uint      `json:"finalizedById,omitempty"`
	FinalizedByName *string  `json:"finalizedBy,omitempty" gorm:"size:255"`
	FinalizedAt   *time.Time `json:"finalizedAt,omitempty"`
	IsLocked      bool       `json:"isLocked" gorm:"default:false"`
	// Timestamps
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (PayrollRun) TableName() string {
	return "payroll_runs"
}

// PayrollRunEmployee represents an employee's data in a payroll run
type PayrollRunEmployee struct {
	ID              uint    `json:"id" gorm:"primaryKey"`
	PayrollRunID    uint    `json:"payrollRunId" gorm:"index;not null"`
	EmployeeID      uint    `json:"employeeId" gorm:"index;not null"`
	EmployeeCode    string  `json:"employeeCode" gorm:"size:50"`
	EmployeeName    string  `json:"employeeName" gorm:"size:255"`
	EmployeePhoto   *string `json:"employeePhoto,omitempty" gorm:"size:500"`
	Department      string  `json:"department" gorm:"size:255"`
	DepartmentID    *uint   `json:"departmentId,omitempty"`
	Designation     string  `json:"designation" gorm:"size:255"`
	BasicSalary     float64 `json:"basicSalary" gorm:"type:decimal(15,2)"`
	GrossSalary     float64 `json:"grossSalary" gorm:"type:decimal(15,2)"`
	TotalDeductions float64 `json:"totalDeductions" gorm:"type:decimal(15,2)"`
	NetPay          float64 `json:"netPay" gorm:"type:decimal(15,2)"`
	DaysWorked      int     `json:"daysWorked" gorm:"default:0"`
	LOPDays         int     `json:"lopDays" gorm:"default:0"`
	OvertimeHours   float64 `json:"overtimeHours" gorm:"type:decimal(8,2);default:0"`
	HasChanges      bool    `json:"hasChanges" gorm:"default:false"`
	ChangeReason    *string `json:"changeReason,omitempty" gorm:"size:500"`
	// Timestamps
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName specifies the table name
func (PayrollRunEmployee) TableName() string {
	return "payroll_run_employees"
}

// PreCheckIssue represents a validation issue during pre-check
type PreCheckIssue struct {
	ID             string   `json:"id"`
	Type           string   `json:"type"`
	Category       string   `json:"category"` // attendance, leave, overtime, employee, proration
	Count          int      `json:"count"`
	Severity       string   `json:"severity"` // High, Medium, Low, Info
	EmployeeIDs    []string `json:"employeeIds"`
	Employees      []PreCheckEmployee `json:"employees,omitempty"`
	Description    string   `json:"description"`
	ActionRequired *string  `json:"actionRequired,omitempty"`
	ResolutionURL  *string  `json:"resolutionUrl,omitempty"`
}

// PreCheckEmployee represents employee info in a pre-check issue
type PreCheckEmployee struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
}

// PreCheckResult represents the result of a payroll pre-check
type PreCheckResult struct {
	PayrollRunID   uint            `json:"payrollRunId"`
	RunName        string          `json:"runName"`
	TotalEmployees int             `json:"totalEmployees"`
	ReadyCount     int             `json:"readyCount"`
	IssuesCount    int             `json:"issuesCount"`
	CanProceed     bool            `json:"canProceed"`
	Summary        string          `json:"summary"`
	Issues         []PreCheckIssue `json:"issues"`
}

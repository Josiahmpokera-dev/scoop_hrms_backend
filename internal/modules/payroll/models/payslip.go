package models

import (
	"time"

	"gorm.io/gorm"
)

// PayslipStatus represents the status of a payslip
type PayslipStatus string

const (
	PayslipStatusDraft    PayslipStatus = "Draft"
	PayslipStatusReleased PayslipStatus = "Released"
	PayslipStatusPaid     PayslipStatus = "Paid"
)

// Payslip represents an employee's payslip for a pay period
type Payslip struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	PayrollRunID    uint          `json:"payrollRunId" gorm:"index;not null"`
	EmployeeID      uint          `json:"employeeId" gorm:"index;not null"`
	EmployeeCode    string        `json:"empId" gorm:"size:50"`
	EmployeeName    string        `json:"employeeName" gorm:"size:255"`
	EmployeePhoto   *string       `json:"employeePhoto,omitempty" gorm:"size:500"`
	Department      string        `json:"department" gorm:"size:255"`
	DepartmentID    *uint         `json:"departmentId,omitempty"`
	Designation     string        `json:"designation" gorm:"size:255"`
	DateOfJoining   *time.Time    `json:"dateOfJoining,omitempty"`
	PayPeriod       string        `json:"payPeriod" gorm:"size:100"` // e.g., "October 2024"
	PayMonth        int           `json:"payMonth"`
	PayYear         int           `json:"payYear"`
	// Bank Info
	BankName        *string       `json:"bankName,omitempty" gorm:"size:255"`
	BankAccount     *string       `json:"bankAccount,omitempty" gorm:"size:100"`
	// Statutory Numbers
	TINNumber       *string       `json:"tinNumber,omitempty" gorm:"size:100"`
	NSSFNumber      *string       `json:"nssfNumber,omitempty" gorm:"size:100"`
	NHIFNumber      *string       `json:"nhifNumber,omitempty" gorm:"size:100"`
	// Amounts
	EarningsTotal              float64 `json:"earningsTotal" gorm:"type:decimal(15,2)"`
	DeductionsTotal            float64 `json:"deductionsTotal" gorm:"type:decimal(15,2)"`
	EmployerContributionsTotal float64 `json:"employerContributionsTotal" gorm:"type:decimal(15,2)"`
	GrossSalary                float64 `json:"grossSalary" gorm:"type:decimal(15,2)"`
	TotalDeductions            float64 `json:"totalDeductions" gorm:"type:decimal(15,2)"`
	NetPay                     float64 `json:"netPay" gorm:"type:decimal(15,2)"`
	// YTD (Year-to-Date)
	YTDGross  float64 `json:"ytdGross" gorm:"type:decimal(15,2)"`
	YTDTax    float64 `json:"ytdTax" gorm:"type:decimal(15,2)"`
	YTDNSSF   float64 `json:"ytdNssf" gorm:"type:decimal(15,2)"`
	YTDNHIF   float64 `json:"ytdNhif" gorm:"type:decimal(15,2)"`
	YTDNet    float64 `json:"ytdNet" gorm:"type:decimal(15,2)"`
	// Attendance
	WorkingDays   int     `json:"workingDays" gorm:"default:0"`
	DaysWorked    int     `json:"daysWorked" gorm:"default:0"`
	LOPDays       int     `json:"lopDays" gorm:"default:0"`
	OvertimeHours float64 `json:"overtimeHours" gorm:"type:decimal(8,2)"`
	// Status
	Status      PayslipStatus `json:"status" gorm:"type:varchar(30);default:'Draft'"`
	ReleasedAt  *time.Time    `json:"releasedAt,omitempty"`
	PaidAt      *time.Time    `json:"paidAt,omitempty"`
	// Payment reference
	PaymentReference *string `json:"paymentReference,omitempty" gorm:"size:100"`
	// Timestamps
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (Payslip) TableName() string {
	return "payslips"
}

// PayslipItem represents a line item in a payslip (earning, deduction, or employer contribution)
type PayslipItem struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	PayslipID     uint          `json:"payslipId" gorm:"index;not null"`
	ComponentCode string        `json:"code" gorm:"size:50"`
	ComponentName string        `json:"name" gorm:"size:255"`
	ItemType      ComponentType `json:"itemType" gorm:"type:varchar(30)"` // Earning, Deduction, Employer Contribution
	Amount        float64       `json:"amount" gorm:"type:decimal(15,2)"`
	SortOrder     int           `json:"sortOrder" gorm:"default:0"`
}

// TableName specifies the table name
func (PayslipItem) TableName() string {
	return "payslip_items"
}

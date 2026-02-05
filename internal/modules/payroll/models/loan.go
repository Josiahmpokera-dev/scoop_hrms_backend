package models

import (
	"time"

	"gorm.io/gorm"
)

// LoanType represents the type of loan
type LoanType string

const (
	LoanTypeEducation     LoanType = "Education Loan"
	LoanTypeEmergency     LoanType = "Emergency Loan"
	LoanTypePersonal      LoanType = "Personal Loan"
	LoanTypeSalaryAdvance LoanType = "Salary Advance"
)

// LoanStatus represents the status of a loan
type LoanStatus string

const (
	LoanStatusPendingApproval LoanStatus = "Pending Approval"
	LoanStatusActive          LoanStatus = "Active"
	LoanStatusClosed          LoanStatus = "Closed"
	LoanStatusDefaulted       LoanStatus = "Defaulted"
	LoanStatusRejected        LoanStatus = "Rejected"
)

// Loan represents a loan or salary advance for an employee
type Loan struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	TenantID           *uint          `json:"tenantId,omitempty" gorm:"index"`
	EmployeeID         uint           `json:"employeeId" gorm:"index;not null"`
	EmployeeCode       string         `json:"empId" gorm:"size:50"`
	EmployeeName       string         `json:"employeeName" gorm:"size:255"`
	EmployeePhoto      *string        `json:"employeePhoto,omitempty" gorm:"size:500"`
	Department         string         `json:"department" gorm:"size:255"`
	DepartmentID       *uint          `json:"departmentId,omitempty"`
	Designation        *string        `json:"designation,omitempty" gorm:"size:255"`
	LoanType           LoanType       `json:"loanType" gorm:"type:varchar(50);not null"`
	Amount             float64        `json:"amount" gorm:"type:decimal(15,2);not null"`
	InterestRate       float64        `json:"interestRate" gorm:"type:decimal(5,2);default:0"` // Annual %
	Tenure             int            `json:"tenure" gorm:"not null"`                          // Months
	EMIAmount          float64        `json:"emiAmount" gorm:"type:decimal(15,2)"`
	DisbursedDate      *time.Time     `json:"disbursedDate,omitempty"`
	StartDate          *time.Time     `json:"startDate,omitempty"` // First EMI date
	EndDate            *time.Time     `json:"endDate,omitempty"`
	Status             LoanStatus     `json:"status" gorm:"type:varchar(30);default:'Pending Approval'"`
	TotalPaid          float64        `json:"totalPaid" gorm:"type:decimal(15,2);default:0"`
	OutstandingBalance float64        `json:"outstandingBalance" gorm:"type:decimal(15,2)"`
	NextEMIDate        *time.Time     `json:"nextEmiDate,omitempty"`
	PaymentsMade       int            `json:"paymentsMade" gorm:"default:0"`
	PaymentsRemaining  int            `json:"paymentsRemaining" gorm:"default:0"`
	Purpose            *string        `json:"purpose,omitempty" gorm:"type:text"`
	// Approval
	ApprovedByID       *uint          `json:"approvedById,omitempty"`
	ApprovedByName     *string        `json:"approvedBy,omitempty" gorm:"size:255"`
	ApprovedAt         *time.Time     `json:"approvedAt,omitempty"`
	// Rejection
	RejectedByID       *uint          `json:"rejectedById,omitempty"`
	RejectedByName     *string        `json:"rejectedBy,omitempty" gorm:"size:255"`
	RejectedAt         *time.Time     `json:"rejectedAt,omitempty"`
	RejectionReason    *string        `json:"rejectionReason,omitempty" gorm:"type:text"`
	Remarks            *string        `json:"remarks,omitempty" gorm:"type:text"`
	// Timestamps
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (Loan) TableName() string {
	return "loans"
}

// LoanRepaymentStatus represents the status of a loan repayment
type LoanRepaymentStatus string

const (
	LoanRepaymentStatusPending LoanRepaymentStatus = "Pending"
	LoanRepaymentStatusPaid    LoanRepaymentStatus = "Paid"
	LoanRepaymentStatusOverdue LoanRepaymentStatus = "Overdue"
)

// LoanRepayment represents a scheduled loan repayment
type LoanRepayment struct {
	ID          uint                `json:"id" gorm:"primaryKey"`
	LoanID      uint                `json:"loanId" gorm:"index;not null"`
	Installment int                 `json:"installment" gorm:"not null"`
	DueDate     time.Time           `json:"dueDate" gorm:"not null"`
	EMIAmount   float64             `json:"emiAmount" gorm:"type:decimal(15,2)"`
	Principal   float64             `json:"principal" gorm:"type:decimal(15,2)"`
	Interest    float64             `json:"interest" gorm:"type:decimal(15,2)"`
	Status      LoanRepaymentStatus `json:"status" gorm:"type:varchar(30);default:'Pending'"`
	PaidDate    *time.Time          `json:"paidDate,omitempty"`
	PaidAmount  *float64            `json:"paidAmount,omitempty" gorm:"type:decimal(15,2)"`
	PayslipID   *uint               `json:"payslipId,omitempty"` // Link to payslip where deducted
	// Timestamps
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName specifies the table name
func (LoanRepayment) TableName() string {
	return "loan_repayments"
}

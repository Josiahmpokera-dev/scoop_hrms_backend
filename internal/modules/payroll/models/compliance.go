package models

import (
	"time"

	"gorm.io/gorm"
)

// TaxSlab represents a PAYE tax slab for a country
type TaxSlab struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Country      string         `json:"country" gorm:"size:100;default:'Tanzania'"`
	TaxYear      int            `json:"taxYear" gorm:"not null"`
	Period       string         `json:"period" gorm:"size:30;default:'Monthly'"` // Monthly, Annual
	SlabFrom     float64        `json:"slabFrom" gorm:"type:decimal(15,2)"`
	SlabTo       float64        `json:"slabTo" gorm:"type:decimal(15,2)"`
	TaxRate      float64        `json:"taxRate" gorm:"type:decimal(5,2)"` // Percentage
	FixedAmount  float64        `json:"fixedAmount" gorm:"type:decimal(15,2)"`
	Description  string         `json:"description" gorm:"size:500"`
	SortOrder    int            `json:"sortOrder" gorm:"default:0"`
	IsActive     bool           `json:"isActive" gorm:"default:true"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (TaxSlab) TableName() string {
	return "tax_slabs"
}

// StatutoryRule represents statutory contribution rules (NSSF, NHIF, SDL, WCF)
type StatutoryRule struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Country         string         `json:"country" gorm:"size:100;default:'Tanzania'"`
	Code            string         `json:"code" gorm:"size:20;not null"` // NSSF, NHIF, SDL, WCF
	Name            string         `json:"name" gorm:"size:255"`
	EmployeeRate    *float64       `json:"employeeRate,omitempty" gorm:"type:decimal(5,2)"` // Percentage
	EmployerRate    *float64       `json:"employerRate,omitempty" gorm:"type:decimal(5,2)"`
	Basis           string         `json:"basis" gorm:"size:100"` // Gross salary, Schedule, Total payroll
	Cap             *float64       `json:"cap,omitempty" gorm:"type:decimal(15,2)"`
	Description     string         `json:"description" gorm:"type:text"`
	Authority       string         `json:"authority" gorm:"size:255"` // e.g., NSSF Tanzania, TRA
	DueDay          int            `json:"dueDay" gorm:"default:7"`   // Day of month
	IsActive        bool           `json:"isActive" gorm:"default:true"`
	EffectiveFrom   *time.Time     `json:"effectiveFrom,omitempty"`
	EffectiveTo     *time.Time     `json:"effectiveTo,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (StatutoryRule) TableName() string {
	return "statutory_rules"
}

// NHIFSchedule represents the NHIF contribution schedule based on salary bands
type NHIFSchedule struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Country          string         `json:"country" gorm:"size:100;default:'Tanzania'"`
	SalaryFrom       float64        `json:"from" gorm:"type:decimal(15,2)"`
	SalaryTo         float64        `json:"to" gorm:"type:decimal(15,2)"`
	EmployeeAmount   float64        `json:"employee" gorm:"type:decimal(15,2)"`
	EmployerAmount   float64        `json:"employer" gorm:"type:decimal(15,2)"`
	IsActive         bool           `json:"isActive" gorm:"default:true"`
	SortOrder        int            `json:"sortOrder" gorm:"default:0"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (NHIFSchedule) TableName() string {
	return "nhif_schedules"
}

// CompliancePayment represents a statutory payment record
type CompliancePayment struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	PayrollRunID    *uint          `json:"payrollRunId,omitempty" gorm:"index"`
	PayMonth        int            `json:"payMonth"`
	PayYear         int            `json:"payYear"`
	StatutoryCode   string         `json:"statutoryCode" gorm:"size:20"` // PAYE, NSSF, NHIF, SDL, WCF
	EmployeeAmount  float64        `json:"employeeAmount" gorm:"type:decimal(15,2)"`
	EmployerAmount  float64        `json:"employerAmount" gorm:"type:decimal(15,2)"`
	TotalAmount     float64        `json:"totalAmount" gorm:"type:decimal(15,2)"`
	EmployeesCount  int            `json:"employeesCount"`
	DueDate         time.Time      `json:"dueDate"`
	Status          string         `json:"status" gorm:"size:30;default:'Pending'"` // Pending, Paid, Overdue
	PaidDate        *time.Time     `json:"paidDate,omitempty"`
	PaymentRef      *string        `json:"paymentRef,omitempty" gorm:"size:100"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (CompliancePayment) TableName() string {
	return "compliance_payments"
}

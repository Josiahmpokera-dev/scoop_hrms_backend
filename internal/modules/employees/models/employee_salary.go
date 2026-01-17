package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeeSalaryComponent represents detailed salary breakdown for an employee
type EmployeeSalaryComponent struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	EmployeeID        *uint          `json:"employee_id,omitempty" gorm:"index"` // Nullable for drafts
	EmployeeIDString  *string        `json:"employee_id_string,omitempty" gorm:"size:50;index"` // Employee ID string (e.g., "EMP466088") - for tracking during draft
	DraftID           *uint          `json:"draft_id,omitempty" gorm:"index"`   // References employee_onboarding_drafts(id)
	
	// Cost to Company (CTC)
	AnnualCTC         *float64       `json:"annual_ctc,omitempty" gorm:"type:decimal(15,2)"` // Annual Cost to Company
	CTCEffectiveDate  *time.Time     `json:"ctc_effective_date,omitempty"`
	Currency          string         `json:"currency" gorm:"size:10;default:'TZS'"`

	// Salary Components (Allowances)
	BasicSalary       *float64       `json:"basic_salary,omitempty" gorm:"type:decimal(15,2)"`
	HouseRentAllowance *float64     `json:"house_rent_allowance,omitempty" gorm:"type:decimal(15,2)"`
	TransportAllowance *float64      `json:"transport_allowance,omitempty" gorm:"type:decimal(15,2)"`
	SpecialAllowance   *float64      `json:"special_allowance,omitempty" gorm:"type:decimal(15,2)"`
	OtherAllowances    *float64      `json:"other_allowances,omitempty" gorm:"type:decimal(15,2)"`
	GrossSalary        *float64      `json:"gross_salary,omitempty" gorm:"type:decimal(15,2)"` // Calculated: sum of all allowances

	// Deductions
	IncomeTax          *float64      `json:"income_tax,omitempty" gorm:"type:decimal(15,2)"`
	ProvidentFund      *float64      `json:"provident_fund,omitempty" gorm:"type:decimal(15,2)"`
	ProfessionalTax    *float64      `json:"professional_tax,omitempty" gorm:"type:decimal(15,2)"`
	OtherDeductions    *float64      `json:"other_deductions,omitempty" gorm:"type:decimal(15,2)"`
	NetMonthlySalary   *float64      `json:"net_monthly_salary,omitempty" gorm:"type:decimal(15,2)"` // Calculated: Gross - Deductions

	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (EmployeeSalaryComponent) TableName() string {
	return "employee_salary_components"
}

// CalculateGrossSalary calculates gross salary from components
func (s *EmployeeSalaryComponent) CalculateGrossSalary() {
	var gross float64
	if s.BasicSalary != nil {
		gross += *s.BasicSalary
	}
	if s.HouseRentAllowance != nil {
		gross += *s.HouseRentAllowance
	}
	if s.TransportAllowance != nil {
		gross += *s.TransportAllowance
	}
	if s.SpecialAllowance != nil {
		gross += *s.SpecialAllowance
	}
	if s.OtherAllowances != nil {
		gross += *s.OtherAllowances
	}
	s.GrossSalary = &gross
}

// CalculateNetSalary calculates net monthly salary
func (s *EmployeeSalaryComponent) CalculateNetSalary() {
	if s.GrossSalary == nil {
		s.CalculateGrossSalary()
	}
	if s.GrossSalary == nil {
		return
	}

	gross := *s.GrossSalary
	var deductions float64

	if s.IncomeTax != nil {
		deductions += *s.IncomeTax
	}
	if s.ProvidentFund != nil {
		deductions += *s.ProvidentFund
	}
	if s.ProfessionalTax != nil {
		deductions += *s.ProfessionalTax
	}
	if s.OtherDeductions != nil {
		deductions += *s.OtherDeductions
	}

	net := gross - deductions
	s.NetMonthlySalary = &net
}

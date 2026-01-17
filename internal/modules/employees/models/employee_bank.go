package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeeBankAccount represents an employee's bank account details
type EmployeeBankAccount struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	EmployeeID   *uint          `json:"employee_id,omitempty" gorm:"index"` // Nullable for drafts
	EmployeeIDString *string     `json:"employee_id_string,omitempty" gorm:"size:50;index"` // Employee ID string (e.g., "EMP466088") - for tracking during draft
	DraftID      *uint          `json:"draft_id,omitempty" gorm:"index"`   // References employee_onboarding_drafts(id)
	BankName     *string        `json:"bank_name,omitempty" gorm:"size:100"`
	AccountHolderName *string    `json:"account_holder_name,omitempty" gorm:"size:100"`
	AccountNumber     *string    `json:"account_number,omitempty" gorm:"size:50"`
	AccountType       *string    `json:"account_type,omitempty" gorm:"size:50"` // savings, current, etc.
	BranchName        *string    `json:"branch_name,omitempty" gorm:"size:100"`
	SWIFTCode         *string    `json:"swift_code,omitempty" gorm:"size:20"`
	IsPrimary         bool       `json:"is_primary" gorm:"default:true"` // Primary bank account
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name
func (EmployeeBankAccount) TableName() string {
	return "employee_bank_accounts"
}

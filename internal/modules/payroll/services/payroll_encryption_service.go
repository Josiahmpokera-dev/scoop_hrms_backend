package services

import (
	"fmt"
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/encryption"
)

// PayrollEncryptionService handles encryption of sensitive payroll data
type PayrollEncryptionService struct {
	encryptionService *encryption.EncryptionService
}

// NewPayrollEncryptionService creates a new payroll encryption service
func NewPayrollEncryptionService() (*PayrollEncryptionService, error) {
	encryptionService, err := encryption.NewEncryptionService()
	if err != nil {
		return nil, fmt.Errorf("failed to create encryption service: %v", err)
	}
	
	return &PayrollEncryptionService{
		encryptionService: encryptionService,
	}, nil
}

// EncryptPayslip encrypts sensitive payslip data
func (s *PayrollEncryptionService) EncryptPayslip(payslip *models.Payslip) error {
	if payslip == nil {
		return nil
	}

	// Encrypt bank account information
	if payslip.BankAccount != nil && *payslip.BankAccount != "" {
		encryptedBankAccount, err := s.encryptionService.EncryptField(payslip.EmployeeID, *payslip.BankAccount)
		if err != nil {
			return fmt.Errorf("failed to encrypt bank account: %v", err)
		}
		payslip.BankAccount = &encryptedBankAccount
	}

	// Encrypt statutory numbers
	if payslip.TINNumber != nil && *payslip.TINNumber != "" {
		encryptedTIN, err := s.encryptionService.EncryptField(payslip.EmployeeID, *payslip.TINNumber)
		if err != nil {
			return fmt.Errorf("failed to encrypt TIN number: %v", err)
		}
		payslip.TINNumber = &encryptedTIN
	}

	if payslip.NSSFNumber != nil && *payslip.NSSFNumber != "" {
		encryptedNSSF, err := s.encryptionService.EncryptField(payslip.EmployeeID, *payslip.NSSFNumber)
		if err != nil {
			return fmt.Errorf("failed to encrypt NSSF number: %v", err)
		}
		payslip.NSSFNumber = &encryptedNSSF
	}

	if payslip.NHIFNumber != nil && *payslip.NHIFNumber != "" {
		encryptedNHIF, err := s.encryptionService.EncryptField(payslip.EmployeeID, *payslip.NHIFNumber)
		if err != nil {
			return fmt.Errorf("failed to encrypt NHIF number: %v", err)
		}
		payslip.NHIFNumber = &encryptedNHIF
	}

	return nil
}

// DecryptPayslip decrypts sensitive payslip data
func (s *PayrollEncryptionService) DecryptPayslip(payslip *models.Payslip) error {
	if payslip == nil {
		return nil
	}

	// Decrypt bank account information
	if payslip.BankAccount != nil && *payslip.BankAccount != "" {
		decryptedBankAccount, err := s.encryptionService.DecryptField(payslip.EmployeeID, *payslip.BankAccount)
		if err != nil {
			log.Printf("Warning: Failed to decrypt bank account for employee %d: %v", payslip.EmployeeID, err)
			// Don't fail entire operation, just log and continue
		} else {
			payslip.BankAccount = &decryptedBankAccount
		}
	}

	// Decrypt statutory numbers
	if payslip.TINNumber != nil && *payslip.TINNumber != "" {
		decryptedTIN, err := s.encryptionService.DecryptField(payslip.EmployeeID, *payslip.TINNumber)
		if err != nil {
			log.Printf("Warning: Failed to decrypt TIN number for employee %d: %v", payslip.EmployeeID, err)
		} else {
			payslip.TINNumber = &decryptedTIN
		}
	}

	if payslip.NSSFNumber != nil && *payslip.NSSFNumber != "" {
		decryptedNSSF, err := s.encryptionService.DecryptField(payslip.EmployeeID, *payslip.NSSFNumber)
		if err != nil {
			log.Printf("Warning: Failed to decrypt NSSF number for employee %d: %v", payslip.EmployeeID, err)
		} else {
			payslip.NSSFNumber = &decryptedNSSF
		}
	}

	if payslip.NHIFNumber != nil && *payslip.NHIFNumber != "" {
		decryptedNHIF, err := s.encryptionService.DecryptField(payslip.EmployeeID, *payslip.NHIFNumber)
		if err != nil {
			log.Printf("Warning: Failed to decrypt NHIF number for employee %d: %v", payslip.EmployeeID, err)
		} else {
			payslip.NHIFNumber = &decryptedNHIF
		}
	}

	return nil
}

// EncryptSalaryStructure encrypts sensitive salary structure data
func (s *PayrollEncryptionService) EncryptSalaryStructure(salary *models.SalaryStructure) error {
	// Salary structures typically don't contain highly sensitive individual data
	// but we can encrypt any sensitive fields here if needed
	return nil
}

// DecryptSalaryStructure decrypts sensitive salary structure data
func (s *PayrollEncryptionService) DecryptSalaryStructure(salary *models.SalaryStructure) error {
	// No sensitive fields to decrypt in salary structure currently
	return nil
}

// EncryptPayslipItems encrypts sensitive payslip item data if needed
func (s *PayrollEncryptionService) EncryptPayslipItems(items []models.PayslipItem) error {
	// Payslip items typically don't contain sensitive data that needs encryption
	return nil
}

// DecryptPayslipItems decrypts sensitive payslip item data if needed
func (s *PayrollEncryptionService) DecryptPayslipItems(items []models.PayslipItem) error {
	// No sensitive data to decrypt in payslip items currently
	return nil
}

// EncryptAllPayslips encrypts sensitive data in multiple payslips
func (s *PayrollEncryptionService) EncryptAllPayslips(payslips []models.Payslip) error {
	for i := range payslips {
		if err := s.EncryptPayslip(&payslips[i]); err != nil {
			return fmt.Errorf("failed to encrypt payslip %d: %v", payslips[i].ID, err)
		}
	}
	return nil
}

// DecryptAllPayslips decrypts sensitive data in multiple payslips
func (s *PayrollEncryptionService) DecryptAllPayslips(payslips []models.Payslip) error {
	for i := range payslips {
		if err := s.DecryptPayslip(&payslips[i]); err != nil {
			return fmt.Errorf("failed to decrypt payslip %d: %v", payslips[i].ID, err)
		}
	}
	return nil
}
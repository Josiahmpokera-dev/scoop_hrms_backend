//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
)

func main() {
	fmt.Println("🌱 Payroll Config Seeder")
	fmt.Println("========================")
	fmt.Println("This will seed all Tanzania statutory payroll configurations.")
	fmt.Println("Existing configs will be skipped (not overwritten).\n")

	// Load configuration
	if _, err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	if _, err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Warning: failed to close database: %v", err)
		}
	}()

	// Run seeder
	if err := seedPayrollConfigs(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	fmt.Println("\n✅ Payroll config seeding completed!")
}

func seedPayrollConfigs() error {
	db := database.GetDB()
	now := time.Now()

	// Define all configurations
	configurations := []models.PayrollConfig{
		// NSSF Rates
		{
			ConfigKey:   "NSSF_EMPLOYEE_RATE",
			Value:       "0.10",
			Label:       "NSSF Employee Contribution",
			Description: "10% of Gross Salary deducted from employee. Formula: =H*10% Paid to: NSSF Tanzania",
			Category:    models.ConfigCategoryNSSF,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "NSSF_EMPLOYER_RATE",
			Value:       "0.10",
			Label:       "NSSF Employer Contribution",
			Description: "10% of Gross Salary paid by employer. Formula: =H*10% Not deducted from employee net pay. Paid to: NSSF Tanzania",
			Category:    models.ConfigCategoryNSSF,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "NSSF_TOTAL_RATE",
			Value:       "0.20",
			Label:       "Total NSSF Payable",
			Description: "Employee 10% + Employer 10% = 20% of Gross. Total amount remitted to NSSF each month.",
			Category:    models.ConfigCategoryNSSF,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		// WCF Rates
		{
			ConfigKey:   "WCF_EMPLOYER_RATE",
			Value:       "0.005",
			Label:       "WCF Employer Contribution",
			Description: "0.5% of Gross Salary paid by employer only. Formula: =H*0.5% Not deducted from employee. Paid to: Workers Compensation Fund Tanzania",
			Category:    models.ConfigCategoryWCF,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		// SDL Rates
		{
			ConfigKey:   "SDL_EMPLOYER_RATE",
			Value:       "0.035",
			Label:       "SDL Employer Contribution",
			Description: "3.5% of Gross Salary paid by employer only. Formula: =H*3.5% Not deducted from employee. Paid to: Tanzania Revenue Authority (TRA)",
			Category:    models.ConfigCategorySDL,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		// HESLB Rates
		{
			ConfigKey:   "HESLB_RATE",
			Value:       "0.15",
			Label:       "HESLB Loan Repayment Rate",
			Description: "15% of Basic Pay deducted from employee. Formula: =BASIC_PAY*15% Applied only to employees with HESLB flag = true. Paid to: Higher Education Students Loans Board",
			Category:    models.ConfigCategoryHESLB,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		// PAYE Tax Bands
		{
			ConfigKey:   "PAYE_BAND_1_MIN",
			Value:       "0",
			Label:       "PAYE Band 1 Lower Limit",
			Description: "Taxable pay from TZS 0",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_1_MAX",
			Value:       "270000",
			Label:       "PAYE Band 1 Upper Limit",
			Description: "Taxable pay up to TZS 270,000 — rate 0%",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_1_RATE",
			Value:       "0.00",
			Label:       "PAYE Band 1 Rate",
			Description: "0% — nil band, no tax payable",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_2_MIN",
			Value:       "270001",
			Label:       "PAYE Band 2 Lower Limit",
			Description: "Taxable pay from TZS 270,001",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_2_MAX",
			Value:       "520000",
			Label:       "PAYE Band 2 Upper Limit",
			Description: "Taxable pay up to TZS 520,000",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_2_RATE",
			Value:       "0.08",
			Label:       "PAYE Band 2 Rate",
			Description: "8% on amount exceeding TZS 270,000. Formula: 8% × (Taxable Pay − 270,000)",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_3_MIN",
			Value:       "520001",
			Label:       "PAYE Band 3 Lower Limit",
			Description: "Taxable pay from TZS 520,001",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_3_MAX",
			Value:       "760000",
			Label:       "PAYE Band 3 Upper Limit",
			Description: "Taxable pay up to TZS 760,000",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_3_BASE_TAX",
			Value:       "20000",
			Label:       "PAYE Band 3 Base Tax",
			Description: "TZS 20,000 base tax for Band 3",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_3_RATE",
			Value:       "0.20",
			Label:       "PAYE Band 3 Rate",
			Description: "TZS 20,000 + 20% on amount exceeding TZS 520,000. Formula: 20,000 + 20% × (Taxable Pay − 520,000)",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_4_MIN",
			Value:       "760001",
			Label:       "PAYE Band 4 Lower Limit",
			Description: "Taxable pay from TZS 760,001",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_4_MAX",
			Value:       "1000000",
			Label:       "PAYE Band 4 Upper Limit",
			Description: "Taxable pay up to TZS 1,000,000",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_4_BASE_TAX",
			Value:       "68000",
			Label:       "PAYE Band 4 Base Tax",
			Description: "TZS 68,000 base tax for Band 4",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_4_RATE",
			Value:       "0.25",
			Label:       "PAYE Band 4 Rate",
			Description: "TZS 68,000 + 25% on amount exceeding TZS 760,000. Formula: 68,000 + 25% × (Taxable Pay − 760,000)",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_5_MIN",
			Value:       "1000001",
			Label:       "PAYE Band 5 Lower Limit",
			Description: "Taxable pay from TZS 1,000,001",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_5_BASE_TAX",
			Value:       "128000",
			Label:       "PAYE Band 5 Base Tax",
			Description: "TZS 128,000 base tax for Band 5",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYE_BAND_5_RATE",
			Value:       "0.30",
			Label:       "PAYE Band 5 Rate",
			Description: "TZS 128,000 + 30% on amount exceeding TZS 1,000,000. Formula: 128,000 + 30% × (Taxable Pay − 1,000,000). This is the top rate — no upper limit.",
			Category:    models.ConfigCategoryPAYE,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		// Formulas
		{
			ConfigKey:   "TAXABLE_PAY_FORMULA",
			Value:       "GROSS_SALARY - NSSF_EMPLOYEE",
			Label:       "Taxable Pay Calculation",
			Description: "Formula: =Gross − NSSF_Employee (10%). Taxable Pay is the base used for PAYE calculation.",
			Category:    models.ConfigCategoryFormula,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "NET_PAY_FORMULA",
			Value:       "TAXABLE_PAY - PAYE - HESLB - LOAN_DEDUCTION",
			Label:       "Net Pay Calculation",
			Description: "Formula: =Taxable Pay − PAYE − HESLB − Loan. This is the final amount paid to the employee.",
			Category:    models.ConfigCategoryFormula,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "CTC_FORMULA",
			Value:       "GROSS + NSSF_EMPLOYER + WCF + SDL",
			Label:       "Total Cost to Company",
			Description: "Formula: =Gross + NSSF_Employer + WCF + SDL. Full cost per employee per month.",
			Category:    models.ConfigCategoryFormula,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		// General Settings
		{
			ConfigKey:   "PAYROLL_CURRENCY",
			Value:       "TZS",
			Label:       "Payroll Currency",
			Description: "All payroll values are in Tanzanian Shillings (TZS). No decimal places. Integer values only.",
			Category:    models.ConfigCategoryGeneral,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ConfigKey:   "PAYROLL_FREQUENCY",
			Value:       "MONTHLY",
			Label:       "Payroll Frequency",
			Description: "Payroll runs once per calendar month. One payroll run per month per company.",
			Category:    models.ConfigCategoryGeneral,
			Editable:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	// Seed each config (skip if exists)
	created := 0
	skipped := 0
	failed := 0

	for _, cfg := range configurations {
		// Check if config exists
		var existing models.PayrollConfig
		err := db.Where("config_key = ?", cfg.ConfigKey).First(&existing).Error
		if err == nil {
			fmt.Printf("  ⏭️  Skipped: %s (already exists)\n", cfg.ConfigKey)
			skipped++
			continue
		}

		// Create config
		if err := db.Create(&cfg).Error; err != nil {
			fmt.Printf("  ❌ Failed:  %s - %v\n", cfg.ConfigKey, err)
			failed++
			continue
		}

		fmt.Printf("  ✅ Created: %s\n", cfg.ConfigKey)
		created++
	}

	fmt.Printf("\n📊 Summary: %d created, %d skipped, %d failed\n", created, skipped, failed)
	return nil
}

package main

//go:build ignore

import (
	"fmt"
	"log"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"gorm.io/gorm"
)

func main() {
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

	// Run migration
	if err := migratePayrollConfig(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("\n✅ Payroll config migration completed successfully!")
	fmt.Println("You can now use the /api/v1/payroll/config/paye/bands endpoint.")
}

// migratePayrollConfig safely migrates the payroll config table
func migratePayrollConfig() error {
	db := database.GetDB()

	// Step 1: Create table if not exists (safe operation)
	fmt.Println("📦 Step 1: Creating payroll_config table if not exists...")
	if err := db.AutoMigrate(&models.PayrollConfig{}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	fmt.Println("   ✓ Table is ready")

	// Step 2: Check existing configs
	var count int64
	if err := db.Model(&models.PayrollConfig{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count configs: %w", err)
	}
	fmt.Printf("   Found %d existing config records\n", count)

	// Step 3: Seed missing PAYE Band 1 MIN if needed
	fmt.Println("\n📊 Step 2: Checking PAYE Band 1 MIN config...")
	if err := seedMissingPAYEConfig(db); err != nil {
		return fmt.Errorf("failed to seed PAYE config: %w", err)
	}

	// Step 4: Verify all required configs exist
	fmt.Println("\n🔍 Step 3: Verifying all required PAYE configs...")
	if err := verifyPAYEConfigs(db); err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}

	// Step 5: Show summary
	fmt.Println("\n📋 Migration Summary:")
	var finalCount int64
	db.Model(&models.PayrollConfig{}).Count(&finalCount)
	fmt.Printf("   Total configs: %d\n", finalCount)

	// Count PAYE configs
	var payeCount int64
	db.Model(&models.PayrollConfig{}).Where("category = ?", models.ConfigCategoryPAYE).Count(&payeCount)
	fmt.Printf("   PAYE configs: %d\n", payeCount)

	return nil
}

// seedMissingPAYEConfig seeds the missing PAYE_BAND_1_MIN config if it doesn't exist
func seedMissingPAYEConfig(db *gorm.DB) error {
	var existing models.PayrollConfig
	err := db.Where("config_key = ?", "PAYE_BAND_1_MIN").First(&existing).Error

	if err == nil {
		fmt.Println("   ✓ PAYE_BAND_1_MIN already exists, skipping")
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing config: %w", err)
	}

	// Config doesn't exist, create it
	now := time.Now()
	config := models.PayrollConfig{
		ConfigKey:   "PAYE_BAND_1_MIN",
		Value:       "0",
		Label:       "PAYE Band 1 Lower Limit",
		Description: "Taxable pay from TZS 0",
		Category:    models.ConfigCategoryPAYE,
		Editable:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := db.Create(&config).Error; err != nil {
		return fmt.Errorf("failed to create PAYE_BAND_1_MIN: %w", err)
	}

	fmt.Println("   ✓ Created PAYE_BAND_1_MIN (value: 0)")
	return nil
}

// verifyPAYEConfigs verifies that all required PAYE configs exist
func verifyPAYEConfigs(db *gorm.DB) error {
	requiredConfigs := []string{
		"PAYE_BAND_1_MIN",
		"PAYE_BAND_1_MAX",
		"PAYE_BAND_1_RATE",
		"PAYE_BAND_2_MIN",
		"PAYE_BAND_2_MAX",
		"PAYE_BAND_2_RATE",
		"PAYE_BAND_3_MIN",
		"PAYE_BAND_3_MAX",
		"PAYE_BAND_3_BASE_TAX",
		"PAYE_BAND_3_RATE",
		"PAYE_BAND_4_MIN",
		"PAYE_BAND_4_MAX",
		"PAYE_BAND_4_BASE_TAX",
		"PAYE_BAND_4_RATE",
		"PAYE_BAND_5_MIN",
		"PAYE_BAND_5_BASE_TAX",
		"PAYE_BAND_5_RATE",
	}

	missing := []string{}
	for _, key := range requiredConfigs {
		var count int64
		if err := db.Model(&models.PayrollConfig{}).Where("config_key = ?", key).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check %s: %w", key, err)
		}
		if count == 0 {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		fmt.Printf("   ⚠️  Missing %d PAYE configs:\n", len(missing))
		for _, key := range missing {
			fmt.Printf("      - %s\n", key)
		}
		fmt.Println("\n   Run the full seed to add all configs:")
		fmt.Println("   go run scripts/seed_payroll_config.go")
		return nil // Don't fail, just warn
	}

	fmt.Println("   ✓ All required PAYE configs present")
	return nil
}


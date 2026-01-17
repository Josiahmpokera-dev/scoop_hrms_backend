package main

import (
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Connect to database
	dsn := cfg.Database.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Create generator
	g := gen.NewGenerator(gen.Config{
		OutPath:       "internal/modules/users/repositories/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.UseDB(db)

	// Generate ONLY users table
	g.ApplyBasic(
		g.GenerateModel("users"),
	)

	// Execute generation
	g.Execute()

	fmt.Println("✅ Successfully generated models for 'users' table")
	fmt.Println("📁 Output: internal/modules/users/repositories/query")
}

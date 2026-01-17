package main

import (
	"fmt"
	"os"

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

	// Get module name from command line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/generators/generate_models.go <module_name> <table_name>")
		fmt.Println("Example: go run scripts/generators/generate_models.go users users")
		fmt.Println("Example: go run scripts/generators/generate_models.go employees employees")
		os.Exit(1)
	}

	moduleName := os.Args[1]
	tableName := os.Args[2]

	// Connect to database
	dsn := cfg.Database.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Create generator
	g := gen.NewGenerator(gen.Config{
		OutPath:       fmt.Sprintf("internal/modules/%s/repositories/query", moduleName),
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.UseDB(db)

	// Generate model for the specified table
	g.ApplyBasic(
		g.GenerateModel(tableName),
	)

	// Execute generation
	g.Execute()

	fmt.Printf("✅ Successfully generated models for table '%s' in module '%s'\n", tableName, moduleName)
	fmt.Printf("📁 Output: internal/modules/%s/repositories/query\n", moduleName)
}

package main

import (
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	dsn := cfg.Database.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:       "internal/modules/users/repositories/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.UseDB(db)
	g.ApplyBasic(
		g.GenerateModel("users"),
	)
	g.Execute()

	fmt.Println("✅ Successfully generated models for 'users' table")
	fmt.Println("📁 Output: internal/modules/users/repositories/query")
}

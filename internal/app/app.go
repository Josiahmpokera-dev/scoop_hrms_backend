package app

import (
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
)

// Initialize initializes the application (config and database)
func Initialize() error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	log.Printf("✅ Configuration loaded successfully (Environment: %s)", cfg.Server.Env)

	// Connect to database (GORM - for migrations)
	_, err = database.Connect()
	if err != nil {
		return err
	}

	// SQLC connection is optional (moved to internal/database/optional/)
	// Uncomment and import if you want to use SQLC alongside GORM
	// See docs/guides/SQLC_SETUP.md for details

	return nil
}

// Shutdown gracefully shuts down the application
func Shutdown() error {
	log.Println("Shutting down application...")

	// Close GORM connection
	if err := database.Close(); err != nil {
		return err
	}

	// SQLC connection cleanup (if used)
	// database.CloseSQLC()

	log.Println("✅ Application shut down successfully")
	return nil
}

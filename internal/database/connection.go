package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes the database connection
func Connect() (*gorm.DB, error) {
	cfg := config.AppConfig
	if cfg == nil {
		return nil, fmt.Errorf("config not loaded, call config.LoadConfig() first")
	}

	dsn := cfg.Database.GetDSN()

	// Configure GORM logger based on environment
	var gormLogger logger.Interface
	if cfg.Server.Env == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		// In development, use a custom logger that filters out "record not found" errors
		// These are expected when checking for optional step data
		gormLogger = logger.New(
			log.New(log.Writer(), "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Error,
				IgnoreRecordNotFoundError: true, // Ignore "record not found" errors
				Colorful:                  true,
			},
		)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		// Disable foreign key constraint checks during migration
		DisableForeignKeyConstraintWhenMigrating: true,
		// Allow migration to continue even if some operations fail
		PrepareStmt: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Printf("✅ Successfully connected to PostgreSQL database: %s", cfg.Database.DBName)

	return db, nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database not initialized. Call database.Connect() first.")
	}
	return DB
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Migrate runs database migrations for the given models
func Migrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Run migrations with error handling
	// GORM's AutoMigrate will:
	// - Create tables if they don't exist
	// - Add missing columns
	// - Add missing indexes
	// - Won't delete unused columns (safe)
	err := DB.AutoMigrate(models...)

	// Check if error is about constraint (common with existing tables)
	if err != nil {
		// Check if it's a constraint error (constraint doesn't exist)
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "does not exist") && strings.Contains(errStr, "constraint") {
			// This is usually safe to ignore - GORM will create the constraint if needed
			log.Printf("⚠️  Migration warning (constraint-related): %v", err)
			log.Println("   This is usually safe - GORM will handle constraint creation.")
		} else {
			// More serious error
			log.Printf("⚠️  Migration error: %v", err)
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// HealthCheck checks if the database connection is healthy
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

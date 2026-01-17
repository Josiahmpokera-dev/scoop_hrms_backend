package database

import (
	"context"
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var SQLCDB *pgxpool.Pool

// ConnectSQLC initializes the sqlc database connection using pgxpool
func ConnectSQLC() (*pgxpool.Pool, error) {
	cfg := config.AppConfig
	if cfg == nil {
		return nil, fmt.Errorf("config not loaded, call config.LoadConfig() first")
	}

	// Build connection string for pgx
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// Parse connection string
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Configure connection pool
	poolConfig.MaxConns = 100
	poolConfig.MinConns = 10

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	SQLCDB = pool
	return pool, nil
}

// GetSQLCDB returns the sqlc database connection pool
// Returns nil if SQLC is not available (code not generated)
func GetSQLCDB() *pgxpool.Pool {
	if SQLCDB == nil {
		return nil
	}
	return SQLCDB
}

// CloseSQLC closes the sqlc database connection pool
func CloseSQLC() {
	if SQLCDB != nil {
		SQLCDB.Close()
	}
}

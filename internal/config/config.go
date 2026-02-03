package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server         ServerConfig
	Database       DatabaseConfig
	JWT            JWTConfig
	CORS           CORSConfig
	BioTime        BioTimeConfig
	RabbitMQ       RabbitMQConfig
	LoginRateLimit LoginRateLimitConfig
}

// LoginRateLimitConfig holds login rate limiting (failed attempts → block) configuration
type LoginRateLimitConfig struct {
	MaxAttempts    int // Max failed login attempts before blocking (default 5)
	LockoutMinutes int // Minutes to lock account (0 = block until admin unblocks)
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret        string
	Expiry        string // Access token expiry, e.g. "24h"
	RefreshExpiry string // Refresh token expiry, e.g. "168h" (7 days)
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// BioTimeConfig holds BioTime biometric device configuration
type BioTimeConfig struct {
	BaseURL  string
	Username string
	Password string
	Enabled  bool
}

// RabbitMQConfig holds RabbitMQ configuration
type RabbitMQConfig struct {
	URL                string
	Enabled            bool
	Exchange           string
	Queue              string
	ProcessingInterval int // Processing interval in seconds (0 = process immediately as messages arrive)
	BatchSize          int // Number of messages to process in a batch (0 = process one at a time)
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "hrms_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiry:        getEnv("JWT_EXPIRY", "24h"),
			RefreshExpiry: getEnv("JWT_REFRESH_EXPIRY", "168h"), // 7 days default
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:3000", "http://127.0.0.1:5173"}),
			AllowedMethods: getEnvSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
			AllowedHeaders: getEnvSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization", "X-Tenant-ID", "Accept", "Accept-Language"}),
		},
		BioTime: BioTimeConfig{
			BaseURL:  getEnv("BIOTIME_BASE_URL", "http://10.4.9.24:8087"),
			Username: getEnv("BIOTIME_USERNAME", "Developer"),
			Password: getEnv("BIOTIME_PASSWORD", "Developer@123"),
			Enabled:  getEnvBool("BIOTIME_ENABLED", true),
		},
		RabbitMQ: RabbitMQConfig{
			URL:                getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			Enabled:            getEnvBool("RABBITMQ_ENABLED", true),
			Exchange:           getEnv("RABBITMQ_EXCHANGE", "biotime_exchange"),
			Queue:              getEnv("RABBITMQ_QUEUE", "biotime_transactions"),
			ProcessingInterval: getEnvInt("RABBITMQ_PROCESSING_INTERVAL", 0), // 0 = immediate, or seconds (e.g., 5, 300)
			BatchSize:          getEnvInt("RABBITMQ_BATCH_SIZE", 0),          // 0 = one at a time, or batch size
		},
		LoginRateLimit: LoginRateLimitConfig{
			MaxAttempts:    getEnvInt("LOGIN_MAX_ATTEMPTS", 5),    // Block after N failed attempts
			LockoutMinutes: getEnvInt("LOGIN_LOCKOUT_MINUTES", 0), // 0 = block until admin unblocks
		},
	}

	// Validate required configuration
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	AppConfig = config
	return config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.Host,
		c.User,
		c.Password,
		c.DBName,
		c.Port,
		c.SSLMode,
	)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as integer or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool gets an environment variable as boolean or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvSlice gets an environment variable as a slice (comma-separated) or returns default values
func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim spaces
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, item := range parts {
			if trimmed := strings.TrimSpace(item); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

// validateConfig validates required configuration values
func validateConfig(config *Config) error {
	if config.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// JWT Secret validation
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required. Please set it in your .env file")
	}

	// In production, require a minimum length for security
	if config.Server.Env == "production" {
		if len(config.JWT.Secret) < 32 {
			return fmt.Errorf("JWT_SECRET must be at least 32 characters long in production")
		}
		// Check if it's still the default placeholder
		if config.JWT.Secret == "your-secret-key-change-in-production" ||
			config.JWT.Secret == "your_super_secret_jwt_key_change_in_production_min_32_chars" {
			return fmt.Errorf("JWT_SECRET must be changed from the default value in production")
		}
	}

	return nil
}

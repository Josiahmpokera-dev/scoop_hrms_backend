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
	SMTP           SMTPConfig
	LoginRateLimit LoginRateLimitConfig
	AWS            AWSConfig
	Encryption     EncryptionConfig
}

// AWSConfig holds AWS S3 storage configuration
type AWSConfig struct {
	Region          string
	S3Bucket        string
	S3BaseURL       string
	S3DocsPrefix    string
	AccessKeyID     string
	SecretAccessKey string
}

// EncryptionConfig holds encryption configuration for sensitive data
type EncryptionConfig struct {
	Key string // Encryption key for sensitive payroll data (32 bytes recommended)
}

// IsS3Enabled returns true when all required S3 fields are set.
func (a *AWSConfig) IsS3Enabled() bool {
	return a.S3Bucket != "" && a.AccessKeyID != "" && a.SecretAccessKey != "" && a.Region != ""
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
	// PublicRecruitmentCareersURL is the front-end careers site base (no trailing slash), e.g. https://careers.example.com
	// Used to build shareable apply links. If empty, the API falls back to this server's public API URLs.
	PublicRecruitmentCareersURL string
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
	// AutoSync pulls BioTime transactions into biotime_transactions on a timer (see cmd/api).
	AutoSyncEnabled         bool
	AutoSyncIntervalMinutes int
	AutoSyncLookbackHours   int
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

// SMTPConfig holds email configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	HREmail  string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:                        getEnv("PORT", "8080"),
			Env:                         getEnv("ENV", "development"),
			PublicRecruitmentCareersURL: strings.TrimSuffix(getEnv("PUBLIC_RECRUITMENT_CAREERS_URL", ""), "/"),
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
			AllowedHeaders: getEnvSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization", "Accept", "Accept-Language"}),
		},
		BioTime: BioTimeConfig{
			BaseURL:                 getEnv("BIOTIME_BASE_URL", "http://10.4.9.24:8087"),
			Username:                getEnv("BIOTIME_USERNAME", "Developer"),
			Password:                getEnv("BIOTIME_PASSWORD", "Developer@123"),
			Enabled:                   getEnvBool("BIOTIME_ENABLED", true),
			AutoSyncEnabled:           getEnvBool("BIOTIME_AUTO_SYNC_ENABLED", true),
			AutoSyncIntervalMinutes: getEnvInt("BIOTIME_AUTO_SYNC_INTERVAL_MINUTES", 1),
			AutoSyncLookbackHours:   getEnvInt("BIOTIME_AUTO_SYNC_LOOKBACK_HOURS", 48),
		},
		RabbitMQ: RabbitMQConfig{
			URL:                getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			Enabled:            getEnvBool("RABBITMQ_ENABLED", true),
			Exchange:           getEnv("RABBITMQ_EXCHANGE", "biotime_exchange"),
			Queue:              getEnv("RABBITMQ_QUEUE", "biotime_transactions"),
			ProcessingInterval: getEnvInt("RABBITMQ_PROCESSING_INTERVAL", 0), // 0 = immediate, or seconds (e.g., 5, 300)
			BatchSize:          getEnvInt("RABBITMQ_BATCH_SIZE", 0),          // 0 = one at a time, or batch size
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.mailtrap.io"),
			Port:     getEnvInt("SMTP_PORT", 587),
			Username: getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASS", ""),
			From:     getEnv("SMTP_FROM_EMAIL", "hrms@scoop.co.tz"),
			FromName: getEnv("SMTP_FROM_NAME", "HRMS"),
			HREmail:  getEnv("HR_EMAIL", "hr@scoop.co.tz"),
		},
		LoginRateLimit: LoginRateLimitConfig{
			MaxAttempts:    getEnvInt("LOGIN_MAX_ATTEMPTS", 5),    // Block after N failed attempts
			LockoutMinutes: getEnvInt("LOGIN_LOCKOUT_MINUTES", 0), // 0 = block until admin unblocks
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", ""),
			S3Bucket:        getEnv("AWS_S3_BUCKET", ""),
			S3BaseURL:       getEnv("AWS_S3_BASE_URL", ""),
			S3DocsPrefix:    getEnv("AWS_S3_DOCS_PREFIX", ""),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		},
		Encryption: EncryptionConfig{
			Key: getEnv("ENCRYPTION_KEY", "hrms-payroll-encryption-key-2024!"),
		},
	}

	if config.BioTime.AutoSyncIntervalMinutes < 1 {
		config.BioTime.AutoSyncIntervalMinutes = 1
	}
	if config.BioTime.AutoSyncIntervalMinutes > 1440 {
		config.BioTime.AutoSyncIntervalMinutes = 1440
	}
	if config.BioTime.AutoSyncLookbackHours < 1 {
		config.BioTime.AutoSyncLookbackHours = 1
	}
	if config.BioTime.AutoSyncLookbackHours > 168 {
		config.BioTime.AutoSyncLookbackHours = 168
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

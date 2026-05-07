package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/app"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	attendanceWorkers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/workers"
	biometricWorkers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/workers"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/pkg/scheduler"
	appRouter "github.com/Josiahmpokera-dev/hrms-backend/internal/router"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/types"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// Swagger documentation
	// @title HRMS Backend API
	// @version 1.0
	// @description Human Resource Management System - REST API

	// @contact.name API Support
	// @contact.url https://github.com/Josiahmpokera-dev/hrms-backend
	// @contact.email support@hrms.com

	// @license.name MIT
	// @license.url https://opensource.org/licenses/MIT

	// @host localhost:8080
	// @BasePath /api/v1
	// @schemes http
	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization
	// @security BearerAuth

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Initialize application (config and database)
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Setup graceful shutdown
	setupGracefulShutdown()

	// Set Gin mode based on environment
	if config.AppConfig.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Configure CORS middleware — accepts any origin.
	// AllowOriginFunc mirrors the request Origin back in the response header,
	// which is required when AllowCredentials is true (browsers reject "*").
	// NOTE: When AllowCredentials is true, browsers also reject "*" for
	// AllowHeaders, AllowMethods, and ExposeHeaders — they MUST be explicit.
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // Allow every origin (local dev, Docker, production)
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Content-Length", "Accept", "Accept-Encoding",
			"Accept-Language", "Authorization", "Cache-Control", "X-Requested-With",
			"X-Tenant-ID", "X-Request-Id",
		},
		ExposeHeaders: []string{
			"Content-Length", "Content-Type", "Authorization", "X-Request-Id",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure trusted proxies (important for security)
	// In production, set this to your actual proxy/load balancer IPs
	if config.AppConfig.Server.Env == "production" {
		router.SetTrustedProxies([]string{"127.0.0.1"}) // Add your actual proxy IPs
	} else {
		router.SetTrustedProxies(nil) // Trust all in development (not recommended for production)
	}

	// Get base URL from request
	getBaseURL := func(c *gin.Context) string {
		scheme := "http"
		if c.GetHeader("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}
		host := c.GetHeader("Host")
		if host == "" {
			host = "localhost:" + config.AppConfig.Server.Port
		}
		return scheme + "://" + host
	}

	// Root endpoint - API Information
	router.GET("/", func(c *gin.Context) {
		baseURL := getBaseURL(c)

		apiInfo := types.APIInfo{
			Title:       "HRMS Backend API",
			Description: "Human Resource Management System (HRMS) Backend API",
			Version:     "1.0.0",
			Status:      "running",
			Environment: config.AppConfig.Server.Env,
			Endpoints: &types.Endpoints{
				Health:        baseURL + "/health",
				HealthDB:      baseURL + "/health/db",
				API:           baseURL + "/api/v1",
				Documentation: baseURL + "/docs",
			},
			Links: &types.APILinks{
				Self:          baseURL + "/",
				Health:        baseURL + "/health",
				HealthDB:      baseURL + "/health/db",
				Documentation: baseURL + "/docs",
			},
		}

		response.Success(c, "HRMS Backend API is running successfully", apiInfo)
	})

	// API v1 root endpoint
	router.GET("/api", func(c *gin.Context) {
		baseURL := getBaseURL(c)

		apiInfo := types.APIInfo{
			Title:       "HRMS Backend API",
			Description: "Human Resource Management System (HRMS) Backend API ",
			Version:     "1.0.0",
			Status:      "running",
			Environment: config.AppConfig.Server.Env,
			Endpoints: &types.Endpoints{
				Health:        baseURL + "/health",
				HealthDB:      baseURL + "/health/db",
				API:           baseURL + "/api/v1",
				Documentation: baseURL + "/docs",
			},
			Links: &types.APILinks{
				Self:          baseURL + "/api",
				Health:        baseURL + "/health",
				HealthDB:      baseURL + "/health/db",
				Documentation: baseURL + "/docs",
			},
		}

		response.Success(c, "HRMS Backend API is running successfully", apiInfo)
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, "HRMS Backend is running", gin.H{
			"status":    "healthy",
			"service":   "hrms-backend",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Database health check
	router.GET("/health/db", func(c *gin.Context) {
		if err := database.HealthCheck(); err != nil {
			response.ServiceUnavailable(c, "Database connection failed", err.Error())
			return
		}
		response.Success(c, "Database connection is healthy", gin.H{
			"status":    "connected",
			"database":  config.AppConfig.Database.DBName,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Setup static file serving for storage
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}
	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		log.Printf("Warning: Failed to create storage directory: %v", err)
	}
	// Serve static files from storage directory
	router.Static("/storage", storagePath)

	// Serve API documentation as static files
	// NOTE: /docs/index.html must be served as a real HTML file, not markdown.
	// The ./docs folder contains index.html and API_DOCUMENTATION.md.
	router.Static("/docs", "./docs")

	// Swagger UI endpoint - with custom config to support Bearer token
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/docs/swagger/swagger.json"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	// Setup API routes
	appRouter.SetupRoutes(router)

	// Start transaction worker if RabbitMQ is enabled
	cfg := config.AppConfig
	if cfg != nil && cfg.RabbitMQ.Enabled {
		worker, err := biometricWorkers.NewTransactionWorker()
		if err != nil {
			log.Printf("Warning: Failed to start transaction worker: %v. Transactions will not be synced to database.", err)
		} else {
			if err := worker.Start(); err != nil {
				log.Printf("Warning: Failed to start transaction worker: %v. Transactions will not be synced to database.", err)
			} else {
				log.Println("✅ Transaction worker started successfully")
				defer worker.Stop()
			}
		}
	}

	// Start Attendance Report worker if RabbitMQ is enabled
	if cfg != nil && cfg.RabbitMQ.Enabled {
		reportWorker, err := attendanceWorkers.NewReportWorker()
		if err != nil {
			log.Printf("Warning: Failed to start attendance report worker: %v", err)
		} else {
			if err := reportWorker.Start(); err != nil {
				log.Printf("Warning: Failed to start attendance report worker: %v", err)
			} else {
				log.Println("✅ Attendance report worker started successfully")
			}
		}
	}

	// Start Scheduler
	s := scheduler.NewScheduler()
	s.Start()

	// Start server
	port := config.AppConfig.Server.Port
	log.Printf("🚀 Server starting on port %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupGracefulShutdown handles graceful shutdown on SIGINT or SIGTERM
func setupGracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Received shutdown signal...")
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		os.Exit(0)
	}()
}

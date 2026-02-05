package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/app"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	assetModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	auditModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/models"
	biometricModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/models"
	biometricWorkers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/workers"
	costCenterModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/cost_centers/models"
	dashboardModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/dashboard/models"
	departmentModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/models"
	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	helpdeskModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	leaveModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	locationModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/models"
	organizationUnitModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization_units/models"
	payrollModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	organizationModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/models"
	positionModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/models"
	roleModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/models"
	shiftModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	teamModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/models"
	tenantModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/tenants/models"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	appRouter "github.com/Josiahmpokera-dev/hrms-backend/internal/router"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/seed"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/types"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize application (config and database)
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Run database migrations for all models
	if err := database.Migrate(
		// Core models
		&tenantModels.Tenant{},
		&roleModels.Role{},
		&roleModels.Permission{},
		&roleModels.UserRole{},
		&roleModels.RolePermission{},
		// User model
		&userModels.User{},
		// Organizational models
		&organizationModels.Organization{},
		&organizationUnitModels.OrganizationUnit{},
		&departmentModels.Department{},
		&teamModels.Team{},
		&positionModels.JobPosition{},
		&locationModels.Location{},
		&costCenterModels.CostCenter{},
		// Asset models
		&assetModels.Asset{},
		// Employee models
		&employeeModels.Employee{},
		&employeeModels.EmployeeOnboardingDraft{},
		&employeeModels.EmployeeBasicInformation{},
		&employeeModels.EmployeeEmploymentDetails{},
		&employeeModels.EmployeeAddress{},
		&employeeModels.EmployeeSalaryComponent{},
		&employeeModels.EmployeeBankAccount{},
		&employeeModels.EmployeeStatutoryInfo{},
		&employeeModels.EmployeeDocument{},
		&employeeModels.EmployeeAsset{},
		&employeeModels.EmployeeEmergencyContact{},
		&employeeModels.EmployeePolicy{},
		&employeeModels.PostOnboardingTask{},
		// Offboarding models
		&employeeModels.OffboardingWorkflow{},
		&employeeModels.OffboardingClearance{},
		&employeeModels.OffboardingAssetReturn{},
		&employeeModels.FinalSettlement{},
		// Self-Service models
		&employeeModels.ServiceRequest{},
		&employeeModels.ProfileUpdateRequest{},
		// Asset self-service models
		&assetModels.AssetRequest{},
		&assetModels.AssetIssue{},
		// Helpdesk models
		&helpdeskModels.Ticket{},
		&helpdeskModels.Comment{},
		&helpdeskModels.Attachment{},
		&helpdeskModels.RoutingRule{},
		&helpdeskModels.KnowledgeBaseArticle{},
		&helpdeskModels.KBArticleFeedback{},
		&helpdeskModels.TicketCategory{},
		// Biometric models
		&biometricModels.BioTimeConfig{},
		&biometricModels.BioTimeTransaction{},
		// Shifts & Rosters models
		&shiftModels.Shift{},
		&shiftModels.ShiftLocation{}, // Join table for shifts and locations
		&shiftModels.RosterAssignment{},
		&shiftModels.SwapRequest{},
		&shiftModels.RosterChangeRequest{},
		// Leave Management models
		&leaveModels.LeaveType{},
		&leaveModels.LeavePolicy{},
		&leaveModels.LeaveRequest{},
		&leaveModels.LeaveApproval{},
		&leaveModels.LeaveDocument{},
		&leaveModels.LeaveBalance{},
		&leaveModels.Holiday{},
		// Security & Audit
		&auditModels.AuditLog{},
		// Payroll models
		&payrollModels.PayrollRun{},
		&payrollModels.PayrollRunEmployee{},
		&payrollModels.SalaryStructure{},
		&payrollModels.SalaryComponent{},
		&payrollModels.Payslip{},
		&payrollModels.PayslipItem{},
		&payrollModels.Loan{},
		&payrollModels.LoanRepayment{},
		&payrollModels.TaxSlab{},
		&payrollModels.StatutoryRule{},
		&payrollModels.NHIFSchedule{},
		&payrollModels.CompliancePayment{},
		// Dashboard models
		&dashboardModels.Announcement{},
		&dashboardModels.AnnouncementAttachment{},
		&dashboardModels.AnnouncementRead{},
	); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed initial data (only in development)
	seed.Run()

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

	// Configure CORS middleware — allow all origins
	corsConfig := cors.Config{
		AllowOrigins:     config.AppConfig.CORS.AllowedOrigins,
		AllowMethods:     config.AppConfig.CORS.AllowedMethods,
		AllowHeaders:     config.AppConfig.CORS.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc:  func(origin string) bool { return true },
	}
	router.Use(cors.New(corsConfig))

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

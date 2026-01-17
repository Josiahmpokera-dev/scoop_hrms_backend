package router

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	assetHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/handlers"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/auth/handlers"
	costCenterHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/cost_centers/handlers"
	departmentHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/handlers"
	employeeHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/handlers"
	locationHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/handlers"
	orgChartHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization/handlers"
	organizationUnitHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization_units/handlers"
	organizationHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/handlers"
	positionHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/handlers"
	teamHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(r *gin.Engine) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	onboardingHandler := handlers.NewOnboardingHandler()
	employeeHandler := employeeHandlers.NewEmployeeHandler()
	organizationHandler := organizationHandlers.NewOrganizationHandler()
	organizationUnitHandler := organizationUnitHandlers.NewOrganizationUnitHandler()
	departmentHandler := departmentHandlers.NewDepartmentHandler()
	orgChartHandler := orgChartHandlers.NewOrgChartHandler()
	teamHandler := teamHandlers.NewTeamHandler()
	positionHandler := positionHandlers.NewJobPositionHandler()
	locationHandler := locationHandlers.NewLocationHandler()
	costCenterHandler := costCenterHandlers.NewCostCenterHandler()

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			// Onboarding (complete signup + organization creation)
			auth.POST("/onboard", onboardingHandler.CompleteOnboarding)

			// Regular auth endpoints
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)

			// Protected routes
			auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
			auth.GET("/setup-wizard/status", middleware.AuthMiddleware(), onboardingHandler.GetSetupWizardStatus)
		}

		// Admin routes (require admin role)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// Admin routes will be added here
			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Admin dashboard",
				})
			})
		}

		// HR routes (require HR or Admin role)
		hr := v1.Group("/hr")
		hr.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			// HR routes will be added here
			hr.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "HR dashboard",
				})
			})
		}

		// Organization routes (require authentication - HR or Admin)
		organizations := v1.Group("/organizations")
		organizations.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			organizations.GET("/me", organizationHandler.GetMyOrganization)      // Get current user's organization (created during onboarding)
			organizations.POST("", organizationHandler.CreateOrganization)       // Create organization
			organizations.GET("", organizationHandler.ListOrganizations)         // List organizations
			organizations.GET("/:id", organizationHandler.GetOrganization)       // Get organization by ID
			organizations.PUT("/:id", organizationHandler.UpdateOrganization)    // Update organization
			organizations.DELETE("/:id", organizationHandler.DeleteOrganization) // Delete organization
			// POST-only action-based endpoint
			organizations.POST("/action", organizationHandler.HandleAction) // Action-based API
		}

		// Organization Chart route (require authentication - HR or Admin)
		organization := v1.Group("/organization")
		organization.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			organization.GET("/chart", orgChartHandler.GetOrganizationChart)            // Get organization chart
			organization.GET("/chart/debug", orgChartHandler.GetOrganizationChartDebug) // Debug/diagnostic endpoint
		}

		// Organization Units routes (require authentication - HR or Admin)
		organizationUnits := v1.Group("/organization-units")
		organizationUnits.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			organizationUnits.POST("", organizationUnitHandler.CreateOrganizationUnit)       // Create organization unit
			organizationUnits.GET("", organizationUnitHandler.ListOrganizationUnits)         // List organization units
			organizationUnits.GET("/root", organizationUnitHandler.GetRootOrganizationUnits) // Get root organization units
			organizationUnits.GET("/:id", organizationUnitHandler.GetOrganizationUnit)       // Get organization unit by ID
			organizationUnits.PUT("/:id", organizationUnitHandler.UpdateOrganizationUnit)    // Update organization unit
			organizationUnits.DELETE("/:id", organizationUnitHandler.DeleteOrganizationUnit) // Delete organization unit
			// POST-only action-based endpoint
			organizationUnits.POST("/action", organizationUnitHandler.HandleAction) // Action-based API
		}

		// Employee routes (require authentication - HR or Admin)
		employees := v1.Group("/employees")
		employees.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			employees.POST("", employeeHandler.OnboardEmployee)                                 // Onboard new employee (legacy - single step)
			employees.GET("", employeeHandler.ListEmployees)                                    // List all employees
			employees.GET("/:id", employeeHandler.GetEmployee)                                  // Get employee by ID
			employees.PUT("/:id", employeeHandler.UpdateEmployee)                               // Update employee
			employees.DELETE("/:id", employeeHandler.DeleteEmployee)                            // Delete employee
			employees.GET("/employee-id/:employee_id", employeeHandler.GetEmployeeByEmployeeID) // Get by employee ID
			employees.GET("/department/:department_id", employeeHandler.ListByDepartment)       // List by department
			employees.GET("/status/:status", employeeHandler.ListByStatus)                      // List by status

			// Employee status management (POST with ID in body)
			employees.POST("/terminate", employeeHandler.TerminateEmployee)   // Terminate employee
			employees.POST("/suspend", employeeHandler.SuspendEmployee)       // Suspend employee
			employees.POST("/archive", employeeHandler.ArchiveEmployee)       // Archive employee
			employees.POST("/reactivate", employeeHandler.ReactivateEmployee) // Reactivate employee

			// Multi-step onboarding routes
			onboardingHandler := employeeHandlers.NewOnboardingHandler()
			fileUploadHandler := employeeHandlers.NewFileUploadHandler()
			onboarding := employees.Group("/onboarding")
			{
				// File upload routes
				onboarding.POST("/upload-document", fileUploadHandler.UploadDocument) // Upload document (Step 6) - replaces if same type exists
				onboarding.POST("/upload-photo", fileUploadHandler.UploadPhoto)       // Upload photo (Step 1)
				onboarding.POST("/delete-file", fileUploadHandler.DeleteFile)         // Delete uploaded file from storage
				onboarding.POST("/delete-document", fileUploadHandler.DeleteDocument) // Delete document by ID (database + storage)

				// Specific routes first (to avoid conflicts with :employee_id)
				onboarding.POST("/draft", onboardingHandler.CreateDraft)                                 // Create new draft
				onboarding.GET("/drafts", onboardingHandler.ListDrafts)                                  // List all drafts
				onboarding.GET("/draft-employees", onboardingHandler.ListDraftEmployees)                 // List incomplete draft employees with details
				onboarding.GET("/draft-employee/:employee_id", onboardingHandler.GetDraftEmployeeByID)   // Get draft employee by employee_id to continue
				onboarding.POST("/completed-employee", onboardingHandler.GetCompletedEmployeeOnboarding) // Get completed employee onboarding data
				onboarding.POST("/complete", onboardingHandler.CompleteOnboarding)                       // Complete onboarding (legacy - by draft_id)

				// Draft ID-based routes (legacy support)
				onboarding.GET("/draft/:draft_id", onboardingHandler.GetDraft)             // Get draft with progress
				onboarding.POST("/draft/:draft_id/step/:step", onboardingHandler.SaveStep) // Save step data

				// Employee ID-based routes (primary - recommended) - must come after specific routes
				onboarding.GET("/:employee_id", onboardingHandler.GetDraftByEmployeeID)                             // Get draft by employee ID
				onboarding.POST("/:employee_id/step/:step", onboardingHandler.SaveStepByEmployeeID)                 // Save step by employee ID (JSON)
				onboarding.POST("/:employee_id/step/:step/upload", onboardingHandler.SaveStepByEmployeeIDWithFiles) // Save step with file uploads (multipart/form-data)
				onboarding.POST("/:employee_id/complete", onboardingHandler.CompleteOnboardingByEmployeeID)         // Complete by employee ID
			}

			// Post-onboarding task routes
			postOnboardingHandler := employeeHandlers.NewPostOnboardingTaskHandler()
			postOnboarding := employees.Group("/post-onboarding")
			{
				postOnboarding.GET("/types", postOnboardingHandler.GetTaskTypes)                                 // Get available task types
				postOnboarding.GET("/employees", postOnboardingHandler.ListEmployeesWithTaskCompletion)          // List all employees with task completion percentages
				postOnboarding.POST("/tasks", postOnboardingHandler.CreateTask)                                  // Create a task
				postOnboarding.POST("/tasks/bulk", postOnboardingHandler.BulkCreateTasks)                        // Create multiple tasks
				postOnboarding.GET("/tasks", postOnboardingHandler.ListTasks)                                    // List tasks with filters
				postOnboarding.POST("/tasks/get", postOnboardingHandler.GetTask)                                 // Get task by ID
				postOnboarding.POST("/tasks/update", postOnboardingHandler.UpdateTask)                           // Update task
				postOnboarding.POST("/tasks/complete", postOnboardingHandler.CompleteTask)                       // Complete task
				postOnboarding.POST("/tasks/delete", postOnboardingHandler.DeleteTask)                           // Delete task
				postOnboarding.GET("/tasks/employee/:employee_id", postOnboardingHandler.GetTasksByEmployeeID)   // Get tasks by employee ID
				postOnboarding.GET("/tasks/employee/:employee_id/summary", postOnboardingHandler.GetTaskSummary) // Get task summary for employee
			}

			// Employee Offboarding routes (require authentication - HR or Admin)
			offboardingHandler := employeeHandlers.NewOffboardingHandler()
			offboarding := employees.Group("/offboarding")
			{
				// Statistics
				offboarding.GET("/statistics", offboardingHandler.GetStatistics) // Get offboarding statistics

				// Workflow management
				offboarding.POST("/initiate", offboardingHandler.InitiateSeparation)                 // Initiate separation
				offboarding.GET("/workflows", offboardingHandler.ListWorkflows)                      // List workflows
				offboarding.GET("/workflows/:offboarding_id", offboardingHandler.GetWorkflowDetails) // Get workflow details
				offboarding.PATCH("/workflows/:offboarding_id", offboardingHandler.UpdateWorkflow)   // Update workflow

				// Clearances
				offboarding.GET("/workflows/:offboarding_id/clearances", offboardingHandler.GetClearances)                   // Get clearances
				offboarding.PATCH("/workflows/:offboarding_id/clearances/:clearance_id", offboardingHandler.UpdateClearance) // Update clearance

				// Asset returns
				offboarding.GET("/workflows/:offboarding_id/assets", offboardingHandler.GetAssetReturns)                     // Get assets
				offboarding.POST("/workflows/:offboarding_id/assets/:asset_id/return", offboardingHandler.RecordAssetReturn) // Record asset return
				offboarding.POST("/workflows/:offboarding_id/assets/:asset_id/issue", offboardingHandler.RecordAssetIssue)   // Record asset issue

				// Exit interview
				offboarding.POST("/workflows/:offboarding_id/exit-interview/schedule", offboardingHandler.ScheduleExitInterview) // Schedule exit interview
				offboarding.POST("/workflows/:offboarding_id/exit-interview/complete", offboardingHandler.CompleteExitInterview) // Complete exit interview
				offboarding.POST("/workflows/:offboarding_id/exit-interview/cancel", offboardingHandler.CancelExitInterview)     // Cancel exit interview

				// Final settlement
				offboarding.POST("/workflows/:offboarding_id/settlement/calculate", offboardingHandler.CalculateSettlement) // Calculate settlement
				offboarding.GET("/workflows/:offboarding_id/settlement", offboardingHandler.GetSettlement)                  // Get settlement
				offboarding.POST("/workflows/:offboarding_id/settlement/approve", offboardingHandler.ApproveSettlement)     // Approve settlement
				offboarding.POST("/workflows/:offboarding_id/settlement/pay", offboardingHandler.PaySettlement)             // Mark settlement as paid

				// Complete offboarding
				offboarding.POST("/workflows/:offboarding_id/complete", offboardingHandler.CompleteOffboarding) // Complete offboarding
			}
		}

		// Department routes (require authentication - HR or Admin)
		departments := v1.Group("/departments")
		departments.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			departments.POST("", departmentHandler.CreateDepartment)       // Create department
			departments.GET("", departmentHandler.ListDepartments)         // List departments
			departments.GET("/root", departmentHandler.GetRootDepartments) // Get root departments
			departments.GET("/:id", departmentHandler.GetDepartment)       // Get department by ID
			departments.PUT("/:id", departmentHandler.UpdateDepartment)    // Update department
			departments.DELETE("/:id", departmentHandler.DeleteDepartment) // Delete department
			// POST-only action-based endpoint
			departments.POST("/action", departmentHandler.HandleAction) // Action-based API
		}

		// Team routes (require authentication - HR or Admin)
		teams := v1.Group("/teams")
		teams.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			teams.POST("", teamHandler.CreateTeam)                                    // Create team
			teams.GET("", teamHandler.ListTeams)                                      // List teams
			teams.GET("/department/:department_id", teamHandler.GetTeamsByDepartment) // Get teams by department
			teams.GET("/:id", teamHandler.GetTeam)                                    // Get team by ID
			teams.PUT("/:id", teamHandler.UpdateTeam)                                 // Update team
			teams.DELETE("/:id", teamHandler.DeleteTeam)                              // Delete team
			// POST-only action-based endpoint
			teams.POST("/action", teamHandler.HandleAction) // Action-based API
		}

		// Job Position routes (require authentication - HR or Admin)
		positions := v1.Group("/job-positions")
		positions.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			positions.POST("", positionHandler.CreateJobPosition)        // Create job position
			positions.GET("", positionHandler.ListJobPositions)          // List job positions
			positions.POST("/get", positionHandler.GetJobPosition)       // Get job position by ID (POST with ID in body)
			positions.POST("/update", positionHandler.UpdateJobPosition) // Update job position (POST with ID in body)
			positions.POST("/delete", positionHandler.DeleteJobPosition) // Delete job position (POST with ID in body)
			// POST-only action-based endpoint
			positions.POST("/action", positionHandler.HandleAction) // Action-based API
		}

		// Location routes (require authentication - HR or Admin)
		locations := v1.Group("/locations")
		locations.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			locations.POST("", locationHandler.CreateLocation)           // Create location
			locations.GET("", locationHandler.ListLocations)             // List locations
			locations.GET("/head-office", locationHandler.GetHeadOffice) // Get head office
			locations.GET("/:id", locationHandler.GetLocation)           // Get location by ID
			locations.PUT("/:id", locationHandler.UpdateLocation)        // Update location
			locations.DELETE("/:id", locationHandler.DeleteLocation)     // Delete location
			// POST-only action-based endpoint
			locations.POST("/action", locationHandler.HandleAction) // Action-based API
		}

		// Cost Center routes (require authentication - HR or Admin)
		costCenters := v1.Group("/cost-centers")
		costCenters.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			costCenters.POST("", costCenterHandler.CreateCostCenter)       // Create cost center
			costCenters.GET("", costCenterHandler.ListCostCenters)         // List cost centers
			costCenters.GET("/root", costCenterHandler.GetRootCostCenters) // Get root cost centers
			costCenters.GET("/:id", costCenterHandler.GetCostCenter)       // Get cost center by ID
			costCenters.PUT("/:id", costCenterHandler.UpdateCostCenter)    // Update cost center
			costCenters.DELETE("/:id", costCenterHandler.DeleteCostCenter) // Delete cost center
			// POST-only action-based endpoint
			costCenters.POST("/action", costCenterHandler.HandleAction) // Action-based API
		}

		// Asset routes (require authentication - HR or Admin)
		assetHandler := assetHandlers.NewAssetHandler()
		assets := v1.Group("/assets")
		assets.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			assets.GET("", assetHandler.ListAssets)                      // List assets with pagination and filters
			assets.GET("/get", assetHandler.GetAsset)                    // Get asset by ID (query param: ?id=1)
			assets.GET("/types", assetHandler.GetAssetTypes)             // Get asset types
			assets.POST("", assetHandler.CreateAsset)                    // Create asset
			assets.POST("/update", assetHandler.UpdateAsset)             // Update asset
			assets.POST("/assign", assetHandler.AssignAsset)             // Assign asset to employee
			assets.POST("/reassign", assetHandler.ReassignAsset)         // Reassign asset to another employee
			assets.POST("/return", assetHandler.ReturnAsset)             // Return asset from employee
			assets.POST("/mark-for-repair", assetHandler.MarkForRepair)  // Mark asset for repair
			assets.POST("/complete-repair", assetHandler.CompleteRepair) // Complete asset repair
			assets.POST("/retire", assetHandler.RetireAsset)             // Retire asset
			assets.POST("/delete", assetHandler.DeleteAsset)             // Delete asset (only retired assets)
		}
	}
}

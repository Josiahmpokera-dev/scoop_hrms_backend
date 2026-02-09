package router

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	assetHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/handlers"
	auditHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/handlers"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/auth/handlers"
	biometricHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/handlers"
	costCenterHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/cost_centers/handlers"
	dashboardHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/dashboard/handlers"
	departmentHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/handlers"
	employeeHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/handlers"
	helpdeskHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/handlers"
	leaveHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/handlers"
	locationHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/handlers"
	orgChartHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization/handlers"
	organizationUnitHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization_units/handlers"
	organizationHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/handlers"
	payrollHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/handlers"
	positionHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/handlers"
	roleHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/handlers"
	shiftHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/handlers"
	teamHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/handlers"
	userHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/handlers"
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
	userHandler := userHandlers.NewUserHandler()
	roleHandler := roleHandlers.NewRoleHandler()
	auditHandler := auditHandlers.NewAuditHandler()
	dashboardHandler := dashboardHandlers.NewDashboardHandler()

	// API v1 routes (audit middleware logs every request)
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuditMiddleware())
	{
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			// Onboarding (complete signup + organization creation)
			auth.POST("/onboard", onboardingHandler.CompleteOnboarding)

			// Regular auth endpoints
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)

			// Protected routes
			auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
			auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
			auth.GET("/setup-wizard/status", middleware.AuthMiddleware(), onboardingHandler.GetSetupWizardStatus)
		}

		// Users routes (require authentication; list requires HR/Admin, transfer-role allows admin or role holder)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/special-roles", middleware.HRMiddleware(), userHandler.ListSpecialRoleUsers) // List IT, HR, Admin users (HR/Admin only)
			users.GET("", middleware.HRMiddleware(), userHandler.ListUsers)                          // List users (HR/Admin only, for transfer-role pickers)
			users.POST("/transfer-role", userHandler.TransferRole)                                   // Transfer role from one user to another (unchanged)
			users.POST("/assign-role", middleware.AdminMiddleware(), userHandler.AssignRole)         // Assign admin/hr/it to a normal user (Admin only)
			users.POST("/add-role", middleware.AdminMiddleware(), userHandler.AddRole)               // Add role to user (supports multiple roles) (Admin only)
			users.POST("/remove-role", middleware.AdminMiddleware(), userHandler.RemoveRole)         // Remove role from user (Admin only)
			users.POST("/set-roles", middleware.AdminMiddleware(), userHandler.SetRoles)             // Replace all roles for user (Admin only)
			users.GET("/:id/roles", middleware.HRMiddleware(), userHandler.GetUserRoles)             // Get all roles for a user (HR/Admin only)
			users.POST("/:id/suspend", middleware.HRMiddleware(), userHandler.SuspendUser)           // Suspend user (HR/Admin only; user cannot login)
			users.POST("/:id/unsuspend", middleware.HRMiddleware(), userHandler.UnsuspendUser)       // Unsuspend user (HR/Admin only; restores login)
			users.POST("/:id/block", middleware.HRMiddleware(), userHandler.BlockUser)               // Block user (HR/Admin only; user cannot login)
			users.POST("/:id/unblock", middleware.HRMiddleware(), userHandler.UnblockUser)           // Unblock user (HR/Admin only; restores login)
		}

		// Roles routes (RBAC roles list; HR/Admin only)
		roles := v1.Group("/roles")
		roles.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			roles.GET("", roleHandler.ListRoles) // List roles
		}

		// Feature Access Control routes
		featureHandler := roleHandlers.NewFeatureHandler()
		features := v1.Group("/features")
		features.Use(middleware.AuthMiddleware())
		{
			features.GET("", middleware.HRMiddleware(), featureHandler.GetAllFeatures)                   // List all features (HR/Admin)
			features.GET("/my-access", featureHandler.GetMyFeatureAccess)                                // Get my feature access (any authenticated user)
			features.GET("/roles/:id", middleware.HRMiddleware(), featureHandler.GetRoleFeatureAccess)    // Get role feature access (HR/Admin)
			features.PUT("/roles/:id", middleware.AdminMiddleware(), featureHandler.UpdateRolePermissions) // Update role permissions (Admin only)
		}

		// Security & Audit routes (Admin only)
		security := v1.Group("/security")
		security.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			security.GET("/audit", auditHandler.ListAuditLogs)                 // List audit logs with pagination, search, filters
			security.GET("/blocked-users", userHandler.ListBlockedUsers)       // List users blocked from login (rate-limit or manual)
			security.GET("/rate-limit/status", userHandler.GetRateLimitStatus) // Get rate-limit status for an email
		}

		// Dashboard routes (require authentication; role-based content)
		dashboard := v1.Group("/dashboard")
		dashboard.Use(middleware.AuthMiddleware())
		{
			dashboard.GET("/statistics", dashboardHandler.GetStatistics)                  // KPI stats (role-based: admin vs employee view)
			dashboard.GET("/announcements", dashboardHandler.GetAnnouncements)            // Company announcements
			dashboard.POST("/announcements/:id/read", dashboardHandler.MarkAnnouncementAsRead) // Mark announcement as read
			dashboard.GET("/quick-actions", dashboardHandler.GetQuickActions)             // Personalized quick action links
			dashboard.GET("/my-activity", dashboardHandler.GetMyActivity)                 // Employee's recent activities

			// Admin/HR only dashboard routes
			dashboard.GET("/events", middleware.HRMiddleware(), dashboardHandler.GetEvents)                      // Birthdays, anniversaries (HR/Admin)
			dashboard.GET("/pending-approvals", middleware.HRMiddleware(), dashboardHandler.GetPendingApprovals) // Pending approvals (HR/Admin)
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
		// Self-Service routes (require authentication - all employees)
		selfServiceHandler := employeeHandlers.NewSelfServiceHandler()
		selfService := v1.Group("/self-service")
		selfService.Use(middleware.AuthMiddleware()) // Only authentication required, not HR/Admin
		{
			// Profile
			selfService.GET("/profile", selfServiceHandler.GetProfile)
			selfService.POST("/profile/update", selfServiceHandler.UpdateProfile)
			selfService.GET("/profile/update-status", selfServiceHandler.GetProfileUpdateStatus)
			selfService.GET("/profile/documents", selfServiceHandler.GetDocuments)
			selfService.GET("/profile/id-card", selfServiceHandler.DownloadIDCard) // ID card generation (placeholder)

			// Payslips (placeholder - to be implemented with payroll module)
			selfService.GET("/payslips", selfServiceHandler.ListPayslips)
			selfService.GET("/payslips/:payslip_id", selfServiceHandler.GetPayslipDetails)
			selfService.GET("/payslips/:payslip_id/download", selfServiceHandler.DownloadPayslip)
			selfService.POST("/payslips/:payslip_id/email", selfServiceHandler.EmailPayslip)
			selfService.GET("/payslips/ytd-summary", selfServiceHandler.GetYTDSummary)
			selfService.POST("/payslips/query", selfServiceHandler.RaiseSalaryQuery)

			// Service Requests
			selfService.GET("/requests", selfServiceHandler.ListServiceRequests)
			selfService.GET("/requests/:request_id", selfServiceHandler.GetServiceRequestDetails)
			selfService.POST("/requests", selfServiceHandler.CreateServiceRequest)
			selfService.POST("/requests/:request_id/cancel", selfServiceHandler.CancelServiceRequest)
			// TODO: selfService.GET("/requests/:request_id/download", selfServiceHandler.DownloadRequestDocument)

			// People Directory
			selfService.GET("/directory", selfServiceHandler.SearchDirectory)
			selfService.GET("/directory/:employee_id", selfServiceHandler.GetDirectoryEmployeeDetails)

			// Assets (self-service)
			assetSelfServiceHandler := assetHandlers.NewAssetSelfServiceHandler()
			selfService.GET("/assets", assetSelfServiceHandler.GetMyAssignedAssets)
			selfService.GET("/assets/requests", assetSelfServiceHandler.ListAssetRequests)
			selfService.POST("/assets/requests", assetSelfServiceHandler.CreateAssetRequest)
			selfService.POST("/assets/requests/:request_id/cancel", assetSelfServiceHandler.CancelAssetRequest) // More specific route first
			selfService.GET("/assets/requests/:request_id", assetSelfServiceHandler.GetAssetRequestDetails)
			selfService.GET("/assets/issues", assetSelfServiceHandler.ListAssetIssues)
			selfService.POST("/assets/issues", assetSelfServiceHandler.CreateAssetIssue)
			selfService.GET("/assets/issues/:issue_id", assetSelfServiceHandler.GetAssetIssueDetails)
		}

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
			employees.GET("/managers", employeeHandler.ListManagers)                            // List all potential reporting managers

			// Employee status management (POST with ID in body)
			employees.POST("/terminate", employeeHandler.TerminateEmployee)   // Terminate employee
			employees.POST("/suspend", employeeHandler.SuspendEmployee)       // Suspend employee
			employees.POST("/archive", employeeHandler.ArchiveEmployee)       // Archive employee
			employees.POST("/reactivate", employeeHandler.ReactivateEmployee) // Reactivate employee

			// Employee document routes
			documentHandler := employeeHandlers.NewDocumentHandler()
			documents := employees.Group("/documents")
			{
				documents.GET("", documentHandler.ListEmployeesWithDocuments)                 // List all employees with documents
				documents.GET("/statistics", documentHandler.GetDocumentStatistics)           // Get document statistics
				documents.GET("/employee/:employee_id", documentHandler.GetEmployeeDocuments) // Get documents for specific employee
				documents.POST("/view", documentHandler.ViewDocument)                         // View document by ID
			}

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
				onboarding.GET("/non-employee-users", onboardingHandler.ListNonEmployeeUsers)            // List users not yet employees (for onboard existing user)
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
				postOnboarding.GET("/statistics", postOnboardingHandler.GetStatistics)                           // Get post-onboarding statistics
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
			positions.POST("", positionHandler.CreateJobPosition)                                 // Create job position
			positions.GET("", positionHandler.ListJobPositions)                                   // List job positions
			positions.GET("/department/:department_id", positionHandler.GetPositionsByDepartment) // Get positions by department
			positions.POST("/get", positionHandler.GetJobPosition)                                // Get job position by ID (POST with ID in body)
			positions.POST("/update", positionHandler.UpdateJobPosition)                          // Update job position (POST with ID in body)
			positions.POST("/delete", positionHandler.DeleteJobPosition)                          // Delete job position (POST with ID in body)
			// POST-only action-based endpoint
			positions.POST("/action", positionHandler.HandleAction) // Action-based API
		}

		// Location routes (require authentication - HR or Admin)
		locations := v1.Group("/locations")
		locations.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			locations.POST("", locationHandler.CreateLocation)               // Create location
			locations.GET("", locationHandler.ListLocations)                 // List locations
			locations.GET("/me", locationHandler.GetMyOrganizationLocations) // Get locations for current user's organization
			locations.GET("/head-office", locationHandler.GetHeadOffice)     // Get head office
			locations.GET("/:id", locationHandler.GetLocation)               // Get location by ID
			locations.PUT("/:id", locationHandler.UpdateLocation)            // Update location
			locations.DELETE("/:id", locationHandler.DeleteLocation)         // Delete location
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

			// Asset Request Management (HR/Admin)
			assetHRHandler := assetHandlers.NewAssetHRHandler()
			assets.GET("/requests", assetHRHandler.ListAssetRequests)
			assets.GET("/requests/:request_id", assetHRHandler.GetAssetRequestDetails)
			assets.POST("/requests/:request_id/approve", assetHRHandler.ApproveAssetRequest)
			assets.POST("/requests/:request_id/reject", assetHRHandler.RejectAssetRequest)
			assets.POST("/requests/:request_id/fulfill", assetHRHandler.FulfillAssetRequest)
			assets.POST("/:asset_id/reassign", assetHRHandler.ReassignAsset)
		}

		// Helpdesk routes
		helpdeskHandler := helpdeskHandlers.NewTicketHandler()
		helpdeskAgentHandler := helpdeskHandlers.NewTicketAgentHandler()
		helpdeskKBHandler := helpdeskHandlers.NewKnowledgeBaseHandler()

		// Employee helpdesk routes (self-service)
		helpdesk := v1.Group("/helpdesk")
		helpdesk.Use(middleware.AuthMiddleware())
		{
			helpdesk.GET("/ticket-categories", helpdeskHandler.GetTicketCategories)
			helpdesk.GET("/tickets", helpdeskHandler.ListMyTickets)
			helpdesk.GET("/tickets/recent", helpdeskAgentHandler.GetRecentTickets)
			helpdesk.GET("/tickets/:ticket_id", helpdeskHandler.GetTicketDetails)
			helpdesk.POST("/tickets", helpdeskHandler.CreateTicket)
			helpdesk.POST("/tickets/:ticket_id/comments", helpdeskHandler.AddComment)
			helpdesk.POST("/tickets/:ticket_id/attachments", helpdeskHandler.UploadAttachments)
			helpdesk.POST("/tickets/:ticket_id/close", helpdeskHandler.CloseTicket)
			helpdesk.POST("/tickets/:ticket_id/csat", helpdeskHandler.SubmitCSAT)
			// Knowledge Base
			helpdesk.GET("/knowledge-base/categories", helpdeskKBHandler.ListCategories)
			helpdesk.GET("/knowledge-base/articles", helpdeskKBHandler.ListArticles)
			helpdesk.GET("/knowledge-base/articles/:article_id", helpdeskKBHandler.GetArticle)
			helpdesk.POST("/knowledge-base/articles/:article_id/feedback", helpdeskKBHandler.SubmitFeedback)
		}

		// Agent/Admin helpdesk routes
		helpdeskAgent := v1.Group("/helpdesk/agent")
		helpdeskAgent.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			helpdeskAgent.GET("/tickets", helpdeskAgentHandler.ListTickets)
			helpdeskAgent.POST("/tickets/:ticket_id/assign", helpdeskAgentHandler.AssignTicket)
			helpdeskAgent.PATCH("/tickets/:ticket_id/status", helpdeskAgentHandler.UpdateTicketStatus)
			helpdeskAgent.POST("/tickets/:ticket_id/resolve", helpdeskAgentHandler.ResolveTicket)
		}

		// Helpdesk dashboard (Agent/Admin)
		helpdeskDashboard := v1.Group("/helpdesk/dashboard")
		helpdeskDashboard.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			helpdeskDashboard.GET("/statistics", helpdeskAgentHandler.GetStatistics)
		}

		// Helpdesk reports export (Agent/Admin)
		helpdeskReports := v1.Group("/helpdesk/reports")
		helpdeskReports.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			helpdeskReports.GET("/export", helpdeskAgentHandler.ExportReport)
		}

		// Biometric/BioTime routes
		biotimeHandler := biometricHandlers.NewBioTimeHandler()
		biometric := v1.Group("/biometric")
		biometric.Use(middleware.AuthMiddleware(), middleware.HRMiddleware()) // Require HR/Admin for biometric operations
		{
			biometric.GET("/biotime/test-connection", biotimeHandler.TestConnection)
			biometric.GET("/biotime/token", biotimeHandler.GetToken)
			biometric.POST("/biotime/refresh-token", biotimeHandler.RefreshToken)
			biometric.GET("/biotime/terminals", biotimeHandler.GetTerminals)
			biometric.GET("/biotime/device-status", biotimeHandler.GetDeviceStatus)
			biometric.GET("/biotime/transactions", biotimeHandler.GetTransactions)
			biometric.GET("/biotime/transactions/:id", biotimeHandler.GetTransaction)
			biometric.POST("/biotime/backfill", biotimeHandler.BackfillTransactions)
			// Daily attendance from database
			biometric.GET("/attendance/daily", biotimeHandler.GetDailyAttendance)
			// Late arrivals (exceptional cases)
			biometric.GET("/attendance/exceptional", biotimeHandler.GetExceptional)
		}

		// Shifts & Rosters routes
		shiftHandler := shiftHandlers.NewShiftHandler()
		rosterHandler := shiftHandlers.NewRosterHandler()
		swapRequestHandler := shiftHandlers.NewSwapRequestHandler()
		changeRequestHandler := shiftHandlers.NewRosterChangeRequestHandler()

		shifts := v1.Group("/shifts")
		shifts.Use(middleware.AuthMiddleware(), middleware.HRMiddleware()) // Require HR/Admin
		{
			// Statistics
			shifts.GET("/statistics", shiftHandler.GetStatistics)

			// Shift management
			shifts.GET("", shiftHandler.ListShifts)
			shifts.GET("/:shift_id", shiftHandler.GetShift)
			shifts.POST("", shiftHandler.CreateShift)
			shifts.PUT("/:shift_id", shiftHandler.UpdateShift)
			shifts.POST("/:shift_id/duplicate", shiftHandler.DuplicateShift)
			shifts.DELETE("/:shift_id", shiftHandler.DeleteShift)
		}

		rosters := v1.Group("/rosters")
		rosters.Use(middleware.AuthMiddleware(), middleware.HRMiddleware()) // Require HR/Admin
		{
			// Roster assignments
			rosters.GET("/assignments", rosterHandler.ListRosterAssignments)
			rosters.GET("/assignments/:assignment_id", rosterHandler.GetRosterAssignment)
			rosters.POST("/assignments", rosterHandler.CreateRosterAssignment)
			rosters.POST("/assignments/bulk", rosterHandler.BulkCreateRosterAssignments)
			rosters.PUT("/assignments/:assignment_id", rosterHandler.UpdateRosterAssignment)
			rosters.DELETE("/assignments/:assignment_id", rosterHandler.DeleteRosterAssignment)

			// Weekly roster view
			rosters.GET("/weekly", rosterHandler.GetWeeklyRosterView)

			// Auto-schedule
			rosters.POST("/auto-schedule", rosterHandler.AutoSchedule)

			// Publish roster
			rosters.POST("/publish", rosterHandler.PublishRoster)

			// Swap requests
			rosters.GET("/swap-requests", swapRequestHandler.ListSwapRequests)
			rosters.GET("/swap-requests/:request_id", swapRequestHandler.GetSwapRequest)
			rosters.POST("/swap-requests", swapRequestHandler.CreateSwapRequest)
			rosters.POST("/swap-requests/:request_id/approve", swapRequestHandler.ApproveSwapRequest)
			rosters.POST("/swap-requests/:request_id/reject", swapRequestHandler.RejectSwapRequest)

			// Roster change requests (HR/Admin - view all)
			rosters.GET("/change-requests", changeRequestHandler.ListRosterChangeRequests)
			rosters.GET("/change-requests/:request_id", changeRequestHandler.GetRosterChangeRequest)
			rosters.POST("/change-requests/:request_id/approve", changeRequestHandler.ApproveRosterChangeRequest)
			rosters.POST("/change-requests/:request_id/reject", changeRequestHandler.RejectRosterChangeRequest)
		}

		// Employee self-service routes for roster change requests
		employeeRosters := v1.Group("/self-service/rosters")
		employeeRosters.Use(middleware.AuthMiddleware()) // Only authentication required (all employees)
		{
			// Employee can view their own roster assignments
			employeeRosters.GET("/assignments", rosterHandler.GetMyRosterAssignments)
			employeeRosters.GET("/assignments/:assignment_id", rosterHandler.GetMyRosterAssignment)

			// Employee can create roster change requests
			employeeRosters.POST("/change-requests", changeRequestHandler.CreateRosterChangeRequest)
			// Employee can view their own change requests
			employeeRosters.GET("/change-requests", changeRequestHandler.ListRosterChangeRequests)
			employeeRosters.GET("/change-requests/:request_id", changeRequestHandler.GetRosterChangeRequest)
		}

		// Leave Management routes
		leaveTypeHandler := leaveHandlers.NewLeaveTypeHandler()
		leavePolicyHandler := leaveHandlers.NewLeavePolicyHandler()
		leaveRequestHandler := leaveHandlers.NewLeaveRequestHandler()
		holidayHandler := leaveHandlers.NewHolidayHandler()
		leaveCalendarHandler := leaveHandlers.NewLeaveCalendarHandler()

		// Leave Types (Admin/HR)
		leaveTypes := v1.Group("/leave/types")
		leaveTypes.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			leaveTypes.GET("", leaveTypeHandler.ListLeaveTypes)
			leaveTypes.GET("/:type_id", leaveTypeHandler.GetLeaveType)
			leaveTypes.POST("", leaveTypeHandler.CreateLeaveType)
			leaveTypes.PUT("/:type_id", leaveTypeHandler.UpdateLeaveType)
			leaveTypes.DELETE("/:type_id", leaveTypeHandler.DeleteLeaveType)
		}

		// Leave Policies (Admin/HR)
		leavePolicies := v1.Group("/leave/policies")
		leavePolicies.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			leavePolicies.GET("", leavePolicyHandler.ListLeavePolicies)
			leavePolicies.GET("/:policy_id", leavePolicyHandler.GetLeavePolicy)
			leavePolicies.POST("", leavePolicyHandler.CreateLeavePolicy)
			leavePolicies.PUT("/:policy_id", leavePolicyHandler.UpdateLeavePolicy)
			leavePolicies.DELETE("/:policy_id", leavePolicyHandler.DeleteLeavePolicy)
		}

		// Leave Requests - Employee endpoints
		leaveRequests := v1.Group("/leave")
		leaveRequests.Use(middleware.AuthMiddleware())
		{
			// Employee info and balances
			leaveRequests.GET("/employee-info", leaveRequestHandler.GetEmployeeInfo)
			leaveRequests.GET("/balances", leaveRequestHandler.GetEmployeeLeaveBalances)
			leaveRequests.GET("/policies/guidelines", leavePolicyHandler.GetPolicyGuidelines)

			// Calculate days
			leaveRequests.POST("/calculate-days", leaveRequestHandler.CalculateLeaveDays)

			// Leave requests
			leaveRequests.POST("/applications", leaveRequestHandler.CreateLeaveRequest)
			leaveRequests.GET("/requests", leaveRequestHandler.ListLeaveRequests)
			leaveRequests.GET("/requests/:request_id", leaveRequestHandler.GetLeaveRequest)
			leaveRequests.PUT("/requests/:request_id", leaveRequestHandler.UpdateLeaveRequest)
			leaveRequests.POST("/requests/:request_id/cancel", leaveRequestHandler.CancelLeaveRequest)
			leaveRequests.DELETE("/requests/:request_id", leaveRequestHandler.DeleteLeaveRequest)
		}

		// Leave Requests - HR/Admin approval endpoints
		leaveApprovals := v1.Group("/leave/requests")
		leaveApprovals.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			leaveApprovals.POST("/:request_id/approve", leaveRequestHandler.ApproveLeaveRequest)
			leaveApprovals.POST("/:request_id/reject", leaveRequestHandler.RejectLeaveRequest)
		}

		// Holidays (Admin/HR)
		holidays := v1.Group("/holidays")
		holidays.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			holidays.GET("", holidayHandler.ListHolidays)
			holidays.GET("/:holiday_id", holidayHandler.GetHoliday)
			holidays.POST("", holidayHandler.CreateHoliday)
			holidays.PUT("/:holiday_id", holidayHandler.UpdateHoliday)
			holidays.DELETE("/:holiday_id", holidayHandler.DeleteHoliday)
		}

		// Leave Calendar (All authenticated users)
		leaveCalendar := v1.Group("/leave/calendar")
		leaveCalendar.Use(middleware.AuthMiddleware())
		{
			leaveCalendar.GET("", leaveCalendarHandler.GetLeaveCalendar)
			leaveCalendar.GET("/today", leaveCalendarHandler.GetEmployeesOnLeaveToday)
			leaveCalendar.GET("/week", leaveCalendarHandler.GetWeeklyCalendar)
		}

		// ============ Payroll Management Routes ============
		payrollHandler := payrollHandlers.NewPayrollHandler()
		salaryStructureHandler := payrollHandlers.NewSalaryStructureHandler()
		payslipHandler := payrollHandlers.NewPayslipHandler()
		loanHandler := payrollHandlers.NewLoanHandler()
		complianceHandler := payrollHandlers.NewComplianceHandler()
		reportsHandler := payrollHandlers.NewReportsHandler()

		// Payroll Dashboard (HR/Admin)
		payrollDashboard := v1.Group("/payroll/dashboard")
		payrollDashboard.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			payrollDashboard.GET("", payrollHandler.GetDashboard)
		}

		// Payroll Runs (HR/Admin)
		payrollRuns := v1.Group("/payroll/runs")
		payrollRuns.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			payrollRuns.GET("", payrollHandler.ListPayrollRuns)
			payrollRuns.POST("", payrollHandler.CreatePayrollRun)
			payrollRuns.GET("/current", payrollHandler.GetCurrentPayrollRun)
			payrollRuns.GET("/:id", payrollHandler.GetPayrollRun)
			payrollRuns.PUT("/:id", payrollHandler.UpdatePayrollRun)
			payrollRuns.DELETE("/:id", payrollHandler.DeletePayrollRun)
			payrollRuns.GET("/:id/summary", payrollHandler.GetPayrollRunSummary)
			payrollRuns.GET("/:id/employees", payrollHandler.GetPayrollEmployees)
			payrollRuns.POST("/:id/pre-check", payrollHandler.RunPreCheck)
			payrollRuns.POST("/:id/calculate", payrollHandler.CalculatePayroll)
			payrollRuns.POST("/:id/advance", payrollHandler.AdvancePayrollStep)
			payrollRuns.POST("/:id/revert", payrollHandler.RevertPayrollStep)
			payrollRuns.POST("/:id/generate-payslips", payrollHandler.GeneratePayslips)
		}

		// Salary Structures (HR/Admin)
		salaryStructures := v1.Group("/payroll/salary-structures")
		salaryStructures.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			salaryStructures.GET("", salaryStructureHandler.ListSalaryStructures)
			salaryStructures.POST("", salaryStructureHandler.CreateSalaryStructure)
			salaryStructures.GET("/:id", salaryStructureHandler.GetSalaryStructure)
			salaryStructures.PUT("/:id", salaryStructureHandler.UpdateSalaryStructure)
			salaryStructures.DELETE("/:id", salaryStructureHandler.DeleteSalaryStructure)
			salaryStructures.POST("/simulate", salaryStructureHandler.SimulateSalary)
		}

		// Salary Components (HR/Admin)
		salaryComponents := v1.Group("/payroll/salary-components")
		salaryComponents.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			salaryComponents.GET("", salaryStructureHandler.ListSalaryComponents)
			salaryComponents.POST("", salaryStructureHandler.CreateSalaryComponent)
			salaryComponents.GET("/:id", salaryStructureHandler.GetSalaryComponent)
			salaryComponents.PUT("/:id", salaryStructureHandler.UpdateSalaryComponent)
			salaryComponents.DELETE("/:id", salaryStructureHandler.DeleteSalaryComponent)
		}

		// Payslips (HR/Admin)
		payslips := v1.Group("/payroll/payslips")
		payslips.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			payslips.GET("", payslipHandler.ListPayslips)
			payslips.GET("/summary", payslipHandler.GetPayslipSummary)
			payslips.GET("/:id", payslipHandler.GetPayslip)
			payslips.GET("/:id/download", payslipHandler.DownloadPayslip)
			payslips.POST("/:id/email", payslipHandler.SendPayslipEmail)
			payslips.POST("/bulk-download", payslipHandler.BulkDownloadPayslips)
			payslips.POST("/run/:runId/release", payslipHandler.ReleasePayslips)
			payslips.POST("/run/:runId/bulk-email", payslipHandler.BulkEmailPayslips)
		}

		// Loans & Advances (HR/Admin)
		loans := v1.Group("/payroll/loans")
		loans.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			loans.GET("", loanHandler.ListLoans)
			loans.POST("", loanHandler.CreateLoan)
			loans.GET("/summary", loanHandler.GetLoanSummary)
			loans.POST("/calculate-emi", loanHandler.CalculateEMI)
			loans.GET("/:id", loanHandler.GetLoan)
			loans.GET("/:id/schedule", loanHandler.GetRepaymentSchedule)
			loans.POST("/:id/approve", loanHandler.ApproveLoan)
			loans.POST("/:id/reject", loanHandler.RejectLoan)
			loans.POST("/:id/repayment", loanHandler.RecordRepayment)
		}

		// Compliance & Tax (HR/Admin)
		compliance := v1.Group("/payroll/compliance")
		compliance.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			compliance.GET("/summary", complianceHandler.GetComplianceSummary)
			compliance.GET("/rules", complianceHandler.GetAllComplianceRules)
			compliance.GET("/tax-slabs", complianceHandler.GetTaxSlabs)
			compliance.GET("/statutory-rules", complianceHandler.GetStatutoryRules)
			compliance.GET("/nhif-schedule", complianceHandler.GetNHIFSchedule)
			compliance.GET("/payments", complianceHandler.GetCompliancePayments)
			compliance.POST("/simulate-tax", complianceHandler.SimulateTax)
			compliance.POST("/calculate-monthly", complianceHandler.CalculateMonthlyStatutory)
			compliance.GET("/return", complianceHandler.GenerateComplianceReturn)
			compliance.POST("/seed", complianceHandler.SeedComplianceData)
		}

		// Payroll Reports (HR/Admin)
		payrollReports := v1.Group("/payroll/reports")
		payrollReports.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			payrollReports.GET("", reportsHandler.GetPayrollReport)
			payrollReports.GET("/kpis", reportsHandler.GetPayrollKPIs)
			payrollReports.GET("/department-cost", reportsHandler.GetDepartmentCostReport)
			payrollReports.GET("/export", reportsHandler.ExportPayrollReport)
			payrollReports.GET("/run/:runId/bank-file", reportsHandler.GetBankFileReport)
		}

		// Employee Self-Service Payroll Routes
		selfServicePayroll := v1.Group("/self-service/payroll")
		selfServicePayroll.Use(middleware.AuthMiddleware())
		{
			// Payslips
			selfServicePayroll.GET("/payslips", payslipHandler.GetMyPayslips)
			selfServicePayroll.GET("/payslips/latest", payslipHandler.GetMyLatestPayslip)
			selfServicePayroll.GET("/payslips/summary", payslipHandler.GetMySalarySlipSummary)
			selfServicePayroll.GET("/payslips/:id/download", payslipHandler.DownloadMyPayslip)

			// Loans
			selfServicePayroll.GET("/loans", loanHandler.GetMyLoans)
			selfServicePayroll.POST("/loans/apply", loanHandler.ApplyForLoan)
		}
	}
}

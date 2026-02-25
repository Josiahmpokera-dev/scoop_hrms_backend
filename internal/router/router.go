package router

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	assetHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/handlers"
	attendanceHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/handlers"
	auditHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/handlers"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/auth/handlers"
	biometricHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/handlers"
	bulkImportHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/bulk_import/handlers"
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
	performanceHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/handlers"
	positionHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/handlers"
	projectHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/handlers"
	recruitmentHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/handlers"
	roleHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/handlers"
	settingsHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/settings/handlers"
	shiftHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/handlers"
	teamHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/handlers"
	uploadHandlers "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/uploads/handlers"
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
	notificationHandler := dashboardHandlers.NewNotificationHandler()
	peopleDirectoryHandler := employeeHandlers.NewPeopleDirectoryHandler()

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
			users.POST("", middleware.HRMiddleware(), userHandler.CreateUser)                                  // Create user from employee (Admin/HR only)
			users.GET("/available-employees", middleware.HRMiddleware(), userHandler.ListEmployeesWithoutUser) // List employees without user accounts (for picker)
			users.GET("/special-roles", middleware.HRMiddleware(), userHandler.ListSpecialRoleUsers)           // List IT, HR, Admin users (HR/Admin only)
			users.GET("", middleware.HRMiddleware(), userHandler.ListUsers)                                    // List users (HR/Admin only, for transfer-role pickers)
			users.POST("/transfer-role", userHandler.TransferRole)                                             // Transfer role from one user to another (unchanged)
			users.POST("/assign-role", middleware.AdminMiddleware(), userHandler.AssignRole)                   // Assign admin/hr/it to a normal user (Admin only)
			users.POST("/add-role", middleware.AdminMiddleware(), userHandler.AddRole)                         // Add role to user (supports multiple roles) (Admin only)
			users.POST("/remove-role", middleware.AdminMiddleware(), userHandler.RemoveRole)                   // Remove role from user (Admin only)
			users.POST("/set-roles", middleware.AdminMiddleware(), userHandler.SetRoles)                       // Replace all roles for user (Admin only)
			users.PUT("/change-roles", middleware.AdminMiddleware(), userHandler.ChangeRoles)                  // Checkbox-style change roles (Admin only)
			users.GET("/:id", middleware.HRMiddleware(), userHandler.GetUser)                                  // Get user detail with roles (HR/Admin only)
			users.GET("/:id/roles", middleware.HRMiddleware(), userHandler.GetUserRoles)                       // Get all roles for a user (HR/Admin only)
			users.POST("/:id/suspend", middleware.HRMiddleware(), userHandler.SuspendUser)                     // Suspend user (HR/Admin only; user cannot login)
			users.POST("/:id/unsuspend", middleware.HRMiddleware(), userHandler.UnsuspendUser)                 // Unsuspend user (HR/Admin only; restores login)
			users.POST("/:id/block", middleware.HRMiddleware(), userHandler.BlockUser)                         // Block user (HR/Admin only; user cannot login)
			users.POST("/:id/unblock", middleware.HRMiddleware(), userHandler.UnblockUser)                     // Unblock user (HR/Admin only; restores login)
			users.POST("/reset-password", middleware.HRMiddleware(), userHandler.ResetPassword)                // Reset user password (HR/Admin; default: GreenTelecom@2026)
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
			features.GET("", middleware.HRMiddleware(), featureHandler.GetAllFeatures)                     // List all features (HR/Admin)
			features.GET("/my-access", featureHandler.GetMyFeatureAccess)                                  // Get my feature access (any authenticated user)
			features.GET("/roles/:id", middleware.HRMiddleware(), featureHandler.GetRoleFeatureAccess)     // Get role feature access (HR/Admin)
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
			dashboard.GET("/statistics", dashboardHandler.GetStatistics)                       // KPI stats (role-based: admin vs employee view)
			dashboard.GET("/announcements", dashboardHandler.GetAnnouncements)                 // Company announcements
			dashboard.POST("/announcements/:id/read", dashboardHandler.MarkAnnouncementAsRead) // Mark announcement as read
			dashboard.GET("/quick-actions", dashboardHandler.GetQuickActions)                  // Personalized quick action links
			dashboard.GET("/my-activity", dashboardHandler.GetMyActivity)                      // Employee's recent activities

			// Admin/HR only dashboard routes
			dashboard.GET("/events", middleware.HRMiddleware(), dashboardHandler.GetEvents)                      // Birthdays, anniversaries (HR/Admin)
			dashboard.GET("/pending-approvals", middleware.HRMiddleware(), dashboardHandler.GetPendingApprovals) // Pending approvals (HR/Admin)

			// Employee-specific dashboard routes
			employee := dashboard.Group("/employee")
			{
				employee.GET("/statistics", dashboardHandler.GetEmployeeStatistics) // Personal KPI stats (leave, hours, requests, payday)
				employee.GET("/my-activity", dashboardHandler.GetEmployeeActivity)  // Employee's recent activity feed
			}
		}

		// Notification routes (require authentication)
		notifications := v1.Group("/notifications")
		notifications.Use(middleware.AuthMiddleware())
		{
			notifications.GET("/count", notificationHandler.GetNotificationCount) // Get notification counts (badges)
			notifications.GET("", notificationHandler.GetNotifications)           // Get notifications list
		}

		// People Directory routes (require authentication)
		peopleDirectory := v1.Group("/people-directory")
		peopleDirectory.Use(middleware.AuthMiddleware())
		{
			peopleDirectory.GET("", peopleDirectoryHandler.SearchDirectory)        // Search/browse directory
			peopleDirectory.GET("/filters", peopleDirectoryHandler.GetFilters)     // Get filter options
			peopleDirectory.GET("/stats", peopleDirectoryHandler.GetStatistics)    // Get directory statistics
			peopleDirectory.GET("/:id", peopleDirectoryHandler.GetEmployeeProfile) // Get employee profile
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
			employees.POST("", employeeHandler.OnboardEmployee)                                       // Onboard new employee (legacy - single step)
			employees.GET("", employeeHandler.ListEmployees)                                          // List all employees
			employees.GET("/:id", employeeHandler.GetEmployee)                                        // Get employee by ID
			employees.PUT("/:id", employeeHandler.UpdateEmployee)                                     // Update employee
			employees.DELETE("/:id", employeeHandler.DeleteEmployee)                                  // Delete employee
			employees.GET("/employee-id/:employee_id", employeeHandler.GetEmployeeByEmployeeID)       // Get by employee ID
			employees.GET("/department/:department_id", employeeHandler.ListByDepartment)             // List by department
			employees.GET("/status/:status", employeeHandler.ListByStatus)                            // List by status
			employees.GET("/managers", employeeHandler.ListManagers)                                  // List potential reporting managers (optional: ?department_id=X to filter & flag department head)
			employees.GET("/department-manager/:department_id", employeeHandler.GetDepartmentManager) // Get suggested manager for a department (department head)

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
			departments.POST("", departmentHandler.CreateDepartment)                     // Create department
			departments.GET("", departmentHandler.ListDepartments)                       // List departments
			departments.GET("/root", departmentHandler.GetRootDepartments)               // Get root departments
			departments.GET("/:id", departmentHandler.GetDepartment)                     // Get department by ID
			departments.PUT("/:id", departmentHandler.UpdateDepartment)                  // Update department
			departments.DELETE("/:id", departmentHandler.DeleteDepartment)               // Delete department
			departments.POST("/:id/assign-head", departmentHandler.AssignDepartmentHead) // Assign employee as department head
			departments.POST("/:id/remove-head", departmentHandler.RemoveDepartmentHead) // Remove department head assignment
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
		enrollmentHandler := biometricHandlers.NewEnrollmentHandler()

		biometric := v1.Group("/biometric")
		biometric.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			biometric.GET("/biotime/test-connection", biotimeHandler.TestConnection)
			biometric.GET("/biotime/token", biotimeHandler.GetToken)
			biometric.POST("/biotime/refresh-token", biotimeHandler.RefreshToken)
			biometric.GET("/biotime/terminals", biotimeHandler.GetTerminals)
			biometric.GET("/biotime/device-status", biotimeHandler.GetDeviceStatus)
			biometric.GET("/biotime/transactions", biotimeHandler.GetTransactions)
			biometric.GET("/biotime/transactions/:id", biotimeHandler.GetTransaction)
			biometric.POST("/biotime/backfill", biotimeHandler.BackfillTransactions)
			biometric.GET("/attendance/daily", biotimeHandler.GetDailyAttendance)
			biometric.GET("/attendance/exceptional", biotimeHandler.GetExceptional)

			// Enrollment — link employees to biometric device users
			biometric.GET("/enrollments", enrollmentHandler.List)
			biometric.GET("/enrollments/statistics", enrollmentHandler.GetStatistics)
			biometric.GET("/enrollments/unlinked", enrollmentHandler.GetUnlinked)
			biometric.GET("/enrollments/device-users", enrollmentHandler.GetDeviceUsers)
			biometric.POST("/enrollments/link", enrollmentHandler.Link)
			biometric.POST("/enrollments/bulk-link", enrollmentHandler.BulkLink)
			biometric.POST("/enrollments/auto-link", enrollmentHandler.AutoLink)
			biometric.DELETE("/enrollments/:employeeId/unlink", enrollmentHandler.Unlink)

			// Merged attendance (biometric + employee data)
			biometric.GET("/attendance/merged", enrollmentHandler.GetMergedAttendance)
			biometric.GET("/attendance/merged/:employeeId", enrollmentHandler.GetEmployeeAttendance)
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

		// Leave Types — Read (any authenticated user)
		leaveTypesRead := v1.Group("/leave/types")
		leaveTypesRead.Use(middleware.AuthMiddleware())
		{
			leaveTypesRead.GET("", leaveTypeHandler.ListLeaveTypes)
			leaveTypesRead.GET("/:type_id", leaveTypeHandler.GetLeaveType)
		}

		// Leave Types — Write (Admin/HR only)
		leaveTypesWrite := v1.Group("/leave/types")
		leaveTypesWrite.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			leaveTypesWrite.POST("", leaveTypeHandler.CreateLeaveType)
			leaveTypesWrite.PUT("/:type_id", leaveTypeHandler.UpdateLeaveType)
			leaveTypesWrite.DELETE("/:type_id", leaveTypeHandler.DeleteLeaveType)
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

		// Leave Requests - Employee endpoints (any authenticated user)
		leaveRequests := v1.Group("/leave")
		leaveRequests.Use(middleware.AuthMiddleware())
		{
			// Employee info and balances
			leaveRequests.GET("/employee-info", leaveRequestHandler.GetEmployeeInfo)
			leaveRequests.GET("/balances", leaveRequestHandler.GetEmployeeLeaveBalances)
			leaveRequests.GET("/policies/guidelines", leavePolicyHandler.GetPolicyGuidelines)

			// Active leave types (read-only for form dropdown)
			leaveRequests.GET("/active-types", leaveRequestHandler.GetActiveLeaveTypes)

			// Holidays for employee calendar / day calculation
			leaveRequests.GET("/holidays", leaveRequestHandler.GetLeaveHolidays)

			// Calculate days
			leaveRequests.POST("/calculate-days", leaveRequestHandler.CalculateLeaveDays)

			// Leave requests (employee's own)
			leaveRequests.POST("/applications", leaveRequestHandler.CreateLeaveRequest)
			leaveRequests.GET("/requests", leaveRequestHandler.ListLeaveRequests)
			leaveRequests.GET("/requests/:request_id", leaveRequestHandler.GetLeaveRequest)
			leaveRequests.PUT("/requests/:request_id", leaveRequestHandler.UpdateLeaveRequest)
			leaveRequests.POST("/requests/:request_id/cancel", leaveRequestHandler.CancelLeaveRequest)
			leaveRequests.DELETE("/requests/:request_id", leaveRequestHandler.DeleteLeaveRequest)
		}

		// Leave Requests - HR/Admin management & approval endpoints
		leaveApprovals := v1.Group("/leave")
		leaveApprovals.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			// List all leave requests across all employees (HR/Admin view)
			leaveApprovals.GET("/admin/requests", leaveRequestHandler.ListAllLeaveRequests)

			// Approval actions
			leaveApprovals.POST("/requests/:request_id/approve", leaveRequestHandler.ApproveLeaveRequest)
			leaveApprovals.POST("/requests/:request_id/reject", leaveRequestHandler.RejectLeaveRequest)
			leaveApprovals.POST("/requests/:request_id/return-for-info", leaveRequestHandler.ReturnForInfo)
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

		// ============ Attendance Module: Timesheets & Overtime ============
		timesheetHandler := attendanceHandlers.NewTimesheetHandler()
		overtimeHandler := attendanceHandlers.NewOvertimeHandler()
		attendanceReportsHandler := attendanceHandlers.NewAttendanceReportsHandler()

		// Timesheet - Employee Self-Service (all authenticated users)
		timesheets := v1.Group("/attendance/timesheets")
		timesheets.Use(middleware.AuthMiddleware())
		{
			// My timesheets
			timesheets.GET("", timesheetHandler.GetMyTimesheets)                       // List my weekly timesheets
			timesheets.GET("/:timesheet_id", timesheetHandler.GetTimesheetByID)        // Get timesheet by ID
			timesheets.POST("/entries", timesheetHandler.CreateEntry)                  // Create single entry
			timesheets.POST("/entries/bulk", timesheetHandler.BulkCreateEntries)       // Bulk create entries for a week
			timesheets.PUT("/entries/:entry_id", timesheetHandler.UpdateEntry)         // Update entry
			timesheets.DELETE("/entries/:entry_id", timesheetHandler.DeleteEntry)      // Delete entry
			timesheets.POST("/:timesheet_id/submit", timesheetHandler.SubmitTimesheet) // Submit for approval
			timesheets.POST("/:timesheet_id/recall", timesheetHandler.RecallTimesheet) // Recall submitted timesheet
			timesheets.POST("/copy-last-week", timesheetHandler.CopyLastWeek)          // Copy entries from last week
			timesheets.GET("/stats", timesheetHandler.GetEmployeeStats)                // Employee utilization stats
			timesheets.GET("/team", timesheetHandler.GetTeamTimesheets)                // Manager: team timesheets
		}

		// Timesheet - Approvals (Admin/HR)
		timesheetApprovals := v1.Group("/attendance/timesheets/approvals")
		timesheetApprovals.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			timesheetApprovals.GET("", timesheetHandler.GetPendingApprovals)             // List pending approvals
			timesheetApprovals.POST("/:timesheet_id", timesheetHandler.ApproveTimesheet) // Approve/reject timesheet
		}

		// Overtime - Employee Self-Service (all authenticated users)
		overtime := v1.Group("/attendance/overtime")
		overtime.Use(middleware.AuthMiddleware())
		{
			overtime.POST("/requests", overtimeHandler.CreateOTRequest)                    // Submit OT request
			overtime.GET("/requests", overtimeHandler.GetMyOTRequests)                     // List my OT requests
			overtime.GET("/requests/:request_id", overtimeHandler.GetOTRequestByID)        // Get OT request by ID
			overtime.POST("/requests/:request_id/cancel", overtimeHandler.CancelOTRequest) // Cancel OT request
			overtime.GET("/policies", overtimeHandler.ListPolicies)                        // View OT policies
			overtime.GET("/policies/:policy_id", overtimeHandler.GetPolicy)                // View specific policy
			overtime.GET("/stats", overtimeHandler.GetOvertimeStats)                       // OT statistics
		}

		// Overtime - Admin/HR Management
		overtimeAdmin := v1.Group("/attendance/overtime/admin")
		overtimeAdmin.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			// Policy management
			overtimeAdmin.POST("/policies", overtimeHandler.CreatePolicy)              // Create OT policy
			overtimeAdmin.PUT("/policies/:policy_id", overtimeHandler.UpdatePolicy)    // Update OT policy
			overtimeAdmin.DELETE("/policies/:policy_id", overtimeHandler.DeletePolicy) // Delete OT policy

			// Request management
			overtimeAdmin.GET("/requests", overtimeHandler.ListAllRequests)                       // List all OT requests
			overtimeAdmin.GET("/pending", overtimeHandler.GetPendingApprovals)                    // List pending approvals
			overtimeAdmin.POST("/requests/:request_id/approve", overtimeHandler.ApproveOTRequest) // Approve/reject OT request
		}

		// Attendance Reports (Admin/HR)
		attendanceReports := v1.Group("/attendance/reports")
		attendanceReports.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			// New endpoints per ATTENDANCE_REPORTS_API.md
			attendanceReports.GET("/summary", attendanceReportsHandler.GetSummary)
			attendanceReports.GET("/trends", attendanceReportsHandler.GetTrends)
			attendanceReports.GET("/by-department", attendanceReportsHandler.GetDepartmentStats)
			attendanceReports.GET("/compliance", attendanceReportsHandler.GetComplianceViolations)
			attendanceReports.GET("/overtime", attendanceReportsHandler.GetOvertimeAnalysis)          // Replaces previous overtime endpoint
			attendanceReports.GET("/overtime-analysis", attendanceReportsHandler.GetOvertimeAnalysis) // Alias for frontend compatibility
			attendanceReports.POST("/export", attendanceReportsHandler.ExportReport)

			// Legacy/Other endpoints
			attendanceReports.GET("/timesheets", attendanceReportsHandler.GetTimesheetSummaryReport)              // Timesheet summary
			attendanceReports.GET("/project-utilization", attendanceReportsHandler.GetProjectUtilizationReport)   // Project utilization
			attendanceReports.GET("/employee-utilization", attendanceReportsHandler.GetEmployeeUtilizationReport) // Employee utilization
		}

		// ============ Projects & Daily Tasks ============
		projectHandler := projectHandlers.NewProjectHandler()
		dailyTaskHandler := projectHandlers.NewDailyTaskHandler()
		recruitmentHandler := recruitmentHandlers.NewRecruitmentHandler()

		// Projects - Management (HR/Admin/Manager)
		projects := v1.Group("/projects")
		projects.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			projects.GET("/statistics", projectHandler.GetStatistics)         // Project statistics
			projects.GET("", projectHandler.ListProjects)                     // List all projects
			projects.POST("", projectHandler.CreateProject)                   // Create project
			projects.GET("/:id", projectHandler.GetProject)                   // Get project details
			projects.PUT("/:id", projectHandler.UpdateProject)                // Update project
			projects.DELETE("/:id", projectHandler.DeleteProject)             // Delete project
			projects.POST("/:id/members", projectHandler.AddMembers)          // Add members to project
			projects.POST("/:id/members/remove", projectHandler.RemoveMember) // Remove member from project
			projects.GET("/:id/members", projectHandler.ListMembers)          // List project members
			projects.GET("/:id/progress", projectHandler.GetProjectProgress)  // Get project progress
			projects.POST("/assign", projectHandler.AssignProject)            // Assign project to employee(s)
			projects.POST("/unassign", projectHandler.UnassignProject)        // Unassign employee from project
		}

		// Daily Tasks - Management (HR/Admin view all tasks)
		dailyTasks := v1.Group("/daily-tasks")
		dailyTasks.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			dailyTasks.GET("", dailyTaskHandler.ListDailyTasks)                                    // List all daily tasks
			dailyTasks.GET("/categories", dailyTaskHandler.GetTaskCategories)                      // Get task categories
			dailyTasks.GET("/project/:project_id/summary", dailyTaskHandler.GetProjectTaskSummary) // Project task summary
			dailyTasks.GET("/:id", dailyTaskHandler.GetDailyTask)                                  // Get task details
		}

		// Self-Service: Projects & Daily Tasks (all authenticated employees)
		selfServiceProjects := v1.Group("/self-service/projects")
		selfServiceProjects.Use(middleware.AuthMiddleware())
		{
			selfServiceProjects.GET("", projectHandler.ListMyProjects)          // List my assigned projects
			selfServiceProjects.GET("/:id", projectHandler.GetMyProjectDetails) // Get assigned project details + progress
		}

		selfServiceDailyTasks := v1.Group("/self-service/daily-tasks")
		selfServiceDailyTasks.Use(middleware.AuthMiddleware())
		{
			selfServiceDailyTasks.POST("", dailyTaskHandler.CreateDailyTask)             // Log a daily task
			selfServiceDailyTasks.GET("", dailyTaskHandler.GetMyTasks)                   // List my daily tasks
			selfServiceDailyTasks.GET("/today", dailyTaskHandler.GetMyTodayTasks)        // Get today's tasks
			selfServiceDailyTasks.GET("/summary", dailyTaskHandler.GetMyTaskSummary)     // Get my task summary
			selfServiceDailyTasks.GET("/categories", dailyTaskHandler.GetTaskCategories) // Get task categories
			selfServiceDailyTasks.GET("/:id", dailyTaskHandler.GetDailyTask)             // Get task details
			selfServiceDailyTasks.PUT("/:id", dailyTaskHandler.UpdateDailyTask)          // Update my task
			selfServiceDailyTasks.DELETE("/:id", dailyTaskHandler.DeleteDailyTask)       // Delete my task
		}

		// =====================================================================
		// Recruitment Module
		// =====================================================================
		
		// Public Routes (No Auth)
		recruitmentPublic := v1.Group("/recruitment/public")
		{
			recruitmentPublic.POST("/apply", recruitmentHandler.SubmitApplication)
			recruitmentPublic.GET("/openings", recruitmentHandler.ListPublicJobOpenings)
		}

		// Protected Routes
		recruitment := v1.Group("/recruitment")
		recruitment.Use(middleware.AuthMiddleware())
		{
			// Requisitions
			recruitment.POST("/requisitions", middleware.PermissionMiddleware("recruitment:create"), recruitmentHandler.CreateRequisition)
			recruitment.GET("/requisitions", middleware.PermissionMiddleware("recruitment:read"), recruitmentHandler.ListRequisitions)
			recruitment.GET("/requisitions/:id", middleware.PermissionMiddleware("recruitment:read"), recruitmentHandler.GetRequisition)
			recruitment.POST("/requisitions/:id/approval", middleware.PermissionMiddleware("recruitment:approve"), recruitmentHandler.ApproveRequisition)
			recruitment.PUT("/requisitions/:id", middleware.PermissionMiddleware("recruitment:update"), recruitmentHandler.UpdateRequisition)

			// Job Openings
			recruitment.POST("/openings", middleware.PermissionMiddleware("recruitment:create"), recruitmentHandler.CreateJobOpening)
			recruitment.GET("/openings", middleware.PermissionMiddleware("recruitment:read"), recruitmentHandler.ListJobOpenings)
			recruitment.POST("/openings/:id/publish", middleware.PermissionMiddleware("recruitment:publish"), recruitmentHandler.PublishJobOpening)

			// Candidates
			recruitment.GET("/candidates", middleware.PermissionMiddleware("recruitment:manage_candidates"), recruitmentHandler.ListCandidates)
			recruitment.PATCH("/candidates/:id/stage", middleware.PermissionMiddleware("recruitment:manage_candidates"), recruitmentHandler.UpdateCandidateStage)

			// Interviews
			recruitment.POST("/interviews", middleware.PermissionMiddleware("recruitment:manage_candidates"), recruitmentHandler.ScheduleInterview)
			recruitment.POST("/interviews/:id/feedback", middleware.PermissionMiddleware("recruitment:manage_candidates"), recruitmentHandler.SubmitFeedback)

			// Offers
			recruitment.POST("/offers", middleware.PermissionMiddleware("recruitment:manage_candidates"), recruitmentHandler.CreateOffer)
			recruitment.POST("/offers/:id/approval", middleware.PermissionMiddleware("recruitment:approve"), recruitmentHandler.ApproveOffer)
			recruitment.POST("/offers/:id/send", middleware.PermissionMiddleware("recruitment:manage_candidates"), recruitmentHandler.SendOffer)

			// Talent Pool
			recruitment.POST("/talent-pool", middleware.PermissionMiddleware("recruitment:manage_candidates"), recruitmentHandler.AddToTalentPool)
			recruitment.GET("/talent-pool/search", middleware.PermissionMiddleware("recruitment:manage_candidates"), recruitmentHandler.SearchTalentPool)
		}

		// =====================================================================
		// Settings — Menu Visibility
		// =====================================================================
		menuVisHandler := settingsHandlers.NewMenuVisibilityHandler()

		// Public endpoint: any authenticated user can fetch hidden keys
		settingsPublic := v1.Group("/settings")
		settingsPublic.Use(middleware.AuthMiddleware())
		{
			settingsPublic.GET("/menu-visibility/active", menuVisHandler.GetActive)
		}

		// Admin-only endpoints
		settingsAdmin := v1.Group("/settings")
		settingsAdmin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			settingsAdmin.GET("/menu-visibility", menuVisHandler.GetAll)
			settingsAdmin.PUT("/menu-visibility", menuVisHandler.BulkUpdate)
			settingsAdmin.PATCH("/menu-visibility/:menuKey", menuVisHandler.ToggleSingle)
			settingsAdmin.POST("/menu-visibility/reset", menuVisHandler.Reset)
		}

		// =====================================================================
		// File Uploads (S3 / Local)
		// =====================================================================
		uploadHandler := uploadHandlers.NewUploadHandler()

		uploads := v1.Group("/uploads")
		uploads.Use(middleware.AuthMiddleware())
		{
			uploads.POST("/file", uploadHandler.UploadFile)
			uploads.POST("/image", uploadHandler.UploadImage)
			uploads.POST("/document", uploadHandler.UploadDocument)
			uploads.DELETE("", uploadHandler.DeleteFile)
			uploads.GET("/info", uploadHandler.GetStorageInfo)
		}

		// =====================================================================
		// Bulk Import (Excel Upload)
		// =====================================================================
		bulkHandler := bulkImportHandlers.NewBulkImportHandler()

		bulkImport := v1.Group("/bulk-import")
		bulkImport.Use(middleware.AuthMiddleware(), middleware.HRMiddleware())
		{
			bulkImport.GET("/templates/employees", bulkHandler.DownloadEmployeeTemplate)
			bulkImport.GET("/templates/departments", bulkHandler.DownloadDepartmentTemplate)
			bulkImport.GET("/templates/positions", bulkHandler.DownloadPositionTemplate)
			bulkImport.POST("/employees", bulkHandler.ImportEmployees)
			bulkImport.POST("/departments", bulkHandler.ImportDepartments)
			bulkImport.POST("/positions", bulkHandler.ImportPositions)
		}

		// =====================================================================
		// Performance Management Module
		// =====================================================================
		goalHandler := performanceHandlers.NewGoalHandler()
		deptTargetHandler := performanceHandlers.NewDepartmentTargetHandler()
		empTargetHandler := performanceHandlers.NewEmployeeTargetHandler()
		appraisalHandler := performanceHandlers.NewAppraisalHandler()
		feedback360Handler := performanceHandlers.NewFeedback360Handler()
		talentReviewHandler := performanceHandlers.NewTalentReviewHandler()
		reportHandler := performanceHandlers.NewReportHandler()

		performance := v1.Group("/performance")
		performance.Use(middleware.AuthMiddleware())
		{
			// --- Dashboard (any authenticated user) ---
			dashboard := performance.Group("/dashboard")
			{
				dashboard.GET("/stats", goalHandler.GetDashboardStats)
				dashboard.GET("/upcoming-actions", goalHandler.GetUpcomingActions)
			}

			// --- Goals & OKRs ---
			goals := performance.Group("/goals")
			{
				goals.GET("", goalHandler.ListGoals)                                  // List goals
				goals.GET("/stats", goalHandler.GetGoalStats)                         // Goal statistics
				goals.GET("/alignment", goalHandler.GetAlignmentMap)                  // Goal alignment map (Manager+)
				goals.POST("", goalHandler.CreateGoal)                                // Create goal
				goals.POST("/assign", goalHandler.AssignGoal)                         // Assign goal to employee (Manager+)
				goals.GET("/:id", goalHandler.GetGoal)                                // Get goal detail
				goals.PATCH("/:id", goalHandler.UpdateGoal)                           // Update goal
				goals.DELETE("/:id", goalHandler.DeleteGoal)                          // Delete goal
				goals.POST("/:id/submit-for-approval", goalHandler.SubmitForApproval) // Submit for manager approval
				goals.PUT("/:id/approve", goalHandler.ApproveGoal)                    // Approve/reject goal
				goals.POST("/:id/request-completion", goalHandler.RequestCompletion)  // Request completion verification
				goals.PUT("/:id/verify-completion", goalHandler.VerifyCompletion)     // Verify/reject completion
				goals.POST("/:id/link-project", goalHandler.LinkProject)              // Link goal to project
				goals.DELETE("/:id/unlink-project", goalHandler.UnlinkProject)        // Unlink goal from project
				// Key Results
				goals.POST("/:id/key-results", goalHandler.CreateKeyResult)
				goals.PATCH("/:id/key-results/:krId", goalHandler.UpdateKeyResult)
				goals.DELETE("/:id/key-results/:krId", goalHandler.DeleteKeyResult)
				// Check-ins
				goals.POST("/:id/check-ins", goalHandler.CreateCheckIn)
			}

			// --- Department Targets (Manager+) ---
			deptTargets := performance.Group("/department-targets")
			{
				deptTargets.GET("", deptTargetHandler.ListTargets)
				deptTargets.GET("/:id", deptTargetHandler.GetTarget)
				deptTargets.POST("", middleware.ManagerMiddleware(), deptTargetHandler.CreateTarget)
				deptTargets.PATCH("/:id", middleware.ManagerMiddleware(), deptTargetHandler.UpdateTarget)
				deptTargets.DELETE("/:id", middleware.ManagerMiddleware(), deptTargetHandler.DeleteTarget)
				deptTargets.POST("/:id/progress", middleware.ManagerMiddleware(), deptTargetHandler.UpdateProgress)
				deptTargets.PUT("/:id/milestones/:milestoneId/complete", middleware.ManagerMiddleware(), deptTargetHandler.CompleteMilestone)
				deptTargets.POST("/:id/link-goal", middleware.ManagerMiddleware(), deptTargetHandler.LinkGoal)
				deptTargets.POST("/:id/link-project", middleware.ManagerMiddleware(), deptTargetHandler.LinkProject)
				deptTargets.DELETE("/:id/unlink-project", middleware.ManagerMiddleware(), deptTargetHandler.UnlinkProject)
			}

			// --- Employee Targets (Manager+) ---
			empTargets := performance.Group("/employee-targets")
			{
				empTargets.GET("", empTargetHandler.ListTargets)
				empTargets.GET("/:id", empTargetHandler.GetTarget)
				empTargets.POST("", middleware.ManagerMiddleware(), empTargetHandler.CreateTarget)
				empTargets.POST("/bulk-assign", middleware.ManagerMiddleware(), empTargetHandler.BulkAssignTargets)
				empTargets.PATCH("/:id", middleware.ManagerMiddleware(), empTargetHandler.UpdateTarget)
				empTargets.DELETE("/:id", middleware.ManagerMiddleware(), empTargetHandler.DeleteTarget)
				empTargets.POST("/:id/progress", empTargetHandler.UpdateProgress)
			}

			// --- Appraisal Cycles (HR/Admin) ---
			appraisalCycles := performance.Group("/appraisal-cycles")
			{
				appraisalCycles.GET("", appraisalHandler.ListCycles)
				appraisalCycles.GET("/active", appraisalHandler.GetActiveCycle)
				appraisalCycles.GET("/:id", appraisalHandler.GetCycle)
				appraisalCycles.POST("", middleware.HRMiddleware(), appraisalHandler.CreateCycle)
				appraisalCycles.PUT("/:id", middleware.HRMiddleware(), appraisalHandler.UpdateCycle)
				appraisalCycles.PATCH("/:id/status", middleware.HRMiddleware(), appraisalHandler.ChangeCycleStatus)
				appraisalCycles.DELETE("/:id", middleware.AdminMiddleware(), appraisalHandler.DeleteCycle)
			}

			// --- Appraisals ---
			appraisals := performance.Group("/appraisals")
			{
				appraisals.GET("", appraisalHandler.ListAppraisals)
				appraisals.GET("/summary", appraisalHandler.GetAppraisalSummary)
				appraisals.GET("/:id", appraisalHandler.GetAppraisal)
				appraisals.POST("/:id/self-review", appraisalHandler.SubmitSelfReview)
				appraisals.POST("/:id/manager-review", appraisalHandler.SubmitManagerReview)
				appraisals.POST("/:id/calibration", middleware.HRMiddleware(), appraisalHandler.SubmitCalibration)
				appraisals.PUT("/:id/send-back", appraisalHandler.SendBack)
				appraisals.PUT("/:id/finalize", middleware.HRMiddleware(), appraisalHandler.Finalize)
			}

			// --- 360° Feedback ---
			feedback360 := performance.Group("/feedback-360")
			{
				feedback360.GET("", feedback360Handler.ListCampaigns)
				feedback360.GET("/:id", feedback360Handler.GetCampaign)
				feedback360.POST("", middleware.HRMiddleware(), feedback360Handler.LaunchCampaign)
			}

			// --- Talent Review / 9-Box ---
			talentReview := performance.Group("/talent-review")
			{
				talentReview.GET("", middleware.ManagerMiddleware(), talentReviewHandler.ListReviews)
				talentReview.GET("/employees/:id", middleware.ManagerMiddleware(), talentReviewHandler.GetReview)
				talentReview.PATCH("/employees/:id", middleware.HRMiddleware(), talentReviewHandler.UpdateReview)
				talentReview.POST("/calibration-sessions", middleware.HRMiddleware(), talentReviewHandler.CreateCalibrationSession)
				talentReview.POST("/employees/:id/succession", middleware.HRMiddleware(), talentReviewHandler.AddSuccessionPlan)
			}

			// --- Performance Reports (HR/Admin) ---
			reports := performance.Group("/reports")
			reports.Use(middleware.HRMiddleware())
			{
				reports.GET("/summary", reportHandler.GetSummary)
				reports.GET("/rating-distribution", reportHandler.GetRatingDistribution)
				reports.GET("/department-summary", reportHandler.GetDepartmentSummary)
			}

			// --- Project-Performance Alignment ---
			projectAlignment := performance.Group("/project-alignment")
			{
				projectAlignment.GET("/:projectId", goalHandler.GetDashboardStats)   // Reuse for now
				projectAlignment.GET("/summary", reportHandler.GetDepartmentSummary) // Reuse for now
			}
		}
	}
}

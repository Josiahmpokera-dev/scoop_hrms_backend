package handlers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/storage"
	"github.com/gin-gonic/gin"
)

// OnboardingHandler handles multi-step employee onboarding
type OnboardingHandler struct {
	onboardingService *services.OnboardingService
}

// NewOnboardingHandler creates a new onboarding handler
func NewOnboardingHandler() *OnboardingHandler {
	return &OnboardingHandler{
		onboardingService: services.NewOnboardingService(),
	}
}

// CreateDraft creates a new onboarding draft
// @Summary Create onboarding draft
// @Description Create a new draft for multi-step employee onboarding. Optionally accepts step 1 data to save immediately.
// @Tags Employee Onboarding
// @Accept json
// @Produce json
// @Param request body models.CreateDraftRequest false "Optional step 1 data"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/draft [post]
func (h *OnboardingHandler) CreateDraft(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var createdBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		createdBy = &userObj.ID
	}

	// Parse optional request body for step data and existing-user onboarding opts
	var req models.CreateDraftRequest
	var stepData map[string]interface{}
	var stepNumber int = 0
	var createDraftOpts *services.CreateDraftOpts

	if c.Request.ContentLength > 0 {
		// Try to parse as CreateDraftRequest or SaveDraftRequest format
		var rawData map[string]interface{}
		if err := c.ShouldBindJSON(&rawData); err == nil {
			// Optional: onboard existing user (no credentials at end)
			if lid, ok := rawData["linked_user_id"]; ok && lid != nil {
				switch v := lid.(type) {
				case float64:
					uid := uint(v)
					createDraftOpts = &services.CreateDraftOpts{LinkedUserID: &uid}
				}
			} else if em, ok := rawData["existing_user_email"]; ok && em != nil {
				if s, ok := em.(string); ok && s != "" {
					createDraftOpts = &services.CreateDraftOpts{ExistingUserEmail: &s}
				}
			}
			// Check if it's in SaveDraftRequest format (has "step" and "data" fields)
			if stepVal, hasStep := rawData["step"]; hasStep {
				// It's in SaveDraftRequest format
				if stepNum, ok := stepVal.(float64); ok {
					stepNumber = int(stepNum)
					if stepNumber < 1 || stepNumber > 10 {
						response.BadRequest(c, "Invalid step number. Must be between 1 and 10", nil)
						return
					}
					if dataVal, hasData := rawData["data"]; hasData {
						if dataMap, ok := dataVal.(map[string]interface{}); ok {
							stepData = dataMap
						}
					}
				}
			} else if err := c.ShouldBindJSON(&req); err == nil {
				// Try to parse as CreateDraftRequest
				if req.LinkedUserID != nil {
					createDraftOpts = &services.CreateDraftOpts{LinkedUserID: req.LinkedUserID}
				} else if req.ExistingUserEmail != nil && *req.ExistingUserEmail != "" {
					createDraftOpts = &services.CreateDraftOpts{ExistingUserEmail: req.ExistingUserEmail}
				}
				if req.Step != nil && *req.Step >= 1 && *req.Step <= 10 {
					stepNumber = *req.Step
					if len(req.Data) > 0 {
						stepData = req.Data
					}
				} else if len(req.Step1Data) > 0 {
					// Legacy: step1_data field
					stepNumber = 1
					stepData = req.Step1Data
				}
			} else {
				// Try to detect step 1 data by checking for first_name field
				if _, hasFirstName := rawData["first_name"]; hasFirstName {
					stepNumber = 1
					stepData = rawData
				}
			}
		}
	}

	// Create draft (opts = nil for new person; opts set for existing user onboarding)
	draft, err := h.onboardingService.CreateDraft(tenantID, createdBy, createDraftOpts)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// If step data is provided, save it immediately
	if stepNumber > 0 && stepData != nil && len(stepData) > 0 {
		updatedDraft, err := h.onboardingService.SaveStep(draft.ID, stepNumber, stepData, createdBy)
		if err != nil {
			// Draft was created but step save failed - return draft with error message
			response.Success(c, fmt.Sprintf("Draft created but step %d save failed: %s", stepNumber, err.Error()), draft)
			return
		}
		// Get full draft with step data to return (use employee_id if available)
		var draftResponse *models.GetDraftResponse
		if updatedDraft.EmployeeID != nil {
			tenantID := middleware.GetTenantID(c)
			draftResponse, err = h.onboardingService.GetDraftByEmployeeID(*updatedDraft.EmployeeID, tenantID)
		} else {
			draftResponse, err = h.onboardingService.GetDraft(updatedDraft.ID)
		}
		if err == nil {
			// Return full draft response with step data
			response.Created(c, fmt.Sprintf("Onboarding draft created and step %d saved successfully. Employee ID: %s", stepNumber, *updatedDraft.EmployeeID), draftResponse)
			return
		}
		// Fallback to returning updated draft if GetDraft fails
		response.Created(c, fmt.Sprintf("Onboarding draft created and step %d saved successfully", stepNumber), updatedDraft)
		return
	}

	response.Created(c, "Onboarding draft created successfully", draft)
}

// SaveStep saves data for a specific onboarding step (by draft_id - legacy support)
// @Summary Save onboarding step by draft ID
// @Description Save data for a specific step in the onboarding process using draft ID
// @Tags Employee Onboarding
// @Accept json
// @Produce json
// @Param draft_id path int true "Draft ID"
// @Param step path int true "Step number (1-10)"
// @Param request body models.SaveDraftRequest true "Step data"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/draft/{draft_id}/step/{step} [post]
func (h *OnboardingHandler) SaveStep(c *gin.Context) {
	draftIDStr := c.Param("draft_id")
	draftID, err := strconv.ParseUint(draftIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid draft ID", nil)
		return
	}

	stepStr := c.Param("step")
	step, err := strconv.Atoi(stepStr)
	if err != nil || step < 1 || step > 10 {
		response.BadRequest(c, "Invalid step number. Must be between 1 and 10", nil)
		return
	}

	var req models.SaveDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	// Validate step matches URL parameter
	if req.Step != step {
		response.BadRequest(c, "Step number in URL must match step in request body", nil)
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	draft, err := h.onboardingService.SaveStep(uint(draftID), step, req.Data, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Step saved successfully", draft)
}

// SaveStepByEmployeeID saves data for a specific onboarding step by employee ID
// @Summary Save onboarding step by employee ID
// @Description Save data for a specific step in the onboarding process using employee ID
// @Tags Employee Onboarding
// @Accept json
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Param step path int true "Step number (1-10)"
// @Param request body models.SaveDraftRequest true "Step data"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/{employee_id}/step/{step} [post]
func (h *OnboardingHandler) SaveStepByEmployeeID(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		response.BadRequest(c, "Employee ID is required", nil)
		return
	}

	stepStr := c.Param("step")
	step, err := strconv.Atoi(stepStr)
	if err != nil || step < 1 || step > 10 {
		response.BadRequest(c, "Invalid step number. Must be between 1 and 10", nil)
		return
	}

	var req models.SaveDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	// Validate step matches URL parameter
	if req.Step != step {
		response.BadRequest(c, "Step number in URL must match step in request body", nil)
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	draft, err := h.onboardingService.SaveStepByEmployeeID(employeeID, tenantID, step, req.Data, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Step saved successfully", draft)
}

func (h *OnboardingHandler) UpdateStep(c *gin.Context) {
	draftIDStr := c.Param("draft_id")
	draftID, err := strconv.ParseUint(draftIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid draft ID", nil)
		return
	}

	stepStr := c.Param("step")
	step, err := strconv.Atoi(stepStr)
	if err != nil || step < 1 || step > 10 {
		response.BadRequest(c, "Invalid step number. Must be between 1 and 10", nil)
		return
	}

	var req models.SaveDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	if req.Step != step {
		response.BadRequest(c, "Step number in URL must match step in request body", nil)
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	draft, err := h.onboardingService.SaveStep(uint(draftID), step, req.Data, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Step updated successfully", draft)
}

func (h *OnboardingHandler) UpdateStepByEmployeeID(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		response.BadRequest(c, "Employee ID is required", nil)
		return
	}

	stepStr := c.Param("step")
	step, err := strconv.Atoi(stepStr)
	if err != nil || step < 1 || step > 10 {
		response.BadRequest(c, "Invalid step number. Must be between 1 and 10", nil)
		return
	}

	var req models.SaveDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	if req.Step != step {
		response.BadRequest(c, "Step number in URL must match step in request body", nil)
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	draft, err := h.onboardingService.SaveStepByEmployeeID(employeeID, tenantID, step, req.Data, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Step updated successfully", draft)
}

// GetDraft retrieves a draft with all step data and progress (by draft_id - legacy support)
// @Summary Get onboarding draft by draft ID
// @Description Get draft with all step data and completion progress using draft ID
// @Tags Employee Onboarding
// @Produce json
// @Param draft_id path int true "Draft ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/draft/{draft_id} [get]
func (h *OnboardingHandler) GetDraft(c *gin.Context) {
	draftIDStr := c.Param("draft_id")
	draftID, err := strconv.ParseUint(draftIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid draft ID", nil)
		return
	}

	draftResponse, err := h.onboardingService.GetDraft(uint(draftID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Draft retrieved successfully", draftResponse)
}

// GetDraftByEmployeeID retrieves a draft with all step data and progress by employee ID
// @Summary Get onboarding draft by employee ID
// @Description Get draft with all step data, progress, finished and unfinished steps using employee ID
// @Tags Employee Onboarding
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/{employee_id} [get]
func (h *OnboardingHandler) GetDraftByEmployeeID(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		response.BadRequest(c, "Employee ID is required", nil)
		return
	}

	tenantID := middleware.GetTenantID(c)
	draftResponse, err := h.onboardingService.GetDraftByEmployeeID(employeeID, tenantID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Draft retrieved successfully", draftResponse)
}

// CompleteOnboarding finalizes the onboarding and creates the employee (by draft_id - legacy support)
// @Summary Complete onboarding by draft ID
// @Description Finalize onboarding and create employee record using draft ID
// @Tags Employee Onboarding
// @Accept json
// @Produce json
// @Param request body models.CompleteOnboardingRequest true "Complete onboarding request"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/complete [post]
func (h *OnboardingHandler) CompleteOnboarding(c *gin.Context) {
	var req models.CompleteOnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	employee, credentials, err := h.onboardingService.CompleteOnboarding(req.DraftID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Prepare response: credentials_created and linked_to_existing_user so frontend can show the right message
	credentialsCreated := credentials != nil
	responseData := map[string]interface{}{
		"employee":                 employee,
		"credentials_created":      credentialsCreated,
		"linked_to_existing_user":  !credentialsCreated,
	}
	if credentials != nil {
		responseData["credentials"] = credentials
	}

	response.Created(c, "Employee onboarding completed successfully", responseData)
}

// CompleteOnboardingByEmployeeID finalizes the onboarding and creates the employee by employee ID
// @Summary Complete onboarding by employee ID
// @Description Finalize onboarding and create employee record using employee ID
// @Tags Employee Onboarding
// @Accept json
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/{employee_id}/complete [post]
func (h *OnboardingHandler) CompleteOnboardingByEmployeeID(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		response.BadRequest(c, "Employee ID is required", nil)
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	employee, credentials, err := h.onboardingService.CompleteOnboardingByEmployeeID(employeeID, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	credentialsCreated := credentials != nil
	responseData := map[string]interface{}{
		"employee":                employee,
		"credentials_created":     credentialsCreated,
		"linked_to_existing_user": !credentialsCreated,
	}
	if credentials != nil {
		responseData["credentials"] = credentials
	}

	response.Created(c, "Employee onboarding completed successfully", responseData)
}

// ListNonEmployeeUsers lists users who are not yet linked to any employee (for starting "onboard existing user" flow).
// Use the returned user id as linked_user_id when calling POST /employees/onboarding/draft with { "linked_user_id": <id> }.
// @Summary List non-employee users
// @Description Get users who are not yet employees; use their id to start onboarding with linked_user_id
// @Tags Employee Onboarding
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20) maximum(100)
// @Param search query string false "Search by email, first name, last name, username"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/non-employee-users [get]
func (h *OnboardingHandler) ListNonEmployeeUsers(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}
	search := strings.TrimSpace(c.Query("search"))

	users, total, err := h.onboardingService.ListNonEmployeeUsers(tenantID, page, pageSize, search)
	if err != nil {
		response.InternalServerError(c, "Failed to list non-employee users", err.Error())
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	response.SuccessWithMeta(c, "Non-employee users retrieved successfully", users, meta)
}

// ListDrafts lists all onboarding drafts for the tenant
// @Summary List onboarding drafts
// @Description Get list of all onboarding drafts
// @Tags Employee Onboarding
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/drafts [get]
func (h *OnboardingHandler) ListDrafts(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	
	drafts, err := h.onboardingService.ListDrafts(tenantID)
	if err != nil {
		response.InternalServerError(c, "Failed to list drafts", err.Error())
		return
	}

	response.Success(c, "Drafts retrieved successfully", drafts)
}

// ListDraftEmployees lists all incomplete draft employees with their filled details
// @Summary List draft employees
// @Description List all employees in draft (incomplete onboarding) with their filled details (with pagination)
// @Tags Employee Onboarding
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20) maximum(100)
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/draft-employees [get]
func (h *OnboardingHandler) ListDraftEmployees(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)
	
	draftEmployees, total, err := h.onboardingService.ListDraftEmployeesWithPagination(tenantID, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list draft employees", err.Error())
		return
	}

	// Calculate total pages
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	response.SuccessWithMeta(c, "Draft employees retrieved successfully", draftEmployees, meta)
}

// GetCompletedEmployeeOnboarding gets a completed employee's onboarding data
// @Summary Get completed employee onboarding data
// @Description Get all onboarding data for a completed employee by employee_id
// @Tags Employee Onboarding
// @Accept json
// @Produce json
// @Param request body models.GetDraftEmployeeRequest true "Employee ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/completed-employee [post]
func (h *OnboardingHandler) GetCompletedEmployeeOnboarding(c *gin.Context) {
	var req models.GetDraftEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", "employee_id is required in request body")
		return
	}

	if req.EmployeeID == "" {
		response.ValidationError(c, "Validation failed", "employee_id is required")
		return
	}

	tenantID := middleware.GetTenantID(c)
	
	completedData, err := h.onboardingService.GetCompletedEmployeeOnboarding(req.EmployeeID, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Completed employee onboarding data retrieved successfully", completedData)
}

// SaveStepByEmployeeIDWithFiles handles saving step data with file uploads (multipart/form-data)
// Supports Step 1 (photo) and Step 6 (documents) with file uploads
// @Summary Save onboarding step with file uploads
// @Description Save step data with file uploads for Step 1 (photo) or Step 6 (documents)
// @Tags Employee Onboarding
// @Accept multipart/form-data
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Param step path int true "Step number (1 for photo, 6 for documents)"
// @Param file formData file false "File to upload (for Step 1: photo, for Step 6: document)"
// @Param files formData file false "Multiple files (for Step 6: documents)"
// @Param data formData string false "JSON string of step data (excluding file fields)"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/:employee_id/step/:step/upload [post]
func (h *OnboardingHandler) SaveStepByEmployeeIDWithFiles(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		response.ValidationError(c, "Validation failed", "employee_id is required")
		return
	}

	stepStr := c.Param("step")
	step, err := strconv.Atoi(stepStr)
	if err != nil || step < 1 || step > 10 {
		response.ValidationError(c, "Validation failed", "step must be between 1 and 10")
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	// Get step data from form (if provided as JSON string)
	var stepData map[string]interface{}
	dataStr := c.PostForm("data")
	if dataStr != "" {
		if err := json.Unmarshal([]byte(dataStr), &stepData); err != nil {
			response.ValidationError(c, "Validation failed", "invalid data JSON format")
			return
		}
	} else {
		stepData = make(map[string]interface{})
	}

	// Handle file uploads based on step
	storageService := storage.NewStorageService()

	if step == 1 {
		// Step 1: Photo upload
		file, err := c.FormFile("file")
		if err == nil && file != nil {
			// Upload photo
			photoURL, fileSize, mimeType, uploadErr := storageService.UploadFile(file, "photos", employeeID)
			if uploadErr != nil {
				response.BadRequest(c, uploadErr.Error(), nil)
				return
			}
			stepData["photo_url"] = storageService.ResolveURL(c.Request, photoURL)
			stepData["photo_file_size"] = fileSize
			stepData["photo_mime_type"] = mimeType
		}
	} else if step == 6 {
		// Step 6: Document uploads
		form, err := c.MultipartForm()
		if err == nil && form != nil {
			files := form.File["files"]
			if len(files) == 0 {
				// Try single file
				file, err := c.FormFile("file")
				if err == nil && file != nil {
					files = []*multipart.FileHeader{file}
				}
			}

			var documents []map[string]interface{}
			for i, file := range files {
				// Get document type from form (files_document_type[i] or document_type)
				documentType := c.PostForm(fmt.Sprintf("files_document_type_%d", i))
				if documentType == "" {
					documentType = c.PostForm("document_type")
				}
				if documentType == "" {
					documentType = "other" // Default
				}

				// Get description
				description := c.PostForm(fmt.Sprintf("files_description_%d", i))
				if description == "" {
					description = c.PostForm("description")
				}

				// Upload file
				fileURL, fileSize, mimeType, uploadErr := storageService.UploadFile(file, "documents", employeeID)
				if uploadErr != nil {
					response.BadRequest(c, uploadErr.Error(), nil)
					return
				}

				documents = append(documents, map[string]interface{}{
					"document_type": documentType,
					"file_name":     file.Filename,
					"file_url":      storageService.ResolveURL(c.Request, fileURL),
					"file_size":     fileSize,
					"mime_type":     mimeType,
					"description":   description,
				})
			}

			if len(documents) > 0 {
				stepData["documents"] = documents
			}
		}
	}

	// Save step data
	draft, err := h.onboardingService.SaveStepByEmployeeID(employeeID, tenantID, step, stepData, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Step saved successfully with file uploads", draft)
}

// GetDraftEmployeeByID gets a draft employee by employee_id to continue onboarding
// @Summary Get draft employee by employee ID
// @Description Get a draft employee by employee_id with all filled details to continue onboarding
// @Tags Employee Onboarding
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/draft-employee/:employee_id [get]
func (h *OnboardingHandler) GetDraftEmployeeByID(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		response.ValidationError(c, "Validation failed", "employee_id is required")
		return
	}

	tenantID := middleware.GetTenantID(c)
	
	// Use existing GetDraftByEmployeeID which returns full draft with all step data
	draftResponse, err := h.onboardingService.GetDraftByEmployeeID(employeeID, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Check if already completed
	if draftResponse.IsCompleted {
		response.BadRequest(c, "Onboarding already completed for this employee", nil)
		return
	}

	response.Success(c, "Draft employee retrieved successfully", draftResponse)
}

func (h *OnboardingHandler) EditEmployeeStep(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		response.ValidationError(c, "Validation failed", "employee_id is required")
		return
	}

	stepStr := c.Param("step")
	step, err := strconv.Atoi(stepStr)
	if err != nil || step < 1 || step > 10 {
		response.ValidationError(c, "Validation failed", "step must be between 1 and 10")
		return
	}

	var req models.EmployeeEditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	if req.Step != step {
		response.BadRequest(c, "Step number in URL must match step in request body", nil)
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	result, err := h.onboardingService.EditEmployeeStep(employeeID, tenantID, step, req.Data, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Employee step updated successfully", result)
}

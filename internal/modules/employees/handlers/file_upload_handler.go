package handlers

import (
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/storage"
	"github.com/gin-gonic/gin"
)

// FileUploadHandler handles file upload operations
type FileUploadHandler struct {
	storageService    *storage.StorageService
	onboardingService *services.OnboardingService
	draftRepo         *employeeRepos.OnboardingDraftRepository
	documentRepo      *employeeRepos.EmployeeDocumentRepository
}

// NewFileUploadHandler creates a new file upload handler
func NewFileUploadHandler() *FileUploadHandler {
	return &FileUploadHandler{
		storageService:    storage.NewStorageService(),
		onboardingService: services.NewOnboardingService(),
		draftRepo:         employeeRepos.NewOnboardingDraftRepository(),
		documentRepo:      employeeRepos.NewEmployeeDocumentRepository(),
	}
}

// UploadDocument handles document upload for onboarding
// @Summary Upload document for employee onboarding
// @Description Upload a document file for employee onboarding (Step 6)
// @Tags Employee Onboarding
// @Accept multipart/form-data
// @Produce json
// @Param employee_id formData string true "Employee ID (e.g., EMP001)"
// @Param document_type formData string true "Document type (identity, work_permit, education, contract, tax_statutory, other)"
// @Param file formData file true "Document file"
// @Param description formData string false "Document description"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/upload-document [post]
func (h *FileUploadHandler) UploadDocument(c *gin.Context) {
	// Get employee ID
	employeeID := c.PostForm("employee_id")
	if employeeID == "" {
		response.ValidationError(c, "Validation failed", "employee_id is required")
		return
	}

	// Get document type
	documentType := c.PostForm("document_type")
	if documentType == "" {
		response.ValidationError(c, "Validation failed", "document_type is required")
		return
	}

	// Validate document type
	validTypes := []string{"identity", "work_permit", "education", "contract", "tax_statutory", "other"}
	isValidType := false
	for _, t := range validTypes {
		if documentType == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		response.ValidationError(c, "Validation failed", "invalid document_type. Must be one of: identity, work_permit, education, contract, tax_statutory, other")
		return
	}

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		response.ValidationError(c, "Validation failed", "file is required")
		return
	}

	// Get description (optional)
	description := c.PostForm("description")

	// Get user ID for uploaded_by
	user, _ := c.Get("user")
	var uploadedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		uploadedBy = &userObj.ID
	}

	// Upload file
	fileURL, fileSize, mimeType, err := h.storageService.UploadFile(file, "documents", employeeID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert relative URL to full URL
	fullFileURL := h.getFullURL(c, fileURL)

	// Get tenant ID
	tenantID := middleware.GetTenantID(c)

	// Get or create draft for this employee
	draft, err := h.draftRepo.FindByEmployeeIDString(employeeID, tenantID)
	if err != nil {
		// Draft doesn't exist, create one
		draft = &models.EmployeeOnboardingDraft{
			TenantID:    tenantID,
			EmployeeID:  &employeeID,
			Progress:    0,
			IsCompleted: false,
		}
		if err := h.draftRepo.Create(draft); err != nil {
			response.BadRequest(c, "Failed to create draft: "+err.Error(), nil)
			return
		}
	}

	// Check if a document with the same type already exists for this draft
	docType := models.DocumentType(documentType)
	existingDoc, err := h.documentRepo.FindByDraftIDAndType(draft.ID, docType)

	var doc *models.EmployeeDocument
	var isUpdate bool

	if err == nil && existingDoc != nil {
		// Document with this type already exists - update/replace it
		isUpdate = true

		// Delete the old file from storage
		if existingDoc.FileURL != "" {
			if deleteErr := h.storageService.DeleteFile(existingDoc.FileURL); deleteErr != nil {
				// Log error but continue with update
				// The old file might not exist or might be in use
			}
		}

		// Update existing document
		existingDoc.FileName = file.Filename
		existingDoc.FileURL = fileURL
		existingDoc.FileSize = &fileSize
		existingDoc.MimeType = &mimeType
		existingDoc.UploadedBy = uploadedBy
		if description != "" {
			existingDoc.Description = &description
		} else {
			existingDoc.Description = nil
		}

		// Update document in database
		if err := h.documentRepo.Update(existingDoc); err != nil {
			response.BadRequest(c, "Failed to update document: "+err.Error(), nil)
			return
		}

		doc = existingDoc
	} else {
		// No existing document with this type - create new one
		isUpdate = false
		doc = &models.EmployeeDocument{
			DraftID:          &draft.ID,
			EmployeeIDString: &employeeID,
			DocumentType:     docType,
			FileName:         file.Filename,
			FileURL:          fileURL, // Store relative URL in database
			FileSize:         &fileSize,
			MimeType:         &mimeType,
			UploadedBy:       uploadedBy,
		}
		if description != "" {
			doc.Description = &description
		}

		// Save new document to database
		if err := h.documentRepo.Create(doc); err != nil {
			response.BadRequest(c, "Failed to save document: "+err.Error(), nil)
			return
		}
	}

	// Update draft progress - mark step 6 as completed if not already
	if !draft.HasStepCompleted(models.StepDocuments) {
		draft.AddCompletedStep(models.StepDocuments)
		draft.Progress = draft.CalculateProgress()
		if err := h.draftRepo.Update(draft); err != nil {
			// Log error but don't fail the request
			// The document is already saved
		}
	}

	// Return file information with document ID
	message := "Document uploaded successfully"
	if isUpdate {
		message = "Document updated successfully (replaced existing document of the same type)"
	}

	response.Success(c, message, gin.H{
		"id":            doc.ID,
		"file_url":      fullFileURL, // Return full URL in response
		"file_name":     file.Filename,
		"file_size":     fileSize,
		"mime_type":     mimeType,
		"document_type": documentType,
		"description":   description,
		"employee_id":   employeeID,
		"draft_id":      draft.ID,
		"uploaded_by":   uploadedBy,
		"is_update":     isUpdate,
		"created_at":    doc.CreatedAt,
		"updated_at":    doc.UpdatedAt,
	})
}

// getFullURL converts a relative URL to a full URL using the request's scheme and host
func (h *FileUploadHandler) getFullURL(c *gin.Context, relativeURL string) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	// Check X-Forwarded-Proto header for reverse proxy setups
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	host := c.Request.Host
	if host == "" {
		host = "localhost:8080" // Default fallback
	}

	// Remove leading slash from relativeURL if present
	if len(relativeURL) > 0 && relativeURL[0] == '/' {
		relativeURL = relativeURL[1:]
	}

	return fmt.Sprintf("%s://%s/%s", scheme, host, relativeURL)
}

// UploadPhoto handles photo upload for employee onboarding
// @Summary Upload photo for employee onboarding
// @Description Upload a photo file for employee onboarding (Step 1)
// @Tags Employee Onboarding
// @Accept multipart/form-data
// @Produce json
// @Param employee_id formData string true "Employee ID (e.g., EMP001)"
// @Param file formData file true "Photo file (image)"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/upload-photo [post]
func (h *FileUploadHandler) UploadPhoto(c *gin.Context) {
	// Get employee ID
	employeeID := c.PostForm("employee_id")
	if employeeID == "" {
		response.ValidationError(c, "Validation failed", "employee_id is required")
		return
	}

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		response.ValidationError(c, "Validation failed", "file is required")
		return
	}

	// Upload file
	fileURL, fileSize, mimeType, err := h.storageService.UploadFile(file, "photos", employeeID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert relative URL to full URL
	fullFileURL := h.getFullURL(c, fileURL)

	// Return file information
	response.Success(c, "Photo uploaded successfully", gin.H{
		"photo_url":   fullFileURL, // Return full URL in response
		"file_name":   file.Filename,
		"file_size":   fileSize,
		"mime_type":   mimeType,
		"employee_id": employeeID,
	})
}

// DeleteFile handles file deletion
// @Summary Delete uploaded file
// @Description Delete a file from storage
// @Tags Employee Onboarding
// @Accept json
// @Produce json
// @Param request body map[string]string true "File URL"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/delete-file [post]
func (h *FileUploadHandler) DeleteFile(c *gin.Context) {
	var req struct {
		FileURL string `json:"file_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", "file_url is required")
		return
	}

	// Delete file
	if err := h.storageService.DeleteFile(req.FileURL); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "File deleted successfully", nil)
}

// DeleteDocument handles document deletion by ID
// @Summary Delete uploaded document
// @Description Delete a document by ID (removes from database and storage)
// @Tags Employee Onboarding
// @Accept json
// @Produce json
// @Param request body map[string]uint true "Document ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/onboarding/delete-document [post]
func (h *FileUploadHandler) DeleteDocument(c *gin.Context) {
	var req struct {
		DocumentID uint `json:"document_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", "document_id is required")
		return
	}

	// Get tenant ID
	tenantID := middleware.GetTenantID(c)

	// Find document
	doc, err := h.documentRepo.FindByID(req.DocumentID)
	if err != nil {
		response.NotFound(c, "Document not found")
		return
	}

	// Verify tenant ownership through draft
	if doc.DraftID != nil {
		draft, err := h.draftRepo.FindByID(*doc.DraftID)
		if err == nil && draft != nil {
			if tenantID != nil && draft.TenantID != nil && *draft.TenantID != *tenantID {
				response.BadRequest(c, "Document does not belong to your tenant", nil)
				return
			}
		}
	}

	// Delete file from storage
	if doc.FileURL != "" {
		if err := h.storageService.DeleteFile(doc.FileURL); err != nil {
			// Log error but continue with database deletion
		}
	}

	// Delete document from database
	if err := h.documentRepo.Delete(req.DocumentID); err != nil {
		response.BadRequest(c, "Failed to delete document: "+err.Error(), nil)
		return
	}

	response.Success(c, "Document deleted successfully", gin.H{
		"document_id": req.DocumentID,
		"deleted_at":  time.Now(),
	})
}

// DocumentHandler handles employee document operations
type DocumentHandler struct {
	documentService *services.DocumentService
	storageService  *storage.StorageService
}

// NewDocumentHandler creates a new document handler
func NewDocumentHandler() *DocumentHandler {
	return &DocumentHandler{
		documentService: services.NewDocumentService(),
		storageService:  storage.NewStorageService(),
	}
}

// ListEmployeesWithDocuments lists all employees with their uploaded documents
func (h *DocumentHandler) ListEmployeesWithDocuments(c *gin.Context) {
	var req struct {
		Page     int `form:"page" binding:"omitempty,min=1"`
		PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	// Set defaults
	page := 1
	if req.Page > 0 {
		page = req.Page
	}
	pageSize := 20
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	tenantID := middleware.GetTenantID(c)
	employees, total, err := h.documentService.ListEmployeesWithDocuments(tenantID, page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Calculate pagination metadata
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response.SuccessWithMeta(c, "Employees with documents retrieved successfully", employees, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetEmployeeDocuments gets all documents for a specific employee
func (h *DocumentHandler) GetEmployeeDocuments(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		response.ValidationError(c, "Validation failed", "employee_id is required")
		return
	}

	tenantID := middleware.GetTenantID(c)
	documents, err := h.documentService.GetEmployeeDocuments(employeeID, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert to response format with full URLs
	docsResponse := make([]map[string]interface{}, len(documents))
	for i, doc := range documents {
		// Convert relative URL to full URL
		fullFileURL := h.getFullURL(c, doc.FileURL)
		docsResponse[i] = map[string]interface{}{
			"id":            doc.ID,
			"document_type": doc.DocumentType,
			"file_name":     doc.FileName,
			"file_url":      fullFileURL, // Return full URL
			"file_size":     doc.FileSize,
			"mime_type":     doc.MimeType,
			"description":   doc.Description,
			"uploaded_at":   doc.CreatedAt,
		}
	}

	response.Success(c, "Documents retrieved successfully", docsResponse)
}

// ViewDocument handles viewing/downloading a document
func (h *DocumentHandler) ViewDocument(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", "id is required in request body")
		return
	}

	tenantID := middleware.GetTenantID(c)
	doc, err := h.documentService.GetDocumentByID(req.ID, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert relative URL to full URL
	fullFileURL := h.getFullURL(c, doc.FileURL)

	// Return document information with full file URL
	response.Success(c, "Document retrieved successfully", map[string]interface{}{
		"id":            doc.ID,
		"employee_id":   doc.EmployeeIDString,
		"document_type": doc.DocumentType,
		"file_name":     doc.FileName,
		"file_url":      fullFileURL, // Return full URL
		"file_size":     doc.FileSize,
		"mime_type":     doc.MimeType,
		"description":   doc.Description,
		"uploaded_at":   doc.CreatedAt,
	})
}

// GetDocumentStatistics gets statistics about employee documents
func (h *DocumentHandler) GetDocumentStatistics(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	statistics, err := h.documentService.GetDocumentStatistics(tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Document statistics retrieved successfully", statistics)
}

// getFullURL converts a relative URL to a full URL using the request's scheme and host
func (h *DocumentHandler) getFullURL(c *gin.Context, relativeURL string) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	// Check X-Forwarded-Proto header for reverse proxy setups
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	host := c.Request.Host
	if host == "" {
		host = "localhost:8080" // Default fallback
	}

	// Remove leading slash from relativeURL if present (it will be added back)
	if len(relativeURL) > 0 && relativeURL[0] == '/' {
		relativeURL = relativeURL[1:]
	}

	return fmt.Sprintf("%s://%s/%s", scheme, host, relativeURL)
}

package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/bulk_import/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type BulkImportHandler struct{}

func NewBulkImportHandler() *BulkImportHandler {
	return &BulkImportHandler{}
}

// ── Template Downloads ─────────────────────────────────────────────────

// GET /bulk-import/templates/employees
func (h *BulkImportHandler) DownloadEmployeeTemplate(c *gin.Context) {
	f, err := services.GenerateEmployeeTemplate()
	if err != nil {
		response.InternalServerError(c, "Failed to generate template", err.Error())
		return
	}
	defer f.Close()

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=employee_bulk_upload_template.xlsx")
	if err := f.Write(c.Writer); err != nil {
		response.InternalServerError(c, "Failed to write template", err.Error())
	}
}

// GET /bulk-import/templates/departments
func (h *BulkImportHandler) DownloadDepartmentTemplate(c *gin.Context) {
	f, err := services.GenerateDepartmentTemplate()
	if err != nil {
		response.InternalServerError(c, "Failed to generate template", err.Error())
		return
	}
	defer f.Close()

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=department_bulk_upload_template.xlsx")
	if err := f.Write(c.Writer); err != nil {
		response.InternalServerError(c, "Failed to write template", err.Error())
	}
}

// GET /bulk-import/templates/positions
func (h *BulkImportHandler) DownloadPositionTemplate(c *gin.Context) {
	f, err := services.GeneratePositionTemplate()
	if err != nil {
		response.InternalServerError(c, "Failed to generate template", err.Error())
		return
	}
	defer f.Close()

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=position_bulk_upload_template.xlsx")
	if err := f.Write(c.Writer); err != nil {
		response.InternalServerError(c, "Failed to write template", err.Error())
	}
}

// ── Bulk Imports ───────────────────────────────────────────────────────

// POST /bulk-import/employees
func (h *BulkImportHandler) ImportEmployees(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Excel file is required (multipart field 'file')", nil)
		return
	}

	if err := validateExcelFile(file.Filename, file.Size); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := services.ImportEmployees(file)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	status := http.StatusOK
	msg := fmt.Sprintf("Employee import completed: %d imported, %d skipped out of %d rows", result.Imported, result.Skipped, result.TotalRows)
	if result.Imported == 0 && result.TotalRows > 0 {
		status = http.StatusUnprocessableEntity
		msg = "No employees imported — check errors"
	}

	c.JSON(status, response.APIResponse{
		Success: result.Imported > 0 || result.TotalRows == 0,
		Message: msg,
		Data:    result,
	})
}

// POST /bulk-import/departments
func (h *BulkImportHandler) ImportDepartments(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Excel file is required (multipart field 'file')", nil)
		return
	}

	if err := validateExcelFile(file.Filename, file.Size); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := services.ImportDepartments(file)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	status := http.StatusOK
	msg := fmt.Sprintf("Department import completed: %d imported, %d skipped out of %d rows", result.Imported, result.Skipped, result.TotalRows)
	if result.Imported == 0 && result.TotalRows > 0 {
		status = http.StatusUnprocessableEntity
		msg = "No departments imported — check errors"
	}

	c.JSON(status, response.APIResponse{
		Success: result.Imported > 0 || result.TotalRows == 0,
		Message: msg,
		Data:    result,
	})
}

// POST /bulk-import/positions
func (h *BulkImportHandler) ImportPositions(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Excel file is required (multipart field 'file')", nil)
		return
	}

	if err := validateExcelFile(file.Filename, file.Size); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := services.ImportPositions(file)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	status := http.StatusOK
	msg := fmt.Sprintf("Position import completed: %d imported, %d skipped out of %d rows", result.Imported, result.Skipped, result.TotalRows)
	if result.Imported == 0 && result.TotalRows > 0 {
		status = http.StatusUnprocessableEntity
		msg = "No positions imported — check errors"
	}

	c.JSON(status, response.APIResponse{
		Success: result.Imported > 0 || result.TotalRows == 0,
		Message: msg,
		Data:    result,
	})
}

// ── Helpers ────────────────────────────────────────────────────────────

func validateExcelFile(filename string, size int64) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".xlsx" && ext != ".xls" {
		return fmt.Errorf("only Excel files (.xlsx, .xls) are accepted")
	}
	if size > 10<<20 {
		return fmt.Errorf("file too large (max 10 MB)")
	}
	return nil
}

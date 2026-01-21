package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// BioTimeHandler handles BioTime API HTTP requests
type BioTimeHandler struct {
	// Service will be created per request with tenant ID
}

// NewBioTimeHandler creates a new BioTime handler
func NewBioTimeHandler() *BioTimeHandler {
	return &BioTimeHandler{}
}

// getService creates a BioTime service instance with tenant context
func (h *BioTimeHandler) getService(c *gin.Context) *services.BioTimeService {
	tenantID := middleware.GetTenantID(c)
	return services.NewBioTimeService(tenantID)
}

// TestConnection handles testing the BioTime API connection
// @Summary Test BioTime connection
// @Description Test connection and authentication with BioTime biometric device API
// @Tags Biometric
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/biometric/biotime/test-connection [get]
func (h *BioTimeHandler) TestConnection(c *gin.Context) {
	service := h.getService(c)
	result, err := service.TestConnection()
	if err != nil {
		response.InternalServerError(c, "Failed to test BioTime connection", err.Error())
		return
	}

	if result.Success {
		response.Success(c, result.Message, result)
	} else {
		response.BadRequest(c, result.Message, result)
	}
}

// GetToken handles getting the BioTime authentication token
// @Summary Get BioTime token
// @Description Get the current authentication token for BioTime API (automatically managed)
// @Tags Biometric
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/biometric/biotime/token [get]
func (h *BioTimeHandler) GetToken(c *gin.Context) {
	service := h.getService(c)
	token, err := service.GetToken()
	if err != nil {
		response.InternalServerError(c, "Failed to get BioTime token", err.Error())
		return
	}

	response.Success(c, "Token retrieved successfully", map[string]interface{}{
		"token": token,
		"note": "Token is automatically managed and cached. You don't need to store this manually.",
	})
}

// GetTerminals handles fetching terminals from BioTime API
// @Summary Get BioTime terminals
// @Description Fetch the list of terminal devices from BioTime API
// @Tags Biometric
// @Produce json
// @Success 200 {object} response.APIResponse{data=services.TerminalsResponse}
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/biometric/biotime/terminals [get]
func (h *BioTimeHandler) GetTerminals(c *gin.Context) {
	service := h.getService(c)
	terminalsResp, err := service.GetTerminals()
	if err != nil {
		response.InternalServerError(c, "Failed to fetch terminals from BioTime", err.Error())
		return
	}

	response.Success(c, "Terminals retrieved successfully", terminalsResp)
}

// GetTransactions handles fetching transactions from BioTime API
// @Summary Get BioTime transactions
// @Description Fetch attendance/transaction records from BioTime API with optional filters
// @Tags Biometric
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param page_size query int false "Number of records per page"
// @Param emp_code query string false "Filter by employee code"
// @Param terminal_sn query string false "Filter by terminal serial number"
// @Param terminal_alias query string false "Filter by terminal alias"
// @Param start_time query string false "Start time filter (format: YYYY-MM-DD HH:MM:SS)"
// @Param end_time query string false "End time filter (format: YYYY-MM-DD HH:MM:SS)"
// @Success 200 {object} response.APIResponse{data=services.TransactionsResponse}
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/biometric/biotime/transactions [get]
func (h *BioTimeHandler) GetTransactions(c *gin.Context) {
	service := h.getService(c)

	// Parse query parameters
	params := &services.GetTransactionsParams{}

	// Page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			params.Page = &page
		}
	}

	// Page size
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			params.PageSize = &pageSize
		}
	}

	// Employee code
	if empCode := c.Query("emp_code"); empCode != "" {
		params.EmpCode = &empCode
	}

	// Terminal SN
	if terminalSN := c.Query("terminal_sn"); terminalSN != "" {
		params.TerminalSN = &terminalSN
	}

	// Terminal alias
	if terminalAlias := c.Query("terminal_alias"); terminalAlias != "" {
		params.TerminalAlias = &terminalAlias
	}

	// Start time
	if startTime := c.Query("start_time"); startTime != "" {
		params.StartTime = &startTime
	}

	// End time
	if endTime := c.Query("end_time"); endTime != "" {
		params.EndTime = &endTime
	}

	transactionsResp, err := service.GetTransactions(params)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch transactions from BioTime", err.Error())
		return
	}

	// Queue transactions for background processing (non-blocking)
	// This doesn't affect the API response - it happens in the background
	tenantID := middleware.GetTenantID(c)
	go func() {
		if transactionsResp != nil && len(transactionsResp.Data) > 0 {
			// Convert to interface slice for queue
			transactions := make([]interface{}, len(transactionsResp.Data))
			for i, txn := range transactionsResp.Data {
				transactions[i] = txn
			}
			
			// Publish to queue (if RabbitMQ is enabled)
			if err := service.QueueTransactions(tenantID, transactions); err != nil {
				// Log error but don't fail the request
				fmt.Printf("Failed to queue transactions: %v\n", err)
			}
		}
	}()

	response.Success(c, "Transactions retrieved successfully", transactionsResp)
}

// GetTransaction handles fetching a single transaction by ID from BioTime API
// @Summary Get BioTime transaction by ID
// @Description Fetch a single attendance/transaction record by ID from BioTime API
// @Tags Biometric
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} response.APIResponse{data=services.Transaction}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/biometric/biotime/transactions/{id} [get]
func (h *BioTimeHandler) GetTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	if transactionID == "" {
		response.BadRequest(c, "Transaction ID is required", nil)
		return
	}

	service := h.getService(c)
	transaction, err := service.GetTransactionByID(transactionID)
	if err != nil {
		if strings.Contains(err.Error(), "status 404") || strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "Transaction not found")
		} else {
			response.InternalServerError(c, "Failed to fetch transaction from BioTime", err.Error())
		}
		return
	}

	response.Success(c, "Transaction retrieved successfully", transaction)
}

// RefreshToken handles manually refreshing the BioTime authentication token
// @Summary Refresh BioTime token
// @Description Force a refresh of the BioTime authentication token (even if current token is still valid)
// @Tags Biometric
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/biometric/biotime/refresh-token [post]
func (h *BioTimeHandler) RefreshToken(c *gin.Context) {
	service := h.getService(c)
	token, err := service.RefreshToken()
	if err != nil {
		response.InternalServerError(c, "Failed to refresh BioTime token", err.Error())
		return
	}

	response.Success(c, "Token refreshed successfully", map[string]interface{}{
		"token": token,
		"note": "New token has been stored in database and will be used for future requests.",
	})
}

// BackfillTransactions handles backfilling historical transactions from BioTime
// @Summary Backfill historical BioTime transactions
// @Description Pulls all transactions from 2025-01-01 to now and queues them for database storage
// @Tags Biometric
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD). Default: 2025-01-01"
// @Param end_date query string false "End date (YYYY-MM-DD). Default: now"
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/biometric/biotime/backfill [post]
func (h *BioTimeHandler) BackfillTransactions(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	backfillService := services.NewBioTimeBackfillService(tenantID)

	// Parse date parameters
	startDateStr := c.DefaultQuery("start_date", "2025-01-01")
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.BadRequest(c, "Invalid start_date format. Use YYYY-MM-DD", nil)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.BadRequest(c, "Invalid end_date format. Use YYYY-MM-DD", nil)
		return
	}

	if startDate.After(endDate) {
		response.BadRequest(c, "start_date must be before end_date", nil)
		return
	}

	// Run backfill in background
	go func() {
		if err := backfillService.BackfillTransactions(startDate, endDate); err != nil {
			log.Printf("Backfill error: %v", err)
		}
	}()

	response.Success(c, "Backfill started. Transactions are being queued in the background. Check logs for progress.", map[string]interface{}{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"note":       "This process runs in the background. Monitor your application logs for progress.",
	})
}

// GetDailyAttendance handles getting daily attendance from database
// @Summary Get daily attendance
// @Description Get daily attendance records with check-in, check-out, and working hours from biotime_transactions table
// @Tags Biometric
// @Produce json
// @Param start_time query string true "Start time (YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)"
// @Param end_time query string true "End time (YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)"
// @Param emp_code query string false "Filter by employee code"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 50, max: 100)"
// @Success 200 {object} response.APIResponse{data=[]services.DailyAttendanceResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/biometric/attendance/daily [get]
func (h *BioTimeHandler) GetDailyAttendance(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	attendanceService := services.NewAttendanceService()

	// Parse query parameters
	params := services.GetDailyAttendanceParams{}

	// Start time (required) - supports both date and datetime
	startTime := c.Query("start_time")
	if startTime == "" {
		// Try legacy parameter name for backward compatibility
		if startDate := c.Query("start_date"); startDate != "" {
			startTime = startDate
		} else {
			response.BadRequest(c, "start_time is required (format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)", nil)
			return
		}
	}
	params.StartTime = startTime

	// End time (required) - supports both date and datetime
	endTime := c.Query("end_time")
	if endTime == "" {
		// Try legacy parameter name for backward compatibility
		if endDate := c.Query("end_date"); endDate != "" {
			endTime = endDate
		} else {
			response.BadRequest(c, "end_time is required (format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)", nil)
			return
		}
	}
	params.EndTime = endTime

	// Employee code (optional)
	if empCode := c.Query("emp_code"); empCode != "" {
		params.EmpCode = &empCode
	}

	// Page (optional, default: 1)
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		} else {
			params.Page = 1
		}
	} else {
		params.Page = 1
	}

	// Page size (optional, default: 50)
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			params.PageSize = pageSize
		} else {
			params.PageSize = 50
		}
	} else {
		params.PageSize = 50
	}

	// Get attendance data
	attendance, total, err := attendanceService.GetDailyAttendance(tenantID, params)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Calculate pagination metadata
	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	response.Success(c, "Daily attendance retrieved successfully", map[string]interface{}{
		"data":       attendance,
		"pagination": map[string]interface{}{
			"page":        params.Page,
			"page_size":   params.PageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

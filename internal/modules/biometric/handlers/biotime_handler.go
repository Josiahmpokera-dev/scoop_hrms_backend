package handlers

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type DailyAttendanceSummary struct {
	EmpCode       string  `json:"emp_code"`
	FirstName     string  `json:"first_name"`
	LastName      *string `json:"last_name,omitempty"`
	Department    string  `json:"department"`
	Position      *string `json:"position,omitempty"`
	Date          string  `json:"date"`
	CheckIn       string  `json:"checkin"`
	CheckOut      string  `json:"checkout"`
	WorkingHours  float64 `json:"working_hours"`
	Transactions  int     `json:"transactions"`
	TerminalSN    string  `json:"terminal_sn"`
	TerminalAlias *string `json:"terminal_alias,omitempty"`
}

type DailyTransactionsResponse struct {
	Count    int                      `json:"count"`
	Next     *string                  `json:"next"`
	Previous *string                  `json:"previous"`
	Message  string                   `json:"msg"`
	Code     int                      `json:"code"`
	Data     []DailyAttendanceSummary `json:"data"`
}

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
// @Router/biometric/biotime/test-connection [get]
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
// @Router/biometric/biotime/token [get]
func (h *BioTimeHandler) GetToken(c *gin.Context) {
	service := h.getService(c)
	token, err := service.GetToken()
	if err != nil {
		response.InternalServerError(c, "Failed to get BioTime token", err.Error())
		return
	}

	response.Success(c, "Token retrieved successfully", map[string]interface{}{
		"token": token,
		"note":  "Token is automatically managed and cached. You don't need to store this manually.",
	})
}

// GetTerminals handles fetching terminals from BioTime API
// @Summary Get BioTime terminals
// @Description Fetch the list of terminal devices from BioTime API
// @Tags Biometric
// @Produce json
// @Success 200 {object} response.APIResponse{data=services.TerminalsResponse}
// @Failure 500 {object} response.APIResponse
// @Router/biometric/biotime/terminals [get]
func (h *BioTimeHandler) GetTerminals(c *gin.Context) {
	service := h.getService(c)
	terminalsResp, err := service.GetTerminals()
	if err != nil {
		response.InternalServerError(c, "Failed to fetch terminals from BioTime", err.Error())
		return
	}

	response.Success(c, "Terminals retrieved successfully", terminalsResp)
}

// GetDeviceStatus handles fetching device status from BioTime API
// @Summary Get device status
// @Description Get the status of all BioTime terminal devices
// @Tags Biometric
// @Produce json
// @Success 200 {object} response.APIResponse{data=services.TerminalsResponse}
// @Failure 500 {object} response.APIResponse
// @Router/biometric/biotime/device-status [get]
func (h *BioTimeHandler) GetDeviceStatus(c *gin.Context) {
	service := h.getService(c)
	terminalsResp, err := service.GetTerminals()
	if err != nil {
		response.InternalServerError(c, "Failed to fetch device status from BioTime", err.Error())
		return
	}

	response.Success(c, "Device status retrieved successfully", terminalsResp)
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
// @Router/biometric/biotime/transactions [get]
func (h *BioTimeHandler) GetTransactions(c *gin.Context) {
	if strings.EqualFold(c.Query("view"), "daily") {
		tenantID := middleware.GetTenantID(c)
		attendanceService := services.NewAttendanceService()
		source := strings.ToLower(strings.TrimSpace(c.DefaultQuery("source", "db")))

		page := 1
		if pageStr := c.Query("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		pageSize := 20
		if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
				pageSize = ps
			}
		}

		now := time.Now()
		defaultStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
		defaultEnd := now.Format("2006-01-02 15:04:05")

		startTime := c.DefaultQuery("start_time", c.DefaultQuery("start_date", defaultStart))
		endTime := c.DefaultQuery("end_time", c.DefaultQuery("end_date", defaultEnd))

		params := services.GetDailyAttendanceParams{
			StartTime: startTime,
			EndTime:   endTime,
			Page:      page,
			PageSize:  pageSize,
		}
		if empCode := c.Query("emp_code"); empCode != "" {
			params.EmpCode = &empCode
		}

		var dailyRows []services.DailyAttendanceResponse
		var total int64
		var err error
		if source == "biotime" {
			dailyRows, total, err = attendanceService.GetDailyAttendanceFromBioTime(tenantID, params)
		} else {
			dailyRows, total, err = attendanceService.GetDailyAttendance(tenantID, params)
		}
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}

		daily := make([]DailyAttendanceSummary, 0, len(dailyRows))
		for _, row := range dailyRows {
			firstName, lastName := splitName(row.Name)
			checkIn := ""
			if row.CheckIn != nil {
				checkIn = *row.CheckIn
			}
			checkOut := ""
			if row.CheckOut != nil {
				checkOut = *row.CheckOut
			}

			workingHours := 0.0
			if row.WorkingHours != nil {
				var hh, mm int
				if _, scanErr := fmt.Sscanf(*row.WorkingHours, "%d:%d", &hh, &mm); scanErr == nil {
					workingHours = roundTo2DP(float64(hh) + float64(mm)/60.0)
				}
			}

			daily = append(daily, DailyAttendanceSummary{
				EmpCode:      row.EmpCode,
				FirstName:    firstName,
				LastName:     lastName,
				Department:   row.Department,
				Date:         row.Date,
				CheckIn:      checkIn,
				CheckOut:     checkOut,
				WorkingHours: workingHours,
				Transactions: row.PunchCount,
			})
		}

		totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
		if totalPages == 0 {
			totalPages = 1
		}

		dailyResp := &DailyTransactionsResponse{
			Count:    len(daily),
			Next:     nil,
			Previous: nil,
			Message:  "Success",
			Code:     0,
			Data:     daily,
		}
		response.Success(c, "Daily transactions retrieved successfully", map[string]interface{}{
			"count":      dailyResp.Count,
			"next":       dailyResp.Next,
			"previous":   dailyResp.Previous,
			"msg":        dailyResp.Message,
			"code":       dailyResp.Code,
			"data":       dailyResp.Data,
			"start_time": startTime,
			"end_time":   endTime,
			"source":     source,
			"pagination": map[string]interface{}{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": totalPages,
			},
		})
		return
	}

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

func splitName(fullName string) (string, *string) {
	parts := strings.Fields(strings.TrimSpace(fullName))
	if len(parts) == 0 {
		return "", nil
	}
	if len(parts) == 1 {
		return parts[0], nil
	}
	last := strings.Join(parts[1:], " ")
	return parts[0], &last
}

func buildDailyAttendanceSummaries(transactionsResp *services.TransactionsResponse) []DailyAttendanceSummary {
	if transactionsResp == nil || len(transactionsResp.Data) == 0 {
		return []DailyAttendanceSummary{}
	}

	type agg struct {
		first         services.Transaction
		date          string
		checkInTime   *time.Time
		checkOutTime  *time.Time
		checkInRaw    string
		checkOutRaw   string
		transactions  int
		terminalSN    string
		terminalAlias *string
	}

	byKey := map[string]*agg{}

	for _, txn := range transactionsResp.Data {
		t, err := parseBioTimePunchTime(txn.PunchTime)
		dateStr := ""
		if err == nil {
			dateStr = t.Format("2006-01-02")
		} else if len(txn.PunchTime) >= 10 {
			dateStr = txn.PunchTime[:10]
		}

		key := txn.EmpCode + "|" + dateStr
		existing, ok := byKey[key]
		if !ok {
			existing = &agg{
				first:         txn,
				date:          dateStr,
				transactions:  0,
				terminalSN:    txn.TerminalSN,
				terminalAlias: txn.TerminalAlias,
			}
			byKey[key] = existing
		}

		existing.transactions++

		if err != nil {
			continue
		}

		if existing.checkInTime == nil || t.Before(*existing.checkInTime) {
			existing.checkInTime = &t
			existing.checkInRaw = t.Format("2006-01-02 15:04:05")
		}
		if existing.checkOutTime == nil || t.After(*existing.checkOutTime) {
			existing.checkOutTime = &t
			existing.checkOutRaw = t.Format("2006-01-02 15:04:05")
		}
	}

	out := make([]DailyAttendanceSummary, 0, len(byKey))
	for _, a := range byKey {
		checkIn := a.checkInRaw
		checkOut := a.checkOutRaw
		workingHours := 0.0

		if a.checkInTime != nil && a.checkOutTime != nil && a.checkOutTime.After(*a.checkInTime) {
			workingHours = a.checkOutTime.Sub(*a.checkInTime).Hours()
		}

		out = append(out, DailyAttendanceSummary{
			EmpCode:       a.first.EmpCode,
			FirstName:     a.first.FirstName,
			LastName:      a.first.LastName,
			Department:    a.first.Department,
			Position:      a.first.Position,
			Date:          a.date,
			CheckIn:       checkIn,
			CheckOut:      checkOut,
			WorkingHours:  roundTo2DP(workingHours),
			Transactions:  a.transactions,
			TerminalSN:    a.terminalSN,
			TerminalAlias: a.terminalAlias,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Date != out[j].Date {
			return out[i].Date > out[j].Date
		}
		return out[i].CheckIn > out[j].CheckIn
	})

	return out
}

func parseBioTimePunchTime(timeStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.000000",
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

func roundTo2DP(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
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
// @Router/biometric/biotime/transactions/{id} [get]
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
// @Router/biometric/biotime/refresh-token [post]
func (h *BioTimeHandler) RefreshToken(c *gin.Context) {
	service := h.getService(c)
	token, err := service.RefreshToken()
	if err != nil {
		response.InternalServerError(c, "Failed to refresh BioTime token", err.Error())
		return
	}

	response.Success(c, "Token refreshed successfully", map[string]interface{}{
		"token": token,
		"note":  "New token has been stored in database and will be used for future requests.",
	})
}

// ManualSyncToDatabase pulls transactions from the BioTime API and inserts new rows into biotime_transactions (no RabbitMQ).
// Defaults when start_time and end_time are omitted: today 00:00 through now in BIOMETRIC_TIMEZONE / APP_TIMEZONE.
func (h *BioTimeHandler) ManualSyncToDatabase(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req struct {
		StartTime string `json:"start_time" form:"start_time"`
		EndTime   string `json:"end_time" form:"end_time"`
	}
	_ = c.ShouldBind(&req)

	start, end, err := services.ResolveManualSyncWindow(req.StartTime, req.EndTime)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	syncSvc := services.NewManualBioTimeSyncService(tenantID)
	result, err := syncSvc.SyncWindow(start, end)
	if err != nil {
		response.InternalServerError(c, "Manual biometric sync failed", err.Error())
		return
	}

	response.Success(c, "Biometric transactions synced to database", result)
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
// @Router/biometric/biotime/backfill [post]
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
// @Router/biometric/attendance/daily [get]
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

	// Source selection:
	// - source=biotime -> fetch direct from BioTime API (real-time)
	// - source=db      -> fetch from local database (synced data)
	// Default: biotime (to avoid empty results when DB is not yet synced)
	source := strings.ToLower(strings.TrimSpace(c.DefaultQuery("source", "biotime")))

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
	var attendance []services.DailyAttendanceResponse
	var total int64
	var err error
	if source == "db" {
		attendance, total, err = attendanceService.GetDailyAttendance(tenantID, params)
	} else {
		attendance, total, err = attendanceService.GetDailyAttendanceFromBioTime(tenantID, params)
	}
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
		"data": attendance,
		"pagination": map[string]interface{}{
			"page":        params.Page,
			"page_size":   params.PageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetExceptional handles getting late arrivals (exceptional cases)
// @Summary Get late arrivals (exceptional)
// @Description Get employees who checked in more than 30 minutes after 08:00 (after 08:30) from biotime_transactions table
// @Tags Biometric
// @Produce json
// @Param start_time query string true "Start time (YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)"
// @Param end_time query string true "End time (YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)"
// @Param emp_code query string false "Filter by employee code"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 50, max: 100)"
// @Success 200 {object} response.APIResponse{data=object{data=[]services.LateArrivalResponse,pagination=object{page=int,page_size=int,total=int,total_pages=int}}}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router/biometric/attendance/exceptional [get]
func (h *BioTimeHandler) GetExceptional(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	attendanceService := services.NewAttendanceService()

	params := services.GetLateArrivalsParams{}

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

	// Get late arrivals data
	lateArrivals, total, err := attendanceService.GetLateArrivals(tenantID, params)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Calculate pagination metadata
	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	response.Success(c, "Late arrivals retrieved successfully", map[string]interface{}{
		"data": lateArrivals,
		"pagination": map[string]interface{}{
			"page":        params.Page,
			"page_size":   params.PageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetAttendanceCalendar handles monthly attendance calendar data for one employee.
// @Summary Get monthly attendance calendar
// @Description Returns calendar-ready daily attendance events for an employee in a given month
// @Tags Biometric
// @Produce json
// @Param emp_code query string true "Employee code"
// @Param month query string true "Month in YYYY-MM format"
// @Success 200 {object} response.APIResponse{data=services.AttendanceCalendarResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /biometric/attendance/calendar [get]
func (h *BioTimeHandler) GetAttendanceCalendar(c *gin.Context) {
	empCode := strings.TrimSpace(c.Query("emp_code"))
	month := strings.TrimSpace(c.Query("month"))

	if empCode == "" {
		response.BadRequest(c, "emp_code is required", map[string]string{"emp_code": "required"})
		return
	}
	if month == "" {
		response.BadRequest(c, "month is required", map[string]string{"month": "required"})
		return
	}
	if _, err := time.Parse("2006-01", month); err != nil {
		response.BadRequest(c, "month must be in YYYY-MM format", map[string]string{"month": "invalid_format"})
		return
	}

	tenantID := middleware.GetTenantID(c)
	attendanceService := services.NewAttendanceService()

	data, err := attendanceService.GetAttendanceCalendar(tenantID, empCode, month)
	if err != nil {
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "employee not found") {
			c.JSON(404, response.APIResponse{
				Success: false,
				Message: "Employee not found",
				Error:   map[string]string{"emp_code": "not_found"},
			})
			return
		}
		if strings.Contains(errMsg, "month must be in yyyy-mm format") {
			response.BadRequest(c, "month must be in YYYY-MM format", map[string]string{"month": "invalid_format"})
			return
		}
		response.InternalServerError(c, "Failed to retrieve attendance calendar", err.Error())
		return
	}

	response.Success(c, "Attendance calendar retrieved successfully", data)
}

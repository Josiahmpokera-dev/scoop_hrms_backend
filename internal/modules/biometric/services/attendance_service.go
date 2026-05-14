package services

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	biometricRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
)

// formatBiometricWallClock renders clock time in BIOMETRIC_TIMEZONE / APP_TIMEZONE.
// Postgres/GORM typically returns instants in UTC; users expect the same wall clock as the device/office.
func formatBiometricWallClock(t time.Time) string {
	return t.In(getBiometricLocation()).Format("15:04:05")
}

// AttendanceService handles attendance-related business logic
type AttendanceService struct {
	transactionRepo *biometricRepos.BioTimeTransactionRepository
	employeeRepo    *employeeRepos.EmployeeRepository
}

// NewAttendanceService creates a new attendance service
func NewAttendanceService() *AttendanceService {
	return &AttendanceService{
		transactionRepo: biometricRepos.NewBioTimeTransactionRepository(),
		employeeRepo:    employeeRepos.NewEmployeeRepository(),
	}
}

type AttendanceCalendarEvent struct {
	Date         string  `json:"date"`
	Day          string  `json:"day"`
	Type         string  `json:"type"`
	Title        string  `json:"title"`
	StatusColor  string  `json:"status_color"`
	CheckIn      *string `json:"checkin"`
	CheckOut     *string `json:"checkout"`
	WorkingHours *string `json:"working_hours"`
	PunchCount   int     `json:"punch_count"`
}

type AttendanceCalendarResponse struct {
	EmpCode string                    `json:"emp_code"`
	Year    int                       `json:"year"`
	Month   int                       `json:"month"`
	Events  []AttendanceCalendarEvent `json:"events"`
}

// GetAttendanceCalendar returns calendar-ready attendance events for one employee and month.
func (s *AttendanceService) GetAttendanceCalendar(tenantID *uint, empCode, month string) (*AttendanceCalendarResponse, error) {
	empCode = strings.TrimSpace(empCode)
	if empCode == "" {
		return nil, fmt.Errorf("emp_code is required")
	}
	month = strings.TrimSpace(month)
	loc := getBiometricLocation()
	monthStart, err := time.ParseInLocation("2006-01", month, loc)
	if err != nil {
		return nil, fmt.Errorf("month must be in YYYY-MM format")
	}

	startTime := time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, loc)
	endTime := startTime.AddDate(0, 1, 0).Add(-time.Nanosecond)

	employee, err := s.employeeRepo.FindByEmployeeID(empCode)
	if err != nil {
		employee = nil
	}

	rows, err := s.transactionRepo.GetEmployeeDailyAttendanceForMonth(tenantID, empCode, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance calendar data: %w", err)
	}

	if len(rows) == 0 && config.AppConfig != nil && config.AppConfig.BioTime.Enabled {
		syncSvc := NewManualBioTimeSyncService(tenantID)
		if _, syncErr := syncSvc.SyncWindow(startTime, endTime); syncErr == nil {
			rows, err = s.transactionRepo.GetEmployeeDailyAttendanceForMonth(tenantID, empCode, startTime, endTime)
			if err != nil {
				return nil, fmt.Errorf("failed to get attendance calendar data: %w", err)
			}
		}
	}

	if employee == nil {
		hasTx, txErr := s.transactionRepo.HasTransactionsForEmpCode(tenantID, empCode, startTime, endTime)
		if txErr != nil {
			return nil, fmt.Errorf("failed to validate employee code: %w", txErr)
		}
		if !hasTx {
			return nil, fmt.Errorf("employee not found")
		}
	}

	byDate := make(map[string]biometricRepos.CalendarAttendanceRecord, len(rows))
	for _, r := range rows {
		var key string
		if r.CheckIn != nil {
			key = r.CheckIn.In(loc).Format("2006-01-02")
		} else {
			key = r.Date.In(loc).Format("2006-01-02")
		}
		byDate[key] = r
	}

	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	events := make([]AttendanceCalendarEvent, 0, 32)
	for day := startTime; day.Month() == startTime.Month(); day = day.AddDate(0, 0, 1) {
		dateKey := day.Format("2006-01-02")
		rec, ok := byDate[dateKey]

		var checkInStr *string
		var checkOutStr *string
		var workingHoursStr *string
		punchCount := 0

		if ok {
			punchCount = rec.PunchCount
			if rec.CheckIn != nil {
				s := formatBiometricWallClock(*rec.CheckIn)
				checkInStr = &s
			}
			// Checkout comes from explicit OUT timestamps or last punch when multiple rows exist (see repository SQL).
			if rec.CheckOut != nil && rec.CheckIn != nil && rec.CheckOut.After(*rec.CheckIn) {
				s := formatBiometricWallClock(*rec.CheckOut)
				checkOutStr = &s
			}
			if rec.CheckIn != nil && checkOutStr != nil && rec.CheckOut != nil {
				d := rec.CheckOut.Sub(*rec.CheckIn)
				if d > 0 {
					h := int(d.Hours())
					m := int(d.Minutes()) % 60
					sec := int(d.Seconds()) % 60
					ws := fmt.Sprintf("%02d:%02d:%02d", h, m, sec)
					workingHoursStr = &ws
				}
			}
		}

		eventType := "absent"
		title := "Absent"
		color := "#ef4444"

		dayLocal := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
		switch {
		case dayLocal.After(today):
			eventType = "future"
			title = "Future"
			color = "#d4d4d8"
		case day.Weekday() == time.Saturday || day.Weekday() == time.Sunday:
			eventType = "weekend"
			title = "Weekend"
			color = "#a1a1aa"
		case punchCount > 0 && workingHoursStr != nil:
			parts := strings.Split(*workingHoursStr, ":")
			hh := 0
			mm := 0
			if len(parts) >= 2 {
				fmt.Sscanf(parts[0], "%d", &hh)
				fmt.Sscanf(parts[1], "%d", &mm)
			}
			totalMinutes := hh*60 + mm
			if totalMinutes < 8*60 {
				eventType = "short_day"
				title = "Short Day"
				color = "#f59e0b"
			} else {
				eventType = "present"
				title = "Present"
				color = "#22c55e"
			}
		case punchCount > 0:
			eventType = "present"
			title = "Present"
			color = "#22c55e"
		}

		events = append(events, AttendanceCalendarEvent{
			Date:         dateKey,
			Day:          day.Weekday().String(),
			Type:         eventType,
			Title:        title,
			StatusColor:  color,
			CheckIn:      checkInStr,
			CheckOut:     checkOutStr,
			WorkingHours: workingHoursStr,
			PunchCount:   punchCount,
		})
	}

	return &AttendanceCalendarResponse{
		EmpCode: empCode,
		Year:    startTime.Year(),
		Month:   int(startTime.Month()),
		Events:  events,
	}, nil
}

// DailyAttendanceResponse represents the response for daily attendance
type DailyAttendanceResponse struct {
	Name         string  `json:"name"`          // Full name (FirstName + LastName)
	EmpCode      string  `json:"emp_code"`      // Employee code
	Department   string  `json:"department"`    // Department
	Date         string  `json:"date"`          // Date in YYYY-MM-DD format
	CheckIn      *string `json:"checkin"`       // Check-in time in HH:MM:SS format (nullable)
	CheckOut     *string `json:"checkout"`      // Check-out time in HH:MM:SS format (nullable)
	WorkingHours *string `json:"working_hours"` // Working hours in HH:MM format (nullable)
	PunchCount   int     `json:"punch_count"`   // Number of punches for the day
}

// GetDailyAttendanceParams represents query parameters for getting daily attendance
type GetDailyAttendanceParams struct {
	StartTime string  `json:"start_time"` // Format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS
	EndTime   string  `json:"end_time"`   // Format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS
	EmpCode   *string `json:"emp_code"`   // Optional employee code filter
	Page      int     `json:"page"`       // Page number (default: 1)
	PageSize  int     `json:"page_size"`  // Page size (default: 50)
}

// parseDateTime parses a datetime string supporting multiple formats
func parseDateTime(dateTimeStr string) (time.Time, error) {
	if t, err := parseBiometricDateTime(dateTimeStr); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid datetime format. Supported formats: YYYY-MM-DD, YYYY-MM-DD HH:MM:SS, YYYY-MM-DD HH:MM, or RFC3339")
}

// ParseAttendanceRange parses start_time/end_time (or start_date/end_date strings) using the same rules as GetDailyAttendance.
func ParseAttendanceRange(startTimeStr, endTimeStr string) (time.Time, time.Time, error) {
	startTime, err := parseDateTime(strings.TrimSpace(startTimeStr))
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start_time/start_date format: %w", err)
	}
	endTime, err := parseDateTime(strings.TrimSpace(endTimeStr))
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end_time/end_date format: %w", err)
	}
	trimEnd := strings.TrimSpace(endTimeStr)
	if len(trimEnd) == 10 && !strings.Contains(trimEnd, " ") {
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, endTime.Location())
	}
	if startTime.After(endTime) {
		return time.Time{}, time.Time{}, fmt.Errorf("start_time must be before or equal to end_time")
	}
	return startTime, endTime, nil
}

// GetDailyAttendance gets daily attendance records with check-in, check-out, and working hours
func (s *AttendanceService) GetDailyAttendance(tenantID *uint, params GetDailyAttendanceParams) ([]DailyAttendanceResponse, int64, error) {
	startTime, endTime, err := ParseAttendanceRange(params.StartTime, params.EndTime)
	if err != nil {
		return nil, 0, err
	}

	// Set default pagination
	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100 // Limit max page size
	}

	// Get records from repository
	records, total, err := s.transactionRepo.GetDailyAttendance(tenantID, startTime, endTime, params.EmpCode, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get daily attendance: %w", err)
	}

	// Convert to response format
	response := make([]DailyAttendanceResponse, len(records))
	for i, record := range records {
		// Build full name
		name := record.FirstName
		if record.LastName != "" {
			if name != "" {
				name += " " + record.LastName
			} else {
				name = record.LastName
			}
		}
		if name == "" {
			name = record.EmpCode // Fallback to emp_code if no name
		}

		// Format check-in / check-out in office timezone (DB stores instants, often as UTC).
		var checkInStr *string
		if record.CheckIn != nil {
			formatted := formatBiometricWallClock(*record.CheckIn)
			checkInStr = &formatted
		}

		// Format check-out: repository leaves it null for true single-punch days.
		var checkOutStr *string
		if record.CheckOut != nil && record.CheckIn != nil && record.CheckOut.After(*record.CheckIn) {
			formatted := formatBiometricWallClock(*record.CheckOut)
			checkOutStr = &formatted
		}

		// Calculate working hours
		var workingHoursStr *string
		if record.CheckIn != nil && checkOutStr != nil && record.CheckOut != nil {
			duration := record.CheckOut.Sub(*record.CheckIn)
			if duration > 0 {
				hours := int(duration.Hours())
				minutes := int(duration.Minutes()) % 60
				formatted := fmt.Sprintf("%02d:%02d", hours, minutes)
				workingHoursStr = &formatted
			}
		}

		// Workday date aligned with first punch in business TZ (matches checkin/checkout display).
		dateStr := record.Date.Format("2006-01-02")
		if record.CheckIn != nil {
			dateStr = record.CheckIn.In(getBiometricLocation()).Format("2006-01-02")
		}

		response[i] = DailyAttendanceResponse{
			Name:         name,
			EmpCode:      record.EmpCode,
			Department:   record.Department,
			Date:         dateStr,
			CheckIn:      checkInStr,
			CheckOut:     checkOutStr,
			WorkingHours: workingHoursStr,
			PunchCount:   record.PunchCount,
		}
	}

	return response, total, nil
}

// GetDailyAttendanceFromBioTime fetches transactions from BioTime API directly and aggregates them per employee per day.
// This is useful when the database sync hasn't run yet or when you want real-time data.
func (s *AttendanceService) GetDailyAttendanceFromBioTime(tenantID *uint, params GetDailyAttendanceParams) ([]DailyAttendanceResponse, int64, error) {
	startTime, endTime, err := ParseAttendanceRange(params.StartTime, params.EndTime)
	if err != nil {
		return nil, 0, err
	}

	// Pagination for final aggregated output
	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}

	bio := NewBioTimeService(tenantID)

	// Fetch ALL pages from BioTime API for the given time range, then paginate after aggregation.
	// This ensures a complete daily summary per employee.
	startStr := startTime.Format("2006-01-02 15:04:05")
	endStr := endTime.Format("2006-01-02 15:04:05")

	const maxPages = 200
	apiPage := 1
	apiPageSize := 100
	if apiPageSize < 1 {
		apiPageSize = 50
	}

	type agg struct {
		empCode    string
		firstName  string
		lastName   string
		department string
		date       string // YYYY-MM-DD
		checkIn    *time.Time
		checkOut   *time.Time
		punchCount int
	}

	aggs := map[string]*agg{}

	for apiPage <= maxPages {
		p := apiPage
		ps := apiPageSize
		st := startStr
		et := endStr
		reqParams := &GetTransactionsParams{
			Page:      &p,
			PageSize:  &ps,
			StartTime: &st,
			EndTime:   &et,
		}
		if params.EmpCode != nil && *params.EmpCode != "" {
			reqParams.EmpCode = params.EmpCode
		}

		resp, err := bio.GetTransactions(reqParams)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to fetch transactions from BioTime: %w", err)
		}
		if resp == nil || len(resp.Data) == 0 {
			break
		}

		for _, t := range resp.Data {
			pt, err := parsePunchTime(t.PunchTime)
			if err != nil || pt.IsZero() {
				continue
			}

			dateKey := pt.Format("2006-01-02")
			key := t.EmpCode + "|" + dateKey

			lastName := ""
			if t.LastName != nil {
				lastName = *t.LastName
			}

			a, ok := aggs[key]
			if !ok {
				a = &agg{
					empCode:    t.EmpCode,
					firstName:  t.FirstName,
					lastName:   lastName,
					department: t.Department,
					date:       dateKey,
				}
				aggs[key] = a
			}

			a.punchCount++
			if a.checkIn == nil || pt.Before(*a.checkIn) {
				tmp := pt
				a.checkIn = &tmp
			}
			if a.checkOut == nil || pt.After(*a.checkOut) {
				tmp := pt
				a.checkOut = &tmp
			}
		}

		if resp.Next == nil || *resp.Next == "" {
			break
		}
		apiPage++
	}

	// Convert map to slice
	all := make([]DailyAttendanceResponse, 0, len(aggs))
	for _, a := range aggs {
		name := strings.TrimSpace(strings.TrimSpace(a.firstName) + " " + strings.TrimSpace(a.lastName))
		if name == "" {
			name = a.empCode
		}

		var checkInStr *string
		var checkOutStr *string
		var workingHoursStr *string

		if a.checkIn != nil {
			s := formatBiometricWallClock(*a.checkIn)
			checkInStr = &s
		}
		if a.checkOut != nil && a.checkIn != nil && a.checkOut.After(*a.checkIn) {
			s := formatBiometricWallClock(*a.checkOut)
			checkOutStr = &s
		}
		if a.checkIn != nil && checkOutStr != nil && a.checkOut != nil {
			d := a.checkOut.Sub(*a.checkIn)
			if d > 0 {
				h := int(d.Hours())
				m := int(d.Minutes()) % 60
				s := fmt.Sprintf("%02d:%02d", h, m)
				workingHoursStr = &s
			}
		}

		all = append(all, DailyAttendanceResponse{
			Name:         name,
			EmpCode:      a.empCode,
			Department:   a.department,
			Date:         a.date,
			CheckIn:      checkInStr,
			CheckOut:     checkOutStr,
			WorkingHours: workingHoursStr,
			PunchCount:   a.punchCount,
		})
	}

	// Sort by date desc then emp_code asc
	sort.Slice(all, func(i, j int) bool {
		if all[i].Date != all[j].Date {
			return all[i].Date > all[j].Date
		}
		return all[i].EmpCode < all[j].EmpCode
	})

	total := int64(len(all))

	// Apply pagination
	start := (page - 1) * pageSize
	if start >= len(all) {
		return []DailyAttendanceResponse{}, total, nil
	}
	end := start + pageSize
	if end > len(all) {
		end = len(all)
	}

	return all[start:end], total, nil
}

// LateArrivalResponse represents the response for late arrivals
type LateArrivalResponse struct {
	Name        string  `json:"name"`         // Full name (FirstName + LastName)
	EmpCode     string  `json:"emp_code"`     // Employee code
	Department  string  `json:"department"`   // Department
	Date        string  `json:"date"`         // Date in YYYY-MM-DD format
	CheckIn     *string `json:"checkin"`      // Check-in time in HH:MM:SS format (nullable)
	MinutesLate int     `json:"minutes_late"` // Minutes late (after 08:30)
}

// GetLateArrivalsParams represents query parameters for getting late arrivals
type GetLateArrivalsParams struct {
	StartTime string  `json:"start_time"` // Format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS
	EndTime   string  `json:"end_time"`   // Format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS
	EmpCode   *string `json:"emp_code"`   // Optional employee code filter
	Page      int     `json:"page"`       // Page number (default: 1)
	PageSize  int     `json:"page_size"`  // Page size (default: 50)
}

// GetLateArrivals gets employees who checked in more than 30 minutes after 08:00 (after 08:30)
func (s *AttendanceService) GetLateArrivals(tenantID *uint, params GetLateArrivalsParams) ([]LateArrivalResponse, int64, error) {
	// Parse start time
	startTime, err := parseDateTime(params.StartTime)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid start_time format: %w", err)
	}

	// Parse end time
	endTime, err := parseDateTime(params.EndTime)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid end_time format: %w", err)
	}

	// If only date was provided (no time component), set end_time to end of day
	if len(params.EndTime) == 10 && !strings.Contains(params.EndTime, " ") {
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, endTime.Location())
	}

	if startTime.After(endTime) {
		return nil, 0, fmt.Errorf("start_time must be before or equal to end_time")
	}

	// Set default pagination
	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100 // Limit max page size
	}

	// Get records from repository
	records, total, err := s.transactionRepo.GetLateArrivals(tenantID, startTime, endTime, params.EmpCode, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get late arrivals: %w", err)
	}

	// Convert to response format
	response := make([]LateArrivalResponse, len(records))
	for i, record := range records {
		// Build full name
		name := record.FirstName
		if record.LastName != "" {
			if name != "" {
				name += " " + record.LastName
			} else {
				name = record.LastName
			}
		}
		if name == "" {
			name = record.EmpCode // Fallback to emp_code if no name
		}

		// Format check-in time in office timezone
		var checkInStr *string
		if record.CheckIn != nil {
			formatted := formatBiometricWallClock(*record.CheckIn)
			checkInStr = &formatted
		}

		// Format date (workday in business TZ when check-in exists)
		dateStr := record.Date.Format("2006-01-02")
		if record.CheckIn != nil {
			dateStr = record.CheckIn.In(getBiometricLocation()).Format("2006-01-02")
		}

		response[i] = LateArrivalResponse{
			Name:        name,
			EmpCode:     record.EmpCode,
			Department:  record.Department,
			Date:        dateStr,
			CheckIn:     checkInStr,
			MinutesLate: record.MinutesLate,
		}
	}

	return response, total, nil
}

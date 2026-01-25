package services

import (
	"fmt"
	"sort"
	"strings"
	"time"

	biometricRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/repositories"
)

// AttendanceService handles attendance-related business logic
type AttendanceService struct {
	transactionRepo *biometricRepos.BioTimeTransactionRepository
}

// NewAttendanceService creates a new attendance service
func NewAttendanceService() *AttendanceService {
	return &AttendanceService{
		transactionRepo: biometricRepos.NewBioTimeTransactionRepository(),
	}
}

// DailyAttendanceResponse represents the response for daily attendance
type DailyAttendanceResponse struct {
	Name        string  `json:"name"`         // Full name (FirstName + LastName)
	EmpCode     string  `json:"emp_code"`     // Employee code
	Date        string  `json:"date"`         // Date in YYYY-MM-DD format
	CheckIn     *string `json:"checkin"`      // Check-in time in HH:MM:SS format (nullable)
	CheckOut    *string `json:"checkout"`     // Check-out time in HH:MM:SS format (nullable)
	WorkingHours *string `json:"working_hours"` // Working hours in HH:MM format (nullable)
	PunchCount  int     `json:"punch_count"`  // Number of punches for the day
}

// GetDailyAttendanceParams represents query parameters for getting daily attendance
type GetDailyAttendanceParams struct {
	StartTime string  `json:"start_time"` // Format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS
	EndTime   string  `json:"end_time"`   // Format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS
	EmpCode   *string `json:"emp_code"`  // Optional employee code filter
	Page      int     `json:"page"`       // Page number (default: 1)
	PageSize  int     `json:"page_size"`  // Page size (default: 50)
}

// parseDateTime parses a datetime string supporting multiple formats
func parseDateTime(dateTimeStr string) (time.Time, error) {
	// Try parsing as full datetime first (YYYY-MM-DD HH:MM:SS)
	if t, err := time.Parse("2006-01-02 15:04:05", dateTimeStr); err == nil {
		return t, nil
	}
	// Try parsing as datetime with seconds (YYYY-MM-DD HH:MM:SS)
	if t, err := time.Parse("2006-01-02T15:04:05", dateTimeStr); err == nil {
		return t, nil
	}
	// Try parsing as datetime without seconds (YYYY-MM-DD HH:MM)
	if t, err := time.Parse("2006-01-02 15:04", dateTimeStr); err == nil {
		return t, nil
	}
	// Try parsing as date only (YYYY-MM-DD) - set to start of day
	if t, err := time.Parse("2006-01-02", dateTimeStr); err == nil {
		return t, nil
	}
	// Try parsing as date with timezone
	if t, err := time.Parse(time.RFC3339, dateTimeStr); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid datetime format. Supported formats: YYYY-MM-DD, YYYY-MM-DD HH:MM:SS, YYYY-MM-DD HH:MM, or RFC3339")
}

// GetDailyAttendance gets daily attendance records with check-in, check-out, and working hours
func (s *AttendanceService) GetDailyAttendance(tenantID *uint, params GetDailyAttendanceParams) ([]DailyAttendanceResponse, int64, error) {
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
	// Check if the original string was date-only format (length 10 and no space)
	if len(params.EndTime) == 10 && !strings.Contains(params.EndTime, " ") {
		// Date only provided, set to end of day
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

		// Format check-in time
		var checkInStr *string
		if record.CheckIn != nil {
			formatted := record.CheckIn.Format("15:04:05")
			checkInStr = &formatted
		}

		// Format check-out time
		var checkOutStr *string
		if record.CheckOut != nil {
			formatted := record.CheckOut.Format("15:04:05")
			checkOutStr = &formatted
		}

		// Calculate working hours
		var workingHoursStr *string
		if record.CheckIn != nil && record.CheckOut != nil {
			duration := record.CheckOut.Sub(*record.CheckIn)
			if duration > 0 {
				hours := int(duration.Hours())
				minutes := int(duration.Minutes()) % 60
				formatted := fmt.Sprintf("%02d:%02d", hours, minutes)
				workingHoursStr = &formatted
			}
		}

		// Format date
		dateStr := record.Date.Format("2006-01-02")

		response[i] = DailyAttendanceResponse{
			Name:         name,
			EmpCode:      record.EmpCode,
			Date:         dateStr,
			CheckIn:      checkInStr,
			CheckOut:     checkOutStr,
			WorkingHours: workingHoursStr,
			PunchCount:    record.PunchCount,
		}
	}

	return response, total, nil
}

// GetDailyAttendanceFromBioTime fetches transactions from BioTime API directly and aggregates them per employee per day.
// This is useful when the database sync hasn't run yet or when you want real-time data.
func (s *AttendanceService) GetDailyAttendanceFromBioTime(tenantID *uint, params GetDailyAttendanceParams) ([]DailyAttendanceResponse, int64, error) {
	// Parse start/end using same rules as DB method
	startTime, err := parseDateTime(params.StartTime)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid start_time format: %w", err)
	}
	endTime, err := parseDateTime(params.EndTime)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid end_time format: %w", err)
	}
	if len(params.EndTime) == 10 && !strings.Contains(params.EndTime, " ") {
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, endTime.Location())
	}
	if startTime.After(endTime) {
		return nil, 0, fmt.Errorf("start_time must be before or equal to end_time")
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
					empCode:   t.EmpCode,
					firstName: t.FirstName,
					lastName:  lastName,
					date:      dateKey,
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
			s := a.checkIn.Format("15:04:05")
			checkInStr = &s
		}
		if a.checkOut != nil {
			s := a.checkOut.Format("15:04:05")
			checkOutStr = &s
		}
		if a.checkIn != nil && a.checkOut != nil {
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
	Date        string  `json:"date"`          // Date in YYYY-MM-DD format
	CheckIn     *string `json:"checkin"`     // Check-in time in HH:MM:SS format (nullable)
	MinutesLate int     `json:"minutes_late"` // Minutes late (after 08:30)
}

// GetLateArrivalsParams represents query parameters for getting late arrivals
type GetLateArrivalsParams struct {
	StartTime string  `json:"start_time"` // Format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS
	EndTime   string  `json:"end_time"`  // Format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS
	EmpCode   *string `json:"emp_code"`   // Optional employee code filter
	Page      int     `json:"page"`      // Page number (default: 1)
	PageSize  int     `json:"page_size"` // Page size (default: 50)
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

		// Format check-in time
		var checkInStr *string
		if record.CheckIn != nil {
			formatted := record.CheckIn.Format("15:04:05")
			checkInStr = &formatted
		}

		// Format date
		dateStr := record.Date.Format("2006-01-02")

		response[i] = LateArrivalResponse{
			Name:        name,
			EmpCode:     record.EmpCode,
			Date:        dateStr,
			CheckIn:     checkInStr,
			MinutesLate: record.MinutesLate,
		}
	}

	return response, total, nil
}

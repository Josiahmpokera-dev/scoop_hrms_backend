package services

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/repositories"
)

type LeaveCalendarService struct {
	requestRepo  *repositories.LeaveRequestRepository
	holidayRepo  *repositories.HolidayRepository
}

func NewLeaveCalendarService() *LeaveCalendarService {
	return &LeaveCalendarService{
		requestRepo: repositories.NewLeaveRequestRepository(),
		holidayRepo: repositories.NewHolidayRepository(),
	}
}

// GetMonthlyCalendar gets leave calendar for a specific month
func (s *LeaveCalendarService) GetMonthlyCalendar(year, month int, tenantID *uint, filters map[string]interface{}) (map[string]interface{}, error) {
	// Create date range for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // Last day of month

	// Get all leave requests in this period
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")
	
	filters["from_date"] = startDateStr
	filters["to_date"] = endDateStr
	
	requests, _, err := s.requestRepo.List(tenantID, 1, 1000, filters) // Get all for the month
	if err != nil {
		return nil, err
	}

	// Get holidays
	holidays, err := s.holidayRepo.FindByDateRange(startDate, endDate, tenantID)
	if err != nil {
		return nil, err
	}

	// Organize leaves by date
	leavesByDate := make(map[string][]models.LeaveRequest)
	for _, req := range requests {
		// Only include approved or pending requests
		if req.Status == "approved" || req.Status == "pending" {
			current := req.FromDate
			for !current.After(req.ToDate) {
				dateStr := current.Format("2006-01-02")
				leavesByDate[dateStr] = append(leavesByDate[dateStr], req)
				current = current.AddDate(0, 0, 1)
			}
		}
	}

	// Build calendar structure
	calendar := []map[string]interface{}{}
	current := startDate
	for !current.After(endDate) {
		dateStr := current.Format("2006-01-02")
		dayLeaves := leavesByDate[dateStr]
		
		dayEntry := map[string]interface{}{
			"date": dateStr,
			"leaves": dayLeaves,
		}
		calendar = append(calendar, dayEntry)
		current = current.AddDate(0, 0, 1)
	}

	// Get weekends
	weekends := []string{}
	current = startDate
	for !current.After(endDate) {
		weekday := current.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			weekends = append(weekends, current.Format("2006-01-02"))
		}
		current = current.AddDate(0, 0, 1)
	}

	// Format holidays
	holidayList := []map[string]interface{}{}
	for _, holiday := range holidays {
		holidayList = append(holidayList, map[string]interface{}{
			"date": holiday.Date.Format("2006-01-02"),
			"name": holiday.Name,
			"type": holiday.Type,
		})
	}

	return map[string]interface{}{
		"year":    year,
		"month":   month,
		"leaves":  calendar,
		"holidays": holidayList,
		"weekends": weekends,
	}, nil
}

// GetEmployeesOnLeaveToday gets employees on leave today
func (s *LeaveCalendarService) GetEmployeesOnLeaveToday(tenantID *uint, filters map[string]interface{}) ([]models.LeaveRequest, error) {
	today := time.Now()
	todayStr := today.Format("2006-01-02")
	
	filters["from_date"] = todayStr
	filters["to_date"] = todayStr
	
	requests, _, err := s.requestRepo.List(tenantID, 1, 1000, filters)
	if err != nil {
		return nil, err
	}

	// Filter to only approved requests that include today
	result := []models.LeaveRequest{}
	for _, req := range requests {
		if req.Status == "approved" {
			if (req.FromDate.Before(today) || req.FromDate.Equal(today)) &&
				(req.ToDate.After(today) || req.ToDate.Equal(today)) {
				result = append(result, req)
			}
		}
	}

	return result, nil
}

// GetWeeklyCalendar gets leave calendar for a specific week
func (s *LeaveCalendarService) GetWeeklyCalendar(year, week int, tenantID *uint, filters map[string]interface{}) (map[string]interface{}, error) {
	// Calculate week start date
	jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	weekStart := jan1.AddDate(0, 0, (week-1)*7)
	
	// Adjust to Monday (week starts on Monday)
	weekday := int(weekStart.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	weekStart = weekStart.AddDate(0, 0, -(weekday - 1))
	
	weekEnd := weekStart.AddDate(0, 0, 6)

	startDateStr := weekStart.Format("2006-01-02")
	endDateStr := weekEnd.Format("2006-01-02")
	
	filters["from_date"] = startDateStr
	filters["to_date"] = endDateStr
	
	requests, _, err := s.requestRepo.List(tenantID, 1, 1000, filters)
	if err != nil {
		return nil, err
	}

	// Organize by day
	dailyLeaves := []map[string]interface{}{}
	current := weekStart
	for !current.After(weekEnd) {
		dateStr := current.Format("2006-01-02")
		dayLeaves := []models.LeaveRequest{}
		
		for _, req := range requests {
			if req.Status == "approved" || req.Status == "pending" {
				if (req.FromDate.Before(current) || req.FromDate.Equal(current)) &&
					(req.ToDate.After(current) || req.ToDate.Equal(current)) {
					dayLeaves = append(dayLeaves, req)
				}
			}
		}
		
		dailyLeaves = append(dailyLeaves, map[string]interface{}{
			"date":   dateStr,
			"leaves": dayLeaves,
		})
		
		current = current.AddDate(0, 0, 1)
	}

	return map[string]interface{}{
		"year":        year,
		"week":        week,
		"week_start":  weekStart.Format("2006-01-02"),
		"week_end":    weekEnd.Format("2006-01-02"),
		"daily_leaves": dailyLeaves,
	}, nil
}

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	locationModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/models"
	locationRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/repositories"
)

type ShiftService struct {
	repo         *repositories.ShiftRepository
	locationRepo *locationRepos.LocationRepository
}

func NewShiftService() *ShiftService {
	return &ShiftService{
		repo:         repositories.NewShiftRepository(),
		locationRepo: locationRepos.NewLocationRepository(),
	}
}

// CreateShift creates a new shift
func (s *ShiftService) CreateShift(req *models.CreateShiftRequest, tenantID *uint, createdBy *uint) (*models.Shift, error) {
	// Validate shift code uniqueness
	if s.repo.ExistsByCode(req.ShiftCode) {
		return nil, errors.New("shift with this code already exists")
	}

	// Validate time format
	if _, err := time.Parse("15:04:05", req.StartTime); err != nil {
		return nil, errors.New("invalid start_time format. Use HH:MM:SS")
	}
	if _, err := time.Parse("15:04:05", req.EndTime); err != nil {
		return nil, errors.New("invalid end_time format. Use HH:MM:SS")
	}

	// Validate late threshold time format if provided
	if req.LateThresholdTime != "" {
		if _, err := time.Parse("15:04:05", req.LateThresholdTime); err != nil {
			return nil, errors.New("invalid late_threshold_time format. Use HH:MM:SS")
		}
	}

	// Calculate working hours
	startTime, _ := time.Parse("15:04:05", req.StartTime)
	endTime, _ := time.Parse("15:04:05", req.EndTime)
	workingHours := endTime.Sub(startTime).Hours()
	if workingHours < 0 {
		workingHours += 24 // Handle cross-day shifts
	}

	// Serialize weekly_off to JSON
	weeklyOffJSON, err := json.Marshal(req.WeeklyOff)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize weekly_off: %w", err)
	}

	shift := &models.Shift{
		ShiftName:         req.ShiftName,
		ShiftCode:         req.ShiftCode,
		ShiftType:         req.ShiftType,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
		LateThresholdTime: req.LateThresholdTime,
		WorkingHours:      workingHours,
		BreakDuration:     getIntValue(req.BreakDuration, 0),
		GraceMinutes:      getIntValue(req.GraceMinutes, 0),
		LateMarkAfter:     getIntValue(req.LateMarkAfter, 0),
		EarlyGoingMinutes: getIntValue(req.EarlyGoingMinutes, 0),
		HalfDayHours:      getFloatValue(req.HalfDayHours, workingHours/2),
		MinimumHours:      getFloatValue(req.MinimumHours, 4.0),
		CrossDay:          getBoolValue(req.CrossDay, false),
		NightShift:        getBoolValue(req.NightShift, false),
		WeeklyOff:         string(weeklyOffJSON),
		IsActive:          getBoolValue(req.IsActive, true),
		CreatedBy:         createdBy,
		UpdatedBy:         createdBy,
	}

	// Validate and associate locations
	if len(req.LocationIDs) > 0 {
		locations := make([]locationModels.Location, 0)
		for _, locID := range req.LocationIDs {
			location, err := s.locationRepo.FindByID(locID)
			if err != nil {
				return nil, fmt.Errorf("location with ID %d not found", locID)
			}
			locations = append(locations, *location)
		}
		shift.Locations = locations
	}

	if err := s.repo.Create(shift); err != nil {
		return nil, fmt.Errorf("failed to create shift: %w", err)
	}

	return s.repo.FindByID(shift.ID)
}

// GetShiftByID retrieves a shift by ID
func (s *ShiftService) GetShiftByID(id uint) (*models.Shift, error) {
	return s.repo.FindByID(id)
}

// UpdateShift updates an existing shift
func (s *ShiftService) UpdateShift(id uint, req *models.UpdateShiftRequest, updatedBy *uint) (*models.Shift, error) {
	shift, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("shift not found")
	}

	// Check code uniqueness if code is being updated
	if req.ShiftCode != nil && *req.ShiftCode != shift.ShiftCode {
		if s.repo.ExistsByCode(*req.ShiftCode) {
			return nil, errors.New("shift with this code already exists")
		}
		shift.ShiftCode = *req.ShiftCode
	}

	if req.ShiftName != nil {
		shift.ShiftName = *req.ShiftName
	}
	if req.ShiftType != nil {
		shift.ShiftType = *req.ShiftType
	}

	// Update times if provided
	if req.StartTime != nil {
		if _, err := time.Parse("15:04:05", *req.StartTime); err != nil {
			return nil, errors.New("invalid start_time format. Use HH:MM:SS")
		}
		shift.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		if _, err := time.Parse("15:04:05", *req.EndTime); err != nil {
			return nil, errors.New("invalid end_time format. Use HH:MM:SS")
		}
		shift.EndTime = *req.EndTime
	}

	// Update late threshold time if provided
	if req.LateThresholdTime != nil {
		if *req.LateThresholdTime != "" {
			if _, err := time.Parse("15:04:05", *req.LateThresholdTime); err != nil {
				return nil, errors.New("invalid late_threshold_time format. Use HH:MM:SS")
			}
		}
		shift.LateThresholdTime = *req.LateThresholdTime
	}

	// Recalculate working hours if times changed
	if req.StartTime != nil || req.EndTime != nil {
		startTime, _ := time.Parse("15:04:05", shift.StartTime)
		endTime, _ := time.Parse("15:04:05", shift.EndTime)
		workingHours := endTime.Sub(startTime).Hours()
		if workingHours < 0 {
			workingHours += 24
		}
		shift.WorkingHours = workingHours
	}

	if req.BreakDuration != nil {
		shift.BreakDuration = *req.BreakDuration
	}
	if req.GraceMinutes != nil {
		shift.GraceMinutes = *req.GraceMinutes
	}
	if req.LateMarkAfter != nil {
		shift.LateMarkAfter = *req.LateMarkAfter
	}
	if req.EarlyGoingMinutes != nil {
		shift.EarlyGoingMinutes = *req.EarlyGoingMinutes
	}
	if req.HalfDayHours != nil {
		shift.HalfDayHours = *req.HalfDayHours
	}
	if req.MinimumHours != nil {
		shift.MinimumHours = *req.MinimumHours
	}
	if req.CrossDay != nil {
		shift.CrossDay = *req.CrossDay
	}
	if req.NightShift != nil {
		shift.NightShift = *req.NightShift
	}
	if req.WeeklyOff != nil {
		weeklyOffJSON, err := json.Marshal(req.WeeklyOff)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize weekly_off: %w", err)
		}
		shift.WeeklyOff = string(weeklyOffJSON)
	}
	if req.IsActive != nil {
		shift.IsActive = *req.IsActive
	}

	// Update locations if provided
	if req.LocationIDs != nil {
		locations := make([]locationModels.Location, 0)
		for _, locID := range req.LocationIDs {
			location, err := s.locationRepo.FindByID(locID)
			if err != nil {
				return nil, fmt.Errorf("location with ID %d not found", locID)
			}
			locations = append(locations, *location)
		}
		shift.Locations = locations
	}

	shift.UpdatedBy = updatedBy

	if err := s.repo.Update(shift); err != nil {
		return nil, fmt.Errorf("failed to update shift: %w", err)
	}

	return s.repo.FindByID(id)
}

// DeleteShift deletes a shift
func (s *ShiftService) DeleteShift(id uint) error {
	// Check if shift has active rosters
	if s.repo.HasActiveRosters(id) {
		return errors.New("shift is assigned to active rosters and cannot be deleted")
	}

	return s.repo.Delete(id)
}

// ListShifts lists shifts with pagination and filters
func (s *ShiftService) ListShifts(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Shift, int64, error) {
	return s.repo.List(tenantID, page, pageSize, filters)
}

// DuplicateShift creates a copy of an existing shift
func (s *ShiftService) DuplicateShift(id uint, req *models.DuplicateShiftRequest, tenantID *uint, createdBy *uint) (*models.Shift, error) {
	original, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("shift not found")
	}

	// Generate new code if not provided
	newCode := original.ShiftCode + "-COPY"
	if req.ShiftCode != nil && *req.ShiftCode != "" {
		newCode = *req.ShiftCode
	}

	// Check code uniqueness
	if s.repo.ExistsByCode(newCode) {
		return nil, errors.New("shift with this code already exists")
	}

	// Generate new name if not provided
	newName := original.ShiftName + " (Copy)"
	if req.ShiftName != nil && *req.ShiftName != "" {
		newName = *req.ShiftName
	}

	// Create new shift with copied data
	createReq := &models.CreateShiftRequest{
		ShiftName:         newName,
		ShiftCode:         newCode,
		ShiftType:         original.ShiftType,
		StartTime:         original.StartTime,
		EndTime:           original.EndTime,
		LateThresholdTime: original.LateThresholdTime,
		BreakDuration:     &original.BreakDuration,
		GraceMinutes:      &original.GraceMinutes,
		LateMarkAfter:     &original.LateMarkAfter,
		EarlyGoingMinutes: &original.EarlyGoingMinutes,
		HalfDayHours:      &original.HalfDayHours,
		MinimumHours:      &original.MinimumHours,
		CrossDay:          &original.CrossDay,
		NightShift:        &original.NightShift,
		IsActive:          &original.IsActive,
	}

	// Parse weekly_off from JSON
	var weeklyOff []string
	if original.WeeklyOff != "" {
		json.Unmarshal([]byte(original.WeeklyOff), &weeklyOff)
	}
	createReq.WeeklyOff = weeklyOff

	// Copy location IDs
	locationIDs := make([]uint, len(original.Locations))
	for i, loc := range original.Locations {
		locationIDs[i] = loc.ID
	}
	createReq.LocationIDs = locationIDs

	return s.CreateShift(createReq, tenantID, createdBy)
}

// GetStatistics returns statistics for shifts and rosters
// Note: This is a placeholder - full statistics will be implemented in a combined service
func (s *ShiftService) GetStatistics(tenantID *uint) (map[string]interface{}, error) {
	// Get all shifts to calculate statistics
	shifts, _, err := s.repo.List(tenantID, 1, 1000, map[string]interface{}{"status": "all"})
	if err != nil {
		return nil, err
	}

	var totalShifts, activeShifts, inactiveShifts int64
	totalShifts = int64(len(shifts))
	for _, shift := range shifts {
		if shift.IsActive {
			activeShifts++
		} else {
			inactiveShifts++
		}
	}

	return map[string]interface{}{
		"total_shifts":    totalShifts,
		"active_shifts":   activeShifts,
		"inactive_shifts": inactiveShifts,
	}, nil
}

// Helper functions
func getIntValue(ptr *int, defaultValue int) int {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

func getFloatValue(ptr *float64, defaultValue float64) float64 {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

func getBoolValue(ptr *bool, defaultValue bool) bool {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

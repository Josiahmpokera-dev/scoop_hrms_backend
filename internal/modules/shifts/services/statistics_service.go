package services

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/repositories"
)

type StatisticsService struct {
	shiftRepo      *repositories.ShiftRepository
	rosterRepo     *repositories.RosterRepository
	swapRequestRepo *repositories.SwapRequestRepository
}

func NewStatisticsService() *StatisticsService {
	return &StatisticsService{
		shiftRepo:       repositories.NewShiftRepository(),
		rosterRepo:      repositories.NewRosterRepository(),
		swapRequestRepo: repositories.NewSwapRequestRepository(),
	}
}

// GetStatistics returns comprehensive statistics for shifts and rosters
func (s *StatisticsService) GetStatistics(tenantID *uint) (map[string]interface{}, error) {
	// Get shift statistics
	shifts, _, _ := s.shiftRepo.List(tenantID, 1, 1000, map[string]interface{}{"status": "all"})
	var totalShifts, activeShifts, inactiveShifts int64
	totalShifts = int64(len(shifts))
	for _, shift := range shifts {
		if shift.IsActive {
			activeShifts++
		} else {
			inactiveShifts++
		}
	}

	// Get roster statistics
	rosters, _, _ := s.rosterRepo.List(tenantID, 1, 1000, map[string]interface{}{})
	var totalRostered, publishedRosters, draftRosters int64
	totalRostered = int64(len(rosters))
	for _, roster := range rosters {
		if roster.Status == "published" {
			publishedRosters++
		} else if roster.Status == "draft" {
			draftRosters++
		}
	}

	// Get swap request statistics
	swapRequests, _, _ := s.swapRequestRepo.List(tenantID, 1, 1000, map[string]interface{}{})
	var pendingSwaps, approvedSwaps, rejectedSwaps int64
	for _, swap := range swapRequests {
		switch swap.Status {
		case "pending":
			pendingSwaps++
		case "approved":
			approvedSwaps++
		case "rejected":
			rejectedSwaps++
		}
	}

	return map[string]interface{}{
		"total_shifts":          totalShifts,
		"active_shifts":         activeShifts,
		"inactive_shifts":       inactiveShifts,
		"total_rostered":        totalRostered,
		"published_rosters":     publishedRosters,
		"draft_rosters":         draftRosters,
		"pending_swap_requests": pendingSwaps,
		"approved_swap_requests": approvedSwaps,
		"rejected_swap_requests": rejectedSwaps,
	}, nil
}

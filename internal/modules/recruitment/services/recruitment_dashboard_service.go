package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	deptRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/repositories"
)

const recruitmentDashboardCacheTTLSeconds = 60
const offersExpiringWithinDaysThreshold = 3
const recruitmentAnalyticsWindowMonths = 12

// RecruitmentDashboardService builds GET /recruitment/dashboard/summary.
type RecruitmentDashboardService struct {
	repo     *repositories.RecruitmentDashboardRepository
	deptRepo *deptRepos.DepartmentRepository
}

func NewRecruitmentDashboardService() *RecruitmentDashboardService {
	return &RecruitmentDashboardService{
		repo:     repositories.NewRecruitmentDashboardRepository(),
		deptRepo: deptRepos.NewDepartmentRepository(),
	}
}

func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}

func ptrFloat64(v float64) *float64 {
	x := round2(v)
	return &x
}

func dashboardDayBoundsUTC(t time.Time) (start, end time.Time) {
	u := t.UTC()
	start = time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
	end = start.Add(24 * time.Hour)
	return start, end
}

// GetSummary builds the consolidated dashboard document. departmentID filters by matching requisition.department to the department's name. asOf is optional RFC3339 or YYYY-MM-DD (UTC end-of-day for date-only).
func (s *RecruitmentDashboardService) GetSummary(departmentIDStr, asOfStr string) (*models.RecruitmentDashboardSummary, error) {
	asOf := time.Now().UTC()
	if strings.TrimSpace(asOfStr) != "" {
		if t, err := time.Parse(time.RFC3339, strings.TrimSpace(asOfStr)); err == nil {
			asOf = t.UTC()
		} else if t, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(asOfStr), time.UTC); err == nil {
			asOf = t.Add(24*time.Hour - time.Nanosecond).UTC()
		}
	}

	departmentName := ""
	if strings.TrimSpace(departmentIDStr) != "" {
		if id, err := strconv.ParseUint(strings.TrimSpace(departmentIDStr), 10, 32); err == nil && id > 0 {
			if d, err := s.deptRepo.FindByID(uint(id)); err == nil && d != nil {
				departmentName = strings.TrimSpace(d.Name)
			}
		}
	}

	windowEnd := asOf
	windowStart := asOf.AddDate(0, -recruitmentAnalyticsWindowMonths, 0)
	dayStart, _ := dashboardDayBoundsUTC(asOf)

	kpis := models.RecruitmentDashboardKpis{OffersExpiringWithinDaysThreshold: offersExpiringWithinDaysThreshold}
	var err error
	if kpis.ActiveOpeningsCount, err = s.repo.CountActiveOpenings(asOf, departmentName); err != nil {
		return nil, fmt.Errorf("active openings: %w", err)
	}
	if kpis.ActiveOpeningsHeadcountSum, err = s.repo.SumActiveOpeningsHeadcount(asOf, departmentName); err != nil {
		return nil, fmt.Errorf("headcount sum: %w", err)
	}
	if kpis.TotalCandidates, err = s.repo.CountTotalCandidates(); err != nil {
		return nil, fmt.Errorf("total candidates: %w", err)
	}
	if kpis.NewCandidatesLast7Days, err = s.repo.CountNewCandidatesLast7Days(asOf); err != nil {
		return nil, fmt.Errorf("new candidates 7d: %w", err)
	}
	if kpis.ScheduledInterviewsCount, err = s.repo.CountScheduledInterviewsFromDay(dayStart, departmentName); err != nil {
		return nil, fmt.Errorf("scheduled interviews: %w", err)
	}
	if kpis.ScheduledInterviewsToday, err = s.repo.CountScheduledInterviewsOnDay(asOf, departmentName); err != nil {
		return nil, fmt.Errorf("scheduled interviews today: %w", err)
	}
	if kpis.PendingOffersCount, err = s.repo.CountPendingOffers(asOf); err != nil {
		return nil, fmt.Errorf("pending offers: %w", err)
	}
	if kpis.OffersExpiringWithinDays, err = s.repo.CountOffersExpiringWithin(asOf, offersExpiringWithinDaysThreshold); err != nil {
		return nil, fmt.Errorf("offers expiring: %w", err)
	}

	rawPipeline, err := s.repo.PipelineStageCounts(departmentName)
	if err != nil {
		return nil, fmt.Errorf("pipeline: %w", err)
	}
	pipeline := models.RecruitmentDashboardPipeline{
		Applied:    rawPipeline[models.StageApplied],
		Screening:  rawPipeline[models.StageScreening],
		Assessment: 0,
		Interview:  rawPipeline[models.StageInterview] + rawPipeline[models.StageInterviewCompleted],
		Offer:      rawPipeline[models.StageOffer],
		Hired:      rawPipeline[models.StageHired],
	}

	qs := models.RecruitmentDashboardQuickStats{}
	if avgFill, err := s.repo.AvgTimeToFillDays(windowStart, windowEnd, departmentName); err == nil && avgFill != nil {
		qs.AvgTimeToFillDays = ptrFloat64(*avgFill)
	}
	if avgHire, err := s.repo.AvgTimeToHireDays(windowStart, windowEnd, departmentName); err == nil && avgHire != nil {
		qs.AvgTimeToHireDays = ptrFloat64(*avgHire)
	}
	if acc, rej, err := s.repo.OfferAcceptanceStats(windowStart, windowEnd); err == nil {
		den := acc + rej
		if den > 0 {
			qs.OfferAcceptanceRatePercent = ptrFloat64(float64(acc) / float64(den) * 100)
		}
	}
	if hired, total, err := s.repo.PipelineConversionStats(windowStart, windowEnd, departmentName); err == nil && total > 0 {
		qs.PipelineConversionRatePercent = ptrFloat64(float64(hired) / float64(total) * 100)
	}

	top, err := s.repo.TopOpeningsByApplicantCount(asOf, departmentName, 5)
	if err != nil {
		return nil, fmt.Errorf("top openings: %w", err)
	}
	upcoming, err := s.repo.ListUpcomingInterviewsDashboard(dayStart, departmentName, 5)
	if err != nil {
		return nil, fmt.Errorf("upcoming interviews: %w", err)
	}
	recentApps, err := s.repo.ListRecentApplications(departmentName, 5)
	if err != nil {
		return nil, fmt.Errorf("recent applications: %w", err)
	}
	recent := make([]models.RecruitmentDashboardRecentCandidate, 0, len(recentApps))
	for _, a := range recentApps {
		rc := models.RecruitmentDashboardRecentCandidate{
			CandidateID:                 a.CandidateID,
			PrimaryApplicationJobTitle:  "",
			PrimaryApplicationAppliedAt: a.AppliedDate,
			Stage:                       a.Stage,
		}
		if a.Candidate != nil {
			rc.FullName = strings.TrimSpace(a.Candidate.FirstName + " " + a.Candidate.LastName)
		}
		if a.JobOpening.ID != "" {
			rc.PrimaryApplicationJobTitle = a.JobOpening.JobTitle
		}
		if a.Score > 0 {
			s := a.Score
			rc.Score = &s
		}
		recent = append(recent, rc)
	}

	return &models.RecruitmentDashboardSummary{
		Kpis:               kpis,
		Pipeline:           pipeline,
		QuickStats:         qs,
		TopOpenings:        top,
		UpcomingInterviews: upcoming,
		RecentCandidates:   recent,
		Meta: models.RecruitmentDashboardMeta{
			GeneratedAt:     time.Now().UTC(),
			CacheTtlSeconds: recruitmentDashboardCacheTTLSeconds,
		},
	}, nil
}

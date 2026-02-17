package services

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/repositories"
)

type ReportService struct {
	goalRepo      *repositories.GoalRepository
	appraisalRepo *repositories.AppraisalRepository
	talentRepo    *repositories.TalentReviewRepository
}

func NewReportService() *ReportService {
	return &ReportService{
		goalRepo:      repositories.NewGoalRepository(),
		appraisalRepo: repositories.NewAppraisalRepository(),
		talentRepo:    repositories.NewTalentReviewRepository(),
	}
}

func (s *ReportService) GetSummary(department string) (map[string]interface{}, error) {
	goalCounts, _ := s.goalRepo.CountByStatus(nil)
	appraisalCounts, _ := s.appraisalRepo.CountAppraisalsByStatus(nil)
	appraisals, totalAppraisals, _ := s.appraisalRepo.ListAppraisals(nil, nil, "", department, false, 1, 10000)
	var ratingSum float64
	var ratingCount int
	for _, a := range appraisals {
		if a.OverallRating != nil {
			ratingSum += *a.OverallRating
			ratingCount++
		}
	}
	avgRating := 0.0
	if ratingCount > 0 {
		avgRating = ratingSum / float64(ratingCount)
	}
	return map[string]interface{}{
		"goals":            goalCounts,
		"appraisals":       appraisalCounts,
		"total_appraisals": totalAppraisals,
		"department":       department,
		"avg_rating":       avgRating,
		"rating_count":     ratingCount,
	}, nil
}

func (s *ReportService) GetRatingDistribution(department string) ([]map[string]interface{}, error) {
	appraisals, _, err := s.appraisalRepo.ListAppraisals(nil, nil, "Completed", department, false, 1, 10000)
	if err != nil {
		return nil, err
	}
	buckets := map[string]int64{"1-2": 0, "2-3": 0, "3-4": 0, "4-5": 0}
	for _, a := range appraisals {
		if a.OverallRating == nil {
			continue
		}
		r := *a.OverallRating
		switch {
		case r < 2:
			buckets["1-2"]++
		case r < 3:
			buckets["2-3"]++
		case r < 4:
			buckets["3-4"]++
		default:
			buckets["4-5"]++
		}
	}
	out := make([]map[string]interface{}, 0, len(buckets))
	for k, v := range buckets {
		out = append(out, map[string]interface{}{"range": k, "count": v})
	}
	return out, nil
}

func (s *ReportService) GetDepartmentSummary() ([]map[string]interface{}, error) {
	appraisals, _, err := s.appraisalRepo.ListAppraisals(nil, nil, "", "", false, 1, 10000)
	if err != nil {
		return nil, err
	}
	deptMap := make(map[string]*struct {
		Count int
		Sum   float64
	})
	for _, a := range appraisals {
		dept := ""
		if a.Department != nil {
			dept = *a.Department
		}
		if deptMap[dept] == nil {
			deptMap[dept] = &struct {
				Count int
				Sum   float64
			}{}
		}
		deptMap[dept].Count++
		if a.OverallRating != nil {
			deptMap[dept].Sum += *a.OverallRating
		}
	}
	out := make([]map[string]interface{}, 0, len(deptMap))
	for dept, v := range deptMap {
		avg := 0.0
		if v.Count > 0 {
			avg = v.Sum / float64(v.Count)
		}
		out = append(out, map[string]interface{}{
			"department":   dept,
			"count":       v.Count,
			"avg_rating":  avg,
		})
	}
	return out, nil
}

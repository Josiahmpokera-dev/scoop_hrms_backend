package models

import "time"

// RecruitmentDashboardSummary is the consolidated payload for GET /recruitment/dashboard/summary.
type RecruitmentDashboardSummary struct {
	Kpis            RecruitmentDashboardKpis            `json:"kpis"`
	Pipeline        RecruitmentDashboardPipeline        `json:"pipeline"`
	QuickStats      RecruitmentDashboardQuickStats      `json:"quickStats"`
	TopOpenings     []RecruitmentDashboardTopOpening    `json:"topOpenings"`
	UpcomingInterviews []RecruitmentDashboardUpcomingInterview `json:"upcomingInterviews"`
	RecentCandidates   []RecruitmentDashboardRecentCandidate   `json:"recentCandidates"`
	Meta            RecruitmentDashboardMeta            `json:"meta"`
}

// RecruitmentDashboardKpis maps to dashboard KPI cards.
type RecruitmentDashboardKpis struct {
	ActiveOpeningsCount                 int64 `json:"activeOpeningsCount"`
	ActiveOpeningsHeadcountSum          int64 `json:"activeOpeningsHeadcountSum"`
	TotalCandidates                     int64 `json:"totalCandidates"`
	NewCandidatesLast7Days              int64 `json:"newCandidatesLast7Days"`
	ScheduledInterviewsCount            int64 `json:"scheduledInterviewsCount"`
	ScheduledInterviewsToday            int64 `json:"scheduledInterviewsToday"`
	PendingOffersCount                  int64 `json:"pendingOffersCount"`
	OffersExpiringWithinDays            int64 `json:"offersExpiringWithinDays"`
	OffersExpiringWithinDaysThreshold  int   `json:"offersExpiringWithinDaysThreshold"`
}

// RecruitmentDashboardPipeline uses UI-aligned keys; assessment is reserved (no backend stage yet).
type RecruitmentDashboardPipeline struct {
	Applied     int64 `json:"applied"`
	Screening   int64 `json:"screening"`
	Assessment  int64 `json:"assessment"`
	Interview   int64 `json:"interview"`
	Offer       int64 `json:"offer"`
	Hired       int64 `json:"hired"`
}

// RecruitmentDashboardQuickStats holds analytic metrics (nullable when not computable).
type RecruitmentDashboardQuickStats struct {
	AvgTimeToFillDays              *float64 `json:"avgTimeToFillDays"`
	AvgTimeToHireDays              *float64 `json:"avgTimeToHireDays"`
	OfferAcceptanceRatePercent     *float64 `json:"offerAcceptanceRatePercent"`
	PipelineConversionRatePercent  *float64 `json:"pipelineConversionRatePercent"`
	TopSourceLabel                 *string  `json:"topSourceLabel"`
	TopSourceSharePercent          *float64 `json:"topSourceSharePercent"`
	ReferralSuccessRatePercent     *float64 `json:"referralSuccessRatePercent"`
}

// RecruitmentDashboardMeta documents when the snapshot was built.
type RecruitmentDashboardMeta struct {
	GeneratedAt      time.Time `json:"generatedAt"`
	CacheTtlSeconds  int       `json:"cacheTtlSeconds"`
}

// RecruitmentDashboardTopOpening is one row for "top performing jobs".
type RecruitmentDashboardTopOpening struct {
	JobOpeningID    string `json:"jobOpeningId"`
	JobTitle        string `json:"jobTitle"`
	Department      string `json:"department"`
	Location        string `json:"location"`
	TotalApplicants int64  `json:"totalApplicants"`
}

// RecruitmentDashboardUpcomingInterview is a list row for the dashboard widget.
type RecruitmentDashboardUpcomingInterview struct {
	InterviewID   string `json:"interviewId"`
	CandidateID   string `json:"candidateId"`
	CandidateName string `json:"candidateName"`
	JobOpeningID  string `json:"jobOpeningId"`
	JobTitle      string `json:"jobTitle"`
	InterviewType string `json:"interviewType"`
	ScheduledDate string `json:"scheduledDate"`
	ScheduledTime string `json:"scheduledTime"`
	Status        string `json:"status"`
}

// RecruitmentDashboardRecentCandidate is a list row for recent activity.
type RecruitmentDashboardRecentCandidate struct {
	CandidateID                   string     `json:"candidateId"`
	FullName                      string     `json:"fullName"`
	CurrentDesignation            string     `json:"currentDesignation,omitempty"`
	Location                      string     `json:"location,omitempty"`
	PrimaryApplicationJobTitle    string     `json:"primaryApplicationJobTitle"`
	PrimaryApplicationAppliedAt   time.Time  `json:"primaryApplicationAppliedAt"`
	Stage                         string     `json:"stage"`
	Score                         *float64   `json:"score,omitempty"`
}

package repositories

import (
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	"gorm.io/gorm"
)

// RecruitmentDashboardRepository runs aggregate queries for the recruitment dashboard.
type RecruitmentDashboardRepository struct {
	db *gorm.DB
}

func NewRecruitmentDashboardRepository() *RecruitmentDashboardRepository {
	return &RecruitmentDashboardRepository{db: database.GetDB()}
}

func dayBoundsUTC(t time.Time) (start, end time.Time) {
	u := t.UTC()
	start = time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
	end = start.Add(24 * time.Hour)
	return start, end
}

func (r *RecruitmentDashboardRepository) CountActiveOpenings(asOf time.Time, departmentName string) (int64, error) {
	var n int64
	q := r.db.Model(&models.JobOpening{}).Where("job_openings.deleted_at IS NULL").
		Where("job_openings.status IN ?", []string{models.JobOpeningStatusActive, models.LegacyStatusPublished}).
		Where("(job_openings.expiry_date IS NULL OR job_openings.expiry_date = ? OR job_openings.expiry_date > ?)", time.Time{}, asOf)
	if strings.TrimSpace(departmentName) != "" {
		q = q.Joins("JOIN job_requisitions jr ON jr.id = job_openings.requisition_id AND jr.deleted_at IS NULL AND LOWER(TRIM(jr.department)) = LOWER(TRIM(?))", strings.TrimSpace(departmentName))
	}
	err := q.Count(&n).Error
	return n, err
}

func (r *RecruitmentDashboardRepository) SumActiveOpeningsHeadcount(asOf time.Time, departmentName string) (int64, error) {
	var sum int64
	q := r.db.Table("job_openings").
		Select("COALESCE(SUM(jr.headcount), 0)").
		Joins("JOIN job_requisitions jr ON jr.id = job_openings.requisition_id AND jr.deleted_at IS NULL").
		Where("job_openings.deleted_at IS NULL").
		Where("job_openings.status IN ?", []string{models.JobOpeningStatusActive, models.LegacyStatusPublished}).
		Where("(job_openings.expiry_date IS NULL OR job_openings.expiry_date = ? OR job_openings.expiry_date > ?)", time.Time{}, asOf)
	if strings.TrimSpace(departmentName) != "" {
		q = q.Where("LOWER(TRIM(jr.department)) = LOWER(TRIM(?))", strings.TrimSpace(departmentName))
	}
	err := q.Scan(&sum).Error
	return sum, err
}

func (r *RecruitmentDashboardRepository) CountTotalCandidates() (int64, error) {
	var n int64
	err := r.db.Model(&models.Candidate{}).Count(&n).Error
	return n, err
}

func (r *RecruitmentDashboardRepository) CountNewCandidatesLast7Days(asOf time.Time) (int64, error) {
	var n int64
	since := asOf.Add(-7 * 24 * time.Hour)
	err := r.db.Model(&models.Candidate{}).
		Where("created_at >= ? AND created_at <= ?", since, asOf).
		Count(&n).Error
	return n, err
}

func (r *RecruitmentDashboardRepository) CountScheduledInterviewsFromDay(dayStart time.Time, departmentName string) (int64, error) {
	var n int64
	q := r.db.Model(&models.Interview{}).
		Where("interviews.deleted_at IS NULL AND interviews.status = ?", models.InterviewStatusScheduled).
		Where("interviews.scheduled_date >= ?", dayStart.UTC())
	if strings.TrimSpace(departmentName) != "" {
		q = q.Joins("JOIN job_openings jo ON jo.id = interviews.job_opening_id AND jo.deleted_at IS NULL").
			Joins("JOIN job_requisitions jr ON jr.id = jo.requisition_id AND jr.deleted_at IS NULL AND LOWER(TRIM(jr.department)) = LOWER(TRIM(?))", strings.TrimSpace(departmentName))
	}
	err := q.Count(&n).Error
	return n, err
}

func (r *RecruitmentDashboardRepository) CountScheduledInterviewsOnDay(asOf time.Time, departmentName string) (int64, error) {
	dayStart, dayEnd := dayBoundsUTC(asOf)
	var n int64
	q := r.db.Model(&models.Interview{}).
		Where("interviews.deleted_at IS NULL AND interviews.status = ?", models.InterviewStatusScheduled).
		Where("interviews.scheduled_date >= ? AND interviews.scheduled_date < ?", dayStart, dayEnd)
	if strings.TrimSpace(departmentName) != "" {
		q = q.Joins("JOIN job_openings jo ON jo.id = interviews.job_opening_id AND jo.deleted_at IS NULL").
			Joins("JOIN job_requisitions jr ON jr.id = jo.requisition_id AND jr.deleted_at IS NULL AND LOWER(TRIM(jr.department)) = LOWER(TRIM(?))", strings.TrimSpace(departmentName))
	}
	err := q.Count(&n).Error
	return n, err
}

func (r *RecruitmentDashboardRepository) CountPendingOffers(asOf time.Time) (int64, error) {
	var n int64
	err := r.db.Model(&models.Offer{}).
		Where("deleted_at IS NULL").
		Where("status IN ?", []string{"Draft", "Sent"}).
		Where("created_at <= ?", asOf).
		Count(&n).Error
	return n, err
}

func (r *RecruitmentDashboardRepository) CountOffersExpiringWithin(asOf time.Time, thresholdDays int) (int64, error) {
	var n int64
	until := asOf.Add(time.Duration(thresholdDays) * 24 * time.Hour)
	err := r.db.Model(&models.Offer{}).
		Where("deleted_at IS NULL").
		Where("status IN ?", []string{"Draft", "Sent"}).
		Where("expiry_date > ? AND expiry_date <= ?", asOf, until).
		Count(&n).Error
	return n, err
}

// PipelineStageCounts returns raw counts per application stage for non-terminal funnel rows.
func (r *RecruitmentDashboardRepository) PipelineStageCounts(departmentName string) (map[string]int64, error) {
	out := make(map[string]int64)
	type row struct {
		Stage string
		Cnt   int64
	}
	var rows []row
	q := r.db.Model(&models.JobApplication{}).
		Select("job_applications.stage, COUNT(*) as cnt").
		Where("job_applications.deleted_at IS NULL").
		Where("job_applications.stage NOT IN ?", []string{models.StageRejected, models.StageTalentPool})
	if strings.TrimSpace(departmentName) != "" {
		q = q.Joins("JOIN job_openings jo ON jo.id = job_applications.job_opening_id AND jo.deleted_at IS NULL").
			Joins("JOIN job_requisitions jr ON jr.id = jo.requisition_id AND jr.deleted_at IS NULL AND LOWER(TRIM(jr.department)) = LOWER(TRIM(?))", strings.TrimSpace(departmentName))
	}
	err := q.Group("job_applications.stage").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, rw := range rows {
		out[rw.Stage] = rw.Cnt
	}
	return out, nil
}

// AvgTimeToFillDays mean days from job opening created_at to first hire on that opening (hires in analytics window).
func (r *RecruitmentDashboardRepository) AvgTimeToFillDays(windowStart, windowEnd time.Time, departmentName string) (*float64, error) {
	var avg *float64
	sql := `
SELECT AVG(EXTRACT(EPOCH FROM (h.hired_at - o.created_at)) / 86400.0)::float AS avg_days
FROM job_openings o
JOIN (
  SELECT job_opening_id, MIN(COALESCE(last_status_date, updated_at)) AS hired_at
  FROM job_applications
  WHERE deleted_at IS NULL AND stage = 'Hired'
    AND COALESCE(last_status_date, updated_at) >= ? AND COALESCE(last_status_date, updated_at) <= ?
  GROUP BY job_opening_id
) h ON h.job_opening_id = o.id
WHERE o.deleted_at IS NULL`
	args := []interface{}{windowStart, windowEnd}
	if strings.TrimSpace(departmentName) != "" {
		sql += `
  AND EXISTS (
    SELECT 1 FROM job_requisitions jr
    WHERE jr.id = o.requisition_id AND jr.deleted_at IS NULL
      AND LOWER(TRIM(jr.department)) = LOWER(TRIM(?))
  )`
		args = append(args, strings.TrimSpace(departmentName))
	}
	err := r.db.Raw(sql, args...).Scan(&avg).Error
	if err != nil {
		return nil, err
	}
	return avg, nil
}

// AvgTimeToHireDays mean days from application applied_date to hire timestamp for hires in window.
func (r *RecruitmentDashboardRepository) AvgTimeToHireDays(windowStart, windowEnd time.Time, departmentName string) (*float64, error) {
	var avg *float64
	q := r.db.Model(&models.JobApplication{}).
		Select("AVG(EXTRACT(EPOCH FROM (COALESCE(last_status_date, updated_at) - applied_date)) / 86400.0)::float").
		Where("deleted_at IS NULL AND stage = ?", models.StageHired).
		Where("COALESCE(last_status_date, updated_at) >= ? AND COALESCE(last_status_date, updated_at) <= ?", windowStart, windowEnd)
	if strings.TrimSpace(departmentName) != "" {
		q = q.Joins("JOIN job_openings jo ON jo.id = job_applications.job_opening_id AND jo.deleted_at IS NULL").
			Joins("JOIN job_requisitions jr ON jr.id = jo.requisition_id AND jr.deleted_at IS NULL AND LOWER(TRIM(jr.department)) = LOWER(TRIM(?))", strings.TrimSpace(departmentName))
	}
	err := q.Scan(&avg).Error
	return avg, err
}

// OfferAcceptanceStats returns accepted and rejected offer counts in window (for acceptance rate).
func (r *RecruitmentDashboardRepository) OfferAcceptanceStats(windowStart, windowEnd time.Time) (accepted, rejected int64, err error) {
	type row struct {
		Status string
		Cnt    int64
	}
	var rows []row
	err2 := r.db.Model(&models.Offer{}).
		Select("status, COUNT(*) as cnt").
		Where("deleted_at IS NULL").
		Where("status IN ?", []string{"Accepted", "Rejected"}).
		Where("updated_at >= ? AND updated_at <= ?", windowStart, windowEnd).
		Group("status").Scan(&rows).Error
	if err2 != nil {
		return 0, 0, err2
	}
	for _, rw := range rows {
		switch rw.Status {
		case "Accepted":
			accepted = rw.Cnt
		case "Rejected":
			rejected = rw.Cnt
		}
	}
	return accepted, rejected, nil
}

// PipelineConversionStats returns hired applications and total applications created in window.
func (r *RecruitmentDashboardRepository) PipelineConversionStats(windowStart, windowEnd time.Time, departmentName string) (hired, totalApps int64, err error) {
	var h, t int64
	qh := r.db.Model(&models.JobApplication{}).
		Where("deleted_at IS NULL AND stage = ?", models.StageHired).
		Where("COALESCE(last_status_date, updated_at) >= ? AND COALESCE(last_status_date, updated_at) <= ?", windowStart, windowEnd)
	if strings.TrimSpace(departmentName) != "" {
		qh = qh.Joins("JOIN job_openings jo ON jo.id = job_applications.job_opening_id AND jo.deleted_at IS NULL").
			Joins("JOIN job_requisitions jr ON jr.id = jo.requisition_id AND jr.deleted_at IS NULL AND LOWER(TRIM(jr.department)) = LOWER(TRIM(?))", strings.TrimSpace(departmentName))
	}
	if err = qh.Count(&h).Error; err != nil {
		return 0, 0, err
	}
	qt := r.db.Model(&models.JobApplication{}).
		Where("deleted_at IS NULL").
		Where("created_at >= ? AND created_at <= ?", windowStart, windowEnd)
	if strings.TrimSpace(departmentName) != "" {
		qt = qt.Joins("JOIN job_openings jo ON jo.id = job_applications.job_opening_id AND jo.deleted_at IS NULL").
			Joins("JOIN job_requisitions jr ON jr.id = jo.requisition_id AND jr.deleted_at IS NULL AND LOWER(TRIM(jr.department)) = LOWER(TRIM(?))", strings.TrimSpace(departmentName))
	}
	if err = qt.Count(&t).Error; err != nil {
		return 0, 0, err
	}
	return h, t, nil
}

// TopOpeningsByApplicantCount returns top N active openings by application count.
func (r *RecruitmentDashboardRepository) TopOpeningsByApplicantCount(asOf time.Time, departmentName string, limit int) ([]models.RecruitmentDashboardTopOpening, error) {
	type row struct {
		JobOpeningID string
		JobTitle     string
		Department   string
		Location     string
		Cnt          int64
	}
	var rows []row
	q := r.db.Table("job_openings jo").
		Select("jo.id as job_opening_id, jo.job_title, jr.department, jr.location, COUNT(ja.id) as cnt").
		Joins("JOIN job_requisitions jr ON jr.id = jo.requisition_id AND jr.deleted_at IS NULL").
		Joins("LEFT JOIN job_applications ja ON ja.job_opening_id = jo.id AND ja.deleted_at IS NULL").
		Where("jo.deleted_at IS NULL").
		Where("jo.status IN ?", []string{models.JobOpeningStatusActive, models.LegacyStatusPublished}).
		Where("(jo.expiry_date IS NULL OR jo.expiry_date = ? OR jo.expiry_date > ?)", time.Time{}, asOf)
	if strings.TrimSpace(departmentName) != "" {
		q = q.Where("LOWER(TRIM(jr.department)) = LOWER(TRIM(?))", strings.TrimSpace(departmentName))
	}
	q = q.Group("jo.id, jo.job_title, jr.department, jr.location").
		Order("cnt DESC, jo.job_title ASC").
		Limit(limit)
	err := q.Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]models.RecruitmentDashboardTopOpening, 0, len(rows))
	for _, rw := range rows {
		out = append(out, models.RecruitmentDashboardTopOpening{
			JobOpeningID:    rw.JobOpeningID,
			JobTitle:        rw.JobTitle,
			Department:      rw.Department,
			Location:        rw.Location,
			TotalApplicants: rw.Cnt,
		})
	}
	return out, nil
}

// ListUpcomingInterviewsDashboard returns next scheduled interviews with names (max limit).
func (r *RecruitmentDashboardRepository) ListUpcomingInterviewsDashboard(dayStart time.Time, departmentName string, limit int) ([]models.RecruitmentDashboardUpcomingInterview, error) {
	type row struct {
		InterviewID   string
		CandidateID   string
		FirstName     string
		LastName      string
		JobOpeningID  string
		JobTitle      string
		InterviewType string
		ScheduledDate time.Time
		ScheduledTime string
		Status        string
	}
	var rows []row
	q := r.db.Table("interviews i").
		Select(`i.id as interview_id, i.candidate_id, c.first_name, c.last_name, i.job_opening_id as job_opening_id,
			jo.job_title, i.interview_type, i.scheduled_date, i.scheduled_time, i.status`).
		Joins("JOIN candidates c ON c.id = i.candidate_id AND c.deleted_at IS NULL").
		Joins("JOIN job_openings jo ON jo.id = i.job_opening_id AND jo.deleted_at IS NULL").
		Where("i.deleted_at IS NULL AND i.status = ?", models.InterviewStatusScheduled).
		Where("i.scheduled_date >= ?", dayStart.UTC())
	if strings.TrimSpace(departmentName) != "" {
		q = q.Joins("JOIN job_requisitions jr ON jr.id = jo.requisition_id AND jr.deleted_at IS NULL AND LOWER(TRIM(jr.department)) = LOWER(TRIM(?))", strings.TrimSpace(departmentName))
	}
	q = q.Order("i.scheduled_date ASC, i.scheduled_time ASC").Limit(limit)
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]models.RecruitmentDashboardUpcomingInterview, 0, len(rows))
	for _, rw := range rows {
		out = append(out, models.RecruitmentDashboardUpcomingInterview{
			InterviewID:     rw.InterviewID,
			CandidateID:   rw.CandidateID,
			CandidateName: strings.TrimSpace(rw.FirstName + " " + rw.LastName),
			JobOpeningID:  rw.JobOpeningID,
			JobTitle:      rw.JobTitle,
			InterviewType: rw.InterviewType,
			ScheduledDate: rw.ScheduledDate.UTC().Format("2006-01-02"),
			ScheduledTime: rw.ScheduledTime,
			Status:        rw.Status,
		})
	}
	return out, nil
}

// ListRecentApplications returns latest applications with candidate and job opening for dashboard rows.
func (r *RecruitmentDashboardRepository) ListRecentApplications(departmentName string, limit int) ([]models.JobApplication, error) {
	var apps []models.JobApplication
	q := r.db.Model(&models.JobApplication{}).
		Preload("Candidate").
		Preload("JobOpening").
		Where("job_applications.deleted_at IS NULL").
		Order("job_applications.applied_date DESC, job_applications.created_at DESC").
		Limit(limit)
	if strings.TrimSpace(departmentName) != "" {
		q = q.Joins("JOIN job_openings jo ON jo.id = job_applications.job_opening_id AND jo.deleted_at IS NULL").
			Joins("JOIN job_requisitions jr ON jr.id = jo.requisition_id AND jr.deleted_at IS NULL AND LOWER(TRIM(jr.department)) = LOWER(TRIM(?))", strings.TrimSpace(departmentName))
	}
	err := q.Find(&apps).Error
	return apps, err
}

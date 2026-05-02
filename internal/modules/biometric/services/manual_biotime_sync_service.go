package services

import (
	"fmt"
	"strings"
	"time"
)

// ManualBioTimeSyncResult aggregates manual pull-from-BioTime + write-to-DB.
type ManualBioTimeSyncResult struct {
	StartTime             string `json:"start_time"`
	EndTime               string `json:"end_time"`
	APIPages              int    `json:"api_pages"`
	TransactionsFetched   int    `json:"transactions_fetched"`
	Inserted              int    `json:"inserted"`
	SkippedDuplicates     int    `json:"skipped_duplicates"`
	ScanErrors            int    `json:"scan_errors"`
}

// ResolveManualSyncWindow parses optional bounds for POST /biometric/biotime/sync.
// Empty start and end → today 00:00 through now (biometric timezone).
// Start only → end is now. End only → start is 00:00 of that local calendar day.
func ResolveManualSyncWindow(startStr, endStr string) (time.Time, time.Time, error) {
	loc := getBiometricLocation()
	now := time.Now().In(loc)
	s := strings.TrimSpace(startStr)
	e := strings.TrimSpace(endStr)

	if s == "" && e == "" {
		st := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		return st, now, nil
	}

	var start, end time.Time
	var err error

	if s != "" {
		start, err = parseBiometricDateTime(s)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start_time: %w", err)
		}
	}
	if e != "" {
		end, err = parseBiometricDateTime(e)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end_time: %w", err)
		}
	}

	if s == "" {
		en := end.In(loc)
		start = time.Date(en.Year(), en.Month(), en.Day(), 0, 0, 0, 0, loc)
	} else {
		start = start.In(loc)
	}
	if e == "" {
		end = now
	} else {
		end = end.In(loc)
	}

	if start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("start_time must be before or equal to end_time")
	}
	return start, end, nil
}

// ManualBioTimeSyncService runs a direct BioTime → database sync (no message queue).
type ManualBioTimeSyncService struct {
	tenantID *uint
	bioTime  *BioTimeService
	sync     *TransactionSyncService
}

// NewManualBioTimeSyncService creates a direct sync helper for the given tenant.
func NewManualBioTimeSyncService(tenantID *uint) *ManualBioTimeSyncService {
	return &ManualBioTimeSyncService{
		tenantID: tenantID,
		bioTime:  NewBioTimeService(tenantID),
		sync:     NewTransactionSyncService(),
	}
}

// SyncWindow pulls all transaction pages from BioTime for [start, end] and persists new rows.
func (s *ManualBioTimeSyncService) SyncWindow(start, end time.Time) (*ManualBioTimeSyncResult, error) {
	if start.After(end) {
		return nil, fmt.Errorf("start_time must be before or equal to end_time")
	}

	loc := getBiometricLocation()
	startInLoc := start.In(loc)
	endInLoc := end.In(loc)

	startStr := startInLoc.Format("2006-01-02 15:04:05")
	endStr := endInLoc.Format("2006-01-02 15:04:05")

	pageSize := 100
	page := 1

	out := &ManualBioTimeSyncResult{
		StartTime: startStr,
		EndTime:   endStr,
	}

	for {
		p := page
		ps := pageSize
		params := &GetTransactionsParams{
			Page:      &p,
			PageSize:  &ps,
			StartTime: &startStr,
			EndTime:   &endStr,
		}

		resp, err := s.bioTime.GetTransactions(params)
		if err != nil {
			return nil, err
		}
		if resp == nil || len(resp.Data) == 0 {
			break
		}

		out.APIPages++
		out.TransactionsFetched += len(resp.Data)

		stats, err := s.sync.SyncTransactionsDetailed(s.tenantID, resp.Data)
		if err != nil {
			return out, err
		}
		out.Inserted += stats.Inserted
		out.SkippedDuplicates += stats.SkippedDuplicates
		out.ScanErrors += stats.ScanErrors

		if resp.Next == nil || *resp.Next == "" {
			break
		}
		page++
		time.Sleep(150 * time.Millisecond)
	}

	return out, nil
}

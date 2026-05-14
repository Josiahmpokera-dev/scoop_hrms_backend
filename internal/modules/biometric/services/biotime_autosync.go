package services

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
)

// RunBioTimeAutoSyncTick pulls a sliding window from BioTime into biotime_transactions.
// Window length comes from config.BioTime.AutoSyncLookbackHours (default 48, clamped 1–168).
func RunBioTimeAutoSyncTick(tenantID *uint) (*ManualBioTimeSyncResult, error) {
	cfg := config.AppConfig
	if cfg == nil || !cfg.BioTime.Enabled {
		return nil, nil
	}

	lookback := cfg.BioTime.AutoSyncLookbackHours
	if lookback < 1 {
		lookback = 1
	}
	if lookback > 168 {
		lookback = 168
	}

	loc := getBiometricLocation()
	now := time.Now().In(loc)
	start := now.Add(-time.Duration(lookback) * time.Hour)
	end := now

	return NewManualBioTimeSyncService(tenantID).SyncWindow(start, end)
}

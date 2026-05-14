package workers

import (
	"log"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/services"
)

// StartBioTimeAutoSync runs periodic BioTime → DB pulls in the background.
// Disabled when BIOTIME_ENABLED=false or BIOTIME_AUTO_SYNC_ENABLED=false.
func StartBioTimeAutoSync() {
	cfg := config.AppConfig
	if cfg == nil || !cfg.BioTime.Enabled || !cfg.BioTime.AutoSyncEnabled {
		log.Println("BioTime auto-sync: disabled (BIOTIME_ENABLED or BIOTIME_AUTO_SYNC_ENABLED)")
		return
	}

	interval := time.Duration(cfg.BioTime.AutoSyncIntervalMinutes) * time.Minute
	if interval < time.Minute {
		interval = time.Minute
	}

	go func() {
		time.Sleep(15 * time.Second)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		run := func() {
			c := config.AppConfig
			if c == nil || !c.BioTime.Enabled || !c.BioTime.AutoSyncEnabled {
				return
			}
			res, err := services.RunBioTimeAutoSyncTick(nil)
			if err != nil {
				log.Printf("BioTime auto-sync error: %v", err)
				return
			}
			if res == nil {
				return
			}
			log.Printf("BioTime auto-sync: fetched=%d inserted=%d skipped_dup=%d pages=%d window=[%s .. %s]",
				res.TransactionsFetched, res.Inserted, res.SkippedDuplicates, res.APIPages, res.StartTime, res.EndTime)
		}

		run()
		for range ticker.C {
			run()
		}
	}()

	log.Printf("BioTime auto-sync: started (every %s, lookback %dh)", interval, cfg.BioTime.AutoSyncLookbackHours)
}

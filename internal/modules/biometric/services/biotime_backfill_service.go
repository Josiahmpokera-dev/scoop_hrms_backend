package services

import (
	"fmt"
	"log"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
)

// BioTimeBackfillService handles backfilling historical BioTime transactions
type BioTimeBackfillService struct {
	bioTimeService *BioTimeService
	queue          QueuePublisher
	tenantID       *uint
}

// NewBioTimeBackfillService creates a new backfill service
func NewBioTimeBackfillService(tenantID *uint) *BioTimeBackfillService {
	var queue QueuePublisher
	cfg := config.AppConfig
	if cfg != nil && cfg.RabbitMQ.Enabled {
		if q, err := NewQueueClient(); err == nil {
			queue = q
		}
	}

	return &BioTimeBackfillService{
		bioTimeService: NewBioTimeService(tenantID),
		queue:          queue,
		tenantID:       tenantID,
	}
}

// BackfillTransactions pulls all transactions from a start date to end date and queues them
func (s *BioTimeBackfillService) BackfillTransactions(startDate, endDate time.Time) error {
	if s.queue == nil {
		return fmt.Errorf("RabbitMQ queue is not available. Please enable RabbitMQ in configuration")
	}

	log.Printf("Starting backfill from %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Process in monthly chunks to avoid overwhelming the API
	currentDate := startDate
	totalQueued := 0
	page := 1
	pageSize := 100

	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		// Calculate end of current month
		nextMonth := currentDate.AddDate(0, 1, 0)
		chunkEnd := nextMonth.AddDate(0, 0, -1) // Last day of current month
		if chunkEnd.After(endDate) {
			chunkEnd = endDate
		}

		log.Printf("Processing chunk: %s to %s", currentDate.Format("2006-01-02"), chunkEnd.Format("2006-01-02"))

		// Fetch all transactions for this month
		startTimeStr := currentDate.Format("2006-01-02 00:00:00")
		endTimeStr := chunkEnd.Format("2006-01-02 23:59:59")

		page = 1
		for {
			params := &GetTransactionsParams{
				Page:      &page,
				PageSize:  &pageSize,
				StartTime: &startTimeStr,
				EndTime:   &endTimeStr,
			}

			resp, err := s.bioTimeService.GetTransactions(params)
			if err != nil {
				return fmt.Errorf("failed to fetch transactions for %s: %w", currentDate.Format("2006-01-02"), err)
			}

			if resp == nil || len(resp.Data) == 0 {
				log.Printf("No more transactions for %s", currentDate.Format("2006-01-02"))
				break
			}

			// Queue transactions
			transactions := make([]interface{}, len(resp.Data))
			for i, txn := range resp.Data {
				transactions[i] = txn
			}

			if err := s.queue.PublishTransaction(s.tenantID, transactions); err != nil {
				log.Printf("Warning: Failed to queue batch: %v", err)
				// Continue with next batch
			} else {
				totalQueued += len(transactions)
				log.Printf("Queued %d transactions (page %d) for %s. Total queued: %d", 
					len(transactions), page, currentDate.Format("2006-01-02"), totalQueued)
			}

			// Check if there are more pages
			if resp.Next == nil || *resp.Next == "" {
				break
			}

			page++
			// Small delay to avoid overwhelming the API
			time.Sleep(500 * time.Millisecond)
		}

		// Move to next month
		currentDate = nextMonth
	}

	log.Printf("Backfill completed. Total transactions queued: %d", totalQueued)
	return nil
}

// BackfillFrom2025 pulls all transactions from 2025-01-01 to now
func (s *BioTimeBackfillService) BackfillFrom2025() error {
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Now()
	return s.BackfillTransactions(startDate, endDate)
}

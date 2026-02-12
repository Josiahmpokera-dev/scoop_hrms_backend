package services

import (
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/models"
	biometricRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/repositories"
)

// Transaction is imported from biotime_service.go - using the same type
// This avoids circular dependencies

// TransactionSyncService handles syncing BioTime transactions to database
type TransactionSyncService struct {
	transactionRepo *biometricRepos.BioTimeTransactionRepository
}

// NewTransactionSyncService creates a new transaction sync service
func NewTransactionSyncService() *TransactionSyncService {
	return &TransactionSyncService{
		transactionRepo: biometricRepos.NewBioTimeTransactionRepository(),
	}
}

// SyncTransactions syncs BioTime transactions to the database
func (s *TransactionSyncService) SyncTransactions(tenantID *uint, transactions []Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	var dbTransactions []models.BioTimeTransaction
	now := time.Now()

	for _, txn := range transactions {
		// Check if transaction already exists
		exists, err := s.transactionRepo.Exists(txn.ID.Value, tenantID, s.parsePunchTime(txn.PunchTime))
		if err != nil {
			// Log error but continue with other transactions
			fmt.Printf("Error checking transaction existence: %v\n", err)
			continue
		}

		if exists {
			// Skip if already exists
			continue
		}

		// Convert BioTime transaction to database model
		lastName := ""
		if txn.LastName != nil {
			lastName = *txn.LastName
		}
		position := ""
		if txn.Position != nil {
			position = *txn.Position
		}
		gpsLocation := ""
		if txn.GPSLocation != nil {
			gpsLocation = *txn.GPSLocation
		}

		dbTxn := models.BioTimeTransaction{
			BioTimeTransactionID: txn.ID.Value,
			EmpCode:             txn.EmpCode,
			FirstName:           txn.FirstName,
			LastName:            lastName,
			Department:          txn.Department,
			Position:            position,
			PunchTime:           s.parsePunchTime(txn.PunchTime),
			PunchState:          txn.PunchState,
			PunchStateDisplay:   txn.PunchStateDisplay,
			VerifyType:          txn.VerifyType.Value,
			VerifyTypeDisplay:   txn.VerifyTypeDisplay,
			WorkCode:            txn.WorkCode,
			GPSLocation:         gpsLocation,
			TerminalSN:          txn.TerminalSN,
			Temperature:         txn.Temperature,
			SyncedAt:            now,
		}

		// Handle nullable fields
		if txn.AreaAlias != nil {
			dbTxn.AreaAlias = txn.AreaAlias
		}
		if txn.TerminalAlias != nil {
			dbTxn.TerminalAlias = txn.TerminalAlias
		}
		if txn.UploadTime != nil && *txn.UploadTime != "" {
			if uploadTime := s.parsePunchTime(*txn.UploadTime); !uploadTime.IsZero() {
				dbTxn.UploadTime = &uploadTime
			}
		}

		dbTransactions = append(dbTransactions, dbTxn)
	}

	// Bulk insert new transactions
	if len(dbTransactions) > 0 {
		if err := s.transactionRepo.BulkCreate(dbTransactions); err != nil {
			return fmt.Errorf("failed to bulk create transactions: %w", err)
		}
		fmt.Printf("Successfully synced %d new transactions to database\n", len(dbTransactions))
	} else {
		fmt.Printf("No new transactions to sync (all %d already exist)\n", len(transactions))
	}

	return nil
}

// parsePunchTime parses the punch time string from BioTime API
func (s *TransactionSyncService) parsePunchTime(timeStr string) time.Time {
	// BioTime format: "2020-06-05 00:00:00"
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return t
		}
	}

	// Return zero time if parsing fails
	return time.Time{}
}

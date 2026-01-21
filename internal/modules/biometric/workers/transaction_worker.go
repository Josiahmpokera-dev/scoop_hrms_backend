package workers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	biometricServices "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

// TransactionWorker processes BioTime transaction messages from RabbitMQ
type TransactionWorker struct {
	queueClient *queue.RabbitMQClient
	syncService *biometricServices.TransactionSyncService
}

// NewTransactionWorker creates a new transaction worker
func NewTransactionWorker() (*TransactionWorker, error) {
	cfg := config.AppConfig
	if cfg == nil || !cfg.RabbitMQ.Enabled {
		return nil, fmt.Errorf("RabbitMQ is not enabled")
	}

	queueClient, err := queue.NewRabbitMQClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create queue client: %w", err)
	}

	return &TransactionWorker{
		queueClient: queueClient,
		syncService: biometricServices.NewTransactionSyncService(),
	}, nil
}

// Start starts consuming messages from the queue
func (w *TransactionWorker) Start() error {
	cfg := config.AppConfig
	if cfg == nil {
		return fmt.Errorf("config not loaded")
	}

	ch := w.queueClient.GetChannel()
	qName := w.queueClient.GetQueueName()

	// Set QoS based on batch size configuration
	batchSize := cfg.RabbitMQ.BatchSize
	if batchSize <= 0 {
		batchSize = 1 // Process one at a time by default
	}

	err := ch.Qos(
		batchSize, // prefetch count (based on batch size)
		0,         // prefetch size
		false,     // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.Consume(
		qName, // queue
		"",    // consumer
		false, // auto-ack (we'll manually ack)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	processingInterval := time.Duration(cfg.RabbitMQ.ProcessingInterval) * time.Second
	if processingInterval > 0 {
		log.Printf("Transaction worker started with %d second processing interval, batch size: %d, queue: %s", 
			cfg.RabbitMQ.ProcessingInterval, batchSize, qName)
	} else {
		log.Printf("Transaction worker started (immediate processing), batch size: %d, queue: %s", batchSize, qName)
	}

	// Process messages with configurable interval
	go func() {
		var messageBatch []amqp.Delivery
		var ticker *time.Ticker
		
		// If processing interval is 0, process immediately
		if processingInterval > 0 {
			ticker = time.NewTicker(processingInterval)
			defer ticker.Stop()
		}

		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					// Channel closed, process remaining batch
					if len(messageBatch) > 0 {
						w.processBatch(messageBatch)
					}
					return
				}
				messageBatch = append(messageBatch, d)

				// Process batch if it reaches the configured size (and interval is 0 or batch is full)
				if processingInterval == 0 || len(messageBatch) >= batchSize {
					w.processBatch(messageBatch)
					messageBatch = nil
				}

			case <-func() <-chan time.Time {
				if ticker != nil {
					return ticker.C
				}
				// Return a channel that never sends if ticker is nil
				c := make(chan time.Time)
				close(c)
				return c
			}():
				// Process batch at interval
				if len(messageBatch) > 0 {
					w.processBatch(messageBatch)
					messageBatch = nil
				}
			}
		}
	}()

	return nil
}

// processBatch processes a batch of messages
func (w *TransactionWorker) processBatch(messages []amqp.Delivery) {
	log.Printf("Processing batch of %d messages", len(messages))
	
	for _, d := range messages {
		if err := w.processMessage(d); err != nil {
			log.Printf("Error processing message in batch: %v", err)
			d.Nack(false, true)
		} else {
			d.Ack(false)
		}
	}
}

// processMessage processes a single message from the queue
func (w *TransactionWorker) processMessage(d amqp.Delivery) error {
	var message struct {
		TenantID    *uint        `json:"tenant_id"`
		Transactions []interface{} `json:"transactions"`
	}

	if err := json.Unmarshal(d.Body, &message); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Convert interface{} to Transaction structs
	transactions := make([]biometricServices.Transaction, 0, len(message.Transactions))
	for _, txnInterface := range message.Transactions {
		// Convert from JSON map to Transaction struct
		txnBytes, err := json.Marshal(txnInterface)
		if err != nil {
			log.Printf("Error marshaling transaction: %v", err)
			continue
		}

		var txn biometricServices.Transaction
		if err := json.Unmarshal(txnBytes, &txn); err != nil {
			log.Printf("Error unmarshaling transaction: %v", err)
			continue
		}

		transactions = append(transactions, txn)
	}

	if len(transactions) == 0 {
		log.Println("No valid transactions to process")
		return nil
	}

	// Sync transactions to database
	if err := w.syncService.SyncTransactions(message.TenantID, transactions); err != nil {
		return fmt.Errorf("failed to sync transactions: %w", err)
	}

	log.Printf("Successfully processed %d transactions for tenant %v", len(transactions), message.TenantID)
	return nil
}

// Stop stops the worker
func (w *TransactionWorker) Stop() error {
	return w.queueClient.Close()
}

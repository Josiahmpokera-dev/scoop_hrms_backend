package services

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/queue"
)

// NewQueueClient creates a new queue client
func NewQueueClient() (QueuePublisher, error) {
	return queue.NewRabbitMQClient()
}

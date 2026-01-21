package queue

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient handles RabbitMQ connections and operations
type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *config.RabbitMQConfig
}

// NewRabbitMQClient creates a new RabbitMQ client
func NewRabbitMQClient() (*RabbitMQClient, error) {
	cfg := config.AppConfig
	if cfg == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	if !cfg.RabbitMQ.Enabled {
		return nil, fmt.Errorf("RabbitMQ is disabled")
	}

	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		cfg.RabbitMQ.Exchange, // name
		"direct",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue
	_, err = ch.QueueDeclare(
		cfg.RabbitMQ.Queue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		cfg.RabbitMQ.Queue,    // queue name
		"transactions",        // routing key
		cfg.RabbitMQ.Exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: ch,
		config:  &cfg.RabbitMQ,
	}, nil
}

// PublishTransaction publishes a transaction to the queue
func (r *RabbitMQClient) PublishTransaction(tenantID *uint, transactions []interface{}) error {
	if !r.config.Enabled {
		return fmt.Errorf("RabbitMQ is disabled")
	}

	message := map[string]interface{}{
		"tenant_id":   tenantID,
		"transactions": transactions,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = r.channel.Publish(
		r.config.Exchange, // exchange
		"transactions",    // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published %d transactions to RabbitMQ queue", len(transactions))
	return nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQClient) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// GetChannel returns the RabbitMQ channel for consuming messages
func (r *RabbitMQClient) GetChannel() *amqp.Channel {
	return r.channel
}

// GetQueueName returns the queue name
func (r *RabbitMQClient) GetQueueName() string {
	return r.config.Queue
}

# RabbitMQ Setup Guide

## Overview
This application uses RabbitMQ to queue BioTime transactions for background processing and storage in the database.

## Installation

### Install RabbitMQ Dependency

```bash
go get github.com/rabbitmq/amqp091-go
```

Or if using Go modules (recommended):
```bash
go get github.com/rabbitmq/amqp091-go@latest
```

**Note:** We use `github.com/rabbitmq/amqp091-go` which is the maintained version. The old `github.com/streadway/amqp` package is deprecated.

### Install RabbitMQ Server

**macOS (using Homebrew):**
```bash
brew install rabbitmq
brew services start rabbitmq
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get update
sudo apt-get install rabbitmq-server
sudo systemctl start rabbitmq-server
sudo systemctl enable rabbitmq-server
```

**Docker:**
```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

## Configuration

Add the following to your `.env` file:

```env
# RabbitMQ Configuration
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_ENABLED=true
RABBITMQ_EXCHANGE=biotime_exchange
RABBITMQ_QUEUE=biotime_transactions
```

### Configuration Parameters

- `RABBITMQ_URL`: RabbitMQ connection URL (default: `amqp://guest:guest@localhost:5672/`)
- `RABBITMQ_ENABLED`: Enable/disable RabbitMQ (default: `true`)
- `RABBITMQ_EXCHANGE`: Exchange name for messages (default: `biotime_exchange`)
- `RABBITMQ_QUEUE`: Queue name for transactions (default: `biotime_transactions`)

## How It Works

1. **API Request**: When you call `GET /api/v1/biometric/biotime/transactions`, the API:
   - Returns data directly from BioTime API (no delay)
   - Queues transactions in the background for database storage

2. **Background Processing**: A worker process:
   - Consumes messages from RabbitMQ queue
   - Checks if transactions already exist in database
   - Stores only new transactions
   - Prevents duplicates

3. **Database Storage**: Transactions are stored in `biotime_transactions` table with:
   - All transaction data from BioTime
   - Tenant isolation
   - Duplicate prevention

## Verification

1. **Check RabbitMQ is running:**
   ```bash
   # macOS/Linux
   rabbitmqctl status
   
   # Or check management UI (if enabled)
   # http://localhost:15672 (guest/guest)
   ```

2. **Check worker is running:**
   - Look for log message: `✅ Transaction worker started successfully`

3. **Test the flow:**
   - Call the transactions API
   - Check database for new records in `biotime_transactions` table

## Troubleshooting

### Worker not starting
- Check RabbitMQ is running: `rabbitmqctl status`
- Verify `RABBITMQ_ENABLED=true` in `.env`
- Check connection URL is correct

### Transactions not being stored
- Check worker logs for errors
- Verify database connection
- Check RabbitMQ queue for messages: `rabbitmqctl list_queues`

### Duplicate transactions
- The system prevents duplicates by checking `biotime_transaction_id` and `punch_time`
- If duplicates appear, check the unique constraint logic

## Production Considerations

1. **RabbitMQ Clustering**: Set up RabbitMQ cluster for high availability
2. **Message Persistence**: Messages are marked as persistent
3. **Error Handling**: Failed messages are requeued automatically
4. **Monitoring**: Monitor queue depth and worker health
5. **Credentials**: Change default `guest/guest` credentials in production

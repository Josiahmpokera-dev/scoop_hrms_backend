# BioTime Transaction Sync with RabbitMQ

## Overview
This document describes the background queue system that automatically syncs BioTime transactions to the database using RabbitMQ.

## Architecture

### Flow Diagram
```
API Request → BioTime API → Response (immediate)
                ↓
         Queue to RabbitMQ (background)
                ↓
         Worker consumes messages
                ↓
         Check for duplicates
                ↓
         Store new transactions in DB
```

### Components

1. **API Handler** (`biotime_handler.go`)
   - Returns data directly from BioTime API (no delay)
   - Queues transactions in background goroutine

2. **RabbitMQ Queue** (`queue/rabbitmq.go`)
   - Publishes transactions to queue
   - Handles connection and message persistence

3. **Worker** (`workers/transaction_worker.go`)
   - Consumes messages from queue
   - Processes transactions asynchronously

4. **Sync Service** (`services/transaction_sync_service.go`)
   - Checks for duplicate transactions
   - Stores only new transactions in database

5. **Database Model** (`models/biotime_transaction.go`)
   - Stores transaction data with tenant isolation
   - Prevents duplicates

## Setup

### 1. Install RabbitMQ Dependency

```bash
go get github.com/streadway/amqp
```

### 2. Install RabbitMQ Server

**macOS:**
```bash
brew install rabbitmq
brew services start rabbitmq
```

**Linux:**
```bash
sudo apt-get install rabbitmq-server
sudo systemctl start rabbitmq-server
```

**Docker:**
```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

### 3. Configure Environment Variables

Add to your `.env` file:

```env
# RabbitMQ Configuration
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_ENABLED=true
RABBITMQ_EXCHANGE=biotime_exchange
RABBITMQ_QUEUE=biotime_transactions
```

### 4. Restart Application

The worker will automatically start when the application starts (if RabbitMQ is enabled).

## How It Works

### API Request Flow

1. **User calls API:**
   ```bash
   GET /api/v1/biometric/biotime/transactions
   ```

2. **API Response (immediate):**
   - Fetches data from BioTime API
   - Returns response immediately
   - **No delay for database operations**

3. **Background Processing:**
   - Transactions are queued to RabbitMQ in a goroutine
   - Worker processes queue asynchronously
   - Only new transactions are stored

### Duplicate Prevention

The system prevents duplicate transactions by checking:
- `biotime_transaction_id` (ID from BioTime)
- `punch_time` (timestamp)
- `tenant_id` (tenant isolation)

If a transaction with the same ID and punch time already exists for the tenant, it's skipped.

### Database Storage

Transactions are stored in the `biotime_transactions` table with:
- All transaction fields from BioTime
- Tenant isolation
- Automatic duplicate prevention
- `synced_at` timestamp

## API Behavior

### Current Behavior (Unchanged)
- ✅ API still returns data directly from BioTime
- ✅ Response time is not affected
- ✅ All existing functionality preserved

### New Background Behavior
- ✅ Transactions are queued automatically
- ✅ Worker processes queue in background
- ✅ Only new transactions are stored
- ✅ No impact on API response time

## Monitoring

### Check Worker Status

Look for this log message on startup:
```
✅ Transaction worker started successfully
```

### Check Queue Status

```bash
# List queues
rabbitmqctl list_queues

# Check queue depth
rabbitmqctl list_queues name messages
```

### Check Database

```sql
-- View synced transactions
SELECT COUNT(*) FROM biotime_transactions;

-- View recent syncs
SELECT * FROM biotime_transactions 
ORDER BY synced_at DESC 
LIMIT 10;
```

## Troubleshooting

### Worker Not Starting

**Symptoms:**
- No "Transaction worker started" message
- Transactions not being stored

**Solutions:**
1. Check RabbitMQ is running: `rabbitmqctl status`
2. Verify `RABBITMQ_ENABLED=true` in `.env`
3. Check connection URL is correct
4. Review application logs for errors

### Transactions Not Being Stored

**Symptoms:**
- API works but no database records

**Solutions:**
1. Check worker is running (see logs)
2. Verify queue has messages: `rabbitmqctl list_queues`
3. Check database connection
4. Review worker logs for processing errors

### Duplicate Transactions

**Symptoms:**
- Same transaction stored multiple times

**Solutions:**
- The system should prevent this automatically
- Check the `Exists()` method in repository
- Verify unique constraint logic

### Queue Messages Piling Up

**Symptoms:**
- Queue depth increasing
- Worker not processing fast enough

**Solutions:**
1. Check worker is running
2. Review worker logs for errors
3. Consider scaling workers (future enhancement)
4. Check database performance

## Configuration Options

### Disable RabbitMQ

Set in `.env`:
```env
RABBITMQ_ENABLED=false
```

When disabled:
- API still works normally
- Transactions are not queued
- No database sync happens

### Custom Queue Settings

```env
RABBITMQ_EXCHANGE=biotime_exchange  # Exchange name
RABBITMQ_QUEUE=biotime_transactions # Queue name
```

## Database Schema

### biotime_transactions Table

```sql
CREATE TABLE biotime_transactions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT,
    biotime_transaction_id INTEGER NOT NULL,
    emp_code VARCHAR(100),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    department VARCHAR(255),
    position VARCHAR(255),
    punch_time TIMESTAMP,
    punch_state VARCHAR(50),
    punch_state_display VARCHAR(100),
    verify_type INTEGER,
    verify_type_display VARCHAR(100),
    work_code VARCHAR(100),
    gps_location TEXT,
    area_alias VARCHAR(255),
    terminal_sn VARCHAR(100),
    temperature DOUBLE PRECISION,
    terminal_alias VARCHAR(255),
    upload_time TIMESTAMP,
    synced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## Benefits

✅ **Non-blocking**: API response is immediate  
✅ **Reliable**: RabbitMQ ensures message delivery  
✅ **Scalable**: Can process large volumes asynchronously  
✅ **Duplicate Prevention**: Only new transactions are stored  
✅ **Tenant Isolation**: Each tenant's data is separate  
✅ **Fault Tolerant**: Failed messages are requeued  
✅ **No Breaking Changes**: Existing API functionality preserved

## Future Enhancements

- Multiple worker instances for scaling
- Retry logic with exponential backoff
- Dead letter queue for failed messages
- Monitoring dashboard
- Batch processing optimizations
- Transaction history API endpoints

# BioTime Biometric Device API Integration

## Overview
This document provides API documentation for the BioTime biometric device integration. BioTime is used for employee attendance tracking through biometric devices.

## Base URL
All API endpoints use the base URL: `{{BASE_URL}}/api/v1/biometric/biotime`

## Authentication
All requests require Bearer token authentication and HR or Admin role:
```
Authorization: Bearer <access_token>
```

**Note:** All endpoints require `HRMiddleware()` which allows both `admin` and `hr` roles.

---

## Configuration

The BioTime integration is configured through environment variables in your `.env` file:

```env
# BioTime Biometric Device Configuration
BIOTIME_BASE_URL=http://10.4.9.24:8087
BIOTIME_USERNAME=Developer
BIOTIME_PASSWORD=Developer@123
BIOTIME_ENABLED=true
```

### Configuration Parameters

- `BIOTIME_BASE_URL` (string, required): Base URL of the BioTime API server (e.g., `http://10.4.9.24:8087`)
- `BIOTIME_USERNAME` (string, required): Username for BioTime API authentication
- `BIOTIME_PASSWORD` (string, required): Password for BioTime API authentication
- `BIOTIME_ENABLED` (boolean, optional): Enable/disable BioTime integration (default: `true`)

---

## 1. Test Connection

### Endpoint
```
GET /api/v1/biometric/biotime/test-connection
```

### Description
Test the connection and authentication with the BioTime API. This endpoint:
1. Attempts to connect to the BioTime server
2. Authenticates using the configured credentials
3. Retrieves a JWT token
4. Returns connection status and authentication result

### Success Response (200)
```json
{
  "success": true,
  "message": "Successfully connected and authenticated with BioTime API",
  "data": {
    "success": true,
    "message": "Successfully connected and authenticated with BioTime API",
    "base_url": "http://10.4.9.24:8087",
    "authenticated": true,
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9..."
  }
}
```

### Error Response (400/500)
```json
{
  "success": false,
  "message": "Failed to authenticate with BioTime API",
  "data": {
    "success": false,
    "message": "Failed to authenticate with BioTime API",
    "base_url": "http://10.4.9.24:8087",
    "authenticated": false,
    "error": "authentication failed with status 401: Invalid credentials"
  }
}
```

### Example Request (cURL)
```bash
curl -X GET "http://localhost:8080/api/v1/biometric/biotime/test-connection" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Example Request (JavaScript)
```javascript
const response = await fetch('http://localhost:8080/api/v1/biometric/biotime/test-connection', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
const data = await response.json();
console.log(data);
```

---

## 2. Get Terminals

### Endpoint
```
GET /api/v1/biometric/biotime/terminals
```

### Description
Fetch the list of terminal devices from the BioTime API. This endpoint:
1. Automatically retrieves a valid BioTime token (from database or authenticates)
2. Makes an authenticated request to BioTime API
3. Returns the list of terminal devices

### Success Response (200)
```json
{
  "success": true,
  "message": "Terminals retrieved successfully",
  "data": {
    "count": 6,
    "next": null,
    "previous": null,
    "msg": "",
    "code": 0,
    "data": [
      {
        "id": 5,
        "sn": "A6KX192060002",
        "ip_address": "172.30.7.162",
        "alias": "Auto add",
        "terminal_name": null,
        "fw_ver": null,
        "push_ver": null,
        "state": 1,
        "terminal_tz": 8,
        "area": {
          "id": 1,
          "area_code": "1",
          "area_name": "Not Authorized"
        },
        "last_activity": "2020-06-02 15:04:38",
        "user_count": null,
        "fp_count": null,
        "face_count": null,
        "palm_count": null,
        "transaction_count": null,
        "push_time": null,
        "transfer_time": "00:00;14:05",
        "transfer_interval": 1,
        "is_attendance": 1,
        "area_name": "Not Authorized"
      }
    ]
  }
}
```

### Response Fields

- `count`: Total number of terminals
- `next`: URL for next page (null if no more pages)
- `previous`: URL for previous page (null if first page)
- `msg`: Response message from BioTime API
- `code`: Response code (0 = success)
- `data`: Array of terminal objects

### Terminal Object Fields

- `id`: Terminal ID
- `sn`: Serial number
- `ip_address`: IP address of the terminal
- `alias`: Terminal alias/name
- `terminal_name`: Terminal name (can be null)
- `fw_ver`: Firmware version (can be null)
- `push_ver`: Push version (can be null)
- `state`: Terminal state (1 = active, etc.)
- `terminal_tz`: Terminal timezone offset
- `area`: Area information object
  - `id`: Area ID
  - `area_code`: Area code
  - `area_name`: Area name
- `last_activity`: Last activity timestamp
- `user_count`: Number of users enrolled (can be null)
- `fp_count`: Fingerprint count (can be null)
- `face_count`: Face count (can be null)
- `palm_count`: Palm count (can be null)
- `transaction_count`: Transaction count (can be null)
- `push_time`: Push time (can be null)
- `transfer_time`: Transfer time schedule
- `transfer_interval`: Transfer interval in hours
- `is_attendance`: Whether terminal is used for attendance (1 = yes, 0 = no)
- `area_name`: Area name (duplicate of area.area_name)

### Error Response (500)
```json
{
  "success": false,
  "message": "Failed to fetch terminals from BioTime",
  "error": "failed to fetch terminals: status 401, response: Unauthorized"
}
```

### Example Request (cURL)
```bash
curl -X GET "http://localhost:8080/api/v1/biometric/biotime/terminals" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Example Request (JavaScript)
```javascript
const response = await fetch('http://localhost:8080/api/v1/biometric/biotime/terminals', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
const data = await response.json();
console.log(data);
```

### Notes
- The BioTime token is automatically retrieved and used - you don't need to manage it
- If the token is expired, it will be automatically renewed
- The endpoint requires HR or Admin role

---

## 3. Get Transactions

### Endpoint
```
GET /api/v1/biometric/biotime/transactions
```

### Description
Fetch attendance/transaction records from the BioTime API. This endpoint supports filtering and pagination. The BioTime token is automatically managed.

### Query Parameters

All parameters are optional:

| Parameter | Type | Description |
|-----------|------|-------------|
| `page` | integer | Page number for pagination (default: 1) |
| `page_size` | integer | Number of records per page |
| `emp_code` | string | Filter by employee code |
| `terminal_sn` | string | Filter by terminal serial number |
| `terminal_alias` | string | Filter by terminal alias |
| `start_time` | string | Start time filter (format: `YYYY-MM-DD HH:MM:SS`) |
| `end_time` | string | End time filter (format: `YYYY-MM-DD HH:MM:SS`) |

### Success Response (200)
```json
{
  "success": true,
  "message": "Transactions retrieved successfully",
  "data": {
    "count": 6,
    "next": null,
    "previous": null,
    "msg": "",
    "code": 0,
    "data": [
      {
        "id": 1,
        "emp_code": "11111111122",
        "first_name": "",
        "last_name": "",
        "department": "Department",
        "position": "Position",
        "punch_time": "2020-06-05 00:00:00",
        "punch_state": "0",
        "punch_state_display": "Check In",
        "verify_type": 0,
        "verify_type_display": "Password",
        "work_code": "",
        "gps_location": "",
        "area_alias": null,
        "terminal_sn": "",
        "temperature": 0.0,
        "terminal_alias": null,
        "upload_time": "2020-06-05 08:47:59"
      }
    ]
  }
}
```

### Response Fields

- `count`: Total number of transactions
- `next`: URL for next page (null if no more pages)
- `previous`: URL for previous page (null if first page)
- `msg`: Response message from BioTime API
- `code`: Response code (0 = success)
- `data`: Array of transaction objects

### Transaction Object Fields

- `id`: Transaction ID
- `emp_code`: Employee code
- `first_name`: Employee first name
- `last_name`: Employee last name
- `department`: Department name
- `position`: Position name
- `punch_time`: Timestamp of the punch/attendance record
- `punch_state`: Punch state code (e.g., "0" for Check In, "1" for Check Out)
- `punch_state_display`: Human-readable punch state (e.g., "Check In", "Check Out")
- `verify_type`: Verification type code (0 = Password, 1 = Fingerprint, etc.)
- `verify_type_display`: Human-readable verification type
- `work_code`: Work code
- `gps_location`: GPS coordinates if available
- `area_alias`: Area alias (can be null)
- `terminal_sn`: Terminal serial number
- `temperature`: Temperature reading if available
- `terminal_alias`: Terminal alias (can be null)
- `upload_time`: When the transaction was uploaded to the server

### Example Request (cURL)
```bash
# Get all transactions
curl -X GET "http://localhost:8080/api/v1/biometric/biotime/transactions" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get transactions with filters
curl -X GET "http://localhost:8080/api/v1/biometric/biotime/transactions?emp_code=11111111122&start_time=2020-06-01 00:00:00&end_time=2020-06-30 23:59:59&page=1&page_size=50" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get transactions for a specific terminal
curl -X GET "http://localhost:8080/api/v1/biometric/biotime/transactions?terminal_sn=A6KX192060002" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Example Request (JavaScript)
```javascript
// Get all transactions
const response = await fetch('http://localhost:8080/api/v1/biometric/biotime/transactions', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
const data = await response.json();
console.log(data);

// Get transactions with filters
const params = new URLSearchParams({
  emp_code: '11111111122',
  start_time: '2020-06-01 00:00:00',
  end_time: '2020-06-30 23:59:59',
  page: '1',
  page_size: '50'
});
const filteredResponse = await fetch(`http://localhost:8080/api/v1/biometric/biotime/transactions?${params}`, {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
const filteredData = await filteredResponse.json();
console.log(filteredData);
```

### Error Response (500)
```json
{
  "success": false,
  "message": "Failed to fetch transactions from BioTime",
  "error": "failed to fetch transactions: status 401, response: Unauthorized"
}
```

### Notes
- The BioTime token is automatically retrieved and used - you don't need to manage it
- If the token is expired, it will be automatically renewed
- The endpoint requires HR or Admin role
- Time format should be `YYYY-MM-DD HH:MM:SS` (e.g., `2020-06-01 00:00:00`)
- All query parameters are optional - you can combine them as needed
- Pagination is handled by BioTime API using `page` and `page_size` parameters

---

## 4. Get Transaction by ID

### Endpoint
```
GET /api/v1/biometric/biotime/transactions/{id}
```

### Description
Fetch a single attendance/transaction record by ID from the BioTime API. The BioTime token is automatically managed.

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Transaction ID |

### Success Response (200)
```json
{
  "success": true,
  "message": "Transaction retrieved successfully",
  "data": {
    "id": 1,
    "emp_code": "11111111122",
    "first_name": "",
    "last_name": "",
    "department": "Department",
    "position": "Position",
    "punch_time": "2020-06-05 00:00:00",
    "punch_state": "0",
    "punch_state_display": "Check In",
    "verify_type": 0,
    "verify_type_display": "Password",
    "work_code": "",
    "gps_location": "",
    "area_alias": null,
    "terminal_sn": "",
    "temperature": 0.0,
    "terminal_alias": null,
    "upload_time": "2020-06-05 08:47:59"
  }
}
```

### Error Response (400)
```json
{
  "success": false,
  "message": "Transaction ID is required",
  "error": null
}
```

### Error Response (404)
```json
{
  "success": false,
  "message": "Transaction not found",
  "error": null
}
```

### Error Response (500)
```json
{
  "success": false,
  "message": "Failed to fetch transaction from BioTime",
  "error": "failed to fetch transaction: status 500, response: Internal Server Error"
}
```

### Example Request (cURL)
```bash
# Get transaction by ID
curl -X GET "http://localhost:8080/api/v1/biometric/biotime/transactions/1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Example Request (JavaScript)
```javascript
const transactionId = "1";
const response = await fetch(`http://localhost:8080/api/v1/biometric/biotime/transactions/${transactionId}`, {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
const data = await response.json();
console.log(data);
```

### Notes
- The BioTime token is automatically retrieved and used - you don't need to manage it
- If the token is expired, it will be automatically renewed
- The endpoint requires HR or Admin role
- The transaction ID should match the ID from the BioTime system

---

## 5. Refresh Authentication Token

### Endpoint
```
POST /api/v1/biometric/biotime/refresh-token
```

### Description
Manually force a refresh of the BioTime authentication token. This endpoint will:
1. Ignore the current stored token (even if it's still valid)
2. Authenticate with BioTime API to get a new token
3. Store the new token in the database
4. Return the new token

**Note:** This is usually not necessary because tokens are automatically refreshed when expired. Use this only if you need to force a refresh.

### Success Response (200)
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "note": "New token has been stored in database and will be used for future requests."
  }
}
```

### Error Response (500)
```json
{
  "success": false,
  "message": "Failed to refresh BioTime token",
  "error": "authentication failed with status 401: Invalid credentials"
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/biometric/biotime/refresh-token" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Example Request (JavaScript)
```javascript
const response = await fetch('http://localhost:8080/api/v1/biometric/biotime/refresh-token', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
const data = await response.json();
console.log(data);
```

### Notes
- This endpoint forces a new token even if the current one is still valid
- The new token is automatically stored in the database
- Future API calls will use the newly refreshed token
- Use this only when you need to manually refresh (e.g., after credential changes)

---

## 6. Get Authentication Token

### Endpoint
```
GET /api/v1/biometric/biotime/token
```

### Description
Retrieve the current authentication token for the BioTime API. The service automatically handles token caching and renewal.

### Success Response (200)
```json
{
  "success": true,
  "message": "Token retrieved successfully",
  "data": {
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9..."
  }
}
```

### Error Response (500)
```json
{
  "success": false,
  "message": "Failed to get BioTime token",
  "error": "authentication failed with status 401: Invalid credentials"
}
```

### Example Request (cURL)
```bash
curl -X GET "http://localhost:8080/api/v1/biometric/biotime/token" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## BioTime API Authentication Flow

The service automatically handles authentication with the BioTime API:

1. **Token Request**: When a token is needed, the service makes a POST request to:
   ```
   POST http://10.4.9.24:8087/jwt-api-token-auth/
   ```

2. **Request Body**:
   ```json
   {
     "username": "Developer",
     "password": "Developer@123"
   }
   ```

3. **Response**: The API returns a JWT token that is cached for 24 hours (or until expiry).

4. **Token Usage**: All subsequent API requests include the token in the Authorization header:
   ```
   Authorization: Bearer <token>
   ```

---

## Service Features

### Automatic Token Management
- **Database Storage**: BioTime tokens are stored in the database (`biotime_configs` table) per tenant
- **Automatic Caching**: Tokens are automatically cached in the database to avoid unnecessary authentication requests
- **Token Validation**: The service checks if stored tokens are still valid before authenticating
- **Automatic Renewal**: Tokens are automatically renewed when they expire (assumed 24-hour expiry)
- **Seamless Integration**: You only need to use your own JWT token - the BioTime token is managed automatically
- **Tenant Isolation**: Each tenant has its own BioTime configuration and token

### Error Handling
- Connection errors are properly handled and returned with descriptive messages
- Authentication failures are clearly identified
- Network timeouts are set to 30 seconds

### Integration Status
- The service can be enabled/disabled via `BIOTIME_ENABLED` environment variable
- When disabled, all operations return appropriate error messages

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Failed to authenticate with BioTime API",
  "data": {
    "success": false,
    "message": "Failed to authenticate with BioTime API",
    "base_url": "http://10.4.9.24:8087",
    "authenticated": false,
    "error": "authentication failed with status 401: Invalid credentials"
  }
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "User not authenticated",
  "error": "Unauthorized access"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "message": "HR or Admin access required",
  "error": "Forbidden"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Failed to test BioTime connection",
  "error": "failed to connect to BioTime API: connection refused"
}
```

---

## Notes

1. **Token Management**:
   - BioTime tokens are **automatically stored in the database** when authenticated
   - Tokens are **automatically retrieved** from the database when needed
   - You **don't need to manually manage** BioTime tokens - just use your own JWT token
   - The service checks token validity and automatically renews expired tokens
   - Each tenant has its own BioTime configuration and token storage

2. **Database Storage**:
   - Configuration and tokens are stored in the `biotime_configs` table
   - First authentication creates the config record automatically
   - Subsequent requests use the stored token if still valid
   - Token expiry is tracked and tokens are renewed automatically

3. **Security**: 
   - BioTime credentials are stored in the database (encrypt in production)
   - Environment variables serve as initial/default configuration
   - Never commit `.env` file to version control
   - Use strong passwords for BioTime API access

2. **Network Access**:
   - Ensure the server can reach the BioTime device IP address (`10.4.9.24:8087`)
   - Check firewall rules if connection fails
   - Verify the BioTime device is powered on and accessible

3. **Token Caching**:
   - Tokens are cached in memory for the service instance
   - If the service restarts, tokens will be re-authenticated
   - Token expiry is assumed to be 24 hours (adjust if BioTime uses different expiry)

4. **Future Enhancements**:
   - The service includes a `MakeRequest` method for making authenticated API calls to other BioTime endpoints
   - This can be extended to fetch attendance data, employee lists, etc.

5. **Testing**:
   - Use the test connection endpoint to verify configuration
   - Check the response to ensure authentication is working
   - Monitor for connection errors in production

---

## Example Usage

### Testing Connection
```bash
# Test the BioTime connection
curl -X GET "http://localhost:8080/api/v1/biometric/biotime/test-connection" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Getting Token for Direct API Calls
```bash
# Get the authentication token
curl -X GET "http://localhost:8080/api/v1/biometric/biotime/token" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Troubleshooting

### Connection Failed
- Verify `BIOTIME_BASE_URL` is correct
- Check if the BioTime device is accessible from the server
- Verify network connectivity and firewall rules

### Authentication Failed
- Verify `BIOTIME_USERNAME` and `BIOTIME_PASSWORD` are correct
- Check if credentials have changed on the BioTime device
- Ensure the user has API access permissions

### Service Disabled
- Check `BIOTIME_ENABLED` environment variable
- Set to `true` to enable the integration

---

## How Token Management Works

### ⚠️ **IMPORTANT: Token Expiration is Handled Automatically!**

**You don't need to do anything when the token expires!** The system automatically handles it:

### Automatic Token Expiration Handling

1. **Token Validation Check**:
   - Every time you call a BioTime API, the system checks if the stored token is still valid
   - A token is considered valid if it expires in more than **5 minutes** (safety buffer)
   - This prevents using tokens that are about to expire

2. **Automatic Renewal**:
   - If the token is **expired or about to expire**, the system automatically:
     - ✅ Authenticates with BioTime API to get a new token
     - ✅ Stores the new token in the database
     - ✅ Uses the new token for your request
   - This happens **transparently** - you don't see any errors or need to do anything
   - **Your request succeeds normally!**

3. **What You Need to Do**:
   - **NOTHING!** Just use your API endpoints normally
   - The BioTime token is managed completely automatically
   - No manual refresh needed
   - No error handling for expired tokens needed

### Automatic Token Flow

1. **First Request**:
   - Service checks database for existing BioTime config
   - If not found, creates config from environment variables
   - Authenticates with BioTime API
   - Stores token and expiry in database
   - Returns token for use

2. **Subsequent Requests (Token Valid)**:
   - Service checks database for stored token
   - If token is valid (expires in > 5 minutes), returns stored token immediately
   - **No API call to BioTime needed - fast response!**

3. **Subsequent Requests (Token Expired)**:
   - Service checks database for stored token
   - **Detects token is expired or about to expire**
   - **Automatically authenticates with BioTime API**
   - **Stores new token in database**
   - **Uses new token for your request**
   - **You get your data successfully - no errors!**

4. **Your API Usage**:
   - You only need to provide **your own JWT token** in the Authorization header
   - The BioTime token is automatically retrieved and used internally
   - **No manual token management required**
   - **No need to handle token expiration** - it's all automatic!

### Example Scenarios

**Scenario 1: Token expires during normal operation**
```
You: GET /api/v1/biometric/biotime/terminals
System: "Token expired 10 minutes ago"
System: [Automatically authenticates] → [Gets new token] → [Stores it] → [Uses it]
You: ✅ Get your terminals data successfully
```

**Scenario 2: Token expires between requests**
```
Token expires at: 2:00 PM
You make request at: 2:15 PM
System: "Token expired, refreshing..."
System: [Automatically gets new token]
You: ✅ Request succeeds without any issues
```

**Scenario 3: Token is about to expire (within 5 minutes)**
```
Token expires at: 2:00 PM
Current time: 1:56 PM (4 minutes until expiry)
System: "Token expires soon, refreshing proactively..."
System: [Gets new token before expiry]
You: ✅ Prevents last-minute failures
```

### Database Storage

The `biotime_configs` table stores:
- `base_url`: BioTime API base URL
- `username`: BioTime API username
- `password`: BioTime API password (stored in database)
- `token`: Current JWT token from BioTime
- `token_expiry`: When the token expires
- `last_sync`: Last successful authentication
- `last_error`: Last error message (if any)
- `enabled`: Whether integration is enabled
- `tenant_id`: Tenant isolation

### Benefits

✅ **No Manual Token Management**: Tokens are automatically stored and retrieved  
✅ **Persistent Storage**: Tokens survive service restarts  
✅ **Automatic Renewal**: Expired tokens are automatically refreshed  
✅ **Tenant Isolation**: Each tenant has its own configuration  
✅ **Error Tracking**: Last errors are stored for debugging  
✅ **Seamless Integration**: Use your JWT token, BioTime token is handled automatically

---

## Next Steps

1. **Configure Environment Variables**:
   Add the BioTime configuration to your `.env` file:
   ```env
   BIOTIME_BASE_URL=http://10.4.9.24:8087
   BIOTIME_USERNAME=Developer
   BIOTIME_PASSWORD=Developer@123
   BIOTIME_ENABLED=true
   ```

2. **Restart Application**:
   The `biotime_configs` table will be created automatically on startup.

3. **Test Connection**:
   Use the test connection endpoint to verify the integration is working:
   ```bash
   GET /api/v1/biometric/biotime/test-connection
   ```
   This will create the config record and authenticate automatically.

4. **Use in Your Code**:
   When calling BioTime APIs, simply use:
   ```go
   service := services.NewBioTimeService(tenantID)
   token, err := service.Authenticate() // Automatically uses stored token or authenticates
   // Use token in BioTime API calls
   ```

5. **Extend Functionality**:
   The service can be extended to:
   - Fetch employee attendance records
   - Sync employee data with BioTime
   - Retrieve device status
   - Manage employee enrollments
   - All using the `MakeRequest()` method with automatic token management
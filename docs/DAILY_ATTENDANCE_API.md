# Daily Attendance API

## Overview
This API retrieves daily attendance records from the `biotime_transactions` database table. It groups transactions by employee and date, calculating check-in, check-out, and working hours for each day.

## Base URL
All API endpoints use the base URL: `{{BASE_URL}}/api/v1/biometric`

## Authentication
All requests require Bearer token authentication and HR or Admin role:
```
Authorization: Bearer <access_token>
```

---

## Get Daily Attendance

### Endpoint
```
GET /api/v1/biometric/attendance/daily
```

### Description
Retrieves daily attendance records with check-in, check-out, and working hours calculated from the `biotime_transactions` table. Data is grouped by employee and date.

### Query Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `start_time` | string | Yes | Start time filter. Supports: `YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS` |
| `end_time` | string | Yes | End time filter. Supports: `YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS` |
| `emp_code` | string | No | Filter by employee code |
| `page` | integer | No | Page number (default: 1) |
| `page_size` | integer | No | Page size (default: 50, max: 100) |

**Note:** If `end_time` is provided as date only (YYYY-MM-DD), it automatically sets to end of day (23:59:59).

### Supported DateTime Formats
- `YYYY-MM-DD` - Date only
- `YYYY-MM-DD HH:MM:SS` - Full datetime
- `YYYY-MM-DD HH:MM` - Datetime without seconds
- `YYYY-MM-DDTHH:MM:SS` - ISO format
- RFC3339 format

### Response Format

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "Daily attendance retrieved successfully",
  "data": {
    "data": [
      {
        "name": "John Doe",
        "emp_code": "EMP001",
        "date": "2025-01-15",
        "checkin": "09:00:00",
        "checkout": "17:30:00",
        "working_hours": "08:30",
        "punch_count": 2
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 50,
      "total": 100,
      "total_pages": 2
    }
  }
}
```

#### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Full name (FirstName + LastName) |
| `emp_code` | string | Employee code |
| `date` | string | Date in YYYY-MM-DD format |
| `checkin` | string\|null | Check-in time (HH:MM:SS) - earliest punch of the day |
| `checkout` | string\|null | Check-out time (HH:MM:SS) - latest punch of the day |
| `working_hours` | string\|null | Working hours (HH:MM) - calculated from check-in to check-out |
| `punch_count` | integer | Number of punches for the day |

#### Sample Response Data Table

| Name | Employee Code | Date | Check-in | Check-out | Working Hours | Punch Count |
|------|---------------|------|----------|-----------|---------------|-------------|
| John Doe | EMP001 | 2025-01-15 | 09:00:00 | 17:30:00 | 08:30 | 2 |
| Jane Smith | EMP002 | 2025-01-15 | 08:45:00 | 18:00:00 | 09:15 | 2 |
| Bob Johnson | EMP003 | 2025-01-15 | 09:15:00 | 17:45:00 | 08:30 | 2 |
| Alice Williams | EMP004 | 2025-01-15 | 08:30:00 | null | null | 1 |
| Charlie Brown | EMP005 | 2025-01-16 | 09:00:00 | 17:00:00 | 08:00 | 2 |

### Error Responses

#### 400 Bad Request
```json
{
  "success": false,
  "message": "start_time is required (format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)",
  "data": null
}
```

#### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized",
  "data": null
}
```

#### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Failed to get daily attendance",
  "data": "error details"
}
```

---

## Examples

### Example 1: Get attendance for a date range
```bash
curl -X GET "http://localhost:8080/api/v1/biometric/attendance/daily?start_time=2025-01-01&end_time=2025-01-31" \
  -H "Authorization: Bearer <token>"
```

### Example 2: Get attendance with datetime range
```bash
curl -X GET "http://localhost:8080/api/v1/biometric/attendance/daily?start_time=2025-01-15 08:00:00&end_time=2025-01-15 18:00:00" \
  -H "Authorization: Bearer <token>"
```

### Example 3: Get attendance for specific employee
```bash
curl -X GET "http://localhost:8080/api/v1/biometric/attendance/daily?start_time=2025-01-01&end_time=2025-01-31&emp_code=EMP001" \
  -H "Authorization: Bearer <token>"
```

### Example 4: Get attendance with pagination
```bash
curl -X GET "http://localhost:8080/api/v1/biometric/attendance/daily?start_time=2025-01-01&end_time=2025-01-31&page=1&page_size=20" \
  -H "Authorization: Bearer <token>"
```

### JavaScript Example
```javascript
const response = await fetch('http://localhost:8080/api/v1/biometric/attendance/daily?start_time=2025-01-01&end_time=2025-01-31', {
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

## Notes

1. **Data Source**: This API reads from the `biotime_transactions` database table, not directly from the BioTime API.

2. **Grouping Logic**: 
   - Transactions are grouped by employee (`emp_code`) and date
   - Check-in is the earliest `punch_time` for the day
   - Check-out is the latest `punch_time` for the day
   - Working hours = Check-out - Check-in

3. **Backward Compatibility**: The API still accepts legacy parameters `start_date` and `end_date` for backward compatibility.

4. **Tenant Isolation**: Results are automatically filtered by the authenticated user's tenant ID.

5. **Date-Only Handling**: If `end_time` is provided as date only (YYYY-MM-DD), it automatically includes the entire day (up to 23:59:59).

---

## Related APIs

- [BioTime API](./BIOTIME_API.md) - BioTime device integration APIs
- [BioTime Transaction Sync](./BIOTIME_TRANSACTION_SYNC.md) - Transaction synchronization documentation

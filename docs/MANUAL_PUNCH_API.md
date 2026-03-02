# Manual Punch API Documentation

## 📋 Overview
The Manual Punch API allows employees to request manual attendance entries when biometric/fingerprint systems fail. This includes creating punch requests, viewing status, and administrative approval workflows.

## 🔐 Authentication
All endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## 📊 Base URL
```
http://localhost:8080/api/v1
```

## 🚀 API Endpoints

### 1. Create Manual Punch Request
**POST** `/attendance/manual-punches`

Create a new manual punch request for when biometric systems fail.

#### Request Body
```typescript
{
  "punch_date": "2024-01-15",       // Date of the punch (YYYY-MM-DD)
  "punch_time": "09:30:00",         // Time of the punch (HH:MM:SS)
  "punch_type": "in",               // "in" or "out"
  "reason": "Fingerprint scanner not working. Had to manually record my arrival time.",  // Min 10 chars, max 500
  "supporting_document": "https://example.com/document.pdf"  // Optional document URL
}
```

#### Response (201 Created)
```typescript
{
  "success": true,
  "message": "Manual punch request submitted successfully",
  "data": {
    "id": 123,
    "employee_id": 456,
    "employee_name": "John Doe",
    "punch_date": "2024-01-15",
    "punch_time": "09:30:00",
    "punch_type": "in",
    "status": "pending",
    "reason": "Fingerprint scanner not working...",
    "supporting_document": "https://example.com/document.pdf",
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

#### Error Responses
- `400 Bad Request`: Invalid date/time format, future dates, or validation errors
- `400 Bad Request`: Too many pending requests (max 5 per employee)
- `401 Unauthorized`: Missing or invalid authentication

---

### 2. Get My Manual Punch Requests
**GET** `/attendance/manual-punches/my`

Retrieve all manual punch requests for the authenticated employee.

#### Query Parameters
| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `status` | string | Filter by status: `pending`, `approved`, `rejected`, `cancelled` | No |
| `start_date` | string | Filter from date (YYYY-MM-DD) | No |
| `end_date` | string | Filter to date (YYYY-MM-DD) | No |

#### Response (200 OK)
```typescript
{
  "success": true,
  "message": "Manual punches retrieved successfully",
  "data": [
    {
      "id": 123,
      "employee_id": 456,
      "punch_date": "2024-01-15",
      "punch_time": "09:30:00",
      "punch_type": "in",
      "status": "approved",
      "reason": "Fingerprint scanner issue",
      "supporting_document": null,
      "approver_id": 789,
      "approver_name": "Jane Smith (HR)",
      "approved_at": "2024-01-15T14:30:00Z",
      "rejection_reason": null,
      "created_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

---

### 3. Get Specific Manual Punch
**GET** `/attendance/manual-punches/{id}`

Retrieve details of a specific manual punch request.

#### Path Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Manual punch request ID |

#### Response (200 OK)
Same structure as "Get My Manual Punch Requests" but for a single item.

#### Error Responses
- `404 Not Found`: Manual punch request not found
- `403 Forbidden`: Not authorized to view this request

---

### 4. Cancel Manual Punch Request
**POST** `/attendance/manual-punches/{id}/cancel`

Cancel a pending manual punch request (only for own requests).

#### Path Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Manual punch request ID |

#### Response (200 OK)
```typescript
{
  "success": true,
  "message": "Manual punch cancelled successfully",
  "data": null
}
```

#### Error Responses
- `400 Bad Request`: Cannot cancel non-pending requests
- `403 Forbidden`: Can only cancel own requests
- `404 Not Found`: Request not found

---

### 5. Get All Manual Punches (Admin/HR)
**GET** `/attendance/manual-punches`

Retrieve all manual punch requests with filtering and pagination (Admin/HR only).

#### Query Parameters
| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `employee_id` | integer | Filter by specific employee | No |
| `status` | string | Filter by status | No |
| `punch_type` | string | Filter by type: `in` or `out` | No |
| `start_date` | string | Filter from date (YYYY-MM-DD) | No |
| `end_date` | string | Filter to date (YYYY-MM-DD) | No |
| `page` | integer | Page number (default: 1) | No |
| `limit` | integer | Items per page (default: 20, max: 100) | No |

#### Response (200 OK)
```typescript
{
  "success": true,
  "message": "Manual punches retrieved successfully",
  "data": [
    // Array of manual punch objects
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

---

### 6. Get Pending Manual Punches (Admin/HR)
**GET** `/attendance/manual-punches/pending`

Retrieve all pending manual punch requests awaiting approval (Admin/HR only).

#### Response (200 OK)
Same structure as "Get All Manual Punches" but only pending requests.

---

### 7. Approve/Reject Manual Punch (Admin/HR)
**PATCH** `/attendance/manual-punches/{id}/status`

Approve or reject a manual punch request (Admin/HR only).

#### Path Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Manual punch request ID |

#### Request Body
```typescript
{
  "status": "approved"  // or "rejected"
  "rejection_reason": "Insufficient documentation"  // Required if status is "rejected"
}
```

#### Response (200 OK)
```typescript
{
  "success": true,
  "message": "Manual punch status updated successfully",
  "data": {
    // Updated manual punch object with approver details
  }
}
```

#### Error Responses
- `400 Bad Request`: Missing rejection reason for rejected status
- `400 Bad Request": Request already processed
- `404 Not Found": Request not found

---

### 8. Delete Manual Punch (Admin/HR)
**DELETE** `/attendance/manual-punches/{id}`

Delete a manual punch request (soft delete, Admin/HR only).

#### Path Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Manual punch request ID |

#### Response (200 OK)
```typescript
{
  "success": true,
  "message": "Manual punch deleted successfully",
  "data": null
}
```

---

## 📋 Status Types

| Status | Description |
|--------|-------------|
| `pending` | Request submitted, awaiting approval |
| `approved` | Request approved by HR/Admin |
| `rejected` | Request rejected by HR/Admin |
| `cancelled` | Request cancelled by employee |

## ⚡ Punch Types

| Type | Description |
|------|-------------|
| `in` | Clock-in/arrival punch |
| `out` | Clock-out/departure punch |

## 🔄 Integration Flow

### Employee Flow:
1. Employee attempts biometric punch (fails)
2. Employee submits manual punch request via mobile app/web
3. System validates request (prevents future dates, limits pending requests)
4. Request appears in "My Requests" with pending status
5. HR receives notification for approval
6. Employee receives notification when approved/rejected

### HR Admin Flow:
1. View all pending requests in dashboard
2. Review each request with employee details and reason
3. Approve or reject with optional comments
4. System updates attendance records accordingly

## 🛡️ Validation Rules

- **Date Validation**: Cannot be future date, max 30 days in past
- **Time Validation**: Cannot be more than 24 hours in future
- **Pending Limit**: Max 5 pending requests per employee
- **Status Changes**: Only pending requests can be modified
- **Ownership**: Employees can only access their own requests

## 📱 Frontend Integration Example

```typescript
// Create manual punch request
const createManualPunch = async (punchData: {
  punch_date: string;
  punch_time: string;
  punch_type: 'in' | 'out';
  reason: string;
  supporting_document?: string;
}) => {
  const response = await fetch('/api/v1/attendance/manual-punches', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(punchData),
  });
  
  return await response.json();
};

// Get employee's manual punches
const getMyManualPunches = async (filters?: {
  status?: string;
  start_date?: string;
  end_date?: string;
}) => {
  const params = new URLSearchParams(filters);
  const response = await fetch(`/api/v1/attendance/manual-punches/my?${params}`, {
    headers: { 'Authorization': `Bearer ${token}` },
  });
  
  return await response.json();
};
```

## 🔔 Notifications

Frontend should implement:
- Success/error toasts for API responses
- Real-time updates for status changes (WebSocket recommended)
- Push notifications for approval decisions
- Form validation for required fields and date formats

## 🚨 Error Handling

Handle these common error scenarios:
- `401 Unauthorized`: Redirect to login
- `403 Forbidden`: Show access denied message
- `400 Bad Request`: Display validation errors to user
- `429 Too Many Requests`: Rate limiting - retry later

## 📊 Performance Considerations

- Implement client-side caching for frequently accessed data
- Use pagination for large result sets
- Debounce search/filter inputs
- Lazy load detailed views

This API provides comprehensive manual attendance management with proper security, validation, and administrative controls.
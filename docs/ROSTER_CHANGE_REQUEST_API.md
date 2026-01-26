# Roster Change Request API

## Overview
This API allows employees to request changes to their roster assignments (shift, date, or location). HR/Admin can then approve or reject these requests.

## Base URL
All endpoints use: `{{BASE_URL}}/api/v1`

## Authentication
All endpoints require Bearer token authentication:
```
Authorization: Bearer {access_token}
```

**Note:**
- Employee endpoints (`/self-service/rosters/change-requests`) require authentication only
- HR/Admin endpoints (`/rosters/change-requests`) require HR or Admin role

---

## Employee Endpoints

### 1. View My Roster Assignments
**Endpoint:** `GET {{BASE_URL}}/api/v1/self-service/rosters/assignments`

**Description:** Employees can view their own roster assignments. The response includes information about whether they can request changes.

**Authentication:** Required (Employee)

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 20, max: 100) - Items per page
- `start_date` (optional) - Filter by start date (YYYY-MM-DD)
- `end_date` (optional) - Filter by end date (YYYY-MM-DD)
- `status` (optional) - Filter by status: `draft`, `published`, `swapped`

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/self-service/rosters/assignments?start_date=2026-01-01&end_date=2026-01-31" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Roster assignments retrieved successfully",
  "data": {
    "data": [
      {
        "id": 5,
        "employee_id": "EMP464350",
        "employee_name": "John Doe",
        "date": "2026-01-24",
        "shift_id": 1,
        "shift_name": "Morning Shift",
        "shift_code": "MORN",
        "shift_timing": "09:00:00 - 18:00:00",
        "shift_type": "fixed",
        "location_id": 1,
        "location_name": "Main Office",
        "status": "published",
        "can_request_change": true,
        "has_pending_change_request": false,
        "employee_photo": "https://example.com/photo.jpg",
        "department": "IT",
        "created_at": "2026-01-20T10:00:00+03:00",
        "updated_at": "2026-01-20T10:00:00+03:00"
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 20,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

**Response Fields:**
- `can_request_change` - `true` if employee can request a change (no pending request and not swapped)
- `has_pending_change_request` - `true` if there's already a pending change request for this assignment

---

### 2. Get My Roster Assignment by ID
**Endpoint:** `GET {{BASE_URL}}/api/v1/self-service/rosters/assignments/:assignment_id`

**Description:** Get details of a specific roster assignment belonging to the employee.

**Authentication:** Required (Employee)

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/self-service/rosters/assignments/5" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Roster assignment retrieved successfully",
  "data": {
    "id": 5,
    "employee_id": "EMP464350",
    "employee_name": "John Doe",
    "date": "2026-01-24",
    "shift_id": 1,
    "shift_name": "Morning Shift",
    "shift_code": "MORN",
    "shift_timing": "09:00:00 - 18:00:00",
    "shift_type": "fixed",
    "location_id": 1,
    "location_name": "Main Office",
    "status": "published",
    "can_request_change": true,
    "has_pending_change_request": false,
    "employee_photo": "https://example.com/photo.jpg",
    "department": "IT",
    "created_at": "2026-01-20T10:00:00+03:00",
    "updated_at": "2026-01-20T10:00:00+03:00"
  }
}
```

**Error Response (if assignment doesn't belong to employee):**
```json
{
  "success": false,
  "message": "You can only view your own roster assignments"
}
```

---

### 3. Create Roster Change Request
**Endpoint:** `POST {{BASE_URL}}/api/v1/self-service/rosters/change-requests`

**Description:** Employees can request to change their roster assignment (shift, date, or location).

**Authentication:** Required (Employee)

**Request Body:**
```json
{
  "assignment_id": 5,
  "requested_shift_id": 2,
  "requested_date": "2026-01-25",
  "requested_location_id": 1,
  "reason": "Need to attend training on 2026-01-24"
}
```

**Request Fields:**
- `assignment_id` (required) - The roster assignment ID to change
- `requested_shift_id` (optional) - New shift ID
- `requested_date` (optional) - New date in YYYY-MM-DD format
- `requested_location_id` (optional) - New location ID
- `reason` (optional) - Reason for the change request

**Note:** At least one of `requested_shift_id`, `requested_date`, or `requested_location_id` must be provided.

**cURL Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/self-service/rosters/change-requests" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "assignment_id": 5,
    "requested_shift_id": 2,
    "requested_date": "2026-01-25",
    "reason": "Need to attend training"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Roster change request created successfully",
  "data": {
    "id": 1,
    "requested_by": "EMP464350",
    "requested_by_name": "John Doe",
    "assignment_id": 5,
    "assignment_date": "2026-01-24",
    "current_shift_id": 1,
    "current_shift_name": "Morning Shift",
    "current_shift_timing": "09:00:00 - 18:00:00",
    "requested_shift_id": 2,
    "requested_shift_name": "Afternoon Shift",
    "requested_shift_timing": "14:00:00 - 23:00:00",
    "requested_date": "2026-01-25",
    "reason": "Need to attend training",
    "status": "pending",
    "requested_at": "2026-01-22T10:00:00+03:00"
  }
}
```

**Error Responses:**

**If assignment doesn't belong to employee:**
```json
{
  "success": false,
  "message": "you can only request changes for your own assignments"
}
```

**If already has pending request:**
```json
{
  "success": false,
  "message": "you already have a pending change request for this assignment"
}
```

---

### 2. List My Change Requests
**Endpoint:** `GET {{BASE_URL}}/api/v1/self-service/rosters/change-requests`

**Description:** Employees can view their own roster change requests.

**Authentication:** Required (Employee)

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 20, max: 100) - Items per page
- `status` (optional) - Filter by status: `pending`, `approved`, `rejected`

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/self-service/rosters/change-requests?status=pending" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Roster change requests retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "requested_by": "EMP464350",
        "requested_by_name": "John Doe",
        "assignment_id": 5,
        "assignment_date": "2026-01-24",
        "current_shift_name": "Morning Shift",
        "requested_shift_id": 2,
        "requested_shift_name": "Afternoon Shift",
        "requested_date": "2026-01-25",
        "reason": "Need to attend training",
        "status": "pending",
        "requested_at": "2026-01-22T10:00:00+03:00"
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 20,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

---

### 3. Get My Change Request by ID
**Endpoint:** `GET {{BASE_URL}}/api/v1/self-service/rosters/change-requests/:request_id`

**Description:** Get details of a specific change request.

**Authentication:** Required (Employee)

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/self-service/rosters/change-requests/1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## HR/Admin Endpoints

### 4. List All Change Requests
**Endpoint:** `GET {{BASE_URL}}/api/v1/rosters/change-requests`

**Description:** HR/Admin can view all roster change requests.

**Authentication:** Required (HR/Admin)

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 20, max: 100) - Items per page
- `status` (optional) - Filter by status: `pending`, `approved`, `rejected`
- `requested_by` (optional) - Filter by employee ID
- `assignment_id` (optional) - Filter by assignment ID

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/change-requests?status=pending" \
  -H "Authorization: Bearer HR_ACCESS_TOKEN"
```

---

### 5. Get Change Request by ID
**Endpoint:** `GET {{BASE_URL}}/api/v1/rosters/change-requests/:request_id`

**Description:** Get details of a specific change request.

**Authentication:** Required (HR/Admin)

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/change-requests/1" \
  -H "Authorization: Bearer HR_ACCESS_TOKEN"
```

---

### 6. Approve Change Request
**Endpoint:** `POST {{BASE_URL}}/api/v1/rosters/change-requests/:request_id/approve`

**Description:** HR/Admin approves a change request and updates the roster assignment.

**Authentication:** Required (HR/Admin)

**Request Body:** (Optional)
```json
{
  "notify_employee": true,
  "notes": "Change approved. Please confirm your new schedule."
}
```

**cURL Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/change-requests/1/approve" \
  -H "Authorization: Bearer HR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "notify_employee": true,
    "notes": "Change approved"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Roster change request approved successfully",
  "data": {
    "id": 1,
    "status": "approved",
    "reviewed_at": "2026-01-22T15:00:00+03:00",
    "reviewed_by": 2,
    "updated_assignment": {
      "assignment_id": 5,
      "employee_id": "EMP464350",
      "date": "2026-01-25",
      "shift_id": 2,
      "shift_name": "Afternoon Shift",
      "location_id": 1,
      "location_name": "Main Office",
      "status": "published"
    }
  }
}
```

**Error Response (if request not pending):**
```json
{
  "success": false,
  "message": "only pending change requests can be approved"
}
```

---

### 7. Reject Change Request
**Endpoint:** `POST {{BASE_URL}}/api/v1/rosters/change-requests/:request_id/reject`

**Description:** HR/Admin rejects a change request.

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "rejection_reason": "Insufficient notice period. Please request at least 3 days in advance.",
  "notify_employee": true
}
```

**cURL Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/change-requests/1/reject" \
  -H "Authorization: Bearer HR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rejection_reason": "Insufficient notice period",
    "notify_employee": true
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Roster change request rejected successfully",
  "data": {
    "id": 1,
    "status": "rejected",
    "rejection_reason": "Insufficient notice period",
    "reviewed_at": "2026-01-22T15:30:00+03:00",
    "reviewed_by": 2
  }
}
```

---

## Complete Workflow Example

### Step 1: Employee Views Their Roster Assignments
```bash
curl -X GET "http://localhost:8080/api/v1/self-service/rosters/assignments" \
  -H "Authorization: Bearer EMPLOYEE_TOKEN"
```

**Response includes `can_request_change: true` for assignments that can be changed.**

### Step 2: Employee Creates Change Request
```bash
curl -X POST "http://localhost:8080/api/v1/self-service/rosters/change-requests" \
  -H "Authorization: Bearer EMPLOYEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "assignment_id": 5,
    "requested_shift_id": 2,
    "requested_date": "2026-01-25",
    "reason": "Need to attend training session"
  }'
```

### Step 2: HR Reviews and Approves
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/change-requests/1/approve" \
  -H "Authorization: Bearer HR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "notify_employee": true,
    "notes": "Change approved. Please confirm your new schedule."
  }'
```

### Step 3: Employee Views Updated Assignment
The assignment is automatically updated when approved. Employee can view it via the roster assignments endpoint.

---

## Status Values

- `pending` - Awaiting HR/Admin review
- `approved` - Approved and assignment updated
- `rejected` - Rejected by HR/Admin

---

## Notes

1. **Employee Ownership:** Employees can only create change requests for their own assignments
2. **Pending Requests:** Only one pending request per assignment is allowed
3. **Automatic Update:** When approved, the assignment is automatically updated with the requested changes
4. **Notifications:** Employee notifications are planned for future implementation
5. **Validation:** The system validates that requested shift exists and assignment belongs to the employee

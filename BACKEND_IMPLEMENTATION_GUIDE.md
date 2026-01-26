# Shifts & Rosters API Documentation

This document outlines the API endpoints for the Shifts & Rosters module, including shift management, roster assignments, swap requests, and roster change requests.

**Related Documentation:**
- [Roster Change Request API](./ROSTER_CHANGE_REQUEST_API.md) - Employee requests to change roster assignments

---

## Base URL
All endpoints use: `{{BASE_URL}}/api/v1/shifts` or `{{BASE_URL}}/api/v1/rosters`

**Example Base URL:** `http://localhost:8080/api/v1`

---

## Quick Reference

### All Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/shifts/statistics` | Get statistics |
| GET | `/api/v1/shifts` | List all shifts |
| GET | `/api/v1/shifts/:shift_id` | Get shift by ID |
| POST | `/api/v1/shifts` | Create shift |
| PUT | `/api/v1/shifts/:shift_id` | Update shift |
| POST | `/api/v1/shifts/:shift_id/duplicate` | Duplicate shift |
| DELETE | `/api/v1/shifts/:shift_id` | Delete shift |
| GET | `/api/v1/rosters/assignments` | List roster assignments |
| GET | `/api/v1/rosters/assignments/:assignment_id` | Get assignment by ID |
| POST | `/api/v1/rosters/assignments` | Create assignment |
| POST | `/api/v1/rosters/assignments/bulk` | Bulk create assignments |
| PUT | `/api/v1/rosters/assignments/:assignment_id` | Update assignment |
| DELETE | `/api/v1/rosters/assignments/:assignment_id` | Delete assignment |
| GET | `/api/v1/rosters/weekly` | Get weekly roster view |
| POST | `/api/v1/rosters/auto-schedule` | Auto-schedule rosters |
| POST | `/api/v1/rosters/publish` | Publish rosters |
| GET | `/api/v1/rosters/swap-requests` | List swap requests |
| GET | `/api/v1/rosters/swap-requests/:request_id` | Get swap request by ID |
| POST | `/api/v1/rosters/swap-requests` | Create swap request |
| POST | `/api/v1/rosters/swap-requests/:request_id/approve` | Approve swap request |
| POST | `/api/v1/rosters/swap-requests/:request_id/reject` | Reject swap request |
| GET | `/api/v1/self-service/rosters/assignments` | List my roster assignments (Employee) |
| GET | `/api/v1/self-service/rosters/assignments/:assignment_id` | Get my roster assignment (Employee) |
| POST | `/api/v1/self-service/rosters/change-requests` | Create change request (Employee) |
| GET | `/api/v1/self-service/rosters/change-requests` | List my change requests (Employee) |
| GET | `/api/v1/rosters/change-requests` | List all change requests (HR/Admin) |
| GET | `/api/v1/rosters/change-requests/:request_id` | Get change request |
| POST | `/api/v1/rosters/change-requests/:request_id/approve` | Approve change request |
| POST | `/api/v1/rosters/change-requests/:request_id/reject` | Reject change request |

---

## Authentication

All endpoints require authentication via Bearer token in the Authorization header. Additionally, all endpoints require HR or Admin role permissions.

```
Authorization: Bearer {access_token}
```

**Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/shifts" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 1. Statistics & Dashboard

### 1.1 Get Shifts & Rosters Statistics
**Endpoint:** `GET {{BASE_URL}}/api/v1/shifts/statistics`

**Description:** Retrieves summary statistics for shifts and rosters, including total shifts, active/inactive shifts, roster counts, and swap request statistics.

**Query Parameters:** None

**Authentication:** Required (HR/Admin)

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/shifts/statistics" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Statistics retrieved successfully",
  "data": {
    "total_shifts": 12,
    "active_shifts": 8,
    "inactive_shifts": 4,
    "total_rostered": 156,
    "published_rosters": 142,
    "draft_rosters": 14,
    "pending_swap_requests": 2,
    "approved_swap_requests": 5,
    "rejected_swap_requests": 1
  }
}
```

---

## 2. Shift Management

### 2.1 Get All Shifts
**Endpoint:** `GET {{BASE_URL}}/api/v1/shifts`

**Description:** Retrieves a list of all shifts with pagination and filtering. Results are sorted by creation date (newest first).

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 20, max: 100) - Items per page
- `status` (optional) - Filter by status: `active`, `inactive`, or `all`
- `shift_type` (optional) - Filter by shift type: `fixed`, `split`, `rotating`, `night`, `flexible`
- `search` (optional) - Search by shift name or code (case-insensitive)

**Authentication:** Required (HR/Admin)

**cURL Examples:**

**Get all shifts (first page):**
```bash
curl -X GET "http://localhost:8080/api/v1/shifts" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get active shifts only:**
```bash
curl -X GET "http://localhost:8080/api/v1/shifts?status=active" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Search for shifts:**
```bash
curl -X GET "http://localhost:8080/api/v1/shifts?search=morning&page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get night shifts:**
```bash
curl -X GET "http://localhost:8080/api/v1/shifts?shift_type=night" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Shifts retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "shift_name": "Morning Shift",
        "shift_code": "MRN",
        "shift_type": "Fixed",
        "start_time": "09:00:00",
        "end_time": "18:00:00",
        "working_hours": 9,
        "break_duration": 60,
        "grace_minutes": 10,
        "late_mark_after": 15,
        "early_going_minutes": 15,
        "half_day_hours": 4.5,
        "minimum_hours": 4,
        "cross_day": false,
        "night_shift": false,
        "weekly_off": ["Sunday"],
        "locations": [1, 2],
        "location_names": ["Main Office", "Branch Office"],
        "is_active": true,
        "created_at": "2026-01-15T10:30:00+03:00",
        "updated_at": "2026-01-20T14:20:00+03:00",
        "created_by": 2,
        "updated_by": 2
      }
    ],
  "meta": {
      "page": 1,
      "per_page": 20,
      "total": 12,
      "total_pages": 1
    }
  }
}
```

### 2.2 Get Shift by ID
**Endpoint:** `GET {{BASE_URL}}/api/v1/shifts/:shift_id`

**Description:** Retrieves detailed information about a specific shift, including all locations associated with the shift.

**Path Parameters:**
- `shift_id` (required) - Shift ID

**Authentication:** Required (HR/Admin)

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/shifts/1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Shift retrieved successfully",
  "data": {
    "id": 1,
    "shift_name": "Morning Shift",
    "shift_code": "MRN",
    "shift_type": "Fixed",
    "start_time": "09:00:00",
    "end_time": "18:00:00",
    "working_hours": 9,
    "break_duration": 60,
    "grace_minutes": 10,
    "late_mark_after": 15,
    "early_going_minutes": 15,
    "half_day_hours": 4.5,
    "minimum_hours": 4,
    "cross_day": false,
    "night_shift": false,
    "weekly_off": ["Sunday"],
    "locations": [
      {
        "id": 1,
        "name": "Main Office",
        "code": "MO-001"
      },
      {
        "id": 2,
        "name": "Branch Office",
        "code": "BO-002"
      }
    ],
    "is_active": true,
    "created_at": "2026-01-15T10:30:00+03:00",
    "updated_at": "2026-01-20T14:20:00+03:00",
    "created_by": 2,
    "updated_by": 2
  }
}
```

### 2.3 Create Shift
**Endpoint:** `POST {{BASE_URL}}/api/v1/shifts`

**Description:** Creates a new shift. Working hours are automatically calculated from start_time and end_time. If end_time is before start_time, the shift is treated as cross-day.

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "shift_name": "Morning Shift",
  "shift_code": "MRN",
  "shift_type": "Fixed",
  "start_time": "09:00:00",
  "end_time": "18:00:00",
  "break_duration": 60,
  "grace_minutes": 10,
  "late_mark_after": 15,
  "early_going_minutes": 15,
  "half_day_hours": 4.5,
  "minimum_hours": 4,
  "cross_day": false,
  "night_shift": false,
  "weekly_off": ["Sunday"],
  "location_ids": [1, 2],
  "is_active": true
}
```

**cURL Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/shifts" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shift_name": "Morning Shift",
    "shift_code": "MRN",
    "shift_type": "fixed",
    "start_time": "09:00:00",
    "end_time": "18:00:00",
    "break_duration": 60,
    "grace_minutes": 10,
    "late_mark_after": 15,
    "early_going_minutes": 15,
    "half_day_hours": 4.5,
    "minimum_hours": 4,
    "cross_day": false,
    "night_shift": false,
    "weekly_off": ["Sunday"],
    "location_ids": [1, 2],
    "is_active": true
  }'
```

**Example: Create Night Shift (Cross-Day):**
```bash
curl -X POST "http://localhost:8080/api/v1/shifts" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shift_name": "Night Shift",
    "shift_code": "NGT",
    "shift_type": "night",
    "start_time": "22:00:00",
    "end_time": "06:00:00",
    "break_duration": 30,
    "grace_minutes": 15,
    "late_mark_after": 20,
    "early_going_minutes": 20,
    "half_day_hours": 4,
    "minimum_hours": 4,
    "cross_day": true,
    "night_shift": true,
    "weekly_off": ["Saturday", "Sunday"],
    "location_ids": [1],
    "is_active": true
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Shift created successfully",
  "data": {
    "id": 13,
    "shift_name": "Morning Shift",
    "shift_code": "MRN",
    "shift_type": "Fixed",
    "start_time": "09:00:00",
    "end_time": "18:00:00",
    "working_hours": 9,
    "break_duration": 60,
    "grace_minutes": 10,
    "late_mark_after": 15,
    "early_going_minutes": 15,
    "half_day_hours": 4.5,
    "minimum_hours": 4,
    "cross_day": false,
    "night_shift": false,
    "weekly_off": ["Sunday"],
    "locations": [1, 2],
    "is_active": true,
    "created_at": "2026-01-22T10:30:00+03:00",
    "updated_at": "2026-01-22T10:30:00+03:00",
    "created_by": 2,
    "updated_by": 2
  }
}
```

### 2.4 Update Shift
**Endpoint:** `PUT {{BASE_URL}}/api/v1/shifts/:shift_id`

**Description:** Updates an existing shift. Only include fields you want to update. Working hours are automatically recalculated if start_time or end_time are updated.

**Path Parameters:**
- `shift_id` (required) - Shift ID

**Authentication:** Required (HR/Admin)

**Request Body:** (All fields optional, only include fields to update)

**cURL Example:**
```bash
curl -X PUT "http://localhost:8080/api/v1/shifts/1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shift_name": "Updated Morning Shift",
    "start_time": "08:00:00",
    "end_time": "17:00:00",
    "break_duration": 45,
    "grace_minutes": 15,
    "is_active": false
  }'
```
```json
{
  "shift_name": "Updated Morning Shift",
  "shift_code": "MRN-UPD",
  "start_time": "08:00:00",
  "end_time": "17:00:00",
  "break_duration": 45,
  "grace_minutes": 15,
  "is_active": false
}
```

**Response:**
```json
{
  "success": true,
  "message": "Shift updated successfully",
  "data": {
    "id": 1,
    "shift_name": "Updated Morning Shift",
    "shift_code": "MRN-UPD",
    "shift_type": "Fixed",
    "start_time": "08:00:00",
    "end_time": "17:00:00",
    "working_hours": 9,
    "break_duration": 45,
    "grace_minutes": 15,
    "late_mark_after": 15,
    "early_going_minutes": 15,
    "half_day_hours": 4.5,
    "minimum_hours": 4,
    "cross_day": false,
    "night_shift": false,
    "weekly_off": ["Sunday"],
    "locations": [1, 2],
    "is_active": false,
    "created_at": "2026-01-15T10:30:00+03:00",
    "updated_at": "2026-01-22T15:45:00+03:00",
    "created_by": 2,
    "updated_by": 2
  }
}
```

### 2.5 Duplicate Shift
**Endpoint:** `POST {{BASE_URL}}/api/v1/shifts/:shift_id/duplicate`

**Description:** Creates a copy of an existing shift. All settings are copied, but you can override the shift name and code. If not provided, the new shift will be named "{Original Name} (Copy)" and code will be "{Original Code}-COPY".

**Path Parameters:**
- `shift_id` (required) - Shift ID to duplicate

**Authentication:** Required (HR/Admin)

**Request Body:** (Optional - override specific fields)

**cURL Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/shifts/1/duplicate" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shift_name": "Copied Morning Shift",
    "shift_code": "MRN-COPY"
  }'
```

**Example: Duplicate without override (uses defaults):**
```bash
curl -X POST "http://localhost:8080/api/v1/shifts/1/duplicate" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}'
```
```json
{
  "shift_name": "Copied Morning Shift",
  "shift_code": "MRN-COPY"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Shift duplicated successfully",
  "data": {
    "id": 14,
    "shift_name": "Copied Morning Shift",
    "shift_code": "MRN-COPY",
    "shift_type": "Fixed",
    "start_time": "09:00:00",
    "end_time": "18:00:00",
    "working_hours": 9,
    "break_duration": 60,
    "grace_minutes": 10,
    "late_mark_after": 15,
    "early_going_minutes": 15,
    "half_day_hours": 4.5,
    "minimum_hours": 4,
    "cross_day": false,
    "night_shift": false,
    "weekly_off": ["Sunday"],
    "locations": [1, 2],
    "is_active": true,
    "created_at": "2026-01-22T16:00:00+03:00",
    "updated_at": "2026-01-22T16:00:00+03:00",
    "created_by": 2,
    "updated_by": 2
  }
}
```

### 2.6 Delete Shift
**Endpoint:** `DELETE {{BASE_URL}}/api/v1/shifts/:shift_id`

**Description:** Soft deletes a shift. Only allowed if shift is not assigned to any published rosters. If the shift has active (published) rosters, deletion will be rejected with a 400 error.

**Path Parameters:**
- `shift_id` (required) - Shift ID

**Authentication:** Required (HR/Admin)

**cURL Example:**
```bash
curl -X DELETE "http://localhost:8080/api/v1/shifts/1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Error Response (if shift has active rosters):**
```json
{
  "success": false,
  "message": "shift is assigned to active rosters and cannot be deleted"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Shift deleted successfully",
  "data": {
    "deleted_at": "2026-01-22T16:30:00+03:00",
    "shift_id": 1
  }
}
```

---

## 3. Roster Assignments

### 3.1 Get Roster Assignments
**Endpoint:** `GET {{BASE_URL}}/api/v1/rosters/assignments`

**Description:** Retrieves roster assignments with filtering and pagination. Results are sorted by date (newest first), then by employee ID.

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 20, max: 100) - Items per page
- `start_date` (optional) - Filter by start date (YYYY-MM-DD)
- `end_date` (optional) - Filter by end date (YYYY-MM-DD)
- `employee_id` (optional) - Filter by employee ID (e.g., "EMP464350")
- `shift_id` (optional) - Filter by shift ID
- `status` (optional) - Filter by status: `draft`, `published`, `swapped`
- `location_id` (optional) - Filter by location ID

**Authentication:** Required (HR/Admin)

**cURL Examples:**

**Get all roster assignments:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/assignments" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get assignments for a date range:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/assignments?start_date=2026-01-22&end_date=2026-01-28" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get assignments for a specific employee:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/assignments?employee_id=EMP464350" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get published rosters only:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/assignments?status=published&page=1&page_size=50" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get assignments for a specific shift:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/assignments?shift_id=1&start_date=2026-01-22&end_date=2026-01-28" \
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
        "id": 1,
        "employee_id": "EMP464350",
        "employee_name": "John Doe",
        "employee_photo": "https://example.com/photo.jpg",
        "department": "Finance Department",
        "date": "2026-01-22",
        "shift_id": 1,
        "shift_name": "Morning Shift",
        "shift_code": "MRN",
        "shift_timing": "09:00 - 18:00",
        "location_id": 1,
        "location_name": "Main Office",
        "status": "published",
        "swap_request": null,
        "created_at": "2026-01-20T10:00:00+03:00",
        "updated_at": "2026-01-20T10:00:00+03:00",
        "created_by": 2,
        "updated_by": 2
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 20,
      "total": 156,
      "total_pages": 8
    }
  }
}
```

### 3.2 Get Roster Assignment by ID
**Endpoint:** `GET {{BASE_URL}}/api/v1/rosters/assignments/:assignment_id`

**Description:** Retrieves detailed information about a specific roster assignment, including employee details, shift information, and location.

**Path Parameters:**
- `assignment_id` (required) - Roster assignment ID

**Authentication:** Required (HR/Admin)

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/assignments/1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Roster assignment retrieved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP464350",
    "employee_name": "John Doe",
    "employee_photo": "https://example.com/photo.jpg",
    "department": "Finance Department",
    "date": "2026-01-22",
    "shift_id": 1,
    "shift_name": "Morning Shift",
    "shift_code": "MRN",
    "shift_timing": "09:00 - 18:00",
    "location_id": 1,
    "location_name": "Main Office",
    "status": "Published",
    "swap_request": null,
    "created_at": "2026-01-20T10:00:00+03:00",
    "updated_at": "2026-01-20T10:00:00+03:00",
    "created_by": 2,
    "updated_by": 2
  }
}
```

### 3.3 Create Roster Assignment
**Endpoint:** `POST {{BASE_URL}}/api/v1/rosters/assignments`

**Description:** Creates a new roster assignment. The employee and shift must exist. If an assignment already exists for the employee on that date, the request will be rejected. Default status is "draft" if not specified.

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "employee_id": "EMP464350",
  "date": "2026-01-23",
  "shift_id": 1,
  "location_id": 1,
  "status": "draft"
}
```

**cURL Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/assignments" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "EMP464350",
    "date": "2026-01-23",
    "shift_id": 1,
    "location_id": 1,
    "status": "draft"
  }'
```

**Error Response (if assignment already exists):**
```json
{
  "success": false,
  "message": "roster assignment already exists for this employee and date"
}
```

**Error Response (if employee not found):**
```json
{
  "success": false,
  "message": "employee not found"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Roster assignment created successfully",
  "data": {
    "id": 157,
    "employee_id": "EMP464350",
    "employee_name": "John Doe",
    "employee_photo": "https://example.com/photo.jpg",
    "department": "Finance Department",
    "date": "2026-01-23",
    "shift_id": 1,
    "shift_name": "Morning Shift",
    "shift_code": "MRN",
    "shift_timing": "09:00 - 18:00",
    "location_id": 1,
    "location_name": "Main Office",
    "status": "draft",
    "swap_request": null,
    "created_at": "2026-01-22T16:00:00+03:00",
    "updated_at": "2026-01-22T16:00:00+03:00",
    "created_by": 2,
    "updated_by": 2
  }
}
```

### 3.4 Bulk Create Roster Assignments
**Endpoint:** `POST {{BASE_URL}}/api/v1/rosters/assignments/bulk`

**Description:** Creates multiple roster assignments in a single request. Invalid assignments (non-existent employees, shifts, or duplicate assignments) are skipped and counted in the "failed" field. Valid assignments are created and returned.

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "assignments": [
    {
      "employee_id": "EMP464350",
      "date": "2026-01-23",
      "shift_id": 1,
      "location_id": 1,
      "status": "draft"
    },
    {
      "employee_id": "EMP567886",
      "date": "2026-01-23",
      "shift_id": 2,
      "location_id": 1,
      "status": "draft"
    }
  ]
}
```

**cURL Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/assignments/bulk" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "assignments": [
      {
        "employee_id": "EMP464350",
        "date": "2026-01-23",
        "shift_id": 1,
        "location_id": 1,
        "status": "draft"
      },
      {
        "employee_id": "EMP567886",
        "date": "2026-01-23",
        "shift_id": 2,
        "location_id": 1,
        "status": "draft"
      }
    ]
  }'
```

**Example Response (with some failures):**
```json
{
  "success": true,
  "message": "Roster assignments created successfully",
  "data": {
    "created": 1,
    "failed": 1,
    "assignments": [
      {
        "id": 157,
        "employee_id": "EMP464350",
        "date": "2026-01-23",
        "shift_id": 1,
        "status": "draft"
      }
    ]
  }
}
```

**Response:**
```json
{
  "success": true,
  "message": "Roster assignments created successfully",
  "data": {
    "created": 2,
    "failed": 0,
    "assignments": [
      {
        "id": 157,
        "employee_id": "EMP464350",
        "date": "2026-01-23",
        "shift_id": 1,
        "status": "draft"
      },
      {
        "id": 158,
        "employee_id": "EMP567886",
        "date": "2026-01-23",
        "shift_id": 2,
        "status": "draft"
      }
    ]
  }
}
```

### 3.5 Update Roster Assignment
**Endpoint:** `PUT {{BASE_URL}}/api/v1/rosters/assignments/:assignment_id`

**Description:** Updates an existing roster assignment. Only include fields you want to update. If updating shift_id, the new shift must exist.

**Path Parameters:**
- `assignment_id` (required) - Roster assignment ID

**Authentication:** Required (HR/Admin)

**Request Body:** (All fields optional)

**cURL Example:**
```bash
curl -X PUT "http://localhost:8080/api/v1/rosters/assignments/1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shift_id": 2,
    "location_id": 2,
    "status": "published"
  }'
```
```json
{
  "shift_id": 2,
  "location_id": 2,
  "status": "published"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Roster assignment updated successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP464350",
    "employee_name": "John Doe",
    "date": "2026-01-22",
    "shift_id": 2,
    "shift_name": "Afternoon Shift",
    "shift_code": "AFT",
    "shift_timing": "14:00 - 23:00",
    "location_id": 2,
    "location_name": "Branch Office",
    "status": "published",
    "updated_at": "2026-01-22T17:00:00+03:00",
    "updated_by": 2
  }
}
```

### 3.6 Delete Roster Assignment
**Endpoint:** `DELETE {{BASE_URL}}/api/v1/rosters/assignments/:assignment_id`

**Description:** Soft deletes a roster assignment.

**Path Parameters:**
- `assignment_id` (required) - Roster assignment ID

**Authentication:** Required (HR/Admin)

**cURL Example:**
```bash
curl -X DELETE "http://localhost:8080/api/v1/rosters/assignments/1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Roster assignment deleted successfully",
  "data": {
    "deleted_at": "2026-01-22T17:30:00+03:00",
    "assignment_id": 1
  }
}
```

### 3.7 Get Weekly Roster View
**Endpoint:** `GET {{BASE_URL}}/api/v1/rosters/weekly`

**Description:** Retrieves roster assignments for a specific week in a grid format, grouped by employee. Each employee's assignments are organized by date for easy viewing in a calendar/grid interface.

**Query Parameters:**
- `week` (required) - Week in format: `YYYY-WW` (e.g., `2026-04`). Week 1 starts from the first Monday of the year.
- `employee_ids` (optional) - Comma-separated list of employee IDs to filter (e.g., `EMP464350,EMP567886`)
- `department_id` (optional) - Filter by department ID
- `location_id` (optional) - Filter by location ID

**Authentication:** Required (HR/Admin)

**cURL Examples:**

**Get weekly roster for a specific week:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/weekly?week=2026-04" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get weekly roster for specific employees:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/weekly?week=2026-04&employee_ids=EMP464350,EMP567886" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get weekly roster for a department:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/weekly?week=2026-04&department_id=3" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get weekly roster for a location:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/weekly?week=2026-04&location_id=1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Weekly roster retrieved successfully",
  "data": {
    "week": "2026-04",
    "start_date": "2026-01-22",
    "end_date": "2026-01-28",
    "employees": [
      {
        "employee_id": "EMP464350",
        "employee_name": "John Doe",
        "employee_photo": "https://example.com/photo.jpg",
        "department": "Finance Department",
        "assignments": {
          "2026-01-22": {
            "id": 1,
            "shift_id": 1,
            "shift_name": "Morning Shift",
            "shift_code": "MRN",
            "shift_timing": "09:00 - 18:00",
            "location_id": 1,
            "location_name": "Main Office",
            "status": "Published"
          },
          "2026-01-23": {
            "id": 2,
            "shift_id": 1,
            "shift_name": "Morning Shift",
            "shift_code": "MRN",
            "shift_timing": "09:00 - 18:00",
            "location_id": 1,
            "location_name": "Main Office",
            "status": "Published"
          },
          "2026-01-24": null,
          "2026-01-25": null,
          "2026-01-26": null,
          "2026-01-27": null,
          "2026-01-28": null
        }
      }
    ]
  }
}
```

### 3.8 Auto-Schedule Roster
**Endpoint:** `POST {{BASE_URL}}/api/v1/rosters/auto-schedule`

**Description:** Automatically generates roster assignments based on rules and employee availability. Uses a round-robin algorithm to distribute shifts among employees. If shift preferences are provided, they are prioritized. The auto-schedule creates assignments with "draft" status, which can later be reviewed and published.

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "start_date": "2026-01-22",
  "end_date": "2026-01-28",
  "employee_ids": ["EMP464350", "EMP567886"],
  "department_id": 3,
  "location_id": 1,
  "shift_preferences": {
    "EMP464350": [1, 2],
    "EMP567886": [2, 3]
  },
  "rules": {
    "max_consecutive_days": 5,
    "min_rest_days_per_week": 1,
    "prefer_weekend_off": true
  }
}
```

**Response:**
```json
{
  "success": true,
  "message": "Roster auto-scheduled successfully",
  "data": {
    "created_assignments": 12,
    "assignments": [
      {
        "id": 159,
        "employee_id": "EMP464350",
        "date": "2026-01-22",
        "shift_id": 1,
        "status": "draft"
      }
    ]
  }
}
```

### 3.9 Publish Roster
**Endpoint:** `POST {{BASE_URL}}/api/v1/rosters/publish`

**Description:** Publishes roster assignments for a date range. All draft assignments within the date range are updated to "published" status. Employee notifications are planned for future implementation.

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "start_date": "2026-01-22",
  "end_date": "2026-01-28",
  "notify_employees": true,
  "notification_message": "Your roster for the week has been published. Please review your schedule."
}
```

**cURL Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/publish" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2026-01-22",
    "end_date": "2026-01-28",
    "notify_employees": true,
    "notification_message": "Your roster for the week has been published. Please review your schedule."
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Roster published successfully",
  "data": {
    "published_count": 35,
    "notified_employees": 25,
    "published_at": "2026-01-22T18:00:00+03:00",
    "published_by": 2
  }
}
```

---

## 4. Swap Requests

### 4.1 Get Swap Requests
**Endpoint:** `GET {{BASE_URL}}/api/v1/rosters/swap-requests`

**Description:** Retrieves shift swap requests with filtering and pagination. Results are sorted by request date (newest first).

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 20, max: 100) - Items per page
- `status` (optional) - Filter by status: `pending`, `approved`, `rejected`
- `requested_by` (optional) - Filter by requester employee ID (e.g., "EMP464350")
- `requested_with` (optional) - Filter by swap partner employee ID

**Authentication:** Required (HR/Admin)

**cURL Examples:**

**Get all swap requests:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/swap-requests" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get pending swap requests:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/swap-requests?status=pending" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get swap requests by requester:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/swap-requests?requested_by=EMP464350" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Swap requests retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "requested_by": "EMP006",
        "requested_by_name": "James Anderson",
        "requested_by_photo": "https://example.com/photo.jpg",
        "requested_by_department": "Sales",
        "requested_with": "EMP010",
        "requested_with_name": "Maria Garcia",
        "requested_with_photo": "https://example.com/photo2.jpg",
        "requested_with_department": "Sales",
        "assignment_id": 5,
        "assignment_date": "2026-01-22",
        "assignment_shift": "Afternoon Shift",
        "swap_assignment_id": 8,
        "swap_assignment_date": "2026-01-25",
        "swap_assignment_shift": "Morning Shift",
        "reason": "Personal commitment",
        "status": "pending",
        "requested_at": "2026-01-20T14:30:00+03:00",
        "reviewed_at": null,
        "reviewed_by": null,
        "reviewer_name": null
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 20,
      "total": 2,
      "total_pages": 1
    }
  }
}
```

### 4.2 Get Swap Request by ID
**Endpoint:** `GET {{BASE_URL}}/api/v1/rosters/swap-requests/:request_id`

**Description:** Retrieves detailed information about a specific swap request, including both assignment details and swap assignment details.

**Path Parameters:**
- `request_id` (required) - Swap request ID

**Authentication:** Required (HR/Admin)

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/rosters/swap-requests/1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Swap request retrieved successfully",
  "data": {
    "id": 1,
    "requested_by": "EMP006",
    "requested_by_name": "James Anderson",
    "requested_by_photo": "https://example.com/photo.jpg",
    "requested_by_department": "Sales",
    "requested_with": "EMP010",
    "requested_with_name": "Maria Garcia",
    "requested_with_photo": "https://example.com/photo2.jpg",
    "requested_with_department": "Sales",
    "assignment_id": 5,
    "assignment_date": "2026-01-22",
    "assignment_shift": "Afternoon Shift",
    "assignment_shift_timing": "14:00 - 23:00",
    "swap_assignment_id": 8,
    "swap_assignment_date": "2026-01-25",
    "swap_assignment_shift": "Morning Shift",
    "swap_assignment_shift_timing": "09:00 - 18:00",
    "reason": "Personal commitment",
    "status": "pending",
    "requested_at": "2026-01-20T14:30:00+03:00",
    "reviewed_at": null,
    "reviewed_by": null,
    "reviewer_name": null
  }
}
```

### 4.3 Create Swap Request
**Endpoint:** `POST {{BASE_URL}}/api/v1/rosters/swap-requests`

**Description:** Creates a new shift swap request. The requester must own the assignment_id. Both assignments must exist and belong to different employees. Both assignments must not already be swapped. The request is created with "pending" status.

**Authentication:** Required (Employee - can only create for their own assignments)

**Request Body:**
```json
{
  "assignment_id": 5,
  "swap_assignment_id": 8,
  "reason": "Personal commitment"
}
```

**cURL Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/swap-requests" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "assignment_id": 5,
    "swap_assignment_id": 8,
    "reason": "Personal commitment"
  }'
```

**Error Responses:**

**If requester doesn't own the assignment:**
```json
{
  "success": false,
  "message": "you can only request swaps for your own assignments"
}
```

**If assignments belong to same employee:**
```json
{
  "success": false,
  "message": "cannot swap with the same employee"
}
```

**If assignment already swapped:**
```json
{
  "success": false,
  "message": "assignment is already swapped"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Swap request created successfully",
  "data": {
    "id": 3,
    "requested_by": "EMP006",
    "requested_by_name": "James Anderson",
    "requested_with": "EMP010",
    "requested_with_name": "Maria Garcia",
    "assignment_id": 5,
    "assignment_date": "2026-01-22",
    "swap_assignment_id": 8,
    "swap_assignment_date": "2026-01-25",
    "reason": "Personal commitment",
    "status": "pending",
    "requested_at": "2026-01-22T10:00:00+03:00"
  }
}
```

### 4.4 Approve Swap Request
**Endpoint:** `POST {{BASE_URL}}/api/v1/rosters/swap-requests/:request_id/approve`

**Description:** Approves a swap request and automatically swaps the employees in both roster assignments. The assignments' status is updated to "swapped", and the swap request status is updated to "approved".

**Path Parameters:**
- `request_id` (required) - Swap request ID

**Authentication:** Required (HR/Admin)

**Request Body:** (Optional)

**cURL Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/swap-requests/1/approve" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "notify_employees": true,
    "notes": "Swap approved. Please confirm your new schedule."
  }'
```

**Error Response (if request not pending):**
```json
{
  "success": false,
  "message": "only pending swap requests can be approved"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Swap request approved successfully",
  "data": {
    "id": 1,
    "status": "approved",
    "reviewed_at": "2026-01-22T15:00:00+03:00",
    "reviewed_by": 2,
    "updated_assignments": [
      {
        "assignment_id": 5,
        "new_employee_id": "EMP010",
        "new_employee_name": "Maria Garcia"
      },
      {
        "assignment_id": 8,
        "new_employee_id": "EMP006",
        "new_employee_name": "James Anderson"
      }
    ]
  }
}
```

### 4.5 Reject Swap Request
**Endpoint:** `POST {{BASE_URL}}/api/v1/rosters/swap-requests/:request_id/reject`

**Description:** Rejects a swap request. A rejection reason is required. The swap request status is updated to "rejected".

**Path Parameters:**
- `request_id` (required) - Swap request ID

**Authentication:** Required (HR/Admin)

**Request Body:**

**cURL Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/swap-requests/1/reject" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rejection_reason": "Insufficient notice period",
    "notify_employees": true
  }'
```

**Error Response (if rejection_reason missing):**
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": "Key: 'RejectSwapRequestRequest.RejectionReason' Error:Field validation for 'RejectionReason' failed on the 'required' tag"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Swap request rejected successfully",
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

## 5. Locations (Helper Endpoint)

**Note:** The locations endpoint is part of the Locations module. Use the existing locations API to retrieve available locations for shift assignment.

**Endpoint:** `GET {{BASE_URL}}/api/v1/locations`

**Description:** Retrieves list of locations for shift assignment. This endpoint is from the Locations module.

**Query Parameters:**
- `is_active` (optional) - Filter by active status: `true`, `false`

**Authentication:** Required (HR/Admin)

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/locations?is_active=true" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Locations retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Main Office",
      "code": "MO-001",
      "address": "123 Business Street",
      "city": "Dar es Salaam",
      "is_active": true
    },
    {
      "id": 2,
      "name": "Branch Office",
      "code": "BO-002",
      "address": "456 Branch Avenue",
      "city": "Arusha",
      "is_active": true
    }
  ]
}
```

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request

**Validation Error:**
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": "Key: 'CreateShiftRequest.ShiftName' Error:Field validation for 'ShiftName' failed on the 'required' tag"
}
```

**Business Logic Error:**
```json
{
  "success": false,
  "message": "shift with this code already exists"
}
```

**Invalid Date Format:**
```json
{
  "success": false,
  "message": "invalid date format. Use YYYY-MM-DD"
}
```

**Invalid Time Format:**
```json
{
  "success": false,
  "message": "invalid start_time format. Use HH:MM:SS"
}
```

### 401 Unauthorized

**Missing Token:**
```json
{
  "success": false,
  "message": "Unauthorized"
}
```

**Invalid Token:**
```json
{
  "success": false,
  "message": "Invalid or expired token"
}
```

### 403 Forbidden

**Insufficient Permissions:**
```json
{
  "success": false,
  "message": "Forbidden - HR or Admin role required"
}
```

### 404 Not Found

**Resource Not Found:**
```json
{
  "success": false,
  "message": "Shift not found"
}
```

```json
{
  "success": false,
  "message": "Roster assignment not found"
}
```

```json
{
  "success": false,
  "message": "Swap request not found"
}
```

### 409 Conflict

**Cannot Delete (Active Rosters):**
```json
{
  "success": false,
  "message": "shift is assigned to active rosters and cannot be deleted"
}
```

**Duplicate Assignment:**
```json
{
  "success": false,
  "message": "roster assignment already exists for this employee and date"
}
```

### 500 Internal Server Error

**Server Error:**
```json
{
  "success": false,
  "message": "Internal server error",
  "error": "Server error details"
}
```

---

## Usage Examples

### Example 1: Complete Shift Management Workflow

**Step 1: Create a Shift**
```bash
curl -X POST "http://localhost:8080/api/v1/shifts" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shift_name": "Morning Shift",
    "shift_code": "MRN",
    "shift_type": "fixed",
    "start_time": "09:00:00",
    "end_time": "18:00:00",
    "break_duration": 60,
    "grace_minutes": 10,
    "late_mark_after": 15,
    "early_going_minutes": 15,
    "half_day_hours": 4.5,
    "minimum_hours": 4,
    "cross_day": false,
    "night_shift": false,
    "weekly_off": ["Sunday"],
    "location_ids": [1, 2],
    "is_active": true
  }'
```

**Step 2: Create Roster Assignments**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/assignments/bulk" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "assignments": [
      {
        "employee_id": "EMP464350",
        "date": "2026-01-22",
        "shift_id": 1,
        "location_id": 1,
        "status": "draft"
      },
      {
        "employee_id": "EMP464350",
        "date": "2026-01-23",
        "shift_id": 1,
        "location_id": 1,
        "status": "draft"
      }
    ]
  }'
```

**Step 3: Publish Roster**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/publish" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2026-01-22",
    "end_date": "2026-01-23",
    "notify_employees": true,
    "notification_message": "Your roster has been published."
  }'
```

### Example 2: Employee Swap Request Workflow

**Step 1: Employee Creates Swap Request**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/swap-requests" \
  -H "Authorization: Bearer EMPLOYEE_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "assignment_id": 5,
    "swap_assignment_id": 8,
    "reason": "Family emergency on 2026-01-22"
  }'
```

**Step 2: HR Approves Swap Request**
```bash
curl -X POST "http://localhost:8080/api/v1/rosters/swap-requests/1/approve" \
  -H "Authorization: Bearer HR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "notify_employees": true,
    "notes": "Swap approved. Please confirm your new schedule."
  }'
```

### Example 3: Weekly Roster View for Dashboard

```bash
curl -X GET "http://localhost:8080/api/v1/rosters/weekly?week=2026-04&department_id=3" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

This returns a grid view perfect for displaying in a calendar interface.

### Example 4: Auto-Schedule for a Week

```bash
curl -X POST "http://localhost:8080/api/v1/rosters/auto-schedule" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2026-01-22",
    "end_date": "2026-01-28",
    "employee_ids": ["EMP464350", "EMP567886", "EMP536795"],
    "location_id": 1,
    "shift_preferences": {
      "EMP464350": [1, 2],
      "EMP567886": [2, 3]
    }
  }'
```

---

## Data Formats

### Date Format
All dates should be in `YYYY-MM-DD` format.
- ✅ Correct: `2026-01-22`
- ❌ Incorrect: `01/22/2026`, `22-01-2026`, `2026/01/22`

### Time Format
All times should be in `HH:mm:ss` format (24-hour).
- ✅ Correct: `09:00:00`, `18:30:00`, `23:59:59`
- ❌ Incorrect: `9:00 AM`, `6:30 PM`, `09:00`

### Weekly Off Days
Array of day names (case-sensitive):
- ✅ Correct: `["Sunday", "Monday"]`
- ❌ Incorrect: `["sunday", "monday"]`, `["Sun", "Mon"]`

### Shift Types
Valid values (lowercase):
- `fixed` - Fixed working hours
- `split` - Split shift (e.g., morning and evening)
- `rotating` - Rotating schedule
- `night` - Night shift
- `flexible` - Flexible working hours

### Status Values

**Roster Assignment Status:**
- `draft` - Not yet published
- `published` - Published and visible to employees
- `swapped` - Has been swapped with another employee

**Swap Request Status:**
- `pending` - Awaiting approval
- `approved` - Approved and swaps executed
- `rejected` - Rejected by HR/Admin

---

## Best Practices

1. **Always use draft status** when creating roster assignments, then publish them in bulk after review.

2. **Validate before publishing:** Review draft assignments before publishing to ensure accuracy.

3. **Use bulk operations** when creating multiple assignments to improve performance.

4. **Check for conflicts:** Before creating assignments, check if employees already have assignments for those dates.

5. **Handle swap requests promptly:** Review and respond to swap requests in a timely manner.

6. **Use weekly view for dashboards:** The weekly roster view is optimized for calendar/grid displays.

7. **Leverage auto-schedule:** Use auto-schedule as a starting point, then manually adjust as needed.

8. **Monitor statistics:** Regularly check statistics to understand shift and roster utilization.

---

## Notes

1. **Date Format:** All dates should be in `YYYY-MM-DD` format.
2. **Time Format:** All times should be in `HH:mm:ss` format (24-hour).
3. **Weekly Off Days:** Array of day names: `["Sunday", "Monday", ...]`
4. **Shift Types:** `fixed`, `split`, `rotating`, `night`, `flexible`
5. **Status Values:** 
   - Roster: `draft`, `published`, `swapped`
   - Swap Request: `pending`, `approved`, `rejected`
6. **Pagination:** Default page size is 20, maximum is 100.
7. **Auto-Schedule:** Uses a basic round-robin algorithm. More sophisticated scheduling rules are planned for future enhancements.
8. **Tenant Isolation:** All endpoints automatically filter data by tenant ID from the authenticated user.
9. **Soft Deletes:** Deleted shifts and roster assignments are soft-deleted and can be recovered if needed.
10. **Working Hours Calculation:** Working hours are automatically calculated from start_time and end_time. Cross-day shifts are supported.

---

## 6. Roster Change Requests

**Note:** For detailed documentation on roster change requests, see [Roster Change Request API](./ROSTER_CHANGE_REQUEST_API.md).

### Quick Reference

**Employee Endpoints:**
- `POST /api/v1/self-service/rosters/change-requests` - Create change request
- `GET /api/v1/self-service/rosters/change-requests` - List my change requests
- `GET /api/v1/self-service/rosters/change-requests/:request_id` - Get my change request

**HR/Admin Endpoints:**
- `GET /api/v1/rosters/change-requests` - List all change requests
- `GET /api/v1/rosters/change-requests/:request_id` - Get change request
- `POST /api/v1/rosters/change-requests/:request_id/approve` - Approve request
- `POST /api/v1/rosters/change-requests/:request_id/reject` - Reject request

### Example: Employee Request Flow

**1. Employee creates change request:**
```bash
POST /api/v1/self-service/rosters/change-requests
{
  "assignment_id": 5,
  "requested_shift_id": 2,
  "requested_date": "2026-01-25",
  "reason": "Need to attend training"
}
```

**2. HR approves:**
```bash
POST /api/v1/rosters/change-requests/1/approve
{
  "notify_employee": true,
  "notes": "Approved"
}
```

**Result:** The assignment is automatically updated with the requested changes.

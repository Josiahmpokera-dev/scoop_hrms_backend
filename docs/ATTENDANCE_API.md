# Attendance Management API

This document provides a comprehensive guide to the Attendance Management API, which includes endpoints for managing employee timesheets, overtime, and manual punch requests.

## Base URL

`/api/v1/attendance`

---

## Table of Contents

### Part A — Overtime Management

1.  [List Overtime Requests](#1-list-overtime-requests)
2.  [Create Overtime Request](#2-create-overtime-request)
3.  [Get Overtime Request](#3-get-overtime-request)
4.  [Approve or Reject Overtime Request](#4-approve-or-reject-overtime-request)

---

## Part A — Overtime Management

### 1. List Overtime Requests

Retrieves a list of all overtime requests, with optional filters for status and employee.

`GET /api/v1/attendance/overtime`

**Query Parameters:**

*   `status` (optional): Filter by status (e.g., `pending`, `approved`, `rejected`)
*   `employeeId` (optional): Filter by a specific employee ID

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "employeeId": 123,
      "date": "2026-03-20",
      "hours": 2,
      "status": "pending",
      "compensationType": "payout"
    }
  ]
}
```

### 2. Create Overtime Request

Allows an employee to submit a new overtime request.

`POST /api/v1/attendance/overtime`

**Request Body:**

```json
{
  "date": "2026-03-20",
  "hours": 2,
  "overtimeType": "weekday",
  "reason": "Urgent project deadline",
  "compensationType": "payout"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Overtime request submitted successfully",
  "data": {
    "id": 1,
    "employeeId": 123,
    "date": "2026-03-20",
    "hours": 2,
    "status": "pending"
  }
}
```

### 3. Get Overtime Request

Retrieves the details of a specific overtime request.

`GET /api/v1/attendance/overtime/{id}`

**Path Parameters:**

*   `id` (required): The ID of the overtime request

**Response:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "employeeId": 123,
    "date": "2026-03-20",
    "hours": 2,
    "status": "pending",
    "reason": "Urgent project deadline",
    "approvals": []
  }
}
```

### 4. Approve or Reject Overtime Request

Allows a manager or HR to approve or reject an overtime request.

`PATCH /api/v1/attendance/overtime/{id}/approve`

**Request Body:**

```json
{
  "status": "approved",
  "comments": "Approved due to project deadline"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Overtime request has been approved"
}
```

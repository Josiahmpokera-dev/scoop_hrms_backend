# Biometric Enrollment & Merged Attendance API

> Link HRMS employee records to biometric device users and retrieve merged attendance data.

## Overview

The biometric device (BioTime) stores attendance records identified by `emp_code`. The HRMS stores employee records identified by `employee_id` (string) and internal `id` (uint). This module provides a formal linking layer between the two systems, enabling:

- **Manual linking** — HR maps an employee to a device emp_code  
- **Bulk linking** — link multiple employees in one call  
- **Auto-linking** — automatically match `Employee.EmployeeID` = `BioTimeTransaction.EmpCode`  
- **Merged attendance** — combine biometric punch data with employee profile data  

---

## Architecture

```
┌─────────────────┐      ┌──────────────────────────┐      ┌─────────────────────┐
│   Employee DB   │      │  biometric_enrollments   │      │ biotime_transactions│
│                 │◄─────┤  employee_id ←→ emp_code │─────►│                     │
│ id, employee_id │  FK  │  status, linked_at       │  ref │ emp_code, punch_time│
│ name, dept, pos │      │  device_user_name        │      │ first/last name     │
└─────────────────┘      └──────────────────────────┘      └─────────────────────┘
```

## Database Table

### `biometric_enrollments`

| Column           | Type         | Constraints                   |
|------------------|--------------|-------------------------------|
| id               | uint         | PK, auto-increment            |
| employee_id      | uint         | NOT NULL, UNIQUE INDEX, FK    |
| emp_code         | varchar(100) | NOT NULL, UNIQUE INDEX        |
| device_user_name | varchar(255) | Name from biometric device    |
| status           | varchar(20)  | default: 'linked'             |
| linked_at        | timestamp    |                               |
| linked_by        | uint         | FK to users                   |
| notes            | text         | Optional                      |
| created_at       | timestamp    |                               |
| updated_at       | timestamp    |                               |
| deleted_at       | timestamp    | Soft delete                   |

---

## API Endpoints — Quick Reference

| # | Method   | Endpoint                                          | Description                          | Access    |
|---|----------|---------------------------------------------------|--------------------------------------|-----------|
| 1 | `GET`    | `/api/v1/biometric/enrollments`                   | List all linked employees            | HR/Admin  |
| 2 | `GET`    | `/api/v1/biometric/enrollments/statistics`        | Enrollment coverage statistics       | HR/Admin  |
| 3 | `GET`    | `/api/v1/biometric/enrollments/unlinked`          | List employees not yet linked        | HR/Admin  |
| 4 | `GET`    | `/api/v1/biometric/enrollments/device-users`      | List device users with linked status | HR/Admin  |
| 5 | `POST`   | `/api/v1/biometric/enrollments/link`              | Link one employee to emp_code        | HR/Admin  |
| 6 | `POST`   | `/api/v1/biometric/enrollments/bulk-link`         | Link multiple employees at once      | HR/Admin  |
| 7 | `POST`   | `/api/v1/biometric/enrollments/auto-link`         | Auto-link by matching codes          | HR/Admin  |
| 8 | `DELETE` | `/api/v1/biometric/enrollments/:employeeId/unlink`| Unlink employee from device          | HR/Admin  |
| 9 | `GET`    | `/api/v1/biometric/attendance/merged`             | Merged daily attendance (all)        | HR/Admin  |
| 10| `GET`    | `/api/v1/biometric/attendance/merged/:employeeId` | Merged attendance for one employee   | HR/Admin  |

---

## Detailed API Reference

---

### 1. List Enrollments

```
GET /api/v1/biometric/enrollments?page=1&page_size=20
```

**Response:**
```json
{
  "success": true,
  "message": "Enrollments retrieved successfully",
  "data": [
    {
      "id": 1,
      "employee_id": 5,
      "employee_code": "EMP001",
      "employee_name": "John Doe",
      "department": "Engineering",
      "position": "Software Engineer",
      "emp_code": "EMP001",
      "device_user_name": "John Doe",
      "status": "linked",
      "linked_at": "2026-02-06 10:30:00"
    }
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

### 2. Enrollment Statistics

```
GET /api/v1/biometric/enrollments/statistics
```

**Response:**
```json
{
  "success": true,
  "message": "Enrollment statistics retrieved successfully",
  "data": {
    "total_employees": 120,
    "linked_employees": 95,
    "unlinked_employees": 25,
    "total_device_users": 100,
    "linked_device_users": 95,
    "unlinked_device_users": 5,
    "coverage_percentage": 79.17
  }
}
```

---

### 3. Get Unlinked Employees

Returns employees that are not yet linked to any biometric device user. If a device user has a matching `emp_code`, it appears as `suggested_emp_code`.

```
GET /api/v1/biometric/enrollments/unlinked
```

**Response:**
```json
{
  "success": true,
  "message": "Unlinked employees retrieved successfully",
  "data": [
    {
      "id": 12,
      "employee_code": "EMP012",
      "first_name": "Jane",
      "last_name": "Smith",
      "department": "Finance",
      "position": "Accountant",
      "suggested_emp_code": "EMP012"
    },
    {
      "id": 15,
      "employee_code": "EMP015",
      "first_name": "Alex",
      "last_name": "Brown",
      "department": "Marketing",
      "position": null,
      "suggested_emp_code": null
    }
  ]
}
```

> When `suggested_emp_code` is non-null, the employee's code matches a biometric device user — the frontend can show a "suggested match" indicator.

---

### 4. List Device Users

Returns all distinct users from the biometric device (from `biotime_transactions`), with their linked status.

```
GET /api/v1/biometric/enrollments/device-users
```

**Response:**
```json
{
  "success": true,
  "message": "Device users retrieved successfully",
  "data": [
    {
      "emp_code": "EMP001",
      "first_name": "John",
      "last_name": "Doe",
      "department": "Engineering",
      "position": "Engineer",
      "linked": true,
      "employee_id": 5
    },
    {
      "emp_code": "BIO042",
      "first_name": "Unknown",
      "last_name": "User",
      "department": "",
      "position": "",
      "linked": false,
      "employee_id": null
    }
  ]
}
```

---

### 5. Link Employee

Link a single employee to a biometric device emp_code.

```
POST /api/v1/biometric/enrollments/link
```

**Request Body:**
```json
{
  "employee_id": 5,
  "emp_code": "EMP001",
  "notes": "Matched by HR during enrollment"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Employee linked to biometric device successfully",
  "data": {
    "id": 1,
    "employee_id": 5,
    "emp_code": "EMP001",
    "device_user_name": "John Doe",
    "status": "linked",
    "linked_at": "2026-02-06T10:30:00Z",
    "linked_by": 1,
    "notes": "Matched by HR during enrollment",
    "created_at": "2026-02-06T10:30:00Z",
    "updated_at": "2026-02-06T10:30:00Z"
  }
}
```

**Error Cases:**
- `400` — Employee not found
- `400` — Employee already linked to a biometric code
- `400` — Biometric code already linked to another employee

---

### 6. Bulk Link

Link multiple employees in a single request.

```
POST /api/v1/biometric/enrollments/bulk-link
```

**Request Body:**
```json
{
  "mappings": [
    { "employee_id": 5, "emp_code": "EMP001" },
    { "employee_id": 6, "emp_code": "EMP002" },
    { "employee_id": 7, "emp_code": "INVALID" }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Bulk link completed",
  "data": {
    "linked": [
      { "id": 1, "employee_id": 5, "emp_code": "EMP001", "status": "linked" },
      { "id": 2, "employee_id": 6, "emp_code": "EMP002", "status": "linked" }
    ],
    "errors": [
      "employee_id=7 emp_code=INVALID: employee with ID 7 not found"
    ]
  }
}
```

---

### 7. Auto-Link

Automatically link employees whose `Employee.EmployeeID` matches `BioTimeTransaction.EmpCode`. This is the fastest way to link when your employee codes are consistent with the biometric device.

```
POST /api/v1/biometric/enrollments/auto-link
```

**No request body required.**

**Response:**
```json
{
  "success": true,
  "message": "Auto-link completed",
  "data": {
    "total_employees": 120,
    "total_device_users": 100,
    "linked": 85,
    "already_linked": 10,
    "no_match": 25,
    "details": [
      {
        "employee_id": 5,
        "employee_code": "EMP001",
        "employee_name": "John Doe",
        "emp_code": "EMP001",
        "action": "linked"
      },
      {
        "employee_id": 6,
        "employee_code": "EMP002",
        "employee_name": "Jane Smith",
        "emp_code": "EMP002",
        "action": "already_linked"
      },
      {
        "employee_id": 15,
        "employee_code": "EMP015",
        "employee_name": "Alex Brown",
        "emp_code": "",
        "action": "no_match"
      }
    ]
  }
}
```

---

### 8. Unlink Employee

Remove the link between an employee and their biometric device user.

```
DELETE /api/v1/biometric/enrollments/:employeeId/unlink
```

**Response:**
```json
{
  "success": true,
  "message": "Employee unlinked from biometric device successfully",
  "data": null
}
```

---

### 9. Merged Daily Attendance (All Employees)

Returns daily attendance for all linked employees, combining biometric punch data with employee profile information.

```
GET /api/v1/biometric/attendance/merged?start_date=2026-02-01&end_date=2026-02-06&page=1&page_size=20
```

**Query Parameters:**

| Parameter    | Required | Default | Description                  |
|-------------|----------|---------|------------------------------|
| start_date  | Yes      | —       | Start date (YYYY-MM-DD)      |
| end_date    | Yes      | —       | End date (YYYY-MM-DD)        |
| employee_id | No       | —       | Filter by specific employee  |
| page        | No       | 1       | Page number                  |
| page_size   | No       | 20      | Items per page (max 100)     |

**Response:**
```json
{
  "success": true,
  "message": "Merged attendance retrieved successfully",
  "data": [
    {
      "employee_id": 5,
      "employee_code": "EMP001",
      "employee_name": "John Doe",
      "department": "Engineering",
      "position": "Software Engineer",
      "photo_url": "https://...",
      "emp_code": "EMP001",
      "date": "2026-02-06",
      "check_in": "08:15:32",
      "check_out": "17:42:18",
      "punch_count": 4,
      "working_hours": "9h 26m",
      "status": "Present"
    },
    {
      "employee_id": 8,
      "employee_code": "EMP008",
      "employee_name": "Sarah Johnson",
      "department": "Finance",
      "position": "Accountant",
      "photo_url": null,
      "emp_code": "EMP008",
      "date": "2026-02-06",
      "check_in": "09:15:00",
      "check_out": "17:30:00",
      "punch_count": 2,
      "working_hours": "8h 15m",
      "status": "Late"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 285,
    "total_pages": 15
  }
}
```

**Status Values:**

| Status    | Meaning                                  |
|-----------|------------------------------------------|
| Present   | Check-in before 08:30, worked >= 4.5h    |
| Late      | Check-in after 08:30                     |
| Half Day  | Worked less than 4.5 hours               |
| Absent    | No check-in and no check-out recorded    |

---

### 10. Merged Attendance for One Employee

Same as above, but scoped to a specific employee.

```
GET /api/v1/biometric/attendance/merged/:employeeId?start_date=2026-02-01&end_date=2026-02-06
```

**Response:** Same structure as endpoint #9.

---

## Workflow Diagrams

### Initial Setup (Recommended Flow)

```
Step 1: GET  /enrollments/statistics        → See coverage gaps
Step 2: GET  /enrollments/device-users      → See who's on the device  
Step 3: GET  /enrollments/unlinked          → See who needs linking
Step 4: POST /enrollments/auto-link         → Match automatically
Step 5: GET  /enrollments/unlinked          → See remaining unmatched
Step 6: POST /enrollments/link (or bulk)    → Manually link the rest
```

### Day-to-Day Usage

```
HR opens attendance page
    │
    ├─► GET /attendance/merged?start_date=...&end_date=...
    │   └─► Returns combined employee + biometric data
    │
    ├─► GET /attendance/merged/:employeeId?start_date=...&end_date=...
    │   └─► Drill into one employee's attendance
    │
    └─► GET /enrollments/statistics
        └─► Dashboard widget showing enrollment coverage
```

### Linking an Employee (Frontend)

```
                  ┌──────────────────────────┐
                  │   Enrollment Dashboard   │
                  │                          │
                  │  ┌────────┐ ┌──────────┐ │
                  │  │Unlinked│ │  Device   │ │
                  │  │  List  │ │  Users    │ │
                  │  └───┬────┘ └────┬─────┘ │
                  │      │           │        │
                  │      └─────┬─────┘        │
                  │            │              │
                  │    Select employee &      │
                  │    select device user     │
                  │            │              │
                  │    POST /enrollments/link │
                  │            │              │
                  │      ┌─────┴─────┐        │
                  │      │  Linked!  │        │
                  │      └───────────┘        │
                  └──────────────────────────┘
```

---

## Frontend Integration Guide

### 1. Enrollment Management Page

Build a page with three sections:

| Section           | Data Source                        | Actions                      |
|-------------------|------------------------------------|------------------------------|
| Statistics Card   | `GET /enrollments/statistics`      | —                            |
| Unlinked List     | `GET /enrollments/unlinked`        | Link button per row          |
| Linked List       | `GET /enrollments`                 | Unlink button per row        |

For each unlinked employee, show the `suggested_emp_code` if available. The frontend can present a dropdown of device users (from `GET /enrollments/device-users`) for manual linking.

### 2. Auto-Link Button

Add a prominent "Auto-Link" button that calls `POST /enrollments/auto-link`. Display the results summary (linked, already linked, no match).

### 3. Attendance Dashboard

Replace the raw biometric attendance view with the merged endpoint:

```
GET /biometric/attendance/merged?start_date=2026-02-01&end_date=2026-02-28
```

This gives you employee names, departments, positions, photos, working hours, and attendance status — all in one call.

### 4. Employee Profile — Attendance Tab

For an individual employee's attendance history:

```
GET /biometric/attendance/merged/:employeeId?start_date=2026-01-01&end_date=2026-01-31
```

---

## Access Control

All endpoints require **HR or Admin** role:

```
middleware.AuthMiddleware() + middleware.HRMiddleware()
```

| Role     | Access                          |
|----------|---------------------------------|
| Admin    | Full access to all endpoints    |
| HR       | Full access to all endpoints    |
| Manager  | No access (use raw biometric)   |
| Employee | No access                       |

---

## Error Responses

| Code | Scenario                                       |
|------|-------------------------------------------------|
| 400  | Invalid request body / missing required fields  |
| 400  | Employee not found                              |
| 400  | Employee already linked                         |
| 400  | Biometric code already linked to another        |
| 400  | Employee not linked (when unlinking)            |
| 400  | Missing start_date or end_date                  |
| 500  | Database / internal error                       |

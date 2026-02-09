# Attendance Module API Documentation

## Overview

The Attendance module provides **Timesheets**, **Overtime**, and **Reports** functionality. It integrates with the existing employee, shift/roster, and payroll modules.

**Base URL:** `/api/v1/attendance`

**Authentication:** All endpoints require a valid JWT token via `Authorization: Bearer <token>`.

---

## Table of Contents

1. [Timesheets](#1-timesheets)
   - [My Timesheets](#11-my-timesheets)
   - [Timesheet Entries](#12-timesheet-entries)
   - [Submit / Recall](#13-submit--recall)
   - [Copy Last Week](#14-copy-last-week)
   - [Employee Stats](#15-employee-stats)
   - [Team Timesheets (Manager)](#16-team-timesheets-manager)
   - [Approvals (Admin/HR)](#17-approvals-adminhr)
2. [Overtime](#2-overtime)
   - [OT Requests (Employee)](#21-ot-requests-employee)
   - [OT Policies (View)](#22-ot-policies-view)
   - [OT Admin (Admin/HR)](#23-ot-admin-adminhr)
3. [Reports (Admin/HR)](#3-reports-adminhr)
4. [Data Models](#4-data-models)
5. [Common Flows](#5-common-flows)

---

## 1. Timesheets

### 1.1 My Timesheets

#### List My Weekly Timesheets

```
GET /api/v1/attendance/timesheets
```

**Access:** Authenticated user (own data)

**Query Parameters:**

| Parameter  | Type   | Required | Description                              |
|------------|--------|----------|------------------------------------------|
| page       | int    | No       | Page number (default: 1)                 |
| page_size  | int    | No       | Items per page (default: 20, max: 100)   |
| status     | string | No       | Filter: draft, submitted, approved, rejected, recalled |

**Response (200):**

```json
{
  "success": true,
  "message": "Timesheets retrieved successfully",
  "data": [
    {
      "id": 1,
      "employee_id": 5,
      "week_start": "2026-02-02",
      "week_end": "2026-02-08",
      "status": "draft",
      "total_hours": 32.5,
      "billable_hours": 24.0,
      "submitted_at": null,
      "approved_at": null,
      "entries": [
        {
          "id": 1,
          "timesheet_id": 1,
          "employee_id": 5,
          "date": "2026-02-02",
          "project_name": "HRMS Platform",
          "client_name": "Internal",
          "task_description": "API development for attendance module",
          "entry_type": "regular",
          "hours": 8.0,
          "is_billable": true,
          "notes": null
        }
      ]
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 4,
    "total_pages": 1
  }
}
```

---

#### Get Timesheet by ID

```
GET /api/v1/attendance/timesheets/:timesheet_id
```

**Access:** Owner or Admin/HR

**Response (200):** Same as single item from list above (with entries and approvals).

---

### 1.2 Timesheet Entries

#### Create Single Entry

```
POST /api/v1/attendance/timesheets/entries
```

**Access:** Authenticated user (own data)

**Request Body:**

```json
{
  "date": "2026-02-03",
  "project_name": "HRMS Platform",
  "client_name": "Internal",
  "task_description": "Implemented timesheet approval workflow",
  "entry_type": "regular",
  "hours": 6.5,
  "is_billable": true,
  "notes": "Completed multi-level approval chain"
}
```

| Field            | Type    | Required | Description                                |
|------------------|---------|----------|--------------------------------------------|
| date             | string  | Yes      | Date (YYYY-MM-DD)                          |
| project_name     | string  | Yes      | Project name                               |
| client_name      | string  | No       | Client name (optional)                     |
| task_description | string  | Yes      | Task description                           |
| entry_type       | string  | No       | Type of work (default: `regular`). Options: `regular`, `meeting`, `training`, `travel`, `support`, `overtime`, `break`, `other` |
| hours            | float   | No       | Hours worked (0.5h precision, max 24). Optional — can be filled in later. Defaults to 0 if omitted. |
| is_billable      | boolean | No       | Whether hours are billable (default: false) |
| notes            | string  | No       | Additional notes                           |

**Response (201):**

```json
{
  "success": true,
  "message": "Timesheet entry created successfully",
  "data": {
    "id": 5,
    "timesheet_id": 1,
    "employee_id": 5,
    "date": "2026-02-03",
    "project_name": "HRMS Platform",
    "client_name": "Internal",
    "task_description": "Implemented timesheet approval workflow",
    "entry_type": "regular",
    "hours": 6.5,
    "is_billable": true,
    "notes": "Completed multi-level approval chain"
  }
}
```

> **Note:** The system automatically creates or retrieves the weekly timesheet (Monday-Sunday) for the given date. Entries can only be added to timesheets in `draft` or `recalled` status.
> 
> **Hours are optional:** You can create an entry without specifying hours (e.g., log a task/project first, then fill in hours later via the Update Entry endpoint). Hours default to `0` if omitted.

---

#### Bulk Create Entries

```
POST /api/v1/attendance/timesheets/entries/bulk
```

**Access:** Authenticated user (own data)

**Request Body:**

```json
{
  "week_start": "2026-02-02",
  "entries": [
    {
      "date": "2026-02-02",
      "project_name": "HRMS Platform",
      "task_description": "Backend API development",
      "entry_type": "regular",
      "hours": 8.0,
      "is_billable": true
    },
    {
      "date": "2026-02-03",
      "project_name": "HRMS Platform",
      "task_description": "Frontend integration",
      "entry_type": "regular",
      "hours": 7.5,
      "is_billable": true
    },
    {
      "date": "2026-02-03",
      "project_name": "Internal Training",
      "task_description": "Go best practices workshop",
      "entry_type": "training",
      "is_billable": false
    }
  ]
}
```

> The third entry above omits `hours` — it will default to `0` and can be filled in later.

**Response (201):** Returns the full timesheet week with all entries.

---

#### Update Entry

```
PUT /api/v1/attendance/timesheets/entries/:entry_id
```

**Access:** Owner only (timesheet must be in `draft` or `recalled` status)

**Request Body:** Same as Create Entry.

**Response (200):** Returns updated entry.

---

#### Delete Entry

```
DELETE /api/v1/attendance/timesheets/entries/:entry_id
```

**Access:** Owner only (timesheet must be in `draft` or `recalled` status)

**Response (200):**

```json
{
  "success": true,
  "message": "Timesheet entry deleted successfully",
  "data": null
}
```

---

### 1.3 Submit / Recall

#### Submit Timesheet for Approval

```
POST /api/v1/attendance/timesheets/:timesheet_id/submit
```

**Access:** Owner only

**Validation:**
- Timesheet must be in `draft` or `recalled` status
- Must have at least one entry

**Response (200):**

```json
{
  "success": true,
  "message": "Timesheet submitted for approval",
  "data": {
    "id": 1,
    "status": "submitted",
    "submitted_at": "2026-02-07T14:30:00Z",
    "total_hours": 40.0,
    "billable_hours": 32.0
  }
}
```

---

#### Recall Submitted Timesheet

```
POST /api/v1/attendance/timesheets/:timesheet_id/recall
```

**Access:** Owner only (timesheet must be in `submitted` status)

**Response (200):**

```json
{
  "success": true,
  "message": "Timesheet recalled successfully",
  "data": {
    "id": 1,
    "status": "recalled"
  }
}
```

---

### 1.4 Copy Last Week

```
POST /api/v1/attendance/timesheets/copy-last-week
```

**Access:** Authenticated user

**Request Body:**

```json
{
  "source_week_start": "2026-01-26",
  "target_week_start": "2026-02-02"
}
```

| Field             | Type   | Required | Description                         |
|-------------------|--------|----------|-------------------------------------|
| source_week_start | string | Yes      | Monday of the week to copy FROM     |
| target_week_start | string | Yes      | Monday of the week to copy TO       |

**Response (200):** Returns the target week with copied entries. Dates are adjusted to the new week.

---

### 1.5 Employee Stats

```
GET /api/v1/attendance/timesheets/stats
```

**Access:** Authenticated user (own data)

**Query Parameters:**

| Parameter  | Type   | Required | Description              |
|------------|--------|----------|--------------------------|
| start_date | string | Yes      | Start date (YYYY-MM-DD)  |
| end_date   | string | Yes      | End date (YYYY-MM-DD)    |

**Response (200):**

```json
{
  "success": true,
  "message": "Timesheet statistics retrieved successfully",
  "data": {
    "employee_id": "EMP001",
    "employee_name": "John Doe",
    "period": {
      "start": "2026-01-01",
      "end": "2026-01-31"
    },
    "total_hours": 168.0,
    "billable_hours": 134.5,
    "non_billable_hours": 33.5,
    "utilization_rate": 80.06,
    "project_distribution": [
      {
        "project_name": "HRMS Platform",
        "total_hours": 100.0,
        "billable_hours": 95.0,
        "entry_count": 22
      },
      {
        "project_name": "Internal Training",
        "total_hours": 33.5,
        "billable_hours": 0.0,
        "entry_count": 8
      }
    ]
  }
}
```

---

### 1.6 Team Timesheets (Manager)

```
GET /api/v1/attendance/timesheets/team
```

**Access:** Authenticated user (shows direct reports)

**Query Parameters:**

| Parameter  | Type   | Required | Description                              |
|------------|--------|----------|------------------------------------------|
| week_start | string | No       | Filter by specific week (YYYY-MM-DD, Monday) |
| status     | string | No       | Filter by status                         |
| page       | int    | No       | Page number (default: 1)                 |
| page_size  | int    | No       | Items per page (default: 20)             |

**Response (200):** Same structure as "My Timesheets" but for direct reports.

---

### 1.7 Approvals (Admin/HR)

#### List Pending Approvals

```
GET /api/v1/attendance/timesheets/approvals
```

**Access:** Admin/HR only

**Query Parameters:** page, page_size

**Response (200):** Paginated list of timesheets with `submitted` status.

---

#### Approve or Reject Timesheet

```
POST /api/v1/attendance/timesheets/approvals/:timesheet_id
```

**Access:** Admin/HR only

**Request Body:**

```json
{
  "status": "approved",
  "comments": "All entries verified. Approved."
}
```

| Field    | Type   | Required | Description                        |
|----------|--------|----------|------------------------------------|
| status   | string | Yes      | `approved` or `rejected`           |
| comments | string | No       | Approval/rejection comments        |

**Response (200):**

```json
{
  "success": true,
  "message": "Timesheet processed successfully",
  "data": {
    "id": 1,
    "status": "approved",
    "approved_at": "2026-02-08T10:00:00Z",
    "approved_by": 1
  }
}
```

---

## 2. Overtime

### 2.1 OT Requests (Employee)

#### Submit Overtime Request

```
POST /api/v1/attendance/overtime/requests
```

**Access:** Authenticated user

**Request Body:**

```json
{
  "date": "2026-02-05",
  "hours": 3.0,
  "overtime_type": "weekday",
  "reason": "Urgent deployment for client release deadline",
  "compensation_type": "payout"
}
```

| Field             | Type   | Required | Description                                 |
|-------------------|--------|----------|---------------------------------------------|
| date              | string | Yes      | Date of overtime (YYYY-MM-DD)               |
| hours             | float  | Yes      | Hours of overtime (max 24)                  |
| overtime_type     | string | Yes      | `weekday`, `weekend`, or `holiday`          |
| reason            | string | Yes      | Reason for overtime (min 3 chars)           |
| compensation_type | string | Yes      | `payout`, `comp_off`, or `both`             |

**Validation (automatic):**
- Checks against daily/weekly/monthly OT limits from the applicable policy
- Validates comp-off is allowed by policy
- Auto-calculates multiplier from policy (weekday 1.5x, weekend 2.0x, holiday 2.5x)
- Auto-calculates payout amount and/or comp-off hours

**Response (201):**

```json
{
  "success": true,
  "message": "Overtime request submitted successfully",
  "data": {
    "id": 1,
    "employee_id": 5,
    "policy_id": 1,
    "date": "2026-02-05",
    "overtime_type": "weekday",
    "hours": 3.0,
    "reason": "Urgent deployment for client release deadline",
    "compensation_type": "payout",
    "status": "pending",
    "multiplier": 1.5,
    "payout_amount": 4.5,
    "comp_off_hours": null,
    "comp_off_expiry": null
  }
}
```

---

#### List My OT Requests

```
GET /api/v1/attendance/overtime/requests
```

**Access:** Authenticated user (own data)

**Query Parameters:**

| Parameter  | Type   | Required | Description                                   |
|------------|--------|----------|-----------------------------------------------|
| status     | string | No       | Filter: pending, approved, rejected, processed, cancelled |
| page       | int    | No       | Page number (default: 1)                      |
| page_size  | int    | No       | Items per page (default: 20)                  |

**Response (200):** Paginated list of OT requests.

---

#### Get OT Request by ID

```
GET /api/v1/attendance/overtime/requests/:request_id
```

**Access:** Owner or Admin/HR

**Response (200):** Single OT request with approval history.

---

#### Cancel OT Request

```
POST /api/v1/attendance/overtime/requests/:request_id/cancel
```

**Access:** Owner only (request must be in `pending` status)

**Response (200):**

```json
{
  "success": true,
  "message": "Overtime request cancelled successfully",
  "data": {
    "id": 1,
    "status": "cancelled"
  }
}
```

---

### 2.2 OT Policies (View)

#### List Overtime Policies

```
GET /api/v1/attendance/overtime/policies
```

**Access:** Authenticated user (view-only for employees)

**Response (200):**

```json
{
  "success": true,
  "message": "Overtime policies retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Default OT Policy",
      "department_id": null,
      "is_active": true,
      "daily_threshold_hours": 8.0,
      "weekly_threshold_hours": 40.0,
      "weekday_multiplier": 1.5,
      "weekend_multiplier": 2.0,
      "holiday_multiplier": 2.5,
      "max_daily_ot_hours": 4.0,
      "max_weekly_ot_hours": 20.0,
      "max_monthly_ot_hours": 60.0,
      "requires_pre_approval": true,
      "allow_comp_off": true,
      "comp_off_ratio": 1.0,
      "comp_off_expiry_days": 90,
      "monthly_budget_cap": null
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 1, "total_pages": 1 }
}
```

---

#### Get Policy by ID

```
GET /api/v1/attendance/overtime/policies/:policy_id
```

**Access:** Authenticated user

---

#### Get Overtime Statistics

```
GET /api/v1/attendance/overtime/stats
```

**Access:** Authenticated user

**Query Parameters:**

| Parameter  | Type   | Required | Description              |
|------------|--------|----------|--------------------------|
| start_date | string | No       | Start date (defaults to current month start) |
| end_date   | string | No       | End date (defaults to current month end)     |

**Response (200):**

```json
{
  "success": true,
  "message": "Overtime statistics retrieved successfully",
  "data": {
    "total_requests": 15,
    "pending_count": 3,
    "approved_hours": 45.5,
    "total_payout": 68.25
  }
}
```

---

### 2.3 OT Admin (Admin/HR)

#### Create OT Policy

```
POST /api/v1/attendance/overtime/admin/policies
```

**Access:** Admin/HR only

**Request Body:**

```json
{
  "name": "Engineering Department OT Policy",
  "department_id": 3,
  "daily_threshold_hours": 8.0,
  "weekly_threshold_hours": 40.0,
  "weekday_multiplier": 1.5,
  "weekend_multiplier": 2.0,
  "holiday_multiplier": 2.5,
  "max_daily_ot_hours": 4.0,
  "max_weekly_ot_hours": 20.0,
  "max_monthly_ot_hours": 60.0,
  "requires_pre_approval": true,
  "allow_comp_off": true,
  "comp_off_ratio": 1.0,
  "comp_off_expiry_days": 90,
  "monthly_budget_cap": 5000000.00
}
```

**Response (201):** Returns the created policy.

---

#### Update OT Policy

```
PUT /api/v1/attendance/overtime/admin/policies/:policy_id
```

**Access:** Admin/HR only

**Request Body:** Same as Create.

**Response (200):** Returns updated policy.

---

#### Delete OT Policy

```
DELETE /api/v1/attendance/overtime/admin/policies/:policy_id
```

**Access:** Admin/HR only

**Response (200):**

```json
{
  "success": true,
  "message": "Overtime policy deleted successfully",
  "data": null
}
```

---

#### List All OT Requests

```
GET /api/v1/attendance/overtime/admin/requests
```

**Access:** Admin/HR only

**Query Parameters:**

| Parameter  | Type   | Required | Description                                   |
|------------|--------|----------|-----------------------------------------------|
| status     | string | No       | Filter by status                              |
| start_date | string | No       | Filter by start date (YYYY-MM-DD)             |
| end_date   | string | No       | Filter by end date (YYYY-MM-DD)               |
| page       | int    | No       | Page number                                   |
| page_size  | int    | No       | Items per page                                |

**Response (200):** Paginated list of all OT requests across the organization.

---

#### List Pending OT Approvals

```
GET /api/v1/attendance/overtime/admin/pending
```

**Access:** Admin/HR only

**Response (200):** Paginated list of pending OT requests.

---

#### Approve or Reject OT Request

```
POST /api/v1/attendance/overtime/admin/requests/:request_id/approve
```

**Access:** Admin/HR only

**Request Body:**

```json
{
  "status": "approved",
  "comments": "Approved for the client release deadline"
}
```

| Field    | Type   | Required | Description                  |
|----------|--------|----------|------------------------------|
| status   | string | Yes      | `approved` or `rejected`     |
| comments | string | No       | Comments                     |

**Response (200):** Returns the updated OT request.

---

## 3. Reports (Admin/HR)

All report endpoints require **Admin/HR** role.

### 3.1 Timesheet Summary Report

```
GET /api/v1/attendance/reports/timesheets
```

**Query Parameters:**

| Parameter  | Type   | Required | Description              |
|------------|--------|----------|--------------------------|
| start_date | string | Yes      | Start date (YYYY-MM-DD)  |
| end_date   | string | Yes      | End date (YYYY-MM-DD)    |

**Response (200):**

```json
{
  "success": true,
  "message": "Timesheet summary report",
  "data": {
    "period": {
      "start_date": "2026-01-01",
      "end_date": "2026-01-31"
    },
    "summary": {
      "total_entries": 450,
      "total_hours": 3600.0,
      "billable_hours": 2880.0,
      "non_billable_hours": 720.0,
      "utilization_rate": 80.0,
      "unique_employees": 25,
      "unique_projects": 8
    },
    "weekly_submission_status": [
      { "status": "approved", "count": 80 },
      { "status": "submitted", "count": 10 },
      { "status": "draft", "count": 5 }
    ],
    "project_distribution": [
      {
        "project_name": "HRMS Platform",
        "total_hours": 1200.0,
        "billable_hours": 1100.0,
        "employee_count": 15
      }
    ]
  }
}
```

---

### 3.2 Overtime Summary Report

```
GET /api/v1/attendance/reports/overtime
```

**Query Parameters:** start_date, end_date (both required)

**Response (200):**

```json
{
  "success": true,
  "message": "Overtime summary report",
  "data": {
    "period": {
      "start_date": "2026-01-01",
      "end_date": "2026-01-31"
    },
    "summary": {
      "total_requests": 45,
      "approved_count": 38,
      "pending_count": 5,
      "rejected_count": 2,
      "total_hours": 135.0,
      "approved_hours": 114.0,
      "total_payout": 171000.0,
      "total_comp_off_hours": 12.0,
      "unique_employees": 18
    },
    "by_overtime_type": [
      { "overtime_type": "weekday", "count": 30, "total_hours": 90.0 },
      { "overtime_type": "weekend", "count": 10, "total_hours": 30.0 },
      { "overtime_type": "holiday", "count": 5, "total_hours": 15.0 }
    ],
    "by_compensation_type": [
      { "compensation_type": "payout", "count": 35, "total_hours": 105.0 },
      { "compensation_type": "comp_off", "count": 8, "total_hours": 24.0 },
      { "compensation_type": "both", "count": 2, "total_hours": 6.0 }
    ],
    "monthly_trend": [
      { "month": "2026-01", "total_hours": 114.0, "count": 38 }
    ]
  }
}
```

---

### 3.3 Project Utilization Report

```
GET /api/v1/attendance/reports/project-utilization
```

**Query Parameters:** start_date, end_date (both required)

**Response (200):**

```json
{
  "success": true,
  "message": "Project utilization report",
  "data": {
    "period": { "start_date": "2026-01-01", "end_date": "2026-01-31" },
    "projects": [
      {
        "project_name": "HRMS Platform",
        "client_name": "Internal",
        "total_hours": 1200.0,
        "billable_hours": 1100.0,
        "non_billable_hours": 100.0,
        "team_size": 15,
        "utilization_rate": 91.67
      }
    ]
  }
}
```

---

### 3.4 Employee Utilization Report

```
GET /api/v1/attendance/reports/employee-utilization
```

**Query Parameters:** start_date, end_date (both required)

**Response (200):**

```json
{
  "success": true,
  "message": "Employee utilization report",
  "data": {
    "period": { "start_date": "2026-01-01", "end_date": "2026-01-31" },
    "employees": [
      {
        "employee_id": 5,
        "total_hours": 168.0,
        "billable_hours": 134.5,
        "non_billable_hours": 33.5,
        "utilization_rate": 80.06,
        "project_count": 3
      }
    ]
  }
}
```

---

## 4. Data Models

### TimesheetWeek

| Field           | Type     | Description                                        |
|-----------------|----------|----------------------------------------------------|
| id              | uint     | Primary key                                        |
| employee_id     | uint     | Employee reference                                 |
| week_start      | date     | Monday of the week                                 |
| week_end        | date     | Sunday of the week                                 |
| status          | string   | draft, submitted, approved, rejected, recalled     |
| total_hours     | float    | Auto-calculated sum of entry hours                 |
| billable_hours  | float    | Auto-calculated sum of billable entry hours        |
| submitted_at    | datetime | When the timesheet was submitted                   |
| approved_at     | datetime | When the timesheet was approved                    |
| approved_by     | uint     | User who approved                                  |
| rejection_reason| string   | Reason for rejection                               |
| entries         | array    | Related TimesheetEntry records                     |
| approvals       | array    | Related TimesheetApproval records                  |

### TimesheetEntry

| Field            | Type    | Description                                |
|------------------|---------|---------------------------------------------|
| id               | uint    | Primary key                                 |
| timesheet_id     | uint    | References TimesheetWeek                    |
| employee_id      | uint    | Employee reference                          |
| date             | date    | Entry date                                  |
| project_name     | string  | Project name                                |
| client_name      | string  | Client name (optional)                      |
| task_description | string  | Task description                            |
| entry_type       | string  | Type of work: `regular`, `meeting`, `training`, `travel`, `support`, `overtime`, `break`, `other` (default: `regular`) |
| hours            | float   | Hours worked (0.5h precision). Optional — defaults to 0 if not provided. |
| is_billable      | boolean | Whether hours are billable                  |
| notes            | string  | Additional notes                            |

### TimesheetApproval

| Field        | Type     | Description                              |
|--------------|----------|------------------------------------------|
| id           | uint     | Primary key                              |
| timesheet_id | uint     | References TimesheetWeek                 |
| approver_id  | uint     | Approver's employee ID                   |
| level        | int      | Approval level (1=Manager, 2=HR/Finance) |
| status       | string   | pending, approved, rejected              |
| comments     | string   | Approval comments                        |
| action_at    | datetime | When the action was taken                |

### OvertimePolicy

| Field                  | Type    | Description                                       |
|------------------------|---------|---------------------------------------------------|
| id                     | uint    | Primary key                                       |
| name                   | string  | Policy name                                       |
| department_id          | uint    | Department (NULL = all departments)               |
| is_active              | boolean | Whether policy is active                          |
| daily_threshold_hours  | float   | Hours beyond which OT kicks in (daily)            |
| weekly_threshold_hours | float   | Hours beyond which OT kicks in (weekly)           |
| weekday_multiplier     | float   | Weekday OT multiplier (e.g., 1.5x)               |
| weekend_multiplier     | float   | Weekend OT multiplier (e.g., 2.0x)               |
| holiday_multiplier     | float   | Holiday OT multiplier (e.g., 2.5x)               |
| max_daily_ot_hours     | float   | Maximum daily OT hours                            |
| max_weekly_ot_hours    | float   | Maximum weekly OT hours                           |
| max_monthly_ot_hours   | float   | Maximum monthly OT hours                          |
| requires_pre_approval  | boolean | Whether OT requires pre-approval                  |
| allow_comp_off         | boolean | Whether comp-off is allowed                       |
| comp_off_ratio         | float   | Comp-off conversion ratio (e.g., 1.0)             |
| comp_off_expiry_days   | int     | Days until comp-off expires                       |
| monthly_budget_cap     | float   | Monthly budget cap (optional)                     |

### OvertimeRequest

| Field             | Type     | Description                                     |
|-------------------|----------|--------------------------------------------------|
| id                | uint     | Primary key                                      |
| employee_id       | uint     | Employee reference                               |
| policy_id         | uint     | Applied OT policy                                |
| date              | date     | OT date                                         |
| overtime_type     | string   | weekday, weekend, holiday                        |
| hours             | float    | OT hours requested                               |
| reason            | string   | Reason for OT                                    |
| compensation_type | string   | payout, comp_off, both                           |
| status            | string   | pending, approved, rejected, processed, cancelled|
| multiplier        | float    | Applied multiplier (auto-calculated)             |
| payout_amount     | float    | Calculated payout (hours * multiplier)           |
| comp_off_hours    | float    | Credited comp-off hours                          |
| comp_off_expiry   | datetime | Comp-off expiry date                             |
| approved_by       | uint     | User who approved                                |
| approved_at       | datetime | When approved                                    |
| rejection_reason  | string   | Reason for rejection                             |
| is_auto_detected  | boolean  | Whether detected from attendance data            |
| approvals         | array    | Approval history                                 |

### OvertimeApproval

| Field               | Type     | Description                         |
|---------------------|----------|--------------------------------------|
| id                  | uint     | Primary key                          |
| overtime_request_id | uint     | References OvertimeRequest           |
| approver_id         | uint     | Approver's employee ID               |
| level               | int      | Approval level (1=Manager, 2=HR)     |
| status              | string   | pending, approved, rejected          |
| comments            | string   | Approval comments                    |
| action_at           | datetime | When the action was taken            |

---

## 5. Common Flows

### Flow 1: Employee Submits Weekly Timesheet

1. Employee creates entries throughout the week:
   - `POST /api/v1/attendance/timesheets/entries` (individual entries)
   - OR `POST /api/v1/attendance/timesheets/entries/bulk` (all at once)
2. System auto-creates the weekly timesheet (TimesheetWeek) in `draft` status.
3. Employee reviews their timesheet:
   - `GET /api/v1/attendance/timesheets` (list all weeks)
4. Employee submits:
   - `POST /api/v1/attendance/timesheets/:id/submit` (status -> `submitted`)
5. Admin/HR reviews pending:
   - `GET /api/v1/attendance/timesheets/approvals`
6. Admin/HR approves or rejects:
   - `POST /api/v1/attendance/timesheets/approvals/:id` with `{"status": "approved"}`
7. If rejected, employee can recall and edit:
   - `POST /api/v1/attendance/timesheets/:id/recall` (status -> `recalled`)

### Flow 2: Copy Last Week's Timesheet

1. `POST /api/v1/attendance/timesheets/copy-last-week` with source and target weeks
2. All entries are duplicated with dates adjusted to the new week.
3. Employee can then edit/add/remove entries before submitting.

### Flow 3: Employee Requests Overtime

1. Employee submits OT request:
   - `POST /api/v1/attendance/overtime/requests`
2. System validates against applicable policy (daily/weekly/monthly limits).
3. System auto-calculates multiplier and compensation.
4. Admin/HR reviews pending requests:
   - `GET /api/v1/attendance/overtime/admin/pending`
5. Admin/HR approves or rejects:
   - `POST /api/v1/attendance/overtime/admin/requests/:id/approve`

### Flow 4: HR Creates OT Policy

1. HR creates a policy (optional: per-department):
   - `POST /api/v1/attendance/overtime/admin/policies`
2. Policy defines multipliers, limits, and compensation rules.
3. When an employee creates an OT request, the system matches the policy by department (falls back to the default if no department-specific policy exists).

---

## Error Responses

All endpoints return standard error responses:

```json
{
  "success": false,
  "message": "Descriptive error message",
  "error": "Error details"
}
```

**Common HTTP Status Codes:**
- `400` - Bad request / validation error
- `401` - Unauthorized (no/invalid token)
- `403` - Forbidden (insufficient role)
- `404` - Resource not found
- `500` - Internal server error

---

## Notes

1. **Timesheet Auto-Creation:** Weekly timesheets (TimesheetWeek) are automatically created when the first entry for a week is added. You don't need to create them manually.

2. **Week Boundaries:** Weeks run Monday to Sunday. The system automatically calculates the correct week for any given date.

3. **Hours Precision:** Timesheet hours support 0.5-hour precision (e.g., 7.5, 8.0, 8.5).

4. **Policy Matching:** OT requests automatically match the most specific policy (department-specific first, then default). If no policy exists, the request is rejected.

5. **Multiplier Calculation:** OT multipliers are automatically applied based on the OT type (weekday/weekend/holiday) and the applicable policy.

6. **Compensation Split:** When `compensation_type` is `"both"`, the system applies a 50-50 split between payout and comp-off.

7. **Soft Deletes:** All records use GORM soft deletes. Deleted entries are hidden but preserved in the database for audit purposes.

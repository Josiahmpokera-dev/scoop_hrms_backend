# Leave Management API — Comprehensive Documentation

Complete API documentation covering **Employee self-service** leave requests, **HR/Admin** management, and the **multi-level approval workflow** (HOD → HR → Director/CEO).

**Base URL:** `/api/v1`

---

## Table of Contents

### Part A — Employee Self-Service APIs
1. [Get Employee Info (Form Pre-fill)](#1-get-employee-info)
2. [Get Leave Balances](#2-get-leave-balances)
3. [Get Active Leave Types (Dropdown)](#3-get-active-leave-types)
4. [Get Holidays](#4-get-holidays)
5. [Get Policy Guidelines](#5-get-policy-guidelines)
6. [Calculate Leave Days](#6-calculate-leave-days)
7. [Submit Leave Application](#7-submit-leave-application)
8. [List My Leave Requests](#8-list-my-leave-requests)
9. [Get Leave Request Detail](#9-get-leave-request-detail)
10. [Update Leave Request (Draft/Returned)](#10-update-leave-request)
11. [Cancel Leave Request](#11-cancel-leave-request)
12. [Delete Draft Leave Request](#12-delete-draft-leave-request)

### Part B — HR/Admin Management APIs
13. [List All Leave Requests (HR View)](#13-list-all-leave-requests-hr-view)
14. [Approve Leave Request](#14-approve-leave-request)
15. [Reject Leave Request](#15-reject-leave-request)
16. [Return for Information](#16-return-for-information)

### Part C — HR/Admin Configuration APIs
17. [Leave Types CRUD](#17-leave-types-crud)
18. [Leave Policies CRUD](#18-leave-policies-crud)
19. [Holidays CRUD](#19-holidays-crud)

### Part D — Leave Calendar (All Users)
20. [Leave Calendar](#20-leave-calendar)
21. [Employees on Leave Today](#21-employees-on-leave-today)
22. [Weekly Calendar](#22-weekly-calendar)

### Appendix
- [Approval Workflow](#approval-workflow)
- [Data Models](#data-models)
- [Quick Reference Table](#quick-reference-table)

---

## Authentication

| Scope              | Middleware                          | Who Can Access              |
|--------------------|-------------------------------------|-----------------------------|
| Employee (Self)    | `AuthMiddleware`                    | Any authenticated user      |
| HR/Admin           | `AuthMiddleware` + `HRMiddleware`   | Admin, Super Admin, HR      |

All endpoints require `Authorization: Bearer <token>`.

---

# Part A — Employee Self-Service APIs

Base URL: `/api/v1/leave`

These endpoints are for **any authenticated employee** to manage their own leave.

---

## 1. Get Employee Info

Pre-fill the leave application form with employee information.

```
GET /api/v1/leave/employee-info
```

**Auth:** Any authenticated user

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Employee information retrieved successfully",
  "data": {
    "employee_id": "EMP060",
    "first_name": "Anna",
    "middle_name": null,
    "last_name": "Juma",
    "job_position": "Software Developer",
    "department": "Engineering",
    "contact_number": "+255712345678",
    "emergency_contact_person": "Jane Doe",
    "emergency_contact_number": "+255712345679"
  }
}
```

---

## 2. Get Leave Balances

Returns leave balance statistics for the stat cards and leave type breakdown.

```
GET /api/v1/leave/balances
```

**Auth:** Any authenticated user

| Query Param   | Type | Default | Description                |
|---------------|------|---------|----------------------------|
| `year`        | int  | current | Leave year to query        |
| `employee_id` | string | self  | Employee ID (HR/Admin only)|

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Leave balances retrieved successfully",
  "data": [
    {
      "leave_type": "annual",
      "leave_type_name": "Annual Leave",
      "entitlement": 21,
      "used": 5,
      "pending": 2,
      "available": 14,
      "carried_forward": 3,
      "expires_on": "2026-03-31",
      "leave_taken_from_previous_year": 2,
      "leave_balance_from_previous_year": 3,
      "previous_days_used": 0
    },
    {
      "leave_type": "sick",
      "leave_type_name": "Sick Leave",
      "entitlement": 14,
      "used": 1,
      "pending": 0,
      "available": 13,
      "carried_forward": 0,
      "expires_on": null,
      "leave_taken_from_previous_year": 0,
      "leave_balance_from_previous_year": 0,
      "previous_days_used": 0
    }
  ]
}
```

### Frontend Stat Card Mapping

| Card               | Source                                                   |
|--------------------|----------------------------------------------------------|
| **Days Remaining** | `SUM(data[].available)`                                  |
| **Days Used**      | `SUM(data[].used)`                                       |
| **Pending**        | `SUM(data[].pending)`                                    |
| **Eligible Year**  | If `SUM(data[].entitlement) > 0` → Eligible, else Not   |

---

## 3. Get Active Leave Types

Returns active leave types for the form dropdown. Read-only endpoint for employees.

```
GET /api/v1/leave/active-types
```

**Auth:** Any authenticated user

| Query Param | Type | Default | Description        |
|-------------|------|---------|--------------------|
| `page`      | int  | 1       | Page number        |
| `page_size` | int  | 20      | Items per page     |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Leave types retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "code": "annual",
        "name": "Annual Leave",
        "icon": "🏖️",
        "category": "earned",
        "paid_leave": true,
        "requires_documentation": false,
        "is_active": true,
        "description": "Annual paid leave"
      },
      {
        "id": 2,
        "code": "sick",
        "name": "Sick Leave",
        "icon": "🤒",
        "category": "statutory",
        "paid_leave": true,
        "requires_documentation": true,
        "is_active": true,
        "description": "Medical/sick leave"
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 20,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

---

## 4. Get Holidays

Returns holidays for the specified year. Used for leave day calculation and calendar display.

```
GET /api/v1/leave/holidays
```

**Auth:** Any authenticated user

| Query Param | Type | Default      | Description |
|-------------|------|--------------|-------------|
| `year`      | int  | current year | Year        |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Holidays retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "New Year's Day",
      "date": "2026-01-01T00:00:00Z",
      "type": "Public",
      "is_floater": false,
      "description": "New Year Holiday"
    },
    {
      "id": 5,
      "name": "Union Day",
      "date": "2026-04-26T00:00:00Z",
      "type": "Public",
      "is_floater": false,
      "description": "Tanzania Union Day"
    }
  ]
}
```

---

## 5. Get Policy Guidelines

Returns leave policy guidelines for employee reference.

```
GET /api/v1/leave/policies/guidelines
```

**Auth:** Any authenticated user

---

## 6. Calculate Leave Days

Calculate working days between two dates, excluding weekends and holidays. Call this as the user selects dates in the form for real-time feedback.

```
POST /api/v1/leave/calculate-days
```

**Auth:** Any authenticated user

### Request Body

```json
{
  "from_date": "2026-02-15",
  "to_date": "2026-02-20",
  "half_day": false
}
```

| Field       | Type    | Required | Description                   |
|-------------|---------|----------|-------------------------------|
| `from_date` | string  | Yes      | Start date (`YYYY-MM-DD`)     |
| `to_date`   | string  | Yes      | End date (`YYYY-MM-DD`)       |
| `half_day`  | boolean | No       | Half-day request (default: false) |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Days calculated successfully",
  "data": {
    "total_days": 4,
    "working_days": 4,
    "weekends": 2,
    "holidays": 0,
    "half_day": false
  }
}
```

---

## 7. Submit Leave Application

Submit a new leave request. The primary endpoint when the employee clicks "Submit Request".

```
POST /api/v1/leave/applications
```

**Auth:** Any authenticated employee

### Request Body

```json
{
  "leave_type_code": "annual",
  "from_date": "2026-02-15",
  "to_date": "2026-02-20",
  "half_day": false,
  "reason": "Family vacation",
  "work_delegated_to": "EMP045",
  "handover_notes": "Please handle the weekly report and client meeting on Feb 18",
  "application_date": "2026-02-06",
  "leave_period_year": 2026,
  "reporting_back_date": "2026-02-21",
  "employee_contact_number": "+255712345678",
  "emergency_contact_person": "Jane Doe",
  "emergency_contact_number": "+255712345679",
  "save_as_draft": false,
  "employee_signature": "data:image/png;base64,...",
  "documents": [
    {
      "file_name": "travel_doc.pdf",
      "file_url": "https://storage.example.com/uploads/travel_doc.pdf",
      "file_type": "application/pdf",
      "file_size": 102400,
      "document_type": "travel_document"
    }
  ],
  "leave_tracker": {
    "start_date": "2026-02-15",
    "end_date": "2026-02-20",
    "number_of_days": 4,
    "previous_days_used": 5,
    "days_remaining_after_request": 10,
    "leave_taken_from_previous_year": 2,
    "leave_balance_from_previous_year": 3
  }
}
```

| Field                      | Type    | Required | Description                                       |
|----------------------------|---------|----------|---------------------------------------------------|
| `leave_type_code`          | string  | Yes      | Leave type code (from `/leave/active-types`)       |
| `from_date`                | string  | Yes      | Start date (`YYYY-MM-DD`)                          |
| `to_date`                  | string  | Yes      | End date (`YYYY-MM-DD`)                            |
| `half_day`                 | boolean | No       | Half-day request (default: false)                  |
| `reason`                   | string  | Yes      | Reason for leave                                   |
| `work_delegated_to`        | string  | No       | Employee ID of duty delegate                       |
| `handover_notes`           | string  | No       | Notes about tasks to hand over                     |
| `application_date`         | string  | Yes      | Date of application (`YYYY-MM-DD`)                 |
| `leave_period_year`        | int     | Yes      | Leave year (e.g. 2026)                             |
| `reporting_back_date`      | string  | Yes      | Expected return date (`YYYY-MM-DD`)                |
| `employee_contact_number`  | string  | Yes      | Contact number during leave                        |
| `emergency_contact_person` | string  | Yes      | Emergency contact name                             |
| `emergency_contact_number` | string  | Yes      | Emergency contact phone                            |
| `save_as_draft`            | boolean | No       | If true, saves as draft (no approval workflow)     |
| `employee_signature`       | string  | No       | Base64 encoded signature                           |
| `documents`                | array   | No       | Supporting documents                               |
| `leave_tracker`            | object  | Yes      | Leave tracking data (see below)                    |

### `leave_tracker` Object

| Field                              | Type   | Description                         |
|------------------------------------|--------|-------------------------------------|
| `start_date`                       | string | Same as `from_date`                 |
| `end_date`                         | string | Same as `to_date`                   |
| `number_of_days`                   | number | Calculated working days             |
| `previous_days_used`               | number | Days already used this year         |
| `days_remaining_after_request`     | number | Available - requested days          |
| `leave_taken_from_previous_year`   | number | Carried forward days used           |
| `leave_balance_from_previous_year` | number | Carried forward balance             |

### What Happens on Submit

1. Leave request is created with status `pending`
2. **Approval workflow is automatically created** with 3 levels:
   - **Level 1:** Head of Department (HOD) — `pending`
   - **Level 2:** HR Department — `pending`
   - **Level 3:** Director/CEO — `not applicable` by default (can be enabled)
3. Leave balance is updated with pending days

### Response — `201 Created`

```json
{
  "success": true,
  "message": "Leave application submitted successfully",
  "data": {
    "id": 42,
    "application_number": "LV-2026-00042",
    "document_number": "HR.FO.04.00",
    "employee_id": "EMP060",
    "leave_type_code": "annual",
    "from_date": "2026-02-15T00:00:00Z",
    "to_date": "2026-02-20T00:00:00Z",
    "total_days": 4,
    "half_day": false,
    "status": "pending",
    "reason": "Family vacation",
    "application_date": "2026-02-06T00:00:00Z",
    "reporting_back_date": "2026-02-21T00:00:00Z",
    "approvals": [
      {
        "id": 1,
        "level": 1,
        "approver_type": "head_of_department",
        "status": "pending",
        "is_applicable": true
      },
      {
        "id": 2,
        "level": 2,
        "approver_type": "hr_department",
        "status": "pending",
        "is_applicable": true
      },
      {
        "id": 3,
        "level": 3,
        "approver_type": "director_ceo",
        "status": "pending",
        "is_applicable": false
      }
    ]
  }
}
```

### Error Responses

| Status | Scenario                                     |
|--------|----------------------------------------------|
| 400    | Insufficient leave balance                   |
| 400    | Invalid date range (from > to)               |
| 400    | Leave period overlaps existing request        |
| 400    | Employee not found                           |
| 401    | Missing or invalid token                     |
| 422    | Validation error (missing required fields)   |

---

## 8. List My Leave Requests

Returns the employee's own leave requests. Used for "My Requests" and "Declined" tabs.

```
GET /api/v1/leave/requests
```

**Auth:** Any authenticated user (returns own requests only; HR/Admin see all)

| Query Param      | Type   | Default | Description                                        |
|------------------|--------|---------|----------------------------------------------------|
| `page`           | int    | 1       | Page number                                        |
| `page_size`      | int    | 10      | Items per page (max 100)                           |
| `status`         | string |         | `pending`, `approved`, `rejected`, `cancelled`, `draft`, `returned_for_info` |
| `leave_type_code`| string |         | Filter by leave type code                          |
| `from_date`      | string |         | Filter from date (`YYYY-MM-DD`)                    |
| `to_date`        | string |         | Filter to date (`YYYY-MM-DD`)                      |

### Frontend Tab Mapping

- **"My Requests" tab:** `status != "rejected"`
- **"Declined" tab:** `status = "rejected"`

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Leave requests retrieved successfully",
  "data": {
    "data": [
      {
        "id": 42,
        "application_number": "LV-2026-00042",
        "employee_id": "EMP060",
        "leave_type_code": "annual",
        "from_date": "2026-02-15T00:00:00Z",
        "to_date": "2026-02-20T00:00:00Z",
        "total_days": 4,
        "half_day": false,
        "status": "pending",
        "reason": "Family vacation",
        "application_date": "2026-02-06T00:00:00Z",
        "approvals": [
          {
            "level": 1,
            "approver_type": "head_of_department",
            "approver_name": null,
            "status": "pending"
          },
          {
            "level": 2,
            "approver_type": "hr_department",
            "approver_name": null,
            "status": "pending"
          }
        ]
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 10,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

---

## 9. Get Leave Request Detail

Full details of a specific leave request, including approval history.

```
GET /api/v1/leave/requests/:request_id
```

**Auth:** Request owner or HR/Admin

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Leave request details retrieved successfully",
  "data": {
    "id": 42,
    "application_number": "LV-2026-00042",
    "document_number": "HR.FO.04.00",
    "employee_id": "EMP060",
    "leave_type_code": "annual",
    "from_date": "2026-02-15T00:00:00Z",
    "to_date": "2026-02-20T00:00:00Z",
    "total_days": 4,
    "half_day": false,
    "status": "pending",
    "reason": "Family vacation",
    "work_delegated_to": "EMP045",
    "handover_notes": "Please handle the weekly report",
    "application_date": "2026-02-06T00:00:00Z",
    "reporting_back_date": "2026-02-21T00:00:00Z",
    "employee_contact_number": "+255712345678",
    "emergency_contact_person": "Jane Doe",
    "emergency_contact_number": "+255712345679",
    "previous_days_used": 5,
    "days_remaining_after_request": 10,
    "leave_taken_from_previous_year": 2,
    "leave_balance_from_previous_year": 3,
    "approvals": [
      {
        "id": 1,
        "level": 1,
        "approver_type": "head_of_department",
        "approver_id": "EMP010",
        "approver_name": "Mike Kamau",
        "status": "approved",
        "approved_at": "2026-02-07T09:00:00Z",
        "remarks": "Approved. Enjoy your vacation.",
        "is_applicable": true
      },
      {
        "id": 2,
        "level": 2,
        "approver_type": "hr_department",
        "approver_id": null,
        "approver_name": null,
        "status": "pending",
        "approved_at": null,
        "remarks": null,
        "is_applicable": true
      }
    ],
    "documents": [],
    "leave_type": {
      "id": 1,
      "code": "annual",
      "name": "Annual Leave",
      "category": "earned"
    }
  }
}
```

---

## 10. Update Leave Request

Update a draft or returned-for-info leave request.

```
PUT /api/v1/leave/requests/:request_id
```

**Auth:** Request owner (only `draft` or `returned_for_info` status)

### Request Body (all fields optional)

```json
{
  "from_date": "2026-02-16",
  "to_date": "2026-02-21",
  "reason": "Updated reason",
  "work_delegated_to": "EMP050",
  "handover_notes": "Updated handover notes",
  "reporting_back_date": "2026-02-22",
  "employee_contact_number": "+255712345000",
  "leave_tracker": {
    "start_date": "2026-02-16",
    "end_date": "2026-02-21",
    "number_of_days": 4,
    "previous_days_used": 5,
    "days_remaining_after_request": 10,
    "leave_taken_from_previous_year": 2,
    "leave_balance_from_previous_year": 3
  }
}
```

### Error Responses

| Status | Scenario                                            |
|--------|-----------------------------------------------------|
| 400    | Only draft or returned_for_info requests can be updated |

---

## 11. Cancel Leave Request

Cancel a pending or approved leave request.

```
POST /api/v1/leave/requests/:request_id/cancel
```

**Auth:** Request owner

### Request Body

```json
{
  "cancellation_reason": "Plans changed, no longer need the leave"
}
```

| Field                | Type   | Required | Description         |
|----------------------|--------|----------|---------------------|
| `cancellation_reason`| string | Yes      | Reason for cancel   |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Leave request cancelled successfully",
  "data": {
    "id": 42,
    "status": "cancelled",
    "cancellation_reason": "Plans changed, no longer need the leave"
  }
}
```

### Business Rules
- Only `pending` or `approved` requests can be cancelled
- Pending balance is restored when a pending request is cancelled

---

## 12. Delete Draft Leave Request

Permanently delete a draft leave request.

```
DELETE /api/v1/leave/requests/:request_id
```

**Auth:** Request owner (`draft` status only)

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Draft leave request deleted successfully",
  "data": null
}
```

---

# Part B — HR/Admin Management APIs

These endpoints require **HR, Admin, or Super Admin** role.

---

## 13. List All Leave Requests (HR View)

Retrieve all leave requests across all employees for review and management.

```
GET /api/v1/leave/admin/requests
```

**Auth:** HR / Admin / Super Admin

| Query Param      | Type   | Default | Description                                        |
|------------------|--------|---------|----------------------------------------------------|
| `page`           | int    | 1       | Page number                                        |
| `page_size`      | int    | 20      | Items per page (max 100)                           |
| `status`         | string |         | `pending`, `approved`, `rejected`, `cancelled`, `draft`, `returned_for_info` |
| `employee_id`    | string |         | Filter by employee ID                              |
| `leave_type_code`| string |         | Filter by leave type                               |
| `from_date`      | string |         | Filter from date                                   |
| `to_date`        | string |         | Filter to date                                     |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Leave requests retrieved successfully",
  "data": {
    "data": [
      {
        "id": 42,
        "application_number": "LV-2026-00042",
        "employee_id": "EMP060",
        "leave_type_code": "annual",
        "from_date": "2026-02-15T00:00:00Z",
        "to_date": "2026-02-20T00:00:00Z",
        "total_days": 4,
        "status": "pending",
        "reason": "Family vacation",
        "approvals": [
          {
            "level": 1,
            "approver_type": "head_of_department",
            "status": "approved",
            "approver_name": "Mike Kamau",
            "approved_at": "2026-02-07T09:00:00Z"
          },
          {
            "level": 2,
            "approver_type": "hr_department",
            "status": "pending"
          }
        ]
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 20,
      "total": 15,
      "total_pages": 1
    }
  }
}
```

---

## 14. Approve Leave Request

Approve a pending leave request at the current approval level.

```
POST /api/v1/leave/requests/:request_id/approve
```

**Auth:** HR / Admin / Super Admin

### Request Body

```json
{
  "action": "approve",
  "approver_type": "head_of_department",
  "approver_name": "Mike Kamau",
  "approver_signature": "data:image/png;base64,...",
  "remarks": "Approved. Enjoy your vacation.",
  "partial_approval": false,
  "partial_days": null,
  "convert_to_unpaid": false,
  "notify_employee": true
}
```

| Field               | Type    | Required | Description                                     |
|---------------------|---------|----------|-------------------------------------------------|
| `action`            | string  | Yes      | Must be `"approve"`                             |
| `approver_type`     | string  | Yes      | `head_of_department`, `hr_department`, `director_ceo` |
| `approver_name`     | string  | Yes      | Name of the approver                            |
| `approver_signature`| string  | Yes      | Base64 encoded signature                        |
| `remarks`           | string  | No       | Approval remarks                                |
| `partial_approval`  | boolean | No       | Approve only partial days                       |
| `partial_days`      | number  | No       | Number of days approved (if partial)            |
| `convert_to_unpaid` | boolean | No       | Convert excess to unpaid leave                  |
| `notify_employee`   | boolean | No       | Send notification to employee                   |

### Approval Flow Logic

The `approver_type` must match the **current pending approval level**:

1. If Level 1 (HOD) is `pending` → `approver_type` must be `"head_of_department"`
2. After HOD approves → Level 2 (HR) becomes the current pending → `approver_type` must be `"hr_department"`
3. After HR approves → If Level 3 is applicable, Director/CEO is next
4. **When ALL applicable levels are approved** → request status changes to `"approved"` and leave balance is updated

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Leave request approved successfully",
  "data": {
    "id": 42,
    "status": "pending",
    "approvals": [
      {
        "level": 1,
        "approver_type": "head_of_department",
        "status": "approved",
        "approver_name": "Mike Kamau",
        "approved_at": "2026-02-07T09:00:00Z"
      },
      {
        "level": 2,
        "approver_type": "hr_department",
        "status": "pending"
      }
    ]
  }
}
```

> **Note:** Status remains `"pending"` until ALL applicable approval levels are completed. When the last level approves, status becomes `"approved"`.

### Error Responses

| Status | Scenario                                      |
|--------|-----------------------------------------------|
| 400    | Only pending/returned requests can be approved|
| 400    | Approver type mismatch (wrong level)          |
| 400    | No pending approvals found                    |

---

## 15. Reject Leave Request

Reject a pending leave request. Rejection at **any level** immediately rejects the entire request.

```
POST /api/v1/leave/requests/:request_id/reject
```

**Auth:** HR / Admin / Super Admin

### Request Body

```json
{
  "approver_type": "head_of_department",
  "approver_name": "Mike Kamau",
  "approver_signature": "data:image/png;base64,...",
  "rejection_reason": "Department is understaffed during that period",
  "remarks": "Please try different dates",
  "notify_employee": true
}
```

| Field              | Type   | Required | Description                    |
|--------------------|--------|----------|--------------------------------|
| `approver_type`    | string | Yes      | Approver's role                |
| `approver_name`    | string | Yes      | Approver's name                |
| `approver_signature`| string| Yes      | Base64 signature               |
| `rejection_reason` | string | Yes      | Reason for rejection           |
| `remarks`          | string | No       | Additional remarks             |
| `notify_employee`  | boolean| No       | Send notification              |

### Business Rules
- Rejection at **any approval level** immediately sets the request status to `"rejected"`
- Pending leave balance is restored

---

## 16. Return for Information

Return a pending leave request to the employee with a request for additional information.

```
POST /api/v1/leave/requests/:request_id/return-for-info
```

**Auth:** HR / Admin / Super Admin

### Request Body

```json
{
  "information_required": "Please attach a medical certificate for sick leave exceeding 3 days",
  "notify_employee": true
}
```

| Field                  | Type    | Required | Description                           |
|------------------------|---------|----------|---------------------------------------|
| `information_required` | string  | Yes      | What information is needed            |
| `notify_employee`      | boolean | No       | Send notification to employee         |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Leave request returned for information successfully",
  "data": {
    "id": 42,
    "status": "returned_for_info",
    "information_required": "Please attach a medical certificate..."
  }
}
```

### Business Rules
- Only `pending` requests can be returned for information
- Employee can then **update** the request (via PUT) and it returns to `pending`

---

# Part C — HR/Admin Configuration APIs

These endpoints manage the leave system configuration and require **HR/Admin** role.

---

## 17. Leave Types CRUD

Manage leave types (Annual, Sick, Emergency, etc.).

| Method | Endpoint                        | Description          |
|--------|---------------------------------|----------------------|
| GET    | `/api/v1/leave/types`           | List all leave types |
| GET    | `/api/v1/leave/types/:type_id`  | Get leave type       |
| POST   | `/api/v1/leave/types`           | Create leave type    |
| PUT    | `/api/v1/leave/types/:type_id`  | Update leave type    |
| DELETE | `/api/v1/leave/types/:type_id`  | Delete leave type    |

### Create Leave Type — Request Body

```json
{
  "code": "annual",
  "name": "Annual Leave",
  "icon": "🏖️",
  "category": "earned",
  "paid_leave": true,
  "requires_documentation": false,
  "is_active": true,
  "description": "Annual paid leave entitlement"
}
```

---

## 18. Leave Policies CRUD

Manage leave policies that define entitlements per leave type.

| Method | Endpoint                              | Description           |
|--------|---------------------------------------|-----------------------|
| GET    | `/api/v1/leave/policies`              | List policies         |
| GET    | `/api/v1/leave/policies/:policy_id`   | Get policy            |
| POST   | `/api/v1/leave/policies`              | Create policy         |
| PUT    | `/api/v1/leave/policies/:policy_id`   | Update policy         |
| DELETE | `/api/v1/leave/policies/:policy_id`   | Delete policy         |

---

## 19. Holidays CRUD

Manage public and company holidays.

| Method | Endpoint                            | Description         |
|--------|-------------------------------------|---------------------|
| GET    | `/api/v1/holidays`                  | List holidays       |
| GET    | `/api/v1/holidays/:holiday_id`      | Get holiday         |
| POST   | `/api/v1/holidays`                  | Create holiday      |
| PUT    | `/api/v1/holidays/:holiday_id`      | Update holiday      |
| DELETE | `/api/v1/holidays/:holiday_id`      | Delete holiday      |

### Create Holiday — Request Body

```json
{
  "name": "Union Day",
  "date": "2026-04-26",
  "type": "Public",
  "is_floater": false,
  "location": ["All Locations"],
  "applicable_for": ["All Employees"],
  "description": "Tanzania Union Day"
}
```

---

# Part D — Leave Calendar (All Users)

Available to all authenticated users.

---

## 20. Leave Calendar

```
GET /api/v1/leave/calendar
```

Returns leave calendar data (who's on leave, when).

---

## 21. Employees on Leave Today

```
GET /api/v1/leave/calendar/today
```

Returns employees who are on leave today.

---

## 22. Weekly Calendar

```
GET /api/v1/leave/calendar/week
```

Returns weekly leave overview.

---

# Approval Workflow

## Flow Diagram

```
Employee Submits Leave Request
         │
         ▼
  ┌──────────────────┐
  │  Status: pending  │
  │  Level 1: HOD     │ ◄── Head of Department
  └────────┬─────────┘
           │
     ┌─────┴─────┐
     │           │
  Approve     Reject ──────► Status: rejected (DONE)
     │
     ▼
  ┌──────────────────┐
  │  Status: pending  │
  │  Level 2: HR      │ ◄── HR Department
  └────────┬─────────┘
           │
     ┌─────┴─────┐
     │           │
  Approve     Reject ──────► Status: rejected (DONE)
     │
     ▼
  ┌──────────────────┐
  │  Level 3: CEO     │ ◄── Director/CEO (if applicable)
  │  (usually N/A)    │
  └────────┬─────────┘
           │
     ALL approved
           │
           ▼
  ┌──────────────────┐
  │ Status: approved  │
  │ Balance updated   │
  └──────────────────┘
```

## Key Rules

1. **Sequential approval:** Level 1 must approve before Level 2 can act
2. **Rejection is immediate:** A rejection at any level immediately rejects the entire request
3. **Level 3 (Director/CEO)** is `is_applicable: false` by default — only needed for special cases
4. **Return for Info** pauses the workflow — employee must update and resubmit
5. **Balance management:** Pending days are tracked; approved days are deducted; cancelled/rejected days are restored

## Request Status Lifecycle

```
draft ──► pending ──► approved
              │            │
              ├──► rejected │
              │            │
              ├──► returned_for_info ──► pending (after employee update)
              │
              └──► cancelled
```

| Status               | Description                                                |
|----------------------|------------------------------------------------------------|
| `draft`              | Saved but not submitted; no approval workflow created       |
| `pending`            | Submitted and awaiting approval(s)                         |
| `approved`           | All applicable approval levels completed                   |
| `rejected`           | Rejected at any approval level                             |
| `returned_for_info`  | Returned to employee for additional information            |
| `partially_approved` | Partial days approved (subset of requested days)           |
| `cancelled`          | Cancelled by the employee                                  |

---

# Data Models

## Leave Types

| Code      | Name              | Category   | Paid | Requires Docs |
|-----------|-------------------|------------|------|---------------|
| `annual`  | Annual Leave      | earned     | Yes  | No            |
| `sick`    | Sick Leave        | statutory  | Yes  | Yes (>3 days) |
| `maternity`| Maternity Leave  | statutory  | Yes  | Yes           |
| `paternity`| Paternity Leave  | statutory  | Yes  | Yes           |
| `comp`    | Compassionate     | emergency  | Yes  | No            |
| `unpaid`  | Unpaid Leave      | other      | No   | No            |

## Approval Levels

| Level | Approver Type         | Description              | Default Applicable |
|-------|-----------------------|--------------------------|--------------------|
| 1     | `head_of_department`  | Head of Department (HOD) | Yes                |
| 2     | `hr_department`       | HR Department            | Yes                |
| 3     | `director_ceo`        | Director / CEO           | No (special cases) |

---

# Quick Reference Table

## Employee Self-Service (`/api/v1/leave`)

| Method | Endpoint                                | Description                       |
|--------|-----------------------------------------|-----------------------------------|
| GET    | `/leave/employee-info`                  | Employee info for form pre-fill   |
| GET    | `/leave/balances`                       | Leave balance statistics          |
| GET    | `/leave/active-types`                   | Active leave types (dropdown)     |
| GET    | `/leave/holidays`                       | Holidays for calendar/calculation |
| GET    | `/leave/policies/guidelines`            | Policy guidelines                 |
| POST   | `/leave/calculate-days`                 | Calculate working days            |
| POST   | `/leave/applications`                   | Submit leave application          |
| GET    | `/leave/requests`                       | List my leave requests            |
| GET    | `/leave/requests/:id`                   | Get request detail                |
| PUT    | `/leave/requests/:id`                   | Update draft/returned request     |
| POST   | `/leave/requests/:id/cancel`            | Cancel request                    |
| DELETE | `/leave/requests/:id`                   | Delete draft request              |

## HR/Admin Management (`/api/v1/leave`)

| Method | Endpoint                                    | Description                       |
|--------|---------------------------------------------|-----------------------------------|
| GET    | `/leave/admin/requests`                     | List ALL leave requests           |
| POST   | `/leave/requests/:id/approve`               | Approve at current level          |
| POST   | `/leave/requests/:id/reject`                | Reject request                    |
| POST   | `/leave/requests/:id/return-for-info`       | Return for more information       |

## HR/Admin Configuration

| Method | Endpoint                       | Description          |
|--------|--------------------------------|----------------------|
| GET    | `/leave/types`                 | List leave types     |
| GET    | `/leave/types/:id`             | Get leave type       |
| POST   | `/leave/types`                 | Create leave type    |
| PUT    | `/leave/types/:id`             | Update leave type    |
| DELETE | `/leave/types/:id`             | Delete leave type    |
| GET    | `/leave/policies`              | List policies        |
| POST   | `/leave/policies`              | Create policy        |
| PUT    | `/leave/policies/:id`          | Update policy        |
| DELETE | `/leave/policies/:id`          | Delete policy        |
| GET    | `/holidays`                    | List holidays        |
| POST   | `/holidays`                    | Create holiday       |
| PUT    | `/holidays/:id`                | Update holiday       |
| DELETE | `/holidays/:id`                | Delete holiday       |

## Leave Calendar (All Users)

| Method | Endpoint                        | Description              |
|--------|---------------------------------|--------------------------|
| GET    | `/leave/calendar`               | Leave calendar data      |
| GET    | `/leave/calendar/today`         | Who's on leave today     |
| GET    | `/leave/calendar/week`          | Weekly leave overview    |

---

## Frontend API Call Flow

### Page Load (Request Leave Page)

```
1. GET /leave/balances?year=2026          → Stat cards + breakdown
2. GET /leave/requests?page=1&page_size=100 → My requests table
3. GET /leave/active-types                → Form dropdown
4. GET /leave/holidays?year=2026          → Calendar holidays
5. GET /leave/employee-info               → Pre-fill form fields
```

### Date Selection in Form

```
POST /leave/calculate-days               → Real-time day calculation
```

### Submit Leave

```
POST /leave/applications                 → Creates request + approval workflow
→ Re-fetch balances and requests on success
```

### HR Approval Flow

```
1. GET /leave/admin/requests?status=pending → List pending requests
2. GET /leave/requests/:id                  → View full details + approval history
3. POST /leave/requests/:id/approve         → Approve (HOD first, then HR)
   OR
   POST /leave/requests/:id/reject          → Reject with reason
   OR
   POST /leave/requests/:id/return-for-info → Ask employee for more info
```

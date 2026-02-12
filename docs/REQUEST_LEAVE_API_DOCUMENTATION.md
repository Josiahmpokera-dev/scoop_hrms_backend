# Request Leave API Documentation

Complete API documentation for the employee-facing **Request Leave** page. This page allows employees to view leave statistics, submit new leave requests (with duty delegation), and track request statuses.

**Base URL:** `/api/v1`

---

## Table of Contents

1. [Leave Balances (Statistics)](#1-leave-balances-statistics)
2. [Leave Types](#2-leave-types)
3. [Calculate Leave Days](#3-calculate-leave-days)
4. [Submit Leave Application](#4-submit-leave-application)
5. [List My Leave Requests](#5-list-my-leave-requests)
6. [Get Leave Request Detail](#6-get-leave-request-detail)
7. [Cancel Leave Request](#7-cancel-leave-request)
8. [Delete Leave Request (Draft)](#8-delete-leave-request-draft)
9. [Employee Info (For Form Pre-fill)](#9-employee-info-for-form-pre-fill)
10. [List Employees (For Delegation Picker)](#10-list-employees-for-delegation-picker)
11. [Frontend Page Specifications](#11-frontend-page-specifications)

---

## 1. Leave Balances (Statistics)

Used to display the **Days Remaining**, **Days Used**, **Pending**, and **Eligibility** stat cards, plus the leave type breakdown section.

```
GET /api/v1/leave/balances
Authorization: Bearer <access_token>
```

**Auth:** Employee (any authenticated user)

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `year` | int | current | Leave year to query |
| `employee_id` | string | self | Employee ID (optional, defaults to authenticated user) |

**Example Request:**

```
GET /api/v1/leave/balances?year=2026
```

**Example Response (200 OK):**

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

| Card | Source |
|------|--------|
| **Days Remaining** | `SUM(data[].available)` |
| **Days Used** | `SUM(data[].used)` |
| **Pending Requests** | `SUM(data[].pending)` |
| **Eligible Year** | If `SUM(data[].entitlement) > 0` → Eligible, else Not Eligible |

---

## 2. Leave Types

Returns active leave types for the leave type dropdown in the request form.

```
GET /api/v1/leave/types
Authorization: Bearer <access_token>
```

**Auth:** Any authenticated user

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `is_active` | boolean | true | Only return active types |
| `page` | int | 1 | Page number |
| `page_size` | int | 20 | Items per page |

**Example Response (200 OK):**

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
        "category": "earned",
        "max_days": 21,
        "is_active": true,
        "carry_forward_allowed": true,
        "requires_approval": true
      },
      {
        "id": 2,
        "code": "sick",
        "name": "Sick Leave",
        "category": "statutory",
        "max_days": 14,
        "is_active": true,
        "carry_forward_allowed": false,
        "requires_approval": true
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

## 3. Calculate Leave Days

Calculate working days between two dates, excluding weekends and holidays. Used for real-time calculation as the employee selects dates in the form.

```
POST /api/v1/leave/calculate-days
Authorization: Bearer <access_token>
```

**Auth:** Any authenticated user

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `leave_type_code` | string | Yes | Leave type code (e.g., `annual`) |
| `from_date` | string | Yes | Start date (YYYY-MM-DD) |
| `to_date` | string | Yes | End date (YYYY-MM-DD) |
| `half_day` | boolean | No | Whether this is a half-day request (default: false) |

**Example Request:**

```json
{
  "leave_type_code": "annual",
  "from_date": "2026-02-15",
  "to_date": "2026-02-20",
  "half_day": false
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Leave days calculated successfully",
  "data": {
    "total_days": 6,
    "working_days": 4,
    "weekends": 2,
    "holidays": 0,
    "half_day": false
  }
}
```

---

## 4. Submit Leave Application

Submit a new leave request. This is the **primary endpoint** called when the employee clicks "Submit Request" in the drawer.

```
POST /api/v1/leave/applications
Authorization: Bearer <access_token>
```

**Auth:** Any authenticated user (employee submitting for themselves)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `leave_type_code` | string | Yes | Leave type code (from types dropdown) |
| `from_date` | string | Yes | Start date (YYYY-MM-DD) |
| `to_date` | string | Yes | End date (YYYY-MM-DD) |
| `half_day` | boolean | No | Half-day request (default: false) |
| `reason` | string | Yes | Reason for leave |
| `work_delegated_to` | string | No | Employee ID of the person covering duties |
| `handover_notes` | string | No | Notes about tasks to be handed over |
| `application_date` | string | Yes | Date of application (today, YYYY-MM-DD) |
| `leave_period_year` | int | Yes | Leave year (e.g., 2026) |
| `reporting_back_date` | string | Yes | Expected return date (YYYY-MM-DD) |
| `employee_contact_number` | string | No | Contact number during leave |
| `emergency_contact_person` | string | No | Emergency contact name |
| `emergency_contact_number` | string | No | Emergency contact phone |
| `save_as_draft` | boolean | No | If true, saves as draft instead of submitting |
| `employee_signature` | string | No | Base64 signature |
| `documents` | array | No | Supporting documents |
| `leave_tracker` | object | Yes | Leave tracking data (see below) |

**`leave_tracker` Object:**

| Field | Type | Description |
|-------|------|-------------|
| `start_date` | string | Same as from_date |
| `end_date` | string | Same as to_date |
| `number_of_days` | number | Calculated working days |
| `previous_days_used` | number | Days already used this year |
| `days_remaining_after_request` | number | Available days minus requested days |
| `leave_taken_from_previous_year` | number | Carried forward days used |
| `leave_balance_from_previous_year` | number | Carried forward balance |

**Example Request:**

```json
{
  "leave_type_code": "annual",
  "from_date": "2026-02-15",
  "to_date": "2026-02-20",
  "half_day": false,
  "reason": "Family vacation",
  "work_delegated_to": "EMP045",
  "handover_notes": "Please handle the weekly report and client meeting on Feb 18",
  "application_date": "2026-01-23",
  "leave_period_year": 2026,
  "reporting_back_date": "2026-02-21",
  "employee_contact_number": "+255712345678",
  "emergency_contact_person": "Jane Doe",
  "emergency_contact_number": "+255712345679",
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

**Example Response (201 Created):**

```json
{
  "success": true,
  "message": "Leave application submitted successfully",
  "data": {
    "id": 42,
    "application_number": "LA-2026-0042",
    "document_number": "DOC-2026-0042",
    "status": "pending",
    "submitted_at": "2026-01-23T10:30:00Z",
    "balance_after_approval": 10,
    "will_result_in_lop": false,
    "lop_days": 0
  }
}
```

**Error Responses:**

| Status | Scenario |
|--------|----------|
| 400 | Insufficient leave balance |
| 400 | Invalid date range (from_date > to_date) |
| 400 | Leave period overlaps with existing request |
| 401 | Missing or invalid access token |
| 422 | Validation error (missing required fields) |

---

## 5. List My Leave Requests

Returns the employee's leave requests. Used to populate the **"My Requests"** and **"Declined"** tabs.

```
GET /api/v1/leave/requests
Authorization: Bearer <access_token>
```

**Auth:** Any authenticated user (returns own requests)

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 100 | Items per page |
| `status` | string | - | Filter: `pending`, `approved`, `rejected`, `cancelled`, `draft` |
| `leave_type_code` | string | - | Filter by leave type |
| `from_date` | string | - | Filter from date (YYYY-MM-DD) |
| `to_date` | string | - | Filter to date (YYYY-MM-DD) |

**Frontend Tab Mapping:**

- **"My Requests" tab**: All requests where `status ≠ "rejected" AND status ≠ "declined"`
- **"Declined" tab**: All requests where `status = "rejected" OR status = "declined"`

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Leave requests retrieved successfully",
  "data": {
    "data": [
      {
        "id": 42,
        "application_number": "LA-2026-0042",
        "employee_id": "EMP060",
        "employee_name": "Anna Juma",
        "department": "Engineering",
        "leave_type_code": "annual",
        "leave_type_name": "Annual Leave",
        "from_date": "2026-02-15",
        "to_date": "2026-02-20",
        "total_days": 4,
        "half_day": false,
        "status": "pending",
        "reason": "Family vacation",
        "applied_date": "2026-01-23",
        "approver": "Mike Kamau",
        "approval_date": null,
        "rejection_reason": null
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 100,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

---

## 6. Get Leave Request Detail

Get full details of a specific leave request. Opened when clicking "View Details" from the table.

```
GET /api/v1/leave/requests/:id
Authorization: Bearer <access_token>
```

**Auth:** Request owner or Admin/HR

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | int | Leave request ID |

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Leave request details retrieved successfully",
  "data": {
    "id": 42,
    "application_number": "LA-2026-0042",
    "document_number": "DOC-2026-0042",
    "employee_id": "EMP060",
    "employee_name": "Anna Juma",
    "department": "Engineering",
    "leave_type_code": "annual",
    "leave_type_name": "Annual Leave",
    "from_date": "2026-02-15",
    "to_date": "2026-02-20",
    "total_days": 4,
    "half_day": false,
    "status": "pending",
    "reason": "Family vacation",
    "applied_date": "2026-01-23",
    "application_date": "2026-01-23",
    "reporting_back_date": "2026-02-21",
    "work_delegated_to": "EMP045",
    "handover_notes": "Please handle the weekly report and client meeting on Feb 18",
    "employee_contact_number": "+255712345678",
    "emergency_contact_person": "Jane Doe",
    "emergency_contact_number": "+255712345679",
    "leave_tracker": {
      "start_date": "2026-02-15",
      "end_date": "2026-02-20",
      "number_of_days": 4,
      "previous_days_used": 5,
      "days_remaining_after_request": 10,
      "leave_taken_from_previous_year": 2,
      "leave_balance_from_previous_year": 3
    },
    "approvals": [
      {
        "approver_type": "supervisor",
        "approver_name": "Mike Kamau",
        "status": "pending",
        "approved_at": null,
        "remarks": null
      }
    ],
    "rejection_reason": null
  }
}
```

---

## 7. Cancel Leave Request

Cancel a pending or draft leave request.

```
POST /api/v1/leave/requests/:id/cancel
Authorization: Bearer <access_token>
```

**Auth:** Request owner

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | int | Leave request ID |

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Leave request cancelled successfully"
}
```

**Error Responses:**

| Status | Scenario |
|--------|----------|
| 400 | Request is already approved/cancelled |
| 403 | Not the owner of this request |

---

## 8. Delete Leave Request (Draft)

Permanently delete a draft leave request.

```
DELETE /api/v1/leave/requests/:id
Authorization: Bearer <access_token>
```

**Auth:** Request owner (draft status only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | int | Leave request ID |

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Leave request deleted successfully"
}
```

---

## 9. Employee Info (For Form Pre-fill)

Returns employee information for pre-filling contact fields in the leave form.

```
GET /api/v1/leave/employee-info
Authorization: Bearer <access_token>
```

**Auth:** Any authenticated user

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Employee info retrieved successfully",
  "data": {
    "employee_id": "EMP060",
    "first_name": "Anna",
    "last_name": "Juma",
    "email": "anna.juma@company.com",
    "department": "Engineering",
    "position": "Software Developer",
    "reporting_manager": "Mike Kamau",
    "phone_number": "+255712345678",
    "emergency_contact_person": "Jane Doe",
    "emergency_contact_number": "+255712345679"
  }
}
```

---

## 10. List Employees (For Delegation Picker)

Returns a list of active employees for the "Delegate Duties To" dropdown in the leave request form.

```
GET /api/v1/employees
Authorization: Bearer <access_token>
```

**Auth:** Any authenticated user

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 200 | Items per page |
| `status` | string | active | Filter by status |
| `search` | string | - | Search by name |

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Employees retrieved successfully",
  "data": {
    "data": [
      {
        "id": 45,
        "employee_id": "EMP045",
        "first_name": "James",
        "last_name": "Techi",
        "email": "james.tech@company.com",
        "department_id": 3,
        "position_id": 7,
        "status": "active"
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 200,
      "total": 50,
      "total_pages": 1
    }
  }
}
```

---

## 11. Frontend Page Specifications

### Page: `/self-service/request-leave`

**Access:** Employee and User roles only

### Layout

```
┌──────────────────────────────────────────────────────┐
│  Header: "Request Leave"              [Refresh] [+Request Leave] │
├──────────────────────────────────────────────────────┤
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────────┐  │
│  │ Days    │ │ Days    │ │ Pending │ │ Eligible │  │
│  │Remaining│ │ Used    │ │Requests │ │Year/Not  │  │
│  │  14     │ │   5     │ │   2     │ │Eligible  │  │
│  └─────────┘ └─────────┘ └─────────┘ └──────────┘  │
├──────────────────────────────────────────────────────┤
│  Leave Balance Breakdown (by type)                    │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌─────────┐│
│  │ Annual   │ │ Sick     │ │Maternity │ │ Comp.   ││
│  │ 14 / 21  │ │ 13 / 14  │ │ 90 / 90  │ │ 5 / 5   ││
│  └──────────┘ └──────────┘ └──────────┘ └─────────┘│
├──────────────────────────────────────────────────────┤
│  Mini Calendar (month view with leave dots)           │
├──────────────────────────────────────────────────────┤
│  Tabs: [My Requests (3)] [Declined (1)]              │
│  ┌──────────────────────────────────────────────────┐│
│  │ MantineReactTable with search                    ││
│  │ Leave Type | From | To | Days | Status | Applied ││
│  │ ...data...                                       ││
│  └──────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────┘
```

### Request Leave Drawer (Right Side)

```
┌────────────────────────────────┐
│ New Leave Request         [X]  │
├────────────────────────────────┤
│ Leave Type *  [Select...]      │
│ Available: 14 days             │
│                                │
│ From Date *  |  To Date *      │
│ [date]       |  [date]         │
│                                │
│ [✓] Half Day                   │
│                                │
│ 4 working days                 │
│ (2 weekends excluded)          │
│                                │
│ Reason *                       │
│ [textarea]                     │
│                                │
│ ── Duty Delegation ──          │
│ Delegate Duties To             │
│ [Select colleague...]          │
│                                │
│ Handover Notes                 │
│ [textarea]                     │
│                                │
│ ── Additional Information ──   │
│ Reporting Back Date            │
│ [date]                         │
│                                │
│ Contact During Leave           │
│ [phone]                        │
│                                │
│ Emergency Contact | Phone      │
│ [name]            | [phone]    │
│                                │
├────────────────────────────────┤
│        [Cancel] [Submit]       │
└────────────────────────────────┘
```

### API Call Flow

1. **Page Load:**
   - `GET /leave/balances?year=2026` → Stats cards + breakdown
   - `GET /leave/requests?page=1&page_size=100` → Tables
   - `GET /leave/types?is_active=true` → Dropdown options
   - `GET /employees?page=1&page_size=200&status=active` → Delegation picker

2. **Date Selection in Form:**
   - `POST /leave/calculate-days` → Real-time day calculation

3. **Submit Leave:**
   - `POST /leave/applications` → Create leave request
   - Re-fetch balances and requests on success

4. **View Details:**
   - `GET /leave/requests/:id` → Full detail in drawer

5. **Cancel/Delete:**
   - `POST /leave/requests/:id/cancel` → Cancel pending request
   - `DELETE /leave/requests/:id` → Delete draft

---

## Quick Reference

| Action | Method | Endpoint | Auth |
|--------|--------|----------|------|
| Get balances | GET | `/leave/balances` | Employee |
| Get leave types | GET | `/leave/types` | Any |
| Calculate days | POST | `/leave/calculate-days` | Any |
| Submit leave | POST | `/leave/applications` | Employee |
| List requests | GET | `/leave/requests` | Employee |
| Request detail | GET | `/leave/requests/:id` | Owner/Admin |
| Cancel request | POST | `/leave/requests/:id/cancel` | Owner |
| Delete draft | DELETE | `/leave/requests/:id` | Owner |
| Employee info | GET | `/leave/employee-info` | Any |
| List employees | GET | `/employees` | Any |

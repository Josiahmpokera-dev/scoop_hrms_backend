# Employee Offboarding API Documentation

## Overview

This document provides complete API specifications for the Employee Offboarding (Separation) management system. The offboarding process manages employee exits, including No Dues Certificate (NDC) clearances, asset returns, exit interviews, and final settlements.

## Base URL

All API endpoints use the base URL: `{{BASE_URL}}/api/v1/employees/offboarding`

**Full base URL:** `http://localhost:8080/api/v1/employees/offboarding`

## Authentication

All requests require Bearer token authentication:
```
Authorization: Bearer <access_token>
```

---

## Table of Contents

1. [Initiate Separation](#1-initiate-separation)
2. [List Offboarding Workflows](#2-list-offboarding-workflows)
3. [Get Offboarding Workflow Details](#3-get-offboarding-workflow-details)
4. [Update Offboarding Workflow](#4-update-offboarding-workflow)
5. [NDC Clearance Management](#5-ndc-clearance-management)
6. [Asset Return Management](#6-asset-return-management)
7. [Exit Interview Management](#7-exit-interview-management)
8. [Final Settlement Management](#8-final-settlement-management)
9. [Complete Offboarding](#9-complete-offboarding)
10. [Get Offboarding Statistics](#10-get-offboarding-statistics)

---

## 1. Initiate Separation

**Endpoint:** `POST /api/v1/employees/offboarding/initiate`

**Description:** Creates a new offboarding workflow for an employee. This starts the separation process and automatically creates initial clearance tasks.

### Request Body

```json
{
  "employee_id": "EMP013",
  "separation_type": "Resignation",
  "resignation_date": "2024-09-15T00:00:00Z",
  "last_working_date": "2024-10-31T00:00:00Z",
  "notice_period_days": 45,
  "reason": "Career Growth - Joining another company",
  "reason_code": "Better Opportunity",
  "exit_interview_required": true,
  "additional_notes": "Employee has provided 45 days notice. Knowledge transfer required.",
  "initiated_by": "EMP002"
}
```

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `employee_id` | string | Yes | Employee ID (e.g., "EMP013") |
| `separation_type` | string | Yes | One of: "Resignation", "Termination", "Retirement", "Mutual Separation", "End of Contract", "Layoff" |
| `resignation_date` | string (ISO 8601) | Yes | Date when resignation/notice was submitted |
| `last_working_date` | string (ISO 8601) | Yes | Employee's last working day |
| `notice_period_days` | integer | Yes | Standard notice period in days (e.g., 30, 45, 60, 90) |
| `reason` | string | Yes | Detailed reason for separation |
| `reason_code` | string | Optional | Predefined reason code: "Better Opportunity", "Higher Studies", "Relocation", "Personal Reasons", "Health Issues", "Performance Issues", "Policy Violation", "Other" |
| `exit_interview_required` | boolean | Optional | Whether exit interview is required (default: true) |
| `additional_notes` | string | Optional | Any additional information |
| `initiated_by` | string | Optional | Employee ID of the person initiating the separation (HR/admin) |

### Success Response (200)

```json
{
  "success": true,
  "message": "Offboarding workflow initiated successfully",
  "data": {
    "offboarding_id": "off-001",
    "employee_id": "EMP013",
    "employee_name": "Thomas Anderson",
    "emp_id": "EMP013",
    "designation": "Senior Developer",
    "department": "Engineering",
    "separation_type": "Resignation",
    "resignation_date": "2024-09-15T00:00:00Z",
    "last_working_date": "2024-10-31T00:00:00Z",
    "notice_period": 45,
    "served_notice_period": 0,
    "reason": "Career Growth - Joining another company",
    "exit_interview_status": "Not Scheduled",
    "clearance_status": "Pending",
    "final_settlement": "Pending",
    "status": "In Progress",
    "clearances": [
      {
        "id": "clear-001",
        "department": "Manager",
        "cleared_by": null,
        "cleared_by_id": null,
        "clearance_date": null,
        "status": "Pending",
        "notes": null,
        "issues": null
      },
      {
        "id": "clear-002",
        "department": "IT",
        "cleared_by": null,
        "cleared_by_id": null,
        "clearance_date": null,
        "status": "Pending",
        "notes": null,
        "issues": null
      },
      {
        "id": "clear-003",
        "department": "HR",
        "cleared_by": null,
        "cleared_by_id": null,
        "clearance_date": null,
        "status": "Pending",
        "notes": null,
        "issues": null
      },
      {
        "id": "clear-004",
        "department": "Finance",
        "cleared_by": null,
        "cleared_by_id": null,
        "clearance_date": null,
        "status": "Pending",
        "notes": null,
        "issues": null
      },
      {
        "id": "clear-005",
        "department": "Assets",
        "cleared_by": null,
        "cleared_by_id": null,
        "clearance_date": null,
        "status": "Pending",
        "notes": null,
        "issues": null
      }
    ],
    "created_at": "2024-09-15T10:30:51.181123+03:00",
    "updated_at": "2024-09-15T10:30:51.181123+03:00"
  }
}
```

### Error Responses

**400 Bad Request** - Invalid input data
```json
{
  "success": false,
  "message": "Validation error",
  "error": "last working date must be after resignation date"
}
```

**404 Not Found** - Employee not found
```json
{
  "success": false,
  "message": "employee not found"
}
```

**400 Bad Request** - Offboarding already exists
```json
{
  "success": false,
  "message": "an active offboarding workflow already exists for this employee"
}
```

---

## 2. List Offboarding Workflows

**Endpoint:** `GET /api/v1/employees/offboarding/workflows`

**Description:** Retrieves a paginated list of all offboarding workflows with filtering and search capabilities.

### Query Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | integer | No | Page number (default: 1) |
| `page_size` | integer | No | Items per page (default: 20, max: 100) |
| `status` | string | No | Filter by status: "In Progress", "Completed", "On Hold" |
| `separation_type` | string | No | Filter by type: "Resignation", "Termination", "Retirement", "Mutual Separation" |
| `clearance_status` | string | No | Filter by clearance status: "Pending", "In Progress", "Completed" |
| `final_settlement` | string | No | Filter by settlement status: "Pending", "Calculated", "Approved", "Paid" |
| `search` | string | No | Search by employee name, employee ID, or designation |
| `department` | string | No | Filter by department |
| `date_from` | string (ISO 8601) | No | Filter by last working date from |
| `date_to` | string (ISO 8601) | No | Filter by last working date to |

### Success Response (200)

```json
{
  "success": true,
  "message": "Offboarding workflows retrieved successfully",
  "data": [
    {
      "offboarding_id": "off-001",
      "employee_id": "EMP013",
      "employee_name": "Thomas Anderson",
      "emp_id": "EMP013",
      "photo_url": "/img/avatars/thumb-13.jpg",
      "designation": "Senior Developer",
      "department": "Engineering",
      "resignation_date": "2024-09-15T00:00:00Z",
      "last_working_date": "2024-10-31T00:00:00Z",
      "notice_period": 45,
      "served_notice_period": 15,
      "separation_type": "Resignation",
      "reason": "Career Growth - Joining another company",
      "exit_interview_status": "Scheduled",
      "exit_interview_date": "2024-10-25T00:00:00Z",
      "clearance_status": "In Progress",
      "final_settlement": "Pending",
      "status": "In Progress",
      "clearance_progress": 20.0,
      "completed_clearances": 1,
      "total_clearances": 5,
      "created_at": "2024-09-15T10:30:51.181123+03:00",
      "updated_at": "2024-09-30T10:30:51.181123+03:00"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 15,
    "total_pages": 1
  }
}
```

---

## 3. Get Offboarding Workflow Details

**Endpoint:** `GET /api/v1/employees/offboarding/workflows/{offboarding_id}`

**Description:** Retrieves complete details of a specific offboarding workflow, including all clearances, asset returns, exit interview details, and settlement information.

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `offboarding_id` | string | Yes | Offboarding workflow ID (e.g., "off-001") |

### Success Response (200)

```json
{
  "success": true,
  "message": "Offboarding workflow retrieved successfully",
  "data": {
    "offboarding_id": "off-001",
    "employee_id": "EMP013",
    "employee_name": "Thomas Anderson",
    "emp_id": "EMP013",
    "photo_url": "/img/avatars/thumb-13.jpg",
    "designation": "Senior Developer",
    "department": "Engineering",
    "resignation_date": "2024-09-15T00:00:00Z",
    "last_working_date": "2024-10-31T00:00:00Z",
    "notice_period": 45,
    "served_notice_period": 45,
    "buyout_amount": null,
    "separation_type": "Resignation",
    "reason": "Career Growth - Joining another company",
    "reason_code": "Better Opportunity",
    "exit_interview_status": "Scheduled",
    "exit_interview_date": "2024-10-25T00:00:00Z",
    "exit_interview_conducted_by": "EMP002",
    "exit_interview_notes": null,
    "clearance_status": "In Progress",
    "final_settlement": "Pending",
    "status": "In Progress",
    "access_revoked_date": null,
    "clearances": [
      {
        "id": "clear-001",
        "department": "Manager",
        "cleared_by": "Sarah Smith",
        "cleared_by_id": "EMP005",
        "clearance_date": "2024-10-15T00:00:00Z",
        "status": "Cleared",
        "notes": "Knowledge transfer completed. Documentation handed over.",
        "issues": null
      },
      {
        "id": "clear-002",
        "department": "IT",
        "cleared_by": "Christopher Miller",
        "cleared_by_id": "EMP008",
        "clearance_date": null,
        "status": "In Progress",
        "notes": "Laptop to be returned on last working day",
        "issues": null
      },
      {
        "id": "clear-003",
        "department": "HR",
        "cleared_by": null,
        "cleared_by_id": null,
        "clearance_date": null,
        "status": "Pending",
        "notes": null,
        "issues": null
      },
      {
        "id": "clear-004",
        "department": "Finance",
        "cleared_by": null,
        "cleared_by_id": null,
        "clearance_date": null,
        "status": "Pending",
        "notes": null,
        "issues": null
      },
      {
        "id": "clear-005",
        "department": "Assets",
        "cleared_by": null,
        "cleared_by_id": null,
        "clearance_date": null,
        "status": "Pending",
        "notes": null,
        "issues": null
      }
    ],
    "assets": [
      {
        "asset_id": 1,
        "asset_code": "LAP-001",
        "asset_type": "laptop",
        "asset_name": "Dell Latitude 5520",
        "serial_number": "DL123456789",
        "assigned_date": "2024-01-15T00:00:00Z",
        "return_status": "Pending",
        "return_date": null,
        "return_condition": null,
        "return_notes": null,
        "expected_return_date": "2024-10-31T00:00:00Z"
      }
    ],
    "final_settlement_data": {
      "settlement_id": null,
      "status": "Pending",
      "calculated_date": null,
      "approved_date": null,
      "paid_date": null,
      "total_amount": null,
      "breakdown": null
    },
    "created_at": "2024-09-15T10:30:51.181123+03:00",
    "updated_at": "2024-09-30T10:30:51.181123+03:00"
  }
}
```

---

## 4. Update Offboarding Workflow

**Endpoint:** `PATCH /api/v1/employees/offboarding/workflows/{offboarding_id}`

**Description:** Updates basic information of an offboarding workflow (dates, reason, etc.). Use specific endpoints for clearances, exit interviews, and settlements.

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `offboarding_id` | string | Yes | Offboarding workflow ID |

### Request Body

```json
{
  "last_working_date": "2024-11-15T00:00:00Z",
  "reason": "Updated reason",
  "additional_notes": "Extended notice period",
  "status": "On Hold"
}
```

### Request Fields (All Optional)

| Field | Type | Description |
|-------|------|-------------|
| `last_working_date` | string (ISO 8601) | Updated last working date |
| `reason` | string | Updated reason for separation |
| `additional_notes` | string | Additional notes |
| `status` | string | Workflow status: "In Progress", "Completed", "On Hold" |

### Success Response (200)

```json
{
  "success": true,
  "message": "Offboarding workflow updated successfully",
  "data": {
    "offboarding_id": "off-001",
    "last_working_date": "2024-11-15T00:00:00Z",
    "status": "On Hold",
    "updated_at": "2024-09-30T15:30:51.181123+03:00"
  }
}
```

---

## 5. NDC Clearance Management

### 5.1 Get Clearances for Workflow

**Endpoint:** `GET /api/v1/employees/offboarding/workflows/{offboarding_id}/clearances`

**Description:** Retrieves all clearance tasks for a specific offboarding workflow.

### Success Response (200)

```json
{
  "success": true,
  "message": "Clearances retrieved successfully",
  "data": [
    {
      "id": "clear-001",
      "offboarding_id": "off-001",
      "department": "Manager",
      "cleared_by": "Sarah Smith",
      "cleared_by_id": "EMP005",
      "clearance_date": "2024-10-15T00:00:00Z",
      "status": "Cleared",
      "notes": "Knowledge transfer completed. Documentation handed over.",
      "issues": null,
      "created_at": "2024-09-15T10:30:51.181123+03:00",
      "updated_at": "2024-10-15T10:30:51.181123+03:00"
    }
  ]
}
```

### 5.2 Update Clearance Status

**Endpoint:** `PATCH /api/v1/employees/offboarding/workflows/{offboarding_id}/clearances/{clearance_id}`

**Description:** Updates the status of a specific clearance (Manager, IT, HR, Finance, Assets).

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `offboarding_id` | string | Yes | Offboarding workflow ID |
| `clearance_id` | string | Yes | Clearance ID (e.g., "clear-001") |

### Request Body

```json
{
  "status": "Cleared",
  "cleared_by_id": "EMP005",
  "notes": "Knowledge transfer completed. Documentation handed over.",
  "issues": null
}
```

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `status` | string | Yes | Clearance status: "Pending", "Cleared", "Issues" |
| `cleared_by_id` | string | Required if status is "Cleared" | Employee ID of person clearing |
| `notes` | string | Optional | Notes about the clearance |
| `issues` | string | Required if status is "Issues" | Description of issues preventing clearance |

### Success Response (200)

```json
{
  "success": true,
  "message": "Clearance updated successfully",
  "data": {
    "id": "clear-001",
    "status": "Cleared",
    "cleared_by": "Sarah Smith",
    "cleared_by_id": "EMP005",
    "clearance_date": "2024-10-15T10:30:51.181123+03:00",
    "notes": "Knowledge transfer completed. Documentation handed over.",
    "updated_at": "2024-10-15T10:30:51.181123+03:00"
  }
}
```

### Clearance Department Values

- `Manager` - Manager/Reporting Manager clearance
- `IT` - IT department clearance (system access, equipment return)
- `HR` - HR department clearance
- `Finance` - Finance department clearance (outstanding payments, loans)
- `Assets` - Assets/Facilities clearance (physical assets return)

---

## 6. Asset Return Management

### 6.1 Get Assets for Offboarding Employee

**Endpoint:** `GET /api/v1/employees/offboarding/workflows/{offboarding_id}/assets`

**Description:** Retrieves all assets assigned to the employee that need to be returned.

### Success Response (200)

```json
{
  "success": true,
  "message": "Assets retrieved successfully",
  "data": [
    {
      "asset_id": 1,
      "asset_code": "LAP-001",
      "asset_type": "laptop",
      "asset_name": "Dell Latitude 5520",
      "brand": "Dell",
      "model": "Latitude 5520",
      "serial_number": "DL123456789",
      "asset_tag": "LAP001",
      "assigned_date": "2024-01-15T00:00:00Z",
      "return_status": "Pending",
      "return_date": null,
      "return_condition": null,
      "return_notes": null,
      "expected_return_date": "2024-10-31T00:00:00Z"
    }
  ]
}
```

### 6.2 Record Asset Return

**Endpoint:** `POST /api/v1/employees/offboarding/workflows/{offboarding_id}/assets/{asset_id}/return`

**Description:** Records the return of an asset. This automatically updates the Assets clearance status if all assets are returned.

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `offboarding_id` | string | Yes | Offboarding workflow ID |
| `asset_id` | integer | Yes | Asset ID |

### Request Body

```json
{
  "return_date": "2024-10-31T00:00:00Z",
  "return_condition": "Good",
  "return_notes": "Asset returned in good condition. Minor scratches on screen.",
  "returned_by_id": "EMP008"
}
```

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `return_date` | string (ISO 8601) | Yes | Date when asset was returned |
| `return_condition` | string | Yes | Condition: "Excellent", "Good", "Fair", "Poor", "Damaged" |
| `return_notes` | string | Optional | Notes about the return |
| `returned_by_id` | string | Yes | Employee ID of person receiving the asset (usually IT/Assets team) |

### Success Response (200)

```json
{
  "success": true,
  "message": "Asset return recorded successfully",
  "data": {
    "asset_id": 1,
    "asset_code": "LAP-001",
    "return_status": "Returned",
    "return_date": "2024-10-31T00:00:00Z",
    "return_condition": "Good",
    "return_notes": "Asset returned in good condition. Minor scratches on screen.",
    "returned_by": "Christopher Miller",
    "returned_by_id": "EMP008",
    "updated_at": "2024-10-31T10:30:51.181123+03:00"
  }
}
```

### 6.3 Mark Asset as Not Returned / Issue

**Endpoint:** `POST /api/v1/employees/offboarding/workflows/{offboarding_id}/assets/{asset_id}/issue`

**Description:** Marks an asset as having issues (not returned, damaged, missing, etc.). This will set the Assets clearance status to "Issues".

### Request Body

```json
{
  "issue_type": "Not Returned",
  "issue_description": "Employee claims asset was lost. Investigation required.",
  "reported_by_id": "EMP008"
}
```

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `issue_type` | string | Yes | One of: "Not Returned", "Damaged", "Missing", "Stolen", "Other" |
| `issue_description` | string | Yes | Detailed description of the issue |
| `reported_by_id` | string | Yes | Employee ID of person reporting the issue |

### Success Response (200)

```json
{
  "success": true,
  "message": "Asset issue recorded successfully",
  "data": {
    "asset_id": 1,
    "asset_code": "LAP-001",
    "return_status": "Issue",
    "issue_type": "Not Returned",
    "issue_description": "Employee claims asset was lost. Investigation required.",
    "reported_by": "Christopher Miller",
    "reported_by_id": "EMP008",
    "updated_at": "2024-10-31T10:30:51.181123+03:00"
  }
}
```

---

## 7. Exit Interview Management

### 7.1 Schedule Exit Interview

**Endpoint:** `POST /api/v1/employees/offboarding/workflows/{offboarding_id}/exit-interview/schedule`

**Description:** Schedules an exit interview for the offboarding employee.

### Request Body

```json
{
  "interview_date": "2024-10-25T14:00:00Z",
  "interviewer_id": "EMP002",
  "location": "Conference Room A",
  "notes": "Focus on feedback about company culture and management"
}
```

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `interview_date` | string (ISO 8601) | Yes | Date and time of exit interview |
| `interviewer_id` | string | Yes | Employee ID of the interviewer (usually HR) |
| `location` | string | Optional | Location of the interview |
| `notes` | string | Optional | Additional notes about the interview |

### Success Response (200)

```json
{
  "success": true,
  "message": "Exit interview scheduled successfully",
  "data": {
    "offboarding_id": "off-001",
    "exit_interview_status": "Scheduled",
    "exit_interview_date": "2024-10-25T14:00:00Z",
    "interviewer": "Jane Wilson",
    "interviewer_id": "EMP002",
    "location": "Conference Room A",
    "notes": "Focus on feedback about company culture and management",
    "updated_at": "2024-09-30T10:30:51.181123+03:00"
  }
}
```

### 7.2 Complete Exit Interview

**Endpoint:** `POST /api/v1/employees/offboarding/workflows/{offboarding_id}/exit-interview/complete`

**Description:** Marks the exit interview as completed and records the interview notes/feedback.

### Request Body

```json
{
  "interview_notes": "Employee provided constructive feedback about communication channels. Overall positive experience. Reason for leaving: Better career opportunity.",
  "feedback_rating": 4,
  "would_recommend": true,
  "conducted_by_id": "EMP002"
}
```

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `interview_notes` | string | Yes | Detailed notes from the exit interview |
| `feedback_rating` | integer | Optional | Rating from 1-5 (1=Very Poor, 5=Excellent) |
| `would_recommend` | boolean | Optional | Would employee recommend company to others |
| `conducted_by_id` | string | Yes | Employee ID of person who conducted the interview |

### Success Response (200)

```json
{
  "success": true,
  "message": "Exit interview completed successfully",
  "data": {
    "offboarding_id": "off-001",
    "exit_interview_status": "Completed",
    "exit_interview_date": "2024-10-25T14:00:00Z",
    "interview_notes": "Employee provided constructive feedback about communication channels. Overall positive experience. Reason for leaving: Better career opportunity.",
    "feedback_rating": 4,
    "would_recommend": true,
    "conducted_by": "Jane Wilson",
    "conducted_by_id": "EMP002",
    "completed_at": "2024-10-25T15:30:51.181123+03:00",
    "updated_at": "2024-10-25T15:30:51.181123+03:00"
  }
}
```

### 7.3 Cancel Exit Interview

**Endpoint:** `POST /api/v1/employees/offboarding/workflows/{offboarding_id}/exit-interview/cancel`

**Description:** Cancels a scheduled exit interview.

### Request Body

```json
{
  "cancellation_reason": "Employee requested to skip exit interview",
  "cancelled_by_id": "EMP002"
}
```

### Success Response (200)

```json
{
  "success": true,
  "message": "Exit interview cancelled successfully",
  "data": {
    "offboarding_id": "off-001",
    "exit_interview_status": "Not Scheduled",
    "exit_interview_date": null,
    "updated_at": "2024-09-30T15:30:51.181123+03:00"
  }
}
```

---

## 8. Final Settlement Management

### 8.1 Calculate Final Settlement

**Endpoint:** `POST /api/v1/employees/offboarding/workflows/{offboarding_id}/settlement/calculate`

**Description:** Calculates the final settlement (Full & Final) amount for the employee. This includes:
- Outstanding salary
- Leave encashment
- Bonus/Incentives
- Deductions (loans, advances, asset damages, etc.)
- Tax calculations

### Request Body (All fields optional - can be calculated automatically)

```json
{
  "outstanding_salary": 1500000.00,
  "leave_encashment": 500000.00,
  "bonus": 0.00,
  "incentives": 0.00,
  "other_earnings": 0.00,
  "outstanding_loans": 0.00,
  "advances": 0.00,
  "asset_damages": 0.00,
  "tax_deductions": 300000.00,
  "other_deductions": 0.00,
  "currency": "TZS",
  "notes": "Settlement calculated based on last working date: 2024-10-31"
}
```

### Success Response (200)

```json
{
  "success": true,
  "message": "Final settlement calculated successfully",
  "data": {
    "settlement_id": "sett-001",
    "offboarding_id": "off-001",
    "employee_id": "EMP013",
    "status": "Calculated",
    "calculated_date": "2024-10-28T10:30:51.181123+03:00",
    "calculated_by": "EMP003",
    "breakdown": {
      "earnings": {
        "outstanding_salary": 1500000.00,
        "leave_encashment": 500000.00,
        "bonus": 0.00,
        "incentives": 0.00,
        "other_earnings": 0.00,
        "total_earnings": 2000000.00
      },
      "deductions": {
        "outstanding_loans": 0.00,
        "advances": 0.00,
        "asset_damages": 0.00,
        "tax_deductions": 300000.00,
        "other_deductions": 0.00,
        "total_deductions": 300000.00
      },
      "net_settlement": 1700000.00,
      "currency": "TZS"
    },
    "notes": "Settlement calculated based on last working date: 2024-10-31",
    "created_at": "2024-10-28T10:30:51.181123+03:00",
    "updated_at": "2024-10-28T10:30:51.181123+03:00"
  }
}
```

### 8.2 Get Settlement Details

**Endpoint:** `GET /api/v1/employees/offboarding/workflows/{offboarding_id}/settlement`

**Description:** Retrieves the final settlement details for an offboarding workflow.

### Success Response (200)

```json
{
  "success": true,
  "message": "Settlement details retrieved successfully",
  "data": {
    "settlement_id": "sett-001",
    "offboarding_id": "off-001",
    "employee_id": "EMP013",
    "status": "Approved",
    "calculated_date": "2024-10-28T10:30:51.181123+03:00",
    "approved_date": "2024-10-29T10:30:51.181123+03:00",
    "approved_by": "John Admin",
    "approved_by_id": "EMP001",
    "paid_date": null,
    "breakdown": {
      "earnings": {
        "outstanding_salary": 1500000.00,
        "leave_encashment": 500000.00,
        "bonus": 0.00,
        "incentives": 0.00,
        "other_earnings": 0.00,
        "total_earnings": 2000000.00
      },
      "deductions": {
        "outstanding_loans": 0.00,
        "advances": 0.00,
        "asset_damages": 0.00,
        "tax_deductions": 300000.00,
        "other_deductions": 0.00,
        "total_deductions": 300000.00
      },
      "net_settlement": 1700000.00,
      "currency": "TZS"
    },
    "payment_method": null,
    "payment_reference": null,
    "notes": "Settlement calculated based on last working date: 2024-10-31",
    "created_at": "2024-10-28T10:30:51.181123+03:00",
    "updated_at": "2024-10-29T10:30:51.181123+03:00"
  }
}
```

### 8.3 Approve Final Settlement

**Endpoint:** `POST /api/v1/employees/offboarding/workflows/{offboarding_id}/settlement/approve`

**Description:** Approves the calculated final settlement. Requires appropriate permissions (usually Finance/HR Manager).

### Request Body

```json
{
  "approved_by_id": "EMP001",
  "approval_notes": "Settlement reviewed and approved. Ready for payment processing."
}
```

### Success Response (200)

```json
{
  "success": true,
  "message": "Final settlement approved successfully",
  "data": {
    "settlement_id": "sett-001",
    "status": "Approved",
    "approved_date": "2024-10-29T10:30:51.181123+03:00",
    "approved_by": "John Admin",
    "approved_by_id": "EMP001",
    "approval_notes": "Settlement reviewed and approved. Ready for payment processing.",
    "updated_at": "2024-10-29T10:30:51.181123+03:00"
  }
}
```

### 8.4 Mark Settlement as Paid

**Endpoint:** `POST /api/v1/employees/offboarding/workflows/{offboarding_id}/settlement/pay`

**Description:** Marks the final settlement as paid. This should be called after the payment has been processed.

### Request Body

```json
{
  "paid_date": "2024-10-31T00:00:00Z",
  "payment_method": "Bank Transfer",
  "payment_reference": "TXN-20241031-001",
  "paid_by_id": "EMP003",
  "payment_notes": "Payment processed via bank transfer. Reference: TXN-20241031-001"
}
```

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `paid_date` | string (ISO 8601) | Yes | Date when payment was made |
| `payment_method` | string | Yes | Payment method: "Bank Transfer", "Cheque", "Cash", "Other" |
| `payment_reference` | string | Optional | Payment reference/transaction ID |
| `paid_by_id` | string | Yes | Employee ID of person processing payment |
| `payment_notes` | string | Optional | Additional payment notes |

### Success Response (200)

```json
{
  "success": true,
  "message": "Settlement marked as paid successfully",
  "data": {
    "settlement_id": "sett-001",
    "status": "Paid",
    "paid_date": "2024-10-31T00:00:00Z",
    "payment_method": "Bank Transfer",
    "payment_reference": "TXN-20241031-001",
    "paid_by": "Robert Johnson",
    "paid_by_id": "EMP003",
    "payment_notes": "Payment processed via bank transfer. Reference: TXN-20241031-001",
    "updated_at": "2024-10-31T10:30:51.181123+03:00"
  }
}
```

---

## 9. Complete Offboarding

**Endpoint:** `POST /api/v1/employees/offboarding/workflows/{offboarding_id}/complete`

**Description:** Completes the offboarding workflow. This can only be done when:
- All clearances are "Cleared" (no "Pending" or "Issues")
- Exit interview is completed (if required)
- Final settlement is "Paid"
- All assets are returned (or issues resolved)

This endpoint will:
- Update employee status to "Separated"
- Revoke system access
- Archive employee data
- Generate final documents (if applicable)

### Request Body

```json
{
  "completed_by_id": "EMP002"
}
```

### Success Response (200)

```json
{
  "success": true,
  "message": "Offboarding workflow completed successfully",
  "data": {
    "offboarding_id": "off-001",
    "status": "Completed",
    "completed_at": "2024-10-31T16:30:51.181123+03:00",
    "completed_by_id": "EMP002",
    "employee_status": "Separated",
    "access_revoked_date": "2024-10-31T16:30:51.181123+03:00",
    "updated_at": "2024-10-31T16:30:51.181123+03:00"
  }
}
```

### Error Response

**400 Bad Request** - Cannot complete offboarding
```json
{
  "success": false,
  "message": "cannot complete offboarding: pending clearances: [IT, Finance]"
}
```

---

## 10. Get Offboarding Statistics

**Endpoint:** `GET /api/v1/employees/offboarding/statistics`

**Description:** Retrieves dashboard statistics for the offboarding page.

### Success Response (200)

```json
{
  "success": true,
  "message": "Statistics retrieved successfully",
  "data": {
    "active_separations": 15,
    "total_clearances": 75,
    "completed_clearances": 45,
    "clearance_completion_percentage": 60.0,
    "scheduled_exit_interviews": 8,
    "completed_exit_interviews": 5,
    "pending_settlements": 12,
    "settlements_by_status": {
      "pending": 5,
      "calculated": 3,
      "approved": 2,
      "paid": 2
    },
    "separations_by_type": {
      "Resignation": 10,
      "Termination": 3,
      "Retirement": 1,
      "Mutual Separation": 1
    },
    "workflows_by_status": {
      "In Progress": 12,
      "Completed": 3,
      "On Hold": 0
    }
  }
}
```

---

## Implementation Notes

1. **Employee Data**: The system uses existing employee data. Use the existing employee ID format (e.g., "EMP013") to link offboarding workflows.

2. **Asset Integration**: The system integrates with the Assets & Equipment API to:
   - Fetch assets assigned to an employee
   - Update asset status when returned
   - Link asset returns to clearance status

3. **Automatic Clearance Creation**: When initiating separation, the system automatically creates clearance tasks for:
   - Manager (assign to employee's reporting manager)
   - IT (assign to IT department head)
   - HR (assign to HR department head)
   - Finance (assign to Finance department head)
   - Assets (assign to Assets/Facilities team)

4. **Clearance Status Calculation**: 
   - "Pending" = No action taken
   - "In Progress" = Some action taken but not cleared
   - "Cleared" = All requirements met
   - "Issues" = Problems preventing clearance

5. **Workflow Status Calculation**:
   - "In Progress" = Workflow active, clearances pending
   - "Completed" = All clearances cleared, settlement paid, workflow closed
   - "On Hold" = Temporarily paused (e.g., employee extended notice period)

6. **Final Settlement Calculation**:
   - Calculate outstanding salary (prorated to last working date)
   - Calculate leave encashment (unused leave days)
   - Check for outstanding loans/advances
   - Calculate tax deductions
   - Include any asset damage charges
   - Net settlement = Total earnings - Total deductions

7. **Access Revocation**: When offboarding is completed, the system:
   - Updates employee status to "terminated"
   - Sets employee.is_active to false
   - Records access_revoked_date

---

**Document Version:** 1.0  
**Last Updated:** 2026-01-16  
**Author:** Backend Development Team

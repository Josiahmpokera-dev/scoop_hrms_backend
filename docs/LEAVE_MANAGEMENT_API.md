# Leave Management API Documentation

This document provides comprehensive API documentation for the Leave Management module.

**Base URL:** All endpoints use `{{BASE_URL}}` which represents `http://localhost:8080/api/v1`

**Authentication:** All endpoints require Bearer token authentication:
```
Authorization: Bearer {access_token}
```

---

## Table of Contents

1. [Leave Types](#1-leave-types)
2. [Leave Policies](#2-leave-policies)
3. [Leave Requests](#3-leave-requests)
4. [Leave Calendar](#4-leave-calendar)
5. [Holidays](#5-holidays)

---

## 1. Leave Types

### 1.1 Get Leave Types
**Endpoint:** `GET {{BASE_URL}}/leave/types`

**Description:** Retrieve all active leave types available for application.

**Authentication:** Required (HR/Admin)

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 10, max: 100) - Items per page
- `is_active` (optional) - Filter by active status (true/false)
- `category` (optional) - Filter by category (Annual, Emergency, Sick, Other)

**cURL Example:**
```bash
curl -X GET "{{BASE_URL}}/leave/types?is_active=true&category=Annual" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Leave types retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "code": "AL",
        "name": "Annual Leave",
        "icon": "✈️",
        "category": "Annual",
        "paid_leave": true,
        "requires_documentation": false,
        "is_active": true,
        "description": "Annual vacation leave",
        "created_at": "2025-01-01T00:00:00+03:00",
        "updated_at": "2025-01-01T00:00:00+03:00"
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 10,
      "total": 8,
      "total_pages": 1
    }
  }
}
```

---

### 1.2 Get Leave Type Details
**Endpoint:** `GET {{BASE_URL}}/leave/types/:type_id`

**Description:** Get detailed information about a specific leave type.

**Authentication:** Required (HR/Admin)

**cURL Example:**
```bash
curl -X GET "{{BASE_URL}}/leave/types/1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### 1.3 Create Leave Type
**Endpoint:** `POST {{BASE_URL}}/leave/types`

**Description:** Create a new leave type (Admin only).

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "code": "BL",
  "name": "Birthday Leave",
  "icon": "🎂",
  "category": "Special",
  "paid_leave": true,
  "requires_documentation": false,
  "is_active": true,
  "description": "Birthday leave entitlement"
}
```

**cURL Example:**
```bash
curl -X POST "{{BASE_URL}}/leave/types" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "BL",
    "name": "Birthday Leave",
    "icon": "🎂",
    "category": "Special",
    "paid_leave": true,
    "is_active": true
  }'
```

---

### 1.4 Update Leave Type
**Endpoint:** `PUT {{BASE_URL}}/leave/types/:type_id`

**Description:** Update an existing leave type (Admin only).

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "name": "Updated Annual Leave",
  "icon": "🏖️",
  "is_active": true,
  "description": "Updated description"
}
```

---

### 1.5 Delete Leave Type
**Endpoint:** `DELETE {{BASE_URL}}/leave/types/:type_id`

**Description:** Delete a leave type (Admin only).

**Authentication:** Required (HR/Admin)

---

## 2. Leave Policies

### 2.1 Get Leave Policies
**Endpoint:** `GET {{BASE_URL}}/leave/policies`

**Description:** Get all leave policies with pagination.

**Authentication:** Required (HR/Admin)

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 10, max: 100) - Items per page
- `country` (optional) - Filter by country
- `leave_type_code` (optional) - Filter by leave type
- `is_active` (optional) - Filter by active status

**cURL Example:**
```bash
curl -X GET "{{BASE_URL}}/leave/policies?country=Tanzania&leave_type_code=AL" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### 2.2 Get Leave Policy Details
**Endpoint:** `GET {{BASE_URL}}/leave/policies/:policy_id`

**Description:** Get detailed information about a specific leave policy.

**Authentication:** Required (HR/Admin)

---

### 2.3 Get Policy Guidelines
**Endpoint:** `GET {{BASE_URL}}/leave/policies/guidelines`

**Description:** Get policy guidelines for a specific leave type.

**Authentication:** Required (Employee)

**Query Parameters:**
- `leave_type_code` (required) - Leave type code (e.g., "AL", "SL")
- `country` (optional, default: "Tanzania") - Country code

**cURL Example:**
```bash
curl -X GET "{{BASE_URL}}/leave/policies/guidelines?leave_type_code=AL&country=Tanzania" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Policy guidelines retrieved successfully",
  "data": {
    "leave_type": "AL",
    "minimum_notice_days": 3,
    "maximum_days_per_request": 14,
    "carry_forward_limit": 14,
    "carry_forward_expiry": "2026-12-31",
    "half_day_allowed": true,
    "requires_documentation": false,
    "documentation_types": [],
    "guidelines": [
      "Minimum 3 days notice required",
      "Carry forward max 14 days",
      "Can be applied in half-days"
    ]
  }
}
```

---

### 2.4 Create Leave Policy
**Endpoint:** `POST {{BASE_URL}}/leave/policies`

**Description:** Create a new leave policy (Admin only).

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "policy_name": "Tanzania Sick Leave Policy",
  "country": "Tanzania",
  "leave_type_code": "SL",
  "entitlement": 63,
  "accrual_frequency": "Annual",
  "proration_on_join": true,
  "proration_on_exit": true,
  "carry_forward": false,
  "encashment_allowed": false,
  "negative_balance_allowed": false,
  "half_day_allowed": false,
  "minimum_notice_days": 0,
  "maximum_days_per_request": 30,
  "is_active": true
}
```

---

### 2.5 Update Leave Policy
**Endpoint:** `PUT {{BASE_URL}}/leave/policies/:policy_id`

**Description:** Update an existing leave policy (Admin only).

**Authentication:** Required (HR/Admin)

---

### 2.6 Delete Leave Policy
**Endpoint:** `DELETE {{BASE_URL}}/leave/policies/:policy_id`

**Description:** Delete a leave policy (Admin only).

**Authentication:** Required (HR/Admin)

---

## 3. Leave Requests

### 3.1 Get Employee Information
**Endpoint:** `GET {{BASE_URL}}/leave/employee-info`

**Description:** Get employee information for auto-filling the leave application form.

**Authentication:** Required (Employee)

**Query Parameters:**
- `employee_id` (optional) - Employee ID (defaults to current user)

**cURL Example:**
```bash
curl -X GET "{{BASE_URL}}/leave/employee-info" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Employee information retrieved successfully",
  "data": {
    "employee_id": "EMP001",
    "first_name": "JOHN",
    "middle_name": "MICHAEL",
    "last_name": "DOE",
    "job_position": "Senior Software Engineer",
    "department": "Engineering",
    "contact_number": "+255712345678",
    "emergency_contact_person": "Jane Doe",
    "emergency_contact_number": "+255712345679"
  }
}
```

---

### 3.2 Get Employee Leave Balances
**Endpoint:** `GET {{BASE_URL}}/leave/balances`

**Description:** Get current leave balances for the logged-in employee.

**Authentication:** Required (Employee)

**Query Parameters:**
- `employee_id` (optional) - Employee ID (defaults to current user)
- `year` (optional) - Year for balance calculation (defaults to current year)

**cURL Example:**
```bash
curl -X GET "{{BASE_URL}}/leave/balances?year=2026" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Leave balances retrieved successfully",
  "data": [
    {
      "leave_type": "AL",
      "leave_type_name": "Annual Leave",
      "entitlement": 28,
      "used": 12,
      "pending": 3,
      "available": 13,
      "carried_forward": 5,
      "expires_on": "2026-12-31"
    }
  ]
}
```

---

### 3.3 Calculate Leave Days
**Endpoint:** `POST {{BASE_URL}}/leave/calculate-days`

**Description:** Calculate working days for a leave period (excluding weekends and holidays).

**Authentication:** Required (Employee)

**Request Body:**
```json
{
  "from_date": "2026-02-15",
  "to_date": "2026-02-20",
  "half_day": false,
  "employee_id": "EMP001"
}
```

**cURL Example:**
```bash
curl -X POST "{{BASE_URL}}/leave/calculate-days" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_date": "2026-02-15",
    "to_date": "2026-02-20",
    "half_day": false
  }'
```

**Response:**
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

### 3.4 Submit Leave Application
**Endpoint:** `POST {{BASE_URL}}/leave/applications`

**Description:** Submit a new leave application. Includes all fields from the manual leave application form (HR.FO.04.00).

**Authentication:** Required (Employee)

**Request Body:**
```json
{
  "leave_type_code": "AL",
  "from_date": "2026-02-15",
  "to_date": "2026-02-20",
  "half_day": false,
  "reason": "Family vacation",
  "work_delegated_to": "EMP011",
  "handover_notes": "Please handle urgent client queries",
  "application_date": "2026-01-26",
  "leave_period_year": 2026,
  "reporting_back_date": "2026-02-21",
  "employee_contact_number": "+255712345678",
  "emergency_contact_person": "Jane Doe",
  "emergency_contact_number": "+255712345679",
  "employee_signature": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
  "leave_tracker": {
    "start_date": "2026-02-15",
    "end_date": "2026-02-20",
    "number_of_days": 4,
    "previous_days_used": 12,
    "days_remaining_after_request": 9,
    "leave_taken_from_previous_year": 0,
    "leave_balance_from_previous_year": 5
  }
}
```

**cURL Example:**
```bash
curl -X POST "{{BASE_URL}}/leave/applications" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "leave_type_code": "AL",
    "from_date": "2026-02-15",
    "to_date": "2026-02-20",
    "reason": "Family vacation",
    "application_date": "2026-01-26",
    "leave_period_year": 2026,
    "reporting_back_date": "2026-02-21",
    "employee_contact_number": "+255712345678",
    "emergency_contact_person": "Jane Doe",
    "emergency_contact_number": "+255712345679",
    "leave_tracker": {
      "start_date": "2026-02-15",
      "end_date": "2026-02-20",
      "number_of_days": 4,
      "previous_days_used": 12,
      "days_remaining_after_request": 9
    }
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Leave application submitted successfully",
  "data": {
    "id": 123,
    "application_number": "LV-2026-00123",
    "document_number": "HR.FO.04.00",
    "status": "pending",
    "submitted_at": "2026-01-26T10:30:00+03:00"
  }
}
```

---

### 3.5 Get Leave Requests
**Endpoint:** `GET {{BASE_URL}}/leave/requests`

**Description:** Get list of leave requests with filtering and pagination.

**Authentication:** Required (Employee/HR/Admin)

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 10, max: 100) - Items per page
- `status` (optional) - Filter by status (draft, pending, approved, rejected, cancelled)
- `employee_id` (optional) - Filter by employee ID
- `leave_type_code` (optional) - Filter by leave type
- `from_date` (optional) - Filter from date
- `to_date` (optional) - Filter to date

**cURL Example:**
```bash
curl -X GET "{{BASE_URL}}/leave/requests?status=pending&page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### 3.6 Get Leave Request Details
**Endpoint:** `GET {{BASE_URL}}/leave/requests/:request_id`

**Description:** Get detailed information about a specific leave request.

**Authentication:** Required (Employee/HR/Admin)

---

### 3.7 Approve Leave Request
**Endpoint:** `POST {{BASE_URL}}/leave/requests/:request_id/approve`

**Description:** Approve a leave request (HR/Admin only).

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "action": "approve",
  "approver_type": "head_of_department",
  "approver_name": "Manager Name",
  "approver_signature": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
  "remarks": "Approved. Enjoy your vacation!",
  "notify_employee": true
}
```

**cURL Example:**
```bash
curl -X POST "{{BASE_URL}}/leave/requests/123/approve" \
  -H "Authorization: Bearer HR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "approve",
    "approver_type": "head_of_department",
    "approver_name": "Manager Name",
    "approver_signature": "data:image/png;base64,...",
    "remarks": "Approved"
  }'
```

---

### 3.8 Reject Leave Request
**Endpoint:** `POST {{BASE_URL}}/leave/requests/:request_id/reject`

**Description:** Reject a leave request (HR/Admin only).

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "approver_type": "head_of_department",
  "approver_name": "Manager Name",
  "approver_signature": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
  "rejection_reason": "Insufficient notice period. Minimum 3 days required.",
  "remarks": "Please reapply with proper notice period.",
  "notify_employee": true
}
```

---

### 3.9 Cancel Leave Request
**Endpoint:** `POST {{BASE_URL}}/leave/requests/:request_id/cancel`

**Description:** Cancel a leave request (only for pending or approved status).

**Authentication:** Required (Employee)

**Request Body:**
```json
{
  "cancellation_reason": "Change of plans",
  "notify_approver": true
}
```

---

### 3.10 Delete Draft Leave Request
**Endpoint:** `DELETE {{BASE_URL}}/leave/requests/:request_id`

**Description:** Delete a draft leave request.

**Authentication:** Required (Employee)

---

## 4. Leave Calendar

### 4.1 Get Leave Calendar
**Endpoint:** `GET {{BASE_URL}}/leave/calendar`

**Description:** Get leave calendar data for a specific month/year.

**Authentication:** Required (All authenticated users)

**Query Parameters:**
- `year` (required) - Year (e.g., 2026)
- `month` (required) - Month (1-12)
- `department_id` (optional) - Filter by department
- `employee_id` (optional) - Filter by specific employee

**cURL Example:**
```bash
curl -X GET "{{BASE_URL}}/leave/calendar?year=2026&month=2" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Leave calendar retrieved successfully",
  "data": {
    "year": 2026,
    "month": 2,
    "leaves": [
      {
        "date": "2026-02-15",
        "leaves": [
          {
            "id": 123,
            "employee_id": "EMP001",
            "employee_name": "John Doe",
            "leave_type_code": "AL",
            "from_date": "2026-02-15",
            "to_date": "2026-02-20"
          }
        ]
      }
    ],
    "holidays": [
      {
        "date": "2026-02-07",
        "name": "Saba Saba Day",
        "type": "Public"
      }
    ],
    "weekends": [
      "2026-02-01",
      "2026-02-02"
    ]
  }
}
```

---

### 4.2 Get Employees on Leave Today
**Endpoint:** `GET {{BASE_URL}}/leave/calendar/today`

**Description:** Get list of employees on leave today.

**Authentication:** Required (All authenticated users)

**Query Parameters:**
- `department_id` (optional) - Filter by department
- `location_id` (optional) - Filter by location

**cURL Example:**
```bash
curl -X GET "{{BASE_URL}}/leave/calendar/today" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### 4.3 Get Leave Calendar by Week
**Endpoint:** `GET {{BASE_URL}}/leave/calendar/week`

**Description:** Get leave calendar data for a specific week.

**Authentication:** Required (All authenticated users)

**Query Parameters:**
- `year` (required) - Year
- `week` (required) - Week number (1-52)
- `department_id` (optional) - Filter by department

---

## 5. Holidays

### 5.1 Get Holidays
**Endpoint:** `GET {{BASE_URL}}/holidays`

**Description:** Get list of holidays with filtering and pagination.

**Authentication:** Required (HR/Admin)

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 10, max: 100) - Items per page
- `year` (optional) - Filter by year
- `type` (optional) - Filter by type (Public, Company, Optional, Restricted)
- `is_floater` (optional) - Filter by floater status

**cURL Example:**
```bash
curl -X GET "{{BASE_URL}}/holidays?year=2026&type=Public" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### 5.2 Get Holiday Details
**Endpoint:** `GET {{BASE_URL}}/holidays/:holiday_id`

**Description:** Get detailed information about a specific holiday.

**Authentication:** Required (HR/Admin)

---

### 5.3 Create Holiday
**Endpoint:** `POST {{BASE_URL}}/holidays`

**Description:** Create a new holiday (Admin only).

**Authentication:** Required (HR/Admin)

**Request Body:**
```json
{
  "name": "Independence Day",
  "date": "2026-12-09",
  "type": "Public",
  "is_floater": false,
  "location": ["Dar es Salaam", "Zanzibar"],
  "applicable_for": ["All Employees"],
  "description": "Tanzania Independence Day"
}
```

**cURL Example:**
```bash
curl -X POST "{{BASE_URL}}/holidays" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Independence Day",
    "date": "2026-12-09",
    "type": "Public",
    "location": ["All Locations"],
    "applicable_for": ["All Employees"]
  }'
```

---

### 5.4 Update Holiday
**Endpoint:** `PUT {{BASE_URL}}/holidays/:holiday_id`

**Description:** Update an existing holiday (Admin only).

**Authentication:** Required (HR/Admin)

---

### 5.5 Delete Holiday
**Endpoint:** `DELETE {{BASE_URL}}/holidays/:holiday_id`

**Description:** Delete a holiday (Admin only).

**Authentication:** Required (HR/Admin)

---

## Complete Workflow Example

### Step 1: Employee Gets Leave Types and Balances
```bash
# Get available leave types
GET {{BASE_URL}}/leave/types

# Get leave balances
GET {{BASE_URL}}/leave/balances?year=2026
```

### Step 2: Employee Calculates Leave Days
```bash
POST {{BASE_URL}}/leave/calculate-days
{
  "from_date": "2026-02-15",
  "to_date": "2026-02-20",
  "half_day": false
}
```

### Step 3: Employee Submits Leave Application
```bash
POST {{BASE_URL}}/leave/applications
{
  "leave_type_code": "AL",
  "from_date": "2026-02-15",
  "to_date": "2026-02-20",
  "reason": "Family vacation",
  ...
}
```

### Step 4: HR Approves Leave Request
```bash
POST {{BASE_URL}}/leave/requests/123/approve
{
  "action": "approve",
  "approver_type": "head_of_department",
  "approver_name": "Manager Name",
  "approver_signature": "...",
  "remarks": "Approved"
}
```

---

## Status Values

**Leave Request Status:**
- `draft` - Saved as draft, not yet submitted
- `pending` - Submitted and awaiting approval
- `returned_for_info` - Returned to employee for additional information
- `approved` - Approved by all required approvers
- `partially_approved` - Partially approved (some days converted to unpaid)
- `rejected` - Rejected by an approver
- `cancelled` - Cancelled by employee

---

## Notes

1. **Employee Endpoints:** Employees can only view and manage their own leave requests
2. **HR/Admin Endpoints:** HR and Admin can view all requests and approve/reject them
3. **Approval Workflow:** Supports 3-level approval (Head of Department → HR → Director/CEO)
4. **Leave Balance:** Automatically updated when requests are approved
5. **Holidays:** Excluded from working days calculation
6. **Weekends:** Excluded from working days calculation

---

## Error Responses

**Common Error Responses:**

**400 Bad Request:**
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "from_date is required"
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "User not authenticated"
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "Insufficient permissions"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Leave request not found"
}
```

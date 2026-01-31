# Leave Management API Documentation

This document provides a comprehensive list of API endpoints required for the Leave Management module in ScoopHRMS.

**Form Reference:** This API documentation is based on the manual Leave Application Form (Document No. HR.FO.04.00) used in the office. All form fields from the manual form are included in the API endpoints.

---

## Table of Contents

1. [Apply Leave](#1-apply-leave)
2. [Leave Requests](#2-leave-requests)
3. [Leave Calendar](#3-leave-calendar)
4. [Policies & Types](#4-policies--types)
5. [Holidays](#5-holidays)
6. [Leave Reports](#6-leave-reports)

---

## 1. Apply Leave

### 1.1 Get Leave Types
**Endpoint:** `GET {{BASE_URL}}/leave/types`

**Description:** Retrieve all active leave types available for application. Includes categories matching the manual form: Annual Leave, Emergency Leave, Sick leave, and Other leaves.

**Query Parameters:**
- `is_active` (optional): Filter by active status (true/false)
- `category` (optional): Filter by category (Annual, Emergency, Sick, Other)

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave types retrieved successfully",
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
      "description": "Annual vacation leave"
    },
    {
      "id": 2,
      "code": "SL",
      "name": "Sick Leave",
      "icon": "🏥",
      "category": "Sick",
      "paid_leave": true,
      "requires_documentation": true,
      "is_active": true,
      "description": "Medical leave with certificate requirement"
    },
    {
      "id": 3,
      "code": "EL",
      "name": "Emergency Leave",
      "icon": "🚨",
      "category": "Emergency",
      "paid_leave": false,
      "requires_documentation": false,
      "is_active": true,
      "description": "Emergency leave for urgent situations"
    }
  ]
}
```

### 1.2 Get Employee Information for Form
**Endpoint:** `GET {{BASE_URL}}/leave/employee-info`

**Description:** Get employee information for auto-filling the leave application form (Section 1 fields).

**Query Parameters:**
- `employee_id` (optional): Employee ID (defaults to current user)

**Sample Response:**
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

### 1.3 Get Employee Leave Balances
**Endpoint:** `GET {{BASE_URL}}/leave/balances`

**Description:** Get current leave balances for the logged-in employee. Includes leave tracker information matching the manual form.

**Query Parameters:**
- `employee_id` (optional): Employee ID (defaults to current user)
- `year` (optional): Year for balance calculation (defaults to current year)

**Sample Response:**
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
      "expires_on": "2026-12-31",
      "leave_taken_from_previous_year": 0,
      "leave_balance_from_previous_year": 5,
      "previous_days_used": 12
    },
    {
      "leave_type": "SL",
      "leave_type_name": "Sick Leave",
      "entitlement": 63,
      "used": 8,
      "pending": 0,
      "available": 55,
      "carried_forward": 0,
      "expires_on": null,
      "leave_taken_from_previous_year": 0,
      "leave_balance_from_previous_year": 0,
      "previous_days_used": 8
    }
  ]
}
```

### 1.4 Calculate Leave Days
**Endpoint:** `POST {{BASE_URL}}/leave/calculate-days`

**Description:** Calculate working days for a leave period (excluding weekends and holidays).

**Request Body:**
```json
{
  "from_date": "2026-02-15",
  "to_date": "2026-02-20",
  "half_day": false,
  "employee_id": "EMP001"
}
```

**Sample Response:**
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

### 1.5 Get Leave Policy Guidelines
**Endpoint:** `GET {{BASE_URL}}/leave/policies/guidelines`

**Description:** Get policy guidelines for a specific leave type.

**Query Parameters:**
- `leave_type_code` (required): Leave type code (e.g., "AL", "SL")
- `country` (optional): Country code for country-specific policies

**Sample Response:**
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

### 1.6 Get Team Members for Delegation
**Endpoint:** `GET {{BASE_URL}}/employees/team-members`

**Description:** Get list of team members for work delegation.

**Query Parameters:**
- `department_id` (optional): Filter by department
- `exclude_employee_id` (optional): Exclude specific employee

**Sample Response:**
```json
{
  "success": true,
  "message": "Team members retrieved successfully",
  "data": [
    {
      "employee_id": "EMP011",
      "name": "Ashley Martinez",
      "department": "Engineering",
      "photo": "/img/avatars/thumb-11.jpg"
    },
    {
      "employee_id": "EMP005",
      "name": "Sarah Smith",
      "department": "Engineering",
      "photo": "/img/avatars/thumb-5.jpg"
    }
  ]
}
```

### 1.7 Submit Leave Application
**Endpoint:** `POST {{BASE_URL}}/leave/applications`

**Description:** Submit a new leave application. This endpoint includes all fields from the manual leave application form (HR.FO.04.00).

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
  "documents": [
    {
      "file_name": "medical_certificate.pdf",
      "file_url": "https://storage.example.com/files/cert.pdf",
      "file_type": "application/pdf",
      "file_size": 245760
    }
  ],
  "save_as_draft": false,
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

**Field Descriptions:**
- `leave_type_code` (required): Category of leave (AL, SL, EL, etc.)
- `from_date` (required): Leave start date
- `to_date` (required): Leave end date
- `half_day` (optional): Whether it's a half-day leave
- `reason` (required): Reason for leave
- `work_delegated_to` (optional): Employee ID to delegate work to
- `handover_notes` (optional): Handover instructions
- `documents` (optional): Array of supporting documents
- `save_as_draft` (optional): Save as draft without submission (default: false)
- `application_date` (required): Date of application submission
- `leave_period_year` (required): Year for which leave is being applied
- `reporting_back_date` (required): Expected date to report back to work
- `employee_contact_number` (required): Employee's contact number
- `emergency_contact_person` (required): Name of emergency contact person
- `emergency_contact_number` (required): Emergency contact person's phone number
- `employee_signature` (required): Base64 encoded signature image or signature data
- `leave_tracker` (required): Leave balance tracking information
  - `start_date`: Leave start date
  - `end_date`: Leave end date
  - `number_of_days`: Total number of leave days requested
  - `previous_days_used`: Number of days already used in current period
  - `days_remaining_after_request`: Calculated remaining balance after this request
  - `leave_taken_from_previous_year`: Days taken from previous year's balance
  - `leave_balance_from_previous_year`: Balance carried forward from previous year

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave application submitted successfully",
  "data": {
    "id": 123,
    "application_number": "LV-2026-00123",
    "document_number": "HR.FO.04.00",
    "status": "pending",
    "submitted_at": "2026-01-26T10:30:00+03:00",
    "balance_after_approval": 9,
    "will_result_in_lop": false,
    "lop_days": 0,
    "employee_info": {
      "employee_id": "EMP001",
      "first_name": "JOHN",
      "middle_name": "MICHAEL",
      "last_name": "DOE",
      "job_position": "Senior Software Engineer",
      "department": "Engineering",
      "contact_number": "+255712345678",
      "emergency_contact_person": "Jane Doe",
      "emergency_contact_number": "+255712345679"
    },
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
}
```

### 1.8 Save Leave Application as Draft
**Endpoint:** `POST {{BASE_URL}}/leave/applications/draft`

**Description:** Save leave application as draft without submission. Includes all form fields from manual form.

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
  "employee_signature": null,
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

**Sample Response:**
```json
{
  "success": true,
  "message": "Draft saved successfully",
  "data": {
    "id": 124,
    "draft_id": "DRAFT-2026-00124",
    "document_number": "HR.FO.04.00",
    "saved_at": "2026-01-26T10:30:00+03:00"
  }
}
```

### 1.9 Upload Leave Document
**Endpoint:** `POST {{BASE_URL}}/leave/documents/upload`

**Description:** Upload supporting documents for leave application.

**Request Body (multipart/form-data):**
- `file` (required): File to upload
- `leave_type_code` (required): Leave type code
- `document_type` (required): Type of document (e.g., "medical_certificate", "travel_document")

**Sample Response:**
```json
{
  "success": true,
  "message": "Document uploaded successfully",
  "data": {
    "file_id": "file_123456",
    "file_name": "medical_certificate.pdf",
    "file_url": "https://storage.example.com/files/cert.pdf",
    "file_type": "application/pdf",
    "file_size": 245760,
    "uploaded_at": "2026-01-26T10:30:00+03:00"
  }
}
```

---

## 2. Leave Requests

### 2.1 Get Leave Requests
**Endpoint:** `GET {{BASE_URL}}/leave/requests`

**Description:** Get list of leave requests with filtering and pagination.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10)
- `status` (optional): Filter by status (draft, pending, approved, rejected, cancelled)
- `employee_id` (optional): Filter by employee ID
- `leave_type_code` (optional): Filter by leave type
- `from_date` (optional): Filter from date
- `to_date` (optional): Filter to date
- `view_type` (optional): View type (my_requests, team_requests, pending_approvals)

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave requests retrieved successfully",
  "data": {
    "data": [
      {
        "id": 123,
        "application_number": "LV-2026-00123",
        "employee_id": "EMP001",
        "employee_name": "John Doe",
        "employee_photo": "/img/avatars/thumb-1.jpg",
        "department": "Engineering",
        "leave_type_code": "AL",
        "leave_type_name": "Annual Leave",
        "from_date": "2026-02-15",
        "to_date": "2026-02-20",
        "total_days": 4,
        "half_day": false,
        "status": "pending",
        "reason": "Family vacation",
        "work_delegated_to": "EMP011",
        "work_delegated_to_name": "Ashley Martinez",
        "handover_notes": "Please handle urgent client queries",
        "applied_date": "2026-01-26",
        "approver": null,
        "approval_date": null,
        "rejection_reason": null,
        "documents": [
          {
            "file_id": "file_123456",
            "file_name": "medical_certificate.pdf",
            "file_url": "https://storage.example.com/files/cert.pdf"
          }
        ]
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 10,
      "total": 45,
      "total_pages": 5
    }
  }
}
```

### 2.2 Get Leave Request Details
**Endpoint:** `GET {{BASE_URL}}/leave/requests/:request_id`

**Description:** Get detailed information about a specific leave request including all form fields from the manual form (HR.FO.04.00).

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave request details retrieved successfully",
  "data": {
    "id": 123,
    "application_number": "LV-2026-00123",
    "document_number": "HR.FO.04.00",
    "employee_id": "EMP001",
    "employee_name": "JOHN MICHAEL DOE",
    "employee_first_name": "JOHN",
    "employee_middle_name": "MICHAEL",
    "employee_last_name": "DOE",
    "employee_photo": "/img/avatars/thumb-1.jpg",
    "job_position": "Senior Software Engineer",
    "department": "Engineering",
    "contact_number": "+255712345678",
    "emergency_contact_person": "Jane Doe",
    "emergency_contact_number": "+255712345679",
    "leave_type_code": "AL",
    "leave_type_name": "Annual Leave",
    "from_date": "2026-02-15",
    "to_date": "2026-02-20",
    "total_days": 4,
    "half_day": false,
    "status": "pending",
    "reason": "Family vacation",
    "work_delegated_to": "EMP011",
    "work_delegated_to_name": "Ashley Martinez",
    "handover_notes": "Please handle urgent client queries",
    "application_date": "2026-01-26",
    "leave_period_year": 2026,
    "reporting_back_date": "2026-02-21",
    "employee_signature": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
    "approver": null,
    "approval_date": null,
    "rejection_reason": null,
    "documents": [
      {
        "file_id": "file_123456",
        "file_name": "medical_certificate.pdf",
        "file_url": "https://storage.example.com/files/cert.pdf"
      }
    ],
    "leave_tracker": {
      "start_date": "2026-02-15",
      "end_date": "2026-02-20",
      "number_of_days": 4,
      "previous_days_used": 12,
      "days_remaining_after_request": 9,
      "leave_taken_from_previous_year": 0,
      "leave_balance_from_previous_year": 5
    },
    "approval_workflow": [
      {
        "level": 1,
        "approver_type": "head_of_department",
        "approver_id": "EMP010",
        "approver_name": "MANAGER NAME",
        "approver_signature": null,
        "status": "pending",
        "approved_at": null,
        "remarks": null
      },
      {
        "level": 2,
        "approver_type": "hr_department",
        "approver_id": null,
        "approver_name": null,
        "approver_signature": null,
        "status": "pending",
        "approved_at": null,
        "remarks": null
      },
      {
        "level": 3,
        "approver_type": "director_ceo",
        "approver_id": null,
        "approver_name": null,
        "approver_signature": null,
        "status": "pending",
        "approved_at": null,
        "remarks": null,
        "applicable": false
      }
    ],
    "approval_remarks": null
  }
}
```

### 2.3 Approve Leave Request
**Endpoint:** `POST {{BASE_URL}}/leave/requests/:request_id/approve`

**Description:** Approve a leave request (Admin/Manager only). Includes approval signature and remarks as per manual form Section 2.

**Request Body:**
```json
{
  "action": "approve",
  "approver_type": "head_of_department",
  "approver_name": "Manager Name",
  "approver_signature": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
  "remarks": "Approved. Enjoy your vacation!",
  "partial_approval": false,
  "partial_days": null,
  "convert_to_unpaid": false,
  "notify_employee": true
}
```

**Field Descriptions:**
- `action` (required): Action type ("approve")
- `approver_type` (required): Type of approver ("head_of_department", "hr_department", "director_ceo")
- `approver_name` (required): Name of the approver (to be filled in capital letters)
- `approver_signature` (required): Base64 encoded signature image
- `remarks` (optional): Approval remarks/comments
- `partial_approval` (optional): Whether this is a partial approval
- `partial_days` (optional): Number of days for partial approval
- `convert_to_unpaid` (optional): Convert excess days to unpaid leave
- `notify_employee` (optional): Send notification to employee (default: true)

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave request approved successfully",
  "data": {
    "id": 123,
    "status": "approved",
    "approved_at": "2026-01-26T12:00:00+03:00",
    "approved_by": "EMP010",
    "approver_name": "MANAGER NAME",
    "approver_type": "head_of_department",
    "approver_signature": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
    "remarks": "Approved. Enjoy your vacation!",
    "balance_after_approval": 9,
    "approval_level": 1,
    "next_approver_required": true,
    "next_approver_type": "hr_department",
    "updated_assignment": {
      "assignment_id": 4,
      "date": "2026-02-15",
      "employee_id": "EMP001",
      "status": "on_leave"
    }
  }
}
```

### 2.4 Reject Leave Request
**Endpoint:** `POST {{BASE_URL}}/leave/requests/:request_id/reject`

**Description:** Reject a leave request (Admin/Manager only). Includes approver signature and remarks.

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

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave request rejected successfully",
  "data": {
    "id": 123,
    "status": "rejected",
    "rejected_at": "2026-01-26T12:00:00+03:00",
    "rejected_by": "EMP010",
    "rejector_name": "MANAGER NAME",
    "approver_type": "head_of_department",
    "approver_signature": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
    "rejection_reason": "Insufficient notice period. Minimum 3 days required.",
    "remarks": "Please reapply with proper notice period."
  }
}
```

### 2.5 Partial Approval (Convert to Unpaid)
**Endpoint:** `POST {{BASE_URL}}/leave/requests/:request_id/partial-approve`

**Description:** Partially approve a leave request by converting excess days to unpaid leave.

**Request Body:**
```json
{
  "approved_days": 2,
  "unpaid_days": 2,
  "remarks": "Approved 2 days. Remaining 2 days will be unpaid.",
  "notify_employee": true
}
```

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave request partially approved successfully",
  "data": {
    "id": 123,
    "status": "partially_approved",
    "approved_days": 2,
    "unpaid_days": 2,
    "approved_at": "2026-01-26T12:00:00+03:00",
    "approved_by": "EMP010",
    "remarks": "Approved 2 days. Remaining 2 days will be unpaid."
  }
}
```

### 2.6 Return for Information
**Endpoint:** `POST {{BASE_URL}}/leave/requests/:request_id/return-for-info`

**Description:** Return a leave request to employee for additional information.

**Request Body:**
```json
{
  "information_required": "Please provide medical certificate for sick leave",
  "notify_employee": true
}
```

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave request returned for information",
  "data": {
    "id": 123,
    "status": "returned_for_info",
    "returned_at": "2026-01-26T12:00:00+03:00",
    "returned_by": "EMP010",
    "information_required": "Please provide medical certificate for sick leave"
  }
}
```

### 2.7 Modify Leave Request
**Endpoint:** `PUT {{BASE_URL}}/leave/requests/:request_id`

**Description:** Modify a leave request (only for draft or returned_for_info status). All form fields can be updated.

**Request Body:**
```json
{
  "from_date": "2026-02-16",
  "to_date": "2026-02-21",
  "reason": "Updated reason",
  "work_delegated_to": "EMP005",
  "handover_notes": "Updated handover notes",
  "reporting_back_date": "2026-02-22",
  "employee_contact_number": "+255712345678",
  "emergency_contact_person": "Updated Contact",
  "emergency_contact_number": "+255712345680",
  "leave_tracker": {
    "start_date": "2026-02-16",
    "end_date": "2026-02-21",
    "number_of_days": 4,
    "previous_days_used": 12,
    "days_remaining_after_request": 9,
    "leave_taken_from_previous_year": 0,
    "leave_balance_from_previous_year": 5
  }
}
```

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave request modified successfully",
  "data": {
    "id": 123,
    "from_date": "2026-02-16",
    "to_date": "2026-02-21",
    "total_days": 4,
    "reporting_back_date": "2026-02-22",
    "updated_at": "2026-01-26T12:00:00+03:00"
  }
}
```

### 2.8 Cancel Leave Request
**Endpoint:** `POST {{BASE_URL}}/leave/requests/:request_id/cancel`

**Description:** Cancel a leave request (only for pending or approved status).

**Request Body:**
```json
{
  "cancellation_reason": "Change of plans",
  "notify_approver": true
}
```

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave request cancelled successfully",
  "data": {
    "id": 123,
    "status": "cancelled",
    "cancelled_at": "2026-01-26T12:00:00+03:00",
    "cancellation_reason": "Change of plans"
  }
}
```

### 2.9 Delete Draft Leave Request
**Endpoint:** `DELETE {{BASE_URL}}/leave/requests/:request_id`

**Description:** Delete a draft leave request.

**Sample Response:**
```json
{
  "success": true,
  "message": "Draft leave request deleted successfully"
}
```

### 2.10 Get Leave Request Statistics
**Endpoint:** `GET {{BASE_URL}}/leave/requests/statistics`

**Description:** Get statistics for leave requests (counts by status).

**Query Parameters:**
- `employee_id` (optional): Filter by employee ID
- `department_id` (optional): Filter by department
- `view_type` (optional): View type (my_requests, team_requests, pending_approvals)

**Sample Response:**
```json
{
  "success": true,
  "message": "Statistics retrieved successfully",
  "data": {
    "total_requests": 45,
    "my_requests": 12,
    "pending_approvals": 8,
    "approved": 25,
    "rejected": 5,
    "draft": 3,
    "cancelled": 2,
    "team_requests": 33
  }
}
```

---

## 3. Leave Calendar

### 3.1 Get Leave Calendar
**Endpoint:** `GET {{BASE_URL}}/leave/calendar`

**Description:** Get leave calendar data for a specific month/year.

**Query Parameters:**
- `year` (required): Year (e.g., 2026)
- `month` (required): Month (1-12)
- `department_id` (optional): Filter by department
- `employee_id` (optional): Filter by specific employee
- `view_type` (optional): View type (month, week, list)

**Sample Response:**
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
            "employee_photo": "/img/avatars/thumb-1.jpg",
            "department": "Engineering",
            "leave_type_code": "AL",
            "leave_type_name": "Annual Leave",
            "from_date": "2026-02-15",
            "to_date": "2026-02-20",
            "is_half_day": false
          },
          {
            "id": 124,
            "employee_id": "EMP005",
            "employee_name": "Sarah Smith",
            "employee_photo": "/img/avatars/thumb-5.jpg",
            "department": "Engineering",
            "leave_type_code": "SL",
            "leave_type_name": "Sick Leave",
            "from_date": "2026-02-15",
            "to_date": "2026-02-15",
            "is_half_day": false
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
      "2026-02-02",
      "2026-02-08",
      "2026-02-09"
    ]
  }
}
```

### 3.2 Get Employees on Leave Today
**Endpoint:** `GET {{BASE_URL}}/leave/calendar/today`

**Description:** Get list of employees on leave today.

**Query Parameters:**
- `department_id` (optional): Filter by department
- `location_id` (optional): Filter by location

**Sample Response:**
```json
{
  "success": true,
  "message": "Today's leave list retrieved successfully",
  "data": {
    "date": "2026-01-26",
    "total_on_leave": 5,
    "employees": [
      {
        "id": 123,
        "employee_id": "EMP001",
        "employee_name": "John Doe",
        "employee_photo": "/img/avatars/thumb-1.jpg",
        "department": "Engineering",
        "leave_type_code": "AL",
        "leave_type_name": "Annual Leave",
        "from_date": "2026-01-25",
        "to_date": "2026-01-30",
        "is_half_day": false
      }
    ]
  }
}
```

### 3.3 Get Leave Calendar by Week
**Endpoint:** `GET {{BASE_URL}}/leave/calendar/week`

**Description:** Get leave calendar data for a specific week.

**Query Parameters:**
- `year` (required): Year
- `week` (required): Week number (1-52)
- `department_id` (optional): Filter by department

**Sample Response:**
```json
{
  "success": true,
  "message": "Week calendar retrieved successfully",
  "data": {
    "year": 2026,
    "week": 8,
    "week_start": "2026-02-16",
    "week_end": "2026-02-22",
    "daily_leaves": [
      {
        "date": "2026-02-16",
        "leaves": [...]
      }
    ]
  }
}
```

---

## 4. Policies & Types

### 4.1 Get Leave Types
**Endpoint:** `GET {{BASE_URL}}/leave/types`

**Description:** Get all leave types with pagination.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10)
- `is_active` (optional): Filter by active status
- `category` (optional): Filter by category

**Sample Response:**
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

### 4.2 Get Leave Type Details
**Endpoint:** `GET {{BASE_URL}}/leave/types/:type_id`

**Description:** Get detailed information about a specific leave type.

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave type details retrieved successfully",
  "data": {
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
}
```

### 4.3 Create Leave Type
**Endpoint:** `POST {{BASE_URL}}/leave/types`

**Description:** Create a new leave type (Admin only).

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

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave type created successfully",
  "data": {
    "id": 9,
    "code": "BL",
    "name": "Birthday Leave",
    "icon": "🎂",
    "category": "Special",
    "paid_leave": true,
    "requires_documentation": false,
    "is_active": true,
    "description": "Birthday leave entitlement",
    "created_at": "2026-01-26T10:30:00+03:00"
  }
}
```

### 4.4 Update Leave Type
**Endpoint:** `PUT {{BASE_URL}}/leave/types/:type_id`

**Description:** Update an existing leave type (Admin only).

**Request Body:**
```json
{
  "name": "Updated Annual Leave",
  "icon": "🏖️",
  "is_active": true,
  "description": "Updated description"
}
```

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave type updated successfully",
  "data": {
    "id": 1,
    "code": "AL",
    "name": "Updated Annual Leave",
    "icon": "🏖️",
    "updated_at": "2026-01-26T10:30:00+03:00"
  }
}
```

### 4.5 Delete Leave Type
**Endpoint:** `DELETE {{BASE_URL}}/leave/types/:type_id`

**Description:** Delete a leave type (Admin only).

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave type deleted successfully"
}
```

### 4.6 Get Leave Policies
**Endpoint:** `GET {{BASE_URL}}/leave/policies`

**Description:** Get all leave policies with pagination.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10)
- `country` (optional): Filter by country
- `leave_type_code` (optional): Filter by leave type
- `is_active` (optional): Filter by active status

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave policies retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "policy_name": "Tanzania Annual Leave Policy",
        "country": "Tanzania",
        "leave_type_code": "AL",
        "leave_type_name": "Annual Leave",
        "entitlement": 28,
        "accrual_frequency": "Annual",
        "proration_on_join": true,
        "proration_on_exit": true,
        "carry_forward": true,
        "carry_forward_limit": 14,
        "carry_forward_expiry": "2026-12-31",
        "encashment_allowed": true,
        "encashment_limit": 7,
        "negative_balance_allowed": false,
        "sandwich_rules": true,
        "half_day_allowed": true,
        "minimum_notice_days": 3,
        "maximum_days_per_request": 14,
        "is_active": true,
        "created_at": "2025-01-01T00:00:00+03:00",
        "updated_at": "2025-01-01T00:00:00+03:00"
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 10,
      "total": 12,
      "total_pages": 2
    }
  }
}
```

### 4.7 Get Leave Policy Details
**Endpoint:** `GET {{BASE_URL}}/leave/policies/:policy_id`

**Description:** Get detailed information about a specific leave policy.

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave policy details retrieved successfully",
  "data": {
    "id": 1,
    "policy_name": "Tanzania Annual Leave Policy",
    "country": "Tanzania",
    "leave_type_code": "AL",
    "leave_type_name": "Annual Leave",
    "entitlement": 28,
    "accrual_frequency": "Annual",
    "proration_on_join": true,
    "proration_on_exit": true,
    "carry_forward": true,
    "carry_forward_limit": 14,
    "carry_forward_expiry": "2026-12-31",
    "encashment_allowed": true,
    "encashment_limit": 7,
    "negative_balance_allowed": false,
    "sandwich_rules": true,
    "half_day_allowed": true,
    "minimum_notice_days": 3,
    "maximum_days_per_request": 14,
    "is_active": true,
    "created_at": "2025-01-01T00:00:00+03:00",
    "updated_at": "2025-01-01T00:00:00+03:00"
  }
}
```

### 4.8 Create Leave Policy
**Endpoint:** `POST {{BASE_URL}}/leave/policies`

**Description:** Create a new leave policy (Admin only).

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
  "carry_forward_limit": 0,
  "encashment_allowed": false,
  "encashment_limit": 0,
  "negative_balance_allowed": false,
  "sandwich_rules": false,
  "half_day_allowed": false,
  "minimum_notice_days": 0,
  "maximum_days_per_request": 30,
  "is_active": true
}
```

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave policy created successfully",
  "data": {
    "id": 13,
    "policy_name": "Tanzania Sick Leave Policy",
    "country": "Tanzania",
    "leave_type_code": "SL",
    "entitlement": 63,
    "is_active": true,
    "created_at": "2026-01-26T10:30:00+03:00"
  }
}
```

### 4.9 Update Leave Policy
**Endpoint:** `PUT {{BASE_URL}}/leave/policies/:policy_id`

**Description:** Update an existing leave policy (Admin only).

**Request Body:**
```json
{
  "entitlement": 30,
  "carry_forward_limit": 15,
  "minimum_notice_days": 5,
  "is_active": true
}
```

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave policy updated successfully",
  "data": {
    "id": 1,
    "entitlement": 30,
    "carry_forward_limit": 15,
    "minimum_notice_days": 5,
    "updated_at": "2026-01-26T10:30:00+03:00"
  }
}
```

### 4.10 Delete Leave Policy
**Endpoint:** `DELETE {{BASE_URL}}/leave/policies/:policy_id`

**Description:** Delete a leave policy (Admin only).

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave policy deleted successfully"
}
```

### 4.11 Get Country Packs
**Endpoint:** `GET {{BASE_URL}}/leave/country-packs`

**Description:** Get available country-specific leave policy packs.

**Sample Response:**
```json
{
  "success": true,
  "message": "Country packs retrieved successfully",
  "data": [
    {
      "id": 1,
      "country": "Tanzania",
      "pack_name": "TZ Standard",
      "description": "Tanzania standard leave policy pack",
      "policies": [
        {
          "leave_type_code": "AL",
          "leave_type_name": "Annual Leave",
          "entitlement": 28
        },
        {
          "leave_type_code": "SL",
          "leave_type_name": "Sick Leave",
          "entitlement": 63
        },
        {
          "leave_type_code": "ML",
          "leave_type_name": "Maternity Leave",
          "entitlement": 84
        },
        {
          "leave_type_code": "PL",
          "leave_type_name": "Paternity Leave",
          "entitlement": 3
        }
      ]
    },
    {
      "id": 2,
      "country": "India",
      "pack_name": "India Standard",
      "description": "India standard leave policy pack",
      "policies": [...]
    }
  ]
}
```

### 4.12 Apply Country Pack
**Endpoint:** `POST {{BASE_URL}}/leave/country-packs/:pack_id/apply`

**Description:** Apply a country pack to create/update policies (Admin only).

**Request Body:**
```json
{
  "overwrite_existing": false,
  "notify_employees": true
}
```

**Sample Response:**
```json
{
  "success": true,
  "message": "Country pack applied successfully",
  "data": {
    "pack_id": 1,
    "policies_created": 4,
    "policies_updated": 0,
    "applied_at": "2026-01-26T10:30:00+03:00"
  }
}
```

---

## 5. Holidays

### 5.1 Get Holidays
**Endpoint:** `GET {{BASE_URL}}/holidays`

**Description:** Get list of holidays with filtering and pagination.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10)
- `year` (optional): Filter by year
- `type` (optional): Filter by type (Public, Company, Optional, Restricted)
- `location` (optional): Filter by location
- `is_floater` (optional): Filter by floater status

**Sample Response:**
```json
{
  "success": true,
  "message": "Holidays retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "New Year's Day",
        "date": "2026-01-01",
        "type": "Public",
        "is_floater": false,
        "location": ["All Locations"],
        "applicable_for": ["All Employees"],
        "description": "New Year celebration",
        "created_at": "2025-01-01T00:00:00+03:00",
        "updated_at": "2025-01-01T00:00:00+03:00"
      },
      {
        "id": 2,
        "name": "Company Foundation Day",
        "date": "2026-03-15",
        "type": "Company",
        "is_floater": false,
        "location": ["All Locations"],
        "applicable_for": ["All Employees"],
        "description": "Company anniversary",
        "created_at": "2025-01-01T00:00:00+03:00",
        "updated_at": "2025-01-01T00:00:00+03:00"
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 10,
      "total": 25,
      "total_pages": 3
    }
  }
}
```

### 5.2 Get Holiday Details
**Endpoint:** `GET {{BASE_URL}}/holidays/:holiday_id`

**Description:** Get detailed information about a specific holiday.

**Sample Response:**
```json
{
  "success": true,
  "message": "Holiday details retrieved successfully",
  "data": {
    "id": 1,
    "name": "New Year's Day",
    "date": "2026-01-01",
    "type": "Public",
    "is_floater": false,
    "location": ["All Locations"],
    "applicable_for": ["All Employees"],
    "description": "New Year celebration",
    "created_at": "2025-01-01T00:00:00+03:00",
    "updated_at": "2025-01-01T00:00:00+03:00"
  }
}
```

### 5.3 Create Holiday
**Endpoint:** `POST {{BASE_URL}}/holidays`

**Description:** Create a new holiday (Admin only).

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

**Sample Response:**
```json
{
  "success": true,
  "message": "Holiday created successfully",
  "data": {
    "id": 26,
    "name": "Independence Day",
    "date": "2026-12-09",
    "type": "Public",
    "is_floater": false,
    "location": ["Dar es Salaam", "Zanzibar"],
    "applicable_for": ["All Employees"],
    "description": "Tanzania Independence Day",
    "created_at": "2026-01-26T10:30:00+03:00"
  }
}
```

### 5.4 Update Holiday
**Endpoint:** `PUT {{BASE_URL}}/holidays/:holiday_id`

**Description:** Update an existing holiday (Admin only).

**Request Body:**
```json
{
  "name": "Updated Independence Day",
  "date": "2026-12-09",
  "type": "Public",
  "location": ["All Locations"]
}
```

**Sample Response:**
```json
{
  "success": true,
  "message": "Holiday updated successfully",
  "data": {
    "id": 26,
    "name": "Updated Independence Day",
    "updated_at": "2026-01-26T10:30:00+03:00"
  }
}
```

### 5.5 Delete Holiday
**Endpoint:** `DELETE {{BASE_URL}}/holidays/:holiday_id`

**Description:** Delete a holiday (Admin only).

**Sample Response:**
```json
{
  "success": true,
  "message": "Holiday deleted successfully"
}
```

### 5.6 Bulk Upload Holidays
**Endpoint:** `POST {{BASE_URL}}/holidays/bulk-upload`

**Description:** Bulk upload holidays from CSV file (Admin only).

**Request Body (multipart/form-data):**
- `file` (required): CSV file
- `year` (required): Year for the holidays
- `overwrite_existing` (optional): Overwrite existing holidays (default: false)

**CSV Format:**
```csv
name,date,type,location,applicable_for,description
New Year's Day,2026-01-01,Public,All Locations,All Employees,New Year celebration
Independence Day,2026-12-09,Public,All Locations,All Employees,Tanzania Independence
```

**Sample Response:**
```json
{
  "success": true,
  "message": "Holidays bulk uploaded successfully",
  "data": {
    "total_rows": 25,
    "created": 20,
    "updated": 5,
    "failed": 0,
    "errors": []
  }
}
```

### 5.7 Import Holidays from API
**Endpoint:** `POST {{BASE_URL}}/holidays/import-from-api`

**Description:** Import public holidays from external API (Admin only).

**Request Body:**
```json
{
  "country": "Tanzania",
  "year": 2026,
  "source": "calendarific",
  "api_key": "your_api_key"
}
```

**Sample Response:**
```json
{
  "success": true,
  "message": "Holidays imported successfully",
  "data": {
    "imported": 15,
    "created": 15,
    "skipped": 0
  }
}
```

### 5.8 Get Holiday Statistics
**Endpoint:** `GET {{BASE_URL}}/holidays/statistics`

**Description:** Get holiday statistics for a specific year.

**Query Parameters:**
- `year` (required): Year

**Sample Response:**
```json
{
  "success": true,
  "message": "Holiday statistics retrieved successfully",
  "data": {
    "year": 2026,
    "total_holidays": 25,
    "public_holidays": 15,
    "company_holidays": 8,
    "optional_holidays": 2,
    "restricted_holidays": 0
  }
}
```

---

## 6. Leave Reports

### 6.1 Get Leave Utilization Report
**Endpoint:** `GET {{BASE_URL}}/leave/reports/utilization`

**Description:** Get leave utilization report with department-wise breakdown.

**Query Parameters:**
- `start_date` (required): Start date (YYYY-MM-DD)
- `end_date` (required): End date (YYYY-MM-DD)
- `department_id` (optional): Filter by department
- `employee_id` (optional): Filter by employee
- `leave_type_code` (optional): Filter by leave type

**Sample Response:**
```json
{
  "success": true,
  "message": "Utilization report retrieved successfully",
  "data": {
    "period": {
      "start_date": "2026-01-01",
      "end_date": "2026-01-31"
    },
    "overall_utilization": 68.5,
    "total_entitlement": 840,
    "total_used": 575,
    "total_available": 265,
    "department_wise": [
      {
        "department_id": 1,
        "department_name": "Engineering",
        "utilization_percent": 72,
        "entitlement": 420,
        "used": 302,
        "available": 118,
        "liability": 3500000
      },
      {
        "department_id": 2,
        "department_name": "Sales",
        "utilization_percent": 65,
        "entitlement": 180,
        "used": 117,
        "available": 63,
        "liability": 1800000
      }
    ],
    "employee_wise": [
      {
        "employee_id": "EMP001",
        "employee_name": "John Doe",
        "department": "Engineering",
        "utilization_percent": 75,
        "entitlement": 28,
        "used": 21,
        "available": 7
      }
    ]
  }
}
```

### 6.2 Get Absenteeism Trend Report
**Endpoint:** `GET {{BASE_URL}}/leave/reports/absenteeism-trend`

**Description:** Get absenteeism trend analysis report.

**Query Parameters:**
- `start_date` (required): Start date (YYYY-MM-DD)
- `end_date` (required): End date (YYYY-MM-DD)
- `department_id` (optional): Filter by department
- `group_by` (optional): Group by (month, week, day) - default: month

**Sample Response:**
```json
{
  "success": true,
  "message": "Absenteeism trend report retrieved successfully",
  "data": {
    "period": {
      "start_date": "2026-01-01",
      "end_date": "2026-01-31"
    },
    "overall_absenteeism_rate": 4.2,
    "trend": "decreasing",
    "trend_percent": -1.1,
    "monthly_trends": [
      {
        "month": "2026-01",
        "absenteeism_rate": 4.2,
        "total_absent_days": 126,
        "total_employees": 50,
        "average_absent_days": 2.52
      }
    ],
    "department_comparison": [
      {
        "department_id": 1,
        "department_name": "Engineering",
        "absenteeism_rate": 4.5,
        "total_absent_days": 63
      }
    ],
    "peak_periods": [
      {
        "date": "2026-01-15",
        "absent_count": 12,
        "reason": "Holiday season"
      }
    ]
  }
}
```

### 6.3 Get Leave Liability Report
**Endpoint:** `GET {{BASE_URL}}/leave/reports/liability`

**Description:** Get leave liability report with financial impact.

**Query Parameters:**
- `as_of_date` (required): As of date (YYYY-MM-DD)
- `department_id` (optional): Filter by department
- `include_carry_forward` (optional): Include carry forward (default: true)
- `include_encashment` (optional): Include encashment projections (default: true)

**Sample Response:**
```json
{
  "success": true,
  "message": "Leave liability report retrieved successfully",
  "data": {
    "as_of_date": "2026-01-26",
    "total_liability_tzs": 12500000,
    "total_liability_usd": 5000,
    "department_wise": [
      {
        "department_id": 1,
        "department_name": "Engineering",
        "liability_tzs": 3500000,
        "liability_usd": 1400,
        "total_leave_days": 118,
        "average_daily_rate": 29661
      }
    ],
    "carry_forward_impact": {
      "total_carry_forward_days": 45,
      "carry_forward_liability_tzs": 1335000,
      "expiring_this_year": 20,
      "expiring_liability_tzs": 593000
    },
    "encashment_projections": {
      "eligible_encashment_days": 35,
      "projected_encashment_tzs": 1038000
    }
  }
}
```

### 6.4 Get Compliance Report
**Endpoint:** `GET {{BASE_URL}}/leave/reports/compliance`

**Description:** Get compliance report for statutory leave requirements.

**Query Parameters:**
- `start_date` (required): Start date (YYYY-MM-DD)
- `end_date` (required): End date (YYYY-MM-DD)
- `country` (optional): Filter by country

**Sample Response:**
```json
{
  "success": true,
  "message": "Compliance report retrieved successfully",
  "data": {
    "period": {
      "start_date": "2026-01-01",
      "end_date": "2026-01-31"
    },
    "country": "Tanzania",
    "statutory_requirements": [
      {
        "leave_type": "Annual Leave",
        "statutory_minimum": 28,
        "company_policy": 28,
        "compliance_status": "compliant",
        "employees_below_minimum": 0
      },
      {
        "leave_type": "Maternity Leave",
        "statutory_minimum": 84,
        "company_policy": 84,
        "compliance_status": "compliant",
        "employees_below_minimum": 0
      }
    ],
    "violations": [],
    "audit_ready": true
  }
}
```

### 6.5 Get Approval SLA Report
**Endpoint:** `GET {{BASE_URL}}/leave/reports/approval-sla`

**Description:** Get approval SLA performance report.

**Query Parameters:**
- `start_date` (required): Start date (YYYY-MM-DD)
- `end_date` (required): End date (YYYY-MM-DD)
- `approver_id` (optional): Filter by approver

**Sample Response:**
```json
{
  "success": true,
  "message": "Approval SLA report retrieved successfully",
  "data": {
    "period": {
      "start_date": "2026-01-01",
      "end_date": "2026-01-31"
    },
    "sla_target_hours": 24,
    "average_approval_time_hours": 18.5,
    "sla_compliance_rate": 92.5,
    "total_requests": 80,
    "approved_within_sla": 74,
    "breached_sla": 6,
    "approver_performance": [
      {
        "approver_id": "EMP010",
        "approver_name": "Manager Name",
        "total_approvals": 25,
        "average_time_hours": 16.2,
        "sla_compliance_rate": 96,
        "breaches": 1
      }
    ],
    "pending_aging": [
      {
        "request_id": 123,
        "pending_since_hours": 48,
        "status": "overdue"
      }
    ]
  }
}
```

### 6.6 Get Carry Forward Report
**Endpoint:** `GET {{BASE_URL}}/leave/reports/carry-forward`

**Description:** Get carry forward eligibility and expiry report.

**Query Parameters:**
- `year` (required): Year for carry forward analysis
- `leave_type_code` (optional): Filter by leave type

**Sample Response:**
```json
{
  "success": true,
  "message": "Carry forward report retrieved successfully",
  "data": {
    "year": 2026,
    "total_eligible_carry_forward": 45,
    "total_carry_forward_liability_tzs": 1335000,
    "expiring_this_year": 20,
    "expiring_liability_tzs": 593000,
    "eligible_for_encashment": 35,
    "encashment_liability_tzs": 1038000,
    "employee_wise": [
      {
        "employee_id": "EMP001",
        "employee_name": "John Doe",
        "leave_type": "AL",
        "carry_forward_days": 5,
        "expires_on": "2026-12-31",
        "eligible_for_encashment": true,
        "encashment_value_tzs": 148305
      }
    ]
  }
}
```

### 6.7 Export Report
**Endpoint:** `POST {{BASE_URL}}/leave/reports/export`

**Description:** Export leave report in various formats.

**Request Body:**
```json
{
  "report_type": "utilization",
  "format": "excel",
  "start_date": "2026-01-01",
  "end_date": "2026-01-31",
  "filters": {
    "department_id": 1,
    "leave_type_code": "AL"
  }
}
```

**Query Parameters:**
- `format` (required): Export format (excel, pdf, csv)

**Sample Response:**
```json
{
  "success": true,
  "message": "Report exported successfully",
  "data": {
    "file_url": "https://storage.example.com/reports/leave_utilization_202601.xlsx",
    "file_name": "leave_utilization_202601.xlsx",
    "file_size": 245760,
    "expires_at": "2026-01-27T10:30:00+03:00"
  }
}
```

### 6.8 Get Report Statistics
**Endpoint:** `GET {{BASE_URL}}/leave/reports/statistics`

**Description:** Get quick statistics for dashboard KPIs.

**Query Parameters:**
- `start_date` (optional): Start date (default: start of current month)
- `end_date` (optional): End date (default: today)

**Sample Response:**
```json
{
  "success": true,
  "message": "Report statistics retrieved successfully",
  "data": {
    "utilization_rate": 68.5,
    "utilization_trend": "+5.2%",
    "absenteeism_rate": 4.2,
    "absenteeism_trend": "-1.1%",
    "approved_leaves": 25,
    "total_leave_days": 98,
    "leave_liability_tzs": 12500000,
    "pending_approvals": 8
  }
}
```

---

## Common Response Format

All API responses follow this structure:

**Success Response:**
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Error message description",
  "error": "Error code or type",
  "errors": {
    "field_name": ["Error message for this field"]
  }
}
```

---

## Authentication

All endpoints require authentication via Bearer token in the Authorization header:

```
Authorization: Bearer {access_token}
```

---

## Pagination

Endpoints that support pagination return data in this format:

```json
{
  "success": true,
  "data": {
    "data": [...],
    "meta": {
      "page": 1,
      "per_page": 10,
      "total": 100,
      "total_pages": 10
    }
  }
}
```

---

## Manual Form Reference

### Leave Application Form (HR.FO.04.00)

The API endpoints are designed to support the manual leave application form used in the office. The form has two main sections:

#### Section 1: Application (To be filled in capital letters by applicant)

**Required Fields:**
- Employee names (first, middle, last) - separate fields, all in CAPITAL LETTERS
- Job position
- Date (application date)
- Department
- Contact No (employee contact number)
- Emergency contact person
- Emergency contact number
- Leave period Year
- Reporting back to work (Date)
- Category of Leave:
  - Annual Leave
  - Emergency Leave
  - Sick leave
  - Other leaves (specify)
- LEAVE TRACKER:
  - Start Date
  - End Date
  - Number of Days
  - Number of Previous days used
  - Number of Days Remaining after this request
  - Leave Taken from previous year
  - Leave Balance from previous year
- Employee signature (required for submission)

#### Section 2: Approval

**Multi-level Approval Fields:**
- Head of Department:
  - Name (in capital letters)
  - Signature
- HR Department:
  - Name (in capital letters)
  - Signature
- Director/CEO (If applicable):
  - Name (in capital letters)
  - Signature
- REMARKS: Approval/rejection remarks

**API Implementation Notes:**
1. All employee names should be stored and displayed in CAPITAL LETTERS as per form requirement
2. Signatures should be captured as base64 encoded images or signature data
3. Multi-level approval workflow supports:
   - Head of Department (Level 1)
   - HR Department (Level 2)
   - Director/CEO (Level 3, if applicable)
4. Each approval level requires:
   - Approver name (in capital letters)
   - Approver signature
   - Remarks (optional)
5. The `document_number` field should always be "HR.FO.04.00" for all leave applications
6. Leave tracker information is automatically calculated but can be manually adjusted if needed

---

## Notes

1. All dates should be in `YYYY-MM-DD` format
2. All timestamps should be in ISO 8601 format with timezone (e.g., `2026-01-26T10:30:00+03:00`)
3. Currency amounts are in TZS (Tanzanian Shilling) unless otherwise specified
4. Employee IDs follow the format `EMP###` (e.g., `EMP001`)
5. Leave type codes are typically 2-3 characters (e.g., `AL`, `SL`, `ML`)
6. Status values: `draft`, `pending`, `approved`, `rejected`, `cancelled`, `returned_for_info`, `partially_approved`
7. Employee names in form fields should be in CAPITAL LETTERS as per manual form requirement
8. Document number for all leave applications: `HR.FO.04.00`
9. Signatures should be provided as base64 encoded images (data URI format: `data:image/png;base64,...`)

---

## End of Documentation

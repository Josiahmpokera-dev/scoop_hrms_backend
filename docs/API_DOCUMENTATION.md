# HRMS Backend API Documentation

## Overview
This document provides comprehensive documentation for the HRMS Backend API. The API is organized around REST principles and uses standard HTTP response codes, authentication, and JSON data formats.

## Base URL
```
Development: http://localhost:8080
Production:  https://your-production-domain.com
```

## Authentication
Most endpoints require authentication via JWT Bearer tokens. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
```
**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "+1234567890"
}
```
**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

#### Login
```http
POST /api/v1/auth/login
```
**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```
**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "employee"
    }
  }
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
```
**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Get User Profile
```http
GET /api/v1/auth/profile
```

#### Logout
```http
POST /api/v1/auth/logout
```

## Employees

### Core Employee Management

#### Onboard New Employee
```http
POST /api/v1/employees
```
**Request Body:**
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane.smith@company.com",
  "phone_number": "+1234567890",
  "department_id": 1,
  "position_id": 2,
  "hire_date": "2026-01-15",
  "employment_type": "full_time"
}
```

#### Get Employee by ID
```http
GET /api/v1/employees/{id}
```

#### Get Employee by Employee ID
```http
GET /api/v1/employees/employee-id/{employee_id}
```

#### Update Employee
```http
PUT /api/v1/employees/{id}
```

#### List Employees
```http
GET /api/v1/employees?page=1&page_size=20&search=john&department_id=1&status=active
```

#### List Employees by Department
```http
GET /api/v1/employees/department/{department_id}
```

#### List Employees by Status
```http
GET /api/v1/employees/status/{status}
```

#### Delete Employee
```http
DELETE /api/v1/employees/{id}
```

### Self-Service

#### Get My Profile
```http
GET /api/v1/self-service/profile
```
**Response (Employee):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "employee_id": "EMP001",
    "personal": {
      "first_name": "John",
      "last_name": "Doe",
      "full_name": "John Doe",
      "email": "john@company.com"
    },
    "employment": {
      "department": "Engineering",
      "position": "Senior Developer",
      "hire_date": "2024-01-15"
    }
  }
}
```
**Response (Non-Employee User):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "user_id": 1,
    "role": "admin",
    "status": "active",
    "personal": {
      "first_name": "Admin",
      "last_name": "User",
      "full_name": "Admin User",
      "username": "admin",
      "email": "admin@company.com"
    },
    "account": {
      "email_verified": true,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### Onboarding

#### Create Onboarding Draft
```http
POST /api/v1/employees/onboarding/draft
```

#### Save Onboarding Step
```http
POST /api/v1/employees/onboarding/draft/{draft_id}/step/{step}
```

#### Get Onboarding Draft
```http
GET /api/v1/employees/onboarding/draft/{draft_id}
```

#### Complete Onboarding
```http
POST /api/v1/employees/onboarding/complete
```

#### List Non-Employee Users
```http
GET /api/v1/employees/onboarding/non-employee-users
```

#### List Onboarding Drafts
```http
GET /api/v1/employees/onboarding/drafts
```

## Attendance

### Timesheets

#### Create Timesheet Entry
```http
POST /api/v1/attendance/timesheets/entries
```
**Request Body:**
```json
{
  "date": "2026-04-27",
  "project_name": "Project Alpha",
  "task_description": "Implemented feature X",
  "hours": 8.0,
  "entry_type": "regular",
  "is_billable": true
}
```

#### Bulk Create Entries
```http
POST /api/v1/attendance/timesheets/entries/bulk
```

#### Update Entry
```http
PUT /api/v1/attendance/timesheets/entries/{entry_id}
```

#### Delete Entry
```http
DELETE /api/v1/attendance/timesheets/entries/{entry_id}
```

#### Get My Timesheets
```http
GET /api/v1/attendance/timesheets?status=submitted&page=1&page_size=20
```

#### Get Timesheet by ID
```http
GET /api/v1/attendance/timesheets/{timesheet_id}
```

#### Submit Timesheet
```http
POST /api/v1/attendance/timesheets/{timesheet_id}/submit
```

#### Recall Timesheet
```http
POST /api/v1/attendance/timesheets/{timesheet_id}/recall
```

#### Copy Last Week
```http
POST /api/v1/attendance/timesheets/copy-last-week
```

#### Get Employee Stats
```http
GET /api/v1/attendance/timesheets/stats?start_date=2026-04-01&end_date=2026-04-30
```

#### Get Team Timesheets (Manager)
```http
GET /api/v1/attendance/timesheets/team?week_start=2026-04-21&status=pending
```

### Timesheet Approvals (Admin/HR)

#### Get Pending Approvals
```http
GET /api/v1/attendance/timesheets/approvals
```

#### Approve/Reject Timesheet
```http
POST /api/v1/attendance/timesheets/approvals/{timesheet_id}
```
**Request Body:**
```json
{
  "status": "approved",
  "remarks": "Good work!"
}
```

### Overtime

#### Create Overtime Request
```http
POST /api/v1/attendance/overtime/requests
```
**Request Body:**
```json
{
  "date": "2026-04-27",
  "hours": 3.5,
  "reason": "Project deadline",
  "overtime_type": "weekday",
  "compensation_type": "payout"
}
```

#### Get My Overtime Requests
```http
GET /api/v1/attendance/overtime/requests?status=pending&page=1&page_size=20
```

#### Get Overtime Request by ID
```http
GET /api/v1/attendance/overtime/requests/{request_id}
```

#### Cancel Overtime Request
```http
POST /api/v1/attendance/overtime/requests/{request_id}/cancel
```

#### List Overtime Policies
```http
GET /api/v1/attendance/overtime/policies
```

#### Get Overtime Stats
```http
GET /api/v1/attendance/overtime/stats?start_date=2026-04-01&end_date=2026-04-30
```

### Overtime Admin (Admin/HR)

#### Create Policy
```http
POST /api/v1/attendance/overtime/admin/policies
```

#### Update Policy
```http
PUT /api/v1/attendance/overtime/admin/policies/{policy_id}
```

#### Delete Policy
```http
DELETE /api/v1/attendance/overtime/admin/policies/{policy_id}
```

#### List All Requests
```http
GET /api/v1/attendance/overtime/admin/requests?status=approved&start_date=2026-04-01&end_date=2026-04-30
```

#### Get Pending Approvals
```http
GET /api/v1/attendance/overtime/admin/pending
```

#### Approve/Reject Request
```http
POST /api/v1/attendance/overtime/admin/requests/{request_id}/approve
```
**Request Body:**
```json
{
  "status": "approved",
  "remarks": "Approved for payout"
}
```

### Manual Punches

#### Create Manual Punch
```http
POST /api/v1/attendance/manual-punches
```
**Request Body:**
```json
{
  "date": "2026-04-27",
  "time": "08:30:00",
  "punch_type": "in",
  "reason": "Biometric device not working",
  "notes": "Forgot to punch in"
}
```

#### Get My Manual Punches
```http
GET /api/v1/attendance/manual-punches/my?status=pending&start_date=2026-04-01&end_date=2026-04-30
```

#### Get Manual Punch by ID
```http
GET /api/v1/attendance/manual-punches/{id}
```

#### Cancel Manual Punch
```http
POST /api/v1/attendance/manual-punches/{id}/cancel
```

### Manual Punches Admin (Admin/HR)

#### List All Manual Punches
```http
GET /api/v1/attendance/manual-punches?status=pending&employee_id=123&page=1&page_size=20
```

#### Get Pending Manual Punches
```http
GET /api/v1/attendance/manual-punches/pending
```

#### Update Manual Punch Status
```http
PATCH /api/v1/attendance/manual-punches/{id}/status
```
**Request Body:**
```json
{
  "status": "approved",
  "remarks": "Verified with manager"
}
```

#### Delete Manual Punch
```http
DELETE /api/v1/attendance/manual-punches/{id}
```

### Attendance Reports (Admin/HR)

#### Get Attendance Summary
```http
GET /api/v1/attendance/reports/summary?startDate=2026-04-01&endDate=2026-04-30&departmentId=1
```

#### Get Attendance Trends
```http
GET /api/v1/attendance/reports/trends?startDate=2026-04-01&endDate=2026-04-30&departmentId=1&interval=daily
```

#### Get Department Stats
```http
GET /api/v1/attendance/reports/by-department?date=2026-04-27
```

#### Get Compliance Violations
```http
GET /api/v1/attendance/reports/compliance?startDate=2026-04-01&endDate=2026-04-30&severity=high
```

#### Get Overtime Analysis
```http
GET /api/v1/attendance/reports/overtime?startDate=2026-04-01&endDate=2026-04-30
```

#### Get Timesheet Summary
```http
GET /api/v1/attendance/reports/timesheets?start_date=2026-04-01&end_date=2026-04-30
```

#### Get Project Utilization
```http
GET /api/v1/attendance/reports/project-utilization?start_date=2026-04-01&end_date=2026-04-30
```

#### Get Employee Utilization
```http
GET /api/v1/attendance/reports/employee-utilization?start_date=2026-04-01&end_date=2026-04-30
```

#### Export Report
```http
POST /api/v1/attendance/reports/export
```
**Request Body:**
```json
{
  "reportType": "daily_attendance",
  "startDate": "2026-04-01",
  "endDate": "2026-04-30",
  "departmentId": "1",
  "format": "csv"
}
```

#### Get Attendance Overview
```http
GET /api/v1/attendance/overview?start_date=2026-04-01&end_date=2026-04-30&department_id=1&location_id=1
```

### Comprehensive Employee Report (Admin/HR)

#### Get Comprehensive Employee Report
```http
GET /api/v1/attendance/reports/employee-comprehensive?start_date=2026-04-01&end_date=2026-04-30&department_id=1&page=1&page_size=20
```
**Response:**
```json
{
  "success": true,
  "message": "Comprehensive employee report retrieved successfully",
  "data": [
    {
      "employee_id": "EMP001",
      "employee_name": "John Doe",
      "department": "Engineering",
      "position": "Senior Developer",
      "period": {
        "start_date": "2026-04-01",
        "end_date": "2026-04-30"
      },
      "attendance": {
        "total_working_days": 22,
        "present_days": 20,
        "absent_days": 2,
        "late_days": 1,
        "present_percentage": 90.9,
        "average_working_hours": 8.2,
        "total_working_hours": 164.0
      },
      "timesheet": {
        "total_hours": 168.0,
        "regular_hours": 160.0,
        "overtime_hours": 8.0,
        "billable_hours": 152.0,
        "non_billable_hours": 16.0
      },
      "leave": {
        "total_leave_days": 2,
        "leave_breakdown": [
          { "leave_type": "annual", "days": 2 }
        ],
        "pending_requests": 0,
        "approved_requests": 1,
        "rejected_requests": 0
      },
      "overtime": {
        "total_requests": 2,
        "approved_hours": 8.0,
        "pending_hours": 0,
        "rejected_hours": 0,
        "total_payout": 120.0,
        "comp_off_hours": 0
      },
      "compliance": {
        "violations_count": 1,
        "violations": [
          {
            "date": "2026-04-15",
            "type": "missing_checkout",
            "details": "Missing checkout punch for the day",
            "severity": "high"
          }
        ]
      }
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

## Leave Management

#### Apply for Leave
```http
POST /api/v1/leave/requests
```
**Request Body:**
```json
{
  "leave_type_id": 1,
  "from_date": "2026-05-01",
  "to_date": "2026-05-03",
  "reason": "Family vacation",
  "contact_during_leave": "john@example.com",
  "attachments": []
}
```

#### Get My Leave Requests
```http
GET /api/v1/leave/requests/my?page=1&page_size=20
```

#### Get Leave Request by ID
```http
GET /api/v1/leave/requests/{request_id}
```

#### Cancel Leave Request
```http
POST /api/v1/leave/requests/{request_id}/cancel
```

#### Get Leave Balance
```http
GET /api/v1/leave/balance
```

#### Get Leave Types
```http
GET /api/v1/leave/types
```

#### Get Holidays
```http
GET /api/v1/leave/holidays?year=2026
```

### Leave Admin (Admin/HR)

#### List All Leave Requests
```http
GET /api/v1/leave/admin/requests?status=pending&page=1&page_size=20
```

#### Approve/Reject Leave
```http
POST /api/v1/leave/requests/{request_id}/approve
```
**Request Body:**
```json
{
  "status": "approved",
  "remarks": "Approved, enjoy your vacation!"
}
```

#### Return for Information
```http
POST /api/v1/leave/requests/{request_id}/return-for-info
```

#### Create Leave Type
```http
POST /api/v1/leave/admin/types
```

#### Update Leave Type
```http
PUT /api/v1/leave/admin/types/{type_id}
```

#### Delete Leave Type
```http
DELETE /api/v1/leave/admin/types/{type_id}
```

#### Create Holiday
```http
POST /api/v1/leave/admin/holidays
```

#### Update Holiday
```http
PUT /api/v1/leave/admin/holidays/{holiday_id}
```

#### Delete Holiday
```http
DELETE /api/v1/leave/admin/holidays/{holiday_id}
```

## Payroll

#### Get Payroll Runs
```http
GET /api/v1/payroll/runs?page=1&page_size=20
```

#### Create Payroll Run
```http
POST /api/v1/payroll/runs
```

#### Get Payslip
```http
GET /api/v1/payroll/payslips/{payslip_id}
```

#### Get My Payslips
```http
GET /api/v1/payroll/payslips/my?page=1&page_size=20
```

## Biometric Integration

#### Test BioTime Connection
```http
GET /api/v1/biometric/biotime/test-connection
```

#### Get BioTime Token
```http
GET /api/v1/biometric/biotime/token
```

#### Get BioTime Terminals
```http
GET /api/v1/biometric/biotime/terminals
```

#### Get Device Status
```http
GET /api/v1/biometric/biotime/device-status
```

#### Get BioTime Transactions
```http
GET /api/v1/biometric/biotime/transactions?start_date=2026-04-01&end_date=2026-04-30&terminal_id=1
```

#### Get Transaction by ID
```http
GET /api/v1/biometric/biotime/transactions/{id}
```

#### Refresh Token
```http
POST /api/v1/biometric/biotime/refresh-token
```

#### Backfill Transactions
```http
POST /api/v1/biometric/biotime/backfill
```

#### Get Daily Attendance
```http
GET /api/v1/biometric/attendance/daily?date=2026-04-27&department_id=1
```

#### Get Exceptional Attendance
```http
GET /api/v1/biometric/attendance/exceptional?date=2026-04-27&type=late
```

### Biometric Enrollments

#### List Enrollments
```http
GET /api/v1/biometric/enrollments?page=1&page_size=20
```

#### Get Enrollment Statistics
```http
GET /api/v1/biometric/enrollments/statistics
```

#### Get Unlinked Employees
```http
GET /api/v1/biometric/enrollments/unlinked
```

#### Get Device Users
```http
GET /api/v1/biometric/enrollments/device-users
```

#### Link Employee
```http
POST /api/v1/biometric/enrollments/link
```
**Request Body:**
```json
{
  "employee_id": "EMP001",
  "device_user_id": "DEV001",
  "terminal_id": 1
}
```

#### Bulk Link
```http
POST /api/v1/biometric/enrollments/bulk-link
```

#### Auto Link
```http
POST /api/v1/biometric/enrollments/auto-link
```

#### Unlink Employee
```http
DELETE /api/v1/biometric/enrollments/{employeeId}/unlink
```

#### Get Merged Attendance
```http
GET /api/v1/biometric/attendance/merged?start_date=2026-04-01&end_date=2026-04-30&employee_id=123&page=1&page_size=20
```

#### Get Employee Attendance
```http
GET /api/v1/biometric/attendance/merged/{employeeId}?start_date=2026-04-01&end_date=2026-04-30&page=1&page_size=20
```

## Dashboard

#### Get Dashboard Statistics
```http
GET /api/v1/dashboard/statistics
```

#### Get Announcements
```http
GET /api/v1/dashboard/announcements?page=1&page_size=20
```

#### Mark Announcement as Read
```http
POST /api/v1/dashboard/announcements/{id}/read
```

#### Get Upcoming Events
```http
GET /api/v1/dashboard/events?start_date=2026-04-01&end_date=2026-04-30
```

#### Get My Recent Activity
```http
GET /api/v1/dashboard/my-activity?page=1&page_size=20
```

#### Get Quick Actions
```http
GET /api/v1/dashboard/quick-actions
```

#### Get Pending Approvals (Manager)
```http
GET /api/v1/dashboard/pending-approvals
```

#### Get Approval Counts
```http
GET /api/v1/dashboard/approval-counts
```

#### Get Employee Dashboard Statistics
```http
GET /api/v1/dashboard/employee/statistics
```

#### Get Employee Recent Activity
```http
GET /api/v1/dashboard/employee/my-activity?page=1&page_size=20
```

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": "Error code or details",
  "errors": {
    "field_name": "Field-specific error message"
  }
}
```

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Request succeeded |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid input or parameters |
| 401 | Unauthorized - Authentication required or failed |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 422 | Unprocessable Entity - Validation errors |
| 500 | Internal Server Error - Server error |

## Pagination

Most list endpoints support pagination:

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| page | int | 1 | - | Page number |
| page_size | int | 20 | 100 | Items per page |

## Filtering

Many endpoints support filtering:

| Parameter | Type | Description |
|-----------|------|-------------|
| search | string | Search by name, email, etc. |
| status | string | Filter by status |
| department_id | uint | Filter by department |
| start_date | string | Start date (YYYY-MM-DD) |
| end_date | string | End date (YYYY-MM-DD) |

## API Versions

| Version | Status | Base URL |
|---------|--------|----------|
| v1 | Current | /api/v1 |

## Rate Limiting

API requests are subject to rate limiting based on user role:
- Standard users: 1000 requests/hour
- Admin users: 5000 requests/hour

## Webhook Events

The system supports webhooks for real-time notifications:

| Event | Description |
|-------|-------------|
| employee.created | New employee onboarded |
| leave.approved | Leave request approved |
| timesheet.submitted | Timesheet submitted |
| overtime.approved | Overtime request approved |

## SDKs and Tools

### Postman Collection
A Postman collection is available at:
```
/docs/HRMS-API-Postman-Collection.json
```

### OpenAPI/Swagger
Interactive API documentation is available at:
```
Development: http://localhost:8080/swagger/index.html
Production:  https://your-production-domain.com/swagger/index.html
```

## Support

For API support and questions:
- Email: api-support@hrms.com
- Documentation: https://docs.hrms.com/api
- Status Page: https://status.hrms.com

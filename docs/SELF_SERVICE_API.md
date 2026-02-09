# Employee Self-Service API

Complete API documentation for all Employee Self-Service endpoints. These endpoints allow authenticated employees to manage their profiles, request leave, track timesheets, request assets, report issues, submit helpdesk tickets, and more — without needing HR/Admin privileges.

**Authentication:** All endpoints require a valid JWT token via `Authorization: Bearer <token>` header.

**Role:** All authenticated employees (no Admin/HR role required unless noted).

---

## Table of Contents

1. [Profile Management](#1-profile-management)
2. [Documents](#2-documents)
3. [Service Requests (HR Letters, IT, Facilities)](#3-service-requests)
4. [People Directory](#4-people-directory)
5. [Leave Management](#5-leave-management)
6. [Leave Calendar](#6-leave-calendar)
7. [Timesheets](#7-timesheets)
8. [Overtime](#8-overtime)
9. [Assets](#9-assets)
10. [Helpdesk & Tickets](#10-helpdesk--tickets)
11. [Knowledge Base](#11-knowledge-base)
12. [Payroll & Payslips](#12-payroll--payslips)
13. [Roster / Shift Management](#13-roster--shift-management)
14. [Dashboard & Notifications](#14-dashboard--notifications)
15. [Endpoints Summary](#15-endpoints-summary)

---

## 1. Profile Management

> **Base URL:** `/api/v1/self-service`

### GET `/self-service/profile`

Retrieve the authenticated employee's complete profile.

**Response includes:** personal info (name, DOB, gender, marital status, blood group, nationality, contacts), addresses (current, permanent), job info (employee ID, email, department, position, location, grade, employment type, shift, reporting manager), emergency contacts, documents, and profile update status.

```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "employee_id": "EMP-001",
    "personal": {
      "first_name": "John",
      "last_name": "Doe",
      "full_name": "John Doe",
      "date_of_birth": "1990-05-15",
      "gender": "male",
      "marital_status": "single",
      "blood_group": "O+",
      "nationality": "Tanzanian",
      "personal_email": "john@gmail.com",
      "personal_phone": "+255712345678",
      "photo": "/photos/john.jpg",
      "current_address": {
        "address_line1": "123 Main Street",
        "city": "Dar es Salaam",
        "state": "Dar es Salaam",
        "country": "Tanzania",
        "postal_code": "11000"
      },
      "permanent_address": null
    },
    "job": {
      "employee_id": "EMP-001",
      "official_email": "john.doe@company.com",
      "date_of_joining": "2023-01-15",
      "grade": "L5",
      "employment_type": "full_time",
      "work_phone": "+255222000111",
      "shift": "Day",
      "department": { "id": 3, "name": "Engineering", "code": "ENG" },
      "position": { "id": 5, "name": "Software Engineer", "code": "SE" },
      "work_location": { "id": 1, "name": "Dar es Salaam HQ", "address": "Plot 123, Bagamoyo Rd" },
      "reporting_manager": { "id": 10, "employee_id": "EMP-010", "full_name": "Jane Smith", "email": "jane@company.com" }
    },
    "emergency_contacts": [
      { "id": 1, "name": "Mary Doe", "relationship": "Spouse", "phone": "+255712000000", "is_primary": true }
    ],
    "documents": [
      { "id": 1, "document_type": "passport", "file_name": "passport.pdf", "file_url": "/docs/passport.pdf", "uploaded_at": "2023-01-20T10:00:00Z" }
    ],
    "profile_update_status": {
      "has_pending_updates": false,
      "last_update_request_date": null,
      "last_update_status": null
    }
  }
}
```

---

### POST `/self-service/profile/update`

Submit a profile update request (requires manager approval).

**Request Body:**

| Field     | Type   | Required | Description                                          |
|-----------|--------|----------|------------------------------------------------------|
| `section` | string | Yes      | `personal`, `emergency_contacts`, or `address`        |
| `updates` | object | Yes      | Key-value pairs of fields to update                  |
| `reason`  | string | No       | Reason for the change                                |

```json
{
  "section": "personal",
  "updates": {
    "phone_number": "+255712999999",
    "marital_status": "married"
  },
  "reason": "Recently married, updating marital status and phone"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Profile update request submitted successfully",
  "data": {
    "update_request_id": "PUR-2026-001",
    "section": "personal",
    "status": "pending",
    "submitted_at": "2026-02-06T10:00:00Z",
    "approver": { "id": 10, "name": "Jane Smith", "email": "jane@company.com" },
    "estimated_processing_time": "24-48 hours"
  }
}
```

---

### GET `/self-service/profile/update-status`

Check status of all profile update requests.

```json
{
  "success": true,
  "data": {
    "pending_requests": [
      { "update_request_id": "PUR-2026-001", "section": "personal", "status": "pending", "submitted_at": "2026-02-06T10:00:00Z", "approver": { "id": 10, "name": "Jane Smith" } }
    ],
    "recent_updates": [
      { "update_request_id": "PUR-2025-042", "section": "address", "status": "approved", "submitted_at": "2025-12-01T09:00:00Z", "approved_at": "2025-12-02T14:00:00Z" }
    ]
  }
}
```

---

### GET `/self-service/profile/id-card`

Download employee ID card. *(Placeholder — not yet implemented)*

---

## 2. Documents

### GET `/self-service/profile/documents`

Retrieve employee documents.

**Query Parameters:**

| Parameter       | Type   | Description              |
|-----------------|--------|--------------------------|
| `document_type` | string | Filter by document type  |
| `page`          | int    | Page number (default: 1) |
| `page_size`     | int    | Items per page (max 100) |

---

## 3. Service Requests

> HR Letters, IT Requests, Facilities Requests

### POST `/self-service/requests`

Create a new service request.

| Field         | Type   | Required | Description                                             |
|---------------|--------|----------|---------------------------------------------------------|
| `type`        | string | Yes      | `hr_letter`, `it_request`, or `facilities`              |
| `category`    | string | Yes      | Category (e.g., "employment_letter", "hardware", etc.)  |
| `subject`     | string | Yes      | Brief subject                                           |
| `description` | string | No       | Detailed description                                    |
| `priority`    | string | No       | `low`, `medium` (default), `high`, `urgent`             |

**For HR Letters additionally:**

| Field              | Type   | Description                    |
|--------------------|--------|--------------------------------|
| `letter_type`      | string | Type of letter                 |
| `purpose`          | string | Purpose of the letter          |
| `addressed_to`     | string | Who it's addressed to          |
| `additional_notes` | string | Extra notes                    |

**For IT/Facilities additionally:**

| Field             | Type    | Description              |
|-------------------|---------|--------------------------|
| `requested_items` | array   | List of items requested  |

**Example — HR Letter:**

```json
{
  "type": "hr_letter",
  "category": "employment_letter",
  "subject": "Employment Confirmation Letter",
  "letter_type": "employment_confirmation",
  "purpose": "Bank loan application",
  "addressed_to": "NMB Bank, Dar es Salaam Branch"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Service request created successfully",
  "data": {
    "id": 5,
    "request_number": "SR-2026-005",
    "type": "hr_letter",
    "category": "employment_letter",
    "subject": "Employment Confirmation Letter",
    "status": "submitted",
    "priority": "medium",
    "requested_date": "2026-02-06T10:00:00Z",
    "sla_hours": 48,
    "assigned_to": { "name": "HR Support Team" },
    "estimated_completion_date": "2026-02-08T10:00:00Z"
  }
}
```

---

### GET `/self-service/requests`

List service requests.

| Parameter  | Type   | Description                                |
|------------|--------|--------------------------------------------|
| `type`     | string | Filter: `hr_letter`, `it_request`, `facilities` |
| `status`   | string | Filter: `submitted`, `in_progress`, `completed`, `cancelled` |
| `priority` | string | Filter: `low`, `medium`, `high`, `urgent`  |
| `page`     | int    | Page number                                |
| `page_size`| int    | Items per page                             |

---

### GET `/self-service/requests/:request_id`

Get service request details.

---

### POST `/self-service/requests/:request_id/cancel`

Cancel a service request.

```json
{ "reason": "No longer needed" }
```

---

## 4. People Directory

### GET `/self-service/directory`

Search and browse the employee directory.

| Parameter       | Type   | Description                               |
|-----------------|--------|-------------------------------------------|
| `search`        | string | Search by name, employee ID, email        |
| `department_id` | int    | Filter by department                      |
| `position_id`   | int    | Filter by position                        |
| `location_id`   | int    | Filter by location                        |
| `status`        | string | Filter by status (default: `active`)      |
| `page`          | int    | Page number                               |
| `page_size`     | int    | Items per page                            |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "employee_id": "EMP-001",
      "full_name": "John Doe",
      "first_name": "John",
      "last_name": "Doe",
      "status": "active",
      "photo": "/photos/john.jpg",
      "official_email": "john@company.com",
      "work_phone": "+255222000111"
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 45, "total_pages": 3 }
}
```

> **Note:** There is also a dedicated People Directory API at `/api/v1/people-directory` with more filters (team, manager, employment type, letter, sort). See `PEOPLE_DIRECTORY_API.md`.

---

### GET `/self-service/directory/:employee_id`

Get detailed directory profile for a specific employee (public contact/work info only).

---

## 5. Leave Management

> **Base URL:** `/api/v1/leave`

### GET `/leave/employee-info`

Get employee info relevant to leave management.

---

### GET `/leave/balances`

Get the employee's leave balances (annual, sick, personal, etc.).

```json
{
  "success": true,
  "data": [
    { "leave_type": "annual", "total": 21, "used": 5, "pending": 2, "available": 14 },
    { "leave_type": "sick", "total": 10, "used": 1, "pending": 0, "available": 9 },
    { "leave_type": "personal", "total": 5, "used": 0, "pending": 0, "available": 5 }
  ]
}
```

---

### GET `/leave/policies/guidelines`

Get leave policy guidelines (rules, entitlements, restrictions).

---

### POST `/leave/calculate-days`

Calculate the number of leave days for a given date range (excluding weekends/holidays).

```json
{
  "start_date": "2026-03-10",
  "end_date": "2026-03-14",
  "leave_type": "annual"
}
```

---

### POST `/leave/applications`

Submit a new leave request.

| Field        | Type   | Required | Description                     |
|--------------|--------|----------|---------------------------------|
| `leave_type` | string | Yes      | Type of leave                   |
| `start_date` | string | Yes      | Start date (YYYY-MM-DD)         |
| `end_date`   | string | Yes      | End date (YYYY-MM-DD)           |
| `reason`     | string | Yes      | Reason for leave                |
| `half_day`   | bool   | No       | Is it a half-day leave?         |
| `attachments`| array  | No       | Supporting documents (e.g. medical certificate) |

```json
{
  "leave_type": "annual",
  "start_date": "2026-03-10",
  "end_date": "2026-03-14",
  "reason": "Family vacation"
}
```

---

### GET `/leave/requests`

List the employee's leave requests with filters.

| Parameter    | Type   | Description                                         |
|--------------|--------|-----------------------------------------------------|
| `status`     | string | Filter: `pending`, `approved`, `rejected`, `cancelled` |
| `leave_type` | string | Filter by leave type                                |
| `year`       | int    | Filter by year                                      |
| `page`       | int    | Page number                                         |
| `page_size`  | int    | Items per page                                      |

---

### GET `/leave/requests/:request_id`

Get details of a specific leave request.

---

### PUT `/leave/requests/:request_id`

Update a pending leave request (change dates, reason, etc.).

---

### POST `/leave/requests/:request_id/cancel`

Cancel a leave request.

---

### DELETE `/leave/requests/:request_id`

Delete a leave request (only if still pending/draft).

---

## 6. Leave Calendar

### GET `/leave/calendar`

Get the leave calendar showing who is on leave and when.

| Parameter    | Type   | Description                  |
|--------------|--------|------------------------------|
| `month`      | int    | Month (1-12)                 |
| `year`       | int    | Year                         |
| `department` | int    | Filter by department ID      |

---

### GET `/leave/calendar/today`

Get the list of employees on leave today.

---

### GET `/leave/calendar/week`

Get the weekly leave calendar.

---

## 7. Timesheets

> **Base URL:** `/api/v1/attendance/timesheets`

### POST `/attendance/timesheets/entries`

Create a single timesheet entry.

| Field          | Type   | Required | Description                                                  |
|----------------|--------|----------|--------------------------------------------------------------|
| `date`         | string | Yes      | Date (YYYY-MM-DD)                                           |
| `project_id`   | uint   | No       | Project ID                                                   |
| `project_name` | string | No       | Project name (if no project_id)                              |
| `task`         | string | No       | Task description                                             |
| `entry_type`   | string | No       | `regular`, `meeting`, `training`, `travel`, `support`, `overtime`, `break`, `other` |
| `hours`        | number | No       | Hours worked (0.5h precision, 0–24, optional)                |
| `is_billable`  | bool   | No       | Billable flag (default: false)                               |
| `notes`        | string | No       | Notes                                                        |

```json
{
  "date": "2026-02-06",
  "project_name": "HRMS Backend",
  "task": "API Development",
  "entry_type": "regular",
  "hours": 8,
  "is_billable": true,
  "notes": "Completed assets module"
}
```

---

### POST `/attendance/timesheets/entries/bulk`

Bulk create entries for a week.

```json
{
  "week_start_date": "2026-02-02",
  "entries": [
    { "date": "2026-02-02", "project_name": "HRMS", "task": "API Dev", "hours": 8 },
    { "date": "2026-02-03", "project_name": "HRMS", "task": "Testing", "hours": 7.5 }
  ]
}
```

---

### PUT `/attendance/timesheets/entries/:entry_id`

Update an existing entry (only if timesheet is in `draft` status).

---

### DELETE `/attendance/timesheets/entries/:entry_id`

Delete a timesheet entry.

---

### GET `/attendance/timesheets`

List the employee's weekly timesheets.

| Parameter | Type   | Description                               |
|-----------|--------|-------------------------------------------|
| `status`  | string | Filter: `draft`, `submitted`, `approved`, `rejected` |
| `page`    | int    | Page number                               |
| `page_size`| int   | Items per page                            |

---

### GET `/attendance/timesheets/:timesheet_id`

Get a specific timesheet with all entries.

---

### POST `/attendance/timesheets/:timesheet_id/submit`

Submit a draft timesheet for approval.

---

### POST `/attendance/timesheets/:timesheet_id/recall`

Recall a submitted timesheet (back to draft).

---

### POST `/attendance/timesheets/copy-last-week`

Copy entries from the previous week into a new week.

```json
{ "target_week_start_date": "2026-02-09" }
```

---

### GET `/attendance/timesheets/stats`

Get personal utilization statistics (total hours, billable hours, etc.).

---

### GET `/attendance/timesheets/team`

Get team timesheets (for managers viewing their team's submissions).

---

## 8. Overtime

> **Base URL:** `/api/v1/attendance/overtime`

### POST `/attendance/overtime/requests`

Submit an overtime request.

| Field             | Type   | Required | Description                               |
|-------------------|--------|----------|-------------------------------------------|
| `date`            | string | Yes      | OT date (YYYY-MM-DD)                      |
| `start_time`      | string | Yes      | Start time (HH:MM)                        |
| `end_time`        | string | Yes      | End time (HH:MM)                          |
| `reason`          | string | Yes      | Reason for overtime                        |
| `compensation_type`| string| No       | `paid`, `comp_off` (default: based on policy) |

---

### GET `/attendance/overtime/requests`

List the employee's overtime requests.

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `status`  | string | Filter: `pending`, `approved`, `rejected`, `cancelled` |
| `page`    | int    | Page number |
| `page_size`| int   | Items per page |

---

### GET `/attendance/overtime/requests/:request_id`

Get specific OT request details.

---

### POST `/attendance/overtime/requests/:request_id/cancel`

Cancel a pending OT request.

---

### GET `/attendance/overtime/policies`

View available overtime policies (rates, limits, thresholds).

---

### GET `/attendance/overtime/policies/:policy_id`

View a specific OT policy.

---

### GET `/attendance/overtime/stats`

Get personal overtime statistics (hours this month, compensation, etc.).

---

## 9. Assets

> **Base URL:** `/api/v1/self-service/assets`
> See also: `ASSETS_API.md` for full details.

### GET `/self-service/assets`

Get all assets currently assigned to the employee.

---

### POST `/self-service/assets/requests`

Request a new, replacement, or additional asset.

| Field                | Type    | Required | Description                              |
|----------------------|---------|----------|------------------------------------------|
| `request_type`       | string  | Yes      | `new`, `replacement`, `additional`        |
| `asset_type`         | string  | Yes      | `laptop`, `mobile_phone`, `desktop`, etc. |
| `justification`      | string  | Yes      | Reason for the request                   |
| `brand`              | string  | No       | Preferred brand                          |
| `model`              | string  | No       | Preferred model                          |
| `priority`           | string  | No       | `low`, `medium`, `high`, `urgent`        |
| `replacing_asset_id` | integer | No       | Asset being replaced (for replacement)   |
| `notes`              | string  | No       | Additional notes                         |

---

### GET `/self-service/assets/requests`

List the employee's asset requests.

| Parameter      | Type   | Description |
|----------------|--------|-------------|
| `request_type` | string | Filter: `new`, `replacement`, `additional` |
| `status`       | string | Filter: `pending`, `approved`, `rejected`, `fulfilled`, `cancelled` |
| `asset_type`   | string | Filter by asset type |
| `page` / `page_size` | int | Pagination |

---

### GET `/self-service/assets/requests/:request_id`

Get asset request details.

---

### POST `/self-service/assets/requests/:request_id/cancel`

Cancel a pending asset request.

---

### POST `/self-service/assets/issues`

Report an issue with an assigned asset.

| Field                  | Type     | Required | Description                          |
|------------------------|----------|----------|--------------------------------------|
| `asset_id`             | integer  | Yes      | Asset with the issue                 |
| `issue_type`           | string   | Yes      | `malfunction`, `damage`, `lost`, `stolen`, `other` |
| `title`                | string   | Yes      | Brief title                          |
| `description`          | string   | Yes      | Detailed description                 |
| `priority`             | string   | No       | `low`, `medium`, `high`, `urgent`    |
| `incident_date`        | string   | No       | When it happened (YYYY-MM-DD)        |
| `incident_location`    | string   | No       | Where it happened                    |
| `police_report_number` | string   | No       | For stolen assets                    |
| `attachment_urls`      | string[] | No       | Photos/documents                     |

---

### GET `/self-service/assets/issues`

List the employee's reported asset issues.

| Parameter    | Type   | Description |
|--------------|--------|-------------|
| `issue_type` | string | Filter by type |
| `status`     | string | Filter: `reported`, `under_review`, `in_progress`, `resolved`, `closed` |
| `asset_id`   | int    | Filter by asset |
| `page` / `page_size` | int | Pagination |

---

### GET `/self-service/assets/issues/:issue_id`

Get asset issue details (including resolution info and attachments).

---

## 10. Helpdesk & Tickets

> **Base URL:** `/api/v1/helpdesk`

### GET `/helpdesk/ticket-categories`

Get available ticket categories for submission.

---

### POST `/helpdesk/tickets`

Create a new helpdesk ticket.

| Field         | Type   | Required | Description                        |
|---------------|--------|----------|------------------------------------|
| `category_id` | int    | Yes      | Ticket category ID                 |
| `subject`     | string | Yes      | Ticket subject                     |
| `description` | string | Yes      | Detailed description               |
| `priority`    | string | No       | `low`, `medium`, `high`, `urgent`  |

---

### GET `/helpdesk/tickets`

List the employee's helpdesk tickets.

| Parameter  | Type   | Description |
|------------|--------|-------------|
| `status`   | string | Filter by status |
| `priority` | string | Filter by priority |
| `page` / `page_size` | int | Pagination |

---

### GET `/helpdesk/tickets/recent`

Get recent tickets (quick view).

---

### GET `/helpdesk/tickets/:ticket_id`

Get full ticket details (with comments and attachments).

---

### POST `/helpdesk/tickets/:ticket_id/comments`

Add a comment/reply to a ticket.

```json
{ "content": "Any update on this issue?" }
```

---

### POST `/helpdesk/tickets/:ticket_id/attachments`

Upload attachments to a ticket.

---

### POST `/helpdesk/tickets/:ticket_id/close`

Close a ticket.

```json
{ "resolution_notes": "Issue resolved" }
```

---

### POST `/helpdesk/tickets/:ticket_id/csat`

Submit customer satisfaction rating after ticket closure.

```json
{ "rating": 5, "feedback": "Very quick resolution!" }
```

---

## 11. Knowledge Base

> **Base URL:** `/api/v1/helpdesk/knowledge-base`

### GET `/helpdesk/knowledge-base/categories`

List knowledge base categories.

---

### GET `/helpdesk/knowledge-base/articles`

Browse/search knowledge base articles.

| Parameter     | Type   | Description                |
|---------------|--------|----------------------------|
| `category_id` | int    | Filter by category         |
| `search`      | string | Search articles            |
| `page` / `page_size` | int | Pagination          |

---

### GET `/helpdesk/knowledge-base/articles/:article_id`

Read a specific article.

---

### POST `/helpdesk/knowledge-base/articles/:article_id/feedback`

Submit feedback on an article (helpful / not helpful).

```json
{ "helpful": true, "comment": "This solved my problem!" }
```

---

## 12. Payroll & Payslips

> **Base URL:** `/api/v1/self-service/payroll`

### GET `/self-service/payroll/payslips`

Get the employee's payslips list.

| Parameter | Type | Description |
|-----------|------|-------------|
| `year`    | int  | Filter by year |
| `month`   | int  | Filter by month |
| `page` / `page_size` | int | Pagination |

---

### GET `/self-service/payroll/payslips/latest`

Get the most recent payslip.

---

### GET `/self-service/payroll/payslips/summary`

Get salary slip summary (YTD breakdown).

---

### GET `/self-service/payroll/payslips/:id/download`

Download a payslip as PDF.

---

### GET `/self-service/payroll/loans`

Get the employee's loan records.

---

### POST `/self-service/payroll/loans/apply`

Apply for a loan.

| Field       | Type   | Required | Description         |
|-------------|--------|----------|---------------------|
| `loan_type` | string | Yes      | Type of loan        |
| `amount`    | number | Yes      | Loan amount         |
| `tenure`    | int    | Yes      | Repayment period (months) |
| `reason`    | string | Yes      | Purpose of the loan |

> **Note:** The legacy payslip endpoints at `/self-service/payslips` are placeholders. Use `/self-service/payroll/payslips` for the active implementation.

---

## 13. Roster / Shift Management

> **Base URL:** `/api/v1/self-service/rosters`

### GET `/self-service/rosters/assignments`

Get the employee's roster/shift assignments.

| Parameter    | Type   | Description          |
|--------------|--------|----------------------|
| `start_date` | string | Filter start date    |
| `end_date`   | string | Filter end date      |
| `status`     | string | Filter by status     |

---

### GET `/self-service/rosters/assignments/:assignment_id`

Get a specific roster assignment details.

---

### POST `/self-service/rosters/change-requests`

Request a roster/shift change.

| Field               | Type   | Required | Description                     |
|---------------------|--------|----------|---------------------------------|
| `assignment_id`     | int    | Yes      | Roster assignment to change     |
| `requested_shift`   | string | No       | Preferred shift                 |
| `requested_date`    | string | No       | Preferred date                  |
| `reason`            | string | Yes      | Reason for the change           |

---

### GET `/self-service/rosters/change-requests`

List the employee's roster change requests.

---

### GET `/self-service/rosters/change-requests/:request_id`

Get roster change request details.

---

## 14. Dashboard & Notifications

These endpoints provide personalized data for the employee's dashboard.

### GET `/api/v1/dashboard/statistics`

Employee dashboard KPIs: leave balance, hours this week, pending requests, next payday.

---

### GET `/api/v1/dashboard/my-activity`

Recent activity feed (leave requests, approvals, etc.).

---

### GET `/api/v1/dashboard/announcements`

Company announcements with read/unread tracking.

---

### POST `/api/v1/dashboard/announcements/:id/read`

Mark an announcement as read.

---

### GET `/api/v1/dashboard/quick-actions`

Get personalized quick action links for the dashboard.

---

### GET `/api/v1/notifications/count`

Get notification badge count (unread announcements, pending items).

```json
{
  "success": true,
  "data": {
    "total": 5,
    "unread_announcements": 2,
    "my_pending": {
      "total": 3,
      "leave_requests": 1,
      "draft_timesheets": 1,
      "overtime_requests": 1
    }
  }
}
```

---

### GET `/api/v1/notifications`

Get notification items list.

---

### GET `/api/v1/features/my-access`

Get the employee's feature access permissions.

---

### GET `/api/v1/auth/profile`

Get the authenticated user's account profile (user-level, not employee-level).

---

## 15. Endpoints Summary

### Profile & Personal (`/api/v1/self-service`)

| Method | Endpoint                                    | Description                  |
|--------|---------------------------------------------|------------------------------|
| GET    | `/self-service/profile`                     | Get my profile               |
| POST   | `/self-service/profile/update`              | Request profile update       |
| GET    | `/self-service/profile/update-status`       | Profile update status        |
| GET    | `/self-service/profile/documents`           | Get my documents             |
| GET    | `/self-service/profile/id-card`             | Download ID card             |

### Service Requests (`/api/v1/self-service/requests`)

| Method | Endpoint                                         | Description           |
|--------|--------------------------------------------------|-----------------------|
| POST   | `/self-service/requests`                         | Create request        |
| GET    | `/self-service/requests`                         | List my requests      |
| GET    | `/self-service/requests/:id`                     | Get request details   |
| POST   | `/self-service/requests/:id/cancel`              | Cancel request        |

### People Directory (`/api/v1/self-service/directory`)

| Method | Endpoint                                    | Description              |
|--------|---------------------------------------------|--------------------------|
| GET    | `/self-service/directory`                   | Search directory         |
| GET    | `/self-service/directory/:employee_id`      | Get employee details     |

### Leave (`/api/v1/leave`)

| Method | Endpoint                                   | Description              |
|--------|--------------------------------------------|--------------------------|
| GET    | `/leave/employee-info`                     | Get employee info        |
| GET    | `/leave/balances`                          | Get leave balances       |
| GET    | `/leave/policies/guidelines`               | Get leave policies       |
| POST   | `/leave/calculate-days`                    | Calculate leave days     |
| POST   | `/leave/applications`                      | Submit leave request     |
| GET    | `/leave/requests`                          | List leave requests      |
| GET    | `/leave/requests/:id`                      | Get leave request        |
| PUT    | `/leave/requests/:id`                      | Update leave request     |
| POST   | `/leave/requests/:id/cancel`               | Cancel leave request     |
| DELETE | `/leave/requests/:id`                      | Delete leave request     |
| GET    | `/leave/calendar`                          | Leave calendar           |
| GET    | `/leave/calendar/today`                    | Who's on leave today     |
| GET    | `/leave/calendar/week`                     | Weekly calendar          |

### Timesheets (`/api/v1/attendance/timesheets`)

| Method | Endpoint                                              | Description            |
|--------|-------------------------------------------------------|------------------------|
| POST   | `/attendance/timesheets/entries`                      | Create entry           |
| POST   | `/attendance/timesheets/entries/bulk`                  | Bulk create entries    |
| PUT    | `/attendance/timesheets/entries/:id`                   | Update entry           |
| DELETE | `/attendance/timesheets/entries/:id`                   | Delete entry           |
| GET    | `/attendance/timesheets`                              | List my timesheets     |
| GET    | `/attendance/timesheets/:id`                          | Get timesheet          |
| POST   | `/attendance/timesheets/:id/submit`                   | Submit for approval    |
| POST   | `/attendance/timesheets/:id/recall`                   | Recall submission      |
| POST   | `/attendance/timesheets/copy-last-week`               | Copy last week         |
| GET    | `/attendance/timesheets/stats`                        | My utilization stats   |
| GET    | `/attendance/timesheets/team`                         | Team timesheets        |

### Overtime (`/api/v1/attendance/overtime`)

| Method | Endpoint                                              | Description            |
|--------|-------------------------------------------------------|------------------------|
| POST   | `/attendance/overtime/requests`                       | Submit OT request      |
| GET    | `/attendance/overtime/requests`                       | List my OT requests    |
| GET    | `/attendance/overtime/requests/:id`                   | Get OT request         |
| POST   | `/attendance/overtime/requests/:id/cancel`            | Cancel OT request      |
| GET    | `/attendance/overtime/policies`                       | View OT policies       |
| GET    | `/attendance/overtime/policies/:id`                   | View specific policy   |
| GET    | `/attendance/overtime/stats`                          | My OT statistics       |

### Assets (`/api/v1/self-service/assets`)

| Method | Endpoint                                              | Description            |
|--------|-------------------------------------------------------|------------------------|
| GET    | `/self-service/assets`                                | My assigned assets     |
| POST   | `/self-service/assets/requests`                       | Request an asset       |
| GET    | `/self-service/assets/requests`                       | List my requests       |
| GET    | `/self-service/assets/requests/:id`                   | Get request details    |
| POST   | `/self-service/assets/requests/:id/cancel`            | Cancel request         |
| POST   | `/self-service/assets/issues`                         | Report an issue        |
| GET    | `/self-service/assets/issues`                         | List my issues         |
| GET    | `/self-service/assets/issues/:id`                     | Get issue details      |

### Helpdesk (`/api/v1/helpdesk`)

| Method | Endpoint                                              | Description            |
|--------|-------------------------------------------------------|------------------------|
| GET    | `/helpdesk/ticket-categories`                         | Get categories         |
| POST   | `/helpdesk/tickets`                                   | Create ticket          |
| GET    | `/helpdesk/tickets`                                   | List my tickets        |
| GET    | `/helpdesk/tickets/recent`                            | Recent tickets         |
| GET    | `/helpdesk/tickets/:id`                               | Get ticket details     |
| POST   | `/helpdesk/tickets/:id/comments`                      | Add comment            |
| POST   | `/helpdesk/tickets/:id/attachments`                   | Upload attachments     |
| POST   | `/helpdesk/tickets/:id/close`                         | Close ticket           |
| POST   | `/helpdesk/tickets/:id/csat`                          | Submit satisfaction    |
| GET    | `/helpdesk/knowledge-base/categories`                 | KB categories          |
| GET    | `/helpdesk/knowledge-base/articles`                   | Browse articles        |
| GET    | `/helpdesk/knowledge-base/articles/:id`               | Read article           |
| POST   | `/helpdesk/knowledge-base/articles/:id/feedback`      | Article feedback       |

### Payroll (`/api/v1/self-service/payroll`)

| Method | Endpoint                                              | Description            |
|--------|-------------------------------------------------------|------------------------|
| GET    | `/self-service/payroll/payslips`                      | My payslips            |
| GET    | `/self-service/payroll/payslips/latest`               | Latest payslip         |
| GET    | `/self-service/payroll/payslips/summary`              | Salary summary         |
| GET    | `/self-service/payroll/payslips/:id/download`         | Download payslip       |
| GET    | `/self-service/payroll/loans`                         | My loans               |
| POST   | `/self-service/payroll/loans/apply`                   | Apply for loan         |

### Rosters (`/api/v1/self-service/rosters`)

| Method | Endpoint                                              | Description              |
|--------|-------------------------------------------------------|--------------------------|
| GET    | `/self-service/rosters/assignments`                   | My roster assignments    |
| GET    | `/self-service/rosters/assignments/:id`               | Assignment details       |
| POST   | `/self-service/rosters/change-requests`               | Request shift change     |
| GET    | `/self-service/rosters/change-requests`               | List my change requests  |
| GET    | `/self-service/rosters/change-requests/:id`           | Change request details   |

### Dashboard & Notifications

| Method | Endpoint                                  | Description                |
|--------|-------------------------------------------|----------------------------|
| GET    | `/dashboard/statistics`                   | My dashboard stats         |
| GET    | `/dashboard/my-activity`                  | My activity feed           |
| GET    | `/dashboard/announcements`                | Announcements              |
| POST   | `/dashboard/announcements/:id/read`       | Mark as read               |
| GET    | `/dashboard/quick-actions`                | Quick actions              |
| GET    | `/notifications/count`                    | Notification badge count   |
| GET    | `/notifications`                          | Notification list          |
| GET    | `/features/my-access`                     | My feature access          |
| GET    | `/auth/profile`                           | My user profile            |

---

**Total Self-Service Endpoints: 75+**

All endpoints require authentication but do **not** require HR/Admin roles, making them accessible to every employee in the organization.

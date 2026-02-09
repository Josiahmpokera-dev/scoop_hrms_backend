# People Directory API

The People Directory API provides a comprehensive employee directory for the organization. It allows searching, filtering, and browsing employees with rich profile data, department/location breakdowns, and filter options for building a user-friendly directory UI.

**Base URL:** `/api/v1/people-directory`

**Authentication:** All endpoints require a valid JWT token via `Authorization: Bearer <token>` header.

---

## Table of Contents

1. [Search / Browse Directory](#1-search--browse-directory)
2. [Get Employee Profile](#2-get-employee-profile)
3. [Get Filter Options](#3-get-filter-options)
4. [Get Directory Statistics](#4-get-directory-statistics)
5. [Data Models](#5-data-models)

---

## 1. Search / Browse Directory

Search and filter employees in the organization directory.

**Endpoint:** `GET /api/v1/people-directory`

### Query Parameters

| Parameter         | Type    | Required | Default      | Description                                                              |
|-------------------|---------|----------|--------------|--------------------------------------------------------------------------|
| `search`          | string  | No       | —            | Search by name, employee ID, email, or phone number                      |
| `department_id`   | integer | No       | —            | Filter by department ID                                                  |
| `position_id`     | integer | No       | —            | Filter by position/designation ID                                        |
| `location_id`     | integer | No       | —            | Filter by work location ID                                               |
| `team_id`         | integer | No       | —            | Filter by team ID                                                        |
| `manager_id`      | integer | No       | —            | Filter by reporting manager ID (shows direct reports)                    |
| `employment_type` | string  | No       | —            | Filter by employment type: `full_time`, `part_time`, `contract`, `intern`|
| `status`          | string  | No       | `active`     | Filter by status: `active`, `inactive`, `on_leave`, `suspended`, `terminated` |
| `letter`          | string  | No       | —            | Filter by first letter of first name (A–Z)                              |
| `sort_by`         | string  | No       | `first_name` | Sort field: `first_name`, `last_name`, `employee_id`, `department`, `hire_date`, `created_at` |
| `sort_order`      | string  | No       | `asc`        | Sort direction: `asc` or `desc`                                         |
| `page`            | integer | No       | `1`          | Page number (1-indexed)                                                  |
| `page_size`       | integer | No       | `20`         | Items per page (max 100)                                                 |

### Example Requests

```bash
# Basic search
GET /api/v1/people-directory?search=john

# Filter by department
GET /api/v1/people-directory?department_id=3

# Browse by letter
GET /api/v1/people-directory?letter=A&page=1&page_size=20

# Combined filters
GET /api/v1/people-directory?department_id=3&employment_type=full_time&sort_by=last_name&sort_order=asc

# Direct reports of a specific manager
GET /api/v1/people-directory?manager_id=15
```

### Success Response (200 OK)

```json
{
  "success": true,
  "message": "People directory retrieved successfully",
  "data": [
    {
      "id": 1,
      "employee_id": "EMP-001",
      "first_name": "John",
      "last_name": "Doe",
      "full_name": "John Doe",
      "initials": "JD",
      "status": "active",
      "photo_url": "https://example.com/photos/john.jpg",
      "email": "john.doe@company.com",
      "phone": "+255712345678",
      "work_phone": "+255222345678",
      "employment_type": "full_time",
      "hire_date": "2023-01-15",
      "designation": "Software Engineer",
      "department_id": 3,
      "department": {
        "id": 3,
        "name": "Engineering",
        "code": "ENG"
      },
      "position_id": 5,
      "position": {
        "id": 5,
        "title": "Software Engineer",
        "code": "SE"
      },
      "location_id": 1,
      "location": {
        "id": 1,
        "name": "Dar es Salaam HQ"
      },
      "team_id": 2,
      "team": {
        "id": 2,
        "name": "Backend Team"
      },
      "reporting_manager": {
        "id": 10,
        "employee_id": "EMP-010",
        "full_name": "Jane Smith"
      }
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

### Notes
- By default, only **active** employees are returned. Use `status` parameter to include other statuses.
- Fields like `photo_url`, `email`, `phone`, `work_phone`, `employment_type`, `department`, `position`, `location`, `team`, `reporting_manager`, `hire_date` are **conditional** — they only appear if data is available.
- The `initials` field is always computed from first and last name (e.g., "JD" for "John Doe").

---

## 2. Get Employee Profile

Get detailed directory profile for a specific employee.

**Endpoint:** `GET /api/v1/people-directory/:id`

### Path Parameters

| Parameter | Type           | Required | Description                              |
|-----------|----------------|----------|------------------------------------------|
| `id`      | integer/string | Yes      | Employee database ID or employee ID code |

### Example Requests

```bash
# By database ID
GET /api/v1/people-directory/1

# By employee ID code
GET /api/v1/people-directory/EMP-001
```

### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Employee profile retrieved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP-001",
    "first_name": "John",
    "last_name": "Doe",
    "full_name": "John Doe",
    "initials": "JD",
    "status": "active",
    "photo": "https://example.com/photos/john.jpg",
    "official_email": "john.doe@company.com",
    "work_phone": "+255222345678",
    "date_of_joining": "2023-01-15",
    "designation": "Software Engineer",
    "department": {
      "id": 3,
      "name": "Engineering",
      "code": "ENG"
    },
    "position": {
      "id": 5,
      "name": "Software Engineer",
      "code": "SE"
    },
    "location": {
      "id": 1,
      "name": "Dar es Salaam HQ",
      "address": "Plot 123, Bagamoyo Road, Dar es Salaam, Tanzania"
    },
    "reporting_manager": {
      "id": 10,
      "employee_id": "EMP-010",
      "full_name": "Jane Smith",
      "designation": null,
      "email": "jane.smith@company.com"
    }
  }
}
```

### Error Response (404 Not Found)

```json
{
  "success": false,
  "message": "Employee not found"
}
```

---

## 3. Get Filter Options

Returns all available filter options for building the directory filter UI (dropdowns, sidebar filters, alphabet bar, etc.).

**Endpoint:** `GET /api/v1/people-directory/filters`

### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Directory filter options retrieved successfully",
  "data": {
    "departments": [
      { "id": 1, "name": "Human Resources", "code": "HR" },
      { "id": 2, "name": "Finance", "code": "FIN" },
      { "id": 3, "name": "Engineering", "code": "ENG" }
    ],
    "positions": [
      { "id": 1, "title": "Software Engineer", "code": "SE" },
      { "id": 2, "title": "Project Manager", "code": "PM" },
      { "id": 3, "title": "HR Manager", "code": "HRM" }
    ],
    "locations": [
      { "id": 1, "name": "Dar es Salaam HQ" },
      { "id": 2, "name": "Arusha Branch" }
    ],
    "teams": [
      { "id": 1, "name": "Frontend Team" },
      { "id": 2, "name": "Backend Team" }
    ],
    "employment_types": [
      { "value": "full_time", "label": "Full Time" },
      { "value": "part_time", "label": "Part Time" },
      { "value": "contract", "label": "Contract" },
      { "value": "intern", "label": "Intern" }
    ],
    "statuses": [
      { "value": "active", "label": "Active" },
      { "value": "inactive", "label": "Inactive" },
      { "value": "on_leave", "label": "On Leave" },
      { "value": "suspended", "label": "Suspended" },
      { "value": "terminated", "label": "Terminated" }
    ],
    "alphabet": [
      { "letter": "A", "count": 12 },
      { "letter": "B", "count": 8 },
      { "letter": "C", "count": 5 },
      { "letter": "D", "count": 3 },
      { "letter": "J", "count": 15 },
      { "letter": "M", "count": 20 },
      { "letter": "S", "count": 10 }
    ]
  }
}
```

### Notes
- The `alphabet` array only includes letters that have at least one employee, with the count of active employees for each letter.
- Use these options to populate filter dropdowns and the alphabet navigation bar in the UI.
- All departments, positions, locations, and teams returned are **active** only.

---

## 4. Get Directory Statistics

Returns summary statistics and breakdowns for the directory dashboard/sidebar.

**Endpoint:** `GET /api/v1/people-directory/stats`

### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Directory statistics retrieved successfully",
  "data": {
    "total_active_employees": 150,
    "by_department": [
      {
        "department_id": 3,
        "department_name": "Engineering",
        "department_code": "ENG",
        "count": 45
      },
      {
        "department_id": 1,
        "department_name": "Human Resources",
        "department_code": "HR",
        "count": 12
      },
      {
        "department_id": 2,
        "department_name": "Finance",
        "department_code": "FIN",
        "count": 18
      }
    ],
    "by_location": [
      {
        "location_id": 1,
        "location_name": "Dar es Salaam HQ",
        "count": 120
      },
      {
        "location_id": 2,
        "location_name": "Arusha Branch",
        "count": 30
      }
    ],
    "by_employment_type": [
      { "employment_type": "full_time", "count": 130 },
      { "employment_type": "part_time", "count": 10 },
      { "employment_type": "contract", "count": 8 },
      { "employment_type": "intern", "count": 2 }
    ]
  }
}
```

### Notes
- Statistics only count **active** employees (`status = 'active'` and `is_active = true`).
- Use these counts to display distribution charts or sidebar summaries in the directory UI.

---

## 5. Data Models

### Directory Employee (List Item)

| Field               | Type            | Description                                      |
|---------------------|-----------------|--------------------------------------------------|
| `id`                | integer         | Employee database ID                             |
| `employee_id`       | string          | Employee ID code (e.g., "EMP-001")               |
| `first_name`        | string          | First name                                       |
| `last_name`         | string          | Last name                                        |
| `full_name`         | string          | Full name (first + last)                         |
| `initials`          | string          | Initials (e.g., "JD")                            |
| `status`            | string          | Employment status                                |
| `photo_url`         | string (opt)    | URL to employee photo                            |
| `email`             | string (opt)    | Work email (or personal if no work email)        |
| `phone`             | string (opt)    | Personal phone number                            |
| `work_phone`        | string (opt)    | Work phone number                                |
| `employment_type`   | string (opt)    | Employment type                                  |
| `hire_date`         | string (opt)    | Hire date (YYYY-MM-DD)                           |
| `designation`       | string (opt)    | Position/job title                               |
| `department_id`     | integer (opt)   | Department ID                                    |
| `department`        | object (opt)    | `{ id, name, code }`                             |
| `position_id`       | integer (opt)   | Position ID                                      |
| `position`          | object (opt)    | `{ id, title, code }`                            |
| `location_id`       | integer (opt)   | Location ID                                      |
| `location`          | object (opt)    | `{ id, name }`                                   |
| `team_id`           | integer (opt)   | Team ID                                          |
| `team`              | object (opt)    | `{ id, name }`                                   |
| `reporting_manager` | object (opt)    | `{ id, employee_id, full_name }`                 |

### Filter Options

| Field              | Type   | Description                                        |
|--------------------|--------|----------------------------------------------------|
| `departments`      | array  | `[{ id, name, code }]`                             |
| `positions`        | array  | `[{ id, title, code }]`                            |
| `locations`        | array  | `[{ id, name }]`                                   |
| `teams`            | array  | `[{ id, name }]`                                   |
| `employment_types` | array  | `[{ value, label }]`                               |
| `statuses`         | array  | `[{ value, label }]`                               |
| `alphabet`         | array  | `[{ letter, count }]` — only letters with employees|

---

## API Endpoints Summary

| Method | Endpoint                            | Description                       | Auth Required |
|--------|-------------------------------------|-----------------------------------|---------------|
| GET    | `/api/v1/people-directory`          | Search/browse employees           | Yes           |
| GET    | `/api/v1/people-directory/filters`  | Get filter options for UI         | Yes           |
| GET    | `/api/v1/people-directory/stats`    | Get directory statistics           | Yes           |
| GET    | `/api/v1/people-directory/:id`      | Get employee detail profile       | Yes           |

---

## Usage Examples

### Building a Directory Page

1. **On page load**, call both `/filters` and `/stats` in parallel to populate the sidebar filters and statistics.
2. **Search/browse**: Call `GET /people-directory` with any combination of filters. The response includes pagination metadata.
3. **Alphabet bar**: Use the `alphabet` data from `/filters` to render an A–Z navigation bar. Letters with `count > 0` are clickable; pass `?letter=A` to filter.
4. **Employee card click**: Navigate to the detail view using `GET /people-directory/:id`.
5. **Filter sidebar**: Use `departments`, `positions`, `locations`, `teams`, `employment_types` from `/filters` to populate dropdown selectors.

### Example: Filter by department and search

```bash
GET /api/v1/people-directory?department_id=3&search=john&page=1&page_size=10
```

### Example: Browse direct reports

```bash
GET /api/v1/people-directory?manager_id=15&sort_by=first_name&sort_order=asc
```

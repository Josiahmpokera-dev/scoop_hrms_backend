# Department Management API Documentation

## Overview

The Department module manages organizational departments with hierarchical structure, head-of-department assignment (linked to employee records), and location association.

**Base URL:** `/api/v1/departments`

**Authentication:** All endpoints require a valid JWT token and HR/Admin role.

---

## Table of Contents

1. [Create Department](#1-create-department)
2. [List Departments](#2-list-departments)
3. [Get Department](#3-get-department)
4. [Update Department](#4-update-department)
5. [Delete Department](#5-delete-department)
6. [Assign Department Head](#6-assign-department-head)
7. [Remove Department Head](#7-remove-department-head)
8. [Get Root Departments](#8-get-root-departments)
9. [Data Models](#9-data-models)
10. [Common Flows](#10-common-flows)

---

## 1. Create Department

```
POST /api/v1/departments
```

**Access:** Admin/HR only

**Request Body:**

```json
{
  "code": "ENG",
  "name": "Engineering",
  "description": "Software Engineering Department",
  "level": "department",
  "department_type": "core",
  "parent_department_id": 1,
  "manager_id": 5,
  "location_id": 2,
  "employee_capacity": 50
}
```

| Field               | Type   | Required | Description                                                  |
|---------------------|--------|----------|--------------------------------------------------------------|
| code                | string | Yes      | Unique department code (2-20 chars)                          |
| name                | string | Yes      | Department name (2-100 chars)                                |
| description         | string | No       | Description                                                  |
| level               | string | No       | `company`, `business_unit`, `department`, `team` (auto-inferred from parent if not provided) |
| department_type     | string | No       | `core`, `support`, `operational`, `strategic`                |
| parent_department_id| uint   | No       | Parent department ID (null = top-level)                      |
| manager_id          | uint   | No       | Head of Department (employee ID). Optional - can be assigned later via `/assign-head` |
| deputy_manager      | string | No       | Deputy manager name                                          |
| location_id         | uint   | No       | Location reference                                           |
| location            | string | No       | Location name (alternative to location_id)                   |
| organization_id     | uint   | No       | Organization reference                                       |
| organization_unit_id| uint   | No       | Business unit/division                                       |
| employee_capacity   | int    | No       | Maximum employee count                                       |
| is_active           | bool   | No       | Active status (default: true)                                |

> **Note:** `budget_allocated`, `budget_currency`, and `cost_center` fields are accepted but optional. They are not required for department creation and can be set later if needed.

**Response (201):**

```json
{
  "success": true,
  "message": "Department created successfully",
  "data": {
    "id": 3,
    "code": "ENG",
    "name": "Engineering",
    "description": "Software Engineering Department",
    "level": "department",
    "department_type": "core",
    "parent_department_id": 1,
    "head_of_department": {
      "id": 5,
      "employee_id": "EMP005",
      "full_name": "John Doe",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@company.com",
      "position": "Engineering Manager"
    },
    "employee_capacity": 50,
    "location": "Dar es Salaam HQ",
    "location_id": 2,
    "is_active": true,
    "created_at": "2026-02-06T10:00:00Z",
    "updated_at": "2026-02-06T10:00:00Z"
  }
}
```

---

## 2. List Departments

```
GET /api/v1/departments
```

**Access:** Admin/HR only

**Query Parameters:**

| Parameter            | Type   | Required | Description                             |
|----------------------|--------|----------|-----------------------------------------|
| page                 | int    | No       | Page number (default: 1)                |
| page_size            | int    | No       | Items per page (default: 20, max: 100)  |
| is_active            | string | No       | Filter: `true` or `false`               |
| parent_department_id | string | No       | Filter by parent department ID (`null` for root departments) |

**Response (200):**

```json
{
  "success": true,
  "message": "Departments retrieved successfully",
  "data": [
    {
      "id": 3,
      "code": "ENG",
      "name": "Engineering",
      "description": "Software Engineering Department",
      "level": "department",
      "department_type": "core",
      "parent_department_id": 1,
      "head_of_department": {
        "id": 5,
        "employee_id": "EMP005",
        "full_name": "John Doe",
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@company.com",
        "position": "Engineering Manager"
      },
      "employee_capacity": 50,
      "location": "Dar es Salaam HQ",
      "location_id": 2,
      "is_active": true,
      "created_at": "2026-02-06T10:00:00Z",
      "updated_at": "2026-02-06T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 5,
    "total_pages": 1
  }
}
```

---

## 3. Get Department

```
GET /api/v1/departments/:id
```

**Access:** Admin/HR only

**Response (200):** Same as single item from list above.

---

## 4. Update Department

```
PUT /api/v1/departments/:id
```

**Access:** Admin/HR only

**Request Body:** Same fields as Create (all optional for update). Only provided fields are updated.

```json
{
  "name": "Software Engineering",
  "employee_capacity": 60
}
```

**Response (200):** Returns updated department with enriched head-of-department data.

---

## 5. Delete Department

```
DELETE /api/v1/departments/:id
```

**Access:** Admin/HR only

**Validation:** Cannot delete a department that has sub-departments.

**Response (200):**

```json
{
  "success": true,
  "message": "Department deleted successfully",
  "data": null
}
```

---

## 6. Assign Department Head

```
POST /api/v1/departments/:id/assign-head
```

**Access:** Admin/HR only

**Description:** Assigns an employee as the Head of Department. The employee must exist and be active. This can be used at any time -- during department creation (via `manager_id`) or later via this dedicated endpoint.

**Request Body:**

```json
{
  "employee_id": 5
}
```

| Field       | Type | Required | Description                    |
|-------------|------|----------|--------------------------------|
| employee_id | uint | Yes      | Employee ID to assign as head  |

**Validation:**
- Employee must exist in the system
- Employee must be active

**Response (200):**

```json
{
  "success": true,
  "message": "Department head assigned successfully",
  "data": {
    "id": 3,
    "code": "ENG",
    "name": "Engineering",
    "head_of_department": {
      "id": 5,
      "employee_id": "EMP005",
      "full_name": "John Doe",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@company.com",
      "position": "Engineering Manager"
    },
    "is_active": true,
    "created_at": "2026-02-06T10:00:00Z",
    "updated_at": "2026-02-06T12:00:00Z"
  }
}
```

---

## 7. Remove Department Head

```
POST /api/v1/departments/:id/remove-head
```

**Access:** Admin/HR only

**Description:** Removes the current Head of Department assignment. The department will have no head until a new one is assigned.

**Request Body:** None required.

**Response (200):**

```json
{
  "success": true,
  "message": "Department head removed successfully",
  "data": {
    "id": 3,
    "code": "ENG",
    "name": "Engineering",
    "head_of_department": null,
    "is_active": true,
    "created_at": "2026-02-06T10:00:00Z",
    "updated_at": "2026-02-06T13:00:00Z"
  }
}
```

---

## 8. Get Root Departments

```
GET /api/v1/departments/root
```

**Access:** Admin/HR only

**Description:** Returns all top-level departments (no parent department).

**Response (200):** Array of departments with sub-departments included.

---

## 9. Data Models

### Department Response

| Field                | Type   | Description                                                  |
|----------------------|--------|--------------------------------------------------------------|
| id                   | uint   | Primary key                                                  |
| code                 | string | Unique department code                                       |
| name                 | string | Department name                                              |
| description          | string | Department description                                       |
| level                | string | `company`, `business_unit`, `department`, `team`             |
| department_type      | string | `core`, `support`, `operational`, `strategic`                |
| parent_department_id | uint   | Parent department reference (null for root)                  |
| head_of_department   | object | Enriched employee data for the department head (see below)   |
| deputy_manager       | string | Deputy manager name                                          |
| employee_capacity    | int    | Maximum employee count                                       |
| location             | string | Location name (resolved from location_id)                    |
| location_id          | uint   | Location reference                                           |
| is_active            | bool   | Active status                                                |
| created_at           | string | Creation timestamp                                           |
| updated_at           | string | Last update timestamp                                        |

### Head of Department (Employee Data)

When a department has a head assigned, the `head_of_department` field contains:

| Field       | Type   | Description                     |
|-------------|--------|---------------------------------|
| id          | uint   | Employee table primary key      |
| employee_id | string | Employee number (e.g., EMP005)  |
| full_name   | string | Full name                       |
| first_name  | string | First name                      |
| last_name   | string | Last name                       |
| email       | string | Work email                      |
| position    | string | Job position title              |

If no head is assigned, `head_of_department` is `null`.

---

## 10. Common Flows

### Flow 1: Create Department and Assign Head Later

1. Create department without a head:

```json
POST /api/v1/departments
{
  "code": "MKT",
  "name": "Marketing",
  "level": "department"
}
```

2. Later, assign an employee as head:

```json
POST /api/v1/departments/4/assign-head
{
  "employee_id": 12
}
```

### Flow 2: Create Department with Head at Creation

```json
POST /api/v1/departments
{
  "code": "FIN",
  "name": "Finance",
  "level": "department",
  "manager_id": 8
}
```

The response will include the enriched `head_of_department` data with the employee's name, email, and position.

### Flow 3: Change Department Head

1. Assign the new head (replaces the old one automatically):

```json
POST /api/v1/departments/3/assign-head
{
  "employee_id": 15
}
```

### Flow 4: Remove Department Head

```json
POST /api/v1/departments/3/remove-head
```

### Flow 5: Auto-Suggest Reporting Manager During Onboarding

When onboarding a new employee (Step 2 - Employment Details):
1. Select a department
2. The system automatically sets the department head as the reporting manager
3. The user can override this with a different manager or clear it entirely

**Related endpoints:**
- `GET /api/v1/employees/department-manager/:department_id` -- Get the suggested manager for a department
- `GET /api/v1/employees/managers?department_id=X` -- List managers filtered by department (flags the head with `is_suggested: true`)

---

## Notes

1. **Budget and Cost Center fields** are still supported in the API (for backward compatibility) but are not required for department creation. These fields (`budget_allocated`, `budget_currency`, `cost_center`) can be omitted entirely from requests.

2. **Head of Department** is linked to real employee data. The API response always returns enriched employee information (name, email, position) instead of just an ID.

3. **Hierarchical Structure:** Departments support parent-child relationships. The `level` field is auto-inferred from the parent if not explicitly provided (company -> business_unit -> department -> team).

4. **Soft Deletes:** Deleted departments are preserved in the database. A department with sub-departments cannot be deleted.

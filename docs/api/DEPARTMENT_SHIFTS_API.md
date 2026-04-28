# Department Shifts API - Changes Documentation

This document defines the API changes for linking **Shifts** to **Departments**.

**Page**: Department Management
**Route**: `/departments`
**Feature**: NOC Department shift support with 07:30 morning grace exception

---

## Base URL
`http://localhost:8080/api/v1`

## Authentication
All endpoints require a valid Bearer token.

`Authorization: Bearer <token>`

---

## Standard Response Envelope

### Success
```json
{
  "success": true,
  "message": "string",
  "data": {}
}
```

### Error
```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "field": "reason"
  }
}
```

---

# 1) Create Department (UPDATED)

## Endpoint
`POST /departments`

## Changes
- **Added**: `shift_ids` field (array of shift IDs to associate with this department)

## Request Body
```json
{
  "organization_id": 1,
  "code": "NOC",
  "name": "NOC Department",
  "level": "department",
  "department_type": "operational",
  "parent_department_id": null,
  "manager_id": 5,
  "deputy_manager": "John Doe",
  "location_id": 1,
  "cost_center": "NOC-001",
  "budget_allocated": 500000.00,
  "budget_currency": "TZS",
  "employee_capacity": 50,
  "shift_ids": [1, 2, 3],
  "is_active": true
}
```

### Field Descriptions
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `shift_ids` | `[]uint` | No | Array of shift IDs to associate with this department. Links department to specific shifts (e.g., NOC has 3 rotating shifts). |

## Validation
- Each `shift_id` must reference an existing shift
- If any `shift_id` is invalid, returns error: `"shift with ID '%d' not found"`

## Success Response (201 Created)
```json
{
  "success": true,
  "message": "Department created successfully",
  "data": {
    "id": 10,
    "organization_id": 1,
    "code": "NOC",
    "name": "NOC Department",
    "level": "department",
    "department_type": "operational",
    "parent_department_id": null,
    "head_of_department": {
      "id": 5,
      "employee_id": "EMP-001",
      "full_name": "Jane Smith",
      "first_name": "Jane",
      "last_name": "Smith",
      "email": "jane.smith@company.com",
      "position": "NOC Manager"
    },
    "employee_capacity": 50,
    "location": "Head Office",
    "location_id": 1,
    "shifts": [
      {
        "id": 1,
        "shift_name": "NOC Morning",
        "shift_code": "NOC-M",
        "start_time": "07:00:00",
        "end_time": "15:00:00"
      },
      {
        "id": 2,
        "shift_name": "NOC Evening",
        "shift_code": "NOC-E",
        "start_time": "15:00:00",
        "end_time": "23:00:00"
      },
      {
        "id": 3,
        "shift_name": "NOC Night",
        "shift_code": "NOC-N",
        "start_time": "23:00:00",
        "end_time": "05:00:00"
      }
    ],
    "is_active": true,
    "created_at": "2026-04-28T10:00:00Z",
    "updated_at": "2026-04-28T10:00:00Z",
    "updated_by": 1
  }
}
```

---

# 2) Update Department (UPDATED)

## Endpoint
`PUT /departments/{id}`

## Changes
- **Added**: `shift_ids` field (array of shift IDs to update department-shift associations)

## Request Body
```json
{
  "name": "NOC Operations Department",
  "shift_ids": [1, 2, 3, 4]
}
```

### Field Descriptions
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `shift_ids` | `[]uint` | No | Array of shift IDs to associate with this department. Replaces existing associations. |

## Validation
- Each `shift_id` must reference an existing shift
- If any `shift_id` is invalid, returns error: `"shift with ID '%d' not found"`
- Setting `shift_ids` to empty array `[]` removes all shift associations

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Department updated successfully",
  "data": {
    "id": 10,
    "code": "NOC",
    "name": "NOC Operations Department",
    "shifts": [
      {
        "id": 1,
        "shift_name": "NOC Morning",
        "shift_code": "NOC-M",
        "start_time": "07:00:00",
        "end_time": "15:00:00"
      },
      {
        "id": 2,
        "shift_name": "NOC Evening",
        "shift_code": "NOC-E",
        "start_time": "15:00:00",
        "end_time": "23:00:00"
      },
      {
        "id": 3,
        "shift_name": "NOC Night",
        "shift_code": "NOC-N",
        "start_time": "23:00:00",
        "end_time": "05:00:00"
      },
      {
        "id": 4,
        "shift_name": "NOC Weekend",
        "shift_code": "NOC-W",
        "start_time": "08:00:00",
        "end_time": "16:00:00"
      }
    ],
    "is_active": true,
    "updated_at": "2026-04-28T11:00:00Z"
  }
}
```

---

# 3) Get Department by ID (UNCHANGED - Response enhanced)

## Endpoint
`GET /departments/{id}`

## Response Changes
- **Enhanced**: Response now includes `shifts` array with shift details

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Department retrieved successfully",
  "data": {
    "id": 10,
    "organization_id": 1,
    "code": "NOC",
    "name": "NOC Department",
    "level": "department",
    "department_type": "operational",
    "parent_department_id": null,
    "head_of_department": {
      "id": 5,
      "employee_id": "EMP-001",
      "full_name": "Jane Smith",
      "first_name": "Jane",
      "last_name": "Smith",
      "email": "jane.smith@company.com",
      "position": "NOC Manager"
    },
    "employee_capacity": 50,
    "location": "Head Office",
    "location_id": 1,
    "shifts": [
      {
        "id": 1,
        "shift_name": "NOC Morning",
        "shift_code": "NOC-M",
        "start_time": "07:00:00",
        "end_time": "15:00:00"
      },
      {
        "id": 2,
        "shift_name": "NOC Evening",
        "shift_code": "NOC-E",
        "start_time": "15:00:00",
        "end_time": "23:00:00"
      },
      {
        "id": 3,
        "shift_name": "NOC Night",
        "shift_code": "NOC-N",
        "start_time": "23:00:00",
        "end_time": "05:00:00"
      }
    ],
    "is_active": true,
    "created_at": "2026-04-28T10:00:00Z",
    "updated_at": "2026-04-28T10:00:00Z",
    "updated_by": 1
  }
}
```

---

# 4) List Departments (UNCHANGED - Response enhanced)

## Endpoint
`GET /departments`

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | number | No | Default `1` |
| `page_size` | number | No | Default `10` |
| `is_active` | boolean | No | Filter by active status |
| `parent_department_id` | number | No | Filter by parent department |

## Response Changes
- **Enhanced**: Each department in list now includes `shifts` array

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Departments retrieved successfully",
  "data": {
    "data": [
      {
        "id": 10,
        "code": "NOC",
        "name": "NOC Department",
        "level": "department",
        "shifts": [
          {
            "id": 1,
            "shift_name": "NOC Morning",
            "shift_code": "NOC-M",
            "start_time": "07:00:00",
            "end_time": "15:00:00"
          },
          {
            "id": 2,
            "shift_name": "NOC Evening",
            "shift_code": "NOC-E",
            "start_time": "15:00:00",
            "end_time": "23:00:00"
          },
          {
            "id": 3,
            "shift_name": "NOC Night",
            "shift_code": "NOC-N",
            "start_time": "23:00:00",
            "end_time": "05:00:00"
          }
        ],
        "is_active": true
      },
      {
        "id": 11,
        "code": "HR",
        "name": "Human Resources",
        "level": "department",
        "shifts": [
          {
            "id": 4,
            "shift_name": "Regular",
            "shift_code": "REG-1",
            "start_time": "08:30:00",
            "end_time": "17:30:00"
          }
        ],
        "is_active": true
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 2,
      "total_pages": 1
    }
  }
}
```

---

# 5) Shift Response Structure

When shifts are returned within department responses, each shift includes:

```json
{
  "id": 1,
  "shift_name": "NOC Morning",
  "shift_code": "NOC-M",
  "start_time": "07:00:00",
  "end_time": "15:00:00",
  "late_threshold_time": "07:30:00"
}
```

### Field Descriptions
| Field | Type | Description |
|-------|------|-------------|
| `id` | uint | Unique shift identifier |
| `shift_name` | string | Human-readable shift name (e.g., "NOC Morning") |
| `shift_code` | string | Short code for the shift (e.g., "NOC-M") |
| `start_time` | string | Shift start time in HH:MM:SS format |
| `end_time` | string | Shift end time in HH:MM:SS format |
| `late_threshold_time` | string | Late exception threshold in HH:MM:SS format. If an employee checks in after this time, they are marked as late. Defaults to `start_time` if not set. |

---

# 6) NOC Department Implementation Example

## Step 1: Create NOC Shifts

First, create the three NOC shifts via the Shifts API (`POST /shifts`):

### NOC Morning Shift (with 07:30 grace exception)
```json
{
  "shift_name": "NOC Morning",
  "shift_code": "NOC-M",
  "shift_type": "rotating",
  "start_time": "07:00:00",
  "end_time": "15:00:00",
  "grace_minutes": 30,
  "break_duration": 60,
  "late_threshold_time": "07:30:00",
  "is_active": true
}
```

### NOC Evening Shift
```json
{
  "shift_name": "NOC Evening",
  "shift_code": "NOC-E",
  "shift_type": "rotating",
  "start_time": "15:00:00",
  "end_time": "23:00:00",
  "grace_minutes": 0,
  "break_duration": 60,
  "is_active": true
}
```

### NOC Night Shift (cross-day)
```json
{
  "shift_name": "NOC Night",
  "shift_code": "NOC-N",
  "shift_type": "rotating",
  "start_time": "23:00:00",
  "end_time": "05:00:00",
  "cross_day": true,
  "night_shift": true,
  "grace_minutes": 0,
  "break_duration": 60,
  "is_active": true
}
```

## Step 2: Create NOC Department with Shift Association

```json
{
  "code": "NOC",
  "name": "NOC Department",
  "level": "department",
  "department_type": "operational",
  "manager_id": 5,
  "location_id": 1,
  "shift_ids": [1, 2, 3],
  "is_active": true
}
```

---

# Summary of Changes for Frontend

## Department-Shift Association Changes

| Endpoint | Method | Change |
|----------|--------|--------|
| `/departments` | POST | New field `shift_ids` added to request |
| `/departments/{id}` | PUT | New field `shift_ids` added to request |
| `/departments/{id}` | GET | Response now includes `shifts` array |
| `/departments` | GET | Response now includes `shifts` array in each department |

### Key Points for Frontend Developer:
1. **Create/Update Department forms**: Add a multi-select field for `shift_ids`
2. **Department list/detail views**: Display associated shifts
3. **Validation**: If a selected shift doesn't exist, API returns error with shift ID
4. **Backward compatible**: Existing departments without shifts will have `shifts: []`

## Shift Model Changes (Late Exception Time)

| Endpoint | Method | Change |
|----------|--------|--------|
| `/shifts` | POST | New field `late_threshold_time` can be set for shift-specific late exception |
| `/shifts/{id}` | PUT | New field `late_threshold_time` can be updated |
| `/shifts/{id}` | GET | Response now includes `late_threshold_time` field |
| `/shifts` | GET | Response now includes `late_threshold_time` field in list |

### Late Threshold Time Field
| Field | Type | Description |
|-------|------|-------------|
| `late_threshold_time` | string | HH:MM:SS format. Defines when an employee is considered late. For NOC Morning shift, this is `07:30:00` (the 07:30 exception you specified). If not set, defaults to the shift's `start_time`. |

## Attendance API Changes (Department & NOC Support)

| Endpoint | Method | Change |
|----------|--------|--------|
| `/attendance/overview` | GET | Response now includes department information for each employee |
| `/attendance/department-stats` | GET | Response includes department breakdowns with shift-specific late thresholds |
| `/attendance/summary` | GET | Uses shift-specific late thresholds via RosterAssignment joins |
| `/attendance/export` | GET | Export now includes department name column |

### Key Points for Attendance Reports:
1. **Late calculation**: Now uses shift-specific `late_threshold_time` instead of hardcoded 09:00:00
2. **Department tracking**: Attendance data is linked to employee departments
3. **NOC support**: NOC employees with different shifts have their late status calculated based on their assigned shift's threshold
4. **Export includes department**: CSV/Excel exports now have a `department_name` column

# Position by Department API Documentation

## Overview

The Position by Department API allows you to retrieve job positions filtered by a specific department. This ensures that when a user selects a department (e.g., "Information Technology"), they only see positions that belong to that department (e.g., "Software Developer", "Network Engineer", etc.), not positions from other departments.

**Base URL:** `/api/v1/job-positions`

**Authentication:** All endpoints require authentication via JWT token in the `Authorization` header and HR/Admin role.

---

## API Endpoint

### Get Positions by Department

Get all active job positions for a specific department. This endpoint ensures that positions are filtered by department, preventing users from seeing irrelevant positions when selecting a department.

**Endpoint:** `GET /api/v1/job-positions/department/:department_id`

**Authentication:** Required

**URL Parameters:**
- `department_id` (required) - Department ID

**Response:**
```json
{
  "success": true,
  "message": "Positions retrieved successfully",
  "data": [
    {
      "id": 1,
      "code": "DEV001",
      "title": "Software Developer",
      "grade": "G5",
      "department_id": 3,
      "reports_to_position": "Senior Software Developer",
      "reports_to_position_id": 2,
      "budgeted_headcount": 10,
      "current_headcount": 5,
      "employment_type": "full_time",
      "key_competencies": "Java, Python, React",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "code": "NET001",
      "title": "Network Engineer",
      "grade": "G4",
      "department_id": 3,
      "reports_to_position": null,
      "reports_to_position_id": null,
      "budgeted_headcount": 5,
      "current_headcount": 2,
      "employment_type": "full_time",
      "key_competencies": "Cisco, Network Security",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 3,
      "code": "MOB001",
      "title": "Mobile Developer",
      "grade": "G5",
      "department_id": 3,
      "reports_to_position": null,
      "reports_to_position_id": null,
      "budgeted_headcount": 8,
      "current_headcount": 3,
      "employment_type": "full_time",
      "key_competencies": "React Native, Flutter",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 4,
      "code": "IT001",
      "title": "IT Technician",
      "grade": "G3",
      "department_id": 3,
      "reports_to_position": "IT Manager",
      "reports_to_position_id": 5,
      "budgeted_headcount": 6,
      "current_headcount": 4,
      "employment_type": "full_time",
      "key_competencies": "Hardware, Software Support",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**Example Request:**
```bash
GET {{BASE_URL}}/api/v1/job-positions/department/3
Authorization: Bearer {{TOKEN}}
```

**Example with cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/job-positions/department/3" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json"
```

---

## Use Case Example

### Scenario: Employee Onboarding - Step 2

When a user is onboarding an employee and selects a department:

1. **User selects Department:** "Information Technology" (ID: 3)
2. **Frontend calls API:** `GET /api/v1/job-positions/department/3`
3. **API returns only IT positions:**
   - Software Developer
   - Network Engineer
   - Mobile Developer
   - IT Technician
   - DevOps Engineer
   - etc.

4. **User will NOT see positions from other departments like:**
   - Human Resources positions (HR Manager, Recruiter, etc.)
   - Finance positions (Accountant, Financial Analyst, etc.)
   - Marketing positions (Marketing Manager, Content Writer, etc.)

This ensures data integrity and a better user experience.

---

## Response Fields

### Position Object

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Position ID |
| `code` | string | Position code (unique identifier) |
| `title` | string | Position title/name |
| `grade` | string | Grade level (e.g., G3, G4, G5) |
| `department_id` | integer | Department ID this position belongs to |
| `reports_to_position` | string | Title of the position this reports to (if any) |
| `reports_to_position_id` | integer | ID of the position this reports to (if any) |
| `budgeted_headcount` | integer | Budgeted number of employees for this position |
| `current_headcount` | integer | Current number of employees in this position |
| `employment_type` | string | Employment type (full_time, part_time, contract, intern) |
| `key_competencies` | string | Key competencies/skills required (comma-separated) |
| `is_active` | boolean | Whether the position is active |
| `created_at` | string | Creation timestamp (ISO 8601) |
| `updated_at` | string | Last update timestamp (ISO 8601) |

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Department ID is required",
  "data": null
}
```

**Or:**
```json
{
  "success": false,
  "message": "Invalid department ID",
  "data": null
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "department with ID 999 not found",
  "data": null
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized",
  "data": null
}
```

### 403 Forbidden
```json
{
  "success": false,
  "message": "department with ID 3 does not belong to your tenant",
  "data": null
}
```

---

## JavaScript Examples

### Get Positions by Department
```javascript
const getPositionsByDepartment = async (departmentId) => {
  const response = await fetch(
    `${BASE_URL}/api/v1/job-positions/department/${departmentId}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    }
  );
  return await response.json();
};

// Example usage
const itPositions = await getPositionsByDepartment(3);
console.log('IT Positions:', itPositions.data);
```

### React Example - Dynamic Position Selection
```javascript
import { useState, useEffect } from 'react';

function PositionSelector({ departmentId, onPositionSelect }) {
  const [positions, setPositions] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (departmentId) {
      setLoading(true);
      fetch(`${BASE_URL}/api/v1/job-positions/department/${departmentId}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      })
        .then(res => res.json())
        .then(data => {
          if (data.success) {
            setPositions(data.data);
          }
        })
        .catch(err => console.error('Error:', err))
        .finally(() => setLoading(false));
    } else {
      setPositions([]);
    }
  }, [departmentId]);

  return (
    <select 
      onChange={(e) => onPositionSelect(e.target.value)}
      disabled={loading || !departmentId}
    >
      <option value="">Select Position</option>
      {positions.map(position => (
        <option key={position.id} value={position.id}>
          {position.title}
        </option>
      ))}
    </select>
  );
}
```

---

## Alternative: Using Query Parameter

You can also filter positions by department using the list endpoint with a query parameter:

**Endpoint:** `GET /api/v1/job-positions?department_id=3&is_active=true`

**Example:**
```bash
GET {{BASE_URL}}/api/v1/job-positions?department_id=3&is_active=true&page=1&page_size=20
Authorization: Bearer {{TOKEN}}
```

This returns paginated results with metadata, while the `/department/:department_id` endpoint returns all active positions for that department without pagination.

---

## Notes

1. **Active Positions Only:** The `/department/:department_id` endpoint returns only active positions (`is_active = true`).

2. **Tenant Isolation:** Positions are automatically filtered by tenant. You can only see positions belonging to your tenant.

3. **Department Validation:** The API validates that:
   - The department exists
   - The department belongs to your tenant
   - The department is active

4. **Sorted Results:** Positions are returned sorted by title in ascending order.

5. **No Pagination:** The `/department/:department_id` endpoint returns all matching positions. For paginated results, use the list endpoint with `department_id` query parameter.

6. **Empty Results:** If a department has no active positions, the API returns an empty array `[]`.

---

## Integration with Employee Onboarding

This API is particularly useful in the employee onboarding flow:

**Step 2 - Employment Details:**
1. User selects `department_id` (e.g., Information Technology)
2. Frontend calls `GET /api/v1/job-positions/department/3`
3. Position dropdown is populated with only IT positions
4. User selects appropriate position (e.g., Software Developer)
5. Form submission includes both `department_id` and `position_id`

This ensures data consistency and prevents invalid combinations like:
- ❌ Department: Information Technology, Position: HR Manager
- ✅ Department: Information Technology, Position: Software Developer

---

## Summary

The Position by Department API provides:

- ✅ Filtered position list by department
- ✅ Only active positions returned
- ✅ Tenant isolation and security
- ✅ Department validation
- ✅ Sorted results (by title)
- ✅ Easy integration with frontend dropdowns
- ✅ Prevents invalid department-position combinations
- ✅ Better user experience in forms

This API solves the problem of showing irrelevant positions when a department is selected, ensuring data integrity and a smoother user experience.

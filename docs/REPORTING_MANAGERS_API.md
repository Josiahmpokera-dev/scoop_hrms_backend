# Reporting Managers API Documentation

## Overview

The Reporting Managers API allows you to list all active employees who can be used as reporting managers during employee onboarding. This is useful when filling out Step 2 (Employment Details) where you need to specify a `reporting_manager_id`.

**Base URL:** `/api/v1/employees/managers`

**Authentication:** All endpoints require authentication via JWT token in the `Authorization` header.

---

## API Endpoint

### List Reporting Managers

Get a list of all active employees who can be reporting managers. This includes their ID, employee ID, name, email, department, and position information.

**Endpoint:** `GET /api/v1/employees/managers`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "message": "Managers retrieved successfully",
  "data": [
    {
      "id": 20,
      "employee_id": "EMP001",
      "full_name": "John Smith",
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@company.com",
      "department_id": 3,
      "department": "Engineering",
      "position_id": 5,
      "position": "Senior Manager"
    },
    {
      "id": 25,
      "employee_id": "EMP002",
      "full_name": "Jane Doe",
      "first_name": "Jane",
      "last_name": "Doe",
      "email": "jane.doe@company.com",
      "department_id": 2,
      "department": "Sales",
      "position_id": 8,
      "position": "Sales Manager"
    },
    {
      "id": 30,
      "employee_id": "EMP003",
      "full_name": "Michael Johnson",
      "first_name": "Michael",
      "last_name": "Johnson",
      "email": "michael.johnson@company.com",
      "department_id": 1,
      "department": "HR",
      "position_id": 3,
      "position": "HR Manager"
    }
  ]
}
```

**Example Request:**
```bash
GET {{BASE_URL}}/api/v1/employees/managers
Authorization: Bearer {{TOKEN}}
```

**Response Fields:**
- `id` - Employee database ID (use this as `reporting_manager_id` in Step 2)
- `employee_id` - Employee ID string (e.g., EMP001)
- `full_name` - Full name (First Name + Last Name)
- `first_name` - First name
- `last_name` - Last name
- `email` - Work email address
- `department_id` - Department ID (optional)
- `department` - Department name (optional)
- `position_id` - Position ID (optional)
- `position` - Position/Job title (optional)

**Note:** Only active employees (`status = 'active'` and `is_active = true`) are returned. The list is filtered by tenant, so you only see managers from your organization.

---

## Usage Example

### Step 1: Get List of Managers

```bash
GET {{BASE_URL}}/api/v1/employees/managers
Authorization: Bearer {{TOKEN}}
```

### Step 2: Use Manager ID in Onboarding Step 2

After getting the list of managers, use the `id` field as the `reporting_manager_id` in your Step 2 request:

```json
{
  "step": 2,
  "data": {
    "official_email": "syncia.daudi@company.com",
    "date_of_joining": "2024-01-15T00:00:00Z",
    "department_id": 3,
    "position_id": 5,
    "grade": "G5",
    "reporting_manager_id": 20,
    "employment_type": "full_time",
    "location_id": 1,
    "shift": "Day",
    "work_phone": "+255222345678",
    "probation_period_days": 90,
    "expected_confirmation_date": "2024-04-15T00:00:00Z"
  }
}
```

**Example:** If you want to assign "John Smith" as the reporting manager, use `"reporting_manager_id": 20` (the `id` from the managers list).

---

## JavaScript Example

```javascript
// Get list of managers
const getManagers = async () => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/managers`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  const result = await response.json();
  return result.data;
};

// Use in onboarding Step 2
const saveStep2 = async (employeeId, step2Data) => {
  // Get managers first
  const managers = await getManagers();
  
  // Find manager by name or employee_id
  const selectedManager = managers.find(m => 
    m.full_name === "John Smith" || m.employee_id === "EMP001"
  );
  
  // Use manager ID in step 2 data
  const step2Payload = {
    step: 2,
    data: {
      ...step2Data,
      reporting_manager_id: selectedManager.id
    }
  };
  
  const response = await fetch(
    `${BASE_URL}/api/v1/employees/onboarding/${employeeId}/step/2`,
    {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(step2Payload)
    }
  );
  
  return await response.json();
};
```

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Failed to get managers: <error details>",
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

---

## Notes

1. **Tenant Isolation:** Only managers from your tenant/organization are returned.

2. **Active Employees Only:** Only employees with `status = 'active'` and `is_active = true` are included in the list.

3. **Sorted by Name:** Managers are sorted alphabetically by first name, then last name.

4. **Optional Fields:** Department and position information may be `null` if not assigned to the employee.

5. **Use ID Field:** Always use the `id` field (not `employee_id`) as the `reporting_manager_id` in Step 2 requests.

---

## Summary

The Reporting Managers API provides a simple way to:
- ✅ List all active employees who can be reporting managers
- ✅ Get their database ID for use in onboarding Step 2
- ✅ View their name, email, department, and position
- ✅ Filter by tenant automatically
- ✅ Use the `id` field as `reporting_manager_id` in Step 2

This makes it easy to select the correct reporting manager when onboarding new employees.

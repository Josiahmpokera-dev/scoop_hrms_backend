# Employee Status Management API Documentation

This document provides comprehensive documentation for managing employee statuses: **Terminate**, **Suspend**, and **Archive**, along with **Reactivate** functionality.

## Table of Contents

1. [Overview](#overview)
2. [Employee Statuses](#employee-statuses)
3. [Update Employee Information](#update-employee-information)
4. [API Endpoints](#api-endpoints)
   - [Terminate Employee](#1-terminate-employee)
   - [Suspend Employee](#2-suspend-employee)
   - [Archive Employee](#3-archive-employee)
   - [Reactivate Employee](#4-reactivate-employee)
5. [Usage Examples](#usage-examples)
6. [Error Handling](#error-handling)

---

## Overview

The Employee Status Management API allows HR administrators to manage employee employment statuses. These operations are critical for maintaining accurate employee records and ensuring proper access control.

**Base URL**: `{{BASE_URL}}/api/v1/employees`

**Authentication**: All endpoints require authentication (JWT token) and HR/Admin role.

---

## Employee Statuses

The system supports the following employee statuses:

| Status | Description | IsActive | Can Reactivate? |
|--------|-------------|-----------|-----------------|
| `active` | Employee is currently active | `true` | N/A |
| `inactive` | Employee is inactive | `false` | Yes |
| `on_leave` | Employee is on leave | `true` | N/A |
| `suspended` | Employee is suspended | `false` | Yes |
| `terminated` | Employee is terminated | `false` | No |
| `archived` | Employee is archived | `false` | Yes |

---

## Update Employee Information

### Option A: Using Dedicated Update API (Recommended)

**Easier and More RESTful**

Use the existing `PUT /api/v1/employees/:id` endpoint to update employee information. This is the recommended approach as it's:
- ✅ Clear and explicit
- ✅ Follows REST conventions
- ✅ Easy to understand and maintain
- ✅ Supports partial updates

**Example:**
```bash
PUT {{BASE_URL}}/api/v1/employees/123
Content-Type: application/json
Authorization: Bearer {{TOKEN}}

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@company.com",
  "department_id": 5,
  "position_id": 10,
  "salary": 50000
}
```

### Option B: Using Onboarding Flow (For Bulk Updates)

**More Complex but Useful for Multi-Step Updates**

You can use the onboarding flow to update employee information, but this requires:
- Creating a draft for an existing employee
- Going through onboarding steps
- Completing the onboarding

**Note**: This approach is more complex and is primarily designed for new employee onboarding. For updates, **Option A is recommended**.

---

## API Endpoints

### 1. Terminate Employee

Terminates an employee, setting their status to `terminated` and `is_active` to `false`. This action is **irreversible** (cannot be reactivated).

**Endpoint**: `POST /api/v1/employees/terminate`

**Request Body:**
```json
{
  "id": 123,
  "reason": "Voluntary resignation",
  "termination_date": "2024-12-31T00:00:00Z"
}
```

**Parameters:**
- `id` (required): Employee ID (database ID, not employee_id string)
- `reason` (optional): Reason for termination (will be added to employee notes)
- `termination_date` (optional): Date of termination (defaults to current date if not provided)

**Response:**
```json
{
  "success": true,
  "message": "Employee terminated successfully",
  "data": {
    "id": 123,
    "employee_id": "EMP001",
    "first_name": "John",
    "last_name": "Doe",
    "status": "terminated",
    "is_active": false,
    "notes": "Termination Reason: Voluntary resignation",
    "updated_at": "2024-12-31T10:00:00Z"
  }
}
```

**cURL Example:**
```bash
curl -X POST "{{BASE_URL}}/api/v1/employees/terminate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{TOKEN}}" \
  -d '{
    "id": 123,
    "reason": "Voluntary resignation",
    "termination_date": "2024-12-31T00:00:00Z"
  }'
```

**JavaScript Example:**
```javascript
const terminateEmployee = async (employeeId, reason, terminationDate) => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/terminate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      id: employeeId,
      reason: reason,
      termination_date: terminationDate
    })
  });
  
  const result = await response.json();
  return result;
};

// Usage
terminateEmployee(123, "Voluntary resignation", "2024-12-31T00:00:00Z");
```

---

### 2. Suspend Employee

Suspends an employee, setting their status to `suspended` and `is_active` to `false`. Suspended employees can be reactivated.

**Endpoint**: `POST /api/v1/employees/suspend`

**Request Body:**
```json
{
  "id": 123,
  "reason": "Pending investigation",
  "suspension_end_date": "2025-01-15T00:00:00Z"
}
```

**Parameters:**
- `id` (required): Employee ID (database ID, not employee_id string)
- `reason` (optional): Reason for suspension (will be added to employee notes)
- `suspension_end_date` (optional): Expected end date of suspension

**Response:**
```json
{
  "success": true,
  "message": "Employee suspended successfully",
  "data": {
    "id": 123,
    "employee_id": "EMP001",
    "first_name": "John",
    "last_name": "Doe",
    "status": "suspended",
    "is_active": false,
    "notes": "Suspension Reason: Pending investigation (Suspension ends: 2025-01-15)",
    "updated_at": "2024-12-31T10:00:00Z"
  }
}
```

**cURL Example:**
```bash
curl -X POST "{{BASE_URL}}/api/v1/employees/suspend" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{TOKEN}}" \
  -d '{
    "id": 123,
    "reason": "Pending investigation",
    "suspension_end_date": "2025-01-15T00:00:00Z"
  }'
```

**JavaScript Example:**
```javascript
const suspendEmployee = async (employeeId, reason, suspensionEndDate) => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/suspend`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      id: employeeId,
      reason: reason,
      suspension_end_date: suspensionEndDate
    })
  });
  
  const result = await response.json();
  return result;
};

// Usage
suspendEmployee(123, "Pending investigation", "2025-01-15T00:00:00Z");
```

---

### 3. Archive Employee

Archives an employee, setting their status to `archived` and `is_active` to `false`. Archived employees can be reactivated.

**Endpoint**: `POST /api/v1/employees/archive`

**Request Body:**
```json
{
  "id": 123,
  "reason": "Long-term leave of absence"
}
```

**Parameters:**
- `id` (required): Employee ID (database ID, not employee_id string)
- `reason` (optional): Reason for archiving (will be added to employee notes)

**Response:**
```json
{
  "success": true,
  "message": "Employee archived successfully",
  "data": {
    "id": 123,
    "employee_id": "EMP001",
    "first_name": "John",
    "last_name": "Doe",
    "status": "archived",
    "is_active": false,
    "notes": "Archived Reason: Long-term leave of absence",
    "updated_at": "2024-12-31T10:00:00Z"
  }
}
```

**cURL Example:**
```bash
curl -X POST "{{BASE_URL}}/api/v1/employees/archive" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{TOKEN}}" \
  -d '{
    "id": 123,
    "reason": "Long-term leave of absence"
  }'
```

**JavaScript Example:**
```javascript
const archiveEmployee = async (employeeId, reason) => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/archive`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      id: employeeId,
      reason: reason
    })
  });
  
  const result = await response.json();
  return result;
};

// Usage
archiveEmployee(123, "Long-term leave of absence");
```

---

### 4. Reactivate Employee

Reactivates a suspended or archived employee, setting their status to `active` and `is_active` to `true`. **Cannot reactivate terminated employees**.

**Endpoint**: `POST /api/v1/employees/reactivate`

**Request Body:**
```json
{
  "id": 123
}
```

**Parameters:**
- `id` (required): Employee ID (database ID, not employee_id string)

**Response:**
```json
{
  "success": true,
  "message": "Employee reactivated successfully",
  "data": {
    "id": 123,
    "employee_id": "EMP001",
    "first_name": "John",
    "last_name": "Doe",
    "status": "active",
    "is_active": true,
    "updated_at": "2024-12-31T10:00:00Z"
  }
}
```

**cURL Example:**
```bash
curl -X POST "{{BASE_URL}}/api/v1/employees/reactivate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{TOKEN}}" \
  -d '{
    "id": 123
  }'
```

**JavaScript Example:**
```javascript
const reactivateEmployee = async (employeeId) => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/reactivate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      id: employeeId
    })
  });
  
  const result = await response.json();
  return result;
};

// Usage
reactivateEmployee(123);
```

---

## Usage Examples

### Complete Workflow Example

```javascript
// 1. Suspend an employee
await suspendEmployee(123, "Pending investigation", "2025-01-15T00:00:00Z");

// 2. Later, reactivate the employee
await reactivateEmployee(123);

// 3. If needed, terminate the employee (irreversible)
await terminateEmployee(123, "Voluntary resignation", "2024-12-31T00:00:00Z");

// 4. Archive an employee for long-term leave
await archiveEmployee(123, "Long-term leave of absence");

// 5. Reactivate archived employee
await reactivateEmployee(123);
```

### React Native Example

```javascript
import axios from 'axios';

const API_BASE_URL = 'https://your-api.com';
const token = 'your-jwt-token';

const employeeAPI = {
  terminate: async (employeeId, reason, terminationDate) => {
    try {
      const response = await axios.post(
        `${API_BASE_URL}/api/v1/employees/terminate`,
        {
          id: employeeId,
          reason: reason,
          termination_date: terminationDate
        },
        {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
      );
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  suspend: async (employeeId, reason, suspensionEndDate) => {
    try {
      const response = await axios.post(
        `${API_BASE_URL}/api/v1/employees/suspend`,
        {
          id: employeeId,
          reason: reason,
          suspension_end_date: suspensionEndDate
        },
        {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
      );
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  archive: async (employeeId, reason) => {
    try {
      const response = await axios.post(
        `${API_BASE_URL}/api/v1/employees/archive`,
        {
          id: employeeId,
          reason: reason
        },
        {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
      );
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  reactivate: async (employeeId) => {
    try {
      const response = await axios.post(
        `${API_BASE_URL}/api/v1/employees/reactivate`,
        {
          id: employeeId
        },
        {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
      );
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  }
};

// Usage
await employeeAPI.suspend(123, "Pending investigation", "2025-01-15T00:00:00Z");
await employeeAPI.reactivate(123);
```

---

## Error Handling

### Common Errors

#### 1. Employee Not Found
```json
{
  "success": false,
  "message": "employee not found"
}
```

#### 2. Employee Already Terminated
```json
{
  "success": false,
  "message": "employee is already terminated"
}
```

#### 3. Employee Already Suspended
```json
{
  "success": false,
  "message": "employee is already suspended"
}
```

#### 4. Cannot Suspend Terminated Employee
```json
{
  "success": false,
  "message": "cannot suspend a terminated employee"
}
```

#### 5. Cannot Reactivate Terminated Employee
```json
{
  "success": false,
  "message": "cannot reactivate a terminated employee"
}
```

#### 6. Employee Already Active
```json
{
  "success": false,
  "message": "employee is already active"
}
```

#### 7. Validation Error
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": "id is required"
}
```

### Error Handling Example

```javascript
const handleEmployeeAction = async (action, employeeId, data) => {
  try {
    const response = await fetch(`${BASE_URL}/api/v1/employees/${action}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        id: employeeId,
        ...data
      })
    });

    const result = await response.json();

    if (!result.success) {
      // Handle specific errors
      switch (result.message) {
        case 'employee not found':
          console.error('Employee does not exist');
          break;
        case 'employee is already terminated':
          console.error('Employee is already terminated');
          break;
        case 'cannot reactivate a terminated employee':
          console.error('Cannot reactivate terminated employees');
          break;
        default:
          console.error(result.message);
      }
      throw new Error(result.message);
    }

    return result.data;
  } catch (error) {
    console.error('Error:', error);
    throw error;
  }
};
```

---

## Quick Reference

| Action | Endpoint | Method | ID in Body | Can Reactivate? |
|--------|----------|--------|------------|-----------------|
| Terminate | `/api/v1/employees/terminate` | POST | ✅ | ❌ |
| Suspend | `/api/v1/employees/suspend` | POST | ✅ | ✅ |
| Archive | `/api/v1/employees/archive` | POST | ✅ | ✅ |
| Reactivate | `/api/v1/employees/reactivate` | POST | ✅ | N/A |

---

## Notes

1. **ID Parameter**: All endpoints use the database `id` (not `employee_id` string) in the request body.

2. **Reasons**: Reasons are automatically appended to the employee's `notes` field for audit purposes.

3. **Status Changes**: 
   - Terminating sets `status = "terminated"` and `is_active = false`
   - Suspending sets `status = "suspended"` and `is_active = false`
   - Archiving sets `status = "archived"` and `is_active = false`
   - Reactivating sets `status = "active"` and `is_active = true`

4. **Termination is Irreversible**: Once an employee is terminated, they cannot be reactivated. Use suspend or archive if you need the ability to reactivate.

5. **Audit Trail**: All actions automatically record the `updated_by` user ID for audit purposes.

---

## Support

For questions or issues, please contact the development team or refer to the main API documentation.

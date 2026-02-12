# Roles & Permissions API Documentation

## Overview

The HRMS system uses a **hybrid role-based access control (RBAC)** model that combines:

1. **Legacy role** (`user_type` column on the `users` table) — backward-compatible single role
2. **RBAC roles** (`user_roles` join table) — allows multiple roles per user with granular permissions

When a user logs in, the response includes **both** a primary `role` (legacy) and a `roles` array (all roles).

---

## Available Roles

| Role Code     | Display Name   | Description                                                                 |
|---------------|----------------|-----------------------------------------------------------------------------|
| `super_admin` | Super Admin    | Full unrestricted system access                                             |
| `admin`       | Administrator  | Full system access with all permissions                                     |
| `hr`          | HR Manager     | Employee, leave, attendance, payroll, and organizational management         |
| `it`          | IT Support     | Helpdesk, asset management, and attendance config                           |
| `manager`     | Manager        | Team oversight with approval access for leave, attendance, and shifts       |
| `employee`    | Employee       | Basic self-service: view own profile, apply leave, view payslips, etc.      |
| `user`        | User (base)    | System base role — automatically assigned to every user                     |

### Role Hierarchy

```
super_admin  →  Full access (everything)
admin        →  Full access (everything)
hr           →  HR management + employee management + organizational
it           →  Helpdesk + assets + attendance config
manager      →  Team approval (leave, attendance, shifts) + read access
employee     →  Self-service (own profile, leave, payslips, attendance)
user         →  Base role (always present for every user)
```

---

## Login Response

### `POST /api/v1/auth/login`

When a user logs in, the response includes both `role` and `roles`:

**Request:**

```json
{
  "email": "john.doe@company.com",
  "password": "password123"
}
```

**Response — Employee with only employee role:**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 5,
      "username": "johndoe",
      "email": "john.doe@company.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user",
      "roles": ["user", "employee"],
      "is_active": true
    },
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "refresh_expires_in": 604800
  }
}
```

**Response — Manager (who is also an employee):**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 8,
      "username": "jane.smith",
      "email": "jane.smith@company.com",
      "first_name": "Jane",
      "last_name": "Smith",
      "role": "user",
      "roles": ["user", "manager", "employee"],
      "is_active": true
    },
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "refresh_expires_in": 604800
  }
}
```

**Response — HR Manager:**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 3,
      "username": "hr.manager",
      "email": "hr@company.com",
      "first_name": "HR",
      "last_name": "Manager",
      "role": "hr",
      "roles": ["hr", "employee", "user"],
      "is_active": true
    },
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "refresh_expires_in": 604800
  }
}
```

**Response — Admin:**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@company.com",
      "first_name": "System",
      "last_name": "Admin",
      "role": "admin",
      "roles": ["admin", "user"],
      "is_active": true
    },
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "refresh_expires_in": 604800
  }
}
```

### How `role` and `roles` Work

| Field   | Type       | Description                                                              |
|---------|------------|--------------------------------------------------------------------------|
| `role`  | `string`   | Primary/legacy role for backward compatibility (first role in array)      |
| `roles` | `string[]` | Full list of all user roles. **Always includes `"user"` as base role.**   |

**Frontend Usage:**

```javascript
// Check if user is an employee
const isEmployee = loginData.roles.includes("employee");

// Check if user is a manager
const isManager = loginData.roles.includes("manager");

// Check if user is HR
const isHR = loginData.roles.includes("hr");

// Check if user is admin
const isAdmin = loginData.roles.includes("admin") || loginData.roles.includes("super_admin");

// A manager may also be an employee
const isManagerAndEmployee = loginData.roles.includes("manager") && loginData.roles.includes("employee");
```

---

## JWT Token Claims

The JWT token includes both `role` and `roles` in the claims:

```json
{
  "user_id": 5,
  "email": "john.doe@company.com",
  "role": "user",
  "roles": ["user", "employee"],
  "exp": 1738972800,
  "iat": 1738886400
}
```

---

## View User API

### Get User Detail

View a user's account data, roles, and status. Does **not** include employee-specific data (use the Employee API for that).

**`GET /api/v1/users/:id`**

**Auth:** HR or Admin

**Response:**

```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": 5,
    "username": "johndoe",
    "email": "john.doe@company.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+255712345678",
    "role": "user",
    "roles": ["user", "employee", "manager"],
    "status": "active",
    "is_active": true,
    "email_verified": false,
    "last_login": "2026-02-05T14:30:00Z",
    "created_at": "2026-01-10T09:00:00Z",
    "updated_at": "2026-02-05T14:30:00Z"
  }
}
```

| Field            | Type       | Description                                               |
|------------------|------------|-----------------------------------------------------------|
| `id`             | `uint`     | User ID                                                   |
| `username`       | `string`   | Login username                                            |
| `email`          | `string`   | Email address                                             |
| `first_name`     | `string`   | First name                                                |
| `last_name`      | `string`   | Last name                                                 |
| `phone_number`   | `string?`  | Phone number (nullable)                                   |
| `role`           | `string`   | Legacy primary role (`admin`, `hr`, `it`, `user`)         |
| `roles`          | `string[]` | All RBAC roles (e.g. `["user", "employee", "manager"]`)   |
| `status`         | `string`   | Account status: `active`, `suspended`, `blocked`          |
| `is_active`      | `bool`     | Whether the account is active                             |
| `email_verified` | `bool`     | Whether email has been verified                           |
| `last_login`     | `string?`  | Last login timestamp (nullable, ISO 8601)                 |
| `created_at`     | `string`   | Account creation timestamp                                |
| `updated_at`     | `string`   | Last update timestamp                                     |

**Error responses:**

| Status | Condition                     |
|--------|-------------------------------|
| 400    | Invalid user ID or not found  |
| 401    | Not authenticated             |
| 403    | Not HR or Admin               |

---

## Role Management APIs (Admin Only)

All role management APIs require **Admin** or **Super Admin** access.

### 1. Assign Role

Replace the user's primary role.

**`POST /api/v1/users/assign-role`**

**Request:**

```json
{
  "user_id": 5,
  "role": "employee",
  "position_id": 3
}
```

| Field         | Type     | Required | Description                                                |
|---------------|----------|----------|------------------------------------------------------------|
| `user_id`     | `uint`   | Yes      | Target user ID                                             |
| `role`        | `string` | Yes      | Role to assign: `super_admin`, `admin`, `hr`, `it`, `manager`, `employee` |
| `position_id` | `uint`   | No       | Optionally update the employee's job position              |

**Response:**

```json
{
  "success": true,
  "message": "Role assigned successfully",
  "data": {
    "user": { "id": 5, "email": "john@company.com", "role": "user" },
    "role": "employee",
    "position_updated": true,
    "rbac_synced": true
  }
}
```

**Notes:**
- `employee` and `manager` roles are stored as `"user"` in the legacy `user_type` column (backward compatibility)
- The actual role is tracked in the RBAC `user_roles` table
- The `roles` array on next login will reflect the new role

---

### 2. Add Role (Multi-Role)

Add an additional role without removing existing ones.

**`POST /api/v1/users/add-role`**

**Request:**

```json
{
  "user_id": 5,
  "role": "manager"
}
```

| Field     | Type     | Required | Description                                                              |
|-----------|----------|----------|--------------------------------------------------------------------------|
| `user_id` | `uint`   | Yes      | Target user ID                                                           |
| `role`    | `string` | Yes      | Role to add: `super_admin`, `admin`, `hr`, `it`, `manager`, `employee`, `user` |

**Response:**

```json
{
  "success": true,
  "message": "Role added successfully",
  "data": {
    "roles": ["user", "employee", "manager"]
  }
}
```

**Use case:** Promoting an employee to manager — they keep `employee` and gain `manager`.

---

### 3. Remove Role

Remove a specific role from a user.

**`POST /api/v1/users/remove-role`**

**Request:**

```json
{
  "user_id": 5,
  "role": "manager"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Role removed successfully",
  "data": {
    "roles": ["user", "employee"]
  }
}
```

**Notes:**
- Cannot remove the last role — user must have at least one role
- The `"user"` base role is always maintained

---

### 4. Set Roles (Replace All)

Replace all roles for a user with a new set.

**`POST /api/v1/users/set-roles`**

**Request:**

```json
{
  "user_id": 5,
  "roles": ["employee", "manager"]
}
```

| Field     | Type       | Required | Description                                                              |
|-----------|------------|----------|--------------------------------------------------------------------------|
| `user_id` | `uint`     | Yes      | Target user ID                                                           |
| `roles`   | `string[]` | Yes      | New roles array: `super_admin`, `admin`, `hr`, `it`, `manager`, `employee`, `user` |

**Response:**

```json
{
  "success": true,
  "message": "User roles updated successfully",
  "data": {
    "roles": ["employee", "manager", "user"]
  }
}
```

---

### 5. Get User Roles

View all roles assigned to a user.

**`GET /api/v1/users/:user_id/roles`**

**Response:**

```json
{
  "success": true,
  "data": {
    "roles": ["user", "employee", "manager"]
  }
}
```

---

### 6. Transfer Role

Transfer a special role from one user to another.

**`POST /api/v1/users/transfer-role`**

**Request:**

```json
{
  "role": "manager",
  "target_user_id": 10,
  "from_user_id": 5
}
```

| Field            | Type     | Required | Description                                                     |
|------------------|----------|----------|-----------------------------------------------------------------|
| `role`           | `string` | Yes      | Role to transfer: `admin`, `hr`, `it`, `manager`               |
| `target_user_id` | `uint`   | Yes      | User who will receive the role                                  |
| `from_user_id`   | `uint`   | No       | User who currently has the role (required when caller is Admin)  |

---

### 7. Create User from Employee

Create a system login account for an enrolled employee.

**`POST /api/v1/users/create-from-employee`**

**Request:**

```json
{
  "employee_id": 15,
  "password": "Welcome@2026",
  "role": "employee",
  "username": "john.doe"
}
```

| Field         | Type     | Required | Description                                                       |
|---------------|----------|----------|-------------------------------------------------------------------|
| `employee_id` | `uint`   | Yes      | Employee record ID to create a user for                           |
| `password`    | `string` | Yes      | Login password (min 6 characters)                                 |
| `role`        | `string` | Yes      | Role: `admin`, `hr`, `it`, `manager`, `employee`, `user`         |
| `username`    | `string` | No       | Username (auto-generated from employee email if not provided)     |

**Notes:**
- Employee's name, email, and phone are automatically copied to the user record
- The RBAC role is assigned immediately
- The user will see their role in the `roles` array on login

---

## Permissions by Role

### Employee Permissions

| Permission        | Description                          |
|-------------------|--------------------------------------|
| `dashboard:view`  | View personal dashboard              |
| `employee:read`   | View own profile                     |
| `leave:read`      | View own leave requests and balance  |
| `leave:create`    | Submit leave requests                |
| `shift:read`      | View own shift schedule              |
| `department:read` | View department info                 |
| `team:read`       | View team info                       |
| `position:read`   | View position info                   |
| `organization:read` | View organization info             |
| `location:read`   | View location info                   |
| `helpdesk:read`   | View own helpdesk tickets            |
| `helpdesk:create` | Create helpdesk tickets              |
| `payroll:read`    | View own payslips                    |
| `asset:read`      | View own assigned assets             |
| `attendance:read` | View own attendance records          |

### Manager Permissions (in addition to Employee)

| Permission             | Description                               |
|------------------------|-------------------------------------------|
| `employee:update`      | Update direct report profiles             |
| `leave:approve`        | Approve/reject team leave requests        |
| `attendance:approve`   | Approve team attendance records           |
| `shift:approve_swaps`  | Approve shift swap requests               |
| `report:view`          | View team reports                         |

### HR Permissions (Full HR Management)

| Permission                    | Description                          |
|-------------------------------|--------------------------------------|
| `dashboard:manage_announcements` | Manage announcements              |
| `employee:create/update/delete`  | Full employee management           |
| `employee:onboard/offboard`     | Onboarding & offboarding           |
| `leave:approve`                 | Approve leave requests              |
| `leave:manage_types/policies/holidays` | Configure leave settings     |
| `attendance:approve`            | Approve attendance                  |
| `shift:create/update/delete`    | Manage shifts                       |
| `shift:manage_rosters`          | Manage shift rosters                |
| `department/team/position:*`    | Full organizational management      |
| `payroll:run/manage_*`          | Payroll management                  |
| `asset:create/update/assign`    | Asset management                    |
| `helpdesk:manage/manage_kb`     | Helpdesk management                 |
| `report:view/export`            | Reports                             |

### IT Permissions

| Permission               | Description                     |
|--------------------------|---------------------------------|
| `asset:create/update/assign` | Full asset management         |
| `helpdesk:manage/manage_kb`  | Full helpdesk management      |
| `attendance:manage_config`   | Attendance device config      |
| `employee:read`             | View employee profiles         |

### Admin / Super Admin Permissions

**All permissions** — full unrestricted access to the entire system.

---

## Middleware Reference (for Frontend)

The backend uses the following middleware to protect API endpoints:

| Middleware             | Roles Allowed                                              | Usage                                      |
|------------------------|------------------------------------------------------------|--------------------------------------------|
| `AuthMiddleware()`     | Any authenticated user                                     | Basic authentication check                 |
| `EmployeeMiddleware()` | `employee`, `manager`, `hr`, `it`, `admin`, `super_admin`, `user` | Employee self-service endpoints     |
| `ManagerMiddleware()`  | `manager`, `hr`, `admin`, `super_admin`                    | Team management & approval endpoints       |
| `HRMiddleware()`       | `hr`, `admin`, `super_admin`                               | HR management endpoints                    |
| `AdminMiddleware()`    | `admin`, `super_admin`                                     | System administration endpoints            |

### API Access by Role

| API Category             | Employee | Manager | HR   | IT   | Admin |
|--------------------------|----------|---------|------|------|-------|
| Self-service (own data)  | Yes      | Yes     | Yes  | Yes  | Yes   |
| Team approvals           | No       | Yes     | Yes  | No   | Yes   |
| Employee management      | No       | No      | Yes  | No   | Yes   |
| Leave configuration      | No       | No      | Yes  | No   | Yes   |
| Payroll management        | No       | No      | Yes  | No   | Yes   |
| Asset management          | No       | No      | Yes  | Yes  | Yes   |
| Helpdesk management       | No       | No      | Yes  | Yes  | Yes   |
| User/Role management      | No       | No      | No   | No   | Yes   |
| System configuration      | No       | No      | No   | No   | Yes   |

---

## Common Scenarios

### Scenario 1: Enrolling a New Employee

1. **HR creates employee** → `POST /api/v1/employees`
2. **HR creates user account** → `POST /api/v1/users/create-from-employee` with `role: "employee"`
3. **Employee logs in** → Gets `roles: ["user", "employee"]`

### Scenario 2: Promoting Employee to Manager

1. **Admin adds manager role** → `POST /api/v1/users/add-role` with `role: "manager"`
2. **Employee now has** → `roles: ["user", "employee", "manager"]`
3. **Manager can now** approve leave, attendance for their team

### Scenario 3: Assigning HR Role

1. **Admin assigns HR role** → `POST /api/v1/users/assign-role` with `role: "hr"`
2. **User now has** → `roles: ["hr", "user"]`
3. **If also employee** → Admin adds employee role: `roles: ["hr", "employee", "user"]`

### Scenario 4: Multi-Role User (HR + Manager)

1. **Admin sets roles** → `POST /api/v1/users/set-roles` with `roles: ["hr", "manager", "employee"]`
2. **User now has** → `roles: ["hr", "manager", "employee", "user"]`
3. **User can** manage HR tasks AND approve team requests

---

## Frontend Integration Guide

### Checking Roles on Login

```javascript
// After login response
const { user } = loginResponse.data;

// Store roles for UI rendering
const userRoles = user.roles; // e.g. ["user", "employee", "manager"]

// Role check helpers
const hasRole = (role) => userRoles.includes(role);

// UI visibility conditions
const showDashboard        = true; // All users
const showSelfService      = hasRole("employee") || hasRole("user");
const showTeamApprovals    = hasRole("manager");
const showHRManagement     = hasRole("hr");
const showAdminPanel       = hasRole("admin") || hasRole("super_admin");
const showITPanel          = hasRole("it");
const showLeaveApproval    = hasRole("manager") || hasRole("hr") || hasRole("admin");
const showPayrollSection   = hasRole("hr") || hasRole("admin");
```

### Navigation Menu by Role

```javascript
const menuItems = [
  // Always visible
  { label: "Dashboard", path: "/dashboard", visible: true },
  
  // Employee self-service
  { label: "My Profile", path: "/profile", visible: hasRole("employee") || hasRole("user") },
  { label: "Apply Leave", path: "/leave/apply", visible: hasRole("employee") || hasRole("user") },
  { label: "My Payslips", path: "/payslips", visible: hasRole("employee") || hasRole("user") },
  { label: "My Attendance", path: "/attendance", visible: hasRole("employee") || hasRole("user") },
  
  // Manager
  { label: "Team Leave Approvals", path: "/leave/approvals", visible: hasRole("manager") },
  { label: "Team Attendance", path: "/attendance/team", visible: hasRole("manager") },
  
  // HR
  { label: "Employee Management", path: "/employees", visible: hasRole("hr") },
  { label: "Leave Management", path: "/leave/manage", visible: hasRole("hr") },
  { label: "Payroll", path: "/payroll", visible: hasRole("hr") },
  
  // Admin
  { label: "User Management", path: "/users", visible: hasRole("admin") || hasRole("super_admin") },
  { label: "System Settings", path: "/settings", visible: hasRole("admin") || hasRole("super_admin") },
];
```

---

## Change Roles API (Checkbox-Style)

This is the **primary API for managing user roles** from the frontend. It works like a checkbox form: send the full desired set of checked roles, and the backend computes the diff.

### `PUT /api/v1/users/change-roles`

**Auth:** Admin or Super Admin only

### How It Works

```
Frontend Checkbox UI:

  [ ] super_admin
  [x] admin          ← checked
  [x] hr             ← checked
  [ ] it
  [x] manager        ← checked
  [x] employee       ← checked
  [x] user           ← always on (auto-added)

    ↓ Submit ↓

PUT /api/v1/users/change-roles
{
  "user_id": 5,
  "roles": ["admin", "hr", "manager", "employee"]
}

    ↓ Backend computes diff ↓

Response:
{
  "previous_roles": ["user", "employee"],
  "current_roles":  ["admin", "hr", "manager", "employee", "user"],
  "added":          ["admin", "hr", "manager"],
  "removed":        [],
  "unchanged":      ["user", "employee"]
}
```

### Request

```json
{
  "user_id": 5,
  "roles": ["employee", "manager"]
}
```

| Field     | Type       | Required | Description                                                                    |
|-----------|------------|----------|--------------------------------------------------------------------------------|
| `user_id` | `uint`     | Yes      | Target user ID                                                                 |
| `roles`   | `string[]` | Yes      | Desired roles (checked items). `"user"` is auto-added if missing.              |

**Valid role values:** `super_admin`, `admin`, `hr`, `it`, `manager`, `employee`, `user`

### Response

```json
{
  "success": true,
  "message": "Roles changed successfully",
  "data": {
    "user_id": 5,
    "email": "john.doe@company.com",
    "username": "johndoe",
    "previous_roles": ["user", "employee"],
    "current_roles": ["user", "employee", "manager"],
    "added": ["manager"],
    "removed": [],
    "unchanged": ["user", "employee"]
  }
}
```

| Field            | Type       | Description                                          |
|------------------|------------|------------------------------------------------------|
| `user_id`        | `uint`     | Target user ID                                       |
| `email`          | `string`   | User email                                           |
| `username`       | `string`   | Username                                             |
| `previous_roles` | `string[]` | Roles before the change                              |
| `current_roles`  | `string[]` | Roles after the change                               |
| `added`          | `string[]` | Roles that were newly added (checked)                |
| `removed`        | `string[]` | Roles that were removed (unchecked)                  |
| `unchanged`      | `string[]` | Roles that stayed the same                           |

### Usage Examples

#### Example 1: Promote Employee to Manager

User currently has: `["user", "employee"]`

```json
// Request — check "employee" and "manager"
PUT /api/v1/users/change-roles
{
  "user_id": 5,
  "roles": ["employee", "manager"]
}

// Response
{
  "previous_roles": ["user", "employee"],
  "current_roles":  ["user", "employee", "manager"],
  "added":          ["manager"],
  "removed":        [],
  "unchanged":      ["user", "employee"]
}
```

#### Example 2: Make Employee an HR Manager

User currently has: `["user", "employee"]`

```json
// Request — check "employee" and "hr"
PUT /api/v1/users/change-roles
{
  "user_id": 5,
  "roles": ["employee", "hr"]
}

// Response
{
  "previous_roles": ["user", "employee"],
  "current_roles":  ["user", "employee", "hr"],
  "added":          ["hr"],
  "removed":        [],
  "unchanged":      ["user", "employee"]
}
```

#### Example 3: Demote Manager back to Employee

User currently has: `["user", "employee", "manager"]`

```json
// Request — only check "employee" (uncheck "manager")
PUT /api/v1/users/change-roles
{
  "user_id": 5,
  "roles": ["employee"]
}

// Response
{
  "previous_roles": ["user", "employee", "manager"],
  "current_roles":  ["user", "employee"],
  "added":          [],
  "removed":        ["manager"],
  "unchanged":      ["user", "employee"]
}
```

#### Example 4: Give Multiple Roles (HR + Manager + Employee)

User currently has: `["user", "employee"]`

```json
// Request — check "employee", "manager", and "hr"
PUT /api/v1/users/change-roles
{
  "user_id": 5,
  "roles": ["employee", "manager", "hr"]
}

// Response
{
  "previous_roles": ["user", "employee"],
  "current_roles":  ["user", "employee", "manager", "hr"],
  "added":          ["manager", "hr"],
  "removed":        [],
  "unchanged":      ["user", "employee"]
}
```

#### Example 5: Replace All Roles Completely

User currently has: `["user", "employee", "manager", "hr"]`

```json
// Request — only check "it" (uncheck everything else)
PUT /api/v1/users/change-roles
{
  "user_id": 5,
  "roles": ["it"]
}

// Response
{
  "previous_roles": ["user", "employee", "manager", "hr"],
  "current_roles":  ["user", "it"],
  "added":          ["it"],
  "removed":        ["employee", "manager", "hr"],
  "unchanged":      ["user"]
}
```

### Frontend Integration (React Example)

```jsx
// Role Checkbox Component
const AVAILABLE_ROLES = [
  { code: "super_admin", label: "Super Admin", description: "Full unrestricted access" },
  { code: "admin",       label: "Administrator", description: "Full system access" },
  { code: "hr",          label: "HR Manager", description: "Employee & leave management" },
  { code: "it",          label: "IT Support", description: "Helpdesk & assets" },
  { code: "manager",     label: "Manager", description: "Team oversight & approvals" },
  { code: "employee",    label: "Employee", description: "Basic self-service" },
];

function ChangeRolesForm({ userId, currentRoles }) {
  const [checkedRoles, setCheckedRoles] = useState(
    currentRoles.filter(r => r !== "user") // "user" is always auto-added
  );

  const handleToggle = (roleCode) => {
    setCheckedRoles(prev =>
      prev.includes(roleCode)
        ? prev.filter(r => r !== roleCode)
        : [...prev, roleCode]
    );
  };

  const handleSubmit = async () => {
    const res = await fetch("/api/v1/users/change-roles", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`,
      },
      body: JSON.stringify({
        user_id: userId,
        roles: checkedRoles,
      }),
    });
    const data = await res.json();
    // data.data.added    → newly added roles
    // data.data.removed  → removed roles
    // data.data.current_roles → final state
  };

  return (
    <div>
      <h3>Change Roles for User #{userId}</h3>
      {AVAILABLE_ROLES.map(role => (
        <label key={role.code}>
          <input
            type="checkbox"
            checked={checkedRoles.includes(role.code)}
            onChange={() => handleToggle(role.code)}
          />
          {role.label} — {role.description}
        </label>
      ))}
      <p><em>"user" role is always included automatically.</em></p>
      <button onClick={handleSubmit}>Save Roles</button>
    </div>
  );
}
```

### Flow Chart

```
┌─────────────────────────────────────────────────────────────┐
│                   CHANGE ROLES FLOW                         │
└─────────────────────────────────────────────────────────────┘

  ┌──────────────────┐
  │  Frontend UI     │
  │                  │
  │  Checkbox Form:  │
  │  [x] employee    │
  │  [x] manager     │
  │  [ ] hr          │
  │  [ ] it          │
  │  [ ] admin       │
  └────────┬─────────┘
           │
           ▼
  PUT /api/v1/users/change-roles
  { "user_id": 5, "roles": ["employee", "manager"] }
           │
           ▼
  ┌────────────────────────────────────────┐
  │  Backend Service                       │
  │                                        │
  │  1. Validate caller is Admin           │
  │  2. Get current roles from DB          │
  │     → ["user", "employee"]             │
  │  3. Compute diff:                      │
  │     desired  = {employee, manager, user}│
  │     current  = {user, employee}        │
  │     added    = [manager]               │
  │     removed  = []                      │
  │     unchanged= [user, employee]        │
  │  4. Apply changes to RBAC table        │
  │  5. Update legacy user_type column     │
  └────────┬───────────────────────────────┘
           │
           ▼
  ┌────────────────────────────────────────┐
  │  Response                              │
  │                                        │
  │  {                                     │
  │    "previous_roles": ["user","employee"]│
  │    "current_roles": ["user","employee",│
  │                       "manager"]       │
  │    "added": ["manager"]                │
  │    "removed": []                       │
  │    "unchanged": ["user","employee"]    │
  │  }                                     │
  └────────────────────────────────────────┘
           │
           ▼
  ┌──────────────────┐
  │  Frontend Shows  │
  │  Success Toast:  │
  │  "Added: manager"│
  │  "Removed: none" │
  └──────────────────┘


  ┌─────────────────────────────────────────────────────────────┐
  │              ROLE ACCESS MATRIX                             │
  ├──────────────┬──────────┬─────────┬─────┬─────┬────────────┤
  │ Feature      │ employee │ manager │ hr  │ it  │ admin      │
  ├──────────────┼──────────┼─────────┼─────┼─────┼────────────┤
  │ Dashboard    │    ✓     │    ✓    │  ✓  │  ✓  │     ✓      │
  │ Own Profile  │    ✓     │    ✓    │  ✓  │  ✓  │     ✓      │
  │ Apply Leave  │    ✓     │    ✓    │  ✓  │  ─  │     ✓      │
  │ View Payslip │    ✓     │    ✓    │  ✓  │  ─  │     ✓      │
  │ Attendance   │    ✓     │    ✓    │  ✓  │  ✓  │     ✓      │
  │ Approve Leave│    ─     │    ✓    │  ✓  │  ─  │     ✓      │
  │ Approve Attn.│    ─     │    ✓    │  ✓  │  ─  │     ✓      │
  │ Manage Emp.  │    ─     │    ─    │  ✓  │  ─  │     ✓      │
  │ Manage Leave │    ─     │    ─    │  ✓  │  ─  │     ✓      │
  │ Run Payroll  │    ─     │    ─    │  ✓  │  ─  │     ✓      │
  │ Manage Assets│    ─     │    ─    │  ✓  │  ✓  │     ✓      │
  │ Helpdesk Mgmt│    ─     │    ─    │  ✓  │  ✓  │     ✓      │
  │ User Mgmt    │    ─     │    ─    │  ─  │  ─  │     ✓      │
  │ System Config│    ─     │    ─    │  ─  │  ─  │     ✓      │
  └──────────────┴──────────┴─────────┴─────┴─────┴────────────┘
```

### Related Endpoints

| Endpoint                               | Method | Description                             | When to Use                           |
|----------------------------------------|--------|-----------------------------------------|---------------------------------------|
| `PUT /api/v1/users/change-roles`       | PUT    | Checkbox-style role change (recommended)| Main role editing UI                  |
| `GET /api/v1/users/:id/roles`          | GET    | Get current roles for a user            | Load checkbox initial state           |
| `POST /api/v1/users/assign-role`       | POST   | Replace primary role                    | Single role assignment                |
| `POST /api/v1/users/add-role`          | POST   | Add one role                            | Quick add without full form           |
| `POST /api/v1/users/remove-role`       | POST   | Remove one role                         | Quick remove without full form        |
| `POST /api/v1/users/set-roles`         | POST   | Replace all roles                       | Programmatic bulk replace             |
| `POST /api/v1/users/transfer-role`     | POST   | Transfer role between users             | Role handover                         |

---

## Changes Summary

### What Changed

1. **New role constants**: `employee` and `manager` are now first-class role constants alongside `admin`, `hr`, `it`, `super_admin`
2. **Login response**: The `roles` array now includes `"employee"` and/or `"manager"` when the user has those RBAC roles
3. **Role assignment APIs**: All role management endpoints (`assign-role`, `add-role`, `set-roles`, `transfer-role`) now accept `"manager"` as a valid role
4. **New Change Roles API**: `PUT /api/v1/users/change-roles` — checkbox-style endpoint to add/remove multiple roles at once with diff response
5. **New middleware**: `ManagerMiddleware()` and `EmployeeMiddleware()` are available for route protection
6. **Legacy compatibility**: The `role` field (legacy `user_type`) maps `employee` and `manager` to `"user"` for backward compatibility; the actual role appears in the `roles` array
7. **RBAC permissions**: Both `employee` and `manager` have dedicated permission sets in the RBAC system

### What Did NOT Change

- All existing API endpoints work exactly as before
- The `"user"` base role is still always included in the `roles` array
- Admin, HR, IT role behavior is unchanged
- Self-service APIs are unaffected
- JWT token structure is the same (just `roles` array may contain new values)

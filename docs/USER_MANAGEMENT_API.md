# User Management & Authentication API

Complete API documentation for authentication, user creation, role management, and account administration.

**Base URL:** `/api/v1`

**How user creation works:** Employees are created first via the HR employee onboarding process. Then, Admin/HR creates a login (user) account for the employee using `POST /api/v1/users`. The user's name and email are pulled from the employee record automatically. The only exception is the initial admin who sets up the system via the onboarding endpoint.

---

## Table of Contents

1. [Authentication](#1-authentication)
   - [Login](#11-login)
   - [Refresh Token](#12-refresh-token)
   - [Get Profile](#13-get-profile)
   - [Logout](#14-logout)
2. [Create Users](#2-create-users)
   - [List Employees Without User Account](#21-list-employees-without-user-account)
   - [Create User from Employee](#22-create-user-from-employee)
   - [Onboarding (First-Time System Setup)](#23-onboarding-first-time-system-setup)
   - [Setup Wizard Status](#24-setup-wizard-status)
3. [User Management](#3-user-management)
   - [List Users](#31-list-users)
   - [List Special Role Users](#32-list-special-role-users)
   - [Get User Roles](#33-get-user-roles)
4. [Role Management](#4-role-management)
   - [Assign Role](#41-assign-role)
   - [Add Role](#42-add-role)
   - [Remove Role](#43-remove-role)
   - [Set Roles](#44-set-roles)
   - [Transfer Role](#45-transfer-role)
5. [Account Administration](#5-account-administration)
   - [Suspend User](#51-suspend-user)
   - [Unsuspend User](#52-unsuspend-user)
   - [Block User](#53-block-user)
   - [Unblock User](#54-unblock-user)
6. [Security](#6-security)
   - [List Blocked Users](#61-list-blocked-users)
   - [Get Rate Limit Status](#62-get-rate-limit-status)
7. [Common Flows](#7-common-flows)
8. [Quick Reference](#8-quick-reference)

---

## 1. Authentication

### 1.1 Login

Authenticate a user and receive access and refresh tokens.

```
POST /api/v1/auth/login
```

**Auth:** None (public)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | User's email address |
| `password` | string | Yes | User's password (min 6 chars) |

**Example Request:**

```json
{
  "email": "james.tech@company.com",
  "password": "SecurePass123"
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 12,
      "username": "jamestech",
      "email": "james.tech@company.com",
      "first_name": "James",
      "last_name": "Techi",
      "role": "it",
      "roles": ["it", "user"],
      "is_active": true
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "refresh_expires_in": 604800
  }
}
```

**Error Responses:**

| Status | Scenario |
|--------|----------|
| 401 | Invalid email or password |
| 401 | Account is suspended or blocked |
| 401 | Account locked due to too many failed attempts |
| 422 | Validation error (missing fields) |

---

### 1.2 Refresh Token

Exchange a valid refresh token for new access and refresh tokens.

```
POST /api/v1/auth/refresh
```

**Auth:** None (uses refresh token in body)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `refresh_token` | string | Yes | Valid refresh token from login |

**Example Request:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Tokens refreshed successfully",
  "data": {
    "user": {
      "id": 12,
      "username": "jamestech",
      "email": "james.tech@company.com",
      "first_name": "James",
      "last_name": "Techi",
      "role": "it",
      "roles": ["it", "user"],
      "is_active": true
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "refresh_expires_in": 604800
  }
}
```

---

### 1.3 Get Profile

Get the authenticated user's profile.

```
GET /api/v1/auth/profile
Authorization: Bearer <access_token>
```

**Auth:** Required (any authenticated user)

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 12,
    "username": "jamestech",
    "email": "james.tech@company.com",
    "first_name": "James",
    "last_name": "Techi",
    "role": "it",
    "roles": ["it", "user"],
    "is_active": true
  }
}
```

---

### 1.4 Logout

Logout the current user. The client must discard the access token after calling this.

```
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
```

**Auth:** Required

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Logged out successfully. Please discard the access token on the client.",
  "data": {
    "logged_out": true
  }
}
```

---

## 2. Create Users

> **Important:** Users are always created from existing employee records. The employee must be onboarded first (via the employee onboarding flow), then Admin/HR creates a login account for them. The only exception is the initial system admin, who is created via the onboarding endpoint.

### 2.1 List Employees Without User Account

Get employees that don't have a login account yet. **Use this to populate the employee picker** when creating a new user on the frontend.

```
GET /api/v1/users/available-employees
Authorization: Bearer <access_token>
```

**Auth:** Admin or HR

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 20 | Items per page (max 100) |
| `search` | string | - | Search by name, employee ID, or email |

**Example Request:**

```
GET /api/v1/users/available-employees?search=james
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Employees without user accounts retrieved successfully",
  "data": [
    {
      "id": 45,
      "employee_id": "EMP045",
      "first_name": "James",
      "last_name": "Techi",
      "email": "james.tech@company.com",
      "phone_number": "+255712345678",
      "department_id": 3,
      "position_id": 7,
      "status": "active"
    },
    {
      "id": 52,
      "employee_id": "EMP052",
      "first_name": "James",
      "last_name": "Mwanga",
      "email": "james.mwanga@company.com",
      "department_id": 2,
      "position_id": 4,
      "status": "active"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

---

### 2.2 Create User from Employee

**This is the primary API to create system users.** It takes an `employee_id`, pulls the employee's name and email automatically, and creates a login account with the specified role.

```
POST /api/v1/users
Authorization: Bearer <access_token>
```

**Auth:** Admin or HR only

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `employee_id` | int | Yes | ID of the existing employee (from the picker above) |
| `password` | string | Yes | Login password (minimum 6 characters) |
| `role` | string | Yes | One of: `admin`, `hr`, `it`, `manager`, `employee`, `user` |
| `username` | string | No | Optional. Auto-generated from employee email if omitted |

**What happens automatically:**
- User's `first_name` and `last_name` are taken from the employee record
- User's `email` is taken from the employee's `work_email` (or `personal_email` if work_email is empty)
- User's `phone_number` is taken from the employee record
- The employee record is linked to the new user (`employee.user_id` = new user ID)
- RBAC role is assigned so permissions work immediately
- The response includes a `credentials` object with `email` and `password` so Admin/HR can share the login details with the employee

**Available Roles:**

| Role | Description |
|------|-------------|
| `admin` | Full system administrator |
| `hr` | HR manager with employee management access |
| `it` | IT support with helpdesk and asset access |
| `manager` | Team manager with approval capabilities |
| `employee` | Regular employee with self-service access |
| `user` | Base-level user |

---

**Example: Create an IT User**

```json
{
  "employee_id": 45,
  "password": "Secure@2026",
  "role": "it"
}
```

**Example Response (201 Created):**

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "user": {
      "id": 12,
      "username": "jamestech",
      "email": "james.tech@company.com",
      "first_name": "James",
      "last_name": "Techi",
      "phone_number": "+255712345678",
      "role": "it",
      "status": "active",
      "is_active": true,
      "created_at": "2026-02-06T15:30:00Z",
      "updated_at": "2026-02-06T15:30:00Z"
    },
    "employee_id": "EMP045",
    "roles": ["it", "user"],
    "credentials": {
      "email": "james.tech@company.com",
      "password": "Secure@2026"
    }
  }
}
```

> **Note:** The `credentials` object contains the login email and password. Share these with the employee so they can login. The password is only returned once at creation time.

---

**Example: Create an HR User**

```json
{
  "employee_id": 30,
  "password": "HrAccess2026!",
  "role": "hr"
}
```

**Example Response (201 Created):**

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "user": {
      "id": 13,
      "username": "sarahm",
      "email": "sarah.m@company.com",
      "first_name": "Sarah",
      "last_name": "Mwangi",
      "role": "hr",
      "status": "active",
      "is_active": true,
      "created_at": "2026-02-06T15:32:00Z",
      "updated_at": "2026-02-06T15:32:00Z"
    },
    "employee_id": "EMP030",
    "roles": ["hr", "user"],
    "credentials": {
      "email": "sarah.m@company.com",
      "password": "HrAccess2026!"
    }
  }
}
```

---

**Example: Create a Manager**

```json
{
  "employee_id": 22,
  "password": "Manager2026!",
  "role": "manager"
}
```

**Example Response (201 Created):**

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "user": {
      "id": 14,
      "username": "mikekamau",
      "email": "mike.kamau@company.com",
      "first_name": "Mike",
      "last_name": "Kamau",
      "role": "user",
      "status": "active",
      "is_active": true,
      "created_at": "2026-02-06T15:34:00Z",
      "updated_at": "2026-02-06T15:34:00Z"
    },
    "employee_id": "EMP022",
    "roles": ["manager", "user"],
    "credentials": {
      "email": "mike.kamau@company.com",
      "password": "Manager2026!"
    }
  }
}
```

---

**Example: Create a Regular Employee User**

```json
{
  "employee_id": 60,
  "password": "Employee2026!",
  "role": "employee"
}
```

**Example Response (201 Created):**

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "user": {
      "id": 15,
      "username": "annajuma",
      "email": "anna.juma@company.com",
      "first_name": "Anna",
      "last_name": "Juma",
      "role": "user",
      "status": "active",
      "is_active": true,
      "created_at": "2026-02-06T15:36:00Z",
      "updated_at": "2026-02-06T15:36:00Z"
    },
    "employee_id": "EMP060",
    "roles": ["employee", "user"],
    "credentials": {
      "email": "anna.juma@company.com",
      "password": "Employee2026!"
    }
  }
}
```

---

**Error Responses:**

| Status | Scenario |
|--------|----------|
| 400 | `"employee not found"` |
| 400 | `"employee is not active"` |
| 400 | `"this employee already has a user account"` |
| 400 | `"employee does not have an email address (work_email or personal_email required)"` |
| 400 | `"a user with email ... already exists"` |
| 400 | `"only an admin or HR can create users"` |
| 401 | Missing or invalid access token |
| 422 | Validation error (missing required fields, invalid role, password too short) |

---

### 2.3 Onboarding (First-Time System Setup)

One-time setup: creates the first admin user and the organization. **Use this only when setting up the system for the first time.** This is the only case where a user is created without an employee record.

```
POST /api/v1/auth/onboard
```

**Auth:** None (public - first-time setup only)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | Admin email address |
| `password` | string | Yes | Minimum 6 characters |
| `first_name` | string | Yes | Minimum 2 characters |
| `last_name` | string | Yes | Minimum 2 characters |
| `username` | string | No | Auto-generated if not provided |
| `name` | string | Yes | Organization/company name (min 2, max 255) |
| `legal_name` | string | No | Legal entity name |
| `registration_number` | string | No | Business registration number |
| `industry` | string | No | Industry sector |
| `company_size` | string | No | `startup`, `small`, `medium`, `large`, or `enterprise` |
| `country` | string | No | Country name |
| `timezone` | string | No | Timezone (e.g., `Africa/Dar_es_Salaam`) |
| `currency_code` | string | No | 3-letter currency code (e.g., `TZS`, `KES`, `USD`) |

**Example Request:**

```json
{
  "email": "admin@mycompany.com",
  "password": "AdminSecure123",
  "first_name": "Jane",
  "last_name": "Smith",
  "name": "My Company Ltd",
  "industry": "Technology",
  "company_size": "small",
  "country": "Tanzania",
  "timezone": "Africa/Dar_es_Salaam",
  "currency_code": "TZS"
}
```

**Example Response (201 Created):**

```json
{
  "success": true,
  "message": "Onboarding completed successfully. Welcome to HRMS!",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@mycompany.com",
      "first_name": "Jane",
      "last_name": "Smith",
      "role": "admin",
      "roles": ["admin"],
      "is_active": true
    },
    "organization": {
      "id": 1,
      "code": "ORG001",
      "name": "My Company Ltd",
      "status": "active"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "setup_complete": false
  }
}
```

---

### 2.4 Setup Wizard Status

Check the progress of the initial setup wizard after onboarding.

```
GET /api/v1/auth/setup-wizard/status
Authorization: Bearer <access_token>
```

**Auth:** Required

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Setup wizard status retrieved successfully",
  "data": {
    "completed": false,
    "current_step": "departments",
    "progress": 40,
    "steps": [
      {
        "step": "organization",
        "title": "Organization Setup",
        "description": "Configure your organization details",
        "completed": true,
        "required": true
      },
      {
        "step": "departments",
        "title": "Departments",
        "description": "Create your departments",
        "completed": false,
        "required": true
      }
    ]
  }
}
```

---

## 3. User Management

### 3.1 List Users

Get a paginated list of all users. Supports filtering by role and search.

```
GET /api/v1/users
Authorization: Bearer <access_token>
```

**Auth:** HR or Admin

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 20 | Items per page (max 100) |
| `role` | string | - | Filter by role: `admin`, `hr`, `it`, `user` |
| `search` | string | - | Search by email, first name, last name, or username |

**Example:**

```
GET /api/v1/users?page=1&page_size=10&role=it
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 12,
      "username": "jamestech",
      "email": "james.tech@company.com",
      "first_name": "James",
      "last_name": "Techi",
      "role": "it",
      "status": "active",
      "is_active": true,
      "created_at": "2026-02-06T15:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 3.2 List Special Role Users

List users who have admin, HR, or IT roles.

```
GET /api/v1/users/special-roles
Authorization: Bearer <access_token>
```

**Auth:** HR or Admin

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 20 | Items per page (max 100) |
| `search` | string | - | Search by email, first name, last name, or username |

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Special role users (IT, HR, Admin) retrieved successfully",
  "data": [
    {
      "id": 1,
      "email": "admin@company.com",
      "first_name": "Jane",
      "last_name": "Smith",
      "role": "admin",
      "status": "active"
    },
    {
      "id": 12,
      "email": "james.tech@company.com",
      "first_name": "James",
      "last_name": "Techi",
      "role": "it",
      "status": "active"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

---

### 3.3 Get User Roles

Get all roles assigned to a specific user.

```
GET /api/v1/users/:id/roles
Authorization: Bearer <access_token>
```

**Auth:** HR or Admin

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "User roles retrieved successfully",
  "data": {
    "user_id": 12,
    "roles": ["it", "user"]
  }
}
```

---

## 4. Role Management

### 4.1 Assign Role

Replace the user's current role with a new one. Useful for promoting an employee.

```
POST /api/v1/users/assign-role
Authorization: Bearer <access_token>
```

**Auth:** Admin only

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | int | Yes | ID of the user |
| `role` | string | Yes | `admin`, `hr`, `it`, or `employee` |
| `position_id` | int | No | Optionally update the user's job position |

**Example Request:**

```json
{
  "user_id": 15,
  "role": "hr",
  "position_id": 3
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Role assigned successfully",
  "data": {
    "user": {
      "id": 15,
      "email": "anna.juma@company.com",
      "first_name": "Anna",
      "last_name": "Juma",
      "role": "hr",
      "is_active": true
    },
    "role": "hr",
    "position_updated": true
  }
}
```

---

### 4.2 Add Role

Add an additional role to a user without removing existing roles. This enables multiple roles per user.

```
POST /api/v1/users/add-role
Authorization: Bearer <access_token>
```

**Auth:** Admin only

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | int | Yes | ID of the user |
| `role` | string | Yes | `admin`, `hr`, `it`, `employee`, or `user` |

**Example Request:**

```json
{
  "user_id": 12,
  "role": "hr"
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Role added successfully",
  "data": {
    "user_id": 12,
    "roles": ["it", "hr", "user"]
  }
}
```

---

### 4.3 Remove Role

Remove a role from a user. Cannot remove the last role.

```
POST /api/v1/users/remove-role
Authorization: Bearer <access_token>
```

**Auth:** Admin only

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | int | Yes | ID of the user |
| `role` | string | Yes | Role to remove |

**Example Request:**

```json
{
  "user_id": 12,
  "role": "hr"
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Role removed successfully",
  "data": {
    "user_id": 12,
    "roles": ["it", "user"]
  }
}
```

---

### 4.4 Set Roles

Replace all roles for a user with a new set. At least one role must be provided.

```
POST /api/v1/users/set-roles
Authorization: Bearer <access_token>
```

**Auth:** Admin only

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | int | Yes | ID of the user |
| `roles` | string[] | Yes | Array of roles (min 1) |

**Example Request:**

```json
{
  "user_id": 12,
  "roles": ["it", "manager"]
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Roles updated successfully",
  "data": {
    "user_id": 12,
    "roles": ["it", "manager"]
  }
}
```

---

### 4.5 Transfer Role

Transfer a special role (admin, hr, it) from one user to another. The source user loses the role.

```
POST /api/v1/users/transfer-role
Authorization: Bearer <access_token>
```

**Auth:** Admin (can transfer any role) or role holder (can transfer own role)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `role` | string | Yes | `admin`, `hr`, or `it` |
| `target_user_id` | int | Yes | User who will receive the role |
| `from_user_id` | int | Admin only | User who currently has the role |

**Example Request:**

```json
{
  "role": "it",
  "target_user_id": 20,
  "from_user_id": 12
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Role transferred successfully",
  "data": {
    "from_user": {
      "id": 12,
      "email": "james.tech@company.com",
      "role": "user"
    },
    "target_user": {
      "id": 20,
      "email": "new.it@company.com",
      "role": "it"
    },
    "role": "it"
  }
}
```

---

## 5. Account Administration

### 5.1 Suspend User

```
POST /api/v1/users/:id/suspend
Authorization: Bearer <access_token>
```

**Auth:** HR or Admin. User will not be able to login until unsuspended.

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "User suspended successfully",
  "data": { "id": 15, "email": "anna.juma@company.com", "status": "suspended" }
}
```

### 5.2 Unsuspend User

```
POST /api/v1/users/:id/unsuspend
Authorization: Bearer <access_token>
```

**Auth:** HR or Admin

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "User unsuspended successfully",
  "data": { "id": 15, "email": "anna.juma@company.com", "status": "active" }
}
```

### 5.3 Block User

```
POST /api/v1/users/:id/block
Authorization: Bearer <access_token>
```

**Auth:** HR or Admin. User will not be able to login until unblocked.

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "User blocked successfully",
  "data": { "id": 15, "email": "anna.juma@company.com", "status": "blocked" }
}
```

### 5.4 Unblock User

```
POST /api/v1/users/:id/unblock
Authorization: Bearer <access_token>
```

**Auth:** HR or Admin

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "User unblocked successfully",
  "data": { "id": 15, "email": "anna.juma@company.com", "status": "active" }
}
```

---

## 6. Security

### 6.1 List Blocked Users

```
GET /api/v1/security/blocked-users
Authorization: Bearer <access_token>
```

**Auth:** Admin only

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 20 | Items per page (max 100) |
| `include_suspended` | bool | false | Also include suspended users |

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Blocked users retrieved successfully",
  "data": [
    {
      "id": 8,
      "email": "blocked.user@company.com",
      "firstName": "Bad",
      "lastName": "Actor",
      "status": "blocked",
      "failedLoginCount": 5,
      "lockedUntil": null
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 1, "total_pages": 1 }
}
```

### 6.2 Get Rate Limit Status

```
GET /api/v1/security/rate-limit/status?email=user@company.com
Authorization: Bearer <access_token>
```

**Auth:** Admin only

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Rate limit status retrieved",
  "data": {
    "email": "user@company.com",
    "blocked": false,
    "suspended": false,
    "failedLoginCount": 3,
    "maxAttempts": 5,
    "attemptsRemaining": 2,
    "canLogin": true
  }
}
```

---

## 7. Common Flows

### Flow 1: First-Time System Setup

```
1. POST /api/v1/auth/onboard            → Creates admin user + organization (no employee needed)
2. GET  /api/v1/auth/setup-wizard/status → Check wizard progress
3. Complete setup: departments, positions, locations, etc.
```

### Flow 2: Create a User for an Employee (Main Flow)

```
Step 1: HR onboards the employee via the employee onboarding flow
        POST /api/v1/employees/onboarding/draft
        POST /api/v1/employees/onboarding/:employee_id/step/:step
        POST /api/v1/employees/onboarding/:employee_id/complete

Step 2: Admin/HR lists employees without user accounts
        GET /api/v1/users/available-employees

Step 3: Admin/HR creates the user account for the employee
        POST /api/v1/users
        {
          "employee_id": 45,
          "password": "Secure2026!",
          "role": "it"
        }

→ Employee #45 now has a login account with role "it"
→ They can login with their work email and the password set above
```

### Flow 3: Create IT, HR, Manager Users

```
# Create IT user
POST /api/v1/users → { "employee_id": 45, "password": "...", "role": "it" }

# Create HR user
POST /api/v1/users → { "employee_id": 30, "password": "...", "role": "hr" }

# Create Manager user
POST /api/v1/users → { "employee_id": 22, "password": "...", "role": "manager" }

# Create regular employee user (self-service access)
POST /api/v1/users → { "employee_id": 60, "password": "...", "role": "employee" }
```

### Flow 4: Add Multiple Roles to an Existing User

```
POST /api/v1/users/add-role
{
  "user_id": 12,
  "role": "hr"
}
→ User 12 now has roles: ["it", "hr", "user"]
```

### Flow 5: Login and Load Permissions (Frontend)

```
1. POST /api/v1/auth/login              → Get access_token, user roles
2. GET  /api/v1/features/my-access      → Get feature permissions
3. Store roles + permissions in client state
4. Render UI based on user's role and permissions
```

### Flow 6: Handle Compromised Account

```
1. POST /api/v1/users/:id/block         → Block user immediately
2. (Investigate)
3. POST /api/v1/users/:id/unblock       → Restore access when resolved
```

---

## 8. Quick Reference

### All Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/auth/login` | Public | Login |
| `POST` | `/api/v1/auth/refresh` | Public | Refresh tokens |
| `GET` | `/api/v1/auth/profile` | Any | Get my profile |
| `POST` | `/api/v1/auth/logout` | Any | Logout |
| `POST` | `/api/v1/auth/onboard` | Public | First-time system setup (creates admin) |
| `GET` | `/api/v1/auth/setup-wizard/status` | Any | Setup wizard progress |
| **`GET`** | **`/api/v1/users/available-employees`** | **Admin/HR** | **List employees without user accounts (picker)** |
| **`POST`** | **`/api/v1/users`** | **Admin/HR** | **Create user from employee record** |
| `GET` | `/api/v1/users` | HR/Admin | List users |
| `GET` | `/api/v1/users/special-roles` | HR/Admin | List admin/hr/it users |
| `GET` | `/api/v1/users/:id/roles` | HR/Admin | Get user's roles |
| `POST` | `/api/v1/users/assign-role` | Admin | Replace user's role |
| `POST` | `/api/v1/users/add-role` | Admin | Add role to user |
| `POST` | `/api/v1/users/remove-role` | Admin | Remove role from user |
| `POST` | `/api/v1/users/set-roles` | Admin | Replace all user roles |
| `POST` | `/api/v1/users/transfer-role` | Admin/Self | Transfer role |
| `POST` | `/api/v1/users/:id/suspend` | HR/Admin | Suspend user |
| `POST` | `/api/v1/users/:id/unsuspend` | HR/Admin | Unsuspend user |
| `POST` | `/api/v1/users/:id/block` | HR/Admin | Block user |
| `POST` | `/api/v1/users/:id/unblock` | HR/Admin | Unblock user |
| `GET` | `/api/v1/security/blocked-users` | Admin | List blocked users |
| `GET` | `/api/v1/security/rate-limit/status` | Admin | Rate limit status |

### Available Roles

| Role | Code | Description |
|------|------|-------------|
| Administrator | `admin` | Full system access |
| HR Manager | `hr` | HR and employee management |
| IT Support | `it` | Helpdesk, assets, and technical support |
| Manager | `manager` | Team oversight and approvals |
| Employee | `employee` | Self-service features |
| User | `user` | Base role (always present) |

### User Statuses

| Status | Can Login | Description |
|--------|-----------|-------------|
| `active` | Yes | Normal active account |
| `suspended` | No | Temporarily suspended by HR/Admin |
| `blocked` | No | Blocked due to security or manual action |

### Error Codes

| Status | Description |
|--------|-------------|
| 400 | Bad request (invalid data or business rule violation) |
| 401 | Unauthorized (missing or invalid token) |
| 403 | Forbidden (insufficient role/permissions) |
| 422 | Validation error (missing required fields) |
| 500 | Internal server error |

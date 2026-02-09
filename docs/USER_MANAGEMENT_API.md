# User Management & Authentication API

Complete API documentation for authentication, user creation, role management, and account administration.

---

## Table of Contents

1. [Authentication](#1-authentication)
   - [Register](#11-register)
   - [Login](#12-login)
   - [Refresh Token](#13-refresh-token)
   - [Get Profile](#14-get-profile)
   - [Logout](#15-logout)
2. [Onboarding (First-Time Setup)](#2-onboarding)
   - [Complete Onboarding](#21-complete-onboarding)
   - [Get Setup Wizard Status](#22-get-setup-wizard-status)
3. [User Management](#3-user-management)
   - [List Users](#31-list-users)
   - [List Special Role Users](#32-list-special-role-users)
   - [Assign Role](#33-assign-role)
   - [Add Role](#34-add-role)
   - [Remove Role](#35-remove-role)
   - [Set Roles](#36-set-roles)
   - [Get User Roles](#37-get-user-roles)
   - [Transfer Role](#38-transfer-role)
4. [Account Administration](#4-account-administration)
   - [Suspend User](#41-suspend-user)
   - [Unsuspend User](#42-unsuspend-user)
   - [Block User](#43-block-user)
   - [Unblock User](#44-unblock-user)
5. [Security](#5-security)
   - [List Blocked Users](#51-list-blocked-users)
   - [Get Rate Limit Status](#52-get-rate-limit-status)
6. [Common Flows](#6-common-flows)

---

## 1. Authentication

### 1.1 Register

Create a new user account.

```
POST /api/v1/auth/register
```

**Auth:** None (public)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | Valid email address |
| `password` | string | Yes | Minimum 6 characters |
| `first_name` | string | Yes | Minimum 2 characters |
| `last_name` | string | Yes | Minimum 2 characters |
| `username` | string | No | Minimum 3 characters. Auto-generated from email if not provided |
| `role` | string | No | `"user"` (default), `"admin"`, or `"hr"` |

**Example Request:**

```json
{
  "email": "john.doe@company.com",
  "password": "SecurePass123",
  "first_name": "John",
  "last_name": "Doe",
  "role": "user"
}
```

**Example Response (201 Created):**

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 5,
      "username": "john.doe",
      "email": "john.doe@company.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user",
      "roles": ["user"],
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

### 1.2 Login

Authenticate a user and receive access tokens.

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
  "email": "john.doe@company.com",
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
      "id": 5,
      "username": "john.doe",
      "email": "john.doe@company.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "hr",
      "roles": ["hr", "user"],
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

### 1.3 Refresh Token

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
      "id": 5,
      "username": "john.doe",
      "email": "john.doe@company.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "hr",
      "roles": ["hr", "user"],
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

### 1.4 Get Profile

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
    "id": 5,
    "username": "john.doe",
    "email": "john.doe@company.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "hr",
    "roles": ["hr", "user"],
    "is_active": true
  }
}
```

---

### 1.5 Logout

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

## 2. Onboarding

### 2.1 Complete Onboarding

First-time setup: creates the admin user, tenant, and organization in one step.

```
POST /api/v1/auth/onboard
```

**Auth:** None (public - used for initial setup)

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
| `timezone` | string | No | Timezone (e.g., `Africa/Nairobi`) |
| `currency_code` | string | No | 3-letter currency code (e.g., `KES`, `USD`) |

**Example Request:**

```json
{
  "email": "admin@mycompany.com",
  "password": "AdminSecure123",
  "first_name": "Jane",
  "last_name": "Smith",
  "name": "My Company Ltd",
  "legal_name": "My Company Limited",
  "registration_number": "BRN-2026-001",
  "industry": "Technology",
  "company_size": "small",
  "country": "Kenya",
  "timezone": "Africa/Nairobi",
  "currency_code": "KES"
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
    "tenant": {
      "id": 1,
      "name": "my_company_ltd",
      "domain": null,
      "status": "active"
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

### 2.2 Get Setup Wizard Status

Check the progress of the initial setup wizard.

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
      },
      {
        "step": "positions",
        "title": "Job Positions",
        "description": "Define job positions",
        "completed": false,
        "required": false
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

**Example Request:**

```
GET /api/v1/users?page=1&page_size=10&role=hr&search=john
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 5,
      "username": "john.doe",
      "email": "john.doe@company.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "hr",
      "status": "active",
      "is_active": true,
      "created_at": "2026-02-06T10:00:00Z"
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
      "id": 3,
      "email": "hr@company.com",
      "first_name": "Bob",
      "last_name": "Jones",
      "role": "hr",
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

### 3.3 Assign Role

Assign a role (admin, hr, it, employee) to a normal user. The target user must currently have the role `"user"`.

```
POST /api/v1/users/assign-role
Authorization: Bearer <access_token>
```

**Auth:** Admin only

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | int | Yes | ID of the user to assign the role to |
| `role` | string | Yes | `admin`, `hr`, `it`, or `employee` |
| `position_id` | int | No | Optionally update the user's job position |

**Example Request:**

```json
{
  "user_id": 5,
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
      "id": 5,
      "username": "john.doe",
      "email": "john.doe@company.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "hr",
      "roles": ["hr"],
      "is_active": true
    },
    "role": "hr",
    "position_updated": true
  }
}
```

---

### 3.4 Add Role

Add an additional role to a user without removing existing roles. This enables multiple roles per user (e.g., `["hr", "it"]`).

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
  "user_id": 5,
  "role": "it"
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Role added successfully",
  "data": {
    "user_id": 5,
    "roles": ["hr", "it", "user"]
  }
}
```

---

### 3.5 Remove Role

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
  "user_id": 5,
  "role": "it"
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Role removed successfully",
  "data": {
    "user_id": 5,
    "roles": ["hr", "user"]
  }
}
```

---

### 3.6 Set Roles

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
  "user_id": 5,
  "roles": ["hr", "employee"]
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Roles updated successfully",
  "data": {
    "user_id": 5,
    "roles": ["hr", "employee"]
  }
}
```

---

### 3.7 Get User Roles

Get all roles assigned to a user.

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
    "user_id": 5,
    "roles": ["hr", "user"]
  }
}
```

---

### 3.8 Transfer Role

Transfer a special role (admin, hr, it) from one user to another. The source user loses the role and the target user gains it.

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
| `from_user_id` | int | Admin only | User who currently has the role (ignored for non-admins) |

**Example Request (Admin transferring HR role):**

```json
{
  "role": "hr",
  "target_user_id": 10,
  "from_user_id": 5
}
```

**Example Request (HR user transferring their own role):**

```json
{
  "role": "hr",
  "target_user_id": 10
}
```

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Role transferred successfully",
  "data": {
    "from_user": {
      "id": 5,
      "email": "old.hr@company.com",
      "role": "user"
    },
    "target_user": {
      "id": 10,
      "email": "new.hr@company.com",
      "role": "hr"
    },
    "role": "hr"
  }
}
```

---

## 4. Account Administration

### 4.1 Suspend User

Suspend a user account. The user will not be able to login until unsuspended.

```
POST /api/v1/users/:id/suspend
Authorization: Bearer <access_token>
```

**Auth:** HR or Admin

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "User suspended successfully",
  "data": {
    "id": 5,
    "email": "john.doe@company.com",
    "status": "suspended",
    "is_active": true
  }
}
```

---

### 4.2 Unsuspend User

Restore a suspended user's account to active.

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
  "data": {
    "id": 5,
    "email": "john.doe@company.com",
    "status": "active",
    "is_active": true
  }
}
```

---

### 4.3 Block User

Block a user account. The user will not be able to login until unblocked.

```
POST /api/v1/users/:id/block
Authorization: Bearer <access_token>
```

**Auth:** HR or Admin

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "User blocked successfully",
  "data": {
    "id": 5,
    "email": "john.doe@company.com",
    "status": "blocked",
    "is_active": true
  }
}
```

---

### 4.4 Unblock User

Restore a blocked user's account to active.

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
  "data": {
    "id": 5,
    "email": "john.doe@company.com",
    "status": "active",
    "is_active": true
  }
}
```

---

## 5. Security

### 5.1 List Blocked Users

List users who are blocked from logging in (rate-limited or manually blocked).

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
      "username": "bad.actor",
      "status": "blocked",
      "failedLoginCount": 5,
      "lockedUntil": null,
      "updatedAt": "2026-02-06T12:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 5.2 Get Rate Limit Status

Check login rate-limit status for a specific email address.

```
GET /api/v1/security/rate-limit/status?email=user@company.com
Authorization: Bearer <access_token>
```

**Auth:** Admin only

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `email` | string | Yes | User's email address |

**Example Response (200 OK):**

```json
{
  "success": true,
  "message": "Rate limit status retrieved",
  "data": {
    "email": "user@company.com",
    "failed_attempts": 3,
    "max_attempts": 5,
    "is_locked": false,
    "locked_until": null,
    "remaining_attempts": 2
  }
}
```

---

## 6. Common Flows

### Flow 1: Initial System Setup

```
1. POST /api/v1/auth/onboard          → Creates admin + organization
2. GET  /api/v1/auth/setup-wizard/status → Check wizard progress
3. (Complete wizard steps: departments, positions, locations, etc.)
```

### Flow 2: Add a New System User (Admin/HR/IT)

```
1. POST /api/v1/auth/register         → Create user account (role: "user")
2. POST /api/v1/users/assign-role     → Promote to admin/hr/it (Admin required)
```

### Flow 3: Add Multiple Roles to Existing User

```
1. POST /api/v1/users/add-role        → Add "hr" role
2. POST /api/v1/users/add-role        → Add "it" role
   → User now has roles: ["hr", "it", "user"]
```

### Flow 4: Login + Load Permissions (Frontend)

```
1. POST /api/v1/auth/login            → Get access_token
2. GET  /api/v1/features/my-access    → Get feature permissions (see Feature Access Control API docs)
3. Store permissions in client state
4. Render UI based on permissions
```

### Flow 5: Handle Compromised Account

```
1. POST /api/v1/users/:id/block       → Block user immediately (HR/Admin)
2. (Investigate the issue)
3. POST /api/v1/users/:id/unblock     → Restore access when resolved
```

### Flow 6: Transfer HR Role to New Person

```
1. POST /api/v1/users/transfer-role   → Transfer "hr" from user A to user B
   → User A becomes "user", User B becomes "hr"
```

---

## Quick Reference

### All Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/auth/register` | Public | Register new user |
| `POST` | `/api/v1/auth/login` | Public | Login |
| `POST` | `/api/v1/auth/refresh` | Public | Refresh tokens |
| `GET` | `/api/v1/auth/profile` | Any | Get my profile |
| `POST` | `/api/v1/auth/logout` | Any | Logout |
| `POST` | `/api/v1/auth/onboard` | Public | First-time setup |
| `GET` | `/api/v1/auth/setup-wizard/status` | Any | Setup wizard progress |
| `GET` | `/api/v1/users` | HR/Admin | List users |
| `GET` | `/api/v1/users/special-roles` | HR/Admin | List admin/hr/it users |
| `POST` | `/api/v1/users/assign-role` | Admin | Assign role to user |
| `POST` | `/api/v1/users/add-role` | Admin | Add role to user |
| `POST` | `/api/v1/users/remove-role` | Admin | Remove role from user |
| `POST` | `/api/v1/users/set-roles` | Admin | Replace all user roles |
| `GET` | `/api/v1/users/:id/roles` | HR/Admin | Get user's roles |
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
| IT Support | `it` | Helpdesk and assets |
| Manager | `manager` | Team oversight and approvals |
| Employee | `employee` | Self-service features |
| User | `user` | Base role (all users have this) |

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
| 404 | Not found |
| 422 | Validation error (missing required fields) |
| 500 | Internal server error |

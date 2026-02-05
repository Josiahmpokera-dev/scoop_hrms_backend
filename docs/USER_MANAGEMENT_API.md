# User Management API

**Base URL:** `BASE_URL` (e.g. `http://localhost:8080/api/v1`)

This document covers all user management endpoints including listing users, role management (including **multi-role support**), and user status management (suspend/block/unblock).

---

## Overview

### Multi-Role Support

Users can have **multiple roles** simultaneously. For example, a user can be both `["admin", "hr"]` or `["employee", "it"]`. The system supports:

- **Legacy role:** Single role stored in `user_type` column (for backward compatibility)
- **RBAC roles:** Multiple roles stored in `user_roles` join table

When a user logs in, all their roles (legacy + RBAC, deduplicated) are returned in the `roles` array in the auth response and embedded in the JWT token.

### Available Roles

| Role | Description |
|------|-------------|
| `admin` | System administrator with full access |
| `hr` | Human Resources - manages employees, leave, etc. |
| `it` | IT department - technical support access |
| `employee` | Regular employee |
| `user` | Basic user (legacy, equivalent to employee) |

### Authentication

All endpoints require Bearer token authentication. Role-specific access is noted per endpoint.

---

## 1. List Users

**`GET BASE_URL/users`**

**Auth:** Bearer token  
**Authorization:** HR or Admin

Returns a paginated list of users with optional filters.

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 20 | Items per page (max 100) |
| `role` | string | - | Filter by role: `admin`, `hr`, `it`, `user` |
| `search` | string | - | Search by email, first name, last name, username |

### Example Request

```http
GET BASE_URL/users?page=1&page_size=20&role=user&search=john
Authorization: Bearer <token>
```

### Example Response (200 OK)

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "email": "john@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "username": "johndoe",
      "role": "user",
      "status": "active"
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

## 2. List Special Role Users (Admin/HR/IT)

**`GET BASE_URL/users/special-roles`**

**Auth:** Bearer token  
**Authorization:** HR or Admin

Returns users who have admin, HR, or IT roles. Useful for role transfer operations.

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 20 | Items per page (max 100) |
| `search` | string | - | Search by email, first name, last name, username |

### Example Request

```http
GET BASE_URL/users/special-roles?page=1&page_size=20
Authorization: Bearer <token>
```

---

## 3. Get User Roles

**`GET BASE_URL/users/:id/roles`**

**Auth:** Bearer token  
**Authorization:** HR or Admin

Returns all roles assigned to a user (combined legacy + RBAC roles, deduplicated).

### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | uint | User ID |

### Example Request

```http
GET BASE_URL/users/42/roles
Authorization: Bearer <token>
```

### Example Response (200 OK)

```json
{
  "success": true,
  "message": "User roles retrieved successfully",
  "data": {
    "user_id": 42,
    "roles": ["admin", "hr"]
  }
}
```

---

## 4. Add Role to User (Multi-Role)

**`POST BASE_URL/users/add-role`**

**Auth:** Bearer token  
**Authorization:** Admin only

Adds an additional role to a user **without removing existing roles**. This enables multi-role assignment (e.g., a user can be both admin and HR).

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | uint | Yes | Target user ID |
| `role` | string | Yes | Role to add: `admin`, `hr`, `it`, `employee`, `user` |

### Example Request

```http
POST BASE_URL/users/add-role
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "user_id": 42,
  "role": "hr"
}
```

### Example Response (200 OK)

```json
{
  "success": true,
  "message": "Role added successfully",
  "data": {
    "user_id": 42,
    "roles": ["admin", "hr"]
  }
}
```

### Error Responses

| Status | Message | Cause |
|--------|---------|-------|
| 400 | "user already has role hr" | Role already assigned |
| 400 | "role xyz not found in system" | Invalid role name |
| 400 | "only an admin can add roles to users" | Caller is not admin |

---

## 5. Remove Role from User

**`POST BASE_URL/users/remove-role`**

**Auth:** Bearer token  
**Authorization:** Admin only

Removes a role from a user. **Cannot remove the last role** - every user must have at least one role.

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | uint | Yes | Target user ID |
| `role` | string | Yes | Role to remove |

### Example Request

```http
POST BASE_URL/users/remove-role
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "user_id": 42,
  "role": "hr"
}
```

### Example Response (200 OK)

```json
{
  "success": true,
  "message": "Role removed successfully",
  "data": {
    "user_id": 42,
    "roles": ["admin"]
  }
}
```

### Error Responses

| Status | Message | Cause |
|--------|---------|-------|
| 400 | "cannot remove the last role - user must have at least one role" | User has only one role |
| 400 | "role xyz not found" | Invalid role name |

---

## 6. Set User Roles (Replace All)

**`POST BASE_URL/users/set-roles`**

**Auth:** Bearer token  
**Authorization:** Admin only

Replaces **all roles** for a user with the specified roles array. At least one role must be provided.

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | uint | Yes | Target user ID |
| `roles` | string[] | Yes | New roles array (minimum 1 role) |

### Example Request

```http
POST BASE_URL/users/set-roles
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "user_id": 42,
  "roles": ["hr", "it"]
}
```

### Example Response (200 OK)

```json
{
  "success": true,
  "message": "Roles updated successfully",
  "data": {
    "user_id": 42,
    "roles": ["hr", "it"]
  }
}
```

### Error Responses

| Status | Message | Cause |
|--------|---------|-------|
| 400 | "at least one role must be specified" | Empty roles array |
| 400 | "invalid role: xyz" | Invalid role in array |

---

## 7. Assign Role (Legacy - Single Role)

**`POST BASE_URL/users/assign-role`**

**Auth:** Bearer token  
**Authorization:** Admin only

Assigns a special role to a normal user. This is the legacy single-role assignment endpoint. For multi-role management, use the **add-role**, **remove-role**, or **set-roles** endpoints instead.

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | uint | Yes | Target user ID (must currently have role "user") |
| `role` | string | Yes | Role to assign: `admin`, `hr`, `it`, `employee` |
| `position_id` | uint | No | Optional: set employee's job position |

### Example Request

```http
POST BASE_URL/users/assign-role
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "user_id": 42,
  "role": "hr",
  "position_id": 5
}
```

---

## 8. Transfer Role

**`POST BASE_URL/users/transfer-role`**

**Auth:** Bearer token  
**Authorization:** Admin (can transfer any role) or Role Holder (can transfer own role)

Transfers a special role (admin, HR, IT) from one user to another.

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `role` | string | Yes | Role to transfer: `admin`, `hr`, `it` |
| `target_user_id` | uint | Yes | User who will receive the role |
| `from_user_id` | uint | Conditional | Required when caller is Admin; the user who currently has the role |

### Example Request (Admin transferring)

```http
POST BASE_URL/users/transfer-role
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "role": "hr",
  "target_user_id": 50,
  "from_user_id": 42
}
```

### Example Request (Role holder transferring own role)

```http
POST BASE_URL/users/transfer-role
Authorization: Bearer <hr_token>
Content-Type: application/json

{
  "role": "hr",
  "target_user_id": 50
}
```

---

## 9. Suspend User

**`POST BASE_URL/users/:id/suspend`**

**Auth:** Bearer token  
**Authorization:** HR or Admin

Sets user status to `suspended`. User cannot login until unsuspended.

### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | uint | User ID to suspend |

### Example Request

```http
POST BASE_URL/users/42/suspend
Authorization: Bearer <token>
```

### Example Response (200 OK)

```json
{
  "success": true,
  "message": "User suspended successfully",
  "data": {
    "id": 42,
    "email": "user@example.com",
    "status": "suspended"
  }
}
```

---

## 10. Unsuspend User

**`POST BASE_URL/users/:id/unsuspend`**

**Auth:** Bearer token  
**Authorization:** HR or Admin

Restores user status to `active` (for suspended users).

### Example Request

```http
POST BASE_URL/users/42/unsuspend
Authorization: Bearer <token>
```

---

## 11. Block User

**`POST BASE_URL/users/:id/block`**

**Auth:** Bearer token  
**Authorization:** HR or Admin

Sets user status to `blocked`. User cannot login until unblocked.

### Example Request

```http
POST BASE_URL/users/42/block
Authorization: Bearer <token>
```

---

## 12. Unblock User

**`POST BASE_URL/users/:id/unblock`**

**Auth:** Bearer token  
**Authorization:** HR or Admin

Restores user status to `active`. Also clears failed login count and lockout.

### Example Request

```http
POST BASE_URL/users/42/unblock
Authorization: Bearer <token>
```

---

## Multi-Role Usage Examples

### Example 1: Make a user both Admin and HR

```http
# User 42 is currently just an admin
# Add HR role
POST BASE_URL/users/add-role
{
  "user_id": 42,
  "role": "hr"
}

# Result: user now has roles ["admin", "hr"]
```

### Example 2: Set multiple roles at once

```http
# Replace all roles with hr and it
POST BASE_URL/users/set-roles
{
  "user_id": 42,
  "roles": ["hr", "it"]
}

# Result: user now has exactly ["hr", "it"]
```

### Example 3: Remove a role

```http
# User 42 has ["admin", "hr", "it"]
# Remove the it role
POST BASE_URL/users/remove-role
{
  "user_id": 42,
  "role": "it"
}

# Result: user now has ["admin", "hr"]
```

### Example 4: Check user's current roles

```http
GET BASE_URL/users/42/roles

# Response:
{
  "success": true,
  "message": "User roles retrieved successfully",
  "data": {
    "user_id": 42,
    "roles": ["admin", "hr"]
  }
}
```

---

## How Roles Work at Login

When a user logs in:

1. The system retrieves the legacy `user_type` role
2. The system retrieves all RBAC roles from `user_roles` table
3. Both are combined and deduplicated into a `roles` array
4. The `roles` array is:
   - Included in the login response (`UserInfo.roles`)
   - Embedded in the JWT token claims
5. Middleware checks the `roles` array when authorizing access to protected routes

### Login Response Example

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG...",
    "user": {
      "id": 42,
      "email": "admin@example.com",
      "firstName": "John",
      "lastName": "Admin",
      "role": "admin",
      "roles": ["admin", "hr"],
      "status": "active"
    }
  }
}
```

Note: `role` (singular) is the primary/legacy role for backward compatibility. `roles` (array) contains all assigned roles.

---

## API Endpoints Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/users` | HR/Admin | List users with filters |
| GET | `/users/special-roles` | HR/Admin | List admin/hr/it users |
| GET | `/users/:id/roles` | HR/Admin | Get all roles for a user |
| POST | `/users/add-role` | Admin | Add role to user (multi-role) |
| POST | `/users/remove-role` | Admin | Remove role from user |
| POST | `/users/set-roles` | Admin | Replace all user roles |
| POST | `/users/assign-role` | Admin | Legacy single role assignment |
| POST | `/users/transfer-role` | Admin/Holder | Transfer role to another user |
| POST | `/users/:id/suspend` | HR/Admin | Suspend user |
| POST | `/users/:id/unsuspend` | HR/Admin | Unsuspend user |
| POST | `/users/:id/block` | HR/Admin | Block user |
| POST | `/users/:id/unblock` | HR/Admin | Unblock user |

---

## Related Documentation

- [Rate Limiter & Blocked Users API](./RATE_LIMITER_AND_BLOCKED_USERS_API.md) - Login rate limiting and security
- [Auth API](./AUTH_API.md) - Authentication endpoints
- [Security & Audit API](./SECURITY_AUDIT_API.md) - Audit logging

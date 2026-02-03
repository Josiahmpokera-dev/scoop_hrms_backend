# Rate Limiter & Blocked Users API

**Base URL:** `BASE_URL` (e.g. `http://localhost:8080/api/v1`)

Login attempts are rate-limited: after too many **failed** login attempts, the user is **blocked from login**. Admins can list blocked users, check rate-limit status by email, and **unblock** (or suspend) users so they can log in again.

---

## Overview

- **Login rate limiting:** Each failed login (wrong password) increments a per-user failure count. When the count reaches the configured **max attempts** (default 5), the user’s status is set to **blocked** and they cannot log in until an admin unblocks them.
- **Blocked users:** Users with status `blocked` (or optionally `suspended`) are listed via **GET blocked-users**.
- **Rate-limit status:** **GET rate-limit/status?email=...** returns for a given email: blocked, failed attempt count, attempts remaining, locked-until (if any).
- **Unblock:** Use **POST /users/:id/unblock** (existing Users API) to restore login. Unblocking also clears the failed-attempt count and lockout window.

**Config (env):**

| Variable | Default | Description |
|----------|--------|-------------|
| `LOGIN_MAX_ATTEMPTS` | 5 | Max failed login attempts before blocking the user |
| `LOGIN_LOCKOUT_MINUTES` | 0 | Minutes to lock (0 = block until admin unblocks; if &gt; 0, `locked_until` is set for display) |

**Database:** The `users` table must have `failed_login_count` (int, default 0) and `locked_until` (timestamp nullable). If your app uses GORM AutoMigrate, these columns are added automatically on startup.

All security endpoints below require **Admin** (Bearer token).

---

## 1. List blocked users (users blocked from login)

**`GET BASE_URL/security/blocked-users`**

**Auth:** Bearer token  
**Authorization:** Admin only

Returns a paginated list of users who are **blocked** from login (rate-limit block or manual block). Optionally include **suspended** users.

### Query parameters

| Parameter | Type | Default | Description |
|-----------|------|--------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 20 | Items per page (max 100) |
| `include_suspended` | bool | false | If true, also return users with status `suspended` |

### Example request

```http
GET BASE_URL/security/blocked-users?page=1&page_size=20
Authorization: Bearer <admin_token>
```

With suspended users:

```http
GET BASE_URL/security/blocked-users?page=1&page_size=20&include_suspended=true
Authorization: Bearer <admin_token>
```

### Example response (200 OK)

```json
{
  "success": true,
  "message": "Blocked users retrieved successfully",
  "data": [
    {
      "id": 42,
      "email": "user@example.com",
      "firstName": "Jane",
      "lastName": "Doe",
      "username": "jane",
      "status": "blocked",
      "failedLoginCount": 5,
      "lockedUntil": "2026-01-24T12:00:00Z",
      "updatedAt": "2026-01-23T10:00:00Z"
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

## 2. Get rate-limit status by email

**`GET BASE_URL/security/rate-limit/status?email=...`**

**Auth:** Bearer token  
**Authorization:** Admin only

Returns login rate-limit status for the given email: blocked, suspended, failed attempt count, attempts remaining, locked-until, and whether the user can log in.

If no user exists for the email, a generic status is returned (e.g. `attemptsRemaining: maxAttempts`, `canLogin: true`) to avoid email enumeration.

### Query parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `email` | string | Yes | User email to look up |

### Example request

```http
GET BASE_URL/security/rate-limit/status?email=user@example.com
Authorization: Bearer <admin_token>
```

### Example response – user exists and is blocked (200 OK)

```json
{
  "success": true,
  "message": "Rate limit status retrieved",
  "data": {
    "email": "user@example.com",
    "blocked": true,
    "suspended": false,
    "failedLoginCount": 5,
    "maxAttempts": 5,
    "attemptsRemaining": 0,
    "lockedUntil": "2026-01-24T12:00:00Z",
    "canLogin": false
  }
}
```

### Example response – user exists and can log in (200 OK)

```json
{
  "success": true,
  "message": "Rate limit status retrieved",
  "data": {
    "email": "user@example.com",
    "blocked": false,
    "suspended": false,
    "failedLoginCount": 2,
    "maxAttempts": 5,
    "attemptsRemaining": 3,
    "lockedUntil": null,
    "canLogin": true
  }
}
```

---

## 3. Unblock a user (restore login)

Use the existing **Users** API to unblock a user so they can log in again. Unblocking also clears the failed-attempt count and any temporary lock.

**`POST BASE_URL/users/:id/unblock`**

**Auth:** Bearer token  
**Authorization:** HR or Admin

**Path:** `id` – user ID to unblock.

### Example request

```http
POST BASE_URL/users/42/unblock
Authorization: Bearer <admin_token>
```

### Example response (200 OK)

```json
{
  "success": true,
  "message": "User unblocked successfully",
  "data": {
    "id": 42,
    "email": "user@example.com",
    "firstName": "Jane",
    "lastName": "Doe",
    "status": "active",
    ...
  }
}
```

After this, the user can log in again and their failed-attempt count is reset.

---

## 4. Suspend a user (block from login)

To **suspend** a user (block from login without rate-limit):

**`POST BASE_URL/users/:id/suspend`**

**Auth:** Bearer token  
**Authorization:** HR or Admin

To **unsuspend** (restore login):

**`POST BASE_URL/users/:id/unsuspend`**

---

## 5. How to use the APIs

1. **See who is blocked from login**  
   `GET /api/v1/security/blocked-users` (optionally `include_suspended=true`).

2. **Check status for one email**  
   `GET /api/v1/security/rate-limit/status?email=user@example.com`.

3. **Unblock a user**  
   `POST /api/v1/users/:id/unblock` (HR or Admin). The user can log in again and failed-attempt state is cleared.

4. **Manually block/suspend**  
   `POST /api/v1/users/:id/block` or `POST /api/v1/users/:id/suspend` (HR or Admin).

5. **Config**  
   Set `LOGIN_MAX_ATTEMPTS` (default 5) and optionally `LOGIN_LOCKOUT_MINUTES` (default 0) in your environment.

All security endpoints (`/security/blocked-users`, `/security/rate-limit/status`) require an **Admin** Bearer token.

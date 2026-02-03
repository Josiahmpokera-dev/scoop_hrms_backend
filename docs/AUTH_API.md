# Authentication API

This document describes the authentication APIs: **login**, **refresh token**, **logout**, **register**, and **profile**.

**Base URL:** `{{BASE_URL}}` (e.g. `http://localhost:8080/api/v1`)

---

## Table of Contents

1. [Login](#1-login)
2. [Refresh Token](#2-refresh-token)
3. [Logout](#3-logout)
4. [Register](#4-register)
5. [Get Profile](#5-get-profile)
6. [Staying logged in (refresh flow)](#6-staying-logged-in-refresh-flow)

---

## 1. Login

**`POST {{BASE_URL}}/auth/login`**

Authenticate with email and password. Returns **access token** and **refresh token** so the user can call protected APIs and stay logged in by refreshing when the access token expires.

### Request body (JSON)

| Field     | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `email`   | string | Yes      | User email  |
| `password`| string | Yes      | Password (min 6 characters) |

### Example

```http
POST {{BASE_URL}}/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

### Response (200 OK)

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "username": "user",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user",
      "roles": ["user"],
      "is_active": true
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "refresh_expires_in": 604800
  }
}
```

| Field                | Description |
|----------------------|-------------|
| `access_token`       | JWT to send in `Authorization: Bearer <access_token>` for protected APIs |
| `refresh_token`      | JWT to send to `/auth/refresh` when the access token expires (e.g. after 24h) |
| `expires_in`         | Access token lifetime in seconds (e.g. 86400 = 24h) |
| `refresh_expires_in` | Refresh token lifetime in seconds (e.g. 604800 = 7 days) |

### Error responses

- **401** – Invalid email or password, or account suspended/blocked/deactivated.

---

## 2. Refresh Token

**`POST {{BASE_URL}}/auth/refresh`**

Exchange a valid **refresh token** for a new **access token** and a new **refresh token**. Use this to keep the user logged in without asking for password again when the access token expires.

- No authentication required (public endpoint).
- Request body must contain the current `refresh_token`.

### Request body (JSON)

| Field          | Type   | Required | Description |
|----------------|--------|----------|-------------|
| `refresh_token`| string | Yes      | The refresh token received from login or a previous refresh |

### Example

```http
POST {{BASE_URL}}/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Response (200 OK)

Same shape as [Login](#1-login):

```json
{
  "success": true,
  "message": "Tokens refreshed successfully",
  "data": {
    "user": {
      "id": 1,
      "username": "user",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user",
      "roles": ["user"],
      "is_active": true
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "refresh_expires_in": 604800
  }
}
```

- **access_token** – New JWT to use for API calls.
- **refresh_token** – New refresh token; use this for the next refresh (old one is not invalidated but best practice is to use the new one).

### Error responses

- **400** – `refresh_token` missing or invalid JSON.
- **401** – Invalid or expired refresh token, or account suspended/blocked/deactivated.

---

## 3. Logout

**`POST {{BASE_URL}}/auth/logout`**

**Auth:** Bearer token (access token) required.

Logout is handled on the client: discard the access token and refresh token (e.g. remove from storage). This endpoint only confirms logout and can be used for audit or cleanup.

### Example

```http
POST {{BASE_URL}}/auth/logout
Authorization: Bearer <access_token>
```

### Response (200 OK)

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

## 4. Register

**`POST {{BASE_URL}}/auth/register`**

Create a new user account. Returns access and refresh tokens like login.

### Request body (JSON)

| Field        | Type   | Required | Description |
|--------------|--------|----------|-------------|
| `email`      | string | Yes      | User email  |
| `password`   | string | Yes      | Password (min 6 characters) |
| `first_name` | string | Yes      | First name (min 2 characters) |
| `last_name`  | string | Yes      | Last name (min 2 characters) |
| `username`   | string | No       | Optional; generated from email if omitted |
| `role`       | string | No       | Optional; e.g. "user", "admin", "hr" |

### Response (201 Created)

Same structure as [Login](#1-login) response (user + `access_token` + `refresh_token` + `expires_in` + `refresh_expires_in`).

---

## 5. Get Profile

**`GET {{BASE_URL}}/auth/profile`**

**Auth:** Bearer token (access token) required.

Returns the authenticated user’s profile.

### Example

```http
GET {{BASE_URL}}/auth/profile
Authorization: Bearer <access_token>
```

### Response (200 OK)

```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "username": "user",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user",
    "roles": ["user"],
    "is_active": true
  }
}
```

---

## 6. Staying logged in (refresh flow)

1. **Login** – Send email/password to `POST /auth/login`. Store both `access_token` and `refresh_token` (e.g. in memory or secure storage).
2. **API calls** – Use `access_token` in the `Authorization: Bearer <access_token>` header for all protected endpoints.
3. **When access token expires** (e.g. 401 from API):
   - Call **`POST /auth/refresh`** with the stored `refresh_token` in the body: `{ "refresh_token": "..." }`.
   - Replace stored `access_token` and `refresh_token` with the new ones from the response.
   - Retry the failed request with the new `access_token`.
4. **When refresh token expires or is invalid** – Redirect the user to the login page (re-authenticate with email/password).
5. **Logout** – Remove both tokens on the client and optionally call `POST /auth/logout` with the current access token.

### Configuration

- **Access token expiry:** `JWT_EXPIRY` (e.g. `24h`). Default: 24 hours.
- **Refresh token expiry:** `JWT_REFRESH_EXPIRY` (e.g. `168h` = 7 days). Default: 168h.

Set these in your `.env` or environment.

---

## Quick reference

| Endpoint              | Method | Auth    | Description |
|-----------------------|--------|--------|-------------|
| `/auth/login`         | POST   | No     | Login; returns access + refresh tokens |
| `/auth/refresh`       | POST   | No     | Exchange refresh token for new tokens |
| `/auth/logout`        | POST   | Bearer | Confirm logout (client discards tokens) |
| `/auth/register`      | POST   | No     | Register; returns access + refresh tokens |
| `/auth/profile`       | GET    | Bearer | Get current user profile |

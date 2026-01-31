# Security & Audit API

**Base URL:** `BASE_URL` (e.g. `http://localhost:8080/api/v1`)

The system logs **every API request** to the `audit_logs` table (who did what, when, where, and outcome). Only **Admin** can list audit logs.

---

## List audit logs

**`GET BASE_URL/security/audit`**

**Auth:** Bearer token  
**Authorization:** Admin only

Returns paginated audit logs with optional filters and search. Results are ordered by **newest first**.

### Query parameters

| Parameter   | Type   | Default | Description |
|-------------|--------|--------|-------------|
| `page`      | int    | 1      | Page number |
| `page_size` | int    | 20     | Items per page (max 100) |
| `user_id`   | int    | —      | Filter by user ID (who performed the action) |
| `action`    | string | —      | Filter by action (partial, case-insensitive) |
| `resource`  | string | —      | Filter by resource (partial, case-insensitive) |
| `method`    | string | —      | Filter by HTTP method: GET, POST, PUT, PATCH, DELETE |
| `date_from` | string | —      | Filter from date (RFC3339 or `2006-01-02`) |
| `date_to`   | string | —      | Filter to date (RFC3339 or `2006-01-02`) |
| `status_code` | int  | —      | Filter by HTTP response status code (e.g. 200, 400) |
| `search`    | string | —      | Search in action, resource, path, details, reason (partial, case-insensitive) |

### Example – all logs (paginated)

```http
GET BASE_URL/security/audit?page=1&page_size=20
Authorization: Bearer <admin_token>
```

### Example – filter by user and date range

```http
GET BASE_URL/security/audit?user_id=5&date_from=2026-01-01&date_to=2026-01-31&page=1&page_size=20
Authorization: Bearer <admin_token>
```

### Example – search and filter by resource

```http
GET BASE_URL/security/audit?resource=users&search=assign-role&method=POST
Authorization: Bearer <admin_token>
```

### Response (200 OK)

```json
{
  "success": true,
  "message": "Audit logs retrieved successfully",
  "data": [
    {
      "id": 1,
      "tenant_id": 1,
      "user_id": 2,
      "action": "assign-role",
      "resource": "users",
      "resource_id": null,
      "method": "POST",
      "path": "/api/v1/users/assign-role",
      "status_code": 200,
      "ip": "192.168.1.10",
      "user_agent": "Mozilla/5.0 ...",
      "details": "",
      "reason": "",
      "created_at": "2026-01-23T14:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

### Audit log fields

| Field        | Description |
|-------------|-------------|
| `id`        | Audit log ID |
| `tenant_id` | Tenant (if any) |
| `user_id`   | User who performed the action (null if unauthenticated) |
| `action`    | Action (e.g. assign-role, login, list) |
| `resource`  | Resource (e.g. users, employees, leave) |
| `resource_id` | Target entity ID if applicable |
| `method`    | HTTP method (GET, POST, etc.) |
| `path`      | Request path |
| `status_code` | HTTP response status |
| `ip`        | Client IP |
| `user_agent` | Client user agent |
| `details`   | Optional JSON or text details |
| `reason`    | Optional reason/notes |
| `created_at` | When the action occurred (UTC) |

---

## How audit logging works

- **Every request** to `/api/v1/*` is logged by the **audit middleware** after the handler runs.
- Each log stores: **who** (user_id), **what** (action, resource, method, path), **when** (created_at), **where** (ip, user_agent), and **outcome** (status_code).
- Action and resource are derived from the request path (e.g. `/api/v1/users/assign-role` → resource `users`, action `assign-role`).
- For unauthenticated requests (e.g. login), `user_id` and `tenant_id` may be null; path, method, IP, and status are still recorded.

---

## Summary

| Item   | Description |
|--------|-------------|
| Endpoint | `GET /security/audit` |
| Auth     | Bearer token, Admin only |
| Pagination | `page`, `page_size` (max 100) |
| Filters  | `user_id`, `action`, `resource`, `method`, `date_from`, `date_to`, `status_code` |
| Search   | `search` (action, resource, path, details, reason) |

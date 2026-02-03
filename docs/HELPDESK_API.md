# Helpdesk & Support — API Documentation

This document describes the APIs for the **Support & Reports** module:

- **My Tickets** — `/api/v1/helpdesk` (tickets, comments, attachments, close, CSAT)
- **Knowledge Base** — `/api/v1/helpdesk/knowledge-base`
- **Helpdesk Dashboard** — `/api/v1/helpdesk/dashboard` and `/api/v1/helpdesk/reports`

---

## Role-based access

| Area | Who can access | Notes |
|------|----------------|--------|
| **My Tickets** (create, list, view, comment, attach, close, CSAT, ticket-categories) | **All authenticated users** | Admin, HR, IT, and Employee can all create their own tickets and see only **their own** tickets. |
| **Knowledge Base** (categories, articles, article detail, feedback) | **All authenticated users** | Browse and search; submit helpful/not helpful. |
| **Helpdesk Dashboard** (statistics, recent tickets, export report) | **HR or Admin** only | KPIs, charts, export. |
| **Agent / Admin** (list all tickets, assign, update status, resolve) | **HR or Admin** only | Triage and manage any ticket in the tenant. |

- **All users** includes: Admin, HR, IT, Employee, and any other role. Every authenticated user can open a ticket (e.g. “I have an issue”) and manage their own tickets.
- **HR or Admin** is enforced by the backend (e.g. `HRMiddleware`) for dashboard and agent routes.

---

## Implementation: tickets for all users (including Admin)

- **Requester linking**
  - If the user has an **employee record**: the ticket is stored with `requester_id` (employee ID) and `requester_employee_id` (e.g. `EMP-001`). The requester in the API is shown with that employee id, name, email, phone, department.
  - If the user has **no employee record** (e.g. Admin): the ticket is stored with `requester_user_id` (user ID), `requester_name`, and `requester_email`. The requester in the API is shown with `id: "USER-{user_id}"`, plus name and email from the user.
- **“My tickets”** list returns tickets where the current user is the requester, either as employee or as user-only (so Admin sees tickets they created, employees see tickets they created).
- **Ownership** for view, comment, attach, close, CSAT: the current user must be the requester (employee or user-only). Agents/Admins can still access any ticket via the Agent APIs.

---

## Table of Contents

1. [Conventions](#1-conventions)
2. [My Tickets](#2-my-tickets)
3. [Knowledge Base](#3-knowledge-base)
4. [Helpdesk Dashboard & Reports](#4-helpdesk-dashboard--reports)
5. [Agent / Admin APIs](#5-agent--admin-apis)
6. [Enums Reference](#6-enums-reference)
7. [Quick reference: pages → APIs](#7-quick-reference-pages--apis)

---

## 1. Conventions

### Base URL & Auth

- **Base URL:** `{{BASE_URL}}` (e.g. `http://localhost:8080/api/v1`)
- **Auth:** All endpoints require `Authorization: Bearer <accessToken>`

### Response Envelope

**Success:**
```json
{
  "success": true,
  "message": "Human readable message",
  "data": {}
}
```

**Error:**
```json
{
  "success": false,
  "message": "Human readable error",
  "error": "ERROR_CODE or details",
  "details": { "field": ["validation message"] }
}
```

### Pagination

List endpoints return data and meta inside `data`:

```json
{
  "success": true,
  "message": "...",
  "data": {
    "data": [],
    "meta": {
      "page": 1,
      "per_page": 20,
      "total": 124,
      "total_pages": 7
    }
  }
}
```

### JSON Keys

All API responses use **camelCase** for field names (e.g. `ticketNumber`, `createdAt`, `assignedTo`).

### Requester in responses

- **Employee requester:** `requester.id` is the employee ID string (e.g. `EMP-001`); `name`, `email`, `phone`, `department` come from the employee record.
- **User-only requester** (e.g. Admin with no employee): `requester.id` is `USER-{user_id}` (e.g. `USER-5`); `name` and `email` come from the user; `phone` and `department` are empty.

---

## 2. My Tickets

Used for: ticket list (with tabs by status), search, detail drawer, create ticket, comments, attachments, close, CSAT.

**Access:** All authenticated users (Admin, HR, IT, Employee) can create tickets and see only their own tickets. No role restriction: only `Authorization: Bearer <token>` is required.

### 2.1 List my tickets

**`GET {{BASE_URL}}/helpdesk/tickets`**

Used for: ticket list, tabs (All / Open / In Progress / Resolved), search.

| Parameter   | Type    | Required | Description |
|------------|---------|----------|-------------|
| `page`     | number  | No       | Default `1` |
| `page_size`| number  | No       | Default `20`, max `100` |
| `search`   | string  | No       | Matches title, description, ticket number |
| `status`   | string  | No       | One of status enums (filter by tab) |
| `priority` | string  | No       | One of priority enums |
| `category` | string  | No       | e.g. `IT`, `HR`, `Payroll` |
| `mine`     | boolean | No       | Default `true` — only requester’s tickets |
| `sort_by`  | string  | No       | `createdAt`, `updatedAt`, `dueDate` |
| `sort_dir` | string  | No       | `asc`, `desc` |

**Example:**
```http
GET {{BASE_URL}}/helpdesk/tickets?page=1&page_size=20&status=Open
Authorization: Bearer <accessToken>
```

**Response (200):**
```json
{
  "success": true,
  "message": "Tickets retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "ticketNumber": "HD-2025-001",
        "title": "Unable to access payslip for January 2025",
        "description": "I am unable to view or download my January 2025 payslip.",
        "category": "Payroll",
        "subCategory": "Payslip Access",
        "priority": "High",
        "status": "In Progress",
        "requester": {
          "id": "EMP-001",
          "name": "Sarah Johnson",
          "email": "sarah.johnson@company.com",
          "phone": "+255 756 123 456",
          "department": ""
        },
        "assignedTo": {
          "id": 2,
          "name": "Michael Chen",
          "email": "michael.chen@company.com",
          "team": "HR Support"
        },
        "createdAt": "2025-01-20T09:30:00Z",
        "updatedAt": "2025-01-20T10:15:00Z",
        "dueDate": "2025-01-21T09:30:00Z",
        "slaStatus": "Within SLA",
        "firstResponseTime": "15 minutes",
        "channel": "Portal",
        "tags": ["payslip", "access-issue"],
        "attachments": [
          {
            "id": 1,
            "name": "error-screenshot.png",
            "url": "/storage/helpdesk/1/error-screenshot_20250120_093200_abc123.png",
            "size": "245 KB",
            "type": "image/png",
            "uploadedAt": "2025-01-20T09:32:00Z"
          }
        ],
        "commentsCount": 2
      }
    ],
    "meta": { "page": 1, "per_page": 20, "total": 124, "total_pages": 7 }
  }
}
```

---

### 2.2 Get ticket details

**`GET {{BASE_URL}}/helpdesk/tickets/:ticket_id`**

Used for: ticket detail drawer (description, attachments, comments).

**Example:**
```http
GET {{BASE_URL}}/helpdesk/tickets/1
Authorization: Bearer <accessToken>
```

**Response (200):**
```json
{
  "success": true,
  "message": "Ticket retrieved successfully",
  "data": {
    "id": 1,
    "ticketNumber": "HD-2025-001",
    "title": "Unable to access payslip for January 2025",
    "description": "I am unable to view or download my January 2025 payslip. The page shows an error message.",
    "category": "Payroll",
    "subCategory": "Payslip Access",
    "priority": "High",
    "status": "In Progress",
    "requester": { "id": "EMP-001", "name": "Sarah Johnson", "email": "...", "phone": "...", "department": "" },
    "assignedTo": { "id": 2, "name": "Michael Chen", "email": "...", "team": "HR Support" },
    "createdAt": "2025-01-20T09:30:00Z",
    "updatedAt": "2025-01-20T10:15:00Z",
    "resolvedAt": null,
    "closedAt": null,
    "dueDate": "2025-01-21T09:30:00Z",
    "slaStatus": "Within SLA",
    "firstResponseTime": "15 minutes",
    "resolutionTime": null,
    "channel": "Portal",
    "tags": ["payslip", "access-issue"],
    "attachments": [
      { "id": 1, "name": "error-screenshot.png", "url": "...", "size": "245 KB", "type": "image/png", "uploadedAt": "2025-01-20T09:32:00Z" }
    ],
    "comments": [
      {
        "id": 1,
        "author": "Sarah Johnson",
        "authorType": "Requester",
        "text": "I tried clearing my browser cache but the issue persists.",
        "timestamp": "2025-01-20T09:32:00Z",
        "isPrivate": false
      },
      {
        "id": 2,
        "author": "Michael Chen",
        "authorType": "Agent",
        "text": "Thank you for reporting this. I am investigating the issue with our IT team.",
        "timestamp": "2025-01-20T10:15:00Z",
        "isPrivate": false
      }
    ],
    "csatRating": null,
    "csatComment": null
  }
}
```

---

### 2.3 Create ticket

**`POST {{BASE_URL}}/helpdesk/tickets`**

Used for: “Create Ticket” dialog. Supports **multipart/form-data** (e.g. for future attachments on create).

| Field         | Type   | Required | Description |
|---------------|--------|----------|-------------|
| `category`    | string | Yes      | e.g. `IT`, `HR`, `Payroll` |
| `sub_category`| string | No       | e.g. `Payslip Access`, `VPN Access` |
| `priority`    | string | Yes      | `Low`, `Medium`, `High`, `Critical` |
| `title`       | string | Yes      | Subject |
| `description`  | string | Yes      | Body |
| `channel`     | string | No       | Default `Portal` |
| `tags`        | string | No       | Comma-separated |

**Example (form-data):**
```http
POST {{BASE_URL}}/helpdesk/tickets
Authorization: Bearer <accessToken>
Content-Type: multipart/form-data

category=Payroll&sub_category=Payslip Access&priority=High&title=Unable to access payslip&description=...
```

**Response (200):**
```json
{
  "success": true,
  "message": "Ticket created successfully",
  "data": {
    "id": 109,
    "ticketNumber": "HD-2026-109",
    "status": "Open",
    "createdAt": "2026-01-21T10:20:00Z"
  }
}
```

---

### 2.4 Add comment

**`POST {{BASE_URL}}/helpdesk/tickets/:ticket_id/comments`**

| Field        | Type    | Required | Description |
|-------------|---------|----------|-------------|
| `text`      | string  | Yes      | Comment body |
| `is_private`| boolean | No      | Default `false` (requester typically cannot set private) |

**Request body (JSON):**
```json
{
  "text": "Any update on this ticket?",
  "is_private": false
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Comment posted successfully",
  "data": {
    "id": 991,
    "ticket_id": 1,
    "author": "Sarah Johnson",
    "authorType": "Requester",
    "text": "Any update on this ticket?",
    "timestamp": "2026-01-21T10:25:00Z",
    "isPrivate": false
  }
}
```

---

### 2.5 Upload attachments to ticket

**`POST {{BASE_URL}}/helpdesk/tickets/:ticket_id/attachments`**

**Request:** `multipart/form-data`

| Field            | Type | Required |
|------------------|------|----------|
| `attachments[]`  | file | Yes (at least one) |
| or `file`        | file | Alternative single file |

**Example:**
```http
POST {{BASE_URL}}/helpdesk/tickets/1/attachments
Authorization: Bearer <accessToken>
Content-Type: multipart/form-data

attachments[]=@log.txt
```

**Response (200):**
```json
{
  "success": true,
  "message": "Attachments uploaded successfully",
  "data": [
    {
      "id": 778,
      "name": "log.txt",
      "url": "/storage/helpdesk/1/log_20260121_103000_xyz.txt",
      "size": "12 KB",
      "type": "text/plain",
      "uploadedAt": "2026-01-21T10:30:00Z"
    }
  ]
}
```

---

### 2.6 Close ticket

**`POST {{BASE_URL}}/helpdesk/tickets/:ticket_id/close`**

| Field     | Type   | Required |
|-----------|--------|----------|
| `comment`| string | No       |

**Request body (JSON):**
```json
{
  "comment": "Issue resolved on my side. Closing ticket."
}
```

**Response (200):** `{ "success": true, "message": "Ticket closed successfully", "data": null }`

---

### 2.7 Submit CSAT

**`POST {{BASE_URL}}/helpdesk/tickets/:ticket_id/csat`**

| Field    | Type   | Required | Description |
|----------|--------|----------|-------------|
| `rating` | number | Yes      | 1–5 |
| `comment`| string | No       | Optional feedback text |

**Request body (JSON):**
```json
{
  "rating": 5,
  "comment": "Excellent and quick support!"
}
```

**Response (200):** `{ "success": true, "message": "CSAT submitted successfully", "data": null }`

---

### 2.8 Get ticket categories (for Create Ticket form)

**`GET {{BASE_URL}}/helpdesk/ticket-categories`**

Used for: category/subcategory dropdowns and SLA info in Create Ticket.

**Response (200):**
```json
{
  "success": true,
  "message": "Ticket categories retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "IT",
      "description": "IT support and technical issues",
      "sla": { "firstResponse": "30", "resolution": "8" },
      "defaultPriority": "High",
      "requiredFields": ["description", "category", "subCategory"],
      "templates": [
        {
          "id": "TPL-IT-001",
          "name": "VPN Issue",
          "fields": { "errorCode": "", "deviceType": "", "internetProvider": "" }
        }
      ]
    }
  ]
}
```

---

## 3. Knowledge Base

Used for: search, featured articles, browse by category, article detail, helpful/not helpful.

### 3.1 List KB categories

**`GET {{BASE_URL}}/helpdesk/knowledge-base/categories`**

Used for: “Browse by Category” (name + article count).

**Response (200):**
```json
{
  "success": true,
  "message": "Categories retrieved successfully",
  "data": [
    { "name": "IT", "count": 12 },
    { "name": "Payroll", "count": 7 },
    { "name": "HR", "count": 9 },
    { "name": "Leave", "count": 6 }
  ]
}
```

---

### 3.2 List KB articles

**`GET {{BASE_URL}}/helpdesk/knowledge-base/articles`**

| Parameter   | Type    | Required | Description |
|------------|---------|----------|-------------|
| `page`     | number  | No       | Default `1` |
| `page_size`| number  | No       | Default `20`, max `100` |
| `search`   | string  | No       | Matches title, content, tags |
| `category`| string  | No       | Filter by category name |
| `featured` | boolean | No       | Only featured articles |
| `status`   | string  | No       | Default `Published` for end users |
| `sort_by`  | string  | No       | `views`, `updatedAt` |
| `sort_dir` | string  | No       | `asc`, `desc` |

**Response (200):**
```json
{
  "success": true,
  "message": "Articles retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "title": "How to Access Your Payslip",
        "content_preview": "How to Access Your Payslip ...",
        "category": "Payroll",
        "tags": ["payslip", "download", "self-service", "how-to"],
        "views": 1247,
        "helpful": 98,
        "notHelpful": 12,
        "author": "HR Team",
        "createdAt": "2024-12-01T10:00:00Z",
        "updatedAt": "2025-01-15T14:30:00Z",
        "status": "Published",
        "featured": true
      }
    ],
    "meta": { "page": 1, "per_page": 20, "total": 34, "total_pages": 2 }
  }
}
```

---

### 3.3 Get KB article details

**`GET {{BASE_URL}}/helpdesk/knowledge-base/articles/:article_id`**

Used for: article detail dialog (full content, tags, “Was this helpful?”).  
**Note:** View count is incremented when the article is opened.

**Response (200):**
```json
{
  "success": true,
  "message": "Article retrieved successfully",
  "data": {
    "id": 1,
    "title": "How to Access Your Payslip",
    "content": "# How to Access Your Payslip\n\n## Step-by-Step Guide\n\n1. Login ...",
    "category": "Payroll",
    "tags": ["payslip", "download", "self-service", "how-to"],
    "views": 1248,
    "helpful": 98,
    "notHelpful": 12,
    "author": "HR Team",
    "createdAt": "2024-12-01T10:00:00Z",
    "updatedAt": "2025-01-15T14:30:00Z",
    "status": "Published",
    "featured": true
  }
}
```

---

### 3.4 Article helpful feedback

**`POST {{BASE_URL}}/helpdesk/knowledge-base/articles/:article_id/feedback`**

Used for: “Was this article helpful?” Yes/No in article detail.

| Field        | Type    | Required | Description |
|-------------|---------|----------|-------------|
| `is_helpful`| boolean | Yes      | `true` = Yes, `false` = No |

**Request body (JSON):**
```json
{
  "is_helpful": true
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Feedback recorded successfully",
  "data": {
    "article_id": 1,
    "helpful": 99,
    "notHelpful": 12
  }
}
```

---

## 4. Helpdesk Dashboard & Reports

Used for: KPI cards, tickets by category/priority, top agents, recent tickets table, export.

**Auth:** Bearer token + **HR or Admin** role.

### 4.1 Dashboard statistics

**`GET {{BASE_URL}}/helpdesk/dashboard/statistics`**

| Parameter    | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `date_range`| string | No       | `today`, `thisWeek`, `thisMonth`, `custom` |
| `from`      | string | If `date_range=custom` | ISO date, e.g. `2026-01-01` |
| `to`        | string | If `date_range=custom` | ISO date, e.g. `2026-01-31` |

**Response (200):**
```json
{
  "success": true,
  "message": "Helpdesk statistics retrieved successfully",
  "data": {
    "totalTickets": 847,
    "openTickets": 124,
    "inProgress": 56,
    "resolved": 589,
    "closed": 78,
    "avgFirstResponseTime": "28 minutes",
    "avgResolutionTime": "4.2 hours",
    "slaCompliance": 94.5,
    "csatScore": 4.3,
    "ticketsByCategory": { "IT": 412, "HR": 245, "Facilities": 98, "Payroll": 56, "Leave": 36 },
    "ticketsByPriority": { "Critical": 12, "High": 45, "Medium": 234, "Low": 556 },
    "topAgents": [
      { "id": "AGT-001", "name": "Michael Chen", "resolved": 156, "avgTime": "3.8h", "csat": 4.6 }
    ]
  }
}
```

*Note: `avgFirstResponseTime`, `avgResolutionTime`, `slaCompliance`, `csatScore`, and `topAgents` may be placeholders (e.g. 0 or empty array) until backend metrics are fully implemented.*

---

### 4.2 Recent tickets

**`GET {{BASE_URL}}/helpdesk/tickets/recent`**

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `limit`  | number | No       | Default `6`, max `50` |

**Response (200):**
```json
{
  "success": true,
  "message": "Recent tickets retrieved successfully",
  "data": [
    {
      "id": 1,
      "ticketNumber": "HD-2025-001",
      "title": "Unable to access payslip for January 2025",
      "requester": { "name": "Sarah Johnson", "department": "Sales" },
      "category": "Payroll",
      "priority": "High",
      "status": "In Progress",
      "assignedTo": { "name": "Michael Chen" },
      "createdAt": "2025-01-20T09:30:00Z"
    }
  ]
}
```

---

### 4.3 Export report

**`GET {{BASE_URL}}/helpdesk/reports/export`**

Returns a **binary file** with `Content-Disposition: attachment; filename="helpdesk_report.csv"` (or `.xlsx`).

| Parameter    | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `format`    | string | No       | `csv` (default), `xlsx` |
| `date_range`| string | No       | Same as statistics |
| `from`      | string | No       | Required if `date_range=custom` |
| `to`        | string | No       | Required if `date_range=custom` |

**Example:**
```http
GET {{BASE_URL}}/helpdesk/reports/export?format=csv&date_range=thisMonth
Authorization: Bearer <accessToken>
```

**Response:** Binary file (CSV or XLSX).  
*Note: XLSX format may be served as CSV with `.xlsx` filename; full XLSX support can be added via a library such as excelize.*

---

## 5. Agent / Admin APIs

**Base path:** `{{BASE_URL}}/helpdesk/agent`  
**Auth:** Bearer token + **HR or Admin** role.

### 5.1 List tickets (agent queue)

**`GET {{BASE_URL}}/helpdesk/agent/tickets`**

| Parameter     | Type   | Required | Description |
|--------------|--------|----------|-------------|
| `page`       | number | No       | Default `1` |
| `page_size`  | number | No       | Default `20`, max `100` |
| `search`     | string | No       | Matches title, description, ticket number |
| `status`     | string | No       | Status filter |
| `priority`   | string | No       | Priority filter |
| `category`   | string | No       | Category filter |
| `assigned_to`| string | No       | `me`, `unassigned`, or user ID |
| `queue`      | string | No       | Queue name |

**Response (200):** Same paginated envelope as [2.1 List my tickets](#21-list-my-tickets), with `data.data` and `data.meta`.

---

### 5.2 Assign ticket

**`POST {{BASE_URL}}/helpdesk/agent/tickets/:ticket_id/assign`**

**Request body (JSON):**
```json
{
  "assignee_user_id": 2,
  "note": "Assigning to you for follow-up."
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Ticket assigned successfully",
  "data": {
    "ticket_id": 1,
    "assignedTo": { "id": 2, "name": "Michael Chen", "email": "...", "team": "HR Support" },
    "updatedAt": "2026-01-21T11:00:00Z"
  }
}
```

---

### 5.3 Update ticket status

**`PATCH {{BASE_URL}}/helpdesk/agent/tickets/:ticket_id/status`**

**Request body (JSON):**
```json
{
  "status": "In Progress",
  "note": "Starting investigation."
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Ticket status updated successfully",
  "data": { "ticket_id": 1, "status": "In Progress", "updatedAt": "2026-01-21T11:05:00Z" }
}
```

---

### 5.4 Resolve ticket

**`POST {{BASE_URL}}/helpdesk/agent/tickets/:ticket_id/resolve`**

**Request body (JSON):**
```json
{
  "resolution_summary": "Issue was due to cache. Cleared and verified.",
  "internal_note": "Internal note for agents only."
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Ticket resolved successfully",
  "data": { "ticket_id": 1, "status": "Resolved", "resolvedAt": "2026-01-21T11:30:00Z" }
}
```

---

## 6. Enums Reference

| Type               | Values |
|--------------------|--------|
| **Ticket status**  | `Open`, `In Progress`, `Pending`, `Resolved`, `Closed`, `Cancelled` |
| **Ticket priority**| `Critical`, `High`, `Medium`, `Low` |
| **Ticket channel** | `Portal`, `Email`, `WhatsApp`, `Phone`, `Chat` |
| **KB article status** | `Published`, `Draft`, `Archived` |
| **Ticket category**| `IT`, `HR`, `Facilities`, `Payroll`, `Leave`, `General` |
| **Comment author type** | `Requester`, `Agent`, `Admin` |

---

## 7. Quick reference: pages → APIs

| Page            | Purpose | APIs used |
|-----------------|--------|-----------|
| **My Tickets**  | List/filter tickets, view detail, create ticket, add comment, attachments, close, CSAT | `GET /helpdesk/tickets`, `GET /helpdesk/tickets/:id`, `POST /helpdesk/tickets`, `POST /helpdesk/tickets/:id/comments`, `POST /helpdesk/tickets/:id/attachments`, `POST /helpdesk/tickets/:id/close`, `POST /helpdesk/tickets/:id/csat`, `GET /helpdesk/ticket-categories` |
| **Knowledge Base** | Browse categories, search articles, view article, submit helpful feedback | `GET /helpdesk/knowledge-base/categories`, `GET /helpdesk/knowledge-base/articles`, `GET /helpdesk/knowledge-base/articles/:id`, `POST /helpdesk/knowledge-base/articles/:id/feedback` |
| **Helpdesk Dashboard** | KPIs, charts, recent tickets, export | `GET /helpdesk/dashboard/statistics`, `GET /helpdesk/tickets/recent`, `GET /helpdesk/reports/export` |

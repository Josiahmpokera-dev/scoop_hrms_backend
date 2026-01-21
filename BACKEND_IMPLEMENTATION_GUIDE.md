# Helpdesk & Support API Documentation (Proposed)

This document proposes API endpoints for the **Helpdesk & Support** module pages:
- **HelpDesk Dashboard** (`/helpdesk/dashboard`)
- **My Tickets** (`/helpdesk/my-tickets`)
- **Knowledge Base** (`/helpdesk/knowledge-base`)

The payloads are aligned to the UI fields currently used in:
- `src/app/(protected-pages)/helpdesk/*`
- `src/mock/data/helpdesk/helpdeskData.ts`

---

## Conventions

### Base URL
- `{{BASE_URL}}`

### Authentication
- **Bearer token**: `Authorization: Bearer <accessToken>`

### Roles & access model (recommended)

This module has three practical access levels:
- **Employee (Requester)**: can create tickets, view **own** tickets, comment, add attachments, close own ticket, submit CSAT.
- **HR/Agent (Resolver)**: can view/triage assigned tickets (and tickets in their queue), comment (public/private), change status, resolve/reopen, request more info, add internal notes, attach files, assign to self/others (if allowed).
- **Admin (Helpdesk Manager)**: full visibility across org, can configure routing rules/queues/SLAs, manage categories, reassign tickets, override SLA, audit.

Implementation note:
- You can model **HR/Agent** as users with permission `HELPDESK_AGENT`, and **Admin** as `HELPDESK_ADMIN`.
- For “HR vs IT vs Facilities”, route by **team/queue** rather than role name. HR users would belong to `HR Support` queue, IT users to `IT Support`, etc.

### Standard response envelope

Success:

```json
{
  "success": true,
  "message": "Human readable message",
  "data": {}
}
```

Error:

```json
{
  "success": false,
  "message": "Human readable error message",
  "error": "Machine readable error code",
  "details": {
    "field": ["optional validation message"]
  }
}
```

### Pagination format

```json
{
  "data": [],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 124,
    "total_pages": 7
  }
}
```

### Common enums (as used by UI)
- **Ticket status**: `Open | In Progress | Pending | Resolved | Closed | Cancelled`
- **Ticket priority**: `Critical | High | Medium | Low`
- **Ticket channel**: `Portal | Email | WhatsApp | Phone | Chat`
- **KB article status**: `Published | Draft | Archived`

---

## 1) Employee (Requester) APIs — “My Tickets”

These endpoints are intended for **normal employees** (requesters).
Rule: requester can only access tickets where `requester.id == current_user.employee_id` (or equivalent identity).

### 1.1 List tickets (for tabs + search)
**GET** `{{BASE_URL}}/helpdesk/tickets`

#### Query parameters
- `page` (optional, default `1`)
- `page_size` (optional, default `20`, max `100`)
- `search` (optional) — matches title/description/ticket number
- `status` (optional) — one of status enums
- `priority` (optional) — one of priority enums
- `category` (optional) — e.g. `IT`, `HR`, `Payroll`, `Facilities`, `Leave`, `General`
- `mine` (optional, boolean, default `true`) — “My Tickets”
- `sort_by` (optional) — `createdAt | updatedAt | dueDate`
- `sort_dir` (optional) — `asc | desc`

#### Sample response

```json
{
  "success": true,
  "message": "Tickets retrieved successfully",
  "data": {
    "data": [
      {
        "id": "TKT-001",
        "ticketNumber": "HD-2025-001",
        "title": "Unable to access payslip for January 2025",
        "description": "I am unable to view or download my January 2025 payslip. The page shows an error message.",
        "category": "Payroll",
        "subCategory": "Payslip Access",
        "priority": "High",
        "status": "In Progress",
        "requester": {
          "id": "EMP-001",
          "name": "Sarah Johnson",
          "email": "sarah.johnson@company.com",
          "phone": "+255 756 123 456",
          "department": "Sales"
        },
        "assignedTo": {
          "id": "AGT-001",
          "name": "Michael Chen",
          "email": "michael.chen@company.com",
          "team": "HR Support"
        },
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
          {
            "id": "ATT-001",
            "name": "error-screenshot.png",
            "url": "https://files.company.com/helpdesk/ATT-001",
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

Notes:
- For list performance, you can return `commentsCount` instead of full comments.
- The UI also displays attachments count; either include `attachments` array or `attachmentsCount`.

---

### 1.2 Get ticket details (for detail dialog/page)
**GET** `{{BASE_URL}}/helpdesk/tickets/:ticket_id`

#### Sample response

```json
{
  "success": true,
  "message": "Ticket retrieved successfully",
  "data": {
    "id": "TKT-001",
    "ticketNumber": "HD-2025-001",
    "title": "Unable to access payslip for January 2025",
    "description": "I am unable to view or download my January 2025 payslip. The page shows an error message.",
    "category": "Payroll",
    "subCategory": "Payslip Access",
    "priority": "High",
    "status": "In Progress",
    "requester": {
      "id": "EMP-001",
      "name": "Sarah Johnson",
      "email": "sarah.johnson@company.com",
      "phone": "+255 756 123 456",
      "department": "Sales"
    },
    "assignedTo": {
      "id": "AGT-001",
      "name": "Michael Chen",
      "email": "michael.chen@company.com",
      "team": "HR Support"
    },
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
      {
        "id": "ATT-001",
        "name": "error-screenshot.png",
        "url": "https://files.company.com/helpdesk/ATT-001",
        "size": "245 KB",
        "type": "image/png",
        "uploadedAt": "2025-01-20T09:32:00Z"
      }
    ],
    "comments": [
      {
        "id": "CMT-001",
        "author": "Sarah Johnson",
        "authorType": "Requester",
        "text": "I tried clearing my browser cache but the issue persists.",
        "timestamp": "2025-01-20T09:32:00Z",
        "isPrivate": false
      },
      {
        "id": "CMT-002",
        "author": "Michael Chen",
        "authorType": "Agent",
        "text": "Thank you for reporting this. I am investigating the issue with our IT team.",
        "timestamp": "2025-01-20T10:15:00Z",
        "isPrivate": false
      }
    ],
    "watchers": ["manager-001", "hr-head-001"],
    "relatedTickets": [],
    "csatRating": null,
    "csatComment": null
  }
}
```

---

### 1.3 Create ticket (Create Ticket dialog)
**POST** `{{BASE_URL}}/helpdesk/tickets`

This should support attachments, so use **multipart/form-data**.

#### Form-data fields
- `category` (required) — e.g. `IT`
- `subCategory` (optional)
- `priority` (required) — `Low|Medium|High|Critical`
- `title` (required)
- `description` (required)
- `channel` (optional, default `Portal`)
- `tags` (optional) — can be comma-separated or repeated fields
- `attachments[]` (optional) — file(s)

#### Sample response

```json
{
  "success": true,
  "message": "Ticket created successfully",
  "data": {
    "id": "TKT-109",
    "ticketNumber": "HD-2026-109",
    "status": "Open",
    "createdAt": "2026-01-21T10:20:00Z"
  }
}
```

---

### 1.4 Add comment (Post Comment)
**POST** `{{BASE_URL}}/helpdesk/tickets/:ticket_id/comments`

#### Body (JSON)
```json
{
  "text": "Any update on this ticket?",
  "is_private": false
}
```

#### Sample response
```json
{
  "success": true,
  "message": "Comment posted successfully",
  "data": {
    "id": "CMT-991",
    "ticket_id": "TKT-001",
    "author": "Sarah Johnson",
    "authorType": "Requester",
    "text": "Any update on this ticket?",
    "timestamp": "2026-01-21T10:25:00Z",
    "isPrivate": false
  }
}
```

---

### 1.5 Upload attachments to an existing ticket
**POST** `{{BASE_URL}}/helpdesk/tickets/:ticket_id/attachments`

#### Form-data fields
- `attachments[]` (required) — file(s)

#### Sample response
```json
{
  "success": true,
  "message": "Attachments uploaded successfully",
  "data": [
    {
      "id": "ATT-778",
      "name": "log.txt",
      "url": "https://files.company.com/helpdesk/ATT-778",
      "size": "12 KB",
      "type": "text/plain",
      "uploadedAt": "2026-01-21T10:30:00Z"
    }
  ]
}
```

---

### 1.6 Close ticket (optional)
**POST** `{{BASE_URL}}/helpdesk/tickets/:ticket_id/close`

#### Body (JSON)
```json
{
  "comment": "Issue resolved on my side. Closing ticket."
}
```

---

### 1.7 Submit CSAT rating (optional)
**POST** `{{BASE_URL}}/helpdesk/tickets/:ticket_id/csat`

#### Body (JSON)
```json
{
  "rating": 5,
  "comment": "Excellent and quick support!"
}
```

---

## 2) HR/Agent APIs — Ticket triage & resolution

These endpoints are intended for **HR/Agents (Resolvers)**.

### 2.1 List tickets in my queue (agent inbox)
**GET** `{{BASE_URL}}/helpdesk/agent/tickets`

#### Query parameters
- `page` (optional, default `1`)
- `page_size` (optional, default `20`, max `100`)
- `search` (optional)
- `status` (optional)
- `priority` (optional)
- `category` (optional)
- `assigned_to` (optional) — `me | unassigned | <agent_user_id>`
- `queue` (optional) — e.g. `HR Support`, `IT Support` (if your system uses queues)
- `sort_by` / `sort_dir` (optional)

#### Sample response
Same shape as ticket list in section 1.1, but includes `requester` always and can include `queue`.

---

### 2.2 Assign / Reassign ticket
**POST** `{{BASE_URL}}/helpdesk/agent/tickets/:ticket_id/assign`

#### Body (JSON)
```json
{
  "assignee_user_id": "AGT-001",
  "note": "Assigning to HR Support for employment letter processing"
}
```

#### Sample response
```json
{
  "success": true,
  "message": "Ticket assigned successfully",
  "data": {
    "ticket_id": "TKT-002",
    "assignedTo": { "id": "AGT-001", "name": "Michael Chen", "email": "michael.chen@company.com", "team": "HR Support" },
    "updatedAt": "2026-01-21T11:00:00Z"
  }
}
```

---

### 2.3 Update ticket status (triage workflow)
**PATCH** `{{BASE_URL}}/helpdesk/agent/tickets/:ticket_id/status`

#### Body (JSON)
```json
{
  "status": "In Progress",
  "note": "Investigating and coordinating with payroll team"
}
```

Allowed transitions (recommended):
- `Open -> In Progress | Pending | Cancelled`
- `Pending -> In Progress | Resolved`
- `In Progress -> Pending | Resolved`
- `Resolved -> Closed | Reopened(Open)`

#### Sample response
```json
{
  "success": true,
  "message": "Ticket status updated successfully",
  "data": {
    "ticket_id": "TKT-001",
    "status": "In Progress",
    "updatedAt": "2026-01-21T11:05:00Z"
  }
}
```

---

### 2.4 Resolve ticket (set resolvedAt, optional resolution summary)
**POST** `{{BASE_URL}}/helpdesk/agent/tickets/:ticket_id/resolve`

#### Body (JSON)
```json
{
  "resolution_summary": "Payslip generation job failed. Re-ran job and verified download works.",
  "internal_note": "Root cause: expired payroll service token"
}
```

---

### 2.5 Add internal note (private comment visible only to agents/admin)
**POST** `{{BASE_URL}}/helpdesk/agent/tickets/:ticket_id/internal-notes`

#### Body (JSON)
```json
{
  "text": "Waiting for user to provide updated passport photo.",
  "visibility": "internal"
}
```

---

### 2.6 Add comment as agent (public or private)
**POST** `{{BASE_URL}}/helpdesk/tickets/:ticket_id/comments`

#### Body (JSON)
```json
{
  "text": "Please bring your laptop to the IT office for inspection.",
  "is_private": false
}
```

Notes:
- Reuse the same endpoint as employee, but enforce permissions for `is_private=true` (agent/admin only).

---

## 3) Admin APIs — Configuration & full visibility

These endpoints are intended for **Helpdesk Admins/Managers**.

### 3.1 List all tickets (org-wide)
**GET** `{{BASE_URL}}/helpdesk/admin/tickets`

Same as Agent list, but unrestricted across org (with filters like `department`, `requester_id`, etc).

---

### 3.2 Configure routing (who tickets are reported to)

This is the key configuration you requested: define **rules** so that when an employee creates a ticket (or when a ticket’s category/subcategory/priority changes), it is automatically routed to the correct **queue/team/agent**.

#### 3.2.1 List routing rules
**GET** `{{BASE_URL}}/helpdesk/admin/routing-rules`

#### Sample response
```json
{
  "success": true,
  "message": "Routing rules retrieved successfully",
  "data": [
    {
      "id": "RR-001",
      "name": "Payroll → HR Support",
      "is_active": true,
      "match": {
        "category": "Payroll",
        "subCategory": null,
        "priority_in": ["High", "Critical"]
      },
      "route_to": {
        "type": "queue",
        "queue": "HR Support"
      },
      "fallback": {
        "type": "user",
        "assignee_user_id": "AGT-001"
      },
      "escalation": {
        "after_minutes_without_first_response": 60,
        "escalate_to_user_id": "AGT-010"
      },
      "createdAt": "2026-01-01T08:00:00Z",
      "updatedAt": "2026-01-21T09:00:00Z"
    }
  ]
}
```

---

#### 3.2.2 Create routing rule
**POST** `{{BASE_URL}}/helpdesk/admin/routing-rules`

#### Body (JSON)
```json
{
  "name": "IT VPN → IT Support",
  "is_active": true,
  "match": {
    "category": "IT",
    "subCategory": "VPN Access",
    "priority_in": ["Low", "Medium", "High", "Critical"]
  },
  "route_to": {
    "type": "queue",
    "queue": "IT Support"
  },
  "fallback": {
    "type": "user",
    "assignee_user_id": "AGT-003"
  },
  "escalation": {
    "after_minutes_without_first_response": 30,
    "escalate_to_user_id": "AGT-020"
  }
}
```

---

#### 3.2.3 Update routing rule
**PUT** `{{BASE_URL}}/helpdesk/admin/routing-rules/:rule_id`

---

#### 3.2.4 Enable/disable routing rule
**PATCH** `{{BASE_URL}}/helpdesk/admin/routing-rules/:rule_id`

#### Body (JSON)
```json
{ "is_active": false }
```

---

### 3.3 Ticket categories + SLA + templates (admin managed)
**GET** `{{BASE_URL}}/helpdesk/admin/ticket-categories`
**POST** `{{BASE_URL}}/helpdesk/admin/ticket-categories`
**PUT** `{{BASE_URL}}/helpdesk/admin/ticket-categories/:category_id`

Note:
- This is where you configure SLA timers (first response + resolution).
- Routing rules can reference category/subCategory/priority to drive assignment.

---

### 3.4 Ticket workflow rules (optional but recommended)
**GET** `{{BASE_URL}}/helpdesk/admin/workflow`

Use this to configure allowed transitions, auto-close after X days, etc.

---

## 4) Knowledge Base APIs

### 2.1 List KB categories (for “Browse by Category”)
**GET** `{{BASE_URL}}/helpdesk/knowledge-base/categories`

#### Sample response
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

### 2.2 List KB articles (search + featured + category)
**GET** `{{BASE_URL}}/helpdesk/knowledge-base/articles`

#### Query parameters
- `page` (optional, default `1`)
- `page_size` (optional, default `20`, max `100`)
- `search` (optional) — matches title/content/tags
- `category` (optional)
- `featured` (optional boolean)
- `status` (optional, default `Published`) — for end-user
- `sort_by` (optional) — `views | updatedAt`
- `sort_dir` (optional) — `asc | desc`

#### Sample response
```json
{
  "success": true,
  "message": "Articles retrieved successfully",
  "data": {
    "data": [
      {
        "id": "KB-001",
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

Notes:
- The UI displays title/category/tags/views/helpful%, featured flag, updatedAt, author.
- `content_preview` is recommended to reduce payload.

---

### 2.3 Get KB article details
**GET** `{{BASE_URL}}/helpdesk/knowledge-base/articles/:article_id`

#### Sample response
```json
{
  "success": true,
  "message": "Article retrieved successfully",
  "data": {
    "id": "KB-001",
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

Implementation note:
- You can increment `views` automatically on GET (recommended) or via a separate endpoint.

---

### 2.4 Article helpful feedback (Yes/No)
**POST** `{{BASE_URL}}/helpdesk/knowledge-base/articles/:article_id/feedback`

#### Body (JSON)
```json
{
  "is_helpful": true
}
```

#### Sample response
```json
{
  "success": true,
  "message": "Feedback recorded successfully",
  "data": {
    "article_id": "KB-001",
    "helpful": 99,
    "notHelpful": 12
  }
}
```

---

## 5) HelpDesk Dashboard APIs

### 3.1 Dashboard statistics (KPI + breakdowns + top agents)
**GET** `{{BASE_URL}}/helpdesk/dashboard/statistics`

#### Query parameters
- `date_range` (optional) — `today | thisWeek | thisMonth | custom`
- `from` (required if `date_range=custom`) — ISO date, e.g. `2026-01-01`
- `to` (required if `date_range=custom`) — ISO date, e.g. `2026-01-31`

#### Sample response
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
    "ticketsByCategory": {
      "IT": 412,
      "HR": 245,
      "Facilities": 98,
      "Payroll": 56,
      "Leave": 36
    },
    "ticketsByPriority": {
      "Critical": 12,
      "High": 45,
      "Medium": 234,
      "Low": 556
    },
    "topAgents": [
      { "id": "AGT-001", "name": "Michael Chen", "resolved": 156, "avgTime": "3.8h", "csat": 4.6 },
      { "id": "AGT-003", "name": "Robert Lee", "resolved": 142, "avgTime": "4.1h", "csat": 4.4 }
    ]
  }
}
```

---

### 3.2 Recent tickets (Dashboard table)
**GET** `{{BASE_URL}}/helpdesk/tickets/recent`

#### Query parameters
- `limit` (optional, default `6`, max `50`)

#### Sample response
```json
{
  "success": true,
  "message": "Recent tickets retrieved successfully",
  "data": [
    {
      "id": "TKT-001",
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

### 3.3 Export report (Dashboard “Export Report”)
**GET** `{{BASE_URL}}/helpdesk/reports/export`

#### Query parameters
- `format` (optional) — `csv | xlsx` (default `csv`)
- `date_range`, `from`, `to` — same as statistics

#### Response
- File download (binary) with `Content-Disposition: attachment; filename="helpdesk_report.csv"`

---

## 6) (Optional) Employee-facing configuration endpoint

### 6.1 Ticket categories + templates + SLA rules (for Create Ticket form)
**GET** `{{BASE_URL}}/helpdesk/ticket-categories`

#### Sample response
```json
{
  "success": true,
  "message": "Ticket categories retrieved successfully",
  "data": [
    {
      "id": "CAT-IT",
      "name": "IT",
      "description": "IT support and technical issues",
      "sla": { "firstResponse": "30", "resolution": "8" },
      "defaultPriority": "High",
      "autoAssignTo": "it-team",
      "requiredFields": ["description", "category", "subCategory"],
      "templates": [
        { "id": "TPL-IT-001", "name": "VPN Issue", "fields": { "errorCode": "", "deviceType": "", "internetProvider": "" } }
      ]
    }
  ]
}
```

---

## Notes on naming / backward compatibility

There is an existing service function in the codebase: `src/services/HelpCenterService.ts` using:
- `GET /helps/articles`
- `DELETE /helps/articles`

If the backend already uses `/helps/articles` for Knowledge Base, you can:
- Keep `/helps/articles` as an alias for `GET /helpdesk/knowledge-base/articles`
- Or update frontend service to the new `/helpdesk/knowledge-base/*` routes consistently.


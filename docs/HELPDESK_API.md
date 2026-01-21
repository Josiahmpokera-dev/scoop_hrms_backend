# Helpdesk & Support API Documentation

## Overview
This document provides comprehensive API documentation for the Helpdesk & Support system, allowing employees to create and manage tickets, agents to triage and resolve tickets, and admins to configure the system.

## Base URL
All API endpoints use the base URL: `{{BASE_URL}}/api/v1/helpdesk`

## Authentication
All requests require Bearer token authentication:
```
Authorization: Bearer <access_token>
```

**Note:** 
- Employee endpoints require authentication only
- Agent/Admin endpoints require `HRMiddleware()` (HR or Admin role)

---

## Employee (Requester) APIs

### 1. Create Ticket

**Endpoint:** `POST /api/v1/helpdesk/tickets`

**Description:** Create a new helpdesk ticket

**Content-Type:** `multipart/form-data`

**Form Fields:**
- `category` (string, required): Category - IT, HR, Payroll, Facilities, Leave, General
- `sub_category` (string, optional): Sub-category
- `priority` (string, required): Priority - Low, Medium, High, Critical
- `title` (string, required): Ticket title
- `description` (string, required): Detailed description
- `channel` (string, optional): Channel - Portal, Email, WhatsApp, Phone, Chat (default: Portal)
- `tags` (string, optional): Comma-separated tags

**Success Response (200):**
```json
{
  "success": true,
  "message": "Ticket created successfully",
  "data": {
    "id": 1,
    "ticket_number": "HD-2026-001",
    "status": "Open",
    "created_at": "2026-01-20T10:30:00+03:00"
  }
}
```

**Example Request (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/helpdesk/tickets" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "category=IT" \
  -F "sub_category=VPN Access" \
  -F "priority=High" \
  -F "title=Unable to access VPN" \
  -F "description=I cannot connect to the company VPN from home." \
  -F "tags=vpn,access,remote"
```

---

### 2. List My Tickets

**Endpoint:** `GET /api/v1/helpdesk/tickets?page=1&page_size=20&search=VPN&status=Open&priority=High&category=IT`

**Description:** Retrieve a paginated list of tickets for the authenticated employee

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Items per page (default: 20, max: 100)
- `search` (string, optional): Search term (matches title, description, ticket number)
- `status` (string, optional): Filter by status - Open, In Progress, Pending, Resolved, Closed, Cancelled
- `priority` (string, optional): Filter by priority - Low, Medium, High, Critical
- `category` (string, optional): Filter by category
- `sort_by` (string, optional): Sort field - createdAt, updatedAt, dueDate
- `sort_dir` (string, optional): Sort direction - asc, desc

**Success Response (200):**
```json
{
  "success": true,
  "message": "Tickets retrieved successfully",
  "data": [
    {
      "id": 1,
      "ticket_number": "HD-2026-001",
      "title": "Unable to access VPN",
      "description": "I cannot connect to the company VPN from home.",
      "category": "IT",
      "sub_category": "VPN Access",
      "priority": "High",
      "status": "Open",
      "channel": "Portal",
      "requester": {
        "id": "EMP001"
      },
      "assigned_to": {
        "name": "Michael Chen",
        "email": "michael.chen@company.com",
        "team": "IT Support"
      },
      "created_at": "2026-01-20T10:30:00+03:00",
      "updated_at": "2026-01-20T10:30:00+03:00",
      "due_date": "2026-01-22T10:30:00+03:00",
      "sla_status": "Within SLA",
      "tags": ["vpn", "access", "remote"]
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

### 3. Get Ticket Details

**Endpoint:** `GET /api/v1/helpdesk/tickets/:ticket_id`

**Description:** Retrieve detailed information about a specific ticket, including comments and attachments

**Path Parameters:**
- `ticket_id` (integer, required): Ticket ID

**Success Response (200):**
```json
{
  "success": true,
  "message": "Ticket details retrieved successfully",
  "data": {
    "id": 1,
    "ticket_number": "HD-2026-001",
    "title": "Unable to access VPN",
    "description": "I cannot connect to the company VPN from home.",
    "category": "IT",
    "priority": "High",
    "status": "In Progress",
    "comments": [
      {
        "id": 1,
        "author": "Sarah Johnson",
        "author_type": "Requester",
        "text": "I tried clearing my browser cache but the issue persists.",
        "timestamp": "2026-01-20T09:32:00+03:00",
        "is_private": false
      },
      {
        "id": 2,
        "author": "Michael Chen",
        "author_type": "Agent",
        "text": "Thank you for reporting this. I am investigating the issue.",
        "timestamp": "2026-01-20T10:15:00+03:00",
        "is_private": false
      }
    ],
    "comments_count": 2,
    "attachments": [
      {
        "id": 1,
        "name": "error-screenshot.png",
        "url": "https://files.company.com/helpdesk/ATT-001",
        "size": 245760,
        "type": "image/png",
        "uploaded_at": "2026-01-20T09:32:00+03:00"
      }
    ]
  }
}
```

---

### 4. Add Comment

**Endpoint:** `POST /api/v1/helpdesk/tickets/:ticket_id/comments`

**Description:** Add a comment to a ticket

**Request Body:**
```json
{
  "text": "Any update on this ticket?",
  "is_private": false
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Comment posted successfully",
  "data": {
    "id": 3,
    "ticket_id": 1,
    "author": "Sarah Johnson",
    "author_type": "Requester",
    "text": "Any update on this ticket?",
    "timestamp": "2026-01-21T10:25:00+03:00",
    "is_private": false
  }
}
```

**Note:** Employees cannot create private comments. Only agents/admin can create private comments.

---

### 5. Close Ticket

**Endpoint:** `POST /api/v1/helpdesk/tickets/:ticket_id/close`

**Description:** Close a ticket (employee can close their own tickets)

**Request Body:**
```json
{
  "comment": "Issue resolved on my side. Closing ticket."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Ticket closed successfully",
  "data": null
}
```

---

### 6. Submit CSAT

**Endpoint:** `POST /api/v1/helpdesk/tickets/:ticket_id/csat`

**Description:** Submit customer satisfaction rating for a resolved/closed ticket

**Request Body:**
```json
{
  "rating": 5,
  "comment": "Excellent and quick support!"
}
```

**Request Parameters:**
- `rating` (integer, required): Rating from 1 to 5
- `comment` (string, optional): Optional feedback comment

**Success Response (200):**
```json
{
  "success": true,
  "message": "CSAT submitted successfully",
  "data": null
}
```

---

## Agent/Admin APIs

### 7. List All Tickets (Agent/Admin)

**Endpoint:** `GET /api/v1/helpdesk/agent/tickets?page=1&page_size=20&assigned_to=me&status=Open`

**Description:** Retrieve a paginated list of all tickets for agent/admin review

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Items per page (default: 20, max: 100)
- `search` (string, optional): Search term
- `status` (string, optional): Filter by status
- `priority` (string, optional): Filter by priority
- `category` (string, optional): Filter by category
- `assigned_to` (string, optional): Filter by assignment - `me`, `unassigned`, or user ID
- `queue` (string, optional): Filter by queue name

**Success Response (200):**
Same format as employee ticket list, but includes all tickets in the system.

---

### 8. Assign Ticket

**Endpoint:** `POST /api/v1/helpdesk/agent/tickets/:ticket_id/assign`

**Description:** Assign or reassign a ticket to an agent

**Request Body:**
```json
{
  "assignee_user_id": 5,
  "note": "Assigning to IT Support for VPN issue resolution"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Ticket assigned successfully",
  "data": {
    "ticket_id": 1,
    "assigned_to": {
      "id": 5,
      "name": "Michael Chen",
      "email": "michael.chen@company.com",
      "team": "IT Support"
    },
    "updated_at": "2026-01-21T11:00:00+03:00"
  }
}
```

---

### 9. Update Ticket Status

**Endpoint:** `PATCH /api/v1/helpdesk/agent/tickets/:ticket_id/status`

**Description:** Update the status of a ticket

**Request Body:**
```json
{
  "status": "In Progress",
  "note": "Investigating and coordinating with IT team"
}
```

**Allowed Status Transitions:**
- `Open` → `In Progress`, `Pending`, `Cancelled`
- `Pending` → `In Progress`, `Resolved`
- `In Progress` → `Pending`, `Resolved`
- `Resolved` → `Closed`, `Open` (reopen)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Ticket status updated successfully",
  "data": {
    "ticket_id": 1,
    "status": "In Progress",
    "updated_at": "2026-01-21T11:05:00+03:00"
  }
}
```

---

### 10. Resolve Ticket

**Endpoint:** `POST /api/v1/helpdesk/agent/tickets/:ticket_id/resolve`

**Description:** Resolve a ticket with resolution summary

**Request Body:**
```json
{
  "resolution_summary": "VPN configuration updated. User can now connect successfully.",
  "internal_note": "Root cause: expired VPN certificate"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Ticket resolved successfully",
  "data": {
    "ticket_id": 1,
    "status": "Resolved",
    "resolved_at": "2026-01-22T14:30:00+03:00"
  }
}
```

---

### 11. Get Helpdesk Statistics

**Endpoint:** `GET /api/v1/helpdesk/dashboard/statistics?date_range=thisMonth&from=2026-01-01&to=2026-01-31`

**Description:** Retrieve helpdesk statistics and KPIs for dashboard

**Query Parameters:**
- `date_range` (string, optional): Date range - `today`, `thisWeek`, `thisMonth`, `custom`
- `from` (string, optional): From date (YYYY-MM-DD) - required if `date_range=custom`
- `to` (string, optional): To date (YYYY-MM-DD) - required if `date_range=custom`

**Success Response (200):**
```json
{
  "success": true,
  "message": "Statistics retrieved successfully",
  "data": {
    "total_tickets": 847,
    "tickets_by_status": {
      "Open": 124,
      "In Progress": 56,
      "Resolved": 589,
      "Closed": 78
    },
    "tickets_by_category": {
      "IT": 412,
      "HR": 245,
      "Facilities": 98,
      "Payroll": 56,
      "Leave": 36
    },
    "tickets_by_priority": {
      "Critical": 12,
      "High": 45,
      "Medium": 234,
      "Low": 556
    }
  }
}
```

---

### 12. Get Recent Tickets

**Endpoint:** `GET /api/v1/helpdesk/tickets/recent?limit=6`

**Description:** Retrieve recent tickets for dashboard

**Query Parameters:**
- `limit` (integer, optional): Number of tickets to return (default: 6, max: 50)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Recent tickets retrieved successfully",
  "data": [
    {
      "id": 1,
      "ticket_number": "HD-2026-001",
      "title": "Unable to access VPN",
      "category": "IT",
      "priority": "High",
      "status": "In Progress",
      "created_at": "2026-01-20T10:30:00+03:00"
    }
  ]
}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "category, priority, title, and description are required",
  "error": null
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "User not authenticated",
  "error": "Unauthorized access"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "message": "HR or Admin access required",
  "error": "Forbidden"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "ticket not found",
  "error": "Resource not found"
}
```

---

## Notes

1. **Authorization:**
   - Employee endpoints: Require authentication only
   - Agent/Admin endpoints: Require HR or Admin role

2. **Ticket Status Flow:**
   - `Open` → `In Progress` / `Pending` / `Cancelled`
   - `Pending` → `In Progress` / `Resolved`
   - `In Progress` → `Pending` / `Resolved`
   - `Resolved` → `Closed` / `Open` (reopen)

3. **Private Comments:**
   - Only agents/admin can create private comments
   - Private comments are visible only to agents/admin, not to the requester

4. **SLA Tracking:**
   - First response time is automatically tracked when an agent/admin adds the first comment
   - Resolution time is tracked when ticket is resolved

5. **CSAT:**
   - Can only be submitted for resolved or closed tickets
   - Rating must be between 1 and 5

6. **Pagination:**
   - All list endpoints support pagination with `page` and `page_size` parameters

7. **Tenant Isolation:**
   - All operations are scoped to the authenticated user's tenant

---

## Implementation Status

✅ **Completed:**
- Employee ticket creation
- Employee ticket listing and details
- Employee comments
- Employee ticket closure
- Employee CSAT submission
- Agent ticket listing
- Agent ticket assignment
- Agent ticket status updates
- Agent ticket resolution
- Dashboard statistics
- Recent tickets

⏳ **To Be Implemented:**
- File attachments upload
- Knowledge Base APIs
- Routing rules management
- Ticket categories management
- Internal notes (separate from comments)
- Ticket watchers
- Related tickets

---

## Next Steps

1. **Restart your application** to run migrations:
   ```bash
   go run cmd/api/main.go
   ```

2. **The following tables will be created:**
   - `helpdesk_tickets` - Main tickets table
   - `helpdesk_comments` - Ticket comments
   - `helpdesk_attachments` - Ticket attachments
   - `helpdesk_routing_rules` - Routing rules
   - `helpdesk_kb_articles` - Knowledge base articles
   - `helpdesk_kb_feedback` - KB article feedback
   - `helpdesk_ticket_categories` - Ticket categories

3. **Test the endpoints** using the examples provided above.

All endpoints follow your existing app patterns and logic. No existing functionality has been modified or destroyed.

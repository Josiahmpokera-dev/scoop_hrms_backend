# Helpdesk Dashboard API Documentation

**Base URL:** `http://localhost:8080/api/v1`

This document outlines the backend API endpoints for the Helpdesk Dashboard module.

---

## Table of Contents

1. [Dashboard Statistics](#1-dashboard-statistics)
2. [Recent Tickets](#2-recent-tickets)
3. [Export Report](#3-export-report)

---

## 1. Dashboard Statistics

Retrieves aggregated statistics for the helpdesk dashboard including ticket counts, resolution times, CSAT scores, SLA compliance, and breakdowns by category/priority/agents.

### Endpoint

```
GET /api/v1/helpdesk/dashboard/statistics
```

### Headers

| Header          | Value              | Required |
| --------------- | ------------------ | -------- |
| `Authorization` | `Bearer <token>`   | Yes      |

### Authorization

Requires **HR**, **IT**, or **Admin** role.

### Query Parameters

| Parameter    | Type   | Required | Description                                                                                                      |
| ------------ | ------ | -------- | ---------------------------------------------------------------------------------------------------------------- |
| `date_range` | string | No       | Predefined range: `today`, `yesterday`, `thisWeek`, `lastWeek`, `thisMonth`, `lastMonth`, `thisYear`, `lastYear` |
| `from`       | string | No       | Start date (ISO 8601 format, e.g. `2026-01-01`). Used with `to` for custom range.                                |
| `to`         | string | No       | End date (ISO 8601 format, e.g. `2026-01-31`). Used with `from` for custom range.                                |

> **Note:** If `date_range` is provided, it takes precedence over `from`/`to`. If neither is provided, all tickets are included.

### Sample Requests

**This month:**
```bash
curl -X GET "http://localhost:8080/api/v1/helpdesk/dashboard/statistics?date_range=thisMonth" \
  -H "Authorization: Bearer <token>"
```

**Custom date range:**
```bash
curl -X GET "http://localhost:8080/api/v1/helpdesk/dashboard/statistics?from=2026-01-01&to=2026-01-31" \
  -H "Authorization: Bearer <token>"
```

### Success Response

**Status Code:** `200 OK`

```json
{
  "success": true,
  "message": "Dashboard statistics retrieved successfully",
  "data": {
    "totalTickets": 1247,
    "openTickets": 89,
    "inProgress": 156,
    "resolved": 892,
    "closed": 98,
    "avgFirstResponseTime": "2.4h",
    "avgResolutionTime": "18.5h",
    "slaCompliance": 94.2,
    "csatScore": 4.6,
    "ticketsByCategory": {
      "IT Support": 423,
      "HR": 312,
      "Facilities": 189,
      "Finance": 156,
      "Other": 167
    },
    "ticketsByPriority": {
      "Critical": 45,
      "High": 198,
      "Medium": 567,
      "Low": 437
    },
    "topAgents": [
      {
        "id": 1,
        "name": "John Smith",
        "resolved": 156,
        "avgTime": "12.3h",
        "csat": 4.8
      },
      {
        "id": 2,
        "name": "Sarah Johnson",
        "resolved": 142,
        "avgTime": "14.1h",
        "csat": 4.7
      }
    ]
  }
}
```

### Response Fields

| Field                  | Type   | Description                                             |
| ---------------------- | ------ | ------------------------------------------------------- |
| `totalTickets`         | number | Total number of tickets in the selected period          |
| `openTickets`          | number | Count of tickets with status "Open"                     |
| `inProgress`           | number | Count of tickets with status "In Progress"              |
| `resolved`             | number | Count of tickets with status "Resolved"                 |
| `closed`               | number | Count of tickets with status "Closed"                   |
| `avgFirstResponseTime` | string | Average time to first response (e.g. "2.4h", "1d 3h")   |
| `avgResolutionTime`    | string | Average time to resolve a ticket (e.g. "18.5h", "2d")   |
| `slaCompliance`        | number | Percentage of tickets resolved within SLA (0-100)       |
| `csatScore`            | number | Customer Satisfaction score (1.0-5.0)                   |
| `ticketsByCategory`    | object | Key-value pairs of category name → ticket count         |
| `ticketsByPriority`    | object | Key-value pairs of priority level → ticket count        |
| `topAgents`            | array  | List of top-performing agents (max 5)                   |
| `topAgents[].id`       | number | Agent's user ID                                         |
| `topAgents[].name`     | string | Agent's display name                                    |
| `topAgents[].resolved` | number | Number of tickets resolved by this agent                |
| `topAgents[].avgTime`  | string | Agent's average resolution time                         |
| `topAgents[].csat`     | number | Agent's average CSAT score (1.0-5.0)                    |

### Time Format

- Under 1 hour: `"45m"`
- Under 24 hours: `"2.4h"` or `"12h"`
- Over 24 hours: `"1d 3h"` or `"2d"`

### Error Responses

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Invalid or expired token"
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "You do not have permission to view dashboard statistics"
}
```

---

## 2. Recent Tickets

Retrieves a list of the most recently created tickets for the dashboard "Recent Tickets" section.

### Endpoint

```
GET /api/v1/helpdesk/tickets/recent
```

### Headers

| Header          | Value              | Required |
| --------------- | ------------------ | -------- |
| `Authorization` | `Bearer <token>`   | Yes      |

### Authorization

Requires authentication (any logged-in user).

### Query Parameters

| Parameter | Type   | Required | Default | Description                           |
| --------- | ------ | -------- | ------- | ------------------------------------- |
| `limit`   | number | No       | 10      | Number of recent tickets to return (max 50) |

### Sample Request

```bash
curl -X GET "http://localhost:8080/api/v1/helpdesk/tickets/recent?limit=10" \
  -H "Authorization: Bearer <token>"
```

### Success Response

**Status Code:** `200 OK`

```json
{
  "success": true,
  "message": "Recent tickets retrieved successfully",
  "data": [
    {
      "id": 1234,
      "ticketNumber": "HD-2026-001",
      "title": "Unable to access email on mobile device",
      "requester": {
        "name": "Alice Thompson",
        "department": ""
      },
      "category": "IT Support",
      "priority": "High",
      "status": "In Progress",
      "assignedTo": {
        "name": "John Smith"
      },
      "createdAt": "2026-01-23T09:30:00Z"
    },
    {
      "id": 1233,
      "ticketNumber": "HD-2026-002",
      "title": "Request for new laptop",
      "requester": {
        "name": "Bob Williams",
        "department": ""
      },
      "category": "IT Support",
      "priority": "Medium",
      "status": "Open",
      "assignedTo": null,
      "createdAt": "2026-01-23T09:15:00Z"
    }
  ]
}
```

### Response Fields

| Field                  | Type           | Description                                      |
| ---------------------- | -------------- | ------------------------------------------------ |
| `id`                   | number         | Unique ticket identifier                         |
| `ticketNumber`         | string         | Human-readable ticket reference (e.g. HD-2026-001) |
| `title`                | string         | Ticket subject/title                             |
| `requester`            | object         | Information about the ticket requester           |
| `requester.name`       | string         | Requester's full name                            |
| `requester.department` | string         | Requester's department (empty if not available)  |
| `category`             | string         | Ticket category (e.g. "IT Support", "HR")        |
| `priority`             | string         | Priority level: `Critical`, `High`, `Medium`, `Low` |
| `status`               | string         | Status: `Open`, `In Progress`, `Pending`, `Resolved`, `Closed`, `Cancelled` |
| `assignedTo`           | object \| null | Information about the assigned agent             |
| `assignedTo.name`      | string         | Assigned agent's name                            |
| `createdAt`            | string         | ISO 8601 timestamp of ticket creation            |

---

## 3. Export Report

Exports helpdesk report data as a downloadable file (CSV or XLSX format).

### Endpoint

```
GET /api/v1/helpdesk/reports/export
```

### Headers

| Header          | Value              | Required |
| --------------- | ------------------ | -------- |
| `Authorization` | `Bearer <token>`   | Yes      |

### Authorization

Requires **HR** or **Admin** role.

### Query Parameters

| Parameter    | Type   | Required | Default | Description                                                                                                      |
| ------------ | ------ | -------- | ------- | ---------------------------------------------------------------------------------------------------------------- |
| `format`     | string | No       | `xlsx`  | Export format: `csv` or `xlsx`                                                                                   |
| `date_range` | string | No       | —       | Predefined range: `today`, `yesterday`, `thisWeek`, `lastWeek`, `thisMonth`, `lastMonth`, `thisYear`, `lastYear` |
| `from`       | string | No       | —       | Start date (ISO 8601 format, e.g. `2026-01-01`)                                                                  |
| `to`         | string | No       | —       | End date (ISO 8601 format, e.g. `2026-01-31`)                                                                    |

### Sample Requests

**XLSX format (default):**
```bash
curl -X GET "http://localhost:8080/api/v1/helpdesk/reports/export?date_range=thisMonth" \
  -H "Authorization: Bearer <token>" \
  -o "helpdesk-report.xlsx"
```

**CSV format:**
```bash
curl -X GET "http://localhost:8080/api/v1/helpdesk/reports/export?format=csv&date_range=thisMonth" \
  -H "Authorization: Bearer <token>" \
  -o "helpdesk-report.csv"
```

### Success Response

**Status Code:** `200 OK`

**Content-Type (XLSX):** `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`  
**Content-Type (CSV):** `text/csv`

**Response Headers:**
```
Content-Disposition: attachment; filename="helpdesk-report-2026-01-23.xlsx"
```

**Response Body:** Binary file data (XLSX or CSV)

### XLSX Report Contents (5 Sheets)

#### Sheet 1: Summary
| Metric              | Value  |
| ------------------- | ------ |
| Total Tickets       | 1247   |
| Open Tickets        | 89     |
| In Progress         | 156    |
| Resolved            | 892    |
| Closed              | 98     |
| Avg First Response  | 2.4h   |
| Avg Resolution Time | 18.5h  |
| SLA Compliance      | 94.2%  |
| CSAT Score          | 4.6    |

#### Sheet 2: Tickets Detail
| Ticket #        | Title                     | Requester      | Department | Category   | Priority | Status      | Assigned To | Created At          | Resolved At         |
| --------------- | ------------------------- | -------------- | ---------- | ---------- | -------- | ----------- | ----------- | ------------------- | ------------------- |
| HD-2026-001     | Unable to access email... | Alice Thompson |            | IT Support | High     | In Progress | John Smith  | 2026-01-23 09:30:00 |                     |
| HD-2026-002     | Request for new laptop    | Bob Williams   |            | IT Support | Medium   | Open        |             | 2026-01-23 09:15:00 |                     |

#### Sheet 3: By Category
| Category    | Count | Percentage |
| ----------- | ----- | ---------- |
| IT Support  | 423   | 33.9%      |
| HR          | 312   | 25.0%      |
| Facilities  | 189   | 15.2%      |

#### Sheet 4: By Priority
| Priority | Count | Percentage |
| -------- | ----- | ---------- |
| Critical | 45    | 3.6%       |
| High     | 198   | 15.9%      |
| Medium   | 567   | 45.5%      |
| Low      | 437   | 35.0%      |

#### Sheet 5: Agent Performance
| Agent          | Tickets Resolved | Avg Resolution Time | CSAT Score |
| -------------- | ---------------- | ------------------- | ---------- |
| John Smith     | 156              | 12.3h               | 4.8        |
| Sarah Johnson  | 142              | 14.1h               | 4.7        |

### CSV Report Contents

CSV exports the "Tickets Detail" data only:

```csv
Ticket #,Title,Requester,Department,Category,Priority,Status,Assigned To,Created At,Resolved At
HD-2026-001,Unable to access email...,Alice Thompson,,IT Support,High,In Progress,John Smith,2026-01-23 09:30:00,
HD-2026-002,Request for new laptop,Bob Williams,,IT Support,Medium,Open,,2026-01-23 09:15:00,
```

### Error Responses

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Invalid or expired token"
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "You do not have permission to export reports"
}
```

---

## Summary of Endpoints

| Method | Endpoint                              | Auth Required | Role Required    | Description                        |
| ------ | ------------------------------------- | ------------- | ---------------- | ---------------------------------- |
| GET    | `/api/v1/helpdesk/dashboard/statistics` | Yes           | HR, IT, Admin    | Get dashboard statistics           |
| GET    | `/api/v1/helpdesk/tickets/recent`       | Yes           | Any              | Get recent tickets list            |
| GET    | `/api/v1/helpdesk/reports/export`       | Yes           | HR, Admin        | Export helpdesk report (CSV/XLSX)  |

---

## Implementation Notes

1. **Authentication:** All endpoints require a valid JWT Bearer token.

2. **Date Range Handling:**
   - If `date_range` is provided, it takes precedence over `from`/`to`.
   - If neither is provided, all tickets are included (no date filter).
   - Supported ranges: `today`, `yesterday`, `thisWeek`, `lastWeek`, `thisMonth`, `lastMonth`, `thisYear`, `lastYear`.

3. **Top Agents:** Limited to 5 agents sorted by resolved ticket count (descending).

4. **Time Formatting:** Resolution times are human-readable:
   - Under 1 hour: `"45m"`
   - Under 24 hours: `"2.4h"`
   - Over 24 hours: `"1d 3h"` or `"2d"`

5. **CSAT Score:** Calculated from CSAT submissions (1-5 scale) for resolved/closed tickets.

6. **SLA Compliance:** Percentage of resolved/closed tickets NOT marked as "Breached" in `sla_status`.

7. **Export:**
   - XLSX: 5 sheets (Summary, Tickets Detail, By Category, By Priority, Agent Performance).
   - CSV: Tickets detail only.
   - Uses `excelize/v2` library for XLSX generation.

---

## TypeScript Types Reference

```typescript
// Dashboard Statistics Response
export type DashboardStatistics = {
  totalTickets: number;
  openTickets: number;
  inProgress: number;
  resolved: number;
  closed: number;
  avgFirstResponseTime: string;
  avgResolutionTime: string;
  slaCompliance: number;
  csatScore: number;
  ticketsByCategory: Record<string, number>;
  ticketsByPriority: Record<string, number>;
  topAgents: {
    id: number;
    name: string;
    resolved: number;
    avgTime: string;
    csat: number;
  }[];
};

// Recent Ticket Item
export type RecentTicketItem = {
  id: number;
  ticketNumber: string;
  title: string;
  requester: {
    name: string;
    department?: string;
  };
  category: string;
  priority: string;
  status: string;
  assignedTo?: {
    name: string;
  } | null;
  createdAt: string;
};

// API Response wrapper
export type APIResponse<T> = {
  success: boolean;
  message: string;
  data: T;
};
```

---

*Document Version: 2.0*  
*Last Updated: January 23, 2026*

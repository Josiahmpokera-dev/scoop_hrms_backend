# Employee Dashboard API Documentation

Base URL: `/api/v1/dashboard/employee`

> **Authentication**: All endpoints require a valid JWT token via `Authorization: Bearer <token>`.
> These endpoints return data specific to the **logged-in employee** only.

---

## Table of Contents

1. [Get Employee Statistics](#1-get-employee-statistics)
2. [Get Employee Recent Activity](#2-get-employee-recent-activity)

---

## 1. Get Employee Statistics

Returns the personal dashboard KPIs for the logged-in employee: leave balance, hours worked this week, pending requests, next payday, and today's attendance status.

### Request

```
GET /api/v1/dashboard/employee/statistics
```

| Query Param | Type   | Required | Default | Description                              |
|-------------|--------|----------|---------|------------------------------------------|
| `date`      | string | No       | today   | Date for statistics (format: `YYYY-MM-DD`) |

### Headers

| Header          | Value              |
|-----------------|--------------------|
| `Authorization` | `Bearer <token>`   |
| `Content-Type`  | `application/json` |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Employee dashboard statistics retrieved successfully",
  "data": {
    "view": "employee",
    "date": "2026-02-06",
    "employeeId": 42,
    "employeeName": "Kevin Massawe",
    "stats": {
      "leaveBalance": {
        "value": 12,
        "unit": "days",
        "breakdown": {
          "annual": 8,
          "sick": 3,
          "personal": 1
        },
        "usedThisYear": 10
      },
      "hoursThisWeek": {
        "value": 38.5,
        "target": 40,
        "overtime": 0,
        "status": "on_track"
      },
      "pendingRequests": {
        "value": 2,
        "breakdown": {
          "leaveRequests": 1,
          "letterRequests": 1
        }
      },
      "nextPayday": {
        "date": "2026-02-28",
        "daysRemaining": 22,
        "estimatedNetPay": 0,
        "currency": "TZS"
      }
    },
    "attendance": {
      "todayStatus": "Not Clocked In",
      "checkInTime": null,
      "checkOutTime": null,
      "workedHoursToday": 0
    }
  }
}
```

### Response Fields

| Field                              | Type    | Description                                          |
|------------------------------------|---------|------------------------------------------------------|
| `view`                             | string  | Always `"employee"`                                  |
| `date`                             | string  | The date the statistics are computed for             |
| `employeeId`                       | integer | Employee record ID                                   |
| `employeeName`                     | string  | Full name of the employee                            |
| `stats.leaveBalance.value`         | number  | Total remaining leave days across all types          |
| `stats.leaveBalance.unit`          | string  | Always `"days"`                                      |
| `stats.leaveBalance.breakdown`     | object  | Leave balance per type (annual, sick, personal)      |
| `stats.leaveBalance.usedThisYear`  | number  | Total leave days used this year                      |
| `stats.hoursThisWeek.value`        | number  | Hours worked this week                               |
| `stats.hoursThisWeek.target`       | number  | Target hours per week (default 40)                   |
| `stats.hoursThisWeek.overtime`     | number  | Overtime hours this week                             |
| `stats.hoursThisWeek.status`       | string  | `"on_track"`, `"behind"`, or `"ahead"`               |
| `stats.pendingRequests.value`      | integer | Total count of pending requests                      |
| `stats.pendingRequests.breakdown`  | object  | Breakdown by request type                            |
| `stats.nextPayday.date`            | string  | Next payday date (`YYYY-MM-DD`)                      |
| `stats.nextPayday.daysRemaining`   | integer | Days until next payday                               |
| `stats.nextPayday.estimatedNetPay` | number  | Estimated net pay (0 if not calculated yet)          |
| `stats.nextPayday.currency`        | string  | Currency code (e.g. `"TZS"`)                         |
| `attendance.todayStatus`           | string  | Today's attendance status                            |
| `attendance.checkInTime`           | string? | Check-in time (null if not clocked in)               |
| `attendance.checkOutTime`          | string? | Check-out time (null if not clocked out)             |
| `attendance.workedHoursToday`      | number  | Hours worked today                                   |

### Error Responses

| Status | Message                                    |
|--------|--------------------------------------------|
| 401    | `User not authenticated`                   |
| 500    | `Failed to retrieve employee statistics`   |

---

## 2. Get Employee Recent Activity

Returns the logged-in employee's recent activity feed — leave requests, timesheet updates, payslip views, and service requests.

### Request

```
GET /api/v1/dashboard/employee/my-activity
```

| Query Param     | Type   | Required | Default | Description                                                      |
|-----------------|--------|----------|---------|------------------------------------------------------------------|
| `page`          | int    | No       | 1       | Page number                                                      |
| `page_size`     | int    | No       | 10      | Items per page (max 50)                                          |
| `activity_type` | string | No       | all     | Filter by type: `leave`, `timesheet`, `payslip`, `request`, `all` |
| `days`          | int    | No       | 30      | Return activities from the last N days (max 365)                 |

### Headers

| Header          | Value              |
|-----------------|--------------------|
| `Authorization` | `Bearer <token>`   |
| `Content-Type`  | `application/json` |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Employee recent activity retrieved successfully",
  "data": {
    "activities": [
      {
        "id": 101,
        "type": "leave_request",
        "title": "Leave Request",
        "description": "Annual leave for Dec 15 - Dec 20",
        "status": "pending",
        "statusLabel": "Pending approval",
        "icon": "calendar",
        "color": "orange",
        "createdAt": "2026-02-04T09:30:00Z",
        "relativeTime": "2 days ago",
        "metadata": {
          "leaveType": "Annual",
          "startDate": "2026-02-15",
          "endDate": "2026-02-20",
          "days": 4
        },
        "actionUrl": "/leave/requests/101"
      },
      {
        "id": 55,
        "type": "timesheet_update",
        "title": "Timesheet Updated",
        "description": "38.5 hours logged this week",
        "status": "completed",
        "statusLabel": "Completed",
        "icon": "clock",
        "color": "green",
        "createdAt": "2026-02-05T17:00:00Z",
        "relativeTime": "1 day ago",
        "metadata": {
          "totalHours": 38.5,
          "weekEnding": "2026-02-07"
        },
        "actionUrl": "/timesheets/55"
      }
    ],
    "summary": {
      "pendingRequests": 2,
      "recentApprovals": 1,
      "newNotifications": 3
    },
    "meta": {
      "page": 1,
      "pageSize": 10,
      "total": 15,
      "totalPages": 2
    }
  }
}
```

### Activity Object Fields

| Field          | Type   | Description                                                   |
|----------------|--------|---------------------------------------------------------------|
| `id`           | integer | Activity record ID                                           |
| `type`         | string | Activity type: `leave_request`, `timesheet_update`, `payslip_generated`, `service_request` |
| `title`        | string | Human-readable title                                          |
| `description`  | string | Brief description of the activity                             |
| `status`       | string | Status: `pending`, `approved`, `rejected`, `completed`        |
| `statusLabel`  | string | User-friendly status label                                    |
| `icon`         | string | Icon name for frontend (e.g. `calendar`, `clock`, `dollar`)   |
| `color`        | string | Color hint: `orange` (pending), `green` (approved/completed), `red` (rejected), `blue` (info) |
| `createdAt`    | string | ISO 8601 timestamp                                            |
| `relativeTime` | string | Human-readable relative time (e.g. `"2 days ago"`)            |
| `metadata`     | object | Type-specific extra data (varies per activity type)           |
| `actionUrl`    | string | Frontend route to view full details                           |

### Summary Object Fields

| Field              | Type    | Description                          |
|--------------------|---------|--------------------------------------|
| `pendingRequests`  | integer | Count of pending requests            |
| `recentApprovals`  | integer | Count of recently approved requests  |
| `newNotifications` | integer | Count of unread notifications        |

### Meta Object Fields

| Field        | Type    | Description          |
|--------------|---------|----------------------|
| `page`       | integer | Current page number  |
| `pageSize`   | integer | Items per page       |
| `total`      | integer | Total items          |
| `totalPages` | integer | Total pages          |

### Error Responses

| Status | Message                                    |
|--------|--------------------------------------------|
| 401    | `User not authenticated`                   |
| 500    | `Failed to retrieve employee activity`     |

---

## Quick Reference

| Method | Endpoint                                     | Description                                |
|--------|----------------------------------------------|--------------------------------------------|
| GET    | `/api/v1/dashboard/employee/statistics`      | Personal KPI stats (leave, hours, payday)  |
| GET    | `/api/v1/dashboard/employee/my-activity`     | Recent activity feed (leave, timesheets…)  |

---

## Frontend Integration Example

```javascript
// Fetch employee dashboard statistics
const statsResponse = await fetch('/api/v1/dashboard/employee/statistics', {
  headers: { 'Authorization': `Bearer ${token}` }
});
const stats = await statsResponse.json();
// stats.data.stats.leaveBalance.value  → 12
// stats.data.stats.hoursThisWeek.value → 38.5
// stats.data.stats.pendingRequests.value → 2
// stats.data.stats.nextPayday.daysRemaining → 22

// Fetch employee recent activity
const activityResponse = await fetch('/api/v1/dashboard/employee/my-activity?page=1&page_size=5&days=30', {
  headers: { 'Authorization': `Bearer ${token}` }
});
const activity = await activityResponse.json();
// activity.data.activities → array of recent activity items
// activity.data.summary.pendingRequests → 2
```

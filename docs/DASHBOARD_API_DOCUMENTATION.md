# Dashboard API Documentation

## Overview

The main Dashboard provides real-time statistics, announcements, events, and activity feeds for both **Admin** and **Employee** views. This document specifies all API endpoints required for the Dashboard module.

**Base URL:** `{{BASE_URL}}/api/v1`

**Authentication:** All endpoints require Bearer token authentication.

---

## Table of Contents

1. [Dashboard Statistics](#1-dashboard-statistics-api)
2. [Announcements](#2-announcements-api)
3. [Upcoming Events](#3-upcoming-events-api-admin)
4. [Recent Activity](#4-recent-activity-api-employee)
5. [Quick Actions](#5-quick-actions-api)
6. [Pending Approvals](#6-pending-approvals-api-admin)
7. [TypeScript Types](#7-typescript-types-reference)
8. [Error Handling](#8-error-handling)

---

## 1. Dashboard Statistics API

### `GET /dashboard/statistics`

Returns KPI statistics based on user role. Admins get organization-wide stats, employees get personal stats.

**Headers:**
```
Authorization: Bearer {{ACCESS_TOKEN}}
Content-Type: application/json
```

**Query Parameters:**

| Parameter | Type   | Required | Description                                              |
|-----------|--------|----------|----------------------------------------------------------|
| `view`    | string | No       | `admin` or `employee` (defaults based on user role)      |
| `date`    | string | No       | Date for statistics (default: today, format: `YYYY-MM-DD`) |

---

### Response (Admin View)

```json
{
  "success": true,
  "message": "Dashboard statistics retrieved successfully",
  "data": {
    "view": "admin",
    "date": "2026-01-23",
    "stats": {
      "totalHeadcount": {
        "value": 156,
        "change": 3,
        "changePercent": 1.96,
        "trend": "up"
      },
      "presentToday": {
        "value": 142,
        "total": 156,
        "percentage": 91.03,
        "absent": 8,
        "onLeave": 6
      },
      "pendingApprovals": {
        "value": 8,
        "breakdown": {
          "leaveRequests": 5,
          "timesheetCorrections": 2,
          "expenseClaims": 1
        }
      },
      "payrollStatus": {
        "status": "Active",
        "currentPeriod": "January 2026",
        "processingDate": "2026-01-25",
        "totalAmount": 485000.00,
        "currency": "TZS"
      }
    },
    "summary": {
      "newHiresThisMonth": 5,
      "terminationsThisMonth": 2,
      "openPositions": 8,
      "avgAttendanceRate": 94.5
    }
  }
}
```

---

### Response (Employee View)

```json
{
  "success": true,
  "message": "Dashboard statistics retrieved successfully",
  "data": {
    "view": "employee",
    "date": "2026-01-23",
    "employeeId": 1234,
    "employeeName": "John Doe",
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
        "date": "2026-01-31",
        "daysRemaining": 8,
        "estimatedNetPay": 2500000.00,
        "currency": "TZS"
      }
    },
    "attendance": {
      "todayStatus": "Present",
      "checkInTime": "08:15:00",
      "checkOutTime": null,
      "workedHoursToday": 4.5
    }
  }
}
```

---

## 2. Announcements API

### `GET /dashboard/announcements`

Returns company announcements for the dashboard.

**Headers:**
```
Authorization: Bearer {{ACCESS_TOKEN}}
Content-Type: application/json
```

**Query Parameters:**

| Parameter     | Type    | Required | Description                                        |
|---------------|---------|----------|----------------------------------------------------|
| `page`        | number  | No       | Page number (default: 1)                           |
| `page_size`   | number  | No       | Items per page (default: 5, max: 20)               |
| `type`        | string  | No       | Filter: `info`, `warning`, `success`, `urgent`     |
| `active_only` | boolean | No       | Only show active announcements (default: true)     |

---

### Response

```json
{
  "success": true,
  "message": "Announcements retrieved successfully",
  "data": {
    "announcements": [
      {
        "id": 1,
        "title": "New Policy Update",
        "description": "Updated leave policy effective from next month. All employees are required to review the changes in the HR portal.",
        "type": "info",
        "priority": "normal",
        "createdAt": "2026-01-23T10:30:00Z",
        "createdBy": {
          "id": 5,
          "name": "HR Admin"
        },
        "expiresAt": "2026-02-23T23:59:59Z",
        "isRead": false,
        "attachments": [
          {
            "id": 1,
            "name": "Leave_Policy_2026.pdf",
            "url": "/attachments/policies/leave_policy_2026.pdf",
            "size": "245KB"
          }
        ],
        "relativeTime": "2 hours ago"
      },
      {
        "id": 2,
        "title": "Upcoming Holiday",
        "description": "Office will be closed on January 25th for National Holiday. Regular operations resume on January 26th.",
        "type": "warning",
        "priority": "high",
        "createdAt": "2026-01-22T14:00:00Z",
        "createdBy": {
          "id": 5,
          "name": "HR Admin"
        },
        "expiresAt": "2026-01-26T00:00:00Z",
        "isRead": true,
        "attachments": [],
        "relativeTime": "1 day ago"
      },
      {
        "id": 3,
        "title": "Team Meeting",
        "description": "All hands meeting scheduled for Friday 3 PM in the main conference room. Attendance is mandatory.",
        "type": "success",
        "priority": "normal",
        "createdAt": "2026-01-21T09:00:00Z",
        "createdBy": {
          "id": 10,
          "name": "CEO Office"
        },
        "expiresAt": "2026-01-24T15:00:00Z",
        "isRead": true,
        "attachments": [],
        "relativeTime": "2 days ago"
      }
    ],
    "meta": {
      "page": 1,
      "pageSize": 5,
      "total": 12,
      "totalPages": 3,
      "unreadCount": 4
    }
  }
}
```

---

### `POST /dashboard/announcements/{id}/read`

Mark an announcement as read.

**Path Parameters:**

| Parameter | Type   | Required | Description        |
|-----------|--------|----------|--------------------|
| `id`      | number | Yes      | Announcement ID    |

**Response:**

```json
{
  "success": true,
  "message": "Announcement marked as read",
  "data": {
    "announcementId": 1,
    "readAt": "2026-01-23T12:45:00Z"
  }
}
```

---

## 3. Upcoming Events API (Admin)

### `GET /dashboard/events`

Returns upcoming birthdays, work anniversaries, and other employee events.

**Headers:**
```
Authorization: Bearer {{ACCESS_TOKEN}}
Content-Type: application/json
```

**Query Parameters:**

| Parameter       | Type   | Required | Description                                                    |
|-----------------|--------|----------|----------------------------------------------------------------|
| `days_ahead`    | number | No       | Number of days to look ahead (default: 30, max: 90)            |
| `event_type`    | string | No       | Filter: `birthday`, `anniversary`, `probation_end`, `all`      |
| `department_id` | number | No       | Filter by department                                           |
| `page`          | number | No       | Page number (default: 1)                                       |
| `page_size`     | number | No       | Items per page (default: 10)                                   |

---

### Response

```json
{
  "success": true,
  "message": "Upcoming events retrieved successfully",
  "data": {
    "events": [
      {
        "id": 1,
        "employeeId": 101,
        "employeeName": "John Doe",
        "employeePhoto": "/photos/employees/101.jpg",
        "department": "Engineering",
        "eventType": "birthday",
        "eventDate": "2026-01-25",
        "displayDate": "Jan 25",
        "daysUntil": 2,
        "age": 32,
        "initials": "JD"
      },
      {
        "id": 2,
        "employeeId": 102,
        "employeeName": "Sarah Smith",
        "employeePhoto": "/photos/employees/102.jpg",
        "department": "Marketing",
        "eventType": "anniversary",
        "eventDate": "2026-01-27",
        "displayDate": "Jan 27",
        "daysUntil": 4,
        "yearsOfService": 5,
        "initials": "SS"
      },
      {
        "id": 3,
        "employeeId": 103,
        "employeeName": "Mike Johnson",
        "employeePhoto": null,
        "department": "Sales",
        "eventType": "birthday",
        "eventDate": "2026-01-30",
        "displayDate": "Jan 30",
        "daysUntil": 7,
        "age": 28,
        "initials": "MJ"
      },
      {
        "id": 4,
        "employeeId": 104,
        "employeeName": "Emily Brown",
        "employeePhoto": "/photos/employees/104.jpg",
        "department": "HR",
        "eventType": "probation_end",
        "eventDate": "2026-02-01",
        "displayDate": "Feb 1",
        "daysUntil": 9,
        "probationMonths": 3,
        "initials": "EB"
      }
    ],
    "summary": {
      "birthdaysThisMonth": 8,
      "anniversariesThisMonth": 3,
      "probationEndsThisMonth": 2
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

---

## 4. Recent Activity API (Employee)

### `GET /dashboard/my-activity`

Returns the logged-in employee's recent activities.

**Headers:**
```
Authorization: Bearer {{ACCESS_TOKEN}}
Content-Type: application/json
```

**Query Parameters:**

| Parameter       | Type   | Required | Description                                              |
|-----------------|--------|----------|----------------------------------------------------------|
| `page`          | number | No       | Page number (default: 1)                                 |
| `page_size`     | number | No       | Items per page (default: 10, max: 50)                    |
| `activity_type` | string | No       | Filter: `leave`, `timesheet`, `payslip`, `request`, `all`|
| `days`          | number | No       | Activities from last N days (default: 30)                |

---

### Response

```json
{
  "success": true,
  "message": "Recent activities retrieved successfully",
  "data": {
    "activities": [
      {
        "id": 1,
        "type": "leave_request",
        "title": "Leave Request Submitted",
        "description": "Annual leave for Dec 15-20",
        "status": "pending",
        "statusLabel": "Pending approval",
        "icon": "calendar",
        "color": "blue",
        "createdAt": "2026-01-21T09:30:00Z",
        "relativeTime": "2 days ago",
        "metadata": {
          "leaveType": "Annual",
          "startDate": "2026-12-15",
          "endDate": "2026-12-20",
          "days": 6,
          "approver": "Jane Manager"
        },
        "actionUrl": "/leave/requests/123"
      },
      {
        "id": 2,
        "type": "timesheet_update",
        "title": "Timesheet Updated",
        "description": "38.5 hours logged this week",
        "status": "completed",
        "statusLabel": "Submitted",
        "icon": "clock",
        "color": "green",
        "createdAt": "2026-01-22T17:00:00Z",
        "relativeTime": "1 day ago",
        "metadata": {
          "weekStarting": "2026-01-20",
          "totalHours": 38.5,
          "regularHours": 38.5,
          "overtimeHours": 0
        },
        "actionUrl": "/attendance/timesheets"
      },
      {
        "id": 3,
        "type": "payslip_available",
        "title": "Payslip Available",
        "description": "December 2025 payslip is ready",
        "status": "new",
        "statusLabel": "View now",
        "icon": "currency",
        "color": "purple",
        "createdAt": "2026-01-15T10:00:00Z",
        "relativeTime": "1 week ago",
        "metadata": {
          "payPeriod": "December 2025",
          "netPay": 2450000.00,
          "currency": "TZS"
        },
        "actionUrl": "/self-service/my-payslips/dec-2025"
      },
      {
        "id": 4,
        "type": "letter_request",
        "title": "Letter Request Approved",
        "description": "Employment verification letter ready for download",
        "status": "approved",
        "statusLabel": "Ready",
        "icon": "document",
        "color": "green",
        "createdAt": "2026-01-18T14:30:00Z",
        "relativeTime": "5 days ago",
        "metadata": {
          "letterType": "Employment Verification",
          "requestedOn": "2026-01-16",
          "approvedBy": "HR Admin"
        },
        "actionUrl": "/self-service/requests/456"
      }
    ],
    "summary": {
      "pendingRequests": 2,
      "recentApprovals": 3,
      "newNotifications": 1
    },
    "meta": {
      "page": 1,
      "pageSize": 10,
      "total": 25,
      "totalPages": 3
    }
  }
}
```

---

## 5. Quick Actions API

### `GET /dashboard/quick-actions`

Returns personalized quick actions based on user role and permissions.

**Headers:**
```
Authorization: Bearer {{ACCESS_TOKEN}}
Content-Type: application/json
```

---

### Response (Admin View)

```json
{
  "success": true,
  "message": "Quick actions retrieved successfully",
  "data": {
    "view": "admin",
    "actions": [
      {
        "id": "add_employee",
        "title": "Add Employee",
        "description": "Register new employee",
        "icon": "user-plus",
        "color": "blue",
        "href": "/employees/add",
        "badge": null,
        "enabled": true
      },
      {
        "id": "view_reports",
        "title": "View Reports",
        "description": "Generate reports",
        "icon": "chart",
        "color": "green",
        "href": "/reports/analytics",
        "badge": null,
        "enabled": true
      },
      {
        "id": "approve_requests",
        "title": "Approve Requests",
        "description": "Pending approvals",
        "icon": "check-circle",
        "color": "orange",
        "href": "/leave/requests",
        "badge": {
          "count": 8,
          "type": "warning"
        },
        "enabled": true
      },
      {
        "id": "manage_payroll",
        "title": "Manage Payroll",
        "description": "Process payroll",
        "icon": "currency",
        "color": "purple",
        "href": "/payroll/dashboard",
        "badge": {
          "text": "Due Soon",
          "type": "info"
        },
        "enabled": true
      }
    ]
  }
}
```

---

### Response (Employee View)

```json
{
  "success": true,
  "message": "Quick actions retrieved successfully",
  "data": {
    "view": "employee",
    "actions": [
      {
        "id": "apply_leave",
        "title": "Apply Leave",
        "description": "Submit leave request",
        "icon": "calendar",
        "color": "blue",
        "href": "/leave/apply",
        "badge": null,
        "enabled": true
      },
      {
        "id": "my_payslips",
        "title": "My Payslips",
        "description": "View payslips",
        "icon": "document",
        "color": "orange",
        "href": "/self-service/my-payslips",
        "badge": {
          "count": 1,
          "type": "info",
          "text": "New"
        },
        "enabled": true
      },
      {
        "id": "requests_letters",
        "title": "Requests & Letters",
        "description": "Submit requests",
        "icon": "mail",
        "color": "green",
        "href": "/self-service/requests",
        "badge": null,
        "enabled": true
      },
      {
        "id": "people_directory",
        "title": "People Directory",
        "description": "Find colleagues",
        "icon": "users",
        "color": "purple",
        "href": "/self-service/directory",
        "badge": null,
        "enabled": true
      }
    ]
  }
}
```

---

## 6. Pending Approvals API (Admin)

### `GET /dashboard/pending-approvals`

Returns detailed breakdown of pending approvals for managers/admins.

**Headers:**
```
Authorization: Bearer {{ACCESS_TOKEN}}
Content-Type: application/json
```

**Query Parameters:**

| Parameter   | Type   | Required | Description                                    |
|-------------|--------|----------|------------------------------------------------|
| `type`      | string | No       | Filter: `leave`, `timesheet`, `expense`, `all` |
| `page`      | number | No       | Page number (default: 1)                       |
| `page_size` | number | No       | Items per page (default: 10)                   |

---

### Response

```json
{
  "success": true,
  "message": "Pending approvals retrieved successfully",
  "data": {
    "summary": {
      "total": 8,
      "leaveRequests": 5,
      "timesheetCorrections": 2,
      "expenseClaims": 1
    },
    "approvals": [
      {
        "id": 101,
        "type": "leave",
        "employeeId": 201,
        "employeeName": "Alice Johnson",
        "employeePhoto": "/photos/employees/201.jpg",
        "department": "Engineering",
        "title": "Annual Leave Request",
        "description": "Family vacation",
        "dates": "Jan 28 - Feb 2, 2026",
        "days": 5,
        "submittedAt": "2026-01-20T09:00:00Z",
        "relativeTime": "3 days ago",
        "priority": "normal",
        "actionUrl": "/leave/requests/101"
      },
      {
        "id": 102,
        "type": "timesheet",
        "employeeId": 202,
        "employeeName": "Bob Smith",
        "employeePhoto": null,
        "department": "Sales",
        "title": "Timesheet Correction",
        "description": "Missed clock-out on Jan 18",
        "originalTime": null,
        "correctedTime": "17:30",
        "submittedAt": "2026-01-21T10:30:00Z",
        "relativeTime": "2 days ago",
        "priority": "high",
        "actionUrl": "/attendance/corrections/102"
      },
      {
        "id": 103,
        "type": "expense",
        "employeeId": 203,
        "employeeName": "Carol Davis",
        "employeePhoto": "/photos/employees/203.jpg",
        "department": "Marketing",
        "title": "Travel Expense Claim",
        "description": "Client meeting in Arusha",
        "amount": 450000.00,
        "currency": "TZS",
        "submittedAt": "2026-01-22T14:00:00Z",
        "relativeTime": "1 day ago",
        "priority": "normal",
        "actionUrl": "/expenses/claims/103"
      }
    ],
    "meta": {
      "page": 1,
      "pageSize": 10,
      "total": 8,
      "totalPages": 1
    }
  }
}
```

---

## 7. TypeScript Types Reference

```typescript
// ============================================
// Dashboard Statistics Types
// ============================================

export type DashboardView = "admin" | "employee";
export type TrendDirection = "up" | "down" | "stable";

export type AdminStatistics = {
  view: "admin";
  date: string;
  stats: {
    totalHeadcount: {
      value: number;
      change: number;
      changePercent: number;
      trend: TrendDirection;
    };
    presentToday: {
      value: number;
      total: number;
      percentage: number;
      absent: number;
      onLeave: number;
    };
    pendingApprovals: {
      value: number;
      breakdown: {
        leaveRequests: number;
        timesheetCorrections: number;
        expenseClaims: number;
      };
    };
    payrollStatus: {
      status: "Active" | "Processing" | "Completed" | "Pending";
      currentPeriod: string;
      processingDate: string;
      totalAmount: number;
      currency: string;
    };
  };
  summary: {
    newHiresThisMonth: number;
    terminationsThisMonth: number;
    openPositions: number;
    avgAttendanceRate: number;
  };
};

export type EmployeeStatistics = {
  view: "employee";
  date: string;
  employeeId: number;
  employeeName: string;
  stats: {
    leaveBalance: {
      value: number;
      unit: string;
      breakdown: {
        annual: number;
        sick: number;
        personal: number;
      };
      usedThisYear: number;
    };
    hoursThisWeek: {
      value: number;
      target: number;
      overtime: number;
      status: "on_track" | "behind" | "ahead";
    };
    pendingRequests: {
      value: number;
      breakdown: {
        leaveRequests: number;
        letterRequests: number;
      };
    };
    nextPayday: {
      date: string;
      daysRemaining: number;
      estimatedNetPay: number;
      currency: string;
    };
  };
  attendance: {
    todayStatus: "Present" | "Absent" | "On Leave" | "Not Clocked In";
    checkInTime: string | null;
    checkOutTime: string | null;
    workedHoursToday: number;
  };
};

export type DashboardStatistics = AdminStatistics | EmployeeStatistics;

export type DashboardStatisticsResponse = {
  success: boolean;
  message: string;
  data: DashboardStatistics;
};

// ============================================
// Announcements Types
// ============================================

export type AnnouncementType = "info" | "warning" | "success" | "urgent";
export type AnnouncementPriority = "low" | "normal" | "high" | "critical";

export type AnnouncementAttachment = {
  id: number;
  name: string;
  url: string;
  size: string;
};

export type Announcement = {
  id: number;
  title: string;
  description: string;
  type: AnnouncementType;
  priority: AnnouncementPriority;
  createdAt: string;
  createdBy: {
    id: number;
    name: string;
  };
  expiresAt: string | null;
  isRead: boolean;
  attachments: AnnouncementAttachment[];
  relativeTime: string;
};

export type AnnouncementsResponse = {
  success: boolean;
  message: string;
  data: {
    announcements: Announcement[];
    meta: {
      page: number;
      pageSize: number;
      total: number;
      totalPages: number;
      unreadCount: number;
    };
  };
};

// ============================================
// Events Types
// ============================================

export type EventType = "birthday" | "anniversary" | "probation_end";

export type UpcomingEvent = {
  id: number;
  employeeId: number;
  employeeName: string;
  employeePhoto: string | null;
  department: string;
  eventType: EventType;
  eventDate: string;
  displayDate: string;
  daysUntil: number;
  age?: number;
  yearsOfService?: number;
  probationMonths?: number;
  initials: string;
};

export type EventsResponse = {
  success: boolean;
  message: string;
  data: {
    events: UpcomingEvent[];
    summary: {
      birthdaysThisMonth: number;
      anniversariesThisMonth: number;
      probationEndsThisMonth: number;
    };
    meta: {
      page: number;
      pageSize: number;
      total: number;
      totalPages: number;
    };
  };
};

// ============================================
// Activity Types
// ============================================

export type ActivityType =
  | "leave_request"
  | "timesheet_update"
  | "payslip_available"
  | "letter_request"
  | "approval";

export type ActivityStatus =
  | "pending"
  | "approved"
  | "rejected"
  | "completed"
  | "new";

export type RecentActivity = {
  id: number;
  type: ActivityType;
  title: string;
  description: string;
  status: ActivityStatus;
  statusLabel: string;
  icon: string;
  color: string;
  createdAt: string;
  relativeTime: string;
  metadata: Record<string, any>;
  actionUrl: string;
};

export type MyActivityResponse = {
  success: boolean;
  message: string;
  data: {
    activities: RecentActivity[];
    summary: {
      pendingRequests: number;
      recentApprovals: number;
      newNotifications: number;
    };
    meta: {
      page: number;
      pageSize: number;
      total: number;
      totalPages: number;
    };
  };
};

// ============================================
// Quick Actions Types
// ============================================

export type BadgeType = "info" | "warning" | "success" | "error";

export type QuickActionBadge = {
  count?: number;
  text?: string;
  type: BadgeType;
} | null;

export type QuickAction = {
  id: string;
  title: string;
  description: string;
  icon: string;
  color: string;
  href: string;
  badge: QuickActionBadge;
  enabled: boolean;
};

export type QuickActionsResponse = {
  success: boolean;
  message: string;
  data: {
    view: DashboardView;
    actions: QuickAction[];
  };
};

// ============================================
// Pending Approvals Types
// ============================================

export type ApprovalType = "leave" | "timesheet" | "expense";
export type ApprovalPriority = "low" | "normal" | "high" | "urgent";

export type PendingApproval = {
  id: number;
  type: ApprovalType;
  employeeId: number;
  employeeName: string;
  employeePhoto: string | null;
  department: string;
  title: string;
  description: string;
  submittedAt: string;
  relativeTime: string;
  priority: ApprovalPriority;
  actionUrl: string;
  // Type-specific fields (leave)
  dates?: string;
  days?: number;
  // Type-specific fields (expense)
  amount?: number;
  currency?: string;
  // Type-specific fields (timesheet)
  originalTime?: string | null;
  correctedTime?: string;
};

export type PendingApprovalsResponse = {
  success: boolean;
  message: string;
  data: {
    summary: {
      total: number;
      leaveRequests: number;
      timesheetCorrections: number;
      expenseClaims: number;
    };
    approvals: PendingApproval[];
    meta: {
      page: number;
      pageSize: number;
      total: number;
      totalPages: number;
    };
  };
};
```

---

## 8. Error Handling

### Standard Error Response Format

```json
{
  "success": false,
  "message": "Error description",
  "error": {
    "code": "ERROR_CODE",
    "details": "Detailed error message"
  }
}
```

### Common Error Codes

| HTTP Status | Error Code      | Description                                      |
|-------------|-----------------|--------------------------------------------------|
| 400         | `BAD_REQUEST`   | Invalid request parameters                       |
| 401         | `UNAUTHORIZED`  | Missing or invalid authentication token          |
| 403         | `FORBIDDEN`     | User lacks permission for this resource          |
| 404         | `NOT_FOUND`     | Requested resource not found                     |
| 422         | `VALIDATION`    | Request validation failed                        |
| 429         | `RATE_LIMITED`  | Too many requests                                |
| 500         | `SERVER_ERROR`  | Internal server error                            |

### Example Error Responses

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Authentication required",
  "error": {
    "code": "UNAUTHORIZED",
    "details": "Invalid or expired token"
  }
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "Access denied",
  "error": {
    "code": "FORBIDDEN",
    "details": "You do not have permission to view admin statistics"
  }
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Resource not found",
  "error": {
    "code": "NOT_FOUND",
    "details": "Announcement with ID 999 not found"
  }
}
```

---

## API Summary Table

| # | Endpoint                              | Method | Description                                    | Access   |
|---|---------------------------------------|--------|------------------------------------------------|----------|
| 1 | `/dashboard/statistics`               | GET    | KPI stats (headcount, attendance, approvals)   | Both     |
| 2 | `/dashboard/announcements`            | GET    | Company announcements list                     | Both     |
| 3 | `/dashboard/announcements/{id}/read`  | POST   | Mark announcement as read                      | Both     |
| 4 | `/dashboard/events`                   | GET    | Birthdays, anniversaries, events               | Admin    |
| 5 | `/dashboard/my-activity`              | GET    | Employee's recent activities                   | Employee |
| 6 | `/dashboard/quick-actions`            | GET    | Personalized quick action links                | Both     |
| 7 | `/dashboard/pending-approvals`        | GET    | Pending items requiring approval               | Admin    |

---

## Changelog

| Version | Date       | Changes                          |
|---------|------------|----------------------------------|
| 1.0.0   | 2026-01-23 | Initial API documentation        |

---

## Notes

1. **Authentication**: All endpoints require a valid JWT token in the Authorization header.
2. **Role-based Access**: Some endpoints return different data based on user role (admin vs employee).
3. **Pagination**: List endpoints support pagination via `page` and `page_size` parameters.
4. **Date Formats**: All dates use ISO 8601 format (`YYYY-MM-DD` or `YYYY-MM-DDTHH:mm:ssZ`).
5. **Currency**: All monetary values are returned as numbers with the currency code specified separately.

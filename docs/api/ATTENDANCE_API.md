# Attendance Reports API - Complete Documentation

This document provides comprehensive API documentation for the **Attendance Reports** module, including **NOC Department shift handling**, **shift-specific late thresholds**, and **department-based attendance tracking**.

**Base URL**: `http://localhost:8080/api/v1`

**Authentication**: All endpoints require a valid Bearer token.
```
Authorization: Bearer <token>
```

---

## Standard Response Envelope

### Success
```json
{
  "success": true,
  "message": "string",
  "data": {}
}
```

### Error
```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "field": "reason"
  }
}
```

---

# 1) Attendance Overview

Retrieves a high-level overview of attendance statistics for a given date range.

## Endpoint
```
GET /attendance/overview
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `start_date` | string | Yes | Start date in `YYYY-MM-DD HH:MM:SS` or `YYYY-MM-DD` format |
| `end_date` | string | Yes | End date in `YYYY-MM-DD HH:MM:SS` or `YYYY-MM-DD` format |
| `view` | string | No | View type (e.g., "daily", "weekly"). Default: empty |
| `late_threshold` | string | No | Late threshold in `HH:MM:SS` format. Default: `09:00:00` |
| `department_id` | number | No | Filter by department ID |
| `location_id` | number | No | Filter by location ID |
| `emp_code` | string | No | Filter by employee code |

## Example Request
```
GET /attendance/overview?start_date=2026-04-01&end_date=2026-04-28&department_id=5
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Attendance overview retrieved successfully",
  "data": {
    "period": {
      "start_date": "2026-04-01 00:00:00",
      "end_date": "2026-04-28 23:59:59",
      "view": "daily",
      "late_threshold": "09:00:00"
    },
    "filters": {
      "department_id": 5,
      "location_id": null,
      "emp_code": null
    },
    "totals": {
      "total_employees": 150,
      "present": 135,
      "absent": 10,
      "on_leave": 5,
      "late": 8,
      "exceptions": 8
    },
    "meta": {
      "units": {
        "late": "employee_days",
        "exceptions": "employee_days"
      }
    }
  }
}
```

### Key Features:
- **Late calculation**: Uses the `late_threshold` query parameter or shifts-specific `LateThresholdTime` from RosterAssignment
- **NOC Support**: When filtering by NOC department, late thresholds are calculated based on assigned shift (e.g., 07:30 for NOC Morning shift)
- **Department filtering**: Can filter to see attendance for specific departments

---

# 2) Attendance Summary

Returns high-level attendance KPIs and trends.

## Endpoint
```
GET /attendance/reports/summary
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |
| `departmentId` | number | No | Filter by department ID |

## Example Request
```
GET /attendance/reports/summary?startDate=2026-04-01&endDate=2026-04-28&departmentId=5
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Attendance summary retrieved successfully",
  "data": {
    "presentRate": 87.5,
    "absentRate": 8.3,
    "lateRate": 4.2,
    "averageWorkingHours": 8.2,
    "overtimeHours": 245.5,
    "exceptionRate": 4.2,
    "onTimeRate": 95.8,
    "shiftAdherence": 98.1,
    "totalEmployees": 120,
    "trends": {
      "presentChange": 2.5,
      "absentChange": -1.2,
      "lateChange": -0.8
    }
  }
}
```

---

# 3) Attendance Trends

Returns daily or weekly attendance trends.

## Endpoint
```
GET /attendance/reports/trends
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |
| `departmentId` | number | No | Filter by department ID |
| `interval` | string | No | `daily` or `weekly`. Default: `daily` |

## Example Request
```
GET /attendance/reports/trends?startDate=2026-04-01&endDate=2026-04-28&departmentId=5
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Attendance trends retrieved successfully",
  "data": [
    {
      "date": "2026-04-01",
      "day": "Wednesday",
      "present": 115,
      "absent": 3,
      "late": 2,
      "onLeave": 0,
      "totalScheduled": 120
    },
    {
      "date": "2026-04-02",
      "day": "Thursday",
      "present": 118,
      "absent": 1,
      "late": 1,
      "onLeave": 1,
      "totalScheduled": 120
    }
  ]
}
```

---

# 4) Department Attendance Stats (Single Day)

Returns department-wise attendance for a specific date.

## Endpoint
```
GET /attendance/reports/by-department
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `date` | string | Yes | Date in `YYYY-MM-DD` format |

## Example Request
```
GET /attendance/reports/by-department?date=2026-04-28
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Department attendance retrieved successfully",
  "data": [
    {
      "departmentId": "5",
      "departmentName": "NOC Department",
      "present": 25,
      "absent": 2,
      "late": 1,
      "total": 28,
      "presentRate": 89.29
    },
    {
      "departmentId": "3",
      "departmentName": "Human Resources",
      "present": 12,
      "absent": 0,
      "late": 0,
      "total": 12,
      "presentRate": 100.0
    }
  ]
}
```

### Key Features:
- **NOC Department handling**: NOC employees are tracked with their specific shifts
- **Late calculation**: Uses shift-specific `LateThresholdTime` for NOC employees

---

# 5) Department Attendance Stats (Date Range)

Returns department-wise attendance statistics for a date range with **shift-specific late thresholds**.

## Endpoint
```
GET /attendance/reports/department-stats
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |

## Example Request
```
GET /attendance/reports/department-stats?startDate=2026-04-01&endDate=2026-04-28
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Department attendance retrieved successfully",
  "data": [
    {
      "departmentId": "5",
      "departmentName": "NOC Department",
      "present": 650,
      "absent": 45,
      "late": 12,
      "onLeave": 15,
      "attendancePercentage": 89.04
    },
    {
      "departmentId": "3",
      "departmentName": "Human Resources",
      "present": 300,
      "absent": 10,
      "late": 5,
      "onLeave": 5,
      "attendancePercentage": 93.75
    }
  ]
}
```

### Key Features:
- **Shift-specific late thresholds**: Late count is calculated using `LateThresholdTime` from the employee's shift via RosterAssignment
- **NOC Support**: For NOC department employees, late is calculated against their assigned shift's threshold (e.g., 07:30 for NOC Morning shift)
- **Department breakdown**: Each department shows separate present, absent, late, and onLeave counts

---

# 6) Employee Attendance Stats

Returns detailed attendance for a specific employee.

## Endpoint
```
GET /attendance/reports/employee-stats
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `employeeId` | number | Yes | Employee ID |
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |

## Example Request
```
GET /attendance/reports/employee-stats?employeeId=123&startDate=2026-04-01&endDate=2026-04-28
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Employee attendance report retrieved successfully",
  "data": {
    "employee_id": "123",
    "employee_name": "John Doe",
    "department": "NOC Department",
    "position": "NOC Analyst",
    "period": {
      "start_date": "2026-04-01",
      "end_date": "2026-04-28"
    },
    "attendance": {
      "total_working_days": 22,
      "present_days": 20,
      "absent_days": 1,
      "late_days": 1,
      "early_departure_days": 0,
      "present_percentage": 90.91,
      "average_working_hours": 8.5,
      "total_working_hours": 170.0
    },
    "timesheet": {
      "total_hours": 180.5,
      "regular_hours": 170.0,
      "overtime_hours": 10.5,
      "billable_hours": 180.5,
      "non_billable_hours": 0
    },
    "leave": {
      "total_leave_days": 1,
      "leave_breakdown": [
        { "leave_type": "Annual Leave", "days": 1 }
      ],
      "pending_requests": 0,
      "approved_requests": 1,
      "rejected_requests": 0
    },
    "overtime": {
      "total_requests": 2,
      "approved_hours": 10.5,
      "pending_hours": 0,
      "rejected_hours": 0,
      "total_payout": 150.00
    },
    "compliance": {
      "violations_count": 0,
      "violations": []
    }
  }
}
```

---

# 7) Comprehensive Employee Report

Returns comprehensive attendance, timesheet, leave, overtime, and compliance data for employees.

## Endpoint
```
GET /attendance/reports/employee-comprehensive
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |
| `employeeId` | number | No | Filter by specific employee ID |
| `departmentId` | number | No | Filter by department ID |
| `locationId` | number | No | Filter by location ID |
| `page` | number | No | Page number (default: 1) |
| `pageSize` | number | No | Page size (default: 20, max: 100) |

## Example Request
```
GET /attendance/reports/employee-comprehensive?startDate=2026-04-01&endDate=2026-04-28&departmentId=5&page=1&pageSize=20
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Comprehensive employee report retrieved successfully",
  "data": [
    {
      "employee_id": "EMP-001",
      "employee_name": "John Doe",
      "department": "NOC Department",
      "position": "NOC Analyst",
      "period": {
        "start_date": "2026-04-01",
        "end_date": "2026-04-28"
      },
      "attendance": {
        "total_working_days": 22,
        "present_days": 20,
        "absent_days": 1,
        "late_days": 1,
        "early_departure_days": 0,
        "present_percentage": 90.91,
        "average_working_hours": 8.5,
        "total_working_hours": 170.0
      },
      "timesheet": {
        "total_hours": 180.5,
        "regular_hours": 170.0,
        "overtime_hours": 10.5,
        "billable_hours": 180.5,
        "non_billable_hours": 0
      },
      "leave": {
        "total_leave_days": 1,
        "leave_breakdown": [
          { "leave_type": "Annual Leave", "days": 1 }
        ],
        "pending_requests": 0,
        "approved_requests": 1,
        "rejected_requests": 0
      },
      "overtime": {
        "total_requests": 2,
        "approved_hours": 10.5,
        "pending_hours": 0,
        "rejected_hours": 0,
        "total_payout": 150.00
      },
      "compliance": {
        "violations_count": 0,
        "violations": []
      }
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

---

# 8) Compliance Violations

Returns compliance issues within a date range.

## Endpoint
```
GET /attendance/reports/compliance
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |
| `departmentId` | number | No | Filter by department ID |
| `severity` | string | No | Filter by severity (low, medium, high, critical) |

## Example Request
```
GET /attendance/reports/compliance?startDate=2026-04-01&endDate=2026-04-28&departmentId=5
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Compliance violations retrieved successfully",
  "data": [
    {
      "id": "VIO-001",
      "employeeId": "EMP-001",
      "employeeName": "John Doe",
      "department": "NOC Department",
      "violationType": "Late Arrival",
      "violationDate": "2026-04-15",
      "severity": "low",
      "details": "Arrived at 07:45, threshold is 07:30",
      "status": "Open"
    }
  ]
}
```

---

# 9) Overtime Analysis

Returns overtime analysis and top contributors.

## Endpoint
```
GET /attendance/reports/overtime
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |

## Example Request
```
GET /attendance/reports/overtime?startDate=2026-04-01&endDate=2026-04-28
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Overtime analysis retrieved successfully",
  "data": {
    "totalOvertimeHours": 450.5,
    "averageOvertimePerEmployee": 12.5,
    "estimatedCost": 8500.00,
    "currency": "TZS",
    "topContributors": [
      {
        "employeeId": "EMP-042",
        "name": "Jane Smith",
        "hours": 35.0
      },
      {
        "employeeId": "EMP-017",
        "name": "Bob Wilson",
        "hours": 28.5
      }
    ]
  }
}
```

---

# 10) Export Daily Attendance

Exports daily attendance data to CSV or XLSX format with **department information**.

## Endpoint
```
POST /attendance/reports/export
```

## Request Body
```json
{
  "reportType": "daily_attendance",
  "startDate": "2026-04-01",
  "endDate": "2026-04-28",
  "departmentId": "5",
  "format": "csv"
}
```

### Field Descriptions
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `reportType` | string | Yes | Report type: `daily_attendance` |
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |
| `departmentId` | string | No | Filter by department ID |
| `format` | string | No | Export format: `csv` (default) or `xlsx` |

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Report generation started",
  "data": {
    "downloadUrl": "https://localhost:8080/storage/reports/attendance_daily_attendance_2026-04-01_1745846400.csv",
    "jobId": "job_1745846400",
    "expiresAt": "2026-04-29T12:00:00Z"
  }
}
```

### Export File Columns
The exported CSV/XLSX file includes:
- `emp_code` - Employee code
- `first_name` - Employee first name
- `last_name` - Employee last name
- `department_name` - **NEW: Department name**
- `day` - Date
- `check_in` - First punch time
- `check_out` - Last punch time
- `punch_count` - Total punches for the day

---

# 11) Trigger Monthly Report

Manually triggers end-of-month report generation.

## Endpoint
```
POST /attendance/reports/trigger-monthly
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `month` | string | Yes | Month in `YYYY-MM` format |

## Example Request
```
POST /attendance/reports/trigger-monthly?month=2026-04
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Report generation triggered successfully. HR will receive an email shortly.",
  "data": null
}
```

---

# 12) Timesheet Summary Report

Returns timesheet summary with billable hours and utilization.

## Endpoint
```
GET /attendance/reports/timesheets
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `start_date` | string | Yes | Start date in `YYYY-MM-DD` format |
| `end_date` | string | Yes | End date in `YYYY-MM-DD` format |

## Example Request
```
GET /attendance/reports/timesheets?start_date=2026-04-01&end_date=2026-04-28
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Timesheet summary report",
  "data": {
    "period": {
      "start_date": "2026-04-01",
      "end_date": "2026-04-28"
    },
    "summary": {
      "total_entries": 450,
      "total_hours": 3600.5,
      "billable_hours": 3200.0,
      "non_billable_hours": 400.5,
      "utilization_rate": 88.88,
      "unique_employees": 45,
      "unique_projects": 12
    },
    "weekly_submission_status": [
      { "status": "draft", "count": 5 },
      { "status": "submitted", "count": 30 },
      { "status": "approved", "count": 10 }
    ],
    "project_distribution": [
      {
        "project_name": "Project Alpha",
        "total_hours": 800.5,
        "billable_hours": 750.0,
        "employee_count": 10
      }
    ]
  }
}
```

---

# 13) Project Utilization Report

Returns per-project utilization metrics.

## Endpoint
```
GET /attendance/reports/project-utilization
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `start_date` | string | Yes | Start date in `YYYY-MM-DD` format |
| `end_date` | string | Yes | End date in `YYYY-MM-DD` format |

## Example Request
```
GET /attendance/reports/project-utilization?start_date=2026-04-01&end_date=2026-04-28
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Project utilization report",
  "data": {
    "period": {
      "start_date": "2026-04-01",
      "end_date": "2026-04-28"
    },
    "projects": [
      {
        "project_name": "Project Alpha",
        "client_name": "Client A",
        "total_hours": 1200.5,
        "billable_hours": 1100.0,
        "non_billable_hours": 100.5,
        "team_size": 15,
        "utilization_rate": 91.63
      }
    ]
  }
}
```

---

# 14) Employee Utilization Report

Returns per-employee utilization metrics.

## Endpoint
```
GET /attendance/reports/employee-utilization
```

## Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `start_date` | string | Yes | Start date in `YYYY-MM-DD` format |
| `end_date` | string | Yes | End date in `YYYY-MM-DD` format |

## Example Request
```
GET /attendance/reports/employee-utilization?start_date=2026-04-01&end_date=2026-04-28
```

## Success Response (200 OK)
```json
{
  "success": true,
  "message": "Employee utilization report",
  "data": {
    "period": {
      "start_date": "2026-04-01",
      "end_date": "2026-04-28"
    },
    "employees": [
      {
        "employee_id": "EMP-001",
        "employee_name": "John Doe",
        "department": "NOC Department",
        "total_hours": 180.5,
        "billable_hours": 175.0,
        "non_billable_hours": 5.5,
        "utilization_rate": 96.95,
        "overtime_hours": 10.5
      }
    ]
  }
}
```

---

# NOC Department Implementation Guide

## Overview

The NOC (Network Operations Center) department has unique shift requirements that differ from regular departments:

### Shift Schedule
| Shift | Check-In | Check-Out | Late Threshold |
|-------|----------|-----------|----------------|
| NOC Morning | 07:00 | 15:00 | **07:30** |
| NOC Evening | 15:00 | 23:00 | 15:00 |
| NOC Night | 23:00 | 05:00 | 23:00 |
| Regular (Non-NOC) | 08:30 | 17:30 | 08:30 |

### Key Implementation Details

1. **Department-Shift Association**: NOC Department is associated with 3 rotating shifts via `shift_ids` in the Department API.

2. **RosterAssignment**: Each employee in NOC has a RosterAssignment that links them to a specific shift for a given date.

3. **Late Threshold Logic**:
   - The system joins `roster_assignments` → `shifts` to get the `late_threshold_time`
   - If `late_threshold_time` is not set on the shift, it falls back to `start_time`
   - Example: NOC Morning shift employee checking in at 07:35 is NOT late (threshold is 07:30), but checking in at 07:31 IS late

4. **Attendance Reports**:
   - `department-stats` endpoint returns late counts based on shift-specific thresholds
   - `overview` endpoint accepts `late_threshold` parameter to override default
   - Export includes department name for filtering

### Example: Creating NOC Department with Shifts

```bash
curl -X POST http://localhost:8080/api/v1/departments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "NOC",
    "name": "NOC Department",
    "level": "department",
    "department_type": "operational",
    "manager_id": 5,
    "location_id": 1,
    "shift_ids": [1, 2, 3],
    "is_active": true
  }'
```

### Example: NOC Shift with Late Exception

```bash
curl -X POST http://localhost:8080/api/v1/shifts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "shift_name": "NOC Morning",
    "shift_code": "NOC-M",
    "shift_type": "rotating",
    "start_time": "07:00:00",
    "end_time": "15:00:00",
    "grace_minutes": 30,
    "break_duration": 60,
    "late_threshold_time": "07:30:00",
    "is_active": true
  }'
```

---

# Summary of Changes

## API Endpoints Updated

| Endpoint | Method | Change |
|----------|--------|--------|
| `/attendance/reports/department-stats` | GET | Now uses shift-specific `LateThresholdTime` via RosterAssignment |
| `/attendance/reports/by-department` | GET | Returns department-wise attendance with late counts |
| `/attendance/reports/overview` | GET | Supports department filtering and custom late thresholds |
| `/attendance/reports/export` | POST | Export now includes `department_name` column |
| `/attendance/reports/employee-comprehensive` | GET | Returns department for each employee |

## Key Features Added

1. **Shift-Specific Late Thresholds**: Late calculation uses `LateThresholdTime` from the employee's assigned shift
2. **Department Tracking**: Attendance data is linked to employee departments
3. **NOC Support**: Different shifts and thresholds for NOC employees
4. **Export Enhancement**: CSV/XLSX exports include department name
5. **Department Filtering**: Most endpoints support filtering by `departmentId`

## Frontend Integration Points

1. **Department Dropdown**: Add department filter to attendance dashboard
2. **Late Threshold Display**: Show shift-specific late threshold in employee details
3. **NOC Badge**: Display NOC badge for NOC department employees
4. **Export**: Support new export format with department column

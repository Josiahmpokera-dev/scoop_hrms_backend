# Attendance Reports & Analytics API

This document outlines the API endpoints required to power the **Attendance Reports & Analytics** page. These endpoints enable fetching key performance indicators (KPIs), trend analysis, department-wise breakdowns, and compliance alerts.

## Base URL
`http://localhost:8080/api/v1`

## Authentication
All endpoints require a valid Bearer token in the Authorization header.
`Authorization: Bearer <token>`

---

## 1. Get Attendance Summary (KPIs)

Fetches the high-level Key Performance Indicators for attendance within a specified date range.

### Endpoint
`GET /attendance/reports/summary`

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |
| `departmentId` | string | No | Filter by specific department ID |

### Response Body
```json
{
  "success": true,
  "data": {
    "presentRate": 92.5,
    "absentRate": 3.2,
    "lateRate": 4.3,
    "averageWorkingHours": 8.7,
    "overtimeHours": 145,
    "exceptionRate": 2.1,
    "onTimeRate": 95.7,
    "shiftAdherence": 97.3,
    "totalEmployees": 150,
    "trends": {
        "presentChange": 2.3, // Percentage change compared to previous period
        "absentChange": -1.1,
        "lateChange": 0.8
    }
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "INVALID_DATE_RANGE",
    "message": "End date cannot be before start date."
  }
}
```

---

## 2. Get Attendance Trends

Retrieves daily or weekly attendance trends for visualization (e.g., bar/line charts).

### Endpoint
`GET /attendance/reports/trends`

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |
| `departmentId` | string | No | Filter by specific department |
| `interval` | string | No | Data grouping: `daily` (default), `weekly`, `monthly` |

### Response Body
```json
{
  "success": true,
  "data": [
    {
      "date": "2024-03-01",
      "day": "Mon",
      "present": 95,
      "absent": 5,
      "late": 3,
      "onLeave": 2,
      "totalScheduled": 105
    },
    {
      "date": "2024-03-02",
      "day": "Tue",
      "present": 93,
      "absent": 4,
      "late": 3,
      "onLeave": 3,
      "totalScheduled": 103
    }
  ]
}
```

---

## 3. Get Department-wise Attendance

Provides a breakdown of attendance statistics aggregated by department.

### Endpoint
`GET /attendance/reports/by-department`

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `date` | string | Yes | The specific date to analyze (`YYYY-MM-DD`) |
| `organizationId` | string | No | Filter by organization unit |

### Response Body
```json
{
  "success": true,
  "data": [
    {
      "departmentId": "dept_001",
      "departmentName": "Engineering",
      "present": 45,
      "absent": 2,
      "late": 3,
      "total": 50,
      "presentRate": 90.0
    },
    {
      "departmentId": "dept_002",
      "departmentName": "Sales",
      "present": 28,
      "absent": 1,
      "late": 1,
      "total": 30,
      "presentRate": 93.3
    }
  ]
}
```

---

## 4. Get Compliance Alerts & Violations

Fetches a list of attendance compliance issues, such as exceeding working hours, missing rest days, or consistent lateness.

### Endpoint
`GET /attendance/reports/compliance`

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |
| `departmentId` | string | No | Filter by department |
| `severity` | string | No | Filter by severity: `critical`, `high`, `medium`, `low` |

### Response Body
```json
{
  "success": true,
  "data": [
    {
      "id": "viol_123",
      "employeeId": "emp_456",
      "employeeName": "John Doe",
      "department": "Engineering",
      "violationType": "Weekly Hours Exceeded",
      "violationDate": "2026-03-15",
      "details": "Worked 52 hours in week 42 (Limit: 45h)",
      "severity": "High",
      "status": "Open"
    },
    {
      "id": "viol_124",
      "employeeId": "emp_789",
      "employeeName": "Ashley Martinez",
      "department": "Sales",
      "violationType": "No Rest Day",
      "violationDate": "2026-03-14",
      "details": "Worked 7 consecutive days without a break",
      "severity": "Critical",
      "status": "Open"
    }
  ]
}
```

---

## 5. Get Overtime Analysis

Detailed breakdown of overtime hours and costs.

### Endpoint
`GET /attendance/reports/overtime-analysis`

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | Yes | Start date in `YYYY-MM-DD` format |
| `endDate` | string | Yes | End date in `YYYY-MM-DD` format |

### Response Body
```json
{
  "success": true,
  "data": {
    "totalOvertimeHours": 145.5,
    "overtimeCost": 4200000,
    "employeesWithOvertime": 12,
    "topOvertimeDepartments": [
      {
        "departmentName": "Engineering",
        "hours": 65.5
      }
    ],
    "overtimeTrends": [
      { "date": "2026-03-01", "hours": 10 }
    ]
  }
}
```

---

## 6. Export Report

Triggers the generation of a downloadable report file (PDF, CSV, Excel).

### Endpoint
`POST /attendance/reports/export`

### Request Body
```json
{
  "reportType": "daily", // Supported: daily (alias: daily_attendance)
  "startDate": "2024-03-01",
  "endDate": "2024-03-31",
  "departmentId": "12", // Optional (numeric department id as string)
  "format": "csv" // Options: csv, xlsx
}
```

### Response Body
```json
{
  "success": true,
  "message": "Report generation started",
  "data": {
    "downloadUrl": "http://localhost:8080/storage/reports/attendance_daily_2024-03-01_1711773180.csv",
    "jobId": "job_998877",
    "expiresAt": "2026-04-10T10:15:00.000Z"
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "EXPORT_FAILED",
    "message": "Unable to generate report. Data unavailable for selected range."
  }
}
```

---

## Common Error Codes

| Code | Description |
|------|-------------|
| `UNAUTHORIZED` | Invalid or missing authentication token |
| `FORBIDDEN` | User does not have permission to view these reports |
| `INVALID_DATE_RANGE` | Start date is after End date or range is too large |
| `RESOURCE_NOT_FOUND` | Department or Employee ID not found |
| `INTERNAL_SERVER_ERROR` | Unexpected system error |

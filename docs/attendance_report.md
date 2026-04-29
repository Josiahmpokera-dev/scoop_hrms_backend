# Attendance Reports & Analytics API (Actual Integration)

This is the rebuilt API documentation for the **Attendance Reports & Analytics** page based on the current implementation in:
- `src/app/(protected-pages)/attendance/reports/page.tsx`
- `src/services/AttendanceService.ts`

## Base URL
`/api/v1`

## Auth
All endpoints require Bearer token:
`Authorization: Bearer <token>`

## Standard Envelope
```json
{
  "success": true,
  "message": "string",
  "data": {}
}
```

---

## UI -> API Mapping

| UI Section | API(s) Used |
|---|---|
| KPIs | `GET /attendance/reports/summary` |
| Weekly Trend | `GET /attendance/reports/trends` |
| Department Attendance | `GET /attendance/reports/department-stats` |
| Overtime & Compliance | `GET /attendance/reports/overtime`, `GET /attendance/reports/compliance` |
| Attendance By Employee | `GET /attendance/reports/employee-stats` |
| Generate Report (header/export button) | `POST /attendance/reports/export` |
| Generate Report tab (employee comprehensive) | `GET /attendance/reports/employee-comprehensive` |
| Monthly Reports -> Generate | `POST /attendance/reports/trigger-monthly?month=YYYY-MM` |
| Monthly Reports -> Email Logs | `GET /attendance/reports/email-logs` |

---

## 1) Attendance Summary (KPIs)

### Endpoint
`GET /attendance/reports/summary`

### Query
- `startDate` (required, `YYYY-MM-DD`)
- `endDate` (required, `YYYY-MM-DD`)
- `departmentId` (optional, number)

### Example
`GET /attendance/reports/summary?startDate=2026-04-01&endDate=2026-04-30&departmentId=5`

### Response (example)
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

## 2) Attendance Trends

### Endpoint
`GET /attendance/reports/trends`

### Query
- `startDate` (required)
- `endDate` (required)
- `departmentId` (optional)
- `interval` (optional: `daily`/`weekly`)

### Example
`GET /attendance/reports/trends?startDate=2026-04-01&endDate=2026-04-30&departmentId=5`

### Response (example)
```json
{
  "success": true,
  "data": [
    {
      "date": "2026-04-01",
      "day": "Wednesday",
      "present": 115,
      "absent": 3,
      "late": 2,
      "onLeave": 0,
      "totalScheduled": 120
    }
  ]
}
```

---

## 3) Department Attendance Stats (Date Range)

### Endpoint
`GET /attendance/reports/department-stats`

### Query
- `startDate` (required)
- `endDate` (required)

### Example
`GET /attendance/reports/department-stats?startDate=2026-04-01&endDate=2026-04-30`

### Response (example)
```json
{
  "success": true,
  "data": [
    {
      "departmentId": "5",
      "departmentName": "NOC Department",
      "present": 25,
      "absent": 2,
      "late": 1,
      "onLeave": 0,
      "attendancePercentage": 89.29
    }
  ]
}
```

---

## 4) Employee Attendance Stats

### Endpoint
`GET /attendance/reports/employee-stats`

### Query
- `employeeId` (required, number)
- `startDate` (required)
- `endDate` (required)

### Example
`GET /attendance/reports/employee-stats?employeeId=101&startDate=2026-04-01&endDate=2026-04-30`

### Response (example)
```json
{
  "success": true,
  "data": {
    "employee_id": "EMP001",
    "employee_name": "John Doe",
    "department": "NOC",
    "position": "Engineer",
    "period": { "start_date": "2026-04-01", "end_date": "2026-04-30" },
    "attendance": {
      "total_working_days": 22,
      "present_days": 20,
      "absent_days": 1,
      "late_days": 1,
      "average_working_hours": 8.3,
      "total_working_hours": 182.6
    },
    "timesheet": { "overtime_hours": 12.5 },
    "leave": { "total_leave_days": 1, "leave_breakdown": [] },
    "compliance": { "violations_count": 0, "violations": [] }
  }
}
```

---

## 5) Employee Comprehensive Report

### Endpoint
`GET /attendance/reports/employee-comprehensive`

### Query
- `startDate` (required)
- `endDate` (required)
- `employeeId` (optional)
- `departmentId` (optional)
- `locationId` (optional)
- `page` (optional)
- `pageSize` (optional)

### Example
`GET /attendance/reports/employee-comprehensive?startDate=2026-04-01&endDate=2026-04-30&page=1&pageSize=20`

### Response (example)
```json
{
  "success": true,
  "data": [
    {
      "employee_id": "EMP001",
      "employee_name": "John Doe",
      "department": "NOC",
      "position": "Engineer",
      "attendance": {
        "present_percentage": 90.9,
        "average_working_hours": 8.2,
        "total_working_hours": 180.4
      },
      "overtime": { "approved_hours": 10.5 },
      "leave": { "total_leave_days": 1 },
      "compliance": { "violations_count": 0 }
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 42, "total_pages": 3 }
}
```

---

## 6) Compliance Violations

### Endpoint
`GET /attendance/reports/compliance`

### Query
- `startDate` (required)
- `endDate` (required)
- `departmentId` (optional)
- `severity` (optional)

### Example
`GET /attendance/reports/compliance?startDate=2026-04-01&endDate=2026-04-30&departmentId=5`

### Response (example)
```json
{
  "success": true,
  "data": [
    {
      "id": "V-1001",
      "employeeId": "EMP001",
      "employeeName": "John Doe",
      "department": "NOC",
      "violationType": "Late Check-in",
      "violationDate": "2026-04-11",
      "severity": "high",
      "details": "Checked in 45 minutes late",
      "status": "Open"
    }
  ]
}
```

---

## 7) Overtime Analysis

### Endpoint
`GET /attendance/reports/overtime`

### Query
- `startDate` (required)
- `endDate` (required)

### Example
`GET /attendance/reports/overtime?startDate=2026-04-01&endDate=2026-04-30`

### Response (example)
```json
{
  "success": true,
  "data": {
    "totalOvertimeHours": 145.5,
    "averageOvertimePerEmployee": 2.9,
    "estimatedCost": 4200000,
    "currency": "TZS",
    "topContributors": [
      { "employeeId": "EMP001", "name": "John Doe", "hours": 12.5 }
    ]
  }
}
```

---

## 8) Export Attendance Report

### Endpoint
`POST /attendance/reports/export`

### Request Body (as implemented)
```json
{
  "reportType": "daily",
  "startDate": "2026-04-01",
  "endDate": "2026-04-30",
  "departmentId": 5,
  "format": "pdf"
}
```

### Response (example)
```json
{
  "success": true,
  "message": "Report generation started",
  "data": {
    "downloadUrl": "https://.../attendance_report_2026-04.pdf",
    "jobId": "job_998877",
    "expiresAt": "2026-05-01T12:00:00Z"
  }
}
```

### UI Flow
1. User clicks **Export Report**.
2. UI calls `POST /attendance/reports/export`.
3. If `data.downloadUrl` exists, UI starts immediate browser download.

---

## 9) Trigger Monthly Report

### Endpoint
`POST /attendance/reports/trigger-monthly?month=YYYY-MM`

### Example
`POST /attendance/reports/trigger-monthly?month=2026-04`

### Response (example)
```json
{
  "success": true,
  "message": "Monthly report generation triggered successfully",
  "data": null
}
```

### UI Flow
1. User selects month in Monthly Reports tab.
2. UI calls trigger endpoint.
3. Backend generates/sends report asynchronously (email flow).

---

## 10) Monthly Report Email Logs

### Endpoint
`GET /attendance/reports/email-logs`

### Query
- `page` (optional)
- `page_size` (optional)

### Example
`GET /attendance/reports/email-logs?page=1&page_size=15`

### Response (example)
```json
{
  "success": true,
  "data": [
    {
      "id": 1001,
      "recipient": "hr@scoopworks.com",
      "subject": "Monthly Attendance Report - 2026-04",
      "report_month": "2026-04",
      "status": "Delivered",
      "date": "2026-05-01T08:10:00Z",
      "sent_at": "2026-05-01T08:10:20Z",
      "download_url": "https://.../report.xlsx",
      "file_type": "xlsx",
      "file_size": 248721
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 15,
    "total": 42,
    "total_pages": 3
  }
}
```

### UI Flow
1. User opens **Monthly Reports -> Email Logs** tab.
2. UI calls email logs endpoint with pagination.
3. Download button opens `download_url` in a new tab when available.

---

## Common Error Responses

```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "field": "reason"
  }
}
```

Common statuses:
- `400` validation error
- `401` unauthorized
- `403` forbidden
- `404` resource not found
- `500` server error

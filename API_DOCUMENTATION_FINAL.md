# 🚀 HRMS Backend API - Complete Documentation

## Overview

**Base URL:** `http://localhost:8080`  
**API Version:** v1  
**Total Endpoints:** 200+  
**Authentication:** Bearer Token (JWT)

---

## 🔐 Authentication

### Bearer Token Authentication

Most API endpoints require Bearer token authentication. Include the token in the Authorization header:

```bash
Authorization: Bearer <your_access_token>
```

### Getting a Token

**Login Endpoint:**
```bash
POST /api/v1/auth/login
```

**Request:**
```json
{
  "email": "admin@hrms.com",
  "password": "Admin@2024!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@hrms.com",
      "first_name": "Admin",
      "last_name": "User",
      "role": "admin"
    }
  }
}
```

### Public Endpoints (No Authentication Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | User login |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| POST | `/api/v1/attendance/manual-punches` | Create manual punch |

### Role-Based Access Control

- 🔓 **Public** - No authentication required
- 🔒 **Authenticated** - Any authenticated user
- 👥 **HR/Admin** - HR or Admin users only
- ⚠️ **Admin** - Admin users only

---

## 📊 Key Endpoints

### 1. Comprehensive Employee Report ⭐

**Endpoint:**
```
GET /api/v1/attendance/reports/employee-comprehensive
```

**Authentication:** Bearer Token (HR/Admin)

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_date | string | ✅ Yes | Start date (YYYY-MM-DD) |
| end_date | string | ✅ Yes | End date (YYYY-MM-DD) |
| employee_id | int | ❌ No | Filter by employee |
| department_id | int | ❌ No | Filter by department |
| location_id | int | ❌ No | Filter by location |
| page | int | ❌ No | Page number (default: 1) |
| page_size | int | ❌ No | Items per page (max: 100) |

**Example Usage:**

```bash
# 1 Week Report
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/attendance/reports/employee-comprehensive?start_date=2026-04-21&end_date=2026-04-27"

# 1 Month Report
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/attendance/reports/employee-comprehensive?start_date=2026-04-01&end_date=2026-04-30"

# 3 Months Report
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/attendance/reports/employee-comprehensive?start_date=2026-01-01&end_date=2026-03-31"
```

**Response Includes:**
- ✅ Attendance: Present/absent/late rates, working hours
- ✅ Timesheet: Regular/overtime/billable hours
- ✅ Leave: Leave days, breakdown by type
- ✅ Overtime: Requests, approved hours, payouts
- ✅ Compliance: Violations count and details

---

### 2. Attendance Management

#### Timesheets
- `POST /api/v1/attendance/timesheets/entries` - Create timesheet entry
- `PUT /api/v1/attendance/timesheets/entries/{id}` - Update entry
- `DELETE /api/v1/attendance/timesheets/entries/{id}` - Delete entry
- `GET /api/v1/attendance/timesheets` - Get my timesheets
- `POST /api/v1/attendance/timesheets/{id}/submit` - Submit for approval

#### Overtime
- `POST /api/v1/attendance/overtime/requests` - Create OT request
- `GET /api/v1/attendance/overtime/requests` - Get my OT requests
- `POST /api/v1/attendance/overtime/admin/requests/{id}/approve` - Approve OT (Admin)

#### Manual Punches
- `POST /api/v1/attendance/manual-punches` - Create manual punch
- `GET /api/v1/attendance/manual-punches/my` - Get my manual punches
- `PATCH /api/v1/attendance/manual-punches/{id}/status` - Update status (Admin)

#### Reports
- `GET /api/v1/attendance/reports/summary` - Attendance summary
- `GET /api/v1/attendance/reports/trends` - Attendance trends
- `GET /api/v1/attendance/reports/compliance` - Compliance violations

---

### 3. Employee Management

#### Core Operations
- `POST /api/v1/employees` - Onboard new employee
- `GET /api/v1/employees` - List all employees
- `GET /api/v1/employees/{id}` - Get employee details
- `PUT /api/v1/employees/{id}` - Update employee
- `DELETE /api/v1/employees/{id}` - Delete employee

#### Self Service
- `GET /api/v1/self-service/profile` - Get my profile
- `GET /api/v1/employees/onboarding/non-employee-users` - List non-employee users

#### People Directory
- `GET /api/v1/employees/directory` - Search directory
- `GET /api/v1/employees/directory/filters` - Get filter options

---

### 4. Leave Management

#### Leave Requests
- `POST /api/v1/leave/requests` - Apply for leave
- `GET /api/v1/leave/requests/my` - Get my leave requests
- `POST /api/v1/leave/requests/{id}/cancel` - Cancel leave request
- `POST /api/v1/leave/requests/{id}/approve` - Approve leave (Admin)

#### Leave Balance
- `GET /api/v1/leave/balance` - Get my leave balance
- `GET /api/v1/leave/types` - Get leave types

#### Calendar
- `GET /api/v1/leave/calendar` - Get leave calendar
- `GET /api/v1/leave/calendar/today` - Get today's leaves

---

### 5. Payroll Management

#### Payroll Runs
- `GET /api/v1/payroll/runs` - List payroll runs
- `POST /api/v1/payroll/runs` - Create payroll run
- `GET /api/v1/payroll/runs/{id}` - Get payroll run details

#### Payslips
- `GET /api/v1/payroll/payslips` - List payslips
- `GET /api/v1/payroll/payslips/{id}` - Get payslip
- `GET /api/v1/payroll/payslips/{id}/download` - Download payslip

#### Self Service
- `GET /api/v1/self-service/payroll/payslips` - Get my payslips
- `GET /api/v1/self-service/payroll/payslips/latest` - Get latest payslip

---

### 6. Biometric Integration

#### Device Management
- `GET /api/v1/biometric/biotime/test-connection` - Test connection
- `GET /api/v1/biometric/biotime/terminals` - List terminals
- `GET /api/v1/biometric/biotime/transactions` - Get transactions

#### Attendance
- `GET /api/v1/biometric/attendance/daily` - Get daily attendance
- `GET /api/v1/biometric/attendance/merged` - Get merged attendance

#### Enrollments
- `GET /api/v1/biometric/enrollments` - List enrollments
- `POST /api/v1/biometric/enrollments/link` - Link employee to device
- `POST /api/v1/biometric/enrollments/auto-link` - Auto-link employees

---

### 7. Dashboard & Analytics

#### Dashboard
- `GET /api/v1/dashboard/statistics` - Get dashboard statistics
- `GET /api/v1/dashboard/announcements` - Get announcements
- `GET /api/v1/dashboard/my-activity` - Get my recent activity
- `GET /api/v1/dashboard/pending-approvals` - Get pending approvals

#### Reports
- `GET /api/v1/attendance/reports/employee-utilization` - Employee utilization
- `GET /api/v1/attendance/reports/project-utilization` - Project utilization
- `GET /api/v1/payroll/reports/summary` - Payroll summary

---

## 📖 Interactive Documentation

### Swagger UI
Access interactive API documentation with testing capabilities:

**URL:** http://localhost:8080/swagger/index.html

**Features:**
- Try out endpoints directly in browser
- View request/response schemas
- Test with Bearer token authentication
- Export specifications

### HTML Documentation
Access readable API documentation:

**URL:** http://localhost:8080/docs/index.html

**Features:**
- Beautiful, responsive design
- Syntax highlighting
- Complete examples
- Mobile-friendly

---

## 🔧 Error Handling

### Error Response Format
```json
{
  "success": false,
  "message": "Error description",
  "error": "Error details"
}
```

### Common Error Codes

| Code | Status | Description |
|------|--------|-------------|
| 400 | Bad Request | Invalid input or parameters |
| 401 | Unauthorized | Missing or invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 422 | Unprocessable Entity | Validation errors |
| 500 | Internal Server Error | Server error |

---

## 📐 Pagination

Most list endpoints support pagination:

**Query Parameters:**
- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 20, max: 100)

**Response Format:**
```json
{
  "success": true,
  "message": "Items retrieved successfully",
  "data": [...],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

## 🚀 Quick Start

### 1. Start the Server
```bash
make run
```

### 2. Get Access Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@hrms.com","password":"Admin@2024!"}'
```

### 3. Test an Endpoint
```bash
# Set token
TOKEN="your_access_token_here"

# Get comprehensive report
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/attendance/reports/employee-comprehensive?start_date=2026-04-01&end_date=2026-04-30"
```

### 4. Explore Documentation
- Swagger UI: http://localhost:8080/swagger/index.html
- HTML Docs: http://localhost:8080/docs/index.html

---

## 🎯 API Categories

### Authentication & Authorization
- User registration, login, logout
- Token refresh and validation
- Role-based access control

### Employee Management
- Employee CRUD operations
- Onboarding workflows
- People directory
- Self-service profile

### Attendance Management
- Timesheet tracking
- Overtime requests
- Manual punch requests
- Attendance reports

### Leave Management
- Leave request workflow
- Leave balance tracking
- Leave calendar
- Leave types and policies

### Payroll Management
- Payroll run processing
- Payslip generation
- Salary structures
- Loan management

### Biometric Integration
- Device connectivity
- Transaction sync
- Employee enrollment
- Attendance merge

### Dashboard & Analytics
- KPI statistics
- Reports and analytics
- Pending approvals
- Activity tracking

### Projects & Tasks
- Project management
- Daily task tracking
- Task assignments
- Progress tracking

### Recruitment
- Job openings
- Candidate management
- Interview scheduling
- Offer management

### Assets & Inventory
- Asset tracking
- Asset assignments
- Asset lifecycle

### Helpdesk
- Ticket management
- Knowledge base
- CSAT tracking

---

## 🔒 Security Best Practices

1. **Always use HTTPS** in production
2. **Store tokens securely** - Use secure storage, never in localStorage
3. **Implement token refresh** - Use refresh tokens for long-lived sessions
4. **Validate input** - Always validate and sanitize user input
5. **Use role-based access** - Implement proper authorization checks
6. **Audit logs** - Track all sensitive operations
7. **Rate limiting** - Prevent abuse and DDoS attacks
8. **CORS configuration** - Restrict origins appropriately

---

## 📞 Support

For issues or questions:
- Check Swagger UI for detailed endpoint documentation
- Review HTML documentation for examples
- Verify authentication token is valid
- Check server logs for error details

---

## 📄 License

MIT License - See LICENSE file for details

---

**Status:** ✅ Production Ready  
**Version:** 1.0.0  
**Last Updated:** 2026-04-27

---
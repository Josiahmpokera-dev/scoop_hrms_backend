# Payroll Module - Backend API Specification

**Version:** 2.0  
**Last Updated:** January 23, 2026  
**Target:** Backend Developers  
**Currency:** TZS (Tanzanian Shilling)  
**Country Focus:** Tanzania (with multi-country structure ready)

---

## Overview

This document provides complete API specifications for implementing the Payroll module backend. It covers all 7 sub-modules:

1. **Payroll Dashboard** - Overview & quick actions
2. **Payroll Run** - Monthly payroll processing (6-stage workflow)
3. **Salary Structure** - Compensation templates & components
4. **Payslips** - View & download employee payslips
5. **Loans & Advances** - Loan management with EMI
6. **Compliance & Tax** - Tanzania statutory compliance (PAYE, NSSF, NHIF, SDL, WCF)
7. **Payroll Reports** - Analytics & exports

---

## Authentication

All endpoints require Bearer token authentication:

```
Authorization: Bearer {access_token}
```

---

## Standard Response Envelope

### Success Response
```json
{
  "success": true,
  "message": "Human-readable success message",
  "data": { ... }
}
```

### Error Response (4xx/5xx)
```json
{
  "success": false,
  "message": "Error description",
  "errors": {
    "field_name": ["Validation error message"]
  }
}
```

### Paginated Response
```json
{
  "success": true,
  "message": "Items retrieved successfully",
  "data": {
    "data": [ ... ],
    "meta": {
      "page": 1,
      "per_page": 10,
      "total": 125,
      "total_pages": 13
    }
  }
}
```

---

## API Endpoints Summary

| # | Method | Endpoint | Description |
|---|--------|----------|-------------|
| **DASHBOARD** ||||
| 1.1 | GET | `/payroll/dashboard` | Get dashboard overview |
| **PAYROLL RUN** ||||
| 2.1 | GET | `/payroll/runs` | List payroll runs |
| 2.2 | GET | `/payroll/runs/:run_id` | Get payroll run details |
| 2.3 | POST | `/payroll/runs` | Create new payroll run |
| 2.4 | GET | `/payroll/runs/:run_id/pre-check` | Get pre-check validation |
| 2.5 | POST | `/payroll/runs/:run_id/step` | Execute payroll step |
| 2.6 | GET | `/payroll/runs/:run_id/employees` | Get employees in run |
| 2.7 | DELETE | `/payroll/runs/:run_id` | Delete draft payroll run |
| **SALARY STRUCTURE** ||||
| 3.1 | GET | `/payroll/salary-structures` | List salary templates |
| 3.2 | POST | `/payroll/salary-structures` | Create salary template |
| 3.3 | GET | `/payroll/salary-structures/:id` | Get salary template |
| 3.4 | PUT | `/payroll/salary-structures/:id` | Update salary template |
| 3.5 | DELETE | `/payroll/salary-structures/:id` | Delete salary template |
| 3.6 | GET | `/payroll/salary-components` | List salary components |
| 3.7 | POST | `/payroll/salary-components` | Create salary component |
| 3.8 | PUT | `/payroll/salary-components/:id` | Update salary component |
| 3.9 | DELETE | `/payroll/salary-components/:id` | Delete salary component |
| **PAYSLIPS** ||||
| 4.1 | GET | `/payroll/payslips` | List payslips |
| 4.2 | GET | `/payroll/payslips/:id` | Get payslip detail |
| 4.3 | GET | `/payroll/payslips/:id/download` | Download payslip PDF |
| 4.4 | POST | `/payroll/payslips/:id/email` | Email payslip |
| 4.5 | GET | `/payroll/payslips/summary` | Get payslip summary |
| 4.6 | GET | `/payroll/payslips/my` | Get current user's payslip |
| **LOANS & ADVANCES** ||||
| 5.1 | GET | `/payroll/loans` | List loans & advances |
| 5.2 | GET | `/payroll/loans/:id` | Get loan detail |
| 5.3 | POST | `/payroll/loans` | Request loan/advance |
| 5.4 | POST | `/payroll/loans/:id/approve` | Approve loan |
| 5.5 | POST | `/payroll/loans/:id/reject` | Reject loan |
| 5.6 | GET | `/payroll/loans/summary` | Get loans summary |
| **COMPLIANCE & TAX** ||||
| 6.1 | GET | `/payroll/compliance/summary` | Get compliance summary |
| 6.2 | GET | `/payroll/compliance/tax-slabs` | List PAYE tax slabs |
| 6.3 | GET | `/payroll/compliance/statutory` | Get statutory rules |
| 6.4 | GET | `/payroll/compliance/returns/:type` | Download returns file |
| **REPORTS** ||||
| 7.1 | GET | `/payroll/reports/filters` | Get report filter options |
| 7.2 | GET | `/payroll/reports/kpis` | Get report KPIs |
| 7.3 | GET | `/payroll/reports/department-cost` | Get department cost analysis |
| 7.4 | GET | `/payroll/reports/export` | Export report file |
| **PAYMENT INTEGRATION** ||||
| 8.1 | GET | `/payroll/runs/:run_id/disbursement` | Get disbursement batch |
| 8.2 | GET | `/payroll/runs/:run_id/disbursement/file` | Download bank payment file |
| 8.3 | POST | `/payroll/disbursement/callback` | Payment status callback |

---

## 1. PAYROLL DASHBOARD

### 1.1 Get Dashboard Overview

**Endpoint:** `GET /payroll/dashboard`

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `limit` | number | No | 10 | Max recent runs to return |

**Response:**
```json
{
  "success": true,
  "message": "Dashboard data retrieved successfully",
  "data": {
    "current_run": {
      "id": "run-002",
      "run_name": "November 2024 Payroll",
      "pay_period": "01 Nov - 30 Nov 2024",
      "pay_month": 11,
      "pay_year": 2024,
      "pay_frequency": "Monthly",
      "status": "Draft",
      "current_step": 0,
      "total_employees": 125,
      "total_gross": 187500000,
      "total_deductions": 39375000,
      "total_net": 148125000,
      "cutoff_date": "2024-11-30",
      "disbursement_date": "2024-12-05",
      "created_by": "HR Admin",
      "created_at": "2024-11-01T00:00:00+03:00"
    },
    "last_finalized_run": {
      "id": "run-001",
      "run_name": "October 2024 Payroll",
      "pay_period": "01 Oct - 31 Oct 2024",
      "status": "Finalized",
      "total_employees": 125,
      "total_net": 148125000,
      "disbursement_date": "2024-11-05"
    },
    "recent_runs": [
      {
        "id": "run-001",
        "run_name": "October 2024 Payroll",
        "pay_period": "01 Oct - 31 Oct 2024",
        "pay_month": 10,
        "pay_year": 2024,
        "status": "Finalized",
        "total_employees": 125,
        "total_net": 148125000
      }
    ],
    "statistics": {
      "last_month_payout": 148125000,
      "employees_paid": 125,
      "tax_deducted": 38000000,
      "statutory_dues": 52000000,
      "pending_exceptions": 3
    }
  }
}
```

**Notes:**
- `current_run` is `null` when no run is in Draft/In Review status
- `statistics` shows previous month's summary

---

## 2. PAYROLL RUN

### 2.1 List Payroll Runs

**Endpoint:** `GET /payroll/runs`

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | number | No | 1 | Page number |
| `page_size` | number | No | 10 | Items per page (max 100) |
| `status` | string | No | - | Filter: `Draft`, `In Review`, `Approved`, `Finalized`, `Disbursed`, `Closed` |
| `pay_year` | number | No | - | Filter by year (e.g. 2024) |
| `pay_month` | number | No | - | Filter by month (1-12) |

**Response:**
```json
{
  "success": true,
  "message": "Payroll runs retrieved successfully",
  "data": {
    "data": [
      {
        "id": "run-002",
        "run_name": "November 2024 Payroll",
        "pay_period": "01 Nov - 30 Nov 2024",
        "pay_month": 11,
        "pay_year": 2024,
        "pay_frequency": "Monthly",
        "status": "Draft",
        "current_step": 0,
        "total_employees": 125,
        "total_gross": 187500000,
        "total_deductions": 39375000,
        "total_net": 148125000,
        "cutoff_date": "2024-11-30",
        "disbursement_date": "2024-12-05",
        "created_by": "HR Admin",
        "approved_by": null,
        "finalized_by": null,
        "created_at": "2024-11-01T00:00:00+03:00",
        "updated_at": "2024-11-15T00:00:00+03:00"
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 10,
      "total": 12,
      "total_pages": 2
    }
  }
}
```

---

### 2.2 Get Payroll Run Details

**Endpoint:** `GET /payroll/runs/:run_id`

**Response:**
```json
{
  "success": true,
  "message": "Payroll run retrieved successfully",
  "data": {
    "id": "run-002",
    "run_name": "November 2024 Payroll",
    "pay_period": "01 Nov - 30 Nov 2024",
    "pay_month": 11,
    "pay_year": 2024,
    "pay_frequency": "Monthly",
    "status": "Draft",
    "total_employees": 125,
    "total_gross": 187500000,
    "total_deductions": 39375000,
    "total_net": 148125000,
    "total_employer_contributions": 47000000,
    "cutoff_date": "2024-11-30",
    "disbursement_date": "2024-12-05",
    "created_by": "HR Admin",
    "approved_by": null,
    "finalized_by": null,
    "current_step": 0,
    "steps": [
      { "index": 0, "title": "Pre-Check", "description": "Validate data", "status": "current" },
      { "index": 1, "title": "Calculate", "description": "Process salaries", "status": "pending" },
      { "index": 2, "title": "Review", "description": "Verify calculations", "status": "pending" },
      { "index": 3, "title": "Approve", "description": "Get approvals", "status": "pending" },
      { "index": 4, "title": "Finalize", "description": "Lock payroll", "status": "pending" },
      { "index": 5, "title": "Disburse", "description": "Process payments", "status": "pending" }
    ],
    "breakdown": {
      "total_basic": 125000000,
      "total_allowances": 62500000,
      "total_paye": 38000000,
      "total_nssf_employee": 23000000,
      "total_nhif_employee": 5000000,
      "total_loan_deductions": 8500000,
      "total_nssf_employer": 23000000,
      "total_nhif_employer": 5000000,
      "total_sdl": 11500000,
      "total_wcf": 2300000
    },
    "created_at": "2024-11-01T00:00:00+03:00",
    "updated_at": "2024-11-15T00:00:00+03:00"
  }
}
```

---

### 2.3 Create Payroll Run

**Endpoint:** `POST /payroll/runs`

**Request Body:**
```json
{
  "run_name": "December 2024 Payroll",
  "pay_month": 12,
  "pay_year": 2024,
  "pay_frequency": "Monthly",
  "cutoff_date": "2024-12-31",
  "disbursement_date": "2025-01-05"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `run_name` | string | Yes | Display name for the run |
| `pay_month` | number | Yes | Month (1-12) |
| `pay_year` | number | Yes | Year (e.g. 2024) |
| `pay_frequency` | string | Yes | `Monthly` \| `Bi-Weekly` \| `Weekly` |
| `cutoff_date` | string | Yes | ISO date (YYYY-MM-DD) |
| `disbursement_date` | string | Yes | ISO date (YYYY-MM-DD) |

**Response:**
```json
{
  "success": true,
  "message": "Payroll run created successfully",
  "data": {
    "id": "run-003",
    "run_name": "December 2024 Payroll",
    "pay_period": "01 Dec - 31 Dec 2024",
    "pay_month": 12,
    "pay_year": 2024,
    "pay_frequency": "Monthly",
    "status": "Draft",
    "current_step": 0,
    "total_employees": 0,
    "total_gross": 0,
    "total_deductions": 0,
    "total_net": 0,
    "cutoff_date": "2024-12-31",
    "disbursement_date": "2025-01-05",
    "created_by": "HR Admin",
    "created_at": "2024-11-20T00:00:00+03:00"
  }
}
```

---

### 2.4 Get Pre-Check Validation

**Endpoint:** `GET /payroll/runs/:run_id/pre-check`

**Description:** Returns validation issues that should be resolved before processing payroll.

**Response:**
```json
{
  "success": true,
  "message": "Pre-check results retrieved",
  "data": {
    "run_id": "run-002",
    "is_valid": false,
    "can_proceed": true,
    "issues": [
      {
        "id": "issue-001",
        "type": "Missing Attendance",
        "category": "attendance",
        "count": 2,
        "severity": "High",
        "employee_ids": ["EMP001", "EMP002"],
        "employees": [
          { "id": "EMP001", "name": "John Doe", "department": "Engineering" },
          { "id": "EMP002", "name": "Jane Smith", "department": "Sales" }
        ],
        "description": "No attendance records for pay period",
        "action_required": "Approve attendance or mark as absent",
        "resolution_url": "/attendance/approvals"
      },
      {
        "id": "issue-002",
        "type": "Unapproved Leave (LOP)",
        "category": "leave",
        "count": 1,
        "severity": "Medium",
        "employee_ids": ["EMP005"],
        "employees": [
          { "id": "EMP005", "name": "Mike Wilson", "department": "HR" }
        ],
        "description": "Leave without pay not approved - will deduct salary",
        "action_required": "Approve leave or adjust",
        "resolution_url": "/leave/approvals"
      },
      {
        "id": "issue-003",
        "type": "Pending OT Approval",
        "category": "overtime",
        "count": 3,
        "severity": "Low",
        "employee_ids": ["EMP010", "EMP011", "EMP012"],
        "description": "Overtime pending approval - will not be included",
        "action_required": "Approve overtime claims",
        "resolution_url": "/attendance/overtime"
      },
      {
        "id": "issue-004",
        "type": "Missing Bank Details",
        "category": "employee",
        "count": 1,
        "severity": "High",
        "employee_ids": ["EMP020"],
        "employees": [
          { "id": "EMP020", "name": "New Employee", "department": "Marketing" }
        ],
        "description": "Employee missing bank account details",
        "action_required": "Update employee profile with bank details",
        "resolution_url": "/employee/EMP020/edit"
      },
      {
        "id": "issue-005",
        "type": "New Joiners",
        "category": "proration",
        "count": 2,
        "severity": "Info",
        "employee_ids": ["EMP025", "EMP026"],
        "description": "New joiners will have prorated salary",
        "action_required": null
      }
    ],
    "summary": {
      "total_issues": 5,
      "high_severity": 2,
      "medium_severity": 1,
      "low_severity": 1,
      "info": 1
    }
  }
}
```

**Severity Levels:**
- `High` - Must be resolved (blocks finalization)
- `Medium` - Should be resolved (warning)
- `Low` - Can proceed (informational)
- `Info` - For awareness only

---

### 2.5 Execute Payroll Step

**Endpoint:** `POST /payroll/runs/:run_id/step`

**Request Body:**
```json
{
  "step_index": 1,
  "action": "calculate"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `step_index` | number | Yes | 0=Pre-Check, 1=Calculate, 2=Review, 3=Approve, 4=Finalize, 5=Disburse |
| `action` | string | No | `calculate`, `submit_for_review`, `approve`, `reject`, `finalize`, `disburse` |
| `remarks` | string | No | Comments for approval/rejection |
| `approver_role` | string | No | `HR`, `Finance`, `Director` (for multi-level approval) |

**Action by Step:**

| Step | Actions | Description |
|------|---------|-------------|
| 0 | `validate`, `skip` | Run pre-check validation |
| 1 | `calculate` | Process all salary calculations |
| 2 | `submit_for_review` | Submit for approval |
| 3 | `approve`, `reject` | Approve or reject run |
| 4 | `finalize` | Lock payroll (immutable) |
| 5 | `disburse` | Generate payment instructions |

**Response (after calculate):**
```json
{
  "success": true,
  "message": "Payroll calculated successfully",
  "data": {
    "run_id": "run-002",
    "status": "Draft",
    "current_step": 2,
    "total_employees": 125,
    "total_gross": 187500000,
    "total_deductions": 39375000,
    "total_net": 148125000,
    "total_employer_contributions": 47000000,
    "calculation_summary": {
      "employees_processed": 125,
      "employees_with_changes": 12,
      "new_joiners_prorated": 2,
      "exits_prorated": 1,
      "lop_applied": 3,
      "overtime_included": 8,
      "loans_deducted": 15
    }
  }
}
```

**Response (after finalize):**
```json
{
  "success": true,
  "message": "Payroll finalized and locked successfully",
  "data": {
    "run_id": "run-002",
    "status": "Finalized",
    "current_step": 5,
    "finalized_by": "Finance Director",
    "finalized_at": "2024-11-25T14:30:00+03:00",
    "is_locked": true
  }
}
```

---

### 2.6 Get Employees in Run

**Endpoint:** `GET /payroll/runs/:run_id/employees`

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | number | No | 1 | Page number |
| `page_size` | number | No | 20 | Items per page |
| `search` | string | No | - | Search by name/employee ID |
| `department_id` | string | No | - | Filter by department |
| `has_changes` | boolean | No | - | Only show employees with changes |

**Response:**
```json
{
  "success": true,
  "message": "Employees retrieved successfully",
  "data": {
    "data": [
      {
        "employee_id": "EMP001",
        "employee_name": "John Doe",
        "employee_photo": "/img/avatars/thumb-1.jpg",
        "department": "Engineering",
        "designation": "Senior Software Engineer",
        "basic_salary": 1500000,
        "gross_salary": 2300000,
        "total_deductions": 710000,
        "net_pay": 1590000,
        "days_worked": 30,
        "lop_days": 0,
        "overtime_hours": 8,
        "has_changes": false,
        "change_reason": null
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 20,
      "total": 125,
      "total_pages": 7
    }
  }
}
```

---

## 3. SALARY STRUCTURE

### 3.1 List Salary Templates

**Endpoint:** `GET /payroll/salary-structures`

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | number | No | 1 | Page number |
| `page_size` | number | No | 10 | Items per page |
| `is_active` | boolean | No | - | Filter by active status |
| `grade` | string | No | - | Filter by grade (e.g. G4) |
| `location` | string | No | - | Filter by location |
| `country` | string | No | - | Filter by country |

**Response:**
```json
{
  "success": true,
  "message": "Salary structures retrieved successfully",
  "data": {
    "data": [
      {
        "id": "str-001",
        "template_name": "G4 - Senior Level - Dar es Salaam",
        "grade": "G4",
        "location": "Dar es Salaam",
        "country": "Tanzania",
        "ctc": 3000000,
        "basic": 1500000,
        "hra": 450000,
        "transport": 100000,
        "medical": 50000,
        "other_allowances": 200000,
        "gross_salary": 2300000,
        "employee_deductions": {
          "paye": 380000,
          "nssf": 230000,
          "nhif": 50000
        },
        "employer_contributions": {
          "nssf": 230000,
          "nhif": 50000,
          "sdl": 115000,
          "wcf": 23000
        },
        "total_deductions": 660000,
        "net_pay": 1640000,
        "is_active": true,
        "employees_count": 25,
        "created_at": "2024-01-01T00:00:00+03:00",
        "updated_at": "2024-06-15T00:00:00+03:00"
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 10,
      "total": 8,
      "total_pages": 1
    }
  }
}
```

---

### 3.2 Create Salary Template

**Endpoint:** `POST /payroll/salary-structures`

**Request Body:**
```json
{
  "template_name": "G5 - Lead - Dar es Salaam",
  "grade": "G5",
  "location": "Dar es Salaam",
  "country": "Tanzania",
  "ctc": 4000000,
  "basic": 2000000,
  "hra": 600000,
  "transport": 120000,
  "medical": 60000,
  "other_allowances": 260000,
  "is_active": true
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `template_name` | string | Yes | Display name |
| `grade` | string | Yes | Grade code (e.g. G4, G5) |
| `location` | string | Yes | Office/location |
| `country` | string | Yes | Country (e.g. Tanzania) |
| `ctc` | number | Yes | Cost to company (TZS) |
| `basic` | number | Yes | Basic salary (TZS) |
| `hra` | number | No | House rent allowance (TZS) |
| `transport` | number | No | Transport allowance (TZS) |
| `medical` | number | No | Medical allowance (TZS) |
| `other_allowances` | number | No | Other allowances (TZS) |
| `is_active` | boolean | No | Active status (default: true) |

**Response:**
```json
{
  "success": true,
  "message": "Salary structure created successfully",
  "data": {
    "id": "str-003",
    "template_name": "G5 - Lead - Dar es Salaam",
    "grade": "G5",
    "location": "Dar es Salaam",
    "country": "Tanzania",
    "ctc": 4000000,
    "basic": 2000000,
    "gross_salary": 3040000,
    "total_deductions": 912000,
    "net_pay": 2128000,
    "is_active": true,
    "created_at": "2024-11-20T00:00:00+03:00"
  }
}
```

**Note:** Backend should auto-calculate `gross_salary`, `total_deductions`, `net_pay`, and statutory contributions based on Tanzania rules.

---

### 3.6 List Salary Components

**Endpoint:** `GET /payroll/salary-components`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | number | No | Page number |
| `page_size` | number | No | Items per page |
| `component_type` | string | No | `Earning`, `Deduction`, `Employer Contribution` |
| `is_active` | boolean | No | Filter by active status |
| `is_statutory` | boolean | No | Filter statutory components |
| `country` | string | No | Filter by country |

**Response:**
```json
{
  "success": true,
  "message": "Salary components retrieved successfully",
  "data": {
    "data": [
      {
        "id": "comp-001",
        "component_name": "Basic Salary",
        "component_code": "BASIC",
        "component_type": "Earning",
        "calculation_type": "Fixed",
        "default_amount": null,
        "default_percentage": null,
        "is_taxable": true,
        "is_statutory": false,
        "is_recurring": true,
        "applicable_for": ["All"],
        "display_in_payslip": true,
        "country": "All",
        "is_active": true
      },
      {
        "id": "comp-002",
        "component_name": "House Allowance (HRA)",
        "component_code": "HRA",
        "component_type": "Earning",
        "calculation_type": "Percentage of Basic",
        "default_amount": null,
        "default_percentage": 30,
        "is_taxable": true,
        "is_statutory": false,
        "is_recurring": true,
        "applicable_for": ["All"],
        "display_in_payslip": true,
        "country": "Tanzania",
        "is_active": true
      },
      {
        "id": "comp-005",
        "component_name": "PAYE Tax",
        "component_code": "PAYE",
        "component_type": "Deduction",
        "calculation_type": "Formula",
        "formula": "PAYE_SLAB_CALCULATION",
        "is_taxable": false,
        "is_statutory": true,
        "is_recurring": true,
        "applicable_for": ["All"],
        "display_in_payslip": true,
        "country": "Tanzania",
        "is_active": true
      },
      {
        "id": "comp-006",
        "component_name": "NSSF Employee",
        "component_code": "NSSF_EE",
        "component_type": "Deduction",
        "calculation_type": "Percentage of Gross",
        "default_percentage": 10,
        "is_taxable": false,
        "is_statutory": true,
        "is_recurring": true,
        "applicable_for": ["All"],
        "display_in_payslip": true,
        "country": "Tanzania",
        "is_active": true
      },
      {
        "id": "comp-010",
        "component_name": "NSSF Employer",
        "component_code": "NSSF_ER",
        "component_type": "Employer Contribution",
        "calculation_type": "Percentage of Gross",
        "default_percentage": 10,
        "is_taxable": false,
        "is_statutory": true,
        "is_recurring": true,
        "country": "Tanzania",
        "is_active": true
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 20,
      "total": 12,
      "total_pages": 1
    }
  }
}
```

**Calculation Types:**
- `Fixed` - Fixed amount
- `Percentage of Basic` - % of basic salary
- `Percentage of Gross` - % of gross salary
- `Formula` - Custom formula (e.g. PAYE slabs)

---

### 3.7 Create Salary Component

**Endpoint:** `POST /payroll/salary-components`

**Request Body:**
```json
{
  "component_name": "Performance Bonus",
  "component_code": "PERF_BONUS",
  "component_type": "Earning",
  "calculation_type": "Fixed",
  "default_amount": 0,
  "default_percentage": null,
  "is_taxable": true,
  "is_statutory": false,
  "is_recurring": false,
  "applicable_for": ["All"],
  "display_in_payslip": true,
  "country": "Tanzania"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `component_name` | string | Yes | Display name |
| `component_code` | string | Yes | Unique code (e.g. BASIC, HRA) |
| `component_type` | string | Yes | `Earning`, `Deduction`, `Employer Contribution` |
| `calculation_type` | string | Yes | `Fixed`, `Percentage of Basic`, `Percentage of Gross`, `Formula` |
| `default_amount` | number | No | Default fixed amount |
| `default_percentage` | number | No | Default percentage |
| `is_taxable` | boolean | Yes | Included in taxable income |
| `is_statutory` | boolean | Yes | Statutory component |
| `is_recurring` | boolean | Yes | Recurring every period |
| `applicable_for` | array | Yes | Grades/roles (e.g. ["All"]) |
| `display_in_payslip` | boolean | Yes | Show on payslip |
| `country` | string | Yes | Country code or "All" |

---

## 4. PAYSLIPS

### 4.1 List Payslips

**Endpoint:** `GET /payroll/payslips`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | number | No | Page number |
| `page_size` | number | No | Items per page |
| `pay_month` | number | No | Month (1-12) |
| `pay_year` | number | No | Year |
| `employee_id` | string | No | Filter by employee |
| `department_id` | string | No | Filter by department |
| `status` | string | No | `Draft`, `Released`, `Paid` |

**Response:**
```json
{
  "success": true,
  "message": "Payslips retrieved successfully",
  "data": {
    "data": [
      {
        "id": "slip-001",
        "emp_id": "EMP001",
        "employee_name": "John Doe",
        "employee_photo": "/img/avatars/thumb-1.jpg",
        "department": "Engineering",
        "designation": "Senior Software Engineer",
        "pay_period": "October 2024",
        "pay_month": 10,
        "pay_year": 2024,
        "gross_salary": 2300000,
        "total_deductions": 710000,
        "net_pay": 1590000,
        "status": "Paid",
        "days_worked": 30,
        "lop_days": 0
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 20,
      "total": 125,
      "total_pages": 7
    }
  }
}
```

---

### 4.2 Get Payslip Detail

**Endpoint:** `GET /payroll/payslips/:id`

**Response:**
```json
{
  "success": true,
  "message": "Payslip retrieved successfully",
  "data": {
    "id": "slip-001",
    "emp_id": "EMP001",
    "employee_name": "John Doe",
    "employee_photo": "/img/avatars/thumb-1.jpg",
    "department": "Engineering",
    "designation": "Senior Software Engineer",
    "date_of_joining": "2020-03-15",
    "pay_period": "October 2024",
    "pay_month": 10,
    "pay_year": 2024,
    "bank_name": "NMB Bank",
    "bank_account": "1234567890",
    "tin_number": "TIN-123456789",
    "nssf_number": "NSSF-987654321",
    
    "earnings": [
      { "code": "BASIC", "name": "Basic Salary", "amount": 1500000 },
      { "code": "HRA", "name": "House Allowance", "amount": 450000 },
      { "code": "TRANSPORT", "name": "Transport Allowance", "amount": 100000 },
      { "code": "MEDICAL", "name": "Medical Allowance", "amount": 50000 },
      { "code": "OT", "name": "Overtime Pay", "amount": 75000 },
      { "code": "OTHER", "name": "Other Allowances", "amount": 125000 }
    ],
    "earnings_total": 2300000,
    
    "deductions": [
      { "code": "PAYE", "name": "PAYE Tax", "amount": 380000 },
      { "code": "NSSF_EE", "name": "NSSF (Employee)", "amount": 230000 },
      { "code": "NHIF_EE", "name": "NHIF (Employee)", "amount": 50000 },
      { "code": "LOAN_EMI", "name": "Loan EMI", "amount": 50000 }
    ],
    "deductions_total": 710000,
    
    "employer_contributions": [
      { "code": "NSSF_ER", "name": "NSSF (Employer)", "amount": 230000 },
      { "code": "NHIF_ER", "name": "NHIF (Employer)", "amount": 50000 },
      { "code": "SDL", "name": "SDL", "amount": 115000 },
      { "code": "WCF", "name": "WCF", "amount": 23000 }
    ],
    "employer_contributions_total": 418000,
    
    "gross_salary": 2300000,
    "total_deductions": 710000,
    "net_pay": 1590000,
    
    "ytd": {
      "gross": 23000000,
      "tax": 3800000,
      "nssf": 2300000,
      "nhif": 500000,
      "net": 15900000
    },
    
    "attendance": {
      "working_days": 30,
      "days_worked": 30,
      "lop_days": 0,
      "overtime_hours": 8
    },
    
    "status": "Paid",
    "released_at": "2024-11-01T00:00:00+03:00",
    "paid_at": "2024-11-05T00:00:00+03:00"
  }
}
```

---

### 4.3 Download Payslip PDF

**Endpoint:** `GET /payroll/payslips/:id/download`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `format` | string | pdf | `pdf` |

**Response:** Binary PDF file

**Headers:**
```
Content-Type: application/pdf
Content-Disposition: attachment; filename="payslip-EMP001-2024-10.pdf"
```

---

### 4.4 Email Payslip

**Endpoint:** `POST /payroll/payslips/:id/email`

**Request Body:**
```json
{
  "email": "john.doe@company.com"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | No | Target email (defaults to employee's email) |

**Response:**
```json
{
  "success": true,
  "message": "Payslip sent to john.doe@company.com successfully"
}
```

---

### 4.5 Get Payslip Summary

**Endpoint:** `GET /payroll/payslips/summary`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `pay_month` | number | Yes | Month (1-12) |
| `pay_year` | number | Yes | Year |
| `employee_id` | string | No | For specific employee |

**Response:**
```json
{
  "success": true,
  "message": "Summary retrieved successfully",
  "data": {
    "pay_period": "October 2024",
    "net_pay": 1590000,
    "gross_salary": 2300000,
    "total_deductions": 710000,
    "ytd_gross": 23000000,
    "ytd_tax": 3800000,
    "ytd_net": 15900000
  }
}
```

---

### 4.6 Get My Payslip (Self-Service)

**Endpoint:** `GET /payroll/payslips/my`

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `pay_month` | number | No | Current | Month (1-12) |
| `pay_year` | number | No | Current | Year |

**Description:** Returns current user's payslip (for employee self-service portal).

**Response:** Same as Get Payslip Detail (4.2)

---

## 5. LOANS & ADVANCES

### 5.1 List Loans & Advances

**Endpoint:** `GET /payroll/loans`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | number | No | Page number |
| `page_size` | number | No | Items per page |
| `employee_id` | string | No | Filter by employee |
| `status` | string | No | `Pending Approval`, `Active`, `Closed`, `Defaulted`, `Rejected` |
| `loan_type` | string | No | `Education Loan`, `Emergency Loan`, `Personal Loan`, `Salary Advance` |

**Response:**
```json
{
  "success": true,
  "message": "Loans retrieved successfully",
  "data": {
    "data": [
      {
        "id": "loan-001",
        "emp_id": "EMP001",
        "employee_name": "John Doe",
        "employee_photo": "/img/avatars/thumb-1.jpg",
        "department": "Engineering",
        "loan_type": "Education Loan",
        "amount": 2000000,
        "interest_rate": 8,
        "tenure": 24,
        "emi_amount": 92000,
        "disbursed_date": "2024-01-15",
        "start_date": "2024-02-01",
        "end_date": "2026-01-01",
        "status": "Active",
        "total_paid": 920000,
        "outstanding_balance": 1080000,
        "next_emi_date": "2024-11-05",
        "payments_made": 10,
        "payments_remaining": 14,
        "purpose": "Child's university education",
        "created_at": "2024-01-10T00:00:00+03:00"
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 10,
      "total": 15,
      "total_pages": 2
    }
  }
}
```

---

### 5.2 Get Loan Detail

**Endpoint:** `GET /payroll/loans/:id`

**Response:**
```json
{
  "success": true,
  "message": "Loan retrieved successfully",
  "data": {
    "id": "loan-001",
    "emp_id": "EMP001",
    "employee_name": "John Doe",
    "employee_photo": "/img/avatars/thumb-1.jpg",
    "department": "Engineering",
    "designation": "Senior Software Engineer",
    "loan_type": "Education Loan",
    "amount": 2000000,
    "interest_rate": 8,
    "tenure": 24,
    "emi_amount": 92000,
    "disbursed_date": "2024-01-15",
    "start_date": "2024-02-01",
    "end_date": "2026-01-01",
    "status": "Active",
    "total_paid": 920000,
    "outstanding_balance": 1080000,
    "next_emi_date": "2024-11-05",
    "payments_made": 10,
    "payments_remaining": 14,
    "purpose": "Child's university education",
    "approved_by": "HR Manager",
    "approved_at": "2024-01-12T00:00:00+03:00",
    "repayment_schedule": [
      {
        "installment": 1,
        "due_date": "2024-02-05",
        "emi_amount": 92000,
        "principal": 75000,
        "interest": 17000,
        "status": "Paid",
        "paid_date": "2024-02-05"
      },
      {
        "installment": 2,
        "due_date": "2024-03-05",
        "emi_amount": 92000,
        "principal": 76000,
        "interest": 16000,
        "status": "Paid",
        "paid_date": "2024-03-05"
      }
    ],
    "created_at": "2024-01-10T00:00:00+03:00"
  }
}
```

---

### 5.3 Request Loan/Advance

**Endpoint:** `POST /payroll/loans`

**Request Body:**
```json
{
  "employee_id": "EMP001",
  "loan_type": "Salary Advance",
  "amount": 500000,
  "tenure": 5,
  "purpose": "Medical emergency"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `employee_id` | string | No | Employee ID (omit for self-service) |
| `loan_type` | string | Yes | `Education Loan`, `Emergency Loan`, `Personal Loan`, `Salary Advance` |
| `amount` | number | Yes | Requested amount (TZS) |
| `tenure` | number | Yes | Repayment tenure in months |
| `purpose` | string | No | Reason/purpose |

**Response:**
```json
{
  "success": true,
  "message": "Loan request submitted successfully",
  "data": {
    "id": "loan-003",
    "emp_id": "EMP001",
    "employee_name": "John Doe",
    "loan_type": "Salary Advance",
    "amount": 500000,
    "tenure": 5,
    "emi_amount": 100000,
    "interest_rate": 0,
    "status": "Pending Approval",
    "purpose": "Medical emergency",
    "created_at": "2024-11-20T00:00:00+03:00"
  }
}
```

---

### 5.4 Approve Loan

**Endpoint:** `POST /payroll/loans/:id/approve`

**Request Body:**
```json
{
  "emi_amount": 100000,
  "interest_rate": 0,
  "disbursement_date": "2024-11-25",
  "start_date": "2024-12-05",
  "remarks": "Approved as per policy"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `emi_amount` | number | Yes | Monthly EMI amount |
| `interest_rate` | number | Yes | Annual interest rate (%) |
| `disbursement_date` | string | Yes | Date to disburse loan |
| `start_date` | string | Yes | First EMI deduction date |
| `remarks` | string | No | Approval remarks |

**Response:**
```json
{
  "success": true,
  "message": "Loan approved successfully",
  "data": {
    "id": "loan-003",
    "status": "Active",
    "disbursed_date": "2024-11-25",
    "start_date": "2024-12-05",
    "emi_amount": 100000,
    "next_emi_date": "2024-12-05",
    "approved_by": "HR Manager",
    "approved_at": "2024-11-20T15:30:00+03:00"
  }
}
```

---

### 5.5 Reject Loan

**Endpoint:** `POST /payroll/loans/:id/reject`

**Request Body:**
```json
{
  "reason": "Insufficient tenure in company",
  "remarks": "Policy requires minimum 1 year of service"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `reason` | string | Yes | Rejection reason |
| `remarks` | string | No | Additional remarks |

**Response:**
```json
{
  "success": true,
  "message": "Loan request rejected",
  "data": {
    "id": "loan-003",
    "status": "Rejected",
    "rejection_reason": "Insufficient tenure in company",
    "rejected_by": "HR Manager",
    "rejected_at": "2024-11-20T15:30:00+03:00"
  }
}
```

---

### 5.6 Get Loans Summary

**Endpoint:** `GET /payroll/loans/summary`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `employee_id` | string | No | For specific employee |

**Response:**
```json
{
  "success": true,
  "message": "Summary retrieved successfully",
  "data": {
    "active_loans_count": 15,
    "total_disbursed": 25000000,
    "total_outstanding": 13800000,
    "this_month_emi_total": 1920000,
    "pending_approvals": 3,
    "by_type": {
      "Education Loan": { "count": 5, "amount": 10000000 },
      "Emergency Loan": { "count": 3, "amount": 3000000 },
      "Personal Loan": { "count": 4, "amount": 8000000 },
      "Salary Advance": { "count": 3, "amount": 4000000 }
    }
  }
}
```

---

## 6. COMPLIANCE & TAX

### 6.1 Get Compliance Summary

**Endpoint:** `GET /payroll/compliance/summary`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `pay_month` | number | Yes | Month (1-12) |
| `pay_year` | number | Yes | Year |

**Response:**
```json
{
  "success": true,
  "message": "Compliance summary retrieved successfully",
  "data": {
    "pay_month": 10,
    "pay_year": 2024,
    "paye": {
      "amount": 42000000,
      "employees_count": 125,
      "due_date": "2024-11-07",
      "status": "Pending"
    },
    "nssf": {
      "employee_amount": 23000000,
      "employer_amount": 23000000,
      "total_amount": 46000000,
      "employees_count": 125,
      "due_date": "2024-11-07",
      "status": "Pending"
    },
    "nhif": {
      "employee_amount": 5000000,
      "employer_amount": 5000000,
      "total_amount": 10000000,
      "employees_count": 125,
      "due_date": "2024-11-07",
      "status": "Pending"
    },
    "sdl": {
      "amount": 11500000,
      "due_date": "2024-11-07",
      "status": "Pending"
    },
    "wcf": {
      "amount": 2300000,
      "due_date": "2024-11-07",
      "status": "Pending"
    },
    "total_statutory": 111800000,
    "next_due_date": "2024-11-07",
    "days_until_due": 7
  }
}
```

---

### 6.2 List PAYE Tax Slabs

**Endpoint:** `GET /payroll/compliance/tax-slabs`

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `country` | string | No | Tanzania | Country code |
| `tax_year` | number | No | Current | Tax year |

**Response:**
```json
{
  "success": true,
  "message": "Tax slabs retrieved successfully",
  "data": {
    "country": "Tanzania",
    "tax_year": 2024,
    "currency": "TZS",
    "period": "Monthly",
    "slabs": [
      {
        "id": "tax-001",
        "slab_from": 0,
        "slab_to": 270000,
        "tax_rate": 0,
        "fixed_amount": 0,
        "description": "Tax-free threshold"
      },
      {
        "id": "tax-002",
        "slab_from": 270001,
        "slab_to": 520000,
        "tax_rate": 9,
        "fixed_amount": 0,
        "description": "9% on amount above 270,000"
      },
      {
        "id": "tax-003",
        "slab_from": 520001,
        "slab_to": 760000,
        "tax_rate": 20,
        "fixed_amount": 22500,
        "description": "TZS 22,500 + 20% on amount above 520,000"
      },
      {
        "id": "tax-004",
        "slab_from": 760001,
        "slab_to": 1000000,
        "tax_rate": 25,
        "fixed_amount": 70500,
        "description": "TZS 70,500 + 25% on amount above 760,000"
      },
      {
        "id": "tax-005",
        "slab_from": 1000001,
        "slab_to": 999999999,
        "tax_rate": 30,
        "fixed_amount": 130500,
        "description": "TZS 130,500 + 30% on amount above 1,000,000"
      }
    ]
  }
}
```

---

### 6.3 Get Statutory Contribution Rules

**Endpoint:** `GET /payroll/compliance/statutory`

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `country` | string | No | Tanzania | Country code |

**Response:**
```json
{
  "success": true,
  "message": "Statutory rules retrieved successfully",
  "data": {
    "country": "Tanzania",
    "currency": "TZS",
    "rules": {
      "nssf": {
        "name": "National Social Security Fund",
        "code": "NSSF",
        "employee_rate": 10,
        "employer_rate": 10,
        "basis": "Gross salary",
        "cap": null,
        "description": "Employee: 10% of gross salary. Employer: 10% of gross salary.",
        "authority": "NSSF Tanzania",
        "due_day": 7
      },
      "nhif": {
        "name": "National Health Insurance Fund",
        "code": "NHIF",
        "employee_rate": null,
        "employer_rate": null,
        "basis": "Schedule",
        "description": "Contribution as per NHIF schedule based on salary bands.",
        "schedule": [
          { "from": 0, "to": 150000, "employee": 2500, "employer": 2500 },
          { "from": 150001, "to": 250000, "employee": 5000, "employer": 5000 },
          { "from": 250001, "to": 400000, "employee": 10000, "employer": 10000 },
          { "from": 400001, "to": 600000, "employee": 15000, "employer": 15000 },
          { "from": 600001, "to": 999999999, "employee": 20000, "employer": 20000 }
        ],
        "authority": "NHIF Tanzania",
        "due_day": 7
      },
      "sdl": {
        "name": "Skills Development Levy",
        "code": "SDL",
        "employee_rate": 0,
        "employer_rate": 5,
        "basis": "Total payroll",
        "description": "Employer only: 5% of total gross payroll.",
        "authority": "TRA",
        "due_day": 7
      },
      "wcf": {
        "name": "Workers Compensation Fund",
        "code": "WCF",
        "employee_rate": 0,
        "employer_rate": 1,
        "basis": "Total payroll",
        "description": "Employer only: 1% of total gross payroll.",
        "authority": "WCF Tanzania",
        "due_day": 7
      }
    }
  }
}
```

---

### 6.4 Download Returns File

**Endpoint:** `GET /payroll/compliance/returns/:type`

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `type` | string | Yes | Return type (see below) |

**Return Types:**
- `paye_schedule` - PAYE Schedule for TRA
- `nssf_return` - NSSF Monthly Return
- `nhif_return` - NHIF Monthly Return
- `sdl_return` - SDL Return
- `wcf_return` - WCF Return
- `payment_challans` - All payment challans

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `pay_month` | number | Yes | Month (1-12) |
| `pay_year` | number | Yes | Year |
| `format` | string | No | `xlsx`, `csv`, `pdf` (default: xlsx) |

**Response:** Binary file download

**Headers:**
```
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="paye-schedule-2024-10.xlsx"
```

---

## 7. PAYROLL REPORTS

### 7.1 Get Report Filter Options

**Endpoint:** `GET /payroll/reports/filters`

**Response:**
```json
{
  "success": true,
  "message": "Filter options retrieved",
  "data": {
    "report_types": [
      { "value": "register", "label": "Payroll Register" },
      { "value": "variance", "label": "Variance Analysis" },
      { "value": "cost_by_department", "label": "Cost by Department" },
      { "value": "tax_liability", "label": "Tax Liability Report" },
      { "value": "employer_contributions", "label": "Employer Contributions" },
      { "value": "bank_payment_file", "label": "Bank Payment File" }
    ],
    "departments": [
      { "value": "", "label": "All Departments" },
      { "value": "dept-1", "label": "Engineering" },
      { "value": "dept-2", "label": "Sales" },
      { "value": "dept-3", "label": "HR" },
      { "value": "dept-4", "label": "Finance" },
      { "value": "dept-5", "label": "Marketing" }
    ],
    "date_range": {
      "min_date": "2023-01-01",
      "max_date": "2024-12-31"
    },
    "formats": [
      { "value": "xlsx", "label": "Excel (.xlsx)" },
      { "value": "csv", "label": "CSV (.csv)" },
      { "value": "pdf", "label": "PDF (.pdf)" }
    ]
  }
}
```

---

### 7.2 Get Report KPIs

**Endpoint:** `GET /payroll/reports/kpis`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `start_month` | string | Yes | Start month (YYYY-MM) |
| `end_month` | string | Yes | End month (YYYY-MM) |
| `department_id` | string | No | Filter by department |

**Response:**
```json
{
  "success": true,
  "message": "KPIs retrieved successfully",
  "data": {
    "period": "October 2024 - November 2024",
    "total_payout": 375000000,
    "payout_trend": 3.2,
    "headcount": 125,
    "headcount_trend": 2,
    "avg_salary": 1500000,
    "avg_salary_trend": 5.1,
    "employer_contribution": 94000000,
    "employer_contribution_trend": 3.0,
    "cost_per_employee": 3760000,
    "variance_from_budget": -2.5
  }
}
```

---

### 7.3 Get Department Cost Analysis

**Endpoint:** `GET /payroll/reports/department-cost`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `start_month` | string | Yes | Start month (YYYY-MM) |
| `end_month` | string | Yes | End month (YYYY-MM) |
| `department_id` | string | No | Filter by department |

**Response:**
```json
{
  "success": true,
  "message": "Department cost retrieved successfully",
  "data": {
    "period": "October 2024 - November 2024",
    "total_cost": 375000000,
    "departments": [
      {
        "department_id": "dept-1",
        "department_name": "Engineering",
        "employees": 50,
        "gross_salary": 150000000,
        "deductions": 31500000,
        "employer_contributions": 37500000,
        "total_cost": 156000000,
        "percent": 41.6,
        "cost_per_employee": 3120000
      },
      {
        "department_id": "dept-2",
        "department_name": "Sales",
        "employees": 30,
        "gross_salary": 90000000,
        "deductions": 18900000,
        "employer_contributions": 22500000,
        "total_cost": 93600000,
        "percent": 25.0,
        "cost_per_employee": 3120000
      },
      {
        "department_id": "dept-3",
        "department_name": "HR",
        "employees": 15,
        "gross_salary": 45000000,
        "deductions": 9450000,
        "employer_contributions": 11250000,
        "total_cost": 46800000,
        "percent": 12.5,
        "cost_per_employee": 3120000
      }
    ]
  }
}
```

---

### 7.4 Export Report

**Endpoint:** `GET /payroll/reports/export`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `report_type` | string | Yes | See report types in 7.1 |
| `start_month` | string | Yes | Start month (YYYY-MM) |
| `end_month` | string | Yes | End month (YYYY-MM) |
| `department_id` | string | No | Filter by department |
| `format` | string | No | `xlsx`, `csv`, `pdf` (default: xlsx) |

**Response:** Binary file download

**Headers:**
```
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="payroll-register-2024-10.xlsx"
```

---

## 8. PAYMENT SYSTEM INTEGRATION

### 8.1 Get Disbursement Batch

**Endpoint:** `GET /payroll/runs/:run_id/disbursement`

**Description:** Returns payment instructions for a finalized payroll run. Used by payment system to execute salary transfers.

**Response:**
```json
{
  "success": true,
  "message": "Disbursement batch retrieved successfully",
  "data": {
    "run_id": "run-002",
    "run_name": "November 2024 Payroll",
    "pay_period": "01 Nov - 30 Nov 2024",
    "disbursement_date": "2024-12-05",
    "status": "Finalized",
    "currency": "TZS",
    "total_amount": 148125000,
    "total_employees": 125,
    "payments": [
      {
        "employee_id": "EMP001",
        "employee_name": "John Doe",
        "payslip_id": "slip-001",
        "net_pay": 1590000,
        "bank_name": "NMB Bank",
        "bank_code": "NMB",
        "branch_code": "001",
        "account_number": "1234567890",
        "account_type": "Current",
        "beneficiary_name": "John Doe",
        "reference": "PAY-NOV2024-EMP001"
      }
    ]
  }
}
```

---

### 8.2 Download Bank Payment File

**Endpoint:** `GET /payroll/runs/:run_id/disbursement/file`

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `format` | string | No | csv | `csv`, `xlsx`, `iso20022`, bank-specific |

**Response:** Binary file download

**CSV Structure:**
```csv
employee_id,employee_name,net_pay,bank_code,branch_code,account_number,beneficiary_name,reference
EMP001,John Doe,1590000,NMB,001,1234567890,John Doe,PAY-NOV2024-EMP001
EMP002,Jane Smith,1450000,CRDB,002,9876543210,Jane Smith,PAY-NOV2024-EMP002
```

---

### 8.3 Payment Status Callback

**Endpoint:** `POST /payroll/disbursement/callback`

**Description:** Webhook for payment system to report payment status back to payroll.

**Request Body:**
```json
{
  "run_id": "run-002",
  "batch_reference": "BANK-REF-20241205-001",
  "status": "completed",
  "processed_at": "2024-12-05T14:30:00+03:00",
  "results": [
    {
      "employee_id": "EMP001",
      "payslip_id": "slip-001",
      "status": "success",
      "transaction_reference": "TXN-123456"
    },
    {
      "employee_id": "EMP002",
      "payslip_id": "slip-002",
      "status": "failed",
      "error_code": "INVALID_ACCOUNT",
      "error_message": "Account closed"
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Disbursement status updated",
  "data": {
    "run_id": "run-002",
    "run_status": "Disbursed",
    "processed_count": 124,
    "failed_count": 1
  }
}
```

---

## Appendix A: Tanzania PAYE Calculation

### Monthly PAYE Tax Slabs (FY 2024)

| Income Range (TZS/month) | Rate | Fixed Amount | Formula |
|--------------------------|------|--------------|---------|
| 0 - 270,000 | 0% | 0 | 0 |
| 270,001 - 520,000 | 9% | 0 | (Income - 270,000) × 9% |
| 520,001 - 760,000 | 20% | 22,500 | 22,500 + (Income - 520,000) × 20% |
| 760,001 - 1,000,000 | 25% | 70,500 | 70,500 + (Income - 760,000) × 25% |
| Above 1,000,000 | 30% | 130,500 | 130,500 + (Income - 1,000,000) × 30% |

### Example Calculation

**Gross Salary:** TZS 2,300,000

```
PAYE = 130,500 + (2,300,000 - 1,000,000) × 30%
     = 130,500 + 1,300,000 × 0.30
     = 130,500 + 390,000
     = TZS 520,500
```

---

## Appendix B: Statutory Calculations

### NSSF Calculation
```
Employee Contribution = Gross Salary × 10%
Employer Contribution = Gross Salary × 10%
```

### NHIF Schedule (2024)
| Salary Range (TZS) | Employee | Employer |
|--------------------|----------|----------|
| 0 - 150,000 | 2,500 | 2,500 |
| 150,001 - 250,000 | 5,000 | 5,000 |
| 250,001 - 400,000 | 10,000 | 10,000 |
| 400,001 - 600,000 | 15,000 | 15,000 |
| Above 600,000 | 20,000 | 20,000 |

### SDL & WCF (Employer Only)
```
SDL = Total Gross Payroll × 5%
WCF = Total Gross Payroll × 1%
```

---

## Appendix C: Data Models

### PayrollRun
```typescript
{
  id: string;
  run_name: string;
  pay_period: string;
  pay_month: number;
  pay_year: number;
  pay_frequency: "Monthly" | "Bi-Weekly" | "Weekly";
  status: "Draft" | "In Review" | "Approved" | "Finalized" | "Disbursed" | "Closed";
  current_step: number;
  total_employees: number;
  total_gross: number;
  total_deductions: number;
  total_net: number;
  total_employer_contributions: number;
  cutoff_date: string;
  disbursement_date: string;
  created_by: string;
  approved_by: string | null;
  finalized_by: string | null;
  created_at: string;
  updated_at: string;
}
```

### SalaryStructure
```typescript
{
  id: string;
  template_name: string;
  grade: string;
  location: string;
  country: string;
  ctc: number;
  basic: number;
  hra: number;
  transport: number;
  medical: number;
  other_allowances: number;
  gross_salary: number;
  total_deductions: number;
  net_pay: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
```

### Payslip
```typescript
{
  id: string;
  emp_id: string;
  employee_name: string;
  department: string;
  designation: string;
  pay_period: string;
  pay_month: number;
  pay_year: number;
  earnings: { code: string; name: string; amount: number }[];
  deductions: { code: string; name: string; amount: number }[];
  employer_contributions: { code: string; name: string; amount: number }[];
  gross_salary: number;
  total_deductions: number;
  net_pay: number;
  ytd: { gross: number; tax: number; net: number };
  attendance: { working_days: number; days_worked: number; lop_days: number };
  status: "Draft" | "Released" | "Paid";
}
```

### Loan
```typescript
{
  id: string;
  emp_id: string;
  employee_name: string;
  loan_type: "Education Loan" | "Emergency Loan" | "Personal Loan" | "Salary Advance";
  amount: number;
  interest_rate: number;
  tenure: number;
  emi_amount: number;
  disbursed_date: string;
  status: "Pending Approval" | "Active" | "Closed" | "Defaulted" | "Rejected";
  total_paid: number;
  outstanding_balance: number;
  next_emi_date: string;
  purpose: string;
  created_at: string;
}
```

---

*End of Payroll Backend API Specification*

**Document Version:** 2.0  
**Created:** January 23, 2026  
**Maintainer:** Backend Development Team

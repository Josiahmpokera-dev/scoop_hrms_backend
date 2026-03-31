# Payroll Management API — Comprehensive Documentation

Complete API documentation covering **Payroll Dashboard**, **Payroll Runs**, **Salary Structures**, **Payslips**, **Loans & Advances**, **Compliance**, and **Employee Self-Service** payroll features.



## Payroll Workflow Overview

The payroll system implements a comprehensive 7-phase workflow to ensure accuracy, compliance, and proper approval processes:

### 🔷 Phase 1: Employee Master Data & Inputs Collection
- **HR Responsibility**: Maintain employee records, salary changes, terminations, bank details, tax setup
- **System Validation**: New employee setup, salary structure assignment, bank account verification
- **Inputs**: Attendance (timesheets), leave records, overtime approvals, bonuses/commissions

### 🔷 Phase 2: Payroll Run Creation & Pre-Validation
- **Pre-check Validation**: Comprehensive validation of all inputs and master data
- **System Checks**: Approved timesheets, pending overtime, salary structures, compliance
- **Output**: Draft payroll register with identified issues

### 🔷 Phase 3: HR Review & Validation
- **HR Review**: Validate attendance corrections, missing employees, salary changes, exceptions
- **Approval Process**: HR can approve, reject, or request corrections
- **Validation Focus**: Employee-level accuracy, policy compliance, exception handling

### 🔷 Phase 4: Finance Review & Approval
- **Finance Validation**: Budget availability, cost center allocation, tax correctness, statutory compliance
- **Budget Verification**: Ensure sufficient funds and proper accounting allocation
- **Approval Process**: Finance can approve, reject, or request corrections

### 🔷 Phase 5: Management Final Approval
- **Management Sign-off**: Final approval by HR Manager, Finance Manager, or CEO (based on amount)
- **Authorization**: Multi-level approval based on payroll amount and company policy
- **Final Lock**: Payroll becomes locked and ready for processing

### 🔷 Phase 6: Payment Execution
- **Bank File Generation**: Create payment files (CSV/ISO20022 format)
- **Payment Processing**: Upload to bank systems and execute payments
- **Status Tracking**: Monitor payment status and handle failures

### 🔷 Phase 7: Post-Payroll Processing
- **Accounting Entries**: Automatic GL posting for salary expenses, tax liabilities, pension contributions
- **Payslip Distribution**: Generate and distribute payslips to employees
- **Compliance Reporting**: Generate statutory reports for tax authorities, pension funds
- **Audit Trail**: Complete audit trail of all approvals and changes

---

## Workflow Status Transitions

```
Draft → Pending HR Review → HR Reviewed → Pending Finance Review → Finance Reviewed → Pending Management Approval → Approved → Finalized → Disbursed → Closed
```

**Status Descriptions:**
- `Draft`: Initial creation, editable
- `Pending HR Review`: Submitted for HR validation
- `HR Reviewed`: HR approved, ready for Finance
- `Pending Finance Review`: Submitted for Finance validation
- `Finance Reviewed`: Finance approved, ready for Management
- `Pending Management Approval`: Awaiting final management sign-off
- `Approved`: Management approved, locked for payment
- `Finalized`: Ready for bank file generation
- `Disbursed`: Payments processed
- `Closed`: Complete, archived

---

---

## Table of Contents

### Part A — Payroll Dashboard & Overview
1. [Get Payroll Dashboard](#1-get-payroll-dashboard)

### Part B — Payroll Runs Management
2. [List Payroll Runs](#2-list-payroll-runs)
3. [Create Payroll Run](#3-create-payroll-run)
4. [Get Current Payroll Run](#4-get-current-payroll-run)
5. [Get Payroll Run Details](#5-get-payroll-run-details)
6. [Update Payroll Run](#6-update-payroll-run)
7. [Delete Payroll Run](#7-delete-payroll-run)
8. [Get Payroll Run Summary](#8-get-payroll-run-summary)
9. [Get Payroll Employees](#9-get-payroll-employees)
10. [Run Pre-Check](#10-run-pre-check)
11. [Calculate Payroll](#11-calculate-payroll)
12. [Submit for HR Review](#12-submit-for-hr-review)
13. [HR Review](#13-hr-review)
14. [Finance Review](#14-finance-review)
15. [Management Approval](#15-management-approval)
16. [Advance Payroll Step](#16-advance-payroll-step)
17. [Revert Payroll Step](#17-revert-payroll-step)
18. [Generate Payslips](#18-generate-payslips)

### Part C — Salary Structures & Components
15. [List Salary Structures](#15-list-salary-structures)
16. [Create Salary Structure](#16-create-salary-structure)
17. [Get Salary Structure](#17-get-salary-structure)
18. [Update Salary Structure](#18-update-salary-structure)
19. [Delete Salary Structure](#19-delete-salary-structure)
20. [Simulate Salary](#20-simulate-salary)
21. [List Salary Components](#21-list-salary-components)
22. [Create Salary Component](#22-create-salary-component)
23. [Get Salary Component](#23-get-salary-component)
24. [Update Salary Component](#24-update-salary-component)
25. [Delete Salary Component](#25-delete-salary-component)

### Part D — Payslips Management
26. [List Payslips](#26-list-payslips)
27. [Get Payslip Summary](#27-get-payslip-summary)
28. [Get Payslip Details](#28-get-payslip-details)
29. [Download Payslip](#29-download-payslip)
30. [Send Payslip Email](#30-send-payslip-email)
31. [Bulk Download Payslips](#31-bulk-download-payslips)
32. [Release Payslips](#32-release-payslips)
33. [Bulk Email Payslips](#33-bulk-email-payslips)

### Part E — Loans & Advances Management
34. [List Loans](#34-list-loans)
35. [Create Loan](#35-create-loan)
36. [Get Loan Summary](#36-get-loan-summary)
37. [Calculate EMI](#37-calculate-emi)
38. [Get Loan Details](#38-get-loan-details)
39. [Get Repayment Schedule](#39-get-repayment-schedule)
40. [Approve Loan](#40-approve-loan)
41. [Reject Loan](#41-reject-loan)
42. [Record Repayment](#42-record-repayment)

### Part F — Compliance & Tax Management
43. [Get Compliance Summary](#43-get-compliance-summary)
44. [Get All Compliance Rules](#44-get-all-compliance-rules)
45. [Get Tax Slabs](#45-get-tax-slabs)
46. [Get Statutory Rules](#46-get-statutory-rules)
47. [Get NHIF Schedule](#47-get-nhif-schedule)
48. [Get Compliance Payments](#48-get-compliance-payments)
49. [Simulate Tax](#49-simulate-tax)
50. [Calculate Monthly Statutory](#50-calculate-monthly-statutory)
51. [Generate Compliance Return](#51-generate-compliance-return)
52. [Seed Compliance Data](#52-seed-compliance-data)

### Part G — Payroll Reports
53. [Get Payroll Report](#53-get-payroll-report)
54. [Get Payroll KPIs](#54-get-payroll-kpis)
55. [Get Department Cost Report](#55-get-department-cost-report)
56. [Export Payroll Report](#56-export-payroll-report)
57. [Get Bank File Report](#57-get-bank-file-report)

### Part H — Employee Self-Service Payroll
58. [Get My Payslips](#58-get-my-payslips)
59. [Get My Latest Payslip](#59-get-my-latest-payslip)
60. [Get My Salary Slip Summary](#60-get-my-salary-slip-summary)
61. [Download My Payslip](#61-download-my-payslip)
62. [Get My Loans](#62-get-my-loans)
63. [Apply For Loan](#63-apply-for-loan)

### Appendix
- [Authentication & Authorization](#authentication--authorization)
- [Encryption-Transparent Behavior](#encryption-transparent-behavior)
- [Data Models](#data-models)
- [Quick Reference Table](#quick-reference-table)

---

## Authentication & Authorization

| Scope              | Middleware                          | Who Can Access              |
|--------------------|-------------------------------------|-----------------------------|
| HR/Admin           | `AuthMiddleware` + `HRMiddleware`   | Admin, Super Admin, HR      |
| Employee (Self)    | `AuthMiddleware`                    | Any authenticated user      |

All endpoints require `Authorization: Bearer <token>`.

---

# Part A — Payroll Dashboard & Overview

Base URL: `/api/v1/payroll/dashboard`

These endpoints provide payroll overview and dashboard metrics for HR/Admin users.

---

## 1. Get Payroll Dashboard

Retrieve payroll dashboard metrics and KPIs.

```
GET /api/v1/payroll/dashboard
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "total_payroll_runs": 12,
    "active_payroll_run": {
      "id": 1,
      "name": "January 2025 Payroll",
      "status": "completed"
    },
    "total_employees_paid": 145,
    "total_payroll_amount": 125000000,
    "pending_approvals": 3,
    "recent_payslips_generated": 145,
    "monthly_trend": [
      {"month": "Oct 2024", "amount": 115000000},
      {"month": "Nov 2024", "amount": 118000000},
      {"month": "Dec 2024", "amount": 122000000},
      {"month": "Jan 2025", "amount": 125000000}
    ]
  }
}
```

---

# Part B — Payroll Runs Management

Base URL: `/api/v1/payroll/runs`

These endpoints manage payroll runs for HR/Admin users.

---

## 2. List Payroll Runs

List all payroll runs with pagination and filtering.

```
GET /api/v1/payroll/runs
```

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `pageSize` (optional) - Items per page (default: 20, max: 100)
- `year` (optional) - Filter by year
- `month` (optional) - Filter by month
- `status` (optional) - Filter by status

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "January 2025 Payroll",
      "pay_month": 1,
      "pay_year": 2025,
      "status": "completed",
      "total_employees": 145,
      "total_amount": 125000000,
      "created_at": "2025-01-15T10:00:00Z",
      "completed_at": "2025-01-20T15:30:00Z"
    }
  ],
  "meta": {
    "total": 12,
    "page": 1,
    "pageSize": 20
  }
}
```

---

## 3. Create Payroll Run

Create a new payroll run.

```
POST /api/v1/payroll/runs
```

**Auth:** HR/Admin only

**Request Body:**
```json
{
  "name": "February 2025 Payroll",
  "pay_month": 2,
  "pay_year": 2025,
  "description": "Regular monthly payroll run"
}
```

### Response — `201 Created`

```json
{
  "success": true,
  "message": "Payroll run created successfully",
  "data": {
    "id": 2,
    "name": "February 2025 Payroll",
    "pay_month": 2,
    "pay_year": 2025,
    "status": "draft",
    "created_at": "2025-02-01T09:00:00Z"
  }
}
```

---

## 4. Get Current Payroll Run

Get the currently active payroll run.

```
GET /api/v1/payroll/runs/current
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "id": 2,
    "name": "February 2025 Payroll",
    "pay_month": 2,
    "pay_year": 2025,
    "status": "draft",
    "total_employees": 0,
    "total_amount": 0,
    "created_at": "2025-02-01T09:00:00Z"
  }
}
```

---

## 5. Get Payroll Run Details

Get detailed information about a specific payroll run.

```
GET /api/v1/payroll/runs/:id
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "January 2025 Payroll",
    "pay_month": 1,
    "pay_year": 2025,
    "status": "completed",
    "total_employees": 145,
    "total_amount": 125000000,
    "created_at": "2025-01-15T10:00:00Z",
    "completed_at": "2025-01-20T15:30:00Z",
    "payslips_generated": 145,
    "payslips_released": 145,
    "employees": [
      {
        "employee_id": "EMP001",
        "name": "John Doe",
        "department": "Engineering",
        "gross_salary": 850000,
        "net_salary": 720000,
        "status": "paid"
      }
    ]
  }
}
```

---

## 6. Update Payroll Run

Update an existing payroll run.

```
PUT /api/v1/payroll/runs/:id
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** HR/Admin only

**Request Body:**
```json
{
  "name": "February 2025 Payroll - Updated",
  "description": "Updated monthly payroll run"
}
```

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payroll run updated successfully",
  "data": {
    "id": 2,
    "name": "February 2025 Payroll - Updated",
    "pay_month": 2,
    "pay_year": 2025,
    "status": "draft",
    "updated_at": "2025-02-01T10:30:00Z"
  }
}
```

---

## 7. Delete Payroll Run

Delete a payroll run (only allowed for draft runs).

```
DELETE /api/v1/payroll/runs/:id
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payroll run deleted successfully"
}
```

---

## 8. Get Payroll Run Summary

Get summary statistics for a payroll run.

```
GET /api/v1/payroll/runs/:id/summary
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "total_employees": 145,
    "total_gross_salary": 125000000,
    "total_deductions": 25000000,
    "total_net_salary": 100000000,
    "total_tax": 18000000,
    "total_nhif": 700000,
    "total_nssf": 6200000,
    "department_breakdown": [
      {
        "department": "Engineering",
        "employee_count": 45,
        "total_salary": 45000000
      },
      {
        "department": "Sales",
        "employee_count": 30,
        "total_salary": 30000000
      }
    ]
  }
}
```

---

## 9. Get Payroll Employees

Get employees included in a payroll run.

```
GET /api/v1/payroll/runs/:id/employees
```

**Path Parameters:**
- `id` - Payroll run ID

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `pageSize` (optional) - Items per page (default: 20, max: 100)
- `department` (optional) - Filter by department

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "employee_id": "EMP001",
      "name": "John Doe",
      "department": "Engineering",
      "position": "Senior Developer",
      "gross_salary": 850000,
      "net_salary": 720000,
      "status": "included"
    }
  ],
  "meta": {
    "total": 145,
    "page": 1,
    "pageSize": 20
  }
}
```

---

## 10. Run Pre-Check

Run pre-check validation for a payroll run.

```
POST /api/v1/payroll/runs/:id/pre-check
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Pre-check completed successfully",
  "data": {
    "total_employees": 145,
    "valid_employees": 142,
    "errors": 3,
    "validation_results": [
      {
        "employee_id": "EMP001",
        "name": "John Doe",
        "status": "valid",
        "issues": []
      },
      {
        "employee_id": "EMP002",
        "name": "Jane Smith",
        "status": "invalid",
        "issues": ["Missing bank account details"]
      }
    ]
  }
}
```

---

## 11. Calculate Payroll

Calculate payroll for a specific run.

**Note:** This process now includes an attendance check. Employees without approved timesheets for the pay period will be skipped. Approved overtime earnings are automatically included in the gross salary calculation.

```
POST /api/v1/payroll/runs/:id/calculate
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payroll calculated successfully",
  "data": {
    "total_employees": 145,
    "total_gross_salary": 125000000,
    "total_net_salary": 100000000,
    "calculation_time": "2.5s",
    "status": "calculated"
  }
}
```

---

## 12. Submit for HR Review

Submit payroll run for HR review. This triggers the formal approval workflow.

**Note:** Pre-check validation is performed before submission. All critical issues must be resolved.

```
POST /api/v1/payroll/runs/:id/submit-hr-review
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payroll run submitted for HR review",
  "data": {
    "id": 2,
    "status": "Pending HR Review",
    "currentStep": 1,
    "hrReviewDeadline": "2025-03-22T10:00:00Z"
  }
}
```

---

## 13. HR Review

Perform HR review of payroll run. HR validates attendance, salary changes, and exception cases.

```
POST /api/v1/payroll/runs/:id/hr-review
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** HR/Admin only

**Request Body:**
```json
{
  "status": "Approved",
  "comments": "All attendance records verified. Salary changes approved."
}
```

**Status Options:**
- `Approved` - Payroll approved by HR, moves to Finance review
- `Rejected` - Payroll rejected, returns to Draft status
- `Needs Correction` - Issues identified, returns to Draft with comments

### Response — `200 OK`

```json
{
  "success": true,
  "message": "HR review completed",
  "data": {
    "id": 2,
    "status": "HR Reviewed",
    "currentStep": 2,
    "hrReviewStatus": "Approved",
    "hrReviewedBy": "Jane Smith",
    "hrReviewedAt": "2025-03-20T14:30:00Z",
    "hrReviewComments": "All attendance records verified. Salary changes approved."
  }
}
```

---

## 14. Finance Review

Perform Finance review of payroll run. Finance validates budget availability, cost allocation, and statutory compliance.

```
POST /api/v1/payroll/runs/:id/finance-review
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** Finance/Admin only

**Request Body:**
```json
{
  "status": "Approved",
  "comments": "Budget verified. Tax calculations correct.",
  "budgetVerified": true
}
```

**Status Options:**
- `Approved` - Payroll approved by Finance, moves to Management approval
- `Rejected` - Payroll rejected, returns to HR Reviewed status
- `Needs Correction` - Issues identified, returns to HR Reviewed status

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Finance review completed",
  "data": {
    "id": 2,
    "status": "Finance Reviewed",
    "currentStep": 3,
    "financeReviewStatus": "Approved",
    "financeReviewedBy": "John Finance",
    "financeReviewedAt": "2025-03-21T09:15:00Z",
    "financeReviewComments": "Budget verified. Tax calculations correct.",
    "budgetVerified": true
  }
}
```

---

## 15. Management Approval

Final management approval for payroll run. Management provides final sign-off before payment processing.

```
POST /api/v1/payroll/runs/:id/management-approval
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** Management/Admin only

**Request Body:**
```json
{
  "approved": true,
  "comments": "Approved for payment processing."
}
```

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Management approval completed",
  "data": {
    "id": 2,
    "status": "Pending Management Approval",
    "currentStep": 4,
    "managementApproved": true,
    "managementApprovedBy": "CEO Name",
    "managementApprovedAt": "2025-03-21T16:45:00Z",
    "managementApprovalComments": "Approved for payment processing."
  }
}
```

---

## 16. Advance Payroll Step

Advance payroll run to the next step.

```
POST /api/v1/payroll/runs/:id/advance
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payroll advanced to next step",
  "data": {
    "current_step": "payslip_generation",
    "previous_step": "calculation",
    "status": "in_progress"
  }
}
```

---

## 13. Revert Payroll Step

Revert payroll run to the previous step.

```
POST /api/v1/payroll/runs/:id/revert
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payroll reverted to previous step",
  "data": {
    "current_step": "calculation",
    "previous_step": "payslip_generation",
    "status": "in_progress"
  }
}
```

---

## 14. Generate Payslips

Generate payslips for a payroll run.

```
POST /api/v1/payroll/runs/:id/generate-payslips
```

**Path Parameters:**
- `id` - Payroll run ID

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payslips generated successfully",
  "data": {
    "total_payslips": 145,
    "generated_payslips": 145,
    "generation_time": "45s",
    "status": "payslips_generated"
  }
}
```

---
    "id": 2,
    "name": "February 2025 Payroll",
    "pay_month": 2,
    "pay_year": 2025,
    "status": "draft",
    "created_at": "2025-02-01T09:00:00Z"
  }
}
```

---

## 4. Get Current Payroll Run

Get the current active payroll run.

```
GET /api/v1/payroll/runs/current
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "id": 2,
    "name": "February 2025 Payroll",
    "pay_month": 2,
    "pay_year": 2025,
    "status": "in_progress",
    "current_step": "calculation",
    "total_employees": 150,
    "processed_employees": 120,
    "created_at": "2025-02-01T09:00:00Z",
    "started_at": "2025-02-01T10:00:00Z"
  }
}
```

---

## 5. Get Payroll Run Details

Get detailed information about a specific payroll run.

```
GET /api/v1/payroll/runs/:id
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "January 2025 Payroll",
    "pay_month": 1,
    "pay_year": 2025,
    "status": "completed",
    "description": "Regular monthly payroll run",
    "total_employees": 145,
    "total_amount": 125000000,
    "total_deductions": 25000000,
    "total_net_pay": 100000000,
    "created_by": "admin@example.com",
    "created_at": "2025-01-15T10:00:00Z",
    "started_at": "2025-01-15T11:00:00Z",
    "completed_at": "2025-01-20T15:30:00Z",
    "employees": [
      {
        "employee_id": "EMP001",
        "name": "John Doe",
        "department": "Engineering",
        "gross_salary": 2500000,
        "deductions": 500000,
        "net_pay": 2000000,
        "status": "processed"
      }
    ]
  }
}
```

---

## 6. Update Payroll Run

Update a payroll run.

```
PUT /api/v1/payroll/runs/:id
```

**Auth:** HR/Admin only

**Request Body:**
```json
{
  "name": "February 2025 Payroll - Updated",
  "description": "Updated payroll run description"
}
```

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payroll run updated successfully",
  "data": {
    "id": 2,
    "name": "February 2025 Payroll - Updated",
    "pay_month": 2,
    "pay_year": 2025,
    "status": "draft",
    "description": "Updated payroll run description",
    "updated_at": "2025-02-01T11:00:00Z"
  }
}
```

---

## 7. Delete Payroll Run

Delete a payroll run (only allowed for draft status).

```
DELETE /api/v1/payroll/runs/:id
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payroll run deleted successfully"
}
```

---

## 8. Get Payroll Run Summary

Get summary statistics for a payroll run.

```
GET /api/v1/payroll/runs/:id/summary
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "total_employees": 145,
    "processed_employees": 145,
    "total_gross_salary": 125000000,
    "total_deductions": 25000000,
    "total_net_pay": 100000000,
    "average_salary": 862069,
    "department_breakdown": [
      {
        "department": "Engineering",
        "employee_count": 45,
        "total_amount": 45000000,
        "percentage": 45.0
      },
      {
        "department": "Sales",
        "employee_count": 30,
        "total_amount": 30000000,
        "percentage": 30.0
      }
    ],
    "deduction_breakdown": [
      {
        "type": "Tax",
        "amount": 15000000,
        "percentage": 60.0
      },
      {
        "type": "NHIF",
        "amount": 5000000,
        "percentage": 20.0
      },
      {
        "type": "NSSF",
        "amount": 5000000,
        "percentage": 20.0
      }
    ]
  }
}
```

---

## 9. Get Payroll Employees

Get employees included in a payroll run.

```
GET /api/v1/payroll/runs/:id/employees
```

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `pageSize` (optional) - Items per page (default: 20, max: 100)
- `department` (optional) - Filter by department
- `status` (optional) - Filter by processing status

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "employee_id": "EMP001",
      "name": "John Doe",
      "department": "Engineering",
      "job_position": "Senior Developer",
      "gross_salary": 2500000,
      "basic_salary": 1500000,
      "allowances": 1000000,
      "deductions": 500000,
      "net_pay": 2000000,
      "status": "processed",
      "processed_at": "2025-01-20T14:30:00Z"
    }
  ],
  "meta": {
    "total": 145,
    "page": 1,
    "pageSize": 20
  }
}
```

---

## 10. Run Pre-Check

Run comprehensive pre-check validation for a payroll run. This validates employee master data, attendance records, overtime approvals, salary structures, and compliance requirements.

**Note:** This check includes:
- Employee master data validation (new hires, terminations, salary changes)
- Attendance validation (approved timesheets for pay period)
- Overtime approval validation
- Salary structure and grade validation
- Basic salary configuration validation

```
POST /api/v1/payroll/runs/:id/pre-check
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Pre-check completed successfully",
  "data": {
    "total_employees": 145,
    "valid_employees": 142,
    "errors": 3,
    "issues": [
      {
        "id": "new_employees_incomplete_setup",
        "type": "employee_master_data",
        "category": "employee",
        "count": 2,
        "severity": "High",
        "employeeIds": ["EMP023", "EMP024"],
        "employees": [
          {
            "id": "EMP023",
            "name": "John Newhire",
            "department": "Engineering"
          }
        ],
        "description": "2 new employee(s) have incomplete setup.",
        "actionRequired": "Complete employee setup including salary structure assignment.",
        "resolutionUrl": "/employees"
      },
      {
        "id": "missing_approved_timesheets",
        "type": "attendance_validation",
        "category": "attendance",
        "count": 1,
        "severity": "High",
        "employeeIds": ["EMP045"],
        "employees": [
          {
            "id": "EMP045",
            "name": "Jane Timesheet",
            "department": "Sales"
          }
        ],
        "description": "1 employee(s) missing approved timesheets for pay period.",
        "actionRequired": "Ensure all employees have approved timesheets for the pay period.",
        "resolutionUrl": "/attendance/timesheets"
      },
      {
        "id": "pending_overtime_approvals",
        "type": "overtime_validation",
        "category": "overtime",
        "count": 3,
        "severity": "Medium",
        "employeeIds": ["EMP067", "EMP068", "EMP069"],
        "employees": [
          {
            "id": "EMP067",
            "name": "Bob Overtime",
            "department": "Operations"
          }
        ],
        "description": "3 employee(s) have pending overtime approvals.",
        "actionRequired": "Review and approve pending overtime requests.",
        "resolutionUrl": "/attendance/overtime"
      }
    ],
    "canProceed": false,
    "summary": "Payroll cannot proceed. 3 employee(s) do not have salary information."
  }
}
```

---

## 11. Calculate Payroll

Calculate payroll for all employees in the run.

```
POST /api/v1/payroll/runs/:id/calculate
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payroll calculation completed successfully",
  "data": {
    "total_employees": 145,
    "processed_employees": 145,
    "total_gross_salary": 125000000,
    "total_deductions": 25000000,
    "total_net_pay": 100000000,
    "calculation_time": "00:02:30",
    "status": "calculated"
  }
}
```

---

## 12. Advance Payroll Step

Advance to the next step in payroll processing.

```
POST /api/v1/payroll/runs/:id/advance
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payroll advanced to next step successfully",
  "data": {
    "current_step": "payslip_generation",
    "previous_step": "calculation",
    "status": "in_progress"
  }
}
```

---

## 13. Revert Payroll Step

Revert to the previous step in payroll processing.

```
POST /api/v1/payroll/runs/:id/revert
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payroll reverted to previous step successfully",
  "data": {
    "current_step": "calculation",
    "previous_step": "payslip_generation",
    "status": "in_progress"
  }
}
```

---

## 14. Generate Payslips

Generate payslips for all employees in the payroll run.

```
POST /api/v1/payroll/runs/:id/generate-payslips
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Payslips generated successfully",
  "data": {
    "total_payslips": 145,
    "generated_payslips": 145,
    "failed_payslips": 0,
    "generation_time": "00:01:45",
    "status": "completed"
  }
}
```

---

# Part C — Salary Structures & Components

Base URL: `/api/v1/payroll/salary-structures`

These endpoints manage salary structures and components for HR/Admin users.

---

## 15. List Salary Structures

List all salary structure templates.

```
GET /api/v1/payroll/salary-structures
```

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `pageSize` (optional) - Items per page (default: 20, max: 100)
- `grade` (optional) - Filter by grade
- `location` (optional) - Filter by location
- `country` (optional) - Filter by country (default: Tanzania)
- `isActive` (optional) - Filter by active status

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Engineering Grade A",
      "grade": "A",
      "location": "Dar es Salaam",
      "country": "Tanzania",
      "basic_salary": 1500000,
      "housing_allowance": 500000,
      "transport_allowance": 300000,
      "other_allowances": 200000,
      "total_allowances": 1000000,
      "gross_salary": 2500000,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "total": 15,
    "page": 1,
    "pageSize": 20
  }
}
```

---

## 16. Create Salary Structure

Create a new salary structure template.

```
POST /api/v1/payroll/salary-structures
```

**Auth:** HR/Admin only

**Request Body:**
```json
{
  "name": "Sales Grade B",
  "grade": "B",
  "location": "Dar es Salaam",
  "country": "Tanzania",
  "basic_salary": 1200000,
  "housing_allowance": 400000,
  "transport_allowance": 250000,
  "other_allowances": 150000,
  "is_active": true
}
```

### Response — `201 Created`

```json
{
  "success": true,
  "message": "Salary structure created successfully",
  "data": {
    "id": 16,
    "name": "Sales Grade B",
    "grade": "B",
    "location": "Dar es Salaam",
    "country": "Tanzania",
    "basic_salary": 1200000,
    "housing_allowance": 400000,
    "transport_allowance": 250000,
    "other_allowances": 150000,
    "total_allowances": 800000,
    "gross_salary": 2000000,
    "is_active": true,
    "created_at": "2025-02-01T12:00:00Z"
  }
}
```

---

## 17. Get Salary Structure

Get details of a specific salary structure.

```
GET /api/v1/payroll/salary-structures/:id
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Engineering Grade A",
    "grade": "A",
    "location": "Dar es Salaam",
    "country": "Tanzania",
    "basic_salary": 1500000,
    "housing_allowance": 500000,
    "transport_allowance": 300000,
    "communication_allowance": 100000,
    "medical_allowance": 100000,
    "other_allowances": 0,
    "total_allowances": 1000000,
    "gross_salary": 2500000,
    "components": [
      {
        "name": "Basic Salary",
        "type": "earning",
        "amount": 1500000,
        "is_taxable": true,
        "is_statutory": true
      },
      {
        "name": "Housing Allowance",
        "type": "earning",
        "amount": 500000,
        "is_taxable": true,
        "is_statutory": false
      }
    ],
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-06-01T00:00:00Z"
  }
}
```

---

## 18. Update Salary Structure

Update a salary structure template.

```
PUT /api/v1/payroll/salary-structures/:id
```

**Auth:** HR/Admin only

**Request Body:**
```json
{
  "name": "Engineering Grade A - Updated",
  "housing_allowance": 550000,
  "transport_allowance": 350000
}
```

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Salary structure updated successfully",
  "data": {
    "id": 1,
    "name": "Engineering Grade A - Updated",
    "grade": "A",
    "location": "Dar es Salaam",
    "country": "Tanzania",
    "basic_salary": 1500000,
    "housing_allowance": 550000,
    "transport_allowance": 350000,
    "communication_allowance": 100000,
    "medical_allowance": 100000,
    "other_allowances": 0,
    "total_allowances": 1100000,
    "gross_salary": 2600000,
    "is_active": true,
    "updated_at": "2025-02-01T13:00:00Z"
  }
}
```

---

## 19. Delete Salary Structure

Delete a salary structure template.

```
DELETE /api/v1/payroll/salary-structures/:id
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Salary structure deleted successfully"
}
```

---

## 20. Simulate Salary

Simulate salary calculation for a given amount and structure.

```
POST /api/v1/payroll/salary-structures/simulate
```

**Auth:** HR/Admin only

**Request Body:**
```json
{
  "basic_salary": 2000000,
  "housing_allowance": 600000,
  "transport_allowance": 400000,
  "other_allowances": 200000,
  "apply_tax": true,
  "apply_nhif": true,
  "apply_nssf": true
}
```

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "gross_salary": 3200000,
    "tax_amount": 480000,
    "nhif_amount": 27000,
    "nssf_amount": 200000,
    "total_deductions": 707000,
    "net_salary": 2493000,
    "breakdown": {
      "earnings": {
        "basic_salary": 2000000,
        "housing_allowance": 600000,
        "transport_allowance": 400000,
        "other_allowances": 200000,
        "total_earnings": 3200000
      },
      "deductions": {
        "tax": 480000,
        "nhif": 27000,
        "nssf": 200000,
        "total_deductions": 707000
      }
    }
  }
}
```

---

## 21. List Salary Components

List all salary components.

```
GET /api/v1/payroll/salary-components
```

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `pageSize` (optional) - Items per page (default: 20, max: 100)
- `type` (optional) - Filter by type (earning/deduction)

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Basic Salary",
      "type": "earning",
      "description": "Basic monthly salary",
      "is_taxable": true,
      "is_statutory": true,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "name": "Housing Allowance",
      "type": "earning",
      "description": "Housing allowance",
      "is_taxable": true,
      "is_statutory": false,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "total": 25,
    "page": 1,
    "pageSize": 20
  }
}
```

---

## 22. Create Salary Component

Create a new salary component.

```
POST /api/v1/payroll/salary-components
```

**Auth:** HR/Admin only

**Request Body:**
```json
{
  "name": "Performance Bonus",
  "type": "earning",
  "description": "Quarterly performance bonus",
  "is_taxable": true,
  "is_statutory": false,
  "is_active": true
}
```

### Response — `201 Created`

```json
{
  "success": true,
  "message": "Salary component created successfully",
  "data": {
    "id": 26,
    "name": "Performance Bonus",
    "type": "earning",
    "description": "Quarterly performance bonus",
    "is_taxable": true,
    "is_statutory": false,
    "is_active": true,
    "created_at": "2025-02-01T14:00:00Z"
  }
}
```

---

## 23. Get Salary Component

Get details of a specific salary component.

```
GET /api/v1/payroll/salary-components/:id
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Basic Salary",
    "type": "earning",
    "description": "Basic monthly salary",
    "is_taxable": true,
    "is_statutory": true,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-06-01T00:00:00Z"
  }
}
```

---

## 24. Update Salary Component

Update a salary component.

```
PUT /api/v1/payroll/salary-components/:id
```

**Auth:** HR/Admin only

**Request Body:**
```json
{
  "name": "Basic Salary - Updated",
  "description": "Updated basic monthly salary description",
  "is_taxable": true,
  "is_statutory": true
}
```

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Salary component updated successfully",
  "data": {
    "id": 1,
    "name": "Basic Salary - Updated",
    "type": "earning",
    "description": "Updated basic monthly salary description",
    "is_taxable": true,
    "is_statutory": true,
    "is_active": true,
    "updated_at": "2025-02-01T15:00:00Z"
  }
}
```

---

## 25. Delete Salary Component

Delete a salary component.

```
DELETE /api/v1/payroll/salary-components/:id
```

**Auth:** HR/Admin only

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Salary component deleted successfully"
}
```

---

# Part D — Payslips Management

Base URL: `/api/v1/payroll/payslips`

These endpoints manage payslips for HR/Admin users. **
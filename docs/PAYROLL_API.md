# Payroll Management API Documentation

This document provides comprehensive documentation for the Payroll Management API endpoints, including payroll runs, salary structures, payslips, loans, compliance, and reports.

## Table of Contents

1. [Authentication](#authentication)
2. [Payroll Dashboard](#payroll-dashboard)
3. [Payroll Runs](#payroll-runs)
4. [Salary Structures](#salary-structures)
5. [Salary Components](#salary-components)
6. [Payslips](#payslips)
7. [Loans & Advances](#loans--advances)
8. [Compliance & Tax](#compliance--tax)
9. [Payroll Reports](#payroll-reports)
10. [Employee Self-Service](#employee-self-service)
11. [Tanzania Statutory Calculations](#tanzania-statutory-calculations)

---

## Authentication

All payroll endpoints require authentication via JWT token. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

Most endpoints require HR or Admin role access. Employee self-service endpoints only require authentication.

---

## Payroll Dashboard

### Get Dashboard

Retrieves payroll dashboard metrics including current run, last finalized run, and loan summary.

**Endpoint:** `GET /api/v1/payroll/dashboard`

**Authorization:** HR/Admin

**Response:**
```json
{
  "success": true,
  "data": {
    "currentRun": {
      "id": 1,
      "runName": "Payroll Run January 2024",
      "status": "Draft",
      "currentStep": 0,
      "payMonth": 1,
      "payYear": 2024
    },
    "lastRun": {
      "id": 2,
      "runName": "Payroll Run December 2023",
      "payPeriod": "Dec 2023",
      "totalEmployees": 150,
      "totalNet": 450000000,
      "disbursementDate": "2023-12-28T00:00:00Z"
    },
    "recentRuns": [...],
    "currentMonth": 1,
    "currentYear": 2024,
    "loanSummary": {
      "activeLoansCount": 25,
      "totalDisbursed": 50000000,
      "totalOutstanding": 35000000,
      "thisMonthEmiTotal": 2500000,
      "pendingApprovals": 3
    }
  }
}
```

---

## Payroll Runs

### List Payroll Runs

**Endpoint:** `GET /api/v1/payroll/runs`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| status | string | Filter by status (Draft, In Review, Approved, Finalized, Disbursed, Closed) |
| year | int | Filter by year |
| month | int | Filter by month |
| page | int | Page number (default: 1) |
| pageSize | int | Items per page (default: 20, max: 100) |

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "runName": "Payroll Run January 2024",
      "payPeriod": "Jan 2024",
      "payMonth": 1,
      "payYear": 2024,
      "status": "Draft",
      "currentStep": 0,
      "totalEmployees": 150,
      "totalGross": 500000000,
      "totalDeductions": 75000000,
      "totalNet": 425000000,
      "totalEmployerContributions": 60000000
    }
  ],
  "meta": {
    "total": 12,
    "page": 1,
    "pageSize": 20,
    "pages": 1
  }
}
```

### Create Payroll Run

**Endpoint:** `POST /api/v1/payroll/runs`

**Request Body:**
```json
{
  "payMonth": 1,
  "payYear": 2024,
  "payFrequency": "Monthly",
  "runName": "Payroll Run January 2024"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Payroll run created successfully",
  "data": {
    "id": 1,
    "runName": "Payroll Run January 2024",
    "payPeriod": "Jan 2024",
    "status": "Draft",
    "currentStep": 0
  }
}
```

### Get Payroll Run

**Endpoint:** `GET /api/v1/payroll/runs/:id`

### Get Current Active Run

**Endpoint:** `GET /api/v1/payroll/runs/current`

### Get Payroll Run Summary

**Endpoint:** `GET /api/v1/payroll/runs/:id/summary`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "runName": "Payroll Run January 2024",
    "payPeriod": "Jan 2024",
    "status": "Finalized",
    "financialSummary": {
      "totalGross": 500000000,
      "totalBasic": 300000000,
      "totalAllowances": 200000000,
      "totalDeductions": 75000000,
      "totalNet": 425000000
    },
    "deductionsBreakdown": {
      "paye": 45000000,
      "nssfEmployee": 20000000,
      "nhifEmployee": 7500000,
      "loanDeductions": 2500000
    },
    "employerContributions": {
      "nssfEmployer": 20000000,
      "nhifEmployer": 7500000,
      "sdl": 25000000,
      "wcf": 5000000,
      "total": 57500000
    }
  }
}
```

### Get Payroll Employees

**Endpoint:** `GET /api/v1/payroll/runs/:id/employees`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| search | string | Search by name or employee code |
| departmentId | string | Filter by department |
| hasChanges | bool | Filter employees with changes |
| page | int | Page number |
| pageSize | int | Items per page |

### Run Pre-Check

Validates payroll data before processing.

**Endpoint:** `POST /api/v1/payroll/runs/:id/pre-check`

**Response:**
```json
{
  "success": true,
  "data": {
    "payrollRunId": 1,
    "runName": "Payroll Run January 2024",
    "totalEmployees": 150,
    "readyCount": 148,
    "issuesCount": 2,
    "canProceed": true,
    "summary": "2 employees have missing bank details",
    "issues": [
      {
        "employeeId": 45,
        "employeeName": "John Doe",
        "issueType": "MISSING_BANK",
        "severity": "Warning",
        "message": "Bank account details not configured"
      }
    ]
  }
}
```

### Calculate Payroll

Calculates salaries, deductions, and contributions for all employees.

**Endpoint:** `POST /api/v1/payroll/runs/:id/calculate`

### Advance Payroll Step

Moves payroll to the next workflow step.

**Endpoint:** `POST /api/v1/payroll/runs/:id/advance`

**Workflow Steps:**
1. Draft → In Review
2. In Review → Approved (requires approval)
3. Approved → Finalized (locks payroll)
4. Finalized → Disbursed (marks as paid)
5. Disbursed → Closed

### Revert Payroll Step

Reverts payroll to previous step (only for Draft/In Review/Approved).

**Endpoint:** `POST /api/v1/payroll/runs/:id/revert`

### Generate Payslips

Creates payslips for all employees in a finalized run.

**Endpoint:** `POST /api/v1/payroll/runs/:id/generate-payslips`

---

## Salary Structures

### List Salary Structures

**Endpoint:** `GET /api/v1/payroll/salary-structures`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| grade | string | Filter by grade level |
| location | string | Filter by location |
| country | string | Country (default: Tanzania) |
| isActive | bool | Filter active/inactive |
| page | int | Page number |
| pageSize | int | Items per page |

### Create Salary Structure

**Endpoint:** `POST /api/v1/payroll/salary-structures`

**Request Body:**
```json
{
  "templateName": "Grade A - Manager",
  "grade": "A",
  "location": "Dar es Salaam",
  "country": "Tanzania",
  "ctc": 5000000,
  "basic": 2500000,
  "hra": 500000,
  "transport": 300000,
  "medical": 200000,
  "otherAllowances": 500000
}
```

**Response:**
```json
{
  "success": true,
  "message": "Salary structure created successfully",
  "data": {
    "id": 1,
    "templateName": "Grade A - Manager",
    "grossSalary": 4000000,
    "payeDeduction": 525000,
    "nssfEmployee": 400000,
    "nhifEmployee": 20000,
    "totalDeductions": 945000,
    "nssfEmployer": 400000,
    "nhifEmployer": 20000,
    "sdl": 200000,
    "wcf": 40000,
    "totalEmployerContributions": 660000,
    "netPay": 3055000
  }
}
```

### Simulate Salary

Calculates salary breakdown from CTC without saving.

**Endpoint:** `POST /api/v1/payroll/salary-structures/simulate`

**Request Body:**
```json
{
  "ctc": 3000000
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "ctc": 3000000,
    "earnings": {
      "basic": 1500000,
      "hra": 600000,
      "transport": 300000,
      "medical": 300000,
      "otherAllowances": 300000,
      "grossSalary": 3000000
    },
    "deductions": {
      "paye": 289500,
      "nssfEmployee": 300000,
      "nhifEmployee": 15000,
      "total": 604500
    },
    "employerContributions": {
      "nssfEmployer": 300000,
      "nhifEmployer": 15000,
      "sdl": 150000,
      "wcf": 30000,
      "total": 495000
    },
    "netPay": 2395500
  }
}
```

---

## Salary Components

### List Salary Components

**Endpoint:** `GET /api/v1/payroll/salary-components`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| type | string | Filter by type (Earning, Deduction, Employer Contribution) |
| isActive | bool | Filter active/inactive |
| isStatutory | bool | Filter statutory components |
| country | string | Country (default: Tanzania) |

### Create Salary Component

**Endpoint:** `POST /api/v1/payroll/salary-components`

**Request Body:**
```json
{
  "componentName": "Housing Allowance",
  "componentCode": "HRA",
  "componentType": "Earning",
  "calculationType": "Percentage of Basic",
  "defaultPercentage": 20,
  "isTaxable": true,
  "isStatutory": false,
  "isRecurring": true,
  "displayInPayslip": true,
  "country": "Tanzania"
}
```

**Component Types:**
- `Earning` - Added to gross salary
- `Deduction` - Deducted from gross salary
- `Employer Contribution` - Not deducted from employee, paid by employer

**Calculation Types:**
- `Fixed` - Fixed amount
- `Percentage of Basic` - Calculated as percentage of basic salary
- `Percentage of Gross` - Calculated as percentage of gross salary
- `Formula` - Custom formula

---

## Payslips

### List Payslips

**Endpoint:** `GET /api/v1/payroll/payslips`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| month | int | Filter by month |
| year | int | Filter by year |
| employeeId | uint | Filter by employee |
| departmentId | uint | Filter by department |
| status | string | Filter by status (Draft, Released, Paid) |
| page | int | Page number |
| pageSize | int | Items per page |

### Get Payslip Details

**Endpoint:** `GET /api/v1/payroll/payslips/:id`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "employeeId": 45,
    "empId": "EMP001",
    "employeeName": "John Doe",
    "department": "Engineering",
    "designation": "Senior Developer",
    "payPeriod": "Jan 2024",
    "bankName": "CRDB Bank",
    "bankAccount": "****5678",
    "tinNumber": "123-456-789",
    "nssfNumber": "NSSF12345",
    "nhifNumber": "NHIF67890",
    "earnings": [
      {"code": "BASIC", "name": "Basic Salary", "amount": 2000000},
      {"code": "HRA", "name": "Housing Allowance", "amount": 400000},
      {"code": "TRANSPORT", "name": "Transport Allowance", "amount": 200000}
    ],
    "deductions": [
      {"code": "PAYE", "name": "PAYE Tax", "amount": 196500},
      {"code": "NSSF_EMP", "name": "NSSF (Employee)", "amount": 260000},
      {"code": "NHIF_EMP", "name": "NHIF (Employee)", "amount": 15000}
    ],
    "employerContributions": [
      {"code": "NSSF_EMPLOYER", "name": "NSSF (Employer)", "amount": 260000},
      {"code": "NHIF_EMPLOYER", "name": "NHIF (Employer)", "amount": 15000},
      {"code": "SDL", "name": "Skills Development Levy", "amount": 130000},
      {"code": "WCF", "name": "Workers Compensation Fund", "amount": 26000}
    ],
    "summary": {
      "grossSalary": 2600000,
      "totalDeductions": 471500,
      "netPay": 2128500
    },
    "ytd": {
      "grossSalary": 2600000,
      "tax": 196500,
      "nssf": 260000,
      "nhif": 15000,
      "netPay": 2128500
    },
    "attendance": {
      "workingDays": 22,
      "daysWorked": 22,
      "lopDays": 0,
      "overtimeHours": 0
    },
    "status": "Released"
  }
}
```

### Download Payslip

**Endpoint:** `GET /api/v1/payroll/payslips/:id/download`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| format | string | Output format: `pdf` (default) or `xlsx` |

### Send Payslip Email

**Endpoint:** `POST /api/v1/payroll/payslips/:id/email`

### Bulk Download Payslips

**Endpoint:** `POST /api/v1/payroll/payslips/bulk-download`

**Request Body:**
```json
{
  "payslipIds": [1, 2, 3, 4, 5]
}
```

### Release Payslips

Makes payslips visible to employees.

**Endpoint:** `POST /api/v1/payroll/payslips/run/:runId/release`

### Bulk Email Payslips

**Endpoint:** `POST /api/v1/payroll/payslips/run/:runId/bulk-email`

---

## Loans & Advances

### List Loans

**Endpoint:** `GET /api/v1/payroll/loans`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| employeeId | uint | Filter by employee |
| status | string | Filter by status |
| type | string | Filter by loan type |
| page | int | Page number |
| pageSize | int | Items per page |

**Loan Types:**
- `Education Loan`
- `Emergency Loan`
- `Personal Loan`
- `Salary Advance`

**Loan Statuses:**
- `Pending Approval`
- `Active`
- `Closed`
- `Defaulted`
- `Rejected`

### Create Loan Application

**Endpoint:** `POST /api/v1/payroll/loans`

**Request Body:**
```json
{
  "employeeId": 45,
  "empId": "EMP001",
  "employeeName": "John Doe",
  "department": "Engineering",
  "loanType": "Personal Loan",
  "amount": 5000000,
  "interestRate": 10,
  "tenure": 12,
  "purpose": "Home renovation"
}
```

### Get Loan Details

**Endpoint:** `GET /api/v1/payroll/loans/:id`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "employeeId": 45,
    "employeeName": "John Doe",
    "loanType": "Personal Loan",
    "amount": 5000000,
    "interestRate": 10,
    "tenure": 12,
    "emiAmount": 439583.33,
    "status": "Active",
    "totalPaid": 1318750,
    "outstandingBalance": 4181250,
    "paymentsMade": 3,
    "paymentsRemaining": 9,
    "nextEmiDate": "2024-02-01",
    "repaymentSchedule": [
      {
        "installment": 1,
        "dueDate": "2024-01-01",
        "emiAmount": 439583.33,
        "principal": 397916.67,
        "interest": 41666.67,
        "status": "Paid",
        "paidDate": "2024-01-01"
      },
      {
        "installment": 2,
        "dueDate": "2024-02-01",
        "emiAmount": 439583.33,
        "principal": 401232.64,
        "interest": 38350.69,
        "status": "Pending"
      }
    ]
  }
}
```

### Approve Loan

**Endpoint:** `POST /api/v1/payroll/loans/:id/approve`

**Request Body:**
```json
{
  "disbursementDate": "2024-01-15T00:00:00Z",
  "remarks": "Approved based on employment tenure"
}
```

### Reject Loan

**Endpoint:** `POST /api/v1/payroll/loans/:id/reject`

**Request Body:**
```json
{
  "reason": "Outstanding loan already exists"
}
```

### Get Repayment Schedule

**Endpoint:** `GET /api/v1/payroll/loans/:id/schedule`

### Record Repayment

**Endpoint:** `POST /api/v1/payroll/loans/:id/repayment`

**Request Body:**
```json
{
  "installment": 4,
  "paidAmount": 439583.33,
  "payslipId": 156
}
```

### Get Loan Summary

**Endpoint:** `GET /api/v1/payroll/loans/summary`

### Calculate EMI

**Endpoint:** `POST /api/v1/payroll/loans/calculate-emi`

**Request Body:**
```json
{
  "principal": 5000000,
  "interestRate": 10,
  "tenure": 12
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "principal": 5000000,
    "interestRate": 10,
    "tenure": 12,
    "emi": 439583.33,
    "totalPayable": 5275000,
    "totalInterest": 275000
  }
}
```

---

## Compliance & Tax

### Get Compliance Summary

**Endpoint:** `GET /api/v1/payroll/compliance/summary`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| month | int | Pay month (default: current) |
| year | int | Pay year (default: current) |

**Response:**
```json
{
  "success": true,
  "data": {
    "payMonth": 1,
    "payYear": 2024,
    "nextDueDate": "2024-02-07",
    "daysUntilDue": 15,
    "paye": {
      "amount": 45000000,
      "employeesCount": 150,
      "dueDate": "2024-02-07",
      "status": "Pending"
    },
    "nssf": {
      "employeeAmount": 20000000,
      "employerAmount": 20000000,
      "totalAmount": 40000000,
      "employeesCount": 150,
      "dueDate": "2024-02-07",
      "status": "Pending"
    },
    "nhif": {
      "employeeAmount": 2250000,
      "employerAmount": 2250000,
      "totalAmount": 4500000,
      "employeesCount": 150,
      "dueDate": "2024-02-07",
      "status": "Pending"
    },
    "sdl": {
      "amount": 25000000,
      "dueDate": "2024-02-07",
      "status": "Pending"
    },
    "wcf": {
      "amount": 5000000,
      "dueDate": "2024-02-07",
      "status": "Pending"
    },
    "totalStatutory": 119500000
  }
}
```

### Get All Compliance Rules

**Endpoint:** `GET /api/v1/payroll/compliance/rules`

**Response:**
```json
{
  "success": true,
  "data": {
    "country": "Tanzania",
    "taxYear": 2024,
    "payeSlabs": [
      {"from": 0, "to": 270000, "rate": 0, "fixedAmount": 0, "description": "Tax-free threshold"},
      {"from": 270001, "to": 520000, "rate": 9, "fixedAmount": 0, "description": "9% on amount above 270,000"},
      {"from": 520001, "to": 760000, "rate": 20, "fixedAmount": 22500, "description": "TZS 22,500 + 20% on amount above 520,000"},
      {"from": 760001, "to": 1000000, "rate": 25, "fixedAmount": 70500, "description": "TZS 70,500 + 25% on amount above 760,000"},
      {"from": 1000001, "to": 999999999, "rate": 30, "fixedAmount": 130500, "description": "TZS 130,500 + 30% on amount above 1,000,000"}
    ],
    "statutoryRules": {
      "NSSF": {
        "name": "National Social Security Fund",
        "employeeRate": 10,
        "employerRate": 10,
        "basis": "Gross salary",
        "description": "Employee: 10% of gross salary. Employer: 10% of gross salary."
      },
      "NHIF": {
        "name": "National Health Insurance Fund",
        "basis": "Schedule",
        "description": "Contribution as per NHIF schedule based on salary bands."
      },
      "SDL": {
        "name": "Skills Development Levy",
        "employerRate": 5,
        "basis": "Total payroll",
        "description": "Employer only: 5% of total gross payroll."
      },
      "WCF": {
        "name": "Workers Compensation Fund",
        "employerRate": 1,
        "basis": "Total payroll",
        "description": "Employer only: 1% of total gross payroll."
      }
    },
    "nhifSchedule": [
      {"from": 0, "to": 150000, "employee": 2500, "employer": 2500},
      {"from": 150001, "to": 250000, "employee": 5000, "employer": 5000},
      {"from": 250001, "to": 400000, "employee": 10000, "employer": 10000},
      {"from": 400001, "to": 600000, "employee": 15000, "employer": 15000},
      {"from": 600001, "to": 999999999, "employee": 20000, "employer": 20000}
    ]
  }
}
```

### Get Tax Slabs

**Endpoint:** `GET /api/v1/payroll/compliance/tax-slabs`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| country | string | Country (default: Tanzania) |
| year | int | Tax year (default: current) |

### Get Statutory Rules

**Endpoint:** `GET /api/v1/payroll/compliance/statutory-rules`

### Get NHIF Schedule

**Endpoint:** `GET /api/v1/payroll/compliance/nhif-schedule`

### Simulate Tax

**Endpoint:** `POST /api/v1/payroll/compliance/simulate-tax`

**Request Body:**
```json
{
  "grossSalary": 2500000,
  "taxYear": 2024
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "grossSalary": 2500000,
    "taxableIncome": 2250000,
    "paye": 305000,
    "nssfEmployee": 250000,
    "nssfEmployer": 250000,
    "nhifEmployee": 15000,
    "nhifEmployer": 15000,
    "sdl": 125000,
    "wcf": 25000,
    "totalEmployeeDeductions": 570000,
    "totalEmployerContributions": 415000,
    "netPay": 1930000,
    "effectiveTaxRate": 12.2
  }
}
```

### Calculate Monthly Statutory

**Endpoint:** `POST /api/v1/payroll/compliance/calculate-monthly`

**Request Body:**
```json
{
  "totalGross": 500000000,
  "employeeCount": 150,
  "payMonth": 1,
  "payYear": 2024
}
```

### Generate Compliance Return

**Endpoint:** `GET /api/v1/payroll/compliance/return`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| month | int | Pay month |
| year | int | Pay year |
| type | string | Return type: PAYE, NSSF, NHIF, SDL, WCF |

### Seed Compliance Data

Seeds Tanzania compliance rules (tax slabs, statutory rules, NHIF schedule).

**Endpoint:** `POST /api/v1/payroll/compliance/seed`

---

## Payroll Reports

### Get Payroll Report

**Endpoint:** `GET /api/v1/payroll/reports`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| month | int | Pay month (default: current) |
| year | int | Pay year (default: current) |

### Get Payroll KPIs

**Endpoint:** `GET /api/v1/payroll/reports/kpis`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| year | int | Year (default: current) |

**Response:**
```json
{
  "success": true,
  "data": {
    "year": 2024,
    "totalRuns": 1,
    "totalGross": 500000000,
    "totalNet": 425000000,
    "totalDeductions": 75000000,
    "totalEmployerContributions": 60000000,
    "averageMonthlyCost": 560000000,
    "monthlyBreakdown": [
      {
        "month": 1,
        "totalGross": 500000000,
        "totalNet": 425000000,
        "totalDeductions": 75000000,
        "employerContributions": 60000000,
        "totalEmployees": 150
      }
    ]
  }
}
```

### Get Department Cost Report

**Endpoint:** `GET /api/v1/payroll/reports/department-cost`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| month | int | Pay month |
| year | int | Pay year |

**Response:**
```json
{
  "success": true,
  "data": {
    "payMonth": 1,
    "payYear": 2024,
    "departments": [
      {
        "department": "Engineering",
        "employeeCount": 45,
        "totalGross": 180000000,
        "totalDeductions": 27000000,
        "totalNet": 153000000
      },
      {
        "department": "Sales",
        "employeeCount": 30,
        "totalGross": 120000000,
        "totalDeductions": 18000000,
        "totalNet": 102000000
      }
    ]
  }
}
```

### Export Payroll Report

**Endpoint:** `GET /api/v1/payroll/reports/export`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| month | int | Pay month |
| year | int | Pay year |
| format | string | Output format: xlsx (default) |

Downloads an Excel file with Summary and Employee Details sheets.

### Get Bank File Report

**Endpoint:** `GET /api/v1/payroll/reports/run/:runId/bank-file`

Downloads an Excel file formatted for bank salary disbursement.

---

## Employee Self-Service

### Get My Payslips

**Endpoint:** `GET /api/v1/self-service/payroll/payslips`

**Authorization:** Any authenticated employee

### Get Latest Payslip

**Endpoint:** `GET /api/v1/self-service/payroll/payslips/latest`

### Get Salary Summary

**Endpoint:** `GET /api/v1/self-service/payroll/payslips/summary`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| month | int | Pay month |
| year | int | Pay year |

### Download My Payslip

**Endpoint:** `GET /api/v1/self-service/payroll/payslips/:id/download`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| format | string | Output format: pdf (default) or xlsx |

### Get My Loans

**Endpoint:** `GET /api/v1/self-service/payroll/loans`

### Apply for Loan

**Endpoint:** `POST /api/v1/self-service/payroll/loans/apply`

**Request Body:**
```json
{
  "loanType": "Personal Loan",
  "amount": 2000000,
  "tenure": 6,
  "purpose": "Family emergency"
}
```

---

## Tanzania Statutory Calculations

### PAYE (Pay As You Earn) Tax

Tanzania uses a progressive tax system with monthly brackets:

| Annual Income (TZS) | Monthly Income (TZS) | Tax Rate |
|---------------------|----------------------|----------|
| 0 - 3,240,000 | 0 - 270,000 | 0% |
| 3,240,001 - 6,240,000 | 270,001 - 520,000 | 9% |
| 6,240,001 - 9,120,000 | 520,001 - 760,000 | 20% |
| 9,120,001 - 12,000,000 | 760,001 - 1,000,000 | 25% |
| Above 12,000,000 | Above 1,000,000 | 30% |

**Note:** PAYE is calculated on taxable income (gross salary minus NSSF employee contribution).

### NSSF (National Social Security Fund)

- **Employee Contribution:** 10% of gross salary
- **Employer Contribution:** 10% of gross salary
- **Total:** 20% of gross salary

### NHIF (National Health Insurance Fund)

Fixed amounts based on salary bands:

| Gross Salary (TZS) | Employee | Employer |
|--------------------|----------|----------|
| 0 - 150,000 | 2,500 | 2,500 |
| 150,001 - 250,000 | 5,000 | 5,000 |
| 250,001 - 400,000 | 10,000 | 10,000 |
| 400,001 - 600,000 | 15,000 | 15,000 |
| Above 600,000 | 20,000 | 20,000 |

### SDL (Skills Development Levy)

- **Employer Only:** 5% of total gross payroll
- Remitted monthly to TRA

### WCF (Workers Compensation Fund)

- **Employer Only:** 1% of total gross payroll
- Provides insurance coverage for workplace injuries

### Filing Deadlines

All statutory contributions must be remitted by the **7th day of the following month**.

---

## Error Responses

All endpoints return consistent error responses:

```json
{
  "success": false,
  "error": "Error message description"
}
```

**Common HTTP Status Codes:**
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid/missing token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error

---

## Rate Limiting

API endpoints are subject to rate limiting:
- 100 requests per minute for standard endpoints
- 10 requests per minute for report generation/export

---

## Changelog

### Version 1.0.0 (January 2024)
- Initial release
- Complete payroll run workflow (6-step process)
- Tanzania statutory compliance (PAYE, NSSF, NHIF, SDL, WCF)
- Salary structure templates and components
- Payslip generation and distribution
- Loan management with EMI calculations
- Comprehensive reporting and exports
- Employee self-service portal

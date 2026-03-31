# Payroll Configuration API Documentation

This document describes the Payroll Configuration API endpoints for Tanzania statutory rates and calculations. These configurations are **read-only** and represent Tanzania's legal statutory rates set by law.

**Base URL:** `/api/v1`

## ⚠️ Important Notice

These rates are **set by Tanzanian law** and **cannot be edited** through the application UI:
- TRA (Tanzania Revenue Authority) sets PAYE bands
- NSSF Act sets contribution percentages  
- SDL Act sets Skills and Development Levy rate
- WCF Act sets Workers Compensation Fund rate
- HESLB Act sets student loan repayment rate

**No user** — including Super Admin — can edit these rates through the UI. Changes require database migration by a developer.

---

## Table of Contents

1. [Get All Configurations](#1-get-all-configurations)
2. [Get Single Configuration](#2-get-single-configuration)
3. [Calculate PAYE Tax](#3-calculate-paye-tax)
4. [Calculate NSSF Contributions](#4-calculate-nssf-contributions)
5. [Calculate Cost to Company (CTC)](#5-calculate-cost-to-company-ctc)
6. [Get PAYE Tax Bands](#6-get-paye-tax-bands)
7. [Get NSSF Rates](#7-get-nssf-rates)
8. [Get Employer Contributions](#8-get-employer-contributions)
9. [Get Formula Configurations](#9-get-formula-configurations)

---

## 1. Get All Configurations

Returns all payroll configuration values as a single JSON object.

```
GET /api/v1/payroll/config
```

**Auth:** Required (any authenticated user)

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "NSSF_EMPLOYEE_RATE": 0.10,
    "NSSF_EMPLOYER_RATE": 0.10,
    "NSSF_TOTAL_RATE": 0.20,
    "WCF_EMPLOYER_RATE": 0.005,
    "SDL_EMPLOYER_RATE": 0.035,
    "HESLB_RATE": 0.15,
    "PAYROLL_CURRENCY": "TZS",
    "PAYROLL_FREQUENCY": "MONTHLY",
    "TAXABLE_PAY_FORMULA": "GROSS_SALARY - NSSF_EMPLOYEE",
    "NET_PAY_FORMULA": "TAXABLE_PAY - PAYE - HESLB - LOAN_DEDUCTION",
    "CTC_FORMULA": "GROSS + NSSF_EMPLOYER + WCF + SDL",
    "PAYE_BANDS": [
      {
        "band": 1,
        "min": 0,
        "max": 270000,
        "rate": 0.00,
        "base_tax": 0,
        "label": "Up to TZS 270,000 — 0%"
      },
      {
        "band": 2,
        "min": 270001,
        "max": 520000,
        "rate": 0.08,
        "base_tax": 0,
        "label": "TZS 270,001 - 520,000 — 8%"
      },
      {
        "band": 3,
        "min": 520001,
        "max": 760000,
        "rate": 0.20,
        "base_tax": 20000,
        "label": "TZS 520,001 - 760,000 — 20%"
      },
      {
        "band": 4,
        "min": 760001,
        "max": 1000000,
        "rate": 0.25,
        "base_tax": 68000,
        "label": "TZS 760,001 - 1,000,000 — 25%"
      },
      {
        "band": 5,
        "min": 1000001,
        "rate": 0.30,
        "base_tax": 128000,
        "label": "Above TZS 1,000,000 — 30%"
      }
    ]
  }
}
```

---

## 2. Get Single Configuration

Returns a specific configuration by key.

```
GET /api/v1/payroll/config/:key
```

**Path Parameters:**
- `key` - Configuration key (e.g., `NSSF_EMPLOYEE_RATE`)

**Auth:** Required (any authenticated user)

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "config_key": "NSSF_EMPLOYEE_RATE",
    "value": 0.10,
    "label": "NSSF Employee Contribution",
    "description": "10% of Gross Salary deducted from employee. Formula from file: =H*10% Paid to: NSSF Tanzania",
    "category": "NSSF",
    "editable": false
  }
}
```

### Response — `404 Not Found`

```json
{
  "success": false,
  "error": "Configuration not found"
}
```

---

## 3. Calculate PAYE Tax

Calculates PAYE tax for a given taxable pay amount using Tanzania tax bands.

```
GET /api/v1/payroll/config/paye/calculate?taxable_pay=1350000
```

**Query Parameters:**
- `taxable_pay` (required) - Taxable pay amount in TZS

**Auth:** Required (any authenticated user)

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "taxable_pay": 1350000,
    "band_applied": 5,
    "band_label": "Above TZS 1,000,000 — 30%",
    "paye_tax": 233000,
    "breakdown": {
      "base_tax": 128000,
      "rate": 0.30,
      "excess_amount": 350000,
      "rate_tax": 105000
    }
  }
}
```

---

## 4. Calculate NSSF Contributions

Calculates NSSF contributions (employee and employer) for a given gross salary.

```
GET /api/v1/payroll/config/nssf/calculate?gross=1500000
```

**Query Parameters:**
- `gross` (required) - Gross salary amount in TZS

**Auth:** Required (any authenticated user)

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "gross_salary": 1500000,
    "nssf_employee": 150000,
    "nssf_employer": 150000,
    "nssf_total": 300000,
    "rate_employee": 0.10,
    "rate_employer": 0.10
  }
}
```

---

## 5. Calculate Cost to Company (CTC)

Calculates full cost to company including employer contributions.

```
GET /api/v1/payroll/config/ctc/calculate?gross=1500000
```

**Query Parameters:**
- `gross` (required) - Gross salary amount in TZS

**Auth:** Required (any authenticated user)

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "gross_salary": 1500000,
    "nssf_employer": 150000,
    "wcf": 7500,
    "sdl": 52500,
    "total_ctc": 1710000
  }
}
```

---

## 6. Get PAYE Tax Bands

Returns all PAYE tax bands with their rates and limits.

```
GET /api/v1/payroll/config/paye/bands
```

**Auth:** Required (any authenticated user)

### Response — `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "band": 1,
      "min": 0,
      "max": 270000,
      "rate": 0.00,
      "base_tax": 0,
      "label": "Up to TZS 270,000 — 0%"
    },
    {
      "band": 2,
      "min": 270001,
      "max": 520000,
      "rate": 0.08,
      "base_tax": 0,
      "label": "TZS 270,001 - 520,000 — 8%"
    },
    {
      "band": 3,
      "min": 520001,
      "max": 760000,
      "rate": 0.20,
      "base_tax": 20000,
      "label": "TZS 520,001 - 760,000 — 20%"
    },
    {
      "band": 4,
      "min": 760001,
      "max": 1000000,
      "rate": 0.25,
      "base_tax": 68000,
      "label": "TZS 760,001 - 1,000,000 — 25%"
    },
    {
      "band": 5,
      "min": 1000001,
      "rate": 0.30,
      "base_tax": 128000,
      "label": "Above TZS 1,000,000 — 30%"
    }
  ]
}
```

---

## 7. Get NSSF Rates

Returns all NSSF contribution rates.

```
GET /api/v1/payroll/config/nssf/rates
```

**Auth:** Required (any authenticated user)

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "NSSF_EMPLOYEE_RATE": 0.10,
    "NSSF_EMPLOYER_RATE": 0.10,
    "NSSF_TOTAL_RATE": 0.20
  }
}
```

---

## 8. Get Employer Contributions

Returns all employer contribution rates (NSSF, WCF, SDL).

```
GET /api/v1/payroll/config/employer-contributions
```

**Auth:** Required (any authenticated user)

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "NSSF_EMPLOYER_RATE": 0.10,
    "WCF_EMPLOYER_RATE": 0.005,
    "SDL_EMPLOYER_RATE": 0.035
  }
}
```

---

## 9. Get Formula Configurations

Returns all payroll calculation formulas.

```
GET /api/v1/payroll/config/formulas
```

**Auth:** Required (any authenticated user)

### Response — `200 OK`

```json
{
  "success": true,
  "data": {
    "TAXABLE_PAY_FORMULA": "GROSS_SALARY - NSSF_EMPLOYEE",
    "NET_PAY_FORMULA": "TAXABLE_PAY - PAYE - HESLB - LOAN_DEDUCTION",
    "CTC_FORMULA": "GROSS + NSSF_EMPLOYER + WCF + SDL"
  }
}
```

---

## 🔒 Modification Prevention

Any attempt to modify configurations via POST, PUT, or DELETE will be blocked:

```json
{
  "success": false,
  "error": "Configuration modifications are not allowed",
  "details": "These rates are set by Tanzanian law and can only be changed by a system administrator via database migration"
}
```

---

## 📊 Tanzania Statutory Rates Reference

### NSSF (National Social Security Fund)
- **Employee Rate**: 10% of Gross Salary
- **Employer Rate**: 10% of Gross Salary  
- **Total Rate**: 20% of Gross Salary

### PAYE Tax Bands (Monthly)
| Band | Taxable Pay Range | Rate | Base Tax |
|------|-------------------|------|----------|
| 1 | Up to TZS 270,000 | 0% | TZS 0 |
| 2 | TZS 270,001 - 520,000 | 8% | TZS 0 |
| 3 | TZS 520,001 - 760,000 | 20% | TZS 20,000 |
| 4 | TZS 760,001 - 1,000,000 | 25% | TZS 68,000 |
| 5 | Above TZS 1,000,000 | 30% | TZS 128,000 |

### Employer Contributions Only
- **WCF (Workers Compensation Fund)**: 0.5% of Gross Salary
- **SDL (Skills Development Levy)**: 3.5% of Gross Salary

### HESLB (Higher Education Student Loans Board)
- **Rate**: 15% of Basic Pay (only for employees with HESLB flag)

---

## 💡 Usage Notes

- All monetary values are in **TZS** (Tanzanian Shillings)
- All calculations return **integer values** (no decimal places)
- Values are **rounded using standard rounding rules**
- Configuration values are **cached for 1 hour** for performance
- The payroll calculation engine **must always** read rates from these APIs
- If the config API is unreachable, payroll calculation should be blocked with message: "Cannot load payroll configuration. Contact admin."
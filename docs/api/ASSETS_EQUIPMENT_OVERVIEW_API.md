# Assets & Equipment — Overview Cards API

This document defines the API required to power the **Assets & Equipment** page top summary cards:
- **Total Assets**
- **In Use**
- **Available**
- **Total Value**

## Current Frontend Status
On the current page implementation, **these cards are not fully driven by API**:
- `Total Assets` uses the API `total` count when available, otherwise falls back to mock data.
- `In Use`, `Available`, and `Total Value` are calculated from mock data, not from API.

This API is intended to make all four cards fully backend-driven.

## Base URL
`http://localhost:8080/api/v1`

## Authentication
All endpoints require a valid Bearer token in the Authorization header.

`Authorization: Bearer <token>`

## Tenancy (System Flow Alignment)
- Current backend runs in single-tenant mode (no tenant/org scoping applied at middleware level).
- Access is protected by Auth + HR middleware under `/api/v1/assets/*`.

---

## 1) Get Asset Overview (Cards)

### Endpoint
`GET /assets/overview`

### Query Parameters (Optional)
| Parameter | Type | Required | Description |
|---|---|---:|---|
| `asset_type` | string | No | Filter by asset type (e.g. `laptop`, `mobile_phone`) |
| `status` | string | No | Filter by status (`available`, `in_use`, `under_repair`, `retired`) |
| `department` | string | No | Filter by department (uses the same value stored in `assets.department`, typically the department ID as string) |
| `currency` | string | No | Currency code for totals (echoed back in response; no conversion performed) |

### Success Response
```json
{
  "success": true,
  "message": "Asset overview retrieved successfully",
  "data": {
    "filters": {
      "asset_type": null,
      "status": null,
      "department": null,
      "currency": "TZS"
    },
    "totals": {
      "total_assets": 120,
      "in_use": 85,
      "available": 25,
      "under_repair": 7,
      "retired": 3,
      "total_value": 950000000
    },
    "meta": {
      "value_unit": "decimal",
      "currency": "TZS"
    }
  }
}
```

### Notes
- `total_value` is computed as `SUM(assets.value)`; current schema stores it as decimal mapped to float in Go.
- `currency` is included for UI display only (no currency conversion is performed).
- The UI cards only need:
  - `totals.total_assets`
  - `totals.in_use`
  - `totals.available`
  - `totals.total_value`
- Returning additional totals (`under_repair`, `retired`) is useful for future UI expansion.

---

## Error Responses

### 400 Bad Request (invalid filters)
```json
{
  "success": false,
  "message": "Invalid request",
  "error": {
    "status": "invalid value"
  }
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "message": "Forbidden"
}
```

---

## Implementation Guidance (Backend)
- Compute counts from the same canonical `assets` table used by `/assets` listing.
- Status normalization recommended:
  - `available`
  - `in_use`
  - `under_repair`
  - `retired`
- Ensure the overview respects the same filters/visibility rules as the assets list (RBAC).

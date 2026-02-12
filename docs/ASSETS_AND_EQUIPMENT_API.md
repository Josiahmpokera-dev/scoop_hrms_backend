# Assets & Equipment API Documentation

This document covers **all** Assets & Equipment endpoints — Admin/HR management and Employee self-service (asset requests and issue reports).

---

## Table of Contents

### Part A — Admin/HR Asset Management
1. [List Assets](#1-list-assets)
2. [Get Asset by ID](#2-get-asset-by-id)
3. [Get Asset Types](#3-get-asset-types)
4. [Create Asset](#4-create-asset)
5. [Update Asset](#5-update-asset)
6. [Assign Asset to Employee](#6-assign-asset-to-employee)
7. [Reassign Asset (Asset Handler)](#7-reassign-asset-asset-handler)
8. [Return Asset](#8-return-asset)
9. [Mark Asset for Repair](#9-mark-asset-for-repair)
10. [Complete Repair](#10-complete-repair)
11. [Retire Asset](#11-retire-asset)
12. [Delete Asset](#12-delete-asset)

### Part B — HR/Admin: Asset Request Management
13. [List All Asset Requests (HR)](#13-list-all-asset-requests-hr)
14. [Get Asset Request Details (HR)](#14-get-asset-request-details-hr)
15. [Approve Asset Request](#15-approve-asset-request)
16. [Reject Asset Request](#16-reject-asset-request)
17. [Fulfill Asset Request](#17-fulfill-asset-request)
18. [Reassign Asset (HR Handler)](#18-reassign-asset-hr-handler)

### Part C — Employee Self-Service
19. [Get My Assigned Assets](#19-get-my-assigned-assets)
20. [Create Asset Request (Employee)](#20-create-asset-request-employee)
21. [List My Asset Requests](#21-list-my-asset-requests)
22. [Get My Asset Request Details](#22-get-my-asset-request-details)
23. [Cancel Asset Request](#23-cancel-asset-request)
24. [Report Asset Issue](#24-report-asset-issue)
25. [List My Asset Issues](#25-list-my-asset-issues)
26. [Get My Asset Issue Details](#26-get-my-asset-issue-details)

### Appendix
- [Data Models](#data-models)
- [Quick Reference Table](#quick-reference-table)

---

## Authentication

All endpoints require `Authorization: Bearer <token>`.

| Scope      | Middleware             | Who Can Access             |
|------------|------------------------|----------------------------|
| Admin/HR   | `AuthMiddleware` + `HRMiddleware` | Admin, Super Admin, HR    |
| Self-Service | `AuthMiddleware`     | Any authenticated employee |

---

# Part A — Admin/HR Asset Management

Base URL: `/api/v1/assets`

---

## 1. List Assets

Retrieve a paginated list of all assets with optional filters.

```
GET /api/v1/assets
```

| Query Param   | Type   | Required | Default | Description                                  |
|---------------|--------|----------|---------|----------------------------------------------|
| `page`        | int    | No       | 1       | Page number                                  |
| `page_size`   | int    | No       | 20      | Items per page (max 100)                     |
| `search`      | string | No       |         | Search by asset code, brand, model, serial   |
| `status`      | string | No       |         | `available`, `in_use`, `under_repair`, `retired` |
| `asset_type`  | string | No       |         | `laptop`, `mobile_phone`, `desktop`, etc.    |
| `assigned_to` | string | No       |         | Filter by employee ID (e.g. `EMP001`)        |
| `department`  | string | No       |         | Filter by department                         |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Assets retrieved successfully",
  "data": [
    {
      "id": 1,
      "asset_code": "AST-001",
      "asset_type": "laptop",
      "brand": "Dell",
      "model": "Latitude 5520",
      "serial_number": "SN123456",
      "assigned_to": "Kevin Massawe",
      "employee_id": "EMP001",
      "employee_photo": "/img/avatars/default.jpg",
      "department": "IT",
      "assigned_date": "2026-01-15T00:00:00Z",
      "return_date": null,
      "status": "in_use",
      "condition": "good",
      "purchase_date": "2025-06-01T00:00:00Z",
      "warranty_expiry": "2028-06-01T00:00:00Z",
      "value": 1500000,
      "notes": null,
      "repair_reason": null,
      "repair_notes": null,
      "repair_date": null,
      "repair_cost": null,
      "repair_completed_date": null,
      "retirement_reason": null,
      "retirement_date": null,
      "created_at": "2025-06-01T10:00:00Z",
      "updated_at": "2026-01-15T12:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

---

## 2. Get Asset by ID

```
GET /api/v1/assets/get?id=1
```

| Query Param | Type | Required | Description |
|-------------|------|----------|-------------|
| `id`        | int  | Yes      | Asset ID    |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset retrieved successfully",
  "data": {
    "id": 1,
    "asset_code": "AST-001",
    "asset_type": "laptop",
    "brand": "Dell",
    "model": "Latitude 5520",
    "serial_number": "SN123456",
    "status": "in_use",
    "condition": "good",
    "assigned_to": "Kevin Massawe",
    "employee_id": "EMP001",
    "employee_photo": "/img/avatars/default.jpg",
    "department": "IT",
    "assigned_date": "2026-01-15T00:00:00Z",
    "purchase_date": "2025-06-01T00:00:00Z",
    "warranty_expiry": "2028-06-01T00:00:00Z",
    "value": 1500000,
    "created_at": "2025-06-01T10:00:00Z",
    "updated_at": "2026-01-15T12:00:00Z"
  }
}
```

---

## 3. Get Asset Types

Returns all supported asset types.

```
GET /api/v1/assets/types
```

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset types retrieved successfully",
  "data": [
    { "id": 1, "name": "Laptop", "code": "laptop", "description": "Laptops and notebooks" },
    { "id": 2, "name": "Mobile Phone", "code": "mobile_phone", "description": "Mobile phones" },
    { "id": 3, "name": "Desktop", "code": "desktop", "description": "Desktop computers" },
    { "id": 4, "name": "Tablet", "code": "tablet", "description": "Tablets" },
    { "id": 5, "name": "ID Card", "code": "id_card", "description": "ID cards" },
    { "id": 6, "name": "Access Card", "code": "access_card", "description": "Access cards" },
    { "id": 7, "name": "Vehicle", "code": "vehicle", "description": "Company vehicles" },
    { "id": 8, "name": "Office Equipment", "code": "office_equipment", "description": "Office equipment" }
  ]
}
```

---

## 4. Create Asset

Register a new asset/equipment in the system.

```
POST /api/v1/assets
```

### Request Body

```json
{
  "asset_type": "laptop",
  "brand": "Dell",
  "model": "Latitude 5520",
  "serial_number": "SN123456",
  "asset_code": "AST-001",
  "purchase_date": "2025-06-01",
  "value": 1500000,
  "warranty_expiry": "2028-06-01",
  "condition": "excellent",
  "assigned_to_employee_id": "EMP001",
  "notes": "Purchased for IT department"
}
```

| Field                     | Type          | Required | Description                                         |
|---------------------------|---------------|----------|-----------------------------------------------------|
| `asset_type`              | string        | Yes      | `laptop`, `mobile_phone`, `desktop`, `tablet`, etc. |
| `brand`                   | string        | Yes      | Brand name (max 100)                                |
| `model`                   | string        | Yes      | Model name (max 100)                                |
| `serial_number`           | string        | Yes      | Unique serial number (max 100)                      |
| `asset_code`              | string        | No       | Custom code (auto-generated if omitted)             |
| `purchase_date`           | string        | No       | `YYYY-MM-DD`                                        |
| `value`                   | number        | No       | Purchase value                                      |
| `warranty_expiry`         | string        | No       | `YYYY-MM-DD`                                        |
| `condition`               | string        | No       | `excellent`, `good`, `fair`, `poor` (default: `excellent`) |
| `assigned_to_employee_id` | string/number | No       | Employee ID to assign immediately                   |
| `notes`                   | string        | No       | Additional notes                                    |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset created successfully",
  "data": { ... }
}
```

---

## 5. Update Asset

Update asset details.

```
POST /api/v1/assets/update
```

### Request Body

```json
{
  "id": 1,
  "brand": "Dell",
  "model": "Latitude 5530",
  "condition": "good",
  "value": 1600000,
  "notes": "Updated model"
}
```

| Field            | Type   | Required | Description                          |
|------------------|--------|----------|--------------------------------------|
| `id`             | int    | Yes      | Asset ID                             |
| `asset_type`     | string | No       | Updated asset type                   |
| `brand`          | string | No       | Updated brand (max 100)              |
| `model`          | string | No       | Updated model (max 100)              |
| `serial_number`  | string | No       | Updated serial number (max 100)      |
| `condition`      | string | No       | `excellent`, `good`, `fair`, `poor`  |
| `value`          | number | No       | Updated value                        |
| `purchase_date`  | string | No       | `YYYY-MM-DD`                         |
| `warranty_expiry`| string | No       | `YYYY-MM-DD`                         |
| `notes`          | string | No       | Updated notes                        |

---

## 6. Assign Asset to Employee

Assign an available asset to an employee.

```
POST /api/v1/assets/assign
```

### Request Body

```json
{
  "asset_id": 1,
  "employee_id": "EMP001",
  "assigned_date": "2026-02-06",
  "notes": "Assigned for project work"
}
```

| Field           | Type   | Required | Description                              |
|-----------------|--------|----------|------------------------------------------|
| `asset_id`      | int    | Yes      | Asset ID to assign                       |
| `employee_id`   | string | Yes      | Employee ID (e.g. `EMP001`)              |
| `assigned_date` | string | No       | `YYYY-MM-DD` (defaults to today)         |
| `notes`         | string | No       | Assignment notes                         |

### Business Rules
- Asset must be in `available` status
- Asset status changes to `in_use` after assignment

---

## 7. Reassign Asset (Asset Handler)

Transfer an asset from one employee to another.

```
POST /api/v1/assets/reassign
```

### Request Body

```json
{
  "asset_id": 1,
  "new_employee_id": "EMP002",
  "reassign_date": "2026-02-06",
  "return_date": "2026-02-06",
  "condition": "good",
  "notes": "Transferred to new team member"
}
```

| Field             | Type   | Required | Description                                   |
|-------------------|--------|----------|-----------------------------------------------|
| `asset_id`        | int    | Yes      | Asset ID to reassign                          |
| `new_employee_id` | string | Yes      | New employee ID                               |
| `reassign_date`   | string | No       | `YYYY-MM-DD` (defaults to today)              |
| `return_date`     | string | No       | Previous employee return date (defaults to reassign date) |
| `condition`       | string | No       | Updated condition on reassignment             |
| `notes`           | string | No       | Reassignment notes                            |

---

## 8. Return Asset

Return an assigned asset from an employee.

```
POST /api/v1/assets/return
```

### Request Body

```json
{
  "asset_id": 1,
  "return_date": "2026-02-06",
  "condition": "good",
  "notes": "Returned in good condition"
}
```

| Field         | Type   | Required | Description                            |
|---------------|--------|----------|----------------------------------------|
| `asset_id`    | int    | Yes      | Asset ID                               |
| `return_date` | string | No       | `YYYY-MM-DD` (defaults to today)       |
| `condition`   | string | No       | Condition on return                    |
| `notes`       | string | No       | Return notes                           |

### Business Rules
- Asset must be in `in_use` status with an assigned employee
- Asset status changes to `available` after return

---

## 9. Mark Asset for Repair

```
POST /api/v1/assets/mark-for-repair
```

### Request Body

```json
{
  "asset_id": 1,
  "repair_reason": "Screen flickering",
  "repair_notes": "Sent to vendor for screen replacement"
}
```

| Field           | Type   | Required | Description            |
|-----------------|--------|----------|------------------------|
| `asset_id`      | int    | Yes      | Asset ID               |
| `repair_reason` | string | No       | Reason for repair      |
| `repair_notes`  | string | No       | Additional repair notes|

### Business Rules
- Asset cannot be already `under_repair` or `retired`
- Asset status changes to `under_repair`

---

## 10. Complete Repair

```
POST /api/v1/assets/complete-repair
```

### Request Body

```json
{
  "asset_id": 1,
  "repair_cost": 150000,
  "condition": "good",
  "notes": "Screen replaced successfully"
}
```

| Field         | Type   | Required | Description                  |
|---------------|--------|----------|------------------------------|
| `asset_id`    | int    | Yes      | Asset ID                     |
| `repair_cost` | number | No       | Cost of repair               |
| `condition`   | string | No       | Condition after repair       |
| `notes`       | string | No       | Repair completion notes      |

### Business Rules
- Asset status changes to `available` after repair completion

---

## 11. Retire Asset

```
POST /api/v1/assets/retire
```

### Request Body

```json
{
  "asset_id": 1,
  "retirement_reason": "End of life - 5 years old",
  "retirement_date": "2026-02-06",
  "notes": "Replaced by new model"
}
```

| Field               | Type   | Required | Description                        |
|---------------------|--------|----------|------------------------------------|
| `asset_id`          | int    | Yes      | Asset ID                           |
| `retirement_reason` | string | No       | Reason for retirement              |
| `retirement_date`   | string | No       | `YYYY-MM-DD` (defaults to today)   |
| `notes`             | string | No       | Additional notes                   |

### Business Rules
- Asset cannot already be `retired`
- Asset status changes to `retired`

---

## 12. Delete Asset

Permanently delete a retired asset.

```
POST /api/v1/assets/delete
```

### Request Body

```json
{
  "id": 1
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id`  | int  | Yes      | Asset ID    |

### Business Rules
- Only `retired` assets can be deleted (soft delete)

---

# Part B — HR/Admin: Asset Request Management

Base URL: `/api/v1/assets/requests`

These endpoints let HR/Admin review, approve, reject, and fulfill asset requests submitted by employees.

---

## 13. List All Asset Requests (HR)

```
GET /api/v1/assets/requests
```

| Query Param    | Type   | Required | Default | Description                            |
|----------------|--------|----------|---------|----------------------------------------|
| `page`         | int    | No       | 1       | Page number                            |
| `page_size`    | int    | No       | 20      | Items per page (max 100)               |
| `employee_id`  | string | No       |         | Filter by employee ID                  |
| `request_type` | string | No       |         | `new`, `replacement`, `additional`     |
| `status`       | string | No       |         | `pending`, `approved`, `rejected`, `fulfilled`, `cancelled` |
| `asset_type`   | string | No       |         | Filter by asset type                   |
| `priority`     | string | No       |         | `low`, `medium`, `high`, `urgent`      |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset requests retrieved successfully",
  "data": [
    {
      "id": 1,
      "request_number": "AR-2026-001",
      "employee_id": "EMP001",
      "request_type": "new",
      "asset_type": "laptop",
      "status": "pending",
      "priority": "high",
      "justification": "Need a laptop for remote work",
      "requested_date": "2026-02-01T00:00:00Z",
      "brand": "Dell",
      "model": "Latitude 5520"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 5,
    "total_pages": 1
  }
}
```

---

## 14. Get Asset Request Details (HR)

```
GET /api/v1/assets/requests/:request_id
```

| Path Param   | Type | Required | Description |
|--------------|------|----------|-------------|
| `request_id` | int  | Yes      | Request ID  |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset request details retrieved successfully",
  "data": {
    "id": 1,
    "request_number": "AR-2026-001",
    "employee_id": "EMP001",
    "request_type": "new",
    "asset_type": "laptop",
    "status": "pending",
    "priority": "high",
    "justification": "Need a laptop for remote work",
    "requested_date": "2026-02-01T00:00:00Z",
    "brand": "Dell",
    "model": "Latitude 5520",
    "notes": "Preferred 16GB RAM"
  }
}
```

---

## 15. Approve Asset Request

```
POST /api/v1/assets/requests/:request_id/approve
```

| Path Param   | Type | Required | Description |
|--------------|------|----------|-------------|
| `request_id` | int  | Yes      | Request ID  |

No request body required.

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset request approved successfully",
  "data": {
    "id": 1,
    "request_number": "AR-2026-001",
    "status": "approved",
    "approved_at": "2026-02-06T10:00:00Z"
  }
}
```

### Business Rules
- Only `pending` requests can be approved
- Records the approving user and timestamp

---

## 16. Reject Asset Request

```
POST /api/v1/assets/requests/:request_id/reject
```

### Request Body

```json
{
  "rejection_reason": "Budget constraints - try again next quarter"
}
```

| Field              | Type   | Required | Description                        |
|--------------------|--------|----------|------------------------------------|
| `rejection_reason` | string | No       | Reason for rejection (defaults to "No reason provided") |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset request rejected successfully",
  "data": {
    "id": 1,
    "request_number": "AR-2026-001",
    "status": "rejected",
    "rejected_at": "2026-02-06T10:00:00Z",
    "rejection_reason": "Budget constraints - try again next quarter"
  }
}
```

---

## 17. Fulfill Asset Request

Assign a specific asset to fulfill an approved request.

```
POST /api/v1/assets/requests/:request_id/fulfill
```

### Request Body

```json
{
  "asset_id": 5
}
```

| Field      | Type | Required | Description                    |
|------------|------|----------|--------------------------------|
| `asset_id` | int  | Yes      | ID of the asset to assign      |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset request fulfilled successfully",
  "data": {
    "request": {
      "id": 1,
      "request_number": "AR-2026-001",
      "status": "fulfilled",
      "fulfilled_at": "2026-02-06T12:00:00Z",
      "assigned_asset_code": "AST-005"
    },
    "asset": {
      "id": 5,
      "asset_code": "AST-005",
      "asset_type": "laptop",
      "brand": "Dell",
      "model": "Latitude 5520",
      "employee_id": "EMP001",
      "status": "in_use",
      "assigned_date": "2026-02-06T00:00:00Z"
    }
  }
}
```

### Business Rules
- Request must be in `approved` status
- The specified asset must be `available`
- Asset gets assigned to the requesting employee

---

## 18. Reassign Asset (HR Handler)

Reassign an asset from one employee to another (via asset ID path param).

```
POST /api/v1/assets/:asset_id/reassign
```

### Request Body

```json
{
  "employee_id": "EMP002",
  "notes": "Reassigned due to department transfer"
}
```

| Field         | Type   | Required | Description                    |
|---------------|--------|----------|--------------------------------|
| `employee_id` | string | Yes      | New employee ID                |
| `notes`       | string | No       | Reassignment notes             |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset reassigned successfully",
  "data": {
    "id": 5,
    "asset_code": "AST-005",
    "asset_type": "laptop",
    "brand": "Dell",
    "model": "Latitude 5520",
    "employee_id": "EMP002",
    "assigned_to": "Jane Doe",
    "status": "in_use",
    "assigned_date": "2026-02-06T00:00:00Z"
  }
}
```

---

# Part C — Employee Self-Service

Base URL: `/api/v1/self-service/assets`

These endpoints let employees view their assigned assets, request new assets, report issues, and manage their own requests.

---

## 19. Get My Assigned Assets

Returns all assets currently assigned to the authenticated employee.

```
GET /api/v1/self-service/assets
```

No query parameters.

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Assigned assets retrieved successfully",
  "data": [
    {
      "id": 5,
      "asset_code": "AST-005",
      "asset_type": "laptop",
      "brand": "Dell",
      "model": "Latitude 5520",
      "serial_number": "SN789012",
      "status": "in_use",
      "condition": "good",
      "assigned_date": "2026-01-15T00:00:00Z",
      "purchase_date": "2025-06-01T00:00:00Z",
      "warranty_expiry": "2028-06-01T00:00:00Z",
      "value": 1500000
    },
    {
      "id": 12,
      "asset_code": "AST-012",
      "asset_type": "access_card",
      "brand": "HID",
      "model": "ProxCard II",
      "serial_number": "AC456789",
      "status": "in_use",
      "condition": "excellent",
      "assigned_date": "2025-12-01T00:00:00Z"
    }
  ]
}
```

---

## 20. Create Asset Request (Employee)

Request a new, replacement, or additional asset from HR.

```
POST /api/v1/self-service/assets/requests
```

### Request Body

```json
{
  "request_type": "new",
  "asset_type": "laptop",
  "justification": "Need a laptop for remote work assignments",
  "brand": "Dell",
  "model": "Latitude 5530",
  "priority": "high",
  "notes": "Prefer 16GB RAM, 512GB SSD"
}
```

| Field               | Type   | Required | Description                                        |
|---------------------|--------|----------|----------------------------------------------------|
| `request_type`      | string | Yes      | `new`, `replacement`, `additional`                 |
| `asset_type`        | string | Yes      | `laptop`, `mobile_phone`, `desktop`, etc.          |
| `justification`     | string | Yes      | Why the asset is needed                            |
| `brand`             | string | No       | Preferred brand                                    |
| `model`             | string | No       | Preferred model                                    |
| `priority`          | string | No       | `low`, `medium` (default), `high`, `urgent`        |
| `replacing_asset_id`| int    | No       | Asset ID being replaced (for `replacement` type)   |
| `notes`             | string | No       | Additional notes/specifications                    |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset request created successfully",
  "data": {
    "id": 3,
    "request_number": "AR-2026-003",
    "request_type": "new",
    "asset_type": "laptop",
    "status": "pending",
    "priority": "high",
    "requested_date": "2026-02-06T00:00:00Z",
    "brand": "Dell",
    "model": "Latitude 5530"
  }
}
```

---

## 21. List My Asset Requests

Retrieve the employee's own asset requests with optional filters.

```
GET /api/v1/self-service/assets/requests
```

| Query Param    | Type   | Required | Default | Description                            |
|----------------|--------|----------|---------|----------------------------------------|
| `page`         | int    | No       | 1       | Page number                            |
| `page_size`    | int    | No       | 20      | Items per page (max 100)               |
| `request_type` | string | No       |         | `new`, `replacement`, `additional`     |
| `status`       | string | No       |         | `pending`, `approved`, `rejected`, `fulfilled`, `cancelled` |
| `asset_type`   | string | No       |         | Filter by asset type                   |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset requests retrieved successfully",
  "data": [
    {
      "id": 3,
      "request_number": "AR-2026-003",
      "request_type": "new",
      "asset_type": "laptop",
      "status": "pending",
      "priority": "high",
      "justification": "Need a laptop for remote work assignments",
      "requested_date": "2026-02-06T00:00:00Z",
      "brand": "Dell",
      "model": "Latitude 5530"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 3,
    "total_pages": 1
  }
}
```

---

## 22. Get My Asset Request Details

```
GET /api/v1/self-service/assets/requests/:request_id
```

| Path Param   | Type | Required | Description |
|--------------|------|----------|-------------|
| `request_id` | int  | Yes      | Request ID  |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset request details retrieved successfully",
  "data": {
    "id": 3,
    "request_number": "AR-2026-003",
    "request_type": "new",
    "asset_type": "laptop",
    "status": "approved",
    "priority": "high",
    "justification": "Need a laptop for remote work assignments",
    "requested_date": "2026-02-06T00:00:00Z",
    "brand": "Dell",
    "model": "Latitude 5530",
    "notes": "Prefer 16GB RAM",
    "approved_at": "2026-02-07T09:00:00Z"
  }
}
```

### Business Rules
- Employee can only view their own requests

---

## 23. Cancel Asset Request

Cancel a pending asset request.

```
POST /api/v1/self-service/assets/requests/:request_id/cancel
```

### Request Body (optional)

```json
{
  "reason": "No longer needed"
}
```

| Field    | Type   | Required | Description                                        |
|----------|--------|----------|----------------------------------------------------|
| `reason` | string | No       | Cancellation reason (defaults to "No reason provided") |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset request cancelled successfully",
  "data": null
}
```

### Business Rules
- Only `pending` requests can be cancelled by the employee

---

## 24. Report Asset Issue

Report a problem with an assigned asset (malfunction, damage, lost, stolen).

```
POST /api/v1/self-service/assets/issues
```

### Request Body

```json
{
  "asset_id": 5,
  "issue_type": "malfunction",
  "title": "Laptop screen flickering",
  "description": "The screen flickers intermittently, especially when plugged into external monitor",
  "priority": "high",
  "incident_date": "2026-02-05",
  "incident_location": "Office - Floor 3",
  "police_report_number": null,
  "attachment_urls": [
    "https://storage.example.com/uploads/issue-photo-1.jpg"
  ]
}
```

| Field                  | Type     | Required | Description                                          |
|------------------------|----------|----------|------------------------------------------------------|
| `asset_id`             | int      | Yes      | ID of the asset with the issue                       |
| `issue_type`           | string   | Yes      | `malfunction`, `damage`, `lost`, `stolen`, `other`   |
| `title`                | string   | Yes      | Brief title of the issue                             |
| `description`          | string   | Yes      | Detailed description                                 |
| `priority`             | string   | No       | `low`, `medium` (default), `high`, `urgent`          |
| `incident_date`        | string   | No       | `YYYY-MM-DD` — when the incident occurred            |
| `incident_location`    | string   | No       | Where the incident happened                          |
| `police_report_number` | string   | No       | Police report number (for `lost`/`stolen`)           |
| `attachment_urls`      | string[] | No       | Array of attachment/photo URLs                       |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset issue reported successfully",
  "data": {
    "id": 7,
    "issue_number": "AI-2026-007",
    "asset_id": 5,
    "asset_code": "AST-005",
    "issue_type": "malfunction",
    "status": "reported",
    "priority": "high",
    "title": "Laptop screen flickering",
    "description": "The screen flickers intermittently...",
    "reported_date": "2026-02-06T00:00:00Z",
    "incident_date": "2026-02-05T00:00:00Z",
    "incident_location": "Office - Floor 3"
  }
}
```

---

## 25. List My Asset Issues

Retrieve the employee's own reported asset issues.

```
GET /api/v1/self-service/assets/issues
```

| Query Param  | Type   | Required | Default | Description                                     |
|--------------|--------|----------|---------|-------------------------------------------------|
| `page`       | int    | No       | 1       | Page number                                     |
| `page_size`  | int    | No       | 20      | Items per page (max 100)                        |
| `issue_type` | string | No       |         | `malfunction`, `damage`, `lost`, `stolen`, `other` |
| `status`     | string | No       |         | `reported`, `under_review`, `in_progress`, `resolved`, `closed` |
| `asset_id`   | int    | No       |         | Filter by specific asset                        |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset issues retrieved successfully",
  "data": [
    {
      "id": 7,
      "issue_number": "AI-2026-007",
      "asset_id": 5,
      "asset_code": "AST-005",
      "issue_type": "malfunction",
      "status": "in_progress",
      "priority": "high",
      "title": "Laptop screen flickering",
      "description": "The screen flickers intermittently...",
      "reported_date": "2026-02-06T00:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

---

## 26. Get My Asset Issue Details

```
GET /api/v1/self-service/assets/issues/:issue_id
```

| Path Param | Type | Required | Description |
|------------|------|----------|-------------|
| `issue_id` | int  | Yes      | Issue ID    |

### Response — `200 OK`

```json
{
  "success": true,
  "message": "Asset issue details retrieved successfully",
  "data": {
    "id": 7,
    "issue_number": "AI-2026-007",
    "asset_id": 5,
    "asset_code": "AST-005",
    "issue_type": "malfunction",
    "status": "resolved",
    "priority": "high",
    "title": "Laptop screen flickering",
    "description": "The screen flickers intermittently...",
    "reported_date": "2026-02-06T00:00:00Z",
    "attachment_urls": ["https://storage.example.com/uploads/issue-photo-1.jpg"],
    "resolved_at": "2026-02-08T14:00:00Z",
    "resolution_notes": "Screen cable was loose, reseated and tested",
    "action_taken": "Hardware repair - screen cable"
  }
}
```

---

# Data Models

## Asset Status Lifecycle

```
available → in_use → available (returned)
available → under_repair → available (repair completed)
any (except retired) → retired → deleted (soft delete)
```

| Status         | Description                        |
|----------------|------------------------------------|
| `available`    | Ready for assignment               |
| `in_use`       | Currently assigned to an employee  |
| `under_repair` | Sent for repair                    |
| `retired`      | End of life, decommissioned        |

## Asset Types

| Code               | Name              |
|--------------------|-------------------|
| `laptop`           | Laptop            |
| `mobile_phone`     | Mobile Phone      |
| `desktop`          | Desktop           |
| `tablet`           | Tablet            |
| `id_card`          | ID Card           |
| `access_card`      | Access Card       |
| `vehicle`          | Vehicle           |
| `office_equipment` | Office Equipment  |

## Asset Condition

| Value       | Description          |
|-------------|----------------------|
| `excellent` | Like new             |
| `good`      | Normal wear          |
| `fair`      | Noticeable wear      |
| `poor`      | Significant wear     |

## Asset Request Types

| Type          | Description                       |
|---------------|-----------------------------------|
| `new`         | Request for a brand new asset     |
| `replacement` | Replace an existing faulty asset  |
| `additional`  | Request an additional asset       |

## Asset Request Status Flow

```
pending → approved → fulfilled
pending → rejected
pending → cancelled (by employee)
```

| Status      | Description                         |
|-------------|-------------------------------------|
| `pending`   | Awaiting HR review                  |
| `approved`  | Approved by HR, awaiting fulfillment|
| `rejected`  | Rejected by HR                      |
| `fulfilled` | Asset has been assigned             |
| `cancelled` | Cancelled by the employee           |

## Asset Issue Types

| Type           | Description                 |
|----------------|-----------------------------|
| `malfunction`  | Asset not working properly  |
| `damage`       | Asset is damaged            |
| `lost`         | Asset is lost               |
| `stolen`       | Asset is stolen             |
| `other`        | Other issues                |

## Asset Issue Status Flow

```
reported → under_review → in_progress → resolved → closed
```

| Status         | Description             |
|----------------|-------------------------|
| `reported`     | Newly reported          |
| `under_review` | HR/IT reviewing         |
| `in_progress`  | Being resolved          |
| `resolved`     | Issue fixed             |
| `closed`       | Issue closed            |

## Priority Levels

| Priority | Description              |
|----------|--------------------------|
| `low`    | Not time-sensitive       |
| `medium` | Normal priority (default)|
| `high`   | Needs attention soon     |
| `urgent` | Needs immediate action   |

---

# Quick Reference Table

## Admin/HR Asset Management (`/api/v1/assets`)

| Method | Endpoint                                  | Description                          |
|--------|-------------------------------------------|--------------------------------------|
| GET    | `/api/v1/assets`                          | List all assets (paginated, filtered)|
| GET    | `/api/v1/assets/get?id=:id`               | Get asset by ID                      |
| GET    | `/api/v1/assets/types`                    | Get all asset types                  |
| POST   | `/api/v1/assets`                          | Create new asset                     |
| POST   | `/api/v1/assets/update`                   | Update asset details                 |
| POST   | `/api/v1/assets/assign`                   | Assign asset to employee             |
| POST   | `/api/v1/assets/reassign`                 | Reassign asset to another employee   |
| POST   | `/api/v1/assets/return`                   | Return asset from employee           |
| POST   | `/api/v1/assets/mark-for-repair`          | Mark asset for repair                |
| POST   | `/api/v1/assets/complete-repair`          | Complete asset repair                |
| POST   | `/api/v1/assets/retire`                   | Retire asset                         |
| POST   | `/api/v1/assets/delete`                   | Delete retired asset                 |

## HR/Admin Request Management (`/api/v1/assets/requests`)

| Method | Endpoint                                           | Description                  |
|--------|-----------------------------------------------------|------------------------------|
| GET    | `/api/v1/assets/requests`                           | List all asset requests      |
| GET    | `/api/v1/assets/requests/:request_id`               | Get request details          |
| POST   | `/api/v1/assets/requests/:request_id/approve`       | Approve request              |
| POST   | `/api/v1/assets/requests/:request_id/reject`        | Reject request               |
| POST   | `/api/v1/assets/requests/:request_id/fulfill`       | Fulfill request (assign asset)|
| POST   | `/api/v1/assets/:asset_id/reassign`                 | Reassign asset (by asset ID) |

## Employee Self-Service (`/api/v1/self-service/assets`)

| Method | Endpoint                                                    | Description                    |
|--------|--------------------------------------------------------------|--------------------------------|
| GET    | `/api/v1/self-service/assets`                               | Get my assigned assets         |
| POST   | `/api/v1/self-service/assets/requests`                      | Create asset request           |
| GET    | `/api/v1/self-service/assets/requests`                      | List my asset requests         |
| GET    | `/api/v1/self-service/assets/requests/:request_id`          | Get my request details         |
| POST   | `/api/v1/self-service/assets/requests/:request_id/cancel`   | Cancel my pending request      |
| POST   | `/api/v1/self-service/assets/issues`                        | Report an asset issue          |
| GET    | `/api/v1/self-service/assets/issues`                        | List my asset issues           |
| GET    | `/api/v1/self-service/assets/issues/:issue_id`              | Get my issue details           |

---

## Error Responses

All endpoints return errors in this format:

```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error (if validation)"
}
```

| Status | Meaning                  |
|--------|--------------------------|
| 400    | Bad request / validation |
| 401    | Not authenticated        |
| 403    | Insufficient permissions |
| 404    | Resource not found       |
| 422    | Validation failed        |
| 500    | Internal server error    |

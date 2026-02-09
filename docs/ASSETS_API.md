# Assets Management API

Complete API documentation for the Assets module covering employee self-service (requesting assets, reporting issues) and HR/Admin asset management (inventory, approvals, assignments, repairs, retirement).

**Authentication:** All endpoints require a valid JWT token via `Authorization: Bearer <token>` header.

---

## Table of Contents

### Employee Self-Service
1. [Get My Assigned Assets](#1-get-my-assigned-assets)
2. [Request an Asset](#2-request-an-asset)
3. [List My Asset Requests](#3-list-my-asset-requests)
4. [Get My Request Details](#4-get-my-request-details)
5. [Cancel My Asset Request](#5-cancel-my-asset-request)
6. [Report an Asset Issue](#6-report-an-asset-issue)
7. [List My Asset Issues](#7-list-my-asset-issues)
8. [Get My Issue Details](#8-get-my-issue-details)

### HR/Admin — Asset Request Management
9. [List All Asset Requests](#9-list-all-asset-requests-hradmin)
10. [Get Asset Request Details](#10-get-asset-request-details-hradmin)
11. [Approve Asset Request](#11-approve-asset-request)
12. [Reject Asset Request](#12-reject-asset-request)
13. [Fulfill Asset Request](#13-fulfill-asset-request)

### HR/Admin — Asset Inventory Management
14. [List Assets](#14-list-assets)
15. [Get Asset by ID](#15-get-asset-by-id)
16. [Get Asset Types](#16-get-asset-types)
17. [Create Asset](#17-create-asset)
18. [Update Asset](#18-update-asset)
19. [Assign Asset to Employee](#19-assign-asset-to-employee)
20. [Reassign Asset](#20-reassign-asset)
21. [Return Asset](#21-return-asset)
22. [Mark Asset for Repair](#22-mark-asset-for-repair)
23. [Complete Asset Repair](#23-complete-asset-repair)
24. [Retire Asset](#24-retire-asset)
25. [Delete Asset](#25-delete-asset)

### Reference
26. [Enums & Status Values](#26-enums--status-values)
27. [Workflows](#27-workflows)
28. [Data Models](#28-data-models)

---

## Employee Self-Service

> **Base URL:** `/api/v1/self-service/assets`
> **Auth:** Any authenticated employee

---

### 1. Get My Assigned Assets

Retrieve all assets currently assigned to the authenticated employee.

**Endpoint:** `GET /api/v1/self-service/assets`

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Assigned assets retrieved successfully",
  "data": [
    {
      "id": 1,
      "asset_code": "AST-2026-001",
      "asset_type": "laptop",
      "brand": "Dell",
      "model": "Latitude 5540",
      "serial_number": "SN12345678",
      "status": "in_use",
      "condition": "good",
      "assigned_date": "2025-06-15T00:00:00Z",
      "purchase_date": "2025-01-10T00:00:00Z",
      "warranty_expiry": "2028-01-10T00:00:00Z",
      "value": 1500.00,
      "notes": "Issued with charger and laptop bag"
    }
  ]
}
```

---

### 2. Request an Asset

Submit a request for a new, replacement, or additional asset. The request goes to HR for approval.

**Endpoint:** `POST /api/v1/self-service/assets/requests`

#### Request Body

| Field                | Type    | Required | Description                                                    |
|----------------------|---------|----------|----------------------------------------------------------------|
| `request_type`       | string  | Yes      | `new`, `replacement`, or `additional`                          |
| `asset_type`         | string  | Yes      | Asset type needed (see [Asset Types](#asset-types))            |
| `justification`      | string  | Yes      | Reason why the asset is needed                                 |
| `brand`              | string  | No       | Preferred brand                                                |
| `model`              | string  | No       | Preferred model                                                |
| `priority`           | string  | No       | `low`, `medium` (default), `high`, `urgent`                    |
| `replacing_asset_id` | integer | No       | ID of asset being replaced (required if `request_type` = `replacement`) |
| `notes`              | string  | No       | Additional notes                                               |

#### Example Request

```json
{
  "request_type": "new",
  "asset_type": "laptop",
  "justification": "New hire starting next week, needs a development laptop",
  "brand": "Dell",
  "model": "Latitude 5540",
  "priority": "high"
}
```

#### Example: Replacement Request

```json
{
  "request_type": "replacement",
  "asset_type": "laptop",
  "justification": "Current laptop is 4 years old and frequently crashes",
  "replacing_asset_id": 5,
  "priority": "medium"
}
```

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset request created successfully",
  "data": {
    "id": 10,
    "request_number": "AR-2026-010",
    "request_type": "new",
    "asset_type": "laptop",
    "status": "pending",
    "priority": "high",
    "requested_date": "2026-02-06T10:30:00Z",
    "brand": "Dell",
    "model": "Latitude 5540"
  }
}
```

---

### 3. List My Asset Requests

Retrieve all asset requests submitted by the authenticated employee.

**Endpoint:** `GET /api/v1/self-service/assets/requests`

#### Query Parameters

| Parameter      | Type   | Default | Description                                |
|----------------|--------|---------|--------------------------------------------|
| `request_type` | string | —       | Filter: `new`, `replacement`, `additional` |
| `status`       | string | —       | Filter: `pending`, `approved`, `rejected`, `fulfilled`, `cancelled` |
| `asset_type`   | string | —       | Filter: `laptop`, `mobile_phone`, etc.     |
| `page`         | int    | 1       | Page number                                |
| `page_size`    | int    | 20      | Items per page (max 100)                   |

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset requests retrieved successfully",
  "data": [
    {
      "id": 10,
      "request_number": "AR-2026-010",
      "request_type": "new",
      "asset_type": "laptop",
      "status": "approved",
      "priority": "high",
      "justification": "New hire starting next week",
      "requested_date": "2026-02-06T10:30:00Z",
      "brand": "Dell",
      "model": "Latitude 5540",
      "approved_at": "2026-02-06T14:00:00Z"
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

### 4. Get My Request Details

Get detailed information about a specific asset request.

**Endpoint:** `GET /api/v1/self-service/assets/requests/:request_id`

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset request details retrieved successfully",
  "data": {
    "id": 10,
    "request_number": "AR-2026-010",
    "request_type": "replacement",
    "asset_type": "laptop",
    "status": "fulfilled",
    "priority": "medium",
    "justification": "Current laptop crashes frequently",
    "requested_date": "2026-02-01T08:00:00Z",
    "brand": "Dell",
    "notes": "Need at least 16GB RAM",
    "replacing_asset_code": "AST-2024-005",
    "approved_at": "2026-02-02T09:00:00Z",
    "fulfilled_at": "2026-02-04T11:00:00Z",
    "assigned_asset_code": "AST-2026-015"
  }
}
```

---

### 5. Cancel My Asset Request

Cancel a pending asset request. Only requests with `pending` status can be cancelled.

**Endpoint:** `POST /api/v1/self-service/assets/requests/:request_id/cancel`

#### Request Body (optional)

```json
{
  "reason": "No longer needed"
}
```

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset request cancelled successfully"
}
```

#### Error (400 Bad Request)

```json
{
  "success": false,
  "message": "cannot cancel a request that is not pending"
}
```

---

### 6. Report an Asset Issue

Report a problem with an assigned asset (malfunction, damage, lost, stolen, etc.).

**Endpoint:** `POST /api/v1/self-service/assets/issues`

#### Request Body

| Field                  | Type     | Required | Description                                               |
|------------------------|----------|----------|-----------------------------------------------------------|
| `asset_id`             | integer  | Yes      | ID of the asset with the issue                            |
| `issue_type`           | string   | Yes      | `malfunction`, `damage`, `lost`, `stolen`, `other`        |
| `title`                | string   | Yes      | Brief title of the issue                                  |
| `description`          | string   | Yes      | Detailed description of the problem                       |
| `priority`             | string   | No       | `low`, `medium` (default), `high`, `urgent`               |
| `incident_date`        | string   | No       | When the incident occurred (YYYY-MM-DD)                   |
| `incident_location`    | string   | No       | Where it happened (for lost/stolen)                       |
| `police_report_number` | string   | No       | Police report number (for stolen assets)                  |
| `attachment_urls`      | string[] | No       | Array of file URLs (photos, documents)                    |

#### Example: Malfunction Report

```json
{
  "asset_id": 1,
  "issue_type": "malfunction",
  "title": "Laptop screen flickering",
  "description": "The screen flickers intermittently, especially when running multiple applications. Started happening 3 days ago.",
  "priority": "high"
}
```

#### Example: Lost/Stolen Report

```json
{
  "asset_id": 3,
  "issue_type": "stolen",
  "title": "Mobile phone stolen",
  "description": "Company phone was stolen from my bag at a conference venue.",
  "priority": "urgent",
  "incident_date": "2026-02-05",
  "incident_location": "Mlimani City Conference Hall, Dar es Salaam",
  "police_report_number": "PR-2026-00456",
  "attachment_urls": ["https://storage.example.com/police-report-scan.pdf"]
}
```

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset issue reported successfully",
  "data": {
    "id": 5,
    "issue_number": "AI-2026-005",
    "asset_id": 3,
    "asset_code": "AST-2025-003",
    "issue_type": "stolen",
    "status": "reported",
    "priority": "urgent",
    "title": "Mobile phone stolen",
    "description": "Company phone was stolen from my bag...",
    "reported_date": "2026-02-06T15:30:00Z",
    "incident_date": "2026-02-05T00:00:00Z",
    "incident_location": "Mlimani City Conference Hall, Dar es Salaam"
  }
}
```

---

### 7. List My Asset Issues

Retrieve all asset issues reported by the authenticated employee.

**Endpoint:** `GET /api/v1/self-service/assets/issues`

#### Query Parameters

| Parameter    | Type   | Default | Description                                               |
|--------------|--------|---------|------------------------------------------------------------|
| `issue_type` | string | —       | Filter: `malfunction`, `damage`, `lost`, `stolen`, `other` |
| `status`     | string | —       | Filter: `reported`, `under_review`, `in_progress`, `resolved`, `closed` |
| `asset_id`   | int    | —       | Filter by specific asset                                    |
| `page`       | int    | 1       | Page number                                                 |
| `page_size`  | int    | 20      | Items per page (max 100)                                    |

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset issues retrieved successfully",
  "data": [
    {
      "id": 5,
      "issue_number": "AI-2026-005",
      "asset_id": 3,
      "asset_code": "AST-2025-003",
      "issue_type": "stolen",
      "status": "under_review",
      "priority": "urgent",
      "title": "Mobile phone stolen",
      "description": "Company phone was stolen...",
      "reported_date": "2026-02-06T15:30:00Z",
      "incident_date": "2026-02-05T00:00:00Z",
      "incident_location": "Mlimani City Conference Hall"
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

### 8. Get My Issue Details

Get detailed information about a specific asset issue, including resolution details and attachments.

**Endpoint:** `GET /api/v1/self-service/assets/issues/:issue_id`

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset issue details retrieved successfully",
  "data": {
    "id": 5,
    "issue_number": "AI-2026-005",
    "asset_id": 3,
    "asset_code": "AST-2025-003",
    "issue_type": "malfunction",
    "status": "resolved",
    "priority": "high",
    "title": "Laptop screen flickering",
    "description": "The screen flickers intermittently...",
    "reported_date": "2026-02-01T10:00:00Z",
    "resolved_at": "2026-02-04T16:00:00Z",
    "resolution_notes": "Replaced the display cable",
    "action_taken": "Sent to vendor for screen cable replacement under warranty",
    "attachment_urls": ["https://storage.example.com/screen-photo.jpg"]
  }
}
```

---

## HR/Admin — Asset Request Management

> **Base URL:** `/api/v1/assets/requests`
> **Auth:** HR or Admin role required

---

### 9. List All Asset Requests (HR/Admin)

View all asset requests from all employees with filters.

**Endpoint:** `GET /api/v1/assets/requests`

#### Query Parameters

| Parameter      | Type   | Default | Description                                 |
|----------------|--------|---------|---------------------------------------------|
| `employee_id`  | string | —       | Filter by employee ID (e.g., "EMP-001")     |
| `request_type` | string | —       | Filter: `new`, `replacement`, `additional`  |
| `status`       | string | —       | Filter: `pending`, `approved`, `rejected`, `fulfilled`, `cancelled` |
| `asset_type`   | string | —       | Filter: `laptop`, `mobile_phone`, etc.      |
| `priority`     | string | —       | Filter: `low`, `medium`, `high`, `urgent`   |
| `page`         | int    | 1       | Page number                                 |
| `page_size`    | int    | 20      | Items per page (max 100)                    |

#### Example

```bash
# Get all pending requests
GET /api/v1/assets/requests?status=pending

# Get pending laptop requests with high priority
GET /api/v1/assets/requests?status=pending&asset_type=laptop&priority=high
```

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset requests retrieved successfully",
  "data": [
    {
      "id": 10,
      "request_number": "AR-2026-010",
      "employee_id": "EMP-001",
      "request_type": "new",
      "asset_type": "laptop",
      "status": "pending",
      "priority": "high",
      "justification": "New hire starting next week",
      "requested_date": "2026-02-06T10:30:00Z",
      "brand": "Dell",
      "model": "Latitude 5540",
      "notes": "Need at least 16GB RAM"
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

### 10. Get Asset Request Details (HR/Admin)

**Endpoint:** `GET /api/v1/assets/requests/:request_id`

Returns the same structure as the employee view but also includes `approved_by` (user ID) field.

---

### 11. Approve Asset Request

Approve a pending asset request. After approval, the request can be fulfilled by assigning an asset.

**Endpoint:** `POST /api/v1/assets/requests/:request_id/approve`

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset request approved successfully",
  "data": {
    "id": 10,
    "request_number": "AR-2026-010",
    "status": "approved",
    "approved_at": "2026-02-06T14:00:00Z"
  }
}
```

---

### 12. Reject Asset Request

Reject a pending asset request with a reason.

**Endpoint:** `POST /api/v1/assets/requests/:request_id/reject`

#### Request Body

| Field              | Type   | Required | Description     |
|--------------------|--------|----------|-----------------|
| `rejection_reason` | string | No       | Reason for rejection (defaults to "No reason provided") |

```json
{
  "rejection_reason": "Budget constraints - request again next quarter"
}
```

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset request rejected successfully",
  "data": {
    "id": 10,
    "request_number": "AR-2026-010",
    "status": "rejected",
    "rejected_at": "2026-02-06T14:30:00Z",
    "rejection_reason": "Budget constraints - request again next quarter"
  }
}
```

---

### 13. Fulfill Asset Request

Assign a specific asset to fulfill an **approved** request. This assigns the asset to the employee and marks the request as fulfilled.

**Endpoint:** `POST /api/v1/assets/requests/:request_id/fulfill`

#### Request Body

| Field      | Type    | Required | Description                              |
|------------|---------|----------|------------------------------------------|
| `asset_id` | integer | Yes      | ID of the available asset to assign      |

```json
{
  "asset_id": 15
}
```

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset request fulfilled successfully",
  "data": {
    "request": {
      "id": 10,
      "request_number": "AR-2026-010",
      "status": "fulfilled",
      "fulfilled_at": "2026-02-07T09:00:00Z",
      "assigned_asset_code": "AST-2026-015"
    },
    "asset": {
      "id": 15,
      "asset_code": "AST-2026-015",
      "asset_type": "laptop",
      "brand": "Dell",
      "model": "Latitude 5540",
      "employee_id": "EMP-001",
      "status": "in_use",
      "assigned_date": "2026-02-07T09:00:00Z"
    }
  }
}
```

---

## HR/Admin — Asset Inventory Management

> **Base URL:** `/api/v1/assets`
> **Auth:** HR or Admin role required

---

### 14. List Assets

Retrieve a paginated list of all assets in the inventory with filters.

**Endpoint:** `GET /api/v1/assets`

#### Query Parameters

| Parameter     | Type   | Default | Description                                    |
|---------------|--------|---------|------------------------------------------------|
| `search`      | string | —       | Search by asset code, brand, model, serial no. |
| `status`      | string | —       | Filter: `available`, `in_use`, `under_repair`, `retired` |
| `asset_type`  | string | —       | Filter: `laptop`, `mobile_phone`, etc.         |
| `assigned_to` | string | —       | Filter by employee ID                          |
| `department`  | string | —       | Filter by department name                      |
| `page`        | int    | 1       | Page number                                    |
| `page_size`   | int    | 20      | Items per page (max 100)                       |

#### Example

```bash
# Get all available laptops
GET /api/v1/assets?status=available&asset_type=laptop

# Search by serial number
GET /api/v1/assets?search=SN12345
```

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Assets retrieved successfully",
  "data": [
    {
      "id": 1,
      "asset_code": "AST-2026-001",
      "asset_type": "laptop",
      "brand": "Dell",
      "model": "Latitude 5540",
      "serial_number": "SN12345678",
      "assigned_to": "John Doe",
      "employee_id": "EMP-001",
      "employee_photo": "/photos/john.jpg",
      "department": "Engineering",
      "assigned_date": "2025-06-15T00:00:00Z",
      "status": "in_use",
      "condition": "good",
      "purchase_date": "2025-01-10T00:00:00Z",
      "warranty_expiry": "2028-01-10T00:00:00Z",
      "value": 1500.00,
      "created_at": "2025-01-12T08:00:00Z",
      "updated_at": "2025-06-15T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 50,
    "total_pages": 3
  }
}
```

---

### 15. Get Asset by ID

**Endpoint:** `GET /api/v1/assets/get?id={id}`

Returns a single asset with all details including repair and retirement info.

---

### 16. Get Asset Types

Returns all available asset types.

**Endpoint:** `GET /api/v1/assets/types`

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Asset types retrieved successfully",
  "data": [
    { "id": 1, "name": "Laptop", "code": "laptop", "description": "Portable computer" },
    { "id": 2, "name": "Mobile Phone", "code": "mobile_phone", "description": "Company mobile phone" },
    { "id": 3, "name": "Desktop", "code": "desktop", "description": "Desktop computer" },
    { "id": 4, "name": "Tablet", "code": "tablet", "description": "Tablet device" },
    { "id": 5, "name": "ID Card", "code": "id_card", "description": "Employee identification card" },
    { "id": 6, "name": "Access Card", "code": "access_card", "description": "Building access card" },
    { "id": 7, "name": "Vehicle", "code": "vehicle", "description": "Company vehicle" },
    { "id": 8, "name": "Office Equipment", "code": "office_equipment", "description": "Office furniture and equipment" }
  ]
}
```

---

### 17. Create Asset

Add a new asset to the inventory. Optionally assign it immediately to an employee.

**Endpoint:** `POST /api/v1/assets`

#### Request Body

| Field                     | Type          | Required | Description                                   |
|---------------------------|---------------|----------|-----------------------------------------------|
| `asset_type`              | string        | Yes      | Asset type (see [Asset Types](#asset-types))   |
| `brand`                   | string        | Yes      | Brand/manufacturer                            |
| `model`                   | string        | Yes      | Model name/number                             |
| `serial_number`           | string        | Yes      | Unique serial number                          |
| `asset_code`              | string        | No       | Custom code (auto-generated if omitted)       |
| `purchase_date`           | string        | No       | Purchase date (YYYY-MM-DD)                    |
| `value`                   | number        | No       | Purchase value                                |
| `warranty_expiry`         | string        | No       | Warranty expiry date (YYYY-MM-DD)             |
| `condition`               | string        | No       | `excellent` (default), `good`, `fair`, `poor` |
| `assigned_to_employee_id` | number/string | No       | Employee ID to assign immediately             |
| `notes`                   | string        | No       | Notes                                         |

```json
{
  "asset_type": "laptop",
  "brand": "Dell",
  "model": "Latitude 5540",
  "serial_number": "SN-2026-ABC123",
  "purchase_date": "2026-02-01",
  "value": 1500.00,
  "warranty_expiry": "2029-02-01",
  "condition": "excellent",
  "notes": "Includes charger and docking station"
}
```

---

### 18. Update Asset

Update asset details (brand, model, condition, value, etc.).

**Endpoint:** `POST /api/v1/assets/update`

#### Request Body

| Field            | Type   | Required | Description            |
|------------------|--------|----------|------------------------|
| `id`             | int    | Yes      | Asset ID               |
| `asset_type`     | string | No       | Updated asset type     |
| `brand`          | string | No       | Updated brand          |
| `model`          | string | No       | Updated model          |
| `serial_number`  | string | No       | Updated serial number  |
| `condition`      | string | No       | Updated condition      |
| `value`          | number | No       | Updated value          |
| `purchase_date`  | string | No       | Updated purchase date  |
| `warranty_expiry`| string | No       | Updated warranty date  |
| `notes`          | string | No       | Updated notes          |

---

### 19. Assign Asset to Employee

Assign an available asset to an employee.

**Endpoint:** `POST /api/v1/assets/assign`

#### Request Body

| Field           | Type   | Required | Description                           |
|-----------------|--------|----------|---------------------------------------|
| `asset_id`      | int    | Yes      | ID of the asset to assign             |
| `employee_id`   | string | Yes      | Employee ID (e.g., "EMP-001")         |
| `assigned_date` | string | No       | Assignment date (YYYY-MM-DD, default: today) |
| `notes`         | string | No       | Assignment notes                      |

```json
{
  "asset_id": 15,
  "employee_id": "EMP-001",
  "notes": "Issued with charger and laptop bag"
}
```

---

### 20. Reassign Asset

Transfer an asset from one employee to another.

**Endpoint:** `POST /api/v1/assets/:asset_id/reassign`

#### Request Body

| Field         | Type   | Required | Description                        |
|---------------|--------|----------|------------------------------------|
| `employee_id` | string | Yes      | New employee ID to reassign to     |
| `notes`       | string | No       | Reassignment notes                 |

```json
{
  "employee_id": "EMP-005",
  "notes": "Employee EMP-001 transferred to another department"
}
```

There is also a bulk reassign endpoint:

**Endpoint:** `POST /api/v1/assets/reassign`

| Field            | Type   | Required | Description                              |
|------------------|--------|----------|------------------------------------------|
| `asset_id`       | int    | Yes      | Asset ID                                 |
| `new_employee_id`| string | Yes      | New employee ID                          |
| `reassign_date`  | string | No       | Reassignment date (YYYY-MM-DD)           |
| `return_date`    | string | No       | Previous employee return date            |
| `condition`      | string | No       | Update condition on reassignment         |
| `notes`          | string | No       | Reassignment notes                       |

---

### 21. Return Asset

Return an asset from an employee (makes it available again).

**Endpoint:** `POST /api/v1/assets/return`

#### Request Body

| Field         | Type   | Required | Description                              |
|---------------|--------|----------|------------------------------------------|
| `asset_id`    | int    | Yes      | Asset ID                                 |
| `return_date` | string | No       | Return date (YYYY-MM-DD, default: today) |
| `condition`   | string | No       | Condition on return                      |
| `notes`       | string | No       | Return notes                             |

```json
{
  "asset_id": 15,
  "condition": "good",
  "notes": "Returned in good condition with all accessories"
}
```

---

### 22. Mark Asset for Repair

Send an asset for repair. Changes status to `under_repair`.

**Endpoint:** `POST /api/v1/assets/mark-for-repair`

#### Request Body

| Field           | Type   | Required | Description        |
|-----------------|--------|----------|--------------------|
| `asset_id`      | int    | Yes      | Asset ID           |
| `repair_reason` | string | No       | Reason for repair  |
| `repair_notes`  | string | No       | Additional notes   |

```json
{
  "asset_id": 1,
  "repair_reason": "Screen flickering - display cable issue",
  "repair_notes": "Sent to Dell authorized service center"
}
```

---

### 23. Complete Asset Repair

Mark a repair as complete. Changes status back to `available`.

**Endpoint:** `POST /api/v1/assets/complete-repair`

#### Request Body

| Field        | Type   | Required | Description               |
|--------------|--------|----------|---------------------------|
| `asset_id`   | int    | Yes      | Asset ID                  |
| `repair_cost`| number | No       | Repair cost               |
| `condition`  | string | No       | Condition after repair    |
| `notes`      | string | No       | Repair completion notes   |

```json
{
  "asset_id": 1,
  "repair_cost": 150.00,
  "condition": "good",
  "notes": "Display cable replaced under warranty"
}
```

---

### 24. Retire Asset

Permanently retire an asset. Only retired assets can be deleted.

**Endpoint:** `POST /api/v1/assets/retire`

#### Request Body

| Field              | Type   | Required | Description                          |
|--------------------|--------|----------|--------------------------------------|
| `asset_id`         | int    | Yes      | Asset ID                             |
| `retirement_reason`| string | No       | Reason for retirement                |
| `retirement_date`  | string | No       | Retirement date (YYYY-MM-DD)         |
| `notes`            | string | No       | Additional notes                     |

```json
{
  "asset_id": 2,
  "retirement_reason": "Beyond economical repair - 6 years old",
  "notes": "Disposed per IT asset disposal policy"
}
```

---

### 25. Delete Asset

Permanently delete an asset record. **Only retired assets can be deleted.**

**Endpoint:** `POST /api/v1/assets/delete`

#### Request Body

```json
{
  "id": 2
}
```

---

## 26. Enums & Status Values

### Asset Types

| Code               | Description              |
|--------------------|--------------------------|
| `laptop`           | Portable computer        |
| `mobile_phone`     | Company mobile phone     |
| `desktop`          | Desktop computer         |
| `tablet`           | Tablet device            |
| `id_card`          | Employee ID card         |
| `access_card`      | Building access card     |
| `vehicle`          | Company vehicle          |
| `office_equipment` | Office furniture/equipment |

### Asset Statuses

| Status         | Description                         |
|----------------|-------------------------------------|
| `available`    | Available for assignment            |
| `in_use`       | Currently assigned to an employee   |
| `under_repair` | Sent for repair                     |
| `retired`      | Permanently retired                 |

### Asset Conditions

| Condition   | Description    |
|-------------|----------------|
| `excellent` | Like new       |
| `good`      | Normal wear    |
| `fair`      | Visible wear   |
| `poor`      | Needs attention|

### Asset Request Types

| Type          | Description                  |
|---------------|------------------------------|
| `new`         | Request for a new asset      |
| `replacement` | Replace an existing asset    |
| `additional`  | Request an additional asset  |

### Asset Request Statuses

| Status      | Description                      |
|-------------|----------------------------------|
| `pending`   | Awaiting HR/Admin review         |
| `approved`  | Approved, waiting for fulfillment|
| `rejected`  | Rejected by HR/Admin             |
| `fulfilled` | Asset has been assigned          |
| `cancelled` | Cancelled by the employee        |

### Asset Issue Types

| Type          | Description             |
|---------------|-------------------------|
| `malfunction` | Asset not working       |
| `damage`      | Asset is damaged        |
| `lost`        | Asset is lost           |
| `stolen`      | Asset is stolen         |
| `other`       | Other issues            |

### Asset Issue Statuses

| Status         | Description          |
|----------------|----------------------|
| `reported`     | Newly reported       |
| `under_review` | HR/IT reviewing      |
| `in_progress`  | Being resolved       |
| `resolved`     | Issue resolved       |
| `closed`       | Issue closed         |

### Priority Levels (Requests & Issues)

| Priority | Description |
|----------|-------------|
| `low`    | Low priority|
| `medium` | Normal (default) |
| `high`   | High priority |
| `urgent` | Requires immediate attention |

---

## 27. Workflows

### Asset Request Workflow

```
Employee creates request  -->  [pending]
                                  |
                         HR reviews request
                        /                    \
                   Approves                 Rejects
                      |                        |
                 [approved]               [rejected]
                      |                   (reason provided)
              HR assigns asset
                      |
                 [fulfilled]
                 (asset assigned)
```

The employee can **cancel** a request at any time while it is `pending`.

### Asset Issue Workflow

```
Employee reports issue  -->  [reported]
                                |
                         HR/IT reviews
                                |
                         [under_review]
                                |
                        Work begins
                                |
                         [in_progress]
                                |
                        Issue fixed
                                |
                          [resolved]
                          (notes + action taken)
                                |
                           [closed]
```

### Asset Lifecycle

```
[Create Asset]  -->  available
                        |
                   Assign to employee
                        |
                      in_use  <---------> Reassign to another employee
                        |
                   Return asset
                        |
                    available
                        |
              Mark for repair
                        |
                  under_repair
                        |
              Complete repair
                        |
                    available
                        |
                 Retire asset
                        |
                     retired
                        |
                 Delete (permanent)
```

---

## 28. Data Models

### Asset

| Field                 | Type      | Description                             |
|-----------------------|-----------|-----------------------------------------|
| `id`                  | integer   | Primary key                             |
| `asset_code`          | string    | Unique asset code (auto-generated)      |
| `asset_type`          | string    | Type of asset                           |
| `brand`               | string    | Brand/manufacturer                      |
| `model`               | string    | Model name                              |
| `serial_number`       | string    | Unique serial number                    |
| `assigned_to`         | string    | Employee name (null if unassigned)      |
| `employee_id`         | string    | Employee ID code (null if unassigned)   |
| `department`          | string    | Employee's department                   |
| `assigned_date`       | datetime  | When assigned                           |
| `return_date`         | datetime  | When returned                           |
| `status`              | string    | Current status                          |
| `condition`           | string    | Physical condition                      |
| `purchase_date`       | datetime  | Purchase date                           |
| `warranty_expiry`     | datetime  | Warranty expiry                         |
| `value`               | number    | Purchase value                          |
| `notes`               | string    | General notes                           |
| `repair_reason`       | string    | Why sent for repair                     |
| `repair_notes`        | string    | Repair details                          |
| `repair_date`         | datetime  | When sent for repair                    |
| `repair_cost`         | number    | Cost of repair                          |
| `repair_completed_date`| datetime | When repair was completed               |
| `retirement_reason`   | string    | Why retired                             |
| `retirement_date`     | datetime  | When retired                            |

### Asset Request

| Field                | Type     | Description                               |
|----------------------|----------|-------------------------------------------|
| `id`                 | integer  | Primary key                               |
| `request_number`     | string   | Unique request number (e.g., AR-2026-001) |
| `employee_id`        | string   | Requesting employee's ID                  |
| `request_type`       | string   | new / replacement / additional            |
| `asset_type`         | string   | Type of asset requested                   |
| `brand`              | string   | Preferred brand (optional)                |
| `model`              | string   | Preferred model (optional)                |
| `priority`           | string   | low / medium / high / urgent              |
| `status`             | string   | Current request status                    |
| `justification`      | string   | Why the asset is needed                   |
| `replacing_asset_id` | integer  | Asset being replaced (if replacement)     |
| `replacing_asset_code`| string  | Code of asset being replaced              |
| `approved_by`        | integer  | HR user who approved                      |
| `approved_at`        | datetime | When approved                             |
| `rejected_at`        | datetime | When rejected                             |
| `rejection_reason`   | string   | Reason for rejection                      |
| `fulfilled_at`       | datetime | When asset was assigned                   |
| `assigned_asset_id`  | integer  | Asset that was assigned                   |
| `assigned_asset_code`| string   | Code of the assigned asset                |
| `requested_date`     | datetime | When request was submitted                |

### Asset Issue

| Field                  | Type     | Description                              |
|------------------------|----------|------------------------------------------|
| `id`                   | integer  | Primary key                              |
| `issue_number`         | string   | Unique issue number (e.g., AI-2026-001)  |
| `employee_id`          | string   | Reporting employee's ID                  |
| `asset_id`             | integer  | Related asset ID                         |
| `asset_code`           | string   | Related asset code                       |
| `issue_type`           | string   | malfunction / damage / lost / stolen / other |
| `priority`             | string   | low / medium / high / urgent             |
| `status`               | string   | Current issue status                     |
| `title`                | string   | Brief issue title                        |
| `description`          | string   | Detailed description                     |
| `reported_date`        | datetime | When reported                            |
| `resolved_at`          | datetime | When resolved                            |
| `resolved_by`          | integer  | User who resolved                        |
| `resolution_notes`     | string   | How it was resolved                      |
| `action_taken`         | string   | What action was taken                    |
| `incident_date`        | datetime | When incident occurred (lost/stolen)     |
| `incident_location`    | string   | Where it happened (lost/stolen)          |
| `police_report_number` | string   | Police report reference (stolen)         |
| `attachment_urls`      | string[] | Photos/documents                         |

---

## API Endpoints Summary

### Employee Self-Service (`/api/v1/self-service/assets`)

| Method | Endpoint                                        | Description              |
|--------|-------------------------------------------------|--------------------------|
| GET    | `/api/v1/self-service/assets`                   | Get my assigned assets   |
| POST   | `/api/v1/self-service/assets/requests`          | Request an asset         |
| GET    | `/api/v1/self-service/assets/requests`          | List my requests         |
| GET    | `/api/v1/self-service/assets/requests/:id`      | Get request details      |
| POST   | `/api/v1/self-service/assets/requests/:id/cancel`| Cancel my request       |
| POST   | `/api/v1/self-service/assets/issues`            | Report an issue          |
| GET    | `/api/v1/self-service/assets/issues`            | List my issues           |
| GET    | `/api/v1/self-service/assets/issues/:id`        | Get issue details        |

### HR/Admin — Requests (`/api/v1/assets/requests`)

| Method | Endpoint                                         | Description            |
|--------|--------------------------------------------------|------------------------|
| GET    | `/api/v1/assets/requests`                        | List all requests      |
| GET    | `/api/v1/assets/requests/:id`                    | Get request details    |
| POST   | `/api/v1/assets/requests/:id/approve`            | Approve request        |
| POST   | `/api/v1/assets/requests/:id/reject`             | Reject request         |
| POST   | `/api/v1/assets/requests/:id/fulfill`            | Fulfill request        |

### HR/Admin — Inventory (`/api/v1/assets`)

| Method | Endpoint                              | Description           |
|--------|---------------------------------------|-----------------------|
| GET    | `/api/v1/assets`                      | List assets           |
| GET    | `/api/v1/assets/get?id={id}`          | Get asset by ID       |
| GET    | `/api/v1/assets/types`                | Get asset types       |
| POST   | `/api/v1/assets`                      | Create asset          |
| POST   | `/api/v1/assets/update`               | Update asset          |
| POST   | `/api/v1/assets/assign`               | Assign to employee    |
| POST   | `/api/v1/assets/reassign`             | Reassign asset        |
| POST   | `/api/v1/assets/:asset_id/reassign`   | Reassign (alt route)  |
| POST   | `/api/v1/assets/return`               | Return asset          |
| POST   | `/api/v1/assets/mark-for-repair`      | Mark for repair       |
| POST   | `/api/v1/assets/complete-repair`      | Complete repair       |
| POST   | `/api/v1/assets/retire`               | Retire asset          |
| POST   | `/api/v1/assets/delete`               | Delete asset          |

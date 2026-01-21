# Asset HR/Admin Management API Documentation

## Overview
This document provides comprehensive API documentation for HR and Admin users to manage asset requests, approve/reject requests, fulfill requests by assigning assets, and reassign assets between employees.

## Base URL
All API endpoints use the base URL: `{{BASE_URL}}/api/v1/assets`

## Authentication
All requests require Bearer token authentication and HR or Admin role:
```
Authorization: Bearer <access_token>
```

**Note:** All endpoints require `HRMiddleware()` which allows both `admin` and `hr` roles.

---

## 1. List All Asset Requests

### Endpoint
```
GET /api/v1/assets/requests?page=1&page_size=20&employee_id=EMP001&request_type=new&status=pending&asset_type=laptop&priority=high
```

### Description
Retrieve a paginated list of all asset requests for HR/Admin review.

### Query Parameters
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Items per page (default: 20, max: 100)
- `employee_id` (string, optional): Filter by employee ID
- `request_type` (string, optional): Filter by request type - `"new"`, `"replacement"`, or `"additional"`
- `status` (string, optional): Filter by status - `"pending"`, `"approved"`, `"rejected"`, `"fulfilled"`, or `"cancelled"`
- `asset_type` (string, optional): Filter by asset type
- `priority` (string, optional): Filter by priority - `"low"`, `"medium"`, `"high"`, or `"urgent"`

### Success Response (200)
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
      "brand": "Apple",
      "model": "MacBook Pro 16\"",
      "status": "pending",
      "priority": "high",
      "justification": "Need a new laptop for development work",
      "requested_date": "2026-01-20T10:30:00+03:00",
      "notes": "Prefer 32GB RAM configuration",
      "replacing_asset_code": null,
      "approved_at": null,
      "rejected_at": null,
      "fulfilled_at": null
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

### Example Request (cURL)
```bash
curl -X GET "http://localhost:8080/api/v1/assets/requests?page=1&page_size=20&status=pending" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 2. Get Asset Request Details

### Endpoint
```
GET /api/v1/assets/requests/:request_id
```

### Description
Retrieve detailed information about a specific asset request.

### Path Parameters
- `request_id` (integer, required): Asset request ID

### Success Response (200)
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
    "brand": "Apple",
    "model": "MacBook Pro 16\"",
    "status": "approved",
    "priority": "high",
    "justification": "Need a new laptop for development work",
    "notes": "Prefer 32GB RAM configuration",
    "requested_date": "2026-01-20T10:30:00+03:00",
    "approved_at": "2026-01-21T09:15:00+03:00",
    "approved_by": 5,
    "fulfilled_at": "2026-01-22T14:30:00+03:00",
    "assigned_asset_code": "LAP-002"
  }
}
```

### Example Request (cURL)
```bash
curl -X GET "http://localhost:8080/api/v1/assets/requests/1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 3. Approve Asset Request

### Endpoint
```
POST /api/v1/assets/requests/:request_id/approve
```

### Description
Approve a pending asset request. This changes the status to `approved`, allowing it to be fulfilled later.

### Path Parameters
- `request_id` (integer, required): Asset request ID

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset request approved successfully",
  "data": {
    "id": 1,
    "request_number": "AR-2026-001",
    "status": "approved",
    "approved_at": "2026-01-21T09:15:00+03:00"
  }
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/assets/requests/1/approve" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 4. Reject Asset Request

### Endpoint
```
POST /api/v1/assets/requests/:request_id/reject
```

### Description
Reject a pending asset request with a reason.

### Path Parameters
- `request_id` (integer, required): Asset request ID

### Request Body
```json
{
  "rejection_reason": "Budget constraints. Request will be reviewed in next quarter."
}
```

### Request Parameters
- `rejection_reason` (string, required): Reason for rejection

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset request rejected successfully",
  "data": {
    "id": 1,
    "request_number": "AR-2026-001",
    "status": "rejected",
    "rejected_at": "2026-01-21T10:30:00+03:00",
    "rejection_reason": "Budget constraints. Request will be reviewed in next quarter."
  }
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/assets/requests/1/reject" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rejection_reason": "Budget constraints. Request will be reviewed in next quarter."
  }'
```

---

## 5. Fulfill Asset Request

### Endpoint
```
POST /api/v1/assets/requests/:request_id/fulfill
```

### Description
Fulfill an approved asset request by assigning an available asset to the requesting employee. This automatically:
- Assigns the asset to the employee
- Updates the asset status to `in_use`
- Marks the request as `fulfilled`
- Links the asset to the request

### Path Parameters
- `request_id` (integer, required): Asset request ID

### Request Body
```json
{
  "asset_id": 5
}
```

### Request Parameters
- `asset_id` (integer, required): ID of the asset to assign

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset request fulfilled successfully",
  "data": {
    "request": {
      "id": 1,
      "request_number": "AR-2026-001",
      "status": "fulfilled",
      "fulfilled_at": "2026-01-22T14:30:00+03:00",
      "assigned_asset_code": "LAP-002"
    },
    "asset": {
      "id": 5,
      "asset_code": "LAP-002",
      "asset_type": "laptop",
      "brand": "Apple",
      "model": "MacBook Pro 16\"",
      "employee_id": "EMP001",
      "status": "in_use",
      "assigned_date": "2026-01-22T14:30:00+03:00"
    }
  }
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/assets/requests/1/fulfill" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "asset_id": 5
  }'
```

### Error Responses

#### 400 Bad Request - Request Not Approved
```json
{
  "success": false,
  "message": "can only fulfill approved requests. Current status: pending",
  "error": null
}
```

#### 400 Bad Request - Asset Not Available
```json
{
  "success": false,
  "message": "asset is not available for assignment",
  "error": null
}
```

---

## 6. Reassign Asset

### Endpoint
```
POST /api/v1/assets/:asset_id/reassign
```

### Description
Reassign an asset from one employee to another. The asset must be currently assigned to an employee.

### Path Parameters
- `asset_id` (integer, required): Asset ID

### Request Body
```json
{
  "employee_id": "EMP002",
  "notes": "Reassigned due to employee transfer"
}
```

### Request Parameters
- `employee_id` (string, required): Employee ID of the new assignee
- `notes` (string, optional): Notes about the reassignment

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset reassigned successfully",
  "data": {
    "id": 5,
    "asset_code": "LAP-002",
    "asset_type": "laptop",
    "brand": "Apple",
    "model": "MacBook Pro 16\"",
    "employee_id": "EMP002",
    "assigned_to": "Jane Smith",
    "status": "in_use",
    "assigned_date": "2026-01-25T10:00:00+03:00"
  }
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/assets/5/reassign" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "EMP002",
    "notes": "Reassigned due to employee transfer"
  }'
```

### Error Responses

#### 400 Bad Request - Asset Not Assigned
```json
{
  "success": false,
  "message": "asset is not currently assigned to any employee",
  "error": null
}
```

#### 400 Bad Request - Same Employee
```json
{
  "success": false,
  "message": "asset is already assigned to this employee",
  "error": null
}
```

#### 400 Bad Request - Employee Not Active
```json
{
  "success": false,
  "message": "new employee is not active",
  "error": null
}
```

---

## Workflow Overview

### Asset Request Lifecycle

1. **Employee Creates Request** → Status: `pending`
   - Employee submits asset request via self-service API

2. **HR/Admin Reviews** → Status: `pending`
   - HR/Admin views request details
   - HR/Admin can approve or reject

3. **HR/Admin Approves** → Status: `approved`
   - Request is approved, ready for fulfillment

4. **HR/Admin Fulfills** → Status: `fulfilled`
   - HR/Admin assigns an available asset
   - Asset is assigned to the employee
   - Request is marked as fulfilled

### Alternative Flows

- **Rejection**: `pending` → `rejected` (with reason)
- **Cancellation**: Employee can cancel `pending` requests
- **Reassignment**: HR/Admin can reassign assets at any time

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid request ID",
  "error": null
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "User not authenticated",
  "error": "Unauthorized access"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "message": "HR or Admin access required",
  "error": "Forbidden"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Asset request not found",
  "error": "Resource not found"
}
```

---

## Notes

1. **Authorization**: All endpoints require HR or Admin role. Regular employees cannot access these endpoints.

2. **Tenant Isolation**: All operations are scoped to the authenticated user's tenant.

3. **Status Validation**:
   - Only `pending` requests can be approved or rejected
   - Only `approved` requests can be fulfilled
   - Fulfilled requests cannot be modified

4. **Asset Assignment**:
   - When fulfilling a request, the asset must be `available`
   - The asset is automatically assigned to the requesting employee
   - Asset status changes to `in_use`

5. **Reassignment**:
   - Can reassign assets that are currently `in_use`
   - Cannot reassign to the same employee
   - New employee must be active

6. **Audit Trail**: All operations record the `updated_by` user ID for audit purposes.

7. **Pagination**: List endpoints support pagination with `page` and `page_size` query parameters.

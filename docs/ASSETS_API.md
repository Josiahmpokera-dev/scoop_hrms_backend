# Assets & Equipment Management API Documentation

## Overview
This document provides comprehensive API documentation for the Assets & Equipment Management system in the HRMS application. All endpoints require JWT Bearer token authentication.

## Base URL
All API endpoints use the base URL: `{{BASE_URL}}/api/v1/assets`

## Authentication
All requests require Bearer token authentication:
```
Authorization: Bearer <access_token>
```

---

## 1. List Assets API

### Endpoint
```
GET /api/v1/assets?page=1&page_size=20&search=laptop&status=in_use&asset_type=laptop
```

### Query Parameters
- `page` (integer, optional): Page number for pagination (default: 1)
- `page_size` (integer, optional): Number of items per page (default: 20, max: 100)
- `search` (string, optional): Search term for asset code, serial number, brand, model, or assigned employee name
- `status` (string, optional): Filter by status (`in_use`, `available`, `under_repair`, `retired`)
- `asset_type` (string, optional): Filter by asset type (`laptop`, `mobile_phone`, `desktop`, `tablet`, `id_card`, `access_card`, `vehicle`, `office_equipment`)
- `assigned_to` (string, optional): Filter by assigned employee ID
- `department` (string, optional): Filter by department name

### Success Response (200)
```json
{
  "success": true,
  "message": "Assets retrieved successfully",
  "data": [
    {
      "id": 1,
      "asset_code": "LAP-001",
      "asset_type": "laptop",
      "brand": "Apple",
      "model": "MacBook Pro 16\"",
      "serial_number": "MBP2023001",
      "assigned_to": "John Doe",
      "employee_id": "EMP001",
      "employee_photo": "/img/avatars/thumb-1.jpg",
      "department": "Engineering",
      "assigned_date": "2024-01-15T03:00:00+03:00",
      "status": "in_use",
      "condition": "excellent",
      "purchase_date": "2023-12-10T03:00:00+03:00",
      "warranty_expiry": "2026-12-10T03:00:00+03:00",
      "value": 2500.00,
      "notes": "High-performance laptop for development",
      "created_at": "2024-01-15T10:30:51.181123+03:00",
      "updated_at": "2024-01-15T10:30:51.181123+03:00"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 25,
    "total_pages": 2
  }
}
```

### Example Request (cURL)
```bash
curl -X GET "http://localhost:8080/api/v1/assets?page=1&page_size=20&status=in_use" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Example Request (JavaScript)
```javascript
const response = await fetch('http://localhost:8080/api/v1/assets?page=1&page_size=20&status=in_use', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
const data = await response.json();
```

---

## 2. Get Asset by ID API

### Endpoint
```
GET /api/v1/assets/get?id=1
```

### Query Parameters
- `id` (integer, required): Asset ID

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset retrieved successfully",
  "data": {
    "id": 1,
    "asset_code": "LAP-001",
    "asset_type": "laptop",
    "brand": "Apple",
    "model": "MacBook Pro 16\"",
    "serial_number": "MBP2023001",
    "assigned_to": "John Doe",
    "employee_id": "EMP001",
    "employee_photo": "/img/avatars/thumb-1.jpg",
    "department": "Engineering",
    "assigned_date": "2024-01-15T03:00:00+03:00",
    "status": "in_use",
    "condition": "excellent",
    "purchase_date": "2023-12-10T03:00:00+03:00",
    "warranty_expiry": "2026-12-10T03:00:00+03:00",
    "value": 2500.00,
    "notes": "High-performance laptop for development",
    "created_at": "2024-01-15T10:30:51.181123+03:00",
    "updated_at": "2024-01-15T10:30:51.181123+03:00"
  }
}
```

### Example Request (cURL)
```bash
curl -X GET "http://localhost:8080/api/v1/assets/get?id=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 3. Create Asset API

### Endpoint
```
POST /api/v1/assets
```

### Request Body
```json
{
  "asset_type": "laptop",
  "brand": "Apple",
  "model": "MacBook Pro 16\"",
  "serial_number": "MBP2023001",
  "asset_code": "LAP-001",
  "purchase_date": "2023-12-10",
  "value": 2500.00,
  "warranty_expiry": "2026-12-10",
  "condition": "excellent",
  "assigned_to_employee_id": 1,
  "notes": "High-performance laptop for development"
}
```

**Note:** `assigned_to_employee_id` can be either:
- A number (database ID): `1`, `2`, `3`
- A string (employee_id): `"EMP001"`, `"EMP002"`
- `null` (if not assigning immediately)

### Alternative Request Body (using employee_id string)
```json
{
  "asset_type": "laptop",
  "brand": "Apple",
  "model": "MacBook Pro 16\"",
  "serial_number": "MBP2023001",
  "asset_code": "LAP-001",
  "purchase_date": "2023-12-10",
  "value": 2500.00,
  "warranty_expiry": "2026-12-10",
  "condition": "excellent",
  "assigned_to_employee_id": "EMP001",
  "notes": "High-performance laptop for development"
}
```

### Validation Rules
- `asset_type` (required): Must be one of: `laptop`, `mobile_phone`, `desktop`, `tablet`, `id_card`, `access_card`, `vehicle`, `office_equipment`
- `brand` (required): String, max 100 characters
- `model` (required): String, max 100 characters
- `serial_number` (required): Unique string, max 100 characters
- `asset_code` (optional): Auto-generated if not provided, format: `{TYPE_PREFIX}-{NUMBER}` (e.g., LAP-001, PHN-001)
- `purchase_date` (optional): Date string in YYYY-MM-DD format
- `value` (optional): Decimal number, positive value
- `warranty_expiry` (optional): Date string in YYYY-MM-DD format
- `condition` (optional): Must be one of: `excellent`, `good`, `fair`, `poor` (default: `excellent`)
- `assigned_to_employee_id` (optional): Employee identifier - accepts either:
  - Number (database ID): `1`, `2`, `3`, etc.
  - String (employee_id): `"EMP001"`, `"EMP002"`, etc.
- `notes` (optional): Text field for additional information

### Success Response (201)
```json
{
  "success": true,
  "message": "Asset created successfully",
  "data": {
    "id": 1,
    "asset_code": "LAP-001",
    "asset_type": "laptop",
    "brand": "Apple",
    "model": "MacBook Pro 16\"",
    "serial_number": "MBP2023001",
    "assigned_to": null,
    "employee_id": null,
    "employee_photo": null,
    "department": null,
    "assigned_date": null,
    "status": "available",
    "condition": "excellent",
    "purchase_date": "2023-12-10T03:00:00+03:00",
    "warranty_expiry": "2026-12-10T03:00:00+03:00",
    "value": 2500.00,
    "notes": "High-performance laptop for development",
    "created_at": "2024-01-15T10:30:51.181123+03:00",
    "updated_at": "2024-01-15T10:30:51.181123+03:00"
  }
}
```

### Error Response (400)
```json
{
  "success": false,
  "message": "asset with this serial number already exists"
}
```

### Example Request (cURL) - Using Database ID
```bash
curl -X POST "http://localhost:8080/api/v1/assets" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "asset_type": "laptop",
    "brand": "Apple",
    "model": "MacBook Pro 16\"",
    "serial_number": "MBP2023001",
    "purchase_date": "2023-12-10",
    "value": 2500.00,
    "warranty_expiry": "2026-12-10",
    "condition": "excellent",
    "assigned_to_employee_id": 1,
    "notes": "High-performance laptop for development"
  }'
```

### Example Request (cURL) - Using Employee ID String
```bash
curl -X POST "http://localhost:8080/api/v1/assets" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "asset_type": "laptop",
    "brand": "Apple",
    "model": "MacBook Pro 16\"",
    "serial_number": "MBP2023001",
    "purchase_date": "2023-12-10",
    "value": 2500.00,
    "warranty_expiry": "2026-12-10",
    "condition": "excellent",
    "assigned_to_employee_id": "EMP001",
    "notes": "High-performance laptop for development"
  }'
```

**Note:** `assigned_to_employee_id` accepts either:
- A number (database ID): `1`, `2`, `3` - looks up employee by database ID
- A string (employee_id): `"EMP001"`, `"EMP002"` - looks up employee by employee_id field
- `null` or omitted - asset will be created as "available" (not assigned)

---

## 4. Update Asset API

### Endpoint
```
POST /api/v1/assets/update
```

### Request Body
```json
{
  "id": 1,
  "asset_type": "laptop",
  "brand": "Apple",
  "model": "MacBook Pro 16\"",
  "serial_number": "MBP2023001",
  "condition": "good",
  "value": 2400.00,
  "notes": "Updated notes"
}
```

### Validation Rules
- `id` (required): Asset ID to update
- All other fields are optional and will only update provided values
- Same validation rules as create API for provided fields
- Retired assets cannot be updated

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset updated successfully",
  "data": {
    "id": 1,
    "asset_code": "LAP-001",
    "asset_type": "laptop",
    "brand": "Apple",
    "model": "MacBook Pro 16\"",
    "serial_number": "MBP2023001",
    "assigned_to": "John Doe",
    "employee_id": "EMP001",
    "employee_photo": "/img/avatars/thumb-1.jpg",
    "department": "Engineering",
    "assigned_date": "2024-01-15T03:00:00+03:00",
    "status": "in_use",
    "condition": "good",
    "purchase_date": "2023-12-10T03:00:00+03:00",
    "warranty_expiry": "2026-12-10T03:00:00+03:00",
    "value": 2400.00,
    "notes": "Updated notes",
    "created_at": "2024-01-15T10:30:51.181123+03:00",
    "updated_at": "2024-01-15T12:45:30.123456+03:00"
  }
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/assets/update" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "condition": "good",
    "value": 2400.00,
    "notes": "Updated notes"
  }'
```

---

## 5. Assign Asset to Employee API

### Endpoint
```
POST /api/v1/assets/assign
```

### Request Body
```json
{
  "asset_id": 1,
  "employee_id": "EMP001",
  "assigned_date": "2024-01-15",
  "notes": "Assigned for development work"
}
```

### Validation Rules
- `asset_id` (required): Asset ID to assign
- `employee_id` (required): Employee ID to assign to (e.g., "EMP001")
- `assigned_date` (optional): Date string in YYYY-MM-DD format (defaults to today)
- `notes` (optional): Assignment notes

### Business Rules
- Asset must have status "available"
- Employee must exist and be active
- Asset cannot be already assigned

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset assigned successfully",
  "data": {
    "id": 1,
    "asset_code": "LAP-001",
    "assigned_to": "John Doe",
    "employee_id": "EMP001",
    "employee_photo": "/img/avatars/thumb-1.jpg",
    "department": "Engineering",
    "assigned_date": "2024-01-15T03:00:00+03:00",
    "status": "in_use",
    "updated_at": "2024-01-15T12:45:30.123456+03:00"
  }
}
```

### Error Response (400)
```json
{
  "success": false,
  "message": "asset is not available for assignment"
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/assets/assign" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "asset_id": 1,
    "employee_id": "EMP001",
    "assigned_date": "2024-01-15",
    "notes": "Assigned for development work"
  }'
```

---

## 6. Return Asset API

### Endpoint
```
POST /api/v1/assets/return
```

### Request Body
```json
{
  "asset_id": 1,
  "return_date": "2024-06-15",
  "condition": "good",
  "notes": "Returned in good condition"
}
```

### Validation Rules
- `asset_id` (required): Asset ID to return
- `return_date` (optional): Date string in YYYY-MM-DD format (defaults to today)
- `condition` (optional): Update asset condition (`excellent`, `good`, `fair`, `poor`)
- `notes` (optional): Return notes

### Business Rules
- Asset must be currently assigned (status "in_use")
- Updates asset status to "available"
- Clears assignment fields

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset returned successfully",
  "data": {
    "id": 1,
    "asset_code": "LAP-001",
    "assigned_to": null,
    "employee_id": null,
    "employee_photo": null,
    "department": null,
    "assigned_date": null,
    "return_date": "2024-06-15T03:00:00+03:00",
    "status": "available",
    "condition": "good",
    "updated_at": "2024-06-15T12:45:30.123456+03:00"
  }
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/assets/return" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "asset_id": 1,
    "return_date": "2024-06-15",
    "condition": "good",
    "notes": "Returned in good condition"
  }'
```

---

## 7. Mark Asset for Repair API

### Endpoint
```
POST /api/v1/assets/mark-for-repair
```

### Request Body
```json
{
  "asset_id": 1,
  "repair_reason": "Screen cracked",
  "repair_notes": "Needs screen replacement"
}
```

### Validation Rules
- `asset_id` (required): Asset ID to mark for repair
- `repair_reason` (optional): Reason for repair
- `repair_notes` (optional): Additional repair notes

### Business Rules
- Asset must not already be under repair
- Asset must not be retired
- Updates asset status to "under_repair"

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset marked for repair successfully",
  "data": {
    "id": 1,
    "asset_code": "LAP-001",
    "status": "under_repair",
    "repair_reason": "Screen cracked",
    "repair_notes": "Needs screen replacement",
    "repair_date": "2024-06-15T03:00:00+03:00",
    "updated_at": "2024-06-15T12:45:30.123456+03:00"
  }
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/assets/mark-for-repair" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "asset_id": 1,
    "repair_reason": "Screen cracked",
    "repair_notes": "Needs screen replacement"
  }'
```

---

## 8. Complete Asset Repair API

### Endpoint
```
POST /api/v1/assets/complete-repair
```

### Request Body
```json
{
  "asset_id": 1,
  "repair_cost": 150.00,
  "condition": "good",
  "notes": "Screen replaced successfully"
}
```

### Validation Rules
- `asset_id` (required): Asset ID to complete repair for
- `repair_cost` (optional): Cost of repair
- `condition` (optional): Updated asset condition (`excellent`, `good`, `fair`, `poor`)
- `notes` (optional): Repair completion notes

### Business Rules
- Asset must have status "under_repair"
- Updates asset status back to "available"

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset repair completed successfully",
  "data": {
    "id": 1,
    "asset_code": "LAP-001",
    "status": "available",
    "condition": "good",
    "repair_cost": 150.00,
    "repair_completed_date": "2024-06-20T03:00:00+03:00",
    "updated_at": "2024-06-20T12:45:30.123456+03:00"
  }
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/assets/complete-repair" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "asset_id": 1,
    "repair_cost": 150.00,
    "condition": "good",
    "notes": "Screen replaced successfully"
  }'
```

---

## 9. Retire Asset API

### Endpoint
```
POST /api/v1/assets/retire
```

### Request Body
```json
{
  "asset_id": 1,
  "retirement_reason": "Obsolete equipment",
  "retirement_date": "2024-12-31",
  "notes": "Replaced with newer model"
}
```

### Validation Rules
- `asset_id` (required): Asset ID to retire
- `retirement_reason` (optional): Reason for retirement
- `retirement_date` (optional): Date string in YYYY-MM-DD format (defaults to today)
- `notes` (optional): Retirement notes

### Business Rules
- Asset status becomes "retired"
- Asset can no longer be assigned or used
- Clears assignment if any

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset retired successfully",
  "data": {
    "id": 1,
    "asset_code": "LAP-001",
    "status": "retired",
    "retirement_reason": "Obsolete equipment",
    "retirement_date": "2024-12-31T03:00:00+03:00",
    "updated_at": "2024-12-31T12:45:30.123456+03:00"
  }
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/assets/retire" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "asset_id": 1,
    "retirement_reason": "Obsolete equipment",
    "retirement_date": "2024-12-31",
    "notes": "Replaced with newer model"
  }'
```

---

## 10. Delete Asset API

### Endpoint
```
POST /api/v1/assets/delete
```

### Request Body
```json
{
  "id": 1
}
```

### Validation Rules
- `id` (required): Asset ID to delete

### Business Rules
- Only assets with status "retired" can be deleted
- Permanent deletion - cannot be undone

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset deleted successfully"
}
```

### Error Response (400)
```json
{
  "success": false,
  "message": "only retired assets can be deleted"
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/assets/delete" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1
  }'
```

---

## 11. Get Asset Types API

### Endpoint
```
GET /api/v1/assets/types
```

### Success Response (200)
```json
{
  "success": true,
  "message": "Asset types retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Laptop",
      "code": "laptop",
      "description": "Portable computers"
    },
    {
      "id": 2,
      "name": "Mobile Phone",
      "code": "mobile_phone",
      "description": "Mobile communication devices"
    },
    {
      "id": 3,
      "name": "Desktop",
      "code": "desktop",
      "description": "Desktop computers"
    },
    {
      "id": 4,
      "name": "Tablet",
      "code": "tablet",
      "description": "Tablet devices"
    },
    {
      "id": 5,
      "name": "ID Card",
      "code": "id_card",
      "description": "Employee identification cards"
    },
    {
      "id": 6,
      "name": "Access Card",
      "code": "access_card",
      "description": "Access control cards"
    },
    {
      "id": 7,
      "name": "Vehicle",
      "code": "vehicle",
      "description": "Company vehicles"
    },
    {
      "id": 8,
      "name": "Office Equipment",
      "code": "office_equipment",
      "description": "Office equipment and supplies"
    }
  ]
}
```

### Example Request (cURL)
```bash
curl -X GET "http://localhost:8080/api/v1/assets/types" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Asset Status Values

- `available` - Available for assignment
- `in_use` - Currently assigned to an employee
- `under_repair` - Currently being repaired
- `retired` - No longer in use

## Asset Type Values

- `laptop` - Laptop computers
- `mobile_phone` - Mobile phones
- `desktop` - Desktop computers
- `tablet` - Tablet devices
- `id_card` - Employee identification cards
- `access_card` - Access control cards
- `vehicle` - Company vehicles
- `office_equipment` - Office equipment and supplies

## Condition Values

- `excellent` - Excellent condition
- `good` - Good condition
- `fair` - Fair condition
- `poor` - Poor condition

---

## Asset Code Auto-Generation

Asset codes are auto-generated if not provided during creation. The format is:
```
{TYPE_PREFIX}-{SEQUENTIAL_NUMBER}
```

### Type Prefixes:
- `LAP` - Laptop
- `PHN` - Mobile Phone
- `DSK` - Desktop
- `TAB` - Tablet
- `IDC` - ID Card
- `ACC` - Access Card
- `VEH` - Vehicle
- `EQP` - Office Equipment
- `AST` - Default (for unknown types)

### Examples:
- `LAP-001`, `LAP-002`, `LAP-003`
- `PHN-001`, `PHN-002`
- `IDC-001`

---

## Status Transitions

The following status transitions are supported:

1. **Assignment Flow:**
   - `available` → `in_use` (when assigned to employee)

2. **Return Flow:**
   - `in_use` → `available` (when returned from employee)

3. **Repair Flow:**
   - `any` → `under_repair` (when marked for repair)
   - `under_repair` → `available` (when repair is completed)

4. **Retirement Flow:**
   - `any` → `retired` (when asset is retired)

---

## Business Rules

1. **Serial Numbers:** Must be unique across all assets
2. **Asset Codes:** Must be unique across all assets (auto-generated if not provided)
3. **Assignment:** Only one employee can be assigned to an asset at a time
4. **Assignment Restrictions:**
   - Assets under repair cannot be assigned
   - Retired assets cannot be assigned
   - Only available assets can be assigned
5. **Update Restrictions:**
   - Retired assets cannot be updated
6. **Deletion Restrictions:**
   - Only retired assets can be permanently deleted
7. **Employee Validation:**
   - Employee must exist and be active to receive asset assignment

---

## Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "message": "Error message description",
  "error": "Detailed error information (optional)"
}
```

### Common Error Codes:
- `400 Bad Request` - Invalid request data or business rule violation
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation errors

---

## Quick Reference

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/assets` | GET | List assets with pagination and filters |
| `/api/v1/assets/get` | GET | Get asset by ID |
| `/api/v1/assets/types` | GET | Get asset types |
| `/api/v1/assets` | POST | Create new asset |
| `/api/v1/assets/update` | POST | Update asset |
| `/api/v1/assets/assign` | POST | Assign asset to employee |
| `/api/v1/assets/return` | POST | Return asset from employee |
| `/api/v1/assets/mark-for-repair` | POST | Mark asset for repair |
| `/api/v1/assets/complete-repair` | POST | Complete asset repair |
| `/api/v1/assets/retire` | POST | Retire asset |
| `/api/v1/assets/delete` | POST | Delete asset (only retired) |

---

## Notes

- All date fields accept format: `YYYY-MM-DD`
- All timestamps in responses are in ISO 8601 format with timezone
- Asset codes are auto-generated if not provided
- Department names are automatically resolved from employee's department
- Employee photos default to `/img/avatars/default.jpg` if not available
- Multi-tenancy is enforced - users can only access assets from their tenant

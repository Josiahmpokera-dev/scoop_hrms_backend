# Leave Policies & Types API

This document defines the API contract for the **Leave Policies & Types** page.

Page: `Leave Policies & Types`  
Route: `/leave/policies`  
Frontend file: `src/app/(protected-pages)/leave/policies/page.tsx`

The page has 2 tabs:
- **Leave Types**
- **Leave Policies**

Each tab supports list + pagination + create + update + delete.

---

## Base URL
`http://localhost:8080/api/v1`

## Authentication
All endpoints require a valid Bearer token.

`Authorization: Bearer <token>`

## Authorization (Recommended)
- Leave Types: Read (any authenticated), Write (HR/Admin)
- Leave Policies: HR/Admin (management endpoints)
- Policy Guidelines: any authenticated user (used by employee request flow)

---

## Standard Response Envelope

### Success
```json
{
  "success": true,
  "message": "string",
  "data": {}
}
```

### Error
```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "field": "reason"
  }
}
```

---

# 1) Leave Types

## 1.1 Get Leave Types (Paginated)
Used by **Types** tab table.

**Endpoint**  
`GET /leave/types`

**Query Parameters**
| Parameter | Type | Required | Description |
|---|---|---:|---|
| `page` | number | No | Default `1` |
| `page_size` | number | No | Default `10` |
| `is_active` | boolean | No | Filter active/inactive |
| `category` | string | No | Filter by category (e.g. `earned`, `sick`, `maternity`) |

**Success Response**
```json
{
  "success": true,
  "message": "Leave types retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "code": "ANNUAL",
        "name": "Annual Leave",
        "icon": "🏖️",
        "category": "earned",
        "paid_leave": true,
        "requires_documentation": false,
        "is_active": true,
        "description": "Annual paid leave",
        "created_at": "2026-04-01T10:00:00Z",
        "updated_at": "2026-04-01T10:00:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "per_page": 10,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

---

## 1.2 Create Leave Type
Used by **Add Type** dialog.

**Endpoint**  
`POST /leave/types`

**Request Body**
```json
{
  "code": "ANNUAL",
  "name": "Annual Leave",
  "icon": "🏖️",
  "category": "earned",
  "paid_leave": true,
  "requires_documentation": false,
  "is_active": true,
  "description": "Annual paid leave"
}
```

**Validation**
- `code` required, unique, uppercase recommended
- `name` required
- `category` required

**Success Response**
```json
{
  "success": true,
  "message": "Leave type created successfully",
  "data": {
    "id": 1,
    "code": "ANNUAL",
    "name": "Annual Leave",
    "icon": "🏖️",
    "category": "earned",
    "paid_leave": true,
    "requires_documentation": false,
    "is_active": true,
    "description": "Annual paid leave",
    "created_at": "2026-04-01T10:00:00Z",
    "updated_at": "2026-04-01T10:00:00Z"
  }
}
```

---

## 1.3 Update Leave Type
Used by **Edit Type** action.

**Endpoint**  
`PUT /leave/types/:type_id`

**Request Body (partial allowed)**
```json
{
  "name": "Annual Leave (Updated)",
  "requires_documentation": true,
  "description": "Updated notes"
}
```

**Success Response**
```json
{
  "success": true,
  "message": "Leave type updated successfully",
  "data": {
    "id": 1,
    "code": "ANNUAL",
    "name": "Annual Leave (Updated)",
    "icon": "🏖️",
    "category": "earned",
    "paid_leave": true,
    "requires_documentation": true,
    "is_active": true,
    "description": "Updated notes",
    "created_at": "2026-04-01T10:00:00Z",
    "updated_at": "2026-04-02T08:30:00Z"
  }
}
```

---

## 1.4 Delete Leave Type
Used by **Delete Type** action.

**Endpoint**  
`DELETE /leave/types/:type_id`

**Success Response**
```json
{
  "success": true,
  "message": "Leave type deleted successfully",
  "data": null
}
```

**Constraint Recommendation**
- Block delete if type is referenced by policies; returns:
```json
{
  "success": false,
  "message": "Cannot delete leave type in use",
  "error": {
    "type_id": "referenced_by_policy"
  }
}
```

---

# 2) Leave Policies

## 2.1 Get Leave Policies (Paginated)
Used by **Policies** tab table.

**Endpoint**  
`GET /leave/policies`

**Query Parameters**
| Parameter | Type | Required | Description |
|---|---|---:|---|
| `page` | number | No | Default `1` |
| `page_size` | number | No | Default `10` |
| `country` | string | No | Country filter |
| `leave_type_code` | string | No | Filter by leave type code |
| `is_active` | boolean | No | Filter active/inactive |

**Success Response (UI-Compatible Shape)**
```json
{
  "success": true,
  "message": "Leave policies retrieved successfully",
  "data": {
    "data": [
      {
        "id": 10,
        "policy_name": "Default Annual Policy",
        "country": "Tanzania",
        "leave_type_code": "ANNUAL",
        "leave_type": {
          "id": 1,
          "code": "ANNUAL",
          "name": "Annual Leave",
          "category": "earned",
          "paid_leave": true,
          "requires_documentation": false,
          "is_active": true
        },
        "entitlement": 24,
        "accrual_frequency": "monthly",
        "proration_on_join": true,
        "proration_on_exit": true,
        "carry_forward": true,
        "encashment_allowed": false,
        "negative_balance_allowed": false,
        "half_day_allowed": true,
        "minimum_notice_days": 3,
        "maximum_days_per_request": 30,
        "is_active": true,
        "created_at": "2026-04-01T10:00:00Z",
        "updated_at": "2026-04-01T10:00:00Z"
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

## 2.2 Get Leave Policy By ID
Useful for edit/detail flows.

**Endpoint**  
`GET /leave/policies/:policy_id`

**Success Response (Detail Shape)**
```json
{
  "success": true,
  "message": "Leave policy retrieved successfully",
  "data": {
    "id": 10,
    "policy_name": "Default Annual Policy",
    "country": "Tanzania",
    "leave_type_code": "ANNUAL",
    "leave_type": {
      "id": 1,
      "code": "ANNUAL",
      "name": "Annual Leave"
    },
    "entitlement": 24,
    "accrual_frequency": "monthly",
    "proration_on_join": true,
    "proration_on_exit": true,
    "carry_forward": true,
    "carry_forward_limit": 10,
    "carry_forward_expiry": "2026-07-01T00:00:00Z",
    "encashment_allowed": false,
    "encashment_limit": null,
    "negative_balance_allowed": false,
    "sandwich_rules": false,
    "half_day_allowed": true,
    "minimum_notice_days": 3,
    "maximum_days_per_request": 30,
    "is_active": true
  }
}
```

---

## 2.3 Create Leave Policy
Used by **Add Policy** dialog.

**Endpoint**  
`POST /leave/policies`

**Request Body**
```json
{
  "policy_name": "Default Annual Policy",
  "country": "Tanzania",
  "leave_type_code": "ANNUAL",
  "entitlement": 24,
  "accrual_frequency": "monthly",
  "carry_forward": true,
  "carry_forward_limit": 10,
  "half_day_allowed": true,
  "minimum_notice_days": 3,
  "maximum_days_per_request": 30,
  "is_active": true
}
```

**Validation**
- `policy_name` required
- `country` required
- `leave_type_code` (or `leave_type_id`) required
- `entitlement` required and `> 0`
- `accrual_frequency` required (e.g., `monthly`, `yearly`)

**Success Response**
```json
{
  "success": true,
  "message": "Leave policy created successfully",
  "data": {
    "id": 10,
    "leave_type_code": "ANNUAL",
    "policy_name": "Default Annual Policy",
    "country": "Tanzania",
    "entitlement": 24,
    "accrual_frequency": "monthly",
    "carry_forward": true,
    "carry_forward_limit": 10,
    "half_day_allowed": true,
    "minimum_notice_days": 3,
    "maximum_days_per_request": 30,
    "is_active": true
  }
}
```

---

## 2.4 Update Leave Policy
Used by **Edit Policy** action.

**Endpoint**  
`PUT /leave/policies/:policy_id`

**Request Body (partial allowed)**
```json
{
  "entitlement": 30,
  "carry_forward_limit": 12,
  "is_active": true
}
```

**Success Response**
```json
{
  "success": true,
  "message": "Leave policy updated successfully",
  "data": {
    "id": 10,
    "leave_type_code": "ANNUAL",
    "policy_name": "Default Annual Policy",
    "entitlement": 30,
    "carry_forward_limit": 12,
    "is_active": true,
    "updated_at": "2026-04-02T11:20:00Z"
  }
}
```

---

## 2.5 Delete Leave Policy
Used by **Delete Policy** action.

**Endpoint**  
`DELETE /leave/policies/:policy_id`

**Success Response**
```json
{
  "success": true,
  "message": "Leave policy deleted successfully",
  "data": null
}
```

**Constraint**
- Block delete if policy is assigned to any employee in `employee_policies.leave_policy_id`; returns:
```json
{
  "success": false,
  "message": "Cannot delete leave policy in use",
  "error": {
    "policy_id": "assigned_to_employee"
  }
}
```

---

# 3) Optional Supporting Endpoint

## 3.1 Policy Guidelines
Useful for “policy hints” during request/apply flow.

**Endpoint**  
`GET /leave/policies/guidelines?leave_type_code=ANNUAL&country=Tanzania`

**Success Response**
```json
{
  "success": true,
  "message": "Policy guidelines retrieved successfully",
  "data": {
    "leave_type": "ANNUAL",
    "minimum_notice_days": 3,
    "maximum_days_per_request": 30,
    "carry_forward_limit": 10,
    "carry_forward_expiry": "2026-07-01T00:00:00Z",
    "half_day_allowed": true,
    "requires_documentation": false,
    "documentation_types": [],
    "guidelines": [
      "Submit at least 3 days in advance"
    ]
  }
}
```

---

## Common Error Cases

### 400 Bad Request
```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "entitlement_days": "must be greater than 0"
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

### 404 Not Found
```json
{
  "success": false,
  "message": "Leave policy not found"
}
```

### 409 Conflict
```json
{
  "success": false,
  "message": "Leave type code already exists",
  "error": {
    "code": "duplicate"
  }
}
```

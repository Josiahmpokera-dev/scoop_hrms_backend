# Organization API Documentation

## Overview

The Organization API allows you to manage organizations/companies in the HRMS system. Organizations represent the main company or business entity that employees belong to.

**Base URL:** `/api/v1/organizations`

**Authentication:** All endpoints require authentication via JWT token in the `Authorization` header and HR/Admin role.

---

## API Endpoints

### 1. Get My Organization

Get the current user's organization (created during onboarding).

**Endpoint:** `GET /api/v1/organizations/me`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "message": "Organization retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 1,
    "code": "ORG001",
    "name": "Acme Corporation",
    "legal_name": "Acme Corporation Ltd",
    "registration_number": "REG123456",
    "tax_id": "TAX789012",
    "industry": "Technology",
    "company_size": "medium",
    "description": "A leading technology company",
    "website": "https://www.acme.com",
    "domain": "acme.com",
    "logo_url": "https://www.acme.com/logo.png",
    "address_line1": "123 Main Street",
    "address_line2": "Suite 100",
    "city": "Dar es Salaam",
    "state": "Dar es Salaam",
    "country": "Tanzania",
    "postal_code": "11101",
    "phone_number": "+255222123456",
    "email": "info@acme.com",
    "founded_date": "2010-01-15T00:00:00Z",
    "fiscal_year_start": 1,
    "timezone": "Africa/Dar_es_Salaam",
    "currency_code": "TZS",
    "status": "active",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Example Request:**
```bash
GET {{BASE_URL}}/api/v1/organizations/me
Authorization: Bearer {{TOKEN}}
```

---

### 2. Create Organization

Create a new organization.

**Endpoint:** `POST /api/v1/organizations`

**Authentication:** Required

**Request Body:**
```json
{
  "code": "ORG002",
  "name": "Tech Solutions Inc",
  "legal_name": "Tech Solutions Incorporated",
  "registration_number": "REG789012",
  "tax_id": "TAX345678",
  "industry": "Software Development",
  "company_size": "small",
  "description": "A software development company",
  "website": "https://www.techsolutions.com",
  "domain": "techsolutions.com",
  "logo_url": "https://www.techsolutions.com/logo.png",
  "address_line1": "456 Tech Street",
  "address_line2": "Floor 2",
  "city": "Dar es Salaam",
  "state": "Dar es Salaam",
  "country": "Tanzania",
  "postal_code": "11102",
  "phone_number": "+255222654321",
  "email": "contact@techsolutions.com",
  "founded_date": "2015-06-01T00:00:00Z",
  "fiscal_year_start": 1,
  "timezone": "Africa/Dar_es_Salaam",
  "currency_code": "TZS",
  "status": "active",
  "is_active": true
}
```

**Required Fields:**
- `code` (string) - Unique organization code
- `name` (string) - Organization name

**Optional Fields:**
- `legal_name` (string) - Legal/registered name
- `registration_number` (string) - Business registration number
- `tax_id` (string) - Tax identification number
- `industry` (string) - Industry sector
- `company_size` (string) - startup, small, medium, large, enterprise
- `description` (string) - Organization description
- `website` (string) - Website URL
- `domain` (string) - Organization domain (e.g., company.com)
- `logo_url` (string) - Logo image URL
- `address_line1` (string) - Primary address
- `address_line2` (string) - Secondary address
- `city` (string) - City
- `state` (string) - State/Province
- `country` (string) - Country
- `postal_code` (string) - Postal/ZIP code
- `phone_number` (string) - Contact phone number
- `email` (string) - Contact email
- `founded_date` (string) - Date organization was founded (ISO 8601)
- `fiscal_year_start` (integer) - Fiscal year start month (1-12, default: 1)
- `timezone` (string) - Timezone (default: "Africa/Dar_es_Salaam")
- `currency_code` (string) - Currency code (default: "TZS")
- `status` (string) - active, suspended, inactive (default: "active")
- `is_active` (boolean) - Active status (default: true)

**Response:**
```json
{
  "success": true,
  "message": "Organization created successfully",
  "data": {
    "id": 2,
    "code": "ORG002",
    "name": "Tech Solutions Inc",
    "domain": "techsolutions.com",
    ...
  }
}
```

**Example Request:**
```bash
POST {{BASE_URL}}/api/v1/organizations
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "code": "ORG002",
  "name": "Tech Solutions Inc",
  "domain": "techsolutions.com"
}
```

---

### 3. List Organizations

Get a list of all organizations (paginated).

**Endpoint:** `GET /api/v1/organizations`

**Authentication:** Required

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Items per page (default: 20, max: 100)

**Response:**
```json
{
  "success": true,
  "message": "Organizations retrieved successfully",
  "data": [
    {
      "id": 1,
      "code": "ORG001",
      "name": "Acme Corporation",
      "domain": "acme.com",
      "status": "active",
      "is_active": true,
      ...
    },
    {
      "id": 2,
      "code": "ORG002",
      "name": "Tech Solutions Inc",
      "domain": "techsolutions.com",
      "status": "active",
      "is_active": true,
      ...
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

**Example Request:**
```bash
GET {{BASE_URL}}/api/v1/organizations?page=1&page_size=20
Authorization: Bearer {{TOKEN}}
```

---

### 4. Get Organization by ID

Get a specific organization by its ID.

**Endpoint:** `GET /api/v1/organizations/:id`

**Authentication:** Required

**URL Parameters:**
- `id` (required) - Organization ID

**Response:**
```json
{
  "success": true,
  "message": "Organization retrieved successfully",
  "data": {
    "id": 1,
    "code": "ORG001",
    "name": "Acme Corporation",
    "domain": "acme.com",
    ...
  }
}
```

**Example Request:**
```bash
GET {{BASE_URL}}/api/v1/organizations/1
Authorization: Bearer {{TOKEN}}
```

---

### 5. Update Organization

Update an existing organization.

**Endpoint:** `PUT /api/v1/organizations/:id`

**Authentication:** Required

**URL Parameters:**
- `id` (required) - Organization ID

**Request Body:**
```json
{
  "name": "Acme Corporation Updated",
  "domain": "acme-updated.com",
  "website": "https://www.acme-updated.com",
  "email": "newemail@acme.com",
  "phone_number": "+255222999999"
}
```

**Note:** All fields are optional. Only provided fields will be updated.

**Response:**
```json
{
  "success": true,
  "message": "Organization updated successfully",
  "data": {
    "id": 1,
    "code": "ORG001",
    "name": "Acme Corporation Updated",
    "domain": "acme-updated.com",
    ...
  }
}
```

**Example Request:**
```bash
PUT {{BASE_URL}}/api/v1/organizations/1
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "domain": "acme-updated.com",
  "website": "https://www.acme-updated.com"
}
```

---

### 6. Delete Organization

Delete (soft delete) an organization.

**Endpoint:** `DELETE /api/v1/organizations/:id`

**Authentication:** Required

**URL Parameters:**
- `id` (required) - Organization ID

**Response:**
```json
{
  "success": true,
  "message": "Organization deleted successfully",
  "data": null
}
```

**Example Request:**
```bash
DELETE {{BASE_URL}}/api/v1/organizations/1
Authorization: Bearer {{TOKEN}}
```

---

### 7. Action-Based API (POST)

Alternative POST-only endpoint for all operations.

**Endpoint:** `POST /api/v1/organizations/action`

**Authentication:** Required

**Request Body:**
```json
{
  "action": "create",
  "data": {
    "code": "ORG003",
    "name": "New Company",
    "domain": "newcompany.com"
  }
}
```

**Available Actions:**
- `create` - Create organization
- `read` - Get organization (requires `id` in request)
- `update` - Update organization (requires `id` in request)
- `delete` - Delete organization (requires `id` in request)

**Example - Create:**
```json
{
  "action": "create",
  "data": {
    "code": "ORG003",
    "name": "New Company",
    "domain": "newcompany.com"
  }
}
```

**Example - Read:**
```json
{
  "action": "read",
  "id": "1"
}
```

**Example - Update:**
```json
{
  "action": "update",
  "id": "1",
  "data": {
    "domain": "updated-domain.com"
  }
}
```

**Example - Delete:**
```json
{
  "action": "delete",
  "id": "1"
}
```

---

## Organization Fields

### Required Fields
- `code` - Unique organization code (string, max 50 chars)
- `name` - Organization name (string, max 255 chars)

### Optional Fields
- `legal_name` - Legal/registered name (string, max 255 chars)
- `registration_number` - Business registration number (string, max 100 chars)
- `tax_id` - Tax identification number (string, max 100 chars)
- `industry` - Industry sector (string, max 100 chars)
- `company_size` - Company size: startup, small, medium, large, enterprise (string, max 50 chars)
- `description` - Organization description (text)
- `website` - Website URL (string, max 255 chars)
- `domain` - Organization domain, e.g., "company.com" (string, max 255 chars) **NEW**
- `logo_url` - Logo image URL (string, max 500 chars)
- `address_line1` - Primary address (text)
- `address_line2` - Secondary address (text)
- `city` - City (string, max 100 chars)
- `state` - State/Province (string, max 100 chars)
- `country` - Country (string, max 100 chars)
- `postal_code` - Postal/ZIP code (string, max 20 chars)
- `phone_number` - Contact phone number (string, max 50 chars)
- `email` - Contact email (string, max 191 chars)
- `founded_date` - Date organization was founded (ISO 8601 date string)
- `fiscal_year_start` - Fiscal year start month: 1-12 (integer, default: 1)
- `timezone` - Timezone (string, default: "Africa/Dar_es_Salaam")
- `currency_code` - Currency code (string, 3 chars, default: "TZS")
- `status` - Organization status: active, suspended, inactive (string, default: "active")
- `is_active` - Active status (boolean, default: true)

---

## Status Values

- `active` - Organization is active and operational
- `suspended` - Organization is temporarily suspended
- `inactive` - Organization is inactive

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": "code is required"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized",
  "data": null
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Organization not found",
  "data": null
}
```

### 409 Conflict
```json
{
  "success": false,
  "message": "organization with this code already exists",
  "data": null
}
```

---

## JavaScript Examples

### Get My Organization
```javascript
const getMyOrganization = async () => {
  const response = await fetch(`${BASE_URL}/api/v1/organizations/me`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return await response.json();
};
```

### Create Organization
```javascript
const createOrganization = async (orgData) => {
  const response = await fetch(`${BASE_URL}/api/v1/organizations`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      code: "ORG002",
      name: "Tech Solutions Inc",
      domain: "techsolutions.com",
      website: "https://www.techsolutions.com",
      email: "contact@techsolutions.com"
    })
  });
  return await response.json();
};
```

### Update Organization
```javascript
const updateOrganization = async (orgId, updates) => {
  const response = await fetch(`${BASE_URL}/api/v1/organizations/${orgId}`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(updates)
  });
  return await response.json();
};

// Example usage
await updateOrganization(1, {
  domain: "newdomain.com",
  website: "https://www.newdomain.com"
});
```

### List Organizations
```javascript
const listOrganizations = async (page = 1, pageSize = 20) => {
  const response = await fetch(
    `${BASE_URL}/api/v1/organizations?page=${page}&page_size=${pageSize}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  return await response.json();
};
```

---

## Notes

1. **Tenant Isolation:** Organizations are automatically filtered by tenant. You can only see and manage organizations belonging to your tenant.

2. **Unique Code:** The `code` field must be unique across all organizations.

3. **Domain Field:** The new `domain` field stores the organization's domain name (e.g., "company.com") without the protocol or www prefix.

4. **Soft Delete:** Deleting an organization performs a soft delete. The record is marked as deleted but not removed from the database.

5. **Default Values:**
   - `fiscal_year_start`: 1 (January)
   - `timezone`: "Africa/Dar_es_Salaam"
   - `currency_code`: "TZS"
   - `status`: "active"
   - `is_active`: true

6. **Get My Organization:** The `/me` endpoint returns the organization associated with the current user's tenant (created during onboarding).

---

## Summary

The Organization API provides comprehensive management of organizations:

- ✅ Get current user's organization (`/me`)
- ✅ Create new organizations
- ✅ List all organizations (paginated)
- ✅ Get organization by ID
- ✅ Update organization details
- ✅ Delete organization (soft delete)
- ✅ Action-based API for POST-only operations
- ✅ New `domain` field for organization domain names
- ✅ Tenant isolation and security
- ✅ Full CRUD operations

All endpoints support the new `domain` field for storing organization domain names.

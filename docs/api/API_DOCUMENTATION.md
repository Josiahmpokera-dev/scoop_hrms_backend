# HRMS Backend API Documentation

## Overview

The HRMS Backend API provides comprehensive endpoints for managing human resources, including employees, departments, job positions, locations, and organizational structure. The API supports both RESTful endpoints and POST-only action-based requests for enhanced security.

**Base URL:** `/api/v1`

**Authentication:** Most endpoints require authentication via JWT token and HR/Admin role.

---

## Table of Contents

1. [Authentication](#authentication)
2. [Organizations API](#organizations-api)
3. [Departments API](#departments-api)
4. [Job Positions API](#job-positions-api)
5. [Locations API](#locations-api)
6. [Employees API](#employees-api)
7. [Common Response Formats](#common-response-formats)
8. [Error Handling](#error-handling)

---

## Authentication

The HRMS API uses JWT (JSON Web Tokens) for authentication. Users must authenticate to access protected endpoints.

### Authentication Endpoints

#### Base URL: `/api/v1/auth`

---

### 1. Onboarding (Signup + Organization Creation)

Complete onboarding flow that creates a user account, tenant, and organization in one step.

**Endpoint:** `POST /api/v1/auth/onboard`

**Request Body:**

```json
{
  "email": "john@company.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe",
  "name": "Acme Corporation",
  "legal_name": "Acme Corporation Ltd",
  "registration_number": "REG-123456",
  "industry": "Technology",
  "company_size": "medium",
  "country": "Tanzania",
  "timezone": "Africa/Dar_es_Salaam",
  "currency_code": "TZS"
}
```

**Required Fields:**
- `email` (string, valid email) - User email address
- `password` (string, min 6 characters) - User password
- `first_name` (string, min 2 characters) - User first name
- `last_name` (string, min 2 characters) - User last name
- `name` (string, min 2, max 255) - Organization name

**Optional Fields:**
- `username` (string) - Username (auto-generated from email if not provided)
- `legal_name` - Legal/registered organization name
- `registration_number` - Business registration number
- `industry` - Industry type
- `company_size` - Company size: `startup`, `small`, `medium`, `large`, `enterprise`
- `country` - Country name
- `timezone` - Timezone (default: `Africa/Dar_es_Salaam`)
- `currency_code` - Currency code (default: `TZS`)

**Success Response (201 Created):**

```json
{
  "success": true,
  "message": "Onboarding completed successfully. Welcome to HRMS!",
  "data": {
    "user": {
      "id": 1,
      "username": "john",
      "email": "john@company.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "admin",
      "is_active": true
    },
    "tenant": {
      "id": 1,
      "name": "Acme Corporation",
      "domain": "acme-corporation",
      "status": "active"
    },
    "organization": {
      "id": 1,
      "code": "ACME-2026",
      "name": "Acme Corporation",
      "status": "active"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "setup_complete": false
  }
}
```

**Error Responses:**
- **400 Bad Request** - Validation failed or user already exists
- **422 Unprocessable Entity** - Invalid input data

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/onboard \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@company.com",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe",
    "name": "Acme Corporation",
    "industry": "Technology",
    "company_size": "medium",
    "country": "Tanzania"
  }'
```

---

### 2. Sign In (Login)

Authenticate an existing user and receive an access token.

**Endpoint:** `POST /api/v1/auth/login`

**Request Body:**

```json
{
  "email": "john@company.com",
  "password": "securepassword123"
}
```

**Required Fields:**
- `email` (string, valid email) - User email address
- `password` (string, min 6 characters) - User password

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "username": "john",
      "email": "john@company.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "admin",
      "is_active": true
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

**Error Responses:**
- **400 Bad Request** - Validation failed
- **401 Unauthorized** - Invalid email or password, or account is deactivated

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@company.com",
    "password": "securepassword123"
  }'
```

**Response Fields:**
- `access_token` - JWT token to use for authenticated requests
- `token_type` - Always "Bearer"
- `expires_in` - Token expiration time in seconds (default: 86400 = 24 hours)
- `user` - User information object

---

### 3. Get User Profile

Get the authenticated user's profile information.

**Endpoint:** `GET /api/v1/auth/profile`

**Authentication:** Required (Bearer token)

**Headers:**

```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "username": "john",
    "email": "john@company.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "admin",
    "is_active": true
  }
}
```

**Error Responses:**
- **401 Unauthorized** - Missing or invalid authentication token
- **404 Not Found** - User not found

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json"
```

---

### 4. Get Setup Wizard Status

Get the current status of the organization setup wizard.

**Endpoint:** `GET /api/v1/auth/setup-wizard/status`

**Authentication:** Required (Bearer token)

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Setup wizard status retrieved successfully",
  "data": {
    "completed": false,
    "current_step": "departments",
    "progress": 20,
    "steps": [
      {
        "step": "organization",
        "title": "Organization Setup",
        "description": "Create your organization",
        "completed": true,
        "required": true
      },
      {
        "step": "departments",
        "title": "Departments",
        "description": "Create your first department",
        "completed": false,
        "required": false
      },
      {
        "step": "locations",
        "title": "Locations",
        "description": "Add your office locations",
        "completed": false,
        "required": false
      },
      {
        "step": "job_positions",
        "title": "Job Positions",
        "description": "Define job positions",
        "completed": false,
        "required": false
      },
      {
        "step": "employees",
        "title": "Employees",
        "description": "Add your first employee",
        "completed": false,
        "required": false
      }
    ]
  }
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/auth/setup-wizard/status \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json"
```

---

### 5. Register (Simple User Registration)

Register a new user account without creating an organization. This is for users who will be invited to an existing organization.

**Endpoint:** `POST /api/v1/auth/register`

**Request Body:**

```json
{
  "email": "jane@company.com",
  "password": "securepassword123",
  "first_name": "Jane",
  "last_name": "Smith",
  "username": "jane" // Optional, auto-generated if not provided
}
```

**Required Fields:**
- `email` (string, valid email) - User email address
- `password` (string, min 6 characters) - User password
- `first_name` (string, min 2 characters) - User first name
- `last_name` (string, min 2 characters) - User last name

**Optional Fields:**
- `username` (string, min 3 characters) - Username (auto-generated from email if not provided)
- `role` (string) - User role: `admin`, `hr`, `user` (default: `user`)

**Success Response (201 Created):**

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 2,
      "username": "jane",
      "email": "jane@company.com",
      "first_name": "Jane",
      "last_name": "Smith",
      "role": "user",
      "is_active": true
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane@company.com",
    "password": "securepassword123",
    "first_name": "Jane",
    "last_name": "Smith"
  }'
```

---

### Using Authentication Tokens

After successful login or onboarding, you'll receive an `access_token`. Use this token in the `Authorization` header for all protected endpoints:

```
Authorization: Bearer <your_access_token>
```

**Example:**

```bash
# Get user profile
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"

# Create a department
curl -X POST http://localhost:8080/api/v1/departments \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "code": "IT",
    "name": "Information Technology",
    "organization_id": 1
  }'
```

### Token Expiration

- Default expiration: **24 hours** (86400 seconds)
- Configured via `JWT_EXPIRY` environment variable
- Format: Go duration string (e.g., `24h`, `720h`, `30d`)

### Protected Endpoints

All protected endpoints require:
1. **JWT Token** in the `Authorization` header: `Bearer <token>`
2. **HR or Admin Role** - Users must have either `hr` or `admin` role
3. **Tenant ID** (optional) - Can be provided in `X-Tenant-ID` header for multi-tenant support

### Headers for Protected Endpoints

```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
X-Tenant-ID: <tenant_id> (optional)
```

---

## Organizations API

### Base URL: `/api/v1/organizations`

Organizations represent companies or business entities in the system. Organizations are the top-level entity in the organizational hierarchy, and departments belong to organizations.

**Note:** If you completed onboarding, your organization was already created. Use the update endpoint to fill in any missing information.

### Endpoints

#### 1. Get My Organization (After Onboarding)

Get the organization that was created during onboarding. This is a convenience endpoint to retrieve your organization without needing to know its ID.

**Endpoint:** `GET /api/v1/organizations/me`

**Authentication:** Required (Bearer token)

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Organization retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 1,
    "code": "ACME-2026",
    "name": "Acme Corporation",
    "legal_name": null,
    "registration_number": null,
    "tax_id": null,
    "industry": "Technology",
    "company_size": "medium",
    "description": null,
    "website": null,
    "logo_url": null,
    "address_line1": null,
    "address_line2": null,
    "city": null,
    "state": null,
    "country": "Tanzania",
    "postal_code": null,
    "phone_number": null,
    "email": null,
    "founded_date": null,
    "fiscal_year_start": 1,
    "timezone": "Africa/Dar_es_Salaam",
    "currency_code": "TZS",
    "status": "active",
    "is_active": true,
    "created_at": "2026-01-14T10:00:00Z",
    "updated_at": "2026-01-14T10:00:00Z"
  }
}
```

**Error Responses:**
- **401 Unauthorized** - Missing or invalid authentication token
- **404 Not Found** - No organization found for your account

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/organizations/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json"
```

---

#### 2. Create Organization

Create a new organization.

**Endpoint:** `POST /api/v1/organizations`

**Request Body:**

```json
{
  "code": "ACME-CORP",
  "name": "ACME Corporation",
  "legal_name": "ACME Corporation Limited",
  "registration_number": "REG-123456",
  "tax_id": "TAX-789012",
  "description": "Leading technology company",
  "website": "https://www.acme-corp.com",
  "email": "info@acme-corp.com",
  "phone": "+255712345678",
  "address": "123 Business Street",
  "city": "Dar es Salaam",
  "state": "Dar es Salaam",
  "country": "Tanzania",
  "postal_code": "11101",
  "status": "active",
  "is_active": true,
  "founded_date": "2010-01-15T00:00:00Z"
}
```

**Required Fields:**
- `code` (string, min 2, max 50) - Unique organization code
- `name` (string, min 2, max 255) - Organization name

**Optional Fields:**
- `legal_name` - Legal/registered name
- `registration_number` - Business registration number
- `tax_id` - Tax identification number
- `description` - Organization description
- `website` - Organization website URL
- `email` - Contact email
- `phone` - Contact phone number
- `address`, `city`, `state`, `country`, `postal_code` - Address information
- `status` - Organization status: `active`, `suspended`, `inactive` (default: `active`)
- `is_active` - Whether organization is active (default: true)
- `founded_date` - Date organization was founded

**Success Response (201 Created):**

```json
{
  "success": true,
  "message": "Organization created successfully",
  "data": {
    "id": 1,
    "code": "ACME-CORP",
    "name": "ACME Corporation",
    "legal_name": "ACME Corporation Limited",
    "registration_number": "REG-123456",
    "tax_id": "TAX-789012",
    "status": "active",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### 2. List Organizations

Get a paginated list of organizations.

**Endpoint:** `GET /api/v1/organizations`

**Query Parameters:**
- `page` (integer, optional) - Page number (default: 1)
- `page_size` (integer, optional) - Items per page (default: 20, max: 100)
- `is_active` (boolean, optional) - Filter by active status
- `status` (string, optional) - Filter by status (`active`, `suspended`, `inactive`)
- `country` (string, optional) - Filter by country

**Example Request:**
```
GET /api/v1/organizations?page=1&page_size=20&is_active=true&country=Tanzania
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Organizations retrieved successfully",
  "data": [
    {
      "id": 1,
      "code": "ACME-CORP",
      "name": "ACME Corporation",
      "country": "Tanzania",
      "status": "active",
      "is_active": true
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

#### 3. Get Organization by ID

Retrieve a specific organization by ID.

**Endpoint:** `GET /api/v1/organizations/:id`

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Organization retrieved successfully",
  "data": {
    "id": 1,
    "code": "ACME-CORP",
    "name": "ACME Corporation",
    "legal_name": "ACME Corporation Limited",
    "registration_number": "REG-123456",
    "tax_id": "TAX-789012",
    "description": "Leading technology company",
    "website": "https://www.acme-corp.com",
    "email": "info@acme-corp.com",
    "phone": "+255712345678",
    "address": "123 Business Street",
    "city": "Dar es Salaam",
    "state": "Dar es Salaam",
    "country": "Tanzania",
    "postal_code": "11101",
    "status": "active",
    "is_active": true,
    "founded_date": "2010-01-15T00:00:00Z",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### 4. Get My Organization (After Onboarding)

Get the organization that was created during onboarding. This is a convenience endpoint to retrieve your organization without needing to know its ID.

**Endpoint:** `GET /api/v1/organizations/me`

**Authentication:** Required (Bearer token)

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Organization retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 1,
    "code": "ACME-2026",
    "name": "Acme Corporation",
    "legal_name": null,
    "registration_number": null,
    "tax_id": null,
    "industry": "Technology",
    "company_size": "medium",
    "country": "Tanzania",
    "timezone": "Africa/Dar_es_Salaam",
    "currency_code": "TZS",
    "status": "active",
    "is_active": true,
    "created_at": "2026-01-14T10:00:00Z",
    "updated_at": "2026-01-14T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/organizations/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json"
```

---

#### 5. Update Organization (POST-Only Action-Based)

Update an existing organization. **This is the endpoint to use after onboarding to fill in missing organization data.**

**Endpoint:** `POST /api/v1/organizations/action`

**Authentication:** Required (Bearer token)

**Note:** This API uses POST-only action-based requests. The organization ID is passed in the request body, not in the URL. All fields are optional - only include the fields you want to update. This allows you to fill in missing data from onboarding without providing all fields.

**Request Body Format:**

```json
{
  "action": "update",
  "id": "1",
  "data": {
    // Organization fields to update (all optional)
  }
}
```

**Request Body Examples:**

**Example 1: Fill in missing address and contact information after onboarding**

```json
{
  "action": "update",
  "id": "1",
  "data": {
    "legal_name": "Acme Corporation Limited",
    "registration_number": "REG-123456",
    "tax_id": "TAX-789012",
    "address_line1": "123 Business Street",
    "address_line2": "Suite 100",
    "city": "Dar es Salaam",
    "state": "Dar es Salaam",
    "postal_code": "11101",
    "phone_number": "+255712345678",
    "email": "info@acme-corp.com",
    "website": "https://www.acme-corp.com",
    "description": "Leading technology company in Tanzania"
  }
}
```

**Example 2: Update only specific fields**

```json
{
  "action": "update",
  "id": "1",
  "data": {
    "logo_url": "https://www.acme-corp.com/logo.png",
    "founded_date": "2010-01-15T00:00:00Z",
    "fiscal_year_start": 7
  }
}
```

**All Optional Fields:**
- `code` - Organization code (must be unique if changed)
- `name` - Organization name
- `legal_name` - Legal/registered name
- `registration_number` - Business registration number
- `tax_id` - Tax identification number
- `industry` - Industry type
- `company_size` - Company size: `startup`, `small`, `medium`, `large`, `enterprise`
- `description` - Organization description
- `website` - Organization website URL
- `logo_url` - Logo URL
- `address_line1` - Primary address line
- `address_line2` - Secondary address line
- `city` - City
- `state` - State/Province
- `country` - Country
- `postal_code` - Postal/ZIP code
- `phone_number` - Contact phone number
- `email` - Contact email
- `founded_date` - Date organization was founded (ISO 8601 format)
- `fiscal_year_start` - Fiscal year start month (1-12)
- `timezone` - Timezone (e.g., "Africa/Dar_es_Salaam")
- `currency_code` - Currency code (3 letters, e.g., "TZS")
- `status` - Organization status: `active`, `suspended`, `inactive`
- `is_active` - Whether organization is active (boolean)

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Organization updated successfully",
  "data": {
    "id": 1,
    "code": "ACME-2026",
    "name": "Acme Corporation",
    "legal_name": "Acme Corporation Limited",
    "registration_number": "REG-123456",
    "tax_id": "TAX-789012",
    "address_line1": "123 Business Street",
    "city": "Dar es Salaam",
    "state": "Dar es Salaam",
    "country": "Tanzania",
    "postal_code": "11101",
    "phone_number": "+255712345678",
    "email": "info@acme-corp.com",
    "website": "https://www.acme-corp.com",
    "status": "active",
    "is_active": true,
    "updated_at": "2026-01-14T11:30:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request** - Invalid organization ID, validation failed, or organization code already exists
- **401 Unauthorized** - Missing or invalid authentication token
- **404 Not Found** - Organization not found

**cURL Examples:**

**Step 1: Get your organization ID (after onboarding)**

```bash
curl -X GET http://localhost:8080/api/v1/organizations/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json"
```

**Step 2: Update organization with missing data (POST-only)**

```bash
curl -X POST http://localhost:8080/api/v1/organizations/action \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "update",
    "id": "1",
    "data": {
      "legal_name": "Acme Corporation Limited",
      "registration_number": "REG-123456",
      "tax_id": "TAX-789012",
      "address_line1": "123 Business Street",
      "city": "Dar es Salaam",
      "state": "Dar es Salaam",
      "country": "Tanzania",
      "postal_code": "11101",
      "phone_number": "+255712345678",
      "email": "info@acme-corp.com",
      "website": "https://www.acme-corp.com",
      "description": "Leading technology company in Tanzania"
    }
  }'
```

**Step 3: Update only specific fields (POST-only)**

```bash
curl -X POST http://localhost:8080/api/v1/organizations/action \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "update",
    "id": "1",
    "data": {
      "logo_url": "https://www.acme-corp.com/logo.png",
      "founded_date": "2010-01-15T00:00:00Z"
    }
  }'
```

**Workflow After Onboarding:**

1. **Get your organization:**
   ```bash
   GET /api/v1/organizations/me
   ```
   This returns your organization with its ID.

2. **Update missing fields (POST-only):**
   ```bash
   POST /api/v1/organizations/action
   {
     "action": "update",
     "id": "1",
     "data": {
       // Only include fields you want to update
     }
   }
   ```
   Include only the fields you want to add/update in the `data` object. You can make multiple update calls to fill in data gradually.

#### 5. Delete Organization

Soft delete an organization.

**Endpoint:** `DELETE /api/v1/organizations/:id`

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Organization deleted successfully",
  "data": null
}
```

#### 6. Action-Based API (POST Only) - **PRIMARY METHOD**

Handle all operations via POST with action parameter. **This is the recommended method for all organization operations.**

**Endpoint:** `POST /api/v1/organizations/action`

**Authentication:** Required (Bearer token)

**Request Format:**

```json
{
  "action": "create|read|update|delete|list",
  "id": "1",  // Required for read, update, delete
  "data": {}, // Required for create, update
  "filters": {}, // Optional for list
  "pagination": { // Optional for list
    "page": 1,
    "page_size": 20
  }
}
```

**Request Body Examples:**

**Create:**
```json
{
  "action": "create",
  "data": {
    "code": "ACME-CORP",
    "name": "ACME Corporation",
    "country": "Tanzania"
  }
}
```

**Read:**
```json
{
  "action": "read",
  "id": "1"
}
```

**Update (Fill Missing Data After Onboarding):**
```json
{
  "action": "update",
  "id": "1",
  "data": {
    "legal_name": "Acme Corporation Limited",
    "registration_number": "REG-123456",
    "address_line1": "123 Business Street",
    "city": "Dar es Salaam",
    "phone_number": "+255712345678",
    "email": "info@acme-corp.com"
  }
}
```

**Delete:**
```json
{
  "action": "delete",
  "id": "1"
}
```

**List:**
```json
{
  "action": "list",
  "filters": {
    "is_active": true,
    "country": "Tanzania"
  },
  "pagination": {
    "page": 1,
    "page_size": 20
  }
}
```

**Success Responses:**

**Update Response (200 OK):**
```json
{
  "success": true,
  "message": "Organization updated successfully",
  "data": {
    "id": 1,
    "code": "ACME-2026",
    "name": "Acme Corporation",
    "legal_name": "Acme Corporation Limited",
    "registration_number": "REG-123456",
    "address_line1": "123 Business Street",
    "city": "Dar es Salaam",
    "phone_number": "+255712345678",
    "email": "info@acme-corp.com",
    "status": "active",
    "updated_at": "2026-01-14T11:30:00Z"
  }
}
```

**cURL Examples:**

**Update Organization (POST-only):**
```bash
curl -X POST http://localhost:8080/api/v1/organizations/action \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "update",
    "id": "1",
    "data": {
      "legal_name": "Acme Corporation Limited",
      "registration_number": "REG-123456",
      "address_line1": "123 Business Street",
      "city": "Dar es Salaam",
      "phone_number": "+255712345678",
      "email": "info@acme-corp.com"
    }
  }'
```

**Delete:**
```json
{
  "action": "delete",
  "id": "1"
}
```

**List:**
```json
{
  "action": "list",
  "filters": {
    "is_active": true,
    "country": "Tanzania"
  },
  "pagination": {
    "page": 1,
    "page_size": 20
  }
}
```

---

## Departments API

### Base URL: `/api/v1/departments`

Departments represent organizational units in the company hierarchy. Departments can have parent-child relationships to form organizational trees.

### Endpoints

#### 1. Create Department (POST-Only Action-Based)

Create a new department with all required organizational information.

**Endpoint:** `POST /api/v1/departments/action`

**Authentication:** Required (Bearer token)

**Request Body Format:**

⚠️ **IMPORTANT**: The action-based API requires wrapping your data in an `action` and `data` object. You cannot send the department fields directly.

**Correct Format:**
```json
{
  "action": "create",
  "data": {
    // Department fields go here inside "data"
  }
}
```

**❌ INCORRECT (Will cause error):**
```json
{
  "code": "IT",
  "name": "Information Technology"
}
```

**✅ CORRECT:**
```json
{
  "action": "create",
  "data": {
    "code": "IT",
    "name": "Information Technology"
  }
}
```

**Request Body Example:**

```json
{
  "action": "create",
  "data": {
    "organization_id": 1,
    "code": "IT",
    "level": "department",
    "name": "Information Technology",
    "description": "IT Department responsible for technology infrastructure",
    "parent_department_id": null,
    "deputy_manager": "Jane Smith",
    "location_id": 1,
    "cost_center": "IT Operations",
    "budget_allocated": 500000.00,
    "budget_currency": "TZS",
    "is_active": true
  }
}
```

**Example Matching Your Data:**

```json
{
  "action": "create",
  "data": {
    "code": "IT",
    "level": "department",
    "name": "Information Technology",
    "description": "IT Department responsible for technology infrastructure",
    "deputy_manager": "Jane Smith",
    "location_id": 1,
    "cost_center": "CC-001",
    "budget_allocated": 500000.00,
    "budget_currency": "TZS",
    "is_active": true
  }
}
```

**Required Fields:**
- `code` (string, min 2, max 20) - **Unique department code** (required)
- `name` (string, min 2, max 100) - Department name (required)

**Optional Fields:**
- `organization_id` (integer) - ID of organization this department belongs to (validated - must exist)
- `level` (string) - **Department level**: `company`, `business_unit`, `department`, `team`
  - If not provided and parent is set, level is inferred from parent (one level down)
  - If not provided and no parent, defaults to `company` (top level)
- `parent_department_id` (integer) - **Parent Department ID** (none for top level, default if not provided)
  - Set to `null` or omit for top-level departments
  - If provided, must reference an existing department
- `manager_id` (integer, optional) - **Head of Department** - ID of employee who manages this department
  - Optional: Can be omitted if no manager is assigned yet
  - If provided, must reference an existing employee
- `deputy_manager` (string) - **Deputy Manager** - Deputy manager name stored as a plain string (e.g., "John Doe", "Jane Smith")
  - Stored directly as a string value, no lookup or validation required
- `location_id` (integer, optional) - **Location ID** - Direct reference to location ID from `locations` table (recommended)
- `location` (string, optional) - **Location Name** - Alternative to `location_id`. Name of location (e.g., "Dar es Salaam Office"). Case-insensitive search.
  - The system will search for a location matching this name
- `cost_center` (string) - **Cost Center** - Cost center name stored as a plain string (e.g., "CC-001", "IT Operations")
  - Stored directly as a string value, no lookup or validation required
- `budget_allocated` (float) - **Budget** - Allocated budget amount
- `budget_currency` (string, 3 chars) - Budget currency code (default: "TZS" if budget_allocated is provided)
- `description` (string) - Department description
- `department_type` (string) - Department type: `core`, `support`, `operational`, `strategic`
- `employee_capacity` (integer) - Maximum number of employees
- `is_active` (boolean) - Whether department is active (default: true)

**Level Hierarchy:**
- `company` - Top level (no parent)
- `business_unit` - Second level (parent is company)
- `department` - Third level (parent is business_unit)
- `team` - Fourth level (parent is department)

**Success Response (201 Created):**

```json
{
  "success": true,
  "message": "Department created successfully",
  "data": {
    "id": 1,
    "organization_id": 1,
    "code": "IT",
    "level": "department",
    "name": "Information Technology",
    "description": "IT Department responsible for technology infrastructure",
    "parent_department_id": null,
    "manager_id": null,
    "deputy_manager": "Jane Smith",
    "location": "Dar es Salaam Office",
    "cost_center": "IT Operations",
    "budget_allocated": 500000.00,
    "budget_currency": "TZS",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request** - Validation failed, department code already exists, or referenced entity not found
- **401 Unauthorized** - Missing or invalid authentication token
- **404 Not Found** - Referenced organization, location, or parent department not found

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/departments/action \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "create",
    "data": {
      "organization_id": 1,
      "code": "IT",
      "level": "department",
      "name": "Information Technology",
      "description": "IT Department responsible for technology infrastructure",
      "parent_department_id": null,
      "deputy_manager": "Jane Smith",
      "location_id": 1,
      "cost_center": "IT Operations",
      "budget_allocated": 500000.00,
      "budget_currency": "TZS"
    }
  }'
```

**Example: Create Top-Level Department (Company Level)**

```json
{
  "action": "create",
  "data": {
    "organization_id": 1,
    "code": "HQ",
    "level": "company",
    "name": "Headquarters",
    "parent_department_id": null
  }
}
```

**Example: Create Department with Parent (Level Auto-Inferred)**

```json
{
  "action": "create",
  "data": {
    "organization_id": 1,
    "code": "IT-DEPT",
    "name": "Information Technology",
    "parent_department_id": 1,
    "deputy_manager": "John Doe",
    "location_id": 1,
    "cost_center": "Technology Budget",
    "budget_allocated": 1000000.00
  }
}
```

**Example: Create Department with Location by Name (Alternative)**

```json
{
  "action": "create",
  "data": {
    "organization_id": 1,
    "code": "HR-DEPT",
    "name": "Human Resources",
    "parent_department_id": 1,
    "location": "Dar es Salaam Office",
    "cost_center": "HR Operations",
    "budget_allocated": 500000.00
  }
}
```

**Note:** You can use either `location_id` (recommended) or `location` (name). If both are provided, `location_id` takes precedence.

**Important Notes:**
- **Optional Manager**: `manager_id` is optional - you can create a department without assigning a manager initially
- **Location Options**: Location can be provided in two ways:
  - `location_id` (integer) - **Recommended**: Direct reference to location ID from the `locations` table. Validated to ensure it exists and belongs to your tenant.
  - `location` (string) - Alternative: Location name (case-insensitive). The system will search within your organization first, then tenant-wide if not found.
  - **Priority**: If both `location_id` and `location` are provided, `location_id` takes precedence.
- **Plain String Fields**: `deputy_manager` and `cost_center` are stored as plain string values (no lookup or validation)
  - `deputy_manager`: Simply provide the deputy manager name as a string (e.g., "John Doe", "Jane Smith")
  - `cost_center`: Simply provide the cost center name/code as a string (e.g., "CC-001", "IT Operations")
- **Location Flexibility**: When creating departments, you can use locations from any organization within your tenant (not restricted to the department's organization)
- **Level auto-inference**: If `level` is not provided but `parent_department_id` is set, the level will be automatically inferred from the parent (one level down)
- **Error messages**: If a referenced entity is not found, you'll receive a clear error message like "location 'Dar es Salaam' not found" or "location with ID '5' not found"

#### 2. List Departments

Get a paginated list of departments.

**Endpoint:** `GET /api/v1/departments`

**Query Parameters:**
- `page` (integer, optional) - Page number (default: 1)
- `page_size` (integer, optional) - Items per page (default: 20, max: 100)
- `is_active` (boolean, optional) - Filter by active status
- `parent_department_id` (integer, optional) - Filter by parent department (use "null" for root departments)

**Example Request:**
```
GET /api/v1/departments?page=1&page_size=20&is_active=true
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Departments retrieved successfully",
  "data": [
    {
      "id": 1,
      "code": "IT",
      "name": "Information Technology",
      "is_active": true
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

#### 3. Get Department by ID

Retrieve a specific department by ID.

**Endpoint:** `GET /api/v1/departments/:id`

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Department retrieved successfully",
  "data": {
    "id": 1,
    "code": "IT",
    "name": "Information Technology",
    "description": "IT Department",
    "parent_department_id": null,
    "manager_id": 1,
    "is_active": true,
    "parent_department": null,
    "sub_departments": []
  }
}
```

#### 4. Get Root Departments

Get all root departments (departments without a parent).

**Endpoint:** `GET /api/v1/departments/root`

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Root departments retrieved successfully",
  "data": [
    {
      "id": 1,
      "code": "IT",
      "name": "Information Technology",
      "sub_departments": [...]
    }
  ]
}
```

#### 5. Update Department (POST-Only Action-Based)

Update an existing department. **This is the recommended method using POST-only action-based requests.**

**Endpoint:** `POST /api/v1/departments/action`

**Authentication:** Required (Bearer token)

**Request Body Format:**

```json
{
  "action": "update",
  "id": "1",
  "data": {
    // Only include fields you want to update (all optional)
  }
}
```

**Request Body Examples:**

**Example 1: Update Basic Information**

```json
{
  "action": "update",
  "id": "1",
  "data": {
    "name": "Updated Department Name",
    "description": "Updated description",
    "is_active": true
  }
}
```

**Example 2: Update with Name-Based Fields**

```json
{
  "action": "update",
  "id": "1",
  "data": {
    "deputy_manager": "John Doe",
    "location": "New Location Name",
    "cost_center": "Updated Cost Center",
    "budget_allocated": 750000.00,
    "budget_currency": "TZS"
  }
}
```

**Example 3: Update Manager and Other Fields**

```json
{
  "action": "update",
  "id": "1",
  "data": {
    "manager_id": 5,
    "deputy_manager": "Jane Smith",
    "location": "Dar es Salaam Office",
    "cost_center": "IT Operations",
    "level": "business_unit"
  }
}
```

**All Optional Fields (same as Create):**
- `code` - Department code (must be unique if changed)
- `name` - Department name
- `level` - Department level: `company`, `business_unit`, `department`, `team`
- `description` - Department description
- `department_type` - Department type: `core`, `support`, `operational`, `strategic`
- `parent_department_id` - Parent Department ID
- `manager_id` (integer, optional) - Head of Department - ID of employee
- `deputy_manager` (string) - Deputy Manager name - Stored as a plain string (no lookup required)
- `location` (string) - Location name
- `cost_center` (string) - Cost Center name - Stored as a plain string (no lookup required)
- `budget_allocated` (float) - Budget amount
- `budget_currency` (string) - Budget currency code
- `employee_capacity` (integer) - Maximum number of employees
- `is_active` (boolean) - Whether department is active

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Department updated successfully",
  "data": {
    "id": 1,
    "code": "IT",
    "level": "department",
    "name": "Updated Department Name",
    "description": "Updated description",
    "deputy_manager": "John Doe",
    "location": "New Location Name",
    "cost_center": "Updated Cost Center",
    "budget_allocated": 750000.00,
    "budget_currency": "TZS",
    "is_active": true,
    "updated_at": "2024-01-16T14:20:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request** - Invalid department ID, validation failed, or referenced entity not found
- **401 Unauthorized** - Missing or invalid authentication token
- **404 Not Found** - Department not found

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/departments/action \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "update",
    "id": "1",
    "data": {
      "deputy_manager": "John Doe",
      "location_id": 1,
      "cost_center": "IT Operations",
      "budget_allocated": 750000.00
    }
  }'
```

**Note:** You can also use the RESTful endpoint `PUT /api/v1/departments/:id` if you prefer, but the POST-only action-based method is recommended for consistency.

#### 6. Delete Department

Soft delete a department.

**Endpoint:** `DELETE /api/v1/departments/:id`

**Note:** Cannot delete departments that have sub-departments.

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Department deleted successfully",
  "data": null
}
```

#### 7. Action-Based API (POST Only)

Handle all operations via POST with action parameter.

**Endpoint:** `POST /api/v1/departments/action`

**Request Body Examples:**

**Create:**
```json
{
  "action": "create",
  "data": {
    "code": "IT",
    "name": "Information Technology",
    "description": "IT Department"
  }
}
```

**Read:**
```json
{
  "action": "read",
  "id": "1"
}
```

**Update:**
```json
{
  "action": "update",
  "id": "1",
  "data": {
    "name": "Updated Department Name",
    "deputy_manager": "John Doe",
    "location": "New Location",
    "cost_center": "Updated Cost Center",
    "budget_allocated": 750000.00
  }
}
```

**Note:** All fields in the `data` object are optional. Only include the fields you want to update. The `id` field is required for update action.

**Delete:**
```json
{
  "action": "delete",
  "id": "1"
}
```

**List:**
```json
{
  "action": "list",
  "filters": {
    "is_active": true
  },
  "pagination": {
    "page": 1,
    "page_size": 20
  }
}
```

---

## Job Positions API

### Base URL: `/api/v1/job-positions`

Job positions represent job titles, roles, and grades within the organization.

### Endpoints

#### 1. Create Job Position

Create a new job position.

**Endpoint:** `POST /api/v1/job-positions`

**Request Body:**

```json
{
  "code": "SWE-001",
  "title": "Senior Software Engineer",
  "grade": "P4",
  "level": 4,
  "description": "Senior software engineer role",
  "is_active": true
}
```

**Required Fields:**
- `code` (string, min 2, max 20) - Unique position code
- `title` (string, min 2, max 100) - Job title

**Optional Fields:**
- `grade` - Position grade (e.g., "P4", "M2")
- `level` - Position level (integer)
- `description` - Position description
- `is_active` - Whether position is active (default: true)

**Success Response (201 Created):**

```json
{
  "success": true,
  "message": "Job position created successfully",
  "data": {
    "id": 1,
    "code": "SWE-001",
    "title": "Senior Software Engineer",
    "grade": "P4",
    "level": 4,
    "description": "Senior software engineer role",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### 2. List Job Positions

Get a paginated list of job positions.

**Endpoint:** `GET /api/v1/job-positions`

**Query Parameters:**
- `page` (integer, optional) - Page number (default: 1)
- `page_size` (integer, optional) - Items per page (default: 20, max: 100)
- `is_active` (boolean, optional) - Filter by active status
- `grade` (string, optional) - Filter by grade
- `level` (integer, optional) - Filter by level

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Job positions retrieved successfully",
  "data": [
    {
      "id": 1,
      "code": "SWE-001",
      "title": "Senior Software Engineer",
      "grade": "P4",
      "level": 4,
      "is_active": true
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

#### 3. Get Job Position by ID

**Endpoint:** `GET /api/v1/job-positions/:id`

#### 4. Update Job Position

**Endpoint:** `PUT /api/v1/job-positions/:id`

#### 5. Delete Job Position

**Endpoint:** `DELETE /api/v1/job-positions/:id`

#### 6. Action-Based API

**Endpoint:** `POST /api/v1/job-positions/action`

Follows the same pattern as Departments action-based API.

---

## Locations API

### Base URL: `/api/v1/locations`

Locations represent physical work locations, offices, or branches. **Locations are linked to Organizations** and can be used across departments within the same tenant.

### Endpoints

#### 1. Create Location (POST-Only Action-Based)

Create a new location linked to an organization.

**Endpoint:** `POST /api/v1/locations/action`

**Authentication:** Required (Bearer token)

**Request Body Format:**

⚠️ **IMPORTANT**: The action-based API requires wrapping your data in an `action` and `data` object.

```json
{
  "action": "create",
  "data": {
    // Location fields go here inside "data"
  }
}
```

**Request Body Example:**

```json
{
  "action": "create",
  "data": {
    "organization_id": 1,
    "name": "Dar es Salaam Office",
    "location_type": "branch",
    "address_line1": "123 Main Street",
    "address_line2": "Suite 100",
    "city": "Dar es Salaam",
    "state": "Dar es Salaam",
    "country": "Tanzania",
    "postal_code": "11101",
    "phone_number": "+255712345678",
    "timezone": "Africa/Dar_es_Salaam",
    "capacity": 50,
    "is_head_office": false,
    "is_active": true
  }
}
```

**Required Fields:**
- `name` (string, min 2, max 100) - **Location name** (required)

**Optional Fields:**
- `organization_id` (integer) - **Organization ID** - Links location to an organization (validated - must exist)
- `location_type` (string) - Location type: `head_office`, `branch`, `remote`, `satellite`
- `address_line1` (string) - Primary address line
- `address_line2` (string) - Secondary address line
- `city` (string) - City name
- `state` (string) - State/Province
- `country` (string) - Country name
- `postal_code` (string) - Postal/ZIP code
- `latitude` (float) - Latitude coordinate
- `longitude` (float) - Longitude coordinate
- `timezone` (string) - Timezone (default: "Africa/Dar_es_Salaam")
- `phone_number` (string) - Contact phone number
- `capacity` (integer) - Maximum number of employees
- `facilities` (string) - JSON array string of available facilities
- `is_head_office` (boolean) - Whether this is the head office (default: false)
  - **Note:** Only one location can be marked as head office per tenant. Setting a location as head office will automatically unset others.
- `is_active` (boolean) - Whether location is active (default: true)

**Important Notes:**
- **Organization Linking**: Locations should be linked to an organization using `organization_id`
- **Cross-Organization Usage**: When creating departments, you can use locations from any organization within your tenant (not restricted to the department's organization)
- **Name Matching**: Location names are case-insensitive when searching (e.g., "Dar es Salaam" matches "dar es salaam")

**Success Response (201 Created):**

```json
{
  "success": true,
  "message": "Location created successfully",
  "data": {
    "id": 1,
    "organization_id": 1,
    "name": "Dar es Salaam Office",
    "location_type": "branch",
    "address_line1": "123 Main Street",
    "address_line2": "Suite 100",
    "city": "Dar es Salaam",
    "state": "Dar es Salaam",
    "country": "Tanzania",
    "postal_code": "11101",
    "phone_number": "+255712345678",
    "timezone": "Africa/Dar_es_Salaam",
    "capacity": 50,
    "is_head_office": false,
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request** - Validation failed or organization not found
- **401 Unauthorized** - Missing or invalid authentication token

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/locations/action \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "create",
    "data": {
      "organization_id": 1,
      "name": "Dar es Salaam Office",
      "location_type": "branch",
      "address_line1": "123 Main Street",
      "city": "Dar es Salaam",
      "state": "Dar es Salaam",
      "country": "Tanzania",
      "postal_code": "11101",
      "phone_number": "+255712345678",
      "is_active": true
    }
  }'
```

**Example: Create Head Office Location**

```json
{
  "action": "create",
  "data": {
    "organization_id": 1,
    "name": "Head Office - Dar es Salaam",
    "location_type": "head_office",
    "address_line1": "123 Main Street",
    "city": "Dar es Salaam",
    "country": "Tanzania",
    "is_head_office": true
  }
}
```

#### 2. List Locations

Get a paginated list of locations.

**Endpoint:** `GET /api/v1/locations`

**Query Parameters:**
- `page` (integer, optional) - Page number (default: 1)
- `page_size` (integer, optional) - Items per page (default: 20, max: 100)
- `is_active` (boolean, optional) - Filter by active status
- `is_head_office` (boolean, optional) - Filter by head office status
- `country` (string, optional) - Filter by country
- `city` (string, optional) - Filter by city

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Locations retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Head Office - Dar es Salaam",
      "city": "Dar es Salaam",
      "country": "Tanzania",
      "is_head_office": true,
      "is_active": true
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 10,
    "total_pages": 1
  }
}
```

#### 3. Get Location by ID

**Endpoint:** `GET /api/v1/locations/:id`

#### 4. Get Head Office

Get the head office location for the tenant.

**Endpoint:** `GET /api/v1/locations/head-office`

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Head office location retrieved successfully",
  "data": {
    "id": 1,
    "name": "Head Office - Dar es Salaam",
    "is_head_office": true
  }
}
```

#### 5. Update Location

**Endpoint:** `PUT /api/v1/locations/:id`

#### 6. Delete Location

**Endpoint:** `DELETE /api/v1/locations/:id`

#### 7. Action-Based API

**Endpoint:** `POST /api/v1/locations/action`

Follows the same pattern as Departments action-based API.

---

## Employees API

### Base URL: `/api/v1/employees`

Employee management endpoints with integration to departments, positions, and locations.

### Updated Endpoints

#### 1. Onboard New Employee

Create a new employee record (onboarding).

**Endpoint:** `POST /api/v1/employees`

**Request Body:**

```json
{
  "employee_id": "EMP001",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+255712345678",
  "date_of_birth": "1990-01-15T00:00:00Z",
  "gender": "male",
  "address": "123 Main Street",
  "city": "Dar es Salaam",
  "state": "Dar es Salaam",
  "country": "Tanzania",
  "postal_code": "11101",
  "department_id": 1,
  "position_id": 1,
  "location_id": 1,
  "hire_date": "2024-01-15T00:00:00Z",
  "employment_type": "full-time",
  "salary": 2000000.00,
  "currency": "TZS",
  "emergency_contact_name": "Jane Doe",
  "emergency_contact_phone": "+255712345679",
  "emergency_contact_relation": "spouse",
  "notes": "New employee onboarding"
}
```

**Required Fields:**
- `employee_id` (string, min 3) - Unique employee identifier
- `first_name` (string, min 2)
- `last_name` (string, min 2)
- `email` (string, valid email format)

**Optional Fields:**
- `department_id` - Department ID (validated - must exist)
- `position_id` - Job position ID (validated - must exist)
- `location_id` - Location ID (validated - must exist)
- `phone`, `date_of_birth`, `gender`, address fields
- `hire_date` - Defaults to current date if not provided
- `employment_type` - Defaults to "full-time"
- `salary`, `currency` - Defaults to "USD"
- Emergency contact information
- `notes`

**Validation:**
- Department, position, and location IDs are validated to ensure they exist
- Employee ID and email must be unique

**Success Response (201 Created):**

```json
{
  "success": true,
  "message": "Employee onboarded successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP001",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "department_id": 1,
    "position_id": 1,
    "location_id": 1,
    "hire_date": "2024-01-15T00:00:00Z",
    "status": "active",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

#### 2. Update Employee

Update an existing employee's information.

**Endpoint:** `PUT /api/v1/employees/:id`

**Request Body:** (All fields optional)

```json
{
  "department_id": 2,
  "position_id": 3,
  "location_id": 2,
  "salary": 2500000.00,
  "status": "active"
}
```

**Validation:**
- Department, position, and location IDs are validated if provided
- Email uniqueness is checked if email is being updated

#### 3. List Employees

**Endpoint:** `GET /api/v1/employees`

**Query Parameters:**
- `page` (integer, optional) - Page number (default: 1)
- `page_size` (integer, optional) - Items per page (default: 20, max: 100)

#### 4. Get Employee by ID

**Endpoint:** `GET /api/v1/employees/:id`

#### 5. Get Employee by Employee ID

**Endpoint:** `GET /api/v1/employees/employee-id/:employee_id`

#### 6. List Employees by Department

**Endpoint:** `GET /api/v1/employees/department/:department_id`

#### 7. List Employees by Status

**Endpoint:** `GET /api/v1/employees/status/:status`

Valid status values: `active`, `inactive`, `on_leave`, `terminated`

#### 8. Delete Employee

**Endpoint:** `DELETE /api/v1/employees/:id`

---

## Common Response Formats

### Success Response

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { /* response data */ },
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Error Response

```json
{
  "success": false,
  "message": "Error message description",
  "error": "Additional error details"
}
```

---

## Error Handling

### HTTP Status Codes

- **200 OK** - Request successful
- **201 Created** - Resource created successfully
- **400 Bad Request** - Invalid request data
- **401 Unauthorized** - Missing or invalid authentication token
- **403 Forbidden** - Insufficient permissions (not HR/Admin)
- **404 Not Found** - Resource not found
- **422 Unprocessable Entity** - Validation errors
- **500 Internal Server Error** - Server error

### Common Error Scenarios

#### Validation Error (422)

```json
{
  "success": false,
  "message": "Validation failed",
  "error": "Key: 'CreateDepartmentRequest.Code' Error:Field validation for 'Code' failed on the 'required' tag"
}
```

#### Resource Not Found (404)

```json
{
  "success": false,
  "message": "Department not found",
  "error": "Resource not found"
}
```

#### Duplicate Resource (400)

```json
{
  "success": false,
  "message": "department with this code already exists",
  "error": null
}
```

#### Insufficient Permissions (403)

```json
{
  "success": false,
  "message": "Insufficient permissions",
  "error": "Forbidden"
}
```

---

## Example cURL Requests

### Create Organization

```bash
curl -X POST http://localhost:8080/api/v1/organizations \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "ACME-CORP",
    "name": "ACME Corporation",
    "country": "Tanzania",
    "status": "active"
  }'
```

### Create Department

```bash
curl -X POST http://localhost:8080/api/v1/departments \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": 1,
    "code": "IT",
    "name": "Information Technology",
    "description": "IT Department"
  }'
```

### Create Job Position

```bash
curl -X POST http://localhost:8080/api/v1/job-positions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SWE-001",
    "title": "Senior Software Engineer",
    "grade": "P4",
    "level": 4
  }'
```

### Create Location

```bash
curl -X POST http://localhost:8080/api/v1/locations \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Head Office - Dar es Salaam",
    "city": "Dar es Salaam",
    "country": "Tanzania",
    "is_head_office": true
  }'
```

### Onboard Employee with Department, Position, and Location

```bash
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "EMP001",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "department_id": 1,
    "position_id": 1,
    "location_id": 1,
    "employment_type": "full-time"
  }'
```

### Action-Based Request (POST Only)

```bash
curl -X POST http://localhost:8080/api/v1/departments/action \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "list",
    "filters": {
      "is_active": true
    },
    "pagination": {
      "page": 1,
      "page_size": 20
    }
  }'
```

---

## Data Models

### Organization Model

```typescript
interface Organization {
  id: number;
  tenant_id?: number;
  code: string;
  name: string;
  legal_name?: string;
  registration_number?: string;
  tax_id?: string;
  description?: string;
  website?: string;
  email?: string;
  phone?: string;
  address?: string;
  city?: string;
  state?: string;
  country?: string;
  postal_code?: string;
  status: "active" | "suspended" | "inactive";
  is_active: boolean;
  founded_date?: string;
  created_at: string;
  updated_at: string;
}
```

### Department Model

```typescript
interface Department {
  id: number;
  tenant_id?: number;
  organization_id?: number;  // NEW - References organization
  code: string;
  name: string;
  description?: string;
  parent_department_id?: number;
  manager_id?: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  parent_department?: Department;
  sub_departments?: Department[];
}
```

### Job Position Model

```typescript
interface JobPosition {
  id: number;
  tenant_id?: number;
  code: string;
  title: string;
  grade?: string;
  level?: number;
  description?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
```

### Location Model

```typescript
interface Location {
  id: number;
  tenant_id?: number;
  name: string;
  address?: string;
  city?: string;
  state?: string;
  country?: string;
  postal_code?: string;
  is_head_office: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
```

### Employee Model (Updated)

```typescript
interface Employee {
  id: number;
  tenant_id?: number;
  employee_id: string;
  first_name: string;
  last_name: string;
  email: string;
  department_id?: number;
  position_id?: number;
  location_id?: number;  // NEW
  hire_date?: string;
  employment_type: string;
  status: "active" | "inactive" | "on_leave" | "terminated";
  salary?: number;
  currency: string;
  created_at: string;
  updated_at: string;
}
```

---

## Notes

1. **Organizational Hierarchy**: 
   - Organizations are the top-level entity
   - Departments belong to organizations
   - Departments can have parent-child relationships within an organization
   - Employees belong to departments, positions, and locations

2. **Multi-Tenant Support**: All endpoints support tenant isolation via `X-Tenant-ID` header or authenticated user's tenant.

3. **Validation**: 
   - Organization IDs are validated when creating/updating departments
   - Department, position, and location IDs are validated when creating/updating employees
   - Organization codes, department codes, and position codes must be unique
   - Only one location can be marked as head office per tenant

4. **Hierarchical Departments**: 
   - Departments can have parent-child relationships
   - Cannot delete departments with sub-departments
   - Root departments are those without a parent within an organization

4. **Head Office**: 
   - Only one location can be marked as head office
   - Setting a new head office automatically unsets the previous one

5. **Pagination**: 
   - Default page size is 20
   - Maximum page size is 100
   - Page numbers start at 1

6. **Action-Based API**: 
   - All modules support POST-only action-based requests
   - Actions: `create`, `read`, `update`, `delete`, `list`
   - Provides enhanced security by avoiding URL parameter exposure

---

## Support

For issues or questions, please contact the development team or refer to the main project documentation.

**Last Updated**: 2024-01-15

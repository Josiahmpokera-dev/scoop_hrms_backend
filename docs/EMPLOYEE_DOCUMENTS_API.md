# Employee Documents API Documentation

## Overview

The Employee Documents API allows you to list employees with their uploaded documents and view individual documents. This API provides access to all documents that have been uploaded during employee onboarding or after onboarding completion.

**Base URL:** `/api/v1/employees/documents`

**Authentication:** All endpoints require authentication via JWT token in the `Authorization` header.

---

## Document Types

The system supports the following document types:

- `identity` - Identity documents (National ID, Passport, etc.)
- `work_permit` - Work permit documents
- `education` - Educational certificates
- `contract` - Employment contracts
- `tax_statutory` - Tax and statutory documents
- `other` - Other documents

---

## API Endpoints

### 1. Get Document Statistics

Get comprehensive statistics about employee documents including totals, expired documents, expiring soon, and breakdown by document type.

**Endpoint:** `GET /api/v1/employees/documents/statistics`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "message": "Document statistics retrieved successfully",
  "data": {
    "total_documents": 150,
    "expired": 12,
    "expiring_soon": 8,
    "document_types": [
      {
        "type": "identity",
        "count": 45
      },
      {
        "type": "work_permit",
        "count": 20
      },
      {
        "type": "education",
        "count": 35
      },
      {
        "type": "contract",
        "count": 30
      },
      {
        "type": "tax_statutory",
        "count": 15
      },
      {
        "type": "other",
        "count": 5
      }
    ]
  }
}
```

**Example Request:**
```bash
GET {{BASE_URL}}/api/v1/employees/documents/statistics
Authorization: Bearer {{TOKEN}}
```

**Response Fields:**
- `total_documents` - Total number of documents across all employees
- `expired` - Number of documents that have expired
- `expiring_soon` - Number of documents expiring within 30 days
- `document_types` - Array of document type breakdowns
  - `type` - Document type (identity, work_permit, education, contract, tax_statutory, other)
  - `count` - Number of documents of this type

**Expiration Logic:**
- **Work Permits**: Uses `WorkPermitExpiryDate` from employee statutory information
- **Identity Documents**: Uses `PassportExpiryDate` from employee statutory information
- **Other Documents**: Uses default validity period of 2 years from upload date
- **Expiring Soon**: Documents expiring within 30 days from today
- **Expired**: Documents past their expiry date

---

### 2. List All Employees with Documents

Get a paginated list of all employees with their uploaded documents.

**Endpoint:** `GET /api/v1/employees/documents`

**Authentication:** Required

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Items per page (default: 20, max: 100)

**Response:**
```json
{
  "success": true,
  "message": "Employees with documents retrieved successfully",
  "data": [
    {
      "employee_id": "EMP567073",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@company.com",
      "department_id": 5,
      "position_id": 10,
      "status": "active",
      "hire_date": "2024-01-15T00:00:00Z",
      "document_count": 3,
      "documents": "identity, education, contract"
    },
    {
      "employee_id": "EMP567074",
      "first_name": "Jane",
      "last_name": "Smith",
      "email": "jane.smith@company.com",
      "department_id": 5,
      "position_id": 11,
      "status": "active",
      "hire_date": "2024-01-20T00:00:00Z",
      "document_count": 0,
      "documents": ""
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

**Example Request:**
```bash
GET {{BASE_URL}}/api/v1/employees/documents?page=1&page_size=20
Authorization: Bearer {{TOKEN}}
```

**Response Fields:**
- `employee_id` - Employee ID (e.g., EMP567073)
- `first_name` - Employee's first name
- `last_name` - Employee's last name
- `email` - Employee's work email
- `department_id` - Department ID
- `position_id` - Position ID
- `status` - Employee status (active, inactive, etc.)
- `hire_date` - Date when employee was hired
- `document_count` - Number of documents uploaded for this employee
- `documents` - Comma-separated list of document types (e.g., "identity, education, contract"). Empty string if no documents.

---

### 2. Get Documents for Specific Employee

Get all documents uploaded for a specific employee.

**Endpoint:** `GET /api/v1/employees/documents/employee/:employee_id`

**Authentication:** Required

**URL Parameters:**
- `employee_id` (required) - Employee ID (e.g., EMP567073)

**Response:**
```json
{
  "success": true,
  "message": "Documents retrieved successfully",
  "data": [
    {
      "id": 1,
      "document_type": "identity",
      "file_name": "national_id.pdf",
      "file_url": "/storage/documents/EMP567073/national_id_20240115_123456.pdf",
      "file_size": 245678,
      "mime_type": "application/pdf",
      "description": "National ID card",
      "uploaded_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": 2,
      "document_type": "education",
      "file_name": "degree_certificate.pdf",
      "file_url": "/storage/documents/EMP567073/degree_certificate_20240115_123500.pdf",
      "file_size": 512345,
      "mime_type": "application/pdf",
      "description": "Bachelor's Degree Certificate",
      "uploaded_at": "2024-01-15T10:35:00Z"
    }
  ]
}
```

**Example Request:**
```bash
GET {{BASE_URL}}/api/v1/employees/documents/employee/EMP567073
Authorization: Bearer {{TOKEN}}
```

**Use Cases:**
- View all documents for a specific employee
- Check which documents have been uploaded
- Verify document completeness for an employee

---

### 3. View Document by ID

Get detailed information about a specific document by its ID, including the file URL for viewing/downloading.

**Endpoint:** `POST /api/v1/employees/documents/view`

**Authentication:** Required

**Request Body:**
```json
{
  "id": 1
}
```

**Response:**
```json
{
  "success": true,
  "message": "Document retrieved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP567073",
    "document_type": "identity",
    "file_name": "national_id.pdf",
    "file_url": "/storage/documents/EMP567073/national_id_20240115_123456.pdf",
    "file_size": 245678,
    "mime_type": "application/pdf",
    "description": "National ID card",
    "uploaded_at": "2024-01-15T10:30:00Z"
  }
}
```

**Example Request:**
```bash
POST {{BASE_URL}}/api/v1/employees/documents/view
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "id": 1
}
```

**Response Fields:**
- `id` - Document ID
- `employee_id` - Employee ID string (e.g., EMP567073)
- `document_type` - Type of document
- `file_name` - Original file name
- `file_url` - URL to access the document file (use this to view/download)
- `file_size` - File size in bytes
- `mime_type` - MIME type of the file
- `description` - Document description (optional)
- `uploaded_at` - Timestamp when document was uploaded

**Use Cases:**
- Get document details by ID
- Retrieve file URL for viewing/downloading
- Verify document information before accessing the file

---

## Accessing Document Files

To view or download a document file, use the `file_url` from the API response. The file URL is relative to your base URL.

**Example:**
If the API returns:
```json
{
  "file_url": "/storage/documents/EMP567073/national_id_20240115_123456.pdf"
}
```

You can access the file at:
```
GET {{BASE_URL}}/storage/documents/EMP567073/national_id_20240115_123456.pdf
```

**Note:** The storage endpoint should be configured to serve static files. Make sure your server is set up to serve files from the `/storage` path.

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": "employee_id is required"
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
  "message": "Document not found",
  "data": null
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Internal server error",
  "data": null
}
```

---

## Notes

1. **Tenant Isolation:** All endpoints respect tenant isolation. You can only access documents for employees belonging to your tenant.

2. **Document Sources:** Documents can be uploaded during:
   - Employee onboarding (Step 6 - Documents)
   - After onboarding completion

3. **Document Types:** Each employee can have multiple documents, but only one document per type during onboarding (duplicates are replaced).

4. **File URLs:** The `file_url` field contains the path to access the document. Make sure your server is configured to serve static files from the storage directory.

5. **Pagination:** When listing employees with documents, use `page` and `page_size` query parameters. Default page size is 20, maximum is 100.

6. **Empty Documents:** Employees without documents will have an empty `documents` string ("") and `document_count` of 0.

7. **Document Types List:** The `documents` field returns a comma-separated list of document types (e.g., "identity, education, contract"). This provides a quick overview of what documents have been uploaded for each employee.

---

## Integration Example

Here's a complete flow for working with employee documents:

```javascript
// 1. List all employees with their documents
const listEmployees = async () => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/documents?page=1&page_size=20`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
};

// 2. Get documents for a specific employee
const getEmployeeDocuments = async (employeeId) => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/documents/employee/${employeeId}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
};

// 3. View a specific document
const viewDocument = async (documentId) => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/documents/view`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ id: documentId })
  });
  return response.json();
};

// 4. Access the document file
const downloadDocument = (fileUrl) => {
  // fileUrl from the API response (e.g., "/storage/documents/EMP567073/file.pdf")
  window.open(`${BASE_URL}${fileUrl}`, '_blank');
};
```

---

## Summary

The Employee Documents API provides comprehensive access to employee documents:

- ✅ List all employees with their documents (paginated)
- ✅ Get all documents for a specific employee
- ✅ View document details by ID
- ✅ Access document files via file URLs
- ✅ Tenant isolation and security
- ✅ Support for multiple document types
- ✅ Document metadata (size, type, description, upload date)

This system ensures that all employee documents are properly tracked, accessible, and secure.

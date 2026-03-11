# Requests & Letters Module - API Documentation

## Overview
This documentation covers the API endpoints for the Requests & Letters module in Scoop HRMS. This module handles service requests, HR letters, and facilities requests through employee self-service.

## Base URL
```
https://api.yourdomain.com/api/self-service
```

## Authentication
All endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Error Responses
All endpoints return standardized error responses:

### 400 Bad Request
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "fieldName": "Error description"
  }
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "Authentication required"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "message": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Internal server error"
}
```

---

## Service Requests

### Data Types

#### ServiceRequest Object
```typescript
interface ServiceRequest {
  id: string;
  requestNumber: string;
  type: "HR Letter" | "IT Request" | "Facilities" | "Document Request" | "Other";
  category: string;
  subject: string;
  description: string;
  priority: "Low" | "Medium" | "High" | "Urgent";
  status: "Draft" | "Submitted" | "In Progress" | "Approved" | "Rejected" | "Completed";
  requestedDate: string; // ISO 8601
  completedDate?: string; // ISO 8601
  sla: number; // hours
  assignedTo?: string;
  attachments?: string[];
  comments?: string;
  approver?: string;
  approvalDate?: string; // ISO 8601
  createdAt: string; // ISO 8601
  updatedAt: string; // ISO 8601
}
```

#### CreateServiceRequest DTO
```typescript
interface CreateServiceRequest {
  type: "HR Letter" | "IT Request" | "Facilities" | "Document Request" | "Other";
  category: string;
  subject: string;
  description: string;
  priority: "Low" | "Medium" | "High" | "Urgent";
  attachments?: File[];
}
```

#### UpdateServiceRequest DTO
```typescript
interface UpdateServiceRequest {
  subject?: string;
  description?: string;
  priority?: "Low" | "Medium" | "High" | "Urgent";
  status?: "Draft" | "Submitted" | "In Progress" | "Approved" | "Rejected" | "Completed";
  comments?: string;
}
```

### Endpoints

#### 1. Get Service Requests
```http
GET /requests
```

**Query Parameters:**
- `page` (number): Page number (default: 1)
- `limit` (number): Items per page (default: 20)
- `status` (string): Filter by status
- `type` (string): Filter by request type
- `priority` (string): Filter by priority
- `search` (string): Search in subject/description

**Response:**
```json
{
  "success": true,
  "data": {
    "requests": ServiceRequest[],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "pages": 5
    }
  }
}
```

#### 2. Get Single Service Request
```http
GET /requests/{id}
```

**Response:**
```json
{
  "success": true,
  "data": ServiceRequest
}
```

#### 3. Create Service Request
```http
POST /requests
Content-Type: multipart/form-data
```

**Request Body:**
- Form data with CreateServiceRequest fields
- Optional file attachments

**Response:**
```json
{
  "success": true,
  "message": "Service request created successfully",
  "data": ServiceRequest
}
```

#### 4. Update Service Request
```http
PUT /requests/{id}
Content-Type: application/json
```

**Request Body:** UpdateServiceRequest

**Response:**
```json
{
  "success": true,
  "message": "Service request updated successfully",
  "data": ServiceRequest
}
```

#### 5. Delete Service Request (Draft only)
```http
DELETE /requests/{id}
```

**Response:**
```json
{
  "success": true,
  "message": "Service request deleted successfully"
}
```

#### 6. Submit Service Request
```http
POST /requests/{id}/submit
```

**Response:**
```json
{
  "success": true,
  "message": "Service request submitted successfully",
  "data": ServiceRequest
}
```

---

## HR Letters

### Data Types

#### HRLetterRequest Object
```typescript
interface HRLetterRequest {
  id: string;
  letterType: "Employment" | "Salary" | "Experience" | "NOC" | "Visa Support" | "Onboarding" | "Offboarding";
  purpose: string;
  addressedTo?: string;
  deliveryMethod: "Email" | "Hard Copy" | "Both";
  status: "Pending" | "Approved" | "Rejected" | "Generated" | "Delivered";
  requestedDate: string; // ISO 8601
  approvedDate?: string; // ISO 8601
  generatedDate?: string; // ISO 8601
  letterUrl?: string;
  createdAt: string; // ISO 8601
  updatedAt: string; // ISO 8601
}
```

#### CreateHRLetterRequest DTO
```typescript
interface CreateHRLetterRequest {
  letterType: "Employment" | "Salary" | "Experience" | "NOC" | "Visa Support" | "Onboarding" | "Offboarding";
  purpose: string;
  addressedTo?: string;
  deliveryMethod: "Email" | "Hard Copy" | "Both";
}
```

### Endpoints

#### 1. Get HR Letter Requests
```http
GET /hr-letters
```

**Query Parameters:**
- `page` (number): Page number (default: 1)
- `limit` (number): Items per page (default: 20)
- `status` (string): Filter by status
- `letterType` (string): Filter by letter type

**Response:**
```json
{
  "success": true,
  "data": {
    "letters": HRLetterRequest[],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 50,
      "pages": 3
    }
  }
}
```

#### 2. Get Single HR Letter Request
```http
GET /hr-letters/{id}
```

**Response:**
```json
{
  "success": true,
  "data": HRLetterRequest
}
```

#### 3. Request HR Letter
```http
POST /hr-letters
Content-Type: application/json
```

**Request Body:** CreateHRLetterRequest

**Response:**
```json
{
  "success": true,
  "message": "HR letter requested successfully",
  "data": HRLetterRequest
}
```

#### 4. Download HR Letter
```http
GET /hr-letters/{id}/download
```

**Response:** PDF file with appropriate Content-Type headers

---

## Facilities Requests

### Data Types

#### FacilitiesRequest Object
```typescript
interface FacilitiesRequest {
  id: string;
  requestNumber: string;
  category: "Access Cards" | "Parking" | "Seating" | "Supplies" | "Maintenance" | "Other";
  subject: string;
  description: string;
  priority: "Low" | "Medium" | "High" | "Urgent";
  status: "Draft" | "Submitted" | "In Progress" | "Approved" | "Rejected" | "Completed";
  location?: string;
  requestedDate: string; // ISO 8601
  completedDate?: string; // ISO 8601
  assignedTo?: string;
  attachments?: string[];
  comments?: string;
  createdAt: string; // ISO 8601
  updatedAt: string; // ISO 8601
}
```

### Endpoints

Facilities requests use the same endpoints as Service Requests with specific category filtering.

---

## Usage Examples

### Example 1: Create IT Access Request
```http
POST /api/self-service/requests
Content-Type: multipart/form-data
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Form Data:**
- type: "IT Request"
- category: "Access Request"
- subject: "GitHub Repository Access"
- description: "Need access to company GitHub organization for new project"
- priority: "High"

**Response:**
```json
{
  "success": true,
  "message": "Service request created successfully",
  "data": {
    "id": "SR-2024-1001",
    "requestNumber": "SR-2024-1001",
    "type": "IT Request",
    "category": "Access Request",
    "subject": "GitHub Repository Access",
    "description": "Need access to company GitHub organization for new project",
    "priority": "High",
    "status": "Draft",
    "requestedDate": "2024-01-15T10:30:00Z",
    "sla": 24,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  }
}
```

### Example 2: Request Employment Letter
```http
POST /api/self-service/hr-letters
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Request Body:**
```json
{
  "letterType": "Employment",
  "purpose": "Bank loan application",
  "addressedTo": "To Whom It May Concern",
  "deliveryMethod": "Email"
}
```

**Response:**
```json
{
  "success": true,
  "message": "HR letter requested successfully",
  "data": {
    "id": "LTR-2024-1001",
    "letterType": "Employment",
    "purpose": "Bank loan application",
    "addressedTo": "To Whom It May Concern",
    "deliveryMethod": "Email",
    "status": "Pending",
    "requestedDate": "2024-01-15T11:45:00Z",
    "createdAt": "2024-01-15T11:45:00Z",
    "updatedAt": "2024-01-15T11:45:00Z"
  }
}
```

### Example 3: Get User's Service Requests
```http
GET /api/self-service/requests?status=Submitted&type=IT Request&page=1&limit=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response:**
```json
{
  "success": true,
  "data": {
    "requests": [
      {
        "id": "SR-2024-1001",
        "requestNumber": "SR-2024-1001",
        "type": "IT Request",
        "category": "Access Request",
        "subject": "GitHub Repository Access",
        "description": "Need access to company GitHub organization for new project",
        "priority": "High",
        "status": "Submitted",
        "requestedDate": "2024-01-15T10:30:00Z",
        "sla": 24,
        "assignedTo": "IT Helpdesk",
        "createdAt": "2024-01-15T10:30:00Z",
        "updatedAt": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "pages": 1
    }
  }
}
```

---

## Rate Limiting
- Maximum 100 requests per minute per user
- Maximum 10 concurrent requests per user

## Data Validation
All endpoints validate:
- Required fields
- Field formats (email, dates, etc.)
- Enum values for type/status fields
- File types and sizes for attachments
- User permissions and ownership

## Webhook Events
- `request.created` - When a new request is created
- `request.updated` - When a request status changes
- `request.completed` - When a request is completed
- `letter.requested` - When an HR letter is requested
- `letter.generated` - When an HR letter is generated

## Versioning
API version is included in the URL path:
```
/api/v1/self-service/requests
```

## Support
For API support, contact: api-support@yourdomain.com

---

*Last Updated: 2024-01-15*
*Version: 1.0.0*
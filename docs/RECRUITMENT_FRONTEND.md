# Recruitment Module - Frontend API Reference

This document provides a comprehensive reference for the Recruitment Module API, tailored for frontend developers building the **Career Portal** (Public) and the **HR Dashboard** (Protected).

## 🌍 Base Configuration

- **Development Base URL**: `http://localhost:8080/api/v1/recruitment`
- **Production Base URL**: `https://api.scoop-hrms.com/v1/recruitment`

### Authentication
- **Public Endpoints**: No authentication required.
- **Protected Endpoints**: Require `Authorization: Bearer <token>` header.

---

## 🔓 Public Endpoints (Career Portal)

These endpoints are for external candidates to view jobs and apply.

### 1. List Published Jobs
Get a list of all active job openings for the career page.

- **Endpoint**: `GET /public/openings`
- **Query Params**: 
  - `department` (optional): Filter by department name.
- **Response**:
  ```json
  {
    "success": true,
    "message": "Job openings retrieved successfully",
    "data": [
      {
        "id": "JOB-1735001234",
        "jobTitle": "Senior Backend Engineer",
        "jobDescription": "We are looking for...",
        "department": "Engineering",
        "location": "Dar es Salaam",
        "employmentType": "Full-time",
        "experienceMin": 3,
        "experienceMax": 5,
        "mustHaveSkills": ["Go", "PostgreSQL"],
        "postedAt": "2024-02-25T10:00:00Z"
      }
    ]
  }
  ```

### 2. Submit Application
Allow a candidate to apply for a specific job. Supports Base64 resume upload.

- **Endpoint**: `POST /public/apply`
- **Headers**: `Content-Type: application/json`
- **Request Body**:
  ```json
  {
    "jobId": "JOB-1735001234",
    "firstName": "Juma",
    "lastName": "Mkapes",
    "email": "juma.m@example.com",
    "phone": "+255712345678",
    "resume": "data:application/pdf;base64,JVBERi0xLjQKJ...", // Base64 encoded file
    "resumeName": "Juma_CV_2024.pdf", // Optional: Original filename
    "linkedInUrl": "https://linkedin.com/in/juma",
    "portfolioUrl": "https://github.com/juma",
    "coverLetter": "I am excited to apply..."
  }
  ```
- **Response**:
  ```json
  {
    "success": true,
    "message": "Application submitted successfully",
    "data": {
      "id": "APP-12345",
      "stage": "Applied",
      "appliedDate": "2024-02-25T12:00:00Z"
    }
  }
  ```

---

## 🔒 Protected Endpoints (HR Dashboard)

These endpoints require a valid JWT token with Recruitment permissions.

### 📋 Job Requisitions

#### List Requisitions
- **Endpoint**: `GET /requisitions`
- **Query Params**: 
  - `status` (optional): Filter by status (e.g., `Pending`, `Approved`)
  - `department` (optional): Filter by department
  - `page`, `limit`: Pagination controls
- **Response**: List of requisitions with pagination metadata.

#### Get Requisition Details
- **Endpoint**: `GET /requisitions/{id}`
- **Response**: Detailed view of a single requisition.

#### Create Requisition
- **Endpoint**: `POST /requisitions`
- **Body**:
  ```json
  {
    "requisitionNo": "REQ-2024-001",
    "jobTitle": "Product Manager",
    "department": "Product",
    "headcount": 1,
    "priority": "High",
    "justification": "Expansion of product line",
    "targetStartDate": "2024-04-01"
  }
  ```

#### Approve/Reject Requisition
- **Endpoint**: `POST /requisitions/{id}/approval`
- **Body**:
  ```json
  {
    "action": "Approve", // or "Reject"
    "comments": "Budget approved for Q2"
  }
  ```

### 📢 Job Openings

#### Create Job Opening (from Requisition)
- **Endpoint**: `POST /openings`
- **Body**:
  ```json
  {
    "requisitionId": "REQ-UUID",
    "jobTitle": "Product Manager",
    "jobDescription": "Full markdown description...",
    "responsibilities": ["Roadmap planning", "Stakeholder management"],
    "mustHaveSkills": ["Agile", "JIRA"],
    "salaryRange": {
      "min": 2000000,
      "max": 4000000,
      "currency": "TZS"
    },
    "expiryDate": "2024-03-30T23:59:59Z"
  }
  ```

#### Publish Job
- **Endpoint**: `POST /openings/{id}/publish`
- **Body**:
  ```json
  {
    "platforms": ["LinkedIn", "Website", "Internal"],
    "expiryDate": "2024-04-15T23:59:59Z" // Optional override
  }
  ```

### 👥 Candidates & Applications

#### List Candidates
- **Endpoint**: `GET /candidates`
- **Query Params**: `jobId`, `stage` (e.g., `Screening`, `Interview`)

#### Update Application Stage
- **Endpoint**: `PATCH /candidates/{id}/stage`
- **Body**:
  ```json
  {
    "stage": "Interview",
    "notes": "Passed screening, moving to technical round"
  }
  ```

### 🗓 Interviews

#### Schedule Interview
- **Endpoint**: `POST /interviews`
- **Body**:
  ```json
  {
    "candidateId": "CAND-UUID",
    "jobId": "JOB-UUID",
    "interviewType": "Technical",
    "round": 1,
    "scheduledDate": "2024-02-28",
    "scheduledTime": "14:00",
    "duration": 60,
    "mode": "Google Meet",
    "interviewers": ["tech.lead@company.com"]
  }
  ```

#### Submit Feedback
- **Endpoint**: `POST /interviews/{id}/feedback`
- **Body**:
  ```json
  {
    "interviewerEmail": "tech.lead@company.com",
    "rating": 4, // 1-5
    "strengths": ["Strong coding skills"],
    "concerns": ["Communication could be clearer"],
    "recommendation": "Hire"
  }
  ```

### 🤝 Offers

#### Create Offer
- **Endpoint**: `POST /offers`
- **Body**:
  ```json
  {
    "candidateId": "CAND-UUID",
    "jobId": "JOB-UUID",
    "joiningDate": "2024-04-01",
    "probationPeriod": 3,
    "compensation": {
      "annualCTC": 36000000,
      "currency": "TZS",
      "components": [
        { "component": "Basic", "amount": 2000000 },
        { "component": "Allowance", "amount": 1000000 }
      ]
    },
    "expiryDate": "2024-03-10T23:59:59Z"
  }
  ```

#### Send Offer (Email)
- **Endpoint**: `POST /offers/{id}/send`
- **Response**: Triggers email to candidate with offer details.

---

## 🛠 TypeScript Interfaces

Use these definitions to type your frontend data.

```typescript
// --- Enums ---
export enum RequisitionStatus {
  Draft = "Draft",
  PendingApproval = "Pending Approval",
  Approved = "Approved",
  Rejected = "Rejected"
}

export enum ApplicationStage {
  Applied = "Applied",
  Screening = "Screening",
  Interview = "Interview",
  Offer = "Offer",
  Hired = "Hired",
  Rejected = "Rejected"
}

// --- Requests ---
export interface ApplyJobRequest {
  jobId: string;
  firstName: string;
  lastName: string;
  email: string;
  phone?: string;
  resume: string; // Base64 string starting with "data:..."
  resumeName?: string; // "my_cv.pdf"
  linkedInUrl?: string;
  coverLetter?: string;
}

// --- Models ---
export interface JobOpening {
  id: string;
  jobTitle: string;
  jobDescription: string;
  department: string; // Populated from Requisition
  status: "Draft" | "Published" | "Closed";
  mustHaveSkills: string[];
  niceToHaveSkills: string[];
  salaryMin?: number;
  salaryMax?: number;
  salaryCurrency?: string;
}

export interface Candidate {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  resumeUrl: string;
  applications: JobApplication[];
}
```

## ⚠️ Error Handling

All API errors follow this format:

```json
{
  "success": false,
  "message": "Error description here",
  "error": "Detailed technical error (optional)"
}
```

**Common Status Codes:**
- `200 OK`: Success
- `400 Bad Request`: Validation failed (check body/params)
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: User lacks permission for this action
- `404 Not Found`: Resource (Job, Candidate) not found
- `500 Internal Server Error`: Server-side issue

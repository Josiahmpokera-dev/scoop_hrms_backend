# Recruitment Module API Documentation

This document outlines the API endpoints for the Recruitment module, covering Job Requisitions, Job Openings, Candidates, Interviews, Offers, Talent Pools, and Job Applications.

## Base URL
`https://api.scoop-hrms.com/v1/recruitment`

## Authentication
All endpoints require a valid Bearer Token in the Authorization header, except for public job board endpoints.
`Authorization: Bearer <token>`

---

## 1. Job Requisitions
Manage internal job requisitions and their approval workflows.

### 1.1. List Job Requisitions
**GET** `/requisitions`

**Query Parameters:**
- `status` (optional): Filter by status (e.g., `Draft`, `Pending Approval`, `Approved`)
- `department` (optional): Filter by department
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10)

**Response:**
```json
{
  "data": [
    {
      "id": "REQ-001",
      "requisitionNumber": "REQ-2024-001",
      "jobTitle": "Senior Software Engineer",
      "department": "Engineering",
      "status": "Approved",
      "hiringManager": "John Mwita",
      "createdDate": "2024-12-15T09:00:00"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

### 1.2. Get Job Requisition Details
**GET** `/requisitions/{id}`

**Response:**
```json
{
  "id": "REQ-001",
  "requisitionNumber": "REQ-2024-001",
  "jobTitle": "Senior Software Engineer",
  "department": "Engineering",
  "location": "Dar es Salaam",
  "grade": "G4",
  "headcount": 2,
  "employmentType": "Full-time",
  "hiringManager": "John Mwita",
  "requestedBy": "Sarah Kamau",
  "targetStartDate": "2025-01-15",
  "estimatedBudget": 80000000,
  "currency": "TZS",
  "justification": "Expanding engineering team...",
  "status": "Approved",
  "approvalFlow": [
    {
      "level": 1,
      "approver": "John Mwita",
      "role": "Hiring Manager",
      "status": "Approved",
      "comments": "Critical role, approved",
      "timestamp": "2024-12-15T10:30:00"
    }
  ]
}
```

### 1.3. Create Job Requisition
**POST** `/requisitions`

**Request Body:**
```json
{
  "jobTitle": "Senior Software Engineer",
  "department": "Engineering",
  "location": "Dar es Salaam",
  "grade": "G4",
  "headcount": 2,
  "employmentType": "Full-time",
  "hiringManager": "John Mwita",
  "targetStartDate": "2025-01-15",
  "estimatedBudget": 80000000,
  "currency": "TZS",
  "justification": "Expanding engineering team..."
}
```

**Response:**
```json
{
  "id": "REQ-001",
  "message": "Requisition created successfully"
}
```

### 1.4. Update Job Requisition
**PUT** `/requisitions/{id}`

**Request Body:**
```json
{
  "headcount": 3,
  "justification": "Updated headcount requirement"
}
```

### 1.5. Approve/Reject Requisition
**POST** `/requisitions/{id}/approval`

**Request Body:**
```json
{
  "action": "Approve", // or "Reject"
  "comments": "Approved as per budget"
}
```

---

## 2. Job Openings
Manage public job postings derived from approved requisitions.

### 2.1. List Job Openings
**GET** `/openings`

**Query Parameters:**
- `status`: `Draft`, `Published`, `Closed`
- `department`: Filter by department

### 2.2. Create Job Opening
**POST** `/openings`

**Request Body:**
```json
{
  "requisitionId": "REQ-001",
  "jobTitle": "Senior Software Engineer",
  "jobDescription": "We are looking for...",
  "responsibilities": ["Design...", "Develop..."],
  "mustHaveSkills": ["React", "Node.js"],
  "niceToHaveSkills": ["AWS"],
  "experience": { "min": 5, "max": 8 },
  "education": ["Bachelor in CS"],
  "salaryRange": { "min": 35000000, "max": 45000000, "currency": "TZS" },
  "expiryDate": "2025-01-18T23:59:59"
}
```

### 2.3. Publish Job Opening
**POST** `/openings/{id}/publish`

**Request Body:**
```json
{
  "platforms": ["LinkedIn", "Indeed", "Company Website"],
  "expiryDate": "2025-02-18T23:59:59"
}
```

---

## 3. Candidates & Applications
Manage candidate profiles and their applications.

### 3.1. List Candidates
**GET** `/candidates`

**Query Parameters:**
- `jobId`: Filter by job opening
- `stage`: `Applied`, `Screening`, `Interview`, `Offer`, `Hired`

### 3.2. Submit Job Application (Public)
**POST** `/public/apply`

**Request Body:**
```json
{
  "jobId": "JOB-001",
  "firstName": "Peter",
  "lastName": "Ochieng",
  "email": "peter.ochieng@email.com",
  "phone": "+254712345678",
  "resume": "(Binary or Base64 content)",
  "linkedInUrl": "https://linkedin.com/in/peterochieng",
  "portfolioUrl": "https://peter.dev",
  "coverLetter": "I am writing to apply..."
}
```

### 3.3. Update Candidate Stage
**PATCH** `/candidates/{id}/stage`

**Request Body:**
```json
{
  "jobId": "JOB-001",
  "stage": "Interview",
  "notes": "Moved to technical interview round"
}
```

---

## 4. Interviews
Manage interview scheduling and feedback.

### 4.1. Schedule Interview
**POST** `/interviews`

**Request Body:**
```json
{
  "candidateId": "CAN-001",
  "jobId": "JOB-001",
  "interviewType": "Technical",
  "round": 2,
  "scheduledDate": "2024-12-23",
  "scheduledTime": "10:00",
  "duration": 90,
  "mode": "Video Call",
  "interviewers": ["john.mwita@scoop.com", "sarah.tech@scoop.com"]
}
```

### 4.2. Submit Interview Feedback
**POST** `/interviews/{id}/feedback`

**Request Body:**
```json
{
  "interviewerEmail": "john.mwita@scoop.com",
  "rating": 4,
  "strengths": ["Strong problem solving", "Good communication"],
  "concerns": ["Limited cloud experience"],
  "recommendation": "Hire",
  "comments": "Strong candidate, recommended for next round."
}
```

---

## 5. Offers
Manage job offers and approvals.

### 5.1. Create Offer
**POST** `/offers`

**Request Body:**
```json
{
  "candidateId": "CAN-003",
  "jobId": "JOB-002",
  "joiningDate": "2025-01-15",
  "probationPeriod": 3,
  "compensation": {
    "annualCTC": 18000000,
    "currency": "TZS",
    "components": [
      { "component": "Basic Salary", "amount": 900000 },
      { "component": "House Allowance", "amount": 350000 }
    ]
  },
  "benefits": ["Health Insurance", "Annual Leave"],
  "expiryDate": "2024-12-27T23:59:59"
}
```

### 5.2. Approve/Reject Offer
**POST** `/offers/{id}/approval`

**Request Body:**
```json
{
  "action": "Approve",
  "comments": "Offer details look good."
}
```

### 5.3. Send Offer to Candidate
**POST** `/offers/{id}/send`

**Response:**
```json
{
  "message": "Offer letter sent to candidate email."
}
```

---

## 6. Talent Pool
Manage potential candidates for future roles.

### 6.1. Add to Talent Pool
**POST** `/talent-pool`

**Request Body:**
```json
{
  "name": "Alice Muthoni",
  "email": "alice.muthoni@email.com",
  "skills": ["Product Management", "Agile"],
  "preferredRole": ["Product Manager"],
  "resumeUrl": "/files/resumes/alice.pdf"
}
```

### 6.2. Search Talent Pool
**GET** `/talent-pool/search`

**Query Parameters:**
- `skills`: Comma-separated list of skills (e.g., "React,Node.js")
- `location`: Filter by location
- `experienceMin`: Minimum years of experience

---

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| `VALIDATION_ERROR` | Invalid input data | The request body or parameters failed validation. |
| `NOT_FOUND` | Resource not found | The requested resource (job, candidate, etc.) does not exist. |
| `UNAUTHORIZED` | Unauthorized access | Invalid or missing authentication token. |
| `FORBIDDEN` | Access denied | User does not have permission to perform this action. |
| `INTERNAL_SERVER_ERROR` | Unexpected system error | An unexpected error occurred on the server. |
